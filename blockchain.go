package blockchain

import (
	"fmt"
)

// Name is a string identifier of different blockchains
type Name = string

// Name of all supported chains.
// format = ("%v_%v", chain, network), all lower cases and network will be omitted if it's mainnet
const (
	Bitcoin         Name = "bitcoin"
	BitcoinTestnet3 Name = "bitcoin_testnet3"
	BitcoinTestnet4 Name = "bitcoin_testnet4"
	BitcoinSignet   Name = "bitcoin_signet"
	BitcoinRegtest  Name = "bitcoin_regtest"

	Ethereum         Name = "ethereum"
	EthereumSepolia  Name = "ethereum_sepolia"
	EthereumLocalnet Name = "ethereum_localnet"

	Base        Name = "base"
	BaseSepolia Name = "base_sepolia"

	Arbitrum         Name = "arbitrum"
	ArbitrumSepolia  Name = "arbitrum_sepolia"
	ArbitrumLocalnet Name = "arbitrum_localnet"

	Bera        Name = "bera"
	BeraBepolia Name = "bera_bepolia"

	HyperEvm        Name = "hyperevm"
	HyperEvmTestnet Name = "hyperevm_testnet"

	CitreaTestnet Name = "citrea_testnet"

	MonadTestnet Name = "monad_testnet"

	Starknet        Name = "starknet"
	StarknetSepolia Name = "starknet_sepolia"
	StarknetDevnet  Name = "starknet_localnet"

	Solana         Name = "solana"
	SolanaDevnet   Name = "solana_devnet"
	SolanaLocalnet Name = "solana_localnet"
)

type Type string

const (
	// TypeEvm is an identifier for all evm-compatible chains,
	// namely, Ethereum, BinanceSmartChain and so on.
	TypeEvm = Type("evm")

	// TypeSolana is an identifier for all Solana chains.
	TypeSolana = Type("solana")

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
	case Bitcoin, BitcoinTestnet3, BitcoinTestnet4, BitcoinSignet, BitcoinRegtest:
		return NewUtxoChain(name), nil
	case "bitcoin_testnet": // alias for testnet4
		return NewUtxoChain(BitcoinTestnet4), nil
	case Ethereum, EthereumSepolia, EthereumLocalnet,
		Arbitrum, ArbitrumSepolia, ArbitrumLocalnet,
		Base, BaseSepolia,
		Bera, BeraBepolia,
		HyperEvm, HyperEvmTestnet,
		CitreaTestnet,
		MonadTestnet,
		Starknet, StarknetSepolia, StarknetDevnet:
		return NewEvmChain(name), nil
	case Solana, SolanaDevnet, SolanaLocalnet:
		return NewSolanaChain(name), nil
	default:
		return nil, fmt.Errorf("unsupported chain = %v", name)
	}
}
