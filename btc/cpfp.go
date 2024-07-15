package btc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"go.uber.org/zap"
)

type FeeStats struct {
	MaxFeeRate int
	TotalSize  int
	FeeDelta   int
}

// createCPFPBatch creates a CPFP (Child Pays For Parent) batch using the pending requests
// and stores the batch in the cache
func (w *batcherWallet) createCPFPBatch(c context.Context) error {
	var requests []BatcherRequest
	var err error

	// Read all pending requests added to the cache
	// All requests are executed in a single batch
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		requests, err = w.cache.ReadPendingRequests(ctx)
		return err
	})
	if err != nil {
		return err
	}

	// Filter requests to get spend and send requests
	spendRequests, sendRequests, reqIds := func() ([]SpendRequest, []SendRequest, map[string]bool) {
		spendRequests := []SpendRequest{}
		sendRequests := []SendRequest{}
		reqIds := make(map[string]bool)

		for _, req := range requests {
			spendRequests = append(spendRequests, req.Spends...)
			sendRequests = append(sendRequests, req.Sends...)
			reqIds[req.ID] = true
		}

		return spendRequests, sendRequests, reqIds
	}()

	// Return error if no requests found
	if len(sendRequests) == 0 && len(spendRequests) == 0 {
		return ErrBatchParametersNotMet
	}

	// Validate spend requests
	err = validateSpendRequest(spendRequests)
	if err != nil {
		return err
	}

	// Fetch fee rates and select the appropriate fee rate based on the wallet's options
	feeRates, err := w.feeEstimator.FeeSuggestion()
	if err != nil {
		return err
	}
	requiredFeeRate := selectFee(feeRates, w.opts.TxOptions.FeeLevel)

	// Read pending batches from the cache
	var batches []Batch
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		batches, err = w.cache.ReadPendingBatches(ctx, w.opts.Strategy)
		return err
	})
	if err != nil {
		return err
	}

	// Filter pending batches and update the status of confirmed transactions
	pendingBatches, confirmedTxs, _, err := filterPendingBatches(batches, w.indexer)
	if err != nil {
		return err
	}

	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		return w.cache.UpdateBatchStatuses(ctx, confirmedTxs, true)
	})
	if err != nil {
		return err
	}

	// Calculate fee stats based on the required fee rate
	feeStats, err := getFeeStats(requiredFeeRate, pendingBatches, w.opts)
	if err != nil && !errors.Is(err, ErrFeeUpdateNotNeeded) {
		return err
	}

	// Fetch UTXOs from the indexer
	var utxos []UTXO
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		utxos, err = w.indexer.GetUTXOs(ctx, w.address)
		return err
	})
	if err != nil {
		return err
	}

	// Build the CPFP transaction
	tx, err := w.buildCPFPTx(
		c,     // parent context
		utxos, // all utxos available in the wallet
		spendRequests,
		sendRequests,
		0,                 // will be calculated in the buildCPFPTx function
		feeStats.FeeDelta, // fee needed to bump the existing batches
		requiredFeeRate,
		1, // recursion depth
	)
	if err != nil {
		return err
	}

	// Submit the CPFP transaction to the indexer
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		return w.indexer.SubmitTx(ctx, tx)
	})
	if err != nil {
		return err
	}

	// Retrieve the transaction details from the indexer
	var transaction Transaction
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		transaction, err = getTransaction(w.indexer, tx.TxHash().String())
		return err
	})
	if err != nil {
		return err
	}

	// Create a new batch and save it to the cache
	batch := Batch{
		Tx:          transaction,
		RequestIds:  reqIds,
		IsStable:    true,
		IsConfirmed: false,
		Strategy:    CPFP,
		ChangeUtxo: UTXO{
			TxID:   tx.TxHash().String(),
			Vout:   uint32(len(tx.TxOut) - 1),
			Amount: tx.TxOut[len(tx.TxOut)-1].Value,
		},
	}

	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		return w.cache.SaveBatch(ctx, batch)
	})
	if err != nil {
		return ErrSavingBatch
	}

	w.logger.Info("submitted CPFP batch", zap.String("txid", tx.TxHash().String()))
	return nil
}

