package btc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

// createRBFBatch creates a new RBF (Replace-By-Fee) batch or re-submits an existing one based on pending requests.
func (w *batcherWallet) createRBFBatch(c context.Context) error {
	// Read pending requests from the cache .
	pendingRequests, err := w.cache.ReadPendingRequests(c)
	if err != nil {
		return fmt.Errorf("failed to read pending requests: %w", err)
	}

	// If there are no pending requests, return an error indicating that batch parameters are not met.
	if len(pendingRequests) == 0 {
		return ErrBatchParametersNotMet
	}

	// Read the latest RBF batch from the cache .
	latestBatch, err := w.cache.ReadLatestBatch(c)
	if err != nil {
		// If no batch is found, create a new RBF batch.
		if err == ErrStoreNotFound {
			return w.createNewRBFBatch(c, pendingRequests, 0, 0)
		}
		return fmt.Errorf("failed to read latest batch: %w", err)
	}

	// Fetch the transaction details for the latest batch.
	var tx Transaction
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		tx, err = w.indexer.GetTx(ctx, latestBatch.Tx.TxID)
		return err
	})

	// If the transaction is confirmed, create a new RBF batch.
	if tx.Status.Confirmed {
		w.logger.Info("latest batch is confirmed, creating new rbf batch", zap.String("txid", tx.TxID))
		return w.createNewRBFBatch(c, pendingRequests, 0, 0)
	}

	// Update the latest batch with the transaction details.
	latestBatch.Tx = tx

	// Re-submit the existing RBF batch with pending requests.
	return w.reSubmitBatchWithNewRequests(c, latestBatch, pendingRequests, 0)
}

// reSubmitBatchWithNewRequests re-submits an existing RBF batch with updated fee rate if necessary.
func (w *batcherWallet) reSubmitBatchWithNewRequests(c context.Context, batch Batch, pendingRequests []BatcherRequest, requiredFeeRate int) error {

	// Read requests from the cache .
	batchedRequests, err := w.cache.ReadRequests(c, maps.Keys(batch.RequestIds)...)
	if err != nil {
		w.logger.Error("failed to read requests", zap.Error(err), zap.Strings("request_ids", maps.Keys(batch.RequestIds)))
		return fmt.Errorf("failed to read requests: %w", err)
	}

	if batch.Tx.Weight == 0 {
		// Something went wrong, mostly batch.Tx is not populated well
		return fmt.Errorf("transaction %s in the batch has no weight", batch.Tx.TxID)
	}

	// Calculate the current fee rate for the batch transaction.
	currentFeeRate := int(batch.Tx.Fee) * blockchain.WitnessScaleFactor / (batch.Tx.Weight)

	// Attempt to create a new RBF batch with combined requests.
	if err = w.createNewRBFBatch(c, append(batchedRequests, pendingRequests...), currentFeeRate, 0); err != ErrTxInputsMissingOrSpent {
		if err != nil {
			w.logger.Error("failed to create new rbf batch", zap.Error(err), zap.String("txid", batch.Tx.TxID))
		}
		return err
	}

	// Get the confirmed batch.
	confirmedBatch, err := w.getConfirmedBatch(c)
	if err != nil {
		w.logger.Error("failed to get confirmed batch", zap.Error(err))
		return err
	}

	// Delete the pending batch from the cache.
	err = w.cache.DeletePendingBatches(c)
	if err != nil {
		w.logger.Error("failed to delete pending batches", zap.Error(err))
		return err
	}

	// Read the missing requests from the cache.
	missingRequestIds := getMissingRequestIds(batch.RequestIds, confirmedBatch.RequestIds)
	missingRequests, err := w.cache.ReadRequests(c, missingRequestIds...)
	if err != nil {
		w.logger.Error("failed to read missing requests", zap.Error(err), zap.Strings("request_ids", missingRequestIds))
		return err
	}

	// Create a new RBF batch with missing and pending requests.
	return w.createNewRBFBatch(c, append(missingRequests, pendingRequests...), 0, requiredFeeRate)
}

