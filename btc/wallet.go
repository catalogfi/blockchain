package btc

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

var (
	// ErrAmountLessThanDust indicates that the amount is less than the dust limit (546 satoshis).
	ErrAmountLessThanDust = fmt.Errorf("amount is less than dust")

	// ErrSACPInvalidHashType indicates that the hash type for the SACP is invalid.
	// Use SigHashSingleAnyoneCanPay for generating SACP.
	ErrSACPInvalidHashType = fmt.Errorf("invalid hash type, use SigHashSingleAnyoneCanPay")

	// ErrFeeExceedsValue indicates that the fee exceeds the value of the script.
	// Occurs when the fee required to spend the script is more than the value of the script.
	ErrFeeExceedsValue = fmt.Errorf("fee exceeds the value of the script")

	// ErrNoUTXOsForRequests indicates that even after multiple attempts, could not find enough utxos to satisfy the requests.
	ErrNoUTXOsForRequests = fmt.Errorf("could not find enough utxos to satisfy the requests")

	// ErrNoFundsToSpend indicates that the scripts have no funds to spend.
	// Occurs when the spend requests have no funds to spend.
	ErrNoFundsToSpend = fmt.Errorf("scripts have no funds to spend")

	// ErrScriptNotFoundInWitness indicates that the script is not found in the witness while validating the witness.
	// Script should be always present at the end of the witness for p2wsh scripts and
	// at the second last for p2tr scripts.
	ErrScriptNotFoundInWitness = fmt.Errorf("script not found in the witness")

	// ErrInvalidP2wpkhScript indicates that the script is invalid for p2wpkh spend.
	ErrInvalidP2wpkhScript = fmt.Errorf("invalid script")

	// ErrInvalidP2wpkhWitness indicates that the witness is invalid for the given p2wpkh script.
	// use AddSignatureSegwitOp and AddPubkeyCompressedOp for p2wpkh scripts.
	ErrInvalidP2wpkhWitness = fmt.Errorf("invalid p2wpkh witness")

	// ErrInvalidSigOp indicates that the signature operation is invalid.
	// Use AddSignatureSegwitOp for segwit v0 scripts and AddSignatureSchnorrOp for segwit v1 scripts.
	ErrInvalidSigOp = func(op string) error {
		return fmt.Errorf("invalid sig, use %s", op)
	}

	// ErrInvalidPubkeyOp indicates that the pubkey operation is invalid.
	// Use AddPubkeyCompressedOp for compressed pubkeys and AddXOnlyPubkeyOp for xonly pubkeys.
	ErrInvalidPubkeyOp = func(op string) error {
		return fmt.Errorf("invalid pubkey, use %s", op)
	}

	// ErrInsufficientFunds indicates that the funds are insufficient to spend the script.
	ErrInsufficientFunds = func(have int64, need int64) error {
		return fmt.Errorf("insufficient funds: have %d, need %d", have, need)
	}

	// ErrNoUTXOsFoundForAddress indicates that no utxos are found for the given address.
	ErrNoUTXOsFoundForAddress = func(addr string) error {
		return fmt.Errorf("utxos not found for address %s", addr)
	}

	// ErrNoInputsFound indicates that no inputs are found in the sacp.
	ErrSCAPNoInputsFound = fmt.Errorf("no inputs found in sacp")

	// ErrSCAPInputsNotEqualOutputs indicates that the number of inputs and outputs are not equal in the sacp.
	ErrSCAPInputsNotEqualOutputs = fmt.Errorf("number of inputs and outputs are not equal in sacp")
)

var (
	// Adds a 64 or 65 byte schnorr signature based on hash type
	AddSignatureSchnorrOp = []byte("add_signature_segwit_v1")
	// Adds a 72 byte ECDSA signature
	AddSignatureSegwitOp = []byte("add_signature_segwit_v0")
	// Adds a 33 byte compressed pubkey
	AddPubkeyCompressedOp = []byte("add_pubkey_compressed")
	// Adds a 32 byte xonly pubkey
	AddXOnlyPubkeyOp = []byte("add_xonly_pubkey")
)

type utxoMap map[string]UTXOs

