package btc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"go.uber.org/zap"
)

var (
	AddSignatureOp        = []byte("add_signature")
	SegwitSpendWeight int = txsizes.RedeemP2WPKHInputWitnessWeight 
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
	ErrCPFPBatchingCorrupted     = errors.New("CPFP batching corrupted")
	ErrSavingBatch               = errors.New("failed to save batch")
	ErrStrategyNotSupported      = errors.New("strategy not supported")
	ErrBuildCPFPDepthExceeded    = errors.New("build CPFP depth exceeded")
	ErrBuildRBFDepthExceeded     = errors.New("build RBF depth exceeded")
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
	// ReadBatch reads a batch based on the transaction ID.
	ReadBatch(ctx context.Context, txId string) (Batch, error)
	// ReadBatchByReqId reads a batch based on the request ID.
	ReadBatchByReqId(ctx context.Context, reqId string) (Batch, error)
	// ReadPendingBatches reads all pending batches for a given strategy.
	ReadPendingBatches(ctx context.Context, strategy Strategy) ([]Batch, error)
	// ReadLatestBatch reads the latest batch for a given strategy.
	ReadLatestBatch(ctx context.Context, strategy Strategy) (Batch, error)
	// ReadPendingChangeUtxos reads all pending change UTXOs for a given strategy.
	ReadPendingChangeUtxos(ctx context.Context, strategy Strategy) ([]UTXO, error)
	// ReadPendingFundingUtxos reads all pending funding UTXOs for a given strategy.
	ReadPendingFundingUtxos(ctx context.Context, strategy Strategy) ([]UTXO, error)
	// UpdateBatchStatuses updates the status of multiple batches and delete pending batches based on confirmed transaction IDs.
	UpdateBatchStatuses(ctx context.Context, txId []string, status bool, strategy Strategy) error
	// UpdateBatchFees updates the fees for multiple batches.
	UpdateBatchFees(ctx context.Context, txId []string, fee int64) error
	// SaveBatch saves a batch.
	SaveBatch(ctx context.Context, batch Batch) error
	// DeletePendingBatches deletes pending batches based on confirmed transaction IDs and strategy.
	DeletePendingBatches(ctx context.Context, confirmedTxIds map[string]bool, strategy Strategy) error

	// ReadRequest reads a request based on its ID.
	ReadRequest(ctx context.Context, id string) (BatcherRequest, error)
	// ReadRequests reads multiple requests based on their IDs.
	ReadRequests(ctx context.Context, id []string) ([]BatcherRequest, error)
	// ReadPendingRequests reads all pending requests.
	ReadPendingRequests(ctx context.Context) ([]BatcherRequest, error)
	// SaveRequest saves a request.
	SaveRequest(ctx context.Context, id string, req BatcherRequest) error
}

// Batcher store spend and send requests in a batched request
// and returns a tracking id
type BatcherRequest struct {
	ID     string
	Spends []SpendRequest
	Sends  []SendRequest
	SACPs  [][]byte
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
	wg   sync.WaitGroup

	address    btcutil.Address
	privateKey *secp256k1.PrivateKey
	logger     *zap.Logger

	sw           Wallet
	opts         BatcherOptions
	indexer      IndexerClient
	feeEstimator FeeEstimator
	cache        Cache
}
type Batch struct {
	Tx           Transaction
	RequestIds   map[string]bool
	IsStable     bool
	IsConfirmed  bool
	Strategy     Strategy
	ChangeUtxo   UTXO
	FundingUtxos []UTXO
}

func NewBatcherWallet(privateKey *secp256k1.PrivateKey, indexer IndexerClient, feeEstimator FeeEstimator, chainParams *chaincfg.Params, cache Cache, logger *zap.Logger, opts ...func(*batcherWallet) error) (BatcherWallet, error) {
	address, err := PublicKeyAddress(chainParams, waddrmgr.WitnessPubKey, privateKey.PubKey())
	if err != nil {
		return nil, err
	}

	wallet := &batcherWallet{
		indexer:      indexer,
		address:      address,
		privateKey:   privateKey,
		cache:        cache,
		logger:       logger,
		feeEstimator: feeEstimator,
		opts:         defaultBatcherOptions(),
	}
	for _, opt := range opts {
		err := opt(wallet)
		if err != nil {
			return nil, err
		}
	}

	simpleWallet, err := NewSimpleWallet(privateKey, chainParams, indexer, feeEstimator, wallet.opts.TxOptions.FeeLevel)
	if err != nil {
		return nil, err
	}

	wallet.sw = simpleWallet
	return wallet, nil
}

