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
		client := localnet.EVMClient()
		evmClient, ok := client.EvmClient(blockchain.NewEvmChain(blockchain.EthereumLocalnet))
		Expect(ok).Should(BeTrue())
		chainID, err := evmClient.ChainID(context.Background())
		Expect(err).Should(BeNil())
		Expect(chainID.Cmp(big.NewInt(31337))).Should(Equal(0))
	})

	It("should be able to get balance", func() {
		chain := blockchain.NewEvmChain(blockchain.EthereumLocalnet)
		client := localnet.EVMClient()
		balance, err := client.Balance(context.Background(), blockchain.NewETH(chain, common.HexToAddress("")), common.HexToAddress("0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199"), nil)
		Expect(err).Should(BeNil())
		expectedBalance, ok := new(big.Int).SetString("10000000000000000000000", 10)
		Expect(ok).Should(BeTrue())
		Expect(balance.Cmp(expectedBalance)).To(Equal(0))
	})
})
