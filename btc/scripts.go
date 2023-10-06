package btc

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
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

// IsHtlc returns if the given script is a HTLC script.
func IsHtlc(script []byte) bool {
	if len(script) != 89 {
		return false
	}

	return script[0] == txscript.OP_IF &&
		script[1] == txscript.OP_SHA256 &&
		// Skip script[2:34] which is the secret hash
		script[35] == txscript.OP_EQUALVERIFY &&
		script[36] == txscript.OP_DUP &&
		script[37] == txscript.OP_HASH160 &&
		// script[38:59] is the public key hash
		script[38] == 0x14 &&
		script[59] == txscript.OP_ELSE &&
		// Skip script[60:61] which is the wait time
		script[61] == txscript.OP_CHECKSEQUENCEVERIFY &&
		script[62] == txscript.OP_DROP &&
		script[63] == txscript.OP_DUP &&
		script[64] == txscript.OP_HASH160 &&
		// script[65:86] is the public key hash
		script[65] == 0x14 &&
		script[86] == txscript.OP_ENDIF &&
		script[87] == txscript.OP_EQUALVERIFY &&
		script[88] == txscript.OP_CHECKSIG
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

// SetMultisigWitness used for setting the witness script for any txs using the instant wallet utxo(redeem, refund)
func SetMultisigWitness(multisigTx *wire.MsgTx, pubkeyA, sigA, pubKeyB, sigB []byte) error {
	// assumed that there's only 1 txIn (instant wallet utxo)
	instantWalletScript, err := MultisigScript(pubkeyA, pubKeyB)
	if err != nil {
		return err
	}
	witnessStack := wire.TxWitness(make([][]byte, 4))
	witnessStack[0] = nil
	witnessStack[1] = sigA
	witnessStack[2] = sigB
	witnessStack[3] = instantWalletScript

	multisigTx.TxIn[0].Witness = witnessStack
	return nil
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

// SetRefundSpendWitness used for setting witness signature for spending the refunded utxo from the refund script,
// has 2 possible paths either refund immediately through user secret or refund through timelock.
func SetRefundSpendWitness(refundSpend *wire.MsgTx, ownerPubkey, revokerPubkey, signature, refundSecretHash, refundSecret, op []byte, waitTime int64) error {
	// assumed that there's only 1 txIn (instant wallet utxo)
	refundScript, err := HtlcScript(ownerPubkey, revokerPubkey, refundSecretHash, waitTime)
	if err != nil {
		return err
	}

	var witnessStack wire.TxWitness

	// else condition for instant revokation
	if op == nil {
		witnessStack = make([][]byte, 4)
		witnessStack[0] = signature
		witnessStack[1] = refundSecret
		witnessStack[2] = nil
		witnessStack[3] = refundScript
	} else {
		// for users normal refund flow
		witnessStack = make([][]byte, 3)
		witnessStack[0] = signature
		witnessStack[1] = []byte{0x1}
		witnessStack[2] = refundScript
	}
	refundSpend.TxIn[0].Witness = witnessStack
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

// SpendingWitness loops through the txs of a particular address and return the spending witness of the address.
// It will return a nil string slice if no spending tx is found.
func SpendingWitness(address btcutil.Address, txs []Transaction) []string {
	for _, tx := range txs {
		for _, vin := range tx.VINs {
			if vin.Prevout.ScriptPubKeyAddress == address.EncodeAddress() && vin.Witness != nil {
				return *vin.Witness
			}
		}
	}
	return nil
}

// P2wshAddress returns the P2WSH address of the give script
func P2wshAddress(script []byte, network *chaincfg.Params) (btcutil.Address, error) {
	scriptHash := sha256.Sum256(script)
	return btcutil.NewAddressWitnessScriptHash(scriptHash[:], network)
}
