package guardian

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/catalogfi/blockchain/btc"
	"go.uber.org/zap"
)

type Wallet struct {
	client         *GuardianClient
	privKey        *btcec.PrivateKey
	cache          Cache
	indexer        btc.IndexerClient
	chainParams    *chaincfg.Params
	feeEstimator   btc.FeeEstimator
	feeLevel       btc.FeeLevel
	guardianPubkey *btcec.PublicKey
	logger         *zap.Logger
	addr           btcutil.Address
	pkScript       []byte
	rpc            btc.BitcoinRPCClient
}

const indexerTimeout = 3 * time.Second

func NewWallet(client *GuardianClient, privKey *btcec.PrivateKey, cache Cache, indexer btc.IndexerClient, chainParams *chaincfg.Params, feeEstimator btc.FeeEstimator, logger *zap.Logger, bitcoinRPC btc.BitcoinRPCClient, feeLevel ...btc.FeeLevel) (*Wallet, error) {
	defaultFeeLevel := btc.MediumFee
	if len(feeLevel) > 0 {
		defaultFeeLevel = feeLevel[0]
	}

	logger.Debug("Guardian wallet created or fetched", zap.String("address", client.GetAccount().Address))

	account := client.GetAccount()

	guardianPubkeyBytes, err := hex.DecodeString(account.GuardianPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode guardian pubkey: %w", err)
	}

	guardianPubkey, err := btcec.ParsePubKey(guardianPubkeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse guardian pubkey: %w", err)
	}

	addr, err := btcutil.DecodeAddress(account.Address, chainParams)
	if err != nil {
		return nil, fmt.Errorf("failed to decode address: %w", err)
	}

	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to get pk script: %w", err)
	}

	return &Wallet{
		client:         client,
		privKey:        privKey,
		cache:          cache,
		indexer:        indexer,
		chainParams:    chainParams,
		feeEstimator:   feeEstimator,
		feeLevel:       defaultFeeLevel,
		guardianPubkey: guardianPubkey,
		logger:         logger,
		addr:           addr,
		pkScript:       pkScript,
		rpc:            bitcoinRPC,
	}, nil
}

func checkForDups(req []btc.SendRequest) error {
	reqMap := make(map[string]bool)
	for _, r := range req {
		ok, id := r.ID()
		if !ok {
			return fmt.Errorf("request id not found")
		}
		reqMap[id] = true
	}

	if len(reqMap) != len(req) {
		return fmt.Errorf("request ids contain duplicates")
	}
	return nil
}

func validateRequests(req []btc.SendRequest, blacklistedAddrs ...btcutil.Address) error {

	// lets make sure reqs contain no duplicates
	if err := checkForDups(req); err != nil {
		return fmt.Errorf("request ids contain duplicates: %w", err)
	}

	for _, r := range req {
		if r.Amount < btc.DustAmount {
			return fmt.Errorf("amount smaller than dust amount")
		}
		for _, addr := range blacklistedAddrs {
			if r.To.String() == addr.String() {
				return fmt.Errorf("can not send to self")
			}
		}
	}
	return nil
}

func (w *Wallet) submitTx(ctx context.Context, tx *wire.MsgTx) error {
	// log the tx hex
	txHex, err := btc.GetTxRawBytes(tx)
	if err != nil {
		return fmt.Errorf("failed to get tx raw bytes: %w", err)
	}
	w.logger.Info("Submitting tx", zap.String("txHex", hex.EncodeToString(txHex)))

	err = w.indexer.SubmitTx(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to submit tx: %w", err)
	}

	return nil
}

func (w *Wallet) shouldCreateNewBatch(ctx context.Context) (*Batch, bool, error) {
	// get the latest batch
	batch, err := w.cache.ReadLatestBatch(ctx)
	if err != nil {
		// if batch is not found, we need to broadcast the tx
		if strings.Contains(err.Error(), "not found") {
			w.logger.Debug("Ongoing batch not found, broadcasting tx")
			return nil, true, nil
		}
		return nil, false, fmt.Errorf("failed to get latest batch: %w", err)
	}

	confirmed, _, err := w.batchStatus(ctx, batch)
	if err != nil {
		return nil, false, fmt.Errorf("failed to check if batch is confirmed: %w", err)
	}

	return batch, confirmed, nil
}

func (w *Wallet) getVoutIndicesForMergeRequests(mergeRequestIDs []string, onGoingBatch *Batch) ([]int, error) {
	indices := []int{}
	for _, id := range mergeRequestIDs {
		vout, ok := onGoingBatch.RequestIds[id]
		if !ok {
			return nil, fmt.Errorf("failed to read request vout from batch")
		}
		indices = append(indices, vout)
	}
	return indices, nil
}

