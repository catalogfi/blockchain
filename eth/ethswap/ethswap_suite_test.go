package ethswap_test

import (
	"context"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/catalogfi/blockchain/eth/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	swapAddr  common.Address
	tokenAddr common.Address
)

var _ = BeforeSuite(func() {
	By("Required envs")
	Expect(os.Getenv("ETH_URL")).ShouldNot(BeEmpty())
	Expect(os.Getenv("ETH_KEY_1")).ShouldNot(BeEmpty())

	By("Initialise client")
	url := os.Getenv("ETH_URL")
	client, err := ethclient.Dial(url)
	Expect(err).Should(BeNil())
	chainID, err := client.ChainID(context.Background())
	Expect(err).Should(BeNil())

	By("Initialise transactor")
	keyStr := strings.TrimPrefix(os.Getenv("ETH_KEY_1"), "0x")
	keyBytes, err := hex.DecodeString(keyStr)
	Expect(err).Should(BeNil())
	key, err := crypto.ToECDSA(keyBytes)
	Expect(err).Should(BeNil())
	transactor, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	Expect(err).Should(BeNil())

	By("Deploy ERC20 contract")
	var tx *types.Transaction
	tokenAddr, tx, _, err = bindings.DeployTestERC20(transactor, client)
	Expect(err).Should(BeNil())
	_, err = bind.WaitMined(context.Background(), client, tx)
	Expect(err).Should(BeNil())
	By(color.GreenString("ERC20 deployed to %v", tokenAddr.Hex()))

	By("Deploy atomic swap contract")
	swapAddr, tx, _, err = bindings.DeployAtomicSwap(transactor, client, tokenAddr)
	Expect(err).Should(BeNil())
	_, err = bind.WaitMined(context.Background(), client, tx)
	Expect(err).Should(BeNil())
	By(color.GreenString("Atomic swap deployed to %v", swapAddr.Hex()))
})

func TestEthswap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ethswap Suite")
}
