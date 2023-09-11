package btc

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"go.uber.org/zap"
)

const (
	DefaultElectrsIndexerURL = "http://0.0.0.0:30000"

	DefaultRetryInterval = 5 * time.Second
)

// IndexerClient provides some rpc functions which usually cannot be achieved by the standard bitcoin json-rpc methods.
// The actual implementation will normally rely on an indexer, or third-party API sitting in front of the indexer.
type IndexerClient interface {

	// GetUTXOs return all utxos of the given address.
	GetUTXOs(ctx context.Context, address btcutil.Address) (UTXOs, error)

	// GetTipBlockHeight returns the tip block height.
	GetTipBlockHeight(ctx context.Context) (uint64, error)

	// GetTx returns the tx details with the given id.
	GetTx(ctx context.Context, txid string) (Transaction, error)

	// SubmitTx submits the given tx to the blockchain. The tx needs to be signed.
	SubmitTx(ctx context.Context, tx *wire.MsgTx) error
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
			// todo : check error message if failed

			return fmt.Errorf("failed to send transaction: %s", data)
		}
		return nil
	})
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
