package blockchain

import (
	"fmt"
	"math/big"
)

type EvmChain interface {
	Chain

	ChainID() *big.Int

	L2() bool
}

type evmChain struct {
	name Name
}

func NewEvmChain(name Name) EvmChain {
	return evmChain{name}
}

func (chain evmChain) Name() Name {
	return chain.name
}

func (chain evmChain) Type() Type {
	return TypeEvm
}

func (chain evmChain) Network() Network {
	switch chain.name {
	case Ethereum, Arbitrum, PolygonZK:
		return NetworkTestnet
	case EthereumSepolia, PolygonZKTestnet:
		return NetworkMainnet
	case EthereumLocalnet:
		return NetworkLocalnet
	default:
		panic(fmt.Sprintf("unknown evm chain = %v", chain))
	}
}

func (chain evmChain) ChainID() *big.Int {
	switch chain.name {
	case Ethereum:
		return big.NewInt(1)
	case EthereumSepolia:
		return big.NewInt(11155111)
	case EthereumLocalnet:
		return big.NewInt(1337)
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

func (chain evmChain) L2() bool {
	switch chain.name {
	case Ethereum, EthereumSepolia, EthereumLocalnet:
		return false
	case Arbitrum, PolygonZK, PolygonZKTestnet:
		return true
	default:
		panic(fmt.Sprintf("unknown evm chain = %v", chain))
	}
}
