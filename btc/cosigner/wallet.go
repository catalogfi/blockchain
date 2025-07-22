package cosigner

import (
	"bytes"
	"context"
	"encoding/hex"
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

	Patch(ctx context.Context, outputs []*wire.TxOut) (*wire.MsgTx, error)

	// Execute(ctx context.Context, actions []btc.HtlcAction) (*wire.MsgTx, error)
}

type wallet struct {
	mu      *sync.Mutex
	network *chaincfg.Params

	key      *btcec.PrivateKey
	cosigner *btcec.PublicKey
	script   []byte
	addr     btcutil.Address

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

// func (wal *wallet) Execute(ctx context.Context, actions []btc.HtlcAction) (*wire.MsgTx, error) {
// 	// Fetch address latest tx
// 	latest, err := wal.cosignerClient.GetLatestTransaction(wal.addr.EncodeAddress())
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Check if we need to create a new tx or update the existing tx
// 	if len(latest.TxHashes) == 0 {
//
// 		// Fetch wallet utxos
// 		utxos, err := wal.fetchWalletUtxos(ctx)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// NewSpendSigner(wal.key, wal.script, nil, true)
// 		wal.sizer.AddUtxos(BaseSizeSpend, SegwitSizeSpend, utxos...)
// 		wal.fetcher.AddUtxo(wal.script, utxos...)
//
// 		// Process all htlc actions
// 		inputs, outputs, err := wal.processActions(ctx, actions)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		// Get current fee rate
// 		feeRate, err := wal.feeEstimator.FeeSuggestion()
// 		if err != nil {
// 			return nil, err
// 		}
// 		feeMode := btc.MinFeeRateMode(feeRate.High, wal.sizer)
// 		tx, err := btc.BuildTx(feeMode, inputs, utxos, outputs, wal.addr)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		// Sign the tx
// 		for i := range tx.TxIn {
// 			sig, err := Sign(wal.script, tx, i, wal.fetcher, wal.key)
// 			if err != nil {
// 				return nil, err
// 			}
// 			tx.TxIn[i].Witness = Witness(wal.script, nil, sig, false)
// 		}
//
// 		// Submit to cosigner server
// 		signedTx, err := wal.cosignerClient.NewTransaction(wal.addr.EncodeAddress(), tx)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if err := wal.indexer.SubmitTx(ctx, signedTx); err != nil {
// 			return nil, err
// 		}
// 		return signedTx, nil
// 	} else {
//
// 	}
//
// 	// No new actions. We'll need to check
// 	// 1) if there's any merge tx waiting to be merged
// 	// 2) if the fee has increased and our previous tx is not good enough to be in the next mined block
//
// 	// TODO : FINISH THIS
// 	panic("unimplemented")
// }

func (wal *wallet) Patch(ctx context.Context, outputs []*wire.TxOut) (*wire.MsgTx, error) {
	// Fetch address latest tx
	latest, err := wal.cosignerClient.GetLatestTransaction(wal.addr.EncodeAddress())
	if err != nil {
		return nil, err
	}

	// Check if we need to create a new tx or update the existing tx
	if len(latest.TxHashes) > 0 {
		last := latest.TxHashes[len(latest.TxHashes)-1]
		// lastTxStr, ok := latest.MempoolTransactions[last]
		// if !ok {
		// 	return nil, fmt.Errorf("cannot find mempool transaction = %v", lastTxStr)
		// }
		// lastTx, err := wal.decodeTxFromString(lastTxStr)
		// if err != nil {
		// 	return nil, err
		// }
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
		txOuts := map[string][]*wire.TxOut{}
		for _, mergeTx := range mergeTxs {
			// Get the TxOut from the mergeTx input
			if len(mergeTx.TxIn) != 1 {
				return nil, fmt.Errorf("invalid merge tx %v : inputs more than 1", mergeTx.TxHash().String())
			}
			txid := mergeTx.TxIn[0].PreviousOutPoint.Hash.String()
			vout := mergeTx.TxIn[0].PreviousOutPoint.Index
			prevTxStr, ok := latest.MempoolTransactions[txid]
			if !ok {
				return nil, fmt.Errorf("cannot find mempool transaction = %v", txid)
			}
			prevTxBytes, err := hex.DecodeString(prevTxStr)
			if err != nil {
				return nil, err
			}
			prevTx, err := btcutil.NewTxFromBytes(prevTxBytes)
			if err != nil {
				return nil, err
			}

			txOut := prevTx.MsgTx().TxOut[vout]
			outKey := fmt.Sprintf("%v:%v", hex.EncodeToString(txOut.PkScript), txOut.Value)

			txOuts[outKey] = mergeTx.TxOut
		}

		// If there're no new outputs, nor new merge txs, only thing left to check is the last tx's fee rate.
		if len(outputs) == 0 && len(mergeTxs) == 0 {
			// todo : check last tx position and do a rbf to update fee rate if needed
		}

		// Create a new tx basing on the mergeTx and new outputs
		inputs := make([]btc.UTXO, 0)
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
		}
		utxos, err := wal.fetchWalletUtxos(ctx)
		if err != nil {
			return nil, err
		}

		signer := NewSpendSigner(wal.key, wal.script, nil, true)
		signers, err := btc.NewSigners(signer, append(utxos, inputs...)...)
		if err != nil {
			return nil, err
		}

		feeMode, err := btc.RbfModeFromPrevTx(wal.btcClient, last, signers)
		if err != nil {
			return nil, err
		}

		recipients := make([]*wire.TxOut, 0)
		for i, vout := range lastTx.VOUTs {
			if i == len(lastTx.VOUTs)-1 && vout.ScriptPubKeyAddress == wal.Address().EncodeAddress() {
				continue
			}

			pkScript, err := hex.DecodeString(vout.ScriptPubKey)
			if err != nil {
				return nil, err
			}
			recipients = append(recipients, &wire.TxOut{
				Value:    int64(vout.Value),
				PkScript: pkScript,
			})
		}
		recipients = append(recipients, outputs...)

		tx, err := btc.BuildTx(feeMode, inputs, utxos, recipients, wal.addr)
		if err != nil {
			return nil, err
		}
		utxosMap := map[string]btc.UTXO{}
		for _, utxo := range utxos {
			utxosMap[utxo.TxID] = utxo
		}
		for _, utxo := range inputs {
			utxosMap[utxo.TxID] = utxo
		}

		if err := signers.Sign(tx); err != nil {
			return nil, err
		}

		// todo : build backup tx
		backupTxs := map[string]*wire.MsgTx{}
		for txid, mpTxStr := range latest.MempoolTransactions {
			mpTx, err := wal.decodeTxFromString(mpTxStr)
			if err != nil {
				return nil, err
			}

			backupTx, err := wal.buildBackupTx(tx, mpTx, utxosMap, txOuts, signers)
			if err != nil {
				return nil, err
			}
			backupTxs[txid] = backupTx
		}

		// Submit the tx to cosigner
		signedTx, err := wal.cosignerClient.UpdateTransaction(wal.addr.EncodeAddress(), tx, backupTxs)
		if err != nil {
			return nil, err
		}
		if err := wal.indexer.SubmitTx(ctx, signedTx); err != nil {
			return nil, err
		}
		return signedTx, nil

	} else {
		// Construct a new tx of the cosigner account
		utxos, err := wal.fetchWalletUtxos(ctx)
		if err != nil {
			return nil, err
		}
		signer := NewSpendSigner(wal.key, wal.script, nil, true)
		signers, err := btc.NewSigners(signer, utxos...)
		if err != nil {
			return nil, err
		}
		feeRate, err := wal.feeEstimator.FeeSuggestion()
		if err != nil {
			return nil, err
		}
		feeMode := btc.MinFeeRateMode(feeRate.High, signers)
		tx, err := btc.BuildTx(feeMode, nil, utxos, outputs, wal.addr)
		if err != nil {
			return nil, err
		}

		// Sign the tx
		if err := signers.Sign(tx); err != nil {
			return nil, err
		}

		// Submit to cosigner server
		signedTx, err := wal.cosignerClient.NewTransaction(wal.addr.EncodeAddress(), tx)
		if err != nil {
			log.Print(1111)
			return nil, err
		}
		if err := wal.indexer.SubmitTx(ctx, signedTx); err != nil {
			return nil, err
		}
		return signedTx, nil
	}
}

