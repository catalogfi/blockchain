package btc

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	MempoolFeeAPI = "https://mempool.space/api/v1/fees/recommended"

	MempoolFeeAPITestnet = "https://mempool.space/testnet/api/v1/fees/recommended"

	BlockstreamAPI = "https://blockstream.info/api/fee-estimates"
)

var (
	// RedeemHtlcRefundSigScriptSize is an estimate of the sigScript size when refunding an htlc script
	// stack number + stack size * 4 + signature + public key + script size
	RedeemHtlcRefundSigScriptSize = 1 + 4 + 73 + 33 + NormalHtlcSize

	// RedeemHtlcRedeemSigScriptSize is an estimate of the sigScript size when redeeming an htlc script
	// stack number + stack size * 5 + signature + public key + secret + script size
	RedeemHtlcRedeemSigScriptSize = func(secretSize int) int {
		return 1 + 5 + 73 + 33 + secretSize + +1 + NormalHtlcSize
	}

	// RedeemMultisigSigScriptSize is an estimate of the sigScript size from an 2-of-2 multisig script
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

// TxVirtualSize returns the virtual size of a transaction.
func TxVirtualSize(tx *wire.MsgTx) int {
	size := tx.SerializeSize()
	baseSize := tx.SerializeSizeStripped()
	swSize := size - baseSize
	return baseSize + (swSize+3)/blockchain.WitnessScaleFactor
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

		var res FeeSuggestion
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return FeeSuggestion{}, err
		}
		f.last = res
		f.lastTime = time.Now()
		return res, nil
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
		return FeeSuggestion{1, 1, 1, 1, 1}, nil
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
