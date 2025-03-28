package blockchain

import (
	"fmt"
	"math/big"
)

type EvmChain struct {
	name Name
}

func NewEvmChain(name Name) EvmChain {
	return EvmChain{name: name}
}

func (chain EvmChain) Name() Name {
	return chain.name
}

func (chain EvmChain) Type() Type {
	return TypeEvm
}

func (chain EvmChain) Network() Network {
	switch chain.name {
	case Ethereum, Arbitrum, PolygonZK:
		return NetworkTestnet
	case EthereumSepolia, PolygonZKTestnet:
		return NetworkMainnet
	case EthereumLocalnet, ArbitrumLocalnet:
		return NetworkLocalnet
	default:
		panic(fmt.Sprintf("unknown evm chain = %v", chain))
	}
}

func (chain EvmChain) ChainID() *big.Int {
	switch chain.name {
	case Ethereum:
		return big.NewInt(1)
	case EthereumSepolia:
		return big.NewInt(11155111)
	case EthereumLocalnet:
		return big.NewInt(31337)
	case ArbitrumLocalnet:
		return big.NewInt(31338)
	case Arbitrum:
		return big.NewInt(42161)
	case PolygonZK:
		return big.NewInt(1101)
	case PolygonZKTestnet:
		return big.NewInt(2442)
	default:
		panic(fmt.Sprintf("unknown evm chain = %v", chain))
	}
}

func (chain EvmChain) L2() bool {
	switch chain.name {
	case Ethereum, EthereumSepolia, EthereumLocalnet:
		return false
	case Arbitrum, ArbitrumLocalnet, PolygonZK, PolygonZKTestnet:
		return true
	default:
		panic(fmt.Sprintf("unknown evm chain = %v", chain))
	}
}
