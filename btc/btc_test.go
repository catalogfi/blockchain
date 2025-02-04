package btc_test

import (
	"context"
	"errors"
	"sort"

	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/fatih/color"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bitcoin", func() {
	Context("Build a transaction", func() {
		It("should be able to build a transaction", func(ctx context.Context) {
			By("Initialize keys and addresses")
			addrType := waddrmgr.PubKeyHash
			privKey1, p2pkhAddr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).To(BeNil())
			_, p2pkhAddr2, err := btctest.NewBtcKey(network, addrType)
			Expect(err).To(BeNil())

			By("Construct a transaction which sends money from p2pkhAddr1 to p2pkhAddr2")
			utxos, err := indexer.GetUTXOs(ctx, p2pkhAddr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e7), 4
			recipients := btc.SingleRecipient(p2pkhAddr2.EncodeAddress(), amount)
			sizer := btc.NewSizeEstimator(utxos, btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH)
			transaction, err := btc.BuildTransaction(network, feeRate, nil, utxos, sizer, recipients, p2pkhAddr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.SignTx(addrType, transaction, privKey1, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
			By(color.GreenString("tx hash = %v", transaction.TxHash().String()))
		})

		It("should not use a utxo if its uneconomical", func(ctx context.Context) {
			By("Initialize keys and addresses")
			addrType := waddrmgr.TaprootPubKey
			privKey1, p2trAddr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).To(BeNil())
			_, p2trAddr2, err := btctest.NewBtcKey(network, addrType)
			Expect(err).To(BeNil())

			By("Create a utxo which is less than the dust amount")
			utxos, err := indexer.GetUTXOs(ctx, p2trAddr1)
			Expect(err).To(BeNil())
			recipients := btc.SingleRecipient(p2trAddr1.EncodeAddress(), btc.DustAmount)
			sizer := btc.NewSizeEstimator(utxos, btc.BaseSizeP2TR, btc.SegwitSizeP2TR)
			transaction, err := btc.BuildTransaction(network, 3, nil, utxos, sizer, recipients, p2trAddr1)
			Expect(err).To(BeNil())
			Expect(btc.SignTx(addrType, transaction, privKey1, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
			Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

			By("Construct a transaction try using the dust amount")
			utxos1, err := indexer.GetUTXOs(ctx, p2trAddr1)
			Expect(err).To(BeNil())
			sort.Slice(utxos1, func(i, j int) bool {
				return utxos1[i].Amount < utxos1[j].Amount
			})
			Expect(utxos1[0].TxID).Should(Equal(transaction.TxHash().String()))
			amount, feeRate := int64(1e7), 3
			recipients1 := btc.SingleRecipient(p2trAddr2.EncodeAddress(), amount)
			sizer1 := btc.NewSizeEstimator(utxos1, btc.BaseSizeP2TR, btc.SegwitSizeP2TR)
			transaction1, err := btc.BuildTransaction(network, feeRate, nil, utxos1, sizer1, recipients1, p2trAddr1)
			Expect(err).To(BeNil())

			By("The new tx should not use the uneconomical utxo")
			for _, input := range transaction1.TxIn {
				if input.PreviousOutPoint.String() == utxos1[0].String() {
					Fail("The new tx should not use the utxo which is less than the dust amount")
				}
			}
		})
	})

	Context("RBF", func() {
		It("should be able to build and replace an RBF transaction", func(ctx context.Context) {
			By("Initialize keys and addresses")
			addrType := waddrmgr.PubKeyHash
			privKey1, p2pkhAddr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).To(BeNil())
			_, p2pkhAddr2, err := btctest.NewBtcKey(network, addrType)
			Expect(err).To(BeNil())

			By("Construct a RBF tx which sends money from p2pkhAddr1 to p2pkhAddr2")
			utxos, err := indexer.GetUTXOs(ctx, p2pkhAddr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e7), 4
			recipients := btc.SingleRecipient(p2pkhAddr2.EncodeAddress(), amount)
			sizer := btc.NewSizeEstimator(utxos, btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH)
			transaction, err := btc.BuildRbfTransaction(network, feeRate, nil, utxos, sizer, recipients, p2pkhAddr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.SignTx(addrType, transaction, privKey1, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
			By(color.GreenString("RBF tx hash = %v", transaction.TxHash().String()))

			By("Build a replacement tx with higher fee")
			feeRate += 2
			replaceTx, err := btc.BuildRbfTransaction(network, feeRate, nil, utxos, sizer, recipients, p2pkhAddr1)
			Expect(err).To(BeNil())

			By("Sign and submit the replacement tx")
			Expect(btc.SignTx(addrType, replaceTx, privKey1, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, replaceTx)).Should(Succeed())
			By(color.GreenString("Replaced RBF tx hash = %v", replaceTx.TxHash().String()))
		})

		It("should get an error when trying to replace a mined tx", func(ctx context.Context) {
			By("Initialize keys and addresses")
			addrType := waddrmgr.PubKeyHash
			privKey, pkAddr, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).To(BeNil())
			_, toAddr, err := btctest.NewBtcKey(network, addrType)
			Expect(err).To(BeNil())

			By("Construct a RBF tx which sends money from pkAddr to toAddr")
			utxos, err := indexer.GetUTXOs(ctx, pkAddr)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e7), 5
			recipients := btc.SingleRecipient(toAddr.EncodeAddress(), amount)
			sizer := btc.NewSizeEstimator(utxos, btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH)
			transaction, err := btc.BuildRbfTransaction(network, feeRate, nil, utxos, sizer, recipients, pkAddr)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.SignTx(addrType, transaction, privKey, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
			By(color.GreenString("RBF tx hash = %v", transaction.TxHash().String()))

			By("Mine a new block")
			Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

			By("Build a replacement tx with higher fee")
			feeRate += 2
			replaceTx, err := btc.BuildRbfTransaction(network, feeRate, nil, utxos, sizer, recipients, pkAddr)
			Expect(err).To(BeNil())

			By("Sign and attempt to submit the replacement tx")
			Expect(btc.SignTx(addrType, replaceTx, privKey, utxos)).Should(Succeed())
			err = indexer.SubmitTx(ctx, replaceTx)
			Expect(errors.Is(err, btc.ErrTxInputsMissingOrSpent)).Should(BeTrue())
		})
	})
})