type SpendRequest struct {
	// Witness required to spend the script.
	// If the script requires a signature or pubkey,
	// use AddSignatureSegwitOp and AddPubkeyCompressedOp respectively.
	Witness [][]byte

	// Script to spend.
	// Only used for other than tapscript spend.
	// Use Leaf for tapscript spend.
	Script []byte

	// Revealing leaf in tapscript spend
	Leaf txscript.TapLeaf

	// Address of the script
	ScriptAddress btcutil.Address

	// Hash type for the signature. Default is SigHashDefault (tapscript only)
	HashType txscript.SigHashType

	// Sequence number for the input
	Sequence uint32

	// UTXO to spend
	Utxos UTXOs
}

type SendRequest struct {
	// Amount to send
	Amount int64
	// Recipient address
	To btcutil.Address
}

type Wallet interface {
	// Address returns the Pay-to-Witness-Public-Key-Hash (P2WPKH) address of the wallet.
	Address() btcutil.Address

	// Spend funds from multiple scripts and send funds to multiple recipients at the same time in a
	// single transaction.
	// Returns the transaction ID (txid) and an error.
	//
	// Example:
	// Send funds to multiple recipients
	// 	txid, err := wallet.Send(ctx, []btc.SendRequest{
	// 		{Amount: 1000, To: recipient1},
	// 		{Amount: 2000, To: recipient2},
	// 	}, nil)
	//
	// Spend funds from multiple scripts
	// Spending will move funds from the scripts to the wallet's address.
	// 	txid, err := wallet.Send(ctx, nil, []btc.SpendRequest{
	// 		{
	// 			Witness: [][]byte{
	// 				{btc.AddSignatureSegwitOp},
	// 				{btc.AddPubkeyCompressedOp},
	// 			},
	// 			Script: script,
	// 			ScriptAddress: scriptAddress,
	// 		},
	// 	}, nil)
	Send(ctx context.Context, sendReq []SendRequest, spendReq []SpendRequest, sacps [][]byte) (string, error)

	// Generate the signature necessary to spend the script.
	//
	// Uses the hash type SigHashSingleAnyoneCanPay.
	GenerateSACP(ctx context.Context, spendReq SpendRequest, to btcutil.Address) ([]byte, error)

	// SignSACPTx generates a schnorr signature for the given details.
	// Returns the witness containing the schnorr signature and an error if any.
	//
	// tx is not mutated instead a new copy is created internally to perform the signing.
	SignSACPTx(tx *wire.MsgTx, idx int, amount int64, leaf txscript.TapLeaf, scriptAddr btcutil.Address, witness [][]byte) ([][]byte, error)

	// Status checks the status of a transaction using its transaction ID (txid).
	// Returns the transaction and a boolean indicating whether the transaction is submitted or not and an error
	Status(ctx context.Context, id string) (Transaction, bool, error)

	// Signs cover utxos
	SignCoverUTXOs(tx *wire.MsgTx, utxos UTXOs, startingIdx int) error

	// Weight of the covering UTXO
	CoverUTXOSpendWeight() int
}

// SimpleWallet is a Wallet implementation that can send and spend funds.
type SimpleWallet struct {
	privateKey   *btcec.PrivateKey
	indexer      IndexerClient
	feeEstimator FeeEstimator
	chainParams  *chaincfg.Params
	signerAddr   btcutil.Address
	feeLevel     FeeLevel
}

// Generates a new p2wpkh simple wallet
func NewSimpleWallet(privKey *btcec.PrivateKey, chainParams *chaincfg.Params, indexer IndexerClient, feeEstimator FeeEstimator, feeLevel FeeLevel) (Wallet, error) {
	address, err := PublicKeyAddress(chainParams, waddrmgr.WitnessPubKey, privKey.PubKey())
	if err != nil {
		return nil, err
	}

	return &SimpleWallet{
		indexer:      indexer,
		signerAddr:   address,
		privateKey:   privKey,
		chainParams:  chainParams,
		feeEstimator: feeEstimator,
		feeLevel:     feeLevel,
	}, nil
}

// Returns the address of the wallet.
func (sw *SimpleWallet) Address() btcutil.Address {
	return sw.signerAddr
}

