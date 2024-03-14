package blockchain

import (
	"github.com/btcsuite/btcd/chaincfg"
)

type UtxoChain interface {
	Chain
	Params() *chaincfg.Params
}

type utxoChain struct {
	name Name
}

func NewUtxoChain(name Name) UtxoChain {
	switch name {
	case Bitcoin, BitcoinTestnet, BitcoinRegtest:
		return utxoChain{name}
	default:
		panic("unknown utxoChain network")
	}
}

func (chain utxoChain) Name() Name {
	return chain.name
}

func (chain utxoChain) IsTestnet() bool {
	return chain.name == BitcoinTestnet || chain.name == BitcoinRegtest
}

func (chain utxoChain) Type() Type {
	return TypeUTXOBased
}

func (chain utxoChain) Network() Network {
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

func (chain utxoChain) Params() *chaincfg.Params {
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
