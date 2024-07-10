package btc

import (
	"bytes"
	"context"
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
	// Adds a 64 or 65 byte schnorr signature
	AddSignatureSchnorrOp = []byte("add_signature_segwit_v1")
	// Adds a 72 byte ECDSA signature
	AddSignatureSegwitOp = []byte("add_signature_segwit_v0")
	// Adds a 33 byte compressed pubkey
	AddPubkeyCompressedOp = []byte("add_pubkey_compressed")
	// Adds a 32 byte xonly pubkey
	AddXOnlyPubkeyOp = []byte("add_xonly_pubkey")
)

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

	// Recipient address
	To btcutil.Address
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
	Send(ctx context.Context, sendReq []SendRequest, spendReq []SpendRequest) (string, error)

	// Status checks the status of a transaction using its transaction ID (txid).
	// Returns the transaction and a boolean indicating whether the transaction is submitted or not and an error
	Status(ctx context.Context, id string) (Transaction, bool, error)
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

func (sw *SimpleWallet) Send(ctx context.Context, sendRequests []SendRequest, spendRequests []SpendRequest) (string, error) {
	//TODO: Estimate a fee based on requests
	fee := 1000
	// validate the spend requests
	err := validateSpendRequest(spendRequests)
	if err != nil {
		return "", err
	}

	//validate the send requests
	for _, r := range sendRequests {
		if r.Amount <= DustAmount {
			return "", fmt.Errorf("amount is less than dust")
		}
	}

	return sw.spendAndSend(ctx, sendRequests, spendRequests, fee, 0)
}

func (sw *SimpleWallet) spendAndSend(ctx context.Context, sendRequests []SendRequest, spendRequests []SpendRequest, fee int, depth int) (string, error) {

	// This means we made 100 recursive calls and still could not find enough utxos to send the amount
	if depth > 100 {
		return "", fmt.Errorf("could not find enough utxos to satisfy the requests")
	}

	// These variables will be empty if there are no spend requests
	spendUTXOs, spendUTXOsMap, balanceOfScripts, err := getUTXOsForSpendRequest(ctx, sw.indexer, spendRequests)
	if err != nil {
		return "", err
	}

	// If there are any spend requests, check if the scripts have funds to spend
	if balanceOfScripts == 0 && len(spendRequests) > 0 {
		return "", fmt.Errorf("scripts have no funds to spend")
	}

	// coverUTXOs are the UTXOs used to cover the remaining amount required to send
	var coverUTXOs UTXOs

	totalSendAmount := calculateTotalSendAmount(sendRequests)
	if balanceOfScripts <= totalSendAmount {
		utxos, _, err := sw.indexer.GetUTXOsForAmount(ctx, sw.signerAddr, totalSendAmount-balanceOfScripts+int64(fee))
		if err != nil {
			return "", err
		}
		coverUTXOs = append(coverUTXOs, utxos...)
	}

	// build the transaction
	tx, err := buildTransaction(append(spendUTXOs, coverUTXOs...), sendRequests, sw.signerAddr, fee)
	if err != nil {
		return "", err
	}

	// Sign the spend inputs
	err = signSpendTx(tx, spendRequests, spendUTXOsMap, sw.privateKey)

	// get the send signing script
	script, err := txscript.PayToAddrScript(sw.signerAddr)
	if err != nil {
		return "", err
	}

	// Sign the cover inputs
	// This is a no op if there are no cover utxos
	err = signSendTx(tx, coverUTXOs, len(spendUTXOs), script, sw.privateKey)
	if err != nil {
		return "", err
	}

	// estimate the fee required to make the transaction
	feeToBePaid, err := EstimateSegwitFee(tx, sw.feeEstimator, sw.feeLevel)
	if err != nil {
		return "", err
	}

	// if the fee required is more than the fee we are willing to pay, we try again with the new fee
	if feeToBePaid > fee {
		return sw.spendAndSend(ctx, sendRequests, spendRequests, feeToBePaid, depth+1)
	}

	// submit the transaction
	return tx.TxHash().String(), sw.indexer.SubmitTx(ctx, tx)
}

// Returns the status of submitted transaction either via `send` or `spend`
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

// ------------------ Helper functions ------------------

// Builds an unsigned transaction with the given utxos, recipients, change address and fee.
// If deductFeeFromSendAmount is true, the fee is deducted from the send amount.
func buildTransaction(utxos UTXOs, recipients []SendRequest, changeAddr btcutil.Address, fee int) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(DefaultTxVersion)

	// add inputs to the transaction
	totalUTXOAmount := int64(0)
	for _, utxo := range utxos {
		txid, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}
		vout := utxo.Vout
		txIn := wire.NewTxIn(wire.NewOutPoint(txid, vout), nil, nil)
		tx.AddTxIn(txIn)
		totalUTXOAmount += utxo.Amount
	}

	totalSendAmount := int64(0)
	// add outputs to the transaction
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

	// add change output to the transaction if required
	if totalUTXOAmount >= totalSendAmount+int64(fee) {
		script, err := txscript.PayToAddrScript(changeAddr)
		if err != nil {
			return nil, err
		}
		if totalUTXOAmount >= totalSendAmount+int64(fee)+DustAmount {
			tx.AddTxOut(wire.NewTxOut(totalUTXOAmount-totalSendAmount-int64(fee), script))
		}
	} else {
		return nil, fmt.Errorf("insufficient funds")
	}

	// return the transaction
	return tx, nil
}

