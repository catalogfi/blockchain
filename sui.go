package blockchain

import (
	"fmt"
	"math/big"
)

type SuiChain struct {
	name Name
}

func NewSuiChain(name Name) SuiChain {
	return SuiChain{name: name}
}

func (chain SuiChain) Name() Name {
	return chain.name
}

func (chain SuiChain) Type() Type {
	return TypeSui
}

func (chain SuiChain) Network() Network {
	switch chain.name {
	case SuiTestnet:
		return NetworkTestnet
	default:
		panic(fmt.Sprintf("unknown Sui chain = %v", chain))
	}
}

func (chain SuiChain) ChainID() *big.Int {
	switch chain.name {
	case SuiTestnet:
		return big.NewInt(0)
	default:
		panic(fmt.Sprintf("unknown Sui chain = %v", chain))
	}
}
