package btc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

var (
	ErrTxNotFound = errors.New("no such mempool or blockchain transaction")

	ErrAlreadyInChain = errors.New("transaction already in block chain")

	ErrTxInputsMissingOrSpent = errors.New("bad-txns-inputs-missingorspent")

	ErrMempoolConflict = errors.New("txn-mempool-conflict")
)

// Client to interact with the Bitcoin network. It's implementation uses standard bitcoind JSON-RPC behind the scene.
type Client interface {

	// Net returns the network params.
	Net() *chaincfg.Params

	// LatestBlock returns the height and hash of the latest block.
	LatestBlock(ctx context.Context) (int64, string, error)

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
}

type client struct {
	config    *rpcclient.ConnConfig
	params    *chaincfg.Params
	rpcClient *rpcclient.Client
}

func NewClient(config *rpcclient.ConnConfig) (Client, error) {
	c, err := rpcclient.New(config, nil)
	if err != nil {
		return nil, err
	}

	var param *chaincfg.Params
	switch config.Params {
	case chaincfg.MainNetParams.Name:
		param = &chaincfg.MainNetParams
	case chaincfg.TestNet3Params.Name:
		param = &chaincfg.TestNet3Params
	case chaincfg.RegressionNetParams.Name:
		param = &chaincfg.RegressionNetParams
	default:
		return nil, fmt.Errorf("rpcclient.New: Unknown chain %s", config.Params)
	}

	return &client{
		config:    config,
		params:    param,
		rpcClient: c,
	}, nil
}

func (client *client) Net() *chaincfg.Params {
	return client.params
}

func (client *client) LatestBlock(ctx context.Context) (int64, string, error) {
	future := client.rpcClient.GetBlockChainInfoAsync()
	results := make(chan *btcjson.GetBlockChainInfoResult, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return 0, "", fmt.Errorf("LatestBlock : %w", ctx.Err())
	case err := <-errs:
		return 0, "", err
	case result := <-results:
		return int64(result.Blocks), result.BestBlockHash, nil
	}
}

func (client *client) SubmitTx(ctx context.Context, tx *wire.MsgTx) error {
	// The SendRawTransactionAsync is not technically asynchronous,
	// we need to have an extra channel for it.
	futureChan := make(chan rpcclient.FutureSendRawTransactionResult, 1)
	go func() {
		defer close(futureChan)
		future := client.rpcClient.SendRawTransactionAsync(tx, false)
		futureChan <- future
	}()

	var future rpcclient.FutureSendRawTransactionResult
	select {
	case <-ctx.Done():
		return fmt.Errorf("SubmitTx : %w", ctx.Err())
	case future = <-futureChan:
	}

	results := make(chan *chainhash.Hash, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("SubmitTx : %w", ctx.Err())
	case err := <-errs:
		// Parse the error based on the error code and message
		var rpcErr *btcjson.RPCError
		if errors.As(err, &rpcErr) {
			switch rpcErr.Code {
			case btcjson.ErrRPCVerifyAlreadyInChain:
				return ErrAlreadyInChain
			case btcjson.ErrRPCTxRejected:
				if strings.Contains(err.Error(), "txn-mempool-conflict") {
					return ErrMempoolConflict
				}
			case btcjson.ErrRPCTxError:
				if strings.Contains(err.Error(), "bad-txns-inputs-missingorspent") {
					return ErrTxInputsMissingOrSpent
				}
			}
		}
		return err
	case <-results:
		return nil
	}
}

func (client *client) GetRawTransaction(ctx context.Context, txhash *chainhash.Hash) (*btcjson.TxRawResult, error) {
	future := client.rpcClient.GetRawTransactionVerboseAsync(txhash)
	results := make(chan *btcjson.TxRawResult, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("GetRawTransaction : %w", ctx.Err())
	case err := <-errs:
		// Parse the error based on the error code and message
		var rpcErr *btcjson.RPCError
		if errors.As(err, &rpcErr) {
			switch rpcErr.Code {
			case btcjson.ErrRPCInvalidAddressOrKey:
				if strings.Contains(err.Error(), "No such mempool or blockchain transaction") {
					return nil, ErrTxNotFound
				}
			}
		}

		return nil, err
	case result := <-results:
		return result, nil
	}
}

