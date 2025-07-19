package btc_test

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/localnet"
	"github.com/fatih/color"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("bitcoin client", func() {
	Context("Client initialization", func() {
		It("should return an error if providing an unknown chain params", func() {
			config := &rpcclient.ConnConfig{
				Params:       "",
				Host:         localnet.DefaultRegtestHost,
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
				_, addr, err := localnet.NewBtcKey(network, waddrmgr.PubKeyHash)
				Expect(err).To(BeNil())
				txid, err := localnet.FundBTC(addr.EncodeAddress())
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
			privKey, pkAddr, err := localnet.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())

			By("funding the addresses")
			_, err = localnet.FundBTC(pkAddr.EncodeAddress())
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
			Expect(localnet.MineBTCBlock()).Should(Succeed())
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

		It("should return an error when the utxo has been spent", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			By("Initialization keys ")
			network := &chaincfg.RegressionNetParams
			privKey1, pkAddr1, err := localnet.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())
			_, pkAddr2, err := localnet.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())

			By("Funding the addresses")
			txhash1, err := localnet.FundBTC(pkAddr1.EncodeAddress())
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

	Context("Test BitcoinRPCClient", func() {
		bitcoinClient := btc.NewBitcoinClient("admin1", "123", "http://0.0.0.0:18443")
		chainParams := chaincfg.RegressionNetParams
		indexer := localnet.BTCIndexer()
		feeEstimator := btc.NewFixFeeEstimator(16)

		newWallet := func() btc.Wallet {
			privKey, err := btcec.NewPrivateKey()
			Expect(err).To(BeNil())
			wallet, err := btc.NewSimpleWallet(privKey, &chainParams, indexer, feeEstimator, btc.HighFee)
			Expect(err).To(BeNil())
			return wallet
		}

		alice := newWallet()
		bob := newWallet()

		_, err := localnet.FundBitcoin(alice.Address().EncodeAddress(), indexer)
		Expect(err).To(BeNil())

		send := func(wallet btc.Wallet, to btcutil.Address, amount int64) (string, btc.Transaction) {
			txid, err := wallet.Send(context.Background(), []btc.SendRequest{{Amount: amount, To: to}}, nil, nil)
			Expect(err).To(BeNil())
			tx, _, err := wallet.Status(context.Background(), txid)
			Expect(err).To(BeNil())
			return txid, tx
		}

		It("should return mempool entry", func(ctx context.Context) {
			// Helper to check mempool entry fee
			checkMempoolEntryFee := func(entry *btc.GetMempoolEntryResult, expectedFee int64) {
				fee, err := btcutil.NewAmount(entry.Fees.Descendant)
				Expect(err).To(BeNil())
				Expect(float64(fee)).Should(Equal(float64(expectedFee)))
			}

			// Alice sends to Bob
			txid1, tx1 := send(alice, bob.Address(), 1000000)

			By("tx with no descendants")
			entry1, err := bitcoinClient.GetMempoolEntry(ctx, txid1)
			Expect(err).To(BeNil())
			Expect(entry1).ShouldNot(BeNil())
			checkMempoolEntryFee(entry1, tx1.Fee)

			// Bob sends to a random address
			By("tx with one direct descendant")
			randAddr, err := randomP2wpkhAddress(chainParams)
			Expect(err).To(BeNil())
			_, tx2 := send(bob, randAddr, 10000)

			By("GetMempoolEntry() for Alice's tx after Bob's tx")
			entry1b, err := bitcoinClient.GetMempoolEntry(ctx, txid1)
			Expect(err).To(BeNil())
			Expect(entry1b).ShouldNot(BeNil())
			Expect(len(entry1b.SpendBy)).Should(Equal(1))
			checkMempoolEntryFee(entry1b, tx1.Fee+tx2.Fee)

			// Bob sends another tx to the same random address
			By("tx with indirect descendant")
			_, tx3 := send(bob, randAddr, 20000)

			entry1c, err := bitcoinClient.GetMempoolEntry(ctx, txid1)
			Expect(err).To(BeNil())
			Expect(entry1c).ShouldNot(BeNil())
			Expect(len(entry1c.SpendBy)).Should(Equal(1))
			checkMempoolEntryFee(entry1c, tx1.Fee+tx2.Fee+tx3.Fee)

			By("GetMempoolDescendants() for Alice's tx")
			descendants, err := bitcoinClient.GetMempoolDescendants(ctx, txid1)
			Expect(err).To(BeNil())
			Expect(descendants).ShouldNot(BeNil())
			Expect(len(descendants)).Should(Equal(2))

			By("GetDescendantFees() for Alice's tx")
			descendantFee, err := bitcoinClient.GetDescendantsFee(ctx, txid1)
			Expect(err).To(BeNil())
			Expect(float64(descendantFee)).Should(Equal(float64(tx2.Fee + tx3.Fee)))
		})

		It("should return correct rbf fee tx information", func(ctx context.Context) {
			// Helper to check RBF info against mempool entry
			checkRBFInfo := func(entry *btc.GetMempoolEntryResult, rbfInfo *btc.RBFTxFeeInfo, vsize float64) {
				entryTotalFee := float64(entry.Fees.Descendant * 1e8)
				descendantFee := float64(entry.Fees.Descendant-entry.Fees.Base) * 1e8
				feeRate := entryTotalFee / vsize

				Expect(rbfInfo.TxFeeRate).Should(Equal(feeRate))
				Expect(float64(rbfInfo.TotalFee)).Should(Equal(entryTotalFee))
				Expect(float64(rbfInfo.DescendantFee)).Should(Equal(descendantFee))
			}

			// Alice sends to Bob
			txid1, _ := send(alice, bob.Address(), 1000000)

			By("rbf tx info for tx with no descendants")
			rbfInfo, err := bitcoinClient.GetRBFTxFeeInfo(ctx, txid1)
			Expect(err).To(BeNil())
			Expect(rbfInfo).ShouldNot(BeNil())
			mempoolEntry, err := bitcoinClient.GetMempoolEntry(ctx, txid1)
			Expect(err).To(BeNil())
			checkRBFInfo(mempoolEntry, rbfInfo, float64(mempoolEntry.VSize))

			// get fee info with one direct descendant
			randAddr, err := randomP2wpkhAddress(chainParams)
			Expect(err).To(BeNil())
			_, directDescTx := send(bob, randAddr, 10000)

			By("rbf tx info for tx with one direct descendant")
			entry1b, err := bitcoinClient.GetMempoolEntry(ctx, txid1)
			Expect(err).To(BeNil())
			Expect(entry1b).ShouldNot(BeNil())
			rbfInfo, err = bitcoinClient.GetRBFTxFeeInfo(ctx, txid1)
			Expect(err).To(BeNil())
			checkRBFInfo(entry1b, rbfInfo, float64(entry1b.DescendantSize))

			// get fee info with one indirect descendant
			_, _ = send(bob, randAddr, 20000)
			By("rbf tx info for tx with indirect descendant")
			entry1c, err := bitcoinClient.GetMempoolEntry(ctx, txid1)
			Expect(err).To(BeNil())
			Expect(entry1c).ShouldNot(BeNil())
			rbfInfo, err = bitcoinClient.GetRBFTxFeeInfo(ctx, txid1)
			Expect(err).To(BeNil())
			mempoolDesc, err := bitcoinClient.GetMempoolDescendants(ctx, txid1)
			Expect(err).To(BeNil())
			Expect(mempoolDesc).ShouldNot(BeNil())

			// For indirect descendant, calculate feeRate using base fee + directDescTx.Fee and vsize + descendant vsize
			entryTotalFee := float64(entry1c.Fees.Descendant * 1e8)
			descendantFee := float64(entry1c.Fees.Descendant-entry1c.Fees.Base) * 1e8
			feeRate := float64(int64(entry1c.Fees.Base*1e8)+directDescTx.Fee) /
				float64(int64(entry1c.VSize)+int64(mempoolDesc[directDescTx.TxID].Vsize))

			Expect(rbfInfo.TxFeeRate).Should(Equal(feeRate))
			Expect(float64(rbfInfo.TotalFee)).Should(Equal(math.Ceil(entryTotalFee)))
			Expect(float64(rbfInfo.DescendantFee)).Should(Equal(math.Ceil(descendantFee)))
		})
	})
})
