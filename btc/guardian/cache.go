package guardian

import (
	"context"
	"fmt"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

type Cache interface {
	SaveLatestBatch(ctx context.Context, batch *Batch) error
	ReadLatestBatch(ctx context.Context) (*Batch, error)
	SaveMergeTxFee(ctx context.Context, fee int64) error
	ReadMergeTxFee(ctx context.Context) (int64, error)
	SaveRequestVOUT(ctx context.Context, reqID string, vout int) error
	ReadRequestVOUT(ctx context.Context, reqID string) (int, error)
	ReadBatchByTxID(ctx context.Context, txID string) (Batch, error)
	UpdateFailedTxId(ctx context.Context, txID string, batch *Batch) error
}

const (
	latestBatchKey = "latest_batch"
	mergeTxFeeKey  = "merge_tx_fee"
	requestVoutKey = "request_vout"
	batchTxIDKey   = "batch_txid"
)

type cache struct {
	db *leveldb.DB
}

func NewCache(db *leveldb.DB) *cache {
	return &cache{
		db: db,
	}
}

func (l *cache) ReadBatchByTxID(ctx context.Context, txID string) (Batch, error) {
	batchBytes, err := l.db.Get([]byte(fmt.Sprintf("%s_%s", batchTxIDKey, txID)), nil)
	if err != nil {
		return Batch{}, fmt.Errorf("failed to get batch by tx id: %w", err)
	}
	batch := &Batch{}
	err = batch.Unmarshal(batchBytes)
	if err != nil {
		return Batch{}, fmt.Errorf("failed to unmarshal batch: %w", err)
	}
	return *batch, nil
}

func (c *cache) SaveLatestBatch(ctx context.Context, batch *Batch) error {
	batchBytes, err := batch.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal batch: %w", err)
	}
	// use leveldb batch to save both the latest batch and the batch txid
	levelBatch := leveldb.Batch{}
	levelBatch.Put([]byte(latestBatchKey), batchBytes)
	levelBatch.Put([]byte(fmt.Sprintf("%s_%s", batchTxIDKey, batch.Tx.TxID)), batchBytes)
	err = c.db.Write(&levelBatch, nil)
	if err != nil {
		return fmt.Errorf("failed to save latest batch: %w", err)
	}
	return nil
}

func (c *cache) ReadLatestBatch(ctx context.Context) (*Batch, error) {
	batchBytes, err := c.db.Get([]byte(latestBatchKey), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest batch: %w", err)
	}
	batch := &Batch{}
	err = batch.Unmarshal(batchBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch: %w", err)
	}
	return batch, nil
}

func (c *cache) SaveMergeTxFee(ctx context.Context, fee int64) error {
	return c.db.Put([]byte(mergeTxFeeKey), []byte(fmt.Sprintf("%d", fee)), nil)
}

func (c *cache) ReadMergeTxFee(ctx context.Context) (int64, error) {
	feeBytes, err := c.db.Get([]byte(mergeTxFeeKey), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get merge tx fee: %w", err)
	}
	return strconv.ParseInt(string(feeBytes), 10, 64)
}

func (c *cache) SaveRequestVOUT(ctx context.Context, reqID string, vout int) error {
	return c.db.Put([]byte(fmt.Sprintf("%s_%s_vout", requestVoutKey, reqID)), []byte(fmt.Sprintf("%d", vout)), nil)
}

func (c *cache) ReadRequestVOUT(ctx context.Context, reqID string) (int, error) {
	voutBytes, err := c.db.Get([]byte(fmt.Sprintf("%s_%s_vout", requestVoutKey, reqID)), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get request vout: %w", err)
	}
	return strconv.Atoi(string(voutBytes))
}

func (c *cache) UpdateFailedTxId(ctx context.Context, txID string, batch *Batch) error {
	// update LastFailedTxHash in the batch with the txID
	batch.LastFailedTxHash = txID
	batchBytes, err := batch.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal batch: %w", err)
	}
	// use leveldb batch to save both the latest batch and the batch txid
	levelBatch := leveldb.Batch{}
	levelBatch.Put([]byte(latestBatchKey), batchBytes)
	levelBatch.Put([]byte(fmt.Sprintf("%s_%s", batchTxIDKey, batch.Tx.TxID)), batchBytes)
	err = c.db.Write(&levelBatch, nil)
	if err != nil {
		return fmt.Errorf("failed to save latest batch: %w", err)
	}
	return nil
}
