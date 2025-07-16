package btc_test

//
// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"time"
//
// 	"github.com/btcsuite/btcd/wire"
// 	"github.com/catalogfi/blockchain/btc"
// 	"github.com/catalogfi/blockchain/btc/btctest"
// 	"github.com/fatih/color"
//
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// )
//
// var _ = Describe("Wallet", func() {
// 	Context("Operation on an HTLC", func() {
// 		It("should be able to initiate and redeem an HTLC", func(ctx context.Context) {
// 			for _, addrType := range addrTypes {
// 				By("Init keys and wallets")
// 				feeEstimator := btc.NewFixFeeEstimator(1e3)
// 				wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 				Expect(err).Should(BeNil())
// 				wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 				Expect(err).Should(BeNil())
//
// 				By("Initiate and redeem an HTLC")
// 				amount, timelock := int64(1e7), int64(144)
// 				htlc, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
// 				Expect(err).Should(BeNil())
// 				_, _, err = wal1.Initiate(ctx, htlc)
// 				Expect(err).Should(BeNil())
// 				Expect(btctest.NewBlockWaitMined(1, indexer)).Should(Succeed())
// 				_, err = wal2.Redeem(ctx, htlc)
// 				Expect(err).Should(BeNil())
// 			}
// 		})
//
// 		It("should be able to refund an HTLC after it expires", func(ctx context.Context) {
// 			for _, addrType := range addrTypes {
// 				By("Init keys and wallets")
// 				feeEstimator := btc.NewFixFeeEstimator(1e3)
// 				wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 				Expect(err).Should(BeNil())
// 				wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 				Expect(err).Should(BeNil())
//
// 				By("Initiate an HTLC")
// 				amount, timelock := int64(1e7), int64(17)
// 				htlc, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
// 				_, _, err = wal1.Initiate(context.Background(), htlc)
// 				Expect(err).Should(BeNil())
//
// 				By("Mine expiry number of blocks")
// 				Expect(btctest.NewBlockWaitMined(int(timelock), indexer)).Should(Succeed())
//
// 				By("Refund an HTLC")
// 				_, err = wal1.Refund(ctx, htlc, nil)
// 				Expect(err).Should(BeNil())
// 			}
// 		})
//
// 		It("should be able to refund an HTLC to a different address", func(ctx context.Context) {
// 			for _, addrType := range addrTypes {
// 				By("Init keys and wallets")
// 				feeEstimator := btc.NewFixFeeEstimator(1e3)
// 				wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 				Expect(err).Should(BeNil())
// 				wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 				Expect(err).Should(BeNil())
// 				_, addr3, err := btctest.NewBtcKey(network, addrType)
// 				Expect(err).Should(BeNil())
//
// 				By("Initiate an HTLC")
// 				amount, timelock := int64(1e7), int64(17)
// 				htlc, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
// 				_, _, err = wal1.Initiate(context.Background(), htlc)
// 				Expect(err).Should(BeNil())
//
// 				By("Mine expiry number of blocks")
// 				Expect(btctest.NewBlockWaitMined(int(timelock), indexer)).Should(Succeed())
//
// 				By("Refund an HTLC")
// 				tx, err := wal1.Refund(ctx, htlc, addr3)
// 				Expect(err).Should(BeNil())
//
// 				By("Target address should have the balance")
// 				Eventually(func() error {
// 					utxos, err := indexer.GetUTXOs(ctx, addr3)
// 					if err != nil {
// 						return err
// 					}
// 					for _, utxo := range utxos {
// 						if utxo.TxID == tx.TxHash().String() {
// 							return nil
// 						}
// 					}
// 					return fmt.Errorf("not found")
// 				}, 10*time.Second, 1*time.Second).Should(Succeed())
// 			}
// 		})
//
// 		FIt("should be able to instant refund a HTLC", func(ctx context.Context) {
// 			for _, addrType := range addrTypes {
// 				By("Init keys and wallet")
// 				feeEstimator := btc.NewFixFeeEstimator(1e3)
// 				wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 				Expect(err).Should(BeNil())
// 				wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 				Expect(err).Should(BeNil())
//
// 				By("Initiate an HTLC")
// 				amount, timelock := int64(1e7), int64(6)
// 				htlc, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
// 				_, irTx, err := wal1.Initiate(ctx, htlc)
// 				Expect(err).Should(BeNil())
//
// 				By("Mine a new block")
// 				Expect(btctest.NewBlockWaitMined(1, indexer)).Should(Succeed())
//
// 				By("Use the instant refund leaf to refund the HTLC")
// 				_, err = wal2.InstantRefund(ctx, htlc, irTx)
// 				Expect(err).Should(BeNil())
// 			}
// 		})
// 	})
//
// 	Context("Batch operation", func() {
// 		Context("single action", func() {
// 			It("should be able to initiate and redeem an HTLC", func(ctx context.Context) {
// 				for _, addrType := range addrTypes {
// 					By("Init keys and wallets")
// 					feeEstimator := btc.NewFixFeeEstimator(10e3)
// 					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 					Expect(err).Should(BeNil())
// 					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 					Expect(err).Should(BeNil())
//
// 					By("Initiate an HTLC")
// 					amount, timelock := int64(1e7), int64(144)
// 					htlc, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
// 					Expect(err).Should(BeNil())
// 					actions1 := []btc.HtlcAction{
// 						{
// 							ActionType: btc.HtlcActionInitiate,
// 							Htlc:       htlc,
// 						},
// 					}
// 					_, err = wal1.Execute(ctx, actions1, "")
// 					Expect(err).Should(BeNil())
// 					Expect(btctest.NewBlockWaitMined(1, indexer)).Should(Succeed())
//
// 					By("Redeem an HTLC")
// 					actions2 := []btc.HtlcAction{
// 						{
// 							ActionType: btc.HtlcActionRedeem,
// 							Htlc:       htlc,
// 						},
// 					}
// 					_, err = wal2.Execute(ctx, actions2, "")
// 					Expect(err).Should(BeNil())
// 				}
// 			})
//
// 			It("should be able to refund an HTLC after it expires", func(ctx context.Context) {
// 				for _, addrType := range addrTypes {
// 					By("Init keys and wallets")
// 					feeEstimator := btc.NewFixFeeEstimator(10e3)
// 					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 					Expect(err).Should(BeNil())
// 					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 					Expect(err).Should(BeNil())
//
// 					By("Initiate an HTLC")
// 					amount, timelock := int64(1e7), int64(6)
// 					htlc, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
// 					actions1 := []btc.HtlcAction{
// 						{
// 							ActionType: btc.HtlcActionInitiate,
// 							Htlc:       htlc,
// 						},
// 					}
// 					_, err = wal1.Execute(ctx, actions1, "")
// 					Expect(err).Should(BeNil())
//
// 					By("Mine expiry number of blocks")
// 					Expect(btctest.NewBlockWaitMined(int(timelock), indexer)).Should(Succeed())
//
// 					By("Refund an HTLC")
// 					actions2 := []btc.HtlcAction{
// 						{
// 							ActionType: btc.HtlcActionRefund,
// 							Htlc:       htlc,
// 						},
// 					}
// 					_, err = wal1.Execute(ctx, actions2, "")
// 					Expect(err).Should(BeNil())
// 				}
// 			})
//
// 			It("should be able to refund an HTLC after it expires", func(ctx context.Context) {
// 				for _, addrType := range addrTypes {
// 					By("Init keys and wallets")
// 					feeEstimator := btc.NewFixFeeEstimator(10e3)
// 					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 					Expect(err).Should(BeNil())
// 					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 					Expect(err).Should(BeNil())
// 					_, addr3, err := btctest.NewBtcKey(network, addrType)
// 					Expect(err).Should(BeNil())
// 					log.Print("addr3 = ", addr3.EncodeAddress())
//
// 					By("Initiate an HTLC")
// 					amount, timelock := int64(1e7), int64(6)
// 					htlc, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
// 					Expect(err).Should(BeNil())
// 					htlc.Amount -= 1000
// 					actions1 := []btc.HtlcAction{
// 						{
// 							ActionType: btc.HtlcActionInitiate,
// 							Htlc:       htlc,
// 						},
// 					}
// 					_, err = wal1.Execute(ctx, actions1, "")
// 					Expect(err).Should(BeNil())
//
// 					By("Mine expiry number of blocks")
// 					Expect(btctest.NewBlockWaitMined(int(timelock), indexer)).Should(Succeed())
//
// 					By("Refund an HTLC")
// 					actions2 := []btc.HtlcAction{
// 						{
// 							ActionType: btc.HtlcActionRefund,
// 							Htlc:       htlc,
// 							RefundTo:   addr3,
// 						},
// 					}
// 					_, err = wal1.Execute(ctx, actions2, "")
// 					Expect(err).Should(BeNil())
// 				}
// 			})
//
// 			It("should be able to refund an HTLC to a different address", func(ctx context.Context) {
// 				for _, addrType := range addrTypes {
// 					By("Init keys and wallets")
// 					feeEstimator := btc.NewFixFeeEstimator(10e3)
// 					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 					Expect(err).Should(BeNil())
// 					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 					Expect(err).Should(BeNil())
// 					_, addr3, err := btctest.NewBtcKey(network, addrType)
// 					Expect(err).Should(BeNil())
//
// 					By("Initiate an HTLC")
// 					amount, timelock := int64(1e7), int64(6)
// 					htlc, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
// 					actions1 := []btc.HtlcAction{
// 						{
// 							ActionType: btc.HtlcActionInitiate,
// 							Htlc:       htlc,
// 						},
// 					}
// 					_, err = wal1.Execute(ctx, actions1, "")
// 					Expect(err).Should(BeNil())
//
// 					By("Mine expiry number of blocks")
// 					Expect(btctest.NewBlockWaitMined(int(timelock), indexer)).Should(Succeed())
//
// 					By("Refund an HTLC")
// 					actions2 := []btc.HtlcAction{
// 						{
// 							ActionType: btc.HtlcActionRefund,
// 							Htlc:       htlc,
// 							RefundTo:   addr3,
// 						},
// 					}
// 					tx, err := wal1.Execute(ctx, actions2, "")
// 					Expect(err).Should(BeNil())
//
// 					By("Target address should have the balance")
// 					Eventually(func() error {
// 						utxos, err := indexer.GetUTXOs(ctx, addr3)
// 						if err != nil {
// 							return err
// 						}
// 						for _, utxo := range utxos {
// 							if utxo.TxID == tx.TxHash().String() && utxo.Amount == htlc.Amount {
// 								return nil
// 							}
// 						}
// 						return fmt.Errorf("not found")
// 					}, 10*time.Second, 1*time.Second).Should(Succeed())
// 				}
// 			})
//
// 			It("should be able to instant refund an HTLC", func(ctx context.Context) {
// 				for _, addrType := range addrTypes {
// 					By("Init keys and wallets")
// 					feeEstimator := btc.NewFixFeeEstimator(1e3)
// 					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 					Expect(err).Should(BeNil())
// 					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 					Expect(err).Should(BeNil())
//
// 					By("Initiate an HTLC")
// 					amount, timelock := int64(1e7), int64(6)
// 					htlc, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
// 					_, irTx, err := wal1.Initiate(ctx, htlc)
// 					Expect(err).Should(BeNil())
//
// 					By("Mine a new block")
// 					Expect(btctest.NewBlockWaitMined(1, indexer)).Should(Succeed())
//
// 					By("Use the instant refund leaf to refund the HTLC")
// 					actions2 := []btc.HtlcAction{
// 						{
// 							ActionType:      btc.HtlcActionInstantRefund,
// 							Htlc:            htlc,
// 							InstantRefundTx: irTx,
// 						},
// 					}
// 					_, err = wal2.Execute(ctx, actions2, "")
// 					Expect(err).Should(BeNil())
// 				}
// 			})
//
// 			It("should be able to instant refund with a different value ", func(ctx context.Context) {
// 				for _, addrType := range addrTypes {
// 					By("Init keys and wallets")
// 					feeEstimator := btc.NewFixFeeEstimator(1e3)
// 					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 					Expect(err).Should(BeNil())
// 					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 					Expect(err).Should(BeNil())
//
// 					By("Initiate an HTLC and fund it with a different amount")
// 					amount, timelock := int64(1e7), int64(6)
// 					htlc, err := btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
// 					htlcAddr := htlc.MustAddress(network)
// 					_, err = btctest.Faucet(htlcAddr.EncodeAddress())
// 					Expect(err).Should(BeNil())
// 					time.Sleep(5 * time.Second)
//
// 					By("Create the instant refund tx")
// 					utxos, err := indexer.GetUTXOs(ctx, htlcAddr)
// 					Expect(err).Should(BeNil())
// 					Expect(len(utxos)).Should(Equal(1))
// 					irTx, err := wal1.InstantRefundTx(utxos[0], htlc)
// 					Expect(err).Should(BeNil())
//
// 					By("Mine a new block")
// 					Expect(btctest.NewBlockWaitMined(1, indexer)).Should(Succeed())
//
// 					By("Use the instant refund leaf to refund the HTLC")
// 					actions2 := []btc.HtlcAction{
// 						{
// 							ActionType:      btc.HtlcActionInstantRefund,
// 							Htlc:            htlc,
// 							InstantRefundTx: irTx,
// 						},
// 					}
// 					_, err = wal2.Execute(ctx, actions2, "")
// 					Expect(err).Should(BeNil())
// 				}
// 			})
// 		})
//
// 		Context("multiple actions", func() {
// 			It("should be able to execute multiple execution at once", func(ctx context.Context) {
// 				for _, addrType := range addrTypes {
// 					By("Init keys and wallets")
// 					feeEstimator := btc.NewFixFeeEstimator(10e3)
// 					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 					Expect(err).Should(BeNil())
// 					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 					Expect(err).Should(BeNil())
//
// 					By("Initialise some htlcs")
// 					batch := 4
// 					htlcs1 := make([]*btc.HTLC, batch*2) // inited by wallet1
// 					htlcs2 := make([]*btc.HTLC, batch*2) // inited by wallet2
// 					amount, timelock := int64(1e7), int64(6)
// 					for i := 0; i < batch*2; i++ {
// 						htlcs1[i], err = btctest.NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
// 						Expect(err).Should(BeNil())
// 						htlcs2[i], err = btctest.NewHtlc(wal2.PublicKey(), wal1.PublicKey(), timelock, amount)
// 						Expect(err).Should(BeNil())
// 					}
//
// 					By("Init the first batch of htlcs")
// 					actions1 := make([]btc.HtlcAction, 0, batch)
// 					for i := 0; i < batch; i++ {
// 						actions1 = append(actions1, btc.HtlcAction{
// 							ActionType: btc.HtlcActionInitiate,
// 							Htlc:       htlcs1[i],
// 						})
// 					}
// 					_, err = wal1.Execute(ctx, actions1, "")
// 					Expect(err).Should(BeNil())
//
// 					By("Mine expiry number of blocks")
// 					Expect(btctest.NewBlockWaitMined(int(timelock), indexer)).Should(Succeed())
//
// 					By("Init another two batch of htlcs for redeem and instant refunds")
// 					actions2 := make([]btc.HtlcAction, 0, batch)
// 					for i := 0; i < 2*batch; i++ {
// 						actions2 = append(actions2, btc.HtlcAction{
// 							ActionType: btc.HtlcActionInitiate,
// 							Htlc:       htlcs2[i],
// 						})
// 					}
// 					_, err = wal2.Execute(ctx, actions2, "")
// 					Expect(err).Should(BeNil())
// 					Expect(btctest.NewBlockWaitMined(1, indexer)).Should(Succeed())
//
// 					By("Construct instant refunds txs")
// 					instantRefunds := make([]*wire.MsgTx, batch)
// 					for i := range instantRefunds {
// 						utxo, err := htlcs2[i].Utxo(ctx, network, indexer)
// 						Expect(err).Should(BeNil())
// 						instantRefunds[i], err = wal2.InstantRefundTx(utxo, htlcs2[i])
// 						Expect(err).Should(BeNil())
// 					}
//
// 					By("Construct the last action which does all kinds of actions at once")
// 					actions3 := make([]btc.HtlcAction, 0, batch*4)
// 					for i := 0; i < len(htlcs1); i++ {
// 						switch {
// 						case i < batch:
// 							// Refund the first batch of htlcs1
// 							actions3 = append(actions3, btc.HtlcAction{
// 								ActionType: btc.HtlcActionRefund,
// 								Htlc:       htlcs1[i],
// 							})
// 							// Instant refund the first batch of htlcs2
// 							actions3 = append(actions3, btc.HtlcAction{
// 								ActionType:      btc.HtlcActionInstantRefund,
// 								Htlc:            htlcs2[i],
// 								InstantRefundTx: instantRefunds[i],
// 							})
// 						case i < 2*batch:
// 							// Init the second batch of htlcs1
// 							actions3 = append(actions3, btc.HtlcAction{
// 								ActionType: btc.HtlcActionInitiate,
// 								Htlc:       htlcs1[i],
// 							})
// 							// Redeem the second batch of htlcs2
// 							actions3 = append(actions3, btc.HtlcAction{
// 								ActionType: btc.HtlcActionRedeem,
// 								Htlc:       htlcs2[i],
// 							})
//
// 						}
// 					}
// 					_, err = wal1.Execute(ctx, actions3, "")
// 					Expect(err).Should(BeNil())
// 				}
//
// 			})
// 		})
//
// 		Context("rbf", func() {
// 			It("should be able to do rbf with new actions", func(ctx context.Context) {
// 				for _, addrType := range addrTypes {
// 					By("Init keys and wallets")
// 					feeEstimator := btc.NewFixFeeEstimator(10e3)
// 					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 					Expect(err).Should(BeNil())
// 					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 					Expect(err).Should(BeNil())
//
// 					By("Prepare some init/redeem/refund/instantRefund")
// 					number := 3
// 					inits, err := btctest.PrepareActions(ctx, number, wal1, wal2, indexer, btc.HtlcActionInitiate)
// 					Expect(err).Should(BeNil())
// 					redeems, err := btctest.PrepareActions(ctx, number, wal2, wal1, indexer, btc.HtlcActionRedeem)
// 					Expect(err).Should(BeNil())
// 					refunds, err := btctest.PrepareActions(ctx, number, wal1, wal2, indexer, btc.HtlcActionRefund)
// 					Expect(err).Should(BeNil())
// 					instantRefunds, err := btctest.PrepareActions(ctx, number, wal2, wal1, indexer, btc.HtlcActionInitiate)
// 					Expect(err).Should(BeNil())
//
// 					By("Combined all actions and shuffle the order")
// 					actions := append(inits, append(redeems, append(refunds, instantRefunds...)...)...)
// 					rand.Shuffle(len(actions), func(i, j int) {
// 						actions[i], actions[j] = actions[j], actions[i]
// 					})
//
// 					By("Execute all actions one by one using rbf")
// 					prevTxid := ""
// 					for i := 0; i < len(actions); i++ {
// 						tx, err := wal1.Execute(ctx, []btc.HtlcAction{actions[i]}, prevTxid)
// 						Expect(err).Should(BeNil())
// 						time.Sleep(1 * time.Second)
// 						prevTxid = tx.TxHash().String()
// 					}
// 				}
// 			})
//
// 			It("should make sure the rbf tx is conflicted with all previous txs", func(ctx context.Context) {
// 				for _, addrType := range addrTypes {
// 					By("Init keys and wallets")
// 					feeEstimator := btc.NewFixFeeEstimator(10e3)
// 					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 					Expect(err).Should(BeNil())
// 					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 					Expect(err).Should(BeNil())
//
// 					By("Construct two inits and one redeem")
// 					inits, err := btctest.PrepareActions(ctx, 2, wal1, wal2, indexer, btc.HtlcActionInitiate)
// 					Expect(err).Should(BeNil())
// 					Expect(len(inits)).Should(Equal(2))
// 					init1, init2 := inits[:1], inits[1:]
// 					redeem, err := btctest.PrepareActions(ctx, 1, wal2, wal1, indexer, btc.HtlcActionRedeem)
// 					Expect(err).Should(BeNil())
//
// 					By("Initiate one htlc")
// 					tx1, err := wal1.Execute(ctx, init1, "")
// 					Expect(err).Should(BeNil())
// 					color.Green("tx1: %s", tx1.TxHash().String())
//
// 					By("Redeem one htlc with higher amount")
// 					tx2, err := wal1.Execute(ctx, redeem, tx1.TxHash().String())
// 					Expect(err).Should(BeNil())
// 					color.Green("tx2: %s", tx2.TxHash().String())
//
// 					By("Initiate another htlc")
// 					tx3, err := wal1.Execute(ctx, init2, tx2.TxHash().String())
// 					Expect(err).Should(BeNil())
// 					color.Green("tx3: %s", tx3.TxHash().String())
//
// 					By("tx1 and tx2 should be rejected from the mempool and cannot be submitted again")
// 					Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(tx3.TxHash().String()))).Should(Succeed())
// 					time.Sleep(5 * time.Second)
// 					canceledCtx, cancel := context.WithCancel(ctx)
// 					cancel()
// 					_, err = indexer.GetTx(canceledCtx, tx1.TxHash().String())
// 					Expect(err).ToNot(BeNil())
// 					_, err = indexer.GetTx(canceledCtx, tx2.TxHash().String())
// 					Expect(err).ToNot(BeNil())
// 					_, err = indexer.GetTx(canceledCtx, tx3.TxHash().String())
// 					Expect(err).To(BeNil())
// 					Expect(indexer.SubmitTx(canceledCtx, tx1)).ShouldNot(Succeed())
// 					Expect(indexer.SubmitTx(canceledCtx, tx2)).ShouldNot(Succeed())
// 				}
// 			})
//
// 			It("should be able to do rbf when the replaced tx has descendants", func(ctx context.Context) {
// 				// todo
// 			})
// 		})
//
// 		Context("duplicate actions", func() {
// 			It("should handle duplicate inits", func(ctx context.Context) {
// 				for _, addrType := range addrTypes {
// 					By("Init keys and wallets")
// 					feeEstimator := btc.NewFixFeeEstimator(10e3)
// 					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 					Expect(err).Should(BeNil())
// 					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 					Expect(err).Should(BeNil())
//
// 					By("Initiate one htlc")
// 					inits1, err := btctest.PrepareActions(ctx, 1, wal1, wal2, indexer, btc.HtlcActionInitiate)
// 					Expect(err).Should(BeNil())
// 					tx1, err := wal1.Execute(ctx, inits1, "")
// 					Expect(err).Should(BeNil())
// 					color.Green("tx1: %s", tx1.TxHash().String())
//
// 					By("Initiate again with duplicate actions")
// 					inits2, err := btctest.PrepareActions(ctx, 1, wal1, wal2, indexer, btc.HtlcActionInitiate)
// 					Expect(err).Should(BeNil())
// 					tx2, err := wal1.Execute(ctx, append(inits1, append(inits2, inits2...)...), tx1.TxHash().String())
// 					Expect(err).Should(BeNil())
// 					color.Green("tx2: %s", tx2.TxHash().String())
// 				}
// 			})
//
// 			It("should handle duplicate redeems", func(ctx context.Context) {
// 				for _, addrType := range addrTypes {
// 					By("Init keys and wallets")
// 					feeEstimator := btc.NewFixFeeEstimator(10e3)
// 					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 					Expect(err).Should(BeNil())
// 					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 					Expect(err).Should(BeNil())
//
// 					By("Redeem one htlc")
// 					redeems, err := btctest.PrepareActions(ctx, 2, wal2, wal1, indexer, btc.HtlcActionRedeem)
// 					Expect(err).Should(BeNil())
// 					tx1, err := wal1.Execute(ctx, redeems[:1], "")
// 					Expect(err).Should(BeNil())
// 					color.Green("tx1: %s", tx1.TxHash().String())
//
// 					By("Redeem again with duplicate actions")
// 					tx2, err := wal1.Execute(ctx, append(redeems, redeems...), tx1.TxHash().String())
// 					Expect(err).Should(BeNil())
// 					color.Green("tx2: %s", tx2.TxHash().String())
// 				}
// 			})
//
// 			It("should handle duplicate refunds", func(ctx context.Context) {
// 				for _, addrType := range addrTypes {
// 					By("Init keys and wallets")
// 					feeEstimator := btc.NewFixFeeEstimator(10e3)
// 					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 					Expect(err).Should(BeNil())
// 					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 					Expect(err).Should(BeNil())
//
// 					By("Refund one htlc")
// 					refunds, err := btctest.PrepareActions(ctx, 2, wal1, wal2, indexer, btc.HtlcActionRefund)
// 					Expect(err).Should(BeNil())
// 					tx1, err := wal1.Execute(ctx, refunds[:1], "")
// 					Expect(err).Should(BeNil())
// 					color.Green("tx1: %s", tx1.TxHash().String())
//
// 					By("Refund again with duplicate actions")
// 					tx2, err := wal1.Execute(ctx, append(refunds, refunds...), tx1.TxHash().String())
// 					Expect(err).Should(BeNil())
// 					color.Green("tx2: %s", tx2.TxHash().String())
// 				}
// 			})
//
// 			It("should handle duplicate instant refunds", func(ctx context.Context) {
// 				for _, addrType := range addrTypes {
// 					By("Init keys and wallets")
// 					feeEstimator := btc.NewFixFeeEstimator(10e3)
// 					wal1, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, false)
// 					Expect(err).Should(BeNil())
// 					wal2, err := btctest.NewWallet(network, addrType, indexer, client, feeEstimator, true)
// 					Expect(err).Should(BeNil())
//
// 					By("Instant refunds one htlc")
// 					instantRefunds, err := btctest.PrepareActions(ctx, 2, wal2, wal1, indexer, btc.HtlcActionInstantRefund)
// 					Expect(err).Should(BeNil())
// 					tx1, err := wal1.Execute(ctx, instantRefunds[:1], "")
// 					Expect(err).Should(BeNil())
// 					color.Green("tx1: %s", tx1.TxHash().String())
//
// 					By("Instant refunds again with duplicate actions")
// 					tx2, err := wal1.Execute(ctx, append(instantRefunds, instantRefunds...), tx1.TxHash().String())
// 					Expect(err).Should(BeNil())
// 					color.Green("tx2: %s", tx2.TxHash().String())
// 				}
// 			})
// 		})
// 	})
// })
