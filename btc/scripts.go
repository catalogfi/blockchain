package btc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
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
)

var ErrInvalidLockTime = fmt.Errorf("invalid lock-time")

// RedeemLeaf is one of the leaf scripts in the HTLC script which can be spent by revealing the secret
// by the redeemer.
//
// redeemerPubkey must be x-only pubkey of the redeemer.
func RedeemLeaf(redeemerPubkey, secretHash []byte) (txscript.TapLeaf, error) {
	script, err := txscript.NewScriptBuilder().
		AddOp(txscript.OP_SHA256).
		AddData(secretHash).
		AddOp(txscript.OP_EQUALVERIFY).
		AddData(redeemerPubkey).
		AddOp(txscript.OP_CHECKSIG).
		Script()
	if err != nil {
		return txscript.TapLeaf{}, err
	}
	return toLeaf(script), nil
}

// RefundLeaf is one of the leaf scripts in the HTLC script which can be spent by the initiator
// after the lock time.
//
// initiatorPubkey must be x-only pubkey of the initiator.
func RefundLeaf(initiatorPubkey []byte, lockTime uint32) (txscript.TapLeaf, error) {
	if lockTime > math.MaxUint16 {
		return txscript.TapLeaf{}, ErrInvalidLockTime
	}

	script, err := txscript.NewScriptBuilder().
		AddInt64(int64(lockTime)).
		AddOp(txscript.OP_CHECKSEQUENCEVERIFY).
		AddOp(txscript.OP_DROP).
		AddData(initiatorPubkey).
		AddOp(txscript.OP_CHECKSIG).
		Script()
	if err != nil {
		return txscript.TapLeaf{}, err
	}
	return toLeaf(script), nil
}

// MultiSigLeaf is a 2 on 2 multisig leaf script
//
// pubkeys must be x-only pubkeys of the initiator and the redeemer.
func MultiSigLeaf(initiatorPubkey, redeemerPubkey []byte) (txscript.TapLeaf, error) {

	script, err := txscript.NewScriptBuilder().
		AddData(initiatorPubkey).
		AddOp(txscript.OP_CHECKSIG).
		AddData(redeemerPubkey).
		AddOp(txscript.OP_CHECKSIGADD).
		AddOp(txscript.OP_2).
		AddOp(txscript.OP_NUMEQUAL).
		Script()
	if err != nil {
		return txscript.TapLeaf{}, err
	}
	return toLeaf(script), nil
}

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

