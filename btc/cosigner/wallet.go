package cosigner

import (
	"bytes"
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
	"github.com/catalogfi/blockchain/btc"
)

type ExecutionResult struct {
	InitTx           *wire.MsgTx
	RedeemOrRefundTx *wire.MsgTx
	InstantRefundTx  *wire.MsgTx
}

func (res ExecutionResult) String() string {
	s := ""
	if res.InitTx != nil {
		s += "initTx = " + res.InitTx.TxID() + "\n"
	}
	if res.RedeemOrRefundTx != nil {
		s += "refundTx = " + res.RedeemOrRefundTx.TxID() + "\n"
	}
	if res.InstantRefundTx != nil {
		s += "instantRefundTx = " + res.InstantRefundTx.TxID() + "\n"
	}
	return s
}

type Wallet interface {
	PublicKey() *btcec.PublicKey

	Address() btcutil.Address

	Execute(ctx context.Context, actions []btc.HtlcAction) (ExecutionResult, error)
}

type wallet struct {
	mu              *sync.Mutex
	network         *chaincfg.Params
	addrType        waddrmgr.AddressType
	script          []byte
	addr            btcutil.Address
	cosigner        *btcec.PublicKey
	instantRefundTx *wire.MsgTx
	redeemTx        *wire.MsgTx

	key            *btcec.PrivateKey
	indexer        btc.IndexerClient
	cosignerClient *Client
	btcClient      btc.Client
	feeEstimator   btc.FeeEstimator
}

