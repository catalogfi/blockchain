package btc_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/localnet"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("BatchWallet:CPFP", Ordered, func() {

	chainParams := chaincfg.RegressionNetParams
	logger, err := zap.NewDevelopment()
	Expect(err).To(BeNil())

	indexer := btc.NewElectrsIndexerClient(logger, os.Getenv("BTC_REGNET_INDEXER"), time.Millisecond*500)

	privateKey, err := btcec.NewPrivateKey()
	Expect(err).To(BeNil())

	mockFeeEstimator := NewMockFeeEstimator(10)
	cache := NewMockCache()
	wallet, err := btc.NewBatcherWallet(privateKey, indexer, mockFeeEstimator, &chainParams, cache, logger, btc.WithPTI(5*time.Second), btc.WithStrategy(btc.CPFP))
	Expect(err).To(BeNil())

	BeforeAll(func() {
		_, err := localnet.FundBitcoin(wallet.Address().EncodeAddress(), indexer)
		Expect(err).To(BeNil())
		err = wallet.Start(context.Background())
		Expect(err).To(BeNil())
	})

	AfterAll(func() {
		err := wallet.Stop()
		Expect(err).To(BeNil())
	})

	It("should be able to send funds", func() {
		req := []btc.SendRequest{
			{
				Amount: 100000,
				To:     wallet.Address(),
			},
		}

		id, err := wallet.Send(context.Background(), req, nil)
		Expect(err).To(BeNil())

		var tx btc.Transaction
		var ok bool

		for {
			fmt.Println("waiting for tx", id)
			tx, ok, err = wallet.Status(context.Background(), id)
			Expect(err).To(BeNil())
			if ok {
				Expect(tx).ShouldNot(BeNil())
				break
			}
			time.Sleep(5 * time.Second)
		}

		// to address
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
		// change address
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))

		time.Sleep(10 * time.Second)
	})

	It("should be able to update fee with CPFP", func() {
		mockFeeEstimator.UpdateFee(20)
		time.Sleep(10 * time.Second)
	})
})

type mockCache struct {
	batches  map[string]btc.Batch
	requests map[string]btc.BatcherRequest
}

func NewMockCache() btc.Cache {
	return &mockCache{
		batches:  make(map[string]btc.Batch),
		requests: make(map[string]btc.BatcherRequest),
	}
}

func (m *mockCache) ReadBatchByReqId(ctx context.Context, id string) (btc.Batch, error) {
	for _, batch := range m.batches {
		if _, ok := batch.RequestIds[id]; ok {
			return batch, nil
		}
	}
	return btc.Batch{}, fmt.Errorf("batch not found")
}

func (m *mockCache) ReadBatch(ctx context.Context, txId string) (btc.Batch, error) {
	batch, ok := m.batches[txId]
	if !ok {
		return btc.Batch{}, fmt.Errorf("batch not found")
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
	for id, _ := range batch.RequestIds {
		request := m.requests[id]
		request.Status = true
		m.requests[id] = request
	}
	return nil
}

func (m *mockCache) UpdateBatchStatuses(ctx context.Context, txId []string, status bool) error {
	for _, id := range txId {
		batch, ok := m.batches[id]
		if !ok {
			return fmt.Errorf("batch not found")
		}
		batch.Tx.Status.Confirmed = status
		m.batches[id] = batch
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
	return nil
}

func (m *mockCache) UpdateBatchFees(ctx context.Context, txId []string, feeRate int64) error {
	for _, id := range txId {
		batch, ok := m.batches[id]
		if !ok {
			return fmt.Errorf("batch not found")
		}

		batch.Tx.Fee = int64(batch.Tx.Weight) * feeRate / 4
		m.batches[id] = batch
	}
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
