package cosigner

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
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

	Send(ctx context.Context, recipients []btc.Recipient) (*wire.MsgTx, error)
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

func (wal *wallet) Send(ctx context.Context, recipients []btc.Recipient) (*wire.MsgTx, error) {
	// Fetch address latest tx
	latest, err := wal.cosignerClient.GetLatestTransaction(wal.addr.EncodeAddress())
	if err != nil {
		return nil, err
	}

	// Fetch confirmed utxos
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
	fetcher, err := btc.NewFetcher(wal.script, utxos...)
	if err != nil {
		return nil, err
	}
	feeRate, err := wal.feeEstimator.FeeSuggestion()
	if err != nil {
		return nil, err
	}

	// Start a new tx if no existing txs
	if len(latest.MempoolTransactions) == 0 {
		// Construct the tx
		sizer := btc.NewSizeEstimator(BaseSizeSpend, SegwitSizeSpend, utxos...)

		feeMode := btc.MinFeeRateMode(feeRate.High, sizer)
		tx, err := btc.BuildTx(wal.network, feeMode, nil, utxos, recipients, wal.addr)
		if err != nil {
			return nil, err
		}

		// Sign the tx
		fetcher, err := btc.NewFetcher(wal.script, utxos...)
		if err != nil {
			return nil, err
		}
		for i := range tx.TxIn {
			sig, err := Sign(wal.script, tx, i, fetcher, wal.key)
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
	} else {
		// We need rbf the existing tx
		// Find the latest mempool tx which is the mempool tx doesn't show up in the backup tx
		txids := map[string]bool{}
		for txid := range latest.MempoolTransactions {
			txids[txid] = true
		}
		for txid := range latest.BackupTransactions {
			delete(txids, txid)
		}
		if len(txids) != 1 {
			return nil, fmt.Errorf("cannot find latest mempool tx")
		}
		var latestTxid string
		for txid := range txids {
			latestTxid = txid
			break
		}
		latestTxStr := latest.MempoolTransactions[latestTxid]
		latestTxBytes, err := hex.DecodeString(latestTxStr)
		if err != nil {
			return nil, err
		}
		latestTx, err := btcutil.NewTxFromBytes(latestTxBytes)
		if err != nil {
			return nil, err
		}

		// Build the new rbf tx first
		latestUtxos := make([]btc.UTXO, 0)
		for _, vin := range latestTx.MsgTx().TxIn {
			value, ok := latest.Values[vin.PreviousOutPoint.String()]
			if !ok {
				return nil, fmt.Errorf("utxo [%v] not found in the values field", vin.PreviousOutPoint.String())
			}
			latestUtxos = append(latestUtxos, btc.UTXO{
				TxID:   vin.PreviousOutPoint.Hash.String(),
				Vout:   vin.PreviousOutPoint.Index,
				Amount: value,
			})
		}
		if err := btc.AddUtxosToFetcher(fetcher, wal.script, latestUtxos...); err != nil {
			return nil, err
		}
		sizer := btc.NewSizeEstimator(BaseSizeSpend, SegwitSizeSpend, append(latestUtxos, utxos...)...)
		feeMode, err := btc.RbfModeFromPrevTx(wal.btcClient, latestTxid, sizer)
		if err != nil {
			return nil, err
		}
		tx, err := wal.BuildRbfTx(latestTx.MsgTx(), feeMode, fetcher, nil, utxos, recipients)
		if err != nil {
			return nil, err
		}

		// Sign the new tx
		for i := range tx.TxIn {
			sig, err := Sign(wal.script, tx, i, fetcher, wal.key)
			if err != nil {
				return nil, err
			}
			tx.TxIn[i].Witness = Witness(wal.script, nil, sig, false)
		}

		// For each of the mempool tx we need to construct a new backup tx in case it gets mined instead of the latest.
		backups := make(map[string]*wire.MsgTx)
		for mpTxid, mpTxStr := range latest.MempoolTransactions {
			mpTxBytes, err := hex.DecodeString(mpTxStr)
			if err != nil {
				return nil, err
			}
			mpTx, err := btcutil.NewTxFromBytes(mpTxBytes)
			if err != nil {
				return nil, err
			}
			existingPrevOutpoints := map[string]bool{}
			for _, input := range mpTx.MsgTx().TxIn {
				existingPrevOutpoints[input.PreviousOutPoint.String()] = true
			}

			// Add all utxos in the unsigned tx which is not in the mempool tx
			var inputs []btc.UTXO
			for _, vin := range tx.TxIn {
				if exist := existingPrevOutpoints[vin.PreviousOutPoint.String()]; exist {
					continue
				}
				value, ok := latest.Values[vin.PreviousOutPoint.String()]
				if !ok {
					return nil, fmt.Errorf("utxo [%v] not found in the values field", vin.PreviousOutPoint.String())
				}
				utxo := btc.UTXO{
					TxID:   vin.PreviousOutPoint.Hash.String(),
					Vout:   vin.PreviousOutPoint.Index,
					Amount: value,
				}
				inputs = append(inputs, utxo)
			}

			// Add the change utxo if the mempool has one
			hasChange := false
			if len(mpTx.MsgTx().TxOut) != 0 {
				change := mpTx.MsgTx().TxOut[len(mpTx.MsgTx().TxOut)-1]
				_, addrs, numSigs, err := txscript.ExtractPkScriptAddrs(change.PkScript, wal.network)
				if err != nil {
					return nil, err
				}
				if numSigs == 1 && len(addrs) == 1 && addrs[0].EncodeAddress() == wal.addr.EncodeAddress() {
					hasChange = true
					inputs = append(inputs, btc.UTXO{
						TxID:   mpTx.MsgTx().TxHash().String(),
						Vout:   uint32(len(mpTx.MsgTx().TxOut) - 1),
						Amount: change.Value,
						Status: nil,
					})
				}
			}

			// Add all recipients in the unsigned tx which are not in the mempool tx
			var rcpts []btc.Recipient
			startIndex := len(mpTx.MsgTx().TxOut)
			if hasChange {
				startIndex -= 1
			}
			for i := startIndex; i < len(tx.TxOut); i++ {
				_, addrs, numSigs, err := txscript.ExtractPkScriptAddrs(tx.TxOut[i].PkScript, wal.network)
				if err != nil {
					return nil, err
				}
				if numSigs == 1 && len(addrs) == 1 && addrs[0].EncodeAddress() == wal.addr.EncodeAddress() {
					continue
				}
				rcpts = append(rcpts, btc.NewRecipient(addrs[0].EncodeAddress(), tx.TxOut[i].Value))
			}

			sizer := btc.NewSizeEstimator(BaseSizeSpend, SegwitSizeSpend, append(utxos, inputs...)...)
			feeMode := btc.MinFeeRateMode(feeRate.High, sizer)
			backupTx, err := btc.BuildTx(wal.network, feeMode, inputs, utxos, recipients, wal.addr)
			if err != nil {
				return nil, err
			}
			if err := btc.AddUtxosToFetcher(fetcher, wal.script, inputs...); err != nil {
				return nil, err
			}
			for index := range backupTx.TxIn {
				sig, err := Sign(wal.script, backupTx, index, fetcher, wal.key)
				if err != nil {
					return nil, err
				}
				backupTx.TxIn[index].Witness = Witness(wal.script, nil, sig, false)
			}
			backups[mpTxid] = backupTx
		}
		signedTx, err := wal.cosignerClient.UpdateTransaction(wal.addr.EncodeAddress(), tx, backups)
		if err != nil {
			return nil, err
		}
		if err := wal.indexer.SubmitTx(ctx, signedTx); err != nil {
			return nil, err
		}
		return signedTx, nil
	}
}

func (wal *wallet) BuildRbfTx(prevTx *wire.MsgTx, feeReq btc.FeeMode, fetcher txscript.PrevOutputFetcher, inputs, utxos []btc.UTXO, recipients []btc.Recipient) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(prevTx.Version)
	totalIn, totalOut := int64(0), int64(0)

	// Keep all inputs and outputs from the previous tx (except the change)
	for _, in := range prevTx.TxIn {
		tx.AddTxIn(in)
		txOut := fetcher.FetchPrevOutput(in.PreviousOutPoint)
		totalIn += txOut.Value
	}
	walPkScript, err := txscript.PayToAddrScript(wal.addr)
	if err != nil {
		return nil, err
	}
	for _, out := range prevTx.TxOut {
		// Skip the change tx
		if bytes.Equal(out.PkScript, walPkScript) {
			continue
		}
		tx.AddTxOut(out)
		totalOut += out.Value
	}

	// Add all new inputs and recipients
	for _, utxo := range inputs {
		if utxo.Amount == 0 {
			return nil, fmt.Errorf("utxo %v amount is not set", utxo.String())
		}
		txIn, err := utxo.ToTxIn()
		if err != nil {
			return nil, err
		}
		tx.AddTxIn(txIn)
		totalIn += utxo.Amount
	}
	for _, recipient := range recipients {
		txOut, err := recipient.ToTxOut(wal.network)
		if err != nil {
			return nil, err
		}
		tx.AddTxOut(txOut)
		totalOut += recipient.Amount
	}

	// Keep adding utxos to make sure it has enough input to cover the output + fees.
	for i := -1; i < len(utxos); i++ {
		// Add the utxo to the transaction input
		if i >= 0 {
			// Skip utxos that are too small
			// Todo : we might want a better way to do this, since sometimes a utxo with (DustAmount +1) could be
			// uneconomical to spend depending on the fee rate.
			if utxos[i].Amount < btc.DustAmount {
				continue
			}
			txin, err := utxos[i].ToTxIn()
			if err != nil {
				return nil, err
			}
			tx.AddTxIn(txin)
			totalIn += utxos[i].Amount
		}

		// Check if the input amount is enough to cover the fees and output
		fees, err := feeReq(tx)
		if err != nil {
			return nil, err
		}
		if totalIn >= totalOut+fees {
			// Add a change utxo to the output if the change amount is greater than the dust
			if totalIn-totalOut-fees > btc.DustAmount {
				changeScript, err := txscript.PayToAddrScript(wal.addr)
				if err != nil {
					return nil, err
				}
				tx.AddTxOut(wire.NewTxOut(0, changeScript)) // adjust the amount later

				// Fees will be changed since we add a new output field to the tx
				fees, err := feeReq(tx)
				if err != nil {
					return nil, err
				}

				// Adjust the change utxo amount if it's still enough, delete it otherwise
				if totalIn-totalOut-fees > btc.DustAmount {
					tx.TxOut[len(tx.TxOut)-1].Value = totalIn - totalOut - fees
				} else {
					tx.TxOut = tx.TxOut[:len(tx.TxOut)-1]
				}
			}
			return tx, nil
		}
	}

	return nil, fmt.Errorf("funds not enough")
}

