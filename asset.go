package blockchain

import "github.com/ethereum/go-ethereum/common"

var Native = ""

type Asset struct {
	Chain    Chain
	Address  string
	Symbol   string
	Decimals int
}

func NewAsset(chain Chain, address, symbol string, decimals int) Asset {
	return Asset{
		Chain:    chain,
		Address:  address,
		Decimals: decimals,
		Symbol:   symbol,
	}
}

// Assets on bitcoin chains
var (
	AssetBtcMainnet = Asset{
		Chain:    NewUtxoChain(Bitcoin),
		Address:  Native,
		Decimals: 8,
		Symbol:   "BTC",
	}
	AssetBtcTestnet3 = Asset{
		Chain:    NewUtxoChain(BitcoinTestnet3),
		Address:  Native,
		Decimals: 8,
		Symbol:   "BTC",
	}
	AssetBtcTestnet4 = Asset{
		Chain:    NewUtxoChain(BitcoinTestnet4),
		Address:  Native,
		Decimals: 8,
		Symbol:   "BTC",
	}
	AssetBtcSignet = Asset{
		Chain:    NewUtxoChain(BitcoinSignet),
		Address:  Native,
		Decimals: 8,
		Symbol:   "BTC",
	}
	AssetBtcLocalnet = Asset{
		Chain:    NewUtxoChain(BitcoinRegtest),
		Address:  Native,
		Decimals: 8,
		Symbol:   "BTC",
	}
)

// Assets on Ethereum chains
var (
	// Mainnet

	AssetEthEthereum = Asset{
		Chain:    NewEvmChain(Ethereum),
		Address:  Native,
		Decimals: 18,
		Symbol:   "ETH",
	}
	AssetWbtcEthereum = Asset{
		Chain:    NewEvmChain(Ethereum),
		Address:  "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599",
		Decimals: 8,
		Symbol:   "WBTC",
	}
	AssetUsdcEthereum = Asset{
		Chain:    NewEvmChain(Ethereum),
		Address:  "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		Decimals: 6,
		Symbol:   "USDC",
	}
	AssetCbbtcEthereum = Asset{
		Chain:    NewEvmChain(Ethereum),
		Address:  "0xcbB7C0000aB88B473b1f5aFd9ef808440eed33Bf",
		Decimals: 8,
		Symbol:   "cbBTC",
	}
	AssetIbtcEthereum = Asset{
		Chain:    NewEvmChain(Ethereum),
		Address:  "0x20157DBAbb84e3BBFE68C349d0d44E48AE7B5AD2",
		Decimals: 8,
		Symbol:   "IBTC",
	}

	// Testnet

	AssetEthSepolia = Asset{
		Chain:    NewEvmChain(EthereumSepolia),
		Address:  Native,
		Decimals: 18,
		Symbol:   "ETH",
	}
	AssetWbtcSepolia = Asset{
		Chain:    NewEvmChain(EthereumSepolia),
		Address:  "0x4D68da063577F98C55166c7AF6955cF58a97b20A",
		Decimals: 8,
		Symbol:   "WBTC",
	}
)

// Asset on Arbitrum
var (
	// Mainnet

	AssetEthArbitrum = Asset{
		Chain:    NewEvmChain(Arbitrum),
		Address:  Native,
		Decimals: 18,
		Symbol:   "ETH",
	}
	AssetWbtcArbitrum = Asset{
		Chain:    NewEvmChain(Arbitrum),
		Address:  "0x2f2a2543B76A4166549F7aaB2e75Bef0aefC5B0f",
		Decimals: 8,
		Symbol:   "WBTC",
	}
	AssetUsdcArbitrum = Asset{
		Chain:    NewEvmChain(Arbitrum),
		Address:  "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
		Decimals: 6,
		Symbol:   "USDC",
	}
	AssetIbtcArbitrum = Asset{
		Chain:    NewEvmChain(Arbitrum),
		Address:  "0x050C24dBf1eEc17babE5fc585F06116A259CC77A",
		Decimals: 8,
		Symbol:   "IBTC",
	}

	// Testnet

	AssetEthArbitrumSepolia = Asset{
		Chain:    NewEvmChain(ArbitrumSepolia),
		Address:  Native,
		Decimals: 18,
		Symbol:   "ETH",
	}
)