func (client *client) GetBlockHash(ctx context.Context, height int64) (*chainhash.Hash, error) {
	future := client.rpcClient.GetBlockHashAsync(height)
	results := make(chan *chainhash.Hash, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("GetBlockByHeight : %w", ctx.Err())
	case err := <-errs:
		return nil, err
	case hash := <-results:
		return hash, nil
	}
}

func (client *client) GetBlock(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	future := client.rpcClient.GetBlockVerboseAsync(hash)
	results := make(chan *btcjson.GetBlockVerboseResult, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("GetBlockByHash : %w", ctx.Err())
	case err := <-errs:
		return nil, err
	case result := <-results:
		return result, nil
	}
}

func (client *client) GetBlockVerbose(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error) {
	future := client.rpcClient.GetBlockVerboseTxAsync(hash)
	results := make(chan *btcjson.GetBlockVerboseTxResult, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("GetBlockByHash : %w", ctx.Err())
	case err := <-errs:
		return nil, err
	case result := <-results:
		return result, nil
	}
}

func (client *client) GetTxOut(ctx context.Context, hash *chainhash.Hash, vout uint32) (*btcjson.GetTxOutResult, error) {
	future := client.rpcClient.GetTxOutAsync(hash, vout, true)
	results := make(chan *btcjson.GetTxOutResult, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("GetTxOut : %w", ctx.Err())
	case err := <-errs:
		return nil, err
	case result := <-results:
		return result, nil
	}
}

func (client *client) GetNetworkInfo(ctx context.Context) (*btcjson.GetNetworkInfoResult, error) {
	future := client.rpcClient.GetNetworkInfoAsync()
	results := make(chan *btcjson.GetNetworkInfoResult, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("GetNetworkInfo : %w", ctx.Err())
	case err := <-errs:
		return nil, err
	case result := <-results:
		return result, nil
	}
}

type Fees struct {
	Base       float64 `json:"base"`
	Modified   float64 `json:"modified"`
	Ancestor   float64 `json:"ancestor"`
	Descendant float64 `json:"descendant"`
}

// Transaction represents a single transaction entry
type DescendantTransaction struct {
	Vsize             int      `json:"vsize"`
	Weight            int      `json:"weight"`
	Time              int64    `json:"time"`
	Height            int      `json:"height"`
	DescendantCount   int      `json:"descendantcount"`
	DescendantSize    int      `json:"descendantsize"`
	AncestorCount     int      `json:"ancestorcount"`
	AncestorSize      int      `json:"ancestorsize"`
	WtxID             string   `json:"wtxid"`
	Fees              Fees     `json:"fees"`
	Depends           []string `json:"depends"`
	SpentBy           []string `json:"spentby"`
	BIP125Replaceable bool     `json:"bip125-replaceable"`
	Unbroadcast       bool     `json:"unbroadcast"`
}

type Request struct {
	Version string            `json:"jsonrpc"`
	ID      uint32            `json:"id"`
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

type BitcoinClient struct {
	RpcUser    string
	RpcPass    string
	RpcURL     string
	httpClient *http.Client
}

func NewBitcoinClient(rpcUser string, rpcPass string, rpcURL string) BitcoinClient {
	return BitcoinClient{
		RpcUser:    rpcUser,
		RpcPass:    rpcPass,
		RpcURL:     rpcURL,
		httpClient: new(http.Client),
	}
}

// GetMempoolDescendants returns all descendants (both direct and indirect) of a transaction in the mempool.
func (client *BitcoinClient) GetMempoolDescendants(ctx context.Context, txid string) (map[string]DescendantTransaction, error) {
	method := "getmempooldescendants"
	params, err := client.packParams(txid, true)
	if err != nil {
		return nil, err
	}
	result, err := client.send(ctx, method, params)
	if err != nil {
		if strings.Contains(err.Error(), "Transaction not in mempool") {
			return nil, ErrTxNotFound
		}
		return nil, err
	}

	var apiResponse map[string]DescendantTransaction
	if err := json.Unmarshal(result, &apiResponse); err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %w", err)
	}
	return apiResponse, nil
}