// updateCPFP updates the fee rate of the pending batches to the required fee rate
func (w *batcherWallet) updateCPFP(c context.Context, requiredFeeRate int) error {
	var batches []Batch
	var err error

	// Read pending batches from the cache
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		batches, err = w.cache.ReadPendingBatches(ctx, w.opts.Strategy)
		return err
	})
	if err != nil {
		return err
	}

	// Filter pending batches and update the status of confirmed transactions
	pendingBatches, confirmedTxs, pendingTxs, err := filterPendingBatches(batches, w.indexer)
	if err != nil {
		return err
	}

	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		return w.cache.UpdateBatchStatuses(ctx, confirmedTxs, true)
	})
	if err != nil {
		return err
	}

	// Return if no pending batches are found
	if len(pendingBatches) == 0 {
		return ErrFeeUpdateNotNeeded
	}

	// Fetch UTXOs from the indexer
	var utxos []UTXO
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		utxos, err = w.indexer.GetUTXOs(ctx, w.address)
		return err
	})
	if err != nil {
		return err
	}

	// Verify CPFP conditions
	if err := verifyCPFPConditions(utxos, pendingBatches, w.address); err != nil {
		return fmt.Errorf("failed to verify CPFP conditions: %w", err)
	}

	// Calculate fee stats based on the required fee rate
	feeStats, err := getFeeStats(requiredFeeRate, pendingBatches, w.opts)
	if err != nil {
		return err
	}

	// Build the CPFP transaction
	tx, err := w.buildCPFPTx(
		c,
		utxos,
		[]SpendRequest{},
		[]SendRequest{},
		0,
		feeStats.FeeDelta,
		requiredFeeRate,
		1,
	)
	if err != nil {
		return err
	}

	// Submit the CPFP transaction to the indexer
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		return w.indexer.SubmitTx(ctx, tx)
	})
	if err != nil {
		return err
	}

	// Update the fee of all batches that got bumped
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		return w.cache.UpdateBatchFees(ctx, pendingTxs, int64(requiredFeeRate))
	})
	if err != nil {
		return err
	}

	// Log the successful submission of the CPFP transaction
	w.logger.Info("submitted CPFP transaction", zap.String("txid", tx.TxHash().String()))
	return nil
}

// buildCPFPTx builds a CPFP transaction
func (w *batcherWallet) buildCPFPTx(c context.Context, utxos []UTXO, spendRequests []SpendRequest, sendRequests []SendRequest, fee, feeOverhead, feeRate int, depth int) (*wire.MsgTx, error) {
	// Check recursion depth to prevent infinite loops
	// 1 depth is optimal for most cases
	if depth < 0 {
		return nil, ErrBuildCPFPDepthExceeded
	}

	var spendUTXOs UTXOs
	var spendUTXOMap map[string]UTXOs
	var balanceOfScripts int64
	var err error

	// Get UTXOs for spend requests
	err = withContextTimeout(c, 5*time.Second, func(ctx context.Context) error {
		spendUTXOs, spendUTXOMap, balanceOfScripts, err = getUTXOsForSpendRequest(ctx, w.indexer, spendRequests)
		return err
	})
	if err != nil {
		return nil, err
	}

	// Check if there are no funds to spend for the given scripts
	if balanceOfScripts == 0 && len(spendRequests) > 0 {
		return nil, fmt.Errorf("scripts have no funds to spend")
	}

	// Temporary send requests for the transaction
	tempSendRequests := append([]SendRequest{}, sendRequests...)

	// If there are UTXOs and no spend/send requests, create a self-send request
	if len(utxos) > 0 && len(tempSendRequests) == 0 && len(spendRequests) == 0 {
		amount := int64(0)
		for _, utxo := range utxos {
			amount += utxo.Amount
		}
		tempSendRequests = append(sendRequests, SendRequest{
			Amount: amount - int64(fee+feeOverhead),
			To:     w.address,
		})
	}

	// Build the transaction with the available UTXOs and requests
	tx, err := buildTransaction(append(spendUTXOs, utxos...), tempSendRequests, w.address, fee+feeOverhead)
	if err != nil {
		return nil, err
	}

	// Sign the spend inputs
	err = signSpendTx(tx, spendRequests, spendUTXOMap, w.privateKey)
	if err != nil {
		return nil, err
	}

	// Get the send signing script
	script, err := txscript.PayToAddrScript(w.address)
	if err != nil {
		return tx, err
	}

	// Sign the fee providing inputs, if any
	err = signSendTx(tx, utxos, len(spendUTXOs), script, w.privateKey)
	if err != nil {
		return tx, err
	}

	// Calculate the true size of the transaction
	txb := btcutil.NewTx(tx)
	trueSize := mempool.GetTxVirtualSize(txb)

	// Estimate the new fee
	newFeeEstimate := (int(trueSize) * (feeRate)) + feeOverhead

	// If the new fee estimate exceeds the current fee, rebuild the CPFP transaction
	if newFeeEstimate > fee+feeOverhead {
		return w.buildCPFPTx(c, utxos, spendRequests, sendRequests, newFeeEstimate, 0, feeRate, depth-1)
	}

	return tx, nil
}

