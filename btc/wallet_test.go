package btc_test

import (
	"context"
	"math/rand"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/wire"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wallet", func() {
	Context("Operation on an HTLC", func() {
		It("should be able to initiate and redeem an HTLC", func(ctx context.Context) {
			for _, addrType := range addrTypes {
				By("Init keys and wallets")
				feeEstimator := btc.NewFixFeeEstimator(10)
				wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
				Expect(err).Should(BeNil())
				wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
				Expect(err).Should(BeNil())

				By("Initiate and redeem an HTLC")
				amount, timelock := int64(1e7), int64(144)
				htlc, secret, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
				Expect(err).Should(BeNil())
				_, _, err = wal1.Initiate(ctx, htlc)
				Expect(err).Should(BeNil())
				Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())
				_, err = wal2.Redeem(ctx, htlc, secret)
				Expect(err).Should(BeNil())

			}
		})

		It("should be able to refund an HTLC after it expires", func(ctx context.Context) {
			for _, addrType := range addrTypes {
				By("Init keys and wallets")
				feeEstimator := btc.NewFixFeeEstimator(10)
				wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
				Expect(err).Should(BeNil())
				wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
				Expect(err).Should(BeNil())

				By("Initiate an HTLC")
				amount, timelock := int64(1e7), int64(6)
				htlc, _, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
				_, _, err = wal1.Initiate(context.Background(), htlc)
				Expect(err).Should(BeNil())

				By("Mine expiry number of blocks")
				for i := 0; i < int(timelock)-1; i++ {
					err = btctest.NewBlock()
					Expect(err).Should(BeNil())
				}
				Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

				By("Refund an HTLC")
				_, err = wal1.Refund(ctx, htlc)
				Expect(err).Should(BeNil())
			}
		})

		It("should be able to instant refund a HTLC", func(ctx context.Context) {
			for _, addrType := range addrTypes {
				By("Init keys and wallet")
				feeEstimator := btc.NewFixFeeEstimator(10)
				wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
				Expect(err).Should(BeNil())
				wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
				Expect(err).Should(BeNil())

				By("Initiate an HTLC")
				amount, timelock := int64(1e7), int64(6)
				htlc, _, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
				_, irTx, err := wal1.Initiate(ctx, htlc)
				Expect(err).Should(BeNil())

				By("Mine a new block")
				Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

				By("Use the instant refund leaf to refund the HTLC")
				_, err = wal2.InstantRefund(ctx, htlc, irTx)
				Expect(err).Should(BeNil())
			}
		})
	})

	Context("Batch operation", func() {
		Context("single action", func() {
			It("should be able to initiate and redeem an HTLC", func(ctx context.Context) {
				for _, addrType := range addrTypes {
					By("Init keys and wallets")
					feeEstimator := btc.NewFixFeeEstimator(10)
					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
					Expect(err).Should(BeNil())
					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
					Expect(err).Should(BeNil())

					By("Initiate an HTLC")
					amount, timelock := int64(1e7), int64(144)
					htlc, secret, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
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
				}
			})

			It("should be able to refund an HTLC after it expires", func(ctx context.Context) {
				for _, addrType := range addrTypes {
					By("Init keys and wallets")
					feeEstimator := btc.NewFixFeeEstimator(10)
					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
					Expect(err).Should(BeNil())
					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
					Expect(err).Should(BeNil())

					By("Initiate an HTLC")
					amount, timelock := int64(1e7), int64(6)
					htlc, _, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
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
				}
			})

			It("should be able to instant refund an HTLC", func(ctx context.Context) {
				for _, addrType := range addrTypes {
					By("Init keys and wallets")
					feeEstimator := btc.NewFixFeeEstimator(10)
					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
					Expect(err).Should(BeNil())
					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
					Expect(err).Should(BeNil())

					By("Initiate an HTLC")
					amount, timelock := int64(1e7), int64(6)
					htlc, _, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
					_, irTx, err := wal1.Initiate(ctx, htlc)
					Expect(err).Should(BeNil())

					By("Mine a new block")
					Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

					By("Use the instant refund leaf to refund the HTLC")
					actions2 := []btc.HtlcAction{
						{
							ActionType:      btc.HtlcActionInstantRefund,
							Htlc:            htlc,
							InstantRefundTx: irTx,
						},
					}
					_, err = wal2.Execute(ctx, actions2)
					Expect(err).Should(BeNil())
				}
			})
		})

		Context("multiple actions", func() {
			It("should be able to execute multiple execution at once", func(ctx context.Context) {
				for _, addrType := range addrTypes {
					By("Init keys and wallets")
					feeEstimator := btc.NewFixFeeEstimator(10)
					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
					Expect(err).Should(BeNil())
					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
					Expect(err).Should(BeNil())

					By("Initialise some htlcs")
					batch := 4
					htlcs1 := make([]*btc.HTLC, batch*2) // inited by wallet1
					htlcs2 := make([]*btc.HTLC, batch*2) // inited by wallet2
					secrets := make([][]byte, batch*4)
					amount, timelock := int64(1e7), int64(6)
					for i := 0; i < batch*2; i++ {
						htlcs1[i], _, err = btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
						Expect(err).Should(BeNil())
						htlcs2[i], secrets[i], err = btctest.NewHtlc(wal2.PublicKey(), wal1.PublicKey(), timelock, amount)
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
						instantRefunds[i], err = wal2.InstantRefundTx(utxo, htlcs2[i])
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
				}

			})
		})

		Context("rbf", func() {
			It("should be able to do rbf with new actions", func(ctx context.Context) {
				for _, addrType := range addrTypes {
					By("Init keys and wallets")
					feeEstimator := btc.NewFixFeeEstimator(10)
					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
					Expect(err).Should(BeNil())
					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
					Expect(err).Should(BeNil())

					By("Prepare some init/redeem/refund/instantRefund")
					number, amount, timelock := 3, int64(1e6), int64(6)
					inits, err := generateInitHtlcs(number, amount, timelock, wal1.PublicKey(), wal2.PublicKey())
					Expect(err).Should(BeNil())
					redeems, err := generateRedeemHtlcs(ctx, number, amount, timelock, wal2.PublicKey(), wal1.PublicKey(), wal2)
					Expect(err).Should(BeNil())
					refunds, err := generateRefundHtlcs(ctx, number, amount, timelock, wal1.PublicKey(), wal2.PublicKey(), wal1)
					Expect(err).Should(BeNil())
					instantRefunds, err := generateInstantRefundHtlcs(ctx, number, amount, timelock, wal2.PublicKey(), wal1.PublicKey(), wal2)
					Expect(err).Should(BeNil())

					By("Combined all actions and shuffle the order")
					actions := append(inits, append(redeems, append(refunds, instantRefunds...)...)...)
					rand.Shuffle(len(actions), func(i, j int) {
						actions[i], actions[j] = actions[j], actions[i]
					})

					By("Execute all actions one by one using rbf")
					opts := []btc.ExecuteOpts{}
					for i := 0; i < len(actions); i++ {
						tx, err := wal1.Execute(ctx, []btc.HtlcAction{actions[i]}, opts...)
						Expect(err).Should(BeNil())
						time.Sleep(1 * time.Second)
						opts = []btc.ExecuteOpts{btc.WithRbfTxid(tx.TxHash().String())}
					}
				}
			})

			It("should make sure the rbf tx is conflicted with all previous txs", func(ctx context.Context) {
				for _, addrType := range addrTypes {
					By("Init keys and wallets")
					feeEstimator := btc.NewFixFeeEstimator(10)
					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
					Expect(err).Should(BeNil())
					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
					Expect(err).Should(BeNil())

					By("Construct two inits and one redeem")
					init1, err := generateInitHtlcs(1, 1e6, 6, wal1.PublicKey(), wal2.PublicKey())
					Expect(err).Should(BeNil())
					init2, err := generateInitHtlcs(1, 1e7, 6, wal1.PublicKey(), wal2.PublicKey())
					Expect(err).Should(BeNil())
					redeem, err := generateRedeemHtlcs(ctx, 1, 1e7, 6, wal2.PublicKey(), wal1.PublicKey(), wal2)
					Expect(err).Should(BeNil())

					By("Initiate one htlc")
					tx1, err := wal1.Execute(ctx, init1)
					Expect(err).Should(BeNil())
					color.Green("tx1: %s", tx1.TxHash().String())

					By("Redeem one htlc with higher amount")
					tx2, err := wal1.Execute(ctx, redeem, btc.WithRbfTxid(tx1.TxHash().String()))
					Expect(err).Should(BeNil())
					color.Green("tx2: %s", tx2.TxHash().String())

					By("Initiate another htlc")
					tx3, err := wal1.Execute(ctx, init2, btc.WithRbfTxid(tx2.TxHash().String()))
					Expect(err).Should(BeNil())
					color.Green("tx3: %s", tx3.TxHash().String())

					By("tx1 and tx2 should be rejected from the mempool and cannot be submitted again")
					Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(tx3.TxHash().String()))).Should(Succeed())
					time.Sleep(5 * time.Second)
					canceledCtx, cancel := context.WithCancel(ctx)
					cancel()
					_, err = indexer.GetTx(canceledCtx, tx1.TxHash().String())
					Expect(err).ToNot(BeNil())
					_, err = indexer.GetTx(canceledCtx, tx2.TxHash().String())
					Expect(err).ToNot(BeNil())
					_, err = indexer.GetTx(canceledCtx, tx3.TxHash().String())
					Expect(err).To(BeNil())
					Expect(indexer.SubmitTx(canceledCtx, tx1)).ShouldNot(Succeed())
					Expect(indexer.SubmitTx(canceledCtx, tx2)).ShouldNot(Succeed())
				}
			})

			It("should be able to do rbf when the replaced tx has descendants", func(ctx context.Context) {

			})
		})

		// Context("duplicate actions", func() {
		// 	It("should handle duplicate inits", func(ctx context.Context) {
		// 		By("Init keys and wallet")
		// 		addrType := waddrmgr.WitnessPubKey
		// 		key1, _, err := btctest.NewBtcAddrWithFunds(network, addrType, nil)
		// 		Expect(err).Should(BeNil())
		// 		key2, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
		// 		Expect(err).Should(BeNil())
		// 		feeEstimator := btc.NewFixFeeEstimator(10)
		// 		wal1, err := btc.NewWallet(network, addrType, key1, indexer, feeEstimator)
		// 		Expect(err).Should(BeNil())
		// 		_, err = btc.NewWallet(network, addrType, key2, indexer, feeEstimator)
		// 		Expect(err).Should(BeNil())
		//
		// 		By("Initiate one htlc")
		// 		init1, err := generateInitHtlcs(1, 1e6, 6, wal1.PublicKey(), key2.PubKey())
		// 		Expect(err).Should(BeNil())
		// 		tx1, err := wal1.Execute(ctx, init1)
		// 		Expect(err).Should(BeNil())
		// 		color.Green("tx1: %s", tx1.TxHash().String())
		//
		// 		By("Initiate again with duplicate actions")
		// 		init2, err := generateInitHtlcs(1, 1e7, 6, key1.PubKey(), key2.PubKey())
		// 		Expect(err).Should(BeNil())
		// 		tx2, err := wal1.Execute(ctx, append(init1, append(init2, init2...)...), btc.WithRbfTxid(tx1.TxHash().String()))
		// 		Expect(err).Should(BeNil())
		// 		color.Green("tx2: %s", tx2.TxHash().String())
		// 	})
		//
		// 	It("should handle duplicate redeems", func(ctx context.Context) {
		// 		By("Init keys and wallet")
		// 		addrType := waddrmgr.WitnessPubKey
		// 		key1, _, err := btctest.NewBtcAddrWithFunds(network, addrType, nil)
		// 		Expect(err).Should(BeNil())
		// 		key2, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
		// 		Expect(err).Should(BeNil())
		// 		feeEstimator := btc.NewFixFeeEstimator(10)
		// 		wal1, err := btc.NewWallet(network, addrType, key1, indexer, feeEstimator)
		// 		Expect(err).Should(BeNil())
		// 		wal2, err := btc.NewWallet(network, addrType, key2, indexer, feeEstimator)
		// 		Expect(err).Should(BeNil())
		//
		// 		By("Redeem one htlc")
		// 		redeems, err := generateRedeemHtlcs(ctx, 2, 1e6, 6, key2.PubKey(), key1.PubKey(), wal2)
		// 		Expect(err).Should(BeNil())
		// 		tx1, err := wal1.Execute(ctx, redeems[:1])
		// 		Expect(err).Should(BeNil())
		// 		color.Green("tx1: %s", tx1.TxHash().String())
		//
		// 		By("Redeem again with duplicate actions")
		// 		tx2, err := wal1.Execute(ctx, append(redeems, redeems...), btc.WithRbfTxid(tx1.TxHash().String()))
		// 		Expect(err).Should(BeNil())
		// 		color.Green("tx2: %s", tx2.TxHash().String())
		// 	})
		//
		// 	It("should handle duplicate refunds", func(ctx context.Context) {
		// 		By("Init keys and wallet")
		// 		addrType := waddrmgr.WitnessPubKey
		// 		key1, _, err := btctest.NewBtcAddrWithFunds(network, addrType, nil)
		// 		Expect(err).Should(BeNil())
		// 		key2, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
		// 		Expect(err).Should(BeNil())
		// 		feeEstimator := btc.NewFixFeeEstimator(10)
		// 		wal1, err := btc.NewWallet(network, addrType, key1, indexer, feeEstimator)
		// 		Expect(err).Should(BeNil())
		//
		// 		By("Refund one htlc")
		// 		refunds, err := generateRefundHtlcs(ctx, 2, 1e6, 6, key1.PubKey(), key2.PubKey(), wal1)
		// 		Expect(err).Should(BeNil())
		// 		tx1, err := wal1.Execute(ctx, refunds[:1])
		// 		Expect(err).Should(BeNil())
		// 		color.Green("tx1: %s", tx1.TxHash().String())
		//
		// 		By("Refund again with duplicate actions")
		// 		tx2, err := wal1.Execute(ctx, append(refunds, refunds...), btc.WithRbfTxid(tx1.TxHash().String()))
		// 		Expect(err).Should(BeNil())
		// 		color.Green("tx2: %s", tx2.TxHash().String())
		// 	})
		//
		// 	It("should handle duplicate instant refunds", func(ctx context.Context) {
		// 		By("Init keys and wallet")
		// 		addrType := waddrmgr.WitnessPubKey
		// 		key1, _, err := btctest.NewBtcAddrWithFunds(network, addrType, nil)
		// 		Expect(err).Should(BeNil())
		// 		key2, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
		// 		Expect(err).Should(BeNil())
		// 		feeEstimator := btc.NewFixFeeEstimator(10)
		// 		wal1, err := btc.NewWallet(network, addrType, key1, indexer, feeEstimator)
		// 		Expect(err).Should(BeNil())
		// 		wal2, err := btc.NewWallet(network, addrType, key2, indexer, feeEstimator)
		// 		Expect(err).Should(BeNil())
		//
		// 		By("Instant refunds one htlc")
		// 		instantRefunds, err := generateInstantRefundHtlcs(ctx, 2, 1e6, 6, key2.PubKey(), key1.PubKey(), wal2, key2)
		// 		Expect(err).Should(BeNil())
		// 		tx1, err := wal1.Execute(ctx, instantRefunds[:1])
		// 		Expect(err).Should(BeNil())
		// 		color.Green("tx1: %s", tx1.TxHash().String())
		//
		// 		By("Instant refunds again with duplicate actions")
		// 		tx2, err := wal1.Execute(ctx, append(instantRefunds, instantRefunds...), btc.WithRbfTxid(tx1.TxHash().String()))
		// 		Expect(err).Should(BeNil())
		// 		color.Green("tx2: %s", tx2.TxHash().String())
		// 	})
		// })
	})
})

