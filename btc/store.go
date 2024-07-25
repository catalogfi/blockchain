package btc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	ErrStoreNotFound        = errors.New("store: not found")
	ErrStoreAlreadyExists   = errors.New("store: already exists")
	ErrStoreNothingToUpdate = errors.New("store: nothing to update")
)

// serializableSpendRequest is a serializable version of SpendRequest
type serializableSpendRequest struct {
	Witness       [][]byte
	Script        []byte
	Leaf          []byte
	ScriptAddress string
	HashType      txscript.SigHashType
	Sequence      uint32
}

// serializableSendRequest is a serializable version of SendRequest
type serializableSendRequest struct {
	Amount int64
	To     string
}

// serializableBatcherRequest is a serializable version of BatcherRequest
type serializableBatcherRequest struct {
	ID     string
	Spends []serializableSpendRequest
	Sends  []serializableSendRequest
	SACPs  [][]byte
	Status bool
}

// serializeBatcherRequest serializes a BatcherRequest to a byte slice
func serializeBatcherRequest(req BatcherRequest) ([]byte, error) {
	primitiveReq := serializableBatcherRequest{
		ID:     req.ID,
		Spends: make([]serializableSpendRequest, len(req.Spends)),
		Sends:  make([]serializableSendRequest, len(req.Sends)),
		SACPs:  req.SACPs,
		Status: req.Status,
	}

	for i, spend := range req.Spends {
		primitiveReq.Spends[i] = serializableSpendRequest{
			Witness:       spend.Witness,
			Script:        spend.Script,
			Leaf:          spend.Leaf.Script,
			ScriptAddress: spend.ScriptAddress.EncodeAddress(),
			HashType:      spend.HashType,
			Sequence:      spend.Sequence,
		}
	}

	for i, send := range req.Sends {
		primitiveReq.Sends[i] = serializableSendRequest{
			Amount: send.Amount,
			To:     send.To.EncodeAddress(),
		}
	}
	return json.Marshal(primitiveReq)
}

// deserializeBatcherRequest deserializes a byte slice to a BatcherRequest
func deserializeBatcherRequest(data []byte) (BatcherRequest, error) {
	var primitiveReq serializableBatcherRequest
	err := json.Unmarshal(data, &primitiveReq)
	if err != nil {
		return BatcherRequest{}, err
	}

	req := BatcherRequest{
		ID:     primitiveReq.ID,
		Spends: make([]SpendRequest, len(primitiveReq.Spends)),
		Sends:  make([]SendRequest, len(primitiveReq.Sends)),
		SACPs:  primitiveReq.SACPs,
		Status: primitiveReq.Status,
	}

	for i, spend := range primitiveReq.Spends {
		addr, err := parseAddress(spend.ScriptAddress)
		if err != nil {
			return BatcherRequest{}, err
		}
		req.Spends[i] = SpendRequest{
			Witness:       spend.Witness,
			Script:        spend.Script,
			Leaf:          txscript.NewTapLeaf(0xc0, spend.Leaf),
			ScriptAddress: addr,
			HashType:      spend.HashType,
			Sequence:      spend.Sequence,
		}
	}

	for i, send := range primitiveReq.Sends {
		addr, err := parseAddress(send.To)
		if err != nil {
			return BatcherRequest{}, err
		}
		req.Sends[i] = SendRequest{
			Amount: send.Amount,
			To:     addr,
		}
	}
	return req, nil
}

// serializableBatch is a serializable version of Batch
type serializableBatch struct {
	Tx          Transaction
	RequestIds  []string
	IsFinalized bool
	Strategy    Strategy
}

// serializeBatch serializes a Batch to a byte slice
func serializeBatch(batch Batch) ([]byte, error) {
	primitiveBatch := serializableBatch{
		Tx:          batch.Tx,
		RequestIds:  make([]string, 0, len(batch.RequestIds)),
		IsFinalized: batch.IsFinalized,
		Strategy:    batch.Strategy,
	}

	for id := range batch.RequestIds {
		primitiveBatch.RequestIds = append(primitiveBatch.RequestIds, id)
	}

	return json.Marshal(primitiveBatch)
}

// deserializeBatch deserializes a byte slice to a Batch
func deserializeBatch(data []byte) (Batch, error) {
	var primitiveBatch serializableBatch
	err := json.Unmarshal(data, &primitiveBatch)
	if err != nil {
		return Batch{}, err
	}

	batch := Batch{
		Tx:          primitiveBatch.Tx,
		RequestIds:  make(map[string]bool, len(primitiveBatch.RequestIds)),
		IsFinalized: primitiveBatch.IsFinalized,
		Strategy:    primitiveBatch.Strategy,
	}

	for _, id := range primitiveBatch.RequestIds {
		batch.RequestIds[id] = true
	}

	return batch, nil
}

