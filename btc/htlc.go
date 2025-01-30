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
	"github.com/btcsuite/btcwallet/waddrmgr"
)

var ErrInvalidLockTime = fmt.Errorf("invalid lock-time")

var GardenNums *btcec.PublicKey

func init() {
	var err error
	GardenNums, err = GenerateGardenNUMS()
	if err != nil {
		panic(fmt.Errorf("generate Garden NUMS failed: %v", err))
	}
}

// GenerateGardenNUMS generates a new Nothing Up My Sleeve (NUMS) point as the taproot internal key.
// We take ths point mentioned in BIP-341 and modified it with a fixed value to get a new NUMS point which isn't on the
// curve, this is to make sure the HTLC can only be redeemed by script path.
func GenerateGardenNUMS() (*btcec.PublicKey, error) {
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
	return btcec.ParsePubKey(compressedNums)
}

var (
	// MaxInitiationUtxoNumber is the maximum number of utxos that can be used to initiate an HTLC.
	MaxInitiationUtxoNumber = 1

	// todo : calculate these
	BaseSizeHtlcRedeem        = 0
	BaseSizeHtlcRefund        = 0
	BaseSizeHtlcInstantRefund = 0

	// todo : calculate these
	SegwitSizeHtlcRedeem        = 10
	SegwitSizeHtlcRefund        = 0
	SegwitSizeHtlcInstantRefund = 0
)

type HTLC struct {
	InitiatorPubKey []byte
	RedeemerPubKey  []byte
	SecretHash      []byte
	Timelock        int64
	Amount          int64

	tree *txscript.IndexedTapScriptTree
}

func NewHTLC(initiatorPubKey, redeemerPubKey, secretHash []byte, timelock, amount int64) (*HTLC, error) {
	// todo : validate the input for generate these leafs
	redeemLeaf, err := RedeemLeaf(redeemerPubKey, secretHash)
	if err != nil {
		return nil, err
	}
	refundLeaf, err := RefundLeaf(initiatorPubKey, timelock)
	if err != nil {
		return nil, err
	}
	instantRefundLeaf, err := InstantRefundLeaf(initiatorPubKey, redeemerPubKey)
	if err != nil {
		return nil, err
	}
	tree := txscript.AssembleTaprootScriptTree(instantRefundLeaf, refundLeaf, redeemLeaf)

	return &HTLC{
		InitiatorPubKey: initiatorPubKey,
		RedeemerPubKey:  redeemerPubKey,
		SecretHash:      secretHash,
		Timelock:        timelock,
		Amount:          amount,

		tree: tree,
	}, nil
}

func (htlc *HTLC) Address(network *chaincfg.Params) (btcutil.Address, error) {
	if htlc.tree == nil {
		return nil, fmt.Errorf("empty htlc tree")
	}
	rootHash := htlc.tree.RootNode.TapHash()
	outputKey := txscript.ComputeTaprootOutputKey(GardenNums, rootHash[:])
	return PublicKeyAddress(network, waddrmgr.TaprootPubKey, outputKey)
}

func (htlc *HTLC) P2trScript() ([]byte, error) {
	rootHash := htlc.tree.RootNode.TapHash()
	outputKey := txscript.ComputeTaprootOutputKey(GardenNums, rootHash[:])
	return txscript.PayToTaprootScript(outputKey)
}

func (htlc *HTLC) Initiated(utxos []UTXO) (bool, uint64, error) {
	if len(utxos) > MaxInitiationUtxoNumber {
		return false, 0, fmt.Errorf("too many utxos for initiation")
	}

	total, blockHeight := int64(0), uint64(0)
	for _, utxo := range utxos {
		if utxo.Status != nil && utxo.Status.Confirmed {
			total += utxo.Amount
			if *utxo.Status.BlockHeight > blockHeight {
				blockHeight = *utxo.Status.BlockHeight
			}
		}
	}
	return total >= htlc.Amount, blockHeight, nil
}

func (htlc *HTLC) Redeemable(utxos []UTXO) (bool, uint64, error) {
	if len(utxos) > MaxInitiationUtxoNumber {
		return false, 0, fmt.Errorf("too many utxos for initiation")
	}

	total, blockHeight := int64(0), uint64(0)
	for _, utxo := range utxos {
		if utxo.Status != nil && utxo.Status.Confirmed {
			total += utxo.Amount
			if *utxo.Status.BlockHeight > blockHeight {
				blockHeight = *utxo.Status.BlockHeight
			}
		}
	}
	return total >= htlc.Amount, blockHeight, nil
}