func (sw *SimpleWallet) Send(ctx context.Context, sendRequests []SendRequest, spendRequests []SpendRequest, sacps [][]byte) (string, error) {
	//TODO: Estimate a fee based on requests
	fee := 1000

	// validate the requests
	err := validateRequests(spendRequests, sendRequests, sacps)
	if err != nil {
		return "", err
	}

	sacpsFee, err := getFeeUsedInSACPs(ctx, sacps, sw.indexer)
	if err != nil {
		return "", err
	}

	tx, err := sw.spendAndSend(ctx, sendRequests, spendRequests, sacps, sacpsFee, fee, 0)
	if err != nil {
		return "", err
	}

	return submitTx(ctx, sw.indexer, tx)
}

// SignSACPTx generates a schnorr signature for the given details.
func (sw *SimpleWallet) SignSACPTx(tx *wire.MsgTx, idx int, amount int64, leaf txscript.TapLeaf, scriptAddr btcutil.Address, witness [][]byte) ([][]byte, error) {
	cTx := tx.Copy()
	script, err := txscript.PayToAddrScript(scriptAddr)
	if err != nil {
		return nil, err
	}

	fetcher := txscript.NewCannedPrevOutputFetcher(script, amount)
	err = signTx(cTx, fetcher, amount, idx, witness, script, &leaf, SigHashSingleAnyoneCanPay, sw.privateKey)
	if err != nil {
		return nil, err
	}
	return cTx.TxIn[idx].Witness, nil
}

// GenerateSACP generates a SACP tx for the given spend request.
func (sw *SimpleWallet) GenerateSACP(ctx context.Context, spendReq SpendRequest, to btcutil.Address) ([]byte, error) {
	if spendReq.HashType != txscript.SigHashDefault && spendReq.HashType != SigHashSingleAnyoneCanPay {
		return nil, ErrSACPInvalidHashType
	}
	if spendReq.HashType == txscript.SigHashDefault {
		spendReq.HashType = SigHashSingleAnyoneCanPay
	}

	if err := validateRequests([]SpendRequest{spendReq}, nil, nil); err != nil {
		return nil, err
	}
	if to == nil {
		to = sw.signerAddr
	}
	return sw.generateSACP(ctx, spendReq, to, 1000)
}

func (sw *SimpleWallet) generateSACP(ctx context.Context, spendRequest SpendRequest, to btcutil.Address, fee int64) ([]byte, error) {

	// get the utxos for the script
	utxos, utxoMap, utxosBalance, err := getUTXOsForSpendRequest(ctx, sw.indexer, []SpendRequest{spendRequest})
	if err != nil {
		return nil, err
	}

	if utxosBalance-int64(fee) <= 0 {
		return nil, ErrFeeExceedsValue
	}

	sequenceMap := generateSequenceMap(utxoMap, []SpendRequest{spendRequest})

	// build the transaction with no recipients or sacps
	tx, _, err := buildTransaction(utxos, nil, nil, to, fee, sequenceMap)
	if err != nil {
		return nil, err
	}

	// sign the transaction
	err = signSpendTx(ctx, tx, 0, []SpendRequest{spendRequest}, utxoMap, sw.indexer, sw.privateKey)
	if err != nil {
		return nil, err
	}

	// estimate the fee required to make the transaction
	feeToBePaid, err := EstimateSegwitFee(tx, sw.feeEstimator, sw.feeLevel)
	if err != nil {
		return nil, err
	}

	if int64(feeToBePaid) > fee {
		return sw.generateSACP(ctx, spendRequest, to, int64(feeToBePaid))
	}

	// serialize the transaction
	var txBytes []byte
	if txBytes, err = GetTxRawBytes(tx); err != nil {
		return nil, err
	}
	return txBytes, nil
}