// BatcherCache is a cache implementation for the batcher
type BatcherCache struct {
	db       *leveldb.DB
	strategy Strategy
	*batcherCacheKeyManager
}

func NewBatcherCache(db *leveldb.DB, strategy Strategy) Cache {
	return &BatcherCache{
		db:                     db,
		strategy:               strategy,
		batcherCacheKeyManager: &batcherCacheKeyManager{strategy: strategy},
	}
}

func (l *BatcherCache) get(key []byte) ([]byte, error) {
	data, err := l.db.Get(key, nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, ErrStoreNotFound
		}
		return nil, err
	}
	return data, nil
}

func (l *BatcherCache) ReadBatchByReqID(_ context.Context, reqID string) (Batch, error) {
	batchID, err := l.get(l.requestIndexKey(reqID))
	if err != nil {
		return Batch{}, err
	}
	return l.searchBatch(string(batchID))
}

func (l *BatcherCache) searchBatch(id string) (Batch, error) {
	// search in pending batches
	batch, err := l.getBatch(id, true)
	if err != nil {
		if !errors.Is(err, ErrStoreNotFound) {
			return Batch{}, err
		}
	} else if batch.Tx.TxID == id {
		return batch, nil
	}

	// search in finalized batches
	batch, err = l.getBatch(id, false)
	if err != nil {
		return Batch{}, err
	}
	return batch, nil
}

func (l *BatcherCache) ReadPendingBatches(_ context.Context) ([]Batch, error) {
	iter := l.db.NewIterator(util.BytesPrefix(l.pendingBatchKey("")), nil)
	defer iter.Release()
	var batches []Batch
	for iter.Next() {
		batch, err := deserializeBatch(iter.Value())
		if err != nil {
			return nil, err
		}
		batches = append(batches, batch)
	}
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return batches, nil
}

func (l *BatcherCache) ReadLatestBatch(_ context.Context) (Batch, error) {
	data, err := l.db.Get(l.latestBatchKey(), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return Batch{}, ErrStoreNotFound
		}
		return Batch{}, err
	}
	return deserializeBatch(data)
}

func (l *BatcherCache) UpdateBatches(_ context.Context, updatedBatches ...Batch) error {
	batch := new(leveldb.Batch)
	for _, b := range updatedBatches {
		data, err := serializeBatch(b)
		if err != nil {
			return err
		}
		if isPending(b) {
			batch.Put(l.pendingBatchKey(b.Tx.TxID), data)
		} else {
			batch.Put(l.batchKey(b.Tx.TxID), data)
			batch.Delete(l.pendingBatchKey(b.Tx.TxID))

		}
	}
	return l.db.Write(batch, nil)
}

func (l *BatcherCache) ReadBatch(_ context.Context, id string) (Batch, error) {
	return l.searchBatch(id)
}

func (l *BatcherCache) UpdateAndDeletePendingBatches(_ context.Context, updatedBatches ...Batch) error {
	batch := new(leveldb.Batch)
	if len(updatedBatches) == 0 {
		return ErrStoreNothingToUpdate
	}

	// delete pending batches
	iter := l.db.NewIterator(util.BytesPrefix(l.pendingBatchKey("")), nil)
	defer iter.Release()
	for iter.Next() {
		batch.Delete(iter.Key())
	}

	// update pending batches
	for _, b := range updatedBatches {
		data, err := serializeBatch(b)
		if err != nil {
			return err
		}
		if isPending(b) {
			batch.Put(l.pendingBatchKey(b.Tx.TxID), data)
		} else {
			batch.Put(l.batchKey(b.Tx.TxID), data)

		}
	}
	return l.db.Write(batch, nil)
}

func (l *BatcherCache) getPendingRequestBytes(id string) ([]byte, error) {
	return l.get(l.pendingRequestKey(id))
}

func (l *BatcherCache) DeletePendingBatches(_ context.Context) error {
	iter := l.db.NewIterator(util.BytesPrefix(l.pendingBatchKey("")), nil)
	defer iter.Release()
	batch := new(leveldb.Batch)
	for iter.Next() {
		batch.Delete(iter.Key())
	}
	if err := iter.Error(); err != nil {
		return err
	}
	return l.db.Write(batch, nil)
}

func (l *BatcherCache) getBatch(id string, isPending bool) (Batch, error) {
	var data []byte
	var err error
	if isPending {
		data, err = l.get(l.pendingBatchKey(id))
	} else {
		data, err = l.get(l.batchKey(id))
	}
	if err != nil {
		return Batch{}, err
	}

	batch, err := deserializeBatch(data)
	if err != nil {
		return Batch{}, err
	}
	return batch, nil
}

