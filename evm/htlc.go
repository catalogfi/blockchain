package evm

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm/bindings/contracts/htlc/gardenhtlc"
	"github.com/catalogfi/blockchain/evm/bindings/openzeppelin/contracts/token/ERC20/erc20"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type HTLCWallet interface {
	Wallet

	OrderID(secretHash [32]byte) [32]byte
	Initiate(ctx context.Context, asset blockchain.EVMAsset, redeemer common.Address, secretHash [32]byte, expiry *big.Int, amount *big.Int, sig []byte) (*types.Receipt, error)
	Redeem(ctx context.Context, asset blockchain.EVMAsset, orderID [32]byte, secret []byte) (*types.Receipt, error)
	Refund(ctx context.Context, asset blockchain.EVMAsset, orderID [32]byte, sig []byte) (*types.Receipt, error)
}

type HTLCClient interface {
	Client

	HTLCEvents(ctx context.Context, asset blockchain.EVMAsset, fromBlock, toBlock *big.Int) ([]HTLCEvent, error)
}

func NewHTLCClient(config Config) (HTLCClient, error) {
	return newClient(config)
}

func NewHTLCWallet(client HTLCClient, key *ecdsa.PrivateKey) HTLCWallet {
	return &wallet{Client: client, privateKey: key}
}

func (w *wallet) SignInitiate(ctx context.Context, asset blockchain.EVMAsset, redeemer common.Address, secretHash [32]byte, expiry *big.Int, amount *big.Int, sig []byte) ([]byte, error) {
	return nil, nil
}

func (w *wallet) OrderID(secretHash [32]byte) [32]byte {
	return sha256.Sum256(append(secretHash[:], common.BytesToHash(w.Address().Bytes()).Bytes()...))
}

func (w *wallet) Initiate(ctx context.Context, asset blockchain.EVMAsset, redeemer common.Address, secretHash [32]byte, expiry *big.Int, amount *big.Int, sig []byte) (*types.Receipt, error) {
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
		erc20, err := erc20.NewERC20(asset.Token, client)
		if err != nil {
			return nil, err
		}
		allowance, err := erc20.Allowance(&bind.CallOpts{}, w.Address(), asset.Swapper())
		if err != nil {
			return nil, err
		}
		if allowance.Cmp(amount) < 0 {
			tx, err := erc20.Approve(tops, asset.Swapper(), MaxETHAmount)
			if err != nil {
				return nil, err
			}
			_, err = bind.WaitMined(ctx, client, tx)
			if err != nil {
				return nil, err
			}
		}
		var tx *types.Transaction
		if sig != nil {
			tx, err = htlc.InitiateWithSignature(tops, redeemer, expiry, amount, secretHash, sig)
		} else {
			tx, err = htlc.Initiate(tops, redeemer, expiry, amount, secretHash)
		}
		if err != nil {
			return nil, err
		}
		return bind.WaitMined(ctx, client, tx)
	default:
		panic(fmt.Sprintf("unsupported asset type: %T", asset))
	}
}

func (w *wallet) Redeem(ctx context.Context, asset blockchain.EVMAsset, orderID [32]byte, secret []byte) (*types.Receipt, error) {
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
		tx, err := htlc.Redeem(tops, orderID, secret)
		if err != nil {
			return nil, err
		}
		return bind.WaitMined(ctx, client, tx)
	default:
		panic(fmt.Sprintf("unsupported asset type: %T", asset))
	}
}

func (w *wallet) Refund(ctx context.Context, asset blockchain.EVMAsset, orderID [32]byte, sig []byte) (*types.Receipt, error) {
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
		var tx *types.Transaction
		if sig != nil {
			tx, err = htlc.InstantRefund(tops, orderID, sig)
		} else {
			tx, err = htlc.Refund(tops, orderID)
		}
		if err != nil {
			return nil, err
		}
		return bind.WaitMined(ctx, client, tx)
	default:
		panic(fmt.Sprintf("unsupported asset type: %T", asset))
	}
}

type HTLCEvent interface {
	OrderID() [32]byte
	Equal(e HTLCEvent) bool
	BlockNumber() uint64
	TxHash() common.Hash
	BlockchainAsset() blockchain.EVMAsset
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

func (e HTLCInitiated) BlockchainAsset() blockchain.EVMAsset {
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

func (e HTLCRedeemed) BlockchainAsset() blockchain.EVMAsset {
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

func (e HTLCRefunded) BlockchainAsset() blockchain.EVMAsset {
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
