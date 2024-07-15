package btc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

// createRBFBatch creates a new RBF batch or re-submits an existing one based on pending requests
func (w *batcherWallet) createRBFBatch(c context.Context) error {
	var pendingRequests []BatcherRequest
	var err error

	// Read pending requests from the cache
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		pendingRequests, err = w.cache.ReadPendingRequests(ctx)
		return err
	})
	if err != nil {
		return err
	}

	// If there are no pending requests, return an error
	if len(pendingRequests) == 0 {
		return ErrBatchParametersNotMet
	}

	// Attempt to read the latest RBF batch from the cache
	var latestBatch Batch
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		latestBatch, err = w.cache.ReadLatestBatch(ctx, RBF)
		return err
	})
	if err != nil {
		// If no batch is found, create a new RBF batch
		if err == ErrBatchNotFound {
			return w.createNewRBFBatch(c, pendingRequests, 0)
		}
		return err
	}

	// Fetch the transaction details for the latest batch
	tx, err := getTransaction(w.indexer, latestBatch.Tx.TxID)
	if err != nil {
		return err
	}

	// If the transaction is confirmed, create a new RBF batch
	if tx.Status.Confirmed {
		return w.createNewRBFBatch(c, pendingRequests, 0)
	}

	// Update the latest batch with the transaction details
	latestBatch.Tx = tx

	// Re-submit the existing RBF batch with pending requests
	return w.reSubmitRBFBatch(c, latestBatch, pendingRequests, 0)
}

// reSubmitRBFBatch re-submits an existing RBF batch with updated fee rate if necessary
func (w *batcherWallet) reSubmitRBFBatch(c context.Context, batch Batch, pendingRequests []BatcherRequest, requiredFeeRate int) error {
	var batchedRequests []BatcherRequest
	var err error

	// Read batched requests from the cache
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		batchedRequests, err = w.cache.ReadRequests(ctx, maps.Keys(batch.RequestIds))
		return err
	})
	if err != nil {
		return err
	}

	// Attempt to create a new RBF batch with combined requests
	if err = w.createNewRBFBatch(c, append(batchedRequests, pendingRequests...), 0); err != ErrTxInputsMissingOrSpent {
		return err
	}

	// Get the confirmed batch
	var confirmedBatch Batch
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		confirmedBatch, err = w.getConfirmedBatch(ctx)
		return err
	})
	if err != nil {
		return err
	}

	// Delete the pending batch from the cache
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		return w.cache.DeletePendingBatches(ctx, map[string]bool{batch.Tx.TxID: true}, RBF)
	})
	if err != nil {
		return err
	}

	// Read the missing requests from the cache
	var missingRequests []BatcherRequest
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		missingRequestIds := getMissingRequestIds(batch.RequestIds, confirmedBatch.RequestIds)
		missingRequests, err = w.cache.ReadRequests(ctx, missingRequestIds)
		return err
	})
	if err != nil {
		return err
	}

	// Create a new RBF batch with missing and pending requests
	return w.createNewRBFBatch(c, append(missingRequests, pendingRequests...), requiredFeeRate)
}

