package btc

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/wire"
	"go.uber.org/zap"
)

type FeeStats struct {
	MaxFeeRate int
	FeeDelta   int
}

// createCPFPBatch creates a CPFP (Child Pays For Parent) batch using the pending requests
// and stores the batch in the cache
func (w *batcherWallet) createCPFPBatch(c context.Context) error {

	// Read all pending requests added to the cache
	// All requests are executed in a single batch
	requests, err := w.cache.ReadPendingRequests(c)
	if err != nil {
		w.logger.Error("failed to read pending requests", zap.Error(err))
		return err
	}

	// Return error if no requests found
	if len(requests) == 0 {
		return ErrBatchParametersNotMet
	}

	// Filter requests to get spend and send requests
	spendRequests, sendRequests, sacps, reqIds := unpackBatcherRequests(requests)

	// Read pending batches from the cache
	batches, err := w.cache.ReadPendingBatches(c)
	if err != nil {
		w.logger.Error("failed to read pending batches", zap.Error(err))
		return err
	}

	// Filter pending batches and update the status of confirmed transactions
	pendingBatches, confirmedTxs, err := filterPendingBatches(batches, w.indexer)
	if err != nil {
		return err
	}

	err = w.cache.UpdateBatches(c, confirmedTxs...)
	if err != nil {
		w.logger.Error("failed to update confirmed batches", zap.Error(err))
		return err
	}

	// Fetch fee rates and select the appropriate fee rate based on the wallet's options
	feeRates, err := w.feeEstimator.FeeSuggestion()
	if err != nil {

		return err
	}
	requiredFeeRate := selectFee(feeRates, w.opts.TxOptions.FeeLevel)

	// Calculate fee stats based on the required fee rate
	feeStats, err := getFeeStats(requiredFeeRate, pendingBatches, w.opts)
	if err != nil && !errors.Is(err, ErrFeeUpdateNotNeeded) {
		return err
	}

	// Fetch UTXOs from the indexer
	var utxos []UTXO
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
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
		sacps,
		nil,
		0,                 // will be calculated in the buildCPFPTx function
		feeStats.FeeDelta, // fee needed to bump the existing batches
		requiredFeeRate,
		1, // recursion depth
	)
	if err != nil {
		return err
	}

	// Submit the CPFP transaction to the indexer
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		return w.indexer.SubmitTx(ctx, tx)
	})
	if err != nil {
		return err
	}

	// Retrieve the transaction details from the indexer
	var transaction Transaction
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		transaction, err = w.indexer.GetTx(ctx, tx.TxHash().String())
		return err
	})
	if err != nil {
		return err
	}

	// Create a new batch and save it to the cache
	batch := Batch{
		Tx:          transaction,
		RequestIds:  reqIds,
		IsFinalized: true,
		Strategy:    CPFP,
	}

	err = w.cache.SaveBatch(c, batch)
	if err != nil {
		w.logger.Error("failed to save CPFP batch", zap.Error(err))
		return ErrSavingBatch
	}

	w.logger.Info("submitted CPFP batch", zap.String("txid", tx.TxHash().String()))
	return nil
}

