package btc

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"go.uber.org/zap"
)

const (
	DefaultElectrsIndexerURL = "http://0.0.0.0:30000"

	DefaultRetryInterval = 5 * time.Second
)

type Transaction struct {
	TxID     string    `json:"txid"`
	Version  int       `json:"version"`
	LockTime int       `json:"lock_time"`
	VINs     []VIN     `json:"vin"`
	VOUTs    []Prevout `json:"vout"`
	Status   Status    `json:"status"`
}

type VIN struct {
	TxID      string    `json:"txid"`
	Vout      int       `json:"vout"`
	Prevout   Prevout   `json:"prevout"`
	ScriptSig string    `json:"scriptsig"`
	Witness   *[]string `json:"witness"`
	Sequence  int       `json:"sequence"`
}

type Prevout struct {
	ScriptPubKeyType    string `json:"scriptpubkey_type"`
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int    `json:"value"`
}

type Status struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight uint64 `json:"block_height"`
	BlockHash   string `json:"block_hash"`
}

// IndexerClient provides some rpc functions which usually cannot be achieved by the standard bitcoin json-rpc methods.
// The actual implementation will normally rely on an indexer, or third-party API sitting in front of the indexer.
type IndexerClient interface {

	// GetAddressTxs returns the tx history of the given address.
	GetAddressTxs(ctx context.Context, address btcutil.Address, lastSeenTxid string) ([]Transaction, error)

	// GetUTXOs return all utxos of the given address.
	GetUTXOs(ctx context.Context, address btcutil.Address) (UTXOs, error)

	// GetTipBlockHeight returns the tip block height.
	GetTipBlockHeight(ctx context.Context) (uint64, error)

	// GetTx returns the tx details with the given id.
	GetTx(ctx context.Context, txid string) (Transaction, error)

	// SubmitTx submits the given tx to the blockchain. The tx needs to be signed.
	SubmitTx(ctx context.Context, tx *wire.MsgTx) error

	// FeeEstimate returns the estimate fees for different confirmation time.
	FeeEstimate(ctx context.Context) (FeeSuggestion, error)
}

type electrsIndexerClient struct {
	logger        *zap.Logger
	url           string
	retryInterval time.Duration
}

func NewElectrsIndexerClient(logger *zap.Logger, url string, retryInterval time.Duration) IndexerClient {
	return &electrsIndexerClient{
		logger:        logger,
		url:           url,
		retryInterval: retryInterval,
	}
}