func (wal *wallet) buildBackupTx(newTx, baseTx *wire.MsgTx, utxoMaps map[string]btc.UTXO, mergeTxOuts map[string][]*wire.TxOut, signers btc.Signers) (*wire.MsgTx, error) {

	// Inputs
	// 1. Get all utxos used in the new tx
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
			inputs = append(inputs, btc.UTXO{
				TxID:     baseTx.TxHash().String(),
				Vout:     uint32(len(baseTx.TxOut) - 1),
				Amount:   change.Value,
				PkScript: change.PkScript,
			})
		}
	}
	for _, utxo := range outpointMap {
		inputs = append(inputs, utxo)
	}

	// Outputs
	// 1. Get all output from previous backup tx
	txOutMap := map[string]int{}
	for _, txout := range baseTx.TxOut {
		outKey := fmt.Sprintf("%v:%v", hex.EncodeToString(txout.PkScript), txout.Value)
		mergeOuts, ok := mergeTxOuts[outKey]
		if ok {
			for _, mergetOut := range mergeOuts {
				mergetOutKey := fmt.Sprintf("%v:%v", hex.EncodeToString(mergetOut.PkScript), mergetOut.Value)
				txOutMap[mergetOutKey]++
			}
			continue
		}
		txOutMap[outKey]++
	}
	// 2. Add txout which are missing from base backup tx
	recipients := make([]*wire.TxOut, 0)
	for i, txout := range newTx.TxOut {
		// ignore change output
		if i == len(newTx.TxIn)-1 {
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

// func (wal *wallet) processActions(ctx context.Context, actions []btc.HtlcAction) ([]btc.UTXO, []*wire.TxOut, error) {
// 	inputs, outputs := []btc.UTXO{}, []*wire.TxOut{}
// 	inputsMap, outputsMaps := map[string]bool{}, map[string]bool{} // make sure no double executions
// 	for _, action := range actions {
// 		addr := action.Htlc.MustAddress(wal.network)
// 		switch action.ActionType {
// 		case btc.HtlcActionInitiate:
// 			if ok := outputsMaps[addr.String()]; ok {
// 				continue
// 			}
// 			// todo : check if the address has been initiated
// 			txOut, err := btc.NewTxOutFromAddress(addr, action.Htlc.Amount)
// 			if err != nil {
// 				return nil, nil, err
// 			}
// 			outputs = append(outputs, txOut)
// 			outputsMaps[addr.String()] = true
// 		case btc.HtlcActionRedeem, btc.HtlcActionRefund:
// 			if ok := inputsMap[addr.EncodeAddress()]; ok {
// 				continue
// 			}
//
// 			var htlcUtxos []btc.UTXO
// 			if action.ActionType == btc.HtlcActionRedeem {
// 				if !bytes.Equal(schnorr.SerializePubKey(wal.PublicKey()), action.Htlc.RedeemerPubKey) {
// 					return nil, nil, fmt.Errorf("cannot redeem the htlc with a different key")
// 				}
// 				utxo, err := action.Htlc.Utxo(ctx, wal.network, wal.indexer)
// 				if err != nil {
// 					return nil, nil, err
// 				}
// 				htlcUtxos = []btc.UTXO{utxo}
// 				wal.sizer.AddUtxos(btc.BaseSizeHtlcRedeem, btc.SegwitSizeHtlcRedeem(len(action.Htlc.Secret())), utxo)
// 			} else if action.ActionType == btc.HtlcActionRefund {
// 				if !bytes.Equal(schnorr.SerializePubKey(wal.PublicKey()), action.Htlc.InitiatorPubKey) {
// 					return nil, nil, fmt.Errorf("cannot refund the htlc with a different key")
// 				}
// 				var err error
// 				htlcUtxos, err = action.Htlc.RefundableUtxos(ctx, wal.network, wal.indexer)
// 				if err != nil {
// 					return nil, nil, err
// 				}
// 				wal.sizer.AddUtxos(btc.BaseSizeHtlcRefund, btc.SegwitSizeHtlcRefund(action.Htlc.Timelock), htlcUtxos...)
// 			}
// 			inputs = append(inputs, htlcUtxos...)
// 			inputsMap[addr.String()] = true
// 			amount := int64(0)
// 			for _, utxo := range htlcUtxos {
// 				inputActions[utxo.String()] = action
// 				amount += utxo.Amount
// 			}
//
// 			// Add utxo to the fetcher
// 			fromScript, err := action.Htlc.P2trScript()
// 			if err != nil {
// 				return nil, nil, err
// 			}
// 			wal.fetcher.AddUtxo(fromScript, htlcUtxos...)
//
// 			// If we want to refund to a different address
// 			if action.ActionType == btc.HtlcActionRefund && action.RefundTo != nil {
// 				txOut, err := btc.NewTxOutFromAddress(action.RefundTo, amount)
// 				if err != nil {
// 					return nil, nil, err
// 				}
// 				outputs = append(outputs, txOut)
// 			}
// 		case btc.HtlcActionInstantRefund:
// 			utxo, recipient, err := btc.ValidateInstantRefundTx(action.Htlc, action.InstantRefundTx, wal.network)
// 			if err != nil {
// 				return nil, nil, err
// 			}
// 			if ok := inputsMap[addr.EncodeAddress()]; ok {
// 				continue
// 			}
// 			inputs = append([]btc.UTXO{utxo}, inputs...)
// 			outputs = append([]*wire.TxOut{recipient}, outputs...)
// 			inputsMap[addr.String()] = true
// 			inputActions[utxo.String()] = action
// 			wal.sizer.AddUtxos(btc.BaseSizeHtlcInstantRefund, btc.SegwitSizeHtlcInstantRefund, utxo)
//
// 			// Add to the fetcher
// 			fromScript, err := action.Htlc.P2trScript()
// 			if err != nil {
// 				return nil, nil, err
// 			}
// 			wal.fetcher.AddPrevOut(action.InstantRefundTx.TxIn[0].PreviousOutPoint, wire.NewTxOut(utxo.Amount, fromScript))
// 		default:
// 			return nil, nil, errors.New("invalid action type")
// 		}
// 	}
// 	return inputs, outputs, nil
// }
