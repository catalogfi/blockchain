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
	case Ethereum, Arbitrum, Base, Bera, HyperEvm:
		return NetworkMainnet
	case EthereumSepolia, ArbitrumSepolia, BaseSepolia, BeraBepolia, HyperEvmTestnet, CitreaTestnet, MonadTestnet, StarknetSepolia:
		return NetworkTestnet
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
	case Arbitrum:
		return big.NewInt(42161)
	case ArbitrumSepolia:
		return big.NewInt(421614)
	case ArbitrumLocalnet:
		return big.NewInt(31338)
	case Base:
		return big.NewInt(8453)
	case BaseSepolia:
		return big.NewInt(84532)
	case Bera:
		return big.NewInt(80094)
	case BeraBepolia:
		return big.NewInt(80069)
	case HyperEvm:
		return big.NewInt(999)
	case HyperEvmTestnet:
		return big.NewInt(998)
	case CitreaTestnet:
		return big.NewInt(5115)
	case MonadTestnet:
		return big.NewInt(10143)
	case StarknetSepolia:
		return big.NewInt(393402133025997798000961)
	default:
		panic(fmt.Sprintf("unknown evm chain = %v", chain))
	}
}

func (chain EvmChain) L2() bool {
	switch chain.name {
	case Ethereum, EthereumSepolia, EthereumLocalnet:
		return false
	case Arbitrum, ArbitrumLocalnet, ArbitrumSepolia, Base, BaseSepolia:
		return true
	case Bera, BeraBepolia: // bera is l1
		return false
	case HyperEvm, HyperEvmTestnet: // hyperEvm is l2 for Hyperliquid
		return true
	case CitreaTestnet:
		return false
	case MonadTestnet:
		return false
	default:
		panic(fmt.Sprintf("unknown evm chain = %v", chain))
	}
}
