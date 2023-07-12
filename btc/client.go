package btc

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"go.uber.org/zap"
)

var ErrNilResult = errors.New("nil result")

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

// Client to interact with the bitcoin network. We'll follow the standard bitcoind rpc interface if that meet our
// requirements.
type Client interface {
	IndexerClient

	// Net returns the network params.
	Net() *chaincfg.Params

	// LatestBlock returns the height and hash of the latest block.
	LatestBlock(ctx context.Context) (int64, string, error)

	// SubmitTx to the bitcoin network
	SubmitTx(ctx context.Context, tx wire.MsgTx) error

	// GetRawTransaction returns the raw transaction of the given hash.
	GetRawTransaction(ctx context.Context, txhash []byte) (btcjson.TxRawResult, error)

	// GetBlockByHeight returns the block detail of the given height.
	GetBlockByHeight(ctx context.Context, height int64) (*btcjson.GetBlockVerboseResult, error)

	// GetBlockByHash returns the block detail with the given hash.
	GetBlockByHash(ctx context.Context, hash string) (*btcjson.GetBlockVerboseResult, error)

	// GetTxOut returns details about an unspent transaction output
	GetTxOut(ctx context.Context, hash string, vout uint32) (*btcjson.GetTxOutResult, error)
}

type client struct {
	opts          ClientOptions
	logger        *zap.Logger
	httpClient    *http.Client
	indexerClient IndexerClient
}

func NewClient(opts ClientOptions, logger *zap.Logger, indexerClient IndexerClient) Client {
	httpClient := new(http.Client)
	httpClient.Timeout = opts.Timeout

	return &client{
		opts:          opts,
		logger:        logger,
		httpClient:    httpClient,
		indexerClient: indexerClient,
	}
}

func (client *client) Net() *chaincfg.Params {
	return client.opts.Net
}

// GetUTXOs gets all utxos of the given address. This implementation is specified to quickNode.
func (client *client) GetUTXOs(ctx context.Context, address btcutil.Address) (UTXOs, error) {
	if client.indexerClient == nil {
		return nil, fmt.Errorf("GetUTXOs not supported")
	}
	return client.indexerClient.GetUTXOs(ctx, address)
}

// LatestBlock returns the height and hash of the longest blockchain.
func (client *client) LatestBlock(ctx context.Context) (int64, string, error) {
	var resp btcjson.GetBlockChainInfoResult
	if err := client.send(ctx, &resp, "getblockchaininfo"); err != nil {
		return 0, "", fmt.Errorf("get block count: %w", err)
	}
	return int64(resp.Blocks), resp.BestBlockHash, nil
}

// SubmitTx to the Bitcoin network.
func (client *client) SubmitTx(ctx context.Context, tx wire.MsgTx) error {
	var serial bytes.Buffer
	if err := tx.Serialize(&serial); err != nil {
		return fmt.Errorf("bad tx: %w", err)
	}
	resp := ""
	if err := client.send(ctx, &resp, "sendrawtransaction", hex.EncodeToString(serial.Bytes())); err != nil {
		return fmt.Errorf("bad \"sendrawtransaction\": %w", err)
	}
	return nil
}

func (client *client) GetRawTransaction(ctx context.Context, txhash []byte) (btcjson.TxRawResult, error) {
	resp := btcjson.TxRawResult{}
	hash := chainhash.Hash{}
	copy(hash[:], txhash)
	if err := client.send(ctx, &resp, "getrawtransaction", hash.String(), 1); err != nil {
		return resp, fmt.Errorf("bad \"getrawtransaction\": %w", err)
	}
	return resp, nil
}

func (client *client) GetBlockByHeight(ctx context.Context, height int64) (*btcjson.GetBlockVerboseResult, error) {
	resp := ""
	if err := client.send(ctx, &resp, "getblockhash", height); err != nil {
		return nil, fmt.Errorf("bad \"getblockhash\": %w", err)
	}
	block := btcjson.GetBlockVerboseResult{}
	if err := client.send(ctx, &block, "getblock", resp); err != nil {
		return nil, fmt.Errorf("bad \"getblock\": %w", err)
	}
	return &block, nil
}

