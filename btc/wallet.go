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

type SpendScripts struct {
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

	// Hash type for the signature. Default is SigHashDefault
	HashType txscript.SigHashType
	IsSACP   bool
}

type SpendRequest struct {
	// Inputs to spend.
	// Each input contains the witness and script required to spend the script.
	Inputs    []SpendScripts
	ToAddress btcutil.Address
}

type SendRequest struct {
	Amount int64
	To     btcutil.Address
}

type Wallet interface {
	// Address returns the Pay-to-Witness-Public-Key-Hash (P2WPKH) address of the wallet.
	Address() btcutil.Address

	// Send initiates a transaction to send funds to one or more addresses specified in the req parameter.
	// It automatically calculates the necessary transaction fee based on the selected UTXOs.
	//
	// Example:
	// To send 1000 satoshis to an address, the SendRequest should have the following values:
	// 	txId, err := wallet.Send(context.Background(), []btc.SendRequest{
	// 		{
	// 			Amount: 1000,
	// 			To:     address,
	// 		},
	// 	})
	Send(ctx context.Context, req []SendRequest) (string, error)

	// Spend initiates a transaction to spend the Unspent Transaction Outputs (UTXOs)
	// of scripts specified in the req parameter.
	//
	// Example:
	// To spend from a P2WPKH address, the SpendRequest should have a single SpendScripts with the following values:
	// 	txId, err := wallet.Spend(context.Background(), &btc.SpendRequest{
	// 	Inputs: []btc.SpendScripts{
	// 		{
	// 			Witness: [][]byte{
	// 				btc.AddSignatureSegwitOp, // signature
	// 				btc.AddPubkeyCompressedOp, // pubkey
	// 			},
	// 			Script:        pkScript, // p2wpkh script
	// 			ScriptAddress: address, // p2wpkh address
	// 			HashType:      txscript.SigHashAll, // hash type, by default SigHashDefault is used (p2tr only)
	// 		},
	// 	},
	// 	ToAddress: address, // the recipient address
	// })
	Spend(ctx context.Context, req *SpendRequest) (string, error)

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

// Builds an unsigned transaction with the given utxos, recipients, change address and fee.
// If deductFeeFromSendAmount is true, the fee is deducted from the send amount.
func buildTransaction(utxos UTXOs, recipients []SendRequest, changeAddr btcutil.Address, fee int, deductFeeFromSendAmount bool) (*wire.MsgTx, error) {
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
		recipientFee := 0
		if deductFeeFromSendAmount {
			recipientFee = fee / len(recipients)
		}
		tx.AddTxOut(wire.NewTxOut(r.Amount-int64(recipientFee), script))
		totalSendAmount += r.Amount
	}

	// add change output to the transaction if required
	if !deductFeeFromSendAmount && totalUTXOAmount > totalSendAmount+int64(fee)+DustAmount {
		script, err := txscript.PayToAddrScript(changeAddr)
		if err != nil {
			return nil, err
		}
		tx.AddTxOut(wire.NewTxOut(totalUTXOAmount-totalSendAmount-int64(fee), script))
	}

	// return the transaction
	return tx, nil
}

// Send funds to the given address.
// The fee is calculated and adjusted based on the utxos available.
// depth is the number of times the function has been called recursively to find enough utxos to send the amount.
func (sw *SimpleWallet) send(ctx context.Context, recipients []SendRequest, fee int, depth int) (string, error) {

	// This means we made 100 recursive calls and still could not find enough utxos to send the amount
	if depth > 100 {
		return "", fmt.Errorf("could not find enough utxos to send amount")
	}

	totalSendAmount := calculateTotalSendAmount(recipients)
	// Get utxos for the given amount + fee
	// If there are not enough utxos, `GetUTXOsForAmount` will return an error
	utxos, _, err := sw.indexer.GetUTXOsForAmount(ctx, sw.signerAddr, totalSendAmount+int64(fee))
	if err != nil {
		return "", err
	}

	// build the transaction
	tx, err := buildTransaction(utxos, recipients, sw.signerAddr, fee, false)
	if err != nil {
		return "", err
	}

	// generate the signing script
	script, err := txscript.PayToAddrScript(sw.signerAddr)
	if err != nil {
		return "", err
	}

	// sign the transaction
	err = signSendTx(tx, utxos, script, sw.privateKey)
	if err != nil {
		return "", err
	}

	// estimate the fee required to make the transaction
	feeToBePaid, err := EstimateSegwitFee(tx, sw.feeEstimator, sw.feeLevel)
	if err != nil {
		return "", err
	}
	// if the fee required is more than the fee we are willing to pay, we try again
	if feeToBePaid > fee {
		return sw.send(ctx, recipients, feeToBePaid, depth+1)
	}

	// submit the transaction
	err = sw.indexer.SubmitTx(ctx, tx)
	if err != nil {
		return "", err
	}

	return tx.TxHash().String(), nil
}

