package cosigner_test

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/catalogfi/blockchain/btc/cosigner"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("wallet", func() {
	Context("htlc actions", func() {
		It("should be able to initiate a few htlcs", func(ctx context.Context) {
			By("Create a new wallet")
			key1, err := btcec.NewPrivateKey()
			Expect(err).Should(BeNil())
			feeEstimator := btc.NewFixFeeEstimator(1e3)
			wallet, err := cosigner.NewWallet(network, key1, cosignerPub, indexer, client, btcClient, feeEstimator)
			Expect(err).Should(BeNil())

			By("Fund the wallet")
			fundTx, err := btctest.Faucet(wallet.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(fundTx.String()))).Should(Succeed())

			By("Create a new htlc")
			key2, _, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).Should(BeNil())
			for i := 0; i < 5; i++ {
				time.Sleep(5 * time.Second)
				amount, timelock := int64(1e7), int64(144)
				htlc, err := btctest.NewHtlc(key1.PubKey(), key2.PubKey(), timelock, amount)
				Expect(err).Should(BeNil())

				action := []btc.HtlcAction{
					{
						ActionType: btc.HtlcActionInitiate,
						Htlc:       htlc,
					},
				}
				tx, _, err := wallet.Execute(ctx, action)
				Expect(err).Should(BeNil())
				log.Printf("tx = %v", tx.TxHash().String())
			}
		})

		It("should be able to initiate a few htlcs and include merges", func(ctx context.Context) {
			By("Create new wallets")
			key1, err := btcec.NewPrivateKey()
			Expect(err).Should(BeNil())
			key2, _, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).Should(BeNil())
			feeEstimator := btc.NewFixFeeEstimator(1e3)
			wal1, err := btc.NewWallet(network, waddrmgr.PubKeyHash, key1, indexer, btcClient, feeEstimator)
			Expect(err).Should(BeNil())
			wal2, err := btc.NewWallet(network, waddrmgr.PubKeyHash, key2, indexer, btcClient, feeEstimator)
			Expect(err).Should(BeNil())

			By("Create a new cosigner wallet")
			wallet, err := cosigner.NewWallet(network, key1, cosignerPub, indexer, client, btcClient, feeEstimator)
			Expect(err).Should(BeNil())

			By("fund the wallets")
			_, err = btctest.Faucet(wallet.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			_, err = btctest.Faucet(wal1.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			fundTx, err := btctest.Faucet(wal2.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(fundTx.String()))).Should(Succeed())

			By("Create a few htlc redeem actions")
			actions, err := btctest.PrepareActions(ctx, 5, wal1, wal2, indexer, btc.HtlcActionInitiate)
			Expect(err).Should(BeNil())
			for i := 0; i < 5; i++ {
				// Wal1 inits the htlc
				time.Sleep(5 * time.Second)
				initTx, _, err := wallet.Execute(ctx, []btc.HtlcAction{actions[i]})
				Expect(err).Should(BeNil())
				log.Printf("tx = %v", initTx.TxHash().String())

				// Wal2 redeems the htlc
				time.Sleep(time.Second)
				redeemTx, err := wal2.Redeem(ctx, actions[i].Htlc)
				Expect(err).Should(BeNil())
				Expect(client.NewMergeTx(redeemTx)).Should(Succeed())
				log.Printf("tx = %v", redeemTx.TxHash().String())
			}
		})

		It("should be able to redeem a few htlcs", func(ctx context.Context) {
			By("Create new wallets")
			key1, err := btcec.NewPrivateKey()
			Expect(err).Should(BeNil())
			key2, _, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).Should(BeNil())
			feeEstimator := btc.NewFixFeeEstimator(1e3)
			wal1, err := btc.NewWallet(network, waddrmgr.PubKeyHash, key1, indexer, btcClient, feeEstimator)
			Expect(err).Should(BeNil())
			wal2, err := btc.NewWallet(network, waddrmgr.PubKeyHash, key2, indexer, btcClient, feeEstimator)
			Expect(err).Should(BeNil())

			By("Create a new cosigner wallet")
			wallet, err := cosigner.NewWallet(network, key1, cosignerPub, indexer, client, btcClient, feeEstimator)
			Expect(err).Should(BeNil())

			By("fund the wallets")
			_, err = btctest.Faucet(wallet.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			fundTx, err := btctest.Faucet(wal2.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(fundTx.String()))).Should(Succeed())

			By("Create a few htlc redeem actions")
			actions, err := btctest.PrepareActions(ctx, 5, wal2, wal1, indexer, btc.HtlcActionRedeem)
			Expect(err).Should(BeNil())
			for i := 0; i < 5; i++ {
				time.Sleep(5 * time.Second)
				_, tx, err := wallet.Execute(ctx, []btc.HtlcAction{actions[i]})
				Expect(err).Should(BeNil())
				log.Printf("tx = %v", tx.TxHash().String())
			}
		})

		It("should be able to refund a few htlcs", func(ctx context.Context) {
			By("Create new wallets")
			key1, err := btcec.NewPrivateKey()
			Expect(err).Should(BeNil())
			key2, _, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).Should(BeNil())
			feeEstimator := btc.NewFixFeeEstimator(1e3)
			wal1, err := btc.NewWallet(network, waddrmgr.PubKeyHash, key1, indexer, btcClient, feeEstimator)
			Expect(err).Should(BeNil())
			wal2, err := btc.NewWallet(network, waddrmgr.PubKeyHash, key2, indexer, btcClient, feeEstimator)
			Expect(err).Should(BeNil())

			By("Create a new cosigner wallet")
			wallet, err := cosigner.NewWallet(network, key1, cosignerPub, indexer, client, btcClient, feeEstimator)
			Expect(err).Should(BeNil())

			By("fund the wallets")
			_, err = btctest.Faucet(wallet.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			_, err = btctest.Faucet(wal1.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			fundTx, err := btctest.Faucet(wal2.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(fundTx.String()))).Should(Succeed())

			By("Create a few htlc redeem actions")
			actions, err := btctest.PrepareActions(ctx, 5, wal1, wal2, indexer, btc.HtlcActionRefund)
			Expect(err).Should(BeNil())
			for i := 0; i < 5; i++ {
				time.Sleep(5 * time.Second)
				_, tx, err := wallet.Execute(ctx, []btc.HtlcAction{actions[i]})
				Expect(err).Should(BeNil())
				log.Printf("tx = %v", tx.TxHash().String())
			}
		})

		It("should be able to do a mix of actions ", func(ctx context.Context) {
			By("Create new wallets")
			key1, err := btcec.NewPrivateKey()
			Expect(err).Should(BeNil())
			key2, _, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).Should(BeNil())
			feeEstimator := btc.NewFixFeeEstimator(1e3)
			wal1, err := btc.NewWallet(network, waddrmgr.PubKeyHash, key1, indexer, btcClient, feeEstimator)
			Expect(err).Should(BeNil())
			wal2, err := btc.NewWallet(network, waddrmgr.PubKeyHash, key2, indexer, btcClient, feeEstimator)
			Expect(err).Should(BeNil())

			By("Create a new cosigner wallet")
			wallet, err := cosigner.NewWallet(network, key1, cosignerPub, indexer, client, btcClient, feeEstimator)
			Expect(err).Should(BeNil())

			By("fund the wallets")
			_, err = btctest.Faucet(wallet.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			_, err = btctest.Faucet(wal1.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			fundTx, err := btctest.Faucet(wal2.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(fundTx.String()))).Should(Succeed())

			By("Prepare some init/redeem/refund/instantRefund")
			number := 3
			inits, err := btctest.PrepareActions(ctx, number, wal1, wal2, indexer, btc.HtlcActionInitiate)
			Expect(err).Should(BeNil())
			redeems, err := btctest.PrepareActions(ctx, number, wal2, wal1, indexer, btc.HtlcActionRedeem)
			Expect(err).Should(BeNil())
			refunds, err := btctest.PrepareActions(ctx, number, wal1, wal2, indexer, btc.HtlcActionRefund)
			Expect(err).Should(BeNil())
			// instantRefunds, err := btctest.PrepareActions(ctx, number, wal2, wal1, indexer, btc.HtlcActionInstantRefund)
			// Expect(err).Should(BeNil())

			By("Combined all actions and shuffle the order")
			actions := append(inits, append(redeems, refunds...)...)
			rand.Shuffle(len(actions), func(i, j int) {
				actions[i], actions[j] = actions[j], actions[i]
			})

			By("Execute all actions one by one using rbf")
			for i := 0; i < len(actions); i++ {
				tx1, tx2, err := wallet.Execute(ctx, []btc.HtlcAction{actions[i]})
				Expect(err).Should(BeNil())
				Expect(tx1 != nil || tx2 != nil).Should(BeTrue())
				time.Sleep(1 * time.Second)
			}
		})
	})
})