// getConfirmedBatch retrieves the confirmed RBF batch from the cache
func (w *batcherWallet) getConfirmedBatch(c context.Context) (Batch, error) {

	// Read pending batches from the cache
	batches, err := w.cache.ReadPendingBatches(c)
	if err != nil {
		w.logger.Error("failed to read pending batches", zap.Error(err))
		return Batch{}, err
	}

	confirmedBatch := Batch{}

	// Loop through the batches to find a confirmed batch
	for _, batch := range batches {
		var tx Transaction
		err := withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
			tx, err = w.indexer.GetTx(ctx, batch.Tx.TxID)
			return err
		})
		if err != nil {
			return Batch{}, err
		}

		if tx.Status.Confirmed {
			if confirmedBatch.Tx.TxID == "" {
				confirmedBatch = batch
			} else {
				return Batch{}, errors.New("multiple confirmed batches found")
			}
		}
	}

	// If no confirmed batch is found, return an error.
	if confirmedBatch.Tx.TxID == "" {
		return Batch{}, errors.New("no confirmed batch found")
	}

	return confirmedBatch, nil
}

// getMissingRequestIds identifies request IDs that are missing from the confirmed batch
func getMissingRequestIds(batchedIds, confirmedIds map[string]bool) []string {
	missingIds := []string{}
	for id := range batchedIds {
		if !confirmedIds[id] {
			missingIds = append(missingIds, id)
		}
	}
	return missingIds
}

// createNewRBFBatch creates a new RBF batch transaction and saves it to the cache
func (w *batcherWallet) createNewRBFBatch(c context.Context, pendingRequests []BatcherRequest, currentFeeRate, requiredFeeRate int) error {
	// Filter requests to get spend and send requests
	spendRequests, sendRequests, sacps, reqIds := unpackBatcherRequests(pendingRequests)

	// Get unconfirmed UTXOs to avoid them in the new transaction
	avoidUtxos, err := w.getUnconfirmedUtxos(c)
	if err != nil {
		w.logger.Error("failed to get unconfirmed utxos from cache", zap.Error(err))
		return err
	}

	// Determine the required fee rate if not provided
	if requiredFeeRate == 0 {
		var feeRates FeeSuggestion
		err := withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
			feeRates, err = w.feeEstimator.FeeSuggestion()
			return err
		})
		if err != nil {
			w.logger.Error("failed to get fee suggestion", zap.Error(err))
			return err
		}

		requiredFeeRate = selectFee(feeRates, w.opts.TxOptions.FeeLevel)
	}

	// Ensure the required fee rate is higher than the current fee rate
	// RBF cannot be performed with reduced or same fee rate
	if currentFeeRate+10 >= requiredFeeRate {
		requiredFeeRate = currentFeeRate + 10
	}

	tx, err := w.createRBFTx(
		c,
		nil,
		spendRequests,
		sendRequests,
		sacps,
		nil,
		avoidUtxos,
		0, // will be calculated in the function
		requiredFeeRate,
		false,
		2,
	)
	if err != nil {
		return err
	}

	// Submit the new RBF transaction to the indexer
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		return w.indexer.SubmitTx(ctx, tx)
	})
	if err != nil {
		return err
	}

	w.logger.Info("submitted rbf tx", zap.String("txid", tx.TxHash().String()))

	var transaction Transaction
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		transaction, err = w.indexer.GetTx(ctx, tx.TxHash().String())
		return err
	})
	if err != nil {
		return err
	}

	// Create a new batch with the transaction details and save it to the cache
	batch := Batch{
		Tx:          transaction,
		RequestIds:  reqIds,
		IsFinalized: false, // RBF transactions are not stable meaning they can be replaced
		Strategy:    RBF,
	}

	// Save the new RBF batch to the cache
	err = w.cache.SaveBatch(c, batch)
	if err != nil {
		w.logger.Error("failed to save batch to cache", zap.Error(err), zap.String("id", batch.Tx.TxID))
		return err
	}

	return nil
}

