package btc_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/catalogfi/multichain/btc"
	"github.com/catalogfi/multichain/testutil"
	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bitcoin scripts", func() {
	It("should create a multisig script and redeem from it", func() {
		By("Initialization (Update these fields if testing on testnet/mainnet)")
		network := &chaincfg.RegressionNetParams
		privKey1, pubKey1, pkAddr1, err := testutil.ParseKeys("PRIV_KEY_1", network)
		Expect(err).To(BeNil())
		privKey2, pubKey2, pkAddr2, err := testutil.ParseKeys("PRIV_KEY_2", network)
		Expect(err).To(BeNil())
		client, err := RegtestClient()
		Expect(err).To(BeNil())
		indexer := RegtestIndexer()

		By("funding the addresses")
		txhash1, err := testutil.NigiriFaucet(pkAddr1.EncodeAddress())
		Expect(err).To(BeNil())
		By(fmt.Sprintf("Funding address1 %v , txid = %v", pkAddr1.EncodeAddress(), txhash1))
		txhash2, err := testutil.NigiriFaucet(pkAddr2.EncodeAddress())
		Expect(err).To(BeNil())
		By(fmt.Sprintf("Funding address2 %v , txid = %v", pkAddr2.EncodeAddress(), txhash2))

		By("Create the multisig script using both public keys")
		walletScript, err := btc.MultisigScript(pubKey1.SerializeCompressed(), pubKey2.SerializeCompressed())
		Expect(err).Should(BeNil())
		walletAddr, err := btc.P2wshAddress(walletScript, network)
		Expect(err).Should(BeNil())
		walletWitnessScript, _ := btc.WitnessScriptHash(walletScript)
		Expect(txscript.IsPayToWitnessScriptHash(walletWitnessScript)).Should(BeTrue())

		By("Construct the funding tx (pk1 -> instantWallet)")
		utxos, err := indexer.GetUTXOs(context.Background(), pkAddr1)
		Expect(err).To(BeNil())
		amount, fee := int64(1e5), int64(500)
		fundingInputs, err := btc.PickUTXOs(utxos, amount, fee)
		Expect(err).To(BeNil())
		fundingRecipients := []btc.Recipient{
			{
				To:     walletAddr,
				Amount: amount,
			},
		}

		fundingTx, _, err := btc.BuildTx(network, fundingInputs, fundingRecipients, fee, pkAddr1)
		Expect(err).To(BeNil())

		By("Sign and submit the fund tx")
		for i := range fundingTx.TxIn {
			// We assume the pkScript is P2PKH
			pkScript, err := txscript.PayToAddrScript(pkAddr1)
			Expect(err).To(BeNil())

			sigScript, err := txscript.SignatureScript(fundingTx, i, pkScript, txscript.SigHashAll, privKey1, true)
			Expect(err).To(BeNil())
			fundingTx.TxIn[i].SignatureScript = sigScript
		}

		Expect(client.SubmitTx(fundingTx)).Should(Succeed())
		By(fmt.Sprintf("Funding tx hash = %v", color.YellowString(fundingTx.TxHash().String())))
		time.Sleep(time.Second)

		By("Redeem the funds in the instant wallet")
		redeemInput := []btc.UTXO{
			{
				TxID:   fundingTx.TxHash().String(),
				Vout:   0,
				Amount: amount,
			},
		}
		redeemRecipients := []btc.Recipient{
			{
				To:     pkAddr2.EncodeAddress(),
				Amount: amount - fee,
			},
		}
		redeemTx, hasChange, err := btc.BuildTx(network, redeemInput, redeemRecipients, fee, pkAddr1)
		Expect(err).To(BeNil())
		Expect(hasChange).Should(BeFalse())

		By("Sign and submit the redeem tx")
		fetcher := txscript.NewCannedPrevOutputFetcher(walletScript, fundingTx.TxOut[0].Value)
		sig1, err := txscript.RawTxInWitnessSignature(redeemTx, txscript.NewTxSigHashes(redeemTx, fetcher), 0, fundingTx.TxOut[0].Value, walletScript, txscript.SigHashAll, privKey1)
		Expect(err).To(BeNil())
		sig2, err := txscript.RawTxInWitnessSignature(redeemTx, txscript.NewTxSigHashes(redeemTx, fetcher), 0, fundingTx.TxOut[0].Value, walletScript, txscript.SigHashAll, privKey2)
		Expect(err).To(BeNil())
		err = btc.SetMultisigWitness(redeemTx, pubKey1.SerializeCompressed(), sig1, pubKey2.SerializeCompressed(), sig2)
		Expect(err).To(BeNil())
		err = client.SubmitTx(redeemTx)
		Expect(err).To(BeNil())
		By(fmt.Sprintf("Redeem tx hash = %v", color.YellowString(redeemTx.TxHash().String())))
	})

	It("should create a HTLC script and redeem it using secret", func() {
		By("Initialization")
		network := &chaincfg.RegressionNetParams
		privKey1, pubKey1, pkAddr1, err := testutil.ParseKeys("PRIV_KEY_1", network)
		Expect(err).To(BeNil())
		privKey2, pubKey2, pkAddr2, err := testutil.ParseKeys("PRIV_KEY_2", network)
		Expect(err).To(BeNil())
		client, err := RegtestClient()
		Expect(err).To(BeNil())
		indexer := RegtestIndexer()

		By("funding the addresses")
		txhash1, err := testutil.NigiriFaucet(pkAddr1.EncodeAddress())
		Expect(err).To(BeNil())
		By(fmt.Sprintf("Funding address1 %v , txid = %v", pkAddr1.EncodeAddress(), txhash1))
		txhash2, err := testutil.NigiriFaucet(pkAddr2.EncodeAddress())
		Expect(err).To(BeNil())
		By(fmt.Sprintf("Funding address2 %v , txid = %v", pkAddr2.EncodeAddress(), txhash2))

		By("Create the multi sig address using both public keys")
		walletScript, err := btc.MultisigScript(pubKey1.SerializeCompressed(), pubKey2.SerializeCompressed())
		Expect(err).Should(BeNil())
		walletAddr, err := btc.P2wshAddress(walletScript, network)
		Expect(err).Should(BeNil())
		walletWitnessScript, err := btc.WitnessScriptHash(walletScript)
		Expect(err).Should(BeNil())
		Expect(txscript.IsPayToWitnessScriptHash(walletWitnessScript)).Should(BeTrue())

		By("Construct the funding tx (pk1 -> multi sig)")
		utxos, err := indexer.GetUTXOs(context.Background(), pkAddr1)
		Expect(err).To(BeNil())
		amount, fee := int64(1e5), int64(500)
		fundingInputs, err := btc.PickUTXOs(utxos, amount, fee)
		Expect(err).To(BeNil())
		fundingRecipients := []btc.Recipient{
			{
				To:     walletAddr,
				Amount: amount,
			},
		}
		fundingTx, _, err := btc.BuildTx(network, fundingInputs, fundingRecipients, fee, pkAddr1)
		Expect(err).To(BeNil())

		By("Sign and submit the fund tx")
		for i := range fundingTx.TxIn {
			// We assume the pkScript is P2PKH
			pkScript, err := txscript.PayToAddrScript(pkAddr1)
			Expect(err).To(BeNil())

			sigScript, err := txscript.SignatureScript(fundingTx, i, pkScript, txscript.SigHashAll, privKey1, true)
			Expect(err).To(BeNil())
			fundingTx.TxIn[i].SignatureScript = sigScript
		}
		Expect(client.SubmitTx(fundingTx)).Should(Succeed())
		By(fmt.Sprintf("Funding tx hash = %v", color.YellowString(fundingTx.TxHash().String())))
		time.Sleep(time.Second)

		By("Build refund transaction")
		secret := testutil.RandomSecret()
		secretHash := sha256.Sum256(secret)
		waitTime := int64(6)
		// pk1 can redeem the funds after expiry
		// pk2 can redeem the funds if knowing the secret
		refundScript, err := btc.HtlcScript(btcutil.Hash160(pubKey1.SerializeCompressed()), btcutil.Hash160(pubKey2.SerializeCompressed()), secretHash[:], waitTime)
		Expect(err).To(BeNil())
		inputUtxo := btc.UTXO{
			TxID:   fundingTx.TxHash().String(),
			Vout:   0,
			Amount: amount,
		}
		refundTx, err := btc.NewRefundTx(inputUtxo, refundScript)
		Expect(err).To(BeNil())

		By("System sign the tx with sighash single")
		fetcher := txscript.NewCannedPrevOutputFetcher(walletScript, fundingTx.TxOut[0].Value)
		sig1, err := txscript.RawTxInWitnessSignature(refundTx, txscript.NewTxSigHashes(refundTx, fetcher), 0, fundingTx.TxOut[0].Value, walletScript, btc.SigHashSingleAnyoneCanPay, privKey1)
		Expect(err).To(BeNil())
		sig2, err := txscript.RawTxInWitnessSignature(refundTx, txscript.NewTxSigHashes(refundTx, fetcher), 0, fundingTx.TxOut[0].Value, walletScript, btc.SigHashSingleAnyoneCanPay, privKey2)
		Expect(err).To(BeNil())

		By("Fill the tx to cover the fees")
		fundingTxHash := fundingTx.TxHash()
		refundTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&fundingTxHash, 1), nil, nil))
		refundTx.AddTxOut(wire.NewTxOut(fundingTx.TxOut[1].Value-fee, fundingTx.TxOut[1].PkScript))

		By("Sign and submit refund transaction")
		err = btc.SetMultisigWitness(refundTx, pubKey1.SerializeCompressed(), sig1, pubKey2.SerializeCompressed(), sig2)
		Expect(err).To(BeNil())
		pkScript, err := txscript.PayToAddrScript(pkAddr1)
		Expect(err).To(BeNil())
		sigScript, err := txscript.SignatureScript(refundTx, 1, pkScript, txscript.SigHashAll, privKey1, true)
		Expect(err).To(BeNil())
		refundTx.TxIn[1].SignatureScript = sigScript
		Expect(client.SubmitTx(refundTx)).Should(Succeed())
		By(fmt.Sprintf("Refund tx hash = %v", color.YellowString(refundTx.TxHash().String())))

		time.Sleep(2 * time.Second)

		By("Spend the refund using secret")
		refundSpendInput := []btc.UTXO{
			{
				TxID:   refundTx.TxHash().String(),
				Vout:   0,
				Amount: amount,
			},
		}
		refundSpendRecipients := []btc.Recipient{
			{
				To:     pkAddr2.EncodeAddress(),
				Amount: amount - fee,
			},
		}
		refundSpendTx, hasChange, err := btc.BuildTx(network, refundSpendInput, refundSpendRecipients, fee, pkAddr1)
		Expect(err).To(BeNil())
		Expect(hasChange).Should(BeFalse())

		By("Sign and submit the refund spend tx")
		fetcher = txscript.NewCannedPrevOutputFetcher(refundScript, refundTx.TxOut[0].Value)
		refundSpendSig, err := txscript.RawTxInWitnessSignature(refundSpendTx, txscript.NewTxSigHashes(refundSpendTx, fetcher), 0, amount, refundScript, txscript.SigHashAll, privKey2)
		Expect(err).To(BeNil())
		err = btc.SetRefundSpendWitness(refundSpendTx, pubKey1.SerializeCompressed(), pubKey2.SerializeCompressed(), refundSpendSig, secretHash[:], secret, false, waitTime)
		Expect(err).To(BeNil())
		err = client.SubmitTx(refundSpendTx)
		Expect(err).To(BeNil())

		By(fmt.Sprintf("RefundSpendTx tx hash = %v", color.YellowString(refundSpendTx.TxHash().String())))
	})

	It("should create instant wallet and refund from it after timelock", func() {
		By("Initialization")
		network := &chaincfg.RegressionNetParams
		privKey1, pubKey1, pkAddr1, err := testutil.ParseKeys("PRIV_KEY_1", network)
		Expect(err).To(BeNil())
		privKey2, pubKey2, pkAddr2, err := testutil.ParseKeys("PRIV_KEY_2", network)
		Expect(err).To(BeNil())
		client, err := RegtestClient()
		Expect(err).To(BeNil())
		indexer := RegtestIndexer()

		By("funding the addresses")
		txhash1, err := testutil.NigiriFaucet(pkAddr1.EncodeAddress())
		Expect(err).To(BeNil())
		By(fmt.Sprintf("Funding address1 %v , txid = %v", pkAddr1.EncodeAddress(), txhash1))
		txhash2, err := testutil.NigiriFaucet(pkAddr2.EncodeAddress())
		Expect(err).To(BeNil())
		By(fmt.Sprintf("Funding address2 %v , txid = %v", pkAddr2.EncodeAddress(), txhash2))

		By("Create multi sig script using both public keys")
		walletScript, err := btc.MultisigScript(pubKey1.SerializeCompressed(), pubKey2.SerializeCompressed())
		Expect(err).Should(BeNil())
		walletAddr, err := btc.P2wshAddress(walletScript, network)
		Expect(err).Should(BeNil())
		walletWitnessScript, _ := btc.WitnessScriptHash(walletScript)
		Expect(txscript.IsPayToWitnessScriptHash(walletWitnessScript)).Should(BeTrue())

		By("Construct the funding tx (pk1 -> instantWallet)")
		utxos, err := indexer.GetUTXOs(context.Background(), pkAddr1)
		Expect(err).To(BeNil())
		amount, fee := int64(1e5), int64(500)
		fundingInputs, err := btc.PickUTXOs(utxos, amount, fee)
		Expect(err).To(BeNil())
		fundingRecipients := []btc.Recipient{
			{
				To:     walletAddr,
				Amount: amount,
			},
		}
		fundingTx, _, err := btc.BuildTx(network, fundingInputs, fundingRecipients, fee, pkAddr1)
		Expect(err).To(BeNil())

		By("Sign and submit the fund tx")
		for i := range fundingTx.TxIn {
			// We assume the pkScript is P2PKH
			pkScript, err := txscript.PayToAddrScript(pkAddr1)
			Expect(err).To(BeNil())

			sigScript, err := txscript.SignatureScript(fundingTx, i, pkScript, txscript.SigHashAll, privKey1, true)
			Expect(err).To(BeNil())
			fundingTx.TxIn[i].SignatureScript = sigScript
		}
		Expect(client.SubmitTx(fundingTx)).Should(Succeed())
		By(fmt.Sprintf("Funding tx hash = %v", color.YellowString(fundingTx.TxHash().String())))
		time.Sleep(time.Second)

		By("Build refund transaction")
		secret := testutil.RandomSecret()
		secretHash := sha256.Sum256(secret)
		waitTime := int64(6)
		// pk1 can redeem the funds after expiry
		// pk2 can redeem the funds if knowing the secret
		refundScript, err := btc.HtlcScript(btcutil.Hash160(pubKey1.SerializeCompressed()), btcutil.Hash160(pubKey2.SerializeCompressed()), secretHash[:], waitTime)
		Expect(err).To(BeNil())
		inputUtxo := btc.UTXO{
			TxID:   fundingTx.TxHash().String(),
			Vout:   0,
			Amount: amount,
		}
		refundTx, err := btc.NewRefundTx(inputUtxo, refundScript)
		Expect(err).To(BeNil())

		By("System sign the tx with sighash single")
		fetcher := txscript.NewCannedPrevOutputFetcher(walletScript, fundingTx.TxOut[0].Value)
		sig1, err := txscript.RawTxInWitnessSignature(refundTx, txscript.NewTxSigHashes(refundTx, fetcher), 0, fundingTx.TxOut[0].Value, walletScript, btc.SigHashSingleAnyoneCanPay, privKey1)
		Expect(err).To(BeNil())
		sig2, err := txscript.RawTxInWitnessSignature(refundTx, txscript.NewTxSigHashes(refundTx, fetcher), 0, fundingTx.TxOut[0].Value, walletScript, btc.SigHashSingleAnyoneCanPay, privKey2)
		Expect(err).To(BeNil())

		By("Fill the tx to cover the fees")
		fundingTxHash := fundingTx.TxHash()
		refundTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&fundingTxHash, 1), nil, nil))
		refundTx.AddTxOut(wire.NewTxOut(fundingTx.TxOut[1].Value-fee, fundingTx.TxOut[1].PkScript))

		By("Sign and submit refund transaction")
		err = btc.SetMultisigWitness(refundTx, pubKey1.SerializeCompressed(), sig1, pubKey2.SerializeCompressed(), sig2)
		Expect(err).To(BeNil())
		pkScript, err := txscript.PayToAddrScript(pkAddr1)
		Expect(err).To(BeNil())
		sigScript, err := txscript.SignatureScript(refundTx, 1, pkScript, txscript.SigHashAll, privKey1, true)
		Expect(err).To(BeNil())
		refundTx.TxIn[1].SignatureScript = sigScript
		Expect(client.SubmitTx(refundTx)).Should(Succeed())
		By(fmt.Sprintf("Refund tx hash = %v", color.YellowString(refundTx.TxHash().String())))

		By("Mine some blocks")
		for i := int64(0); i < waitTime-1; i++ {
			Expect(testutil.NigiriNewBlock()).Should(Succeed())
		}
		time.Sleep(time.Second)

		By("Construct the refund spend tx ")
		refundSpendInput := []btc.UTXO{
			{
				TxID:   refundTx.TxHash().String(),
				Vout:   0,
				Amount: amount,
			},
		}
		refundSpendRecipients := []btc.Recipient{
			{
				To:     pkAddr2.EncodeAddress(),
				Amount: amount - fee,
			},
		}
		refundSpendTx, hasChange, err := btc.BuildTx(network, refundSpendInput, refundSpendRecipients, fee, pkAddr1)
		Expect(err).To(BeNil())
		Expect(hasChange).Should(BeFalse())

		By("Sign the refund spend tx")
		refundSpendTx.TxIn[0].Sequence = uint32(waitTime)
		fetcher = txscript.NewCannedPrevOutputFetcher(refundScript, refundTx.TxOut[0].Value)
		refundSpendSig, err := txscript.RawTxInWitnessSignature(refundSpendTx, txscript.NewTxSigHashes(refundSpendTx, fetcher), 0, amount, refundScript, txscript.SigHashAll, privKey1)
		Expect(err).To(BeNil())
		err = btc.SetRefundSpendWitness(refundSpendTx, pubKey1.SerializeCompressed(), pubKey2.SerializeCompressed(), refundSpendSig, secretHash[:], secret, true, waitTime)
		Expect(err).To(BeNil())

		By("You shouldn't be able to spend the refund tx since it doesnt have enough confirmations")
		refundTxHash := refundTx.TxHash()
		rawRefund, err := client.GetRawTransaction(refundTxHash.CloneBytes())
		Expect(err).To(BeNil())
		Expect(rawRefund.Confirmations).Should(Equal(uint64(waitTime) - 1))
		err = client.SubmitTx(refundSpendTx)
		Expect(err.Error()).Should(ContainSubstring("non-BIP68-final"))

		By("Mine a new block and we should be able to submit the refund spend tx")
		Expect(testutil.NigiriNewBlock()).Should(Succeed())
		time.Sleep(time.Second)
		err = client.SubmitTx(refundSpendTx)
		Expect(err).To(BeNil())
		By(fmt.Sprintf("RefundSpendTx tx hash = %v", color.YellowString(refundSpendTx.TxHash().String())))
	})
})
