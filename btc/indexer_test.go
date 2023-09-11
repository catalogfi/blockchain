package btc_test

import (
	"context"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/multichain/btc"
	"github.com/catalogfi/multichain/testutil"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Indexer client", func() {

	Context("When using electrs API", func() {
		It("should be able to fetch the utxos of an address ", func() {
			By("Initialise the quickNode client")
			logger, err := zap.NewDevelopment()
			Expect(err).To(BeNil())
			url := testutil.ParseStringEnv("BTC_INDEXER_ELECTRS_REGNET", "")
			client := btc.NewElectrsIndexerClient(logger, url, btc.DefaultRetryInterval)

			By("GetUTXOs()")
			txid, err := testutil.NigiriFaucet("mynaqNW3Pa2t3WC3jb9hr3hcFQYeVmyQNb")
			Expect(err).To(BeNil())
			time.Sleep(5 * time.Second)
			addr, err := btcutil.DecodeAddress("mynaqNW3Pa2t3WC3jb9hr3hcFQYeVmyQNb", &chaincfg.RegressionNetParams)
			Expect(err).To(BeNil())
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
		})
	})
})