func (w *Wallet) handleMergeRequests(tx *wire.MsgTx, req []btc.SendRequest, onGoingBatch *Batch) (*wire.MsgTx, []*TxOutput, []string, error) {
	// check for merges, and if there are there, we need to remove them from the tx

	remainingRequests, mergeIDs, mergeTxHexes, err := extractMergeIDs(req, onGoingBatch)
	if err != nil {
		w.logger.Error("Failed to extract merge ids", zap.Error(err))
		return nil, nil, nil, fmt.Errorf("failed to extract merge ids: %w", err)
	}

	prettyPrint(mergeIDs)

	indices, err := w.getVoutIndicesForMergeRequests(mergeIDs, onGoingBatch)
	if err != nil {
		w.logger.Error("Failed to get vout indices for merge requests", zap.Error(err))
		return nil, nil, nil, fmt.Errorf("failed to get vout indices for merge requests: %w", err)
	}

	// this tx will have more inputs as we remove some outputs
	tx, err = removeIndicesFromTx(tx, indices)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to remove indices from tx: %w", err)
	}

	newOuts, err := parseIntoTxOutputs(remainingRequests)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to extract ins and outs from tx: %w", err)
	}

	return tx, newOuts, mergeTxHexes, nil
}

func (w *Wallet) Send(ctx context.Context, req []btc.SendRequest) (chainhash.Hash, error) {
	if len(req) == 0 {
		return chainhash.Hash{}, fmt.Errorf("no request found")
	}

	if err := validateRequests(req, w.addr); err != nil {
		return chainhash.Hash{}, err
	}

	onGoingBatch, shouldCreateNewBatch, err := w.shouldCreateNewBatch(ctx)
	if err != nil {
		return chainhash.Hash{}, fmt.Errorf("failed to check if we should create a new batch: %w", err)
	}
	if shouldCreateNewBatch {
		return w.batchAndBroadcast(ctx, req, onGoingBatch)
	}

	tx, err := w.txFromBatch(ctx, onGoingBatch)
	if err != nil {
		if strings.Contains(err.Error(), "Transaction not found") {
			return w.batchAndBroadcast(ctx, req, onGoingBatch)
		}
		return chainhash.Hash{}, fmt.Errorf("failed to get raw tx from batch: %w", err)
	}

	tx, newTxOuts, mergeTxHexes, err := w.handleMergeRequests(tx, req, onGoingBatch)
	if err != nil {
		return chainhash.Hash{}, fmt.Errorf("failed to handle merge requests: %w", err)
	}

	tx, err = w.tryAdjustingAmountsForNewOuts(ctx, tx, newTxOuts)
	if err != nil {
		if strings.Contains(err.Error(), "Transaction not found") {
			return w.batchAndBroadcast(ctx, req, onGoingBatch)
		}
		return chainhash.Hash{}, fmt.Errorf("failed to adjust amounts for new outs: %w", err)
	}

	values, totalInAmount, err := w.getInAmounts(ctx, tx)
	if err != nil {
		if strings.Contains(err.Error(), "Transaction not found") {
			return w.batchAndBroadcast(ctx, req, onGoingBatch)
		}
		return chainhash.Hash{}, fmt.Errorf("failed to get in amounts: %w", err)
	}

	previousFeeRate := calculateFeeRate(int64(onGoingBatch.Tx.Weight), onGoingBatch.Tx.Fee)
	tx, err = w.adjustFee(ctx, tx, totalInAmount, previousFeeRate, onGoingBatch.Tx.Fee)
	if err != nil {
		return chainhash.Hash{}, fmt.Errorf("failed to adjust fee: %w", err)
	}

	fee, err := w.rpc.GetDescendants(ctx, onGoingBatch.Tx.TxID)
	if err != nil {
		return chainhash.Hash{}, fmt.Errorf("failed to get descendants: %w", err)
	}

	tx, err = w.includeMergeTxFee(tx, fee)
	if err != nil {
		if strings.Contains(err.Error(), "Transaction not found") {
			return w.batchAndBroadcast(ctx, req, onGoingBatch)
		}
		return chainhash.Hash{}, fmt.Errorf("failed to include merge tx fee: %w", err)
	}

	// at this point, we have added new requests to the tx and adjusted the fee
	tx, err = w.signTx(tx, values)
	if err != nil {
		return chainhash.Hash{}, fmt.Errorf("failed to sign tx: %w", err)
	}

	prevTxID, err := chainhash.NewHashFromStr(onGoingBatch.Tx.TxID)
	if err != nil {
		return chainhash.Hash{}, fmt.Errorf("failed to get prev tx id: %w", err)
	}

	err = w.client.UpdateTransaction(ctx, *prevTxID, tx, mergeTxHexes, values)
	if err != nil {
		return chainhash.Hash{}, fmt.Errorf("failed to sign tx: %w", err)
	}

	err = w.submitTx(ctx, tx)
	if err != nil {
		time.Sleep(1 * time.Second)
		// could be that previous batch got confirmed.
		confirmed, notFound, statusErr := w.batchStatus(ctx, onGoingBatch)
		if statusErr != nil {
			return chainhash.Hash{}, fmt.Errorf("failed to get batch status: %w", statusErr)
		}
		w.logger.Info("Batch status", zap.Bool("confirmed", confirmed), zap.Bool("notFound", notFound))
		if strings.Contains(err.Error(), "bad-txns-inputs-missingorspent") || notFound || confirmed {
			return w.batchAndBroadcast(ctx, req, onGoingBatch)
		}
		return chainhash.Hash{}, fmt.Errorf("failed to submit tx: %w", err)
	}

	err = w.saveLatestBatch(ctx, req, tx, onGoingBatch)
	return tx.TxHash(), err

}

