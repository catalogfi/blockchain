package btc_test

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/fatih/color"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("bitcoin client", func() {
	Context("Client initialization", func() {
		It("should return an error if providing an unknown chain params", func() {
			config := &rpcclient.ConnConfig{
				Params:       "",
				Host:         btctest.DefaultRegtestHost,
				HTTPPostMode: true,
				DisableTLS:   true,
			}
			_, err := btc.NewClient(config)
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("regression local testnet", func() {
		Context("when using the Client", func() {
			It("should be able to get all the data without any error", func(ctx context.Context) {
				By("Net()")
				Expect(reflect.DeepEqual(client.Net(), &chaincfg.RegressionNetParams)).Should(BeTrue())

				By("LatestBlock()")
				height, hash, err := client.LatestBlock(ctx)
				Expect(err).To(BeNil())
				Expect(height).Should(BeNumerically(">=", 100))
				Expect(len(hash)).Should(Equal(64))
				Expect(hash).ShouldNot(Equal("0000000000000000000000000000000000000000000000000000000000000000"))

				By("Create a new tx")
				_, addr, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
				Expect(err).To(BeNil())
				txid, err := btctest.NigiriFaucet(addr.EncodeAddress())
				Expect(err).To(BeNil())
				time.Sleep(500 * time.Millisecond)

				By("GetRawTransaction()")
				rawTx, err := client.GetRawTransaction(ctx, txid)
				Expect(err).To(BeNil())
				Expect(rawTx.Txid).Should(Equal(txid.String()))

				By("GetBlockByHash()")
				blockHash, err := chainhash.NewHashFromStr(rawTx.BlockHash)
				Expect(err).To(BeNil())
				block, err := client.GetBlock(ctx, blockHash)
				Expect(err).To(BeNil())
				Expect(block.Hash).Should(Equal(rawTx.BlockHash))

				By("GetBlockHash()")
				bHash, err := client.GetBlockHash(ctx, block.Height)
				Expect(err).To(BeNil())
				Expect(bHash.String()).Should(Equal(rawTx.BlockHash))

				By("GetBlockVerbose()")
				blockVerbose, err := client.GetBlockVerbose(ctx, blockHash)
				Expect(err).To(BeNil())
				Expect(blockVerbose.Hash).Should(Equal(rawTx.BlockHash))

				By("GetTxOut()")
				vout := 0
				for index, out := range rawTx.Vout {
					if out.Value == 1 {
						vout = index
						break
					}
				}
				txout, err := client.GetTxOut(ctx, txid, uint32(vout))
				Expect(err).To(BeNil())
				Expect(txout.Value).Should(Equal(float64(1)))
				scriptPubkey, err := txscript.PayToAddrScript(addr)
				Expect(err).To(BeNil())
				Expect(txout.ScriptPubKey.Hex).Should(Equal(hex.EncodeToString(scriptPubkey)))

				By("GetNetworkInfo()")
				netInfo, err := client.GetNetworkInfo(ctx)
				Expect(err).To(BeNil())
				Expect(netInfo.RelayFee).Should(BeNumerically(">=", 0))
			})
		})
	})

	Context("errors", func() {
		It("should return specific errors", func(ctx context.Context) {
			By("New address")
			privKey, pkAddr, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())

			By("funding the addresses")
			_, err = btctest.NigiriFaucet(pkAddr.EncodeAddress())
			Expect(err).To(BeNil())
			time.Sleep(5 * time.Second)

			By("Construct a new tx")
			utxos, err := indexer.GetUTXOs(ctx, pkAddr)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e5), 5
			recipients := []btc.Recipient{
				{
					To:     pkAddr.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildTransaction(network, feeRate, btc.NewRawInputs(), utxos, btc.P2pkhUpdater, recipients, pkAddr)
			Expect(err).To(BeNil())

			By("Sign the transaction inputs")
			Expect(btc.SignP2pkhTx(network, privKey, transaction)).Should(Succeed())

			By("Expect `ErrTxNotFound` before submitting the tx")
			txid := transaction.TxHash()
			_, err = client.GetRawTransaction(ctx, &txid)
			Expect(errors.Is(err, btc.ErrTxNotFound)).Should(BeTrue())

			By("Submit the transaction")
			Expect(client.SubmitTx(ctx, transaction)).Should(Succeed())
			By(fmt.Sprintf("Funding tx hash = %v", color.YellowString(transaction.TxHash().String())))
			time.Sleep(time.Second)

			By("We should not get any error fetching the tx details")
			_, err = client.GetRawTransaction(ctx, &txid)
			Expect(err).Should(BeNil())

			By("Expect a `ErrAlreadyInChain` error if the tx is already in a block")
			Expect(btctest.NigiriNewBlock()).Should(Succeed())
			time.Sleep(1 * time.Second)
			err = client.SubmitTx(ctx, transaction)
			Expect(errors.Is(err, btc.ErrAlreadyInChain)).Should(BeTrue())

			By("Try construct a new transaction spending the same input")
			recipients1 := []btc.Recipient{
				{
					To:     pkAddr.EncodeAddress(),
					Amount: 2 * amount,
				},
			}

			transaction1, err := btc.BuildTransaction(network, feeRate, btc.NewRawInputs(), utxos, btc.P2pkhUpdater, recipients1, pkAddr)
			Expect(err).To(BeNil())
			Expect(btc.SignP2pkhTx(network, privKey, transaction1)).Should(Succeed())

			By("Expect a `ErrTxInputsMissingOrSpent` error if the tx is already in a block")
			err = client.SubmitTx(ctx, transaction1)
			Expect(errors.Is(err, btc.ErrTxInputsMissingOrSpent)).Should(BeTrue())
		})

		It("should return an error when the utxo has been spent", func(ctx context.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			By("Initialization keys ")
			network := &chaincfg.RegressionNetParams
			privKey1, pkAddr1, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())
			_, pkAddr2, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())

			By("Funding the addresses")
			txhash1, err := btctest.NigiriFaucet(pkAddr1.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address1 %v , txid = %v", pkAddr1.EncodeAddress(), txhash1))
			time.Sleep(5 * time.Second)

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(ctx, pkAddr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e5), 5
			recipients := []btc.Recipient{
				{
					To:     pkAddr2.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildTransaction(network, feeRate, btc.NewRawInputs(), utxos, btc.P2pkhUpdater, recipients, pkAddr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.SignP2pkhTx(network, privKey1, transaction)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())

			By("Expect an error if the utxo is spent")
			time.Sleep(time.Second)
			for _, input := range transaction.TxIn {
				res, err := client.GetTxOut(ctx, &input.PreviousOutPoint.Hash, input.PreviousOutPoint.Index)
				Expect(err).Should(BeNil())
				Expect(res).Should(BeNil())
			}
		})
	})

	Context("when the server is offline", func() {
		It("should err out when the context is done", func() {
			By("Simulate a client pointing to a offline server")
			config := &rpcclient.ConnConfig{
				Params:       network.Name,
				Host:         "0.0.0.0:18444",
				User:         btcUsername,
				Pass:         btcPassword,
				HTTPPostMode: true,
				DisableTLS:   true,
			}
			client, err := btc.NewClient(config)
			Expect(err).Should(BeNil())

			By("LatestBlock()")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			_, _, err = client.LatestBlock(ctx)
			Expect(errors.Is(err, context.DeadlineExceeded)).Should(BeTrue())
			cancel()

			By("GetRawTransaction()")
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			hash, err := chainhash.NewHashFromStr("0000000000000000000000000000000000000000000000000000000000000000")
			Expect(err).To(BeNil())
			_, err = client.GetRawTransaction(ctx, hash)
			Expect(errors.Is(err, context.DeadlineExceeded)).Should(BeTrue())
			cancel()

			By("GetBlock()")
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			_, err = client.GetBlock(ctx, hash)
			Expect(errors.Is(err, context.DeadlineExceeded)).Should(BeTrue())
			cancel()

			By("GetBlockVerbose()")
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			_, err = client.GetBlockVerbose(ctx, hash)
			Expect(errors.Is(err, context.DeadlineExceeded)).Should(BeTrue())
			cancel()

			By("GetBlockByHeight()")
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			_, err = client.GetBlockHash(ctx, 1)
			Expect(errors.Is(err, context.DeadlineExceeded)).Should(BeTrue())
			cancel()

			By("GetTxOut()")
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			_, err = client.GetTxOut(ctx, hash, 0)
			Expect(errors.Is(err, context.DeadlineExceeded)).Should(BeTrue())
			cancel()

			By("SubmitTx()")
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			err = client.SubmitTx(ctx, new(wire.MsgTx))
			Expect(errors.Is(err, context.DeadlineExceeded)).Should(BeTrue())
			cancel()
		})
	})
})
