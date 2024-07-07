package blockchain

import (
	"fmt"
)

// Name is a string identifier of different blockchains
type Name string

// Name of all supported chains.
// format = ("%v_%v", chain, network), all lower cases and network will be omitted if it's mainnet
const (
	Bitcoin          Name = "bitcoin"
	BitcoinTestnet   Name = "bitcoin_testnet"
	BitcoinRegtest   Name = "bitcoin_regtest"
	Ethereum         Name = "ethereum"
	EthereumSepolia  Name = "ethereum_sepolia"
	EthereumLocalnet Name = "ethereum_localnet"
	Arbitrum         Name = "arbitrum"
	ArbitrumLocalnet Name = "arbitrum_localnet"
	PolygonZK        Name = "polygonzk"
	PolygonZKTestnet Name = "polygonzk_testnet"

	// TODO : uncomment this when we start supporting them
	// Optimism         Name = "optimism"
	// Polygon          Name = "polygon"
	// Avalanche        Name = "avalanche"
	// BNB              Name = "bnb"
)

type Type string

const (
	// TypeEvm is an identifier for all evm-compatible chains,
	// namely, Ethereum, BinanceSmartChain and so on.
	TypeEvm = Type("evm")

	// TypeUTXOBased is an identifier for all utxo-based chains, namely,
	// Bitcoin, BitcoinCash, Dogecoin, and so on.
	TypeUTXOBased = Type("utxo")
)

type Network string

const (
	// NetworkMainnet usually refers to the production network of the chain. The token of the network usually has real
	// value and needs to be handled carefully.
	NetworkMainnet = Network("mainnet")

	// NetworkTestnet refers to the public testnet of the chain. It usually has same config as mainnet and serves for
	// testing purpose. Token on this network don't have real value, but might not easy to get large amount.
	NetworkTestnet = Network("testnet")

	// NetworkLocalnet refers to the network for local testing. This usually not public accessible and we have the
	// ability to mint infinite amount of tokens.
	NetworkLocalnet = Network("localnet")
)

type Chain interface {
	// Name of the chain, serve as an identifier of the chain.
	Name() Name

	// Type returns the type of the chain, either an evm chain or utxo based chain for now.
	Type() Type

	// Network returns which network type of this chain.
	Network() Network

	// Validate address and return error if the address is invalid for the given chain
	ValidateAddress(string) error
}

type Asset interface {
	String() string

	Chain() Chain
}

func ChainFromName(name Name) Chain {
	chain, err := ParseChainName(name)
	if err != nil {
		panic(err)
	}
	return chain
}

func ParseChainName(name Name) (Chain, error) {
	switch name {
	case Bitcoin, BitcoinTestnet, BitcoinRegtest:
		return NewUtxoChain(name), nil
	case Ethereum, EthereumSepolia, EthereumLocalnet, Arbitrum, ArbitrumLocalnet, PolygonZK, PolygonZKTestnet:
		return NewEvmChain(name), nil
	default:
		return nil, fmt.Errorf("unsupported chain = %v", name)
	}
}