// GetDescendantsFee calculates the total fee of all descendants of a transaction in the mempool.
// it fetches the descendants using `getmempooldescendants` and sums their fees.
// If the transaction is not found in the mempool, it returns ErrTxNotFound.
func (client *BitcoinClient) GetDescendantsFee(ctx context.Context, txid string) (int64, error) {
	apiResponse, err := client.GetMempoolDescendants(ctx, txid)
	if err != nil {
		return 0, fmt.Errorf("Error getting mempool descendants: %w", err)
	}

	if len(apiResponse) == 0 {
		return 0, nil
	}

	var sum float64
	for _, tx := range apiResponse {
		sum += tx.Fees.Base
	}

	amount, err := btcutil.NewAmount(sum)
	if err != nil {
		return 0, fmt.Errorf("Error converting sum to btcutil.Amount: %w", err)
	}
	return int64(amount), nil
}

// GetMempoolEntryResult models the data returned from the getmempoolentry's
// fee field
type MempoolFees struct {
	// Base is the fee of the transaction itself, in BTC.
	Base     float64 `json:"base"`
	Modified float64 `json:"modified"`
	Ancestor float64 `json:"ancestor"`
	// Descendant is the total fee of the transaction and all its descendants (both direct and indirect), in BTC.
	Descendant float64 `json:"descendant"`
}

// GetMempoolEntryResult models the data returned from the getmempoolentry
type GetMempoolEntryResult struct {
	// Vsize of the actual transaction
	VSize           int32   `json:"vsize"`
	Size            int32   `json:"size"`
	Weight          int64   `json:"weight"`
	Fee             float64 `json:"fee"`
	ModifiedFee     float64 `json:"modifiedfee"`
	Time            int64   `json:"time"`
	Height          int64   `json:"height"`
	DescendantCount int64   `json:"descendantcount"`
	// DescendantSize is the total Vsize of the transaction and all its descendants
	DescendantSize int64 `json:"descendantsize"`
	// DescendantFees is the total fee of the transaction and all its descendants(both direct and indirect)
	DescendantFees float64     `json:"descendantfees"`
	AncestorCount  int64       `json:"ancestorcount"`
	AncestorSize   int64       `json:"ancestorsize"`
	AncestorFees   float64     `json:"ancestorfees"`
	WTxId          string      `json:"wtxid"`
	Fees           MempoolFees `json:"fees"`
	Depends        []string    `json:"depends"`
	// SpentBy is a list of transaction IDs that spend outputs from this transaction.
	// This is useful for tracking direct descendants of a transaction in the mempool.
	SpendBy []string `json:"spentby"`
}

// GetMempoolEntry returns mempool data for given transaction
func (client *BitcoinClient) GetMempoolEntry(ctx context.Context, txid string) (*GetMempoolEntryResult, error) {
	method := "getmempoolentry"
	params, err := client.packParams(txid)
	if err != nil {
		return nil, err
	}
	result, err := client.send(ctx, method, params)
	if err != nil {
		if strings.Contains(err.Error(), "Transaction not in mempool") {
			return nil, ErrTxNotFound
		}
		return nil, err
	}
	var res GetMempoolEntryResult
	err = json.Unmarshal(result, &res)
	return &res, err
}

// GetTxOut returns details about an unspent transaction output. It returns nil
// when the outpoint is already spent.
func (client *BitcoinClient) GetTxOut(ctx context.Context, hash *chainhash.Hash, vout uint32) (*btcjson.GetTxOutResult, error) {
	method := "gettxout"
	params, err := client.packParams(hash.String(), vout, true)
	if err != nil {
		return nil, err
	}
	result, err := client.send(ctx, method, params)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(bytes.TrimSpace(result), []byte("null")) {
		return nil, nil
	}

	var txOut btcjson.GetTxOutResult
	if err := json.Unmarshal(result, &txOut); err != nil {
		return nil, fmt.Errorf("gettxout: failed to unmarshal response: %w", err)
	}
	return &txOut, nil
}

type RBFTxFeeInfo struct {
	// Total fees of in-mempool descendants (including this transaction) in sats
	TotalFee float64 `json:"total_fee"`
	// Total fee of the descendants in sats
	DescendantFee float64 `json:"descendant_fee"`
	// current fee of the tx to be replaced ,
	// the fee rate is calculated based on the actual tx and its direct descendants
	// source : https://github.dev/bitcoin/bitcoin/blob/ed060e01e756fc8d4c1c590e79f9555ae86c433f/src/policy/rbf.cpp#L145
	TxFeeRate float64 `json:"tx_fee_rate"`
}