func (client *client) GetBlockByHash(ctx context.Context, hash string) (*btcjson.GetBlockVerboseResult, error) {
	block := btcjson.GetBlockVerboseResult{}
	if err := client.send(ctx, &block, "getblock", hash); err != nil {
		return nil, fmt.Errorf("bad \"getblock\": %w", err)
	}
	return &block, nil
}

func (client *client) GetTxOut(ctx context.Context, hash string, vout uint32) (*btcjson.GetTxOutResult, error) {
	result := btcjson.GetTxOutResult{}
	if err := client.send(ctx, &result, "gettxout", hash, vout); err != nil {
		return nil, fmt.Errorf("bad \"gettxout\": %w", err)
	}
	return &result, nil
}

func (client *client) send(ctx context.Context, resp interface{}, method string, params ...interface{}) error {
	// Encode the request.
	data, err := encodeRequest(method, params)
	if err != nil {
		return err
	}

	return retry(client.logger, ctx, client.opts.TimeoutRetry, func() error {
		// Create request and add basic authentication headers. The context is
		// not attached to the request, and instead we all each attempt to run
		// for the timeout duration, and we keep attempting until success, or
		// the context is done.
		req, err := http.NewRequest("POST", client.opts.Host, bytes.NewBuffer(data))
		if err != nil {
			return fmt.Errorf("building http request: %v", err)
		}
		req.SetBasicAuth(client.opts.User, client.opts.Password)

		// Send the request and decode the response.
		res, err := client.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("sending http request: %v", err)
		}
		defer res.Body.Close()

		if err := decodeResponse(resp, res.Body); err != nil {
			return NewNoRetryError(err)
		}
		return nil
	})
}

func encodeRequest(method string, params []interface{}) ([]byte, error) {
	rawParams, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("encoding params: %v", err)
	}
	req := struct {
		Version string          `json:"version"`
		ID      int             `json:"id"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params"`
	}{
		Version: "2.0",
		ID:      rand.Int(),
		Method:  method,
		Params:  rawParams,
	}
	rawReq, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encoding request: %v", err)
	}
	return rawReq, nil
}

func decodeResponse(resp interface{}, r io.Reader) error {
	res := struct {
		Version string           `json:"version"`
		ID      int              `json:"id"`
		Result  *json.RawMessage `json:"result"`
		Error   *json.RawMessage `json:"error"`
	}{}
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return fmt.Errorf("decoding response: %v", err)
	}
	if res.Error != nil {
		return errors.New(string(*res.Error))
	}
	if res.Result == nil {
		return ErrNilResult
	}
	if err := json.Unmarshal(*res.Result, resp); err != nil {
		return fmt.Errorf("decoding result: %v", err)
	}
	return nil
}

func retry(logger *zap.Logger, ctx context.Context, dur time.Duration, f func() error) error {
	ticker := time.NewTicker(dur)
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

type MockClient struct {
	FuncNet               func() *chaincfg.Params
	FuncGetUTXOs          func(context.Context, btcutil.Address) (UTXOs, error)
	FuncLatestBlock       func(context.Context) (int64, string, error)
	FuncSubmitTx          func(context.Context, wire.MsgTx) error
	FuncGetRawTransaction func(context.Context, []byte) (btcjson.TxRawResult, error)
	FuncGetBlockByHeight  func(context.Context, int64) (*btcjson.GetBlockVerboseResult, error)
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) Net() *chaincfg.Params {
	return m.FuncNet()
}

func (m *MockClient) GetUTXOs(ctx context.Context, address btcutil.Address) (UTXOs, error) {
	return m.FuncGetUTXOs(ctx, address)
}

func (m *MockClient) LatestBlock(ctx context.Context) (int64, string, error) {
	return m.FuncLatestBlock(ctx)
}

func (m *MockClient) SubmitTx(ctx context.Context, tx wire.MsgTx) error {
	return m.FuncSubmitTx(ctx, tx)
}

func (m *MockClient) GetRawTransaction(ctx context.Context, txhash []byte) (btcjson.TxRawResult, error) {
	return m.FuncGetRawTransaction(ctx, txhash)
}

func (m *MockClient) GetBlockByHeight(ctx context.Context, height int64) (*btcjson.GetBlockVerboseResult, error) {
	return m.FuncGetBlockByHeight(ctx, height)
}