// removeIndicesFromTx removes the indices from the tx and pushes the remaining txouts to the new tx
func removeIndicesFromTx(tx *wire.MsgTx, indices []int) (*wire.MsgTx, error) {

	newTx := wire.NewMsgTx(tx.Version)

	// add all inputs to the new tx
	for _, in := range tx.TxIn {
		newTx.AddTxIn(in)
	}

	// add all outputs to the new tx
	for i, out := range tx.TxOut {
		skip := false
		for _, index := range indices {
			if i == index {
				skip = true
				break
			}
		}
		if !skip {
			newTx.AddTxOut(out)
		}
	}

	return newTx, nil
}

// Returns the filtered valid requests and merge requests
func extractMergeIDs(req []btc.SendRequest, batch *Batch) ([]btc.SendRequest, []string, []string, error) {
	// if a request contains a invalidateID, and if it is present in the batch, keep a note of it (output index) and else if
	// it is not present in the batch, then we need to ignore the send request all together

	mergeIDs := []string{}
	mergeTxHexes := []string{}
	remainingRequests := []btc.SendRequest{}

	for _, r := range req {
		exits, id := r.InvalidateTxID()
		if !exits {
			remainingRequests = append(remainingRequests, r)
			continue
		}

		for rID := range batch.RequestIds {
			if rID == id {
				mergeIDs = append(mergeIDs, id)
				ok, mergeTxHex := r.MergeTxHex()
				if !ok {
					return nil, nil, nil, fmt.Errorf("merge tx hex not found")
				}
				mergeTxHexes = append(mergeTxHexes, mergeTxHex)
				remainingRequests = append(remainingRequests, r)
				continue
			}
		}
	}

	return remainingRequests, mergeIDs, mergeTxHexes, nil
}

func NewGaurdianScript(guardian, user *btcec.PublicKey) ([]byte, error) {
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_IF).
		AddInt64(2).
		AddData(guardian.SerializeCompressed()).
		AddData(user.SerializeCompressed()).
		AddInt64(2).
		AddOp(txscript.OP_CHECKMULTISIG).
		AddOp(txscript.OP_ELSE).
		AddInt64(144 * 180).
		AddOp(txscript.OP_CHECKLOCKTIMEVERIFY).
		AddOp(txscript.OP_DROP).
		AddData(user.SerializeCompressed()).
		AddOp(txscript.OP_CHECKSIG).
		AddOp(txscript.OP_ENDIF).
		Script()
}

func (w *Wallet) signTx(tx *wire.MsgTx, values []int64) (*wire.MsgTx, error) {

	guardianScript, err := NewGaurdianScript(w.guardianPubkey, w.privKey.PubKey())
	if err != nil {
		return nil, fmt.Errorf("failed to create guardian script: %w", err)
	}

	// for p2wpkh, we only need to add the signature and pubkey
	witness := [][]byte{
		{}, {},
		btc.AddSignatureSegwitOp,
		{0x01},
		guardianScript,
	}

	for i := range tx.TxIn {
		fetcher := txscript.NewCannedPrevOutputFetcher(guardianScript, values[i])
		err := btc.SignTx(tx, fetcher, values[i], i, witness, guardianScript, nil, txscript.SigHashAll, w.privKey)
		if err != nil {
			return nil, fmt.Errorf("failed to sign tx: %w", err)
		}
	}

	return tx, nil
}

