package btc

import (
	"crypto/sha256"

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
		AddInt64(waitTime).
		AddOp(txscript.OP_CHECKSEQUENCEVERIFY).
		AddOp(txscript.OP_DROP).
		AddData(ownerPub).
		AddOp(txscript.OP_ELSE).
		AddOp(txscript.OP_SHA256).
		AddData(refundSecretHash).
		AddOp(txscript.OP_EQUALVERIFY).
		AddData(revokerPub).
		AddOp(txscript.OP_ENDIF).
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

func WitnessScriptHash(witnessScript []byte) ([]byte, error) {
	bldr := txscript.NewScriptBuilder()

	bldr.AddOp(txscript.OP_0)
	scriptHash := sha256.Sum256(witnessScript)
	bldr.AddData(scriptHash[:])
	return bldr.Script()
}
