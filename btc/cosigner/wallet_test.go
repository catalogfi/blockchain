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
		FIt("Make multiple spends in multiple batch ", func(ctx context.Context) {
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
				_, err = wallet.Send(ctx, recipients)
				Expect(err).Should(BeNil())
			}
		})
	})
})
