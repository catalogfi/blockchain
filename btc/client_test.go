package btc_test

import (
	"context"
	"encoding/hex"
	"os"
	"reflect"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/catalogfi/multichain/btc"
	"github.com/catalogfi/multichain/testutil"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("bitcoin client", func() {
	Context("regression local testnet", func() {
		Context("when using electrs indexer", func() {
			It("should be able to get all the data without any error", func() {
				By("Initialise the logger")
				logger, err := zap.NewDevelopment()
				Expect(err).To(BeNil())

				By("Initialise a local regnet client")
				user := testutil.ParseStringEnv("BTC_USER", "")
				password := testutil.ParseStringEnv("BTC_PASSWORD", "")
				opts := btc.DefaultClientOptions().WithUser(user).WithPassword(password)
				indexerClient := btc.NewElectrsIndexerClient(logger, btc.DefaultElectrsIndexerURL)
				client := btc.NewClient(opts, logger, indexerClient)

				By("Net()")
				Expect(reflect.DeepEqual(client.Net(), &chaincfg.RegressionNetParams)).Should(BeTrue())

				By("LatestBlock()")
				height, hash, err := client.LatestBlock(context.Background())
				Expect(err).To(BeNil())
				Expect(height).Should(BeNumerically(">", 100))
				Expect(len(hash)).Should(Equal(64))
				Expect(hash).ShouldNot(Equal("0000000000000000000000000000000000000000000000000000000000000000"))

				By("GetBlockByHeight()")
				block, err := client.GetBlockByHeight(context.Background(), 1)
				Expect(err).To(BeNil())
				Expect(block.Hash).ShouldNot(Equal("0000000000000000000000000000000000000000000000000000000000000000"))
				Expect(len(block.Tx)).Should(Equal(1))

				By("GetRawTransaction()")
				txHash, err := chainhash.NewHashFromStr(block.Tx[0])
				Expect(err).To(BeNil())
				rawTx, err := client.GetRawTransaction(context.Background(), txHash.CloneBytes())
				Expect(err).To(BeNil())
				Expect(len(rawTx.Hex)).Should(BeNumerically(">", 1))

				By("GetUTXOs()")
				scriptPubkey, err := hex.DecodeString(rawTx.Vout[0].ScriptPubKey.Hex)
				Expect(err).To(BeNil())
				_, addrs, _, err := txscript.ExtractPkScriptAddrs(scriptPubkey, &chaincfg.RegressionNetParams)
				Expect(len(addrs)).Should(Equal(1))
				utxos, err := client.GetUTXOs(context.Background(), addrs[0])
				Expect(err).To(BeNil())
				Expect(len(utxos)).Should(BeNumerically(">", 0))
			})
		})
	})

	Context("bitcoin testnet", func() {
		Context("when using quicknode api", func() {
			It("should get the data without any error", func(ctx SpecContext) {
				By("Initialise the logger")
				logger, err := zap.NewDevelopment()
				Expect(err).To(BeNil())

				By("Initialise a testnet client")
				host := os.Getenv("BTC_RPC")
				indexerURL := os.Getenv("BTC_INDEXER_QUICKNODE")
				indexerClient := btc.NewQuickNodeIndexerClient(logger, indexerURL)
				opts := btc.DefaultClientOptions().WithNet(&chaincfg.TestNet3Params).WithHost(host)
				client := btc.NewClient(opts, logger, indexerClient)

				By("Net()")
				Expect(client.Net()).Should(Equal(&chaincfg.TestNet3Params))

				By("LatestBlock()")
				height, hash, err := client.LatestBlock(context.Background())
				Expect(err).To(BeNil())
				Expect(height).Should(BeNumerically(">", 1))
				Expect(len(hash)).Should(Equal(64))
				Expect(hash).ShouldNot(Equal("0000000000000000000000000000000000000000000000000000000000000000"))

				By("GetBlockByHeight()")
				block, err := client.GetBlockByHeight(context.Background(), 2431641)
				Expect(err).To(BeNil())
				Expect(block.Hash).Should(Equal("00000000792fb6eaaf9591d0b809c43e6db4ca4a7a34e2543a95b931f6097729"))
				Expect(len(block.Tx)).Should(Equal(96))

				By("GetRawTransaction()")
				txHashStr := "0351dbc63c2331cf7a7f504979a946a9fbff18dfd2dd5e73cb82a9dd0b8eef36"
				txHash, err := chainhash.NewHashFromStr(txHashStr)
				Expect(err).To(BeNil())
				rawTx, err := client.GetRawTransaction(context.Background(), txHash.CloneBytes())
				Expect(err).To(BeNil())
				Expect(rawTx.Hex).Should(Equal("0200000001935f334905b4786f296c067408e3ab028c28670dbad8fc6f84b2564645293018000000006a473044022048cc0fb22526365e08c9f71580ac5a9554784d8b2c0b2ee0c6b08a16fac636a802206c680787f9d23ba650e3171434ea8ca7775a305ade07dc00bd7d3d04fd127b76012102ef96baa1bff1890335bd9c0d88d428932e7cb5aee3ddfe1b79be06cd01aa733dfeffffff0269a77857000000001600140712fe0355630e4ffad1003a285bf2f69c13832999f10400000000001600145ce5a037f931b963f87421f447f126d86d781fc8971a2500"))

				By("GetUTXOs()")
				addr, err := btcutil.DecodeAddress("mvb4nE8firm7abnss9r7j5tra2w7M7uktF", &chaincfg.TestNet3Params)
				Expect(err).To(BeNil())
				utxos, err := client.GetUTXOs(ctx, addr)
				Expect(err).To(BeNil())
				Expect(len(utxos)).Should(BeNumerically(">", 1))
			})
		})
	})
})
