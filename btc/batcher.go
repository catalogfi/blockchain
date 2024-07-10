package btc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/log"
)

var (
	AddSignatureOp        = []byte("add_signature")
	SegwitSpendWeight int = txsizes.RedeemP2WPKHInputWitnessWeight - 10 // removes 10vb of overhead for segwit
)
var (
	ErrBatchNotFound             = errors.New("batch not found")
	ErrBatcherStillRunning       = errors.New("batcher is still running")
	ErrBatcherNotRunning         = errors.New("batcher is not running")
	ErrBatchParametersNotMet     = errors.New("batch parameters not met")
	ErrHighFeeEstimate           = errors.New("estimated fee too high")
	ErrFeeDeltaHigh              = errors.New("fee delta too high")
	ErrFeeUpdateNotNeeded        = errors.New("fee update not needed")
	ErrMaxBatchLimitReached      = errors.New("max batch limit reached")
	ErrCPFPFeeUpdateParamsNotMet = errors.New("CPFP fee update parameters not met")
	ErrCPFPBatchingParamsNotMet  = errors.New("CPFP batching parameters not met")
	ErrCPFPBatchingCorrupted     = errors.New("CPFP batching corrupted")
	ErrSavingBatch               = errors.New("failed to save batch")
)

// Batcher is a wallet that runs as a service and batches requests
// into transactions based on the strategy provided
// It is responsible for creating, signing and submitting transactions
// to the network.
type BatcherWallet interface {
	Wallet
	Lifecycle
}

// Lifecycle interface defines the lifecycle of a BatcherWallet
// It provides methods to start, stop and restart the wallet service
type Lifecycle interface {
	Start(ctx context.Context) error
	Stop() error
	Restart(ctx context.Context) error
}

// Cache interface defines the methods that a BatcherWallet's state
// should implement example implementations include in-memory cache and
// rdbs cache
type Cache interface {
	ReadBatch(ctx context.Context, id string) (Batch, error)
	ReadBatchByReqId(ctx context.Context, id string) (Batch, error)
	ReadPendingBatches(ctx context.Context, strategy Strategy) ([]Batch, error)
	SaveBatch(ctx context.Context, batch Batch) error

	ReadRequest(ctx context.Context, id string) (BatcherRequest, error)
	ReadPendingRequests(ctx context.Context) ([]BatcherRequest, error)
	SaveRequest(ctx context.Context, id string, req BatcherRequest) error
}

// Batcher store spend and send requests in a batched request
// and returns a tracking id
type BatcherRequest struct {
	ID     string
	Spends []SpendRequest
	Sends  []SendRequest
	Status bool
}

type BatcherOptions struct {
	PTI       time.Duration // Periodic Time Interval for batching
	TxOptions TxOptions
	Strategy  Strategy
}

// Strategy defines the batching strategy to be used by the BatcherWallet
// It can be one of RBF, CPFP, RBF_CPFP, Multi_CPFP
// RBF - Replace By Fee
// CPFP - Child Pays For Parent
// Multi_CPFP - Multiple CPFP threads are maintained across multiple addresses
type Strategy string

var (
	RBF        Strategy = "RBF"
	CPFP       Strategy = "CPFP"
	RBF_CPFP   Strategy = "RBF_CPFP"
	Multi_CPFP Strategy = "Multi_CPFP"
)

type TxOptions struct {
	MaxOutputs int
	MaxInputs  int

	MaxUnconfirmedAge int

	MaxBatches   int
	MaxBatchSize int

	FeeLevel    FeeLevel
	MaxFeeRate  int
	MinFeeDelta int
	MaxFeeDelta int
}

type batcherWallet struct {
	quit chan struct{}

	address    btcutil.Address
	privateKey *secp256k1.PrivateKey

	opts    BatcherOptions
	indexer IndexerClient
	cache   Cache
}

type Inputs map[string][]RawInputs

