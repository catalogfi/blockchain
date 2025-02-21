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

type ExecuteOpts func(*executeOpts)

type executeOpts struct {
	rbfTxid string
}

func defaultExecuteOpts() *executeOpts {
	return &executeOpts{
		rbfTxid: "",
	}
}

func WithRbfTxid(txid string) ExecuteOpts {
	return func(o *executeOpts) {
		o.rbfTxid = txid
	}
}

type Wallet interface {
	PublicKey() *btcec.PublicKey

	Address() btcutil.Address

	InstantRefundTx(utxo UTXO, htlc *HTLC) (*wire.MsgTx, error)

	Initiate(ctx context.Context, htlc *HTLC) (*wire.MsgTx, *wire.MsgTx, error)

	Redeem(ctx context.Context, htlc *HTLC, secret []byte) (*wire.MsgTx, error)

	Refund(ctx context.Context, htlc *HTLC) (*wire.MsgTx, error)

	InstantRefund(ctx context.Context, htlc *HTLC, tx *wire.MsgTx) (*wire.MsgTx, error)

	Execute(ctx context.Context, actions []HtlcAction, opts ...ExecuteOpts) (*wire.MsgTx, error)
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
	sizer := NewSizeEstimatorOfAddrType(wal.addrType, utxos...)

	// Recipients
	htlcAddr, err := htlc.Address(wal.network)
	if err != nil {
		return nil, nil, err
	}
	recipients := []Recipient{NewRecipient(htlcAddr.EncodeAddress(), htlc.Amount)}

	// Fees
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, nil, err
	}
	feeMode := MinFeeRateMode(feeRate.High*1000, sizer)

	// Build tx
	tx, err := BuildTx(wal.network, feeMode, nil, utxos, recipients, wal.Address())
	if err != nil {
		return nil, nil, err
	}

	// Sign tx
	if err := SignTx(wal.addrType, tx, wal.internalKey, utxos); err != nil {
		return nil, nil, err
	}

	// Submit tx
	if err := wal.indexer.SubmitTx(ctx, tx); err != nil {
		return nil, nil, err
	}

	utxo := NewUtxo(tx.TxHash().String(), 0, htlc.Amount)
	irTx, err := wal.InstantRefundTx(utxo, htlc)
	if err != nil {
		return nil, nil, err
	}
	return tx, irTx, nil
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
	sizer := NewSizeEstimator(BaseSizeHtlcRedeem, SegwitSizeHtlcRedeem(len(secret)), utxos...)
	feeMode := MinFeeRateMode(feeRate.High*1000, sizer)
	tx, err := BuildTx(wal.network, feeMode, utxos, nil, nil, wal.Address())
	if err != nil {
		return nil, err
	}

	// Sign tx
	script, err := htlc.P2trScript()
	if err != nil {
		return nil, err
	}
	fetcher, err := NewFetcher(script, utxos...)
	if err != nil {
		return nil, err
	}

	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
	for i, utxo := range tx.TxIn {
		out := fetcher.FetchPrevOutput(utxo.PreviousOutPoint)
		leaf, ctrBlk := htlc.Leaf(HtlcActionRedeem)
		ctrBlkBytes, err := ctrBlk.ToBytes()
		if err != nil {
			return nil, err
		}
		sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, out.Value, out.PkScript, leaf, txscript.SigHashDefault, wal.externalKey)
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
	sizer := NewSizeEstimator(BaseSizeHtlcRefund, SegwitSizeHtlcRefund, utxos...)
	feeMode := MinFeeRateMode(feeRate.High*1000, sizer)
	tx, err := BuildTx(wal.network, feeMode, utxos, nil, nil, wal.Address())
	if err != nil {
		return nil, err
	}

	// Sign tx
	script, err := htlc.P2trScript()
	if err != nil {
		return nil, err
	}
	fetcher, err := NewFetcher(script, utxos...)
	if err != nil {
		return nil, err
	}

	// Update tx inputs sequence to htlc timelock
	for i := range tx.TxIn {
		tx.TxIn[i].Sequence = uint32(htlc.Timelock)
	}
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
	leaf, ctrBlk := htlc.Leaf(HtlcActionRefund)
	ctrBlkBytes, err := ctrBlk.ToBytes()
	if err != nil {
		return nil, err
	}
	for i, utxo := range tx.TxIn {
		out := fetcher.FetchPrevOutput(utxo.PreviousOutPoint)
		sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, out.Value, out.PkScript, leaf, txscript.SigHashDefault, wal.externalKey)
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
	fetcher, err := NewFetcher(pkScript, utxos...)
	if err != nil {
		return nil, err
	}
	p2trScript, err := htlc.P2trScript()
	if err != nil {
		return nil, err
	}
	fetcher.AddPrevOut(tx.TxIn[0].PreviousOutPoint, wire.NewTxOut(htlc.Amount, p2trScript))
	sizer := NewSizeEstimatorOfAddrType(wal.addrType, utxos...)
	sizer.AddUtxos(BaseSizeHtlcInstantRefund, SegwitSizeHtlcInstantRefund, utxo)
	feeMode := MinFeeRateMode(feeRate.High*1000, sizer)
	transaction, err := BuildTx(wal.network, feeMode, []UTXO{utxo}, utxos, []Recipient{recipient}, wal.Address())
	if err != nil {
		return nil, err
	}

	// Sign tx
	sigHashes := txscript.NewTxSigHashes(transaction, fetcher)
	for i, input := range transaction.TxIn {
		if i == 0 {
			// Sign the instant refund leaf
			leaf, ctrBlk := htlc.Leaf(HtlcActionInstantRefund)
			ctrBlkBytes, err := ctrBlk.ToBytes()
			if err != nil {
				return nil, err
			}

			out := fetcher.FetchPrevOutput(input.PreviousOutPoint)
			redeemerSig, err := txscript.RawTxInTapscriptSignature(transaction, sigHashes, i, out.Value, out.PkScript, leaf, txscript.SigHashAll, wal.externalKey)
			if err != nil {
				return nil, err
			}
			initiatorSig := tx.TxIn[i].Witness[0]
			transaction.TxIn[i].Witness = append(wire.TxWitness{}, redeemerSig, initiatorSig, leaf.Script, ctrBlkBytes)
		} else {
			if err := SignUtxos(wal.addrType, transaction, i, wal.internalKey, fetcher, sigHashes); err != nil {
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

func (wal *wallet) Execute(ctx context.Context, actions []HtlcAction, exeOpts ...ExecuteOpts) (*wire.MsgTx, error) {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Parse extra execution options
	opts := defaultExecuteOpts()
	for _, exeOpt := range exeOpts {
		exeOpt(opts)
	}

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

	// Add wallet utxos which are used in the previous tx
	var replacedTx Transaction
	var conflictUtxo *UTXO
	inputsMap, outputsMaps := map[string]bool{}, map[string]bool{} // make sure no double executions
	if opts.rbfTxid != "" {
		replacedTx, err = wal.indexer.GetTx(ctx, opts.rbfTxid)
		if err != nil {
			return nil, err
		}
		if replacedTx.Status.Confirmed {
			return nil, ErrTxNotInMempool
		}

		for _, vin := range replacedTx.VINs {
			if vin.Prevout.ScriptPubKeyAddress == wal.addr.EncodeAddress() {
				utxo := UTXO{
					TxID:   vin.TxID,
					Vout:   uint32(vin.Vout),
					Amount: int64(vin.Prevout.Value),
				}

				// Mark the first utxo from the previous tx as the conflict utxo. This conflict utxo will always present
				// in the inputs to make sure it will be conflicted with all replaced txs.
				if conflictUtxo == nil {
					conflictUtxo = &utxo
					inputsMap[vin.Prevout.ScriptPubKeyAddress] = true
					continue
				}
				utxos = append(utxos, utxo)
			}
		}
	}

	// Add all wallet utxos info to the size estimator
	sizer := NewSizeEstimatorOfAddrType(wal.addrType, utxos...)
	if conflictUtxo != nil {
		sizer = NewSizeEstimatorOfAddrType(wal.addrType, append(utxos, *conflictUtxo)...)
	}

	// Fetcher
	pkScript, err := PkScript(wal.addrType, wal.externalKey.PubKey())
	if err != nil {
		return nil, err
	}
	fetcher, err := NewFetcher(pkScript, utxos...)
	if err != nil {
		return nil, err
	}
	if conflictUtxo != nil {
		if err := AddUtxosToFetcher(fetcher, pkScript, *conflictUtxo); err != nil {
			return nil, err
		}
	}

	// Include all actions from the previous tx
	recipients := []Recipient{}
	inputs := []UTXO{}
	prevWitnesses := map[string]wire.TxWitness{}
	prevSequences := map[string]int{}
	if replacedTx.TxID != "" {
		// Inputs
		for _, vin := range replacedTx.VINs {
			if vin.Prevout.ScriptPubKeyAddress == wal.addr.EncodeAddress() {
				continue
			}

			// Store the witness and sequence from previous tx for future signing
			utxo := UTXO{
				TxID:   vin.TxID,
				Vout:   uint32(vin.Vout),
				Amount: int64(vin.Prevout.Value),
			}
			inputs = append(inputs, utxo)
			inputsMap[vin.Prevout.ScriptPubKeyAddress] = true
			witness, err := decodeWitness(*vin.Witness)
			if err != nil {
				return nil, err
			}
			prevWitnesses[utxo.String()] = witness
			prevSequences[utxo.String()] = vin.Sequence

			// Add utxo to fetcher
			script, err := hex.DecodeString(vin.Prevout.ScriptPubKey)
			if err != nil {
				return nil, err
			}
			if err := AddUtxosToFetcher(fetcher, script, utxo); err != nil {
				return nil, err
			}

			// Add utxo to the size estimator depending on the action type
			action, err := HtlcActionFromWitness(witness)
			if err != nil {
				return nil, err
			}
			switch action {
			case HtlcActionRedeem:
				sizer.AddUtxos(BaseSizeHtlcRedeem, SegwitSizeHtlcRedeem(len(witness[1])), utxo)
			case HtlcActionRefund:
				sizer.AddUtxos(BaseSizeHtlcRefund, SegwitSizeHtlcRefund, utxo)
			case HtlcActionInstantRefund:
				sizer.AddUtxos(BaseSizeHtlcInstantRefund, SegwitSizeHtlcInstantRefund, utxo)
			}
		}

		// Outputs except change utxo
		for _, vout := range replacedTx.VOUTs {
			if vout.ScriptPubKeyAddress == wal.Address().EncodeAddress() {
				continue
			}
			recipients = append(recipients, Recipient{
				To:     vout.ScriptPubKeyAddress,
				Amount: int64(vout.Value),
			})
			outputsMaps[vout.ScriptPubKeyAddress] = true
		}

		// Add the conflict
		if conflictUtxo != nil {
			inputs = append(inputs, *conflictUtxo)
			inputsMap[conflictUtxo.String()] = true
		}
	}

	// Parse the actions
	inputActions := map[string]HtlcAction{}
	for _, action := range actions {
		addr, err := action.Htlc.Address(wal.network)
		if err != nil {
			return nil, err
		}
		switch action.ActionType {
		case HtlcActionInitiate:
			if ok := outputsMaps[addr.String()]; ok {
				continue
			}
			recipients = append(recipients, NewRecipient(addr.String(), action.Htlc.Amount))
			outputsMaps[addr.String()] = true
		case HtlcActionRedeem, HtlcActionRefund:
			if ok := inputsMap[addr.EncodeAddress()]; ok {
				continue
			}
			utxo, err := action.Htlc.Utxo(ctx, wal.network, wal.indexer)
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, utxo)
			inputsMap[addr.String()] = true
			inputActions[utxo.String()] = action
			if action.ActionType == HtlcActionRedeem {
				sizer.AddUtxos(BaseSizeHtlcRedeem, SegwitSizeHtlcRedeem(len(action.Secret)), utxo)
			} else if action.ActionType == HtlcActionRefund {
				sizer.AddUtxos(BaseSizeHtlcRefund, SegwitSizeHtlcRefund, utxo)
			}

			// Add utxo to the fetcher
			fromScript, err := action.Htlc.P2trScript()
			if err != nil {
				return nil, err
			}
			if err := AddUtxosToFetcher(fetcher, fromScript, utxo); err != nil {
				return nil, err
			}
		case HtlcActionInstantRefund:
			utxo, recipient, err := ValidateInstantRefundTx(action.Htlc, action.InstantRefundTx, wal.network)
			if err != nil {
				return nil, err
			}
			if ok := inputsMap[addr.EncodeAddress()]; ok {
				continue
			}
			inputs = append([]UTXO{utxo}, inputs...)
			recipients = append([]Recipient{recipient}, recipients...)
			inputsMap[addr.String()] = true
			inputActions[utxo.String()] = action
			sizer.AddUtxos(BaseSizeHtlcInstantRefund, SegwitSizeHtlcInstantRefund, utxo)

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
	feeMode := MinFeeRateMode(feeRate*1000, sizer)
	if replacedTx.TxID != "" {
		entry, err := wal.client.GetMempoolEntry(ctx, replacedTx.TxID)
		if err != nil {
			return nil, err
		}
		prevFees, err := btcutil.NewAmount(entry.Fees.Descendant)
		if err != nil {
			return nil, err
		}

		prevFeeRate := int(prevFees*1e3) / int(entry.DescendantSize)
		feeMode = RbfMode(feeRate*1000, prevFeeRate, int64(prevFees), sizer)
	}

	// Build the tx
	tx, err := BuildTx(wal.network, feeMode, inputs, utxos, recipients, wal.Address())
	if err != nil {
		return nil, err
	}

	// Check if we have meet the fee requirement for rbf
	if replacedTx.TxID != "" {
		for {
			vsize, err := sizer.EstimateTxVirtualSize(tx)
			if err != nil {
				return nil, err
			}
			if TotalFee(tx, fetcher) > int(replacedTx.Fee)+vsize {
				break
			}
			feeRate++

			// Build tx again with new fee rate
			feeMode := MinFeeRateMode(feeRate*1000, sizer)
			tx, err = BuildTx(wal.network, feeMode, inputs, utxos, recipients, wal.Address())
			if err != nil {
				return nil, err
			}
		}
	}

	// Set sequence number for refund inputs
	for i := range tx.TxIn {
		action, ok := inputActions[tx.TxIn[i].PreviousOutPoint.String()]
		if ok && action.ActionType == HtlcActionRefund {
			tx.TxIn[i].Sequence = uint32(action.Htlc.Timelock)

		}
		if sequence, ok := prevSequences[tx.TxIn[i].PreviousOutPoint.String()]; ok {
			tx.TxIn[i].Sequence = uint32(sequence)
		}
	}

	// Sign the tx
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
	for i, input := range tx.TxIn {
		outpoint := fetcher.FetchPrevOutput(input.PreviousOutPoint)

		// If the utxo is from the previous tx, we only need to sign our signature again and reuse the rest parts of the
		// witness.
		witness, ok := prevWitnesses[input.PreviousOutPoint.String()]
		if ok {
			// Second last witness should be the script and the first witness is our signature.
			// We have validated the witness before storing the map, so we can confidently use it.
			leaf := txscript.NewTapLeaf(txscript.BaseLeafVersion, witness[len(witness)-2])
			sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, outpoint.Value, outpoint.PkScript, leaf, txscript.SigHashDefault, wal.externalKey)
			if err != nil {
				return nil, err
			}
			witness[0] = sig
			tx.TxIn[i].Witness = witness
			continue
		}

		// Rest utxos should be the wallet utxos or from new actions.
		action, ok := inputActions[tx.TxIn[i].PreviousOutPoint.String()]
		if !ok {
			if err := SignUtxos(wal.addrType, tx, i, wal.internalKey, fetcher, sigHashes); err != nil {
				return nil, err
			}
			continue
		}

		// Sign and build the witness
		leaf, ctrBlk := action.Htlc.Leaf(action.ActionType)
		ctrBlkBytes, err := ctrBlk.ToBytes()
		if err != nil {
			return nil, err
		}
		sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, outpoint.Value, outpoint.PkScript, leaf, txscript.SigHashDefault, wal.externalKey)
		if err != nil {
			return nil, err
		}
		switch action.ActionType {
		case HtlcActionRedeem:
			tx.TxIn[i].Witness = append(tx.TxIn[i].Witness, sig, action.Secret, leaf.Script, ctrBlkBytes)
		case HtlcActionRefund:
			tx.TxIn[i].Witness = append(tx.TxIn[i].Witness, sig, leaf.Script, ctrBlkBytes)
		case HtlcActionInstantRefund:
			initiatorSig := action.InstantRefundTx.TxIn[0].Witness[0]
			tx.TxIn[i].Witness = append(wire.TxWitness{}, sig, initiatorSig, leaf.Script, ctrBlkBytes)
		default:
			return nil, fmt.Errorf("unknown action type: %v", action.ActionType)
		}
	}

	// Submit tx
	return tx, wal.indexer.SubmitTx(ctx, tx)
}

func (wal *wallet) InstantRefundTx(utxo UTXO, htlc *HTLC) (*wire.MsgTx, error) {
	recipient := NewRecipient(wal.addr.EncodeAddress(), htlc.Amount)
	irTx, err := BuildTx(wal.network, GaslessMode(), []UTXO{utxo}, nil, []Recipient{recipient}, nil)
	if err != nil {
		return nil, err
	}

	// Sign the tx
	leaf, _ := htlc.Leaf(HtlcActionInstantRefund)
	script, err := htlc.P2trScript()
	if err != nil {
		return nil, err
	}
	fetcher, err := NewFetcher(script, utxo)
	if err != nil {
		return nil, err
	}
	sigHashes := txscript.NewTxSigHashes(irTx, fetcher)
	for i, input := range irTx.TxIn {
		out := fetcher.FetchPrevOutput(input.PreviousOutPoint)
		sig, err := txscript.RawTxInTapscriptSignature(irTx, sigHashes, i, out.Value, out.PkScript, leaf, SigHashSingleAnyoneCanPay, wal.externalKey)
		if err != nil {
			return nil, err
		}

		irTx.TxIn[i].Witness = append(irTx.TxIn[i].Witness, sig)
	}
	return irTx, nil
}

func decodeWitness(witnessStr []string) (wire.TxWitness, error) {
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
