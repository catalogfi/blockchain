package btc_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/catalogfi/blockchain/btc"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/exp/maps"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func saveBatch(batch btc.Batch, cache btc.Cache) error {
	for reqId := range batch.RequestIds {
		req := dummyRequest()
		req.ID = reqId
		err := cache.SaveRequest(context.Background(), req)
		if err != nil {
			return err
		}
	}
	return cache.SaveBatch(context.Background(), batch)
}

var _ = Describe("BatcherCache", Ordered, func() {

	dbPath := "./testdb"
	var db *leveldb.DB
	var cache btc.Cache
	var err error

	ctx := context.Background()

	BeforeAll(func() {
		var err error
		db, err = leveldb.OpenFile(dbPath, nil)
		Expect(err).To(BeNil())

		cache = btc.NewBatcherCache(db, btc.CPFP)
	})

	AfterAll(func() {
		db.Close()
		// remove db path
		err := os.RemoveAll(dbPath)
		Expect(err).To(BeNil())
	})

	Describe("SaveBatch", func() {
		It("should save a pending batch", func() {
			batchToSave := dummyBatch()
			err := saveBatch(batchToSave, cache)
			Expect(err).To(BeNil())

			// check if batch is saved
			b, err := cache.ReadPendingBatches(ctx)
			Expect(err).To(BeNil())
			Expect(len(b)).To(Equal(1))

			Expect(b[0].Tx.TxID).To(Equal(batchToSave.Tx.TxID))
			Expect(b[0].IsFinalized).To(Equal(batchToSave.IsFinalized))
			Expect(b[0].Strategy).To(Equal(batchToSave.Strategy))
		})
		It("should not save a batch that already exists", func() {
			batchToSave := dummyBatch()
			err = saveBatch(batchToSave, cache)
			Expect(err).To(BeNil())

			err = cache.SaveBatch(ctx, batchToSave)
			Expect(err).To(Equal(btc.ErrStoreAlreadyExists))
		})
		It("should save a finalized batch", func() {
			err = cache.DeletePendingBatches(ctx)
			Expect(err).To(BeNil())

			batchToSave := dummyBatch()
			batchToSave.Tx.Status.Confirmed = true
			err = saveBatch(batchToSave, cache)

			Expect(err).To(BeNil())

			// batch should not be saved in pending
			b, err := cache.ReadPendingBatches(ctx)
			Expect(err).To(BeNil())
			Expect(len(b)).To(Equal(0))

			batch, err := cache.ReadLatestBatch(ctx)
			Expect(err).To(BeNil())
			Expect(batch.Tx.TxID).To(Equal(batchToSave.Tx.TxID))
		})
	})

	Describe("ReadBatchByReqId", func() {

		It("should read a batch by request id", func() {
			batchToSave := dummyBatch()
			reqId := maps.Keys(batchToSave.RequestIds)[0]
			err := saveBatch(batchToSave, cache)
			Expect(err).To(BeNil())

			batch, err := cache.ReadBatchByReqID(ctx, reqId)
			Expect(err).To(BeNil())
			Expect(batch.Tx.TxID).To(Equal(batchToSave.Tx.TxID))

			// check for invalid request id
			_, err = cache.ReadBatchByReqID(ctx, "invalid")
			Expect(err).To(Equal(btc.ErrStoreNotFound))
		})
	})

	Describe("ReadPendingBatches", func() {

		It("should read pending batches", func() {
			err := cache.DeletePendingBatches(ctx)
			Expect(err).To(BeNil())

			batch1ToSave := dummyBatch()
			err = saveBatch(batch1ToSave, cache)
			Expect(err).To(BeNil())

			batch2ToSave := dummyBatch()
			err = saveBatch(batch2ToSave, cache)
			Expect(err).To(BeNil())

			b, err := cache.ReadPendingBatches(ctx)
			Expect(err).To(BeNil())

			Expect(len(b)).To(Equal(2))

			Expect(b[0].Tx.TxID).To(Equal(batch1ToSave.Tx.TxID))
			Expect(b[1].Tx.TxID).To(Equal(batch2ToSave.Tx.TxID))
		})
	})

	Describe("ReadLatestBatch", func() {
		It("should read the latest batch", func() {
			dummyBatch1 := dummyBatch()
			err := saveBatch(dummyBatch1, cache)
			Expect(err).To(BeNil())

			batch, err := cache.ReadLatestBatch(ctx)
			Expect(err).To(BeNil())
			Expect(batch.Tx.TxID).To(Equal(dummyBatch1.Tx.TxID))

			dummyBatch2 := dummyBatch()
			err = saveBatch(dummyBatch2, cache)
			Expect(err).To(BeNil())

			batch, err = cache.ReadLatestBatch(ctx)
			Expect(err).To(BeNil())
			Expect(batch.Tx.TxID).To(Equal(dummyBatch2.Tx.TxID))
		})
	})

	Describe("UpdateBatches", func() {
		It("should update batches", func() {
			batchToSave := dummyBatch()
			err = saveBatch(batchToSave, cache)
			Expect(err).To(BeNil())

			batchToSave.Tx.Status.Confirmed = true
			err = cache.UpdateBatches(ctx, batchToSave)
			Expect(err).To(BeNil())

			b, err := cache.ReadBatch(ctx, batchToSave.Tx.TxID)
			Expect(err).To(BeNil())
			Expect(b.Tx.Status.Confirmed).To(BeTrue())
		})
		It("should be able to create a batch that does not exist", func() {
			batchToSave := dummyBatch()
			err := cache.UpdateBatches(ctx, batchToSave)
			Expect(err).To(BeNil())

			b, err := cache.ReadBatch(ctx, batchToSave.Tx.TxID)
			Expect(err).To(BeNil())
			Expect(b.Tx.TxID).To(Equal(batchToSave.Tx.TxID))

		})
	})

	Describe("UpdateAndDeletePendingBatches", func() {
		It("should update and delete pending batches", func() {
			batch := dummyBatch()
			err := saveBatch(batch, cache)
			Expect(err).To(BeNil())

			batch2 := dummyBatch()
			err = saveBatch(batch2, cache)
			Expect(err).To(BeNil())

			err = cache.UpdateAndDeletePendingBatches(ctx, batch)
			Expect(err).To(BeNil())

			// batch2 should not exist
			b, err := cache.ReadPendingBatches(ctx)
			Expect(err).To(BeNil())
			for _, b := range b {
				Expect(b.Tx.TxID).NotTo(Equal(batch2.Tx.TxID))
			}
		})

		It("should return an error if nothing to update", func() {
			err := cache.UpdateAndDeletePendingBatches(ctx)
			Expect(err).To(Equal(btc.ErrStoreNothingToUpdate))
		})

		It("should create a batch that does not exist", func() {
			batch := dummyBatch()
			err := cache.UpdateAndDeletePendingBatches(ctx, batch)
			Expect(err).To(BeNil())

			b, err := cache.ReadBatch(ctx, batch.Tx.TxID)
			Expect(err).To(BeNil())
			Expect(b.Tx.TxID).To(Equal(batch.Tx.TxID))
		})

		It("should delete all pending batches && requests", func() {

			err := cache.DeletePendingBatches(ctx)
			Expect(err).To(BeNil())

			request := dummyRequest()
			err = cache.SaveRequest(ctx, request)
			Expect(err).To(BeNil())

			batch := dummyBatch()
			batch.RequestIds[request.ID] = true
			err = saveBatch(batch, cache)
			Expect(err).To(BeNil())

			batch2 := dummyBatch()
			err = saveBatch(batch2, cache)
			Expect(err).To(BeNil())

			batch.IsFinalized = true
			batch.Tx.Status.Confirmed = true
			err = cache.UpdateAndDeletePendingBatches(ctx, batch)
			Expect(err).To(BeNil())

			b, err := cache.ReadPendingBatches(ctx)
			Expect(err).To(BeNil())
			Expect(len(b)).To(Equal(0))

			reqs, err := cache.ReadPendingRequests(ctx)
			Expect(err).To(BeNil())
			Expect(len(reqs)).To(Equal(0))

			//make sure batch is updated
			updatedBatch, err := cache.ReadBatch(ctx, batch.Tx.TxID)
			Expect(err).To(BeNil())

			Expect(updatedBatch.IsFinalized).To(BeTrue())

		})
	})

	Describe("DeletePendingBatches", func() {
		It("should delete all pending batches", func() {
			batch := dummyBatch()
			err = saveBatch(batch, cache)
			Expect(err).To(BeNil())

			confirmedBatch := dummyBatch()
			confirmedBatch.Tx.Status.Confirmed = true
			err = saveBatch(confirmedBatch, cache)
			Expect(err).To(BeNil())

			err = cache.DeletePendingBatches(ctx)
			Expect(err).To(BeNil())

			b, err := cache.ReadPendingBatches(ctx)
			Expect(err).To(BeNil())
			Expect(len(b)).To(Equal(0))

			cBatch, err := cache.ReadBatch(ctx, confirmedBatch.Tx.TxID)
			Expect(err).To(BeNil())

			Expect(cBatch.Tx.Status.Confirmed).To(Equal(true))
		})
	})

	Describe("ReadRequests", func() {
		It("should read requests", func() {
			request := dummyRequest()
			err := cache.SaveRequest(ctx, request)
			Expect(err).To(BeNil())

			reqs, err := cache.ReadRequests(ctx, request.ID)
			Expect(err).To(BeNil())
			Expect(len(reqs)).To(Equal(1))
			Expect(reqs[0].ID).To(Equal(request.ID))
			Expect(reqs[0].Spends[0].Utxos[0].TxID).To(Equal(request.Spends[0].Utxos[0].TxID))

		})
		It("should return an error if request not found", func() {
			_, err := cache.ReadRequests(ctx, "invalid")
			Expect(err).To(Equal(btc.ErrStoreNotFound))
		})
	})

	Describe("ReadPendingRequests", func() {
		It("should read pending requests", func() {

			existingReqs, err := cache.ReadPendingRequests(ctx)
			Expect(err).To(BeNil())

			request := dummyRequest()

			err = cache.SaveRequest(ctx, request)
			Expect(err).To(BeNil())

			reqs, err := cache.ReadPendingRequests(ctx)
			Expect(err).To(BeNil())
			Expect(len(reqs)).To(Equal(len(existingReqs) + 1))
			found := false
			for _, req := range reqs {
				if req.ID == request.ID {
					found = true
				}
			}
			Expect(found).To(BeTrue())

		})
	})

})

func dummyBatch() btc.Batch {
	reqIds := make(map[string]bool)
	reqIds["id1"] = true
	reqIds["id2"] = true
	return btc.Batch{
		Tx: btc.Transaction{
			TxID: fmt.Sprintf("%d", time.Now().UnixNano()),
			Status: btc.Status{
				Confirmed: false,
			},
		},
		RequestIds:  reqIds,
		IsFinalized: false,
		Strategy:    btc.CPFP,
	}
}

func confirmBatch(b btc.Batch) btc.Batch {
	batch := b
	batch.Tx.Status.Confirmed = true
	return batch
}

func dummyRequest() btc.BatcherRequest {
	return btc.BatcherRequest{
		ID:     fmt.Sprintf("%d", time.Now().UnixNano()),
		Status: false,
		Spends: []btc.SpendRequest{
			{
				Utxos: []btc.UTXO{
					{
						TxID:   "txid",
						Vout:   0,
						Amount: 1000,
					},
				},
				ScriptAddress: &btcutil.AddressPubKeyHash{},
			},
		},
	}
}