func (w *Wallet) mergeTxFee(ctx context.Context, txHex string) (int, error) {

	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return 0, fmt.Errorf("failed to decode tx hex: %w", err)
	}

	// deserialize the tx
	txx, err := btcutil.NewTxFromBytes(txBytes)
	if err != nil {
		return 0, err
	}

	msgTx := txx.MsgTx()

	_, totalInAmount, err := w.getInAmounts(ctx, msgTx)
	if err != nil {
		return 0, fmt.Errorf("failed to get in amounts: %w", err)
	}

	totalOutAmount := w.getTotalOutAmount(msgTx)

	fee := totalInAmount - totalOutAmount

	return int(fee), nil
}

func (w *Wallet) getConfirmedBatch(ctx context.Context, lastBatchID string) (*Batch, error) {
	batch, err := w.cache.ReadBatchByTxID(ctx, lastBatchID)

	if err != nil {
		return nil, fmt.Errorf("failed to read batch by tx id: %w", err)
	}

	if batch.PreviousBatchID == "coinbase" {
		return &batch, nil
	}

	confirmed, _, err := w.batchStatus(ctx, &batch)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch status: %w", err)
	}
	if !confirmed {
		return w.getConfirmedBatch(ctx, batch.PreviousBatchID)
	}

	return &batch, nil
}

func (w *Wallet) getUnconfirmedRequests(ctx context.Context, previousBatch *Batch) ([]btc.SendRequest, error) {
	confirmedBatch, err := w.getConfirmedBatch(ctx, previousBatch.Tx.TxID)
	if err != nil {
		return nil, fmt.Errorf("failed to get confirmed batch: %w", err)
	}

	tx, err := w.txFromBatch(ctx, confirmedBatch)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx from batch: %w", err)
	}

	requestsToAdd := []btc.SendRequest{}

	for _, pr := range previousBatch.Tx.VOUTs {
		// see if pr is present in the tx
		found := false
		for _, out := range tx.TxOut {

			scriptPubKey, err := hex.DecodeString(pr.ScriptPubKey)
			if err != nil {
				return nil, fmt.Errorf("failed to decode script pubkey: %w", err)
			}
			if bytes.Equal(out.PkScript, scriptPubKey) {
				found = true
				break
			}
		}
		if !found {
			addr, err := btcutil.DecodeAddress(pr.ScriptPubKeyAddress, w.chainParams)
			if err != nil {
				return nil, fmt.Errorf("failed to decode address: %w", err)
			}
			requestsToAdd = append(requestsToAdd, btc.NewSendRequest(pr.ScriptPubKeyAddress, int64(pr.Value), addr))
		}
	}
	return requestsToAdd, nil
}

func (w *Wallet) getReqIDs(req []btc.SendRequest) ([]string, error) {
	reqIDs := []string{}
	for _, r := range req {
		ok, id := r.ID()
		if !ok {
			return nil, fmt.Errorf("request id not found")
		}
		reqIDs = append(reqIDs, id)
	}
	return reqIDs, nil
}

