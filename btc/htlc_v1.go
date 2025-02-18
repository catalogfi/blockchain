package btc

import (
	"crypto/sha256"
	"math"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

var (
	// MaxmiumHtlcSize is the maximum size of an HTLC script. It assumes each public key hash is 20 bytes and the secret
	// hash is 32 bytes. Details calculation is
	//      OpCode * 13 + PublicKeyHash *2 + SecretHash + timelock(maximum 65535)
	MaxmiumHtlcSize = 13 + (20+1)*2 + (32 + 1) + (3 + 1)

	// NormalHtlcSize is a more accurate way to estimate the htlc script size, since the timelock will be between
	// [128, 32767] most of our use cases.
	NormalHtlcSize = 13 + (20+1)*2 + (32 + 1) + (2 + 1)

	// RedeemHtlcRefundSigScriptSize is an estimate of the sigScript size when refunding an htlc script
	// stack number + stack size * 4 + signature + public key + script size
	RedeemHtlcRefundSigScriptSize = 1 + 4 + 73 + 33 + NormalHtlcSize

	// RedeemHtlcRedeemSigScriptSize is an estimate of the sigScript size when redeeming an htlc script
	// stack number + stack size * 5 + signature + public key + secret + script size
	RedeemHtlcRedeemSigScriptSize = func(secretSize int) int {
		return 1 + 5 + 73 + 33 + secretSize + +1 + NormalHtlcSize
	}

	// RedeemMultisigSigScriptSize is an estimate of the sigScript size from an 2-of-2 multisig script
	// stack number + stack size * 4 + signature * 2 + script size
	RedeemMultisigSigScriptSize = 1 + 4 + 73*2 + 71
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
		if tokenizer.Opcode() != opCode {
			return false
		}
	}
	return tokenizer.Done()
}

// HtlcWitness returns the witness for spending the htlc script, it has 2 possible paths either refund immediately
// through user secret or refund after timelock.
func HtlcWitness(script, pub, signature, secret []byte) wire.TxWitness {
	var witnessStack wire.TxWitness
	if len(secret) != 0 {
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

// MultisigWitness used for generating the witness script for spending a multisig utxo.
func MultisigWitness(script, sigA, sigB []byte) wire.TxWitness {
	witnessStack := wire.TxWitness(make([][]byte, 4))
	witnessStack[0] = nil
	witnessStack[1] = sigA
	witnessStack[2] = sigB
	witnessStack[3] = script
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

// P2wshAddress returns the P2WSH address of the give script
func P2wshAddress(script []byte, network *chaincfg.Params) (btcutil.Address, error) {
	scriptHash := sha256.Sum256(script)
	return btcutil.NewAddressWitnessScriptHash(scriptHash[:], network)
}
