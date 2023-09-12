package btc_test

import (
	"encoding/hex"
	"reflect"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/catalogfi/multichain/btc"
	"github.com/catalogfi/multichain/testutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("bitcoin client", func() {
	Context("Client initialization", func() {
		It("should return an error if providing an unknown chain params", func() {
			config := &rpcclient.ConnConfig{
				Params:       "",
				Host:         "0.0.0.0:18443",
				HTTPPostMode: true,
				DisableTLS:   true,
			}
			_, err := btc.NewClient(config)
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("regression local testnet", func() {
		Context("when using electrs indexer", func() {
			It("should be able to get all the data without any error", func() {
				By("Initialise a local regnet client")
				network := &chaincfg.RegressionNetParams
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

				By("Create a new tx")
				addr, err := testutil.RandomBtcAddressP2PKH(network)
				Expect(err).To(BeNil())
				txid, err := testutil.NigiriFaucet(addr.EncodeAddress())
				Expect(err).To(BeNil())
				time.Sleep(500 * time.Millisecond)

				By("GetRawTransaction()")
				txhash, err := chainhash.NewHashFromStr(txid)
				Expect(err).To(BeNil())
				rawTx, err := client.GetRawTransaction(txhash.CloneBytes())
				Expect(err).To(BeNil())
				Expect(rawTx.Txid).Should(Equal(txid))

				By("GetBlockByHash()")
				blockByHash, err := client.GetBlockByHash(rawTx.BlockHash)
				Expect(err).To(BeNil())
				Expect(blockByHash.Hash).Should(Equal(rawTx.BlockHash))

				By("GetBlockByHeight()")
				blockByHeight, err := client.GetBlockByHeight(blockByHash.Height)
				Expect(err).To(BeNil())
				Expect(*blockByHeight).Should(Equal(*blockByHash))

				By("GetTxOut()")
				vout := 0
				for index, out := range rawTx.Vout {
					if out.Value == 1 {
						vout = index
						break
					}
				}
				txout, err := client.GetTxOut(txid, uint32(vout))
				Expect(err).To(BeNil())
				Expect(txout.Value).Should(Equal(float64(1)))
				scriptPubkey, err := txscript.PayToAddrScript(addr)
				Expect(err).To(BeNil())
				Expect(txout.ScriptPubKey.Hex).Should(Equal(hex.EncodeToString(scriptPubkey)))
			})
		})
	})
})
