package btc

import (
	"crypto/sha256"
	"math"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

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
func RefundLeaf(initiatorPubkey []byte, lockTime int64) (txscript.TapLeaf, error) {
	if lockTime > math.MaxUint16 || lockTime < 0 {
		return txscript.TapLeaf{}, ErrInvalidLockTime
	}
	script, err := txscript.NewScriptBuilder().
		AddOp(byte(lockTime)).
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

// InstantRefundLeaf is one of the leaf scripts in the HTLC script which can be spent by the initiator
// if both parties agree to refund instantly.
//
// pubkeys must be x-only pubkeys of the initiator and the redeemer.
func InstantRefundLeaf(initiatorPubkey, redeemerPubkey []byte) (txscript.TapLeaf, error) {
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

type HTLC struct {
	InitiatorPubkey []byte
	RedeemerPubkey  []byte
	SecretHash      []byte
	LockTime        int64
}

type HTLCWallet interface {
	Initiate(htlc *HTLC, amount uint64) (string, error)
	Redeem(htlc *HTLC, secret []byte) (string, error)
	Refund(htlc *HTLC, sig []byte) (string, error)
	Address(htlc *HTLC) (string, error)

	Status(id string) (*btcjson.TxRawResult, bool, error)
}

type htlcWallet struct {
	// Wallet here could be a batcher wallet or RBF wallet or simple wallet
	wallet      Wallet
	chain       *chaincfg.Params
	internalKey *btcec.PublicKey
}

func NewHTLCWallet(wallet Wallet, chain *chaincfg.Params) *htlcWallet {
	return &htlcWallet{
		wallet:      wallet,
		chain:       chain,
		internalKey: gardenNums(),
	}

}

func (w *htlcWallet) Initiate(htlc *HTLC, amount uint64) (string, error) {
	addr, err := w.Address(htlc)
	if err != nil {
		return "", err
	}
	return w.wallet.AddOutput(TxOutputRequest{
		Value:   amount,
		Address: addr,
	})
}

// Generates the p2tr output key using the provided htlc
func (w *htlcWallet) Address(htlc *HTLC) (string, error) {

	redeemLeaf, err := RedeemLeaf(htlc.RedeemerPubkey, htlc.SecretHash)
	if err != nil {
		return "", err
	}
	refundLeaf, err := RefundLeaf(htlc.InitiatorPubkey, htlc.LockTime)
	if err != nil {
		return "", err
	}
	instantRefundLeaf, err := InstantRefundLeaf(htlc.InitiatorPubkey, htlc.RedeemerPubkey)
	if err != nil {
		return "", err
	}

	refundsBranch := txscript.NewTapBranch(refundLeaf, instantRefundLeaf)

	rootBranch := txscript.NewTapBranch(redeemLeaf, refundsBranch)
	rootHash := rootBranch.TapHash()

	outputKey := txscript.ComputeTaprootOutputKey(w.internalKey, rootHash[:])

	addr, err := btcutil.NewAddressTaproot(outputKey.X().Bytes(), w.chain)
	if err != nil {
		return "", err
	}

	return addr.String(), nil

}

func gardenNums() *btcec.PublicKey {
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
		panic(err)
	}

	r := sha256.Sum256([]byte("GardenHTLC"))
	_, rG := btcec.PrivKeyFromBytes(r[:])

	numsX, numsY := curve.Add(H.X(), H.Y(), rG.X(), rG.Y())
	if !curve.IsOnCurve(numsX, numsY) {
		panic("invalid nums")
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
	return nums
}

// Helper function to convert a script to a leaf with version 0xc0.
func toLeaf(script []byte) txscript.TapLeaf {
	return txscript.NewTapLeaf(0xc0, script)
}
