package wallet

import (
	"context"
	"fmt"
	"sync"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
)

type Wallet interface {
	Address() btcutil.Address

	Initiate(ctx context.Context, htlc *btc.HTLC) (string, error)

	Redeem(ctx context.Context, htlc *btc.HTLC, secret []byte) (string, error)

	Refund(ctx context.Context, htlc *btc.HTLC) (string, error)

	// InstantRefund(ctx context.Context, htlc btc.HTLC) (string, error)

	// Execute(ctx context.Context, htlcActions) (string, error)
}

type wallet struct {
	mu           *sync.Mutex
	key          *btcec.PrivateKey
	addrType     waddrmgr.AddressType
	addr         btcutil.Address
	network      *chaincfg.Params
	indexer      btc.IndexerClient
	feeEstimator btc.FeeEstimator
}

func NewWallet(network *chaincfg.Params, addrType waddrmgr.AddressType, key *btcec.PrivateKey, indexer btc.IndexerClient, estimator btc.FeeEstimator) (Wallet, error) {
	addr, err := btc.PublicKeyAddress(network, addrType, key.PubKey())
	if err != nil {
		return nil, err
	}
	return &wallet{
		mu:           new(sync.Mutex),
		key:          key,
		addrType:     addrType,
		addr:         addr,
		network:      network,
		indexer:      indexer,
		feeEstimator: estimator,
	}, nil
}

func (wal *wallet) Address() btcutil.Address {
	return wal.addr
}

func (wal *wallet) Initiate(ctx context.Context, htlc *btc.HTLC) (string, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Inputs
	utxos, err := wal.indexer.GetUTXOs(ctx, wal.addr)
	if err != nil {
		return "", err
	}
	sizeEstimator := wal.sizeEstimator(utxos)

	// Recipients
	htlcAddr, err := htlc.Address(wal.network)
	if err != nil {
		return "", err
	}
	recipients := []btc.Recipient{
		{
			To:     htlcAddr.EncodeAddress(),
			Amount: htlc.Amount,
		},
	}

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return "", err
	}

	// Build tx
	tx, err := btc.BuildTransaction(wal.network, feeRate.High, nil, utxos, sizeEstimator, recipients, wal.Address())
	if err != nil {
		return "", err
	}

	// Sign tx
	if err := wal.signTx(tx, utxos); err != nil {
		return "", err
	}

	// Submit tx
	if err := wal.indexer.SubmitTx(ctx, tx); err != nil {
		return "", err
	}
	return tx.TxHash().String(), nil
}

func (wal *wallet) Redeem(ctx context.Context, htlc *btc.HTLC, secret []byte) (string, error) {
	// todo : might not need the lock, since we're collecting funds from external source
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Make sure the HTLC is initiated and not redeemed yet
	addr, err := htlc.Address(wal.network)
	if err != nil {
		return "", err
	}
	utxos, err := wal.indexer.GetUTXOs(ctx, addr)
	if err != nil {
		return "", err
	}
	redeemable, _, err := htlc.Redeemable(utxos)
	if err != nil {
		return "", err
	}
	if !redeemable {
		return "", fmt.Errorf("HTLC is not redeemable")
	}

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return "", err
	}

	// Build tx
	sizeEstimator := btc.NewSizeEstimator(utxos, btc.BaseSizeHtlcRedeem, btc.SegwitSizeHtlcRedeem)
	tx, err := btc.BuildTransaction(wal.network, feeRate.High, utxos, nil, sizeEstimator, nil, wal.Address())
	if err != nil {
		return "", err
	}

	// Sign tx
	script, err := htlc.P2trScript()
	if err != nil {
		panic(err)
	}
	fetcher, err := btc.InitFetcher(utxos, script)
	if err != nil {
		return "", err
	}

	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
	for i, utxo := range tx.TxIn {
		out := fetcher.FetchPrevOutput(utxo.PreviousOutPoint)
		leaf, ctrBlk := htlc.RedeemLeaf()
		ctrBlkBytes, err := ctrBlk.ToBytes()
		if err != nil {
			return "", err
		}
		sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, out.Value, out.PkScript, leaf, txscript.SigHashAll, wal.key)
		if err != nil {
			return "", err
		}
		tx.TxIn[i].Witness = append(tx.TxIn[i].Witness, sig, secret, leaf.Script, ctrBlkBytes)
	}

	// Submit tx
	if err := wal.indexer.SubmitTx(ctx, tx); err != nil {
		return "", err
	}
	return tx.TxHash().String(), nil
}

