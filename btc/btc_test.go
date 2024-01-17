package btc_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing/quick"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/catalogfi/blockchain/testutil"
	"github.com/fatih/color"
	"github.com/tyler-smith/go-bip39"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bitcoin", func() {
	Context("keys", func() {
		It("should generate deterministic keys from mnemonic and user public key", func() {
			test := func() bool {
				entropy, err := bip39.NewEntropy(256)
				Expect(err).To(BeNil())
				mnemonic, err := bip39.NewMnemonic(entropy)
				Expect(err).To(BeNil())
				key, err := btcec.NewPrivateKey()
				Expect(err).To(BeNil())

				extendedKey1, err := btc.GenerateSystemPrivKey(mnemonic, key.PubKey().SerializeCompressed())
				Expect(err).To(BeNil())
				extendedKey2, err := btc.GenerateSystemPrivKey(mnemonic, key.PubKey().SerializeCompressed())
				Expect(err).To(BeNil())

				Expect(bytes.Equal(extendedKey1.Serialize(), extendedKey2.Serialize())).Should(BeTrue())
				return true
			}

			Expect(quick.Check(test, nil)).NotTo(HaveOccurred())
		})
	})

	Context("RBF", func() {
		It("should be able to build transaction which support rbf", func(ctx context.Context) {
			By("Initialization (Update these fields if testing on testnet/mainnet)")
			network := &chaincfg.RegressionNetParams
			privKey1, p2pkhAddr1, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())
			_, p2pkhAddr2, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())
			indexer := btctest.RegtestIndexer()

			By("Funding the addresses")
			txhash1, err := testutil.NigiriFaucet(p2pkhAddr1.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address1 %v , txid = %v", p2pkhAddr1.EncodeAddress(), txhash1))
			time.Sleep(5 * time.Second)

			By("Construct a RBF tx which sends money from key1 to key2")
			utxos, err := indexer.GetUTXOs(ctx, p2pkhAddr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e7), 4
			recipients := []btc.Recipient{
				{
					To:     p2pkhAddr2.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildRbfTransaction(feeRate, network, btc.NewRawInputs(), utxos, recipients, btc.P2pkhUpdater, p2pkhAddr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			for i := range transaction.TxIn {
				pkScript, err := txscript.PayToAddrScript(p2pkhAddr1)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(transaction, i, pkScript, txscript.SigHashAll, privKey1, true)
				Expect(err).To(BeNil())
				transaction.TxIn[i].SignatureScript = sigScript
			}
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
			By(color.GreenString("RBF tx hash = %v", transaction.TxHash().String()))

			By("Construct a tx with higher fees")
			transaction1, err := btc.BuildTransaction(feeRate+5, network, btc.NewRawInputs(), utxos, recipients, btc.P2pkhUpdater, p2pkhAddr1)
			Expect(err).To(BeNil())
			transaction1.TxOut[1].Value++

			By("Sign and submit replacement tx")
			for i := range transaction1.TxIn {
				pkScript, err := txscript.PayToAddrScript(p2pkhAddr1)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(transaction1, i, pkScript, txscript.SigHashAll, privKey1, true)
				Expect(err).To(BeNil())
				transaction1.TxIn[i].SignatureScript = sigScript
			}
			Expect(indexer.SubmitTx(ctx, transaction1)).Should(Succeed())
			color.Green("Replacement tx = %v", transaction1.TxHash().String())
		})

		It("should get an error when trying to replace a mined tx", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			By("Initialise a local regnet client")
			network := &chaincfg.RegressionNetParams
			client, err := btctest.RegtestClient()
			Expect(err).To(BeNil())
			indexer := btctest.RegtestIndexer()

			By("New address")
			privKey, pkAddr, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())

			_, toAddr, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())

			By("funding the addresses")
			txhash, err := testutil.NigiriFaucet(pkAddr.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address1 %v , txid = %v", pkAddr.EncodeAddress(), txhash))
			time.Sleep(5 * time.Second)

			By("Construct a new tx")
			utxos, err := indexer.GetUTXOs(ctx, pkAddr)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e5), 5
			recipients := []btc.Recipient{
				{
					To:     toAddr.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildRbfTransaction(feeRate, network, btc.NewRawInputs(), utxos, recipients, btc.P2pkhUpdater, pkAddr)
			Expect(err).To(BeNil())

			By("Sign the transaction inputs")
			for i := range transaction.TxIn {
				pkScript, err := txscript.PayToAddrScript(pkAddr)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(transaction, i, pkScript, txscript.SigHashAll, privKey, true)
				Expect(err).To(BeNil())
				transaction.TxIn[i].SignatureScript = sigScript
			}

			By("Submit the transaction")
			Expect(client.SubmitTx(ctx, transaction)).Should(Succeed())
			By(fmt.Sprintf("Funding tx hash = %v", color.YellowString(transaction.TxHash().String())))
			time.Sleep(time.Second)

			By("Build a new tx with higher fee")
			feeRate += 2
			transaction, err = btc.BuildRbfTransaction(feeRate, network, btc.NewRawInputs(), utxos, recipients, btc.P2pkhUpdater, pkAddr)
			Expect(err).To(BeNil())

			By("Sign the transaction inputs")
			for i := range transaction.TxIn {
				pkScript, err := txscript.PayToAddrScript(pkAddr)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(transaction, i, pkScript, txscript.SigHashAll, privKey, true)
				Expect(err).To(BeNil())
				transaction.TxIn[i].SignatureScript = sigScript
				Expect(testutil.NigiriNewBlock()).Should(Succeed())
			}

			By("Submit the transaction again and it should be rejected")
			err = client.SubmitTx(ctx, transaction)
			Expect(errors.Is(err, btc.ErrTxInputsMissingOrSpent)).Should(BeTrue())
		})
	})
})
