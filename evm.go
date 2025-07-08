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
	case Ethereum, Base, Arbitrum, Bera, HyperEvm, Citrea, Monad, Starknet, Unichain, Botanix, Corn:
		return NetworkMainnet
	case EthereumSepolia, BaseSepolia, ArbitrumSepolia, BeraBepolia, HyperEvmTestnet, CitreaTestnet, MonadTestnet, StarknetSepolia, UnichainSepolia, BotanixTestnet, CornTestnet:
		return NetworkTestnet
	case EthereumLocalnet, ArbitrumLocalnet, StarknetDevnet:
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

	case Base:
		return big.NewInt(8453)
	case BaseSepolia:
		return big.NewInt(84532)

	case Arbitrum:
		return big.NewInt(42161)
	case ArbitrumSepolia:
		return big.NewInt(421614)
	case ArbitrumLocalnet:
		return big.NewInt(31338)

	case Bera:
		return big.NewInt(80094)
	case BeraBepolia:
		return big.NewInt(80069)

	case HyperEvm:
		return big.NewInt(999)
	case HyperEvmTestnet:
		return big.NewInt(998)

	case Citrea:
		panic("unimplemented")
	case CitreaTestnet:
		return big.NewInt(5115)

	case Monad:
		panic("unimplemented")
	case MonadTestnet:
		return big.NewInt(10143)

	case Starknet:
		value, _ := big.NewInt(0).SetString("23448594291968334", 10)
		return value
	case StarknetSepolia, StarknetDevnet:
		value, _ := big.NewInt(0).SetString("393402133025997798000961", 10)
		return value

	case Unichain:
		return big.NewInt(130)
	case UnichainSepolia:
		return big.NewInt(1301)

	case Botanix:
		return big.NewInt(3637)
	case BotanixTestnet:
		return big.NewInt(3633)

	case Corn:
		return big.NewInt(21000000)
	case CornTestnet:
		return big.NewInt(21000001)
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
	case Citrea, CitreaTestnet:
		return true
	case Monad, MonadTestnet:
		return false
	case Starknet, StarknetSepolia, StarknetDevnet:
		return true
	case Unichain, UnichainSepolia:
		return true
	case Botanix, BotanixTestnet:
		return true
	case Corn, CornTestnet:
		return true
	default:
		panic(fmt.Sprintf("unknown evm chain = %v", chain))
	}
}