// updateRBF updates the fee rate of the latest RBF batch transaction
func (w *batcherWallet) updateRBF(c context.Context, requiredFeeRate int) error {

	// Read the latest RBF batch from the cache
	latestBatch, err := w.cache.ReadLatestBatch(c)
	if err != nil {
		if err == ErrStoreNotFound {
			return ErrFeeUpdateNotNeeded
		}
		return err
	}

	var tx Transaction
	// Check if the transaction is already confirmed
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		tx, err = w.indexer.GetTx(ctx, latestBatch.Tx.TxID)
		return err
	})
	if err != nil {
		w.logger.Error("updateRBF: failed to get tx", zap.Error(err))
		return err
	}

	if tx.Status.Confirmed && !latestBatch.Tx.Status.Confirmed {
		latestBatch.Tx = tx
		err = w.cache.UpdateAndDeletePendingBatches(c, latestBatch)
		if err == nil {
			return ErrFeeUpdateNotNeeded
		}
		w.logger.Error("updateRBF: failed to update batch", zap.Error(err))
		return err
	}

	currentFeeRate := int(tx.Fee) * 4 / tx.Weight

	// Validate the fee rate update according to the wallet options
	err = validateUpdate(currentFeeRate, requiredFeeRate, w.opts)
	if err != nil {
		return err
	}

	latestBatch.Tx = tx

	// Re-submit the RBF batch with the updated fee rate
	return w.reSubmitBatchWithNewRequests(c, latestBatch, nil, requiredFeeRate)
}

