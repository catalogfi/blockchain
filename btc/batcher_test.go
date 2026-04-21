package btc_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"go.uber.org/zap"
)

type mockCache struct {
	batches     map[string]btc.Batch
	batchList   []string
	requests    map[string]btc.BatcherRequest
	requestList []string
	mode        btc.Strategy
}

func NewTestCache(mode btc.Strategy) btc.Cache {
	return &mockCache{
		batches:  make(map[string]btc.Batch),
		requests: make(map[string]btc.BatcherRequest),
		mode:     mode,
	}
}

func (m *mockCache) ReadBatchByReqID(ctx context.Context, id string) (btc.Batch, error) {
	for _, batchId := range m.batchList {
		batch, ok := m.batches[batchId]
		if !ok {
			return btc.Batch{}, fmt.Errorf("ReadBatchByReqId, batch not recorded")
		}
		if _, ok := batch.RequestIds[id]; ok {
			return batch, nil
		}
	}
	return btc.Batch{}, fmt.Errorf("ReadBatchByReqId, batch not found")
}

func (m *mockCache) ReadBatch(ctx context.Context, txId string) (btc.Batch, error) {
	batch, ok := m.batches[txId]
	if !ok {
		return btc.Batch{}, fmt.Errorf("bReadBatch, batch not found")
	}
	return batch, nil
}

func (m *mockCache) ReadPendingBatches(ctx context.Context) ([]btc.Batch, error) {
	batches := []btc.Batch{}
	for _, batch := range m.batches {
		if batch.Tx.Status.Confirmed == false {
			batches = append(batches, batch)
		}
	}
	return batches, nil
}
func (m *mockCache) SaveBatch(ctx context.Context, batch btc.Batch) error {
	if _, ok := m.batches[batch.Tx.TxID]; ok {
		return fmt.Errorf("batch already exists")
	}
	m.batches[batch.Tx.TxID] = batch
	for id := range batch.RequestIds {
		request := m.requests[id]
		request.Status = true
		m.requests[id] = request
	}
	m.batchList = append(m.batchList, batch.Tx.TxID)
	return nil
}

func (m *mockCache) ReadRequest(ctx context.Context, id string) (btc.BatcherRequest, error) {
	request, ok := m.requests[id]
	if !ok {
		return btc.BatcherRequest{}, fmt.Errorf("request not found")
	}
	return request, nil
}
func (m *mockCache) ReadPendingRequests(ctx context.Context) ([]btc.BatcherRequest, error) {
	requests := []btc.BatcherRequest{}
	for _, request := range m.requests {
		if request.Status == false {
			requests = append(requests, request)
		}
	}
	return requests, nil
}

func (m *mockCache) DeletePendingRequest(ctx context.Context, id string) error {
	delete(m.requests, id)
	for i, rid := range m.requestList {
		if rid == id {
			m.requestList = append(m.requestList[:i], m.requestList[i+1:]...)
			break
		}
	}
	return nil
}

func (m *mockCache) SaveRequest(ctx context.Context, req btc.BatcherRequest) error {
	if _, ok := m.requests[req.ID]; ok {
		return fmt.Errorf("request already exists")
	}
	m.requests[req.ID] = req
	m.requestList = append(m.requestList, req.ID)
	return nil
}

func (m *mockCache) UpdateBatchFees(ctx context.Context, txId []string, feeRate int64) error {
	for _, id := range txId {
		batch, ok := m.batches[id]
		if !ok {
			return fmt.Errorf("UpdateBatchFees, batch not found")
		}

		batch.Tx.Fee = int64(batch.Tx.Weight) * feeRate / 4
		m.batches[id] = batch
	}
	return nil
}

func (m *mockCache) UpdateAndDeletePendingBatches(ctx context.Context, updatedBatches ...btc.Batch) error {
	for _, batch := range updatedBatches {
		if _, ok := m.batches[batch.Tx.TxID]; !ok {
			return fmt.Errorf("UpdateAndDeleteBatches, batch not found")
		}
		m.batches[batch.Tx.TxID] = batch
	}

	// delete pending batches
	for _, id := range m.batchList {
		if m.batches[id].Tx.Status.Confirmed == false {
			delete(m.batches, id)
		}
	}
	return nil
}

func (m *mockCache) UpdateBatches(ctx context.Context, updatedBatches ...btc.Batch) error {
	for _, batch := range updatedBatches {
		if _, ok := m.batches[batch.Tx.TxID]; !ok {
			return fmt.Errorf("UpdateBatches, batch not found")
		}
		m.batches[batch.Tx.TxID] = batch
	}
	return nil
}

func (m *mockCache) ReadLatestBatch(ctx context.Context) (btc.Batch, error) {
	if len(m.batchList) == 0 {
		return btc.Batch{}, btc.ErrStoreNotFound
	}
	nbatches := len(m.batchList) - 1
	for nbatches >= 0 {
		batch, ok := m.batches[m.batchList[nbatches]]
		if ok && batch.Strategy == m.mode {
			return batch, nil
		}
		nbatches--
	}
	return btc.Batch{}, fmt.Errorf("no batch found")
}

func (m *mockCache) ReadRequests(ctx context.Context, ids ...string) ([]btc.BatcherRequest, error) {
	requests := []btc.BatcherRequest{}
	for _, id := range ids {
		request, ok := m.requests[id]
		if !ok {
			return nil, fmt.Errorf("request not found")
		}
		requests = append(requests, request)
	}
	return requests, nil
}

