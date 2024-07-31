package btc

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type HTLCAction string

const (
	InitiateHTLCAction      HTLCAction = "initiate"
	RedeemHTLCAction        HTLCAction = "redeem"
	RefundHTLCAction        HTLCAction = "refund"
	InstantRefundHTLCAction HTLCAction = "instantRefund"
)

type RawHTLCAction struct {
	Action HTLCAction
	HTLC   HTLC
	// Amount to initiate the HTLC.
	// Only used in the case of InitiateHTLCAction
	Amount int64
	// Only used in the case of RedeemHTLCAction.
	Secret []byte
	// Only used in the case of RefundHTLCAction.
	InsantRefundSACPTxBytes []byte
}

var (
	// ErrSACPInvalidInputsLen is returned when the signed inputs are less than the available inputs
	// Occurs in the case of instant refunds
	ErrSACPInvalidInputsLen = fmt.Errorf("signed inputs less than available inputs")

	// ErrSACPInvalidInput is returned when the input in the signed tx is not the same as the utxo
	// Occurs in the case of instant refunds
	ErrSACPInvalidInput = fmt.Errorf("invalid signed input")

	// ErrSACPInvalidOutput is returned when the output script in the signed tx is not the same as the wallet address
	ErrSACPInvalidOutput = fmt.Errorf("invalid output script in the signed tx")

	// ErrInvalidSecret is returned when the secret is invalid while redeeming the HTLC
	ErrInvalidSecret = fmt.Errorf("invalid secret")

	// ErrInvalidInstantRefundSACPWitnessLen is returned when the witness length of the SACP tx is invalid
	ErrInvalidInstantRefundSACPWitnessLen = fmt.Errorf("invalid witness length")

	// ErrInvalidControlBlock is returned when the control block in the signed SACP tx is invalid
	ErrInvalidControlBlock = fmt.Errorf("invalid control block")

	// ErrInvalidInstantRefundScript is returned when the instant refund script is invalid in the SACP tx
	ErrInvalidInstantRefundScript = fmt.Errorf("invalid instant refund script")

	ErrHTLCNeedMoreBlocks = func(blocks uint64) error { return fmt.Errorf("need more %d blocks to refund", blocks) }
)

type HTLC struct {
	// X-only pubkey of the initiator
	InitiatorPubkey []byte
	// X-only pubkey of the redeemer
	RedeemerPubkey []byte
	SecretHash     []byte
	// Locktime in blocks
	Timelock uint32
}

type HTLCWallet interface {
	// Initiate sends the amount to the HTLC address
	Initiate(ctx context.Context, htlc *HTLC, amount int64) (string, error)
	// Redeem redeems the HTLC with the secret
	Redeem(ctx context.Context, htlc *HTLC, secret []byte) (string, error)
	// Refund refunds the HTLC if the htlc is expired.
	// For instant refunds, the SACP tx signed by counterparty should be passed
	Refund(ctx context.Context, htlc *HTLC, instantRefundSACPTx []byte) (string, error)
	// GenerateInstantRefundSACP generates the SACP tx needed for the instant refunds.
	//
	// Signature is added at the first index of the witness of the transaction inputs.
	GenerateInstantRefundSACP(ctx context.Context, htlc *HTLC, recipient btcutil.Address) ([]byte, error)
	// Address returns the tapscript address of the HTLC
	Address(htlc *HTLC) (btcutil.Address, error)
	// Status returns the transaction if submitted and bool indicating whether the transaction
	// is submitted or not
	Status(ctx context.Context, id string) (Transaction, bool, error)

	// Execute executes passed HTLC actions
	Execute(ctx context.Context, htlcActions []RawHTLCAction) (string, error)
}

type htlcWallet struct {
	// Wallet here could be a batcher wallet or RBF wallet or simple wallet
	wallet      Wallet
	chain       *chaincfg.Params
	internalKey *btcec.PublicKey
	indexer     IndexerClient
}

