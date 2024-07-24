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
	SegwitSpendWeight = txsizes.RedeemP2WPKHInputWitnessWeight
	DefaultAPITimeout = 5 * time.Second
)
var (
	ErrBatcherStillRunning        = errors.New("batcher is still running")
	ErrBatcherNotRunning          = errors.New("batcher is not running")
	ErrBatchParametersNotMet      = errors.New("batch parameters not met")
	ErrHighFeeEstimate            = errors.New("estimated fee too high")
	ErrFeeDeltaHigh               = errors.New("fee delta too high")
	ErrFeeUpdateNotNeeded         = errors.New("fee update not needed")
	ErrMaxBatchLimitReached       = errors.New("max batch limit reached")
	ErrCPFPFeeUpdateParamsNotMet  = errors.New("CPFP fee update parameters not met")
	ErrCPFPBatchingCorrupted      = errors.New("CPFP batching corrupted")
	ErrSavingBatch                = errors.New("failed to save batch")
	ErrStrategyNotSupported       = errors.New("strategy not supported")
	ErrBuildCPFPDepthExceeded     = errors.New("build CPFP depth exceeded")
	ErrBuildRBFDepthExceeded      = errors.New("build RBF depth exceeded")
	ErrTxIDEmpty                  = errors.New("txid is empty")
	ErrInsufficientFundsInRequest = func(have, need int64) error {
		return fmt.Errorf("%v , have :%v, need at least : %v", ErrBatchParametersNotMet, have, need)
	}
	ErrUnconfirmedUTXO = func(utxo UTXO) error {
		return fmt.Errorf("%v ,unconfirmed utxo :%v", ErrBatchParametersNotMet, utxo)
	}
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
	// ReadBatchByReqID reads a batch based on the request ID.
	ReadBatchByReqID(ctx context.Context, reqID string) (Batch, error)
	// ReadPendingBatches reads all pending batches for a given strategy.
	ReadPendingBatches(ctx context.Context) ([]Batch, error)
	// ReadLatestBatch reads the latest batch for a given strategy.
	ReadLatestBatch(ctx context.Context) (Batch, error)

	ReadBatch(ctx context.Context, id string) (Batch, error)

	// UpdateBatches updates multiple batches.
	// It completely overwrites the existing batches with the newer ones.
	// If no batch exists, it will create a new one.
	//
	// Note: All the requests which are a part of this batch will be removed from pending requests.
	UpdateBatches(ctx context.Context, updatedBatches ...Batch) error

	// UpdateAndDeletePendingBatches overwrites the existing batches with newer ones and also deletes pending batches.
	// If no batch exists, it will create a new one and delete pending batches.
	//
	// Even if updating batch is a pending batch, it will not be deleted.
	UpdateAndDeletePendingBatches(ctx context.Context, updatedBatches ...Batch) error

	// DeletePendingBatches deletes pending batches based on confirmed transaction IDs and strategy.
	DeletePendingBatches(ctx context.Context) error

	// SaveBatch saves a batch.
	SaveBatch(ctx context.Context, batch Batch) error

	// ReadRequest reads a request based on its ID.	// ReadRequests reads multiple requests based on their IDs.
	ReadRequests(ctx context.Context, id ...string) ([]BatcherRequest, error)
	// ReadPendingRequests reads all pending requests.
	ReadPendingRequests(ctx context.Context) ([]BatcherRequest, error)
	// SaveRequest saves a request.
	SaveRequest(ctx context.Context, req BatcherRequest) error
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
	// Periodic Time Interval for batching
	PTI       time.Duration
	TxOptions TxOptions
	Strategy  Strategy
}

// Strategy defines the batching strategy to be used by the BatcherWallet.
// It can be one of RBF, CPFP, RBF_CPFP, Multi_CPFP.
//
// 1. RBF - Replace By Fee
//
// 2. CPFP - Child Pays For Parent
//
// 3. Multi_CPFP - Multiple CPFP threads are maintained across multiple addresses
type Strategy string

var (
	RBF        Strategy = "RBF"
	CPFP       Strategy = "CPFP"
	RBF_CPFP   Strategy = "RBF_CPFP"
	Multi_CPFP Strategy = "Multi_CPFP"
)

type TxOptions struct {
	FeeLevel    FeeLevel
	MaxFeeRate  int
	MinFeeDelta int
	MaxFeeDelta int
}

type batcherWallet struct {
	quit chan struct{}
	wg   sync.WaitGroup

	chainParams *chaincfg.Params
	address     btcutil.Address
	privateKey  *secp256k1.PrivateKey
	logger      *zap.Logger

	sw           Wallet
	opts         BatcherOptions
	indexer      IndexerClient
	feeEstimator FeeEstimator
	cache        Cache
}

type Batch struct {
	Tx         Transaction
	RequestIds map[string]bool
	// true indicates that the batch is finalized and will not be replaced by more fee.
	IsFinalized bool
	Strategy    Strategy
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
		chainParams:  chainParams,
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
	if err := w.validateBatchRequest(ctx, w.opts.Strategy, spends, sends, sacps); err != nil {
		return "", err
	}

	// generate random id
	id := chainhash.HashH([]byte(fmt.Sprintf("%v", time.Now().UnixNano()))).String()

	req := BatcherRequest{
		ID:     id,
		Spends: spends,
		Sends:  sends,
		SACPs:  sacps,
		Status: false,
	}
	return id, w.cache.SaveRequest(ctx, req)
}

