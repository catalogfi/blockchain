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

var _ = Describe("BatchWallet:RBF", Ordered, func() {

	chainParams := chaincfg.RegressionNetParams
	logger, err := zap.NewDevelopment()
	Expect(err).To(BeNil())

	indexer := btc.NewElectrsIndexerClient(logger, os.Getenv("BTC_REGNET_INDEXER"), time.Millisecond*500)

	privateKey, err := btcec.NewPrivateKey()
	Expect(err).To(BeNil())

	mockFeeEstimator := NewMockFeeEstimator(10)
	cache := NewTestCache()
	wallet, err := btc.NewBatcherWallet(privateKey, indexer, mockFeeEstimator, &chainParams, cache, logger, btc.WithPTI(5*time.Second), btc.WithStrategy(btc.RBF))
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

	It("should be able to update fee with RBF", func() {
		mockFeeEstimator.UpdateFee(20)
		time.Sleep(10 * time.Second)
	})
})
