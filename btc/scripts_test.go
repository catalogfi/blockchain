package btc_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/catalogfi/blockchain/testutil"
	"github.com/fatih/color"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bitcoin scripts", func() {
	Context("Multisig script", func() {
		It("should create a multisig script and can be used by user", func(ctx context.Context) {
			By("Initialization (Update these fields if testing on testnet/mainnet)")
			network := &chaincfg.RegressionNetParams
			privKey1, p2pkhAddr1, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())
			privKey2, p2pkhAddr2, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())
			pubKey1, pubKey2 := privKey1.PubKey(), privKey2.PubKey()
			client, err := btctest.RegtestClient()
			Expect(err).To(BeNil())
			indexer := btctest.RegtestIndexer()

			By("Funding the addresses")
			txhash1, err := testutil.NigiriFaucet(p2pkhAddr1.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address1 %v , txid = %v", p2pkhAddr1.EncodeAddress(), txhash1))
			txhash2, err := testutil.NigiriFaucet(p2pkhAddr2.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address2 %v , txid = %v", p2pkhAddr2.EncodeAddress(), txhash2))
			time.Sleep(5 * time.Second)

			By("Create the multisig script using both public keys")
			script, err := btc.MultisigScript(pubKey1.SerializeCompressed(), pubKey2.SerializeCompressed())
			Expect(err).Should(BeNil())
			addr, err := btc.P2wshAddress(script, network)
			Expect(err).Should(BeNil())
			witnessScript, _ := btc.WitnessScriptHash(script)
			Expect(txscript.IsPayToWitnessScriptHash(witnessScript)).Should(BeTrue())

			By("Construct the funding tx (p2pkhAddr1 -> multi-sig)")
			utxos, err := indexer.GetUTXOs(context.Background(), p2pkhAddr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e5), 10
			fundingRecipients := []btc.Recipient{
				{
					To:     addr.EncodeAddress(),
					Amount: amount,
				},
			}
			fundingTx, err := btc.BuildTransaction(feeRate, network, btc.NewRawInputs(), utxos, fundingRecipients, btc.P2pkhUpdater, p2pkhAddr1)
			Expect(err).To(BeNil())

			By("Sign and submit the funding tx")
			for i := range fundingTx.TxIn {
				pkScript, err := txscript.PayToAddrScript(p2pkhAddr1)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(fundingTx, i, pkScript, txscript.SigHashAll, privKey1, true)
				Expect(err).To(BeNil())
				fundingTx.TxIn[i].SignatureScript = sigScript
			}

			Expect(client.SubmitTx(ctx, fundingTx)).Should(Succeed())
			By(fmt.Sprintf("Funding tx hash = %v", color.YellowString(fundingTx.TxHash().String())))
			time.Sleep(time.Second)

			By("Redeem the funds in the multi-sig")
			redeemInput := []btc.UTXO{
				{
					TxID:   fundingTx.TxHash().String(),
					Vout:   0,
					Amount: amount,
				},
			}
			redeemRecipients := []btc.Recipient{
				{
					To:     p2pkhAddr2.EncodeAddress(),
					Amount: 0,
				},
			}
			rawInputs := btc.RawInputs{
				VIN:        redeemInput,
				BaseSize:   0,
				SegwitSize: btc.RedeemMultisigSigScriptSize,
			}
			redeemTx, err := btc.BuildTransaction(feeRate, network, rawInputs, nil, redeemRecipients, nil, nil)
			Expect(err).To(BeNil())

			By("Sign and submit the redeem tx")
			outpoints := map[wire.OutPoint]*wire.TxOut{
				redeemTx.TxIn[0].PreviousOutPoint: &wire.TxOut{
					Value: amount,
				},
			}
			fetcher := txscript.NewMultiPrevOutFetcher(outpoints)
			for i := range redeemInput {
				sig1, err := txscript.RawTxInWitnessSignature(redeemTx, txscript.NewTxSigHashes(redeemTx, fetcher), i, amount, script, txscript.SigHashAll, privKey1)
				Expect(err).To(BeNil())
				sig2, err := txscript.RawTxInWitnessSignature(redeemTx, txscript.NewTxSigHashes(redeemTx, fetcher), i, amount, script, txscript.SigHashAll, privKey2)
				Expect(err).To(BeNil())
				redeemTx.TxIn[i].Witness = btc.MultisigWitness(script, sig1, sig2)
			}
			err = client.SubmitTx(ctx, redeemTx)
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Redeem tx hash = %v", color.YellowString(redeemTx.TxHash().String())))
		})
	})

	Context("HTLC script", func() {
		It("should create a HTLC script and redeem it using secret", func(ctx context.Context) {
			By("Initialization (Update these fields if testing on testnet/mainnet)")
			network := &chaincfg.RegressionNetParams
			privKey1, p2pkhAddr1, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())
			privKey2, p2pkhAddr2, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())
			pubKey1, pubKey2 := privKey1.PubKey(), privKey2.PubKey()
			client, err := btctest.RegtestClient()
			Expect(err).To(BeNil())
			indexer := btctest.RegtestIndexer()

			By("Funding the addresses")
			txhash1, err := testutil.NigiriFaucet(p2pkhAddr1.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address1 %v , txid = %v", p2pkhAddr1.EncodeAddress(), txhash1))
			txhash2, err := testutil.NigiriFaucet(p2pkhAddr2.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address2 %v , txid = %v", p2pkhAddr2.EncodeAddress(), txhash2))
			time.Sleep(5 * time.Second)

			By("Create the multi-sig address using both public keys")
			multisigScript, err := btc.MultisigScript(pubKey1.SerializeCompressed(), pubKey2.SerializeCompressed())
			Expect(err).Should(BeNil())
			multisigAddr, err := btc.P2wshAddress(multisigScript, network)
			Expect(err).Should(BeNil())
			multisigWitnessScript, err := btc.WitnessScriptHash(multisigScript)
			Expect(err).Should(BeNil())
			Expect(txscript.IsPayToWitnessScriptHash(multisigWitnessScript)).Should(BeTrue())

			By("Construct the funding tx (p2pkhAddr1 -> multi-sig)")
			utxos, err := indexer.GetUTXOs(ctx, p2pkhAddr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e5), 10
			fundingRecipients := []btc.Recipient{
				{
					To:     multisigAddr.EncodeAddress(),
					Amount: amount,
				},
			}
			fundingTx, err := btc.BuildTransaction(feeRate, network, btc.NewRawInputs(), utxos, fundingRecipients, btc.P2pkhUpdater, p2pkhAddr1)
			Expect(err).To(BeNil())

			By("Sign and submit the funding tx")
			for i := range fundingTx.TxIn {
				pkScript, err := txscript.PayToAddrScript(p2pkhAddr1)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(fundingTx, i, pkScript, txscript.SigHashAll, privKey1, true)
				Expect(err).To(BeNil())
				fundingTx.TxIn[i].SignatureScript = sigScript
			}
			Expect(client.SubmitTx(ctx, fundingTx)).Should(Succeed())
			By(fmt.Sprintf("Funding tx hash = %v", color.YellowString(fundingTx.TxHash().String())))
			time.Sleep(time.Second)

			By("Construct a tx which transferring funds from multisig to the HTLC")
			secret := testutil.RandomSecret()
			secretHash := sha256.Sum256(secret)
			waitTime := int64(6)
			// pk1 can redeem the funds after expiry
			// pk2 can redeem the funds if knowing the secret
			htlcScript, err := btc.HtlcScript(btcutil.Hash160(pubKey1.SerializeCompressed()), btcutil.Hash160(pubKey2.SerializeCompressed()), secretHash[:], waitTime)
			Expect(err).To(BeNil())
			Expect(btc.IsHtlc(htlcScript)).Should(BeTrue())
			inputUtxo := btc.UTXO{
				TxID:   fundingTx.TxHash().String(),
				Vout:   0,
				Amount: amount,
			}
			transferTx := wire.NewMsgTx(btc.DefaultBTCVersion)
			txID, err := chainhash.NewHashFromStr(inputUtxo.TxID)
			Expect(err).To(BeNil())
			transferTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(txID, inputUtxo.Vout), nil, nil))
			payScript, err := btc.WitnessScriptHash(htlcScript)
			Expect(err).To(BeNil())
			transferTx.AddTxOut(wire.NewTxOut(inputUtxo.Amount, payScript))

			By("Sign the tx with sighash single")
			outpoints := map[wire.OutPoint]*wire.TxOut{
				transferTx.TxIn[0].PreviousOutPoint: &wire.TxOut{
					Value: amount,
				},
			}
			fetcher := txscript.NewMultiPrevOutFetcher(outpoints)
			sig1, err := txscript.RawTxInWitnessSignature(transferTx, txscript.NewTxSigHashes(transferTx, fetcher), 0, amount, multisigScript, btc.SigHashSingleAnyoneCanPay, privKey1)
			Expect(err).To(BeNil())
			sig2, err := txscript.RawTxInWitnessSignature(transferTx, txscript.NewTxSigHashes(transferTx, fetcher), 0, amount, multisigScript, btc.SigHashSingleAnyoneCanPay, privKey2)
			Expect(err).To(BeNil())
			transferTx.TxIn[0].Witness = btc.MultisigWitness(multisigScript, sig1, sig2)

			By("Add more input to cover the transaction fees")
			fundingTxHash := fundingTx.TxHash()
			fee := int64(1000)
			transferTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&fundingTxHash, 1), nil, nil))
			transferTx.AddTxOut(wire.NewTxOut(fundingTx.TxOut[1].Value-fee, fundingTx.TxOut[1].PkScript))

			By("Sign the new utxo and submit the transaction")
			pkScript, err := txscript.PayToAddrScript(p2pkhAddr1)
			Expect(err).To(BeNil())
			sigScript, err := txscript.SignatureScript(transferTx, 1, pkScript, txscript.SigHashAll, privKey1, true)
			Expect(err).To(BeNil())
			transferTx.TxIn[1].SignatureScript = sigScript
			Expect(client.SubmitTx(ctx, transferTx)).Should(Succeed())
			By(fmt.Sprintf("Refund tx hash = %v", color.YellowString(transferTx.TxHash().String())))
			time.Sleep(2 * time.Second)

			By("Spent the funds in the htlc")
			htlcSpendInputs := []btc.UTXO{
				{
					TxID:   transferTx.TxHash().String(),
					Vout:   0,
					Amount: amount,
				},
			}
			htlcSpendRecipients := []btc.Recipient{
				{
					To:     p2pkhAddr2.EncodeAddress(),
					Amount: 0,
				},
			}
			rawInputs := btc.RawInputs{
				VIN:        htlcSpendInputs,
				BaseSize:   0,
				SegwitSize: btc.RedeemHtlcRedeemSigScriptSize(len(secret)),
			}
			htlcSpendTx, err := btc.BuildTransaction(feeRate, network, rawInputs, nil, htlcSpendRecipients, nil, nil)
			Expect(err).To(BeNil())

			By("Sign and submit the tx")
			fetcher.AddPrevOut(htlcSpendTx.TxIn[0].PreviousOutPoint, &wire.TxOut{Value: amount})
			sig, err := txscript.RawTxInWitnessSignature(htlcSpendTx, txscript.NewTxSigHashes(htlcSpendTx, fetcher), 0, amount, htlcScript, txscript.SigHashAll, privKey2)
			Expect(err).To(BeNil())

			htlcSpendTx.TxIn[0].Witness = btc.HtlcWitness(htlcScript, privKey2.PubKey().SerializeCompressed(), sig, secret)
			err = client.SubmitTx(ctx, htlcSpendTx)
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Htlc SpendTx tx hash = %v", color.YellowString(htlcSpendTx.TxHash().String())))
		})

		It("should create a HTLC script and refund from it after timelock", func(ctx context.Context) {
			By("Initialization (Update these fields if testing on testnet/mainnet)")
			network := &chaincfg.RegressionNetParams
			privKey1, p2pkhAddr1, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())
			privKey2, p2pkhAddr2, err := btctest.NewBtcKey(network)
			Expect(err).To(BeNil())
			pubKey1, pubKey2 := privKey1.PubKey(), privKey2.PubKey()
			client, err := btctest.RegtestClient()
			Expect(err).To(BeNil())
			indexer := btctest.RegtestIndexer()

			By("Funding the addresses")
			txhash1, err := testutil.NigiriFaucet(p2pkhAddr1.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address1 %v , txid = %v", p2pkhAddr1.EncodeAddress(), txhash1))
			txhash2, err := testutil.NigiriFaucet(p2pkhAddr2.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address2 %v , txid = %v", p2pkhAddr2.EncodeAddress(), txhash2))
			time.Sleep(5 * time.Second)

			By("Create the multi-sig address using both public keys")
			multisigScript, err := btc.MultisigScript(pubKey1.SerializeCompressed(), pubKey2.SerializeCompressed())
			Expect(err).Should(BeNil())
			multisigAddr, err := btc.P2wshAddress(multisigScript, network)
			Expect(err).Should(BeNil())
			multisigWitnessScript, err := btc.WitnessScriptHash(multisigScript)
			Expect(err).Should(BeNil())
			Expect(txscript.IsPayToWitnessScriptHash(multisigWitnessScript)).Should(BeTrue())

			By("Construct the funding tx (p2pkhAddr1 -> multi-sig)")
			utxos, err := indexer.GetUTXOs(ctx, p2pkhAddr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e5), 10
			fundingRecipients := []btc.Recipient{
				{
					To:     multisigAddr.EncodeAddress(),
					Amount: amount,
				},
			}
			fundingTx, err := btc.BuildTransaction(feeRate, network, btc.NewRawInputs(), utxos, fundingRecipients, btc.P2pkhUpdater, p2pkhAddr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			for i := range fundingTx.TxIn {
				pkScript, err := txscript.PayToAddrScript(p2pkhAddr1)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(fundingTx, i, pkScript, txscript.SigHashAll, privKey1, true)
				Expect(err).To(BeNil())
				fundingTx.TxIn[i].SignatureScript = sigScript
			}
			Expect(client.SubmitTx(ctx, fundingTx)).Should(Succeed())
			By(fmt.Sprintf("Funding tx hash = %v", color.YellowString(fundingTx.TxHash().String())))
			time.Sleep(time.Second)

			By("Construct a tx which transferring funds from multisig to the HTLC")
			secret := testutil.RandomSecret()
			secretHash := sha256.Sum256(secret)
			waitTime := int64(6)
			// pk1 can redeem the funds after expiry
			// pk2 can redeem the funds if knowing the secret
			htlcScript, err := btc.HtlcScript(btcutil.Hash160(pubKey1.SerializeCompressed()), btcutil.Hash160(pubKey2.SerializeCompressed()), secretHash[:], waitTime)
			Expect(err).To(BeNil())
			Expect(btc.IsHtlc(htlcScript)).Should(BeTrue())
			inputUtxo := btc.UTXO{
				TxID:   fundingTx.TxHash().String(),
				Vout:   0,
				Amount: amount,
			}
			transferTx := wire.NewMsgTx(btc.DefaultBTCVersion)
			txID, err := chainhash.NewHashFromStr(inputUtxo.TxID)
			Expect(err).To(BeNil())
			transferTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(txID, inputUtxo.Vout), nil, nil))
			payScript, err := btc.WitnessScriptHash(htlcScript)
			Expect(err).To(BeNil())
			transferTx.AddTxOut(wire.NewTxOut(inputUtxo.Amount, payScript))

			By("Sign the tx with sighash single")
			outpoints := map[wire.OutPoint]*wire.TxOut{
				transferTx.TxIn[0].PreviousOutPoint: &wire.TxOut{
					Value: amount,
				},
			}
			fetcher := txscript.NewMultiPrevOutFetcher(outpoints)
			sig1, err := txscript.RawTxInWitnessSignature(transferTx, txscript.NewTxSigHashes(transferTx, fetcher), 0, amount, multisigScript, btc.SigHashSingleAnyoneCanPay, privKey1)
			Expect(err).To(BeNil())
			sig2, err := txscript.RawTxInWitnessSignature(transferTx, txscript.NewTxSigHashes(transferTx, fetcher), 0, amount, multisigScript, btc.SigHashSingleAnyoneCanPay, privKey2)
			Expect(err).To(BeNil())
			transferTx.TxIn[0].Witness = btc.MultisigWitness(multisigScript, sig1, sig2)

			By("Add more input to cover the transaction fees")
			fundingTxHash := fundingTx.TxHash()
			fee := int64(1000)
			transferTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&fundingTxHash, 1), nil, nil))
			transferTx.AddTxOut(wire.NewTxOut(fundingTx.TxOut[1].Value-fee, fundingTx.TxOut[1].PkScript))

			By("Sign the new utxo and submit the transaction")
			pkScript, err := txscript.PayToAddrScript(p2pkhAddr1)
			Expect(err).To(BeNil())
			sigScript, err := txscript.SignatureScript(transferTx, 1, pkScript, txscript.SigHashAll, privKey1, true)
			Expect(err).To(BeNil())
			transferTx.TxIn[1].SignatureScript = sigScript
			Expect(client.SubmitTx(ctx, transferTx)).Should(Succeed())
			By(fmt.Sprintf("Refund tx hash = %v", color.YellowString(transferTx.TxHash().String())))
			time.Sleep(2 * time.Second)

			By("Mine some blocks")
			for i := int64(0); i < waitTime-1; i++ {
				Expect(testutil.NigiriNewBlock()).Should(Succeed())
			}
			time.Sleep(time.Second)

			By("Construct the tx which spends the funds in htlc")
			htlcSpendInput := []btc.UTXO{
				{
					TxID:   transferTx.TxHash().String(),
					Vout:   0,
					Amount: amount,
				},
			}
			htlcSpendRecipients := []btc.Recipient{
				{
					To:     p2pkhAddr2.EncodeAddress(),
					Amount: 0,
				},
			}
			rawInputs := btc.RawInputs{
				VIN:        htlcSpendInput,
				BaseSize:   0,
				SegwitSize: btc.RedeemHtlcRefundSigScriptSize,
			}
			htlcSpendTx, err := btc.BuildTransaction(feeRate, network, rawInputs, nil, htlcSpendRecipients, nil, nil)
			Expect(err).To(BeNil())

			By("Sign the htlc spend tx")
			htlcSpendTx.TxIn[0].Sequence = uint32(waitTime)
			fetcher.AddPrevOut(htlcSpendTx.TxIn[0].PreviousOutPoint, &wire.TxOut{Value: amount})
			sig, err := txscript.RawTxInWitnessSignature(htlcSpendTx, txscript.NewTxSigHashes(htlcSpendTx, fetcher), 0, amount, htlcScript, txscript.SigHashAll, privKey1)
			Expect(err).To(BeNil())
			htlcSpendTx.TxIn[0].Witness = btc.HtlcWitness(htlcScript, privKey1.PubKey().SerializeCompressed(), sig, nil)

			By("You shouldn't be able to spend the htlc tx since it doesnt have enough confirmations")
			htlcSpendTxHash := transferTx.TxHash()
			rawRefund, err := client.GetRawTransaction(ctx, &htlcSpendTxHash)
			Expect(err).To(BeNil())
			Expect(rawRefund.Confirmations).Should(Equal(uint64(waitTime) - 1))
			err = client.SubmitTx(ctx, htlcSpendTx)
			Expect(err.Error()).Should(ContainSubstring("non-BIP68-final"))

			By("Mine a new block and we should be able to submit the htlc spend tx")
			Expect(testutil.NigiriNewBlock()).Should(Succeed())
			time.Sleep(time.Second)

			err = client.SubmitTx(ctx, htlcSpendTx)
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Htlc SpendTx tx hash = %v", color.YellowString(htlcSpendTx.TxHash().String())))
		})
	})
})