func (sw *SimpleWallet) spendAndSend(ctx context.Context, sendRequests []SendRequest, spendRequests []SpendRequest, sacps [][]byte, sacpFee, fee int, depth int) (*wire.MsgTx, error) {

	// This means we made 100 recursive calls and still could not find enough utxos to send the amount
	if depth > 100 {
		return nil, ErrNoUTXOsForRequests
	}

	// spendUTXOs are the UTXOs used to spend the scripts
	// coverUTXOs are the UTXOs used to cover the remaining amount required to send
	// utxoMap is a map of script address to UTXOs
	spendUTXOs, coverUTXOs, utxoMap, err := getUTXOsForRequests(ctx, sw.indexer, spendRequests, sendRequests, sw.signerAddr, fee, sacpFee)
	if err != nil {
		return nil, err
	}

	// generate sequence map (used to set sequence number for each input)
	sequenceMap := generateSequenceMap(utxoMap, spendRequests)

	// build the transaction
	// Signing index indicates from which index we need to start signing the transaction
	tx, signingIdx, err := buildTransaction(append(spendUTXOs, coverUTXOs...), sacps, sendRequests, sw.signerAddr, int64(fee), sequenceMap)
	if err != nil {
		return nil, err
	}

	// Sign the spend inputs
	err = signSpendTx(ctx, tx, signingIdx, spendRequests, utxoMap, sw.indexer, sw.privateKey)
	if err != nil {
		return nil, err
	}

	// Sign the cover inputs
	// This is a no op if there are no cover utxos
	err = sw.SignCoverUTXOs(tx, coverUTXOs, signingIdx+len(spendRequests))
	if err != nil {
		return nil, err
	}

	// estimate the fee required to make the transaction
	feeToBePaid, err := EstimateSegwitFee(tx, sw.feeEstimator, sw.feeLevel)
	if err != nil {
		return nil, err
	}

	// sacpFee is the fee used in the SACPs
	// This could be zero if there are no SACPs or SACPs have no fee
	feeToBePaid -= sacpFee

	if feeToBePaid > fee {
		return sw.spendAndSend(ctx, sendRequests, spendRequests, sacps, sacpFee, feeToBePaid, depth+1)
	}

	return tx, nil
}

// Status checks the status of a transaction using its transaction ID (txid).
func (sw *SimpleWallet) Status(ctx context.Context, id string) (Transaction, bool, error) {
	tx, err := sw.indexer.GetTx(ctx, id)

	if err != nil {
		if strings.Contains(err.Error(), "Transaction not found") {
			return Transaction{}, false, nil
		}
		return Transaction{}, false, err
	}
	return tx, true, nil
}

// Signs the send transaction (p2wpkh spend).
// Use startingIdx to start signing from a specific index
func (w *SimpleWallet) SignCoverUTXOs(tx *wire.MsgTx, utxos UTXOs, startingIdx int) error {
	// get the send signing script
	script, err := txscript.PayToAddrScript(w.signerAddr)
	if err != nil {
		return err
	}

	// for p2wpkh, we only need to add the signature and pubkey
	witness := [][]byte{
		AddSignatureSegwitOp,
		AddPubkeyCompressedOp,
	}
	idx := startingIdx
	for i := range utxos {
		fetcher := txscript.NewCannedPrevOutputFetcher(script, utxos[i].Amount)
		err := signTx(tx, fetcher, utxos[i].Amount, idx, witness, script, nil, txscript.SigHashAll, w.privateKey)
		if err != nil {
			return err
		}
		idx++
	}
	return nil
}

func (w *SimpleWallet) CoverUTXOSpendWeight() int {
	return SegwitSpendWeight
}

// ------------------ Helper functions ------------------

// getSACPAmounts returns the total input and output amounts for the given SACPs
func getSACPAmounts(ctx context.Context, sacps [][]byte, indexer IndexerClient) (int64, int64, error) {
	tx, _, err := buildTxFromSacps(sacps)
	if err != nil {
		return 0, 0, err
	}

	// go through each input and get the amount it holds
	// add all the inputs and subtract the outputs to get the fee
	totalInputAmount := int64(0)
	for _, in := range tx.TxIn {
		txFromIndexer, err := indexer.GetTx(ctx, in.PreviousOutPoint.Hash.String())
		if err != nil {
			return 0, 0, err
		}
		totalInputAmount += int64(txFromIndexer.VOUTs[in.PreviousOutPoint.Index].Value)
	}

	totalOutputAmount := int64(0)
	for _, out := range tx.TxOut {
		totalOutputAmount += out.Value
	}

	return totalInputAmount, totalOutputAmount, nil
}

