package btc

import (
	"encoding/json"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
)

const (
	MempoolFeeAPI = "https://mempool.space/api/v1/fees/recommended"

	BlockstreamAPI = "https://blockstream.info/api/fee-estimates"
)

var (
	// stack number + stack size * 4 + signature + public key + script size
	RedeemHtlcRefundSigScriptSize = 1 + 4 + 73 + 33 + 89

	// stack number + stack size * 5 + signature + public key + secret + script size
	RedeemHtlcRedeemSigScriptSize = func(secretSize int) int {
		return 1 + 5 + 73 + 33 + secretSize + +1 + 89
	}

	// stack number + stack size * 4 + signature * 2 + script size
	RedeemMultisigSigScriptSize = 1 + 4 + 73*2 + 71
)

// EstimateVirtualSize will return an estimate virtual size of the given unsigned tx. The extraBaseSize will be the signature
// size of all legacy type utxos and extraSegwitSize will be signature size of all segwit utxos.
func EstimateVirtualSize(tx *wire.MsgTx, extraBaseSize, extraSegwitSize int) int {
	baseSize := tx.SerializeSizeStripped()
	baseSize += extraBaseSize
	swSize := 2 + extraSegwitSize
	return baseSize + (swSize+3)/blockchain.WitnessScaleFactor
}

// EstimateUtxoSize returns the estimated signature size for a few fixed-length utxo types. `n` is the number of
// utxos from the transaction input.
func EstimateUtxoSize(address btcutil.Address, n int) int {
	switch address.(type) {
	case *btcutil.AddressPubKeyHash:
		return txsizes.RedeemP2PKHSigScriptSize * n
	case *btcutil.AddressWitnessPubKeyHash:
		return txsizes.RedeemP2WPKHInputWitnessWeight * n
	case *btcutil.AddressTaproot:
		return txsizes.RedeemP2TRInputWitnessWeight * n
	default:
		return 0
	}
}

// TxVirtualSize returns the virtual size of a transaction
func TxVirtualSize(tx *wire.MsgTx) int {
	size := tx.SerializeSize()
	baseSize := tx.SerializeSizeStripped()
	swSize := size - baseSize
	return baseSize + (swSize+3)/blockchain.WitnessScaleFactor
}

type FeeSuggestion struct {
	Minimum int `json:"minimumFee"`
	Economy int `json:"economyFee"`
	Low     int `json:"hourFee"`
	Medium  int `json:"halfHourFee"`
	High    int `json:"fastestFee"`
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
	return &mempoolFeeEstimator{
		params: params,
		url:    url,
		mu:     new(sync.Mutex),
		ttl:    ttl,
	}
}

func (f *mempoolFeeEstimator) FeeSuggestion() (FeeSuggestion, error) {
	if f.params.Name == "mainnet" && f.url != "" {
		f.mu.Lock()
		defer f.mu.Unlock()

		if f.lastTime.IsZero() || time.Since(f.lastTime) >= f.ttl {
			resp, err := http.Get(f.url)
			if err != nil {
				return FeeSuggestion{}, err
			}
			defer resp.Body.Close()

			var res FeeSuggestion
			if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
				return FeeSuggestion{}, err
			}
			f.last = res
			f.lastTime = time.Now()
			return res, nil
		}
		return f.last, nil
	} else {
		return FeeSuggestion{
			2, 2, 2, 2, 2,
		}, nil
	}
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
	if f.params.Name == "mainnet" && f.url != "" {
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
				return FeeSuggestion{2, 2, 2, 2, 2}, nil
			}

			feerates := FeeSuggestion{
				Minimum: int(math.Ceil(fees["504"])),
				Economy: int(math.Ceil(fees["144"])),
				Low:     int(math.Ceil(fees["6"])),
				Medium:  int(math.Ceil(fees["3"])),
				High:    int(math.Ceil(fees["1"])),
			}

			f.last = feerates
			f.lastTime = time.Now()
			return feerates, nil
		}
		return f.last, nil
	} else {
		return FeeSuggestion{2, 2, 2, 2, 2}, nil
	}
}

type fixFeeEstimator struct {
	fee int
}

func (f fixFeeEstimator) FeeSuggestion() (FeeSuggestion, error) {
	return FeeSuggestion{
		Minimum: f.fee,
		Economy: f.fee,
		Low:     f.fee,
		Medium:  f.fee,
		High:    f.fee,
	}, nil
}

func NewFixFeeEstimator(fee int) FeeEstimator {
	return fixFeeEstimator{
		fee: fee,
	}
}