// func (wal *wallet) Execute(ctx context.Context, actions []Action, prevTxid string) (*wire.MsgTx, error) {
// 	wal.mu.Lock()
// 	defer wal.mu.Unlock()
//
// 	// Fetch wallet utxos which are confirmed
// 	walUtxos, err := wal.indexer.GetUTXOs(ctx, wal.addr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	utxos := make([]btc.UTXO, 0, len(walUtxos))
// 	for _, utxo := range walUtxos {
// 		if utxo.Status != nil && utxo.Status.Confirmed {
// 			utxos = append(utxos, utxo)
// 		}
// 	}
//
// 	// Add wallet utxos which are used in the previous tx
// 	var replacedTx btc.Transaction
// 	var conflictUtxo *btc.UTXO
// 	inputsMap, outputsMaps := map[string]bool{}, map[string]bool{} // make sure no double executions
// 	if prevTxid != "" {
// 		replacedTx, err = wal.indexer.GetTx(ctx, prevTxid)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if replacedTx.Status.Confirmed {
// 			return nil, btc.ErrTxNotInMempool
// 		}
//
// 		for _, vin := range replacedTx.VINs {
// 			if vin.Prevout.ScriptPubKeyAddress == wal.addr.EncodeAddress() {
// 				utxo := btc.UTXO{
// 					TxID:   vin.TxID,
// 					Vout:   uint32(vin.Vout),
// 					Amount: int64(vin.Prevout.Value),
// 				}
//
// 				// Mark the first utxo from the previous tx as the conflict utxo. This conflict utxo will always present
// 				// in the inputs to make sure it will be conflicted with all replaced txs.
// 				if conflictUtxo == nil {
// 					conflictUtxo = &utxo
// 					inputsMap[vin.Prevout.ScriptPubKeyAddress] = true
// 					continue
// 				}
// 				utxos = append(utxos, utxo)
// 			}
// 		}
// 	}
//
// 	// Add all wallet utxos info to the size estimator
// 	sizer := btc.NewSizeEstimatorOfAddrType(wal.addrType, utxos...)
// 	if conflictUtxo != nil {
// 		sizer = btc.NewSizeEstimatorOfAddrType(wal.addrType, append(utxos, *conflictUtxo)...)
// 	}
//
// 	// Fetcher
// 	pkScript, err := btc.PkScript(wal.addrType, wal.externalKey.PubKey())
// 	if err != nil {
// 		return nil, err
// 	}
// 	fetcher, err := btc.NewFetcher(pkScript, utxos...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if conflictUtxo != nil {
// 		if err := btc.AddUtxosToFetcher(fetcher, pkScript, *conflictUtxo); err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	// Include all actions from the previous tx
// 	recipients := []btc.Recipient{}
// 	inputs := []btc.UTXO{}
// 	prevWitnesses := map[string]wire.TxWitness{}
// 	prevSequences := map[string]int{}
// 	if replacedTx.TxID != "" {
// 		// Inputs
// 		for _, vin := range replacedTx.VINs {
// 			if vin.Prevout.ScriptPubKeyAddress == wal.addr.EncodeAddress() {
// 				continue
// 			}
//
// 			// Store the witness and sequence from previous tx for future signing
// 			utxo := btc.UTXO{
// 				TxID:   vin.TxID,
// 				Vout:   uint32(vin.Vout),
// 				Amount: int64(vin.Prevout.Value),
// 			}
// 			inputs = append(inputs, utxo)
// 			inputsMap[vin.Prevout.ScriptPubKeyAddress] = true
// 			witness, err := decodeWitness(*vin.Witness)
// 			if err != nil {
// 				return nil, err
// 			}
// 			prevWitnesses[utxo.String()] = witness
// 			prevSequences[utxo.String()] = vin.Sequence
//
// 			// Add utxo to fetcher
// 			script, err := hex.DecodeString(vin.Prevout.ScriptPubKey)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if err := btc.AddUtxosToFetcher(fetcher, script, utxo); err != nil {
// 				return nil, err
// 			}
//
// 			// Add utxo to the size estimator depending on the action type
// 			action, err := btc.HtlcActionFromWitness(witness)
// 			if err != nil {
// 				return nil, err
// 			}
// 			switch action {
// 			case btc.HtlcActionRedeem:
// 				sizer.AddUtxos(btc.BaseSizeHtlcRedeem, btc.SegwitSizeHtlcRedeem(len(witness[1])), utxo)
// 			case btc.HtlcActionRefund:
// 				sizer.AddUtxos(btc.BaseSizeHtlcRefund, btc.SegwitSizeHtlcRefund(int64(vin.Sequence)), utxo)
// 			case btc.HtlcActionInstantRefund:
// 				sizer.AddUtxos(btc.BaseSizeHtlcInstantRefund, btc.SegwitSizeHtlcInstantRefund, utxo)
// 			}
// 		}
//
// 		// Outputs except change utxo
// 		for _, vout := range replacedTx.VOUTs {
// 			if vout.ScriptPubKeyAddress == wal.Address().EncodeAddress() {
// 				continue
// 			}
// 			recipients = append(recipients, btc.Recipient{
// 				To:     vout.ScriptPubKeyAddress,
// 				Amount: int64(vout.Value),
// 			})
// 			outputsMaps[vout.ScriptPubKeyAddress] = true
// 		}
//
// 		// Add the conflict
// 		if conflictUtxo != nil {
// 			inputs = append(inputs, *conflictUtxo)
// 			inputsMap[conflictUtxo.String()] = true
// 		}
// 	}
//
// 	// Parse the actions
// 	inputActions := map[string]btc.HtlcAction{}
// 	for _, action := range actions {
// 		addr := action.Htlc.Address(wal.network)
// 		switch action.ActionType {
// 		case btc.HtlcActionInitiate:
// 			if ok := outputsMaps[addr.String()]; ok {
// 				continue
// 			}
// 			recipients = append(recipients, btc.NewRecipient(addr.String(), action.Htlc.Amount))
// 			outputsMaps[addr.String()] = true
// 		case btc.HtlcActionRedeem, btc.HtlcActionRefund:
// 			if ok := inputsMap[addr.EncodeAddress()]; ok {
// 				continue
// 			}
//
// 			var htlcUtxos []btc.UTXO
// 			if action.ActionType == btc.HtlcActionRedeem {
// 				if !bytes.Equal(schnorr.SerializePubKey(wal.PublicKey()), action.Htlc.RedeemerPubKey) {
// 					return nil, fmt.Errorf("cannot redeem the htlc with a different key")
// 				}
// 				utxo, err := action.Htlc.Utxo(ctx, wal.network, wal.indexer)
// 				if err != nil {
// 					return nil, err
// 				}
// 				htlcUtxos = []btc.UTXO{utxo}
// 				sizer.AddUtxos(btc.BaseSizeHtlcRedeem, btc.SegwitSizeHtlcRedeem(len(action.Htlc.Secret())), utxo)
// 			} else if action.ActionType == btc.HtlcActionRefund {
// 				if !bytes.Equal(schnorr.SerializePubKey(wal.PublicKey()), action.Htlc.InitiatorPubKey) {
// 					return nil, fmt.Errorf("cannot refund the htlc with a different key")
// 				}
// 				htlcUtxos, err = action.Htlc.RefundableUtxos(ctx, wal.network, wal.indexer)
// 				if err != nil {
// 					return nil, err
// 				}
// 				sizer.AddUtxos(btc.BaseSizeHtlcRefund, SegwitSizeHtlcRefund(action.Htlc.Timelock), htlcUtxos...)
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
// 				return nil, err
// 			}
// 			if err := AddUtxosToFetcher(fetcher, fromScript, htlcUtxos...); err != nil {
// 				return nil, err
// 			}
//
// 			// If we want to refund to a different address
// 			if action.ActionType == HtlcActionRefund && action.RefundTo != nil {
// 				recipients = append(recipients, NewRecipient(action.RefundTo.String(), amount))
// 			}
// 		case HtlcActionInstantRefund:
// 			utxo, recipient, err := ValidateInstantRefundTx(action.Htlc, action.InstantRefundTx, wal.network)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if ok := inputsMap[addr.EncodeAddress()]; ok {
// 				continue
// 			}
// 			inputs = append([]UTXO{utxo}, inputs...)
// 			recipients = append([]Recipient{recipient}, recipients...)
// 			inputsMap[addr.String()] = true
// 			inputActions[utxo.String()] = action
// 			sizer.AddUtxos(BaseSizeHtlcInstantRefund, SegwitSizeHtlcInstantRefund, utxo)
//
// 			// Add to the fetcher
// 			fromScript, err := action.Htlc.P2trScript()
// 			if err != nil {
// 				return nil, err
// 			}
// 			fetcher.AddPrevOut(action.InstantRefundTx.TxIn[0].PreviousOutPoint, wire.NewTxOut(utxo.Amount, fromScript))
// 		default:
// 			return nil, errors.New("invalid action type")
// 		}
// 	}
//
// 	// Fees
// 	feeRates, err := wal.feeEstimator.FeeSuggestion()
// 	if err != nil {
// 		return nil, err
// 	}
// 	feeRate := feeRates.High
// 	feeMode := MinFeeRateMode(feeRate, sizer)
// 	if replacedTx.TxID != "" {
// 		entry, err := wal.client.GetMempoolEntry(ctx, replacedTx.TxID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		prevFees, err := btcutil.NewAmount(entry.Fees.Descendant)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		prevFeeRate := NewSatoshiPerKb(int64(prevFees), int(entry.DescendantSize))
// 		feeMode = RbfMode(feeRate, prevFeeRate, int64(prevFees), sizer)
// 	}
//
// 	// Build the tx
// 	tx, err := BuildTx(wal.network, feeMode, inputs, utxos, recipients, wal.Address())
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Set sequence number for refund inputs
// 	for i := range tx.TxIn {
// 		action, ok := inputActions[tx.TxIn[i].PreviousOutPoint.String()]
// 		if ok && action.ActionType == HtlcActionRefund {
// 			tx.TxIn[i].Sequence = uint32(action.Htlc.Timelock)
//
// 		}
// 		if sequence, ok := prevSequences[tx.TxIn[i].PreviousOutPoint.String()]; ok {
// 			tx.TxIn[i].Sequence = uint32(sequence)
// 		}
// 	}
//
// 	// Sign the tx
// 	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
// 	for i, input := range tx.TxIn {
// 		outpoint := fetcher.FetchPrevOutput(input.PreviousOutPoint)
//
// 		// If the utxo is from the previous tx, we only need to sign our signature again and reuse the rest parts of the
// 		// witness.
// 		witness, ok := prevWitnesses[input.PreviousOutPoint.String()]
// 		if ok {
// 			// Second last witness should be the script and the first witness is our signature.
// 			// We have validated the witness before storing the map, so we can confidently use it.
// 			leaf := txscript.NewTapLeaf(txscript.BaseLeafVersion, witness[len(witness)-2])
// 			sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, outpoint.Value, outpoint.PkScript, leaf, txscript.SigHashDefault, wal.externalKey)
// 			if err != nil {
// 				return nil, err
// 			}
// 			witness[0] = sig
// 			tx.TxIn[i].Witness = witness
// 			continue
// 		}
//
// 		// Rest utxos should be the wallet utxos or from new actions.
// 		action, ok := inputActions[tx.TxIn[i].PreviousOutPoint.String()]
// 		if !ok {
// 			if err := SignInput(wal.addrType, tx, i, wal.internalKey, fetcher, sigHashes); err != nil {
// 				return nil, err
// 			}
// 			continue
// 		}
//
// 		// Sign and build the witness
// 		leaf, ctrBlk := action.Htlc.Leaf(action.ActionType)
// 		ctrBlkBytes, err := ctrBlk.ToBytes()
// 		if err != nil {
// 			return nil, err
// 		}
// 		sig, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, i, outpoint.Value, outpoint.PkScript, leaf, txscript.SigHashDefault, wal.externalKey)
// 		if err != nil {
// 			return nil, err
// 		}
// 		switch action.ActionType {
// 		case HtlcActionRedeem:
// 			tx.TxIn[i].Witness = append(tx.TxIn[i].Witness, sig, action.Htlc.Secret(), leaf.Script, ctrBlkBytes)
// 		case HtlcActionRefund:
// 			tx.TxIn[i].Witness = append(tx.TxIn[i].Witness, sig, leaf.Script, ctrBlkBytes)
// 		case HtlcActionInstantRefund:
// 			initiatorSig := action.InstantRefundTx.TxIn[0].Witness[0]
// 			tx.TxIn[i].Witness = append(wire.TxWitness{}, sig, initiatorSig, leaf.Script, ctrBlkBytes)
// 		default:
// 			return nil, fmt.Errorf("unknown action type: %v", action.ActionType)
// 		}
// 	}
//
// 	// Submit tx
// 	return tx, wal.indexer.SubmitTx(ctx, tx)
// }

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