// Status returns the status of a transaction based on the tracking id
func (w *batcherWallet) Status(ctx context.Context, id string) (Transaction, bool, error) {
	request, err := w.cache.ReadRequests(ctx, id)
	if err != nil {
		return Transaction{}, false, err
	}
	if !request[0].Status {
		return Transaction{}, false, nil
	}
	batch, err := w.cache.ReadBatchByReqID(ctx, id)
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

	w.logger.Info("--------starting batcher wallet--------")
	return w.run(ctx)
}

// Stop gracefully stops the batcher wallet service
func (w *batcherWallet) Stop() error {
	if w.quit == nil {
		return ErrBatcherNotRunning
	}

	w.logger.Info("--------stopping batcher wallet--------")
	close(w.quit)

	w.logger.Info("waiting for batcher wallet to stop")
	w.wg.Wait()

	w.logger.Info("batcher stopped")
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
func (w *batcherWallet) run(ctx context.Context) error {
	switch w.opts.Strategy {
	case CPFP, RBF:
		w.runPeriodicBatcher(ctx)
	default:
		return ErrStrategyNotSupported
	}
	return nil
}

// runPeriodicBatcher is used by strategies that require
// triggering the batching process at regular intervals
func (w *batcherWallet) runPeriodicBatcher(ctx context.Context) {
	ticker := time.NewTicker(w.opts.PTI)
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.quit:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.processBatch()
			}
		}
	}()
}

// processBatch contains the core logic of the batcher
//  1. It creates a batch at regular intervals
//  2. It also updates the fee rate at regular intervals
//     if fee rate increases more than threshold and there are
//     no batches to create
func (w *batcherWallet) processBatch() {
	if err := w.createBatch(); err != nil {
		if !errors.Is(err, ErrBatchParametersNotMet) {
			w.logger.Error("failed to create batch", zap.Error(err))
		} else {
			w.logger.Info("waiting for new batch")
		}

		if err := w.updateBatchFeeRate(); err != nil {
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

// updateBatchFeeRate updates the fee rate based on the strategy
func (w *batcherWallet) updateBatchFeeRate() error {
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
		return ErrStrategyNotSupported
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
		return ErrStrategyNotSupported
	}
}

func (w *batcherWallet) validateBatchRequest(ctx context.Context, strategy Strategy, spends []SpendRequest, sends []SendRequest, sacps [][]byte) error {
	if len(spends) == 0 && len(sends) == 0 && len(sacps) == 0 {
		return ErrBatchParametersNotMet
	}

	err := validateRequests(spends, sends, sacps)
	if err != nil {
		return err
	}

	for _, spend := range spends {
		if spend.ScriptAddress == w.address {
			return ErrBatchParametersNotMet
		}
	}

	utxos, err := w.indexer.GetUTXOs(ctx, w.address)
	if err != nil {
		return err
	}

	walletBalance := int64(0)
	for _, utxo := range utxos {
		walletBalance += utxo.Amount
	}

	spendsAmount := int64(0)
	spendsUtxos := UTXOs{}
	err = withContextTimeout(ctx, DefaultAPITimeout, func(ctx context.Context) error {
		spendsUtxos, _, spendsAmount, err = getUTXOsForSpendRequest(ctx, w.indexer, spends)
		return err
	})
	if err != nil {
		return err
	}

	var sacpsIn int64
	var sacpOut int64
	err = withContextTimeout(ctx, DefaultAPITimeout, func(ctx context.Context) error {
		sacpsIn, sacpOut, err = getSACPAmounts(ctx, sacps, w.indexer)
		return err
	})

	totalSendAmt := calculateTotalSendAmount(sends)

	in := walletBalance + spendsAmount + int64(sacpsIn)
	out := totalSendAmt + int64(sacpOut)
	if in < out+1000 {
		return ErrInsufficientFundsInRequest(in, out)
	}

	switch strategy {
	case RBF:
		for _, utxo := range spendsUtxos {
			if !utxo.Status.Confirmed {
				return ErrUnconfirmedUTXO(utxo)
			}
		}
	}

	return nil
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

func filterPendingBatches(batches []Batch, indexer IndexerClient) (pendingBatches []Batch, confirmedBatches []Batch, err error) {
	for _, batch := range batches {
		ctx, cancel := context.WithTimeout(context.Background(), DefaultAPITimeout)
		defer cancel()
		tx, err := indexer.GetTx(ctx, batch.Tx.TxID)
		if err != nil {
			return nil, nil, err
		}
		if tx.Status.Confirmed {
			confirmedBatches = append(confirmedBatches, batch)
			continue
		}
		pendingBatches = append(pendingBatches, batch)
	}
	return pendingBatches, confirmedBatches, nil
}

func withContextTimeout(parentContext context.Context, duration time.Duration, fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(parentContext, duration)
	defer cancel()
	return fn(ctx)
}

// unpackBatcherRequests unpacks the batcher requests into spend requests, send requests and SACPs
func unpackBatcherRequests(reqs []BatcherRequest) ([]SpendRequest, []SendRequest, [][]byte, map[string]bool) {
	spendRequests := []SpendRequest{}
	sendRequests := []SendRequest{}
	sacps := [][]byte{}
	reqIds := make(map[string]bool)

	for _, req := range reqs {
		spendRequests = append(spendRequests, req.Spends...)
		sendRequests = append(sendRequests, req.Sends...)
		sacps = append(sacps, req.SACPs...)
		reqIds[req.ID] = true
	}

	return spendRequests, sendRequests, sacps, reqIds
}