// getFeeUsedInSACPs returns the amount of fee used in the given SACPs
func getFeeUsedInSACPs(ctx context.Context, sacps [][]byte, indexer IndexerClient) (int, error) {
	totalInputAmount, totalOutputAmount, err := getSACPAmounts(ctx, sacps, indexer)
	if err != nil {
		return 0, err
	}
	return int(totalInputAmount - totalOutputAmount), nil
}

func submitTx(ctx context.Context, indexer IndexerClient, tx *wire.MsgTx) (string, error) {
	txid := tx.TxHash().String()
	err := indexer.SubmitTx(ctx, tx)
	if err != nil {
		return "", err
	}
	return txid, nil
}

// getPrevoutsForSACPs returns the previous outputs and txouts for the given SACPs used to build the prevOutFetcher
func getPrevoutsForSACPs(ctx context.Context, tx *wire.MsgTx, endingSACPIdx int, indexer IndexerClient) ([]wire.OutPoint, []*wire.TxOut, error) {
	prevouts := []wire.OutPoint{}
	txOuts := []*wire.TxOut{}
	for i := 0; i < endingSACPIdx; i++ {
		outpoint := tx.TxIn[i].PreviousOutPoint
		prevouts = append(prevouts, outpoint)

		// The only way to get the txOut is to get the transaction from the indexer
		txFromIndexer, err := indexer.GetTx(ctx, outpoint.Hash.String())
		if err != nil {
			return nil, nil, err
		}
		pkScriptHex := txFromIndexer.VOUTs[outpoint.Index].ScriptPubKey
		pkScript, err := hex.DecodeString(pkScriptHex)
		if err != nil {
			return nil, nil, err
		}
		txOuts = append(txOuts, &wire.TxOut{PkScript: pkScript, Value: int64(txFromIndexer.VOUTs[outpoint.Index].Value)})
	}
	return prevouts, txOuts, nil
}

// getUTXOsForRequests returns the UTXOs required to spend the scripts and cover the send amount.
func getUTXOsForRequests(ctx context.Context, indexer IndexerClient, spendReqs []SpendRequest, sendReqs []SendRequest, feePayer btcutil.Address, fee, sacpFee int) (UTXOs, UTXOs, utxoMap, error) {

	spendUTXOs, spendUTXOsMap, balanceOfScripts, err := getUTXOsForSpendRequest(ctx, indexer, spendReqs)
	if err != nil {
		return nil, nil, nil, err
	}

	// coverUTXOs are the UTXOs used to cover the remaining amount required to send
	var coverUTXOs UTXOs
	totalSendAmount := calculateTotalSendAmount(sendReqs)
	if balanceOfScripts <= totalSendAmount && sacpFee <= fee {
		utxos, _, err := indexer.GetUTXOsForAmount(ctx, feePayer, totalSendAmount-balanceOfScripts+int64(fee))
		if err != nil {
			return nil, nil, nil, err
		}
		coverUTXOs = append(coverUTXOs, utxos...)
	}

	// SpendUTXOMap only contains the UTXOs of the scripts but not the cover UTXOs
	// We need to add the cover UTXOs to the map
	spendUTXOsMap[feePayer.EncodeAddress()] = append(spendUTXOsMap[feePayer.EncodeAddress()], coverUTXOs...)

	return spendUTXOs, coverUTXOs, spendUTXOsMap, nil
}

// generateSequenceMap returns a map of txid to sequence number for the given spend requests
func generateSequenceMap(utxosMap utxoMap, spendRequest []SpendRequest) map[string]uint32 {
	sequencesMap := make(map[string]uint32)
	for _, req := range spendRequest {
		utxos, ok := utxosMap[req.ScriptAddress.EncodeAddress()]
		if !ok {
			continue
		}
		for _, utxo := range utxos {
			if req.Sequence != 0 {
				sequencesMap[utxo.TxID] = req.Sequence
			}
		}
	}
	return sequencesMap
}

