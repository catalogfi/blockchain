package evm_test

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm"
)

var _ = Describe("Wallet", func() {
	It("should be able to build a client", func() {
		chain := blockchain.NewEvmChain(blockchain.EthereumLocalnet)
		client, err := evm.NewClient(evm.Config{
			RPC: map[string]string{
				string(chain.Name()): "http://127.0.0.1:8545",
			},
		})
		Expect(err).Should(BeNil())
		key, err := crypto.HexToECDSA("de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0")
		Expect(err).Should(BeNil())

		wallet := evm.NewGardenWallet(client, key)

		tx, err := wallet.Send(context.Background(), blockchain.NewETH(chain, common.HexToAddress("")), common.HexToAddress("0xbDA5747bFD65F08deb54cb465eB87D40e51B197E"), big.NewInt(100000))
		Expect(err).Should(BeNil())
		fmt.Println("tx_hash", tx.Hash().Hex())
	})
})
