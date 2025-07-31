package btc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
)

type Wallet interface {
	PublicKey() *btcec.PublicKey

	Address() btcutil.Address

	InstantRefundTx(utxo UTXO, htlc *HTLC) (*wire.MsgTx, error)

	Initiate(ctx context.Context, htlc *HTLC) (*wire.MsgTx, *wire.MsgTx, error)

	Redeem(ctx context.Context, htlc *HTLC) (*wire.MsgTx, error)

	Refund(ctx context.Context, htlc *HTLC, target btcutil.Address) (*wire.MsgTx, error)

	InstantRefund(ctx context.Context, htlc *HTLC, tx *wire.MsgTx) (*wire.MsgTx, error)

	Execute(ctx context.Context, actions []HtlcAction, prevTxid string) (*wire.MsgTx, error)
}

type wallet struct {
	mu           *sync.Mutex
	network      *chaincfg.Params
	internalKey  *btcec.PrivateKey
	externalKey  *btcec.PrivateKey
	addrType     waddrmgr.AddressType
	addr         btcutil.Address
	indexer      IndexerClient
	client       Client
	feeEstimator FeeEstimator
}

func NewWallet(network *chaincfg.Params, addrType waddrmgr.AddressType, privateKey *btcec.PrivateKey, indexer IndexerClient, client Client, estimator FeeEstimator) (Wallet, error) {
	externalKey := privateKey
	if addrType == waddrmgr.TaprootPubKey {
		externalKey = txscript.TweakTaprootPrivKey(*privateKey, nil)
	}
	addr, err := PublicKeyAddress(network, addrType, externalKey.PubKey())
	if err != nil {
		return nil, err
	}
	return &wallet{
		mu:           new(sync.Mutex),
		network:      network,
		internalKey:  privateKey,
		externalKey:  externalKey,
		addrType:     addrType,
		addr:         addr,
		indexer:      indexer,
		client:       client,
		feeEstimator: estimator,
	}, nil
}

func (wal *wallet) PublicKey() *btcec.PublicKey {
	return wal.externalKey.PubKey()
}

func (wal *wallet) Address() btcutil.Address {
	return wal.addr
}

func (wal *wallet) Initiate(ctx context.Context, htlc *HTLC) (*wire.MsgTx, *wire.MsgTx, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Inputs
	utxos, err := wal.indexer.GetUTXOs(ctx, wal.addr)
	if err != nil {
		return nil, nil, err
	}
	signers, _, err := NewSignersByAddrType(wal.addrType, wal.internalKey, utxos)
	if err != nil {
		return nil, nil, err
	}

	// Recipients
	htlcAddr := htlc.MustAddress(wal.network)
	recipient, err := NewTxOutFromAddress(htlcAddr, htlc.Amount)
	if err != nil {
		return nil, nil, err
	}
	recipients := []*wire.TxOut{recipient}

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, nil, err
	}
	feeMode := MinFeeRateMode(feeRate.High, signers)

	// Build tx
	tx, err := BuildTx(feeMode, nil, utxos, recipients, wal.Address())
	if err != nil {
		return nil, nil, err
	}

	// Sign tx
	if err := signers.Sign(tx); err != nil {
		return nil, nil, err
	}

	// Submit tx
	if err := wal.indexer.SubmitTx(ctx, tx); err != nil {
		return nil, nil, err
	}
	utxo := NewUtxo(tx.TxHash().String(), 0, htlc.Amount, recipient.PkScript)
	irTx, err := wal.InstantRefundTx(utxo, htlc)
	if err != nil {
		return nil, nil, err
	}
	return tx, irTx, nil
}