// Asset on Base
var (
	// Mainnet

	AssetEthBase = Asset{
		Chain:    NewEvmChain(Base),
		Address:  Native,
		Decimals: 18,
		Symbol:   "ETH",
	}
	AssetUsdcBase = Asset{
		Chain:    NewEvmChain(Base),
		Address:  "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
		Decimals: 6,
		Symbol:   "USDC",
	}
	AssetCbbtcBase = Asset{
		Chain:    NewEvmChain(Base),
		Address:  "0xcbB7C0000aB88B473b1f5aFd9ef808440eed33Bf",
		Decimals: 8,
		Symbol:   "cbBTC",
	}

	// Testnet

	AssetEthBaseSepolia = Asset{
		Chain:    NewEvmChain(BaseSepolia),
		Address:  Native,
		Decimals: 18,
		Symbol:   "ETH",
	}
)

// Asset on Bera

var (
	// Mainnet

	AssetBeraBera = Asset{
		Chain:    NewEvmChain(Bera),
		Address:  Native,
		Decimals: 18,
		Symbol:   "BERA",
	}
	AssetLbtcBera = Asset{
		Chain:    NewEvmChain(Bera),
		Address:  "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
		Decimals: 8,
		Symbol:   "LBTC",
	}

	// Testnet

	AssetBeraBepolia = Asset{
		Chain:    NewEvmChain(BeraBepolia),
		Address:  Native,
		Decimals: 18,
		Symbol:   "BERA",
	}
)

// Asset on HyperEvm

var (
	// Mainnet

	AssetHypeHyperEvm = Asset{
		Chain:    NewEvmChain(HyperEvm),
		Address:  Native,
		Decimals: 18,
		Symbol:   "HYPE",
	}
	AssetUbtcHyperEvm = Asset{
		Chain:    NewEvmChain(HyperEvm),
		Address:  "0x9FDBdA0A5e284c32744D2f17Ee5c74B284993463",
		Decimals: 8,
		Symbol:   "UBTC",
	}

	// Testnet

	AssetHypeHyperEvmTestnet = Asset{
		Chain:    NewEvmChain(HyperEvmTestnet),
		Address:  Native,
		Decimals: 18,
		Symbol:   "HYPE",
	}
)

// Asset on Unichain
var (
	// Mainnet

	AssetEthUnichain = Asset{
		Chain:    NewEvmChain(Unichain),
		Address:  Native,
		Decimals: 18,
		Symbol:   "ETH",
	}
	AssetWbtcUnichain = Asset{
		Chain:    NewEvmChain(Unichain),
		Address:  "0x927B51f251480a681271180DA4de28D44EC4AfB8",
		Decimals: 8,
		Symbol:   "WBTC",
	}
)

func HtlcContractAddress(asset Asset) common.Address {
	switch {
	case asset.Chain.Name() == Ethereum:
		switch asset.Symbol {
		case AssetWbtcEthereum.Symbol:
			return common.HexToAddress("0x795Dcb58d1cd4789169D5F938Ea05E17ecEB68cA")
		case AssetUsdcEthereum.Symbol:
			return common.HexToAddress("0xD8a6E3FCA403d79b6AD6216b60527F51cc967D39")
		case AssetCbbtcEthereum.Symbol:
			return common.HexToAddress("0xeaE7721d779276eb0f5837e2fE260118724a2Ba4")
		case AssetIbtcEthereum.Symbol:
			return common.HexToAddress("0xDC74a45e86DEdf1fF7c6dac77e0c2F082f9E4F72")
		}
	case asset.Chain.Name() == Arbitrum:
		switch asset.Symbol {
		case AssetWbtcArbitrum.Symbol:
			return common.HexToAddress("0x6b6303fAb8eC7232b4f2a7b9fa58E5216F608fcb")
		case AssetUsdcArbitrum.Symbol:
			return common.HexToAddress("0xeaE7721d779276eb0f5837e2fE260118724a2Ba4")
		case AssetIbtcArbitrum.Symbol:
			return common.HexToAddress("0xDC74a45e86DEdf1fF7c6dac77e0c2F082f9E4F72")
		}
	case asset.Chain.Name() == Base:
		switch asset.Symbol {
		case AssetUsdcBase.Symbol:
			return common.HexToAddress("0xD8a6E3FCA403d79b6AD6216b60527F51cc967D39")
		case AssetCbbtcBase.Symbol:
			return common.HexToAddress("0xeaE7721d779276eb0f5837e2fE260118724a2Ba4")
		}
	case asset.Chain.Name() == Bera:
		switch asset.Symbol {
		case AssetLbtcBera.Symbol:
			return common.HexToAddress("0x39f3294352208905fc6ebf033954E6c6455CdB4C")
		}
	case asset.Chain.Name() == HyperEvm:
		switch asset.Symbol {
		case AssetUbtcHyperEvm.Symbol:
			return common.HexToAddress("0x3fDEe07b0756651152BF11c8D170D72d7eBbEc49")
		}
	}

	return common.Address{}
}
