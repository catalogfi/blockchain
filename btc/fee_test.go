package btc_test

import (
	"math/rand"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/multichain/btc"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fee estimator", func() {
	Context("mempool estimator", func() {
		It("should return the live data from mempool", func() {
			estimator := btc.NewMempoolFeeEstimator(&chaincfg.MainNetParams, btc.MempoolFeeAPI, 15*time.Second)
			fees, err := estimator.FeeSuggestion()
			Expect(err).Should(BeNil())

			Expect(fees.Minimum).Should(BeNumerically(">", 1))

			Expect(fees.Economy).Should(BeNumerically(">", 1))
			Expect(fees.Economy).Should(BeNumerically(">=", fees.Minimum))

			Expect(fees.Low).Should(BeNumerically(">", 1))
			Expect(fees.Low).Should(BeNumerically(">=", fees.Economy))

			Expect(fees.Medium).Should(BeNumerically(">", 1))
			Expect(fees.Medium).Should(BeNumerically(">=", fees.Low))

			Expect(fees.High).Should(BeNumerically(">", 1))
			Expect(fees.High).Should(BeNumerically(">=", fees.Medium))
		})

		It("should use the cached result if the more requests are made during certain time period", func() {
			estimator := btc.NewMempoolFeeEstimator(&chaincfg.MainNetParams, btc.MempoolFeeAPI, time.Second)
			fees, err := estimator.FeeSuggestion()
			Expect(err).Should(BeNil())

			time.Sleep(500 * time.Millisecond)

			fees1, err := estimator.FeeSuggestion()
			Expect(fees.Minimum).Should(Equal(fees1.Minimum))
			Expect(fees.Economy).Should(Equal(fees1.Economy))
			Expect(fees.Low).Should(Equal(fees1.Low))
			Expect(fees.Medium).Should(Equal(fees1.Medium))
			Expect(fees.High).Should(Equal(fees1.High))

			time.Sleep(500 * time.Millisecond)
			_, err = estimator.FeeSuggestion()
			Expect(err).Should(BeNil())
		})

		It("should return a fixed fee for non-mainnet network", func() {
			By("Testnet")
			estimatorTestnet := btc.NewMempoolFeeEstimator(&chaincfg.TestNet3Params, "", 15*time.Second)
			fees, err := estimatorTestnet.FeeSuggestion()
			Expect(err).Should(BeNil())
			Expect(fees.Minimum).Should(Equal(2))
			Expect(fees.Economy).Should(Equal(2))
			Expect(fees.Low).Should(Equal(2))
			Expect(fees.Medium).Should(Equal(2))
			Expect(fees.High).Should(Equal(2))

			By("Regnet")
			estimatorRegnet := btc.NewMempoolFeeEstimator(&chaincfg.RegressionNetParams, "", 15*time.Second)
			fees, err = estimatorRegnet.FeeSuggestion()
			Expect(err).Should(BeNil())
			Expect(fees.Minimum).Should(Equal(2))
			Expect(fees.Economy).Should(Equal(2))
			Expect(fees.Low).Should(Equal(2))
			Expect(fees.Medium).Should(Equal(2))
			Expect(fees.High).Should(Equal(2))
		})
	})

	Context("blockstream estimator", func() {
		It("should return the live data from blockstream api", func() {
			estimator := btc.NewBlockstreamFeeEstimator(&chaincfg.MainNetParams, btc.BlockstreamAPI, 15*time.Second)
			fees, err := estimator.FeeSuggestion()
			Expect(err).Should(BeNil())

			Expect(fees.Minimum).Should(BeNumerically(">", 1))

			Expect(fees.Economy).Should(BeNumerically(">", 1))
			Expect(fees.Economy).Should(BeNumerically(">=", fees.Minimum))

			Expect(fees.Low).Should(BeNumerically(">", 1))
			Expect(fees.Low).Should(BeNumerically(">=", fees.Economy))

			Expect(fees.Medium).Should(BeNumerically(">", 1))
			Expect(fees.Medium).Should(BeNumerically(">=", fees.Low))

			Expect(fees.High).Should(BeNumerically(">", 1))
			Expect(fees.High).Should(BeNumerically(">=", fees.Medium))
		})

		It("should use the cached result if the more requests are made during certain time period", func() {
			estimator := btc.NewBlockstreamFeeEstimator(&chaincfg.MainNetParams, btc.BlockstreamAPI, time.Second)
			fees, err := estimator.FeeSuggestion()
			Expect(err).Should(BeNil())

			time.Sleep(500 * time.Millisecond)

			fees1, err := estimator.FeeSuggestion()
			Expect(fees.Minimum).Should(Equal(fees1.Minimum))
			Expect(fees.Economy).Should(Equal(fees1.Economy))
			Expect(fees.Low).Should(Equal(fees1.Low))
			Expect(fees.Medium).Should(Equal(fees1.Medium))
			Expect(fees.High).Should(Equal(fees1.High))

			time.Sleep(500 * time.Millisecond)
			_, err = estimator.FeeSuggestion()
			Expect(err).Should(BeNil())
		})

		It("should return a fixed fee for non-mainnet network", func() {
			By("Testnet")
			estimatorTestnet := btc.NewBlockstreamFeeEstimator(&chaincfg.TestNet3Params, "", 15*time.Second)
			fees, err := estimatorTestnet.FeeSuggestion()
			Expect(err).Should(BeNil())
			Expect(fees.Minimum).Should(Equal(2))
			Expect(fees.Economy).Should(Equal(2))
			Expect(fees.Low).Should(Equal(2))
			Expect(fees.Medium).Should(Equal(2))
			Expect(fees.High).Should(Equal(2))

			By("Regnet")
			estimatorRegnet := btc.NewBlockstreamFeeEstimator(&chaincfg.RegressionNetParams, "", 15*time.Second)
			fees, err = estimatorRegnet.FeeSuggestion()
			Expect(err).Should(BeNil())
			Expect(fees.Minimum).Should(Equal(2))
			Expect(fees.Economy).Should(Equal(2))
			Expect(fees.Low).Should(Equal(2))
			Expect(fees.Medium).Should(Equal(2))
			Expect(fees.High).Should(Equal(2))
		})
	})

	Context("fix fee estimator", func() {
		It("should return a fixed fee for fixFeeEstimator", func() {
			fee := rand.Intn(100)
			estimator := btc.NewFixFeeEstimator(fee)
			fees, err := estimator.FeeSuggestion()
			Expect(err).Should(BeNil())
			Expect(fees.Minimum).Should(Equal(fee))
			Expect(fees.Economy).Should(Equal(fee))
			Expect(fees.Low).Should(Equal(fee))
			Expect(fees.Medium).Should(Equal(fee))
			Expect(fees.High).Should(Equal(fee))
		})
	})
})