// getConfirmedBatch retrieves the confirmed RBF batch from the cache
func (w *batcherWallet) getConfirmedBatch(c context.Context) (Batch, error) {
	var batches []Batch
	var err error

	// Read pending batches from the cache
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
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
		err := withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
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
func (w *batcherWallet) createNewRBFBatch(c context.Context, pendingRequests []BatcherRequest, requiredFeeRate int) error {
	// Filter requests to get spend and send requests
	spendRequests, sendRequests, reqIds := func() ([]SpendRequest, []SendRequest, map[string]bool) {
		spendRequests := []SpendRequest{}
		sendRequests := []SendRequest{}
		reqIds := make(map[string]bool)

		for _, req := range pendingRequests {
			spendRequests = append(spendRequests, req.Spends...)
			sendRequests = append(sendRequests, req.Sends...)
			reqIds[req.ID] = true
		}

		return spendRequests, sendRequests, reqIds
	}()

	var avoidUtxos map[string]bool
	var err error

	// Get unconfirmed UTXOs to avoid them in the new transaction
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		avoidUtxos, err = w.getUnconfirmedUtxos(ctx, RBF)
		return err
	})
	if err != nil {
		return err
	}

	// Determine the required fee rate if not provided
	if requiredFeeRate == 0 {
		var feeRates FeeSuggestion
		err := withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
			feeRates, err = w.feeEstimator.FeeSuggestion()
			return err
		})
		if err != nil {
			return err
		}

		requiredFeeRate = selectFee(feeRates, w.opts.TxOptions.FeeLevel)
	}

	var tx *wire.MsgTx
	var fundingUtxos UTXOs
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		tx, fundingUtxos, err = w.createRBFTx(
			c,
			nil,
			spendRequests,
			sendRequests,
			avoidUtxos,
			0,
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
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		return w.indexer.SubmitTx(ctx, tx)
	})
	if err != nil {
		return err
	}

	w.logger.Info("submitted rbf tx", zap.String("txid", tx.TxHash().String()))

	var transaction Transaction
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		transaction, err = getTransaction(w.indexer, tx.TxHash().String())
		return err
	})
	if err != nil {
		return err
	}

	// Create a new batch with the transaction details and save it to the cache
	batch := Batch{
		Tx:          transaction,
		RequestIds:  reqIds,
		IsStable:    false,
		IsConfirmed: false,
		Strategy:    RBF,
		ChangeUtxo: UTXO{
			TxID:   tx.TxHash().String(),
			Vout:   uint32(len(tx.TxOut) - 1),
			Amount: tx.TxOut[len(tx.TxOut)-1].Value,
		},
		FundingUtxos: fundingUtxos,
	}

	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
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
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
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
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		tx, err = getTransaction(w.indexer, latestBatch.Tx.TxID)
		return err
	})
	if err != nil {
		return err
	}

	if tx.Status.Confirmed {
		err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
			return w.cache.UpdateBatchStatuses(ctx, []string{tx.TxID}, true)
		})
		return err
	}

	size := tx.Weight / 4
	currentFeeRate := int(tx.Fee) / size

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
func (w *batcherWallet) createRBFTx(c context.Context, utxos UTXOs, spendRequests []SpendRequest, sendRequests []SendRequest, avoidUtxos map[string]bool, fee uint, requiredFeeRate int, checkValidity bool, depth int) (*wire.MsgTx, UTXOs, error) {
	if depth < 0 {
		return nil, nil, ErrBuildRBFDepthExceeded
	}

	var tx *wire.MsgTx
	var spendUTXOs UTXOs
	var spendUTXOsMap map[string]UTXOs
	var balanceOfSpendScripts int64
	var err error

	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		spendUTXOs, spendUTXOsMap, balanceOfSpendScripts, err = getUTXOsForSpendRequest(ctx, w.indexer, spendRequests)
		return err
	})
	if err != nil {
		return nil, nil, err
	}

	if balanceOfSpendScripts == 0 && len(spendRequests) > 0 {
		return nil, nil, fmt.Errorf("scripts have no funds to spend")
	}

	tx, err = buildRBFTransaction(append(spendUTXOs, utxos...), sendRequests, w.address, int(fee), checkValidity)
	if err != nil {
		return nil, nil, err
	}

	err = signSpendTx(tx, spendRequests, spendUTXOsMap, w.privateKey)
	if err != nil {
		return nil, nil, err
	}

	script, err := txscript.PayToAddrScript(w.address)
	if err != nil {
		return nil, nil, err
	}

	err = signSendTx(tx, utxos, len(spendUTXOs), script, w.privateKey)
	if err != nil {
		return nil, nil, err
	}

	txb := btcutil.NewTx(tx)
	trueSize := mempool.GetTxVirtualSize(txb)
	newFeeEstimate := int(trueSize) * requiredFeeRate

	if newFeeEstimate > int(fee) {
		var utxos UTXOs
		err := withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
			utxos, _, err = w.getUtxosForFee(ctx, newFeeEstimate, requiredFeeRate, avoidUtxos)
			return err
		})
		if err != nil {
			return nil, nil, err
		}

		return w.createRBFTx(c, utxos, spendRequests, sendRequests, avoidUtxos, uint(newFeeEstimate), requiredFeeRate, true, depth-1)
	}

	return tx, utxos, nil
}

