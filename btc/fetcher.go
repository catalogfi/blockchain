package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type PrevOutFetcherBuilder struct {
	fetcher *txscript.MultiPrevOutFetcher
}

func NewPrevOutFetcherBuilder() *PrevOutFetcherBuilder {
	return &PrevOutFetcherBuilder{
		fetcher: txscript.NewMultiPrevOutFetcher(nil),
	}
}

func (pf *PrevOutFetcherBuilder) AddPrevOut(outpoint wire.OutPoint, txOut *wire.TxOut) {
	pf.fetcher.AddPrevOut(outpoint, txOut)
}

func (pf *PrevOutFetcherBuilder) AddPrevOuts(outpoints []wire.OutPoint, txOuts []*wire.TxOut) error {
	if len(outpoints) != len(txOuts) {
		return fmt.Errorf("outpoints and txOuts length mismatch")
	}
	for i := range outpoints {
		pf.AddPrevOut(outpoints[i], txOuts[i])
	}
	return nil
}

func (pf *PrevOutFetcherBuilder) AddFromUTXOs(utxos UTXOs, scriptAddr btcutil.Address) error {
	script, err := txscript.PayToAddrScript(scriptAddr)
	if err != nil {
		return err
	}
	for _, utxo := range utxos {
		txId, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return err
		}
		pf.AddPrevOut(wire.OutPoint{Hash: *txId, Index: utxo.Vout}, &wire.TxOut{PkScript: script, Value: utxo.Amount})
	}
	return nil
}

func (pf *PrevOutFetcherBuilder) Build() *txscript.MultiPrevOutFetcher {
	return pf.fetcher
}
