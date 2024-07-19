package btc

import (
	"context"
	"encoding/hex"
	"errors"

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
	var pendingRequests []BatcherRequest
	var err error

	// Read pending requests from the cache .
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		pendingRequests, err = w.cache.ReadPendingRequests(ctx)
		return err
	})
	if err != nil {
		return err
	}

	// If there are no pending requests, return an error indicating that batch parameters are not met.
	if len(pendingRequests) == 0 {
		return ErrBatchParametersNotMet
	}

	var latestBatch Batch
	// Read the latest RBF batch from the cache .
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		latestBatch, err = w.cache.ReadLatestBatch(ctx, RBF)
		return err
	})
	if err != nil {
		// If no batch is found, create a new RBF batch.
		if err == ErrBatchNotFound {
			return w.createNewRBFBatch(c, pendingRequests, 0, 0)
		}
		return err
	}

	// Fetch the transaction details for the latest batch.
	tx, err := getTransaction(w.indexer, latestBatch.Tx.TxID)
	if err != nil {
		return err
	}

	// If the transaction is confirmed, create a new RBF batch.
	if tx.Status.Confirmed {
		return w.createNewRBFBatch(c, pendingRequests, 0, 0)
	}

	// Update the latest batch with the transaction details.
	latestBatch.Tx = tx

	// Re-submit the existing RBF batch with pending requests.
	return w.reSubmitRBFBatch(c, latestBatch, pendingRequests, 0)
}

// reSubmitRBFBatch re-submits an existing RBF batch with updated fee rate if necessary.
func (w *batcherWallet) reSubmitRBFBatch(c context.Context, batch Batch, pendingRequests []BatcherRequest, requiredFeeRate int) error {
	var batchedRequests []BatcherRequest
	var err error

	// Read batched requests from the cache .
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		batchedRequests, err = w.cache.ReadRequests(ctx, maps.Keys(batch.RequestIds))
		return err
	})
	if err != nil {
		return err
	}

	// Calculate the current fee rate for the batch transaction.
	currentFeeRate := int(batch.Tx.Fee) * blockchain.WitnessScaleFactor / (batch.Tx.Weight)

	// Attempt to create a new RBF batch with combined requests.
	if err = w.createNewRBFBatch(c, append(batchedRequests, pendingRequests...), currentFeeRate, 0); err != ErrTxInputsMissingOrSpent {
		return err
	}

	// Get the confirmed batch.
	var confirmedBatch Batch
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		confirmedBatch, err = w.getConfirmedBatch(ctx)
		return err
	})
	if err != nil {
		return err
	}

	// Delete the pending batch from the cache.
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		return w.cache.DeletePendingBatches(ctx, map[string]bool{batch.Tx.TxID: true}, RBF)
	})
	if err != nil {
		return err
	}

	// Read the missing requests from the cache.
	var missingRequests []BatcherRequest
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		missingRequestIds := getMissingRequestIds(batch.RequestIds, confirmedBatch.RequestIds)
		missingRequests, err = w.cache.ReadRequests(ctx, missingRequestIds)
		return err
	})
	if err != nil {
		return err
	}

	// Create a new RBF batch with missing and pending requests.
	return w.createNewRBFBatch(c, append(missingRequests, pendingRequests...), 0, requiredFeeRate)
}

// getConfirmedBatch retrieves the confirmed RBF batch from the cache
func (w *batcherWallet) getConfirmedBatch(c context.Context) (Batch, error) {
	var batches []Batch
	var err error

	// Read pending batches from the cache
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		batches, err = w.cache.ReadPendingBatches(ctx, RBF)
		return err
	})
	if err != nil {
		return Batch{}, err
	}

	confirmedBatch := Batch{}

	// Loop through the batches to find a confirmed batch
	for _, batch := range batches {
		var tx Transaction
		err := withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
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

	var avoidUtxos map[string]bool
	var err error

	// Get unconfirmed UTXOs to avoid them in the new transaction
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		avoidUtxos, err = w.getUnconfirmedUtxos(ctx, RBF)
		return err
	})
	if err != nil {
		return err
	}

	// Determine the required fee rate if not provided
	if requiredFeeRate == 0 {
		var feeRates FeeSuggestion
		err := withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
			feeRates, err = w.feeEstimator.FeeSuggestion()
			return err
		})
		if err != nil {
			return err
		}

		requiredFeeRate = selectFee(feeRates, w.opts.TxOptions.FeeLevel)
	}

	// Ensure the required fee rate is higher than the current fee rate
	// RBF cannot be performed with reduced or same fee rate
	if currentFeeRate >= requiredFeeRate {
		requiredFeeRate = currentFeeRate + 10
	}

	var tx *wire.MsgTx
	var fundingUtxos UTXOs
	var selfUtxos UTXOs
	// Create a new RBF transaction
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		tx, fundingUtxos, selfUtxos, err = w.createRBFTx(
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
		return err
	})
	if err != nil {
		return err
	}

	// Submit the new RBF transaction to the indexer
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		return w.indexer.SubmitTx(ctx, tx)
	})
	if err != nil {
		return err
	}

	w.logger.Info("submitted rbf tx", zap.String("txid", tx.TxHash().String()))

	var transaction Transaction
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		transaction, err = getTransaction(w.indexer, tx.TxHash().String())
		return err
	})
	if err != nil {
		return err
	}

	// Create a new batch with the transaction details and save it to the cache
	batch := Batch{
		Tx:           transaction,
		RequestIds:   reqIds,
		IsStable:     false, // RBF transactions are not stable meaning they can be replaced
		IsConfirmed:  false,
		Strategy:     RBF,
		SelfUtxos:    selfUtxos,
		FundingUtxos: fundingUtxos,
	}

	// Save the new RBF batch to the cache
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		return w.cache.SaveBatch(ctx, batch)
	})
	if err != nil {
		return err
	}

	return nil
}

