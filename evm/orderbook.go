package evm

import (
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm/bindings/contracts/orderbook"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	uint256Ty, _   = abi.NewType("uint256", "uint256", nil)
	bytes32Ty, _   = abi.NewType("bytes32", "bytes32", nil)
	stringTy, _    = abi.NewType("string", "string", nil)
	CreateOrderAbi = abi.Arguments{
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
	SecretHash                  [32]byte
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

	data, err := GetCreateOrderTypedData(order)
	if err != nil {
		return "", err
	}

	tx, err := contract.CreateOrder(tops, data)

	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func GetCreateOrderTypedData(order *CreateOrder) ([]byte, error) {

	sourceAssetId, err := ParseAssetId(order.SourceAsset)
	if err != nil {
		return nil, err
	}

	destinationAssetId, err := ParseAssetId(order.DestinationAsset)
	if err != nil {
		return nil, err
	}

	data, err := CreateOrderAbi.Pack(
		string(order.SourceChain.Name()),
		string(order.DestinationChain.Name()),
		sourceAssetId,
		destinationAssetId,
		order.InitiatorSourceAddress,
		order.InitiatorDestinationAddress,
		order.SecretHash,
		big.NewInt(int64(order.MinDestinationConfirmations)),
		big.NewInt(int64(order.Timelock)),
		order.SourceAmount,
		order.DestinationAmount,
		order.Fee,
		big.NewInt(int64(order.Nonce)),
	)
	if err != nil {
		return nil, err
	}

	data = append(make([]byte, 32), data...)
	binary.BigEndian.PutUint16(data[24:32], uint16(32))

	return data, nil
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