// Merge multiple sacps into a single transaction
func buildTxFromSacps(sacps [][]byte) (*wire.MsgTx, int, error) {
	idx := 0
	tx := wire.NewMsgTx(DefaultTxVersion)
	for _, sacp := range sacps {
		sacpTx, err := buildAndValidateSacpTx(sacp)
		if err != nil {
			return nil, 0, err
		}
		for _, in := range sacpTx.TxIn {
			tx.AddTxIn(in)
			idx++
		}

		for _, out := range sacpTx.TxOut {
			tx.AddTxOut(out)
		}
	}
	return tx, idx, nil
}

// Builds an unsigned transaction with the given utxos, recipients, change address and fee.
func buildTransaction(utxos UTXOs, sacps [][]byte, recipients []SendRequest, changeAddr btcutil.Address, fee int64, sequencesMap map[string]uint32) (*wire.MsgTx, int, error) {

	tx, idx, err := buildTxFromSacps(sacps)
	if err != nil {
		return nil, 0, err
	}

	// add inputs to the transaction
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

	totalSendAmount := int64(0)
	// add outputs to the transaction
	for _, r := range recipients {
		script, err := txscript.PayToAddrScript(r.To)
		if err != nil {
			return nil, 0, err
		}

		tx.AddTxOut(wire.NewTxOut(r.Amount, script))
		totalSendAmount += r.Amount
	}
	// add change output to the transaction if required
	if totalUTXOAmount >= totalSendAmount+fee {
		script, err := txscript.PayToAddrScript(changeAddr)
		if err != nil {
			return nil, 0, err
		}
		if totalUTXOAmount >= totalSendAmount+fee+DustAmount {
			tx.AddTxOut(wire.NewTxOut(totalUTXOAmount-totalSendAmount-fee, script))
		}
	} else if len(recipients) != 0 && len(utxos) != 0 {
		return nil, 0, ErrInsufficientFunds(totalUTXOAmount, totalSendAmount+fee)
	}

	// return the transaction
	return tx, idx, nil
}

func getUTXOsForSpendRequest(ctx context.Context, indexer IndexerClient, spendReq []SpendRequest) (UTXOs, utxoMap, int64, error) {
	utxos := UTXOs{}
	totalValue := int64(0)
	utxoMap := make(utxoMap)

	for _, req := range spendReq {
		utxosForAddress, err := indexer.GetUTXOs(ctx, req.ScriptAddress)
		if err != nil {
			return nil, nil, 0, err
		}

		utxos = append(utxos, utxosForAddress...)
		for _, utxo := range utxosForAddress {
			totalValue += utxo.Amount
		}
		utxoMap[req.ScriptAddress.EncodeAddress()] = utxosForAddress
	}

	// If there are any spend requests, check if the scripts have funds to spend
	if totalValue == 0 && len(spendReq) > 0 {
		return nil, nil, 0, ErrNoFundsToSpend
	}

	return utxos, utxoMap, totalValue, nil
}

func parseAddress(addr string) (btcutil.Address, error) {
	chainParams := chaincfg.MainNetParams
	address, err := btcutil.DecodeAddress(addr, &chainParams)
	if err == nil {
		return address, nil
	}
	address, err = btcutil.DecodeAddress(addr, &chaincfg.TestNet3Params)
	if err == nil {
		return address, nil
	}
	address, err = btcutil.DecodeAddress(addr, &chaincfg.RegressionNetParams)
	if err == nil {
		return address, nil
	}
	return nil, fmt.Errorf("error decoding address %s: %w", addr, err)
}

// Builds a prevOutFetcher for the given utxoMap, outpoints and txouts.
//
// outpoints should have corresponding txouts in the same order.
func buildPrevOutFetcher(utxosByAddressMap utxoMap, outpoints []wire.OutPoint, txouts []*wire.TxOut) (txscript.PrevOutputFetcher, error) {
	fetcher := NewPrevOutFetcherBuilder()
	for addrStr, utxos := range utxosByAddressMap {
		addr, err := parseAddress(addrStr)
		if err != nil {
			return nil, err
		}
		err = fetcher.AddFromUTXOs(utxos, addr)
		if err != nil {
			return nil, err
		}
	}

	err := fetcher.AddPrevOuts(outpoints, txouts)
	if err != nil {
		return nil, err
	}
	return fetcher.Build(), nil
}

