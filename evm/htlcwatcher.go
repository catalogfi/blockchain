package evm

import (
	"bytes"
	"context"
	"math/big"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm/bindings/contracts/htlc/gardenhtlc"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type HTLCEvent interface {
	OrderID() [32]byte
	Equal(e HTLCEvent) bool
}

type htlcwatcher struct {
	client Client

	logger *zap.Logger
}

type HTLCWatcher interface {
	Watch(ctx context.Context, start uint64, asset blockchain.EVMAsset) <-chan HTLCEvent
}

func NewHTLCWatcher(logger *zap.Logger, client Client) HTLCWatcher {
	return &htlcwatcher{client: client, logger: logger}
}

func (w *htlcwatcher) Watch(ctx context.Context, start uint64, asset blockchain.EVMAsset) <-chan HTLCEvent {
	logger := w.logger.With(zap.String("name", string(asset.Chain().Name())))

	events := make(chan HTLCEvent)
	client, ok := w.client.EvmClient(asset.Chain())
	if !ok {
		logger.Error("failed to create a new evm client for the given chain")
		close(events)
		return nil
	}

	filterer, err := gardenhtlc.NewGardenHTLCFilterer(asset.Swapper(), client)
	if err != nil {
		logger.Error("failed to build filterer", zap.Error(err))
		close(events)
		return nil
	}

	initiatedSink := make(chan *gardenhtlc.GardenHTLCInitiated)
	iSub, err := filterer.WatchInitiated(&bind.WatchOpts{Start: &start, Context: ctx}, initiatedSink, nil, nil)
	if err != nil {
		logger.Error("failed to listen to initiated events", zap.Error(err))
		close(events)
		return nil
	}

	redeemedSink := make(chan *gardenhtlc.GardenHTLCRedeemed)
	rSub, err := filterer.WatchRedeemed(&bind.WatchOpts{Start: &start, Context: ctx}, redeemedSink, nil, nil)
	if err != nil {
		logger.Error("failed to listen to redeemed events", zap.Error(err))
		close(events)
		return nil
	}

	refundedSink := make(chan *gardenhtlc.GardenHTLCRefunded)
	rfSub, err := filterer.WatchRefunded(&bind.WatchOpts{Start: &start, Context: ctx}, refundedSink, nil)
	if err != nil {
		logger.Error("failed to listen to refunded events", zap.Error(err))
		close(events)
		return nil
	}

	select {
	case err, ok := <-iSub.Err():
		if !ok {
			logger.Error("initiated event subscription ended")
		}
		if err != nil {
			logger.Error("initiated event subscription ended", zap.Error(err))
		}
		close(events)
		return nil
	case err, ok := <-rSub.Err():
		if !ok {
			logger.Error("redeemed event subscription ended")
		}
		if err != nil {
			logger.Error("redeemed event subscription ended", zap.Error(err))
		}
		close(events)
		return nil
	case err, ok := <-rfSub.Err():
		if !ok {
			logger.Error("refunded event subscription ended")
		}
		if err != nil {
			logger.Error("refunded event subscription ended", zap.Error(err))
		}
		close(events)
		return nil
	case iEvent, ok := <-initiatedSink:
		if !ok {
			logger.Error("initiated event sink closed")
			close(events)
			return nil
		}
		logger.Info("new initiated event detected")
		events <- HTLCInitiated{
			ID:                    iEvent.OrderID,
			InitiateTxHash:        iEvent.Raw.TxHash,
			InitiateTxBlockNumber: iEvent.Raw.BlockNumber,
			Initiator:             iEvent.Initiator,
			Redeemer:              iEvent.Redeemer,
			Amount:                iEvent.Amount,
			Expiry:                iEvent.Expiry,
		}
	case rEvent, ok := <-redeemedSink:
		if !ok {
			logger.Error("redeemed event sink closed")
			close(events)
			return nil
		}
		logger.Info("new redeemed event detected")
		events <- HTLCRedeemed{
			ID:                  rEvent.OrderID,
			Secret:              rEvent.Secret,
			RedeemTxHash:        rEvent.Raw.TxHash,
			RedeemTxBlockNumber: rEvent.Raw.BlockNumber,
		}
	case rfEvent, ok := <-refundedSink:
		if !ok {
			logger.Error("refunded event sink closed")
			close(events)
			return nil
		}
		logger.Info("new refunded event detected")
		events <- HTLCRefunded{
			ID:                  rfEvent.OrderID,
			RefundTxHash:        rfEvent.Raw.TxHash,
			RefundTxBlockNumber: rfEvent.Raw.BlockNumber,
		}
	case <-ctx.Done():
		close(events)
		return nil
	}
	return events
}

type HTLCInitiated struct {
	ID                    [32]byte
	InitiateTxHash        common.Hash
	InitiateTxBlockNumber uint64
	Initiator             common.Address
	Redeemer              common.Address
	Amount                *big.Int
	Expiry                *big.Int
}

func (e HTLCInitiated) OrderID() [32]byte {
	return e.ID
}

func (e HTLCInitiated) Equal(o HTLCEvent) bool {
	other, ok := o.(HTLCInitiated)
	if !ok {
		return ok
	}

	return bytes.Equal(e.ID[:], other.ID[:]) &&
		e.Redeemer.Cmp(other.Redeemer) == 0 &&
		e.Amount.Cmp(other.Amount) == 0 &&
		e.Expiry.Cmp(other.Expiry) == 0 &&
		bytes.Equal(e.InitiateTxHash[:], other.InitiateTxHash[:]) &&
		e.InitiateTxBlockNumber == other.InitiateTxBlockNumber
}

type HTLCRedeemed struct {
	ID                  [32]byte
	Secret              []byte
	RedeemTxHash        common.Hash
	RedeemTxBlockNumber uint64
}

func (e HTLCRedeemed) OrderID() [32]byte {
	return e.ID
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
	ID                  [32]byte
	RefundTxHash        common.Hash
	RefundTxBlockNumber uint64
}

func (e HTLCRefunded) OrderID() [32]byte {
	return e.ID
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
