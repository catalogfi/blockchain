package evm

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm/bindings/contracts/htlc/gardenhtlc"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type HTLCEvent interface {
	OrderID() [32]byte
	Equal(e HTLCEvent) bool
	BlockNumber() uint64
	TxHash() common.Hash
	EVMAsset() blockchain.EVMAsset
}

func (client *client) HTLCEvents(ctx context.Context, asset blockchain.EVMAsset, fromBlock, toBlock *big.Int) ([]HTLCEvent, error) {
	parsed, err := gardenhtlc.GardenHTLCMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get abi: %v", err)
	}

	ethClient, ok := client.EvmClient(asset.Chain())
	if !ok {
		return nil, fmt.Errorf("unsupported chain: %v", asset.Chain())
	}

	logs, err := ethClient.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []common.Address{asset.Swapper()},
		Topics:    [][]common.Hash{{parsed.Events["Initiated"].ID, parsed.Events["Redeemed"].ID, parsed.Events["Refunded"].ID}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to filter logs: %v", err)
	}

	htlc, err := gardenhtlc.NewGardenHTLCFilterer(asset.Swapper(), ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to build filterer: %v", err)
	}

	events := []HTLCEvent{}
	for _, log := range logs {
		switch log.Topics[0].Hex() {
		case parsed.Events["Initiated"].ID.Hex():
			initated, err := htlc.ParseInitiated(log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse initiated events: %v", err)
			}
			events = append(events, HTLCInitiated{
				Asset:                 asset,
				ID:                    initated.OrderID,
				SecretHash:            initated.SecretHash,
				InitiateTxHash:        initated.Raw.TxHash,
				InitiateTxBlockNumber: initated.Raw.BlockNumber,
				Amount:                initated.Amount,
			})
		case parsed.Events["Redeemed"].ID.Hex():
			redeemed, err := htlc.ParseRedeemed(log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse redeemed events: %v", err)
			}
			events = append(events, HTLCRedeemed{
				Asset:               asset,
				ID:                  redeemed.OrderID,
				Secret:              redeemed.Secret,
				RedeemTxHash:        redeemed.Raw.TxHash,
				RedeemTxBlockNumber: redeemed.Raw.BlockNumber,
				SecretHash:          redeemed.SecretHash,
			})
		case parsed.Events["Refunded"].ID.Hex():
			refunded, err := htlc.ParseRefunded(log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse refunded events: %v", err)
			}
			events = append(events, HTLCRefunded{
				Asset:               asset,
				ID:                  refunded.OrderID,
				RefundTxHash:        refunded.Raw.TxHash,
				RefundTxBlockNumber: refunded.Raw.BlockNumber,
			})
		}
	}
	return events, nil
}

type HTLCInitiated struct {
	Asset                 blockchain.EVMAsset
	ID                    [32]byte
	SecretHash            [32]byte
	InitiateTxHash        common.Hash
	InitiateTxBlockNumber uint64
	Amount                *big.Int
}

func (e HTLCInitiated) OrderID() [32]byte {
	return e.ID
}

func (e HTLCInitiated) BlockNumber() uint64 {
	return e.InitiateTxBlockNumber
}

func (e HTLCInitiated) EVMAsset() blockchain.EVMAsset {
	return e.Asset
}

func (e HTLCInitiated) TxHash() common.Hash {
	return e.InitiateTxHash
}

func (e HTLCInitiated) Equal(o HTLCEvent) bool {
	other, ok := o.(HTLCInitiated)
	if !ok {
		return ok
	}

	return bytes.Equal(e.ID[:], other.ID[:]) &&
		e.Amount.Cmp(other.Amount) == 0 &&
		bytes.Equal(e.InitiateTxHash[:], other.InitiateTxHash[:]) &&
		e.InitiateTxBlockNumber == other.InitiateTxBlockNumber
}

type HTLCRedeemed struct {
	Asset               blockchain.EVMAsset
	ID                  [32]byte
	SecretHash          [32]byte
	Secret              []byte
	RedeemTxHash        common.Hash
	RedeemTxBlockNumber uint64
}

func (e HTLCRedeemed) OrderID() [32]byte {
	return e.ID
}

func (e HTLCRedeemed) BlockNumber() uint64 {
	return e.RedeemTxBlockNumber
}

func (e HTLCRedeemed) EVMAsset() blockchain.EVMAsset {
	return e.Asset
}

func (e HTLCRedeemed) TxHash() common.Hash {
	return e.RedeemTxHash
}

func (event HTLCRedeemed) Equal(o HTLCEvent) bool {
	other, ok := o.(HTLCRedeemed)
	if !ok {
		return ok
	}
	return bytes.Equal(event.ID[:], other.ID[:]) &&
		event.RedeemTxHash.Cmp(other.RedeemTxHash) == 0 &&
		bytes.Equal(event.Secret[:], other.Secret[:]) &&
		event.RedeemTxBlockNumber == other.RedeemTxBlockNumber
}

type HTLCRefunded struct {
	Asset               blockchain.EVMAsset
	ID                  [32]byte
	RefundTxHash        common.Hash
	RefundTxBlockNumber uint64
}

func (e HTLCRefunded) OrderID() [32]byte {
	return e.ID
}

func (e HTLCRefunded) BlockNumber() uint64 {
	return e.RefundTxBlockNumber
}

func (e HTLCRefunded) EVMAsset() blockchain.EVMAsset {
	return e.Asset
}

func (e HTLCRefunded) TxHash() common.Hash {
	return e.RefundTxHash
}

func (event HTLCRefunded) Equal(o HTLCEvent) bool {
	other, ok := o.(HTLCRefunded)
	if !ok {
		return ok
	}
	return bytes.Equal(event.ID[:], other.ID[:]) &&
		event.RefundTxHash.Cmp(other.RefundTxHash) == 0 &&
		event.RefundTxBlockNumber == other.RefundTxBlockNumber
}
