package evm_test

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"math/big"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm"
)

func RandomSecret() []byte {
	secret := [32]byte{}
	_, err := rand.Read(secret[:])
	if err != nil {
		panic(err)
	}
	return secret[:]
}

var _ = Describe("Orderbook", func() {
	It("should be able pack and unpack create order", func() {
		secret := RandomSecret()
		secretHash := sha256.Sum256(secret)

		createOrder := &evm.CreateOrder{
			SourceChain:                 string(blockchain.ArbitrumLocalnet),
			DestinationChain:            string(blockchain.EthereumLocalnet),
			SourceAsset:                 "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9",
			DestinationAsset:            "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512",
			InitiatorSourceAddress:      "0xffffffffffffffffffffffffffffffffffffffff",
			InitiatorDestinationAddress: "0xffffffffffffffffffffffffffffffffffffffff",
			SecretHash:                  secretHash,
			MinDestinationConfirmations: big.NewInt(1),
			Timelock:                    big.NewInt(1000),
			SourceAmount:                big.NewInt(10000),
			DestinationAmount:           big.NewInt(9000),
			Fee:                         big.NewInt(1),
			Nonce:                       big.NewInt(1),
			AdditionalData:              json.RawMessage([]byte(`"bitcoin":false`)),
		}

		typedData, err := evm.PackCreateOrder(createOrder)
		if err != nil {
			Expect(err).Should(BeNil())
		}

		if typedData == nil {
			Expect(typedData).ShouldNot(BeNil())
		}

		unpackedOrder, err := evm.UnpackCreateOrder(typedData[32:])
		if err != nil {
			Expect(err).Should(BeNil())
		}

		Expect(unpackedOrder).ShouldNot(BeNil())
		Expect(unpackedOrder.SourceChain).Should(Equal(createOrder.SourceChain))
		Expect(unpackedOrder.DestinationChain).Should(Equal(createOrder.DestinationChain))
		Expect(unpackedOrder.SourceAsset).Should(Equal(createOrder.SourceAsset))
		Expect(unpackedOrder.DestinationAsset).Should(Equal(createOrder.DestinationAsset))
		Expect(unpackedOrder.InitiatorSourceAddress).Should(Equal(createOrder.InitiatorSourceAddress))
		Expect(unpackedOrder.InitiatorDestinationAddress).Should(Equal(createOrder.InitiatorDestinationAddress))
		Expect(unpackedOrder.SecretHash).Should(Equal(createOrder.SecretHash))
		Expect(unpackedOrder.MinDestinationConfirmations).Should(Equal(createOrder.MinDestinationConfirmations))
		Expect(unpackedOrder.Timelock).Should(Equal(createOrder.Timelock))
		Expect(unpackedOrder.SourceAmount).Should(Equal(createOrder.SourceAmount))
		Expect(unpackedOrder.DestinationAmount).Should(Equal(createOrder.DestinationAmount))
		Expect(unpackedOrder.Fee).Should(Equal(createOrder.Fee))
		Expect(unpackedOrder.Nonce).Should(Equal(createOrder.Nonce))
		Expect(string(unpackedOrder.AdditionalData)).Should(Equal(string(createOrder.AdditionalData)))
	})
})