// TODO: never maintain duplicates
func (w *Wallet) batchAndBroadcast(ctx context.Context, req []btc.SendRequest, previousBatch *Batch) (chainhash.Hash, error) {
	if previousBatch != nil {
		requestsToAdd, err := w.getUnconfirmedRequests(ctx, previousBatch)
		if err != nil {
			return chainhash.Hash{}, fmt.Errorf("failed to get unconfirmed requests: %w", err)
		}
		req = append(req, requestsToAdd...)
	}

	// remove all the requests which have mergeIDs
	remainingRequests := []btc.SendRequest{}
	for _, r := range req {
		if ok, _ := r.InvalidateTxID(); !ok {
			remainingRequests = append(remainingRequests, r)
		}
	}

	if len(remainingRequests) == 0 {
		return chainhash.Hash{}, fmt.Errorf("no remaining requests to process")
	}

	outs, err := parseIntoTxOutputs(remainingRequests)
	if err != nil {
		return chainhash.Hash{}, fmt.Errorf("failed to parse into tx outs: %w", err)
	}

	tx := wire.NewMsgTx(2)
	totalOutAmount := int64(0)
	for _, out := range outs {
		tx.AddTxOut(out.TxOut)
		totalOutAmount += out.Value
	}

	tx, err = w.selectAndAddUTXOsForNewOuts(ctx, tx, outs)
	if err != nil {
		return chainhash.Hash{}, fmt.Errorf("failed to select and add utxos for new outs: %w", err)
	}

	inValues, totalInAmount, err := w.getInAmounts(ctx, tx)
	if err != nil {
		return chainhash.Hash{}, fmt.Errorf("failed to get in amounts: %w", err)
	}

	tx, err = w.adjustFee(ctx, tx, totalInAmount, 0, 0)
	if err != nil {
		return chainhash.Hash{}, fmt.Errorf("failed to adjust fee: %w", err)
	}

	tx, err = w.signTx(tx, inValues)
	if err != nil {
		return chainhash.Hash{}, fmt.Errorf("failed to sign tx: %w", err)
	}

	txHex, err := btc.GetTxRawBytes(tx)
	if err != nil {
		return chainhash.Hash{}, fmt.Errorf("failed to get tx raw bytes: %w", err)
	}
	fmt.Println(hex.EncodeToString(txHex))

	currentTxHash := tx.TxHash().String()

	//  check to verify if  sign or update based on the LastFailedTxHash in previousBatch
	if previousBatch != nil && previousBatch.LastFailedTxHash != "" && previousBatch.LastFailedTxHash != currentTxHash {
		prevFailedTxHash, err := chainhash.NewHashFromStr(previousBatch.LastFailedTxHash)
		if err != nil {
			return chainhash.Hash{}, fmt.Errorf("failed to parse prevFailedTxHash: %w", err)
		}
		err = w.client.UpdateTransaction(ctx, *prevFailedTxHash, tx, []string{}, inValues)
		if err != nil {
			return chainhash.Hash{}, fmt.Errorf("failed to update tx: %w", err)
		}
	} else {
		err = w.client.SignTransaction(ctx, tx, inValues, []string{})
		if err != nil {
			return chainhash.Hash{}, fmt.Errorf("failed to sign tx: %w", err)
		}
	}

	err = w.submitTx(ctx, tx)
	if err != nil {
		if previousBatch != nil {
			previousBatch.LastFailedTxHash = tx.TxHash().String()
			UpdateErr := w.cache.SaveLatestBatch(ctx, previousBatch)
			fmt.Println(w.cache.ReadLatestBatch(ctx))
			if UpdateErr != nil {
				return chainhash.Hash{}, fmt.Errorf("failed to update failed tx id: %w", UpdateErr)
			}
		}
		return chainhash.Hash{}, fmt.Errorf("failed to submit tx: %w", err)
	}

	if previousBatch != nil {
		fmt.Println(previousBatch.RequestIds)
		previousBatch.RequestIds = map[string]int{}
	}

	err = w.saveLatestBatch(ctx, req, tx, previousBatch)
	return tx.TxHash(), err
}

func (w *Wallet) getVoutsFromTx(tx *wire.MsgTx, req []btc.SendRequest) (map[string]int, error) {
	reqToVout := make(map[string]int)
	for _, r := range req {
		ok, id := r.ID()
		if !ok {
			return reqToVout, fmt.Errorf("request id not found")
		}
		toPkScript, err := txscript.PayToAddrScript(r.To)
		if err != nil {
			return reqToVout, fmt.Errorf("failed to get to pkscript: %w", err)
		}
		idx := 0
		for i, out := range tx.TxOut {
			if bytes.Equal(out.PkScript, toPkScript) {
				idx = i
				break
			}
		}
		reqToVout[id] = idx
	}
	return reqToVout, nil
}

func calculateFeeRate(weight, feePaid int64) int64 {
	vsize := math.Ceil(float64(weight) / 4.0)
	return int64(math.Ceil(float64(feePaid) / vsize))
}

func (w *Wallet) addChangeOutput(tx *wire.MsgTx, changeAmt int64) (*wire.MsgTx, error) {
	tx.AddTxOut(wire.NewTxOut(changeAmt, w.pkScript))
	return tx, nil
}

func (w *Wallet) increaseChangeAmount(tx *wire.MsgTx, increase int64) (*wire.MsgTx, error) {
	foundOnce := false
	for i, out := range tx.TxOut {

		if bytes.Equal(out.PkScript, w.pkScript) {
			if foundOnce {
				return nil, fmt.Errorf("multiple change outputs found")
			}
			tx.TxOut[i].Value += increase
			foundOnce = true
		}

	}

	return tx, nil
}