// updateCPFP updates the fee rate of the pending batches to the required fee rate
func (w *batcherWallet) updateCPFP(c context.Context, requiredFeeRate int) error {

	// Read pending batches from the cache
	batches, err := w.cache.ReadPendingBatches(c)
	if err != nil {
		w.logger.Error("failed to read pending batches", zap.Error(err))
		return err
	}

	// Filter pending batches and update the status of confirmed transactions
	pendingBatches, confirmedTxs, err := filterPendingBatches(batches, w.indexer)
	if err != nil {
		return err
	}

	err = w.cache.UpdateBatches(c, confirmedTxs...)
	if err != nil {
		w.logger.Error("failed to update confirmed batches", zap.Error(err))
		return err
	}

	// Return if no pending batches are found
	if len(pendingBatches) == 0 {
		return ErrFeeUpdateNotNeeded
	}

	// Fetch UTXOs from the indexer
	var utxos []UTXO
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		utxos, err = w.indexer.GetUTXOs(ctx, w.address)
		return err
	})
	if err != nil {
		return err
	}

	// Verify CPFP conditions
	// if err := verifyCPFPConditions(utxos, pendingBatches, w.address); err != nil {
	// 	return fmt.Errorf("failed to verify CPFP conditions: %w", err)
	// }

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
		nil,
		nil,
		0,
		feeStats.FeeDelta,
		requiredFeeRate,
		1,
	)
	if err != nil {
		return err
	}

	// Submit the CPFP transaction to the indexer
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		return w.indexer.SubmitTx(ctx, tx)
	})
	if err != nil {
		return err
	}

	// Update the fee of all batches that got bumped

	newBatches := []Batch{}
	for _, batch := range pendingBatches {
		batch.Tx.Fee = int64(requiredFeeRate) * int64(batch.Tx.Weight) / blockchain.WitnessScaleFactor
		newBatches = append(newBatches, batch)
	}
	err = w.cache.UpdateBatches(c, newBatches...)
	if err != nil {
		w.logger.Error("failed to update CPFP batches", zap.Error(err))
		return err
	}

	// Log the successful submission of the CPFP transaction
	w.logger.Info("submitted CPFP transaction", zap.String("txid", tx.TxHash().String()))
	return nil
}

// buildCPFPTx builds a CPFP transaction
func (w *batcherWallet) buildCPFPTx(c context.Context, utxos []UTXO, spendRequests []SpendRequest, sendRequests []SendRequest, sacps [][]byte, sequencesMap map[string]uint32, fee, feeOverhead, feeRate int, depth int) (*wire.MsgTx, error) {
	// Check recursion depth to prevent infinite loops
	// 1 depth is optimal for most cases
	if depth < 0 {
		w.logger.Debug(
			ErrBuildCPFPDepthExceeded.Error(),
			zap.Any("utxos", utxos),
			zap.Any("spendRequests", spendRequests),
			zap.Any("sendRequests", sendRequests),
			zap.Any("sacps", sacps),
			zap.Any("sequencesMap", sequencesMap),
			zap.Int("fee", fee),
			zap.Int("feeOverhead", feeOverhead),
			zap.Int("feeRate", feeRate),
			zap.Int("depth", depth),
		)
		return nil, ErrBuildCPFPDepthExceeded
	}

	var spendUTXOs UTXOs
	var spendUTXOsMap map[string]UTXOs
	var balanceOfScripts int64
	var err error

	// Get UTXOs for spend requests
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		spendUTXOs, spendUTXOsMap, balanceOfScripts, err = getUTXOsForSpendRequest(ctx, w.indexer, spendRequests)
		return err
	})
	if err != nil {
		return nil, err
	}

	utxos, err = removeDoubleSpends(spendUTXOsMap[w.address.EncodeAddress()], utxos)
	if err != nil {
		return nil, err
	}

	spendUTXOsMap[w.address.EncodeAddress()] = append(spendUTXOsMap[w.address.EncodeAddress()], utxos...)
	if sequencesMap == nil {
		sequencesMap = generateSequenceMap(spendUTXOsMap, spendRequests)
	}

	// Check if there are no funds to spend for the given scripts
	if balanceOfScripts == 0 && len(spendRequests) > 0 {
		return nil, ErrNoFundsToSpend
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
	tx, signIdx, err := buildTransaction(append(spendUTXOs, utxos...), sacps, tempSendRequests, w.address, int64(fee+feeOverhead), sequencesMap)
	if err != nil {
		return nil, err
	}

	swSigs, trSigs := getNumberOfSigs(spendRequests)
	bufferFee := 0
	if depth > 0 {
		// buffer fee accounts for
		// 4 bytes for each sewgwit signature in the spend requests
		// 4 bytes per each UTXO in the spend requests
		// 1 byte for each taproot signature in the spend requests
		// /2 is to convert the bytes to virtual size
		bufferFee = ((4*(swSigs+len(utxos)) + trSigs) / 2) * feeRate
	}

	// Sign the spend inputs
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		return signSpendTx(ctx, tx, signIdx, spendRequests, spendUTXOsMap, w.indexer, w.privateKey)
	})
	if err != nil {
		return nil, err
	}

	// Sign the fee providing inputs, if any
	err = signSendTx(tx, utxos, signIdx+len(spendUTXOs), w.address, w.privateKey)
	if err != nil {
		return tx, err
	}

	// Calculate the true size of the transaction
	txb := btcutil.NewTx(tx)
	trueSize := mempool.GetTxVirtualSize(txb)

	var sacpsInAmount int64
	var sacpOutAmount int64
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		sacpsInAmount, sacpOutAmount, err = getSACPAmounts(ctx, sacps, w.indexer)
		return err
	})

	// Estimate the new fee
	newFeeEstimate := (int(trueSize) * (feeRate)) + feeOverhead + bufferFee - int(sacpsInAmount-sacpOutAmount)

	// If the new fee estimate exceeds the current fee, rebuild the CPFP transaction
	if newFeeEstimate > fee+feeOverhead {
		var txBytes []byte
		if txBytes, err = GetTxRawBytes(tx); err != nil {
			return nil, err
		}
		w.logger.Info(
			"rebuilding CPFP transaction",
			zap.Int("depth", depth),
			zap.Int("fee", fee),
			zap.Int("feeOverhead", feeOverhead),
			zap.Int("feeBuffer", bufferFee),
			zap.Int("required", newFeeEstimate),
			zap.Int("coverUtxos", len(utxos)),
			zap.Int("TxIns", len(tx.TxIn)),
			zap.Int("TxOuts", len(tx.TxOut)),
			zap.String("TxData", hex.EncodeToString(txBytes)),
		)
		return w.buildCPFPTx(c, utxos, spendRequests, sendRequests, sacps, sequencesMap, newFeeEstimate, 0, feeRate, depth-1)
	}

	return tx, nil
}

