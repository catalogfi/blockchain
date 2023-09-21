package btc_test

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/testutil"
	"github.com/fatih/color"
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

	Context("errors", func() {
		It("should return specific errors", func() {
			By("Initialise a local regnet client")
			network := &chaincfg.RegressionNetParams
			client, err := RegtestClient()
			Expect(err).To(BeNil())
			indexer := RegtestIndexer()

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
			amount, fee := int64(1e5), int64(500)
			inputs, err := btc.PickUTXOs(utxos, amount, fee)
			Expect(err).To(BeNil())
			recipients := []btc.Recipient{
				{
					To:     pkAddr.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, _, err := btc.BuildTx(network, inputs, recipients, fee, pkAddr)
			Expect(err).To(BeNil())
			for i := range transaction.TxIn {
				pkScript, err := txscript.PayToAddrScript(pkAddr)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(transaction, i, pkScript, txscript.SigHashAll, privKey, true)
				Expect(err).To(BeNil())
				transaction.TxIn[i].SignatureScript = sigScript
			}

			By("Expect `ErrTxNotFound` before submitting the tx")
			txid := transaction.TxHash()
			_, err = client.GetRawTransaction(txid.CloneBytes())
			Expect(errors.Is(err, btc.ErrTxNotFound)).Should(BeTrue())

			By("Submit the transaction")
			Expect(client.SubmitTx(transaction)).Should(Succeed())
			By(fmt.Sprintf("Funding tx hash = %v", color.YellowString(transaction.TxHash().String())))
			time.Sleep(time.Second)

			By("We should not get any error fetching the tx details")
			_, err = client.GetRawTransaction(txid.CloneBytes())
			Expect(err).Should(BeNil())

			By("Expect a `ErrAlreadyInChain` error if the tx is already in a block")
			Expect(testutil.NigiriNewBlock()).Should(Succeed())
			time.Sleep(1 * time.Second)
			err = client.SubmitTx(transaction)
			Expect(errors.Is(err, btc.ErrAlreadyInChain)).Should(BeTrue())

			By("Try construct a new transaction spending the same input")
			recipients1 := []btc.Recipient{
				{
					To:     pkAddr.EncodeAddress(),
					Amount: 2 * amount,
				},
			}
			transaction1, _, err := btc.BuildTx(network, inputs, recipients1, fee, pkAddr)
			Expect(err).To(BeNil())
			for i := range transaction1.TxIn {
				pkScript, err := txscript.PayToAddrScript(pkAddr)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(transaction1, i, pkScript, txscript.SigHashAll, privKey, true)
				Expect(err).To(BeNil())
				transaction1.TxIn[i].SignatureScript = sigScript
			}
			By("Expect a `ErrAlreadyInChain` error if the tx is already in a block")
			err = client.SubmitTx(transaction1)
			Expect(errors.Is(err, btc.ErrTxInputsMissingOrSpent)).Should(BeTrue())
		})

		It("should return an error when the utxo has been spent", func() {
			By("Initialization keys ")
			network := &chaincfg.RegressionNetParams
			privKey1, err := btcec.NewPrivateKey()
			Expect(err).To(BeNil())
			pubKey1 := privKey1.PubKey()
			pkAddr1, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey1.SerializeCompressed()), network)
			Expect(err).To(BeNil())
			privKey2, err := btcec.NewPrivateKey()
			Expect(err).To(BeNil())
			pubKey2 := privKey2.PubKey()
			pkAddr2, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey2.SerializeCompressed()), network)
			Expect(err).To(BeNil())
			client, err := RegtestClient()
			Expect(err).To(BeNil())
			indexer := RegtestIndexer()

			By("Funding the addresses")
			txhash1, err := testutil.NigiriFaucet(pkAddr1.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address1 %v , txid = %v", pkAddr1.EncodeAddress(), txhash1))
			time.Sleep(5 * time.Second)

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(context.Background(), pkAddr1)
			Expect(err).To(BeNil())
			amount, fee := int64(1e5), int64(500)
			inputs, err := btc.PickUTXOs(utxos, amount, fee)
			Expect(err).To(BeNil())
			recipients := []btc.Recipient{
				{
					To:     pkAddr2.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, _, err := btc.BuildTx(network, inputs, recipients, fee, pkAddr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			for i := range transaction.TxIn {
				pkScript, err := txscript.PayToAddrScript(pkAddr1)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(transaction, i, pkScript, txscript.SigHashAll, privKey1, true)
				Expect(err).To(BeNil())
				transaction.TxIn[i].SignatureScript = sigScript
			}
			Expect(indexer.SubmitTx(context.Background(), transaction)).Should(Succeed())

			By("Expect an error if the utxo is spent")
			time.Sleep(time.Second)
			for _, input := range inputs {
				res, err := client.GetTxOut(input.TxID, input.Vout)
				Expect(err).Should(BeNil())
				Expect(res).Should(BeNil())
			}
		})
	})
})
