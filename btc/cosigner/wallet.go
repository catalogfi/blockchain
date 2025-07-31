package cosigner

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/catalogfi/blockchain/btc"
)

type Wallet interface {
	PublicKey() *btcec.PublicKey

	Address() btcutil.Address

	Execute(ctx context.Context, actions []btc.HtlcAction) (*wire.MsgTx, *wire.MsgTx, error)
}

type wallet struct {
	mu      *sync.Mutex
	network *chaincfg.Params

	key      *btcec.PrivateKey
	cosigner *btcec.PublicKey
	script   []byte
	addr     btcutil.Address
	prevTx   *wire.MsgTx

	indexer        btc.IndexerClient
	cosignerClient *Client
	btcClient      btc.Client
	feeEstimator   btc.FeeEstimator
}

func NewWallet(network *chaincfg.Params, privateKey *btcec.PrivateKey, cosignerPub *btcec.PublicKey, indexer btc.IndexerClient, client *Client, btcClient btc.Client, estimator btc.FeeEstimator) (Wallet, error) {

	// Calculate the address
	script, err := Script(cosignerPub.SerializeCompressed(), privateKey.PubKey().SerializeCompressed(), DefaultTimelock)
	if err != nil {
		return nil, err
	}
	addr, err := btc.P2wshAddress(script, network)
	if err != nil {
		return nil, err
	}

	// todo : Create account if first time initiated
	response, err := client.NewAccount(privateKey.PubKey())
	if err != nil {
		return nil, err
	}
	if response != addr.EncodeAddress() {
		return nil, fmt.Errorf("cosigner: cosigner address mismatch")
	}

	return &wallet{
		mu:      new(sync.Mutex),
		network: network,

		key:      privateKey,
		cosigner: cosignerPub,
		script:   script,
		addr:     addr,

		indexer:        indexer,
		cosignerClient: client,
		btcClient:      btcClient,
		feeEstimator:   estimator,
	}, nil
}

func (wal *wallet) PublicKey() *btcec.PublicKey {
	return wal.key.PubKey()
}

func (wal *wallet) Address() btcutil.Address {
	return wal.addr
}

func (wal *wallet) Execute(ctx context.Context, actions []btc.HtlcAction) (*wire.MsgTx, *wire.MsgTx, error) {
	// Separate initiate from other actions
	inits, others := []btc.HtlcAction{}, []btc.HtlcAction{}
	for _, action := range actions {
		if action.ActionType == btc.HtlcActionInitiate {
			inits = append(inits, action)
		} else {
			others = append(others, action)
		}
	}
	log.Print("have %v init and %v others", len(inits), len(others))

	// Process actions
	var initTx, otherTx *wire.MsgTx
	var err error
	if len(inits) != 0 {
		initTx, err = wal.processInits(ctx, inits)
		if err != nil {
			return nil, nil, err
		}
	}

	// Process other actions
	if len(others) != 0 {
		otherTx, err = wal.processOthers(ctx, others)
		if err != nil {
			return nil, nil, err
		}
	}

	return initTx, otherTx, nil
}

