package btc

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
)

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
}

func NewFeeEstimator(params *chaincfg.Params, url string) FeeEstimator {
	return &mempoolFeeEstimator{
		params: params,
		url:    url,
		mu:     new(sync.Mutex),
	}
}

func (f *mempoolFeeEstimator) FeeSuggestion() (FeeSuggestion, error) {
	if f.params.Name == "mainnet" && f.url != "" {
		f.mu.Lock()
		defer f.mu.Unlock()
		if f.lastTime.IsZero() || time.Since(f.lastTime) >= time.Minute {
			resp, err := http.Get(f.url)
			if err != nil {
				return FeeSuggestion{}, err
			}
			var res FeeSuggestion
			if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
				return FeeSuggestion{}, err
			}
			f.last = res
			return res, nil
		}
		return f.last, nil
	} else {
		return FeeSuggestion{
			2, 2, 2, 2, 2,
		}, nil
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
