package btc_test

import (
	"context"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/multichain/btc"
	"github.com/catalogfi/multichain/testutil"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Indexer client", func() {
	Context("When using quicknode API", func() {
		It("should be able to fetch the utxos of an address ", func() {
			By("Initialise the quickNode client")
			logger, err := zap.NewDevelopment()
			Expect(err).To(BeNil())
			url := testutil.ParseStringEnv("BTC_INDEXER_QUICKNODE", "")
			client := btc.NewQuickNodeIndexerClient(logger, url)

			By("GetUTXOs()")
			addr, err := btcutil.DecodeAddress("mvb4nE8firm7abnss9r7j5tra2w7M7uktF", &chaincfg.TestNet3Params)
			Expect(err).To(BeNil())
			utxos, err := client.GetUTXOs(context.Background(), addr)
			Expect(err).To(BeNil())
			Expect(len(utxos)).Should(BeNumerically(">", 1))
		})
	})

	Context("When using electrs API", func() {
		It("should be able to fetch the utxos of an address ", func() {
			By("Initialise the quickNode client")
			logger, err := zap.NewDevelopment()
			Expect(err).To(BeNil())
			url := testutil.ParseStringEnv("BTC_INDEXER_QUICKNODE", "")
			client := btc.NewQuickNodeIndexerClient(logger, url)

			By("GetUTXOs()")
			addr, err := btcutil.DecodeAddress("mynaqNW3Pa2t3WC3jb9hr3hcFQYeVmyQNb", &chaincfg.TestNet3Params)
			Expect(err).To(BeNil())
			utxos, err := client.GetUTXOs(context.Background(), addr)
			Expect(err).To(BeNil())
			Expect(len(utxos)).Should(BeNumerically(">", 1))
		})
	})
})