func defaultBatcherOptions() BatcherOptions {
	return BatcherOptions{
		PTI: 1 * time.Minute,
		TxOptions: TxOptions{
			MaxOutputs: 0,
			MaxInputs:  0,

			MaxUnconfirmedAge: 0,

			MaxBatches:   0,
			MaxBatchSize: 0,

			FeeLevel:    HighFee,
			MaxFeeRate:  0,
			MinFeeDelta: 0,
			MaxFeeDelta: 0,
		},
		Strategy: RBF,
	}
}

func WithStrategy(strategy Strategy) func(*batcherWallet) error {
	return func(w *batcherWallet) error {
		err := parseStrategy(strategy)
		if err != nil {
			return err
		}
		w.opts.Strategy = strategy
		return nil
	}
}

func WithPTI(pti time.Duration) func(*batcherWallet) error {
	return func(w *batcherWallet) error {
		w.opts.PTI = pti
		return nil
	}
}

func parseStrategy(strategy Strategy) error {
	switch strategy {
	case RBF, CPFP, RBF_CPFP, Multi_CPFP:
		return nil
	default:
		return ErrStrategyNotSupported
	}
}

func (w *batcherWallet) Address() btcutil.Address {
	return w.address
}

func (w *batcherWallet) GenerateSACP(ctx context.Context, spendReq SpendRequest, to btcutil.Address) ([]byte, error) {
	return w.sw.GenerateSACP(ctx, spendReq, to)
}

func (w *batcherWallet) SignSACPTx(tx *wire.MsgTx, idx int, amount int64, leaf txscript.TapLeaf, scriptAddr btcutil.Address, witness [][]byte) ([][]byte, error) {
	return w.sw.SignSACPTx(tx, idx, amount, leaf, scriptAddr, witness)
}

// Send creates a batch request , saves it in the cache and returns a tracking id
func (w *batcherWallet) Send(ctx context.Context, sends []SendRequest, spends []SpendRequest, sacps [][]byte) (string, error) {
	if err := w.validateBatchRequest(ctx, spends, sends, sacps); err != nil {
		return "", err
	}

	id := chainhash.HashH([]byte(fmt.Sprintf("%v_%v", spends, sends))).String()
	req := BatcherRequest{
		ID:     id,
		Spends: spends,
		Sends:  sends,
		SACPs:  sacps,
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
	w.quit = make(chan struct{})

	w.logger.Info("starting batcher wallet")
	w.run(ctx)
	return nil
}

// Stop gracefully stops the batcher wallet service
func (w *batcherWallet) Stop() error {
	if w.quit == nil {
		return ErrBatcherNotRunning
	}

	w.logger.Info("stopping batcher wallet")
	close(w.quit)
	w.quit = nil

	w.logger.Info("waiting for batcher wallet to stop")
	w.wg.Wait()
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
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			ticker := time.NewTicker(w.opts.PTI)
			select {
			case <-w.quit:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := w.createBatch(); err != nil {
					if !errors.Is(err, ErrBatchParametersNotMet) {
						w.logger.Error("failed to create batch", zap.Error(err))
					} else {
						w.logger.Info("waiting for new batch")
					}

					if err := w.updateFeeRate(); err != nil {
						if !errors.Is(err, ErrFeeUpdateNotNeeded) {
							w.logger.Error("failed to update fee rate", zap.Error(err))
						} else {
							w.logger.Info("fee update skipped")
						}
					} else {
						w.logger.Info("batch fee updated", zap.String("strategy", string(w.opts.Strategy)))
					}
				} else {
					w.logger.Info("new batch created", zap.String("strategy", string(w.opts.Strategy)))
				}
			}

		}

	}()
}

