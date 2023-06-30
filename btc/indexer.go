package btc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"go.uber.org/zap"
)

const (
	DefaultElectrsIndexerURL = "http://0.0.0.0:30000"
)

// IndexerClient provides some rpc functions which usually cannot be achieved by the standard bitcoin json-rpc methods.
// The actual implementation will normally rely on an indexer, or third-party API sitting in front of the indexer.
type IndexerClient interface {

	// GetUTXOs return all utxos of the given address.
	GetUTXOs(ctx context.Context, address btcutil.Address) (UTXOs, error)
}

type quickNodeClient struct {
	logger *zap.Logger
	url    string
}

func NewQuickNodeIndexerClient(logger *zap.Logger, url string) IndexerClient {
	return &quickNodeClient{logger: logger, url: url}
}

// GetUTXOs implements the IndexerClient using the quicknode bitcoin API.
// See https://www.quicknode.com/docs/bitcoin/bb_getutxos
func (client *quickNodeClient) GetUTXOs(ctx context.Context, address btcutil.Address) (UTXOs, error) {
	type UtxoJson struct {
		TxID   string `json:"txid"`
		Vout   uint32 `json:"vout"`
		Amount int64  `json:"value,string"`
	}
	utxoJsons := []UtxoJson{}

	// Encode the request
	data, err := encodeRequest("bb_getutxos", []interface{}{address.EncodeAddress()})
	if err != nil {
		return nil, err
	}

	// Send the request
	if err := retry(client.logger, ctx, 5*time.Second, func() error {
		// Create request
		req, err := http.NewRequest("POST", client.url, bytes.NewBuffer(data))
		if err != nil {
			return fmt.Errorf("building http request: %v", err)
		}

		// Send the request and decode the response.
		httpClient := new(http.Client)
		res, err := httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("sending http request: %v", err)
		}
		defer res.Body.Close()
		if err := decodeResponse(&utxoJsons, res.Body); err != nil {
			return fmt.Errorf("decoding http response: %v", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	utxos := make([]UTXO, len(utxoJsons))
	for i := range utxos {
		utxos[i] = UTXO{
			TxID:   utxoJsons[i].TxID,
			Vout:   utxoJsons[i].Vout,
			Amount: utxoJsons[i].Amount,
		}
	}
	return utxos, nil
}

type electrsIndexerClient struct {
	logger *zap.Logger
	url    string
}

func NewElectrsIndexerClient(logger *zap.Logger, url string) IndexerClient {
	return &electrsIndexerClient{logger: logger, url: url}
}

// GetUTXOs implements the IndexerClient basing on the electrs indexer API.
// See https://github.com/Blockstream/esplora/blob/master/API.md
func (client *electrsIndexerClient) GetUTXOs(ctx context.Context, address btcutil.Address) (UTXOs, error) {
	url, err := url.JoinPath(client.url, "address", address.EncodeAddress(), "utxo")
	if err != nil {
		return nil, err
	}

	// Send the request
	utxos := UTXOs{}
	if err := retry(client.logger, ctx, 5*time.Second, func() error {
		resp, err := http.Get(url)
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
