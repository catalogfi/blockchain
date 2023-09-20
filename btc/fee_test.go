package btc_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
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

var _ = Describe("bitcoin fees", func() {
	Context("fee estimator", func() {
		Context("mempool estimator", func() {
			It("should return the live data from mempool", func() {
				estimator := btc.NewMempoolFeeEstimator(&chaincfg.MainNetParams, btc.MempoolFeeAPI, 15*time.Second)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())

				Expect(fees.Minimum).Should(BeNumerically(">", 1))

				Expect(fees.Economy).Should(BeNumerically(">", 1))
				Expect(fees.Economy).Should(BeNumerically(">=", fees.Minimum))

				Expect(fees.Low).Should(BeNumerically(">", 1))
				Expect(fees.Low).Should(BeNumerically(">=", fees.Economy))

				Expect(fees.Medium).Should(BeNumerically(">", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", fees.Low))

				Expect(fees.High).Should(BeNumerically(">", 1))
				Expect(fees.High).Should(BeNumerically(">=", fees.Medium))
			})

			It("should use the cached result if the more requests are made during certain time period", func() {
				estimator := btc.NewMempoolFeeEstimator(&chaincfg.MainNetParams, btc.MempoolFeeAPI, time.Second)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())

				time.Sleep(500 * time.Millisecond)

				fees1, err := estimator.FeeSuggestion()
				Expect(fees.Minimum).Should(Equal(fees1.Minimum))
				Expect(fees.Economy).Should(Equal(fees1.Economy))
				Expect(fees.Low).Should(Equal(fees1.Low))
				Expect(fees.Medium).Should(Equal(fees1.Medium))
				Expect(fees.High).Should(Equal(fees1.High))

				time.Sleep(500 * time.Millisecond)
				_, err = estimator.FeeSuggestion()
				Expect(err).Should(BeNil())
			})

			It("should return a fixed fee for non-mainnet network", func() {
				By("Testnet")
				estimatorTestnet := btc.NewMempoolFeeEstimator(&chaincfg.TestNet3Params, "", 15*time.Second)
				fees, err := estimatorTestnet.FeeSuggestion()
				Expect(err).Should(BeNil())
				Expect(fees.Minimum).Should(Equal(2))
				Expect(fees.Economy).Should(Equal(2))
				Expect(fees.Low).Should(Equal(2))
				Expect(fees.Medium).Should(Equal(2))
				Expect(fees.High).Should(Equal(2))

				By("Regnet")
				estimatorRegnet := btc.NewMempoolFeeEstimator(&chaincfg.RegressionNetParams, "", 15*time.Second)
				fees, err = estimatorRegnet.FeeSuggestion()
				Expect(err).Should(BeNil())
				Expect(fees.Minimum).Should(Equal(2))
				Expect(fees.Economy).Should(Equal(2))
				Expect(fees.Low).Should(Equal(2))
				Expect(fees.Medium).Should(Equal(2))
				Expect(fees.High).Should(Equal(2))
			})
		})

		Context("blockstream estimator", func() {
			It("should return the live data from blockstream api", func() {
				estimator := btc.NewBlockstreamFeeEstimator(&chaincfg.MainNetParams, btc.BlockstreamAPI, 15*time.Second)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())

				Expect(fees.Minimum).Should(BeNumerically(">", 1))

				Expect(fees.Economy).Should(BeNumerically(">", 1))
				Expect(fees.Economy).Should(BeNumerically(">=", fees.Minimum))

				Expect(fees.Low).Should(BeNumerically(">", 1))
				Expect(fees.Low).Should(BeNumerically(">=", fees.Economy))

				Expect(fees.Medium).Should(BeNumerically(">", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", fees.Low))

				Expect(fees.High).Should(BeNumerically(">", 1))
				Expect(fees.High).Should(BeNumerically(">=", fees.Medium))
			})

			It("should use the cached result if the more requests are made during certain time period", func() {
				estimator := btc.NewBlockstreamFeeEstimator(&chaincfg.MainNetParams, btc.BlockstreamAPI, time.Second)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())

				time.Sleep(500 * time.Millisecond)

				fees1, err := estimator.FeeSuggestion()
				Expect(fees.Minimum).Should(Equal(fees1.Minimum))
				Expect(fees.Economy).Should(Equal(fees1.Economy))
				Expect(fees.Low).Should(Equal(fees1.Low))
				Expect(fees.Medium).Should(Equal(fees1.Medium))
				Expect(fees.High).Should(Equal(fees1.High))

				time.Sleep(500 * time.Millisecond)
				_, err = estimator.FeeSuggestion()
				Expect(err).Should(BeNil())
			})

			It("should return a fixed fee for non-mainnet network", func() {
				By("Testnet")
				estimatorTestnet := btc.NewBlockstreamFeeEstimator(&chaincfg.TestNet3Params, "", 15*time.Second)
				fees, err := estimatorTestnet.FeeSuggestion()
				Expect(err).Should(BeNil())
				Expect(fees.Minimum).Should(Equal(2))
				Expect(fees.Economy).Should(Equal(2))
				Expect(fees.Low).Should(Equal(2))
				Expect(fees.Medium).Should(Equal(2))
				Expect(fees.High).Should(Equal(2))

				By("Regnet")
				estimatorRegnet := btc.NewBlockstreamFeeEstimator(&chaincfg.RegressionNetParams, "", 15*time.Second)
				fees, err = estimatorRegnet.FeeSuggestion()
				Expect(err).Should(BeNil())
				Expect(fees.Minimum).Should(Equal(2))
				Expect(fees.Economy).Should(Equal(2))
				Expect(fees.Low).Should(Equal(2))
				Expect(fees.Medium).Should(Equal(2))
				Expect(fees.High).Should(Equal(2))
			})
		})

		Context("fix fee estimator", func() {
			It("should return a fixed fee for fixFeeEstimator", func() {
				fee := rand.Intn(100)
				estimator := btc.NewFixFeeEstimator(fee)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())
				Expect(fees.Minimum).Should(Equal(fee))
				Expect(fees.Economy).Should(Equal(fee))
				Expect(fees.Low).Should(Equal(fee))
				Expect(fees.Medium).Should(Equal(fee))
				Expect(fees.High).Should(Equal(fee))
			})
		})
	})

	Context("estimate transaction fees", func() {
		It("should return a proper estimate of tx size for a simple transfer", func() {
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

			By("Estimate the tx size before signing")
			extraBaseSize := btc.EstimateUtxoSize(pkAddr1, len(transaction.TxIn))
			estimatedSize := btc.EstimateVirtualSize(transaction, extraBaseSize, 0)

			By("Sign and submit the fund tx")
			for i := range transaction.TxIn {
				pkScript, err := txscript.PayToAddrScript(pkAddr1)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(transaction, i, pkScript, txscript.SigHashAll, privKey1, true)
				Expect(err).To(BeNil())
				transaction.TxIn[i].SignatureScript = sigScript
			}
			actualSize := btc.TxVirtualSize(transaction)
			Expect(estimatedSize).Should(BeNumerically(">=", actualSize))
			By(fmt.Sprintf("Txid = %v ,Estimated tx size = %v, actual tx size = %v", transaction.TxHash().String(), estimatedSize, actualSize))
			Expect(estimatedSize - actualSize).Should(BeNumerically("<=", 10))
		})

		It("should return a proper estimate of tx size when spending multisig utxo", func() {
			By("Initialization keys")
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
			indexer := RegtestIndexer()

			By("Create the multisig script using both public keys")
			multisig, err := btc.MultisigScript(pubKey1.SerializeCompressed(), pubKey2.SerializeCompressed())
			Expect(err).Should(BeNil())
			multisigAddr, err := btc.P2wshAddress(multisig, network)
			Expect(err).Should(BeNil())
			multisigWitnessScript, _ := btc.WitnessScriptHash(multisig)
			Expect(txscript.IsPayToWitnessScriptHash(multisigWitnessScript)).Should(BeTrue())

			By("Funding the multisig address")
			txhash, err := testutil.NigiriFaucet(multisigAddr.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding multisig %v , txid = %v", multisigAddr, txhash))
			time.Sleep(5 * time.Second)

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(context.Background(), multisigAddr)
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

			By("Estimate the tx size before signing")
			extraSegwitSize := btc.RedeemMultisigSigScriptSize * len(transaction.TxIn)
			estimatedSize := btc.EstimateVirtualSize(transaction, 0, extraSegwitSize)

			By("Sign and submit the fund tx")
			for i := range transaction.TxIn {
				err = btc.SignMultisig(transaction, 0, inputs[i].Amount, privKey1, privKey2, txscript.SigHashAll)
				Expect(err).To(BeNil())
			}
			actualSize := btc.TxVirtualSize(transaction)
			Expect(estimatedSize).Should(BeNumerically(">=", actualSize))
			By(fmt.Sprintf("Txid = %v ,Estimated tx size = %v, actual tx size = %v", transaction.TxHash().String(), estimatedSize, actualSize))
			Expect(estimatedSize - actualSize).Should(BeNumerically("<=", 10))
			err = indexer.SubmitTx(context.Background(), transaction)
			Expect(err).To(BeNil())
		})

		It("should return a proper estimate of tx size when spending htlc utxo", func() {
			By("Initialization keys")
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
			indexer := RegtestIndexer()

			By("Create the multisig script using both public keys")
			secret := testutil.RandomSecret()
			secretHash := sha256.Sum256(secret)
			waitTime := int64(6)
			htlc, err := btc.HtlcScript(btcutil.Hash160(pubKey1.SerializeCompressed()), btcutil.Hash160(pubKey2.SerializeCompressed()), secretHash[:], waitTime)
			Expect(err).To(BeNil())
			htlcAddr, err := btc.P2wshAddress(htlc, network)
			Expect(err).To(BeNil())

			By("Funding the multisig address")
			txhash, err := testutil.NigiriFaucet(htlcAddr.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding multisig %v , txid = %v", htlcAddr, txhash))
			time.Sleep(5 * time.Second)

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(context.Background(), htlcAddr)
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

			By("Estimate the tx size before signing")
			extraSegwitSize := btc.RedeemHtlcRedeemSigScriptSize(len(secret)) * len(transaction.TxIn)
			estimatedSize := btc.EstimateVirtualSize(transaction, 0, extraSegwitSize)

			By("Sign and submit the fund tx")
			for i := range transaction.TxIn {
				err = btc.SignHtlcScript(transaction, pubKey1.SerializeCompressed(), pubKey2.SerializeCompressed(), secret, secretHash[:], i, inputs[i].Amount, waitTime, privKey2)
				Expect(err).To(BeNil())
			}
			actualSize := btc.TxVirtualSize(transaction)
			Expect(estimatedSize).Should(BeNumerically(">=", actualSize))
			By(fmt.Sprintf("Txid = %v ,Estimated tx size = %v, actual tx size = %v", transaction.TxHash().String(), estimatedSize, actualSize))
			Expect(estimatedSize - actualSize).Should(BeNumerically("<=", 10))
			err = indexer.SubmitTx(context.Background(), transaction)
			Expect(err).To(BeNil())
		})

		It("should return a proper estimate of tx size when refunding htlc utxo", func() {
			By("Initialization keys")
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
			indexer := RegtestIndexer()

			By("Create the multisig script using both public keys")
			secret := testutil.RandomSecret()
			secretHash := sha256.Sum256(secret)
			waitTime := int64(6)
			htlc, err := btc.HtlcScript(btcutil.Hash160(pubKey1.SerializeCompressed()), btcutil.Hash160(pubKey2.SerializeCompressed()), secretHash[:], waitTime)
			Expect(err).To(BeNil())
			htlcAddr, err := btc.P2wshAddress(htlc, network)
			Expect(err).To(BeNil())

			By("Funding the multisig address")
			txhash, err := testutil.NigiriFaucet(htlcAddr.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding multisig %v , txid = %v", htlcAddr, txhash))
			time.Sleep(5 * time.Second)

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(context.Background(), htlcAddr)
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

			By("Estimate the tx size before signing")
			extraSegwitSize := btc.RedeemHtlcRefundSigScriptSize * len(transaction.TxIn)
			estimatedSize := btc.EstimateVirtualSize(transaction, 0, extraSegwitSize)

			By("Mine some blocks")
			for i := int64(0); i <= waitTime; i++ {
				Expect(testutil.NigiriNewBlock()).Should(Succeed())
			}
			time.Sleep(5 * time.Second)

			By("Sign and submit the fund tx")
			for i := range transaction.TxIn {
				transaction.TxIn[i].Sequence = uint32(waitTime)
				err = btc.SignHtlcScript(transaction, pubKey1.SerializeCompressed(), pubKey2.SerializeCompressed(), nil, secretHash[:], i, inputs[i].Amount, waitTime, privKey1)
				Expect(err).To(BeNil())
			}
			actualSize := btc.TxVirtualSize(transaction)
			Expect(estimatedSize).Should(BeNumerically(">=", actualSize))
			By(fmt.Sprintf("Txid = %v ,Estimated tx size = %v, actual tx size = %v", transaction.TxHash().String(), estimatedSize, actualSize))
			Expect(estimatedSize - actualSize).Should(BeNumerically("<=", 10))

			buffer := bytes.NewBuffer([]byte{})
			err = transaction.Serialize(buffer)
			if err != nil {
				panic(err)
			}
			err = indexer.SubmitTx(context.Background(), transaction)
			Expect(err).To(BeNil())
		})
	})
})
