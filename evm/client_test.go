package evm_test

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/localnet"
)

var _ = Describe("Client", func() {
	It("should be able to build a client", func() {
		client, err := localnet.EVMClient()
		Expect(err).Should(BeNil())
		evmClient, ok := client.EvmClient(blockchain.NewEvmChain(blockchain.EthereumLocalnet))
		Expect(ok).Should(BeTrue())
		chainID, err := evmClient.ChainID(context.Background())
		Expect(err).Should(BeNil())
		Expect(chainID.Cmp(big.NewInt(31337))).Should(Equal(0))
	})

	It("should be able to get balance", func() {
		chain := blockchain.NewEvmChain(blockchain.EthereumLocalnet)
		client, err := localnet.EVMClient()
		Expect(err).Should(BeNil())

		balance, err := client.Balance(context.Background(), blockchain.NewETH(chain, common.HexToAddress("")), common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"), nil)
		Expect(err).Should(BeNil())
		Expect(balance.Cmp(big.NewInt(0))).To(Equal(1))
	})
})