// updateFeeRate updates the fee rate based on the strategy
func (w *batcherWallet) updateFeeRate() error {
	feeRates, err := w.feeEstimator.FeeSuggestion()
	if err != nil {
		return err
	}
	requiredFeeRate := selectFee(feeRates, w.opts.TxOptions.FeeLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch w.opts.Strategy {
	case CPFP:
		return w.updateCPFP(ctx, requiredFeeRate)
	case RBF:
		return w.updateRBF(ctx, requiredFeeRate)
	default:
		panic("fee update for strategy not implemented")
	}
}

// createBatch creates a batch based on the strategy
func (w *batcherWallet) createBatch() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch w.opts.Strategy {
	case CPFP:
		return w.createCPFPBatch(ctx)
	case RBF:
		return w.createRBFBatch(ctx)
	default:
		panic("batch creation for strategy not implemented")
	}
}

func (w *batcherWallet) validateBatchRequest(ctx context.Context, spends []SpendRequest, sends []SendRequest, sacps [][]byte) error {
	if len(spends) == 0 && len(sends) == 0 && len(sacps) == 0 {
		return ErrBatchParametersNotMet
	}

	utxos, err := w.indexer.GetUTXOs(ctx, w.address)
	if err != nil {
		return err
	}

	walletBalance := int64(0)
	for _, utxo := range utxos {
		walletBalance += utxo.Amount
	}

	sendsAmount := int64(0)
	for _, send := range sends {
		sendsAmount += send.Amount
	}

	spendsAmount := int64(0)
	err = withContextTimeout(ctx, 5*time.Second, func(ctx context.Context) error {
		_, _, spendsAmount, err = getUTXOsForSpendRequest(ctx, w.indexer, spends)
		return err
	})
	if err != nil {
		return err
	}

	if walletBalance+spendsAmount < sendsAmount {
		return ErrBatchParametersNotMet
	}

	return validateRequests(spends, sends, sacps)
}

// verifies if the fee rate delta is within the threshold
func validateUpdate(currentFeeRate, requiredFeeRate int, opts BatcherOptions) error {
	if currentFeeRate >= requiredFeeRate {
		return ErrFeeUpdateNotNeeded
	}
	if opts.TxOptions.MinFeeDelta > 0 && requiredFeeRate-currentFeeRate < opts.TxOptions.MinFeeDelta {
		return ErrFeeUpdateNotNeeded
	}
	if opts.TxOptions.MaxFeeDelta > 0 && requiredFeeRate-currentFeeRate > opts.TxOptions.MaxFeeDelta {
		return ErrFeeDeltaHigh
	}
	if opts.TxOptions.MaxFeeRate > 0 && requiredFeeRate > opts.TxOptions.MaxFeeRate {
		return ErrHighFeeEstimate
	}
	return nil
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

func filterPendingBatches(batches []Batch, indexer IndexerClient) ([]Batch, []string, []string, error) {
	pendingBatches := []Batch{}
	confirmedTxs := []string{}
	pendingTxs := []string{}
	for _, batch := range batches {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tx, err := indexer.GetTx(ctx, batch.Tx.TxID)
		if err != nil {
			return nil, nil, nil, err
		}
		if tx.Status.Confirmed {
			confirmedTxs = append(confirmedTxs, tx.TxID)
			continue
		}
		pendingBatches = append(pendingBatches, batch)
		pendingTxs = append(pendingTxs, tx.TxID)
	}
	return pendingBatches, confirmedTxs, pendingTxs, nil
}

func getTransaction(indexer IndexerClient, txid string) (Transaction, error) {
	if txid == "" {
		return Transaction{}, fmt.Errorf("txid is empty")
	}
	for i := 1; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tx, err := indexer.GetTx(ctx, txid)
		if err != nil {
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}
		return tx, nil
	}
	return Transaction{}, ErrTxNotFound
}

func withContextTimeout(parentContext context.Context, duration time.Duration, fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(parentContext, duration)
	defer cancel()
	return fn(ctx)
}
