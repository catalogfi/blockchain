package btc_test

import (
	"context"
	"crypto/rand"
	"crypto/sha256"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/localnet"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTLC Wallet(p2tr)", Ordered, func() {

	chainParams := chaincfg.RegressionNetParams
	indexer := localnet.BTCIndexer()
	alicePrivKey, err := btcec.NewPrivateKey()
	Expect(err).To(BeNil())
	bobPrivKey, err := btcec.NewPrivateKey()
	Expect(err).To(BeNil())

	fixedFeeEstimator := btc.NewFixFeeEstimator(10)

	aliceSimpleWallet, err := btc.NewSimpleWallet(alicePrivKey, &chainParams, indexer, fixedFeeEstimator, btc.HighFee)
	Expect(err).To(BeNil())

	bobSimpleWallet, err := btc.NewSimpleWallet(bobPrivKey, &chainParams, indexer, fixedFeeEstimator, btc.HighFee)
	Expect(err).To(BeNil())

	initiateAmount := int64(1e6)

	BeforeAll(func() {
		By("Fund Alice and Bob wallets")
		_, err := localnet.FundBitcoin(aliceSimpleWallet.Address().EncodeAddress(), indexer)
		Expect(err).To(BeNil())
		// Send half of the amount to Bob
		_, err = aliceSimpleWallet.Send(context.Background(), []btc.SendRequest{
			{
				To:     bobSimpleWallet.Address(),
				Amount: 50000000,
			},
		}, nil, nil)
		Expect(err).To(BeNil())
	})

	It("should be able to generate HTLC address", func(ctx context.Context) {

		By("Lets create AliceHTLCWallet where Alice is the initiator and Bob is the redeemer")
		aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())

		aliceHTLC, _, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())
		htlcAddr, err := aliceHTLCWallet.Address(aliceHTLC)
		Expect(err).To(BeNil())
		Expect(htlcAddr).NotTo(BeNil())
	})

	It("should be able to initiate HTLC", func(ctx context.Context) {
		aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())
		aliceHTLC, _, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())
		txid, err := aliceHTLCWallet.Initiate(ctx, aliceHTLC, initiateAmount)
		Expect(err).To(BeNil())
		Expect(txid).NotTo(BeEmpty())

		tx, err := indexer.GetTx(ctx, txid)
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())
		Expect(tx.VOUTs[0].Value).To(Equal(int(initiateAmount)))
	})

	It("should be able to initiate and redeem HTLC", func(ctx context.Context) {
		aliceHTLC, secret, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())

		By("Initiate Alice HTLC")
		aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())
		txid, err := aliceHTLCWallet.Initiate(ctx, aliceHTLC, initiateAmount)
		Expect(err).To(BeNil())
		Expect(txid).NotTo(BeEmpty())

		By("Initiate Bob HTLC")
		bobHTLC := turnRolesInHTLC(aliceHTLC)
		bobHTLCWallet, err := btc.NewHTLCWallet(bobSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())
		txid, err = bobHTLCWallet.Initiate(ctx, bobHTLC, initiateAmount)
		Expect(err).To(BeNil())
		Expect(txid).NotTo(BeEmpty())

		By("Redeem Alice HTLC")
		sec, _, err := generateSecret()
		Expect(err).To(BeNil())
		// should fail as the secret is invalid
		_, err = aliceHTLCWallet.Redeem(ctx, bobHTLC, sec)
		Expect(err).Should(Equal(btc.ErrInvalidSecret))

		txid, err = aliceHTLCWallet.Redeem(ctx, bobHTLC, secret)
		Expect(err).To(BeNil())
		Expect(txid).NotTo(BeEmpty())

		By("Redeem Bob HTLC")
		Expect(err).To(BeNil())
		txid, err = bobHTLCWallet.Redeem(ctx, aliceHTLC, secret)
		Expect(err).To(BeNil())
		Expect(txid).NotTo(BeEmpty())

	})

	It("should be able to initiate and refund HTLC", func(ctx context.Context) {
		aliceHTLC, _, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())

		By("Initiate Alice HTLC")
		aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())
		txid, err := aliceHTLCWallet.Initiate(ctx, aliceHTLC, initiateAmount)
		Expect(err).To(BeNil())
		Expect(txid).NotTo(BeEmpty())

		By("Mine expiry no of blocks")
		err = localnet.MineBitcoinBlocks(int(aliceHTLC.LockTime), indexer)
		Expect(err).To(BeNil())

		By("Refund Alice HTLC")
		txid, err = aliceHTLCWallet.Refund(ctx, aliceHTLC, nil)
		Expect(err).To(BeNil())
		Expect(txid).NotTo(BeEmpty())

	})

	It("should be able to initiate and refund HTLC instantly", func(ctx context.Context) {

		aliceHTLC, _, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())

		By("Initiate Alice HTLC")
		aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())
		txid, err := aliceHTLCWallet.Initiate(ctx, aliceHTLC, initiateAmount)
		Expect(err).To(BeNil())
		Expect(txid).NotTo(BeEmpty())

		By("Instant refund Alice HTLC")
		bobHTLCWallet, err := btc.NewHTLCWallet(bobSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())

		instantRefundTxBytes, err := bobHTLCWallet.GenerateInstantRefundSACP(ctx, aliceHTLC, aliceSimpleWallet.Address())
		Expect(err).To(BeNil())

		txid, err = aliceHTLCWallet.Refund(ctx, aliceHTLC, instantRefundTxBytes)
		Expect(err).To(BeNil())
		Expect(txid).NotTo(BeEmpty())

		tx, _, err := aliceHTLCWallet.Status(ctx, txid)
		Expect(err).To(BeNil())

		Expect(tx.VOUTs[0].Value).To(Equal(int(initiateAmount) - int(tx.Fee)))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).To(Equal(aliceSimpleWallet.Address().EncodeAddress()))
	})

	It("should not be able to refund instantly if the SACP is not correct", func(ctx context.Context) {

		aliceHTLC, _, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())

		By("Initiate Alice HTLC")
		aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())
		txid, err := aliceHTLCWallet.Initiate(ctx, aliceHTLC, initiateAmount)
		Expect(err).To(BeNil())
		Expect(txid).NotTo(BeEmpty())

		By("Instant refund should fail")
		bobHTLCWallet, err := btc.NewHTLCWallet(bobSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())

		instantRefundTxBytes, err := bobHTLCWallet.GenerateInstantRefundSACP(ctx, aliceHTLC, bobSimpleWallet.Address())
		Expect(err).To(BeNil())

		_, err = aliceHTLCWallet.Refund(ctx, aliceHTLC, instantRefundTxBytes)
		Expect(err).Should(Equal(btc.ErrSACPInvalidOutput))

		// lets initiate bobHTLC
		bobHTLC := turnRolesInHTLC(aliceHTLC)
		_, err = bobHTLCWallet.Initiate(ctx, bobHTLC, initiateAmount)
		Expect(err).To(BeNil())

		instantRefundTxBytes, err = bobHTLCWallet.GenerateInstantRefundSACP(ctx, bobHTLC, aliceSimpleWallet.Address())
		Expect(err).To(BeNil())

		By("Instant refund Alice HTLC should fail as the inputs are not correct")
		_, err = aliceHTLCWallet.Refund(ctx, aliceHTLC, instantRefundTxBytes)
		Expect(err).Should(Equal(btc.ErrSACPInvalidInput))
	})

	It("should be able to refund instantly even if the fee is low", func(ctx context.Context) {
		aliceHTLC, _, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())

		By("Initiate Alice HTLC")
		aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())
		txid, err := aliceHTLCWallet.Initiate(ctx, aliceHTLC, initiateAmount)
		Expect(err).To(BeNil())
		Expect(txid).NotTo(BeEmpty())

		// Let's create a bobSimpleWallet with lower fees
		feeEstimater := btc.NewFixFeeEstimator(1)
		bobSimpleWallet, err := btc.NewSimpleWallet(bobPrivKey, &chainParams, indexer, feeEstimater, btc.LowFee)
		Expect(err).To(BeNil())
		bobHTLCWallet, err := btc.NewHTLCWallet(bobSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())

		//generate SACP
		instantRefundTxBytes, err := bobHTLCWallet.GenerateInstantRefundSACP(ctx, aliceHTLC, aliceSimpleWallet.Address())
		Expect(err).To(BeNil())

		txId, err := aliceHTLCWallet.Refund(ctx, aliceHTLC, instantRefundTxBytes)
		Expect(err).To(BeNil())
		Expect(txId).NotTo(BeEmpty())

		tx, _, err := aliceHTLCWallet.Status(ctx, txId)
		Expect(err).To(BeNil())

		// One is the refund tx and the other is the change tx to adjust the fee
		Expect(tx.VOUTs).To(HaveLen(2))
	})

	It("should be able to redeem or refund multiple initiates", func(ctx context.Context) {
		aliceHTLC, secret, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())

		By("Initiate Alice HTLC")
		aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())
		_, err = aliceHTLCWallet.Initiate(ctx, aliceHTLC, initiateAmount)
		Expect(err).To(BeNil())
		// Initiate again
		_, err = aliceHTLCWallet.Initiate(ctx, aliceHTLC, initiateAmount)
		Expect(err).To(BeNil())

		By("Redeem Bob HTLC")
		bobHTLCWallet, err := btc.NewHTLCWallet(bobSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())

		txid, err := bobHTLCWallet.Redeem(ctx, aliceHTLC, secret)
		Expect(err).To(BeNil())
		Expect(txid).NotTo(BeEmpty())

		By("Redeemed funds should be equal to the sum of the two initiates - fee")
		tx, _, err := bobHTLCWallet.Status(ctx, txid)
		Expect(err).To(BeNil())
		Expect(tx.VOUTs[0].Value).To(Equal(2*int(initiateAmount) - int(tx.Fee)))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).To(Equal(bobSimpleWallet.Address().EncodeAddress()))

		// Initiate again
		_, err = aliceHTLCWallet.Initiate(ctx, aliceHTLC, initiateAmount)
		Expect(err).To(BeNil())
		_, err = aliceHTLCWallet.Initiate(ctx, aliceHTLC, initiateAmount)
		Expect(err).To(BeNil())

		err = localnet.MineBitcoinBlocks(int(aliceHTLC.LockTime), indexer)
		Expect(err).To(BeNil())

		txid, err = aliceHTLCWallet.Refund(ctx, aliceHTLC, nil)
		Expect(err).To(BeNil())

		tx, _, err = aliceHTLCWallet.Status(ctx, txid)
		Expect(err).To(BeNil())
		Expect(tx.VOUTs[0].Value).To(Equal(2*int(initiateAmount) - int(tx.Fee)))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).To(Equal(aliceSimpleWallet.Address().EncodeAddress()))
	})

	It("should be able to initiate and redeem multiple HTLCs", func(ctx context.Context) {

		aliceHTLC1, secret1, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())

		aliceHTLC2, secret2, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())

		aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())

		By("Initiate HTLCs")

		txID, err := aliceHTLCWallet.Execute(ctx, []btc.RawHTLCAction{
			{
				Action: btc.InitiateHTLCAction,
				HTLC:   *aliceHTLC1,
				Amount: initiateAmount,
			},
			{
				Action: btc.InitiateHTLCAction,
				HTLC:   *aliceHTLC2,
				Amount: initiateAmount,
			},
		})
		Expect(err).To(BeNil())
		Expect(txID).NotTo(BeEmpty())

		bobHTLC1 := turnRolesInHTLC(aliceHTLC1)
		bobHTLC2 := turnRolesInHTLC(aliceHTLC2)

		bobHTLCWallet, err := btc.NewHTLCWallet(bobSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())

		txID, err = bobHTLCWallet.Execute(ctx, []btc.RawHTLCAction{
			{
				Action: btc.InitiateHTLCAction,
				HTLC:   *bobHTLC1,
				Amount: initiateAmount,
			},
			{
				Action: btc.InitiateHTLCAction,
				HTLC:   *bobHTLC2,
				Amount: initiateAmount,
			},
		})
		Expect(err).To(BeNil())
		Expect(txID).NotTo(BeEmpty())

		By("Redeem HTLCs")

		aliceRedeemTxID, err := aliceHTLCWallet.Execute(ctx, []btc.RawHTLCAction{
			{
				Action: btc.RedeemHTLCAction,
				HTLC:   *bobHTLC1,
				Secret: secret1,
			},
			{
				Action: btc.RedeemHTLCAction,
				HTLC:   *bobHTLC2,
				Secret: secret2,
			},
		})
		Expect(err).To(BeNil())
		Expect(aliceRedeemTxID).NotTo(BeEmpty())

		bobRedeemTxID, err := bobHTLCWallet.Execute(ctx, []btc.RawHTLCAction{
			{
				Action: btc.RedeemHTLCAction,
				HTLC:   *aliceHTLC1,
				Secret: secret1,
			},
			{
				Action: btc.RedeemHTLCAction,
				HTLC:   *aliceHTLC2,
				Secret: secret2,
			},
		})
		Expect(err).To(BeNil())
		Expect(bobRedeemTxID).NotTo(BeEmpty())

		By("Verify the status")

		tx, _, err := aliceHTLCWallet.Status(ctx, aliceRedeemTxID)
		Expect(err).To(BeNil())
		Expect(tx.VOUTs).Should(HaveLen(1))
		Expect(tx.VOUTs[0].Value).To(Equal(2*int(initiateAmount) - int(tx.Fee)))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).To(Equal(aliceSimpleWallet.Address().EncodeAddress()))

		tx, _, err = bobHTLCWallet.Status(ctx, bobRedeemTxID)
		Expect(err).To(BeNil())
		Expect(tx.VOUTs).Should(HaveLen(1))
		Expect(tx.VOUTs[0].Value).To(Equal(2*int(initiateAmount) - int(tx.Fee)))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).To(Equal(bobSimpleWallet.Address().EncodeAddress()))
	})

	It("should be able to initiate and refund multiple HTLCs", func(ctx context.Context) {
		aliceHTLC1, _, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())

		aliceHTLC2, _, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())

		aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())

		By("Initiate HTLCs")

		txID, err := aliceHTLCWallet.Execute(ctx, []btc.RawHTLCAction{
			{
				Action: btc.InitiateHTLCAction,
				HTLC:   *aliceHTLC1,
				Amount: initiateAmount,
			},
			{
				Action: btc.InitiateHTLCAction,
				HTLC:   *aliceHTLC2,
				Amount: initiateAmount,
			},
		})
		Expect(err).To(BeNil())
		Expect(txID).NotTo(BeEmpty())

		By("Refund HTLCs")

		// Mine expiry no of blocks
		err = localnet.MineBitcoinBlocks(int(aliceHTLC1.LockTime), indexer)
		Expect(err).To(BeNil())

		aliceHTLC3, _, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())

		// Refund and also initiate another HTLC
		aliceRefundTxID, err := aliceHTLCWallet.Execute(ctx, []btc.RawHTLCAction{
			{
				Action: btc.RefundHTLCAction,
				HTLC:   *aliceHTLC1,
			},
			{
				Action: btc.RefundHTLCAction,
				HTLC:   *aliceHTLC2,
			},
			{
				Action: btc.InitiateHTLCAction,
				HTLC:   *aliceHTLC3,
				Amount: initiateAmount,
			},
		})
		Expect(err).To(BeNil())
		Expect(aliceRefundTxID).NotTo(BeEmpty())

		By("Verify the status")
		tx, _, err := aliceHTLCWallet.Status(ctx, aliceRefundTxID)
		Expect(err).To(BeNil())

		htlc3Address, err := aliceHTLCWallet.Address(aliceHTLC3)
		Expect(err).To(BeNil())

		Expect(tx.VOUTs).Should(HaveLen(2))
		Expect(tx.VOUTs[0].Value).To(Equal(int(initiateAmount)))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).To(Equal(htlc3Address.EncodeAddress()))
		Expect(tx.VOUTs[1].Value).To(Equal(int(initiateAmount - tx.Fee)))
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).To(Equal(aliceSimpleWallet.Address().EncodeAddress()))
	})
	It("should be able to initiate and refund multiple HTLCs instantly", func(ctx context.Context) {
		aliceHTLC1, _, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())

		aliceHTLC2, _, err := generateHTLC(alicePrivKey, bobPrivKey)
		Expect(err).To(BeNil())

		aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())

		By("Initiate HTLCs")

		txID, err := aliceHTLCWallet.Execute(ctx, []btc.RawHTLCAction{
			{
				Action: btc.InitiateHTLCAction,
				HTLC:   *aliceHTLC1,
				Amount: initiateAmount,
			},
			{
				Action: btc.InitiateHTLCAction,
				HTLC:   *aliceHTLC2,
				Amount: initiateAmount,
			},
		})
		Expect(err).To(BeNil())
		Expect(txID).NotTo(BeEmpty())

		bobHTLCWallet, err := btc.NewHTLCWallet(bobSimpleWallet, indexer, &chainParams)
		Expect(err).To(BeNil())

		sacpTx1, err := bobHTLCWallet.GenerateInstantRefundSACP(ctx, aliceHTLC1, aliceSimpleWallet.Address())
		Expect(err).To(BeNil())

		sacpTx2, err := bobHTLCWallet.GenerateInstantRefundSACP(ctx, aliceHTLC2, aliceSimpleWallet.Address())
		Expect(err).To(BeNil())

		By("Instant refund HTLCs")

		txID, err = aliceHTLCWallet.Execute(ctx, []btc.RawHTLCAction{
			{
				Action:                  btc.InstantRefundHTLCAction,
				HTLC:                    *aliceHTLC1,
				InsantRefundSACPTxBytes: sacpTx1,
			},
			{
				Action:                  btc.InstantRefundHTLCAction,
				HTLC:                    *aliceHTLC2,
				InsantRefundSACPTxBytes: sacpTx2,
			},
		})
		Expect(err).To(BeNil())
		Expect(txID).NotTo(BeEmpty())

		By("Verify the status")
		tx, _, err := aliceHTLCWallet.Status(ctx, txID)
		Expect(err).To(BeNil())

		Expect(tx.VOUTs).Should(HaveLen(2))
		Expect(tx.VOUTs[0].Value).To(Equal(int(initiateAmount) - int(tx.Fee)/2))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).To(Equal(aliceSimpleWallet.Address().EncodeAddress()))
		Expect(tx.VOUTs[1].Value).To(Equal(int(initiateAmount) - int(tx.Fee)/2))
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).To(Equal(aliceSimpleWallet.Address().EncodeAddress()))
	})
})