func NewWallet(network *chaincfg.Params, privateKey *btcec.PrivateKey, indexer btc.IndexerClient, client *Client, btcClient btc.Client, estimator btc.FeeEstimator) (Wallet, error) {
	// Fetch the cosigner's public key
	cosignerPub, err := client.CosignerPub()
	if err != nil {
		return nil, err
	}

	// Calculate the address
	script, err := Script(cosignerPub.SerializeCompressed(), privateKey.PubKey().SerializeCompressed(), DefaultTimelock)
	if err != nil {
		return nil, err
	}
	addr, err := btc.P2wshAddress(script, network)
	if err != nil {
		return nil, err
	}

	// Create an account with cosigner server in case we haven't
	response, err := client.NewAccount(privateKey.PubKey())
	if err != nil {
		return nil, err
	}
	if response != addr.EncodeAddress() {
		return nil, fmt.Errorf("cosigner: cosigner address mismatch")
	}

	return &wallet{
		mu:       new(sync.Mutex),
		network:  network,
		addrType: waddrmgr.WitnessPubKey,
		script:   script,
		addr:     addr,
		cosigner: cosignerPub,

		key:            privateKey,
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

func (wal *wallet) Execute(ctx context.Context, actions []btc.HtlcAction) (ExecutionResult, error) {
	// Separate action into different categories
	inits, spends, irs := []btc.HtlcAction{}, []btc.HtlcAction{}, []btc.HtlcAction{}
	for _, action := range actions {
		switch action.ActionType {
		case btc.HtlcActionInitiate:
			inits = append(inits, action)
		case btc.HtlcActionRedeem, btc.HtlcActionRefund:
			spends = append(spends, action)
		case btc.HtlcActionInstantRefund:
			irs = append(irs, action)
		}
	}

	result := ExecutionResult{}

	// Process inits
	if len(inits) != 0 {
		initTx, err := wal.processInits(ctx, inits)
		if err != nil {
			return ExecutionResult{}, err
		}
		result.InitTx = initTx
	}

	// Process redeems/refunds
	if len(spends) != 0 {
		tx, err := wal.processRedeemOrRefund(ctx, spends)
		if err != nil {
			return ExecutionResult{}, err
		}
		result.RedeemOrRefundTx = tx
	}

	// Process instant refunds
	if len(irs) != 0 {
		tx, err := wal.processInstantRefunds(ctx, irs)
		if err != nil {
			return ExecutionResult{}, err
		}
		result.InstantRefundTx = tx
	}

	return result, nil
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

	// Get current fee rate
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, err
	}

	// Prevent double init
	outputMap := map[string]bool{}

	// Create a new tx when there's no existing one
	inputs, outputs := []btc.UTXO{}, []*wire.TxOut{}
	var feeMode btc.FeeMode
	if len(latest.TxHashes) == 0 {
		for _, action := range actions {
			// todo : need to check if the htlc has been initiated in previous txs
			scriptPk, err := action.Htlc.ScriptPubKey()
			if err != nil {
				return nil, err
			}
			scriptPkHex := hex.EncodeToString(scriptPk)
			if ok := outputMap[scriptPkHex]; ok {
				continue
			}
			recipient := wire.NewTxOut(action.Htlc.Amount, scriptPk)
			outputs = append(outputs, recipient)
			outputMap[scriptPkHex] = true
		}

		feeMode = btc.MinFeeRateMode(feeRate.High, signers)
	} else {
		// Get the latest tx status
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
			// outputMap[hex.EncodeToString(pkScript)] = true // mark the target has been initiated
		}
		for i, vout := range lastTx.VOUTs {
			// Ignore the change output
			if i == len(lastTx.VOUTs)-1 && vout.ScriptPubKeyAddress == wal.Address().EncodeAddress() {
				continue
			}
			outputMap[vout.ScriptPubKey] = true // mark the target has been initiated
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
			scriptPk, err := action.Htlc.ScriptPubKey()
			if err != nil {
				return nil, err
			}
			scriptPkHex := hex.EncodeToString(scriptPk)
			if ok := outputMap[scriptPkHex]; ok {
				continue
			}
			recipient := wire.NewTxOut(action.Htlc.Amount, scriptPk)
			outputs = append(outputs, recipient)
			outputMap[scriptPkHex] = true
		}

		// Rbf fee mode
		feeMode, err = btc.RbfModeFromPrevTx(wal.btcClient, last, signers)
		if err != nil {
			return nil, err
		}
	}

	// Build the tx
	tx, err := btc.BuildTx(feeMode, inputs, utxos, outputs, wal.addr)
	if err != nil {
		return nil, err
	}

	// Sign tx
	if err := signers.Sign(tx); err != nil {
		return nil, err
	}

	// Submit the partial signed tx to cosigner
	var signedTx *wire.MsgTx
	if len(latest.TxHashes) == 0 {
		signedTx, err = wal.cosignerClient.NewTransaction(wal.addr.EncodeAddress(), tx)
		if err != nil {
			return nil, err
		}
	} else {
		// Save all utxos to a map for quick query
		utxosMap := map[string]btc.UTXO{}
		for _, utxo := range utxos {
			utxosMap[utxo.String()] = utxo
		}
		for _, utxo := range inputs {
			utxosMap[utxo.String()] = utxo
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

		signedTx, err = wal.cosignerClient.UpdateTransaction(wal.addr.EncodeAddress(), tx, backupTxs)
		if err != nil {
			return nil, err
		}
	}

	return signedTx, wal.indexer.SubmitTx(ctx, signedTx)
}

func (wal *wallet) processRedeemOrRefund(ctx context.Context, actions []btc.HtlcAction) (*wire.MsgTx, error) {
	signers, err := btc.NewSigners(nil)
	if err != nil {
		return nil, err
	}

	// Include all actions from the redeemTx
	inputs, outputs := []btc.UTXO{}, []*wire.TxOut{}
	sequenceMap := map[string]int{}
	if wal.redeemTx != nil {
		replacedTx, err := wal.indexer.GetTx(ctx, wal.redeemTx.TxID())
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
			wal.redeemTx = nil
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
	if wal.redeemTx != nil {
		entry, err := wal.btcClient.GetMempoolEntry(ctx, wal.redeemTx.TxID())
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
	wal.redeemTx = tx
	return tx, nil
}

func (wal *wallet) processInstantRefunds(ctx context.Context, irs []btc.HtlcAction) (*wire.MsgTx, error) {
	// todo
	panic("todo")
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
