package btc_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/catalogfi/blockchain/testutil"

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
			indexer := btctest.RegtestIndexer()

			By("Funding the addresses")
			txhash1, err := testutil.NigiriFaucet(pkAddr1.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address1 %v , txid = %v", pkAddr1.EncodeAddress(), txhash1))
			time.Sleep(5 * time.Second)

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(context.Background(), pkAddr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e5), 10
			recipients := []btc.Recipient{
				{
					To:     pkAddr2.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildTransaction(feeRate, network, btc.NewRawInputs(), utxos, recipients, btc.P2pkhUpdater, pkAddr1)
			Expect(err).To(BeNil())

			By("Estimate the tx size before signing")
			extraBaseSize := txsizes.RedeemP2PKHSigScriptSize * len(transaction.TxIn)
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
			indexer := btctest.RegtestIndexer()

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
			amounts := map[string]int64{}
			for _, utxo := range utxos {
				amounts[utxo.TxID] = utxo.Amount
			}
			amount, feeRate := int64(1e5), 10
			recipients := []btc.Recipient{
				{
					To:     pkAddr2.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildTransaction(feeRate, network, btc.NewRawInputs(), utxos, recipients, btc.MultisigUpdater, pkAddr1)
			Expect(err).To(BeNil())

			By("Estimate the tx size before signing")
			extraSegwitSize := btc.RedeemMultisigSigScriptSize * len(transaction.TxIn)
			estimatedSize := btc.EstimateVirtualSize(transaction, 0, extraSegwitSize)

			By("Sign and submit the fund tx")
			fetcher := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{})
			for i := range transaction.TxIn {
				inputAmount := amounts[transaction.TxIn[i].PreviousOutPoint.Hash.String()]
				fetcher.AddPrevOut(transaction.TxIn[i].PreviousOutPoint, &wire.TxOut{Value: inputAmount})
				sig1, err := txscript.RawTxInWitnessSignature(transaction, txscript.NewTxSigHashes(transaction, fetcher), i, inputAmount, multisig, txscript.SigHashAll, privKey1)
				Expect(err).To(BeNil())
				sig2, err := txscript.RawTxInWitnessSignature(transaction, txscript.NewTxSigHashes(transaction, fetcher), i, inputAmount, multisig, txscript.SigHashAll, privKey2)
				Expect(err).To(BeNil())
				transaction.TxIn[i].Witness = btc.MultisigWitness(multisig, sig1, sig2)
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
			privKey1, pkAddr1, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())
			privKey2, pkAddr2, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())
			pubKey1, pubKey2 := privKey1.PubKey(), privKey2.PubKey()
			indexer := btctest.RegtestIndexer()

			By("Create the htlc script using both public keys")
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
			fetcher := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{})
			for _, utxo := range utxos {
				hash, err := chainhash.NewHashFromStr(utxo.TxID)
				Expect(err).To(BeNil())

				fetcher.AddPrevOut(wire.OutPoint{
					Hash:  *hash,
					Index: utxo.Vout,
				}, &wire.TxOut{
					Value: utxo.Amount,
				})
			}
			amount, feeRate := int64(1e5), 10
			recipients := []btc.Recipient{
				{
					To:     pkAddr2.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildTransaction(feeRate, network, btc.NewRawInputs(), utxos, recipients, btc.HtlcUpdater(len(secret)), pkAddr1)
			Expect(err).To(BeNil())

			By("Estimate the tx size before signing")
			extraSegwitSize := btc.RedeemHtlcRedeemSigScriptSize(len(secret)) * len(transaction.TxIn)
			estimatedSize := btc.EstimateVirtualSize(transaction, 0, extraSegwitSize)

			By("Sign and submit the fund tx")
			for i := range transaction.TxIn {
				sig, err := txscript.RawTxInWitnessSignature(transaction, txscript.NewTxSigHashes(transaction, fetcher), i, utxos[i].Amount, htlc, txscript.SigHashAll, privKey2)
				Expect(err).To(BeNil())
				transaction.TxIn[i].Witness = btc.HtlcWitness(htlc, privKey2.PubKey().SerializeCompressed(), sig, secret)
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
			privKey1, pkAddr1, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())
			privKey2, pkAddr2, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())
			pubKey1, pubKey2 := privKey1.PubKey(), privKey2.PubKey()
			indexer := btctest.RegtestIndexer()

			By("Create the htlc script using both public keys")
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
			fetcher := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{})
			for _, utxo := range utxos {
				hash, err := chainhash.NewHashFromStr(utxo.TxID)
				Expect(err).To(BeNil())

				fetcher.AddPrevOut(wire.OutPoint{
					Hash:  *hash,
					Index: utxo.Vout,
				}, &wire.TxOut{
					Value: utxo.Amount,
				})
			}
			amount, feeRate := int64(1e5), 10
			recipients := []btc.Recipient{
				{
					To:     pkAddr2.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildTransaction(feeRate, network, btc.NewRawInputs(), utxos, recipients, btc.HtlcUpdater(0), pkAddr1)
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
			}
			for i := range transaction.TxIn {
				sig, err := txscript.RawTxInWitnessSignature(transaction, txscript.NewTxSigHashes(transaction, fetcher), i, utxos[i].Amount, htlc, txscript.SigHashAll, privKey1)
				Expect(err).To(BeNil())
				transaction.TxIn[i].Witness = btc.HtlcWitness(htlc, privKey1.PubKey().SerializeCompressed(), sig, nil)
			}
			actualSize := btc.TxVirtualSize(transaction)
			Expect(estimatedSize).Should(BeNumerically(">=", actualSize))
			By(fmt.Sprintf("Txid = %v ,Estimated tx size = %v, actual tx size = %v", transaction.TxHash().String(), estimatedSize, actualSize))
			Expect(estimatedSize - actualSize).Should(BeNumerically("<=", 10))
			err = indexer.SubmitTx(context.Background(), transaction)
			Expect(err).To(BeNil())
		})
	})
})
