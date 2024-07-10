package btc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
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
	ErrBatchNotFound         = errors.New("batch not found")
	ErrBatcherStillRunning   = errors.New("batcher is still running")
	ErrBatcherNotRunning     = errors.New("batcher is not running")
	ErrBatchParametersNotMet = errors.New("batch parameters not met")
	ErrHighFeeEstimate       = errors.New("estimated fee too high")
	ErrFeeDeltaHigh          = errors.New("fee delta too high")
	ErrFeeUpdateNotNeeded    = errors.New("fee update not needed")
	ErrMaxBatchLimitReached  = errors.New("max batch limit reached")
)

type SpendUTXOs struct {
	Witness       [][]byte
	Script        []byte
	ScriptAddress btcutil.Address
	Utxo          UTXO
	HashType      txscript.SigHashType
}

type BatcherWallet interface {
	Wallet
	Lifecycle
}

type Lifecycle interface {
	Start(ctx context.Context) error
	Stop() error
	Restart(ctx context.Context) error
}

type Cache interface {
	ReadBatch(ctx context.Context, id string) (Batch, error)
	ReadBatchByReqId(ctx context.Context, id string) (Batch, error)
	ReadPendingBatches(ctx context.Context) ([]Batch, error)
	SaveBatch(ctx context.Context, batch Batch) error

	ReadRequest(ctx context.Context, id string) (BatcherRequest, error)
	ReadPendingRequests(ctx context.Context) ([]BatcherRequest, error)
	SaveRequest(ctx context.Context, id string, req BatcherRequest) error
}

type BatcherRequest struct {
	ID     string
	Spends []SpendRequest
	Sends  []SendRequest
	Status bool
}

type BatcherOptions struct {
	PTI       time.Duration
	TxOptions TxOptions
	Strategy  Strategy
}

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
	MaxBatchSize float64

	FeeLevel    FeeLevel
	MaxFeeRate  float64
	MinFeeDelta float64
	MaxFeeDelta float64
}

type batcherWallet struct {
	quit chan struct{}

	address    btcutil.Address
	privateKey *secp256k1.PrivateKey

	opts    BatcherOptions
	indexer IndexerClient
	cache   Cache
}

type Inputs map[btcutil.Address][]RawInputs

type Output struct {
	wire.OutPoint
	Recipient
}

type Outputs map[btcutil.Address][]Output

type FeeUTXOS struct {
	Utxos []UTXO
	Used  map[string]bool
}

type FeeData struct {
	Fee         int64
	FeeRate     float64
	NetFeeRate  float64
	BaseSize    float64
	WitnessSize float64
	FeeUtxos    FeeUTXOS
}

type FeeStats struct {
	MaxFeeRate float64
	TotalSize  float64
	FeeDelta   float64
}

type Batch struct {
	TxId        string
	Inputs      Inputs
	Outputs     Outputs
	TotalIn     int64
	TotalOut    int64
	FeeData     FeeData
	RequestIds  map[string]bool
	IsStable    bool
	IsConfirmed bool
	Transaction Transaction
}

func verifyOptions(opts BatcherOptions) error {
	return nil
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

	tx, err := w.indexer.GetTx(ctx, batch.TxId)
	if err != nil {
		return Transaction{}, false, err
	}
	return tx, true, nil
}

func (w *batcherWallet) Start(ctx context.Context) error {
	if w.quit != nil {
		return ErrBatcherStillRunning
	}
	w.run(ctx)
	return nil
}

func (w *batcherWallet) Stop() error {
	if w.quit == nil {
		return ErrBatcherNotRunning
	}
	close(w.quit)
	w.quit = nil
	return nil
}

func (w *batcherWallet) Restart(ctx context.Context) error {
	if err := w.Stop(); err != nil {
		return err
	}
	return w.Start(ctx)
}

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

func (w *batcherWallet) updateFeeRate() error {
	requiredFeeRate, feeStats, err := w.getFee()
	if err != nil {
		return err
	}

	switch w.opts.Strategy {
	case CPFP:
		return w.updateCPFP(feeStats, requiredFeeRate)
	case RBF:
		return w.updateRBF(feeStats, requiredFeeRate)
	default:
		panic("fee update for strategy not implemented")
	}
}

func (w *batcherWallet) createBatch() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pendingRequests, err := w.cache.ReadPendingRequests(ctx)
	if err != nil {
		return err
	}
	switch w.opts.Strategy {
	case CPFP:
		return w.createCPFPBatch(ctx, pendingRequests)
	case RBF:
		return w.createRBFBatch(ctx, pendingRequests)
	default:
		panic("batch creation for strategy not implemented")
	}
}

func (w *batcherWallet) getFee() (float64, FeeStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	feeRate, err := w.indexer.FeeEstimate(ctx)
	if err != nil {
		return 0, FeeStats{}, err
	}
	requiredFeeRate := selectFee(feeRate, w.opts.TxOptions.FeeLevel)
	pendingBatches, err := w.cache.ReadPendingBatches(ctx)
	if err != nil {
		return 0, FeeStats{}, err
	}

	feeStats := getFeeStats(requiredFeeRate, pendingBatches)
	if err := w.validateUpdate(feeStats.MaxFeeRate, requiredFeeRate); err != nil {
		return 0, FeeStats{}, err
	}
	return requiredFeeRate, feeStats, nil
}

func (w *batcherWallet) validateUpdate(currentFeeRate, requiredFeeRate float64) error {
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

func getFeeStats(feeRate float64, batches []Batch) FeeStats {
	maxFeeRate := float64(0)
	totalSize := float64(0)
	feeDelta := float64(0)

	for _, batch := range batches {
		if batch.FeeData.FeeRate > maxFeeRate {
			maxFeeRate = batch.FeeData.FeeRate
		}
		batchSize := batch.FeeData.BaseSize + batch.FeeData.WitnessSize
		if batch.FeeData.FeeRate > feeRate {
			feeDelta += (batch.FeeData.FeeRate - feeRate) * batchSize
		}
		totalSize += batchSize
	}
	return FeeStats{
		MaxFeeRate: maxFeeRate,
		TotalSize:  totalSize,
		FeeDelta:   feeDelta,
	}
}

func selectFee(feeRate FeeSuggestion, feeLevel FeeLevel) float64 {
	switch feeLevel {
	case MediumFee:
		return float64(feeRate.Medium)
	case HighFee:
		return float64(feeRate.High)
	case LowFee:
		return float64(feeRate.Low)
	default:
		return float64(feeRate.High)
	}
}
