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
	FuncNet               func() *chaincfg.Params
	FuncGetUTXOs          func(context.Context, btcutil.Address) (btc.UTXOs, error)
	FuncLatestBlock       func(context.Context) (int64, string, error)
	FuncSubmitTx          func(context.Context, *wire.MsgTx) error
	FuncGetRawTransaction func(context.Context, *chainhash.Hash) (*btcjson.TxRawResult, error)
	FuncGetBlockByHeight  func(context.Context, int64) (*btcjson.GetBlockVerboseResult, error)
	FuncGetBlockByHash    func(context.Context, *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error)
	FuncGetTxOut          func(context.Context, *chainhash.Hash, uint32) (*btcjson.GetTxOutResult, error)
}

// Make sure the MockClient implements the Client interface in case
var _ btc.Client = NewMockClient()

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) Net() *chaincfg.Params {
	return m.FuncNet()
}

func (m *MockClient) GetUTXOs(ctx context.Context, address btcutil.Address) (btc.UTXOs, error) {
	return m.FuncGetUTXOs(ctx, address)
}

func (m *MockClient) LatestBlock(ctx context.Context) (int64, string, error) {
	return m.FuncLatestBlock(ctx)
}

func (m *MockClient) SubmitTx(ctx context.Context, tx *wire.MsgTx) error {
	return m.FuncSubmitTx(ctx, tx)
}

func (m *MockClient) GetRawTransaction(ctx context.Context, txhash *chainhash.Hash) (*btcjson.TxRawResult, error) {
	return m.FuncGetRawTransaction(ctx, txhash)
}

func (m *MockClient) GetBlockByHeight(ctx context.Context, height int64) (*btcjson.GetBlockVerboseResult, error) {
	return m.FuncGetBlockByHeight(ctx, height)
}

func (m *MockClient) GetBlockByHash(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	return m.FuncGetBlockByHash(ctx, hash)
}

func (m *MockClient) GetTxOut(ctx context.Context, hash *chainhash.Hash, vout uint32) (*btcjson.GetTxOutResult, error) {
	return m.FuncGetTxOut(ctx, hash, vout)
}

// MockIndexerClient for testing purpose
type MockIndexerClient struct {
	FuncGetAddressTxs     func(context.Context, btcutil.Address) ([]btc.Transaction, error)
	FuncGetUTXOs          func(context.Context, btcutil.Address) (btc.UTXOs, error)
	FuncGetTipBlockHeight func(ctx context.Context) (uint64, error)
	FuncGetTx             func(context.Context, string) (btc.Transaction, error)
	FuncSubmitTx          func(ctx context.Context, tx *wire.MsgTx) error
	FuncFeeEstimate       func(ctx context.Context) (btc.FeeSuggestion, error)
}

func (client MockIndexerClient) GetAddressTxs(ctx context.Context, address btcutil.Address) ([]btc.Transaction, error) {
	return client.FuncGetAddressTxs(ctx, address)
}

func (client MockIndexerClient) GetUTXOs(ctx context.Context, address btcutil.Address) (btc.UTXOs, error) {
	return client.FuncGetUTXOs(ctx, address)
}

func (client MockIndexerClient) GetTipBlockHeight(ctx context.Context) (uint64, error) {
	return client.FuncGetTipBlockHeight(ctx)
}

func (client MockIndexerClient) GetTx(ctx context.Context, txid string) (btc.Transaction, error) {
	return client.FuncGetTx(ctx, txid)
}

func (client MockIndexerClient) SubmitTx(ctx context.Context, tx *wire.MsgTx) error {
	return client.FuncSubmitTx(ctx, tx)
}

func (client MockIndexerClient) FeeEstimate(ctx context.Context) (btc.FeeSuggestion, error) {
	return client.FeeEstimate(ctx)
}

// Make sure the MockIndexerClient implements the Client interface
var _ *MockIndexerClient = NewMockIndexerClient()

func NewMockIndexerClient() *MockIndexerClient {
	return &MockIndexerClient{}
}
