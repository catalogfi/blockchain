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
}

func NewTestCache() btc.Cache {
	return &mockCache{
		batches:  make(map[string]btc.Batch),
		requests: make(map[string]btc.BatcherRequest),
	}
}

func (m *mockCache) ReadBatchByReqId(ctx context.Context, id string) (btc.Batch, error) {
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

func (m *mockCache) ReadPendingBatches(ctx context.Context, strategy btc.Strategy) ([]btc.Batch, error) {
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

func (m *mockCache) ConfirmBatchStatuses(ctx context.Context, txIds []string, deletePending bool, strategy btc.Strategy) error {
	confirmedBatchIds := make(map[string]bool)
	for _, id := range txIds {
		batch, ok := m.batches[id]
		if !ok {
			return fmt.Errorf("UpdateBatchStatuses, batch not found")
		}

		confirmedBatchIds[id] = true

		batch.Tx.Status.Confirmed = true
		m.batches[id] = batch
	}
	if deletePending {
		return m.DeletePendingBatches(ctx, confirmedBatchIds, strategy)
	}
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

func (m *mockCache) SaveRequest(ctx context.Context, id string, req btc.BatcherRequest) error {
	if _, ok := m.requests[id]; ok {
		return fmt.Errorf("request already exists")
	}
	m.requests[id] = req
	m.requestList = append(m.requestList, id)
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

func (m *mockCache) ReadLatestBatch(ctx context.Context, strategy btc.Strategy) (btc.Batch, error) {
	if len(m.batchList) == 0 {
		return btc.Batch{}, btc.ErrBatchNotFound
	}
	nbatches := len(m.batchList) - 1
	for nbatches >= 0 {
		batch, ok := m.batches[m.batchList[nbatches]]
		if ok && batch.Strategy == strategy {
			return batch, nil
		}
		nbatches--
	}
	return btc.Batch{}, fmt.Errorf("no batch found")
}

func (m *mockCache) ReadRequests(ctx context.Context, ids []string) ([]btc.BatcherRequest, error) {
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

func (m *mockCache) DeletePendingBatches(ctx context.Context, confirmedBatchIds map[string]bool, strategy btc.Strategy) error {
	newList := m.batchList
	for i, id := range m.batchList {
		if m.batches[id].Strategy != strategy {
			continue
		}

		if _, ok := confirmedBatchIds[id]; ok {
			batch := m.batches[id]
			batch.Tx.Status.Confirmed = true
			m.batches[id] = batch
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

func (m *mockCache) ReadPendingChangeUtxos(ctx context.Context, strategy btc.Strategy) ([]btc.UTXO, error) {
	utxos := []btc.UTXO{}
	for _, id := range m.batchList {
		if m.batches[id].Strategy == strategy && m.batches[id].Tx.Status.Confirmed == false {
			utxos = append(utxos, m.batches[id].SelfUtxos...)
		}
	}
	return utxos, nil
}
func (m *mockCache) ReadPendingFundingUtxos(ctx context.Context, strategy btc.Strategy) ([]btc.UTXO, error) {
	utxos := []btc.UTXO{}
	for _, id := range m.batchList {
		if m.batches[id].Strategy == strategy && (m.batches[id].Tx.Status.Confirmed == false) {
			utxos = append(utxos, m.batches[id].FundingUtxos...)
		}
	}
	return utxos, nil
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
