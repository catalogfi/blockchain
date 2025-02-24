package btc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
)

const (
	MempoolFeeApiMainnet = "https://mempool.space/api/v1/fees/recommended"

	MempoolFeeApiTestnet3 = "https://mempool.space/testnet/api/v1/fees/recommended"

	MempoolFeeApiTestnet4 = "https://mempool.space/testnet4/api/v1/fees/recommended"

	BlockstreamApiMainnet = "https://blockstream.info/api/fee-estimates"
)

var (
	// BaseSizeP2PKH is the worst case (largest) serialize size of a transaction input script that redeems a compressed
	// P2PKH output. This assumes we always use a low-s value signature. The signature size are usually 71 (60%) or
	// 72(40%), with a very small chance of 70 or 69 (<1%). It is calculated as :
	// sigLength(1) + sig(72) + pubKeyLength(1) + compressedPubKey(33)
	BaseSizeP2PKH = 1 + 72 + 1 + 33

	BaseSizeP2WPKH = 0

	BaseSizeP2TR = 0

	SegwitSizeP2PKH = 0

	// SegwitSizeP2WPKH is the worst case weight of a witness for spending P2WPKH outputs. It is calculated as :
	// number of items(1) + sigLength(1) + sig(72) + pubKeyLength(1) + compressedPubKey(33)
	SegwitSizeP2WPKH = 1 + 1 + 72 + 1 + 33

	// SegwitSizeP2TR is the worst case weight of a witness for spending P2TR outputs. It is calculated as :
	// number of items(1) + sigLength(1) + sig(64)
	SegwitSizeP2TR = 1 + 1 + 64

	// SegwitSizeP2trDefault is the witness size when the schnorr signature is signed using the default sighash flag.
	SegwitSizeP2trDefault = 1 + 1 + 64
)

// NewFetcher initializes a txscript.MultiPrevOutFetcher with the given utxos and script.
func NewFetcher(script []byte, utxos ...UTXO) (*txscript.MultiPrevOutFetcher, error) {
	fetcher := txscript.NewMultiPrevOutFetcher(nil)
	for _, utxo := range utxos {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}
		fetcher.AddPrevOut(wire.OutPoint{
			Hash:  *hash,
			Index: utxo.Vout,
		}, wire.NewTxOut(utxo.Amount, script))
	}
	return fetcher, nil
}

// AddUtxosToFetcher adds the list of utxos to the fetcher.
func AddUtxosToFetcher(fetcher *txscript.MultiPrevOutFetcher, script []byte, utxos ...UTXO) error {
	for _, utxo := range utxos {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return err
		}
		fetcher.AddPrevOut(wire.OutPoint{
			Hash:  *hash,
			Index: utxo.Vout,
		}, wire.NewTxOut(utxo.Amount, script))
	}

	return nil
}

// TotalFee returns the total amount fees used by the given tx.
func TotalFee(tx *wire.MsgTx, fetcher txscript.PrevOutputFetcher) int {
	fees := int64(0)
	for _, in := range tx.TxIn {
		output := fetcher.FetchPrevOutput(in.PreviousOutPoint)
		fees += output.Value
	}
	for _, out := range tx.TxOut {
		fees -= out.Value
	}
	return int(fees)
}

// SizeEstimator collects estimated signature size of UTXOs, it then can be used to estimate the transaction size for
// fee purpose.
type SizeEstimator struct {
	mu            *sync.Mutex
	baseSizeMap   map[string]int
	segwitSizeMap map[string]int
}

// NewEmptySizeEstimator returns an empty SizeEstimator
func NewEmptySizeEstimator() *SizeEstimator {
	return &SizeEstimator{
		mu:            new(sync.Mutex),
		baseSizeMap:   make(map[string]int),
		segwitSizeMap: make(map[string]int),
	}
}

// NewSizeEstimator returns an SizeEstimator with some preload UTXOs.
func NewSizeEstimator(base, segwit int, utxos ...UTXO) *SizeEstimator {
	baseSizeMap := make(map[string]int)
	segwitSizeMap := make(map[string]int)
	if base != 0 || segwit != 0 {
		for _, utxo := range utxos {
			baseSizeMap[utxo.String()] = base
			segwitSizeMap[utxo.String()] = segwit
		}
	}

	return &SizeEstimator{
		mu:            new(sync.Mutex),
		baseSizeMap:   baseSizeMap,
		segwitSizeMap: segwitSizeMap,
	}
}

