package wallet

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
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

	Initiate(ctx context.Context, htlc *btc.HTLC) (*wire.MsgTx, error)

	Redeem(ctx context.Context, htlc *btc.HTLC, secret []byte) (*wire.MsgTx, error)

	Refund(ctx context.Context, htlc *btc.HTLC) (*wire.MsgTx, error)

	InstantRefund(ctx context.Context, htlc *btc.HTLC, tx *wire.MsgTx) (*wire.MsgTx, error)

	// Execute(ctx context.Context, htlcActions) (string, error)
}

type wallet struct {
	mu           *sync.Mutex
	network      *chaincfg.Params
	key          *btcec.PrivateKey
	addrType     waddrmgr.AddressType
	addr         btcutil.Address
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

func (wal *wallet) Initiate(ctx context.Context, htlc *btc.HTLC) (*wire.MsgTx, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Inputs
	utxos, err := wal.indexer.GetUTXOs(ctx, wal.addr)
	if err != nil {
		return nil, err
	}
	sizer := btc.NewSizeEstimatorOfAddrType(utxos, wal.addrType)

	// Recipients
	htlcAddr, err := htlc.Address(wal.network)
	if err != nil {
		return nil, err
	}
	recipients := btc.SingleRecipient(htlcAddr.EncodeAddress(), htlc.Amount)

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, err
	}

	// Build tx
	tx, err := btc.BuildTransaction(wal.network, feeRate.High, nil, utxos, sizer, recipients, wal.Address())
	if err != nil {
		return nil, err
	}

	// Sign tx
	if err := wal.signTx(tx, utxos); err != nil {
		return nil, err
	}

	// Submit tx
	if err := wal.indexer.SubmitTx(ctx, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (wal *wallet) Redeem(ctx context.Context, htlc *btc.HTLC, secret []byte) (*wire.MsgTx, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Make sure the HTLC is initiated and not redeemed yet
	addr, err := htlc.Address(wal.network)
	if err != nil {
		return nil, err
	}
	utxos, err := wal.indexer.GetUTXOs(ctx, addr)
	if err != nil {
		return nil, err
	}
	redeemable, _, err := htlc.Redeemable(utxos)
	if err != nil {
		return nil, err
	}
	if !redeemable {
		return nil, fmt.Errorf("HTLC is not redeemable")
	}

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, err
	}

	// Build tx
	sizeEstimator := btc.NewSizeEstimator(utxos, btc.BaseSizeHtlcRedeem, btc.SegwitSizeHtlcRedeem)
	tx, err := btc.BuildTransaction(wal.network, feeRate.High, utxos, nil, sizeEstimator, nil, wal.Address())
	if err != nil {
		return nil, err
	}

	// Sign tx
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
		leaf, ctrBlk := htlc.RedeemLeaf()
		ctrBlkBytes, err := ctrBlk.ToBytes()
		if err != nil {
			return nil, err
		}
		sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, out.Value, out.PkScript, leaf, txscript.SigHashAll, wal.key)
		if err != nil {
			return nil, err
		}
		tx.TxIn[i].Witness = append(tx.TxIn[i].Witness, sig, secret, leaf.Script, ctrBlkBytes)
	}

	// Submit tx
	if err := wal.indexer.SubmitTx(ctx, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (wal *wallet) Refund(ctx context.Context, htlc *btc.HTLC) (*wire.MsgTx, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Make sure the HTLC is refundable
	addr, err := htlc.Address(wal.network)
	if err != nil {
		return nil, err
	}
	utxos, err := wal.indexer.GetUTXOs(ctx, addr)
	if err != nil {
		return nil, err
	}
	latest, err := wal.indexer.GetTipBlockHeight(ctx)
	if err != nil {
		return nil, err
	}
	refundable := htlc.Refundable(utxos, latest)
	if !refundable {
		return nil, fmt.Errorf("HTLC is not refundable")
	}

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, err
	}

	// Build tx
	sizeEstimator := btc.NewSizeEstimator(utxos, btc.BaseSizeHtlcRedeem, btc.SegwitSizeHtlcRedeem)
	tx, err := btc.BuildTransaction(wal.network, feeRate.High, utxos, nil, sizeEstimator, nil, wal.Address())
	if err != nil {
		return nil, err
	}

	// Sign tx
	script, err := htlc.P2trScript()
	if err != nil {
		return nil, err
	}
	fetcher, err := btc.InitFetcher(utxos, script)
	if err != nil {
		return nil, err
	}

	// Update tx inputs sequence to htlc timelock
	for i := range tx.TxIn {
		tx.TxIn[i].Sequence = uint32(htlc.Timelock)
	}
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
	leaf, ctrBlk := htlc.RefundLeaf()
	ctrBlkBytes, err := ctrBlk.ToBytes()
	if err != nil {
		return nil, err
	}
	for i, utxo := range tx.TxIn {
		out := fetcher.FetchPrevOutput(utxo.PreviousOutPoint)
		sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, out.Value, out.PkScript, leaf, txscript.SigHashAll, wal.key)
		if err != nil {
			return nil, err
		}
		tx.TxIn[i].Witness = append(tx.TxIn[i].Witness, sig, leaf.Script, ctrBlkBytes)
	}

	// Submit tx
	if err := wal.indexer.SubmitTx(ctx, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (wal *wallet) InstantRefund(ctx context.Context, htlc *btc.HTLC, tx *wire.MsgTx) (*wire.MsgTx, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Validate tx
	if err := ValidateInstantRefundTx(htlc, tx); err != nil {
		return nil, err
	}

	// Parse the input and output from the pre-signed tx
	inputs := append(btc.UTXOs{}, btc.UTXO{
		TxID:   tx.TxIn[0].PreviousOutPoint.Hash.String(),
		Vout:   tx.TxIn[0].PreviousOutPoint.Index,
		Amount: htlc.Amount,
	})
	_, addrs, _, err := txscript.ExtractPkScriptAddrs(tx.TxOut[0].PkScript, wal.network)
	if err != nil {
		return nil, err
	}
	if len(addrs) != 1 {
		return nil, errors.New("invalid output address")
	}
	recipients := []btc.Recipient{
		{
			To:     addrs[0].EncodeAddress(),
			Amount: tx.TxOut[0].Value,
		},
	}

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, err
	}

	utxos, err := wal.indexer.GetUTXOs(ctx, wal.Address())
	if err != nil {
		return nil, err
	}
	pkScript, err := txscript.PayToAddrScript(wal.addr)
	if err != nil {
		return nil, err
	}
	log.Print("pkScript: ", hex.EncodeToString(pkScript))
	log.Print("pkScript: ", hex.EncodeToString(wal.addr.ScriptAddress()))

	fetcher, err := btc.InitFetcher(utxos, wal.addr.ScriptAddress())
	if err != nil {
		return nil, err
	}
	p2trScript, err := htlc.P2trScript()
	if err != nil {
		return nil, err
	}
	fetcher.AddPrevOut(tx.TxIn[0].PreviousOutPoint, wire.NewTxOut(htlc.Amount, p2trScript))

	sizer := btc.NewSizeEstimatorOfAddrType(utxos, wal.addrType)
	sizer.AddUtxos(inputs, btc.BaseSizeHtlcInstantRefund, btc.SegwitSizeHtlcInstantRefund)
	transaction, err := btc.BuildTransaction(wal.network, feeRate.High, inputs, utxos, sizer, recipients, wal.Address())
	if err != nil {
		return nil, err
	}

	// Sign tx
	sigHashes := txscript.NewTxSigHashes(transaction, fetcher)
	for i, input := range transaction.TxIn {
		if i == 0 {
			// Sign the instant refund leaf
			leaf, ctrBlk := htlc.InstantRefundLeaf()
			ctrBlkBytes, err := ctrBlk.ToBytes()
			if err != nil {
				return nil, err
			}

			out := fetcher.FetchPrevOutput(input.PreviousOutPoint)
			redeemerSig, err := txscript.RawTxInTapscriptSignature(transaction, sigHashes, i, out.Value, out.PkScript, leaf, txscript.SigHashAll, wal.key)
			if err != nil {
				return nil, err
			}
			initiatorSig := tx.TxIn[i].Witness[0]
			transaction.TxIn[i].Witness = append(wire.TxWitness{}, redeemerSig, initiatorSig, leaf.Script, ctrBlkBytes)
		} else {
			if err := btc.SignUtxos(wal.network, wal.addrType, transaction, i, wal.key, fetcher); err != nil {
				return nil, err
			}

		}
	}

	// Submit tx
	if err := wal.indexer.SubmitTx(ctx, transaction); err != nil {
		return nil, err
	}
	return transaction, nil
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

// NewInstantRefundTx builds a new tx to redeem an HTLC instantly using the instantRefund branch. The instant refund tx
// will contain only one input and one output, the input amount should be equal or slightly more than the output amount.
// The output will be sent to the target address. Initiator's signature will be added to the witness.
func NewInstantRefundTx(network *chaincfg.Params, key *btcec.PrivateKey, htlc *btc.HTLC, utxo btc.UTXO, recipient btc.Recipient) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(btc.DefaultTxVersion)

	// Build tx
	hash, err := chainhash.NewHashFromStr(utxo.TxID)
	if err != nil {
		return nil, err
	}
	txIn := wire.NewTxIn(wire.NewOutPoint(hash, utxo.Vout), nil, nil)
	tx.AddTxIn(txIn)

	toAddress, err := btcutil.DecodeAddress(recipient.To, network)
	if err != nil {
		return nil, err
	}
	toScript, err := txscript.PayToAddrScript(toAddress)
	if err != nil {
		return nil, err
	}
	tx.AddTxOut(wire.NewTxOut(utxo.Amount, toScript))

	// Sign the tx
	leaf, _ := htlc.InstantRefundLeaf()
	script, err := htlc.P2trScript()
	if err != nil {
		return nil, err
	}

	fetcher, err := btc.InitFetcher(btc.UTXOs{utxo}, script)
	if err != nil {
		return nil, err
	}
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
	for i, input := range tx.TxIn {
		out := fetcher.FetchPrevOutput(input.PreviousOutPoint)
		sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, out.Value, out.PkScript, leaf, btc.SigHashSingleAnyoneCanPay, key)
		if err != nil {
			return nil, err
		}

		tx.TxIn[i].Witness = append(tx.TxIn[i].Witness, sig)
	}

	return tx, nil
}

