package btc_test

import (
	"context"
	"errors"
	"time"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/fatih/color"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bitcoin", func() {
	Context("RBF", func() {
		It("should be able to build transaction which support rbf", func(ctx context.Context) {
			By("Initialize keys and addresses")
			privKey1, p2pkhAddr1, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())
			_, p2pkhAddr2, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())

			By("Funding the addresses")
			_, err = btctest.NigiriFaucet(p2pkhAddr1.EncodeAddress())
			Expect(err).To(BeNil())
			time.Sleep(5 * time.Second)

			By("Construct a RBF tx which sends money from p2pkhAddr1 to p2pkhAddr2")
			utxos, err := indexer.GetUTXOs(ctx, p2pkhAddr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e7), 4
			recipients := []btc.Recipient{
				{
					To:     p2pkhAddr2.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildRbfTransaction(network, feeRate, btc.NewRawInputs(), utxos, btc.P2pkhUpdater, recipients, p2pkhAddr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.SignP2pkhTx(network, privKey1, transaction)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
			By(color.GreenString("RBF tx hash = %v", transaction.TxHash().String()))

			By("Construct a tx with higher fees")
			transaction1, err := btc.BuildTransaction(network, feeRate+5, btc.NewRawInputs(), utxos, btc.P2pkhUpdater, recipients, p2pkhAddr1)
			Expect(err).To(BeNil())
			// transaction1.TxOut[1].Value++

			By("Sign and submit replacement tx")
			Expect(btc.SignP2pkhTx(network, privKey1, transaction1)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction1)).Should(Succeed())
			color.Green("Replacement tx = %v", transaction1.TxHash().String())
		})

		It("should get an error when trying to replace a mined tx", func(ctx context.Context) {
			By("Initialize keys and addresses")
			privKey, pkAddr, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())
			_, toAddr, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())

			By("Funding the addresses")
			_, err = btctest.NigiriFaucet(pkAddr.EncodeAddress())
			Expect(err).To(BeNil())
			time.Sleep(5 * time.Second)

			By("Construct a new tx")
			utxos, err := indexer.GetUTXOs(ctx, pkAddr)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e5), 5
			recipients := []btc.Recipient{
				{
					To:     toAddr.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildRbfTransaction(network, feeRate, btc.NewRawInputs(), utxos, btc.P2pkhUpdater, recipients, pkAddr)
			Expect(err).To(BeNil())

			By("Sign the transaction inputs")
			for i := range transaction.TxIn {
				pkScript, err := txscript.PayToAddrScript(pkAddr)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(transaction, i, pkScript, txscript.SigHashAll, privKey, true)
				Expect(err).To(BeNil())
				transaction.TxIn[i].SignatureScript = sigScript
			}

			By("Submit the transaction")
			Expect(btc.SignP2pkhTx(network, privKey, transaction)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
			By(color.GreenString("RBF tx hash = %v", transaction.TxHash().String()))
			Expect(btctest.NigiriNewBlock()).Should(Succeed())

			By("Build a new tx with higher fee")
			feeRate += 2
			transaction, err = btc.BuildRbfTransaction(network, feeRate, btc.NewRawInputs(), utxos, btc.P2pkhUpdater, recipients, pkAddr)
			Expect(err).To(BeNil())

			By("Sign the replacement tx and the submission should be rejected")
			Expect(btc.SignP2pkhTx(network, privKey, transaction)).Should(Succeed())
			err = indexer.SubmitTx(ctx, transaction)
			Expect(errors.Is(err, btc.ErrTxInputsMissingOrSpent)).Should(BeTrue())
		})
	})
})
