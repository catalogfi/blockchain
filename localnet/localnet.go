package localnet

import (
	"fmt"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func EVMClient() evm.Client {
	client, err := evm.NewClient(evm.Config{RPC: map[string]string{
		string(blockchain.EthereumLocalnet): "http://localhost:8545",
		string(blockchain.ArbitrumLocalnet): "http://localhost:8546",
	}})
	if err != nil {
		panic(fmt.Errorf("failed to create localnet client: %v", err))
	}
	return client
}

func EVMWallet() evm.Wallet {
	key, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		panic(fmt.Errorf("failed to parse localnet privatekey: %v", err))
	}
	return evm.NewWallet(EVMClient(), key)
}

func ETH() blockchain.EVMAsset {
	return blockchain.NewETH(blockchain.NewEvmChain(blockchain.EthereumLocalnet), common.HexToAddress(""))
}

func WBTC() blockchain.EVMAsset {
	return blockchain.NewERC20(blockchain.NewEvmChain(blockchain.EthereumLocalnet), common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"), common.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"))
}

func ArbitrumETH() blockchain.EVMAsset {
	return blockchain.NewETH(blockchain.NewEvmChain(blockchain.ArbitrumLocalnet), common.HexToAddress(""))
}

func ArbitrumWBTC() blockchain.EVMAsset {
	return blockchain.NewERC20(blockchain.NewEvmChain(blockchain.ArbitrumLocalnet), common.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"), common.HexToAddress("0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"))
}
