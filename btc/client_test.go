package btc_test

import (
	"reflect"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("bitcoin client", func() {
	Context("regression local testnet", func() {
		Context("when using electrs indexer", func() {
			It("should be able to get all the data without any error", func() {
				By("Initialise a local regnet client")
				client, err := RegtestClient()
				Expect(err).To(BeNil())

				By("Net()")
				Expect(reflect.DeepEqual(client.Net(), &chaincfg.RegressionNetParams)).Should(BeTrue())

				By("LatestBlock()")
				height, hash, err := client.LatestBlock()
				Expect(err).To(BeNil())
				Expect(height).Should(BeNumerically(">=", 100))
				Expect(len(hash)).Should(Equal(64))
				Expect(hash).ShouldNot(Equal("0000000000000000000000000000000000000000000000000000000000000000"))

				By("GetBlockByHeight()")
				block, err := client.GetBlockByHeight(1)
				Expect(err).To(BeNil())
				Expect(block.Hash).ShouldNot(Equal("0000000000000000000000000000000000000000000000000000000000000000"))
				Expect(len(block.Tx)).Should(Equal(1))

				By("GetRawTransaction()")
				txHash, err := chainhash.NewHashFromStr(block.Tx[0])
				Expect(err).To(BeNil())
				rawTx, err := client.GetRawTransaction(txHash.CloneBytes())
				Expect(err).To(BeNil())
				Expect(len(rawTx.Hex)).Should(BeNumerically(">", 1))
			})
		})
	})
})
