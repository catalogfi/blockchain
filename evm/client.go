package evm

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm/bindings/openzeppelin/contracts/token/ERC20/erc20"
	"github.com/catalogfi/blockchain/evm/bindings/openzeppelin/contracts/token/ERC721/erc721"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type client struct {
	evmClients map[blockchain.EvmChain]*ethclient.Client
}
type Client interface {
	Balance(ctx context.Context, asset blockchain.EVMAsset, owner common.Address, blockNumber *big.Int) (*big.Int, error)
	EvmClient(chain blockchain.EvmChain) (*ethclient.Client, bool)
}

type HTLCClient interface {
	Client

	HTLCEvents(ctx context.Context, asset blockchain.EVMAsset, fromBlock, toBlock *big.Int) ([]HTLCEvent, error)
}

type Config struct {
	RPC map[string]string
}

func NewClient(config Config) (Client, error) {
	return newClient(config)
}

func NewHTLCClient(config Config) (HTLCClient, error) {
	return newClient(config)
}

func newClient(config Config) (*client, error) {
	evmClients := map[blockchain.EvmChain]*ethclient.Client{}
	for chainName, rpc := range config.RPC {
		chain, ok := blockchain.ChainFromName(blockchain.Name(strings.ToLower(chainName))).(blockchain.EvmChain)
		if !ok {
			return nil, fmt.Errorf("unsupported evm chain: %v", chainName)
		}
		client, err := ethclient.Dial(rpc)
		if err != nil {
			return nil, err
		}
		chainID, err := client.ChainID(context.Background())
		if err != nil {
			return nil, err
		}
		if chainID.Cmp(chain.ChainID()) != 0 {
			return nil, fmt.Errorf("invalid rpc url for: %v, chain id mismatch %v != %v", chain.Name(), chainID, chain.ChainID())
		}
		evmClients[chain] = client
	}
	return &client{evmClients: evmClients}, nil
}

func (client *client) Balance(ctx context.Context, asset blockchain.EVMAsset, owner common.Address, blockNumber *big.Int) (*big.Int, error) {
	evmClient, ok := client.evmClients[asset.Chain()]
	if !ok {
		return nil, fmt.Errorf("unsupported chain: %v", asset.Chain().Name())
	}

	switch asset := asset.(type) {
	case blockchain.ERC20:
		caller, err := erc20.NewERC20Caller(asset.Token, evmClient)
		if err != nil {
			return nil, fmt.Errorf("failed to load erc20 bindings: %v", err)
		}
		return caller.BalanceOf(&bind.CallOpts{BlockNumber: blockNumber, Context: ctx}, owner)
	case blockchain.ETH:
		return evmClient.BalanceAt(ctx, owner, blockNumber)
	case blockchain.ERC721:
		caller, err := erc721.NewERC721Caller(asset.Token, evmClient)
		if err != nil {
			return nil, fmt.Errorf("failed to load erc20 bindings: %v", err)
		}
		return caller.BalanceOf(&bind.CallOpts{BlockNumber: blockNumber, Context: ctx}, owner)
	default:
		panic("constraint violation: unsupported evm asset type")
	}
}

func (client *client) EvmClient(chain blockchain.EvmChain) (*ethclient.Client, bool) {
	c, ok := client.evmClients[chain]
	return c, ok
}
