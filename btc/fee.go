package btc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/mempool"
)

const (
	MempoolFeeApiMainnet = "https://mempool.space/api/v1/fees/recommended"

	MempoolFeeApiTestnet3 = "https://mempool.space/testnet/api/v1/fees/recommended"

	MempoolFeeApiTestnet4 = "https://mempool.space/testnet4/api/v1/fees/recommended"

	BlockstreamApiMainnet = "https://blockstream.info/api/fee-estimates"
)

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
	case chaincfg.MainNetParams.Name, chaincfg.TestNet3Params.Name, chaincfg.TestNet4Params.Name, chaincfg.SigNetParams.Name:
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
