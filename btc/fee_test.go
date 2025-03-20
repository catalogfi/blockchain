package btc_test

import (
	"context"
	"math/rand"
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
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
				estimator := btc.NewMempoolFeeEstimator(&chaincfg.MainNetParams, btc.MempoolFeeApiMainnet, 15*time.Second)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())

				Expect(fees.Low).Should(BeNumerically(">=", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", fees.Low))
				Expect(fees.High).Should(BeNumerically(">=", 1))
				Expect(fees.High).Should(BeNumerically(">=", fees.Medium))
			})

			It("should return the live data from mempool for testnet3", func() {
				estimator := btc.NewMempoolFeeEstimator(&chaincfg.TestNet3Params, btc.MempoolFeeApiTestnet3, 15*time.Second)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())

				Expect(fees.Low).Should(BeNumerically(">=", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", fees.Low))
				Expect(fees.High).Should(BeNumerically(">=", 1))
				Expect(fees.High).Should(BeNumerically(">=", fees.Medium))
			})

			It("should return the live data from mempool for testnet4", func() {
				estimator := btc.NewMempoolFeeEstimator(&chaincfg.TestNet3Params, btc.MempoolFeeApiTestnet4, 15*time.Second)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())

				Expect(fees.Low).Should(BeNumerically(">=", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", fees.Low))
				Expect(fees.High).Should(BeNumerically(">=", 1))
				Expect(fees.High).Should(BeNumerically(">=", fees.Medium))
			})

			It("should use the cached result if the more requests are made during certain time period", func() {
				estimator := btc.NewMempoolFeeEstimator(&chaincfg.MainNetParams, btc.MempoolFeeApiMainnet, time.Second)
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
				estimator := btc.NewBlockstreamFeeEstimator(&chaincfg.MainNetParams, btc.BlockstreamApiMainnet, 15*time.Second)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())

				Expect(fees.Low).Should(BeNumerically(">", 1))
				Expect(fees.Medium).Should(BeNumerically(">", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", fees.Low))
				Expect(fees.High).Should(BeNumerically(">", 1))
				Expect(fees.High).Should(BeNumerically(">=", fees.Medium))
			})

			It("should use the cached result if the more requests are made during certain time period", func() {
				estimator := btc.NewBlockstreamFeeEstimator(&chaincfg.MainNetParams, btc.BlockstreamApiMainnet, time.Second)
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
				Expect(fees.Low).Should(Equal(btc.SatoshiPerKb(1e3)))
				Expect(fees.Medium).Should(Equal(btc.SatoshiPerKb(1e3)))
				Expect(fees.High).Should(Equal(btc.SatoshiPerKb(1e3)))

				By("Regnet")
				estimatorRegnet := btc.NewBlockstreamFeeEstimator(&chaincfg.RegressionNetParams, "", 15*time.Second)
				fees, err = estimatorRegnet.FeeSuggestion()
				Expect(err).Should(BeNil())
				Expect(fees.Low).Should(Equal(btc.SatoshiPerKb(1e3)))
				Expect(fees.Medium).Should(Equal(btc.SatoshiPerKb(1e3)))
				Expect(fees.High).Should(Equal(btc.SatoshiPerKb(1e3)))
			})
		})

		Context("fix fee estimator", func() {
			It("should return a fixed fee for fixFeeEstimator", func() {
				feeRate := btctest.RandomFeeRate()
				estimator := btc.NewFixFeeEstimator(feeRate)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())
				Expect(fees.Low).Should(Equal(feeRate))
				Expect(fees.Medium).Should(Equal(feeRate))
				Expect(fees.High).Should(Equal(feeRate))
			})
		})
	})

	Context("estimate transaction fees", func() {
		It("should return a proper estimate of tx size when spending utxos from known address types", func(ctx context.Context) {
			for _, addrType := range addrTypes {
				By("Initialization keys")
				key1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
				Expect(err).To(BeNil())
				_, addr2, err := btctest.NewBtcKey(network, addrType)
				Expect(err).To(BeNil())

				By("Build some transactions and compare the estimated weight vs actual weight")
				utxos, err := indexer.GetUTXOs(ctx, addr1)
				Expect(err).To(BeNil())
				sizer := btc.NewSizeEstimatorOfAddrType(addrType, utxos...)

				edge := 0
				exact := 0
				for i := 0; i < 10000; i++ {
					amount, feeRate := rand.Int63n(10000)+1e7, btctest.RandomFeeRate()
					recipients := []btc.Recipient{btc.NewRecipient(addr2.EncodeAddress(), amount)}
					feeMode := btc.MinFeeRateMode(feeRate, sizer)
					transaction, err := btc.BuildTx(network, feeMode, nil, utxos, recipients, addr1)
					Expect(err).To(BeNil())

					estWeight, err := sizer.EstimateTxWeight(transaction)
					Expect(err).To(BeNil())
					Expect(btc.SignTx(addrType, transaction, key1, utxos)).Should(Succeed())
					actualWeight := blockchain.GetTransactionWeight(btcutil.NewTx(transaction))
					Expect(estWeight - int(actualWeight)).Should(BeNumerically(">=", 0))

					switch addrType {
					case waddrmgr.PubKeyHash:
						// Most tx should be the same or 1 byte less than the estimation.
						// A very small chance been 2 or 3 bytes less.
						switch estWeight - int(actualWeight) {
						case 0:
							exact++
						case 4:
						case 8, 12:
							edge++
						default:
							Fail("unexpected weight difference")
						}
					case waddrmgr.WitnessPubKey:
						// Most tx should be the same or 1 weight less than the estimation.
						// A very small chance been 2 or 3 weight less.
						switch estWeight - int(actualWeight) {
						case 0:
							exact++
						case 1:
						case 2, 3:
							edge++
						default:
							Fail("unexpected weight difference")
						}
					case waddrmgr.TaprootPubKey:
						// Should be exactly same or 1 weight more
						switch estWeight - int(actualWeight) {
						case 0:
							exact++
						case 1:
						default:
							Fail("unexpected weight difference")
						}
					}
				}

				By("Expect the edge cases to be less than 1% ")
				Expect(edge).Should(BeNumerically("<=", 100))
				Expect(exact).Should(BeNumerically(">", 0))
			}
		})

		It("should return a proper estimate of tx size when spending both legacy and segwit utxos", func(ctx context.Context) {
			By("Initialization keys")
			key1, addr1, err := btctest.NewBtcAddrWithFunds(network, waddrmgr.PubKeyHash, nil)
			Expect(err).To(BeNil())
			key2, addr2, err := btctest.NewBtcAddrWithFunds(network, waddrmgr.WitnessPubKey, indexer)
			Expect(err).To(BeNil())

			By("Fetch all utxos from addr1 and addr2")
			utxos1, err := indexer.GetUTXOs(ctx, addr1)
			Expect(err).To(BeNil())
			utxos2, err := indexer.GetUTXOs(ctx, addr2)
			Expect(err).To(BeNil())
			sizer := btc.NewSizeEstimatorOfAddrType(waddrmgr.PubKeyHash, utxos1...)
			sizer.AddUtxos(btc.BaseSizeP2WPKH, btc.SegwitSizeP2WPKH, utxos2...)

			By("Init the fetcher")
			pkScript1, err := btc.PkScript(waddrmgr.PubKeyHash, key1.PubKey())
			Expect(err).To(BeNil())
			fetcher, err := btc.NewFetcher(pkScript1, utxos1...)
			Expect(err).To(BeNil())
			pkScript2, err := btc.PkScript(waddrmgr.WitnessPubKey, key2.PubKey())
			Expect(err).To(BeNil())
			Expect(btc.AddUtxosToFetcher(fetcher, pkScript2, utxos2...)).Should(Succeed())

			By("Build and sign the transaction")
			diff := map[int]int{}
			for i := 0; i < 10000; i++ {
				feeMode := btc.MinFeeRateMode(btctest.RandomFeeRate(), sizer)
				transaction, err := btc.BuildTx(network, feeMode, append(utxos1, utxos2...), nil, nil, addr1)
				Expect(err).To(BeNil())
				estWeight, err := sizer.EstimateTxWeight(transaction)
				Expect(err).To(BeNil())

				sigHashes := txscript.NewTxSigHashes(transaction, fetcher)
				for i := range transaction.TxIn {
					if i < len(utxos1) {
						Expect(btc.SignInput(waddrmgr.PubKeyHash, transaction, i, key1, fetcher, sigHashes)).Should(Succeed())
					} else {
						Expect(btc.SignInput(waddrmgr.WitnessPubKey, transaction, i, key2, fetcher, sigHashes)).Should(Succeed())
					}
				}
				actualWeight := blockchain.GetTransactionWeight(btcutil.NewTx(transaction))
				Expect(estWeight - int(actualWeight)).Should(BeNumerically(">=", 0))
				diff[estWeight-int(actualWeight)]++
			}

			By("We should have some tx match our estimate, and most of the tx should be within 5 wu of our estimate")
			Expect(diff[0]).Should(BeNumerically(">", 0))
			edge := 0
			for k, v := range diff {
				if k > 5 {
					edge += v
				}
			}
			Expect(edge).Should(BeNumerically("<", 200))
		})
	})
})
