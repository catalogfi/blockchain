package evm_test

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm"
	"github.com/catalogfi/blockchain/localnet"
)

var _ = Describe("HTLC", func() {
	It("should be able to do an atomic swap", func() {
		initiatorWallet, err := localnet.EVMHTLCWallet(0)
		Expect(err).Should(BeNil())
		redeemerWallet, err := localnet.EVMHTLCWallet(1)
		Expect(err).Should(BeNil())

		htlcClient, err := localnet.EVMHTLCClient()
		Expect(err).Should(BeNil())

		cl, ok := htlcClient.EvmClient(blockchain.NewEvmChain(blockchain.ArbitrumLocalnet))
		Expect(ok).Should(BeTrue())

		start, err := cl.BlockNumber(context.Background())
		Expect(err).Should(BeNil())

		secret := [32]byte{}
		rand.Read(secret[:])
		secretHash := sha256.Sum256(secret[:])

		itx, err := initiatorWallet.Initiate(context.Background(), localnet.ArbitrumWBTC(), redeemerWallet.Address(), secretHash, big.NewInt(2), big.NewInt(1000000), nil)
		Expect(err).Should(BeNil())

		end, err := cl.BlockNumber(context.Background())
		Expect(err).Should(BeNil())

		events, err := htlcClient.HTLCEvents(context.Background(), localnet.ArbitrumWBTC(), new(big.Int).SetUint64(start+1), new(big.Int).SetUint64(end))
		Expect(err).Should(BeNil())

		oid := initiatorWallet.OrderID(secretHash)
		Expect(len(events)).Should(Equal(1))

		iEvent, ok := events[0].(evm.HTLCInitiated)
		Expect(ok).Should(BeTrue())

		Expect(iEvent.InitiateTxHash.Hex()).Should(Equal(itx.TxHash.Hex()))
		Expect(iEvent.ID).Should(Equal(oid))
		Expect(iEvent.Amount).Should(Equal(big.NewInt(1000000)))
		Expect(iEvent.InitiateTxBlockNumber).Should(Equal(itx.BlockNumber.Uint64()))
		Expect(iEvent.SecretHash).Should(Equal(secretHash))

		rtx, err := redeemerWallet.Redeem(context.Background(), localnet.ArbitrumWBTC(), oid, secret[:])
		Expect(err).Should(BeNil())

		end2, err := cl.BlockNumber(context.Background())
		Expect(err).Should(BeNil())

		events2, err := htlcClient.HTLCEvents(context.Background(), localnet.ArbitrumWBTC(), new(big.Int).SetUint64(end+1), new(big.Int).SetUint64(end2))
		Expect(err).Should(BeNil())

		rEvent, ok := events2[0].(evm.HTLCRedeemed)
		Expect(ok).Should(BeTrue())
		Expect(rEvent.RedeemTxHash).Should(Equal(rtx.TxHash))
		Expect(rEvent.ID).Should(Equal(oid))
		Expect(rEvent.Secret).Should(Equal(secret[:]))
		Expect(rEvent.RedeemTxBlockNumber).Should(Equal(rtx.BlockNumber.Uint64()))
		Expect(rEvent.SecretHash).Should(Equal(secretHash))
	})
})