func (w *Wallet) decreaseChangeAmount(tx *wire.MsgTx, decrease int64) (*wire.MsgTx, error) {

	foundOnce := false
	for i, out := range tx.TxOut {
		if bytes.Equal(out.PkScript, w.pkScript) {
			if foundOnce {
				return nil, fmt.Errorf("multiple change outputs found")
			}
			tx.TxOut[i].Value -= decrease
			foundOnce = true
		}
	}
	return tx, nil
}

func (w *Wallet) hasChangeOutput(tx *wire.MsgTx) bool {
	for _, out := range tx.TxOut {
		if bytes.Equal(out.PkScript, w.pkScript) {
			return true
		}
	}
	return false
}

func (w *Wallet) getVSize(tx *wire.MsgTx) int64 {
	baseSize := tx.SerializeSizeStripped()
	totalSize := tx.SerializeSize()
	if !tx.HasWitness() {
		totalSize += 267
	}

	if !w.hasChangeOutput(tx) {
		baseSize += 43
		totalSize += 43
	}
	weight := baseSize*3 + totalSize
	return int64(weight / blockchain.WitnessScaleFactor)
}

func (w *Wallet) getChangeAmount(tx *wire.MsgTx) int64 {
	for _, out := range tx.TxOut {
		if bytes.Equal(out.PkScript, w.pkScript) {
			return out.Value
		}
	}
	return 0
}

func (w *Wallet) includeMergeTxFee(tx *wire.MsgTx, mergeTxFee int64) (*wire.MsgTx, error) {
	var err error
	if mergeTxFee > 0 {
		tx, err = w.decreaseChangeAmount(tx, mergeTxFee)
		if err != nil {
			return nil, fmt.Errorf("failed to decrease change amount when merging txs: %w", err)
		}
	}
	return tx, nil
}

// adjustFee adds fee output to the tx if it doesn't exist or adjusts the change output if it does for current fee rate
func (w *Wallet) adjustFee(ctx context.Context, tx *wire.MsgTx, totalInAmount int64, previousFeeRate int64, previousFee int64) (*wire.MsgTx, error) {

	extraBaseSize := 0
	if !w.hasChangeOutput(tx) {
		extraBaseSize = 43
	}

	// this is the fee needed for the current transaction
	feeToBePaid, vsize, err := btc.EstimateGuardianFee(tx, w.feeEstimator, w.feeLevel, GuardianWitnessSize, extraBaseSize)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate fee: %w", err)
	}

	// newFee is basically the fee needed for the current transaction + the previous fee to satisfy the rbf rules
	newFee := previousFee + feeToBePaid

	// sometimes the new fee is less than the previous fee rate, so we need to adjust it
	newFeeWithPreviousFeeRate := previousFeeRate*int64(vsize) + 1

	// feeSelected is the fee that we will use for the current transaction
	feeSelected := int64(0)
	// we select the max of the two
	if newFee > newFeeWithPreviousFeeRate {
		feeSelected = newFee
	} else {
		feeSelected = newFeeWithPreviousFeeRate
	}

	// remove the change output if it exists
	tx = removeChangeOutputs(tx, w.pkScript)

	totalOutAmount := w.getTotalOutAmount(tx)
	// calculate the change needed
	changeAmt := totalInAmount - totalOutAmount - feeSelected

	if changeAmt > btc.DustAmount {
		tx, err = w.addChangeOutput(tx, changeAmt)
		if err != nil {
			return nil, fmt.Errorf("failed to add change output: %w", err)
		}
	} else if changeAmt < 0 {
		feeToFetch := feeToBePaid
		for {
			txIns, amountToBeAdded, err := w.selectUTXOsForAmount(ctx, tx, int64(feeToFetch))
			if err != nil {
				return nil, fmt.Errorf("failed to select utxos for amount: %w", err)
			}
			txCopy := tx.Copy()
			for _, txIn := range txIns {
				found := false
				for _, in := range tx.TxIn {
					if in.PreviousOutPoint.Hash.String() == txIn.PreviousOutPoint.Hash.String() && in.PreviousOutPoint.Index == txIn.PreviousOutPoint.Index {
						found = true
						break
					}
				}
				if !found {
					tx.AddTxIn(txIn)
				}
			}
			newFee, _, err := btc.EstimateGuardianFee(tx, w.feeEstimator, w.feeLevel, GuardianWitnessSize, GuardianChangeSize)
			if err != nil {
				return nil, fmt.Errorf("failed to estimate fee: %w", err)
			}

			changeAmt := totalInAmount + amountToBeAdded - int64(newFee)

			if changeAmt > btc.DustAmount {
				tx, err = w.addChangeOutput(tx, changeAmt)
				if err != nil {
					return nil, fmt.Errorf("failed to add change output: %w", err)
				}
				break
			} else if changeAmt > 0 {
				// is greater 0 but less than dust amount, so ignore it
				break
			} else {
				// there negative change, so we need utxos to fetch for the fee
				feeToFetch = newFee
				// revert the tx to the previous state
				tx = txCopy
			}

		}
	}

	return tx, nil
}

