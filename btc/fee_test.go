package btc_test

import (
	"context"
	"math/rand"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("bitcoin fees", func() {
	Context("fee estimator", func() {
		Context("mempool estimator", func() {
			It("should return the live data from mempool", func() {
				estimator := btc.NewMempoolFeeEstimator(&chaincfg.MainNetParams, btc.MempoolFeeAPI, 15*time.Second)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())

				Expect(fees.Low).Should(BeNumerically(">=", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", fees.Low))
				Expect(fees.High).Should(BeNumerically(">=", 1))
				Expect(fees.High).Should(BeNumerically(">=", fees.Medium))
			})

			It("should return the live data from mempool for testnet", func() {
				estimator := btc.NewMempoolFeeEstimator(&chaincfg.TestNet3Params, btc.MempoolFeeAPITestnet, 15*time.Second)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())

				Expect(fees.Low).Should(BeNumerically(">=", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", fees.Low))
				Expect(fees.High).Should(BeNumerically(">=", 1))
				Expect(fees.High).Should(BeNumerically(">=", fees.Medium))
			})

			It("should use the cached result if the more requests are made during certain time period", func() {
				estimator := btc.NewMempoolFeeEstimator(&chaincfg.MainNetParams, btc.MempoolFeeAPI, time.Second)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())

				time.Sleep(500 * time.Millisecond)

				fees1, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())
				Expect(fees.Low).Should(Equal(fees1.Low))
				Expect(fees.Medium).Should(Equal(fees1.Medium))
				Expect(fees.High).Should(Equal(fees1.High))

				time.Sleep(500 * time.Millisecond)

				_, err = estimator.FeeSuggestion()
				Expect(err).Should(BeNil())
			})

			It("should panic if the network is not testnet or mainnet", func() {
				By("Regnet")
				PanicWith(func() {
					_ = btc.NewMempoolFeeEstimator(&chaincfg.RegressionNetParams, "", 15*time.Second)
				})
			})
		})

		Context("blockstream estimator", func() {
			It("should return the live data from blockstream api", func() {
				estimator := btc.NewBlockstreamFeeEstimator(&chaincfg.MainNetParams, btc.BlockstreamAPI, 15*time.Second)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())

				Expect(fees.Low).Should(BeNumerically(">", 1))
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
				Expect(err).Should(BeNil())
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
				Expect(fees.Low).Should(Equal(1))
				Expect(fees.Medium).Should(Equal(1))
				Expect(fees.High).Should(Equal(1))

				By("Regnet")
				estimatorRegnet := btc.NewBlockstreamFeeEstimator(&chaincfg.RegressionNetParams, "", 15*time.Second)
				fees, err = estimatorRegnet.FeeSuggestion()
				Expect(err).Should(BeNil())
				Expect(fees.Low).Should(Equal(1))
				Expect(fees.Medium).Should(Equal(1))
				Expect(fees.High).Should(Equal(1))
			})
		})

		Context("fix fee estimator", func() {
			It("should return a fixed fee for fixFeeEstimator", func() {
				fee := rand.Intn(100)
				estimator := btc.NewFixFeeEstimator(fee)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())
				Expect(fees.Low).Should(Equal(fee))
				Expect(fees.Medium).Should(Equal(fee))
				Expect(fees.High).Should(Equal(fee))
			})
		})
	})

	Context("estimate transaction fees", func() {
		It("should return a proper estimate of tx size when spending p2pkh utxos", func(ctx context.Context) {
			By("Initialization keys")
			key1, addr1, err := btctest.NewBtcAddrWithFunds(network, waddrmgr.PubKeyHash, indexer)
			Expect(err).To(BeNil())
			_, addr2, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(ctx, addr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e5), 10
			recipients := btc.SingleRecipient(addr2.EncodeAddress(), amount)
			sizeEstimator := btc.NewSizeEstimator(utxos, btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH)
			transaction, err := btc.BuildTransaction(network, feeRate, nil, utxos, sizeEstimator, recipients, addr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.SignTx(waddrmgr.PubKeyHash, transaction, key1, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())

			By("The actual fee rate should be within [feeRate -1 , feeRate] range")
			tx, err := indexer.GetTx(ctx, transaction.TxHash().String())
			Expect(err).To(BeNil())
			actualFeeRate := float64(tx.Fee) * 4 / float64(tx.Weight)
			Expect(float64(feeRate) - actualFeeRate).Should(BeNumerically("<=", float64(1)))
		})

		It("should return a proper estimate of tx size when spending p2wpkh utxos", func(ctx context.Context) {
			By("Initialization keys")
			key1, addr1, err := btctest.NewBtcAddrWithFunds(network, waddrmgr.WitnessPubKey, indexer)
			Expect(err).To(BeNil())
			_, addr2, err := btctest.NewBtcKey(network, waddrmgr.WitnessPubKey)
			Expect(err).To(BeNil())

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(ctx, addr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e5), 10
			recipients := btc.SingleRecipient(addr2.EncodeAddress(), amount)
			sizeEstimator := btc.NewSizeEstimator(utxos, btc.BaseSizeP2WPKH, btc.SegwitSizeP2WPKH)
			transaction, err := btc.BuildTransaction(network, feeRate, nil, utxos, sizeEstimator, recipients, addr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.SignTx(waddrmgr.WitnessPubKey, transaction, key1, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())

			By("The actual fee rate should be within [feeRate -1 , feeRate] range")
			tx, err := indexer.GetTx(ctx, transaction.TxHash().String())
			Expect(err).To(BeNil())
			actualFeeRate := float64(tx.Fee) * 4 / float64(tx.Weight)
			Expect(float64(feeRate) - actualFeeRate).Should(BeNumerically("<=", float64(1)))
		})

		It("should return a proper estimate of tx size when spending p2tr utxos", func(ctx context.Context) {
			By("Initialization keys")
			key1, addr1, err := btctest.NewBtcAddrWithFunds(network, waddrmgr.TaprootPubKey, indexer)
			Expect(err).To(BeNil())
			_, addr2, err := btctest.NewBtcKey(network, waddrmgr.TaprootPubKey)
			Expect(err).To(BeNil())

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(ctx, addr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e5), 10
			recipients := btc.SingleRecipient(addr2.EncodeAddress(), amount)
			sizeEstimator := btc.NewSizeEstimator(utxos, btc.BaseSizeP2TR, btc.SegwitSizeP2TR)
			transaction, err := btc.BuildTransaction(network, feeRate, nil, utxos, sizeEstimator, recipients, addr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.SignTx(waddrmgr.TaprootPubKey, transaction, key1, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())

			By("The actual fee rate should be within [feeRate -1 , feeRate] range")
			tx, err := indexer.GetTx(ctx, transaction.TxHash().String())
			Expect(err).To(BeNil())
			actualFeeRate := float64(tx.Fee) * 4 / float64(tx.Weight)
			Expect(float64(feeRate) - actualFeeRate).Should(BeNumerically("<=", float64(1)))
		})
	})
})
