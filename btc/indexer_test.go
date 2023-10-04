package btc_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/catalogfi/blockchain/testutil"
	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Indexer client", func() {

	Context("When using electrs API", func() {
		It("should be able to fetch the utxos of an address ", func() {
			By("Initialise the electrs client")
			network := &chaincfg.RegressionNetParams
			client := btctest.RegtestIndexer()

			By("GetUTXOs()")
			key, err := btcec.NewPrivateKey()
			Expect(err).To(BeNil())
			addr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(key.PubKey().SerializeCompressed()), network)
			Expect(err).To(BeNil())
			txid, err := testutil.NigiriFaucet(addr.EncodeAddress())
			Expect(err).To(BeNil())
			time.Sleep(5 * time.Second)
			utxos, err := client.GetUTXOs(context.Background(), addr)
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
			tip, err := client.GetTipBlockHeight(context.Background())
			Expect(err).To(BeNil())
			Expect(tip).Should(BeNumerically(">=", 100))

			By("GetTx()")
			tx, err := client.GetTx(context.Background(), txid.String())
			Expect(err).To(BeNil())
			Expect(tx.TxID).Should(Equal(txid.String()))
			Expect(tx.Status.Confirmed).Should(BeTrue())

			By("SubmitTx()")
			amount, feeRate := int64(1e6), 10
			recipients := []btc.Recipient{
				{
					To:     addr.String(),
					Amount: amount,
				},
			}
			rawTx, err := btc.BuildTransaction(feeRate, network, nil, utxos, recipients, 0, 0, addr)
			Expect(err).To(BeNil())
			for i := range rawTx.TxIn {
				pkScript, err := txscript.PayToAddrScript(addr)
				Expect(err).To(BeNil())
				sigScript, err := txscript.SignatureScript(rawTx, i, pkScript, txscript.SigHashAll, key, true)
				Expect(err).To(BeNil())
				rawTx.TxIn[i].SignatureScript = sigScript
			}
			Expect(client.SubmitTx(context.Background(), rawTx)).Should(Succeed())

			By("FeeEstimate()")
			By("    --local env should not have enough data for the estimate")
			_, err = client.FeeEstimate(context.Background())
			Expect(err.Error()).Should(ContainSubstring("not enough data"))

			By("GetAddressTxs()")
			txs, err := client.GetAddressTxs(context.Background(), addr)
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
			By("Initialise a local regnet client")
			network := &chaincfg.RegressionNetParams
			indexer := btctest.RegtestIndexer()

			By("New address")
			privKey, err := btcec.NewPrivateKey()
			Expect(err).To(BeNil())
			pubKey := privKey.PubKey()
			pkAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), network)
			Expect(err).To(BeNil())

			By("funding the addresses")
			txhash, err := testutil.NigiriFaucet(pkAddr.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address1 %v , txid = %v", pkAddr.EncodeAddress(), txhash))
			time.Sleep(5 * time.Second)

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
			transaction, err := btc.BuildTransaction(feeRate, network, nil, utxos, recipients, 0, 0, pkAddr)
			Expect(err).To(BeNil())
			for i := range transaction.TxIn {
				pkScript, err := txscript.PayToAddrScript(pkAddr)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(transaction, i, pkScript, txscript.SigHashAll, privKey, true)
				Expect(err).To(BeNil())
				transaction.TxIn[i].SignatureScript = sigScript
			}

			By("Submit the transaction")
			Expect(indexer.SubmitTx(context.Background(), transaction)).Should(Succeed())
			By(fmt.Sprintf("Funding tx hash = %v", color.YellowString(transaction.TxHash().String())))
			time.Sleep(time.Second)

			By("Expect a `ErrAlreadyInChain` error if the tx is already in a block")
			Expect(testutil.NigiriNewBlock()).Should(Succeed())
			time.Sleep(1 * time.Second)
			err = indexer.SubmitTx(context.Background(), transaction)
			Expect(errors.Is(err, btc.ErrAlreadyInChain)).Should(BeTrue())

			By("Try construct a new transaction spending the same input")
			recipients1 := []btc.Recipient{
				{
					To:     pkAddr.EncodeAddress(),
					Amount: 2 * amount,
				},
			}
			transaction1, err := btc.BuildTransaction(feeRate, network, nil, utxos, recipients1, 0, 0, pkAddr)
			Expect(err).To(BeNil())
			for i := range transaction1.TxIn {
				pkScript, err := txscript.PayToAddrScript(pkAddr)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(transaction1, i, pkScript, txscript.SigHashAll, privKey, true)
				Expect(err).To(BeNil())
				transaction1.TxIn[i].SignatureScript = sigScript
			}
			By("Expect a `ErrAlreadyInChain` error if the tx is already in a block")
			err = indexer.SubmitTx(context.Background(), transaction1)
			Expect(errors.Is(err, btc.ErrTxInputsMissingOrSpent)).Should(BeTrue())
		})
	})
})
