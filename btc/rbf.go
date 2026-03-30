package btc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
	"strings"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
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
			return w.createNewRBFBatch(c, nil, pendingRequests, 0, 0, 0, 0)
		}
		return fmt.Errorf("failed to read latest batch: %w", err)
	}

	// Fetch the transaction details for the latest batch.
	var tx Transaction
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		tx, err = w.indexer.GetTx(ctx, latestBatch.Tx.TxID)
		return err
	})
	if err != nil {
		// which means the tx is not in the mempool
		// and one of the previous batch got mined
		if strings.Contains(err.Error(), "not found") {

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
			missingRequestIds := getMissingRequestIds(latestBatch.RequestIds, confirmedBatch.RequestIds)
			missingRequests, err := w.cache.ReadRequests(c, missingRequestIds...)
			if err != nil {
				w.logger.Error("failed to read missing requests", zap.Error(err), zap.Strings("request_ids", missingRequestIds))
				return err
			}

			// Create a new RBF batch with missing and pending requests.
			return w.createNewRBFBatch(c, nil, append(missingRequests, pendingRequests...), 0, 0, 0, 0)
		}

		return fmt.Errorf("failed to get tx: %w", err)
	}

	// If the transaction is confirmed, create a new RBF batch.
	if tx.Status.Confirmed {
		w.logger.Info("latest batch is confirmed, creating new rbf batch", zap.String("txid", tx.TxID))

		// Delete the pending batch from the cache.
		err = w.cache.DeletePendingBatches(c)
		if err != nil {
			w.logger.Error("failed to delete pending batches", zap.Error(err))
			return err
		}

		return w.createNewRBFBatch(c, nil, pendingRequests, 0, 0, 0, 0)
	}

	// Update the latest batch with the transaction details.
	latestBatch.Tx = tx

	// Re-submit the existing RBF batch with pending requests.
	return w.reSubmitBatchWithNewRequests(c, latestBatch, pendingRequests)
}

