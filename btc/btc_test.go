package btc_test

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/fatih/color"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bitcoin", func() {
	Context("Build transaction", func() {
		Context("MinFeeRateMode", func() {
			It("should build a transaction with min fee rate", func(ctx context.Context) {
				for _, addrType := range addrTypes {
					By("Initialize keys and addresses")
					key1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
					Expect(err).To(BeNil())
					_, addr2, err := btctest.NewBtcKey(network, addrType)
					Expect(err).To(BeNil())

					By("Construct a transaction which sends money from addr1 to addr2")
					utxos, err := indexer.GetUTXOs(ctx, addr1)
					Expect(err).To(BeNil())
					amount, feeRate := int64(1e7), btctest.RandomFeeRate()
					recipient, err := btc.NewTxOutFromAddress(addr2, amount)
					Expect(err).To(BeNil())
					recipients := []*wire.TxOut{recipient}

					sizer := btc.NewSizeEstimatorOfAddrType(addrType, utxos...)
					feeMode := btc.MinFeeRateMode(feeRate, sizer)
					transaction, err := btc.BuildTx(feeMode, nil, utxos, recipients, addr1)
					Expect(err).To(BeNil())

					By("Sign and submit the fund tx")
					Expect(btc.QuickSign(addrType, transaction, key1, utxos)).Should(Succeed())
					Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
					By(color.GreenString("tx hash = %v", transaction.TxHash().String()))

					By("The actual fee rate should not be less than the given fee rate")
					query, err := indexer.GetTx(ctx, transaction.TxHash().String())
					Expect(err).To(BeNil())
					vsize := (query.Weight + 3) / 4
					actualFeeRate := btc.NewSatoshiPerKb(query.Fee, vsize)
					Expect(actualFeeRate).Should(BeNumerically(">=", feeRate))
					Expect(actualFeeRate - feeRate).Should(BeNumerically("<=", 1000))
					By(color.GreenString("Expected fee rate = %v, actual fee rate = %v, diff = %v", feeRate, actualFeeRate, feeRate-actualFeeRate))
				}
			})
		})

		Context("FixedFeesMode", func() {
			It("should build a transaction with fixed fee", func(ctx context.Context) {
				for _, addrType := range addrTypes {
					By("Initialize keys and addresses")
					key1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
					Expect(err).To(BeNil())
					_, addr2, err := btctest.NewBtcKey(network, addrType)
					Expect(err).To(BeNil())

					By("Construct a transaction which sends money from addr1 to addr2")
					utxos, err := indexer.GetUTXOs(ctx, addr1)
					Expect(err).To(BeNil())
					amount, fees := int64(1e7), int64(500+rand.Intn(1000))
					recipient, err := btc.NewTxOutFromAddress(addr2, amount)
					Expect(err).To(BeNil())
					recipients := []*wire.TxOut{recipient}
					feeMode := btc.FixedFeesMode(fees)
					transaction, err := btc.BuildTx(feeMode, nil, utxos, recipients, addr1)
					Expect(err).To(BeNil())

					By("Sign and submit the fund tx")
					Expect(btc.QuickSign(addrType, transaction, key1, utxos)).Should(Succeed())
					Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
					By(color.GreenString("tx hash = %v", transaction.TxHash().String()))

					By("The actual fee rate should not be less than the given fee rate")
					query, err := indexer.GetTx(ctx, transaction.TxHash().String())
					Expect(err).To(BeNil())
					Expect(query.Fee).Should(Equal(fees))
				}
			})
		})

		Context("RbfMode", func() {
			Context("No descendants", func() {
				It("should rbf a tx using minimum extra fees when not providing a required fee rate", func(ctx context.Context) {
					feeRates := []btc.SatoshiPerKb{1e3, 100e3}
					for _, addrType := range addrTypes {
						for _, feeRate := range feeRates {
							By("Initialize keys and addresses")
							key1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, nil)
							Expect(err).To(BeNil())
							fundingTx, err := btctest.Faucet(addr1.EncodeAddress())
							Expect(err).To(BeNil())
							Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(fundingTx.String()))).Should(Succeed())
							_, addr2, err := btctest.NewBtcKey(network, addrType)
							Expect(err).To(BeNil())

							By("Construct a transaction which sends money from addr1 to addr2")
							utxos, err := indexer.GetUTXOs(ctx, addr1)
							Expect(err).To(BeNil())
							amount := int64(1e7)
							recipient, err := btc.NewTxOutFromAddress(addr2, amount)
							Expect(err).To(BeNil())
							recipients := []*wire.TxOut{recipient}
							sizer := btc.NewSizeEstimatorOfAddrType(addrType, utxos...)
							feeMode := btc.MinFeeRateMode(feeRate, sizer)
							tx1, err := btc.BuildTx(feeMode, nil, utxos, recipients, addr1)
							Expect(err).To(BeNil())

							By("Sign and submit the fund tx")
							Expect(btc.QuickSign(addrType, tx1, key1, utxos)).Should(Succeed())
							Expect(indexer.SubmitTx(ctx, tx1)).Should(Succeed())
							By(color.GreenString("tx hash = %v", tx1.TxHash().String()))

							By("Build a new tx which replacing tx1")
							query, err := indexer.GetTx(ctx, tx1.TxHash().String())
							Expect(err).To(BeNil())
							vsize := (query.Weight + 3) / 4
							prevFeeRate := btc.NewSatoshiPerKb(query.Fee, vsize)
							feeMode1 := btc.RbfMode(0, prevFeeRate, query.Fee, sizer)
							tx2, err := btc.BuildTx(feeMode1, utxos, nil, recipients, addr1)
							Expect(err).To(BeNil())

							if addrType == waddrmgr.TaprootPubKey {
								By("Tx should be rejected if we pay one sat less in fee")
								tx2.TxOut[len(tx2.TxOut)-1].Value = tx2.TxOut[len(tx2.TxOut)-1].Value + 1
								Expect(btc.QuickSign(addrType, tx2, key1, utxos)).Should(Succeed())
								Expect(indexer.SubmitTx(ctx, tx2)).ShouldNot(Succeed())
								tx2.TxOut[len(tx2.TxOut)-1].Value = tx2.TxOut[len(tx2.TxOut)-1].Value - 1
							}

							By("Replacement tx should be accepted")
							Expect(btc.QuickSign(addrType, tx2, key1, utxos)).Should(Succeed())
							Expect(indexer.SubmitTx(ctx, tx2)).Should(Succeed())
							By(color.GreenString("tx hash = %v", tx2.TxHash().String()))
						}
					}
				})

				It("should consider the fee rate when providing a non-zero value", func(ctx context.Context) {
					for _, addrType := range addrTypes {
						By("Initialize keys and addresses")
						key1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, nil)
						Expect(err).To(BeNil())
						fundingTx, err := btctest.Faucet(addr1.EncodeAddress())
						Expect(err).To(BeNil())
						Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(fundingTx.String()))).Should(Succeed())
						_, addr2, err := btctest.NewBtcKey(network, addrType)
						Expect(err).To(BeNil())

						By("Construct a transaction which sends money from addr1 to addr2")
						utxos, err := indexer.GetUTXOs(ctx, addr1)
						Expect(err).To(BeNil())
						amount, feeRate := int64(1e7), btctest.RandomFeeRate()
						recipient, err := btc.NewTxOutFromAddress(addr2, amount)
						Expect(err).To(BeNil())
						recipients := []*wire.TxOut{recipient}
						sizer := btc.NewSizeEstimatorOfAddrType(addrType, utxos...)
						feeMode := btc.MinFeeRateMode(feeRate, sizer)
						tx1, err := btc.BuildTx(feeMode, nil, utxos, recipients, addr1)
						Expect(err).To(BeNil())

						By("Sign and submit the fund tx")
						Expect(btc.QuickSign(addrType, tx1, key1, utxos)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx1)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx1.TxHash().String()))

						By("Build a new tx which replacing tx1 and require a minimum fee rate")
						query, err := indexer.GetTx(ctx, tx1.TxHash().String())
						Expect(err).To(BeNil())
						vsize := (query.Weight + 3) / 4
						prevFeeRate := btc.NewSatoshiPerKb(query.Fee, vsize)
						minFeerate := prevFeeRate + 20e3
						feeMode1 := btc.RbfMode(minFeerate, prevFeeRate, query.Fee, sizer)
						tx2, err := btc.BuildTx(feeMode1, utxos, nil, recipients, addr1)
						Expect(err).To(BeNil())

						By("Replacement tx should be accepted")
						Expect(btc.QuickSign(addrType, tx2, key1, utxos)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx2)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx2.TxHash().String()))

						By("The actual fee rate should not be less than the given fee rate")
						query2, err := indexer.GetTx(ctx, tx2.TxHash().String())
						Expect(err).To(BeNil())
						vsize2 := (query2.Weight + 3) / 4
						actualFeeRate := btc.NewSatoshiPerKb(query2.Fee, vsize2)
						Expect(actualFeeRate).Should(BeNumerically(">=", minFeerate))
						Expect(actualFeeRate - minFeerate).Should(BeNumerically("<=", 1000))
						By(color.GreenString("Expected fee rate = %v, actual fee rate = %v, diff = %v", minFeerate, actualFeeRate, actualFeeRate-minFeerate))
					}
				})
			})

			Context("With descendants", func() {
				It("should rbf the tx when the fee is restricted by the fee bandwidth", func(ctx context.Context) {
					By("Initialize keys and addresses")
					for _, addrType := range addrTypes {
						key1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
						Expect(err).To(BeNil())
						key2, addr2, err := btctest.NewBtcKey(network, addrType)
						Expect(err).To(BeNil())

						By("Construct a transaction which sends money from addr1 to addr2")
						utxos1, err := indexer.GetUTXOs(ctx, addr1)
						Expect(err).To(BeNil())
						amount, feeRate := int64(1e7), btc.SatoshiPerKb(100e3)
						recipient, err := btc.NewTxOutFromAddress(addr2, amount)
						Expect(err).To(BeNil())
						recipients := []*wire.TxOut{recipient}
						sizer1 := btc.NewSizeEstimatorOfAddrType(addrType, utxos1...)
						feeMode1 := btc.MinFeeRateMode(feeRate, sizer1)
						tx1, err := btc.BuildTx(feeMode1, nil, utxos1, recipients, addr1)
						Expect(err).To(BeNil())

						By("Sign and submit the fund tx1")
						Expect(btc.QuickSign(addrType, tx1, key1, utxos1)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx1)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx1.TxHash().String()))

						By("Create a descendant tx of tx1 with a lower fee rate")
						time.Sleep(5 * time.Second)
						utxos2, err := indexer.GetUTXOs(ctx, addr2)
						Expect(err).To(BeNil())
						sizer2 := btc.NewSizeEstimatorOfAddrType(addrType, utxos2...)
						feeMode2 := btc.MinFeeRateMode(btc.SatoshiPerKb(1e3), sizer2)
						tx2, err := btc.BuildTx(feeMode2, utxos2, nil, nil, addr2)
						Expect(err).To(BeNil())

						By("Sign and submit the fund tx2")
						Expect(btc.QuickSign(addrType, tx2, key2, utxos2)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx2)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx2.TxHash().String()))

						By("Build a new tx which replacing tx1")
						entry, err := client.GetMempoolEntry(ctx, tx1.TxHash().String())
						Expect(err).To(BeNil())
						prevFees, err := btcutil.NewAmount(entry.Fees.Descendant)
						Expect(err).To(BeNil())
						query, err := indexer.GetTx(ctx, tx1.TxHash().String())
						Expect(err).To(BeNil())
						vsize := (query.Weight + 3) / 4
						prevFeeRate := btc.NewSatoshiPerKb(query.Fee, vsize)
						feeMode3 := btc.RbfMode(0, prevFeeRate, int64(prevFees), sizer1)
						tx3, err := btc.BuildTx(feeMode3, nil, utxos1, recipients, addr1)
						Expect(err).To(BeNil())

						if addrType == waddrmgr.TaprootPubKey {
							By("Tx should be rejected if we pay one sat less in fee")
							tx3.TxOut[len(tx3.TxOut)-1].Value = tx3.TxOut[len(tx3.TxOut)-1].Value + 1
							Expect(btc.QuickSign(addrType, tx3, key1, utxos1)).Should(Succeed())
							Expect(indexer.SubmitTx(ctx, tx3)).ShouldNot(Succeed())
							tx3.TxOut[len(tx3.TxOut)-1].Value = tx3.TxOut[len(tx3.TxOut)-1].Value - 1
						}

						By("Replacement tx should be accepted")
						Expect(btc.QuickSign(addrType, tx3, key1, utxos1)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx3)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx3.TxHash().String()))
					}
				})

				It("should rbf the tx when the fee is restricted by the fee rate", func(ctx context.Context) {
					By("Initialize keys and addresses")
					for _, addrType := range addrTypes {
						key1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
						Expect(err).To(BeNil())
						key2, addr2, err := btctest.NewBtcKey(network, addrType)
						Expect(err).To(BeNil())

						By("Construct a transaction which sends money from addr1 to addr2")
						utxos1, err := indexer.GetUTXOs(ctx, addr1)
						Expect(err).To(BeNil())
						amount, feeRate := int64(1e7), btc.SatoshiPerKb(100e3)
						recipient, err := btc.NewTxOutFromAddress(addr2, amount)
						Expect(err).To(BeNil())
						recipients := []*wire.TxOut{recipient}
						sizer1 := btc.NewSizeEstimatorOfAddrType(addrType, utxos1...)
						feeMode1 := btc.MinFeeRateMode(feeRate, sizer1)
						tx1, err := btc.BuildTx(feeMode1, nil, utxos1, recipients, addr1)
						Expect(err).To(BeNil())

						By("Sign and submit the fund tx1")
						Expect(btc.QuickSign(addrType, tx1, key1, utxos1)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx1)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx1.TxHash().String()))

						By("Create a descendant tx of tx1 with a lower fee rate")
						time.Sleep(5 * time.Second)
						utxos2, err := indexer.GetUTXOs(ctx, addr2)
						Expect(err).To(BeNil())
						sizer2 := btc.NewSizeEstimatorOfAddrType(addrType, utxos2...)
						feeMode2 := btc.MinFeeRateMode(btc.SatoshiPerKb(1e3), sizer2)
						tx2, err := btc.BuildTx(feeMode2, utxos2, nil, nil, addr2)
						Expect(err).To(BeNil())

						By("Sign and submit the fund tx2")
						Expect(btc.QuickSign(addrType, tx2, key2, utxos2)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx2)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx2.TxHash().String()))

						By("Build a new tx which replacing tx1")
						entry, err := client.GetMempoolEntry(ctx, tx1.TxHash().String())
						Expect(err).To(BeNil())
						prevFees, err := btcutil.NewAmount(entry.Fees.Descendant)
						Expect(err).To(BeNil())
						query, err := indexer.GetTx(ctx, tx1.TxHash().String())
						Expect(err).To(BeNil())
						vsize := (query.Weight + 3) / 4
						prevFeeRate := btc.NewSatoshiPerKb(query.Fee, vsize)
						feeMode3 := btc.RbfMode(0, prevFeeRate, int64(prevFees), sizer1)

						recipients1 := make([]*wire.TxOut, 20)
						for i := 0; i < 20; i++ {
							recipients1[i], err = btc.NewTxOutFromAddress(addr2, amount/10)
							Expect(err).To(BeNil())
						}
						tx3, err := btc.BuildTx(feeMode3, nil, utxos1, recipients1, addr1)
						Expect(err).To(BeNil())

						if addrType == waddrmgr.TaprootPubKey {
							By("Tx should be rejected if we pay one sat less in fee")
							tx3.TxOut[len(tx3.TxOut)-1].Value = tx3.TxOut[len(tx3.TxOut)-1].Value + 1
							Expect(btc.QuickSign(addrType, tx3, key1, utxos1)).Should(Succeed())
							Expect(indexer.SubmitTx(ctx, tx3)).ShouldNot(Succeed())
							tx3.TxOut[len(tx3.TxOut)-1].Value = tx3.TxOut[len(tx3.TxOut)-1].Value - 1
						}

						By("Replacement tx should be accepted")
						Expect(btc.QuickSign(addrType, tx3, key1, utxos1)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx3)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx3.TxHash().String()))
					}
				})

				It("should rbf the tx when the fee is restricted by the fee bandwidth", func(ctx context.Context) {
					By("Initialize keys and addresses")
					for _, addrType := range addrTypes {
						key1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
						Expect(err).To(BeNil())
						key2, addr2, err := btctest.NewBtcKey(network, addrType)
						Expect(err).To(BeNil())

						By("Construct a transaction which sends money from addr1 to addr2")
						utxos1, err := indexer.GetUTXOs(ctx, addr1)
						Expect(err).To(BeNil())
						amount, feeRate := int64(1e7), btc.SatoshiPerKb(1e3)
						recipient, err := btc.NewTxOutFromAddress(addr2, amount)
						Expect(err).To(BeNil())
						recipients := []*wire.TxOut{recipient}
						sizer1 := btc.NewSizeEstimatorOfAddrType(addrType, utxos1...)
						feeMode1 := btc.MinFeeRateMode(feeRate, sizer1)
						tx1, err := btc.BuildTx(feeMode1, nil, utxos1, recipients, addr1)
						Expect(err).To(BeNil())

						By("Sign and submit the fund tx1")
						Expect(btc.QuickSign(addrType, tx1, key1, utxos1)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx1)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx1.TxHash().String()))

						By("Create a descendant tx of tx1 with a lower fee rate")
						time.Sleep(5 * time.Second)
						utxos2, err := indexer.GetUTXOs(ctx, addr2)
						Expect(err).To(BeNil())
						sizer2 := btc.NewSizeEstimatorOfAddrType(addrType, utxos2...)
						feeMode2 := btc.MinFeeRateMode(btc.SatoshiPerKb(100e3), sizer2)
						tx2, err := btc.BuildTx(feeMode2, utxos2, nil, nil, addr2)
						Expect(err).To(BeNil())

						By("Sign and submit the fund tx2")
						Expect(btc.QuickSign(addrType, tx2, key2, utxos2)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx2)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx2.TxHash().String()))

						By("Build a new tx which replacing tx1")
						entry, err := client.GetMempoolEntry(ctx, tx1.TxHash().String())
						Expect(err).To(BeNil())
						prevFees, err := btcutil.NewAmount(entry.Fees.Descendant)
						Expect(err).To(BeNil())
						query, err := indexer.GetTx(ctx, tx1.TxHash().String())
						Expect(err).To(BeNil())
						vsize := (query.Weight + 3) / 4
						prevFeeRate := btc.NewSatoshiPerKb(query.Fee, vsize)
						feeMode3 := btc.RbfMode(0, prevFeeRate, int64(prevFees), sizer1)
						tx3, err := btc.BuildTx(feeMode3, nil, utxos1, recipients, addr1)
						Expect(err).To(BeNil())

						if addrType == waddrmgr.TaprootPubKey {
							By("Tx should be rejected if we pay one sat less in fee")
							tx3.TxOut[len(tx3.TxOut)-1].Value = tx3.TxOut[len(tx3.TxOut)-1].Value + 1
							Expect(btc.QuickSign(addrType, tx3, key1, utxos1)).Should(Succeed())
							Expect(indexer.SubmitTx(ctx, tx3)).ShouldNot(Succeed())
							tx3.TxOut[len(tx3.TxOut)-1].Value = tx3.TxOut[len(tx3.TxOut)-1].Value - 1
						}

						By("Replacement tx should be accepted")
						Expect(btc.QuickSign(addrType, tx3, key1, utxos1)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx3)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx3.TxHash().String()))
					}
				})

				It("should rbf the tx when the fee is restricted by the fee rate", func(ctx context.Context) {
					By("Initialize keys and addresses")
					for _, addrType := range addrTypes {
						key1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
						Expect(err).To(BeNil())
						key2, addr2, err := btctest.NewBtcKey(network, addrType)
						Expect(err).To(BeNil())

						By("Construct a transaction which sends money from addr1 to addr2")
						utxos1, err := indexer.GetUTXOs(ctx, addr1)
						Expect(err).To(BeNil())
						amount, feeRate := int64(1e7), btc.SatoshiPerKb(50e3)
						recipient, err := btc.NewTxOutFromAddress(addr2, amount)
						Expect(err).To(BeNil())
						recipients := []*wire.TxOut{recipient}
						sizer1 := btc.NewSizeEstimatorOfAddrType(addrType, utxos1...)
						feeMode1 := btc.MinFeeRateMode(feeRate, sizer1)
						tx1, err := btc.BuildTx(feeMode1, nil, utxos1, recipients, addr1)
						Expect(err).To(BeNil())

						By("Sign and submit the fund tx1")
						Expect(btc.QuickSign(addrType, tx1, key1, utxos1)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx1)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx1.TxHash().String()))

						By("Create a descendant tx of tx1 with a lower fee rate")
						time.Sleep(5 * time.Second)
						utxos2, err := indexer.GetUTXOs(ctx, addr2)
						Expect(err).To(BeNil())
						sizer2 := btc.NewSizeEstimatorOfAddrType(addrType, utxos2...)
						feeMode2 := btc.MinFeeRateMode(btc.SatoshiPerKb(55e3), sizer2)
						tx2, err := btc.BuildTx(feeMode2, utxos2, nil, nil, addr2)
						Expect(err).To(BeNil())

						By("Sign and submit the fund tx2")
						Expect(btc.QuickSign(addrType, tx2, key2, utxos2)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx2)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx2.TxHash().String()))

						By("Build a new tx which replacing tx1")
						entry, err := client.GetMempoolEntry(ctx, tx1.TxHash().String())
						Expect(err).To(BeNil())
						prevFees, err := btcutil.NewAmount(entry.Fees.Descendant)
						Expect(err).To(BeNil())
						query, err := indexer.GetTx(ctx, tx1.TxHash().String())
						Expect(err).To(BeNil())
						vsize := (query.Weight + 3) / 4
						prevFeeRate := btc.NewSatoshiPerKb(query.Fee, vsize)
						recipients1 := make([]*wire.TxOut, 20)
						for i := 0; i < 20; i++ {
							recipients1[i], err = btc.NewTxOutFromAddress(addr2, amount/10)
							Expect(err).To(BeNil())
						}
						feeMode3 := btc.RbfMode(0, prevFeeRate, int64(prevFees), sizer1)
						tx3, err := btc.BuildTx(feeMode3, nil, utxos1, recipients1, addr1)
						Expect(err).To(BeNil())

						if addrType == waddrmgr.TaprootPubKey {
							By("Tx should be rejected if we pay one sat less in fee")
							tx3.TxOut[len(tx3.TxOut)-1].Value = tx3.TxOut[len(tx3.TxOut)-1].Value + 1
							Expect(btc.QuickSign(addrType, tx3, key1, utxos1)).Should(Succeed())
							Expect(indexer.SubmitTx(ctx, tx3)).ShouldNot(Succeed())
							tx3.TxOut[len(tx3.TxOut)-1].Value = tx3.TxOut[len(tx3.TxOut)-1].Value - 1
						}

						By("Replacement tx should be accepted")
						Expect(btc.QuickSign(addrType, tx3, key1, utxos1)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx3)).Should(Succeed())
						By(color.GreenString("tx hash = %v", tx3.TxHash().String()))
					}
				})
			})
		})

		Context("OP_RETURN", func() {
			It("should be able to build an OP_RETURN tx", func(ctx context.Context) {
				for _, addrType := range addrTypes {
					By("Initialize keys and addresses")
					key1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
					Expect(err).To(BeNil())

					By("Construct a transaction which sends money from addr1 to addr2")
					utxos, err := indexer.GetUTXOs(ctx, addr1)
					Expect(err).To(BeNil())
					feeRate := btctest.RandomFeeRate()
					sizer := btc.NewSizeEstimatorOfAddrType(addrType, utxos...)
					feeMode := btc.MinFeeRateMode(feeRate, sizer)
					transaction, err := btc.BuildTx(feeMode, nil, utxos, nil, addr1)
					Expect(err).To(BeNil())

					By("Construct a message in the OP_RETURN script")
					message := "hello world!"
					script, err := txscript.NewScriptBuilder().
						AddOp(txscript.OP_RETURN).
						AddData([]byte(message)).
						Script()
					Expect(err).To(BeNil())

					transaction.TxOut = append(transaction.TxOut, &wire.TxOut{
						Value:    0,
						PkScript: script,
					})

					By("Sign and submit the fund tx")
					Expect(btc.QuickSign(addrType, transaction, key1, utxos)).Should(Succeed())
					Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
					By(color.GreenString("tx hash = %v", transaction.TxHash().String()))
				}
			})
		})

		Context("Utxo selection", func() {
			It("should not use a utxo if its uneconomical", func(ctx context.Context) {
				By("Initialize keys and addresses")
				addrType := waddrmgr.TaprootPubKey
				_, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, nil)
				Expect(err).To(BeNil())
				fundingTx, err := btctest.Faucet(addr1.EncodeAddress())
				Expect(err).To(BeNil())
				Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(fundingTx.String()))).Should(Succeed())

				By("Manually modify one of the utxo value and try use it to build a tx")
				utxos, err := indexer.GetUTXOs(ctx, addr1)
				Expect(err).To(BeNil())
				Expect(len(utxos)).Should(Equal(2))
				utxos[0].Amount = btc.DustAmount - 1

				By("Try using the utxo set to build a tx")
				amount, feeRate := int64(1e7), btctest.RandomFeeRate()
				recipient, err := btc.NewTxOutFromAddress(addr1, amount)
				Expect(err).To(BeNil())
				recipients := []*wire.TxOut{recipient}
				sizer := btc.NewSizeEstimatorOfAddrType(addrType, utxos...)
				feeMode := btc.MinFeeRateMode(feeRate, sizer)
				transaction, err := btc.BuildTx(feeMode, nil, utxos, recipients, addr1)
				Expect(err).To(BeNil())

				By("The built tx should not use the uneconomical utxo")
				for _, input := range transaction.TxIn {
					if input.PreviousOutPoint.String() == utxos[0].String() {
						Fail("The new tx should not use the utxo which is less than the dust amount")
					}
				}
			})
		})
	})

	Context("Tx serialization", func() {
		It("should serialize a tx into a readable json format", func(ctx context.Context) {
			By("Initialize keys and addresses")
			addrType := waddrmgr.TaprootPubKey
			key1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).To(BeNil())
			_, addr2, err := btctest.NewBtcKey(network, addrType)
			Expect(err).To(BeNil())

			By("Construct a transaction which sends money from addr1 to addr2")
			utxos, err := indexer.GetUTXOs(ctx, addr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e7), btctest.RandomFeeRate()
			recipient, err := btc.NewTxOutFromAddress(addr2, amount)
			Expect(err).To(BeNil())
			recipients := []*wire.TxOut{recipient}
			sizer := btc.NewSizeEstimatorOfAddrType(addrType, utxos...)
			feeMode := btc.MinFeeRateMode(feeRate, sizer)
			transaction, err := btc.BuildTx(feeMode, nil, utxos, recipients, addr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.QuickSign(addrType, transaction, key1, utxos)).Should(Succeed())

			By("Decode the tx")
			rawTxResult, err := btc.CreateTxRawResult(network, transaction)
			Expect(err).To(BeNil())
			raw, err := json.MarshalIndent(rawTxResult, "", "  ")
			Expect(err).To(BeNil())
			log.Println(string(raw))
		})
	})
})
