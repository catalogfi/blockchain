package blockchain

import (
	"github.com/btcsuite/btcd/chaincfg"
)

type UtxoChain struct {
	name Name
}

func NewUtxoChain(name Name) UtxoChain {
	switch name {
	case Bitcoin, BitcoinTestnet, BitcoinRegtest:
		return UtxoChain{name}
	default:
		panic("unknown utxoChain network")
	}
}

func (chain UtxoChain) Name() Name {
	return chain.name
}

func (chain UtxoChain) IsTestnet() bool {
	return chain.name == BitcoinTestnet || chain.name == BitcoinRegtest
}

func (chain UtxoChain) Type() Type {
	return TypeUTXOBased
}

func (chain UtxoChain) Network() Network {
	switch chain.name {
	case Bitcoin:
		return NetworkMainnet
	case BitcoinTestnet:
		return NetworkTestnet
	case BitcoinRegtest:
		return NetworkLocalnet
	default:
		panic("unknown utxoChain network")
	}
}

func (chain UtxoChain) Params() *chaincfg.Params {
	switch chain.name {
	case Bitcoin:
		return &chaincfg.MainNetParams
	case BitcoinTestnet:
		return &chaincfg.TestNet3Params
	case BitcoinRegtest:
		return &chaincfg.RegressionNetParams
	default:
		panic("unknown utxoChain network")
	}
}