// NewSizeEstimatorOfAddrType returns an SizeEstimator basing on the provided address type.
func NewSizeEstimatorOfAddrType(addrType waddrmgr.AddressType, utxos ...UTXO) *SizeEstimator {
	switch addrType {
	case waddrmgr.PubKeyHash:
		return NewSizeEstimator(BaseSizeP2PKH, SegwitSizeP2PKH, utxos...)
	case waddrmgr.WitnessPubKey:
		return NewSizeEstimator(BaseSizeP2WPKH, SegwitSizeP2WPKH, utxos...)
	case waddrmgr.TaprootPubKey:
		return NewSizeEstimator(BaseSizeP2TR, SegwitSizeP2TR, utxos...)
	default:
		panic(fmt.Sprintf("unknown address type: %v", addrType))
	}
}

// AddUtxos adds a list of UTXOs and their estimated base and segwit size
func (estimator *SizeEstimator) AddUtxos(base, segwit int, utxos ...UTXO) {
	estimator.mu.Lock()
	defer estimator.mu.Unlock()

	for _, utxo := range utxos {
		estimator.baseSizeMap[utxo.String()] = base
		estimator.segwitSizeMap[utxo.String()] = segwit
	}
}

func (estimator *SizeEstimator) FetchSize(utxo UTXO) (int, int, error) {
	estimator.mu.Lock()
	defer estimator.mu.Unlock()

	base, ok := estimator.baseSizeMap[utxo.String()]
	if !ok {
		return 0, 0, fmt.Errorf("unknown utxo = %v", utxo.String())
	}
	segwit, ok := estimator.segwitSizeMap[utxo.String()]
	if !ok {
		return 0, 0, fmt.Errorf("unknown utxo = %v", utxo.String())
	}
	return base, segwit, nil
}

func (estimator *SizeEstimator) EstimateTxWeight(tx *wire.MsgTx) (int, error) {
	totalBase, totalSegwit := tx.SerializeSizeStripped(), 0
	legacy := 0
	for _, input := range tx.TxIn {
		key := input.PreviousOutPoint.String()
		base, ok := estimator.baseSizeMap[key]
		if !ok {
			return 0, fmt.Errorf("unknown utxo = %v", key)
		}
		totalBase += base
		if base != 0 {
			legacy++
		}
		segwit, ok := estimator.segwitSizeMap[key]
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
		totalBase += legacy
	}

	return totalBase*4 + totalSegwit, nil
}

// EstimateTxVirtualSize returns the estimated size of the given transaction. It would be an upperbound, and usually the
// fees might be a few bytes less. It assumes the tx is not signed at all. It would return an error if one of the utxo
// is unknown in terms of signature size.
func (estimator *SizeEstimator) EstimateTxVirtualSize(tx *wire.MsgTx) (int, error) {
	weight, err := estimator.EstimateTxWeight(tx)
	if err != nil {
		return 0, err
	}
	return (weight + 3) / blockchain.WitnessScaleFactor, nil
}

// SatoshiPerKb is the fee rate in satoshi per kilobyte. This will be the default measure of fee rate in this package.
type SatoshiPerKb int64

// NewSatoshiPerKb returns a new SatoshiPerKb from the given fee and virtual size.
func NewSatoshiPerKb(fee int64, vsize int) SatoshiPerKb {
	return SatoshiPerKb(fee * 1000 / int64(vsize))
}

func (rate SatoshiPerKb) Int() int {
	return int(rate)
}

func (rate SatoshiPerKb) ToSatoshiPerByte() mempool.SatoshiPerByte {
	return mempool.SatoshiPerByte(float64(rate) / 1000)
}

type FeeLevel string

var (
	FeeLow    FeeLevel = "low"
	FeeMedium FeeLevel = "medium"
	FeeHigh   FeeLevel = "high"
)

type FeeSuggestion struct {
	Low    SatoshiPerKb `json:"low"`
	Medium SatoshiPerKb `json:"medium"`
	High   SatoshiPerKb `json:"high"`
}