// getUTXOsForSpendRequest returns UTXOs required to cover amount
// also return change amount if any
func (w *batcherWallet) getUtxosForFee(ctx context.Context, amount, feeRate int, avoidUtxos map[string]bool) ([]UTXO, int, error) {
	var prevUtxos, coverUtxos UTXOs
	var err error
	err = withContextTimeout(ctx, 5*time.Second, func(ctx context.Context) error {
		prevUtxos, err = w.cache.ReadPendingFundingUtxos(ctx, RBF)
		return err
	})
	if err != nil {
		return nil, 0, err
	}

	err = withContextTimeout(ctx, 5*time.Second, func(ctx context.Context) error {
		coverUtxos, err = w.indexer.GetUTXOs(ctx, w.address)
		return err
	})
	if err != nil {
		return nil, 0, err
	}

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
		overHead = (len(selectedUtxos) * SegwitSpendWeight * feeRate)
		if total >= amount+overHead {
			break
		}
	}
	requiredFee := amount + overHead
	if total < requiredFee {
		return nil, 0, errors.New("insufficient funds")
	}
	change := total - requiredFee
	if change < DustAmount {
		change = 0
	}
	return selectedUtxos, change, nil
}

// getUnconfirmedUtxos returns UTXOs that are currently being spent in unconfirmed transactions
// to double spend them in the new transaction
func (w *batcherWallet) getUnconfirmedUtxos(ctx context.Context, strategy Strategy) (map[string]bool, error) {
	var pendingChangeUtxos []UTXO
	var err error
	err = withContextTimeout(ctx, 5*time.Second, func(ctx context.Context) error {
		pendingChangeUtxos, err = w.cache.ReadPendingChangeUtxos(ctx, strategy)
		return err
	})
	if err != nil {
		return nil, err
	}

	avoidUtxos := make(map[string]bool)
	for _, utxo := range pendingChangeUtxos {
		avoidUtxos[utxo.TxID] = true
	}

	return avoidUtxos, nil
}

// buildRBFTransaction is same as buildTransaction but with validity checks
func buildRBFTransaction(utxos UTXOs, recipients []SendRequest, changeAddr btcutil.Address, fee int, checkValidity bool) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(DefaultTxVersion)

	totalUTXOAmount := int64(0)
	for _, utxo := range utxos {
		txid, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}
		vout := utxo.Vout
		txIn := wire.NewTxIn(wire.NewOutPoint(txid, vout), nil, nil)
		txIn.Sequence = wire.MaxTxInSequenceNum - 2
		tx.AddTxIn(txIn)
		totalUTXOAmount += utxo.Amount
	}

	totalSendAmount := int64(0)
	for _, r := range recipients {
		script, err := txscript.PayToAddrScript(r.To)
		if err != nil {
			return nil, err
		}
		if r.Amount < 0 {
			r.Amount = totalUTXOAmount
		}
		tx.AddTxOut(wire.NewTxOut(r.Amount, script))
		totalSendAmount += r.Amount
	}

	if totalUTXOAmount >= totalSendAmount+int64(fee) {
		script, err := txscript.PayToAddrScript(changeAddr)
		if err != nil {
			return nil, err
		}
		if totalUTXOAmount >= totalSendAmount+int64(fee)+DustAmount {
			tx.AddTxOut(wire.NewTxOut(totalUTXOAmount-totalSendAmount-int64(fee), script))
		}
	} else if checkValidity {
		return nil, fmt.Errorf("insufficient funds")
	}

	return tx, nil
}
