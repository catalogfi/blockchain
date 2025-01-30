package btc_test

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/fatih/color"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Indexer client", func() {
	Context("When using electrs API", func(ctx context.Context) {
		It("should be able to fetch the utxos of an address ", func() {
			By("GetUTXOs()")
			key, addr, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())
			txid, err := btctest.Faucet(addr.EncodeAddress())
			Expect(err).To(BeNil())
			Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(txid.String()))).Should(Succeed())

			utxos, err := indexer.GetUTXOs(context.Background(), addr)
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
			tip, err := indexer.GetTipBlockHeight(context.Background())
			Expect(err).To(BeNil())
			Expect(tip).Should(BeNumerically(">=", 100))

			By("GetTx()")
			tx, err := indexer.GetTx(context.Background(), txid.String())
			Expect(err).To(BeNil())
			Expect(tx.TxID).Should(Equal(txid.String()))
			Expect(tx.Status.Confirmed).Should(BeTrue())

			By("GetTxHex()")
			txHex, err := indexer.GetTxHex(context.Background(), txid.String())
			Expect(err).To(BeNil())
			Expect(txHex).ShouldNot(BeEmpty())
			txBytes, err := hex.DecodeString(txHex)
			Expect(err).To(BeNil())
			Expect(txBytes).ShouldNot(BeEmpty())
			btcTx, err := btcutil.NewTxFromBytes(txBytes)
			Expect(err).To(BeNil())
			Expect(btcTx.MsgTx().TxHash().String()).Should(Equal(txid.String()))

			By("SubmitTx()")
			amount, feeRate := int64(1e6), 10
			recipients := []btc.Recipient{
				{
					To:     addr.String(),
					Amount: amount,
				},
			}
			sizer := btc.NewSizeEstimator(utxos, btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH)
			rawTx, err := btc.BuildTransaction(network, feeRate, nil, utxos, sizer, recipients, addr)
			Expect(err).To(BeNil())
			Expect(btc.SignP2pkhTx(network, key, rawTx)).Should(Succeed())
			Expect(client.SubmitTx(context.Background(), rawTx)).Should(Succeed())

			By("GetAddressTxs()")
			txs, err := indexer.GetAddressTxs(context.Background(), addr, "")
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
		It("should return specific errors", func() {
			By("New address")
			key, pkAddr, err := btctest.NewBtcAddrWithFunds(network, waddrmgr.PubKeyHash, indexer)
			Expect(err).To(BeNil())

			By("Construct a new tx")
			utxos, err := indexer.GetUTXOs(context.Background(), pkAddr)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e6), 10
			recipients := []btc.Recipient{
				{
					To:     pkAddr.EncodeAddress(),
					Amount: amount,
				},
			}
			sizer := btc.NewSizeEstimator(utxos, btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH)
			transaction, err := btc.BuildTransaction(network, feeRate, nil, utxos, sizer, recipients, pkAddr)
			Expect(err).To(BeNil())
			Expect(btc.SignP2pkhTx(network, key, transaction)).Should(Succeed())

			By("Submit the transaction")
			Expect(indexer.SubmitTx(context.Background(), transaction)).Should(Succeed())
			By(fmt.Sprintf("Funding tx hash = %v", color.YellowString(transaction.TxHash().String())))
			time.Sleep(time.Second)

			By("Expect a `ErrAlreadyInChain` error if the tx is already in a block")
			Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())
			err = indexer.SubmitTx(context.Background(), transaction)
			Expect(errors.Is(err, btc.ErrAlreadyInChain)).Should(BeTrue())

			By("Try construct a new transaction spending the same input")
			recipients1 := []btc.Recipient{
				{
					To:     pkAddr.EncodeAddress(),
					Amount: 2 * amount,
				},
			}
			transaction1, err := btc.BuildTransaction(network, feeRate, nil, utxos, sizer, recipients1, pkAddr)
			Expect(err).To(BeNil())
			Expect(btc.SignP2pkhTx(network, key, transaction)).Should(Succeed())
			By("Expect a `ErrAlreadyInChain` error if the tx is already in a block")
			err = indexer.SubmitTx(context.Background(), transaction1)
			Expect(errors.Is(err, btc.ErrTxInputsMissingOrSpent)).Should(BeTrue())
		})
	})
})
