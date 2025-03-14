package evm

import (
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm/bindings/contracts/orderbook"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	uint256Ty, _   = abi.NewType("uint256", "uint256", nil)
	bytes32Ty, _   = abi.NewType("bytes32", "bytes32", nil)
	stringTy, _    = abi.NewType("string", "string", nil)
	bytesTy, _     = abi.NewType("bytes", "bytes", nil)
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
		{
			Name: "additionalData",
			Type: bytesTy,
		},
	}
)

type CreateOrder struct {
	SourceChain                 string
	DestinationChain            string
	SourceAsset                 string
	DestinationAsset            string
	InitiatorSourceAddress      string
	InitiatorDestinationAddress string
	SourceAmount                *big.Int
	DestinationAmount           *big.Int
	Fee                         *big.Int
	Nonce                       *big.Int
	MinDestinationConfirmations *big.Int
	Timelock                    *big.Int
	SecretHash                  [32]byte
	AdditionalData              []byte
}

type FillOrder struct {
	SourceChain                string
	DestinationChain           string
	SourceAsset                string
	DestinationAsset           string
	RedeemerSourceAddress      string
	RedeemerDestinationAddress string
	SourceAmount               *big.Int
	DestinationAmount          *big.Int
	Fee                        *big.Int
	Nonce                      string
	MinSourceConfirmations     *big.Int
	TimeLock                   *big.Int
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

	data, err := PackCreateOrder(order)
	if err != nil {
		return "", err
	}

	tx, err := contract.CreateOrder(tops, data)

	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func PackCreateOrder(order *CreateOrder) ([]byte, error) {
	data, err := CreateOrderAbi.Pack(
		order.SourceChain,
		order.DestinationChain,
		order.SourceAsset,
		order.DestinationAsset,
		order.InitiatorSourceAddress,
		order.InitiatorDestinationAddress,
		order.SecretHash,
		order.MinDestinationConfirmations,
		order.Timelock,
		order.SourceAmount,
		order.DestinationAmount,
		order.Fee,
		order.Nonce,
		order.AdditionalData,
	)
	if err != nil {
		return nil, err
	}

	data = append(make([]byte, 32), data...)
	binary.BigEndian.PutUint16(data[24:32], uint16(32))

	return data, nil
}

func UnpackCreateOrder(data []byte) (*CreateOrder, error) {
	values, err := CreateOrderAbi.Unpack(data)
	if err != nil {
		return nil, err
	}

	if len(values) != 14 {
		return nil, errors.New("invalid number of values unpacked")
	}

	return &CreateOrder{
		SourceChain:                 string(values[0].(string)),
		DestinationChain:            string(values[1].(string)),
		SourceAsset:                 string(values[2].(string)),
		DestinationAsset:            string(values[3].(string)),
		InitiatorSourceAddress:      values[4].(string),
		InitiatorDestinationAddress: values[5].(string),
		SecretHash:                  values[6].([32]byte),
		MinDestinationConfirmations: values[7].(*big.Int),
		Timelock:                    values[8].(*big.Int),
		SourceAmount:                values[9].(*big.Int),
		DestinationAmount:           values[10].(*big.Int),
		Fee:                         values[11].(*big.Int),
		Nonce:                       values[12].(*big.Int),
		AdditionalData:              values[13].([]byte),
	}, nil
}
