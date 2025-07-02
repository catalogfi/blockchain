package btc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
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
	BaseSizeHtlcRedeem        = 0
	BaseSizeHtlcRefund        = 0
	BaseSizeHtlcInstantRefund = 0

	// SegwitSizeHtlcRedeem = stack number + stack size * 4 + signature + secret + script size + control block
	SegwitSizeHtlcRedeem = func(secretSize int) int { return 1 + 4 + 64 + secretSize + 69 + 65 }
	// SegwitSizeHtlcRefund = stack number + stack size * 3 + signature + script size + control block
	SegwitSizeHtlcRefund = func(timelock int64) int {
		timelockSize := 0
		switch {
		case timelock <= 16:
			timelockSize = 1
		case timelock < 128:
			timelockSize = 2
		case timelock < 32768:
			timelockSize = 3
		default:
			timelockSize = 4
		}
		return 1 + 3 + 64 + (36 + timelockSize) + 97
	}
	// SegwitSizeHtlcInstantRefund = stack number + stack size * 4 + signature1 + signature2 + script size + control block
	SegwitSizeHtlcInstantRefund = 1 + 4 + 64 + 65 + 70 + 97
)

type HtlcActionType string

const (
	HtlcActionInitiate      HtlcActionType = "initiate"
	HtlcActionRedeem        HtlcActionType = "redeem"
	HtlcActionRefund        HtlcActionType = "refund"
	HtlcActionInstantRefund HtlcActionType = "instantRefund"
)

type HtlcAction struct {
	ActionType      HtlcActionType
	Htlc            *HTLC
	InstantRefundTx *wire.MsgTx
	RefundTo        btcutil.Address
}

type HTLC struct {
	InitiatorPubKey []byte
	RedeemerPubKey  []byte
	SecretHash      []byte
	Timelock        int64
	Amount          int64

	secret []byte
	tree   *txscript.IndexedTapScriptTree
}