// reSubmitBatchWithNewRequests re-submits an existing RBF batch with updated fee rate if necessary.
func (w *batcherWallet) reSubmitBatchWithNewRequests(c context.Context, batch Batch, newRequests []BatcherRequest) error {

	// Read requests from the cache .
	existingRequests, err := w.cache.ReadRequests(c, maps.Keys(batch.RequestIds)...)
	if err != nil {
		w.logger.Error("failed to read requests", zap.Error(err), zap.Strings("request_ids", maps.Keys(batch.RequestIds)))
		return fmt.Errorf("failed to read requests: %w", err)
	}

	if batch.Tx.Weight == 0 {
		// Something went wrong, mostly batch.Tx is not populated well
		return fmt.Errorf("transaction %s in the batch has no weight", batch.Tx.TxID)
	}

	ctx, cancel := context.WithTimeout(c, DefaultAPITimeout)
	defer cancel()
	// Calculate the current fee rate for the batch transaction.
	rbfFeeInfo, err := w.rpc.GetRBFTxFeeInfo(ctx, batch.Tx.TxID)
	if err != nil {
		if !errors.Is(err, ErrTxNotFound) {
			w.logger.Error("failed to get RBF fee info", zap.Error(err), zap.String("txid", batch.Tx.TxID))
			return fmt.Errorf("failed to get RBF fee info: %w", err)
		}
		return ErrFeeUpdateNotNeeded
	}

	// currentFeeRate is the fee rate from the actual and tx and its direct descendants
	currentFeeRate := rbfFeeInfo.TxFeeRate
	w.logger.Info("current batch RBF fee info", zap.Any("rbf_fee_info", rbfFeeInfo))

	previousUTXOs := UTXOs{}
	for _, vin := range batch.Tx.VINs {

		utxoTx, err := w.indexer.GetTx(c, vin.TxID)
		if err != nil {
			return fmt.Errorf("failed to get utxo tx: %w", err)
		}
		previousUTXOs = append(previousUTXOs, UTXO{
			TxID:   vin.TxID,
			Vout:   uint32(vin.Vout),
			Amount: int64(utxoTx.VOUTs[vin.Vout].Value),
			Status: &utxoTx.Status,
		})
	}
	descendantsFee := rbfFeeInfo.DescendantFee

	// Attempt to create a new RBF batch with combined requests.
	return w.createNewRBFBatch(c, previousUTXOs, append(existingRequests, newRequests...), currentFeeRate, int(batch.Tx.Fee), 0, int(descendantsFee))
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

	w.logger.Info("found pending batches", zap.Int("count", len(batches)))

	// Loop through the batches to find a confirmed batch
	for _, batch := range batches {
		w.logger.Info("checking batch for validity", zap.String("txid", batch.Tx.TxID))
		var tx Transaction
		err := withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
			tx, err = w.indexer.GetTx(ctx, batch.Tx.TxID)
			return err
		})
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				continue
			}
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
func (w *batcherWallet) createNewRBFBatch(c context.Context, previousUTXOs UTXOs, pendingRequests []BatcherRequest, currentFeeRate float64, currentFee, requiredFeeRate, descendantsFee int) error {
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

	tx, err := w.createRBFTx(
		c,
		previousUTXOs,
		spendRequests,
		sendRequests,
		sacps,
		nil,
		avoidUtxos,
		100, // will be calculated in the function
		requiredFeeRate,
		false,
		uint(currentFee),
		currentFeeRate,
		descendantsFee,
		10,
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

	txID := tx.TxHash().String()
	w.logger.Info("submitted rbf tx", zap.String("txid", txID))

	// Retrieve the full transaction from the indexer. 
	// This is critical without it the batch won't be persisted and 
	// we'll hit "no confirmed batch found" errors. We use dedicated
	// contexts (detached from the caller) so that upstream cancellation cannot prevent us from
	// fetching and persisting the batch after a successful submission.
	//
	// Each attempt gives the indexer DefaultAPITimeout (60s) where its internal tick-based
	// retry will keep polling. If it fails, we check the Bitcoin node's mempool:
	//   - Node has the tx  → indexer is lagging, retry indexer.GetTx.
	//   - Node doesn't have it → re-submit the tx, then retry.
	const maxRetrievalAttempts = 5

	var transaction Transaction
	indexerTx, err := w.waitForTx(c, tx, maxRetrievalAttempts)
	switch {
	case err != nil:
		// Something dangerous happened (vanished from mempool, resubmit
		// failed to produce a fetchable tx). Do not populate the batch.
		w.logger.Error("dangerous failure waiting for tx", zap.String("txid", txID), zap.Error(err))
		return err
	case indexerTx != nil:
		transaction = *indexerTx
	default:
		// (nil, nil): indexer never returned the tx but nothing went wrong
		// (pure indexer lag). Safe to fall back to local tx details.
		w.logger.Warn("indexer did not return tx, using local tx details", zap.String("txid", txID))
		transaction = Transaction{
			TxID:     txID,
			Version:  int(tx.Version),
			LockTime: int(tx.LockTime),
			Status:   Status{Confirmed: false},
		}
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
		if strings.Contains(err.Error(), "not found") {
			// Get the confirmed batch.
			confirmedBatch, err := w.getConfirmedBatch(c)
			if err != nil {
				w.logger.Error("failed to get confirmed batch", zap.Error(err))
				return err
			}

			// Read the missing requests from the cache.
			missingRequestIds := getMissingRequestIds(latestBatch.RequestIds, confirmedBatch.RequestIds)
			missingRequests, err := w.cache.ReadRequests(c, missingRequestIds...)
			if err != nil {
				w.logger.Error("failed to read missing requests", zap.Error(err), zap.Strings("request_ids", missingRequestIds))
				return err
			}
			if len(missingRequests) > 0 {
				// Delete the pending batch from the cache.
				err = w.cache.DeletePendingBatches(c)
				if err != nil {
					w.logger.Error("failed to delete pending batches", zap.Error(err))
					return err
				}

				return w.createNewRBFBatch(c, nil, missingRequests, 0, 0, 0, 0)
			}
			return nil
		}
		w.logger.Error("updateRBF: failed to get tx", zap.Error(err))
		return err
	}

	if tx.Status.Confirmed && !latestBatch.Tx.Status.Confirmed {
		// Get the confirmed batch.
		confirmedBatch, err := w.getConfirmedBatch(c)
		if err != nil {
			w.logger.Error("failed to get confirmed batch", zap.Error(err))
			return err
		}

		// Read the missing requests from the cache.
		missingRequestIds := getMissingRequestIds(latestBatch.RequestIds, confirmedBatch.RequestIds)
		missingRequests, err := w.cache.ReadRequests(c, missingRequestIds...)
		if err != nil {
			w.logger.Error("failed to read missing requests", zap.Error(err), zap.Strings("request_ids", missingRequestIds))
			return err
		}
		if len(missingRequests) > 0 {
			// Delete the pending batch from the cache.
			err = w.cache.DeletePendingBatches(c)
			if err != nil {
				w.logger.Error("failed to delete pending batches", zap.Error(err))
				return err
			}

			return w.createNewRBFBatch(c, nil, missingRequests, 0, 0, 0, 0)
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(c, DefaultAPITimeout)
	defer cancel()
	// get current tx fee info
	feeInfo, err := w.rpc.GetRBFTxFeeInfo(ctx, tx.TxID)
	if err != nil {
		if !errors.Is(err, ErrTxNotFound) {
			w.logger.Error("failed to get RBF fee info", zap.Error(err), zap.String("txid", tx.TxID))
			return fmt.Errorf("failed to get RBF fee info: %w", err)
		}
		return ErrFeeUpdateNotNeeded
	}

	// Validate the fee rate update according to the wallet options
	err = validateUpdate(int(feeInfo.TxFeeRate), requiredFeeRate, w.opts)
	if err != nil {
		return err
	}

	latestBatch.Tx = tx

	// Re-submit the RBF batch with the updated fee rate
	return w.reSubmitBatchWithNewRequests(c, latestBatch, nil)
}

// Remove duplicates from `a` that are present in `b` UTXOs list
func filterDuplicates(a UTXOs, b UTXOs) UTXOs {
	result := UTXOs{}
	for _, utxo := range a {
		found := false
		for _, utxo2 := range b {
			if utxo.TxID == utxo2.TxID && utxo.Vout == utxo2.Vout {
				found = true
				break
			}
		}
		if !found {
			result = append(result, utxo)
		}
	}
	return result
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
	previousFee uint,

	previousFeeRate float64,
	descendantsFee int,
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
	var sacpsUTXOs UTXOs
	var err error
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		sacpsInAmount, sacpsOutAmount, sacpsUTXOs, err = getSACPAmounts(ctx, sacps, w.indexer)
		return err
	})

	var spendUTXOs UTXOs
	var spendUTXOsMap map[string]UTXOs
	var totalSpendsToMeValue int64

	// Fetch UTXOs for spend requests
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		spendUTXOs, spendUTXOsMap, totalSpendsToMeValue, _, err = getUTXOsFromSpendRequest(spendRequests, w.Address())
		return err
	})
	if err != nil {
		return nil, err
	}

	utxos = filterDuplicates(utxos, spendUTXOs)
	utxos = filterDuplicates(utxos, sacpsUTXOs)

	totalExistingValue := int64(0)
	for _, utxo := range utxos {
		totalExistingValue += utxo.Amount
	}

	totalSendAmount := int64(0)
	for _, r := range sendRequests {
		totalSendAmount += r.Amount
	}

	// Check if the total value of the spend UTXOs and existing UTXOs is enough to cover the fee and send requests
	if totalSpendsToMeValue+totalExistingValue < DustAmount+totalSendAmount+int64(fee) {
		previousUTXOs := utxos
		err := withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
			utxos, _, err = w.getUtxosWithFee(ctx, previousUTXOs, totalSendAmount+int64(fee)-(totalSpendsToMeValue+totalExistingValue), int64(feeRate), avoidUtxos)
			return err
		})
		if err != nil {
			return nil, err
		}
		utxos = append(previousUTXOs, utxos...)
	}

	// Add the provided UTXOs to the spend map
	spendUTXOsMap[w.Address().EncodeAddress()] = append(spendUTXOsMap[w.Address().EncodeAddress()], utxos...)
	if sequencesMap == nil {
		sequencesMap = generateSequenceMap(spendUTXOsMap, spendRequests)
	}
	sequencesMap = getRbfSequenceMap(sequencesMap, utxos)

	// Combine spend UTXOs with provided UTXOs
	totalUtxos := append(spendUTXOs, utxos...)

	// Generate the recipients for the spend requests
	extraSendRequests, err := generateSendRequests(spendRequests, spendUTXOsMap, w.Address())
	if err != nil {
		return nil, err
	}

	// Build the RBF transaction
	tx, signIdx, err := buildRBFTransaction(totalUtxos, sacps, int(sacpsInAmount-sacpsOutAmount), sendRequests, extraSendRequests, w.Address(), int64(fee), sequencesMap, checkValidity)
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
	err = w.SignCoverUTXOs(tx, utxos, signIdx+len(spendUTXOs))
	if err != nil {
		return nil, err
	}

	// Calculate the transaction size
	baseSize := tx.SerializeSizeStripped()
	totalSize := tx.SerializeSize()
	weight := baseSize*3 + totalSize
	vSize := int(math.Ceil(float64(weight) / blockchain.WitnessScaleFactor))

	fees1 := math.Ceil((float64(previousFeeRate) + 0.001) * float64(vSize))
	fees2 := math.Ceil(float64(previousFee+uint(descendantsFee)) + float64(vSize))
	fees3 := float64(feeRate * vSize)
	newFeeEstimate := int64(math.Ceil(math.Max(math.Max(fees1, fees2), fees3)))

	w.logger.Info(
		"new fee estimate for RBF transaction",
		zap.Float64("fees1", fees1),
		zap.Float64("fees2", fees2),
		zap.Float64("fees3", fees3),
		zap.Int64("newFeeEstimate", newFeeEstimate),
		zap.Int("vSize", vSize),
	)

	if newFeeEstimate > int64(fee) {
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
		script, err := txscript.PayToAddrScript(w.Address())
		if err != nil {
			return nil, err
		}
		changeAmount := int64(0)
		for _, txOut := range tx.TxOut {
			if string(txOut.PkScript) == string(script) {
				changeAmount += txOut.Value
			}
		}

		if totalOut+int64(newFeeEstimate) > totalIn && changeAmount-int64(newFeeEstimate)+int64(fee) < DustAmount {
			w.logger.Debug(
				"getting cover utxos",
				zap.Int64("totalIn", totalIn),
				zap.Int64("totalOut", totalOut),
				zap.Int64("changeAmount", changeAmount),
				zap.Int64("newFeeEstimate", newFeeEstimate),
			)
			previousUTXOs := utxos

			err := withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
				utxos, _, err = w.getUtxosWithFee(ctx, previousUTXOs, int64(newFeeEstimate), int64(feeRate), avoidUtxos)
				return err
			})
			if err != nil {
				return nil, err
			}
			utxos = append(previousUTXOs, utxos...)
		}

		var txBytes []byte
		if txBytes, err = GetTxRawBytes(tx); err != nil {
			return nil, err
		}
		w.logger.Info(
			"rebuilding rbf tx",
			zap.Int("depth", depth),
			zap.Uint("fee", fee),
			zap.Int64("newFeeEstimate", newFeeEstimate),
			zap.Int("requiredFeeRate", feeRate),
			zap.Int("TxIns", len(tx.TxIn)),
			zap.Int("TxOuts", len(tx.TxOut)),
			zap.String("TxData", hex.EncodeToString(txBytes)),
		)
		// Recursively call createRBFTx with the updated parameters
		return w.createRBFTx(c, utxos, spendRequests, sendRequests, sacps, sequencesMap, avoidUtxos, uint(newFeeEstimate), feeRate, checkValidity, previousFee, previousFeeRate, descendantsFee, depth-1)
	}

	// Return the created transaction and utxo used to fund the transaction
	return tx, nil
}

