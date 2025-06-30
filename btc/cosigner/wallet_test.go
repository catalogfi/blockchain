package cosigner_test

import (
	"context"
	"log"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/catalogfi/blockchain/btc/cosigner"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("wallet", func() {
	Context("spend", func() {
		It("Make multiple spends in multiple batch ", func(ctx context.Context) {
			By("Create a new wallet")
			key, err := btcec.NewPrivateKey()
			Expect(err).Should(BeNil())
			feeEstimator := btc.NewFixFeeEstimator(1e3)
			wallet, err := cosigner.NewWallet(network, key, cosignerPub, indexer, client, btcClient, feeEstimator)
			Expect(err).Should(BeNil())
			log.Print("wal addr = ", wallet.Address().String())

			By("Fund the wallet")
			fundTx, err := btctest.Faucet(wallet.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(fundTx.String()))).Should(Succeed())

			for i := 0; i < 7; i++ {
				By("Generate a new receiver")
				time.Sleep(5 * time.Second)
				_, to, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
				Expect(err).Should(BeNil())
				recipients := []btc.Recipient{
					btc.NewRecipient(to.EncodeAddress(), 1e7),
				}
				tx, err := wallet.Send(ctx, recipients)
				Expect(err).Should(BeNil())
				log.Printf("txid = %v", tx.TxHash().String())
			}
		})

		It("Merge tx", func(ctx context.Context) {
			By("Create a new wallet")
			key, err := btcec.NewPrivateKey()
			Expect(err).Should(BeNil())
			feeEstimator := btc.NewFixFeeEstimator(1e3)
			wallet, err := cosigner.NewWallet(network, key, cosignerPub, indexer, client, btcClient, feeEstimator)
			Expect(err).Should(BeNil())
			log.Print("wal addr = ", wallet.Address().String())

			By("Fund the wallet")
			fundTx, err := btctest.Faucet(wallet.Address().EncodeAddress())
			Expect(err).Should(BeNil())
			Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(fundTx.String()))).Should(Succeed())

			By("Start a new tx ")
			time.Sleep(5 * time.Second)
			toKey, to, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).Should(BeNil())
			recipients := []btc.Recipient{
				btc.NewRecipient(to.EncodeAddress(), 1e7),
			}
			tx, err := wallet.Send(ctx, recipients)
			Expect(err).Should(BeNil())
			log.Printf("txid = %v", tx.TxHash().String())

			By("Create a merge tx")
			time.Sleep(5 * time.Second)
			_, to1, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).Should(BeNil())
			utxos, err := indexer.GetUTXOs(ctx, to)
			Expect(err).Should(BeNil())
			sizer := btc.NewSizeEstimatorOfAddrType(waddrmgr.PubKeyHash, utxos...)
			feeMode := btc.MinFeeRateMode(1e3, sizer)
			mergeTx, err := btc.BuildTx(network, feeMode, utxos, nil, nil, to1)
			Expect(err).Should(BeNil())
			Expect(btc.SignTx(waddrmgr.PubKeyHash, mergeTx, toKey, utxos)).Should(Succeed())
			Expect(client.NewMergeTx(mergeTx)).Should(Succeed())

			By("Submit a new rbf to include the merge tx")
			signedTx, err := wallet.Merge(ctx)
			Expect(err).Should(BeNil())
			log.Printf("txid = %v", signedTx.TxHash().String())
		})
	})
})
