package blockchain

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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
	case Ethereum, Arbitrum, PolygonZK:
		return NetworkTestnet
	case EthereumSepolia, PolygonZKTestnet:
		return NetworkMainnet
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
	case ArbitrumLocalnet:
		return big.NewInt(31338)
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

func (chain EvmChain) ValidateAddress(address string) error {
	if len(address) == 42 {
		address = address[2:]
	}
	if len(address) == 40 {
		_, err := hex.DecodeString(address)
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("invalid address: %v", address)
}

func (chain EvmChain) L2() bool {
	switch chain.name {
	case Ethereum, EthereumSepolia, EthereumLocalnet:
		return false
	case Arbitrum, PolygonZK, PolygonZKTestnet:
		return true
	default:
		panic(fmt.Sprintf("unknown evm chain = %v", chain))
	}
}

type EVMAsset interface {
	Asset

	Swapper() common.Address
}

func ParseAsset(a string) (Asset, error) {
	return ParseEVMAsset(a)
}

type utxoAsset struct {
	name  string
	chain UtxoChain
}

func (a utxoAsset) String() string {
	return string(a.name)
}

func (a utxoAsset) Chain() Chain {
	return a.chain
}

func ParseBTCAsset(a string) (UtxoAsset, error) {
	return utxoAsset{}, nil
}

func ParseEVMAsset(a string) (EVMAsset, error) {
	vals := strings.Split(a, "-")
	if len(vals) != 2 || len(vals) != 4 {
		return nil, fmt.Errorf("invalid evm asset string: %v", a)
	}
	chain, err := ParseChainName(Name(vals[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid evm asset string: %v", vals[0])
	}
	chain, ok := chain.(EvmChain)
	if !ok {
		return nil, fmt.Errorf("non-evm chain: %v", chain)
	}
	if len(vals) == 2 {
		return NewETH(chain.(EvmChain), common.HexToAddress(vals[1])), nil
	}
	if strings.ToLower(vals[1]) == "erc20" {
		return NewERC20(chain.(EvmChain), common.HexToAddress(vals[2]), common.HexToAddress(vals[3])), nil
	}
	if strings.ToLower(vals[1]) == "erc721" {
		return NewERC721(chain.(EvmChain), common.HexToAddress(vals[2]), common.HexToAddress(vals[3])), nil
	}
	return nil, fmt.Errorf("unsupported evm token standard: %v", vals[1])
}

type ERC20 struct {
	Token common.Address

	chain   EvmChain
	swapper common.Address
}

func NewERC20(chain EvmChain, token, swapper common.Address) EVMAsset {
	return ERC20{token, chain, swapper}
}

func (a ERC20) String() string {
	return fmt.Sprintf("%s-erc20-%s-%s", a.chain.name, a.Token.Hex(), a.swapper.Hex())
}

func (a ERC20) Swapper() common.Address { return a.swapper }
func (a ERC20) Chain() Chain            { return a.chain }

func NewETH(chain EvmChain, swapper common.Address) EVMAsset { return ETH{chain, swapper} }

type ETH struct {
	chain   EvmChain
	swapper common.Address
}

func (a ETH) Swapper() common.Address { return a.swapper }
func (a ETH) Chain() Chain            { return a.chain }
func (a ETH) String() string {
	return fmt.Sprintf("%s-%s", a.chain.name, a.swapper.Hex())
}

type ERC721 struct {
	Token common.Address

	chain   EvmChain
	swapper common.Address
}

func NewERC721(chain EvmChain, token, swapper common.Address) EVMAsset {
	return ERC20{token, chain, swapper}
}

func (a ERC721) Swapper() common.Address { return a.swapper }
func (a ERC721) Chain() Chain            { return a.chain }
func (a ERC721) String() string {
	return fmt.Sprintf("%s-erc721-%s-%s", a.chain.name, a.Token.Hex(), a.swapper.Hex())
}

func NewBTC(chain UtxoChain) UtxoAsset {
	return utxoAsset{name: "btc", chain: chain}
}
