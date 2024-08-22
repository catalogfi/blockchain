package evm_test

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm"
	"github.com/ethereum/go-ethereum/common"
)

func TestCreateOrder(t *testing.T) {
	secret := RandomSecret()
	secretHash := sha256.Sum256(secret)

	ethLocalnetChain := blockchain.NewEvmChain(blockchain.EthereumLocalnet)
	arbLocalnetChain := blockchain.NewEvmChain(blockchain.ArbitrumLocalnet)

	ethWBTCAsset := blockchain.NewERC20(ethLocalnetChain, common.HexToAddress(""), common.HexToAddress("0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"))
	arbWBTCAsset := blockchain.NewERC20(arbLocalnetChain, common.HexToAddress(""), common.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"))

	createOrder := &evm.CreateOrder{
		SourceChain:                 arbLocalnetChain,
		DestinationChain:            ethLocalnetChain,
		SourceAsset:                 arbWBTCAsset,
		DestinationAsset:            ethWBTCAsset,
		InitiatorSourceAddress:      "0xffffffffffffffffffffffffffffffffffffffff",
		InitiatorDestinationAddress: "0xffffffffffffffffffffffffffffffffffffffff",
		SecretHash:                  secretHash,
		MinDestinationConfirmations: 1,
		Timelock:                    2,
		SourceAmount:                big.NewInt(10000),
		DestinationAmount:           big.NewInt(9000),
		Fee:                         big.NewInt(1),
		Nonce:                       1,
	}

	typedData, err := evm.GetCreateOrderTypedData(createOrder)
	if err != nil {
		t.Fatal(err)
	}

	if typedData == nil {
		t.Fatal("typedData is nil")
	}

}

func RandomSecret() []byte {
	secret := [32]byte{}
	_, err := rand.Read(secret[:])
	if err != nil {
		panic(err)
	}
	return secret[:]
}
