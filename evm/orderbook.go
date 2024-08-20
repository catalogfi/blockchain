package evm

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm/bindings/contracts/orderbook"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	uint256Ty, _ = abi.NewType("uint256", "uint256", nil)
	bytes32Ty, _ = abi.NewType("bytes32", "bytes32", nil)
	stringTy, _  = abi.NewType("string", "string", nil)
)

type CreateOrder struct {
	ID                          string
	CreatedAt                   time.Time
	SourceChain                 blockchain.Chain
	DestinationChain            blockchain.Chain
	SourceAsset                 blockchain.Asset
	DestinationAsset            blockchain.Asset
	InitiatorSourceAddress      string
	InitiatorDestinationAddress string
	SourceAmount                *big.Int
	DestinationAmount           *big.Int
	Fee                         *big.Int
	Nonce                       uint64
	MinDestinationConfirmations uint64
	Timelock                    uint64
	SecretHash                  []byte
}

type FillOrder struct {
	ID                         string
	CreatedAt                  time.Time
	SourceChain                blockchain.Chain
	DestinationChain           blockchain.Chain
	SourceAsset                blockchain.Asset
	DestinationAsset           blockchain.Asset
	RedeemerSourceAddress      string
	RedeemerDestinationAddress string
	SourceAmount               *big.Int
	DestinationAmount          *big.Int
	Fee                        string
	Nonce                      string
	MinSourceConfirmations     uint64
	Timelock                   uint64
}

type OrderbookWallet interface {
	Wallet

	CreateOrder(ctx context.Context, chain blockchain.Chain, orderbookAddr common.Address, order *CreateOrder) (string, error)
	// todo : add fill order
	// FillOrder(ctx context.Context, chain blockchain.Chain, orderbookAddr common.Address, order *FillOrder) (string, error)
}

func NewOrderbookWallet(client Client, key *ecdsa.PrivateKey) OrderbookWallet {
	return &wallet{Client: client, privateKey: key}
}

func (c *wallet) CreateOrder(ctx context.Context, chain blockchain.Chain, orderbookAddr common.Address, order *CreateOrder) (string, error) {
	client, tops, err := c.transactor(ctx, chain)
	if err != nil {
		return "", err
	}
	contract, err := orderbook.NewOrderbook(orderbookAddr, client)
	if err != nil {
		return "", err
	}

	data, err := getCreateOrderTypedData(order)
	if err != nil {
		return "", err
	}

	tx, err := contract.CreateOrder(tops, data)

	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func getCreateOrderTypedData(order *CreateOrder) ([]byte, error) {
	createOrderAbi := abi.Arguments{
		{
			Name: "sourceChain",
			Type: stringTy,
		},
		{
			Name: "destinationChain",
			Type: stringTy,
		},
		{
			Name: "sourceAsset",
			Type: stringTy,
		},
		{
			Name: "destinationAsset",
			Type: stringTy,
		},
		{
			Name: "initiatorSourceAddress",
			Type: stringTy,
		},
		{
			Name: "initiatorDestinationAddress",
			Type: stringTy,
		},
		{
			Name: "secretHash",
			Type: bytes32Ty,
		},
		{
			Name: "minConfirmations",
			Type: uint256Ty,
		},
		{
			Name: "timelock",
			Type: uint256Ty,
		},
		{
			Name: "sourceAmount",
			Type: uint256Ty,
		},
		{
			Name: "destinationAmount",
			Type: uint256Ty,
		},
		{
			Name: "fee",
			Type: uint256Ty,
		},
		{
			Name: "nonce",
			Type: uint256Ty,
		},
	}

	sourceAssetId, err := ParseAssetId(order.SourceAsset)
	if err != nil {
		return nil, err
	}

	destinationAssetId, err := ParseAssetId(order.DestinationAsset)
	if err != nil {
		return nil, err
	}

	return createOrderAbi.Pack(
		order.SourceChain.Name(),
		order.DestinationChain.Name(),
		sourceAssetId,
		destinationAssetId,
		order.InitiatorSourceAddress,
		order.InitiatorDestinationAddress,
		order.SecretHash,
		order.MinDestinationConfirmations,
		order.Timelock,
		order.SourceAmount,
		order.DestinationAmount,
		order.Fee,
		order.Nonce,
	)

}

func ParseAssetId(asset blockchain.Asset) (string, error) {
	switch a := asset.(type) {
	case blockchain.EVMAsset:
		return a.Swapper().Hex(), nil
	case blockchain.UtxoAsset:
		return a.String(), nil
	default:
		return "", fmt.Errorf("unsupported asset type: %T", asset)
	}
}