func removeChangeOutputs(tx *wire.MsgTx, pkscript []byte) *wire.MsgTx {
	newTx := wire.NewMsgTx(tx.Version)
	for _, in := range tx.TxIn {
		newTx.AddTxIn(in)
	}
	for _, out := range tx.TxOut {
		if !bytes.Equal(out.PkScript, pkscript) {
			newTx.AddTxOut(out)
		}
	}
	return newTx
}

type TxOutput struct {
	*wire.TxOut
	isMergeOutput bool
}

// tryAdjustingAmountsForNewOuts tries to adjust the amounts for the new outs
// if it is not possible, new utxos are selected and added to the tx
// Note: Fee is ignored here
func (w *Wallet) tryAdjustingAmountsForNewOuts(ctx context.Context, tx *wire.MsgTx, outs []*TxOutput) (*wire.MsgTx, error) {

	_, totalInAmount, err := w.getInAmounts(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get in amounts: %w", err)
	}

	totalOutAmount := int64(0)
	for _, out := range tx.TxOut {
		if !bytes.Equal(out.PkScript, w.pkScript) {
			totalOutAmount += out.Value
		}
	}
	neededAmount := calculateTxOutputAmount(outs)
	availableAmount := totalInAmount - totalOutAmount

	if availableAmount < neededAmount {
		tx, err = w.selectAndAddUTXOsForNewOuts(ctx, tx, outs)
		if err != nil {
			return nil, fmt.Errorf("failed to select and add utxos for new outs: %w", err)
		}
	}

	newOutTotalAmount := int64(0)
	for _, out := range outs {
		tx.AddTxOut(out.TxOut)
		if !out.isMergeOutput {
			newOutTotalAmount += out.Value
		}
	}

	for i, out := range tx.TxOut {
		if bytes.Equal(out.PkScript, w.pkScript) {
			tx.TxOut[i].Value -= newOutTotalAmount
		}
	}

	return tx, nil
}

func calculateTxOutputAmount(outs []*TxOutput) int64 {
	amount := int64(0)
	for _, out := range outs {
		amount += out.Value
	}
	return amount
}

func calculateAmount(outs []*wire.TxOut) int64 {
	amount := int64(0)
	for _, out := range outs {
		amount += out.Value
	}
	return amount
}

func prettyPrint(anything interface{}) {
	json, err := json.MarshalIndent(anything, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
}

func (w *Wallet) getInAmounts(ctx context.Context, tx *wire.MsgTx) ([]int64, int64, error) {
	amounts := []int64{}
	for _, in := range tx.TxIn {
		ctx, cancel := context.WithTimeout(ctx, indexerTimeout)
		prevTx, err := w.indexer.GetTx(ctx, in.PreviousOutPoint.Hash.String())
		cancel()
		if err != nil {
			return nil, 0, err
		}
		amounts = append(amounts, int64(prevTx.VOUTs[in.PreviousOutPoint.Index].Value))
	}
	totalAmount := int64(0)
	for _, amount := range amounts {
		totalAmount += amount
	}
	return amounts, totalAmount, nil
}

func (w *Wallet) getTotalOutAmount(tx *wire.MsgTx) int64 {
	amount := int64(0)
	for _, out := range tx.TxOut {
		amount += out.Value
	}
	return amount
}

func (w *Wallet) batchStatus(ctx context.Context, batch *Batch) (bool, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, indexerTimeout)
	defer cancel()
	tx, err := w.indexer.GetTx(ctx, batch.Tx.TxID)
	if err != nil {
		// sometimes previous batches go through and the tx is not found
		if strings.Contains(err.Error(), "Transaction not found") {
			return false, true, nil
		}
		return false, false, err
	}
	batch.Tx = tx
	return tx.Status.Confirmed, false, nil
}

