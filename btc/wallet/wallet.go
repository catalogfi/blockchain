package wallet

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
	"github.com/catalogfi/blockchain/btc"
)

type Wallet interface {
	Address() btcutil.Address

	Initiate(ctx context.Context, htlc *btc.HTLC) (*wire.MsgTx, error)

	Redeem(ctx context.Context, htlc *btc.HTLC, secret []byte) (*wire.MsgTx, error)

	Refund(ctx context.Context, htlc *btc.HTLC) (*wire.MsgTx, error)

	InstantRefund(ctx context.Context, htlc *btc.HTLC, tx *wire.MsgTx) (*wire.MsgTx, error)

	Execute(ctx context.Context, actions []btc.HtlcAction, conflictUtxo *btc.UTXO) (*wire.MsgTx, error)
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
	if err := btc.SignTx(wal.addrType, tx, wal.key, utxos); err != nil {
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
	sizeEstimator := btc.NewSizeEstimator(utxos, btc.BaseSizeHtlcRedeem, btc.SegwitSizeHtlcRedeem(len(secret)))
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
	sizeEstimator := btc.NewSizeEstimator(utxos, btc.BaseSizeHtlcRefund, btc.SegwitSizeHtlcRefund)
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
	if err := btc.ValidateInstantRefundTx(htlc, tx); err != nil {
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
	recipients := btc.SingleRecipient(addrs[0].EncodeAddress(), tx.TxOut[0].Value)

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

	fetcher, err := btc.InitFetcher(utxos, pkScript)
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
			if err := btc.SignUtxos(wal.addrType, transaction, i, wal.key, fetcher, sigHashes); err != nil {
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

func (wal *wallet) Execute(ctx context.Context, actions []btc.HtlcAction, conflictUtxo *btc.UTXO) (*wire.MsgTx, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Fetch wallet utxos without the conflict utxo and unconfirmed utxo
	rawUtxos, err := wal.indexer.GetUTXOs(ctx, wal.addr)
	if err != nil {
		return nil, err
	}
	utxos := make([]btc.UTXO, 0, len(rawUtxos))
	for _, utxo := range utxos {
		if conflictUtxo != nil && utxo.String() == conflictUtxo.String() {
			continue
		}
		if utxo.Status != nil && !utxo.Status.Confirmed {
			continue
		}
		utxos = append(utxos, utxo)
	}
	sizer := btc.NewSizeEstimatorOfAddrType(utxos, wal.addrType)

	// Fetcher
	pkScript, err := btc.PkScript(wal.addrType, wal.key.PubKey())
	if err != nil {
		return nil, err
	}
	if wal.addrType == waddrmgr.TaprootPubKey {
		tapkey := txscript.ComputeTaprootOutputKey(wal.key.PubKey(), nil)
		pkScript, err = btc.PkScript(wal.addrType, tapkey)
		if err != nil {
			return nil, err
		}
	}
	fetcher, err := btc.InitFetcher(utxos, pkScript)
	if err != nil {
		return nil, err
	}

	// Append the conflict utxo to make sure the replacement txs will be conflicted with each other
	recipients := []btc.Recipient{}
	inputs := []btc.UTXO{}
	if conflictUtxo != nil {
		inputs = append(inputs, *conflictUtxo)
	}

	// Parse the actions
	inputActions := map[int]btc.HtlcAction{}
	for _, action := range actions {
		switch action.ActionType {
		case btc.HtlcActionInitiate:
			addr, err := action.Htlc.Address(wal.network)
			if err != nil {
				return nil, err
			}
			recipients = append(recipients, btc.Recipient{
				To:     addr.String(),
				Amount: action.Htlc.Amount,
			})
		case btc.HtlcActionRedeem, btc.HtlcActionRefund, btc.HtlcActionInstantRefund:
			utxo, err := action.Htlc.Utxo(ctx, wal.network, wal.indexer)
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, utxo)
			inputActions[len(inputs)-1] = action
			if action.ActionType == btc.HtlcActionRedeem {
				sizer.AddUtxos([]btc.UTXO{utxo}, btc.BaseSizeHtlcRedeem, btc.SegwitSizeHtlcRedeem(len(action.Secret)))
			} else if action.ActionType == btc.HtlcActionRefund {
				sizer.AddUtxos([]btc.UTXO{utxo}, btc.BaseSizeHtlcRefund, btc.SegwitSizeHtlcRefund)
			} else if action.ActionType == btc.HtlcActionInstantRefund {
				sizer.AddUtxos([]btc.UTXO{utxo}, btc.BaseSizeHtlcInstantRefund, btc.SegwitSizeHtlcInstantRefund)
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
		default:
			return nil, errors.New("invalid action type")
		}
	}

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, err
	}

	// Build the tx
	tx, err := btc.BuildRbfTransaction(wal.network, feeRate.High, inputs, utxos, sizer, recipients, wal.Address())
	if err != nil {
		return nil, err
	}

	// Set sequence number for refund inputs
	for i := range tx.TxIn {
		action, ok := inputActions[i]
		if ok && action.ActionType == btc.HtlcActionRefund {
			tx.TxIn[i].Sequence = uint32(action.Htlc.Timelock)
		}
	}

	// Sign the tx
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
	for i, input := range tx.TxIn {
		action, ok := inputActions[i]
		if !ok {
			if err := btc.SignUtxos(wal.addrType, tx, i, wal.key, fetcher, sigHashes); err != nil {
				return nil, err
			}
		}

		outpoint := fetcher.FetchPrevOutput(input.PreviousOutPoint)
		switch action.ActionType {
		case btc.HtlcActionRedeem:
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
		case btc.HtlcActionRefund:
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
		case btc.HtlcActionInstantRefund:
			// todo : this is wrong, we need to add the instant refund tx to the corresponding index
			leaf, ctrBlk := action.Htlc.InstantRefundLeaf()
			ctrBlkBytes, err := ctrBlk.ToBytes()
			if err != nil {
				return nil, err
			}
			redeemerSig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, outpoint.Value, outpoint.PkScript, leaf, txscript.SigHashAll, wal.key)
			if err != nil {
				return nil, err
			}
			initiatorSig := tx.TxIn[i].Witness[0]
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