// GetRBFTxFeeInfo returns FeeInfo required to RBF (Replace-By-Fee) a transaction in the mempool
func (client *BitcoinClient) GetRBFTxFeeInfo(ctx context.Context, txid string) (*RBFTxFeeInfo, error) {
	// Get the descendants and the entry for the transaction
	descendants, err := client.GetMempoolDescendants(ctx, txid)
	if err != nil {
		return nil, fmt.Errorf("error getting mempool descendants: %w", err)
	}
	entry, err := client.GetMempoolEntry(ctx, txid)
	if err != nil {
		return nil, fmt.Errorf("error getting mempool entry: %w", err)
	}
	if entry == nil {
		return nil, ErrTxNotFound
	}

	// Fees are in BTC, convert to sats
	totalFee := math.Ceil(entry.Fees.Descendant * 1e8)
	baseFee := math.Ceil(entry.Fees.Base * 1e8)
	descendantFee := math.Ceil((entry.Fees.Descendant - entry.Fees.Base) * 1e8)

	feeInfo := &RBFTxFeeInfo{
		TotalFee:      totalFee,
		DescendantFee: descendantFee,
		TxFeeRate:     float64(totalFee) / float64(entry.DescendantSize),
	}
	// If there are no descendants, fee rate is is the base fee divided by the vsize of the transaction
	if len(descendants) == 0 {
		return feeInfo, nil
	}

	// Calculate the fee rate based on the direct descendants
	// The fee rate is calculated as the total fee of the transaction and its direct descendants divided
	// by the total vsize of the transaction and its direct descendants.
	var directDescendantSize, directDescendantFee float64
	for _, childTxid := range entry.SpendBy {
		if desc, ok := descendants[childTxid]; ok {
			directDescendantSize += float64(desc.Vsize)
			directDescendantFee += float64(desc.Fees.Base * 1e8)
		}
	}
	if directDescendantSize+float64(entry.VSize) > 0 {
		feeInfo.TxFeeRate = float64(directDescendantFee+baseFee) / float64(directDescendantSize+float64(entry.VSize))
	}

	return feeInfo, nil
}

// TestMempoolAcceptResult models a single entry of the testmempoolaccept response.
type TestMempoolAcceptResult struct {
	Txid         string `json:"txid"`
	Wtxid        string `json:"wtxid"`
	Allowed      bool   `json:"allowed"`
	VSize        int    `json:"vsize,omitempty"`
	RejectReason string `json:"reject-reason,omitempty"`
}

// TestMempoolAccept performs a dry-run acceptance check of a raw transaction
// against the local mempool without broadcasting it.
func (client *BitcoinClient) TestMempoolAccept(ctx context.Context, rawTxHex string) (*TestMempoolAcceptResult, error) {
	method := "testmempoolaccept"
	params, err := client.packParams([]string{rawTxHex})
	if err != nil {
		return nil, err
	}
	result, err := client.send(ctx, method, params)
	if err != nil {
		return nil, err
	}
	var res []TestMempoolAcceptResult
	if err := json.Unmarshal(result, &res); err != nil {
		return nil, fmt.Errorf("testmempoolaccept: failed to unmarshal response: %w", err)
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("testmempoolaccept: empty response")
	}
	return &res[0], nil
}

// packParams converts the provided parameters into a slice of json.RawMessage.
func (client *BitcoinClient) packParams(params ...interface{}) ([]json.RawMessage, error) {
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

// send sends a JSON-RPC request to the Bitcoin node and returns the raw response.
func (client *BitcoinClient) send(ctx context.Context, method string, params []json.RawMessage) ([]byte, error) {
	// Construct the request
	jReq := Request{
		Version: "1.0",
		ID:      rand.Uint32(),
		Method:  method,
		Params:  params,
	}
	raw, err := json.Marshal(jReq)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(raw)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", client.RpcURL, bodyReader)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if client.RpcUser != "" && client.RpcPass != "" {
		httpReq.SetBasicAuth(client.RpcUser, client.RpcPass)
	}

	httpResponse, err := client.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	respBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading json reply: %v", err)
	}

	var resp rawResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}
