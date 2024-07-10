package btc

import (
	"context"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/log"
)

// create a CPFP batch using the pending requests
// stores the batch in the cache
func (w *batcherWallet) createCPFPBatch() error {
	ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel1()
	requests, err := w.cache.ReadPendingRequests(context.Background())
	if err != nil {
		return err
	}
	sendRequests := []SendRequest{}
	spendRequests := []SpendRequest{}
	reqIds := make(map[string]bool)

	for _, req := range requests {
		sendRequests = append(sendRequests, req.Sends...)
		spendRequests = append(spendRequests, req.Spends...)
		reqIds[req.ID] = true
	}

	err = validateSpendRequest(spendRequests)
	if err != nil {
		return err
	}

	utxos, err := w.indexer.GetUTXOs(ctx1, w.address)
	if err != nil {
		return err
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	feeRate, err := w.indexer.FeeEstimate(ctx2)
	if err != nil {
		return err
	}
	requiredFeeRate := selectFee(feeRate, w.opts.TxOptions.FeeLevel)

	feeStats, err := w.getFeeStats(requiredFeeRate)
	if err != nil {
		return err
	}

	tx, err := w.buildCPFPTx(
		utxos,
		spendRequests,
		sendRequests,
		(len(utxos)+len(spendRequests)*SegwitSpendWeight)+10,
		feeStats.FeeDelta,
		requiredFeeRate,
	)
	if err != nil {
		return err
	}

	ctx3, cancel3 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel3()
	if err := w.indexer.SubmitTx(ctx3, tx); err != nil {
		return err
	}

	transaction, err := getTransaction(w.indexer, tx.TxHash().String())

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

	ctx4, cancel4 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel4()
	if err := w.cache.SaveBatch(ctx4, batch); err != nil {
		return ErrSavingBatch
	}

	log.Info("submitted CPFP batch", "txid", tx.TxHash().String())
	return nil

}

func (w *batcherWallet) updateCPFP(feeStats FeeStats, requiredFeeRate int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	utxos, err := w.indexer.GetUTXOs(ctx, w.address)
	if err != nil {
		return err
	}

	batches, err := w.cache.ReadPendingBatches(ctx, w.opts.Strategy)
	if err != nil {
		return err
	}

	if err := verifyCPFPConditions(utxos, batches, w.address); err != nil {
		return fmt.Errorf("failed to verify CPFP conditions: %w", err)
	}

	tx, err := w.buildCPFPTx(
		utxos,
		[]SpendRequest{},
		[]SendRequest{{Amount: -1, To: w.address}},
		(len(utxos)*SegwitSpendWeight)+10,
		feeStats.FeeDelta,
		requiredFeeRate,
	)
	if err != nil {
		return err
	}

	if err := w.indexer.SubmitTx(ctx, tx); err != nil {
		return err
	}

	log.Info("submitted CPFP transaction", "txid", tx.TxHash().String())
	return nil
}

func (w *batcherWallet) buildCPFPTx(utxos []UTXO, spendRequests []SpendRequest, sendRequests []SendRequest, fee, feeOverhead, feeRate int) (*wire.MsgTx, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	spendReqUTXOs, spendReqUTXOMap, balanceOfScripts, err := getUTXOsForSpendRequest(ctx, w.indexer, spendRequests)
	if err != nil {
		return nil, err
	}

	if balanceOfScripts == 0 && len(spendRequests) > 0 {
		return nil, fmt.Errorf("scripts have no funds to spend")
	}

	// build the transaction
	tx, err := buildTransaction(append(spendReqUTXOs, utxos...), sendRequests, w.address, feeOverhead+fee)
	if err != nil {
		return nil, err
	}

	// Sign the spend inputs
	err = signSpendTx(tx, spendRequests, spendReqUTXOMap, w.privateKey)

	// get the send signing script
	script, err := txscript.PayToAddrScript(w.address)
	if err != nil {
		return tx, err
	}

	// Sign the cover inputs
	// This is a no op if there are no cover utxos
	err = signSendTx(tx, utxos, len(spendReqUTXOs), script, w.privateKey)
	if err != nil {
		return tx, err
	}

	txb := btcutil.NewTx(tx)
	trueSize := mempool.GetTxVirtualSize(txb)

	newFeeEstimate := (int(trueSize) * (feeRate)) + feeOverhead

	if newFeeEstimate != fee+feeOverhead {
		return w.buildCPFPTx(utxos, spendRequests, sendRequests, newFeeEstimate, 0, feeRate)
	}

	return tx, nil
}

func verifyCPFPConditions(utxos []UTXO, batches []Batch, walletAddr btcutil.Address) error {
	ucUtxos := getUnconfirmedUtxos(utxos)
	if len(ucUtxos) == 0 {
		return ErrCPFPFeeUpdateParamsNotMet
	}
	trailingBatches, err := getTrailingBatches(batches, ucUtxos)
	if err != nil {
		return err
	}

	if len(trailingBatches) == 0 || len(trailingBatches) > 1 {
		return ErrCPFPBatchingCorrupted
	}

	return reconstructCPFPBatches(batches, trailingBatches[0], walletAddr)
}

func getUnconfirmedUtxos(utxos []UTXO) []UTXO {
	var ucUtxos []UTXO
	for _, utxo := range utxos {
		if !utxo.Status.Confirmed {
			ucUtxos = append(ucUtxos, utxo)
		}
	}
	return ucUtxos
}

func getTrailingBatches(batches []Batch, utxos []UTXO) ([]Batch, error) {
	utxomap := make(map[string]bool)
	for _, utxo := range utxos {
		utxomap[utxo.TxID] = true
	}

	batches = []Batch{}

	for _, batch := range batches {
		if _, ok := utxomap[batch.ChangeUtxo.TxID]; ok {
			batches = append(batches, batch)
		}
	}

	return batches, nil
}

func reconstructCPFPBatches(batches []Batch, trailingBatch Batch, walletAddr btcutil.Address) error {
	// todo : verify that the trailing batch can trace back to the funding utxos from wallet address
	return nil
}

func getTransaction(indexer IndexerClient, txid string) (Transaction, error) {
	for i := 1; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tx, err := indexer.GetTx(ctx, txid)
		if err != nil {
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}
		return tx, nil
	}
	return Transaction{}, ErrTxNotFound
}
