package btc

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// MultisigScript generates a 2-out-2 multisig script for instant wallet.
func MultisigScript(pubKeyA, pubKeyB []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_2).
		AddData(pubKeyA).
		AddData(pubKeyB).
		AddOp(txscript.OP_2).
		AddOp(txscript.OP_CHECKMULTISIG).
		Script()
}

// HtlcScript generates a HTLC script for refunding purpose.
func HtlcScript(ownerPub, revokerPub, refundSecretHash []byte, waitTime int64) ([]byte, error) {
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_IF).
		AddOp(txscript.OP_SHA256).
		AddData(refundSecretHash).
		AddOp(txscript.OP_EQUALVERIFY).
		AddOp(txscript.OP_DUP).
		AddOp(txscript.OP_HASH160).
		AddData(revokerPub).
		AddOp(txscript.OP_ELSE).
		AddInt64(waitTime).
		AddOp(txscript.OP_CHECKSEQUENCEVERIFY).
		AddOp(txscript.OP_DROP).
		AddOp(txscript.OP_DUP).
		AddOp(txscript.OP_HASH160).
		AddData(ownerPub).
		AddOp(txscript.OP_ENDIF).
		AddOp(txscript.OP_EQUALVERIFY).
		AddOp(txscript.OP_CHECKSIG).
		Script()
}

// NewRefundTx is helper function to build the refund tx. It assumes the input has only one utxo which is the given
// funding utxo, and the output will be the refundScript.
func NewRefundTx(fundingUtxo UTXO, refundScript []byte) (*wire.MsgTx, error) {
	refundTx := wire.NewMsgTx(DefaultBTCVersion)

	// Input
	txID, err := chainhash.NewHashFromStr(fundingUtxo.TxID)
	if err != nil {
		return nil, err
	}
	refundTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(txID, fundingUtxo.Vout), nil, nil))

	// Output
	payScript, err := WitnessScriptHash(refundScript)
	if err != nil {
		return nil, err
	}
	refundTx.AddTxOut(wire.NewTxOut(fundingUtxo.Amount, payScript))
	return refundTx, nil
}

func SignMultisig(tx *wire.MsgTx, index int, amount int64, key1, key2 *btcec.PrivateKey, hashType txscript.SigHashType) error {
	multiSigScript, err := MultisigScript(key1.PubKey().SerializeCompressed(), key2.PubKey().SerializeCompressed())
	if err != nil {
		return err
	}

	fetcher := txscript.NewCannedPrevOutputFetcher(multiSigScript, amount)
	sig1, err := txscript.RawTxInWitnessSignature(tx, txscript.NewTxSigHashes(tx, fetcher), index, amount, multiSigScript, hashType, key1)
	if err != nil {
		return err
	}
	sig2, err := txscript.RawTxInWitnessSignature(tx, txscript.NewTxSigHashes(tx, fetcher), index, amount, multiSigScript, hashType, key2)
	if err != nil {
		return err
	}
	witness := wire.TxWitness(make([][]byte, 4))
	witness[0] = nil
	witness[1] = sig1
	witness[2] = sig2
	witness[3] = multiSigScript

	tx.TxIn[index].Witness = witness
	return nil
}

func SignHtlcScript(tx *wire.MsgTx, ownerPub, redeemerPub, secret, secretHash []byte, index int, amount, waitTime int64, key *btcec.PrivateKey) error {
	htlcScript, err := HtlcScript(btcutil.Hash160(ownerPub), btcutil.Hash160(redeemerPub), secretHash[:], waitTime)
	if err != nil {
		return err
	}

	// Sign the message
	fetcher := txscript.NewCannedPrevOutputFetcher(htlcScript, amount)
	sig, err := txscript.RawTxInWitnessSignature(tx, txscript.NewTxSigHashes(tx, fetcher), index, amount, htlcScript, txscript.SigHashAll, key)
	if err != nil {
		return err
	}

	// Set the witness
	if bytes.Equal(key.PubKey().SerializeCompressed(), ownerPub) {
		witness := make([][]byte, 4)
		witness[0] = sig
		witness[1] = ownerPub
		witness[2] = nil
		witness[3] = htlcScript
		tx.TxIn[index].Witness = witness
	} else {
		witness := make([][]byte, 5)
		witness[0] = sig
		witness[1] = redeemerPub
		witness[2] = secret
		witness[3] = []byte{0x1}
		witness[4] = htlcScript
		tx.TxIn[index].Witness = witness
	}
	return nil
}

func WitnessScriptHash(witnessScript []byte) ([]byte, error) {
	builder := txscript.NewScriptBuilder()

	builder.AddOp(txscript.OP_0)
	scriptHash := sha256.Sum256(witnessScript)
	builder.AddData(scriptHash[:])
	return builder.Script()
}

// ParseScriptSigP2PKH gets the signature and public key from a P2PKH scriptSig.
func ParseScriptSigP2PKH(scriptSig []byte) ([]byte, []byte, error) {
	// The scriptSig should be (length + （sig 70/71/72 + sighash） + length + pubkey)
	switch len(scriptSig) {
	case 1 + 71 + 1 + 33:
		return scriptSig[1:72], scriptSig[73:], nil
	case 1 + 72 + 1 + 33:
		return scriptSig[1:73], scriptSig[74:], nil
	case 1 + 73 + 1 + 33:
		return scriptSig[1:74], scriptSig[75:], nil
	default:
		return nil, nil, fmt.Errorf("invalid script length %v", len(scriptSig))
	}
}