// CPFP (Child Pays For Parent) helpers

// verifyCPFPConditions verifies the conditions required for CPFP
func verifyCPFPConditions(utxos []UTXO, batches []Batch, walletAddr btcutil.Address) error {
	ucUtxos := getUnconfirmedUtxos(utxos)
	if len(ucUtxos) == 0 {
		return ErrCPFPFeeUpdateParamsNotMet
	}
	trailingBatches, err := getTrailingBatches(batches, ucUtxos)
	if err != nil {
		return err
	}

	if len(trailingBatches) == 0 {
		return nil
	}

	if len(trailingBatches) > 1 {
		return ErrCPFPBatchingCorrupted
	}

	return reconstructCPFPBatches(batches, trailingBatches[0], walletAddr)
}

// getUnconfirmedUtxos filters and returns unconfirmed UTXOs
func getUnconfirmedUtxos(utxos []UTXO) []UTXO {
	var ucUtxos []UTXO
	for _, utxo := range utxos {
		if !utxo.Status.Confirmed {
			ucUtxos = append(ucUtxos, utxo)
		}
	}
	return ucUtxos
}

// getTrailingBatches returns batches that match the provided UTXOs
func getTrailingBatches(batches []Batch, utxos []UTXO) ([]Batch, error) {
	utxomap := make(map[string]bool)
	for _, utxo := range utxos {
		utxomap[utxo.TxID] = true
	}

	trailingBatches := []Batch{}

	for _, batch := range batches {
		if _, ok := utxomap[batch.ChangeUtxo.TxID]; ok {
			trailingBatches = append(trailingBatches, batch)
		}
	}

	return trailingBatches, nil
}

// reconstructCPFPBatches reconstructs the CPFP batches
func reconstructCPFPBatches(batches []Batch, trailingBatch Batch, walletAddr btcutil.Address) error {
	// TODO: Verify that the trailing batch can trace back to the funding UTXOs from the wallet address
	// This is essential to ensure that all the pending transactions are moved to the estimated
	// fee rate and the trailing batch is the only one that needs to be bumped
	// Current implementation assumes that the trailing batch is the last batch in the list
	// It maintains only one thread of CPFP transactions
	return nil
}

// getFeeStats generates fee stats based on the required fee rate
// Fee stats are used to determine how much fee is required
// to bump existing batches
func getFeeStats(requiredFeeRate int, pendingBatches []Batch, opts BatcherOptions) (FeeStats, error) {
	feeStats := calculateFeeStats(requiredFeeRate, pendingBatches)
	if err := validateUpdate(feeStats.MaxFeeRate, requiredFeeRate, opts); err != nil {
		if err == ErrFeeUpdateNotNeeded && feeStats.FeeDelta > 0 {
			return feeStats, nil
		}
		return FeeStats{}, err
	}
	return feeStats, nil
}

// calculateFeeStats calculates the fee stats based on the required fee rate
func calculateFeeStats(reqFeeRate int, batches []Batch) FeeStats {
	maxFeeRate := int(0)
	totalSize := int(0)
	feeDelta := int(0)

	for _, batch := range batches {
		size := batch.Tx.Weight / 4
		feeRate := int(batch.Tx.Fee) / size
		if feeRate > maxFeeRate {
			maxFeeRate = feeRate
		}
		if reqFeeRate > feeRate {
			feeDelta += (reqFeeRate - feeRate) * size
		}
		totalSize += size
	}
	return FeeStats{
		MaxFeeRate: maxFeeRate,
		TotalSize:  totalSize,
		FeeDelta:   feeDelta,
	}
}