// Returns the script to sign based on the script address type.
// For taproot address, the scriptPubKey is returned.
// For other addresses, the script is returned.
func getScriptToSign(scriptAddr btcutil.Address, script []byte) ([]byte, error) {
	// If the script address is a taproot address, we need to sign the scriptPubKey
	_, ok := scriptAddr.(*btcutil.AddressTaproot)
	if ok {
		script, err := txscript.PayToAddrScript(scriptAddr)
		if err != nil {
			return nil, err
		}
		return script, nil
	}
	return script, nil
}

// Signs the spend transaction
//
// Internally signTx is called for each input to sign the transaction.
func signSpendTx(ctx context.Context, tx *wire.MsgTx, startingIdx int, inputs []SpendRequest, utxoMap utxoMap, indexer IndexerClient, privateKey *secp256k1.PrivateKey) error {

	// building the prevOutFetcherBuilder
	// get the prevouts and txouts for the sacps to build the prevOutFetcher
	outpoints, txouts, err := getPrevoutsForSACPs(ctx, tx, startingIdx, indexer)
	if err != nil {
		return err
	}

	// This is an important step if we are spending from p2tr scripts which contain signature checks
	prevOutFetcher, err := buildPrevOutFetcher(utxoMap, outpoints, txouts)
	if err != nil {
		return err
	}

	idx := startingIdx

	// Loop through all the inputs, get utxos and sign the transaction
	for _, in := range inputs {

		// For tapscript spend, we need to sign the scriptPubKey
		// For every other spend, we need to sign the script (compiled script)
		script, err := getScriptToSign(in.ScriptAddress, in.Script)
		if err != nil {
			return err
		}

		utxos, ok := utxoMap[in.ScriptAddress.EncodeAddress()]
		if !ok {
			return ErrNoUTXOsFoundForAddress(in.ScriptAddress.String())
		}

		for _, utxo := range utxos {
			err = signTx(tx, prevOutFetcher, utxo.Amount, idx, in.Witness, script, &in.Leaf, in.HashType, privateKey)
			if err != nil {
				return err
			}
			idx++
		}
	}

	return nil
}

// Signs the transaction with the given witness and script.
// If there are OP Codes in the witness, they are replaced by the actual signature or pubkey.
func signTx(tx *wire.MsgTx, prevOutFetcher txscript.PrevOutputFetcher, amount int64, index int, witness [][]byte, script []byte, leaf *txscript.TapLeaf, hashType txscript.SigHashType, privateKey *secp256k1.PrivateKey) error {
	newWitness := [][]byte{}

	for _, w := range witness {
		switch string(w) {
		// Adds schnorr signature
		// Make sure to use this only with segwit v1 scripts
		case string(AddSignatureSchnorrOp):
			sigHashes := txscript.NewTxSigHashes(tx, prevOutFetcher)
			sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, index, amount, script, *leaf, hashType, privateKey)
			if err != nil {
				return err
			}
			newWitness = append(newWitness, sig)
		// Adds ecdsa signature
		// Make sure to use this only with segwit v0 scripts
		case string(AddSignatureSegwitOp):
			sigHashes := txscript.NewTxSigHashes(tx, prevOutFetcher)
			sig, err := txscript.RawTxInWitnessSignature(tx, sigHashes, index, amount, script, hashType, privateKey)
			if err != nil {
				return err
			}
			newWitness = append(newWitness, sig)
		// Adds 33 byte compressed pubkey
		case string(AddPubkeyCompressedOp):
			newWitness = append(newWitness, privateKey.PubKey().SerializeCompressed())
		// Adds 32 byte xonly pubkey
		// Make sure to use this only with segwit v1 scripts
		case string(AddXOnlyPubkeyOp):
			newWitness = append(newWitness, privateKey.PubKey().X().Bytes())
		// Adds the witness data as is
		default:
			newWitness = append(newWitness, w)
		}
	}
	tx.TxIn[index].Witness = newWitness
	return nil
}

