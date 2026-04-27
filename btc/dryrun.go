package btc

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type rejectCategory int

const (
	rejectScrap rejectCategory = iota
	rejectMissingOrSpent
	rejectScriptVerify
	rejectAddsUnconfirmed
)

func classifyRejectReason(reason string) rejectCategory {
	r := strings.ToLower(reason)
	switch {
	case strings.Contains(r, "bad-txns-inputs-missingorspent"),
		strings.Contains(r, "missingorspent"),
		strings.Contains(r, "missing-or-spent"):
		return rejectMissingOrSpent
	case strings.Contains(r, "mandatory-script-verify-flag-failed"),
		strings.Contains(r, "non-mandatory-script-verify-flag"):
		return rejectScriptVerify
	case strings.Contains(r, "replacement-adds-unconfirmed"):
		return rejectAddsUnconfirmed
	default:
		return rejectScrap
	}
}

// dryRunAccept serializes tx and asks the Bitcoin node to test-accept it into
// the mempool. Returns (rejected, reason, err) — rejected is true only when the
// RPC completed and the node refused the tx.
func (w *batcherWallet) dryRunAccept(c context.Context, tx *wire.MsgTx) (bool, string, error) {
	txBytes, err := GetTxRawBytes(tx)
	if err != nil {
		return false, "", err
	}
	rawHex := hex.EncodeToString(txBytes)
	var acceptRes *TestMempoolAcceptResult
	err = withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		acceptRes, err = w.rpc.TestMempoolAccept(ctx, rawHex)
		return err
	})
	if err != nil {
		return false, "", err
	}
	if acceptRes == nil {
		return false, "", nil
	}
	return !acceptRes.Allowed, acceptRes.RejectReason, nil
}

func (w *batcherWallet) findSpentInputOwner(ctx context.Context, tx *wire.MsgTx, requests []BatcherRequest) (string, bool, error) {
	ownerByOutpoint, err := buildOutpointOwnerMap(requests)
	if err != nil {
		return "", false, err
	}

	for _, vin := range tx.TxIn {
		requestID, ok := ownerByOutpoint[outpointKey(vin.PreviousOutPoint)]
		if !ok {
			continue
		}
		spent, err := w.isOutpointSpent(ctx, vin.PreviousOutPoint)
		if err != nil {
			return "", false, err
		}
		if spent {
			return requestID, true, nil
		}
	}

	return "", false, nil
}

func (w *batcherWallet) findBadScriptInputOwner(ctx context.Context, tx *wire.MsgTx, requests []BatcherRequest) (string, bool, error) {
	prevouts, err := w.buildDryRunPrevouts(ctx, requests)
	if err != nil {
		return "", false, err
	}

	hashCache := txscript.NewTxSigHashes(tx, prevouts.fetcher)
	for i, vin := range tx.TxIn {
		requestID, ok := prevouts.ownerByOutpoint[outpointKey(vin.PreviousOutPoint)]
		if !ok {
			continue
		}

		prevout, ok := prevouts.outputs[outpointKey(vin.PreviousOutPoint)]
		if !ok {
			return "", false, fmt.Errorf("prevout %s not found for script verification", outpointKey(vin.PreviousOutPoint))
		}

		engine, err := txscript.NewEngine(
			prevout.PkScript,
			tx,
			i,
			txscript.StandardVerifyFlags,
			nil,
			hashCache,
			prevout.Value,
			prevouts.fetcher,
		)
		if err == nil {
			err = engine.Execute()
		}
		if err != nil {
			return requestID, true, nil
		}
	}

	return "", false, nil
}

type dryRunPrevouts struct {
	fetcher         txscript.PrevOutputFetcher
	outputs         map[string]*wire.TxOut
	ownerByOutpoint map[string]string
}