func NewHTLC(initiatorPubKey, redeemerPubKey, secretHash []byte, timelock, amount int64) (*HTLC, error) {
	// todo : validate the input before generating these leafs
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

// Address returns the address of the htlc script. We panic for the error since it shouldn't happen, this will make
// things easier for the caller.
func (htlc *HTLC) Address(network *chaincfg.Params) btcutil.Address {
	if htlc.tree == nil {
		panic(fmt.Errorf("empty htlc tree"))
	}
	rootHash := htlc.tree.RootNode.TapHash()
	outputKey := txscript.ComputeTaprootOutputKey(GardenNums, rootHash[:])
	addr, err := PublicKeyAddress(network, waddrmgr.TaprootPubKey, outputKey)
	if err != nil {
		panic(err)
	}
	return addr
}

func (htlc *HTLC) SetSecret(secret []byte) {
	htlc.secret = secret
}

func (htlc *HTLC) Secret() []byte {
	return htlc.secret
}

func (htlc *HTLC) P2trScript() ([]byte, error) {
	rootHash := htlc.tree.RootNode.TapHash()
	outputKey := txscript.ComputeTaprootOutputKey(GardenNums, rootHash[:])
	return txscript.PayToTaprootScript(outputKey)
}

// Redeemable checks the utxos to see if the htlc is redeemable. This means there's a single confirmed utxo which has
// enough amount to cover the htlc amount. We currently don't allow initiation of HTLCs with multiple utxos. The given
// utxos must be from the htlc address.
func (htlc *HTLC) Redeemable(utxos []UTXO) (bool, uint64, error) {
	for _, utxo := range utxos {
		// todo : do we need to force the amount to be equal?
		if utxo.Status != nil && utxo.Status.Confirmed && utxo.Amount >= htlc.Amount {
			return true, *utxo.Status.BlockHeight, nil
		}
	}

	return false, 0, nil
}

// Refundable checks if the htlc is refundable. It takes the utxos and the latest block height as input. It finds the
// utxo with enough amount first and then check if it's expired for refunding.
func (htlc *HTLC) Refundable(utxos []UTXO, latest uint64) bool {
	for _, utxo := range utxos {
		// todo : do we need to force the amount to be equal?
		if utxo.Status != nil && utxo.Status.Confirmed && utxo.Amount >= htlc.Amount {
			if latest-*utxo.Status.BlockHeight+1 >= uint64(htlc.Timelock) {
				return true
			}
		}
	}
	return false
}

// Utxo finds the initiation utxo of the htlc.
func (htlc *HTLC) Utxo(ctx context.Context, network *chaincfg.Params, indexer IndexerClient) (UTXO, error) {
	addr := htlc.Address(network)
	utxos, err := indexer.GetUTXOs(ctx, addr)
	if err != nil {
		return UTXO{}, err
	}
	for _, utxo := range utxos {
		if utxo.Status != nil && utxo.Status.Confirmed && utxo.Amount >= htlc.Amount {
			return utxo, nil
		}
	}

	return UTXO{}, fmt.Errorf("not initiated")
}

func (htlc *HTLC) RefundableUtxos(ctx context.Context, network *chaincfg.Params, indexer IndexerClient) ([]UTXO, error) {
	addr := htlc.Address(network)
	utxos, err := indexer.GetUTXOs(ctx, addr)
	if err != nil {
		return nil, err
	}
	if len(utxos) == 0 {
		return nil, fmt.Errorf("no refundable utxo")
	}
	refundableUtxos := make([]UTXO, 0, len(utxos))
	latest, err := indexer.GetTipBlockHeight(ctx)
	if err != nil {
		return nil, err
	}
	for _, utxo := range utxos {
		if utxo.Status != nil && utxo.Status.Confirmed && utxo.Amount > DustAmount {
			if latest-*utxo.Status.BlockHeight+1 >= uint64(htlc.Timelock) {
				refundableUtxos = append(refundableUtxos, utxo)
			}
		}
	}
	return refundableUtxos, nil
}

// Leaf returns the tapLeaf associated with given action.
func (htlc *HTLC) Leaf(action HtlcActionType) (txscript.TapLeaf, txscript.ControlBlock) {
	var leaf txscript.TapLeaf
	switch action {
	case HtlcActionRedeem:
		leaf = htlc.tree.RootNode.Right().(txscript.TapLeaf)
	case HtlcActionRefund:
		leaf = htlc.tree.RootNode.Left().Right().(txscript.TapLeaf)
	case HtlcActionInstantRefund:
		leaf = htlc.tree.RootNode.Left().Left().(txscript.TapLeaf)
	default:
		panic("invalid action type")
	}

	index := htlc.tree.LeafProofIndex[leaf.TapHash()]
	ctrBlk := htlc.tree.LeafMerkleProofs[index].ToControlBlock(GardenNums)
	return leaf, ctrBlk
}

// HtlcActionFromWitness determines the type of HTLC (Hashed Timelock Contract) action based on the provided witness.
// It returns the corresponding HtlcActionType and an error if the witness does not match any known HTLC action types.
func HtlcActionFromWitness(witness wire.TxWitness) (HtlcActionType, error) {
	switch len(witness) {
	case 4: // redeem or instant refund
		ok, _ := IsMultiSigLeaf(witness[2])
		if ok {
			return HtlcActionInstantRefund, nil
		}
		ok, _ = IsRedeemLeaf(witness[2])
		if ok {
			return HtlcActionRedeem, nil
		}
	case 3: // refund
		ok, _ := IsRefundLeaf(witness[1])
		if ok {
			return HtlcActionRefund, nil
		}
	}
	return "", fmt.Errorf("witness not associated with any htlc action")
}

// RedeemLeaf is one of the leaf scripts in the HTLC script which can be spent by revealing the secret
// by the redeemer. `redeemerPubKey` must be x-only public key of the redeemer.
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
// after the lock time. `initiatorPubKey` must be x-only public key of the initiator.
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

// InstantRefundLeaf is a 2 on 2 multisig leaf script. This is used for fast refunds when both parties agree.
// both the `initiatorPubKey` and `redeemerPubKey` must be x-only pubkeys.
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

// ValidateInstantRefundTx checks if the given tx is a valid instant refund tx of the htlc. It assumes the tx will only
// have one input and one output. Input's witness stack will only contain one item which is the initiator's signature.
func ValidateInstantRefundTx(htlc *HTLC, tx *wire.MsgTx, network *chaincfg.Params) (UTXO, *wire.TxOut, error) {
	if len(tx.TxIn) != 1 {
		return UTXO{}, nil, errors.New("invalid number of inputs")
	}
	if len(tx.TxOut) != 1 {
		return UTXO{}, nil, errors.New("invalid number of outputs")
	}
	if len(tx.TxIn[0].Witness) != 1 {
		return UTXO{}, nil, errors.New("invalid witness length")
	}
	sigBytes := tx.TxIn[0].Witness[0]
	amount := tx.TxOut[0].Value
	if amount < DustAmount {
		return UTXO{}, nil, errors.New("amount lower than dust amount")
	}

	// Verify signature
	script, err := htlc.P2trScript()
	if err != nil {
		return UTXO{}, nil, err
	}
	fetcher := txscript.NewCannedPrevOutputFetcher(script, amount)
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
	leaf, _ := htlc.Leaf(HtlcActionInstantRefund)
	tapSigHashes, err := txscript.CalcTapscriptSignaturehash(sigHashes, SigHashSingleAnyoneCanPay, tx, 0, fetcher, leaf)
	if err != nil {
		return UTXO{}, nil, err
	}
	if len(sigBytes) == schnorr.SignatureSize+1 {
		sigBytes = sigBytes[:len(sigBytes)-1]
	}
	signature, err := schnorr.ParseSignature(sigBytes)
	if err != nil {
		return UTXO{}, nil, err
	}
	pub, err := schnorr.ParsePubKey(htlc.InitiatorPubKey)
	if err != nil {
		return UTXO{}, nil, err
	}
	if ok := signature.Verify(tapSigHashes, pub); !ok {
		return UTXO{}, nil, errors.New("invalid signature")
	}

	// Parse the input's utxo and output address from the tx
	utxo := UTXO{
		TxID:   tx.TxIn[0].PreviousOutPoint.Hash.String(),
		Vout:   tx.TxIn[0].PreviousOutPoint.Index,
		Amount: amount,
	}
	return utxo, tx.TxOut[0], nil
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
