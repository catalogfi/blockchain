package btc_test

import (
	"context"
	"math/rand"
	"time"

	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/fatih/color"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bitcoin", func() {
	Context("Build transaction", func() {
		Context("Using a minimum fee rate mode", func() {
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
					amount, feeRate := int64(1e7), 1+rand.Intn(100)*1000
					recipients := []btc.Recipient{btc.NewRecipient(addr2.EncodeAddress(), amount)}
					sizer := btc.NewSizeEstimatorOfAddrType(addrType, utxos...)
					feeMode := btc.MinFeeRateMode(feeRate, sizer)
					transaction, err := btc.BuildTx(network, feeMode, nil, utxos, recipients, addr1)
					Expect(err).To(BeNil())

					By("Sign and submit the fund tx")
					Expect(btc.SignTx(addrType, transaction, key1, utxos)).Should(Succeed())
					Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
					By(color.GreenString("tx hash = %v", transaction.TxHash().String()))

					By("The actual fee rate should not be less than the given fee rate")
					query, err := indexer.GetTx(ctx, transaction.TxHash().String())
					Expect(err).To(BeNil())
					vsize := (query.Weight + 3) / 4
					actualFeeRate := int(query.Fee) * 1000 / vsize
					Expect(actualFeeRate).Should(BeNumerically(">=", feeRate))
					Expect(actualFeeRate - feeRate).Should(BeNumerically("<=", 1000))
					By(color.GreenString("Expected fee rate = %v, actual fee rate = %v", feeRate, actualFeeRate))
				}
			})
		})

		Context("Using a fixed fee mode", func() {
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
					recipients := []btc.Recipient{btc.NewRecipient(addr2.EncodeAddress(), amount)}
					feeMode := btc.FixedFeesMode(fees)
					transaction, err := btc.BuildTx(network, feeMode, nil, utxos, recipients, addr1)
					Expect(err).To(BeNil())

					By("Sign and submit the fund tx")
					Expect(btc.SignTx(addrType, transaction, key1, utxos)).Should(Succeed())
					Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
					By(color.GreenString("tx hash = %v", transaction.TxHash().String()))

					By("The actual fee rate should not be less than the given fee rate")
					query, err := indexer.GetTx(ctx, transaction.TxHash().String())
					Expect(err).To(BeNil())
					Expect(query.Fee).Should(Equal(fees))
				}
			})
		})

		Context("Using a rbf fee mode", func() {
			It("should rbf a tx using minimum extra fees", func(ctx context.Context) {
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
					amount, feeRate := int64(1e7), 1+rand.Intn(100)*1000
					recipients := []btc.Recipient{btc.NewRecipient(addr2.EncodeAddress(), amount)}
					sizer := btc.NewSizeEstimatorOfAddrType(addrType, utxos...)
					feeMode := btc.MinFeeRateMode(feeRate, sizer)
					tx1, err := btc.BuildTx(network, feeMode, nil, utxos, recipients, addr1)
					Expect(err).To(BeNil())

					By("Sign and submit the fund tx")
					Expect(btc.SignTx(addrType, tx1, key1, utxos)).Should(Succeed())
					Expect(indexer.SubmitTx(ctx, tx1)).Should(Succeed())
					By(color.GreenString("tx hash = %v", tx1.TxHash().String()))

					By("Build a new tx which replacing tx1")
					query, err := indexer.GetTx(ctx, tx1.TxHash().String())
					Expect(err).To(BeNil())
					vsize := (query.Weight + 3) / 4
					prevFeeRate := int(query.Fee) * 1000 / vsize
					feeMode1 := btc.RbfMode(0, prevFeeRate, query.Fee, sizer)
					tx2, err := btc.BuildTx(network, feeMode1, utxos, nil, recipients, addr1)
					Expect(err).To(BeNil())

					if addrType == waddrmgr.TaprootPubKey {
						By("Tx should be rejected if we pay one sat less in fee")
						tx2.TxOut[len(tx2.TxOut)-1].Value = tx2.TxOut[len(tx2.TxOut)-1].Value + 1
						Expect(btc.SignTx(addrType, tx2, key1, utxos)).Should(Succeed())
						Expect(indexer.SubmitTx(ctx, tx2)).ShouldNot(Succeed())
						tx2.TxOut[len(tx2.TxOut)-1].Value = tx2.TxOut[len(tx2.TxOut)-1].Value - 1
					}

					By("Replacement tx should be accepted")
					Expect(btc.SignTx(addrType, tx2, key1, utxos)).Should(Succeed())
					Expect(indexer.SubmitTx(ctx, tx2)).Should(Succeed())
					By(color.GreenString("tx hash = %v", tx2.TxHash().String()))
				}
			})

			It("should work when the replaced tx has descendants", func(ctx context.Context) {
				By("Initialize keys and addresses")
				addrType := waddrmgr.TaprootPubKey
				key1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
				Expect(err).To(BeNil())
				key2, addr2, err := btctest.NewBtcKey(network, addrType)
				Expect(err).To(BeNil())

				By("Construct a transaction which sends money from addr1 to addr2")
				utxos1, err := indexer.GetUTXOs(ctx, addr1)
				Expect(err).To(BeNil())
				amount, feeRate := int64(1e7), 1+rand.Intn(100)*1000
				recipients := []btc.Recipient{btc.NewRecipient(addr2.EncodeAddress(), amount)}
				sizer1 := btc.NewSizeEstimatorOfAddrType(addrType, utxos1...)
				feeMode1 := btc.MinFeeRateMode(feeRate, sizer1)
				tx1, err := btc.BuildTx(network, feeMode1, nil, utxos1, recipients, addr1)
				Expect(err).To(BeNil())

				By("Sign and submit the fund tx1")
				Expect(btc.SignTx(addrType, tx1, key1, utxos1)).Should(Succeed())
				Expect(indexer.SubmitTx(ctx, tx1)).Should(Succeed())
				By(color.GreenString("tx hash = %v", tx1.TxHash().String()))

				By("Create a descendant tx of tx1 with a higher fee rate")
				time.Sleep(5 * time.Second)
				utxos2, err := indexer.GetUTXOs(ctx, addr2)
				Expect(err).To(BeNil())
				sizer2 := btc.NewSizeEstimatorOfAddrType(addrType, utxos2...)
				feeMode2 := btc.MinFeeRateMode(feeRate+10000, sizer2)
				tx2, err := btc.BuildTx(network, feeMode2, utxos2, nil, nil, addr2)
				Expect(err).To(BeNil())

				By("Sign and submit the fund tx2")
				Expect(btc.SignTx(addrType, tx2, key2, utxos2)).Should(Succeed())
				Expect(indexer.SubmitTx(ctx, tx2)).Should(Succeed())
				By(color.GreenString("tx hash = %v", tx2.TxHash().String()))

				time.Sleep(30 * time.Second)
				By("Build a new tx which replacing tx1")
				entry, err := client.GetMempoolEntry(ctx, tx1.TxHash().String())
				Expect(err).To(BeNil())
				prevFeeRate := int(entry.Fees.Descendant*1e8*1e3) / int(entry.DescendantSize)
				feeMode3 := btc.RbfMode(0, prevFeeRate, int64(entry.Fees.Descendant*1e8), sizer1)
				tx3, err := btc.BuildTx(network, feeMode3, nil, utxos1, recipients, addr1)
				Expect(err).To(BeNil())

				By("Tx should be rejected if we pay one sat less in fee")
				tx3.TxOut[len(tx3.TxOut)-1].Value = tx3.TxOut[len(tx3.TxOut)-1].Value + 1
				Expect(btc.SignTx(addrType, tx3, key1, utxos1)).Should(Succeed())
				Expect(indexer.SubmitTx(ctx, tx3)).ShouldNot(Succeed())

				By("Replacement tx should be accepted by pumping 1 sat in fee")
				tx3.TxOut[len(tx3.TxOut)-1].Value = tx3.TxOut[len(tx3.TxOut)-1].Value - 1
				Expect(btc.SignTx(addrType, tx3, key1, utxos1)).Should(Succeed())
				Expect(indexer.SubmitTx(ctx, tx3)).Should(Succeed())
				By(color.GreenString("tx hash = %v", tx3.TxHash().String()))
			})
		})

		// It("should not use a utxo if its uneconomical", func(ctx context.Context) {
		// 	By("Initialize keys and addresses")
		// 	addrType := waddrmgr.TaprootPubKey
		// 	privKey1, p2trAddr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
		// 	Expect(err).To(BeNil())
		// 	_, p2trAddr2, err := btctest.NewBtcKey(network, addrType)
		// 	Expect(err).To(BeNil())
		//
		// 	By("Create a utxo which is less than the dust amount")
		// 	utxos, err := indexer.GetUTXOs(ctx, p2trAddr1)
		// 	Expect(err).To(BeNil())
		// 	recipients := btc.SingleRecipient(p2trAddr1.EncodeAddress(), btc.DustAmount)
		// 	sizer := btc.NewSizeEstimator(utxos, btc.BaseSizeP2TR, btc.SegwitSizeP2TR)
		// 	transaction, err := btc.BuildTx(network, 3, nil, utxos, sizer, recipients, p2trAddr1)
		// 	Expect(err).To(BeNil())
		// 	Expect(btc.SignTx(addrType, transaction, privKey1, utxos)).Should(Succeed())
		// 	Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
		// 	Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())
		//
		// 	By("Construct a transaction try using the dust amount")
		// 	utxos1, err := indexer.GetUTXOs(ctx, p2trAddr1)
		// 	Expect(err).To(BeNil())
		// 	sort.Slice(utxos1, func(i, j int) bool {
		// 		return utxos1[i].Amount < utxos1[j].Amount
		// 	})
		// 	Expect(utxos1[0].TxID).Should(Equal(transaction.TxHash().String()))
		// 	amount, feeRate := int64(1e7), 3
		// 	recipients1 := btc.SingleRecipient(p2trAddr2.EncodeAddress(), amount)
		// 	sizer1 := btc.NewSizeEstimator(utxos1, btc.BaseSizeP2TR, btc.SegwitSizeP2TR)
		// 	transaction1, err := btc.BuildTx(network, feeRate, nil, utxos1, sizer1, recipients1, p2trAddr1)
		// 	Expect(err).To(BeNil())
		//
		// 	By("The new tx should not use the uneconomical utxo")
		// 	for _, input := range transaction1.TxIn {
		// 		if input.PreviousOutPoint.String() == utxos1[0].String() {
		// 			Fail("The new tx should not use the utxo which is less than the dust amount")
		// 		}
		// 	}
		// })
	})
})
