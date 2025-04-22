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
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

var (
	ErrTxNotFound = errors.New("no such mempool or blockchain transaction")

	// replaced by ErrAlreadyInUtxoSet in v28.0
	// see https://github.com/bitcoin/bitcoin/blob/master/doc/release-notes/release-notes-28.0.md
	// ErrAlreadyInChain = errors.New("transaction already in block chain")

	ErrAlreadyInUtxoSet = errors.New("Transaction outputs already in utxo set")

	ErrTxInputsMissingOrSpent = errors.New("bad-txns-inputs-missingorspent")

	ErrMempoolConflict = errors.New("txn-mempool-conflict")

	ErrTxNotInMempool = errors.New("Transaction not in mempool")
)

// Client to interact with the Bitcoin network. It's implementation uses standard bitcoind JSON-RPC behind the scene.
type Client interface {

	// Net returns the network params.
	Net() *chaincfg.Params

	// GetBlockchainInfo returns various state info regarding blockchain processing. Can be used to retrieve the best
	// block height and hash.
	GetBlockchainInfo(ctx context.Context) (*btcjson.GetBlockChainInfoResult, error)

	// SubmitTx to the Bitcoin network.
	SubmitTx(ctx context.Context, tx *wire.MsgTx) error

	// GetRawTransaction returns the raw transaction of the given hash.
	GetRawTransaction(ctx context.Context, hash *chainhash.Hash) (*btcjson.TxRawResult, error)

	// GetBlockHash returns the block hash of the given height.
	GetBlockHash(ctx context.Context, height int64) (*chainhash.Hash, error)

	// GetBlock returns the block detail with the given hash.
	GetBlock(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error)

	// GetBlockVerbose returns the block detail with the given hash. It has more detailed info about txs than `GetBlock`
	GetBlockVerbose(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error)

	// GetTxOut returns details about an unspent transaction output. It will return nil result if the utxo has been
	// spent.
	GetTxOut(ctx context.Context, hash *chainhash.Hash, vout uint32) (*btcjson.GetTxOutResult, error)

	// GetNetworkInfo returns the network configuration of the node we connect to.
	GetNetworkInfo(ctx context.Context) (*btcjson.GetNetworkInfoResult, error)

	// GetMempoolEntry returns details on the active state of the TX memory pool.
	GetMempoolEntry(ctx context.Context, txid string) (*btcjson.GetMempoolEntryResult, error)
}

