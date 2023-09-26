package evm_test

import (
	"math/big"

	"github.com/catalogfi/blockchain/evm"
	"github.com/catalogfi/blockchain/evm/typings/TestERC20"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ethereum", func() {
	Context("Client", func() {
		It("should be able to read data from the blockchain", func() {
			By("Initialise client")
			localRPC := "http://127.0.0.1:7545"
			logger, _ := zap.NewDevelopment()
			client, err := evm.NewClient(logger, localRPC)
			Expect(err).Should(BeNil())

			By("GetCurrentBlock()")
			_, err = client.GetCurrentBlock()
			Expect(err).Should(BeNil())
		})

		It("should read ERC20 related data", func() {
			By("Initialise client")
			localRPC := "http://127.0.0.1:7545"
			logger, _ := zap.NewDevelopment()
			client, err := evm.NewClient(logger, localRPC)
			Expect(err).Should(BeNil())

			By("Deploy an ERC20 contract")
			keyStr := "052b060a565a13bbc656314847129b92e36f0882bf08d4a781daee0db25089c2"
			key, err := crypto.HexToECDSA(keyStr)
			owner := crypto.PubkeyToAddress(key.PublicKey)
			Expect(err).Should(BeNil())
			transactor, err := client.GetTransactOpts(key)
			Expect(err).Should(BeNil())
			tokenAddr, _, erc20, err := TestERC20.DeployTestERC20(transactor, client.GetProvider())
			Expect(err).Should(BeNil())

			By("Approve and balance querying")
			balance, err := client.GetERC20Balance(tokenAddr, owner)
			Expect(err).Should(BeNil())
			Expect(balance.Cmp(big.NewInt(0))).Should(BeNumerically(">", 0))
			toKey, err := crypto.HexToECDSA("8341095aedd412f5c4794ce55489172a8091d59456f69374993a0e44b2b6c8bc")
			Expect(err).Should(BeNil())
			to := crypto.PubkeyToAddress(toKey.PublicKey)
			_, err = client.ApproveERC20(key, balance, tokenAddr, to)
			Expect(err).Should(BeNil())
			allowance, err := erc20.Allowance(&bind.CallOpts{}, owner, to)
			Expect(err).Should(BeNil())
			Expect(allowance.Cmp(balance)).Should(BeZero())
		})
	})
})
