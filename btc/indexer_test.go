package btc_test

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/fatih/color"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Indexer client", func() {
	Context("When using electrs API", func() {
		It("should be able to fetch the utxos of an address ", func(ctx context.Context) {
			By("GetUTXOs()")
			key, addr, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())
			txid, err := btctest.Faucet(addr.EncodeAddress())
			Expect(err).To(BeNil())
			Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(txid.String()))).Should(Succeed())

			utxos, err := indexer.GetUTXOs(ctx, addr)
			Expect(err).To(BeNil())
			Expect(len(utxos)).Should(BeNumerically(">=", 1))
			exist := false
			for _, utxo := range utxos {
				if utxo.TxID == txid.String() {
					exist = true
					break
				}
			}
			Expect(exist).Should(BeTrue())

			By("GetTipBlockHeight()")
			tip, err := indexer.GetTipBlockHeight(ctx)
			Expect(err).To(BeNil())
			Expect(tip).Should(BeNumerically(">=", 100))

			By("GetTx()")
			tx, err := indexer.GetTx(ctx, txid.String())
			Expect(err).To(BeNil())
			Expect(tx.TxID).Should(Equal(txid.String()))
			Expect(tx.Status.Confirmed).Should(BeTrue())

			By("GetTxHex()")
			txHex, err := indexer.GetTxHex(ctx, txid.String())
			Expect(err).To(BeNil())
			Expect(txHex).ShouldNot(BeEmpty())
			txBytes, err := hex.DecodeString(txHex)
			Expect(err).To(BeNil())
			Expect(txBytes).ShouldNot(BeEmpty())
			btcTx, err := btcutil.NewTxFromBytes(txBytes)
			Expect(err).To(BeNil())
			Expect(btcTx.MsgTx().TxHash().String()).Should(Equal(txid.String()))

			By("SubmitTx()")
			amount, feeRate := int64(1e6), btctest.RandomFeeRate()
			recipient, err := btc.NewTxOutFromAddress(addr, amount)
			Expect(err).To(BeNil())
			recipients := []*wire.TxOut{recipient}
			sizer := btc.NewSizeEstimator(btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH, utxos...)
			feeMode := btc.MinFeeRateMode(feeRate, sizer)
			rawTx, err := btc.BuildTx(feeMode, nil, utxos, recipients, addr)
			Expect(err).To(BeNil())
			Expect(btc.SignTx(waddrmgr.PubKeyHash, rawTx, key, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, rawTx)).Should(Succeed())

			By("GetAddressTxs()")
			txs, err := indexer.GetAddressTxs(ctx, addr, "")
			Expect(err).To(BeNil())
			has := false
			for _, tx := range txs {
				if tx.TxID == txid.String() {
					has = true
				}
			}
			Expect(has).Should(BeTrue())
		})
	})

	Context("errors", func() {
		It("should return specific errors", func(ctx context.Context) {
			By("New address")
			key, pkAddr, err := btctest.NewBtcAddrWithFunds(network, waddrmgr.PubKeyHash, indexer)
			Expect(err).To(BeNil())

			By("Construct a new tx")
			utxos, err := indexer.GetUTXOs(ctx, pkAddr)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e6), btctest.RandomFeeRate()
			recipient, err := btc.NewTxOutFromAddress(pkAddr, amount)
			Expect(err).To(BeNil())
			recipients := []*wire.TxOut{recipient}
			sizer := btc.NewSizeEstimator(btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH, utxos...)
			feeMode := btc.MinFeeRateMode(feeRate, sizer)
			transaction, err := btc.BuildTx(feeMode, nil, utxos, recipients, pkAddr)
			Expect(err).To(BeNil())
			Expect(btc.SignTx(waddrmgr.PubKeyHash, transaction, key, utxos)).Should(Succeed())

			By("Submit the transaction")
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
			By(fmt.Sprintf("Funding tx hash = %v", color.YellowString(transaction.TxHash().String())))
			time.Sleep(time.Second)

			By("Expect a `ErrAlreadyInChain` error if the tx is already in a block")
			Expect(btctest.NewBlockWaitMined(1, indexer)).Should(Succeed())
			err = indexer.SubmitTx(ctx, transaction)
			Expect(errors.Is(err, btc.ErrAlreadyInUtxoSet)).Should(BeTrue())

			By("Try construct a new transaction spending the same input")
			recipient1, err := btc.NewTxOutFromAddress(pkAddr, 2*amount)
			Expect(err).To(BeNil())
			recipients1 := []*wire.TxOut{recipient1}
			transaction1, err := btc.BuildTx(feeMode, nil, utxos, recipients1, pkAddr)
			Expect(err).To(BeNil())
			Expect(btc.SignTx(waddrmgr.PubKeyHash, transaction, key, utxos)).Should(Succeed())
			By("Expect a `ErrAlreadyInChain` error if the tx is already in a block")
			err = indexer.SubmitTx(ctx, transaction1)
			Expect(errors.Is(err, btc.ErrTxInputsMissingOrSpent)).Should(BeTrue())
		})
	})
})
