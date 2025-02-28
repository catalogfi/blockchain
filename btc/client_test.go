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
	Context("regression testnet", func() {
		Context("when using the Client", func() {
			It("should be able to get all the data without any error", func(ctx context.Context) {
				By("Net()")
				Expect(reflect.DeepEqual(client.Net(), &chaincfg.RegressionNetParams)).Should(BeTrue())

				By("LatestBlock()")
				info, err := client.GetBlockchainInfo(ctx)
				Expect(err).To(BeNil())
				Expect(info.Blocks).Should(BeNumerically(">=", 100))
				Expect(info.BestBlockHash).ShouldNot(Equal("0000000000000000000000000000000000000000000000000000000000000000"))

				By("Create a new tx")
				key, addr, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
				Expect(err).To(BeNil())
				txid, err := btctest.Faucet(addr.EncodeAddress())
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

				By("GetMempoolEntry()")
				utxos := []btc.UTXO{
					{
						TxID:   txid.String(),
						Vout:   uint32(vout),
						Amount: 1e8,
					},
				}
				sizer := btc.NewSizeEstimator(btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH, utxos...)
				feeMode := btc.MinFeeRateMode(10000, sizer)
				tx1, err := btc.BuildTx(network, feeMode, utxos, nil, nil, addr)
				Expect(err).To(BeNil())
				Expect(btc.SignTx(waddrmgr.PubKeyHash, tx1, key, utxos)).Should(Succeed())
				Expect(client.SubmitTx(ctx, tx1)).Should(Succeed())
				Eventually(func() error {
					_, err := client.GetMempoolEntry(ctx, tx1.TxHash().String())
					return err
				}).Should(Succeed())
			})
		})
	})

	Context("errors", func() {
		It("should return specific errors", func(ctx context.Context) {
			By("New address")
			privKey, pkAddr, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())

			By("funding the addresses")
			_, err = btctest.Faucet(pkAddr.EncodeAddress())
			Expect(err).To(BeNil())
			time.Sleep(5 * time.Second)

			By("Construct a new tx")
			utxos, err := indexer.GetUTXOs(ctx, pkAddr)
			Expect(err).To(BeNil())
			amount := int64(1e5)
			recipients := []btc.Recipient{btc.NewRecipient(pkAddr.EncodeAddress(), amount)}
			sizer := btc.NewSizeEstimator(btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH, utxos...)
			feeMode := btc.MinFeeRateMode(10000, sizer)
			transaction, err := btc.BuildTx(network, feeMode, nil, utxos, recipients, pkAddr)
			Expect(err).To(BeNil())

			By("Sign the transaction inputs")
			Expect(btc.SignTx(waddrmgr.PubKeyHash, transaction, privKey, utxos)).Should(Succeed())

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
			Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())
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

			feeMode1 := btc.FixedFeesMode(1000)
			transaction1, err := btc.BuildTx(network, feeMode1, nil, utxos, recipients1, pkAddr)
			Expect(err).To(BeNil())
			Expect(btc.SignTx(waddrmgr.PubKeyHash, transaction1, privKey, utxos)).Should(Succeed())

			By("Expect a `ErrTxInputsMissingOrSpent` error if the tx is already in a block")
			err = client.SubmitTx(ctx, transaction1)
			Expect(errors.Is(err, btc.ErrTxInputsMissingOrSpent)).Should(BeTrue())

			By("Expect a `ErrTxNotInMempool` error if the tx is not in the mempool")
			_, err = client.GetMempoolEntry(ctx, transaction1.TxHash().String())
			Expect(errors.Is(err, btc.ErrTxNotInMempool)).Should(BeTrue())
		})

		It("should return an error when the utxo has been spent", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			By("Initialization keys ")
			privKey1, pkAddr1, err := btctest.NewBtcAddrWithFunds(network, waddrmgr.PubKeyHash, indexer)
			Expect(err).To(BeNil())
			_, pkAddr2, err := btctest.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(ctx, pkAddr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e5), btctest.RandomFeeRate()
			recipients := []btc.Recipient{btc.NewRecipient(pkAddr2.EncodeAddress(), amount)}
			sizer := btc.NewSizeEstimator(btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH, utxos...)
			feeMode := btc.MinFeeRateMode(feeRate, sizer)
			transaction, err := btc.BuildTx(network, feeMode, nil, utxos, recipients, pkAddr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.SignTx(waddrmgr.PubKeyHash, transaction, privKey1, utxos)).Should(Succeed())
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
			client := btc.NewClient(network, "http://0.0.0.0:18444", btcUsername, btcPassword)

			By("LatestBlock()")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			_, err := client.GetBlockchainInfo(ctx)
			Expect(err).ShouldNot(BeNil())
			cancel()

			By("GetRawTransaction()")
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			hash, err := chainhash.NewHashFromStr("0000000000000000000000000000000000000000000000000000000000000000")
			Expect(err).To(BeNil())
			_, err = client.GetRawTransaction(ctx, hash)
			Expect(err).ShouldNot(BeNil())
			cancel()

			By("GetBlock()")
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			_, err = client.GetBlock(ctx, hash)
			Expect(err).ShouldNot(BeNil())
			cancel()

			By("GetBlockVerbose()")
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			_, err = client.GetBlockVerbose(ctx, hash)
			Expect(err).ShouldNot(BeNil())
			cancel()

			By("GetBlockByHeight()")
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			_, err = client.GetBlockHash(ctx, 1)
			Expect(err).ShouldNot(BeNil())
			cancel()

			By("GetTxOut()")
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			_, err = client.GetTxOut(ctx, hash, 0)
			Expect(err).ShouldNot(BeNil())
			cancel()

			By("SubmitTx()")
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			err = client.SubmitTx(ctx, new(wire.MsgTx))
			Expect(err).ShouldNot(BeNil())
			cancel()
		})
	})
})