// returns the index of the change utxo and true if it exists, or false if it does not
func getChangeUTXOIndex(address btcutil.Address, txOuts []*wire.TxOut) (int64, int64, bool) {
	script, err := txscript.PayToAddrScript(address)
	if err != nil {
		return 0, 0, false
	}
	for i, txOut := range txOuts {
		if string(txOut.PkScript) == string(script) {
			return int64(i), txOut.Value, true
		}
	}
	return 0, 0, false
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
func (w *batcherWallet) getUtxosWithFee(ctx context.Context, usedUTXOS UTXOs, amount, feeRate int64, avoidUtxos map[string]bool) (UTXOs, int64, error) {

	// Read pending funding UTXOs
	prevUtxos, err := getPendingFundingUTXOs(ctx, w.cache, w.Address())
	if err != nil {
		w.logger.Error("failed to get pending funding utxos", zap.Error(err))
		return nil, 0, err
	}
	var coverUtxos UTXOs

	// Get UTXOs from the indexer
	err = withContextTimeout(ctx, DefaultAPITimeout, func(ctx context.Context) error {
		coverUtxos, err = w.indexer.GetUTXOs(ctx, w.Address())
		return err
	})
	if err != nil {
		w.logger.Error("failed to get utxos", zap.Error(err), zap.String("address", w.Address().EncodeAddress()))
		return nil, 0, err
	}

	// Combine previous UTXOs and cover UTXOs
	utxos := append(prevUtxos, coverUtxos...)
	total := int64(0)
	overhead := int64(0)
	selectedUtxos := []UTXO{}
	selectedUtxosMap := make(map[string]bool)

	for _, utxo := range utxos {
		found := false
		for _, utxo2 := range usedUTXOS {
			if utxo.TxID == utxo2.TxID && utxo.Vout == utxo2.Vout {
				found = true
				break
			}
		}
		if found {
			continue
		}
		if utxo.Amount < DustAmount {
			continue
		}
		if avoidUtxos[utxo.TxID] {
			continue
		}
		if selectedUtxosMap[utxo.TxID+strconv.Itoa(int(utxo.Vout))] {
			continue
		}
		total += utxo.Amount
		selectedUtxos = append(selectedUtxos, utxo)
		selectedUtxosMap[utxo.TxID+strconv.Itoa(int(utxo.Vout))] = true
		overhead = int64(len(selectedUtxos)*(w.CoverUTXOSpendWeight())) * feeRate
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

// waitForTx attempts to retrieve a transaction from the indexer after it has
// been submitted. It is the caller's responsibility to ensure the transaction
// has already been broadcast before calling this function.
//
// Return values:
//   - (*Transaction, nil) — indexer returned the tx, use it directly.
//   - (nil, nil)          — the tx is in the node mempool but the indexer
//     never caught up within maxMempoolWaits (~60 s). Nothing dangerous;
//     the caller may safely fall back to local tx details.
//   - (nil, error)        — something dangerous happened and the caller must
//     NOT populate a batch. This includes: tx vanished from the node mempool
//     after previously being seen, tx still not found after a resubmission,
//     or the parent context was cancelled.
//
// Decision tree per iteration:
//  1. Indexer has the tx → return (*tx, nil).
//  2. Indexer doesn't have it, node mempool has it → indexer is lagging.
//     Wait retryDelay and retry without burning an attempt. Capped at
//     maxMempoolWaits iterations; if exceeded → return (nil, nil).
//  3. Indexer doesn't have it, node mempool doesn't have it:
//     a. If previously seen in the mempool → it may have been confirmed
//        and evicted. One final indexer check; if confirmed return (*tx, nil),
//        otherwise return (nil, error).
//     b. If we already resubmitted → return (nil, error). We submitted but
//        still can't verify the tx; populating a batch would be dangerous.
//     c. Otherwise → re-submit the tx once, mark hasResubmitted, increment
//        attempt, wait retryDelay, and retry.
//
// In practice the attempt counter only increments on resubmission (step 3c),
// which can happen at most once. Subsequent iterations where neither the
// indexer nor the mempool has the tx exit immediately via step 3b. The for
// loop's maxAttempts exit is therefore a defensive fallback.
func (w *batcherWallet) waitForTx(c context.Context, tx *wire.MsgTx, maxAttempts int) (*Transaction, error) {
	const retryDelay = 2 * time.Second

	txID := tx.TxHash().String()
	seenInMempool := false
	hasResubmitted := false
	mempoolWaits := 0
	// maxMempoolWaits caps how long we wait for the indexer to catch up when
	// the tx is confirmed in the node mempool. With retryDelay=2s this gives
	// ~60s before we give up and let the caller use local tx details.
	const maxMempoolWaits = 30

	for attempt := 0; attempt < maxAttempts; {
		// 1. Ask the indexer.
		ctx, cancel := context.WithTimeout(c, DefaultAPITimeout)
		transaction, err := w.indexer.GetTx(ctx, txID)
		cancel()
		if err == nil {
			return &transaction, nil
		}

		w.logger.Error("failed to get tx from indexer",
			zap.Error(err),
			zap.String("txid", txID),
			zap.Int("attempt", attempt+1),
			zap.Int("maxAttempts", maxAttempts),
		)

		// 2. Ask the node mempool.
		mempoolCtx, mempoolCancel := context.WithTimeout(c, DefaultAPITimeout)
		_, mempoolErr := w.rpc.GetMempoolEntry(mempoolCtx, txID)
		mempoolCancel()

		if mempoolErr == nil {
			// Node has it, indexer is just lagging — wait and retry without
			// burning an attempt.
			seenInMempool = true
			mempoolWaits++
			if mempoolWaits > maxMempoolWaits {
				w.logger.Warn("indexer still hasn't caught up after max mempool waits, giving up",
					zap.String("txid", txID),
					zap.Int("mempoolWaits", mempoolWaits),
				)
				return nil, nil
			}
			w.logger.Info("tx found in node mempool, waiting for indexer to catch up",
				zap.String("txid", txID),
				zap.Int("mempoolWait", mempoolWaits),
				zap.Int("maxMempoolWaits", maxMempoolWaits),
			)
			select {
			case <-time.After(retryDelay):
			case <-c.Done():
				return nil, fmt.Errorf("context cancelled waiting for indexer to catch up for tx %s: %w", txID, c.Err())
			}
			continue
		}

		// 3. Neither indexer nor node has the tx.
		if seenInMempool {
			// It was in the mempool before but has now vanished — it likely
			// got confirmed and evicted. Do one final indexer check before
			// giving up.
			finalCtx, finalCancel := context.WithTimeout(c, DefaultAPITimeout)
			confirmedTx, finalErr := w.indexer.GetTx(finalCtx, txID)
			finalCancel()
			if finalErr == nil && confirmedTx.Status.Confirmed {
				w.logger.Info("tx confirmed in chain", zap.String("txid", txID))
				return &confirmedTx, nil
			}
			return nil, fmt.Errorf("tx %s vanished from node mempool after being seen: %w", txID, mempoolErr)
		}

		if hasResubmitted {
			// We already resubmitted but still can't find the tx anywhere.
			// Dangerous to proceed — the tx state is unknown.
			return nil, fmt.Errorf("tx %s not found after resubmission, refusing to populate batch", txID)
		}

		// Never seen in mempool and haven't resubmitted yet — submit once.
		w.logger.Warn("tx not found in node mempool, re-submitting",
			zap.Error(mempoolErr),
			zap.String("txid", txID),
			zap.Int("attempt", attempt+1),
		)
		resubmitCtx, resubmitCancel := context.WithTimeout(c, DefaultAPITimeout)
		if resubmitErr := w.indexer.SubmitTx(resubmitCtx, tx); resubmitErr != nil {
			w.logger.Warn("failed to re-submit tx", zap.Error(resubmitErr), zap.String("txid", txID))
		} else {
			w.logger.Info("re-submitted tx", zap.String("txid", txID))
		}
		resubmitCancel()
		hasResubmitted = true
		attempt++

		// Give the resubmitted tx time to propagate before retrying.
		select {
		case <-time.After(retryDelay):
		case <-c.Done():
			return nil, fmt.Errorf("context cancelled after resubmitting tx %s: %w", txID, c.Err())
		}
	}

	// This is unreachable in practice: attempt only increments on resubmit
	// (once), and subsequent loops with hasResubmitted=true exit via path D.
	// Kept as a defensive fallback.
	return nil, fmt.Errorf("exhausted %d attempts to retrieve tx %s", maxAttempts, txID)
}

// buildRBFTransaction builds an unsigned transaction with the given UTXOs, recipients, change address, and fee
//
// checkValidity is used to determine if the transaction should be validated while building
func buildRBFTransaction(utxos UTXOs, sacps [][]byte, sacpsFee int, recipients []SendRequest, redirectedRecipients []RedirectedSendRequest, changeAddr btcutil.Address, fee int64, sequencesMap map[string]uint32, checkValidity bool) (*wire.MsgTx, int, error) {
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

		sequence, ok := sequencesMap[utxo.TxID+strconv.Itoa(int(utxo.Vout))]
		if ok {
			tx.TxIn[len(tx.TxIn)-1].Sequence = sequence
		} else {
			tx.TxIn[len(tx.TxIn)-1].Sequence = wire.MaxTxInSequenceNum - 2
		}

		totalUTXOAmount += utxo.Amount
	}

	// Amount being sent to the change address

	// Add outputs to the transaction
	totalSendAmount := int64(0)
	for _, r := range recipients {

		script, err := txscript.PayToAddrScript(r.To)
		if err != nil {
			return nil, 0, err
		}

		tx.AddTxOut(wire.NewTxOut(r.Amount, script))
		totalSendAmount += r.Amount
	}

	// Add change output to the transaction if required
	if totalUTXOAmount >= totalSendAmount+fee {
		script, err := txscript.PayToAddrScript(changeAddr)
		if err != nil {
			return nil, 0, err
		}

		if len(redirectedRecipients) > 0 {
			fundsLeft := totalUTXOAmount - totalSendAmount
			feePerRecipient := fee / int64(len(redirectedRecipients))
			feeCollected := int64(0)
			for _, r := range redirectedRecipients {
				recipientScript, err := txscript.PayToAddrScript(r.To)
				if err != nil {
					return nil, 0, err
				}
				if r.Amount < feePerRecipient+DustAmount {
					// do not deduct fee here
					tx.AddTxOut(wire.NewTxOut(r.Amount, recipientScript))
				} else {
					tx.AddTxOut(wire.NewTxOut(r.Amount-feePerRecipient, recipientScript))
					feeCollected += feePerRecipient
				}
				fundsLeft -= r.Amount
			}
			remainingFee := fee - feeCollected
			if fundsLeft > DustAmount+remainingFee {
				tx.AddTxOut(wire.NewTxOut(fundsLeft-remainingFee, script))
			}
		} else if totalUTXOAmount >= totalSendAmount+fee+DustAmount {
			tx.AddTxOut(wire.NewTxOut(totalUTXOAmount-totalSendAmount-fee, script))
		}

	} else if checkValidity {
		return nil, 0, ErrInsufficientFunds(totalUTXOAmount, totalSendAmount+fee)
	} else {
		// we need more funds
		return nil, 0, errors.New("need more funds")
	}

	// Return the built transaction and the index of inputs that need to be signed
	return tx, idx, nil
}

// getRbfSequenceMap updates the sequence map with rbf sequences for cover UTXOs
func getRbfSequenceMap(sequencesMap map[string]uint32, coverUtxos UTXOs) map[string]uint32 {
	for _, utxo := range coverUtxos {
		sequencesMap[utxo.TxID+strconv.Itoa(int(utxo.Vout))] = wire.MaxTxInSequenceNum - 2
	}
	return sequencesMap
}

// getUTXOsFromSpendRequest returns UTXOs from spend requests and the total value of the UTXOs
func getUTXOsFromSpendRequest(spendReq []SpendRequest, selfAddress btcutil.Address) (UTXOs, utxoMap, int64, int64, error) {
	utxos := UTXOs{}
	totalValue := int64(0)
	utxoMap := make(utxoMap)
	spendsToMeValue := int64(0)

	for _, req := range spendReq {
		utxos = append(utxos, req.Utxos...)
		currentAmount := int64(0)
		for _, utxo := range req.Utxos {
			totalValue += utxo.Amount
			currentAmount += utxo.Amount
		}
		utxoMap[req.ScriptAddress.EncodeAddress()] = req.Utxos

		if req.Recipient == nil || (req.Recipient != nil && req.Recipient.EncodeAddress() == selfAddress.EncodeAddress()) {
			spendsToMeValue += currentAmount
		}
	}

	// If there are any spend requests, check if the scripts have funds to spend
	if totalValue == 0 && len(spendReq) > 0 {
		return nil, nil, 0, 0, ErrNoFundsToSpend
	}

	return utxos, utxoMap, spendsToMeValue, totalValue, nil
}