// Fee will return the fee rate of the given fee level.
func (feeSuggestion FeeSuggestion) Fee(level FeeLevel) SatoshiPerKb {
	switch level {
	case FeeLow:
		return feeSuggestion.Low
	case FeeMedium:
		return feeSuggestion.Medium
	case FeeHigh:
		return feeSuggestion.High
	default:
		panic(fmt.Sprintf("unknown fee level: %v", level))
	}
}

type FeeEstimator interface {
	FeeSuggestion() (FeeSuggestion, error)
}

type mempoolFeeEstimator struct {
	params *chaincfg.Params
	url    string

	mu       *sync.Mutex
	last     FeeSuggestion
	lastTime time.Time
	ttl      time.Duration
}

func NewMempoolFeeEstimator(params *chaincfg.Params, url string, ttl time.Duration) FeeEstimator {
	switch params.Name {
	case chaincfg.MainNetParams.Name, chaincfg.TestNet3Params.Name:
		// do nothing
	default:
		panic("unsupported network")
	}
	return &mempoolFeeEstimator{
		params: params,
		url:    url,
		mu:     new(sync.Mutex),
		ttl:    ttl,
	}
}

func (f *mempoolFeeEstimator) FeeSuggestion() (FeeSuggestion, error) {
	if f.url == "" {
		return FeeSuggestion{}, errors.New("url is empty")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.lastTime.IsZero() || time.Since(f.lastTime) >= f.ttl {
		resp, err := http.Get(f.url)
		if err != nil {
			return FeeSuggestion{}, err
		}
		defer resp.Body.Close()

		res := struct {
			Minimum float64 `json:"minimumFee"`
			Economy float64 `json:"economyFee"`
			Low     float64 `json:"hourFee"`
			Medium  float64 `json:"halfHourFee"`
			High    float64 `json:"fastestFee"`
		}{}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return FeeSuggestion{}, err
		}
		f.last = FeeSuggestion{
			Low:    SatoshiPerKb(res.Economy * 1000),
			Medium: SatoshiPerKb(res.Medium * 1000),
			High:   SatoshiPerKb(res.High * 1000),
		}
		f.lastTime = time.Now()
		return f.last, nil
	}
	return f.last, nil
}

type blockstreamFeeEstimator struct {
	params *chaincfg.Params
	url    string

	mu       *sync.Mutex
	last     FeeSuggestion
	lastTime time.Time
	ttl      time.Duration
}

func NewBlockstreamFeeEstimator(params *chaincfg.Params, url string, ttl time.Duration) FeeEstimator {
	return &blockstreamFeeEstimator{
		params: params,
		url:    url,
		mu:     new(sync.Mutex),
		ttl:    ttl,
	}
}

func (f *blockstreamFeeEstimator) FeeSuggestion() (FeeSuggestion, error) {
	if f.url != "" {
		f.mu.Lock()
		defer f.mu.Unlock()

		if f.lastTime.IsZero() || time.Since(f.lastTime) >= f.ttl {
			resp, err := http.Get(f.url)
			if err != nil {
				return FeeSuggestion{}, err
			}
			defer resp.Body.Close()

			fees := map[string]float64{}
			if err := json.NewDecoder(resp.Body).Decode(&fees); err != nil {
				return FeeSuggestion{}, err
			}
			if len(fees) == 0 {
				return FeeSuggestion{2, 2, 2}, nil
			}

			feerates := FeeSuggestion{
				Low:    SatoshiPerKb(math.Ceil(fees["6"] * 1000)),
				Medium: SatoshiPerKb(math.Ceil(fees["3"] * 1000)),
				High:   SatoshiPerKb(math.Ceil(fees["1"] * 1000)),
			}

			f.last = feerates
			f.lastTime = time.Now()
			return feerates, nil
		}
		return f.last, nil
	} else {
		return FeeSuggestion{1e3, 1e3, 1e3}, nil
	}
}

type fixFeeEstimator struct {
	fee SatoshiPerKb
}

func NewFixFeeEstimator(fee SatoshiPerKb) FeeEstimator {
	return fixFeeEstimator{
		fee: fee,
	}
}

func (f fixFeeEstimator) FeeSuggestion() (FeeSuggestion, error) {
	return FeeSuggestion{
		Low:    f.fee,
		Medium: f.fee,
		High:   f.fee,
	}, nil
}