// CPFP (Child Pays For Parent) helpers

// reconstructCPFPBatches reconstructs the CPFP batches
func reconstructCPFPBatches([]Batch, Batch, btcutil.Address) error {
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
		size := batch.Tx.Weight / blockchain.WitnessScaleFactor
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
		FeeDelta:   feeDelta,
	}
}

func removeDoubleSpends(spends UTXOs, coverUtxos UTXOs) (UTXOs, error) {
	var utxomap = make(map[string]bool)
	for _, spendUtxo := range spends {
		utxomap[spendUtxo.TxID] = true
	}
	var newCoverUtxos UTXOs
	for _, coverUtxo := range coverUtxos {
		if _, ok := utxomap[coverUtxo.TxID]; !ok {
			newCoverUtxos = append(newCoverUtxos, coverUtxo)
		}
	}
	return newCoverUtxos, nil
}

// verifyCPFPConditions verifies the conditions required for CPFP
// func verifyCPFPConditions(utxos []UTXO, batches []Batch, walletAddr btcutil.Address) error {
// 	ucUtxos := getUnconfirmedUtxos(utxos)
// 	if len(ucUtxos) == 0 {
// 		return ErrCPFPFeeUpdateParamsNotMet
// 	}
// 	trailingBatches, err := getTrailingBatches(batches, ucUtxos)
// 	if err != nil {
// 		return err
// 	}

// 	if len(trailingBatches) == 0 {
// 		return nil
// 	}

// 	if len(trailingBatches) > 1 {
// 		return ErrCPFPBatchingCorrupted
// 	}

// 	return reconstructCPFPBatches(batches, trailingBatches[0], walletAddr)
// }

// getUnconfirmedUtxos filters and returns unconfirmed UTXOs
// func getUnconfirmedUtxos(utxos []UTXO) []UTXO {
// 	var ucUtxos []UTXO
// 	for _, utxo := range utxos {
// 		if !utxo.Status.Confirmed {
// 			ucUtxos = append(ucUtxos, utxo)
// 		}
// 	}
// 	return ucUtxos
// }

// getTrailingBatches returns batches that match the provided UTXOs
// func getTrailingBatches(batches []Batch, utxos []UTXO) ([]Batch, error) {
// 	utxomap := make(map[string]bool)
// 	for _, utxo := range utxos {
// 		utxomap[utxo.TxID] = true
// 	}

// 	trailingBatches := []Batch{}

// 	for _, batch := range batches {
// 		if _, ok := utxomap[batch.ChangeUtxo.TxID]; ok {
// 			trailingBatches = append(trailingBatches, batch)
// 		}
// 	}

// 	return trailingBatches, nil
// }