func validateRequests(spendReqs []SpendRequest, sendReqs []SendRequest, sacps [][]byte) error {
	for _, in := range spendReqs {
		if len(in.Witness) == 0 {
			return fmt.Errorf("witness is required")
		}
		switch in.ScriptAddress.(type) {
		case *btcutil.AddressTaproot:
			// check if leaf exists
			if in.Leaf.LeafVersion == 0 || len(in.Leaf.Script) == 0 {
				return fmt.Errorf("spending leaf not found")
			}

			// check if control block is present
			_, err := txscript.ParseControlBlock(in.Witness[len(in.Witness)-1])
			if err != nil {
				return fmt.Errorf("invalid control block for tapscript spend")
			}

			// make sure witness contains proper signature and pubkey
			for _, w := range in.Witness {
				if string(w) == string(AddSignatureSegwitOp) {
					return fmt.Errorf("invalid signature op for tapscript spend, use schnorr signature op")
				}
				if string(w) == string(AddPubkeyCompressedOp) {
					return fmt.Errorf("invalid pubkey op for tapscript spend, use xonly pubkey op")
				}
			}
		case *btcutil.AddressWitnessPubKeyHash:

			script, err := txscript.PayToAddrScript(in.ScriptAddress)
			if err != nil {
				return err
			}
			if !bytes.Equal(in.Script, script) {
				return ErrInvalidP2wpkhScript
			}

			// len of witness should be 2
			if len(in.Witness) != 2 {
				return ErrInvalidP2wpkhWitness
			}
			// witness[0] should be signature
			if string(in.Witness[0]) != string(AddSignatureSegwitOp) {
				return ErrInvalidSigOp(string(AddSignatureSegwitOp))
			}
			// witness[1] should be pubkey
			if string(in.Witness[1]) != string(AddPubkeyCompressedOp) {
				return ErrInvalidPubkeyOp(string(AddPubkeyCompressedOp))
			}
		case *btcutil.AddressWitnessScriptHash:
			// the script should match the last element in the witness
			scriptInWitness := in.Witness[len(in.Witness)-1]
			if !bytes.Equal(scriptInWitness, in.Script) {
				return ErrScriptNotFoundInWitness
			}
			// should not contain Schnorr signature or xonly pubkey
			for _, w := range in.Witness {
				if string(w) == string(AddSignatureSchnorrOp) {
					return ErrInvalidSigOp(string(AddSignatureSegwitOp))
				}
				if string(w) == string(AddXOnlyPubkeyOp) {
					return ErrInvalidPubkeyOp(string(AddPubkeyCompressedOp))
				}
			}
		}

	}

	for _, r := range sendReqs {
		if r.Amount <= DustAmount {
			return ErrAmountLessThanDust
		}
	}

	for _, sacp := range sacps {
		_, err := buildAndValidateSacpTx(sacp)
		if err != nil {
			return err
		}
	}

	return nil
}

func calculateTotalSendAmount(req []SendRequest) int64 {
	totalSendAmount := int64(0)
	for _, r := range req {
		totalSendAmount += r.Amount
	}
	return totalSendAmount
}

func getNumberOfSigs(spends []SpendRequest) (int, int) {
	numSegwitSigs := 0
	numSchnorrSigs := 0
	for _, spend := range spends {
		for _, w := range spend.Witness {
			if string(w) == string(AddSignatureSegwitOp) {
				numSegwitSigs++
			} else if string(w) == string(AddSignatureSchnorrOp) {
				numSchnorrSigs++
			}
		}
	}
	return numSegwitSigs, numSchnorrSigs
}

func buildAndValidateSacpTx(sacp []byte) (*wire.MsgTx, error) {
	btcTx, err := btcutil.NewTxFromBytes(sacp)
	if err != nil {
		return nil, err
	}
	sacpTx := btcTx.MsgTx()
	err = validateSacp(sacpTx)
	if err != nil {
		return nil, err
	}
	return sacpTx, nil
}

func validateSacp(tx *wire.MsgTx) error {
	// TODO : simulate the tx and check if it is valid
	if len(tx.TxIn) == 0 {
		return ErrSCAPNoInputsFound
	}

	if len(tx.TxIn) != len(tx.TxOut) {
		return ErrSCAPInputsNotEqualOutputs
	}

	return nil
}