func (client *electrsIndexerClient) GetAddressTxs(ctx context.Context, address btcutil.Address, lastSeenTxid string) ([]Transaction, error) {
	endpoint, err := url.JoinPath(client.url, "address", address.EncodeAddress(), "txs")
	if err != nil {
		return nil, err
	}
	if lastSeenTxid != "" {
		endpoint, err = url.JoinPath(client.url, "address", address.EncodeAddress(), "txs", "chain", lastSeenTxid)
		if err != nil {
			return nil, err
		}
	}

	// Send the request
	txs := []Transaction{}
	if err := retry(client.logger, ctx, client.retryInterval, func() error {
		resp, err := http.Get(endpoint)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Decode response
		if err := json.NewDecoder(resp.Body).Decode(&txs); err != nil {
			return fmt.Errorf("failed to decode UTXOs: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return txs, nil
}

// GetUTXOs implements the IndexerClient basing on the electrs indexer API.
// See https://github.com/Blockstream/esplora/blob/master/API.md
func (client *electrsIndexerClient) GetUTXOs(ctx context.Context, address btcutil.Address) (UTXOs, error) {
	endpoint, err := url.JoinPath(client.url, "address", address.EncodeAddress(), "utxo")
	if err != nil {
		return nil, err
	}

	// Send the request
	utxos := UTXOs{}
	if err := retry(client.logger, ctx, client.retryInterval, func() error {
		resp, err := http.Get(endpoint)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Decode response
		if err := json.NewDecoder(resp.Body).Decode(&utxos); err != nil {
			return fmt.Errorf("failed to decode UTXOs: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return utxos, nil
}

func (client *electrsIndexerClient) GetTipBlockHeight(ctx context.Context) (uint64, error) {
	endpoint, err := url.JoinPath(client.url, "blocks", "tip", "height")
	if err != nil {
		return 0, err
	}

	// Send the request
	var height uint64
	if err := retry(client.logger, ctx, client.retryInterval, func() error {
		resp, err := http.Get(endpoint)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Decode response
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		height, err = strconv.ParseUint(string(data), 10, 64)
		return err
	}); err != nil {
		return 0, err
	}

	return height, nil
}

func (client *electrsIndexerClient) GetTx(ctx context.Context, txid string) (Transaction, error) {
	endpoint, err := url.JoinPath(client.url, "tx", txid)
	if err != nil {
		return Transaction{}, err
	}

	// Send the request
	var tx Transaction
	if err := retry(client.logger, ctx, client.retryInterval, func() error {
		resp, err := http.Get(endpoint)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Decode response
		if err := json.NewDecoder(resp.Body).Decode(&tx); err != nil {
			return fmt.Errorf("failed to decode UTXOs: %w", err)
		}
		return nil
	}); err != nil {
		return Transaction{}, err
	}

	return tx, nil
}

func (client *electrsIndexerClient) SubmitTx(ctx context.Context, tx *wire.MsgTx) error {
	endpoint, err := url.JoinPath(client.url, "tx")
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return err
	}
	strBuffer := bytes.NewBufferString(hex.EncodeToString(buf.Bytes()))

	// Send the request
	return retry(client.logger, ctx, client.retryInterval, func() error {
		resp, err := http.Post(endpoint, "application/text", strBuffer)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			errMessage := strings.ToLower(string(data))
			switch {
			case strings.Contains(errMessage, "transaction already in block chain"):
				return NewNoRetryError(ErrAlreadyInChain)
			case strings.Contains(errMessage, "bad-txns-inputs-missingorspent"):
				return NewNoRetryError(ErrTxInputsMissingOrSpent)
			case strings.Contains(errMessage, "txn-mempool-conflict"):
				return NewNoRetryError(ErrMempoolConflict)
			default:
				return NewNoRetryError(errors.New(string(data)))
			}
		}
		return nil
	})
}

func (client *electrsIndexerClient) FeeEstimate(ctx context.Context) (FeeSuggestion, error) {
	endpoint, err := url.JoinPath(client.url, "fee-estimates")
	if err != nil {
		return FeeSuggestion{}, err
	}

	// Send the request
	var fees FeeSuggestion
	err = retry(client.logger, ctx, client.retryInterval, func() error {
		resp, err := http.Get(endpoint)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		feerates := map[string]float64{}
		if err := json.NewDecoder(resp.Body).Decode(&fees); err != nil {
			return err
		}
		if len(feerates) == 0 {
			return NewNoRetryError(fmt.Errorf("not enough data"))
		}

		fees = FeeSuggestion{
			Minimum: int(math.Ceil(feerates["504"])),
			Economy: int(math.Ceil(feerates["144"])),
			Low:     int(math.Ceil(feerates["6"])),
			Medium:  int(math.Ceil(feerates["3"])),
			High:    int(math.Ceil(feerates["1"])),
		}

		return nil
	})
	return fees, err
}

func retry(logger *zap.Logger, ctx context.Context, dur time.Duration, f func() error) error {
	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	err := f()
	for err != nil {
		// Skip retrying if it's a `NoRetryError`
		var e *NoRetryError
		if errors.As(err, &e) {
			return err
		}

		logger.Debug("retrying", zap.Any("error", err.Error()))
		select {
		case <-ctx.Done():
			return fmt.Errorf("%v: %v", ctx.Err(), err)
		case <-ticker.C:
			err = f()
		}
	}
	return nil
}

type NoRetryError struct {
	err error
}

func NewNoRetryError(err error) *NoRetryError {
	return &NoRetryError{
		err: err,
	}
}

func (err *NoRetryError) Error() string {
	if err.err != nil {
		return err.err.Error()
	}
	return ""
}

func (err *NoRetryError) Unwrap() error {
	return err.err
}