func HtlcScriptV2(internalKey *btcec.PublicKey, chain *chaincfg.Params,  htlc *HTLC) (btcutil.Address, error) {

	leaves, err := htlcLeaves(htlc)
	if err != nil {
		return nil, err
	}

	tapScriptTree := txscript.AssembleTaprootScriptTree(leaves.ToArray()...)

	tapScriptRootHash := tapScriptTree.RootNode.TapHash()
	outputKey := txscript.ComputeTaprootOutputKey(
		internalKey, tapScriptRootHash[:],
	)

	addr, err := btcutil.NewAddressTaproot(outputKey.X().Bytes(), chain)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// isWaitTimeOpCode returns if the given opCode is a valid opCode for a `OP_CHECKSEQUENCEVERIFY` params.
// Since we require the timelock to be an integer (max_uint16), so we interpret maximum 3 bytes.
func isWaitTimeOpCode(opCode byte) bool {
	return (opCode >= txscript.OP_1 && opCode <= txscript.OP_16) ||
		(opCode >= txscript.OP_DATA_1 && opCode <= txscript.OP_DATA_3)
}

func IsRedeemLeaf(script []byte) (bool, string) {
	validRedeem := []byte{
		txscript.OP_SHA256,
		txscript.OP_DATA_32,
		txscript.OP_EQUALVERIFY,
		txscript.OP_DATA_32,
		txscript.OP_CHECKSIG,
	}
	tokenizer := txscript.MakeScriptTokenizer(0, script)

	var redeemerPubkey string
	isSecond := false

	for _, opCode := range validRedeem {
		if !tokenizer.Next() {
			return false, ""
		}
		if tokenizer.Opcode() != opCode {
			return false, ""
		}

		if opCode == txscript.OP_DATA_32 {
			if isSecond {
				redeemerPubkey = hex.EncodeToString(tokenizer.Data())
			}
			isSecond = !isSecond
		}
	}
	return tokenizer.Done(), redeemerPubkey
}

func IsRefundLeaf(script []byte) (bool, string) {
	validRefund := []byte{
		0xff,
		txscript.OP_CHECKSEQUENCEVERIFY,
		txscript.OP_DROP,
		txscript.OP_DATA_32,
		txscript.OP_CHECKSIG,
	}
	tokenizer := txscript.MakeScriptTokenizer(0, script)

	var refunderPubkey string

	for _, opCode := range validRefund {
		if !tokenizer.Next() {
			return false, ""
		}
		if opCode == 0xff {
			if !isWaitTimeOpCode(tokenizer.Opcode()) {
				return false, ""
			}
			lockTime := decodeLocktime(tokenizer.Data())
			if lockTime > math.MaxUint16 || lockTime < 0 {
				return false, ""
			}
			continue
		}

		if tokenizer.Opcode() != opCode {
			return false, ""
		}

		if opCode == txscript.OP_DATA_32 {
			refunderPubkey = hex.EncodeToString(tokenizer.Data())
		}
	}

	return tokenizer.Done(), refunderPubkey
}

func IsMultiSigLeaf(script []byte) (bool, string) {
	validMultiSig := []byte{
		txscript.OP_DATA_32,
		txscript.OP_CHECKSIG,
		txscript.OP_DATA_32,
		txscript.OP_CHECKSIGADD,
		txscript.OP_2,
		txscript.OP_NUMEQUAL,
	}
	tokenizer := txscript.MakeScriptTokenizer(0, script)

	var refunderPubkey string
	isFirst := true

	for _, opCode := range validMultiSig {
		if !tokenizer.Next() {
			return false, ""
		}
		if tokenizer.Opcode() != opCode {
			return false, ""
		}

		if opCode == txscript.OP_DATA_32 {
			if isFirst {
				refunderPubkey = hex.EncodeToString(tokenizer.Data())
			}
			isFirst = !isFirst
		}
	}
	return tokenizer.Done(), refunderPubkey
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

func GardenNUMS() (*btcec.PublicKey, error) {
	// H value from BIP-341
	xCoordHex := "0x50929b74c1a04954b78b4b6035e97a5e078a5a0f28ec96d547bfee9ace803ac0"

	// Convert the x coordinate from hex to a big integer
	xCoord, _ := new(big.Int).SetString(xCoordHex[2:], 16)

	// Calculate the y coordinate for the given x coordinate
	curve := btcec.S256()
	yCoord := new(big.Int).ModSqrt(new(big.Int).Exp(xCoord, big.NewInt(3), curve.P), curve.P)
	format := byte(0x03)
	if yCoord.Bit(0) == 0 {
		format = byte(0x02)
	}

	compressedH := append([]byte{format}, xCoord.Bytes()...)
	H, err := btcec.ParsePubKey(compressedH[:])
	if err != nil {
		return nil, err
	}

	r := sha256.Sum256([]byte("GardenHTLC"))
	_, rG := btcec.PrivKeyFromBytes(r[:])

	numsX, numsY := curve.Add(H.X(), H.Y(), rG.X(), rG.Y())
	if !curve.IsOnCurve(numsX, numsY) {
		return nil, fmt.Errorf("invalid NUMS point")
	}

	formatNums := byte(0x03)
	if numsY.Bit(0) == 0 {
		formatNums = byte(0x02)
	}

	compressedNums := append([]byte{formatNums}, numsX.Bytes()...)
	nums, err := btcec.ParsePubKey(compressedNums)
	if err != nil {
		panic(err)
	}
	return nums, nil
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

// Helper function to convert a script to a leaf with version 0xc0.
func toLeaf(script []byte) txscript.TapLeaf {
	return txscript.NewTapLeaf(0xc0, script)
}