// createRBFTx creates a new RBF transaction with the given UTXOs, spend requests, and send requests
// checkValidity is used to determine if the transaction should be validated while building
// depth is used to limit the number of add cover utxos to the transaction
func (w *batcherWallet) createRBFTx(
	c context.Context,
	// Unspent transaction outputs to be used in the transaction
	utxos UTXOs,
	spendRequests []SpendRequest,
	sendRequests []SendRequest,
	sacps [][]byte,
	// Map for sequences of inputs
	sequencesMap map[string]uint32,
	// Map to avoid using certain UTXOs , those which are generated from previous unconfirmed batches
	avoidUtxos map[string]bool,
	// Transaction fee ,if fee is not provided it will dynamically added
	fee uint,
	// required fee rate per vByte
	feeRate int,
	// Flag to check the transaction's validity during construction
	checkValidity bool,
	// Depth to limit the recursion
	depth int,
) (*wire.MsgTx, error) {
	// Check if the recursion depth is exceeded
	if depth < 0 {
		w.logger.Debug(
			ErrBuildRBFDepthExceeded.Error(),
			zap.Any("utxos", utxos),
			zap.Any("spendRequests", spendRequests),
			zap.Any("sendRequests", sendRequests),
			zap.Any("sacps", sacps),
			zap.Any("sequencesMap", sequencesMap),
			zap.Any("avoidUtxos", avoidUtxos),
			zap.Uint("fee", fee),
			zap.Int("requiredFeeRate", feeRate),
			zap.Bool("checkValidity", checkValidity),
			zap.Int("depth", depth),
		)
		return nil, ErrBuildRBFDepthExceeded
	} else if depth == 0 {
		checkValidity = true
	}

	var sacpsInAmount int64
	var sacpsOutAmount int64
	var err error
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		sacpsInAmount, sacpsOutAmount, err = getSACPAmounts(ctx, sacps, w.indexer)
		return err
	})

	var spendUTXOs UTXOs
	var spendUTXOsMap map[string]UTXOs

	// Fetch UTXOs for spend requests
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		spendUTXOs, spendUTXOsMap, _, err = getUTXOsFromSpendRequest(spendRequests)
		return err
	})
	if err != nil {
		return nil, err
	}

	// Add the provided UTXOs to the spend map
	spendUTXOsMap[w.address.EncodeAddress()] = append(spendUTXOsMap[w.address.EncodeAddress()], utxos...)
	if sequencesMap == nil {
		sequencesMap = generateSequenceMap(spendUTXOsMap, spendRequests)
	}
	sequencesMap = getRbfSequenceMap(sequencesMap, utxos)

	// Combine spend UTXOs with provided UTXOs
	totalUtxos := append(spendUTXOs, utxos...)

	// Build the RBF transaction
	tx, signIdx, err := buildRBFTransaction(totalUtxos, sacps, int(sacpsInAmount-sacpsOutAmount), sendRequests, w.address, int64(fee), sequencesMap, checkValidity)
	if err != nil {
		return nil, err
	}

	// Sign the inputs related to spend requests
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		return signSpendTx(ctx, tx, signIdx, spendRequests, spendUTXOsMap, w.indexer, w.privateKey)
	})
	if err != nil {
		return nil, err
	}

	// Sign the inputs related to provided UTXOs
	err = signSendTx(tx, utxos, signIdx+len(spendUTXOs), w.address, w.privateKey)
	if err != nil {
		return nil, err
	}

	// Calculate the transaction size
	txb := btcutil.NewTx(tx)
	trueSize := mempool.GetTxVirtualSize(txb)

	// Calculate the number of signatures required
	swSigs, trSigs := getNumberOfSigs(spendRequests)
	bufferFee := 0
	if depth > 0 {
		bufferFee = ((4*(swSigs+len(utxos)) + trSigs) / 2) * feeRate
	}
	newFeeEstimate := ((int(trueSize)) * feeRate) + bufferFee

	// Check if the new fee estimate exceeds the provided fee
	if newFeeEstimate > int(fee) {
		totalIn, totalOut := func() (int64, int64) {
			totalOut := int64(0)
			for _, txOut := range tx.TxOut {
				totalOut += txOut.Value
			}

			totalIn := int64(0)
			for _, utxo := range totalUtxos {
				totalIn += utxo.Amount
			}

			totalIn += sacpsInAmount

			return totalIn, totalOut
		}()

		// If total inputs are less than the required amount, get additional UTXOs
		if totalIn < totalOut+int64(newFeeEstimate) {
			w.logger.Debug(
				"getting cover utxos",
				zap.Int64("totalIn", totalIn),
				zap.Int64("totalOut", totalOut),
				zap.Int("newFeeEstimate", newFeeEstimate),
			)
			err := withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
				utxos, _, err = w.getUtxosWithFee(ctx, totalOut+int64(newFeeEstimate)-totalIn, int64(feeRate), avoidUtxos)
				return err
			})
			if err != nil {
				return nil, err
			}
		}

		var txBytes []byte
		if txBytes, err = GetTxRawBytes(tx); err != nil {
			return nil, err
		}
		w.logger.Info(
			"rebuilding rbf tx",
			zap.Int("depth", depth),
			zap.Uint("fee", fee),
			zap.Int("newFeeEstimate", newFeeEstimate),
			zap.Int("bufferFee", bufferFee),
			zap.Int("requiredFeeRate", feeRate),
			zap.Int("TxIns", len(tx.TxIn)),
			zap.Int("TxOuts", len(tx.TxOut)),
			zap.String("TxData", hex.EncodeToString(txBytes)),
		)
		// Recursively call createRBFTx with the updated parameters
		return w.createRBFTx(c, utxos, spendRequests, sendRequests, sacps, sequencesMap, avoidUtxos, uint(newFeeEstimate), feeRate, checkValidity, depth-1)
	}

	// Return the created transaction and utxo used to fund the transaction
	return tx, nil
}

