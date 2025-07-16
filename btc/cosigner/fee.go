package cosigner

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/wire"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/tools"
	"github.com/catalogfi/tools/pkg/memcache"
)

// cachedSizeEstimator implements the `btc.SizeEstimator` using an expirable memcache.
type cachedSizeEstimator struct {
	mu     *sync.Mutex
	base   memcache.Cache[int]
	segwit memcache.Cache[int]
}

// NewCachedSizeEstimator returns an SizeEstimator with some preload UTXOs.
func NewCachedSizeEstimator(ttl time.Duration) (btc.SizeEstimator, error) {
	baseCache, err := tools.NewMemCache[int](memcache.WithTtl(ttl))
	if err != nil {
		return nil, err
	}
	segwitCache, err := tools.NewMemCache[int](memcache.WithTtl(ttl))
	if err != nil {
		return nil, err
	}

	return &cachedSizeEstimator{
		mu:     new(sync.Mutex),
		base:   baseCache,
		segwit: segwitCache,
	}, nil
}

// AddUtxos adds a list of UTXOs and their estimated base and segwit size
func (estimator *cachedSizeEstimator) AddUtxos(base, segwit int, utxos ...btc.UTXO) {
	estimator.mu.Lock()
	defer estimator.mu.Unlock()

	for _, utxo := range utxos {
		estimator.base.Set(utxo.String(), base)
		estimator.segwit.Set(utxo.String(), segwit)
	}
}

func (estimator *cachedSizeEstimator) FetchSize(utxo btc.UTXO) (int, int, error) {
	estimator.mu.Lock()
	defer estimator.mu.Unlock()

	base, ok := estimator.base.Get(utxo.String())
	if !ok {
		return 0, 0, fmt.Errorf("unknown utxo = %v", utxo.String())
	}
	segwit, ok := estimator.segwit.Get(utxo.String())
	if !ok {
		return 0, 0, fmt.Errorf("unknown utxo = %v", utxo.String())
	}
	return base, segwit, nil
}

func (estimator *cachedSizeEstimator) EstimateTxWeight(tx *wire.MsgTx) (int, error) {
	totalBase, totalSegwit := tx.SerializeSizeStripped(), 0
	legacy := 0
	for _, input := range tx.TxIn {
		key := input.PreviousOutPoint.String()
		base, ok := estimator.base.Get(key)
		if !ok {
			return 0, fmt.Errorf("unknown utxo = %v", key)
		}
		totalBase += base
		if base != 0 {
			legacy++
		}
		segwit, ok := estimator.segwit.Get(key)
		if !ok {
			return 0, fmt.Errorf("unknown utxo = %v", key)
		}
		totalSegwit += segwit
	}

	// Additional 2 weight units for segwit marker + flag if tx has any witness input
	if totalSegwit > 0 {
		totalSegwit += 2
	}

	// When including both legacy and segwit inputs
	if totalSegwit > 0 && legacy > 0 {
		totalSegwit += legacy
	}

	return totalBase*4 + totalSegwit, nil
}

// EstimateTxVirtualSize returns the estimated size of the given transaction. It would be an upperbound, and usually the
// fees might be a few bytes less. It assumes the tx is not signed at all. It would return an error if one of the utxo
// is unknown in terms of signature size.
func (estimator *cachedSizeEstimator) EstimateTxVirtualSize(tx *wire.MsgTx) (int, error) {
	weight, err := estimator.EstimateTxWeight(tx)
	if err != nil {
		return 0, err
	}
	return (weight + 3) / blockchain.WitnessScaleFactor, nil
}

// CachedFetcher implements the `txscript.PrevOutputFetcher` using an memcache and an indexer. This ensures the data
// will be cached for a certain amount of time (usually longer than the time taking to mine) and an indexer will be used
// to fetch missing data on the fly.
type CachedFetcher struct {
	cache   memcache.Cache[wire.TxOut]
	indexer btc.IndexerClient
}

func NewCachedFetcher(ttl time.Duration, indexer btc.IndexerClient) (*CachedFetcher, error) {
	cache, err := tools.NewMemCache[wire.TxOut](memcache.WithTtl(ttl))
	if err != nil {
		return nil, err
	}

	return &CachedFetcher{
		cache:   cache,
		indexer: indexer,
	}, nil
}

func (fetcher *CachedFetcher) AddPrevOut(op wire.OutPoint, txOut *wire.TxOut) {
	key := op.String()
	fetcher.cache.Set(key, *txOut)
}

func (fetcher *CachedFetcher) AddUtxo(pkScript []byte, utxos ...btc.UTXO) {
	for _, utxo := range utxos {
		key := utxo.String()
		txOut := wire.TxOut{
			Value:    utxo.Amount,
			PkScript: pkScript,
		}
		fetcher.cache.Set(key, txOut)
	}
}

func (fetcher *CachedFetcher) FetchPrevOutput(outpoint wire.OutPoint) *wire.TxOut {
	key := outpoint.String()
	val, ok := fetcher.cache.Get(key)
	if !ok {
		if fetcher.indexer != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			tx, err := fetcher.indexer.GetTx(ctx, outpoint.Hash.String())
			if err != nil {
				return nil
			}
			if len(tx.VOUTs) > int(outpoint.Index) {
				out := tx.VOUTs[outpoint.Index]
				pkScript, err := hex.DecodeString(out.ScriptPubKey)
				if err != nil {
					return nil
				}
				return &wire.TxOut{
					Value:    int64(out.Value),
					PkScript: pkScript,
				}
			}
		}
		return nil
	}
	return &val
}