type Output struct {
	wire.OutPoint
	Recipient
}

type Outputs map[string][]Output

type FeeUTXOS struct {
	Utxos []UTXO
	Used  map[string]bool
}

type FeeData struct {
	Fee  int64
	Size int
}

type FeeStats struct {
	MaxFeeRate int
	TotalSize  int
	FeeDelta   int
}

type Batch struct {
	Tx          Transaction
	RequestIds  map[string]bool
	IsStable    bool
	IsConfirmed bool
	Strategy    Strategy
	ChangeUtxo  UTXO
}

func NewBatcherWallet(indexer IndexerClient, address btcutil.Address, privateKey *secp256k1.PrivateKey, cache Cache, opts ...func(*batcherWallet) error) (BatcherWallet, error) {
	wallet := &batcherWallet{
		indexer:    indexer,
		address:    address,
		privateKey: privateKey,
		cache:      cache,
		quit:       make(chan struct{}),
	}
	for _, opt := range opts {
		err := opt(wallet)
		if err != nil {
			return nil, err
		}
	}
	return wallet, nil
}

func (w *batcherWallet) Address() btcutil.Address {
	return w.address
}

// Send creates a batch request , saves it in the cache and returns a tracking id
func (w *batcherWallet) Send(ctx context.Context, sends []SendRequest, spends []SpendRequest) (string, error) {
	id := chainhash.HashH([]byte(fmt.Sprintf("%v_%v", spends, sends))).String()
	req := BatcherRequest{
		ID:     id,
		Spends: spends,
		Sends:  sends,
		Status: false,
	}
	return id, w.cache.SaveRequest(ctx, id, req)
}

// Status returns the status of a transaction based on the tracking id
func (w *batcherWallet) Status(ctx context.Context, id string) (Transaction, bool, error) {
	request, err := w.cache.ReadRequest(ctx, id)
	if err != nil {
		return Transaction{}, false, err
	}
	if !request.Status {
		return Transaction{}, false, nil
	}
	batch, err := w.cache.ReadBatchByReqId(ctx, id)
	if err != nil {
		return Transaction{}, false, err
	}

	tx, err := w.indexer.GetTx(ctx, batch.Tx.TxID)
	if err != nil {
		return Transaction{}, false, err
	}
	return tx, true, nil
}

// Start starts the batcher wallet service
func (w *batcherWallet) Start(ctx context.Context) error {
	if w.quit != nil {
		return ErrBatcherStillRunning
	}
	w.run(ctx)
	return nil
}

// Stop gracefully stops the batcher wallet service
func (w *batcherWallet) Stop() error {
	if w.quit == nil {
		return ErrBatcherNotRunning
	}
	close(w.quit)
	w.quit = nil
	return nil
}

// Restart restarts the batcher wallet service
func (w *batcherWallet) Restart(ctx context.Context) error {
	if err := w.Stop(); err != nil {
		return err
	}
	return w.Start(ctx)
}

// starts the batcher based on the strategy
// There are two types of batching triggers
// 1. Periodic Time Interval (PTI) - Batches are created at regular intervals
// 2. Pending Request - Batches are created when a certain number of requests are pending
// 3. Exponential Time Interval (ETI) - Batches are created at exponential intervals but the interval is custom
func (w *batcherWallet) run(ctx context.Context) {
	switch w.opts.Strategy {
	case CPFP:
		w.runPTIBatcher(ctx)
	case RBF:
		w.runPTIBatcher(ctx)
	default:
		panic("strategy not implemented")
	}
}

// PTI stands for Periodic time interval
//  1. It creates a batch at regular intervals
//  2. It also updates the fee rate at regular intervals
//     if fee rate increases more than threshold and there are
//     no batches to create
func (w *batcherWallet) runPTIBatcher(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.opts.PTI)
		for {
			select {
			case <-w.quit:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := w.createBatch(); err != nil {
					if !errors.Is(err, ErrBatchParametersNotMet) {
						log.Error("failed to create batch", "error", err)
					}

					if err := w.updateFeeRate(); err != nil && !errors.Is(err, ErrFeeUpdateNotNeeded) {
						log.Error("failed to update fee rate", "error", err)
					}
				}
			}

		}

	}()
}

