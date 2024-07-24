package btc_test

import (
	"context"
	"fmt"

	"github.com/catalogfi/blockchain/btc"
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
