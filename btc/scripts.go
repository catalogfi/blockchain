package btc

import (
	"crypto/sha256"
	"fmt"
	"math"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

var ErrInvalidLockTime = fmt.Errorf("invalid lock-time")

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
	if waitTime > math.MaxUint16 || waitTime < 0 {
		return nil, ErrInvalidLockTime
	}
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

// isWaitTimeOpCode returns if the given opCode is a valid opCode for a `OP_CHECKSEQUENCEVERIFY` params.
// Since we require the timelock to be an integer (max_uint16), so we interpret maximum 5 bytes.
func isWaitTimeOpCode(opCode byte) bool {
	return (opCode >= txscript.OP_1 && opCode <= txscript.OP_16) ||
		(opCode >= txscript.OP_DATA_1 && opCode <= txscript.OP_DATA_3)
}

// IsHtlc returns if the given script is a HTLC script.
func IsHtlc(script []byte) bool {
	// 0xff is used to represent a data of variable length
	validHtlc := []byte{
		txscript.OP_IF,
		txscript.OP_SHA256,
		txscript.OP_DATA_32,
		txscript.OP_EQUALVERIFY,
		txscript.OP_DUP,
		txscript.OP_HASH160,
		txscript.OP_DATA_20,
		txscript.OP_ELSE,
		0xff,
		txscript.OP_CHECKSEQUENCEVERIFY,
		txscript.OP_DROP,
		txscript.OP_DUP,
		txscript.OP_HASH160,
		txscript.OP_DATA_20,
		txscript.OP_ENDIF,
		txscript.OP_EQUALVERIFY,
		txscript.OP_CHECKSIG,
	}
	tokenizer := txscript.MakeScriptTokenizer(0, script)

	for _, opCode := range validHtlc {
		if !tokenizer.Next() {
			return false
		}
		// Extra check for the lock time
		if opCode == 0xff {
			if !isWaitTimeOpCode(tokenizer.Opcode()) {
				return false
			}
			lockTime := decodeLocktime(tokenizer.Data())
			if lockTime > math.MaxUint16 || lockTime < 0 {
				return false
			}
			continue
		}
		if !(tokenizer.Opcode() == opCode) {
			return false
		}
	}
	return tokenizer.Done()
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

// modified from https://github.com/btcsuite/btcd/blob/4171854739fa2590a99c486341209d3aea8404dc/txscript/scriptnum.go#L198
func decodeLocktime(v []byte) int64 {
	// Zero is encoded as an empty byte slice.
	if len(v) == 0 {
		return 0
	}

	// Decode from little endian.
	var result int64
	for i, val := range v {
		result |= int64(val) << uint8(8*i)
	}

	// When the most significant byte of the input bytes has the sign bit
	// set, the result is negative.  So, remove the sign bit from the result
	// and make it negative.
	if v[len(v)-1]&0x80 != 0 {
		// The maximum length of v has already been determined to be 4
		// above, so uint8 is enough to cover the max possible shift
		// value of 24.
		result &= ^(int64(0x80) << uint8(8*(len(v)-1)))
		return -result
	}

	return result
}
