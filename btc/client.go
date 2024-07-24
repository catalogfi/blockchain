package btc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

var (
	ErrTxNotFound = errors.New("no such mempool or blockchain transaction")

	ErrAlreadyInChain = errors.New("transaction already in block chain")

	ErrTxInputsMissingOrSpent = errors.New("bad-txns-inputs-missingorspent")

	ErrMempoolConflict = errors.New("txn-mempool-conflict")
)

// Client to interact with the Bitcoin network. It's implementation uses standard bitcoind JSON-RPC behind the scene.
type Client interface {

	// Net returns the network params.
	Net() *chaincfg.Params

	// LatestBlock returns the height and hash of the latest block.
	LatestBlock(ctx context.Context) (int64, string, error)

	// SubmitTx to the Bitcoin network.
	SubmitTx(ctx context.Context, tx *wire.MsgTx) error

	// GetRawTransaction returns the raw transaction of the given hash.
	GetRawTransaction(ctx context.Context, hash *chainhash.Hash) (*btcjson.TxRawResult, error)

	// GetBlockHash returns the block hash of the given height.
	GetBlockHash(ctx context.Context, height int64) (*chainhash.Hash, error)

	// GetBlock returns the block detail with the given hash.
	GetBlock(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error)

	// GetBlockVerbose returns the block detail with the given hash. It has more detailed info about txs than `GetBlock`
	GetBlockVerbose(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error)

	// GetTxOut returns details about an unspent transaction output. It will return nil result if the utxo has been
	// spent.
	GetTxOut(ctx context.Context, hash *chainhash.Hash, vout uint32) (*btcjson.GetTxOutResult, error)

	// GetNetworkInfo returns the network configuration of the node we connect to.
	GetNetworkInfo(ctx context.Context) (*btcjson.GetNetworkInfoResult, error)
}

type client struct {
	config    *rpcclient.ConnConfig
	params    *chaincfg.Params
	rpcClient *rpcclient.Client
}

func NewClient(config *rpcclient.ConnConfig) (Client, error) {
	c, err := rpcclient.New(config, nil)
	if err != nil {
		return nil, err
	}

	var param *chaincfg.Params
	switch config.Params {
	case chaincfg.MainNetParams.Name:
		param = &chaincfg.MainNetParams
	case chaincfg.TestNet3Params.Name:
		param = &chaincfg.TestNet3Params
	case chaincfg.RegressionNetParams.Name:
		param = &chaincfg.RegressionNetParams
	default:
		return nil, fmt.Errorf("rpcclient.New: Unknown chain %s", config.Params)
	}

	return &client{
		config:    config,
		params:    param,
		rpcClient: c,
	}, nil
}

func (client *client) Net() *chaincfg.Params {
	return client.params
}

func (client *client) LatestBlock(ctx context.Context) (int64, string, error) {
	future := client.rpcClient.GetBlockChainInfoAsync()
	results := make(chan *btcjson.GetBlockChainInfoResult, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return 0, "", fmt.Errorf("LatestBlock : %w", ctx.Err())
	case err := <-errs:
		return 0, "", err
	case result := <-results:
		return int64(result.Blocks), result.BestBlockHash, nil
	}
}

func (client *client) SubmitTx(ctx context.Context, tx *wire.MsgTx) error {
	// The SendRawTransactionAsync is not technically asynchronous,
	// we need to have an extra channel for it.
	futureChan := make(chan rpcclient.FutureSendRawTransactionResult, 1)
	go func() {
		defer close(futureChan)
		future := client.rpcClient.SendRawTransactionAsync(tx, false)
		futureChan <- future
	}()

	var future rpcclient.FutureSendRawTransactionResult
	select {
	case <-ctx.Done():
		return fmt.Errorf("SubmitTx : %w", ctx.Err())
	case future = <-futureChan:
	}

	results := make(chan *chainhash.Hash, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("SubmitTx : %w", ctx.Err())
	case err := <-errs:
		// Parse the error based on the error code and message
		var rpcErr *btcjson.RPCError
		if errors.As(err, &rpcErr) {
			switch rpcErr.Code {
			case btcjson.ErrRPCVerifyAlreadyInChain:
				return ErrAlreadyInChain
			case btcjson.ErrRPCTxRejected:
				if strings.Contains(err.Error(), "txn-mempool-conflict") {
					return ErrMempoolConflict
				}
			case btcjson.ErrRPCTxError:
				if strings.Contains(err.Error(), "bad-txns-inputs-missingorspent") {
					return ErrTxInputsMissingOrSpent
				}
			}
		}
		return err
	case <-results:
		return nil
	}
}

func (client *client) GetRawTransaction(ctx context.Context, txhash *chainhash.Hash) (*btcjson.TxRawResult, error) {
	future := client.rpcClient.GetRawTransactionVerboseAsync(txhash)
	results := make(chan *btcjson.TxRawResult, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("GetRawTransaction : %w", ctx.Err())
	case err := <-errs:
		// Parse the error based on the error code and message
		var rpcErr *btcjson.RPCError
		if errors.As(err, &rpcErr) {
			switch rpcErr.Code {
			case btcjson.ErrRPCInvalidAddressOrKey:
				if strings.Contains(err.Error(), "No such mempool or blockchain transaction") {
					return nil, ErrTxNotFound
				}
			}
		}

		return nil, err
	case result := <-results:
		return result, nil
	}
}

func (client *client) GetBlockHash(ctx context.Context, height int64) (*chainhash.Hash, error) {
	future := client.rpcClient.GetBlockHashAsync(height)
	results := make(chan *chainhash.Hash, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("GetBlockByHeight : %w", ctx.Err())
	case err := <-errs:
		return nil, err
	case hash := <-results:
		return hash, nil
	}
}

func (client *client) GetBlock(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	future := client.rpcClient.GetBlockVerboseAsync(hash)
	results := make(chan *btcjson.GetBlockVerboseResult, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("GetBlockByHash : %w", ctx.Err())
	case err := <-errs:
		return nil, err
	case result := <-results:
		return result, nil
	}
}

func (client *client) GetBlockVerbose(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error) {
	future := client.rpcClient.GetBlockVerboseTxAsync(hash)
	results := make(chan *btcjson.GetBlockVerboseTxResult, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("GetBlockByHash : %w", ctx.Err())
	case err := <-errs:
		return nil, err
	case result := <-results:
		return result, nil
	}
}

func (client *client) GetTxOut(ctx context.Context, hash *chainhash.Hash, vout uint32) (*btcjson.GetTxOutResult, error) {
	future := client.rpcClient.GetTxOutAsync(hash, vout, true)
	results := make(chan *btcjson.GetTxOutResult, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("GetTxOut : %w", ctx.Err())
	case err := <-errs:
		return nil, err
	case result := <-results:
		return result, nil
	}
}

func (client *client) GetNetworkInfo(ctx context.Context) (*btcjson.GetNetworkInfoResult, error) {
	future := client.rpcClient.GetNetworkInfoAsync()
	results := make(chan *btcjson.GetNetworkInfoResult, 1)
	errs := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errs)

		result, err := future.Receive()
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("GetNetworkInfo : %w", ctx.Err())
	case err := <-errs:
		return nil, err
	case result := <-results:
		return result, nil
	}
}