// Sends funds to the given address
func (sw *SimpleWallet) Send(ctx context.Context, recipients []SendRequest) (string, error) {

	// validate the request
	for _, r := range recipients {
		if r.Amount < DustAmount {
			return "", fmt.Errorf("amount should be greater than %d", DustAmount)
		}
	}

	// starting fee is 1000 and it will be incremental gets adjusted based on the utxos selected
	// depth is 0 as we are starting the recursion
	return sw.send(ctx, recipients, 1000, 0)
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

// Spends the utxos from the given scripts.
// Fee is adjusted based on the utxos selected
func (sw *SimpleWallet) spend(ctx context.Context, req *SpendRequest, fee int) (string, error) {

	// get all utxos for the given scripts
	utxos, utxosByAddressMap, totalValue, err := getUTXOsForSpendRequest(ctx, sw.indexer, req)

	// validate spending amount
	if totalValue == 0 {
		return "", fmt.Errorf("no funds to spend")
	}
	if totalValue-int64(fee) < 0 {
		return "", fmt.Errorf("fee exceeds the total value")
	}

	// build the transaction
	tx, err := buildTransaction(utxos, []SendRequest{{
		Amount: totalValue,
		To:     req.ToAddress,
	}}, sw.signerAddr, fee, true)
	if err != nil {
		return "", err
	}

	// sign the transaction
	err = signSpendTx(tx, req.Inputs, utxosByAddressMap, sw.privateKey)
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
		return sw.spend(ctx, req, feeToBePaid)
	}

	// submit the transaction
	err = sw.indexer.SubmitTx(ctx, tx)
	if err != nil {
		return "", err
	}

	return tx.TxHash().String(), nil
}

// Spends the utxos from the given scripts
func (sw *SimpleWallet) Spend(ctx context.Context, req *SpendRequest) (string, error) {
	// input validation
	err := validateSpendRequest(req)
	if err != nil {
		return "", err
	}

	// starting fee is 1000 and it will be incremental gets adjusted based on the utxos selected
	return sw.spend(ctx, req, 1000)
}

// ------------------ Helper functions ------------------

func getUTXOsForSpendRequest(ctx context.Context, indexer IndexerClient, req *SpendRequest) (UTXOs, map[string]UTXOs, int64, error) {
	utxos := UTXOs{}
	totalValue := int64(0)
	utxosByAddress := make(map[string]UTXOs)
	for _, in := range req.Inputs {
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

// Signs the spend transaction
func signSpendTx(tx *wire.MsgTx, inputs []SpendScripts, utxosByAddressMap map[string]UTXOs, privateKey *secp256k1.PrivateKey) error {
	idx := 0
	for _, in := range inputs {

		// For tapscript spend, we need to sign the scriptPubKey
		// For every other spend, we need to sign the script
		script := in.Script
		var err error
		_, ok := in.ScriptAddress.(*btcutil.AddressTaproot)
		if ok {
			script, err = txscript.PayToAddrScript(in.ScriptAddress)
		}
		if err != nil {
			return err
		}

		utxos, ok := utxosByAddressMap[in.ScriptAddress.EncodeAddress()]
		if !ok {
			return fmt.Errorf("utxos not found for address %s", in.ScriptAddress.EncodeAddress())
		}

		for _, utxo := range utxos {
			err = signTx(tx, utxo.Amount, idx, in.Witness, script, &in.Leaf, in.HashType, privateKey)
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
func signTx(tx *wire.MsgTx, amount int64, index int, witness [][]byte, script []byte, leaf *txscript.TapLeaf, hashType txscript.SigHashType, privateKey *secp256k1.PrivateKey) error {
	newWitness := [][]byte{}

	for _, w := range witness {
		switch string(w) {
		// Adds schnorr signature
		// Make sure to use this only with segwit v1 scripts
		case string(AddSignatureSchnorrOp):
			fetcher := txscript.NewCannedPrevOutputFetcher(script, amount)
			sigHashes := txscript.NewTxSigHashes(tx, fetcher)

			sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, index, amount, script, *leaf, hashType, privateKey)
			if err != nil {
				return err
			}

			newWitness = append(newWitness, sig)
		// Adds normal signature
		// Make sure to use this only with segwit v0 scripts
		case string(AddSignatureSegwitOp):
			fetcher := txscript.NewCannedPrevOutputFetcher(script, amount)
			sigHashes := txscript.NewTxSigHashes(tx, fetcher)
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

// signs the send transaction
func signSendTx(tx *wire.MsgTx, utxos UTXOs, script []byte, privateKey *secp256k1.PrivateKey) error {
	// for p2wpkh, we only need to add the signature and pubkey
	witness := [][]byte{
		AddSignatureSegwitOp,
		AddPubkeyCompressedOp,
	}

	for i := range utxos {
		err := signTx(tx, utxos[i].Amount, i, witness, script, nil, txscript.SigHashAll, privateKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateSpendRequest(req *SpendRequest) error {
	for _, in := range req.Inputs {
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