func (wal *wallet) Refund(ctx context.Context, htlc *btc.HTLC) (string, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Make sure the HTLC is refundable
	addr, err := htlc.Address(wal.network)
	if err != nil {
		return "", err
	}
	utxos, err := wal.indexer.GetUTXOs(ctx, addr)
	if err != nil {
		return "", err
	}
	latest, err := wal.indexer.GetTipBlockHeight(ctx)
	if err != nil {
		return "", err
	}
	refundable := htlc.Refundable(utxos, latest)
	if !refundable {
		return "", fmt.Errorf("HTLC is not refundable")
	}

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return "", err
	}

	// Build tx
	sizeEstimator := btc.NewSizeEstimator(utxos, btc.BaseSizeHtlcRedeem, btc.SegwitSizeHtlcRedeem)
	tx, err := btc.BuildTransaction(wal.network, feeRate.High, utxos, nil, sizeEstimator, nil, wal.Address())
	if err != nil {
		return "", err
	}

	// Sign tx
	script, err := htlc.P2trScript()
	if err != nil {
		panic(err)
	}
	fetcher, err := btc.InitFetcher(utxos, script)
	if err != nil {
		return "", err
	}

	// Update tx inputs sequence to htlc timelock
	for i := range tx.TxIn {
		tx.TxIn[i].Sequence = uint32(htlc.Timelock)
	}
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
	leaf, ctrBlk := htlc.RefundLeaf()
	ctrBlkBytes, err := ctrBlk.ToBytes()
	if err != nil {
		return "", err
	}
	for i, utxo := range tx.TxIn {
		out := fetcher.FetchPrevOutput(utxo.PreviousOutPoint)
		sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, out.Value, out.PkScript, leaf, txscript.SigHashAll, wal.key)
		if err != nil {
			return "", err
		}
		tx.TxIn[i].Witness = append(tx.TxIn[i].Witness, sig, leaf.Script, ctrBlkBytes)
	}

	// Submit tx
	if err := wal.indexer.SubmitTx(ctx, tx); err != nil {
		return "", err
	}
	return tx.TxHash().String(), nil
}

//
// func (wal *wallet) InstantRefund(ctx context.Context, tx *wire.MsgTx, htlc *btc.HTLC) (string, error) {
// 	wal.mu.Lock()
// 	defer wal.mu.Unlock()
//
// 	// todo : Validate tx
//
// 	// Adding new utxos to cover the fees
// 	utxos, err := wal.indexer.GetUTXOs(ctx, wal.Address())
// 	if err != nil {
// 		return "", err
// 	}
// 	index := len(tx.TxIn)
//
// 	// Finish the tx to cover the fees
// 	feeRate, err := wal.feeEstimator.FeeSuggestion()
// 	if err != nil {
// 		return "", err
// 	}
// 	if err := btc.AddUtxoToCoverTxFees(tx, utxos, wal.sizeEstimator(utxos), feeRate.High, wal.addr); err != nil {
// 		return "", err
// 	}
//
// 	// Sign tx
// 	for i, input := range tx.TxIn {
// 		if i < index {
//
// 		} else {
// 			switch wal.addrType {
// 			case waddrmgr.PubKeyHash:
//
// 			}
// 		}
// 	}
// 	if err := wal.signTx(tx, utxos); err != nil {
// 		return "", err
// 	}
//
// 	// Make sure the HTLC is refundable
// 	addr, err := htlc.Address(wal.network)
// 	if err != nil {
// 		return "", err
// 	}
// 	utxos, err := wal.indexer.GetUTXOs(ctx, addr)
// 	if err != nil {
// 		return "", err
// 	}
// 	latest, err := wal.indexer.GetTipBlockHeight(ctx)
// 	if err != nil {
// 		return "", err
// 	}
// 	refundable := htlc.Refundable(utxos, latest)
// 	if !refundable {
// 		return "", fmt.Errorf("HTLC is not refundable")
// 	}
//
// 	// Fees
// 	feeRate, err := wal.feeEstimator.FeeSuggestion()
// 	if err != nil {
// 		return "", err
// 	}
//
// 	// Build tx
// 	sizeEstimator := btc.NewSizeEstimator(utxos, btc.BaseSizeHtlcRedeem, btc.SegwitSizeHtlcRedeem)
// 	tx, err := btc.BuildTransaction(wal.network, feeRate.High, utxos, nil, sizeEstimator, nil, wal.Address())
// 	if err != nil {
// 		return "", err
// 	}
//
// 	// Sign tx
// 	script, err := txscript.PayToTaprootScript(htlc.TapPubKey)
// 	if err != nil {
// 		return "", err
// 	}
// 	fetcher, err := btc.InitFetcher(utxos, script)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	// Update tx inputs sequence to htlc timelock
// 	for i := range tx.TxIn {
// 		tx.TxIn[i].Sequence = uint32(htlc.Timelock)
// 	}
// 	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
// 	leaf, ctrBlk := htlc.RefundLeaf()
// 	ctrBlkBytes, err := ctrBlk.ToBytes()
// 	if err != nil {
// 		return "", err
// 	}
// 	for i, utxo := range tx.TxIn {
// 		out := fetcher.FetchPrevOutput(utxo.PreviousOutPoint)
// 		sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, out.Value, out.PkScript, leaf, txscript.SigHashAll, wal.key)
// 		if err != nil {
// 			return "", err
// 		}
// 		tx.TxIn[i].Witness = append(tx.TxIn[i].Witness, sig, leaf.Script, ctrBlkBytes)
// 	}
//
// 	// Submit tx
// 	if err := wal.indexer.SubmitTx(ctx, tx); err != nil {
// 		return "", err
// 	}
// 	return tx.TxHash().String(), nil
// }