// updateFeeRate updates the fee rate based on the strategy
func (w *batcherWallet) updateFeeRate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	feeRate, err := w.indexer.FeeEstimate(ctx)
	if err != nil {
		return err
	}
	requiredFeeRate := selectFee(feeRate, w.opts.TxOptions.FeeLevel)

	feeStats, err := w.getFeeStats(requiredFeeRate)
	if err != nil {
		return err
	}

	switch w.opts.Strategy {
	case CPFP:
		return w.updateCPFP(feeStats, requiredFeeRate)
	// case RBF:
	// 	return w.updateRBF(feeStats, requiredFeeRate)
	default:
		panic("fee update for strategy not implemented")
	}
}

// createBatch creates a batch based on the strategy
func (w *batcherWallet) createBatch() error {
	switch w.opts.Strategy {
	case CPFP:
		return w.createCPFPBatch()
	// case RBF:
	// 	return w.createRBFBatch()
	default:
		panic("batch creation for strategy not implemented")
	}
}

// Generate fee stats based on the required fee rate
// Fee stats are used to determine how much fee is required
// to bump existing batches
func (w *batcherWallet) getFeeStats(requiredFeeRate int) (FeeStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pendingBatches, err := w.cache.ReadPendingBatches(ctx, w.opts.Strategy)
	if err != nil {
		return FeeStats{}, err
	}

	feeStats := calculateFeeStats(requiredFeeRate, pendingBatches)
	if err := w.validateUpdate(feeStats.MaxFeeRate, requiredFeeRate); err != nil {
		return FeeStats{}, err
	}
	return feeStats, nil
}

// verifies if the fee rate delta is within the threshold
func (w *batcherWallet) validateUpdate(currentFeeRate, requiredFeeRate int) error {
	if currentFeeRate > requiredFeeRate {
		return ErrFeeUpdateNotNeeded
	}
	if w.opts.TxOptions.MinFeeDelta > 0 && requiredFeeRate-currentFeeRate < w.opts.TxOptions.MinFeeDelta {
		return ErrFeeUpdateNotNeeded
	}
	if w.opts.TxOptions.MaxFeeDelta > 0 && requiredFeeRate-currentFeeRate > w.opts.TxOptions.MaxFeeDelta {
		return ErrFeeDeltaHigh
	}
	if w.opts.TxOptions.MaxFeeRate > 0 && requiredFeeRate > w.opts.TxOptions.MaxFeeRate {
		return ErrHighFeeEstimate
	}
	return nil
}

// calculates the fee stats based on the required fee rate
func calculateFeeStats(reqFeeRate int, batches []Batch) FeeStats {
	maxFeeRate := int(0)
	totalSize := int(0)
	feeDelta := int(0)

	for _, batch := range batches {
		size := batch.Tx.Weight / 4
		feeRate := int(batch.Tx.Fee) / size
		if feeRate > maxFeeRate {
			maxFeeRate = feeRate
		}
		if reqFeeRate > feeRate {
			feeDelta += (reqFeeRate - feeRate) * size
		}
		totalSize += size
	}
	return FeeStats{
		MaxFeeRate: maxFeeRate,
		TotalSize:  totalSize,
		FeeDelta:   feeDelta,
	}
}

// selects the fee rate based on the fee level option
func selectFee(feeRate FeeSuggestion, feeLevel FeeLevel) int {
	switch feeLevel {
	case MediumFee:
		return feeRate.Medium
	case HighFee:
		return feeRate.High
	case LowFee:
		return feeRate.Low
	default:
		return feeRate.High
	}
}
