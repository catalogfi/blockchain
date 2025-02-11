package btc

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
)

type Wallet interface {
	Address() btcutil.Address

	Initiate(ctx context.Context, htlc *HTLC) (*wire.MsgTx, error)

	Redeem(ctx context.Context, htlc *HTLC, secret []byte) (*wire.MsgTx, error)

	Refund(ctx context.Context, htlc *HTLC) (*wire.MsgTx, error)

	InstantRefund(ctx context.Context, htlc *HTLC, tx *wire.MsgTx) (*wire.MsgTx, error)

	Execute(ctx context.Context, actions []HtlcAction) (*wire.MsgTx, error)
}

type wallet struct {
	mu           *sync.Mutex
	network      *chaincfg.Params
	key          *btcec.PrivateKey
	addrType     waddrmgr.AddressType
	addr         btcutil.Address
	indexer      IndexerClient
	feeEstimator FeeEstimator
}

func NewWallet(network *chaincfg.Params, addrType waddrmgr.AddressType, key *btcec.PrivateKey, indexer IndexerClient, estimator FeeEstimator) (Wallet, error) {
	addr, err := PublicKeyAddress(network, addrType, key.PubKey())
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

func (wal *wallet) Initiate(ctx context.Context, htlc *HTLC) (*wire.MsgTx, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Inputs
	utxos, err := wal.indexer.GetUTXOs(ctx, wal.addr)
	if err != nil {
		return nil, err
	}
	sizer := NewSizeEstimatorOfAddrType(utxos, wal.addrType)

	// Recipients
	htlcAddr, err := htlc.Address(wal.network)
	if err != nil {
		return nil, err
	}
	recipients := SingleRecipient(htlcAddr.EncodeAddress(), htlc.Amount)

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, err
	}

	// Build tx
	tx, err := BuildTransaction(wal.network, feeRate.High, nil, utxos, sizer, recipients, wal.Address())
	if err != nil {
		return nil, err
	}

	// Sign tx
	if err := SignTx(wal.addrType, tx, wal.key, utxos); err != nil {
		return nil, err
	}

	// Submit tx
	if err := wal.indexer.SubmitTx(ctx, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (wal *wallet) Redeem(ctx context.Context, htlc *HTLC, secret []byte) (*wire.MsgTx, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Make sure the HTLC is redeemable
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
	sizeEstimator := NewSizeEstimator(utxos, BaseSizeHtlcRedeem, SegwitSizeHtlcRedeem(len(secret)))
	tx, err := BuildTransaction(wal.network, feeRate.High, utxos, nil, sizeEstimator, nil, wal.Address())
	if err != nil {
		return nil, err
	}

	// Sign tx
	script, err := htlc.P2trScript()
	if err != nil {
		return nil, err
	}
	fetcher, err := InitFetcher(utxos, script)
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

func (wal *wallet) Refund(ctx context.Context, htlc *HTLC) (*wire.MsgTx, error) {
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
	sizeEstimator := NewSizeEstimator(utxos, BaseSizeHtlcRefund, SegwitSizeHtlcRefund)
	tx, err := BuildTransaction(wal.network, feeRate.High, utxos, nil, sizeEstimator, nil, wal.Address())
	if err != nil {
		return nil, err
	}

	// Sign tx
	script, err := htlc.P2trScript()
	if err != nil {
		return nil, err
	}
	fetcher, err := InitFetcher(utxos, script)
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

func (wal *wallet) InstantRefund(ctx context.Context, htlc *HTLC, tx *wire.MsgTx) (*wire.MsgTx, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Validate tx
	utxo, recipient, err := ValidateInstantRefundTx(htlc, tx, wal.network)
	if err != nil {
		return nil, err
	}

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, err
	}

	// Fetcher
	utxos, err := wal.indexer.GetUTXOs(ctx, wal.Address())
	if err != nil {
		return nil, err
	}
	pkScript, err := txscript.PayToAddrScript(wal.addr)
	if err != nil {
		return nil, err
	}
	fetcher, err := InitFetcher(utxos, pkScript)
	if err != nil {
		return nil, err
	}
	p2trScript, err := htlc.P2trScript()
	if err != nil {
		return nil, err
	}
	fetcher.AddPrevOut(tx.TxIn[0].PreviousOutPoint, wire.NewTxOut(htlc.Amount, p2trScript))

	sizer := NewSizeEstimatorOfAddrType(utxos, wal.addrType)
	sizer.AddUtxos([]UTXO{utxo}, BaseSizeHtlcInstantRefund, SegwitSizeHtlcInstantRefund)
	transaction, err := BuildTransaction(wal.network, feeRate.High, []UTXO{utxo}, utxos, sizer, []Recipient{recipient}, wal.Address())
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
			if err := SignUtxos(wal.addrType, transaction, i, wal.key, fetcher, sigHashes); err != nil {
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

func (wal *wallet) Execute(ctx context.Context, actions []HtlcAction) (*wire.MsgTx, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Fetch wallet utxos without the conflict utxo and unconfirmed utxo
	rawUtxos, err := wal.indexer.GetUTXOs(ctx, wal.addr)
	if err != nil {
		return nil, err
	}
	utxos := make([]UTXO, 0, len(rawUtxos))
	for _, utxo := range rawUtxos {
		if utxo.Status != nil && !utxo.Status.Confirmed {
			continue
		}
		utxos = append(utxos, utxo)
	}
	sizer := NewSizeEstimatorOfAddrType(utxos, wal.addrType)

	// Fetcher
	pkScript, err := PkScript(wal.addrType, wal.key.PubKey())
	if err != nil {
		return nil, err
	}
	if wal.addrType == waddrmgr.TaprootPubKey {
		tapkey := txscript.ComputeTaprootOutputKey(wal.key.PubKey(), nil)
		pkScript, err = PkScript(wal.addrType, tapkey)
		if err != nil {
			return nil, err
		}
	}
	fetcher, err := InitFetcher(utxos, pkScript)
	if err != nil {
		return nil, err
	}

	// Parse the actions
	recipients := []Recipient{}
	inputs := []UTXO{}
	inputActions := map[string]HtlcAction{}
	for _, action := range actions {
		switch action.ActionType {
		case HtlcActionInitiate:
			addr, err := action.Htlc.Address(wal.network)
			if err != nil {
				return nil, err
			}
			recipients = append(recipients, Recipient{
				To:     addr.String(),
				Amount: action.Htlc.Amount,
			})
		case HtlcActionRedeem, HtlcActionRefund:
			utxo, err := action.Htlc.Utxo(ctx, wal.network, wal.indexer)
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, utxo)
			inputActions[utxo.String()] = action
			if action.ActionType == HtlcActionRedeem {
				sizer.AddUtxos([]UTXO{utxo}, BaseSizeHtlcRedeem, SegwitSizeHtlcRedeem(len(action.Secret)))
			} else if action.ActionType == HtlcActionRefund {
				sizer.AddUtxos([]UTXO{utxo}, BaseSizeHtlcRefund, SegwitSizeHtlcRefund)
			}

			// Add to the fetcher
			hash, err := chainhash.NewHashFromStr(utxo.TxID)
			if err != nil {
				return nil, err
			}
			fromScript, err := action.Htlc.P2trScript()
			if err != nil {
				return nil, err
			}
			fetcher.AddPrevOut(wire.OutPoint{
				Hash:  *hash,
				Index: utxo.Vout,
			}, wire.NewTxOut(utxo.Amount, fromScript))
		case HtlcActionInstantRefund:
			utxo, recipient, err := ValidateInstantRefundTx(action.Htlc, action.InstantRefundTx, wal.network)
			if err != nil {
				return nil, err
			}
			inputs = append([]UTXO{utxo}, inputs...)
			recipients = append([]Recipient{recipient}, recipients...)
			inputActions[utxo.String()] = action
			sizer.AddUtxos([]UTXO{utxo}, BaseSizeHtlcInstantRefund, SegwitSizeHtlcInstantRefund)

			// Add to the fetcher
			fromScript, err := action.Htlc.P2trScript()
			if err != nil {
				return nil, err
			}
			fetcher.AddPrevOut(action.InstantRefundTx.TxIn[0].PreviousOutPoint, wire.NewTxOut(utxo.Amount, fromScript))
		default:
			return nil, errors.New("invalid action type")
		}
	}

	// Fees
	feeRates, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, err
	}
	feeRate := feeRates.High

	// Build the tx
	tx, err := BuildTransaction(wal.network, feeRate, inputs, utxos, sizer, recipients, wal.Address())
	if err != nil {
		return nil, err
	}

	// Set sequence number for refund inputs
	for i := range tx.TxIn {
		action, ok := inputActions[tx.TxIn[i].PreviousOutPoint.String()]
		if ok && action.ActionType == HtlcActionRefund {
			tx.TxIn[i].Sequence = uint32(action.Htlc.Timelock)
		}
	}

	// Sign the tx
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
	for i, input := range tx.TxIn {
		action, ok := inputActions[tx.TxIn[i].PreviousOutPoint.String()]
		if !ok {
			if err := SignUtxos(wal.addrType, tx, i, wal.key, fetcher, sigHashes); err != nil {
				return nil, err
			}
			continue
		}

		outpoint := fetcher.FetchPrevOutput(input.PreviousOutPoint)
		switch action.ActionType {
		case HtlcActionRedeem:
			leaf, ctrBlk := action.Htlc.RedeemLeaf()
			ctrBlkBytes, err := ctrBlk.ToBytes()
			if err != nil {
				return nil, err
			}
			sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, outpoint.Value, outpoint.PkScript, leaf, txscript.SigHashAll, wal.key)
			if err != nil {
				return nil, err
			}
			tx.TxIn[i].Witness = append(tx.TxIn[i].Witness, sig, action.Secret, leaf.Script, ctrBlkBytes)
		case HtlcActionRefund:
			leaf, ctrBlk := action.Htlc.RefundLeaf()
			ctrBlkBytes, err := ctrBlk.ToBytes()
			if err != nil {
				return nil, err
			}
			sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, outpoint.Value, outpoint.PkScript, leaf, txscript.SigHashAll, wal.key)
			if err != nil {
				return nil, err
			}
			tx.TxIn[i].Witness = append(tx.TxIn[i].Witness, sig, leaf.Script, ctrBlkBytes)
		case HtlcActionInstantRefund:
			leaf, ctrBlk := action.Htlc.InstantRefundLeaf()
			ctrBlkBytes, err := ctrBlk.ToBytes()
			if err != nil {
				return nil, err
			}
			redeemerSig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, outpoint.Value, outpoint.PkScript, leaf, txscript.SigHashAll, wal.key)
			if err != nil {
				return nil, err
			}
			initiatorSig := action.InstantRefundTx.TxIn[0].Witness[0]
			tx.TxIn[i].Witness = append(wire.TxWitness{}, redeemerSig, initiatorSig, leaf.Script, ctrBlkBytes)
		default:
			return nil, fmt.Errorf("unknown action type: %v", action.ActionType)
		}
	}

	// Submit tx
	if err := wal.indexer.SubmitTx(ctx, tx); err != nil {
		return nil, err
	}
	return tx, nil
}
