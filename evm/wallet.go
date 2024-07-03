package evm

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm/bindings/openzeppelin/contracts/token/ERC20/erc20"
	"github.com/catalogfi/blockchain/evm/bindings/openzeppelin/contracts/token/ERC721/erc721"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var MaxETHAmount, _ = new(big.Int).SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)

type wallet struct {
	Client
	privateKey *ecdsa.PrivateKey
}

type Wallet interface {
	Client

	Address() common.Address
	Send(ctx context.Context, asset blockchain.EVMAsset, to common.Address, amount *big.Int) (*types.Transaction, error)
	SendAll(ctx context.Context, asset blockchain.EVMAsset, to common.Address) (*types.Transaction, error)
}

type GardenWallet interface {
	HTLCWallet
}

func NewWallet(client Client, key *ecdsa.PrivateKey) Wallet {
	return &wallet{Client: client, privateKey: key}
}

func NewGardenWallet(client Client, key *ecdsa.PrivateKey) GardenWallet {
	return &wallet{Client: client, privateKey: key}
}

func (w *wallet) Address() common.Address {
	return crypto.PubkeyToAddress(w.privateKey.PublicKey)
}

func (w *wallet) Send(ctx context.Context, asset blockchain.EVMAsset, to common.Address, amount *big.Int) (*types.Transaction, error) {
	evmChain, ok := asset.Chain().(blockchain.EvmChain)
	if !ok {
		return nil, fmt.Errorf("%v is not a btc chain", asset.Chain().Name())
	}

	client, tops, err := w.transactor(ctx, asset.Chain())
	if err != nil {
		return nil, err
	}
	switch asset := asset.(type) {
	case blockchain.ERC20:
		token, err := erc20.NewERC20(asset.Token, client)
		if err != nil {
			return nil, err
		}
		return token.Transfer(tops, to, amount)
	case blockchain.ERC721:
		nft, err := erc721.NewERC721(asset.Token, client)
		if err != nil {
			return nil, err
		}
		return nft.TransferFrom(tops, tops.From, to, amount)
	case blockchain.ETH:
		nonce, err := client.PendingNonceAt(ctx, tops.From)
		if err != nil {
			return nil, err
		}
		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
		gasTip, err := client.SuggestGasTipCap(ctx)
		if err != nil {
			return nil, err
		}
		signedTx, err := tops.Signer(tops.From, types.NewTx(&types.DynamicFeeTx{
			ChainID:   evmChain.ChainID(),
			Nonce:     nonce,
			To:        &to,
			GasFeeCap: gasPrice,
			GasTipCap: gasTip,
			Gas:       21000,
			Value:     amount,
		}))
		if err != nil {
			return nil, err
		}
		return signedTx, client.SendTransaction(ctx, signedTx)
	default:
		panic(fmt.Sprintf("constraint violation: unsupported asset type: %T", asset))
	}
}

func (w *wallet) SendAll(ctx context.Context, asset blockchain.EVMAsset, to common.Address) (*types.Transaction, error) {
	balance, err := w.Client.Balance(ctx, asset, crypto.PubkeyToAddress(w.privateKey.PublicKey), nil)
	if err != nil {
		return nil, err
	}
	return w.Send(ctx, asset, to, balance)
}

func (w *wallet) transactor(ctx context.Context, chain blockchain.Chain) (*ethclient.Client, *bind.TransactOpts, error) {
	client, ok := w.Client.EvmClient(chain)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported evm chain: %v", chain.Name())
	}
	tops, err := bind.NewKeyedTransactorWithChainID(w.privateKey, chain.(blockchain.EvmChain).ChainID())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create a transactor: %v", err)
	}
	tops.Context = ctx
	return client, tops, nil
}