func (wal *wallet) sizeEstimator(utxos []btc.UTXO) *btc.SizeEstimator {
	switch wal.addrType {
	case waddrmgr.PubKeyHash:
		return btc.NewSizeEstimator(utxos, btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH)
	case waddrmgr.TaprootPubKey:
		return btc.NewSizeEstimator(utxos, btc.BaseSizeP2TR, btc.SegwitSizeP2TR)
	case waddrmgr.WitnessPubKey:
		return btc.NewSizeEstimator(utxos, btc.BaseSizeP2WPKH, btc.SegwitSizeP2WPKH)
	default:
		panic(fmt.Sprintf("unknown address type: %v", wal.addrType))
	}
}

func (wal *wallet) signTx(tx *wire.MsgTx, utxos []btc.UTXO) error {
	switch wal.addrType {
	case waddrmgr.PubKeyHash:
		return btc.SignP2pkhTx(wal.network, wal.key, tx)
	case waddrmgr.WitnessPubKey:
		return btc.SignP2wpkhTx(wal.network, utxos, wal.key, tx)
	case waddrmgr.TaprootPubKey:
		return btc.SignP2trTx(utxos, wal.key, tx)
	default:
		panic(fmt.Sprintf("unknown address type: %v", wal.addrType))
	}
}

// NewInstantRefundTx builds a new tx to redeem an HTLC instantly using the instantRefund branch. This requires mutural
// agreement between the two parties.
func NewInstantRefundTx(key *btcec.PrivateKey, htlc btc.HTLC, utxos []btc.UTXO, target btcutil.Address) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(btc.DefaultTxVersion)

	// Build tx
	for _, utxo := range utxos {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}
		txIn := wire.NewTxIn(wire.NewOutPoint(hash, utxo.Vout), nil, nil)
		tx.AddTxIn(txIn)

		toScript, err := txscript.PayToAddrScript(target)
		if err != nil {
			return nil, err
		}
		tx.AddTxOut(wire.NewTxOut(utxo.Amount, toScript))
	}

	// Sign the tx
	leaf, _ := htlc.InstantRefundLeaf()
	script, err := htlc.P2trScript()
	if err != nil {
		return nil, err
	}

	fetcher, err := btc.InitFetcher(utxos, script)
	if err != nil {
		return nil, err
	}
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)

	for i, utxo := range tx.TxIn {
		out := fetcher.FetchPrevOutput(utxo.PreviousOutPoint)
		sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, out.Value, out.PkScript, leaf, btc.SigHashSingleAnyoneCanPay, key)
		if err != nil {
			return nil, err
		}
		tx.TxIn[i].Witness = append(tx.TxIn[i].Witness, sig)
	}

	return tx, nil
}