// ------------------------------Helper functions--------------------------------

func getPubkey(privKey *btcec.PrivateKey) []byte {
	return schnorr.SerializePubKey(privKey.PubKey())
}

// turnRolesInHTLC swaps the roles of the initiator and the redeemer in the HTLC
func turnRolesInHTLC(htlc *btc.HTLC) *btc.HTLC {
	return &btc.HTLC{
		InitiatorPubkey: htlc.RedeemerPubkey,
		RedeemerPubkey:  htlc.InitiatorPubkey,
		SecretHash:      htlc.SecretHash,
		LockTime:        htlc.LockTime,
	}
}

func generateHTLC(alicePrivKey, bobPrivKey *btcec.PrivateKey) (*btc.HTLC, []byte, error) {
	// generate random 32 bytes secret
	secret, secretHash, err := generateSecret()
	if err != nil {
		return nil, secret, err
	}

	aliceHTLC := &btc.HTLC{
		InitiatorPubkey: getPubkey(alicePrivKey),
		RedeemerPubkey:  getPubkey(bobPrivKey),
		SecretHash:      secretHash[:],
		LockTime:        10,
	}
	return aliceHTLC, secret, nil
}

func generateSecret() ([]byte, []byte, error) {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		return nil, nil, err
	}
	sh := sha256.Sum256(secret)
	return secret, sh[:], nil
}