func (wal *wallet) processInits(ctx context.Context, actions []btc.HtlcAction) (*wire.MsgTx, error) {
	// Fetch address latest tx
	latest, err := wal.cosignerClient.GetLatestTransaction(wal.addr.EncodeAddress())
	if err != nil {
		return nil, err
	}

	// Fetch wallet utxos
	utxos, err := wal.fetchWalletUtxos(ctx)
	if err != nil {
		return nil, err
	}
	signer := NewSpendSigner(wal.key, wal.script, nil, true)
	signers, err := btc.NewSigners(signer, utxos...)
	if err != nil {
		return nil, err
	}

	// Create a new tx when there's no existing one
	if len(latest.TxHashes) == 0 {
		outputs := []*wire.TxOut{}
		outputsMaps := map[string]bool{}
		for _, action := range actions {
			addr := action.Htlc.MustAddress(wal.network)
			if ok := outputsMaps[addr.String()]; ok {
				continue
			}
			recipient, err := btc.NewTxOutFromAddress(addr, action.Htlc.Amount)
			if err != nil {
				return nil, err
			}
			outputs = append(outputs, recipient)
			outputsMaps[addr.String()] = true
		}

		// Get current fee rate
		feeRate, err := wal.feeEstimator.FeeSuggestion()
		if err != nil {
			return nil, err
		}

		// Build tx
		feeMode := btc.MinFeeRateMode(feeRate.High, signers)
		tx, err := btc.BuildTx(feeMode, nil, utxos, outputs, wal.addr)
		if err != nil {
			return nil, err
		}

		// Sign tx
		if err := signers.Sign(tx); err != nil {
			return nil, err
		}

		raw, err := btc.TxRawBytes(tx)
		if err != nil {
			return nil, err
		}
		log.Printf("raw = %v", hex.EncodeToString(raw))

		// Submit to cosigner server
		signedTx, err := wal.cosignerClient.NewTransaction(wal.addr.EncodeAddress(), tx)
		if err != nil {
			return nil, err
		}
		if err := wal.indexer.SubmitTx(ctx, signedTx); err != nil {
			return nil, err
		}
		return signedTx, nil
	}

	// Update the existing tx
	last := latest.TxHashes[len(latest.TxHashes)-1]
	lastTx, err := wal.indexer.GetTx(ctx, last)
	if err != nil {
		return nil, err
	}

	// Check if there are pending merge txs which haven't been merged
	mergeTxs := make([]*wire.MsgTx, 0, len(latest.MergeTransactions))
	for _, tx := range latest.MergeTransactions {
		mtx, err := wal.decodeTxFromString(tx)
		if err != nil {
			return nil, err
		}

		// If merge is not from the latest tx, we can safely ignore it
		// todo : check if this is right, especially when non-latest tx is mined
		if len(mtx.TxIn) != 1 {
			return nil, fmt.Errorf("cosigner: merge tx %v has invalid number of inputs", mtx.TxHash().String())
		}
		if mtx.TxIn[0].PreviousOutPoint.Hash.String() == last {
			mergeTxs = append(mergeTxs, mtx)
		}
	}
	txOuts := map[int][]*wire.TxOut{}
	for _, mergeTx := range mergeTxs {
		vout := mergeTx.TxIn[0].PreviousOutPoint.Index
		txOuts[int(vout)] = mergeTx.TxOut
	}

	// Including all inputs and outputs from the last tx
	inputs, outputs := make([]btc.UTXO, 0), make([]*wire.TxOut, 0)
	outputsMaps := map[string]bool{}
	for _, vin := range lastTx.VINs {
		pkScript, err := hex.DecodeString(vin.Prevout.ScriptPubKey)
		if err != nil {
			return nil, err
		}
		utxo := btc.UTXO{
			TxID:     vin.TxID,
			Vout:     uint32(vin.Vout),
			Amount:   int64(vin.Prevout.Value),
			PkScript: pkScript,
		}
		inputs = append(inputs, utxo)
		signers.AddUtxo(signer, utxo)
	}
	for i, vout := range lastTx.VOUTs {
		// Ignore the change output
		if i == len(lastTx.VOUTs)-1 && vout.ScriptPubKeyAddress == wal.Address().EncodeAddress() {
			continue
		}
		outputsMaps[vout.ScriptPubKeyAddress] = true
		pkScript, err := hex.DecodeString(vout.ScriptPubKey)
		if err != nil {
			return nil, err
		}
		mergedOuts, ok := txOuts[i]
		if ok {
			outputs = append(outputs, mergedOuts...)
		} else {
			outputs = append(outputs, &wire.TxOut{
				Value:    int64(vout.Value),
				PkScript: pkScript,
			})
		}
	}

	// Process all htlc actions
	for _, action := range actions {
		addr := action.Htlc.MustAddress(wal.network)
		if ok := outputsMaps[addr.String()]; ok {
			continue
		}
		recipient, err := btc.NewTxOutFromAddress(addr, action.Htlc.Amount)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, recipient)
		outputsMaps[addr.String()] = true
	}

	// Save all utxos to a map for quick query
	utxosMap := map[string]btc.UTXO{}
	for _, utxo := range utxos {
		utxosMap[utxo.String()] = utxo
	}
	for _, utxo := range inputs {
		utxosMap[utxo.String()] = utxo
	}

	// Build tx
	feeMode, err := btc.RbfModeFromPrevTx(wal.btcClient, last, signers)
	if err != nil {
		return nil, err
	}
	tx, err := btc.BuildTx(feeMode, inputs, utxos, outputs, wal.addr)
	if err != nil {
		return nil, err
	}

	// Sign tx
	if err := signers.Sign(tx); err != nil {
		return nil, err
	}

	// Build backup txs
	backupTxs := map[string]*wire.MsgTx{}
	for txid, mpTxStr := range latest.MempoolTransactions {
		mpTx, err := wal.decodeTxFromString(mpTxStr)
		if err != nil {
			return nil, err
		}

		var prevBackupTx *wire.MsgTx
		backupTxStr, ok := latest.BackupTransactions[txid]
		if !ok {
			prevBackupTx = mpTx
		} else {
			prevBackupTx, err = wal.decodeTxFromString(backupTxStr)
			if err != nil {
				return nil, err
			}
		}

		backupTx, err := wal.buildBackupTx(mpTx, prevBackupTx, tx, utxosMap, nil, signers)
		if err != nil {
			return nil, err
		}
		backupTxs[txid] = backupTx
	}
	raw, err := btc.TxRawBytes(tx)
	if err != nil {
		return nil, err
	}
	log.Printf("raw = %v", string(raw))

	// Submit to cosigner server
	signedTx, err := wal.cosignerClient.UpdateTransaction(wal.addr.EncodeAddress(), tx, backupTxs)
	if err != nil {
		return nil, err
	}
	if err := wal.indexer.SubmitTx(ctx, signedTx); err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (wal *wallet) processOthers(ctx context.Context, actions []btc.HtlcAction) (*wire.MsgTx, error) {
	signers, err := btc.NewSigners(nil)
	if err != nil {
		return nil, err
	}

	// Include all actions from the prevTx
	inputs, outputs := []btc.UTXO{}, []*wire.TxOut{}
	sequenceMap := map[string]int{}
	if wal.prevTx != nil {
		replacedTx, err := wal.indexer.GetTx(ctx, wal.prevTx.TxID())
		if err != nil {
			return nil, err
		}

		if !replacedTx.Status.Confirmed {
			for _, vin := range replacedTx.VINs {
				pkScript, err := hex.DecodeString(vin.Prevout.ScriptPubKey)
				if err != nil {
					return nil, err
				}
				utxo := btc.UTXO{
					TxID:     vin.TxID,
					Vout:     uint32(vin.Vout),
					Amount:   int64(vin.Prevout.Value),
					PkScript: pkScript,
				}
				inputs = append(inputs, utxo)

				// Decode the action from its witness
				witness, err := btc.DecodeWitness(*vin.Witness)
				if err != nil {
					return nil, err
				}
				action, err := btc.HtlcActionFromWitness(witness)
				if err != nil {
					return nil, err
				}

				leaf := txscript.NewTapLeaf(txscript.BaseLeafVersion, witness[len(witness)-2])
				ctrBlock, err := txscript.ParseControlBlock(witness[len(witness)-1])
				if err != nil {
					return nil, err
				}

				var signer btc.Signer
				switch action {
				case btc.HtlcActionRedeem:
					secret := witness[1]
					signer = btc.NewHtlcRedeemSigner(wal.key, leaf, *ctrBlock, secret)
					signers.AddUtxo(signer, utxo)
				case btc.HtlcActionRefund:
					timelock := int64(vin.Sequence)
					signer = btc.NewHtlcRefundSigner(wal.key, leaf, *ctrBlock, timelock)
					signers.AddUtxo(signer, utxo)
					sequenceMap[utxo.String()] = vin.Sequence
				case btc.HtlcActionInstantRefund:
					otherSig := witness[0]
					signer = btc.NewHtlcInstantRefundSigner(wal.key, leaf, *ctrBlock, true, otherSig)
					signers.AddUtxo(signer, utxo)
				}
			}

			for i, vout := range replacedTx.VOUTs {
				if i == len(replacedTx.VOUTs)-1 && vout.ScriptPubKeyAddress == wal.Address().EncodeAddress() {
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
			}
		} else {
			wal.prevTx = nil
		}
	}

	for _, action := range actions {
		switch action.ActionType {
		case btc.HtlcActionRedeem:
			// todo : no duplicate protection
			utxo, err := action.Htlc.Utxo(ctx, wal.network, wal.indexer)
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, utxo)
			leaf, ctrBlk := action.Htlc.Leaf(btc.HtlcActionRedeem)
			signer := btc.NewHtlcRedeemSigner(wal.key, leaf, ctrBlk, action.Htlc.Secret())
			signers.AddUtxo(signer, utxo)
		case btc.HtlcActionRefund:
			utxo, err := action.Htlc.Utxo(ctx, wal.network, wal.indexer)
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, utxo)
			leaf, ctrBlk := action.Htlc.Leaf(btc.HtlcActionRefund)
			signer := btc.NewHtlcRefundSigner(wal.key, leaf, ctrBlk, action.Htlc.Timelock)
			signers.AddUtxo(signer, utxo)
			sequenceMap[utxo.String()] = int(action.Htlc.Timelock)
			if action.RefundTo != nil {
				recipient, err := btc.NewTxOutFromAddress(action.RefundTo, action.Htlc.Amount)
				if err != nil {
					return nil, err
				}
				outputs = append(outputs, recipient)
			}
		case btc.HtlcActionInstantRefund:
			utxo, recipient, err := btc.ValidateInstantRefundTx(action.Htlc, action.InstantRefundTx, wal.network)
			if err != nil {
				return nil, err
			}
			inputs = append([]btc.UTXO{utxo}, inputs...)
			outputs = append([]*wire.TxOut{recipient}, outputs...)

			leaf, ctrBlk := action.Htlc.Leaf(btc.HtlcActionInstantRefund)
			signer := btc.NewHtlcInstantRefundSigner(wal.key, leaf, ctrBlk, true, action.InstantRefundTx.TxIn[0].Witness[0])
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
	feeMode := btc.MinFeeRateMode(feeRate, signers)
	if wal.prevTx != nil {
		entry, err := wal.btcClient.GetMempoolEntry(ctx, wal.prevTx.TxID())
		if err != nil {
			return nil, err
		}
		prevFees, err := btcutil.NewAmount(entry.Fees.Descendant)
		if err != nil {
			return nil, err
		}

		prevFeeRate := btc.NewSatoshiPerKb(int64(prevFees), int(entry.DescendantSize))
		feeMode = btc.RbfMode(feeRate, prevFeeRate, int64(prevFees), signers)
	}

	// Build the tx
	tx, err := btc.BuildTx(feeMode, inputs, nil, outputs, wal.Address())
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
	if err := wal.indexer.SubmitTx(ctx, tx); err != nil {
		return nil, err
	}
	wal.prevTx = tx
	return tx, nil
}

func (wal *wallet) buildBackupTx(baseTx, prevBackupTx, newTx *wire.MsgTx, utxoMaps map[string]btc.UTXO, mergeTxOuts map[string][]*wire.TxOut, signers btc.Signers) (*wire.MsgTx, error) {

	// Inputs
	// 1. Get all utxos used in the base tx
	outpointMap := map[string]btc.UTXO{}
	for _, vin := range newTx.TxIn {
		utxo := utxoMaps[vin.PreviousOutPoint.String()]
		outpointMap[utxo.String()] = utxo
	}
	// 2. Find utxos which are missing in the base tx
	for _, txin := range baseTx.TxIn {
		delete(outpointMap, txin.PreviousOutPoint.String())
	}
	// 3. Add the change utxo from the base tx
	var inputs []btc.UTXO
	if len(baseTx.TxOut) != 0 {
		// todo : we assume the change is always the last output
		// todo : what if user trying to make payment to the cosigner wallet address?  how to differentiate?
		change := baseTx.TxOut[len(baseTx.TxOut)-1]
		walPkScript, err := txscript.PayToAddrScript(wal.addr)
		if err != nil {
			return nil, err
		}
		if bytes.Equal(change.PkScript, walPkScript) {
			utxo := btc.UTXO{
				TxID:     baseTx.TxHash().String(),
				Vout:     uint32(len(baseTx.TxOut) - 1),
				Amount:   change.Value,
				PkScript: change.PkScript,
			}
			inputs = append(inputs, utxo)
			signer := NewSpendSigner(wal.key, wal.script, nil, true)
			signers.AddUtxo(signer, utxo)
		}
	}
	for _, utxo := range outpointMap {
		inputs = append(inputs, utxo)
	}

	// Outputs
	// 1. Get all output from previous backup tx
	recipients := make([]*wire.TxOut, 0)
	txOutMap := map[string]int{}
	for _, txout := range prevBackupTx.TxOut {
		outKey := fmt.Sprintf("%v:%v", hex.EncodeToString(txout.PkScript), txout.Value)
		mergeOuts, ok := mergeTxOuts[outKey]
		if ok {
			for _, mergetOut := range mergeOuts {
				mergetOutKey := fmt.Sprintf("%v:%v", hex.EncodeToString(mergetOut.PkScript), mergetOut.Value)
				txOutMap[mergetOutKey]++
				recipients = append(recipients, mergetOut)
			}
			continue
		}
		txOutMap[outKey]++
	}
	// 2. Add txout which are missing from previous backup tx
	for i, txout := range newTx.TxOut {
		// ignore change output
		if i == len(newTx.TxOut)-1 {
			// todo
			continue
		}

		outKey := fmt.Sprintf("%v:%v", hex.EncodeToString(txout.PkScript), txout.Value)
		count, ok := txOutMap[outKey]
		if !ok || count == 1 {
			output := &wire.TxOut{
				Value:    txout.Value,
				PkScript: txout.PkScript,
			}

			recipients = append(recipients, output)
			if ok {
				txOutMap[outKey]--
				if txOutMap[outKey] == 0 {
					delete(txOutMap, outKey)
				}
			}
		}
	}

	// Build the backup tx
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, err
	}
	feeMode := btc.MinFeeRateMode(feeRate.High, signers)

	tx, err := btc.BuildTx(feeMode, inputs, nil, recipients, wal.addr)
	if err != nil {
		return nil, err
	}
	if err := signers.Sign(tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (wal *wallet) fetchWalletUtxos(ctx context.Context) (btc.UTXOs, error) {
	allUtxos, err := wal.indexer.GetUTXOs(ctx, wal.addr)
	if err != nil {
		return nil, err
	}

	utxos := make([]btc.UTXO, 0, len(allUtxos))
	for _, utxo := range allUtxos {
		if utxo.Status != nil && utxo.Status.Confirmed {
			utxos = append(utxos, utxo)
		}
	}
	return utxos, nil
}

func (wal *wallet) decodeTxFromString(raw string) (*wire.MsgTx, error) {
	rawBytes, err := hex.DecodeString(raw)
	if err != nil {
		return nil, err
	}
	latestTx, err := btcutil.NewTxFromBytes(rawBytes)
	if err != nil {
		return nil, err
	}
	return latestTx.MsgTx(), nil
}