func ValidateInstantRefundTx(htlc *btc.HTLC, tx *wire.MsgTx) error {
	if len(tx.TxIn) != 1 {
		return errors.New("invalid number of inputs")
	}
	if len(tx.TxOut) != 1 {
		return errors.New("invalid number of outputs")
	}
	if len(tx.TxIn[0].Witness) != 1 {
		return errors.New("invalid witness length")
	}
	if tx.TxOut[0].Value > htlc.Amount {
		return errors.New("invalid output amount")
	}

	// Verify signature
	script, err := htlc.P2trScript()
	if err != nil {
		return err
	}
	fetcher := txscript.NewCannedPrevOutputFetcher(script, htlc.Amount)
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
	leaf, _ := htlc.InstantRefundLeaf()
	tapSigHashes, err := txscript.CalcTapscriptSignaturehash(sigHashes, btc.SigHashSingleAnyoneCanPay, tx, 0, fetcher, leaf)
	if err != nil {
		return err
	}
	log.Print("sighash before = ", hex.EncodeToString(tapSigHashes))

	sigBytes := tx.TxIn[0].Witness[0]
	if len(sigBytes) == schnorr.SignatureSize+1 {
		sigBytes = sigBytes[:len(sigBytes)-1]
	}

	signature, err := schnorr.ParseSignature(sigBytes)
	if err != nil {
		return err
	}
	pub, err := schnorr.ParsePubKey(htlc.InitiatorPubKey)
	if err != nil {
		return err
	}
	if ok := signature.Verify(tapSigHashes, pub); !ok {
		return errors.New("invalid signature")
	}
	return nil
}
