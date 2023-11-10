package btc

import (
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// MultisigScript generates a 2-out-2 multisig script.
func MultisigScript(pubKeyA, pubKeyB []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_2).
		AddData(pubKeyA).
		AddData(pubKeyB).
		AddOp(txscript.OP_2).
		AddOp(txscript.OP_CHECKMULTISIG).
		Script()
}

// HtlcScript generates a HTLC script following BIP-199.
// (https://github.com/bitcoin/bips/blob/master/bip-0199.mediawiki#summary)
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

func isNumOpCode(opcode byte) bool {
	return opcode >= 0x51 && opcode <= 0x60
}

func isPushDataOpCode(opcode byte) bool {
	return opcode >= 0x01 && opcode <= 0x4e
}

// IsHtlc returns if the given script is a HTLC script.
func IsHtlc(script []byte) bool {
	// 0xff is used to represent a data of variable length
	validHtlc := []byte{
		txscript.OP_IF,
		txscript.OP_SHA256,
		0x20, // to ensure secret hash is not more than 32 bytes
		txscript.OP_EQUALVERIFY,
		txscript.OP_DUP,
		txscript.OP_HASH160,
		0x14, // to ensure address is not more than 20 bytes
		txscript.OP_ELSE,
		0xff, // to accomodate variable encoding for wait time, a number
		txscript.OP_CHECKSEQUENCEVERIFY,
		txscript.OP_DROP,
		txscript.OP_DUP,
		txscript.OP_HASH160,
		0x14, // to ensure address is not more than 20 bytes
		txscript.OP_ENDIF,
		txscript.OP_EQUALVERIFY,
		txscript.OP_CHECKSIG,
	}
	// script version is 0
	tokenizer := txscript.MakeScriptTokenizer(0, script)

	if !tokenizer.Next() {
		return false
	}
	for i, opCode := range validHtlc {
		if !(tokenizer.Opcode() == opCode || opCode == 0xff) {
			return false
		}
		if opCode == 0xff {
			for tokenizer.Opcode() != validHtlc[i+1] {
				if !(isNumOpCode(tokenizer.Opcode()) || isPushDataOpCode(tokenizer.Opcode())) {
					return false
				}
				if !tokenizer.Next() {
					return false
				}
			}
			continue
		}
		if !tokenizer.Next() && i != len(validHtlc)-1 {
			return false
		}
	}
	return true
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

// MultisigWitness used for generating the witness script for spending a multisig utxo.
func MultisigWitness(script, sigA, sigB []byte) wire.TxWitness {
	witnessStack := wire.TxWitness(make([][]byte, 4))
	witnessStack[0] = nil
	witnessStack[1] = sigA
	witnessStack[2] = sigB
	witnessStack[3] = script
	return witnessStack
}

// HtlcWitness returns the witness for spending the htlc script, it has 2 possible paths either refund immediately
// through user secret or refund after timelock.
func HtlcWitness(script, pub, signature, secret []byte, redeem bool) wire.TxWitness {
	var witnessStack wire.TxWitness
	if redeem {
		witnessStack = make([][]byte, 5)
		witnessStack[0] = signature
		witnessStack[1] = pub
		witnessStack[2] = secret
		witnessStack[3] = []byte{0x1}
		witnessStack[4] = script
	} else {
		witnessStack = make([][]byte, 4)
		witnessStack[0] = signature
		witnessStack[1] = pub
		witnessStack[2] = nil
		witnessStack[3] = script
	}
	return witnessStack
}

// WitnessScriptHash returns the hash of the witness script.
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
