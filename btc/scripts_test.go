package btc_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
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
		It("should create a multisig script and redeem from it", func(ctx context.Context) {
			By("Initialization (Update these fields if testing on testnet/mainnet)")
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
			client, err := btctest.RegtestClient()
			Expect(err).To(BeNil())
			indexer := btctest.RegtestIndexer()

			By("funding the addresses")
			txhash1, err := testutil.NigiriFaucet(pkAddr1.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address1 %v , txid = %v", pkAddr1.EncodeAddress(), txhash1))
			txhash2, err := testutil.NigiriFaucet(pkAddr2.EncodeAddress())
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Funding address2 %v , txid = %v", pkAddr2.EncodeAddress(), txhash2))
			time.Sleep(5 * time.Second)

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
			amount, feeRate := int64(1e5), 10
			fundingRecipients := []btc.Recipient{
				{
					To:     walletAddr.EncodeAddress(),
					Amount: amount,
				},
			}
			fundingTx, err := btc.BuildTransaction(feeRate, network, nil, utxos, fundingRecipients, 0, 0, pkAddr1)
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

			Expect(client.SubmitTx(ctx, fundingTx)).Should(Succeed())
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
					Amount: amount - 2e3,
				},
			}
			redeemTx, err := btc.BuildTransaction(feeRate, network, nil, redeemInput, redeemRecipients, 0, 0, pkAddr1)
			Expect(err).To(BeNil())
			extraSegwitSize := btc.RedeemMultisigSigScriptSize
			estimatedSize := btc.EstimateVirtualSize(redeemTx, 0, extraSegwitSize)

			By("Sign and submit the redeem tx")
			for i := range redeemInput {
				err = btc.SignMultisig(redeemTx, i, redeemInput[i].Amount, privKey1, privKey2, txscript.SigHashAll)
				Expect(err).To(BeNil())
			}
			actualSize := btc.TxVirtualSize(redeemTx)
			Expect(estimatedSize).Should(BeNumerically(">=", actualSize))
			By(fmt.Sprintf("Estimated tx size = %v, actual tx size = %v", estimatedSize, actualSize))
			Expect(estimatedSize - actualSize).Should(BeNumerically("<=", 10))
			err = client.SubmitTx(ctx, redeemTx)
			Expect(err).To(BeNil())
			By(fmt.Sprintf("Redeem tx hash = %v", color.YellowString(redeemTx.TxHash().String())))
		})
	})

	It("should create a HTLC script and redeem it using secret", func(ctx context.Context) {
		By("Initialization")
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
		client, err := btctest.RegtestClient()
		Expect(err).To(BeNil())
		indexer := btctest.RegtestIndexer()

		By("funding the addresses")
		txhash1, err := testutil.NigiriFaucet(pkAddr1.EncodeAddress())
		Expect(err).To(BeNil())
		By(fmt.Sprintf("Funding address1 %v , txid = %v", pkAddr1.EncodeAddress(), txhash1))
		txhash2, err := testutil.NigiriFaucet(pkAddr2.EncodeAddress())
		Expect(err).To(BeNil())
		By(fmt.Sprintf("Funding address2 %v , txid = %v", pkAddr2.EncodeAddress(), txhash2))
		time.Sleep(5 * time.Second)

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
		amount, feeRate := int64(1e5), 10
		fundingRecipients := []btc.Recipient{
			{
				To:     walletAddr.EncodeAddress(),
				Amount: amount,
			},
		}
		fundingTx, err := btc.BuildTransaction(feeRate, network, nil, utxos, fundingRecipients, 0, 0, pkAddr1)
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
		Expect(client.SubmitTx(ctx, fundingTx)).Should(Succeed())
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
		err = btc.SignMultisig(refundTx, 0, fundingTx.TxOut[0].Value, privKey1, privKey2, btc.SigHashSingleAnyoneCanPay)
		Expect(err).To(BeNil())

		By("Fill the tx to cover the fees")
		fundingTxHash := fundingTx.TxHash()
		fee := int64(1000)
		refundTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&fundingTxHash, 1), nil, nil))
		refundTx.AddTxOut(wire.NewTxOut(fundingTx.TxOut[1].Value-fee, fundingTx.TxOut[1].PkScript))

		By("Sign and submit refund transaction")
		pkScript, err := txscript.PayToAddrScript(pkAddr1)
		Expect(err).To(BeNil())
		sigScript, err := txscript.SignatureScript(refundTx, 1, pkScript, txscript.SigHashAll, privKey1, true)
		Expect(err).To(BeNil())
		refundTx.TxIn[1].SignatureScript = sigScript
		Expect(client.SubmitTx(ctx, refundTx)).Should(Succeed())
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
				Amount: amount - 2e3,
			},
		}
		refundSpendTx, err := btc.BuildTransaction(feeRate, network, nil, refundSpendInput, refundSpendRecipients, 0, 0, pkAddr1)
		Expect(err).To(BeNil())

		By("Sign and submit the refund spend tx")
		err = btc.SignHtlcScript(refundSpendTx, pubKey1.SerializeCompressed(), pubKey2.SerializeCompressed(), secret, secretHash[:], 0, amount, waitTime, privKey2)
		Expect(err).To(BeNil())
		err = client.SubmitTx(ctx, refundSpendTx)
		Expect(err).To(BeNil())
		By(fmt.Sprintf("RefundSpendTx tx hash = %v", color.YellowString(refundSpendTx.TxHash().String())))
	})

	It("should create instant wallet and refund from it after timelock", func(ctx context.Context) {
		By("Initialization")
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
		client, err := btctest.RegtestClient()
		Expect(err).To(BeNil())
		indexer := btctest.RegtestIndexer()

		By("funding the addresses")
		txhash1, err := testutil.NigiriFaucet(pkAddr1.EncodeAddress())
		Expect(err).To(BeNil())
		By(fmt.Sprintf("Funding address1 %v , txid = %v", pkAddr1.EncodeAddress(), txhash1))
		txhash2, err := testutil.NigiriFaucet(pkAddr2.EncodeAddress())
		Expect(err).To(BeNil())
		By(fmt.Sprintf("Funding address2 %v , txid = %v", pkAddr2.EncodeAddress(), txhash2))
		time.Sleep(5 * time.Second)

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
		amount, feeRate := int64(1e5), 10
		fundingRecipients := []btc.Recipient{
			{
				To:     walletAddr.EncodeAddress(),
				Amount: amount,
			},
		}
		fundingTx, err := btc.BuildTransaction(feeRate, network, nil, utxos, fundingRecipients, 0, 0, pkAddr1)
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
		Expect(client.SubmitTx(ctx, fundingTx)).Should(Succeed())
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
		err = btc.SignMultisig(refundTx, 0, fundingTx.TxOut[0].Value, privKey1, privKey2, btc.SigHashSingleAnyoneCanPay)
		Expect(err).To(BeNil())

		By("Fill the tx to cover the fees")
		fundingTxHash := fundingTx.TxHash()
		fee := int64(1e3)
		refundTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&fundingTxHash, 1), nil, nil))
		refundTx.AddTxOut(wire.NewTxOut(fundingTx.TxOut[1].Value-fee, fundingTx.TxOut[1].PkScript))

		By("Sign and submit refund transaction")
		pkScript, err := txscript.PayToAddrScript(pkAddr1)
		Expect(err).To(BeNil())
		sigScript, err := txscript.SignatureScript(refundTx, 1, pkScript, txscript.SigHashAll, privKey1, true)
		Expect(err).To(BeNil())
		refundTx.TxIn[1].SignatureScript = sigScript
		Expect(client.SubmitTx(ctx, refundTx)).Should(Succeed())
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
				Amount: amount - 2*fee,
			},
		}
		refundSpendTx, err := btc.BuildTransaction(feeRate, network, nil, refundSpendInput, refundSpendRecipients, 0, 0, pkAddr1)
		Expect(err).To(BeNil())

		By("Sign the refund spend tx")
		refundSpendTx.TxIn[0].Sequence = uint32(waitTime)
		err = btc.SignHtlcScript(refundSpendTx, pubKey1.SerializeCompressed(), pubKey2.SerializeCompressed(), secret, secretHash[:], 0, amount, waitTime, privKey1)
		Expect(err).To(BeNil())

		By("You shouldn't be able to spend the refund tx since it doesnt have enough confirmations")
		refundTxHash := refundTx.TxHash()
		rawRefund, err := client.GetRawTransaction(ctx, &refundTxHash)
		Expect(err).To(BeNil())
		Expect(rawRefund.Confirmations).Should(Equal(uint64(waitTime) - 1))
		err = client.SubmitTx(ctx, refundSpendTx)
		Expect(err.Error()).Should(ContainSubstring("non-BIP68-final"))

		By("Mine a new block and we should be able to submit the refund spend tx")
		Expect(testutil.NigiriNewBlock()).Should(Succeed())
		time.Sleep(time.Second)
		err = client.SubmitTx(ctx, refundSpendTx)
		Expect(err).To(BeNil())
		By(fmt.Sprintf("RefundSpendTx tx hash = %v", color.YellowString(refundSpendTx.TxHash().String())))
	})
})
