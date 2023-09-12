package btc_test

import (
	"context"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/catalogfi/multichain/btc"
	"github.com/catalogfi/multichain/testutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Indexer client", func() {

	Context("When using electrs API", func() {
		It("should be able to fetch the utxos of an address ", func() {
			By("Initialise the electrs client")
			network := &chaincfg.RegressionNetParams
			client := RegtestIndexer()

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
				if utxo.TxID == txid {
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
			tx, err := client.GetTx(context.Background(), txid)
			Expect(err).To(BeNil())
			Expect(tx.TxID).Should(Equal(txid))
			Expect(tx.Status.Confirmed).Should(BeTrue())

			By("SubmitTx()")
			total := int64(0)
			for _, utxo := range utxos {
				total += utxo.Amount
			}
			fee := int64(1000)
			recipients := []btc.Recipient{
				{
					To:     addr.String(),
					Amount: total - fee,
				},
			}
			rawTx, _, err := btc.BuildTx(network, utxos, recipients, fee, addr)
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
				if tx.TxID == txid {
					has = true
				}
			}
			Expect(has).Should(BeTrue())
		})
	})
})
