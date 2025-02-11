package btc_test

import (
	"context"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wallet", func() {
	Context("Operation on an HTLC", func() {
		It("should be able to initiate and redeem an HTLC", func(ctx context.Context) {
			By("Init keys and wallet")
			addrType := waddrmgr.WitnessPubKey
			key1, _, err := btctest.NewBtcAddrWithFunds(network, addrType, nil)
			Expect(err).Should(BeNil())
			key2, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).Should(BeNil())
			feeEstimator := btc.NewFixFeeEstimator(10)
			wal1, err := btc.NewWallet(network, addrType, key1, indexer, feeEstimator)
			Expect(err).Should(BeNil())
			wal2, err := btc.NewWallet(network, addrType, key2, indexer, feeEstimator)
			Expect(err).Should(BeNil())

			By("Initiate and redeem an HTLC")
			amount, timelock := int64(1e7), int64(144)
			htlc, secret, err := btctest.NewHtlc(key1.PubKey(), key2.PubKey(), timelock, amount)
			Expect(err).Should(BeNil())
			_, err = wal1.Initiate(ctx, htlc)
			Expect(err).Should(BeNil())
			Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())
			_, err = wal2.Redeem(ctx, htlc, secret)
			Expect(err).Should(BeNil())
		})

		It("should be able to refund an HTLC after it expires", func() {
			By("Init keys and wallet")
			addrType := waddrmgr.WitnessPubKey
			key1, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).Should(BeNil())
			key2, _, err := btctest.NewBtcKey(network, addrType)
			Expect(err).Should(BeNil())
			feeEstimator := btc.NewFixFeeEstimator(10)
			wal1, err := btc.NewWallet(network, addrType, key1, indexer, feeEstimator)
			Expect(err).Should(BeNil())

			By("Initiate an HTLC")
			amount, timelock := int64(1e7), int64(6)
			htlc, _, err := btctest.NewHtlc(key1.PubKey(), key2.PubKey(), timelock, amount)
			_, err = wal1.Initiate(context.Background(), htlc)
			Expect(err).Should(BeNil())

			By("Mine expiry number of blocks")
			for i := 0; i < int(timelock)-1; i++ {
				err = btctest.NewBlock()
				Expect(err).Should(BeNil())
			}
			Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

			By("Refund an HTLC")
			_, err = wal1.Refund(context.Background(), htlc)
			Expect(err).Should(BeNil())
		})

		It("should be able to instant refund a HTLC", func(ctx context.Context) {
			By("Init keys and wallet")
			addrType := waddrmgr.WitnessPubKey
			key1, _, err := btctest.NewBtcAddrWithFunds(network, addrType, nil)
			Expect(err).Should(BeNil())
			key2, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).Should(BeNil())
			feeEstimator := btc.NewFixFeeEstimator(10)
			wal1, err := btc.NewWallet(network, addrType, key1, indexer, feeEstimator)
			Expect(err).Should(BeNil())
			wal2, err := btc.NewWallet(network, addrType, key2, indexer, feeEstimator)
			Expect(err).Should(BeNil())

			By("Initiate an HTLC")
			amount, timelock := int64(1e7), int64(6)
			htlc, _, err := btctest.NewHtlc(key1.PubKey(), key2.PubKey(), timelock, amount)
			initTx, err := wal1.Initiate(context.Background(), htlc)
			Expect(err).Should(BeNil())

			By("Initiate the HTLC with the pre-signed instant refund tx")
			initUtxo := btc.UTXO{
				TxID:   initTx.TxHash().String(),
				Vout:   0,
				Amount: amount,
			}
			recipient := btc.Recipient{
				To:     wal2.Address().String(),
				Amount: amount,
			}
			refundTx, err := btc.NewInstantRefundTx(network, key1, htlc, initUtxo, recipient)
			Expect(err).Should(BeNil())

			By("Mine a new block")
			Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

			By("Use the instant refund leaf to refund the HTLC")
			_, err = wal2.InstantRefund(ctx, htlc, refundTx)
			Expect(err).Should(BeNil())
		})
	})

	Context("Batch operation", func() {
		Context("single action", func() {
			It("should be able to initiate and redeem an HTLC", func(ctx context.Context) {
				By("Init keys and wallet")
				addrType := waddrmgr.WitnessPubKey
				key1, _, err := btctest.NewBtcAddrWithFunds(network, addrType, nil)
				Expect(err).Should(BeNil())
				key2, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
				Expect(err).Should(BeNil())
				feeEstimator := btc.NewFixFeeEstimator(10)
				wal1, err := btc.NewWallet(network, addrType, key1, indexer, feeEstimator)
				Expect(err).Should(BeNil())
				wal2, err := btc.NewWallet(network, addrType, key2, indexer, feeEstimator)
				Expect(err).Should(BeNil())

				By("Initiate an HTLC")
				amount, timelock := int64(1e7), int64(144)
				htlc, secret, err := btctest.NewHtlc(key1.PubKey(), key2.PubKey(), timelock, amount)
				Expect(err).Should(BeNil())
				actions1 := []btc.HtlcAction{
					{
						ActionType: btc.HtlcActionInitiate,
						Htlc:       htlc,
					},
				}
				_, err = wal1.Execute(ctx, actions1)
				Expect(err).Should(BeNil())
				Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

				By("Redeem an HTLC")
				actions2 := []btc.HtlcAction{
					{
						ActionType: btc.HtlcActionRedeem,
						Htlc:       htlc,
						Secret:     secret,
					},
				}
				_, err = wal2.Execute(ctx, actions2)
				Expect(err).Should(BeNil())
			})

			It("should be able to refund an HTLC after it expires", func(ctx context.Context) {
				By("Init keys and wallet")
				addrType := waddrmgr.WitnessPubKey
				key1, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
				Expect(err).Should(BeNil())
				key2, _, err := btctest.NewBtcKey(network, addrType)
				Expect(err).Should(BeNil())
				feeEstimator := btc.NewFixFeeEstimator(10)
				wal1, err := btc.NewWallet(network, addrType, key1, indexer, feeEstimator)
				Expect(err).Should(BeNil())

				By("Initiate an HTLC")
				amount, timelock := int64(1e7), int64(6)
				htlc, _, err := btctest.NewHtlc(key1.PubKey(), key2.PubKey(), timelock, amount)
				actions1 := []btc.HtlcAction{
					{
						ActionType: btc.HtlcActionInitiate,
						Htlc:       htlc,
					},
				}
				_, err = wal1.Execute(ctx, actions1)
				Expect(err).Should(BeNil())

				By("Mine expiry number of blocks")
				for i := 0; i < int(timelock)-1; i++ {
					err = btctest.NewBlock()
					Expect(err).Should(BeNil())
				}
				Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

				By("Refund an HTLC")
				actions2 := []btc.HtlcAction{
					{
						ActionType: btc.HtlcActionRefund,
						Htlc:       htlc,
					},
				}
				_, err = wal1.Execute(ctx, actions2)
				Expect(err).Should(BeNil())
			})

			It("should be able to instant refund an HTLC", func(ctx context.Context) {
				By("Init keys and wallet")
				addrType := waddrmgr.WitnessPubKey
				key1, _, err := btctest.NewBtcAddrWithFunds(network, addrType, nil)
				Expect(err).Should(BeNil())
				key2, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
				Expect(err).Should(BeNil())
				feeEstimator := btc.NewFixFeeEstimator(10)
				wal1, err := btc.NewWallet(network, addrType, key1, indexer, feeEstimator)
				Expect(err).Should(BeNil())
				wal2, err := btc.NewWallet(network, addrType, key2, indexer, feeEstimator)
				Expect(err).Should(BeNil())

				By("Initiate an HTLC")
				amount, timelock := int64(1e7), int64(6)
				htlc, _, err := btctest.NewHtlc(key1.PubKey(), key2.PubKey(), timelock, amount)
				actions1 := []btc.HtlcAction{
					{
						ActionType: btc.HtlcActionInitiate,
						Htlc:       htlc,
					},
				}
				initTx, err := wal1.Execute(ctx, actions1)
				Expect(err).Should(BeNil())

				By("Initiate the HTLC with the pre-signed instant refund tx")
				initUtxo := btc.UTXO{
					TxID:   initTx.TxHash().String(),
					Vout:   0,
					Amount: amount,
				}
				recipient := btc.Recipient{
					To:     wal2.Address().String(),
					Amount: amount,
				}
				refundTx, err := btc.NewInstantRefundTx(network, key1, htlc, initUtxo, recipient)
				Expect(err).Should(BeNil())

				By("Mine a new block")
				Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

				By("Use the instant refund leaf to refund the HTLC")
				actions2 := []btc.HtlcAction{
					{
						ActionType:      btc.HtlcActionInstantRefund,
						Htlc:            htlc,
						InstantRefundTx: refundTx,
					},
				}
				_, err = wal2.Execute(ctx, actions2)
				Expect(err).Should(BeNil())
			})
		})

		Context("multiple actions", func() {
			It("should be able to execute multiple execution at once", func(ctx context.Context) {
				By("Init keys and wallet")
				addrType := waddrmgr.WitnessPubKey
				key1, _, err := btctest.NewBtcAddrWithFunds(network, addrType, nil)
				Expect(err).Should(BeNil())
				key2, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
				Expect(err).Should(BeNil())
				feeEstimator := btc.NewFixFeeEstimator(10)
				wal1, err := btc.NewWallet(network, addrType, key1, indexer, feeEstimator)
				Expect(err).Should(BeNil())
				wal2, err := btc.NewWallet(network, addrType, key2, indexer, feeEstimator)
				Expect(err).Should(BeNil())

				By("Initialise some htlcs")
				batch := 4
				htlcs1 := make([]*btc.HTLC, batch*2) // inited by wallet1
				htlcs2 := make([]*btc.HTLC, batch*2) // inited by wallet2
				secrets := make([][]byte, batch*4)
				amount, timelock := int64(1e7), int64(6)
				for i := 0; i < batch*2; i++ {
					htlcs1[i], _, err = btctest.NewHtlc(key1.PubKey(), key2.PubKey(), timelock, amount)
					Expect(err).Should(BeNil())
					htlcs2[i], secrets[i], err = btctest.NewHtlc(key2.PubKey(), key1.PubKey(), timelock, amount)
					Expect(err).Should(BeNil())
				}

				By("Init the first batch of htlcs")
				actions1 := make([]btc.HtlcAction, 0, batch)
				for i := 0; i < batch; i++ {
					actions1 = append(actions1, btc.HtlcAction{
						ActionType: btc.HtlcActionInitiate,
						Htlc:       htlcs1[i],
					})
				}
				_, err = wal1.Execute(ctx, actions1)
				Expect(err).Should(BeNil())

				By("Mine expiry number of blocks")
				for i := 0; i < int(timelock)-1; i++ {
					err = btctest.NewBlock()
					Expect(err).Should(BeNil())
				}
				Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

				By("Init another two batch of htlcs for redeem and instant refunds")
				actions2 := make([]btc.HtlcAction, 0, batch)
				for i := 0; i < 2*batch; i++ {
					actions2 = append(actions2, btc.HtlcAction{
						ActionType: btc.HtlcActionInitiate,
						Htlc:       htlcs2[i],
					})
				}
				_, err = wal2.Execute(ctx, actions2)
				Expect(err).Should(BeNil())
				Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

				By("Construct instant refunds txs")
				instantRefunds := make([]*wire.MsgTx, batch)
				for i := range instantRefunds {
					utxo, err := htlcs2[i].Utxo(ctx, network, indexer)
					Expect(err).Should(BeNil())
					instantRefunds[i], err = btc.NewInstantRefundTx(network, key2, htlcs2[i], utxo, btc.Recipient{wal1.Address().EncodeAddress(), htlcs2[i].Amount})
					Expect(err).Should(BeNil())
				}

				By("Construct the last action which does all kinds of actions at once")
				actions3 := make([]btc.HtlcAction, 0, batch*4)
				for i := 0; i < len(htlcs1); i++ {
					switch {
					case i < batch:
						// Refund the first batch of htlcs1
						actions3 = append(actions3, btc.HtlcAction{
							ActionType: btc.HtlcActionRefund,
							Htlc:       htlcs1[i],
						})
						// Instant refund the first batch of htlcs2
						actions3 = append(actions3, btc.HtlcAction{
							ActionType:      btc.HtlcActionInstantRefund,
							Htlc:            htlcs2[i],
							InstantRefundTx: instantRefunds[i],
						})
					case i < 2*batch:
						// Init the second batch of htlcs1
						actions3 = append(actions3, btc.HtlcAction{
							ActionType: btc.HtlcActionInitiate,
							Htlc:       htlcs1[i],
						})
						// Redeem the second batch of htlcs2
						actions3 = append(actions3, btc.HtlcAction{
							ActionType: btc.HtlcActionRedeem,
							Htlc:       htlcs2[i],
							Secret:     secrets[i],
						})

					}
				}
				_, err = wal1.Execute(ctx, actions3)
				Expect(err).Should(BeNil())
			})
		})
	})
})
