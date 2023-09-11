package btc

import (
	"crypto/sha256"
	"fmt"

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
func SetRefundSpendWitness(refundSpend *wire.MsgTx, ownerPubkey, revokerPubkey, signature, refundSecretHash, refundSecret []byte, owner bool, waitTime int64) error {
	refundScript, err := HtlcScript(btcutil.Hash160(ownerPubkey), btcutil.Hash160(revokerPubkey), refundSecretHash, waitTime)
	if err != nil {
		return err
	}

	var witnessStack wire.TxWitness

	if owner {
		// owner
		witnessStack = make([][]byte, 4)
		witnessStack[0] = signature
		witnessStack[1] = ownerPubkey
		witnessStack[2] = nil
		witnessStack[3] = refundScript
	} else {
		// revoker
		witnessStack = make([][]byte, 5)
		witnessStack[0] = signature
		witnessStack[1] = revokerPubkey
		witnessStack[2] = refundSecret
		witnessStack[3] = []byte{0x1}
		witnessStack[4] = refundScript
	}
	refundSpend.TxIn[0].Witness = witnessStack
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