// updateRBF updates the fee rate of the latest RBF batch transaction
func (w *batcherWallet) updateRBF(c context.Context, requiredFeeRate int) error {

	var latestBatch Batch
	var err error
	// Read the latest RBF batch from the cache
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		latestBatch, err = w.cache.ReadLatestBatch(ctx, RBF)
		return err
	})
	if err != nil {
		if err == ErrBatchNotFound {
			return ErrFeeUpdateNotNeeded
		}
		return err
	}

	var tx Transaction
	// Check if the transaction is already confirmed
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		tx, err = getTransaction(w.indexer, latestBatch.Tx.TxID)
		return err
	})
	if err != nil {
		return err
	}

	if tx.Status.Confirmed && !latestBatch.Tx.Status.Confirmed {
		err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
			if err = w.cache.ConfirmBatchStatuses(ctx, []string{tx.TxID}, true, RBF); err == nil {
				return ErrFeeUpdateNotNeeded
			}
			return err
		})
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
	return w.reSubmitRBFBatch(c, latestBatch, nil, requiredFeeRate)
}

// createRBFTx creates a new RBF transaction with the given UTXOs, spend requests, and send requests
// checkValidity is used to determine if the transaction should be validated while building
// depth is used to limit the number of add cover utxos to the transaction
func (w *batcherWallet) createRBFTx(
	c context.Context,
	utxos UTXOs, // Unspent transaction outputs to be used in the transaction
	spendRequests []SpendRequest,
	sendRequests []SendRequest,
	sacps [][]byte,
	sequencesMap map[string]uint32, // Map for sequences of inputs
	avoidUtxos map[string]bool, // Map to avoid using certain UTXOs , those which are generated from previous unconfirmed batches
	fee uint, // Transaction fee ,if fee is not provided it will dynamically added
	feeRate int, // required fee rate per vByte
	checkValidity bool, // Flag to check the transaction's validity during construction
	depth int, // Depth to limit the recursion
) (*wire.MsgTx, UTXOs, UTXOs, error) {
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
		return nil, nil, nil, ErrBuildRBFDepthExceeded
	} else if depth == 0 {
		checkValidity = true
	}

	var sacpsInAmount int
	var sacpsOutAmount int
	var err error
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		sacpsInAmount, sacpsOutAmount, err = getTotalInAndOutSACPs(ctx, sacps, w.indexer)
		return err
	})

	var spendUTXOs UTXOs
	var spendUTXOsMap map[btcutil.Address]UTXOs
	var balanceOfSpendScripts int64

	// Fetch UTXOs for spend requests
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		spendUTXOs, spendUTXOsMap, balanceOfSpendScripts, err = getUTXOsForSpendRequest(ctx, w.indexer, spendRequests)
		return err
	})
	if err != nil {
		return nil, nil, nil, err
	}

	// Add the provided UTXOs to the spend map
	spendUTXOsMap[w.address] = append(spendUTXOsMap[w.address], utxos...)
	if sequencesMap == nil {
		sequencesMap = generateSequenceMap(spendUTXOsMap, spendRequests)
	}
	sequencesMap = generateSequenceForCoverUtxos(sequencesMap, utxos)

	// Check if there are funds to spend
	if balanceOfSpendScripts == 0 && len(spendRequests) > 0 {
		return nil, nil, nil, ErrNoFundsToSpend
	}

	// Combine spend UTXOs with provided UTXOs
	totalUtxos := append(spendUTXOs, utxos...)

	// Build the RBF transaction
	tx, signIdx, err := buildRBFTransaction(totalUtxos, sacps, sacpsInAmount-sacpsOutAmount, sendRequests, w.address, int64(fee), sequencesMap, checkValidity)
	if err != nil {
		return nil, nil, nil, err
	}

	var selfUtxos UTXOs
	for i := 0; i < len(tx.TxOut); i++ {
		script := tx.TxOut[i].PkScript
		// convert script to btcutil.Address
		class, addrs, _, err := txscript.ExtractPkScriptAddrs(script, w.chainParams)
		if err != nil {
			return nil, nil, nil, err
		}

		if class == txscript.WitnessV0PubKeyHashTy && len(addrs) > 0 && addrs[0] == w.address {
			selfUtxos = append(selfUtxos, UTXO{
				TxID:   tx.TxHash().String(),
				Vout:   uint32(i),
				Amount: tx.TxOut[i].Value,
			})
		}
	}

	// Sign the inputs related to spend requests
	err = withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
		return signSpendTx(ctx, tx, signIdx, spendRequests, spendUTXOsMap, w.indexer, w.privateKey)
	})
	if err != nil {
		return nil, nil, nil, err
	}

	// Sign the inputs related to provided UTXOs
	err = signSendTx(tx, utxos, signIdx+len(spendUTXOs), w.address, w.privateKey)
	if err != nil {
		return nil, nil, nil, err
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
		totalIn, totalOut := func() (int, int) {
			totalOut := int64(0)
			for _, txOut := range tx.TxOut {
				totalOut += txOut.Value
			}

			totalIn := 0
			for _, utxo := range totalUtxos {
				totalIn += int(utxo.Amount)
			}

			totalIn += sacpsInAmount

			return totalIn, int(totalOut)
		}()

		// If total inputs are less than the required amount, get additional UTXOs
		if totalIn < totalOut+newFeeEstimate {
			w.logger.Debug(
				"getting cover utxos",
				zap.Int("totalIn", totalIn),
				zap.Int("totalOut", totalOut),
				zap.Int("newFeeEstimate", newFeeEstimate),
			)
			err := withContextTimeout(c, DefaultContextTimeout, func(ctx context.Context) error {
				utxos, _, err = w.getUtxosForFee(ctx, totalOut+newFeeEstimate-totalIn, feeRate, avoidUtxos)
				return err
			})
			if err != nil {
				return nil, nil, nil, err
			}
		}

		var txBytes []byte
		if txBytes, err = GetTxRawBytes(tx); err != nil {
			return nil, nil, nil, err
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
	return tx, utxos, selfUtxos, nil
}

// getUTXOsForSpendRequest returns UTXOs required to cover amount and also returns change amount if any
func (w *batcherWallet) getUtxosForFee(ctx context.Context, amount, feeRate int, avoidUtxos map[string]bool) (UTXOs, int, error) {
	var prevUtxos, coverUtxos UTXOs
	var err error

	// Read pending funding UTXOs
	err = withContextTimeout(ctx, DefaultContextTimeout, func(ctx context.Context) error {
		prevUtxos, err = w.cache.ReadPendingFundingUtxos(ctx, RBF)
		return err
	})
	if err != nil {
		return nil, 0, err
	}

	// Get UTXOs from the indexer
	err = withContextTimeout(ctx, DefaultContextTimeout, func(ctx context.Context) error {
		coverUtxos, err = w.indexer.GetUTXOs(ctx, w.address)
		return err
	})
	if err != nil {
		return nil, 0, err
	}

	// Combine previous UTXOs and cover UTXOs
	utxos := append(prevUtxos, coverUtxos...)
	total := 0
	overHead := 0
	selectedUtxos := []UTXO{}
	for _, utxo := range utxos {
		if utxo.Amount < DustAmount {
			continue
		}
		if avoidUtxos[utxo.TxID] {
			continue
		}
		total += int(utxo.Amount)
		selectedUtxos = append(selectedUtxos, utxo)
		overHead = (len(selectedUtxos) * (SegwitSpendWeight) * feeRate)
		if total >= amount+overHead {
			break
		}
	}

	// Calculate the required fee and change
	requiredFee := amount + overHead
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

// getUnconfirmedUtxos returns UTXOs that are currently being spent in unconfirmed transactions to double spend them in the new transaction
func (w *batcherWallet) getUnconfirmedUtxos(ctx context.Context, strategy Strategy) (map[string]bool, error) {
	var pendingChangeUtxos []UTXO
	var err error

	// Read pending change UTXOs
	err = withContextTimeout(ctx, DefaultContextTimeout, func(ctx context.Context) error {
		pendingChangeUtxos, err = w.cache.ReadPendingChangeUtxos(ctx, strategy)
		return err
	})
	if err != nil {
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
// checkValidity is used to determine if the transaction should be validated while building
func buildRBFTransaction(utxos UTXOs, sacps [][]byte, scapsFee int, recipients []SendRequest, changeAddr btcutil.Address, fee int64, sequencesMap map[string]uint32, checkValidity bool) (*wire.MsgTx, int, error) {
	tx, idx, err := buildTxFromSacps(sacps)
	if err != nil {
		return nil, 0, err
	}

	if fee > int64(scapsFee) {
		fee -= int64(scapsFee)
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
		if r.To == changeAddr {
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

// generateSequenceForCoverUtxos updates the sequence map with sequences for cover UTXOs
func generateSequenceForCoverUtxos(sequencesMap map[string]uint32, coverUtxos UTXOs) map[string]uint32 {
	for _, utxo := range coverUtxos {
		sequencesMap[utxo.TxID] = wire.MaxTxInSequenceNum - 2
	}
	return sequencesMap
}