func (w *batcherWallet) buildDryRunPrevouts(ctx context.Context, requests []BatcherRequest) (*dryRunPrevouts, error) {
	builder := NewPrevOutFetcherBuilder()
	outputs := make(map[string]*wire.TxOut)
	ownerByOutpoint := make(map[string]string)

	register := func(requestID string, outpoint wire.OutPoint, txOut *wire.TxOut) error {
		key := outpointKey(outpoint)
		if existing, ok := ownerByOutpoint[key]; ok && existing != requestID {
			return fmt.Errorf("outpoint %s belongs to multiple requests (%s, %s)", key, existing, requestID)
		}
		ownerByOutpoint[key] = requestID
		outputs[key] = txOut
		builder.AddPrevOut(outpoint, txOut)
		return nil
	}

	for _, request := range requests {
		for _, spend := range request.Spends {
			pkScript, err := txscript.PayToAddrScript(spend.ScriptAddress)
			if err != nil {
				return nil, err
			}
			for _, utxo := range spend.Utxos {
				hash, err := chainhash.NewHashFromStr(utxo.TxID)
				if err != nil {
					return nil, err
				}
				outpoint := wire.OutPoint{Hash: *hash, Index: utxo.Vout}
				if err := register(request.ID, outpoint, &wire.TxOut{PkScript: pkScript, Value: utxo.Amount}); err != nil {
					return nil, err
				}
			}
		}

		for _, sacp := range request.SACPs {
			sacpTx, _, err := buildTxFromSacps([][]byte{sacp})
			if err != nil {
				return nil, err
			}
			for _, vin := range sacpTx.TxIn {
				prevTx, err := w.indexer.GetTx(ctx, vin.PreviousOutPoint.Hash.String())
				if err != nil {
					return nil, err
				}
				if int(vin.PreviousOutPoint.Index) >= len(prevTx.VOUTs) {
					return nil, fmt.Errorf("prevout %s missing vout %d", vin.PreviousOutPoint.Hash.String(), vin.PreviousOutPoint.Index)
				}
				prevVout := prevTx.VOUTs[vin.PreviousOutPoint.Index]
				pkScript, err := hex.DecodeString(prevVout.ScriptPubKey)
				if err != nil {
					return nil, err
				}
				if err := register(request.ID, vin.PreviousOutPoint, &wire.TxOut{
					PkScript: pkScript,
					Value:    int64(prevVout.Value),
				}); err != nil {
					return nil, err
				}
			}
		}
	}

	return &dryRunPrevouts{
		fetcher:         builder.Build(),
		outputs:         outputs,
		ownerByOutpoint: ownerByOutpoint,
	}, nil
}

func buildOutpointOwnerMap(requests []BatcherRequest) (map[string]string, error) {
	ownerByOutpoint := make(map[string]string)

	register := func(requestID string, outpoint wire.OutPoint) error {
		key := outpointKey(outpoint)
		if existing, ok := ownerByOutpoint[key]; ok && existing != requestID {
			return fmt.Errorf("outpoint %s belongs to multiple requests (%s, %s)", key, existing, requestID)
		}
		ownerByOutpoint[key] = requestID
		return nil
	}

	for _, request := range requests {
		for _, spend := range request.Spends {
			for _, utxo := range spend.Utxos {
				hash, err := chainhash.NewHashFromStr(utxo.TxID)
				if err != nil {
					return nil, err
				}
				if err := register(request.ID, wire.OutPoint{Hash: *hash, Index: utxo.Vout}); err != nil {
					return nil, err
				}
			}
		}

		for _, sacp := range request.SACPs {
			sacpTx, _, err := buildTxFromSacps([][]byte{sacp})
			if err != nil {
				return nil, err
			}
			for _, vin := range sacpTx.TxIn {
				if err := register(request.ID, vin.PreviousOutPoint); err != nil {
					return nil, err
				}
			}
		}
	}

	return ownerByOutpoint, nil
}

func (w *batcherWallet) isOutpointSpent(c context.Context, outpoint wire.OutPoint) (bool, error) {
	var txOut *btcjson.GetTxOutResult
	err := withContextTimeout(c, DefaultAPITimeout, func(ctx context.Context) error {
		var err error
		txOut, err = w.rpc.GetTxOut(ctx, &outpoint.Hash, outpoint.Index)
		return err
	})
	if err != nil {
		return false, err
	}
	return txOut == nil, nil
}

func outpointKey(outpoint wire.OutPoint) string {
	return fmt.Sprintf("%s:%d", outpoint.Hash.String(), outpoint.Index)
}
