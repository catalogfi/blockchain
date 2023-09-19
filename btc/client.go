package btc

import (
	"errors"
	"fmt"
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

// Client to interact with the bitcoin network. We'll follow the standard bitcoind rpc interface if that meet our
// requirements.
type Client interface {

	// Net returns the network params.
	Net() *chaincfg.Params

	// LatestBlock returns the height and hash of the latest block.
	LatestBlock() (int64, string, error)

	// SubmitTx to the bitcoin network
	SubmitTx(tx *wire.MsgTx) error

	// GetRawTransaction returns the raw transaction of the given hash.
	GetRawTransaction(txhash []byte) (*btcjson.TxRawResult, error)

	// GetBlockByHeight returns the block detail of the given height.
	GetBlockByHeight(height int64) (*btcjson.GetBlockVerboseResult, error)

	// GetBlockByHash returns the block detail with the given hash.
	GetBlockByHash(hash string) (*btcjson.GetBlockVerboseResult, error)

	// GetTxOut returns details about an unspent transaction output. It will return nil result if the utxo has been
	// spent.
	GetTxOut(hash string, vout uint32) (*btcjson.GetTxOutResult, error)
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

// LatestBlock returns the height and hash of the longest blockchain.
func (client *client) LatestBlock() (int64, string, error) {
	res, err := client.rpcClient.GetBlockChainInfo()
	if err != nil {
		return 0, "", err
	}
	return int64(res.Blocks), res.BestBlockHash, nil
}

// SubmitTx to the Bitcoin network.
func (client *client) SubmitTx(tx *wire.MsgTx) error {
	_, err := client.rpcClient.SendRawTransaction(tx, false)
	if err != nil {
		var rpcErr *btcjson.RPCError
		if errors.As(err, &rpcErr) {
			switch rpcErr.Code {
			case btcjson.ErrRPCVerifyAlreadyInChain:
				return fmt.Errorf(`bad "sendrawtransaction": %w`, ErrAlreadyInChain)
			case btcjson.ErrRPCTxRejected:
				if strings.Contains(err.Error(), "txn-mempool-conflict") {
					return fmt.Errorf(`bad "sendrawtransaction": %w`, ErrMempoolConflict)
				}
			case btcjson.ErrRPCTxError:
				if strings.Contains(err.Error(), "bad-txns-inputs-missingorspent") {
					return fmt.Errorf(`bad "sendrawtransaction": %w`, ErrTxInputsMissingOrSpent)
				}
			}
		}
		return err
	}
	return nil
}

func (client *client) GetRawTransaction(txhash []byte) (*btcjson.TxRawResult, error) {
	hash, err := chainhash.NewHash(txhash)
	if err != nil {
		return nil, err
	}
	res, err := client.rpcClient.GetRawTransactionVerbose(hash)
	if err != nil {
		if strings.Contains(err.Error(), "No such mempool or blockchain transaction") {
			return nil, fmt.Errorf(`bad "getrawtransaction": %w`, ErrTxNotFound)
		}
		return nil, err
	}
	return res, nil
}

func (client *client) GetBlockByHeight(height int64) (*btcjson.GetBlockVerboseResult, error) {
	hash, err := client.rpcClient.GetBlockHash(height)
	if err != nil {
		return nil, fmt.Errorf("bad \"getblockhash\": %w", err)
	}
	return client.rpcClient.GetBlockVerbose(hash)
}

func (client *client) GetBlockByHash(hash string) (*btcjson.GetBlockVerboseResult, error) {
	blockHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return nil, err
	}
	return client.rpcClient.GetBlockVerbose(blockHash)
}

func (client *client) GetTxOut(hash string, vout uint32) (*btcjson.GetTxOutResult, error) {
	txhash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return nil, err
	}
	return client.rpcClient.GetTxOut(txhash, vout, true)
}

// MockClient for testing purpose
type MockClient struct {
	FuncNet               func() *chaincfg.Params
	FuncGetUTXOs          func(btcutil.Address) (UTXOs, error)
	FuncLatestBlock       func() (int64, string, error)
	FuncSubmitTx          func(*wire.MsgTx) error
	FuncGetRawTransaction func([]byte) (*btcjson.TxRawResult, error)
	FuncGetBlockByHeight  func(int64) (*btcjson.GetBlockVerboseResult, error)
	FuncGetBlockByHash    func(string) (*btcjson.GetBlockVerboseResult, error)
	FuncGetTxOut          func(string, uint32) (*btcjson.GetTxOutResult, error)
}

// Make sure the MockClient implements the Client interface
var _ Client = NewMockClient()

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) Net() *chaincfg.Params {
	return m.FuncNet()
}

func (m *MockClient) GetUTXOs(address btcutil.Address) (UTXOs, error) {
	return m.FuncGetUTXOs(address)
}

func (m *MockClient) LatestBlock() (int64, string, error) {
	return m.FuncLatestBlock()
}

func (m *MockClient) SubmitTx(tx *wire.MsgTx) error {
	return m.FuncSubmitTx(tx)
}

func (m *MockClient) GetRawTransaction(txhash []byte) (*btcjson.TxRawResult, error) {
	return m.FuncGetRawTransaction(txhash)
}

func (m *MockClient) GetBlockByHeight(height int64) (*btcjson.GetBlockVerboseResult, error) {
	return m.FuncGetBlockByHeight(height)
}

func (m *MockClient) GetBlockByHash(hash string) (*btcjson.GetBlockVerboseResult, error) {
	return m.FuncGetBlockByHash(hash)
}

func (m *MockClient) GetTxOut(hash string, vout uint32) (*btcjson.GetTxOutResult, error) {
	return m.FuncGetTxOut(hash, vout)
}