func NewHTLCWallet(wallet Wallet, indexer IndexerClient, chain *chaincfg.Params) (HTLCWallet, error) {
	internalKey, err := GardenNUMS()
	if err != nil {
		return nil, err
	}

	return &htlcWallet{
		wallet:      wallet,
		chain:       chain,
		internalKey: internalKey,
		indexer:     indexer,
	}, nil
}

func (hw *htlcWallet) Status(ctx context.Context, id string) (Transaction, bool, error) {
	return hw.wallet.Status(ctx, id)
}

// Address returns the tapscript address of the HTLC
func (hw *htlcWallet) Address(htlc *HTLC) (btcutil.Address, error) {

	leaves, err := htlcLeaves(htlc)
	if err != nil {
		return nil, err
	}

	tapScriptTree := txscript.AssembleTaprootScriptTree(leaves.ToArray()...)

	tapScriptRootHash := tapScriptTree.RootNode.TapHash()
	outputKey := txscript.ComputeTaprootOutputKey(
		hw.internalKey, tapScriptRootHash[:],
	)

	addr, err := btcutil.NewAddressTaproot(outputKey.X().Bytes(), hw.chain)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// GenerateInstantRefundSACP generates the SACP tx needed for the instant refunds
func (hw *htlcWallet) GenerateInstantRefundSACP(ctx context.Context, htlc *HTLC, recipient btcutil.Address) ([]byte, error) {
	instantRefundLeaf, cbBytes, err := getControlBlock(hw.internalKey, htlc, LeafInstantRefund)
	if err != nil {
		return nil, err
	}
	witness := [][]byte{
		AddSignatureSchnorrOp,
		// insert random sig placeholder for other parties to insert their signature
		// this is for proper fee calculation
		randomSig(),
		instantRefundLeaf.Script,
		cbBytes,
	}

	scriptAddr, err := hw.Address(htlc)
	if err != nil {
		return nil, err
	}

	txBytes, err := hw.wallet.GenerateSACP(ctx, SpendRequest{
		Witness:       witness,
		Leaf:          instantRefundLeaf,
		ScriptAddress: scriptAddr,
		HashType:      SigHashSingleAnyoneCanPay,
	}, recipient)
	if err != nil {
		return nil, err
	}

	return txBytes, nil
}

// Initiate sends the amount to the HTLC address
func (hw *htlcWallet) Initiate(ctx context.Context, htlc *HTLC, amount int64) (string, error) {
	addr, err := hw.Address(htlc)
	if err != nil {
		return "", err
	}
	return hw.wallet.Send(ctx, []SendRequest{
		{
			To:     addr,
			Amount: amount,
		},
	}, nil, nil)
}

func (hw *htlcWallet) redeem(htlc *HTLC, secret []byte) (SpendRequest, error) {
	if !isSecretValid(secret, htlc) {
		return SpendRequest{}, ErrInvalidSecret
	}

	redeemTapLeaf, cbBytes, err := getControlBlock(hw.internalKey, htlc, LeafRedeem)
	if err != nil {
		return SpendRequest{}, err
	}
	witness := [][]byte{
		AddSignatureSchnorrOp,
		secret,
		redeemTapLeaf.Script,
		cbBytes,
	}

	scriptAddr, err := hw.Address(htlc)
	if err != nil {
		return SpendRequest{}, err
	}

	return SpendRequest{
		Witness:       witness,
		Leaf:          redeemTapLeaf,
		ScriptAddress: scriptAddr,
		HashType:      txscript.SigHashAll,
	}, nil
}

// Redeem redeems the HTLC with the secret
func (hw *htlcWallet) Redeem(ctx context.Context, htlc *HTLC, secret []byte) (string, error) {
	redeemSpendRequest, err := hw.redeem(htlc, secret)
	if err != nil {
		return "", err
	}

	return hw.wallet.Send(ctx, nil, []SpendRequest{
		redeemSpendRequest,
	}, nil)
}

// instantRefund refunds given the counterparty signed SACP tx
func (hw *htlcWallet) instantRefund(ctx context.Context, htlc *HTLC, instantRefundSACPTx []byte) ([]byte, error) {
	if instantRefundSACPTx == nil {
		return nil, fmt.Errorf("instantRefundSACPTx is nil")
	}
	scriptAddr, err := hw.Address(htlc)
	if err != nil {
		return nil, err
	}

	utxos, err := hw.indexer.GetUTXOs(ctx, scriptAddr)
	if err != nil {
		return nil, err
	}
	utxoValue := int64(0)
	for _, utxo := range utxos {
		utxoValue += utxo.Amount
	}
	instandRefundLeaf, cbBytes, err := getControlBlock(hw.internalKey, htlc, LeafInstantRefund)
	if err != nil {
		return nil, err
	}
	tx, err := validateInstantRefundSACP(instantRefundSACPTx, utxos, hw.wallet.Address(), cbBytes, instandRefundLeaf)
	if err != nil {
		return nil, err
	}

	instantRefundWitness := [][]byte{
		AddSignatureSchnorrOp,
		// insert invalid placeholder signature for other parties to insert their signature
		// this is for proper fee calculation
		randomSig(),
		instandRefundLeaf.Script,
		cbBytes,
	}

	for i := range tx.TxIn {
		// 0th index is the signature of this wallet
		witnessWithSig, err := hw.wallet.SignSACPTx(tx, i, utxoValue, instandRefundLeaf, scriptAddr, instantRefundWitness)
		if err != nil {
			return nil, err
		}

		// Format the witness to include the signature of the initiator at the 1st index of the witness
		// 0th index should be the signature of the redeemer
		// 1st index should be the signature of the initiator
		tx.TxIn[i].Witness[1] = witnessWithSig[0]
	}

	return GetTxRawBytes(tx)
}

func (hw *htlcWallet) send(ctx context.Context, sends []SendRequest, spends []SpendRequest, sacps [][]byte) (string, error) {
	return hw.wallet.Send(ctx, sends, spends, sacps)
}

func (hw *htlcWallet) Execute(ctx context.Context, htlcActions []RawHTLCAction) (string, error) {

	// aggregate all sends
	var sends []SendRequest
	var spends []SpendRequest
	var sacps [][]byte

	for _, htlcAction := range htlcActions {
		switch htlcAction.Action {
		case InitiateHTLCAction:
			addr, err := hw.Address(&htlcAction.HTLC)
			if err != nil {
				return "", err
			}
			sends = append(sends, SendRequest{
				To:     addr,
				Amount: htlcAction.Amount,
			})
		case RedeemHTLCAction:
			redeemSpendRequest, err := hw.redeem(&htlcAction.HTLC, htlcAction.Secret)
			if err != nil {
				return "", err
			}
			spends = append(spends, redeemSpendRequest)
		case RefundHTLCAction:
			refundSpendRequest, err := hw.refund(&htlcAction.HTLC)
			if err != nil {
				return "", err
			}
			spends = append(spends, refundSpendRequest)
		case InstantRefundHTLCAction:
			refundSACP, err := hw.instantRefund(ctx, &htlcAction.HTLC, htlcAction.InsantRefundSACPTxBytes)
			if err != nil {
				fmt.Println("dcfgjvhbkjnkmjhkgjfdxfcgvhbjnkm")
				return "", err
			}
			sacps = append(sacps, refundSACP)
		}

	}
	return hw.send(ctx, sends, spends, sacps)
}

func (hw *htlcWallet) refund(htlc *HTLC) (SpendRequest, error) {

	scriptAddr, err := hw.Address(htlc)
	if err != nil {
		return SpendRequest{}, err
	}

	currentTip, err := hw.indexer.GetTipBlockHeight(context.Background())
	if err != nil {
		return SpendRequest{}, err
	}

	utxos, err := hw.indexer.GetUTXOs(context.Background(), scriptAddr)
	if err != nil {
		return SpendRequest{}, err
	}

	canRefund, needMoreBlocks := canRefund(utxos, htlc, currentTip)
	if !canRefund {
		return SpendRequest{}, ErrHTLCNeedMoreBlocks(needMoreBlocks)
	}
	tapLeaf, cbBytes, err := getControlBlock(hw.internalKey, htlc, LeafRefund)
	if err != nil {
		return SpendRequest{}, err
	}
	witness := [][]byte{
		AddSignatureSchnorrOp,
		tapLeaf.Script,
		cbBytes,
	}

	return SpendRequest{
		Witness:       witness,
		Leaf:          tapLeaf,
		ScriptAddress: scriptAddr,
		HashType:      txscript.SigHashAll,
		Sequence:      htlc.Timelock,
	}, nil
}

func (hw *htlcWallet) Refund(ctx context.Context, htlc *HTLC, sigTx []byte) (string, error) {
	if sigTx != nil {
		sacp, err := hw.instantRefund(ctx, htlc, sigTx)
		if err != nil {
			return "", err
		}
		return hw.send(ctx, nil, nil, [][]byte{sacp})
	}

	refundSpendRequest, err := hw.refund(htlc)
	if err != nil {
		return "", err
	}

	return hw.wallet.Send(ctx, nil, []SpendRequest{
		refundSpendRequest,
	}, nil)
}

// ------------------ Helper functions ------------------

func isSecretValid(secret []byte, htlc *HTLC) bool {
	hash := sha256.Sum256(secret)
	return bytes.Equal(hash[:], htlc.SecretHash)
}

func validateInstantRefundSACP(refundSACP []byte, utxos []UTXO, recipient btcutil.Address, cb []byte, instantRefundLeaf txscript.TapLeaf) (*wire.MsgTx, error) {
	btcTx, err := btcutil.NewTxFromBytes(refundSACP)
	if err != nil {
		return nil, err
	}

	tx := btcTx.MsgTx()

	// Validate the tx
	if len(tx.TxIn) != len(utxos) {
		return nil, ErrSACPInvalidInputsLen
	}

	pkScript, err := txscript.PayToAddrScript(recipient)
	if err != nil {
		return nil, err
	}

	// Check if txHashs match with the utxos
	for i, txIn := range tx.TxIn {
		if txIn.PreviousOutPoint.Hash.String() != utxos[i].TxID {
			return nil, ErrSACPInvalidInput
		}
		if !bytes.Equal(tx.TxOut[i].PkScript, pkScript) {
			return nil, ErrSACPInvalidOutput
		}
	}

	// witness should have 4 elements
	if len(tx.TxIn[0].Witness) != 4 {
		return nil, ErrInvalidInstantRefundSACPWitnessLen
	}

	//TODO: check if the signature is valid

	// first two should be signature lens
	if len(tx.TxIn[0].Witness[0]) != 65 || len(tx.TxIn[0].Witness[1]) != 65 {
		return nil, ErrInvalidInstantRefundSACPWitnessLen
	}
	// instant refund script should be the same
	if !bytes.Equal(tx.TxIn[0].Witness[2], instantRefundLeaf.Script) {
		return nil, ErrInvalidInstantRefundScript
	}
	// control block should be the same
	if !bytes.Equal(tx.TxIn[0].Witness[3], cb) {
		return nil, ErrInvalidControlBlock
	}

	return tx, nil
}

// checks if the utxos can be refunded
func canRefund(utxos []UTXO, htlc *HTLC, currentTip uint64) (bool, uint64) {
	for _, utxo := range utxos {
		needMoreBlocks := uint64(0)
		timelock := uint64(htlc.Timelock)

		// check if utxo has been expired
		if utxo.Status.Confirmed && *utxo.Status.BlockHeight+timelock-1 > currentTip {
			needMoreBlocks = *utxo.Status.BlockHeight + timelock - currentTip + 1
		} else if !utxo.Status.Confirmed {
			needMoreBlocks = timelock + 1
		}
		if needMoreBlocks > 0 {
			return false, needMoreBlocks
		}
	}
	return true, 0
}

func randomSig() []byte {
	sig := make([]byte, 65)
	_, _ = rand.Read(sig)
	return sig
}

func htlcLeaves(htlc *HTLC) (*htlcTapLeaves, error) {
	redeemLeaf, err := RedeemLeaf(htlc.RedeemerPubkey, htlc.SecretHash)
	if err != nil {
		return &htlcTapLeaves{}, err
	}
	refundLeaf, err := RefundLeaf(htlc.InitiatorPubkey, htlc.Timelock)
	if err != nil {
		return &htlcTapLeaves{}, err
	}

	instantRefundLeaf, err := MultiSigLeaf(htlc.InitiatorPubkey, htlc.RedeemerPubkey)
	if err != nil {
		return &htlcTapLeaves{}, err
	}
	return newLeaves(redeemLeaf, refundLeaf, instantRefundLeaf)
}

// Helper struct to manage HTLC leaves
type htlcTapLeaves struct {
	redeem        txscript.TapLeaf
	refund        txscript.TapLeaf
	instantRefund txscript.TapLeaf
}

type Leaf string

const (
	LeafRedeem        Leaf = "redeem"
	LeafRefund        Leaf = "refund"
	LeafInstantRefund Leaf = "instantRefund"
)

func newLeaves(redeem, refund, instantRefund txscript.TapLeaf) (*htlcTapLeaves, error) {

	return &htlcTapLeaves{
		redeem:        redeem,
		refund:        refund,
		instantRefund: instantRefund,
	}, nil
}

func (l *htlcTapLeaves) ToArray() []txscript.TapLeaf {
	return []txscript.TapLeaf{l.redeem, l.refund, l.instantRefund}
}

func getControlBlock(internalKey *btcec.PublicKey, htlc *HTLC, leaf Leaf) (txscript.TapLeaf, []byte, error) {
	leaves, err := htlcLeaves(htlc)
	if err != nil {
		return txscript.TapLeaf{}, nil, err
	}

	tapScriptTree := txscript.AssembleTaprootScriptTree(leaves.ToArray()...)

	controlBlock := tapScriptTree.LeafMerkleProofs[leaves.IndexOf(leaf)].ToControlBlock(
		internalKey,
	)

	cbBytes, err := controlBlock.ToBytes()
	if err != nil {
		return txscript.TapLeaf{}, nil, err
	}
	tapLeaf, err := leaves.GetTapLeaf(leaf)
	if err != nil {
		return txscript.TapLeaf{}, nil, err
	}

	return tapLeaf, cbBytes, nil
}

func (l *htlcTapLeaves) GetTapLeaf(leaf Leaf) (txscript.TapLeaf, error) {
	switch leaf {
	case LeafRedeem:
		return l.redeem, nil
	case LeafRefund:
		return l.refund, nil
	case LeafInstantRefund:
		return l.instantRefund, nil
	default:
		return txscript.TapLeaf{}, nil
	}
}

func (l *htlcTapLeaves) IndexOf(leaf Leaf) int {
	var leafScript []byte

	switch leaf {
	case LeafRedeem:
		leafScript = l.redeem.Script
	case LeafRefund:
		leafScript = l.refund.Script
	case LeafInstantRefund:
		leafScript = l.instantRefund.Script
	default:
		return -1
	}

	leaves := l.ToArray()
	for i, l := range leaves {
		if bytes.Equal(l.Script, leafScript) {
			return i
		}
	}
	return -1
}