func (w *Wallet) txFromBatch(ctx context.Context, batch *Batch) (*wire.MsgTx, error) {
	// get the txHex from indexer
	ctx, cancel := context.WithTimeout(ctx, indexerTimeout)
	defer cancel()
	txHex, err := w.indexer.GetTxHex(ctx, batch.Tx.TxID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx hex: %w", err)
	}

	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tx hex: %w", err)
	}

	tx, err := btcutil.NewTxFromBytes(txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx from hex: %w", err)
	}
	return tx.MsgTx(), nil
}

func parseIntoTxOutputs(req []btc.SendRequest) ([]*TxOutput, error) {
	outs := []*TxOutput{}

	isMergeRequest := func(r btc.SendRequest) bool {
		exits, _ := r.InvalidateTxID()
		return exits
	}

	for _, r := range req {
		script, err := txscript.PayToAddrScript(r.To)
		if err != nil {
			return nil, err
		}
		outs = append(outs, &TxOutput{
			TxOut: &wire.TxOut{
				Value:    int64(r.Amount),
				PkScript: script,
			},
			isMergeOutput: isMergeRequest(r),
		})
	}
	return outs, nil
}

func (w *Wallet) selectUTXOsForAmount(ctx context.Context, tx *wire.MsgTx, amount int64) ([]*wire.TxIn, int64, error) {

	utxos, _, err := w.indexer.GetUTXOsForAmount(ctx, w.addr, amount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get utxos: %w", err)
	}

	if len(utxos) == 0 {
		return nil, 0, fmt.Errorf("no utxos found")
	}

	utxosToBeSelected := []btc.UTXO{}
	amountToBeAdded := int64(0)
	for _, utxo := range utxos {
		if amountToBeAdded >= amount {
			break
		}
		found := false
		for _, in := range tx.TxIn {
			if in.PreviousOutPoint.Hash.String() == utxo.TxID && in.PreviousOutPoint.Index == utxo.Vout {
				found = true
				break
			}
		}
		if !found {
			utxosToBeSelected = append(utxosToBeSelected, utxo)
			amountToBeAdded += utxo.Amount
		}
	}

	if amountToBeAdded < amount {
		return nil, 0, fmt.Errorf("insufficient balance : has %d and need %d", amountToBeAdded, amount)
	}

	txIns := []*wire.TxIn{}
	for _, utxo := range utxosToBeSelected {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get hash from str: %w", err)
		}
		txIns = append(txIns, &wire.TxIn{
			PreviousOutPoint: wire.OutPoint{
				Hash:  *hash,
				Index: utxo.Vout,
			},
			Sequence: wire.MaxTxInSequenceNum - 2,
		})
	}
	return txIns, amountToBeAdded, nil
}

func (w *Wallet) Address() btcutil.Address {
	return w.addr
}

func (w *Wallet) selectAndAddUTXOsForNewOuts(ctx context.Context, tx *wire.MsgTx, outs []*TxOutput) (*wire.MsgTx, error) {
	amount := calculateTxOutputAmount(outs)
	txIns, _, err := w.selectUTXOsForAmount(ctx, tx, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to select utxos for amount: %w", err)
	}
	for _, txIn := range txIns {
		tx.AddTxIn(txIn)
	}

	return tx, nil
}

func (w *Wallet) saveLatestBatch(ctx context.Context, req []btc.SendRequest, tx *wire.MsgTx, workingBatch *Batch) error {

	txCtx, cancel := context.WithTimeout(ctx, indexerTimeout)
	defer cancel()
	batchTx, err := w.indexer.GetTx(txCtx, tx.TxHash().String())
	if err != nil {
		if strings.Contains(err.Error(), "Transaction not found") {
			_, err = w.batchAndBroadcast(ctx, req, workingBatch)
			return err
		}
		return fmt.Errorf("failed to get tx: %w", err)
	}
	reqIDs, err := w.getVoutsFromTx(tx, req)
	if err != nil {
		return fmt.Errorf("failed to save vouts: %w", err)
	}

	var finalBatch *Batch
	// Creating A New Batch

	if workingBatch == nil {
		finalBatch = NewBatch(batchTx, reqIDs, CoinbaseBatchID, 0)
	} else {
		workingBatch.PreviousBatchID = workingBatch.Tx.TxID
		workingBatch.Tx = batchTx
		for reqID, vout := range reqIDs {
			workingBatch.RequestIds[reqID] = vout
		}
		finalBatch = workingBatch
	}

	err = w.cache.SaveLatestBatch(txCtx, finalBatch)
	if err != nil {
		return fmt.Errorf("failed to save batch: %w", err)
	}

	return nil
}
