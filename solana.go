package blockchain

import (
	"fmt"
	"math/big"
)

type SolanaChain struct {
	name Name
}

func NewSolanaChain(name Name) SolanaChain {
	return SolanaChain{name: name}
}

func (chain SolanaChain) Name() Name {
	return chain.name
}

func (chain SolanaChain) Type() Type {
	return TypeSolana
}

func (chain SolanaChain) Network() Network {
	switch chain.name {
	case Solana:
		return NetworkMainnet
	case SolanaDevnet:
		return NetworkTestnet
	case SolanaLocalnet:
		return NetworkLocalnet
	default:
		panic(fmt.Sprintf("unknown solana chain = %v", chain))
	}
}

func (chain SolanaChain) ChainID() *big.Int {
	switch chain.name {
	case Solana:
		return big.NewInt(101)
	case SolanaDevnet:
		return big.NewInt(103)
	case SolanaLocalnet:
		return big.NewInt(104)
	default:
		panic(fmt.Sprintf("unknown solana chain = %v", chain))
	}
}
