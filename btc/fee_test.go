package btc_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/localnet"

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

			It("should return the live data from mempool for testnet", func() {
				estimator := btc.NewMempoolFeeEstimator(&chaincfg.TestNet3Params, btc.MempoolFeeAPITestnet, 15*time.Second)
				fees, err := estimator.FeeSuggestion()
				Expect(err).Should(BeNil())

				Expect(fees.Minimum).Should(BeNumerically(">=", 1))
				Expect(fees.Economy).Should(BeNumerically(">=", 1))
				Expect(fees.Low).Should(BeNumerically(">=", 1))
				Expect(fees.Medium).Should(BeNumerically(">=", 1))
				Expect(fees.High).Should(BeNumerically(">=", 1))
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

			It("should panic if the network is not testnet or mainnet", func() {
				By("Regnet")
				PanicWith(func() {
					_ = btc.NewMempoolFeeEstimator(&chaincfg.RegressionNetParams, "", 15*time.Second)
				})
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
				Expect(fees.Minimum).Should(Equal(1))
				Expect(fees.Economy).Should(Equal(1))
				Expect(fees.Low).Should(Equal(1))
				Expect(fees.Medium).Should(Equal(1))
				Expect(fees.High).Should(Equal(1))

				By("Regnet")
				estimatorRegnet := btc.NewBlockstreamFeeEstimator(&chaincfg.RegressionNetParams, "", 15*time.Second)
				fees, err = estimatorRegnet.FeeSuggestion()
				Expect(err).Should(BeNil())
				Expect(fees.Minimum).Should(Equal(1))
				Expect(fees.Economy).Should(Equal(1))
				Expect(fees.Low).Should(Equal(1))
				Expect(fees.Medium).Should(Equal(1))
				Expect(fees.High).Should(Equal(1))
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
		It("should return a proper estimate of tx size for a simple transfer", func(ctx context.Context) {
			By("Initialization keys")
			key1, addr1, err := localnet.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())
			_, addr2, err := localnet.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())

			By("Funding the addresses")
			_, err = localnet.FundBTC(addr1.EncodeAddress())
			Expect(err).To(BeNil())
			time.Sleep(5 * time.Second)

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(ctx, addr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e5), 10
			recipients := []btc.Recipient{
				{
					To:     addr2.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildTransaction(network, feeRate, btc.NewRawInputs(), utxos, btc.P2pkhUpdater, recipients, addr1)
			Expect(err).To(BeNil())

			By("Estimate the tx size before signing")
			extraBaseSize := txsizes.RedeemP2PKHSigScriptSize * len(transaction.TxIn)
			estimatedSize := btc.EstimateVirtualSize(transaction, extraBaseSize, 0)

			By("Sign and submit the fund tx")
			Expect(btc.SignP2pkhTx(network, key1, transaction))
			actualSize := btc.TxVirtualSize(transaction)
			By(fmt.Sprintf("Txid = %v ,Estimated tx size = %v, actual tx size = %v", transaction.TxHash().String(), estimatedSize, actualSize))
			Expect(estimatedSize).Should(BeNumerically(">=", actualSize))
			Expect(estimatedSize - actualSize).Should(BeNumerically("<=", 10))
		})

		It("should return a proper estimate of tx size when spending multisig utxo", func(ctx context.Context) {
			By("Initialization keys")
			key1, addr1, err := localnet.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())
			key2, addr2, err := localnet.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())

			By("Create the multisig script using both public keys")
			multisig, err := btc.MultisigScript(key1.PubKey().SerializeCompressed(), key2.PubKey().SerializeCompressed())
			Expect(err).Should(BeNil())
			multisigAddr, err := btc.P2wshAddress(multisig, network)
			Expect(err).Should(BeNil())
			multisigWitnessScript, _ := btc.WitnessScriptHash(multisig)
			Expect(txscript.IsPayToWitnessScriptHash(multisigWitnessScript)).Should(BeTrue())

			By("Funding the multisig address")
			_, err = localnet.FundBTC(multisigAddr.EncodeAddress())
			Expect(err).To(BeNil())
			time.Sleep(5 * time.Second)

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(ctx, multisigAddr)
			Expect(err).To(BeNil())
			amounts := map[string]int64{}
			for _, utxo := range utxos {
				amounts[utxo.TxID] = utxo.Amount
			}
			amount, feeRate := int64(1e5), 10
			recipients := []btc.Recipient{
				{
					To:     addr2.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildTransaction(network, feeRate, btc.NewRawInputs(), utxos, btc.MultisigUpdater, recipients, addr1)
			Expect(err).To(BeNil())

			By("Estimate the tx size before signing")
			extraSegwitSize := btc.RedeemMultisigSigScriptSize * len(transaction.TxIn)
			estimatedSize := btc.EstimateVirtualSize(transaction, 0, extraSegwitSize)

			By("Sign and submit the fund tx")
			fetcher := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{})
			for i := range transaction.TxIn {
				inputAmount := amounts[transaction.TxIn[i].PreviousOutPoint.Hash.String()]
				fetcher.AddPrevOut(transaction.TxIn[i].PreviousOutPoint, &wire.TxOut{Value: inputAmount})
				sig1, err := txscript.RawTxInWitnessSignature(transaction, txscript.NewTxSigHashes(transaction, fetcher), i, inputAmount, multisig, txscript.SigHashAll, key1)
				Expect(err).To(BeNil())
				sig2, err := txscript.RawTxInWitnessSignature(transaction, txscript.NewTxSigHashes(transaction, fetcher), i, inputAmount, multisig, txscript.SigHashAll, key2)
				Expect(err).To(BeNil())
				transaction.TxIn[i].Witness = btc.MultisigWitness(multisig, sig1, sig2)
			}
			actualSize := btc.TxVirtualSize(transaction)
			Expect(estimatedSize).Should(BeNumerically(">=", actualSize))
			By(fmt.Sprintf("Txid = %v ,Estimated tx size = %v, actual tx size = %v", transaction.TxHash().String(), estimatedSize, actualSize))
			Expect(estimatedSize - actualSize).Should(BeNumerically("<=", 10))
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
		})

		It("should return a proper estimate of tx size when spending htlc utxo", func(ctx context.Context) {
			By("Initialization keys")
			key1, addr1, err := localnet.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())
			key2, addr2, err := localnet.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())

			By("Create the htlc script using both public keys")
			secret := localnet.RandomSecret()
			secretHash := sha256.Sum256(secret)
			waitTime := int64(6)
			htlc, err := btc.HtlcScript(btcutil.Hash160(key1.PubKey().SerializeCompressed()), btcutil.Hash160(key2.PubKey().SerializeCompressed()), secretHash[:], waitTime)
			Expect(err).To(BeNil())
			htlcAddr, err := btc.P2wshAddress(htlc, network)
			Expect(err).To(BeNil())

			By("Funding the multisig address")
			txhash, err := localnet.FundBTC(htlcAddr.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding multisig %v , txid = %v", htlcAddr, txhash))
			time.Sleep(5 * time.Second)

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(ctx, htlcAddr)
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
					To:     addr2.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildTransaction(network, feeRate, btc.NewRawInputs(), utxos, btc.HtlcUpdater(len(secret)), recipients, addr1)
			Expect(err).To(BeNil())

			By("Estimate the tx size before signing")
			extraSegwitSize := btc.RedeemHtlcRedeemSigScriptSize(len(secret)) * len(transaction.TxIn)
			estimatedSize := btc.EstimateVirtualSize(transaction, 0, extraSegwitSize)

			By("Sign and submit the fund tx")
			for i := range transaction.TxIn {
				sig, err := txscript.RawTxInWitnessSignature(transaction, txscript.NewTxSigHashes(transaction, fetcher), i, utxos[i].Amount, htlc, txscript.SigHashAll, key2)
				Expect(err).To(BeNil())
				transaction.TxIn[i].Witness = btc.HtlcWitness(htlc, key2.PubKey().SerializeCompressed(), sig, secret)
			}
			actualSize := btc.TxVirtualSize(transaction)
			Expect(estimatedSize).Should(BeNumerically(">=", actualSize))
			By(fmt.Sprintf("Txid = %v ,Estimated tx size = %v, actual tx size = %v", transaction.TxHash().String(), estimatedSize, actualSize))
			Expect(estimatedSize - actualSize).Should(BeNumerically("<=", 10))
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
		})

		It("should return a proper estimate of tx size when refunding htlc utxo", func(ctx context.Context) {
			By("Initialization keys")
			key1, addr1, err := localnet.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())
			key2, addr2, err := localnet.NewBtcKey(network, waddrmgr.PubKeyHash)
			Expect(err).To(BeNil())

			By("Create the htlc script using both public keys")
			secret := localnet.RandomSecret()
			secretHash := sha256.Sum256(secret)
			waitTime := int64(6)
			htlc, err := btc.HtlcScript(btcutil.Hash160(key1.PubKey().SerializeCompressed()), btcutil.Hash160(key2.PubKey().SerializeCompressed()), secretHash[:], waitTime)
			Expect(err).To(BeNil())
			htlcAddr, err := btc.P2wshAddress(htlc, network)
			Expect(err).To(BeNil())

			By("Funding the multisig address")
			txhash, err := localnet.FundBTC(htlcAddr.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding multisig %v , txid = %v", htlcAddr, txhash))
			time.Sleep(5 * time.Second)

			By("Build the transaction")
			utxos, err := indexer.GetUTXOs(ctx, htlcAddr)
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
					To:     addr2.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildTransaction(network, feeRate, btc.NewRawInputs(), utxos, btc.HtlcUpdater(0), recipients, addr1)
			Expect(err).To(BeNil())

			By("Estimate the tx size before signing")
			extraSegwitSize := btc.RedeemHtlcRefundSigScriptSize * len(transaction.TxIn)
			estimatedSize := btc.EstimateVirtualSize(transaction, 0, extraSegwitSize)

			By("Mine some blocks")
			for i := int64(0); i <= waitTime; i++ {
				Expect(localnet.MineBTCBlock()).Should(Succeed())
			}
			time.Sleep(5 * time.Second)

			By("Sign and submit the fund tx")
			for i := range transaction.TxIn {
				transaction.TxIn[i].Sequence = uint32(waitTime)
			}
			for i := range transaction.TxIn {
				sig, err := txscript.RawTxInWitnessSignature(transaction, txscript.NewTxSigHashes(transaction, fetcher), i, utxos[i].Amount, htlc, txscript.SigHashAll, key1)
				Expect(err).To(BeNil())
				transaction.TxIn[i].Witness = btc.HtlcWitness(htlc, key1.PubKey().SerializeCompressed(), sig, nil)
			}
			actualSize := btc.TxVirtualSize(transaction)
			Expect(estimatedSize).Should(BeNumerically(">=", actualSize))
			By(fmt.Sprintf("Txid = %v ,Estimated tx size = %v, actual tx size = %v", transaction.TxHash().String(), estimatedSize, actualSize))
			Expect(estimatedSize - actualSize).Should(BeNumerically("<=", 10))
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
		})
	})
	Context("estimate segwit transaction fees", func() {
		It("should return a proper segwit fee for a simple transfer", func(ctx context.Context) {
			By("Initialization Wallet")
			pk, _, err := localnet.NewBtcKey(network, waddrmgr.WitnessPubKey)
			Expect(err).To(BeNil())
			fixedFeeEstimator := btc.NewFixFeeEstimator(10)

			wallet, err := btc.NewSimpleWallet(pk, network, indexer, fixedFeeEstimator, btc.HighFee)
			Expect(err).To(BeNil())

			By("Fund the wallet")
			_, err = localnet.FundBitcoin(wallet.Address().EncodeAddress(), indexer)
			Expect(err).To(BeNil())

			By("Build the request")
			req := []btc.SendRequest{
				{
					To:     wallet.Address(),
					Amount: 1e5,
				},
			}
			txId, err := wallet.Send(ctx, req, nil)

			By("Rebuild submitted tx")
			txHex, err := indexer.GetTxHex(ctx, txId)
			Expect(err).To(BeNil())

			txBytes, err := hex.DecodeString(txHex)
			Expect(err).To(BeNil())

			tx, err := btcutil.NewTxFromBytes(txBytes)
			Expect(err).To(BeNil())

			By("Estimate the fee")
			estimatedFee, err := btc.EstimateSegwitFee(tx.MsgTx(), fixedFeeEstimator, btc.HighFee)
			Expect(err).To(BeNil())

			By("Verifying whether the estimated fee is used in the transaction")
			transaction, err := indexer.GetTx(ctx, txId)
			Expect(err).To(BeNil())

			feeRates, err := fixedFeeEstimator.FeeSuggestion()
			Expect(err).To(BeNil())
			Expect(estimatedFee).Should(Equal(int((transaction.Weight / blockchain.WitnessScaleFactor) * feeRates.High)))

		})
	})
})