func getUTXOsForSpendRequest(ctx context.Context, indexer IndexerClient, req []SpendRequest) (UTXOs, map[string]UTXOs, int64, error) {
	utxos := UTXOs{}
	totalValue := int64(0)
	utxosByAddress := make(map[string]UTXOs)

	for _, in := range req {
		utxosForAddress, err := indexer.GetUTXOs(ctx, in.ScriptAddress)
		if err != nil {
			return nil, nil, 0, err
		}
		utxos = append(utxos, utxosForAddress...)
		for _, utxo := range utxosForAddress {
			totalValue += utxo.Amount
		}
		utxosByAddress[in.ScriptAddress.EncodeAddress()] = utxosForAddress
	}
	return utxos, utxosByAddress, totalValue, nil
}

// Useful for tapscript spend where we need to sign all inputs and their respective amounts
func buildPrevOutFetcher(utxosByAddressMap map[string]UTXOs, req []SpendRequest) (txscript.PrevOutputFetcher, error) {
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)
	for _, in := range req {
		utxos, ok := utxosByAddressMap[in.ScriptAddress.EncodeAddress()]
		if !ok {
			return nil, fmt.Errorf("utxos not found for address %s", in.ScriptAddress.EncodeAddress())
		}
		for _, utxo := range utxos {
			script, err := txscript.PayToAddrScript(in.ScriptAddress)
			if err != nil {
				continue
			}
			txId, err := chainhash.NewHashFromStr(utxo.TxID)
			if err != nil {
				continue
			}
			prevOutFetcher.AddPrevOut(*wire.NewOutPoint(txId, utxo.Vout), wire.NewTxOut(utxo.Amount, script))
		}
	}
	return prevOutFetcher, nil
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
func signSpendTx(tx *wire.MsgTx, inputs []SpendRequest, utxosByAddressMap map[string]UTXOs, privateKey *secp256k1.PrivateKey) error {
	idx := 0
	// build the prevOutFetcher
	// This is an important step if we are spending from p2tr scripts which contain signature checks
	prevOutFetcher, err := buildPrevOutFetcher(utxosByAddressMap, inputs)
	if err != nil {
		return err
	}
	// Loop through all the inputs, get utxos and sign the transaction
	for _, in := range inputs {

		// For tapscript spend, we need to sign the scriptPubKey
		// For every other spend, we need to sign the script (compiled script)
		script, err := getScriptToSign(in.ScriptAddress, in.Script)
		if err != nil {
			return err
		}

		utxos, ok := utxosByAddressMap[in.ScriptAddress.EncodeAddress()]
		if !ok {
			return fmt.Errorf("utxos not found for address %s", in.ScriptAddress.EncodeAddress())
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

// Signs the send transaction (p2wpkh spend).
// Use startingIdx to start signing from a specific index
func signSendTx(tx *wire.MsgTx, utxos UTXOs, startingIdx int, script []byte, privateKey *secp256k1.PrivateKey) error {
	// for p2wpkh, we only need to add the signature and pubkey
	witness := [][]byte{
		AddSignatureSegwitOp,
		AddPubkeyCompressedOp,
	}
	idx := startingIdx
	for i := range utxos {
		fetcher := txscript.NewCannedPrevOutputFetcher(script, utxos[i].Amount)
		err := signTx(tx, fetcher, utxos[i].Amount, idx, witness, script, nil, txscript.SigHashAll, privateKey)
		if err != nil {
			return err
		}
		idx++
	}
	return nil
}

func validateSpendRequest(req []SpendRequest) error {
	for _, in := range req {
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
				return fmt.Errorf("invalid script for p2wpkh spend")
			}

			// len of witness should be 2
			if len(in.Witness) != 2 {
				return fmt.Errorf("invalid witness for p2wpkh spend")
			}
			// witness[0] should be signature
			if string(in.Witness[0]) != string(AddSignatureSegwitOp) {
				return fmt.Errorf("invalid signature op for p2wpkh spend")
			}
			// witness[1] should be pubkey
			if string(in.Witness[1]) != string(AddPubkeyCompressedOp) {
				return fmt.Errorf("invalid pubkey op for p2wpkh spend")
			}
		case *btcutil.AddressWitnessScriptHash:
			// the script should match the last element in the witness
			scriptInWitness := in.Witness[len(in.Witness)-1]
			if !bytes.Equal(scriptInWitness, in.Script) {
				return fmt.Errorf("no script found in witness")
			}
			// should not contain Schnorr signature or xonly pubkey
			for _, w := range in.Witness {
				if string(w) == string(AddSignatureSchnorrOp) {
					return fmt.Errorf("invalid signature op for p2wsh spend, use segwit signature op")
				}
				if string(w) == string(AddXOnlyPubkeyOp) {
					return fmt.Errorf("invalid pubkey op for p2wsh spend, use compressed pubkey op")
				}
			}
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
