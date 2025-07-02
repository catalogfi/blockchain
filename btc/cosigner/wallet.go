package cosigner

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

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

	// Send(ctx context.Context, recipients []*wire.TxOut) (*wire.MsgTx, error)

	Patch(ctx context.Context, outputs []*wire.TxOut) (*wire.MsgTx, error)

	// Merge(ctx context.Context) (*wire.MsgTx, error)
	// Execute(ctx context.Context, actions []Action, prevTxid string) (*wire.MsgTx, error)
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
	fetcher        *btc.InMemFetcher
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
	fetcher, err := btc.NewInMemFetcher(time.Hour, indexer)
	if err != nil {
		return nil, err
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
		fetcher:        fetcher,
	}, nil
}

func (wal *wallet) PublicKey() *btcec.PublicKey {
	return wal.key.PubKey()
}

func (wal *wallet) Address() btcutil.Address {
	return wal.addr
}

func (wal *wallet) Patch(ctx context.Context, outputs []*wire.TxOut) (*wire.MsgTx, error) {
	// Fetch address latest tx
	latest, err := wal.cosignerClient.GetLatestTransaction(wal.addr.EncodeAddress())
	if err != nil {
		return nil, err
	}

	// Check if we need to create a new tx or update the existing tx
	if len(latest.TxHashes) > 0 {
		last := latest.TxHashes[len(latest.TxHashes)-1]
		lastTxStr, ok := latest.MempoolTransactions[last]
		if !ok {
			return nil, fmt.Errorf("cannot find mempool transaction = %v", lastTxStr)
		}
		lastTx, err := wal.decodeTxFromString(lastTxStr)
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
		lastUtxos := make([]btc.UTXO, 0)
		for _, vin := range lastTx.TxIn {
			value, ok := latest.Values[vin.PreviousOutPoint.String()]
			if !ok {
				return nil, fmt.Errorf("utxo [%v] not found in the values field", vin.PreviousOutPoint.String())
			}
			lastUtxos = append(lastUtxos, btc.UTXO{
				TxID:   vin.PreviousOutPoint.Hash.String(),
				Vout:   vin.PreviousOutPoint.Index,
				Amount: value,
			})
		}
		utxos, err := wal.fetchWalletUtxos(ctx)
		if err != nil {
			return nil, err
		}
		// todo : we assume all utxos are from the wallet for now
		wal.fetcher.AddUtxo(wal.script, append(utxos, lastUtxos...)...)
		// fetcher, err := btc.NewFetcher(wal.script, append(utxos, lastUtxos...)...)
		// if err != nil {
		// 	return nil, err
		// }
		sizer := btc.NewSizeEstimator(BaseSizeSpend, SegwitSizeSpend, append(lastUtxos, utxos...)...)
		feeMode, err := btc.RbfModeFromPrevTx(wal.btcClient, last, sizer)
		if err != nil {
			return nil, err
		}
		inputs := make([]btc.UTXO, 0)
		for _, input := range lastTx.TxIn {
			out := wal.fetcher.FetchPrevOutput(input.PreviousOutPoint)
			inputs = append(inputs, btc.UtxoFromOutPoint(input.PreviousOutPoint, out.Value))
		}
		recipients := make([]*wire.TxOut, 0)
		walPkSript, err := txscript.PayToAddrScript(wal.addr)
		if err != nil {
			return nil, err
		}
		for i, output := range lastTx.TxOut {
			if i == len(lastTx.TxOut)-1 && bytes.Equal(output.PkScript, walPkSript) {
				continue
			}
			outKey := fmt.Sprintf("%v:%v", hex.EncodeToString(output.PkScript), output.Value)
			newOutputs, ok := txOuts[outKey]
			if ok {
				recipients = append(recipients, newOutputs...)
			} else {
				recipients = append(recipients, output)
			}
		}
		recipients = append(recipients, outputs...)

		tx, err := btc.BuildTx(feeMode, inputs, nil, recipients, wal.addr)
		if err != nil {
			return nil, err
		}

		// Sign the new tx
		for i := range tx.TxIn {
			sig, err := Sign(wal.script, tx, i, wal.fetcher, wal.key)
			if err != nil {
				return nil, err
			}
			tx.TxIn[i].Witness = Witness(wal.script, nil, sig, false)
		}

		// todo : build backup tx
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

			backupTx, err := wal.buildBackupTx(mpTx, prevBackupTx, tx, txOuts)
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
		sizer := btc.NewSizeEstimator(BaseSizeSpend, SegwitSizeSpend, utxos...)
		feeRate, err := wal.feeEstimator.FeeSuggestion()
		if err != nil {
			return nil, err
		}
		feeMode := btc.MinFeeRateMode(feeRate.High, sizer)
		tx, err := btc.BuildTx(feeMode, nil, utxos, outputs, wal.addr)
		if err != nil {
			return nil, err
		}

		// Sign the tx
		wal.fetcher.AddUtxo(wal.script, utxos...)
		for i := range tx.TxIn {
			sig, err := Sign(wal.script, tx, i, wal.fetcher, wal.key)
			if err != nil {
				return nil, err
			}
			tx.TxIn[i].Witness = Witness(wal.script, nil, sig, false)
		}

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
}

func (wal *wallet) buildBackupTx(baseTx, prevBackupTx, newTx *wire.MsgTx, mergeTxOuts map[string][]*wire.TxOut) (*wire.MsgTx, error) {

	// Inputs
	// 1. Get all utxos used in the new tx
	outpointMap := map[string]btc.UTXO{}
	for _, txin := range newTx.TxIn {
		output := wal.fetcher.FetchPrevOutput(txin.PreviousOutPoint)
		outpointMap[txin.PreviousOutPoint.String()] = btc.UTXO{
			TxID:   txin.PreviousOutPoint.Hash.String(),
			Vout:   txin.PreviousOutPoint.Index,
			Amount: output.Value,
		}
	}
	// 2. Find utxos which are missing in the base tx
	var inputs []btc.UTXO
	for _, txin := range prevBackupTx.TxIn {
		_, ok := outpointMap[txin.PreviousOutPoint.String()]
		if !ok {
			output := wal.fetcher.FetchPrevOutput(txin.PreviousOutPoint)
			inputs = append(inputs, btc.UTXO{
				TxID:   txin.PreviousOutPoint.Hash.String(),
				Vout:   txin.PreviousOutPoint.Index,
				Amount: output.Value,
			})
		}
	}
	// 3. Add the change utxo from the base tx
	if len(baseTx.TxOut) != 0 {
		// todo : we assume the change is always the last output
		// todo : what if user trying to make payment to the cosigner wallet address?  how to differentiate?
		change := baseTx.TxOut[len(baseTx.TxOut)-1]
		_, addrs, numSigs, err := txscript.ExtractPkScriptAddrs(change.PkScript, wal.network)
		if err != nil {
			return nil, err
		}
		if numSigs == 1 && len(addrs) == 1 && addrs[0].EncodeAddress() == wal.addr.EncodeAddress() {
			inputs = append(inputs, btc.UTXO{
				TxID:   baseTx.TxHash().String(),
				Vout:   uint32(len(baseTx.TxOut) - 1),
				Amount: change.Value,
			})
		}
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
			_, addrs, numSigs, err := txscript.ExtractPkScriptAddrs(txout.PkScript, wal.network)
			if err != nil {
				return nil, err
			}
			if numSigs == 1 && len(addrs) == 1 && addrs[0].EncodeAddress() == wal.addr.EncodeAddress() {
				continue
			}
		}

		outKey := fmt.Sprintf("%v:%v", hex.EncodeToString(txout.PkScript), txout.Value)
		count, ok := txOutMap[outKey]
		if !ok || count == 1 {
			recipients = append(recipients, txout)
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
	sizer := btc.NewSizeEstimator(BaseSizeSpend, SegwitSizeSpend, inputs...)
	feeMode := btc.MinFeeRateMode(feeRate.High, sizer)
	wal.fetcher.AddUtxo(wal.script, inputs...)
	tx, err := btc.BuildTx(feeMode, inputs, nil, recipients, wal.addr)
	if err != nil {
		return nil, err
	}

	// Sign the tx
	for i := range tx.TxIn {
		sig, err := Sign(wal.script, tx, i, wal.fetcher, wal.key)
		if err != nil {
			return nil, err
		}
		tx.TxIn[i].Witness = Witness(wal.script, nil, sig, false)
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