func getPendingFundingUTXOs(ctx context.Context, cache Cache, funderAddr btcutil.Address) (UTXOs, error) {
	pendingFundingUtxos, err := cache.ReadPendingBatches(ctx)
	if err != nil {
		return nil, err
	}

	script, err := txscript.PayToAddrScript(funderAddr)
	if err != nil {
		return nil, fmt.Errorf("getPendingFundingUTXOs: failed to create script for %s: %w", funderAddr.EncodeAddress(), err)
	}
	scriptHex := hex.EncodeToString(script)

	utxos := UTXOs{}
	for _, batch := range pendingFundingUtxos {
		for _, vin := range batch.Tx.VINs {
			if vin.Prevout.ScriptPubKey == scriptHex {
				utxos = append(utxos, UTXO{
					TxID:   vin.TxID,
					Vout:   uint32(vin.Vout),
					Amount: int64(vin.Prevout.Value),
					Status: &Status{
						Confirmed: false,
					},
				})
			}
		}
	}
	return utxos, nil
}

// getUtxosWithFee is an iterative function that returns self sufficient UTXOs to cover the required fee and change left
func (w *batcherWallet) getUtxosWithFee(ctx context.Context, amount, feeRate int64, avoidUtxos map[string]bool) (UTXOs, int64, error) {

	// Read pending funding UTXOs
	prevUtxos, err := getPendingFundingUTXOs(ctx, w.cache, w.address)
	if err != nil {
		w.logger.Error("failed to get pending funding utxos", zap.Error(err))
		return nil, 0, err
	}

	var coverUtxos UTXOs

	// Get UTXOs from the indexer
	err = withContextTimeout(ctx, DefaultAPITimeout, func(ctx context.Context) error {
		coverUtxos, err = w.indexer.GetUTXOs(ctx, w.address)
		return err
	})
	if err != nil {
		w.logger.Error("failed to get utxos", zap.Error(err), zap.String("address", w.address.EncodeAddress()))
		return nil, 0, err
	}

	// Combine previous UTXOs and cover UTXOs
	utxos := append(prevUtxos, coverUtxos...)
	total := int64(0)
	overhead := int64(0)
	selectedUtxos := []UTXO{}
	for _, utxo := range utxos {
		if utxo.Amount < DustAmount {
			continue
		}
		if avoidUtxos[utxo.TxID] {
			continue
		}
		total += utxo.Amount
		selectedUtxos = append(selectedUtxos, utxo)
		overhead = int64(len(selectedUtxos)*(SegwitSpendWeight)) * feeRate
		if total >= amount+overhead {
			break
		}
	}

	// Calculate the required fee and change
	requiredFee := amount + overhead
	if total < requiredFee {
		return nil, 0, errors.New("insufficient funds")
	}
	change := total - requiredFee
	if change < DustAmount {
		change = 0
	}

	// Return selected UTXOs and change amount
	return selectedUtxos, change, nil
}

func getPendingChangeUTXOs(ctx context.Context, cache Cache) ([]UTXO, error) {
	// Read pending change UTXOs
	pendingChangeUtxos, err := cache.ReadPendingBatches(ctx)
	if err != nil {
		return nil, err
	}

	utxos := []UTXO{}
	for _, batch := range pendingChangeUtxos {
		//last vout is the change output
		idx := len(batch.Tx.VOUTs) - 1
		utxos = append(utxos, UTXO{
			TxID:   batch.Tx.TxID,
			Vout:   uint32(idx),
			Amount: int64(batch.Tx.VOUTs[idx].Value),
		})

	}
	return utxos, nil
}

// getUnconfirmedUtxos returns UTXOs that are currently being spent in unconfirmed transactions to double spend them in the new transaction
func (w *batcherWallet) getUnconfirmedUtxos(ctx context.Context) (map[string]bool, error) {

	// Read pending change UTXOs from cache
	pendingChangeUtxos, err := getPendingChangeUTXOs(ctx, w.cache)
	if err != nil {
		w.logger.Error("failed to get pending change utxos", zap.Error(err))
		return nil, err
	}

	// Create a map to avoid UTXOs that are currently being spent
	avoidUtxos := make(map[string]bool)
	for _, utxo := range pendingChangeUtxos {
		avoidUtxos[utxo.TxID] = true
	}

	// Return the map of UTXOs to avoid
	return avoidUtxos, nil
}