type Request struct {
	Version string            `json:"jsonrpc"`
	ID      uint64            `json:"id"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
}

// Response is the raw bytes of a JSON-RPC result, or the error if the response
// error object was non-null.
type Response struct {
	result []byte
	err    error
}

// rawResponse is a partially-unmarshaled JSON-RPC response.  For this
// to be valid (according to JSON-RPC 1.0 spec), ID may not be nil.
type rawResponse struct {
	Result json.RawMessage   `json:"result"`
	Error  *btcjson.RPCError `json:"error"`
}

func (r rawResponse) Response() Response {
	if r.Error != nil {
		return Response{err: r.Error}
	}
	return Response{r.Result, nil}
}

type rpcClient struct {
	network    *chaincfg.Params
	url        string
	user       string
	password   string
	httpClient *http.Client
}

func NewClient(network *chaincfg.Params, url, user, password string) Client {
	return &rpcClient{
		network:    network,
		url:        url,
		user:       user,
		password:   password,
		httpClient: new(http.Client),
	}
}

func (client *rpcClient) Net() *chaincfg.Params {
	return client.network
}

func (client *rpcClient) GetBlockchainInfo(ctx context.Context) (*btcjson.GetBlockChainInfoResult, error) {
	method := "getblockchaininfo"
	result, err := client.send(ctx, method, nil)
	if err != nil {
		return nil, err
	}
	var res btcjson.GetBlockChainInfoResult
	err = json.Unmarshal(result, &res)
	return &res, err
}

func (client *rpcClient) SubmitTx(ctx context.Context, tx *wire.MsgTx) error {
	method := "sendrawtransaction"
	buff := bytes.NewBuffer([]byte{})
	if err := tx.Serialize(buff); err != nil {
		return err
	}
	params, err := client.packParams(hex.EncodeToString(buff.Bytes()))
	if err != nil {
		return err
	}
	_, err = client.send(ctx, method, params)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "txn-mempool-conflict"):
			return ErrMempoolConflict
		case strings.Contains(err.Error(), "bad-txns-inputs-missingorspent"):
			return ErrTxInputsMissingOrSpent
		case strings.Contains(err.Error(), "Transaction already in block chain") || strings.Contains(err.Error(), "Transaction outputs already in utxo set"):
			return ErrAlreadyInUtxoSet
		}
	}
	return err
}

func (client *rpcClient) GetRawTransaction(ctx context.Context, hash *chainhash.Hash) (*btcjson.TxRawResult, error) {
	method := "getrawtransaction"
	params, err := client.packParams(hash.String(), true)
	if err != nil {
		return nil, err
	}
	result, err := client.send(ctx, method, params)
	if err != nil {
		if strings.Contains(err.Error(), "No such mempool or blockchain transaction") {
			return nil, ErrTxNotFound
		}
		return nil, err
	}
	var res btcjson.TxRawResult
	if err := json.Unmarshal(result, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (client *rpcClient) GetBlockHash(ctx context.Context, height int64) (*chainhash.Hash, error) {
	method := "getblockhash"
	params, err := client.packParams(height)
	if err != nil {
		return nil, err
	}
	result, err := client.send(ctx, method, params)
	if err != nil {
		return nil, err
	}
	var res chainhash.Hash
	err = json.Unmarshal(result, &res)
	return &res, err
}

func (client *rpcClient) GetBlock(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	method := "getblock"
	params, err := client.packParams(hash.String())
	if err != nil {
		return nil, err
	}
	result, err := client.send(ctx, method, params)
	if err != nil {
		return nil, err
	}
	var res btcjson.GetBlockVerboseResult
	err = json.Unmarshal(result, &res)
	return &res, err
}

func (client *rpcClient) GetBlockVerbose(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error) {
	method := "getblock"
	params, err := client.packParams(hash.String(), 2)
	if err != nil {
		return nil, err
	}
	result, err := client.send(ctx, method, params)
	if err != nil {
		return nil, err
	}
	var res btcjson.GetBlockVerboseTxResult
	err = json.Unmarshal(result, &res)
	return &res, err
}

func (client *rpcClient) GetTxOut(ctx context.Context, hash *chainhash.Hash, vout uint32) (*btcjson.GetTxOutResult, error) {
	method := "gettxout"
	params, err := client.packParams(hash.String(), vout)
	if err != nil {
		return nil, err
	}
	result, err := client.send(ctx, method, params)
	if err != nil {
		return nil, err
	}
	var res *btcjson.GetTxOutResult
	err = json.Unmarshal(result, &res)
	return res, err
}

func (client *rpcClient) GetNetworkInfo(ctx context.Context) (*btcjson.GetNetworkInfoResult, error) {
	method := "getnetworkinfo"
	result, err := client.send(ctx, method, nil)
	if err != nil {
		return nil, err
	}
	var res btcjson.GetNetworkInfoResult
	err = json.Unmarshal(result, &res)
	return &res, err
}

func (client *rpcClient) GetMempoolEntry(ctx context.Context, txid string) (*btcjson.GetMempoolEntryResult, error) {
	method := "getmempoolentry"
	params, err := client.packParams(txid)
	if err != nil {
		return nil, err
	}
	result, err := client.send(ctx, method, params)
	if err != nil {
		if strings.Contains(err.Error(), "Transaction not in mempool") {
			return nil, ErrTxNotInMempool
		}
		return nil, err
	}
	var res btcjson.GetMempoolEntryResult
	err = json.Unmarshal(result, &res)
	return &res, nil
}

func (client *rpcClient) send(ctx context.Context, method string, params []json.RawMessage) ([]byte, error) {
	// Construct the request
	jReq := Request{
		Version: "1.0",
		ID:      rand.Uint64(),
		Method:  method,
		Params:  params,
	}
	raw, err := json.Marshal(jReq)
	if err != nil {
		return nil, err
	}

	// Post the request
	var (
		lastErr      error
		backoff      time.Duration
		httpResponse *http.Response
	)

	for i := 0; i < 10; i++ {
		bodyReader := bytes.NewReader(raw)
		httpReq, err := http.NewRequest("POST", client.url, bodyReader)
		if err != nil {
			return nil, err
		}
		httpReq.Close = true
		httpReq.Header.Set("Content-Type", "application/json")
		if client.user != "" && client.password != "" {
			httpReq.SetBasicAuth(client.user, client.password)
		}

		httpResponse, err = client.httpClient.Do(httpReq)

		// Quit the retry loop on success or if we can't retry anymore.
		if err == nil || i == 10 {
			break
		}

		// Save the last error for the case where we backoff further,
		// retry and get an invalid response but no error. If this
		// happens the saved last error will be used to enrich the error
		// message that we pass back to the caller.
		lastErr = err

		// Backoff sleep otherwise.
		backoff = 500 * time.Millisecond * time.Duration(i+1)
		if backoff > 5*time.Second {
			backoff = 5 * time.Second
		}

		select {
		case <-ctx.Done():
			return nil, err
		case <-time.After(backoff):
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}

	// We still want to return an error if for any reason the response
	// remains empty.
	if httpResponse == nil {
		return nil, fmt.Errorf("invalid http POST response (nil), "+
			"method: %s, id: %d, last error=%v",
			jReq.Method, jReq.ID, lastErr)
	}

	// Read the raw bytes and close the response.
	respBytes, err := io.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error reading json reply: %v", err)
	}

	// Try to unmarshal the response as a regular JSON-RPC response.
	var resp rawResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

func (client *rpcClient) packParams(params ...interface{}) ([]json.RawMessage, error) {
	rawParams := make([]json.RawMessage, len(params))
	for i, param := range params {
		var err error
		rawParams[i], err = json.Marshal(param)
		if err != nil {
			return nil, err
		}
	}
	return rawParams, nil
}