func (htlc *HTLC) Refundable(utxos []UTXO, latest uint64) bool {
	for _, utxo := range utxos {
		if utxo.Status != nil && utxo.Status.Confirmed {
			if latest-*utxo.Status.BlockHeight >= uint64(htlc.Timelock) {
				return false
			}
		}
	}
	return true
}

func (htlc *HTLC) Expired(utxos []UTXO, latest uint64) bool {
	for _, utxo := range utxos {
		if utxo.Status != nil && utxo.Status.Confirmed {
			if latest-*utxo.Status.BlockHeight >= uint64(htlc.Timelock) {
				return true
			}
		}
	}
	return false
}

func (htlc *HTLC) RedeemLeaf() (txscript.TapLeaf, txscript.ControlBlock) {
	leaf := htlc.tree.RootNode.Right().(txscript.TapLeaf)
	index := htlc.tree.LeafProofIndex[leaf.TapHash()]
	ctrBlk := htlc.tree.LeafMerkleProofs[index].ToControlBlock(GardenNums)
	return leaf, ctrBlk
}

func (htlc *HTLC) RefundLeaf() (txscript.TapLeaf, txscript.ControlBlock) {
	leaf := htlc.tree.RootNode.Left().Right().(txscript.TapLeaf)
	index := htlc.tree.LeafProofIndex[leaf.TapHash()]
	ctrBlk := htlc.tree.LeafMerkleProofs[index].ToControlBlock(GardenNums)
	return leaf, ctrBlk
}

func (htlc *HTLC) InstantRefundLeaf() (txscript.TapLeaf, txscript.ControlBlock) {
	leaf := htlc.tree.RootNode.Left().Left().(txscript.TapLeaf)
	index := htlc.tree.LeafProofIndex[leaf.TapHash()]
	ctrBlk := htlc.tree.LeafMerkleProofs[index].ToControlBlock(GardenNums)
	return leaf, ctrBlk
}

// RedeemLeaf is one of the leaf scripts in the HTLC script which can be spent by revealing the secret
// by the redeemer.
//
// redeemerPubKey must be x-only public key of the redeemer.
func RedeemLeaf(redeemerPubKey, secretHash []byte) (txscript.TapLeaf, error) {
	script, err := txscript.NewScriptBuilder().
		AddOp(txscript.OP_SHA256).
		AddData(secretHash).
		AddOp(txscript.OP_EQUALVERIFY).
		AddData(redeemerPubKey).
		AddOp(txscript.OP_CHECKSIG).
		Script()
	if err != nil {
		return txscript.TapLeaf{}, err
	}
	return txscript.NewBaseTapLeaf(script), nil
}

// RefundLeaf is one of the leaf scripts in the HTLC script which can be spent by the initiator
// after the lock time.
//
// initiatorPubKey must be x-only public key of the initiator.
func RefundLeaf(initiatorPubKey []byte, lockTime int64) (txscript.TapLeaf, error) {
	if lockTime > math.MaxUint16 {
		return txscript.TapLeaf{}, ErrInvalidLockTime
	}

	script, err := txscript.NewScriptBuilder().
		AddInt64(lockTime).
		AddOp(txscript.OP_CHECKSEQUENCEVERIFY).
		AddOp(txscript.OP_DROP).
		AddData(initiatorPubKey).
		AddOp(txscript.OP_CHECKSIG).
		Script()
	if err != nil {
		return txscript.TapLeaf{}, err
	}
	return txscript.NewBaseTapLeaf(script), nil
}

// InstantRefundLeaf is a 2 on 2 multisig leaf script
//
// pubkeys must be x-only pubkeys of the initiator and the redeemer.
func InstantRefundLeaf(initiatorPubKey, redeemerPubKey []byte) (txscript.TapLeaf, error) {
	script, err := txscript.NewScriptBuilder().
		AddData(initiatorPubKey).
		AddOp(txscript.OP_CHECKSIG).
		AddData(redeemerPubKey).
		AddOp(txscript.OP_CHECKSIGADD).
		AddOp(txscript.OP_2).
		AddOp(txscript.OP_NUMEQUAL).
		Script()
	if err != nil {
		return txscript.TapLeaf{}, err
	}
	return txscript.NewBaseTapLeaf(script), nil
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

// isWaitTimeOpCode returns if the given opCode is a valid opCode for a `OP_CHECKSEQUENCEVERIFY` params.
// Since we require the timelock to be an integer (max_uint16), so we interpret maximum 3 bytes.
func isWaitTimeOpCode(opCode byte) bool {
	return (opCode >= txscript.OP_1 && opCode <= txscript.OP_16) ||
		(opCode >= txscript.OP_DATA_1 && opCode <= txscript.OP_DATA_3)
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