func (wal *wallet) Redeem(ctx context.Context, htlc *HTLC) (*wire.MsgTx, error) {
	if len(htlc.Secret()) == 0 {
		return nil, fmt.Errorf("nil secret")
	}

	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Make sure the HTLC is redeemable
	addr := htlc.MustAddress(wal.network)
	utxos, err := wal.indexer.GetUTXOs(ctx, addr)
	if err != nil {
		return nil, err
	}
	// todo : disable this for now
	// redeemable, _, err := htlc.Redeemable(utxos)
	// if err != nil {
	// 	return nil, err
	// }
	// if !redeemable {
	// 	return nil, fmt.Errorf("HTLC is not redeemable")
	// }

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, err
	}

	// Build tx
	leaf, ctlBlock := htlc.Leaf(HtlcActionRedeem)
	signers, err := NewSigners(NewHtlcRedeemSigner(wal.externalKey, leaf, ctlBlock, htlc.Secret()), utxos...)
	if err != nil {
		return nil, err
	}
	feeMode := MinFeeRateMode(feeRate.High, signers)
	tx, err := BuildTx(feeMode, utxos, nil, nil, wal.Address())
	if err != nil {
		return nil, err
	}

	// Sign tx
	if err := signers.Sign(tx); err != nil {
		return nil, err
	}

	// Submit tx
	if err := wal.indexer.SubmitTx(ctx, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (wal *wallet) Refund(ctx context.Context, htlc *HTLC, target btcutil.Address) (*wire.MsgTx, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Find all utxos which are confirmed and refundable
	utxos, err := htlc.RefundableUtxos(ctx, wal.network, wal.indexer)
	if err != nil {
		return nil, err
	}
	if len(utxos) == 0 {
		return nil, fmt.Errorf("no utxo to refund")
	}

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, err
	}

	// Build tx
	leaf, ctrBlk := htlc.Leaf(HtlcActionRefund)
	signers, err := NewSigners(NewHtlcRefundSigner(wal.externalKey, leaf, ctrBlk, htlc.Timelock), utxos...)
	if err != nil {
		return nil, err
	}
	feeMode := MinFeeRateMode(feeRate.High, signers)
	if target == nil {
		target = wal.addr
	}
	tx, err := BuildTx(feeMode, utxos, nil, nil, target)
	if err != nil {
		return nil, err
	}

	// Update tx inputs sequence to htlc timelock
	for i := range tx.TxIn {
		tx.TxIn[i].Sequence = uint32(htlc.Timelock)
	}

	// Sign tx
	if err := signers.Sign(tx); err != nil {
		return nil, err
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

	// Inputs
	utxos, err := wal.indexer.GetUTXOs(ctx, wal.Address())
	if err != nil {
		return nil, err
	}

	// Signers
	signers, _, err := NewSignersByAddrType(wal.addrType, wal.internalKey, utxos)
	if err != nil {
		return nil, err
	}
	leaf, ctrBlk := htlc.Leaf(HtlcActionInstantRefund)
	signer := NewHtlcInstantRefundSigner(wal.externalKey, leaf, ctrBlk, true, tx.TxIn[0].Witness[0])
	signers.AddUtxo(signer, utxo)

	// Fee
	feeMode := MinFeeRateMode(feeRate.High, signers)
	transaction, err := BuildTx(feeMode, []UTXO{utxo}, utxos, []*wire.TxOut{recipient}, wal.Address())
	if err != nil {
		return nil, err
	}

	// Sign tx
	if err := signers.Sign(transaction); err != nil {
		return nil, err
	}

	// Submit tx
	if err := wal.indexer.SubmitTx(ctx, transaction); err != nil {
		return nil, err
	}
	return transaction, nil
}

func (wal *wallet) Execute(ctx context.Context, actions []HtlcAction, prevTxid string) (*wire.MsgTx, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Fetch wallet utxos which are confirmed
	walUtxos, err := wal.indexer.GetUTXOs(ctx, wal.addr)
	if err != nil {
		return nil, err
	}
	utxos := make([]UTXO, 0, len(walUtxos))
	for _, utxo := range walUtxos {
		if utxo.Status != nil && utxo.Status.Confirmed {
			utxos = append(utxos, utxo)
		}
	}

	// Init a signer
	signers, walSigner, err := NewSignersByAddrType(wal.addrType, wal.internalKey, utxos)
	if err != nil {
		return nil, err
	}

	// Parse inputs and output from prevTx if exists
	inputs, outputs := []UTXO{}, []*wire.TxOut{}
	sequenceMap := map[string]int{}

	var replacedTx Transaction
	inputsMap, outputsMaps := map[string]bool{}, map[string]bool{} // make sure no double executions
	if prevTxid != "" {
		replacedTx, err = wal.indexer.GetTx(ctx, prevTxid)
		if err != nil {
			return nil, err
		}
		if replacedTx.Status.Confirmed {
			return nil, ErrTxNotInMempool
		}

		for _, vin := range replacedTx.VINs {
			pkScript, err := hex.DecodeString(vin.Prevout.ScriptPubKey)
			if err != nil {
				return nil, err
			}
			utxo := UTXO{
				TxID:     vin.TxID,
				Vout:     uint32(vin.Vout),
				Amount:   int64(vin.Prevout.Value),
				PkScript: pkScript,
			}
			inputs = append(inputs, utxo)
			inputsMap[utxo.String()] = true

			if vin.Prevout.ScriptPubKeyAddress == wal.addr.EncodeAddress() {
				signers.AddUtxo(walSigner, utxo)
			} else {
				// Decode the action from its witness
				witness, err := DecodeWitness(*vin.Witness)
				if err != nil {
					return nil, err
				}
				action, err := HtlcActionFromWitness(witness)
				if err != nil {
					return nil, err
				}

				leaf := txscript.NewTapLeaf(txscript.BaseLeafVersion, witness[len(witness)-2])
				ctrBlock, err := txscript.ParseControlBlock(witness[len(witness)-1])
				if err != nil {
					return nil, err
				}

				var signer Signer
				switch action {
				case HtlcActionRedeem:
					secret := witness[1]
					signer = NewHtlcRedeemSigner(wal.externalKey, leaf, *ctrBlock, secret)
					signers.AddUtxo(signer, utxo)
				case HtlcActionRefund:
					timelock := int64(vin.Sequence)
					signer = NewHtlcRefundSigner(wal.externalKey, leaf, *ctrBlock, timelock)
					signers.AddUtxo(signer, utxo)
					sequenceMap[utxo.String()] = vin.Sequence
				case HtlcActionInstantRefund:
					otherSig := witness[1]
					signer = NewHtlcInstantRefundSigner(wal.externalKey, leaf, *ctrBlock, true, otherSig)
					signers.AddUtxo(signer, utxo)
				}
			}
		}

		for _, vout := range replacedTx.VOUTs {
			// todo : restrict the change to be the last output?
			if vout.ScriptPubKeyAddress == wal.Address().EncodeAddress() {
				continue
			}

			pkScript, err := hex.DecodeString(vout.ScriptPubKey)
			if err != nil {
				return nil, err
			}
			outputs = append(outputs, &wire.TxOut{
				Value:    int64(vout.Value),
				PkScript: pkScript,
			})

			outputsMaps[vout.ScriptPubKeyAddress] = true
		}
	}

	for _, action := range actions {
		addr := action.Htlc.MustAddress(wal.network)
		switch action.ActionType {
		case HtlcActionInitiate:
			if ok := outputsMaps[addr.String()]; ok {
				continue
			}
			recipient, err := NewTxOutFromAddress(addr, action.Htlc.Amount)
			if err != nil {
				return nil, err
			}
			outputs = append(outputs, recipient)
			outputsMaps[addr.String()] = true
		case HtlcActionRedeem:
			if ok := inputsMap[addr.EncodeAddress()]; ok {
				continue
			}
			utxo, err := action.Htlc.Utxo(ctx, wal.network, wal.indexer)
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, utxo)
			leaf, ctrBlk := action.Htlc.Leaf(HtlcActionRedeem)
			signer := NewHtlcRedeemSigner(wal.externalKey, leaf, ctrBlk, action.Htlc.Secret())
			signers.AddUtxo(signer, utxo)
			inputsMap[addr.String()] = true
		case HtlcActionRefund:
			if ok := inputsMap[addr.EncodeAddress()]; ok {
				continue
			}
			utxo, err := action.Htlc.Utxo(ctx, wal.network, wal.indexer)
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, utxo)
			leaf, ctrBlk := action.Htlc.Leaf(HtlcActionRefund)
			signer := NewHtlcRefundSigner(wal.externalKey, leaf, ctrBlk, action.Htlc.Timelock)
			signers.AddUtxo(signer, utxo)
			inputsMap[addr.String()] = true

			sequenceMap[utxo.String()] = int(action.Htlc.Timelock)
			if action.RefundTo != nil {
				recipient, err := NewTxOutFromAddress(action.RefundTo, action.Htlc.Amount)
				if err != nil {
					return nil, err
				}
				outputs = append(outputs, recipient)
			}
		case HtlcActionInstantRefund:
			if ok := inputsMap[addr.EncodeAddress()]; ok {
				continue
			}
			utxo, recipient, err := ValidateInstantRefundTx(action.Htlc, action.InstantRefundTx, wal.network)
			if err != nil {
				return nil, err
			}
			if ok := inputsMap[addr.EncodeAddress()]; ok {
				continue
			}
			inputs = append([]UTXO{utxo}, inputs...)
			outputs = append([]*wire.TxOut{recipient}, outputs...)
			inputsMap[addr.String()] = true

			leaf, ctrBlk := action.Htlc.Leaf(HtlcActionInstantRefund)
			signer := NewHtlcInstantRefundSigner(wal.externalKey, leaf, ctrBlk, true, action.InstantRefundTx.TxIn[0].Witness[0])
			signers.AddUtxo(signer, utxo)
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
	feeMode := MinFeeRateMode(feeRate, signers)
	if replacedTx.TxID != "" {
		entry, err := wal.client.GetMempoolEntry(ctx, replacedTx.TxID)
		if err != nil {
			return nil, err
		}
		prevFees, err := btcutil.NewAmount(entry.Fees.Descendant)
		if err != nil {
			return nil, err
		}

		prevFeeRate := NewSatoshiPerKb(int64(prevFees), int(entry.DescendantSize))
		feeMode = RbfMode(feeRate, prevFeeRate, int64(prevFees), signers)
	}

	// Build the tx
	tx, err := BuildTx(feeMode, inputs, utxos, outputs, wal.Address())
	if err != nil {
		return nil, err
	}

	// Set sequence number for refund inputs
	for i := range tx.TxIn {
		seq, ok := sequenceMap[tx.TxIn[i].PreviousOutPoint.String()]
		if ok {
			tx.TxIn[i].Sequence = uint32(seq)
		}
	}

	// Sign the tx
	if err := signers.Sign(tx); err != nil {
		return nil, err
	}

	// Submit tx
	return tx, wal.indexer.SubmitTx(ctx, tx)
}

func (wal *wallet) InstantRefundTx(utxo UTXO, htlc *HTLC) (*wire.MsgTx, error) {
	recipient, err := NewTxOutFromAddress(wal.addr, utxo.Amount)
	if err != nil {
		return nil, err
	}
	irTx, err := BuildTx(GaslessMode(), []UTXO{utxo}, nil, []*wire.TxOut{recipient}, nil)
	if err != nil {
		return nil, err
	}
	leaf, ctrBlk := htlc.Leaf(HtlcActionInstantRefund)
	signers, err := NewSigners(NewHtlcInstantRefundSigner(wal.externalKey, leaf, ctrBlk, false, nil, WithSighashType(SigHashSingleAnyoneCanPay)), utxo)
	if err != nil {
		return nil, err
	}

	if err := signers.Sign(irTx); err != nil {
		return nil, err
	}
	irTx.TxIn[0].Witness = wire.TxWitness{irTx.TxIn[0].Witness[1]}
	return irTx, nil
}

func DecodeWitness(witnessStr []string) (wire.TxWitness, error) {
	witness := make(wire.TxWitness, len(witnessStr))
	for i := range witness {
		decoded, err := hex.DecodeString(witnessStr[i])
		if err != nil {
			return nil, err
		}
		witness[i] = decoded
	}
	return witness, nil
}
