package evm_test

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/localnet"
)

var _ = Describe("Wallet", func() {
	It("should be able to build a client", func() {
		chain := blockchain.NewEvmChain(blockchain.EthereumLocalnet)
		wallet := localnet.EVMWallet()
		tx, err := wallet.Send(context.Background(), blockchain.NewETH(chain, common.HexToAddress("")), common.HexToAddress("0xbDA5747bFD65F08deb54cb465eB87D40e51B197E"), big.NewInt(100000))
		Expect(err).Should(BeNil())
		fmt.Println("tx_hash", tx.Hash().Hex())
	})
})