// buildRBFTransaction builds an unsigned transaction with the given UTXOs, recipients, change address, and fee
//
// checkValidity is used to determine if the transaction should be validated while building
func buildRBFTransaction(utxos UTXOs, sacps [][]byte, sacpsFee int, recipients []SendRequest, changeAddr btcutil.Address, fee int64, sequencesMap map[string]uint32, checkValidity bool) (*wire.MsgTx, int, error) {
	tx, idx, err := buildTxFromSacps(sacps)
	if err != nil {
		return nil, 0, err
	}

	if fee > int64(sacpsFee) {
		fee -= int64(sacpsFee)
	}

	// Add inputs to the transaction
	totalUTXOAmount := int64(0)
	for _, utxo := range utxos {
		txid, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, 0, err
		}
		vout := utxo.Vout
		txIn := wire.NewTxIn(wire.NewOutPoint(txid, vout), nil, nil)
		tx.AddTxIn(txIn)

		sequence, ok := sequencesMap[utxo.TxID]
		if ok {
			tx.TxIn[len(tx.TxIn)-1].Sequence = sequence
		}

		totalUTXOAmount += utxo.Amount
	}

	// Amount being sent to the change address
	pendingAmount := int64(0)

	// Add outputs to the transaction
	totalSendAmount := int64(0)
	for _, r := range recipients {

		if r.To.EncodeAddress() == changeAddr.EncodeAddress() {
			pendingAmount += r.Amount
			continue
		}

		script, err := txscript.PayToAddrScript(r.To)
		if err != nil {
			return nil, 0, err
		}

		tx.AddTxOut(wire.NewTxOut(r.Amount, script))
		totalSendAmount += r.Amount
	}

	// Add change output to the transaction if required
	if totalUTXOAmount >= totalSendAmount+pendingAmount+fee {
		script, err := txscript.PayToAddrScript(changeAddr)
		if err != nil {
			return nil, 0, err
		}
		if totalUTXOAmount >= totalSendAmount+pendingAmount+fee+DustAmount {
			tx.AddTxOut(wire.NewTxOut(totalUTXOAmount-totalSendAmount-fee, script))
		} else if pendingAmount > 0 {
			tx.AddTxOut(wire.NewTxOut(pendingAmount, script))
		}
	} else if checkValidity {
		return nil, 0, ErrInsufficientFunds(totalUTXOAmount, totalSendAmount+pendingAmount+fee)
	}

	// Return the built transaction and the index of inputs that need to be signed
	return tx, idx, nil
}

// getRbfSequenceMap updates the sequence map with rbf sequences for cover UTXOs
func getRbfSequenceMap(sequencesMap map[string]uint32, coverUtxos UTXOs) map[string]uint32 {
	for _, utxo := range coverUtxos {
		sequencesMap[utxo.TxID] = wire.MaxTxInSequenceNum - 2
	}
	return sequencesMap
}

// getUTXOsFromSpendRequest returns UTXOs from spend requests and the total value of the UTXOs
func getUTXOsFromSpendRequest(spendReq []SpendRequest) (UTXOs, utxoMap, int64, error) {
	utxos := UTXOs{}
	totalValue := int64(0)
	utxoMap := make(utxoMap)

	for _, req := range spendReq {
		utxos = append(utxos, req.utxos...)
		for _, utxo := range req.utxos {
			totalValue += utxo.Amount
		}
		utxoMap[req.ScriptAddress.EncodeAddress()] = req.utxos
	}

	// If there are any spend requests, check if the scripts have funds to spend
	if totalValue == 0 && len(spendReq) > 0 {
		return nil, nil, 0, ErrNoFundsToSpend
	}

	return utxos, utxoMap, totalValue, nil
}