func generateInitHtlcs(n int, amount, timelock int64, initiatorPub, redeemerPub *btcec.PublicKey) ([]btc.HtlcAction, error) {
	actions := make([]btc.HtlcAction, n)
	for i := 0; i < n; i++ {
		htlc, secret, err := btctest.NewHtlc(initiatorPub, redeemerPub, timelock, amount)
		if err != nil {
			return nil, err
		}
		actions[i] = btc.HtlcAction{
			Htlc:       htlc,
			Secret:     secret,
			ActionType: btc.HtlcActionInitiate,
		}
	}
	return actions, nil
}

func generateRedeemHtlcs(ctx context.Context, n int, amount, timelock int64, initiatorPub, redeemerPub *btcec.PublicKey, wallet btc.Wallet) ([]btc.HtlcAction, error) {
	actions1 := make([]btc.HtlcAction, n)
	actions2 := make([]btc.HtlcAction, n)
	for i := 0; i < n; i++ {
		var err error
		htlc, secret, err := btctest.NewHtlc(initiatorPub, redeemerPub, timelock, amount)
		if err != nil {
			return nil, err
		}
		actions1[i] = btc.HtlcAction{
			Htlc:       htlc,
			ActionType: btc.HtlcActionInitiate,
		}

		actions2[i] = btc.HtlcAction{
			Htlc:       htlc,
			ActionType: btc.HtlcActionRedeem,
			Secret:     secret,
		}
	}

	_, err := wallet.Execute(ctx, actions1)
	if err != nil {
		return nil, err
	}
	return actions2, btctest.NewBlockWaitMined(indexer)
}