func (l *BatcherCache) SaveBatch(ctx context.Context, batch Batch) error {
	isPending := isPending(batch)
	existinBatch, err := l.getBatch(batch.Tx.TxID, isPending)
	if err == nil && existinBatch.Tx.TxID == batch.Tx.TxID {
		return ErrStoreAlreadyExists
	}
	if !errors.Is(err, ErrStoreNotFound) {
		return err
	}
	data, err := serializeBatch(batch)
	if err != nil {
		return err
	}

	levelDBBatch := new(leveldb.Batch)

	// Save the batch
	if isPending {
		levelDBBatch.Put(l.pendingBatchKey(batch.Tx.TxID), data)
	} else {
		levelDBBatch.Put(l.batchKey(batch.Tx.TxID), data)
	}

	// Once a batch is created, we need to move all the pending
	// requests of the batch to finalized requests
	for id := range batch.RequestIds {
		reqs, err := l.ReadRequests(ctx, id)
		if err != nil {
			return fmt.Errorf("error getting request %s: %w", id, err)
		}
		req := reqs[0]
		if req.Status {
			continue
		}
		req.Status = true
		data, err := serializeBatcherRequest(req)
		if err != nil {
			return err
		}
		levelDBBatch.Delete(l.pendingRequestKey(id))
		levelDBBatch.Put(l.requestKey(id), data)
	}
	// save the latest batch
	levelDBBatch.Put(l.latestBatchKey(), data)

	// Index the requests to the batch
	for id := range batch.RequestIds {
		levelDBBatch.Put(l.requestIndexKey(id), []byte(batch.Tx.TxID))
	}

	return l.db.Write(levelDBBatch, nil)
}

func (l *BatcherCache) getPendingRequest(id string) (BatcherRequest, error) {
	data, err := l.getPendingRequestBytes(id)
	if err != nil {
		return BatcherRequest{}, err
	}
	return deserializeBatcherRequest(data)
}

func (l *BatcherCache) saveLatestBatch(batch Batch) error {
	data, err := serializeBatch(batch)
	if err != nil {
		return err
	}
	return l.db.Put([]byte(fmt.Sprintf(latestBatchPrefix, l.strategy)), data, nil)
}

func (l *BatcherCache) searchRequest(id string) (BatcherRequest, error) {
	// searching in pending requests
	req, err := l.getPendingRequest(id)
	if err != nil {
		if !errors.Is(err, ErrStoreNotFound) {
			return BatcherRequest{}, err
		}
	} else {
		return req, nil
	}

	// searching in finalized requests
	data, err := l.get(l.requestKey(id))
	if err != nil {
		return BatcherRequest{}, err
	}
	return deserializeBatcherRequest(data)
}

func (l *BatcherCache) ReadRequests(_ context.Context, ids ...string) ([]BatcherRequest, error) {
	var requests []BatcherRequest
	for _, id := range ids {
		req, err := l.searchRequest(id)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (l *BatcherCache) ReadPendingRequests(_ context.Context) ([]BatcherRequest, error) {
	iter := l.db.NewIterator(util.BytesPrefix(l.pendingRequestKey("")), nil)
	defer iter.Release()
	var requests []BatcherRequest
	for iter.Next() {
		req, err := deserializeBatcherRequest(iter.Value())
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return requests, nil
}

func (l *BatcherCache) SaveRequest(_ context.Context, req BatcherRequest) error {
	data, err := serializeBatcherRequest(req)
	if err != nil {
		return err
	}
	return l.db.Put(l.pendingRequestKey(req.ID), data, nil)
}

const (
	pendingBatchPrefix = "%s_pending_batch_%s"
	batchPrefix        = "%s_batch_%s"
	latestBatchPrefix  = "%s_latest_batch"
	reqIndexPrefix     = "%s_request_idx_%s"
	requestKey         = "%s_request"
	pendingRequestKey  = "%s_pending_request_%s"
)

func isPending(batch Batch) bool {
	return !batch.Tx.Status.Confirmed
}

// All levelDB keys are managed by this struct
type batcherCacheKeyManager struct {
	strategy Strategy
}

func (b *batcherCacheKeyManager) requestIndexKey(reqID string) []byte {
	return []byte(fmt.Sprintf(reqIndexPrefix, b.strategy, reqID))
}

func (b *batcherCacheKeyManager) pendingBatchKey(batchID string) []byte {
	return []byte(fmt.Sprintf(pendingBatchPrefix, b.strategy, batchID))
}

func (b *batcherCacheKeyManager) batchKey(batchID string) []byte {
	return []byte(fmt.Sprintf(batchPrefix, b.strategy, batchID))
}

func (b *batcherCacheKeyManager) latestBatchKey() []byte {
	return []byte(fmt.Sprintf(latestBatchPrefix, b.strategy))
}

func (b *batcherCacheKeyManager) requestKey(reqID string) []byte {
	return []byte(fmt.Sprintf(requestKey, reqID))
}

func (b *batcherCacheKeyManager) pendingRequestKey(reqID string) []byte {
	return []byte(fmt.Sprintf(pendingRequestKey, b.strategy, reqID))
}
