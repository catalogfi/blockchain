package btctest

import (
	"context"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/catalogfi/blockchain/btc"
)

// MockClient is a mock implementation of the `btc.Client` for testing purposes. All the function can be customised to
// test different scenarios.
type MockClient struct {
	Client                btc.Client
	FuncNet               func() *chaincfg.Params
	FuncGetUTXOs          func(context.Context, btcutil.Address) (btc.UTXOs, error)
	FuncLatestBlock       func(context.Context) (int64, string, error)
	FuncSubmitTx          func(context.Context, *wire.MsgTx) error
	FuncGetRawTransaction func(context.Context, *chainhash.Hash) (*btcjson.TxRawResult, error)
	FuncGetBlockByHeight  func(context.Context, int64) (*btcjson.GetBlockVerboseResult, error)
	FuncGetBlockByHash    func(context.Context, *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error)
	FuncGetTxOut          func(context.Context, *chainhash.Hash, uint32) (*btcjson.GetTxOutResult, error)
}

// Make sure the MockClient implements the Client interface
var _ btc.Client = NewMockClient()

func NewMockClient() *MockClient {
	return &MockClient{}
}

func NewMockClientWrapper(client btc.Client) *MockClient {
	return &MockClient{
		Client: client,
	}
}

func (m *MockClient) Net() *chaincfg.Params {
	if m.FuncNet != nil {
		return m.FuncNet()
	}
	return m.Client.Net()
}

func (m *MockClient) LatestBlock(ctx context.Context) (int64, string, error) {
	if m.FuncLatestBlock != nil {
		return m.FuncLatestBlock(ctx)
	}
	return m.Client.LatestBlock(ctx)
}

func (m *MockClient) SubmitTx(ctx context.Context, tx *wire.MsgTx) error {
	if m.FuncSubmitTx != nil {
		return m.FuncSubmitTx(ctx, tx)
	}
	return m.SubmitTx(ctx, tx)
}

func (m *MockClient) GetRawTransaction(ctx context.Context, txhash *chainhash.Hash) (*btcjson.TxRawResult, error) {
	if m.FuncGetRawTransaction != nil {
		return m.FuncGetRawTransaction(ctx, txhash)
	}
	return m.Client.GetRawTransaction(ctx, txhash)
}

func (m *MockClient) GetBlockByHeight(ctx context.Context, height int64) (*btcjson.GetBlockVerboseResult, error) {
	if m.FuncGetBlockByHeight != nil {
		return m.FuncGetBlockByHeight(ctx, height)
	}
	return m.Client.GetBlockByHeight(ctx, height)
}

func (m *MockClient) GetBlockByHash(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	if m.FuncGetBlockByHash != nil {
		return m.FuncGetBlockByHash(ctx, hash)
	}
	return m.Client.GetBlockByHash(ctx, hash)
}

func (m *MockClient) GetTxOut(ctx context.Context, hash *chainhash.Hash, vout uint32) (*btcjson.GetTxOutResult, error) {
	if m.FuncGetTxOut != nil {
		return m.FuncGetTxOut(ctx, hash, vout)
	}
	return m.Client.GetTxOut(ctx, hash, vout)
}

// MockIndexerClient for testing purpose
type MockIndexerClient struct {
	Indexer               btc.IndexerClient
	FuncGetAddressTxs     func(context.Context, btcutil.Address) ([]btc.Transaction, error)
	FuncGetUTXOs          func(context.Context, btcutil.Address) (btc.UTXOs, error)
	FuncGetTipBlockHeight func(ctx context.Context) (uint64, error)
	FuncGetTx             func(context.Context, string) (btc.Transaction, error)
	FuncSubmitTx          func(ctx context.Context, tx *wire.MsgTx) error
	FuncFeeEstimate       func(ctx context.Context) (btc.FeeSuggestion, error)
}

// Make sure the MockIndexerClient implements the Client interface
var _ btc.IndexerClient = NewMockIndexerClient()

func NewMockIndexerClient() *MockIndexerClient {
	return &MockIndexerClient{}
}

func (client MockIndexerClient) GetAddressTxs(ctx context.Context, address btcutil.Address) ([]btc.Transaction, error) {
	if client.FuncGetAddressTxs != nil {
		return client.FuncGetAddressTxs(ctx, address)
	}
	return client.Indexer.GetAddressTxs(ctx, address)
}

func (client MockIndexerClient) GetUTXOs(ctx context.Context, address btcutil.Address) (btc.UTXOs, error) {
	if client.FuncGetUTXOs != nil {
		return client.FuncGetUTXOs(ctx, address)
	}
	return client.Indexer.GetUTXOs(ctx, address)
}

func (client MockIndexerClient) GetTipBlockHeight(ctx context.Context) (uint64, error) {
	if client.FuncGetTipBlockHeight != nil {
		return client.Indexer.GetTipBlockHeight(ctx)
	}
	return client.Indexer.GetTipBlockHeight(ctx)
}

func (client MockIndexerClient) GetTx(ctx context.Context, txid string) (btc.Transaction, error) {
	if client.FuncGetTx != nil {
		return client.FuncGetTx(ctx, txid)
	}
	return client.Indexer.GetTx(ctx, txid)
}

func (client MockIndexerClient) SubmitTx(ctx context.Context, tx *wire.MsgTx) error {
	if client.FuncSubmitTx != nil {
		return client.FuncSubmitTx(ctx, tx)
	}
	return client.Indexer.SubmitTx(ctx, tx)
}

func (client MockIndexerClient) FeeEstimate(ctx context.Context) (btc.FeeSuggestion, error) {
	if client.FuncFeeEstimate != nil {
		return client.FuncFeeEstimate(ctx)
	}
	return client.Indexer.FeeEstimate(ctx)
}