func generateRefundHtlcs(ctx context.Context, n int, amount, timelock int64, initiatorPub, redeemerPub *btcec.PublicKey, wallet btc.Wallet) ([]btc.HtlcAction, error) {
	actions1 := make([]btc.HtlcAction, n)
	actions2 := make([]btc.HtlcAction, n)
	for i := 0; i < n; i++ {
		var err error
		htlc, secret, err := btctest.NewHtlc(initiatorPub, redeemerPub, timelock, amount)
		if err != nil {
			return nil, err
		}
		actions1[i] = btc.HtlcAction{
			Htlc:       htlc,
			ActionType: btc.HtlcActionInitiate,
		}

		actions2[i] = btc.HtlcAction{
			Htlc:       htlc,
			ActionType: btc.HtlcActionRefund,
			Secret:     secret,
		}
	}

	_, err := wallet.Execute(ctx, actions1)
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(timelock)-1; i++ {
		if err := btctest.NewBlock(); err != nil {
			return nil, err
		}
	}
	if err := btctest.NewBlockWaitMined(indexer); err != nil {
		return nil, err
	}
	return actions2, nil
}

func generateInstantRefundHtlcs(ctx context.Context, n int, amount, timelock int64, initiatorPub, redeemerPub *btcec.PublicKey, wallet btc.Wallet) ([]btc.HtlcAction, error) {
	actions1 := make([]btc.HtlcAction, n)
	for i := 0; i < n; i++ {
		var err error
		htlc, _, err := btctest.NewHtlc(initiatorPub, redeemerPub, timelock, amount)
		if err != nil {
			return nil, err
		}
		actions1[i] = btc.HtlcAction{
			Htlc:       htlc,
			ActionType: btc.HtlcActionInitiate,
		}
	}

	_, err := wallet.Execute(ctx, actions1)
	if err != nil {
		return nil, err
	}
	if err := btctest.NewBlockWaitMined(indexer); err != nil {
		return nil, err
	}
	actions2 := make([]btc.HtlcAction, n)
	for i := 0; i < n; i++ {
		utxo, err := actions1[i].Htlc.Utxo(ctx, network, indexer)
		if err != nil {
			return nil, err
		}
		instantRefunds, err := wallet.InstantRefundTx(utxo, actions1[i].Htlc)
		actions2[i] = btc.HtlcAction{
			Htlc:            actions1[i].Htlc,
			ActionType:      btc.HtlcActionInstantRefund,
			Secret:          actions1[i].Secret,
			InstantRefundTx: instantRefunds,
		}
	}

	return actions2, nil
}