func (m *mockCache) DeletePendingBatches(ctx context.Context) error {
	newList := m.batchList
	for i, id := range m.batchList {
		if m.batches[id].Strategy != m.mode {
			continue
		}

		if m.batches[id].Tx.Status.Confirmed == false {
			delete(m.batches, id)
			newList = append(newList[:i], newList[i+1:]...)
		}
	}

	m.batchList = newList
	return nil
}

type mockFeeEstimator struct {
	fee int
}

func (f *mockFeeEstimator) UpdateFee(newFee int) {
	f.fee = newFee
}

func (f *mockFeeEstimator) FeeSuggestion() (btc.FeeSuggestion, error) {
	return btc.FeeSuggestion{
		Minimum: f.fee,
		Economy: f.fee,
		Low:     f.fee,
		Medium:  f.fee,
		High:    f.fee,
	}, nil
}

func NewMockFeeEstimator(fee int) *mockFeeEstimator {
	return &mockFeeEstimator{
		fee: fee,
	}
}

// mockIndexer is a minimal IndexerClient: only GetUTXOs returns meaningful
// data. validateBatchRequest with empty spends/sacps never touches the
// other methods, so they are benign stubs.
type mockIndexer struct {
	utxos btc.UTXOs
}

func (m *mockIndexer) GetUTXOs(_ context.Context, _ btcutil.Address) (btc.UTXOs, error) {
	return m.utxos, nil
}
func (m *mockIndexer) GetAddressTxs(_ context.Context, _ btcutil.Address, _ string) ([]btc.Transaction, error) {
	return nil, nil
}
func (m *mockIndexer) GetUTXOsForAmount(_ context.Context, _ btcutil.Address, _ int64) (btc.UTXOs, int64, error) {
	return nil, 0, nil
}
func (m *mockIndexer) GetTipBlockHeight(_ context.Context) (uint64, error) { return 0, nil }
func (m *mockIndexer) GetTx(_ context.Context, _ string) (btc.Transaction, error) {
	return btc.Transaction{}, nil
}
func (m *mockIndexer) GetTxHex(_ context.Context, _ string) (string, error) { return "", nil }
func (m *mockIndexer) SubmitTx(_ context.Context, _ *wire.MsgTx) error      { return nil }
func (m *mockIndexer) FeeEstimate(_ context.Context) (btc.FeeSuggestion, error) {
	return btc.FeeSuggestion{}, nil
}

// TestValidateBatchRequestExcludesCommittedUTXOs verifies that UTXOs
// already spent by the latest in-flight batch transaction are excluded
// from the wallet balance computed inside validateBatchRequest. Without
// the exclusion, the committed UTXO is double-counted and the send is
// admitted; with the exclusion, only the free UTXO is counted and the
// request is rejected as insufficient funds.
func TestValidateBatchRequestExcludesCommittedUTXOs(t *testing.T) {
	ctx := context.Background()
	chainParams := &chaincfg.RegressionNetParams

	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatalf("generate priv key: %v", err)
	}
	recipKey, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatalf("generate recipient key: %v", err)
	}
	recipAddr, err := btc.PublicKeyAddress(chainParams, waddrmgr.WitnessPubKey, recipKey.PubKey())
	if err != nil {
		t.Fatalf("recipient address: %v", err)
	}

	// Two UTXOs on the wallet: tx1:0 is already spent by the latest batch,
	// tx2:0 is the only free UTXO.
	indexer := &mockIndexer{
		utxos: btc.UTXOs{
			{TxID: "tx1", Vout: 0, Amount: 10_000, Status: &btc.Status{Confirmed: true}},
			{TxID: "tx2", Vout: 0, Amount: 5_000, Status: &btc.Status{Confirmed: true}},
		},
	}

	cache := NewTestCache(btc.CPFP)
	latestBatch := btc.Batch{
		Tx: btc.Transaction{
			TxID: "latestBatchTx",
			VINs: []btc.VIN{{TxID: "tx1", Vout: 0}},
		},
		Strategy:   btc.CPFP,
		RequestIds: map[string]bool{},
	}
	if err := cache.SaveBatch(ctx, latestBatch); err != nil {
		t.Fatalf("seed latest batch: %v", err)
	}

	rpc := btc.NewBitcoinClient("", "", "")
	wallet, err := btc.NewBatcherWallet(
		privKey, indexer, NewMockFeeEstimator(10), chainParams,
		cache, zap.NewNop(), &rpc, btc.WithStrategy(btc.CPFP),
	)
	if err != nil {
		t.Fatalf("new batcher wallet: %v", err)
	}

	// Totals (5_000 is the 1_000-buffer check in validateBatchRequest):
	//   unfiltered: in = 15_000, out = 6_000 → passes (bug)
	//   filtered:   in =  5_000, out = 6_000 → rejects (fixed)
	_, err = wallet.Send(ctx, []btc.SendRequest{{Amount: 6_000, To: recipAddr}}, nil, nil)

	if err == nil {
		t.Fatalf("expected insufficient-funds error (committed UTXO should have been excluded), got nil")
	}
	if !strings.Contains(err.Error(), "batch parameters not met") {
		t.Fatalf("expected ErrBatchParametersNotMet, got: %v", err)
	}
}
