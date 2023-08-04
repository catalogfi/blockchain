package eth_test

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"

	"github.com/catalogfi/multichain/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("etheruem", func() {
	Context("Gas estimation for user operation", func() {
		It("should return an estimation of the given userOp", func() {
			client, err := eth.NewAccountAbstractionClient(eth.NetworkSepolia, "https://eth-sepolia.g.alchemy.com/v2/Vze6SOH6IOj-IGxICn7DHuslhAP-kt8g")
			Expect(err).Should(BeNil())

			callData, err := hex.DecodeString("b61d27f60000000000000000000000000be71941d041a32fe7df4a61eb2fcff3b03502c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000004d087d28800000000000000000000000000000000000000000000000000000000")
			Expect(err).Should(BeNil())

			userOp := userop.UserOperation{
				Sender:   common.HexToAddress("0x0be71941D041a32fe7DF4A61Eb2fCff3b03502C0"),
				Nonce:    big.NewInt(43),
				CallData: callData,
			}
			fee, err := client.EstimateUserOperationGas(context.Background(), userOp)
			Expect(err).Should(BeNil())
			log.Print(fee)

		})
	})
})
