package evm

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm/bindings/contracts/htlc/gardenhtlc"
	"github.com/catalogfi/blockchain/evm/bindings/openzeppelin/contracts/token/ERC20/erc20"
	"github.com/catalogfi/blockchain/evm/bindings/openzeppelin/contracts/token/ERC721/erc721"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type wallet struct {
	client     Client
	privateKey *ecdsa.PrivateKey
}

type Wallet interface {
	Send(ctx context.Context, asset blockchain.EVMAsset, to common.Address, amount *big.Int) (*types.Transaction, error)
	SendAll(ctx context.Context, asset blockchain.EVMAsset, to common.Address) (*types.Transaction, error)
}

type HTLCWallet interface {
	Wallet

	Initiate(ctx context.Context, asset blockchain.EVMAsset, redeemer common.Address, secretHash [32]byte, expiry *big.Int, amount *big.Int, sig []byte) (*types.Transaction, error)
	Redeem(ctx context.Context, asset blockchain.EVMAsset, orderID [32]byte, secret []byte) (*types.Transaction, error)
	Refund(ctx context.Context, asset blockchain.EVMAsset, orderID [32]byte, sig []byte) (*types.Transaction, error)
}

type GardenWallet interface {
	HTLCWallet
}

func NewWallet(client Client, key *ecdsa.PrivateKey) Wallet {
	return &wallet{client: client, privateKey: key}
}

func NewHTLCWallet(client Client, key *ecdsa.PrivateKey) HTLCWallet {
	return &wallet{client: client, privateKey: key}
}

func NewGardenWallet(client Client, key *ecdsa.PrivateKey) GardenWallet {
	return &wallet{client: client, privateKey: key}
}

func (w *wallet) Send(ctx context.Context, asset blockchain.EVMAsset, to common.Address, amount *big.Int) (*types.Transaction, error) {
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
			ChainID:   asset.Chain().ChainID(),
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
	balance, err := w.client.Balance(ctx, asset, crypto.PubkeyToAddress(w.privateKey.PublicKey), nil)
	if err != nil {
		return nil, err
	}
	return w.Send(ctx, asset, to, balance)
}

func (w *wallet) Initiate(ctx context.Context, asset blockchain.EVMAsset, redeemer common.Address, secretHash [32]byte, expiry *big.Int, amount *big.Int, sig []byte) (*types.Transaction, error) {
	client, tops, err := w.transactor(ctx, asset.Chain())
	if err != nil {
		return nil, err
	}
	switch asset := asset.(type) {
	case blockchain.ERC20:
		htlc, err := gardenhtlc.NewGardenHTLC(asset.Swapper(), client)
		if err != nil {
			return nil, err
		}
		if sig != nil {
			return htlc.InitiateWithSignature(tops, redeemer, expiry, amount, secretHash, sig)
		}
		return htlc.Initiate(tops, redeemer, expiry, amount, secretHash)
	default:
		panic(fmt.Sprintf("unsupported asset type: %T", asset))
	}
}

func (w *wallet) Redeem(ctx context.Context, asset blockchain.EVMAsset, orderID [32]byte, secret []byte) (*types.Transaction, error) {
	client, tops, err := w.transactor(ctx, asset.Chain())
	if err != nil {
		return nil, err
	}
	switch asset := asset.(type) {
	case blockchain.ERC20:
		htlc, err := gardenhtlc.NewGardenHTLC(asset.Swapper(), client)
		if err != nil {
			return nil, err
		}
		return htlc.Redeem(tops, orderID, secret)
	default:
		panic(fmt.Sprintf("unsupported asset type: %T", asset))
	}
}

func (w *wallet) Refund(ctx context.Context, asset blockchain.EVMAsset, orderID [32]byte, sig []byte) (*types.Transaction, error) {
	client, tops, err := w.transactor(ctx, asset.Chain())
	if err != nil {
		return nil, err
	}
	switch asset := asset.(type) {
	case blockchain.ERC20:
		htlc, err := gardenhtlc.NewGardenHTLC(asset.Swapper(), client)
		if err != nil {
			return nil, err
		}
		if sig != nil {
			return htlc.InstantRefund(tops, orderID, sig)
		}
		return htlc.Refund(tops, orderID)
	default:
		panic(fmt.Sprintf("unsupported asset type: %T", asset))
	}
}

func (w *wallet) transactor(ctx context.Context, chain blockchain.EvmChain) (*ethclient.Client, *bind.TransactOpts, error) {
	client, ok := w.client.EvmClient(chain)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported evm chain: %v", chain.Name())
	}
	tops, err := bind.NewKeyedTransactorWithChainID(w.privateKey, chain.ChainID())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create a transactor: %v", err)
	}
	tops.Context = ctx
	return client, tops, nil
}
