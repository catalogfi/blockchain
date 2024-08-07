package btc_test

import (
	"context"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/localnet"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("--- Event ---", Ordered, func() {
	chainParams := &chaincfg.RegressionNetParams
	indexer := localnet.BTCIndexer()
	fixedFeeEstimator := btc.NewFixFeeEstimator(16)

	htlcClient := btc.NewHTLCClient(indexer)

	alicePrivKey, err := btcec.NewPrivateKey()
	Expect(err).To(BeNil())
	aliceSimpleWallet, err := btc.NewSimpleWallet(alicePrivKey, chainParams, indexer, fixedFeeEstimator, btc.HighFee)
	Expect(err).To(BeNil())

	bobPrivKey, err := btcec.NewPrivateKey()
	Expect(err).To(BeNil())
	bobSimpleWallet, err := btc.NewSimpleWallet(bobPrivKey, chainParams, indexer, fixedFeeEstimator, btc.HighFee)
	Expect(err).To(BeNil())

	BeforeAll(func() {
		By("Fund Alice's and Bob's wallet")
		_, err := localnet.FundBitcoin(aliceSimpleWallet.Address().EncodeAddress(), indexer)
		Expect(err).To(BeNil())

		_, err = aliceSimpleWallet.Send(context.Background(), []btc.SendRequest{
			{
				To:     bobSimpleWallet.Address(),
				Amount: 50000000,
			},
		}, nil, nil)
		Expect(err).To(BeNil())
	})

	Describe("- Initiate -", func() {
		It("should able to get the initiate event", func(ctx context.Context) {
			aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, chainParams)
			Expect(err).To(BeNil())
			aliceHTLC, _, err := generateHTLC(alicePrivKey, bobPrivKey)
			Expect(err).To(BeNil())
			aliceHTLCAddress, err := btc.HTLCWallet.Address(aliceHTLCWallet, aliceHTLC)
			Expect(err).To(BeNil())
			aliceHTLCAsset := btc.NewBTCAsset(aliceHTLCAddress, blockchain.UtxoChain{})

			txID, err := aliceHTLCWallet.Initiate(ctx, aliceHTLC, 1000000)
			Expect(err).To(BeNil())

			aliceHTLCEvent, err := htlcClient.HTLCEvents(ctx, aliceHTLCAsset, 0, 0)
			Expect(err).To(BeNil())

			Expect(aliceHTLCEvent).To(HaveLen(1))
			Expect(aliceHTLCEvent[0].TxHash()).To(Equal(txID))
			Expect(aliceHTLCEvent[0].BlockNumber()).To(Equal(uint64(0)))

			err = localnet.MineBitcoinBlocks(1, indexer)
			Expect(err).To(BeNil())

			aliceHTLCEvent, err = htlcClient.HTLCEvents(ctx, aliceHTLCAsset, 0, 0)
			Expect(err).To(BeNil())

			Expect(aliceHTLCEvent).To(HaveLen(1))
			Expect(aliceHTLCEvent[0].TxHash()).To(Equal(txID))
			Expect(aliceHTLCEvent[0].BlockNumber()).ToNot(Equal(uint64(0)))

			switch e := aliceHTLCEvent[0].(type) {
			case btc.HTLCInitiated:
				Expect(e.Amount()).To(Equal(uint64(1000000)))
			}
		})
	})

	Describe("- Redeem -", func() {
		It("should able to get the redeem event", func(ctx context.Context) {
			aliceHTLC, secret, err := generateHTLC(alicePrivKey, bobPrivKey)
			Expect(err).To(BeNil())

			aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, chainParams)
			Expect(err).To(BeNil())
			txid, err := aliceHTLCWallet.Initiate(ctx, aliceHTLC, 1000000)
			Expect(err).To(BeNil())
			Expect(txid).ToNot(BeEmpty())

			err = localnet.MineBitcoinBlocks(1, indexer)
			Expect(err).To(BeNil())

			aliceHTLCAddress, err := btc.HTLCWallet.Address(aliceHTLCWallet, aliceHTLC)
			Expect(err).To(BeNil())
			aliceHTLCAsset := btc.NewBTCAsset(aliceHTLCAddress, blockchain.UtxoChain{})

			bobHTLC := turnRolesInHTLC(aliceHTLC)
			bobHTLCWallet, err := btc.NewHTLCWallet(bobSimpleWallet, indexer, chainParams)
			Expect(err).To(BeNil())
			txid, err = bobHTLCWallet.Initiate(ctx, bobHTLC, 1000000)
			Expect(err).To(BeNil())
			Expect(txid).ToNot(BeEmpty())

			err = localnet.MineBitcoinBlocks(1, indexer)
			Expect(err).To(BeNil())

			bobHTLCAddress, err := btc.HTLCWallet.Address(bobHTLCWallet, bobHTLC)
			Expect(err).To(BeNil())
			bobHTLCAsset := btc.NewBTCAsset(bobHTLCAddress, blockchain.UtxoChain{})

			txid, err = aliceHTLCWallet.Redeem(ctx, bobHTLC, secret)
			Expect(err).To(BeNil())
			Expect(txid).ToNot(BeEmpty())

			bobHTLCEvent, err := htlcClient.HTLCEvents(ctx, bobHTLCAsset, 0, 0)
			Expect(err).To(BeNil())
			Expect(bobHTLCEvent).To(HaveLen(2))
			Expect(bobHTLCEvent[0].TxHash()).To(Equal(txid))
			Expect(bobHTLCEvent[0].BlockNumber()).To(Equal(uint64(0)))

			err = localnet.MineBitcoinBlocks(1, indexer)
			Expect(err).To(BeNil())

			bobHTLCEvent, err = htlcClient.HTLCEvents(ctx, bobHTLCAsset, 0, 0)
			Expect(err).To(BeNil())
			Expect(bobHTLCEvent).To(HaveLen(2))
			Expect(bobHTLCEvent[0].TxHash()).To(Equal(txid))
			Expect(bobHTLCEvent[0].BlockNumber()).ToNot(Equal(uint64(0)))

			switch e := bobHTLCEvent[0].(type) {
			case btc.HTLCRedeemed:
				Expect(hex.EncodeToString(e.Secret())).To(Equal(hex.EncodeToString(secret)))
				Expect(e.RedeemerPubkey()).To(Equal(hex.EncodeToString(alicePrivKey.PubKey().SerializeCompressed())[2:]))
			}

			txid, err = bobHTLCWallet.Redeem(ctx, aliceHTLC, secret)
			Expect(err).To(BeNil())
			Expect(txid).ToNot(BeEmpty())

			aliceHTLCEvent, err := htlcClient.HTLCEvents(ctx, aliceHTLCAsset, 0, 0)
			Expect(err).To(BeNil())
			Expect(aliceHTLCEvent).To(HaveLen(2))
			Expect(aliceHTLCEvent[0].TxHash()).To(Equal(txid))
			Expect(aliceHTLCEvent[0].BlockNumber()).To(Equal(uint64(0)))

			err = localnet.MineBitcoinBlocks(1, indexer)
			Expect(err).To(BeNil())

			aliceHTLCEvent, err = htlcClient.HTLCEvents(ctx, aliceHTLCAsset, 0, 0)
			Expect(err).To(BeNil())
			Expect(aliceHTLCEvent).To(HaveLen(2))
			Expect(aliceHTLCEvent[0].TxHash()).To(Equal(txid))
			Expect(aliceHTLCEvent[0].BlockNumber()).ToNot(Equal(uint64(0)))

			switch e := aliceHTLCEvent[0].(type) {
			case btc.HTLCRedeemed:
				Expect(hex.EncodeToString(e.Secret())).To(Equal(hex.EncodeToString(secret)))
				Expect(e.RedeemerPubkey()).To(Equal(hex.EncodeToString(bobPrivKey.PubKey().SerializeCompressed())[2:]))
			}
		})
	})

	Describe("- Refund -", func() {
		It("should able to get the refund event", func(ctx context.Context) {
			aliceHTLC, _, err := generateHTLC(alicePrivKey, bobPrivKey)
			Expect(err).To(BeNil())

			aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, chainParams)
			Expect(err).To(BeNil())
			txid, err := aliceHTLCWallet.Initiate(ctx, aliceHTLC, 1000000)
			Expect(err).To(BeNil())
			Expect(txid).ToNot(BeEmpty())

			err = localnet.MineBitcoinBlocks(1, indexer)
			Expect(err).To(BeNil())

			aliceHTLCAddress, err := btc.HTLCWallet.Address(aliceHTLCWallet, aliceHTLC)
			Expect(err).To(BeNil())
			aliceHTLCAsset := btc.NewBTCAsset(aliceHTLCAddress, blockchain.UtxoChain{})

			err = localnet.MineBitcoinBlocks(int(aliceHTLC.Timelock), indexer)
			Expect(err).To(BeNil())

			txid, err = aliceHTLCWallet.Refund(ctx, aliceHTLC, nil)
			Expect(err).To(BeNil())
			Expect(txid).ToNot(BeEmpty())

			aliceHTLCEvent, err := htlcClient.HTLCEvents(ctx, aliceHTLCAsset, 0, 0)
			Expect(err).To(BeNil())
			Expect(aliceHTLCEvent).To(HaveLen(2))
			Expect(aliceHTLCEvent[0].TxHash()).To(Equal(txid))
			Expect(aliceHTLCEvent[0].BlockNumber()).To(Equal(uint64(0)))

			err = localnet.MineBitcoinBlocks(1, indexer)
			Expect(err).To(BeNil())

			aliceHTLCEvent, err = htlcClient.HTLCEvents(ctx, aliceHTLCAsset, 0, 0)
			Expect(err).To(BeNil())
			Expect(aliceHTLCEvent).To(HaveLen(2))
			Expect(aliceHTLCEvent[0].TxHash()).To(Equal(txid))
			Expect(aliceHTLCEvent[0].BlockNumber()).ToNot(Equal(uint64(0)))

			switch e := aliceHTLCEvent[0].(type) {
			case btc.HTLCRefunded:
				Expect(e.RefunderPubkey()).To(Equal(hex.EncodeToString(alicePrivKey.PubKey().SerializeCompressed())[2:]))
			}
		})

		It("should able to get the instant refund event", func(ctx context.Context) {
			aliceHTLC, _, err := generateHTLC(alicePrivKey, bobPrivKey)
			Expect(err).To(BeNil())

			aliceHTLCWallet, err := btc.NewHTLCWallet(aliceSimpleWallet, indexer, chainParams)
			Expect(err).To(BeNil())
			txid, err := aliceHTLCWallet.Initiate(ctx, aliceHTLC, 1000000)
			Expect(err).To(BeNil())
			Expect(txid).ToNot(BeEmpty())

			err = localnet.MineBitcoinBlocks(1, indexer)
			Expect(err).To(BeNil())

			aliceHTLCAddress, err := btc.HTLCWallet.Address(aliceHTLCWallet, aliceHTLC)
			Expect(err).To(BeNil())
			aliceHTLCAsset := btc.NewBTCAsset(aliceHTLCAddress, blockchain.UtxoChain{})

			bobHTLCWallet, err := btc.NewHTLCWallet(bobSimpleWallet, indexer, chainParams)
			Expect(err).To(BeNil())
			instantRefundTxBytes, err := bobHTLCWallet.GenerateInstantRefundSACP(ctx, aliceHTLC, aliceSimpleWallet.Address())
			Expect(err).To(BeNil())

			txid, err = aliceHTLCWallet.Refund(ctx, aliceHTLC, instantRefundTxBytes)
			Expect(err).To(BeNil())
			Expect(txid).ToNot(BeEmpty())

			aliceHTLCEvent, err := htlcClient.HTLCEvents(ctx, aliceHTLCAsset, 0, 0)
			Expect(err).To(BeNil())
			Expect(aliceHTLCEvent).To(HaveLen(2))
			Expect(aliceHTLCEvent[0].TxHash()).To(Equal(txid))
			Expect(aliceHTLCEvent[0].BlockNumber()).To(Equal(uint64(0)))

			err = localnet.MineBitcoinBlocks(1, indexer)
			Expect(err).To(BeNil())

			aliceHTLCEvent, err = htlcClient.HTLCEvents(ctx, aliceHTLCAsset, 0, 0)
			Expect(err).To(BeNil())
			Expect(aliceHTLCEvent).To(HaveLen(2))
			Expect(aliceHTLCEvent[0].TxHash()).To(Equal(txid))
			Expect(aliceHTLCEvent[0].BlockNumber()).ToNot(Equal(uint64(0)))

			switch e := aliceHTLCEvent[0].(type) {
			case btc.HTLCRefunded:
				Expect(e.RefunderPubkey()).To(Equal(hex.EncodeToString(alicePrivKey.PubKey().SerializeCompressed())[2:]))
			}
		})
	})
})
