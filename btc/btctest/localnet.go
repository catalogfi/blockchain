package btctest

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/catalogfi/blockchain/btc"
	"github.com/fatih/color"
)

var dummyAddr = "mwt4FeMsGv6Ua3WrfuhypPtqDUse9CoJev"

// Faucet funds the given address using the `merry faucet` command. It will transfer 1 BTC to the target address and
// automatically generate a new block for the tx. It returns the txid of the funding transaction.
func Faucet(addr string) (*chainhash.Hash, error) {
	res, err := RunOutput("merry", "faucet", "--to", addr)
	if err != nil {
		return nil, err
	}
	txid := strings.TrimSpace(strings.TrimPrefix(string(res), "Successfully submitted at http://localhost:5050/tx/"))
	color.Green(fmt.Sprintf("Funding address %v , txid = %v", addr, txid))

	return chainhash.NewHashFromStr(txid)
}

// NewBlock will mine a new block in the reg testnet. This is useful when we need to test something with confirmations.
// It uses the `merry faucet` command to generate a new block, the receiver address is a dummy address
// which shouldn't affect our testing
func NewBlock() error {
	_, err := RunOutput("merry", "faucet", "--to", dummyAddr)
	color.Green("Mined a new block")
	return err
}

// NewBlockWaitMined does the same thing as NewBlock, but it will wait until the block been detected by the indexer.
func NewBlockWaitMined(indexer btc.IndexerClient) error {
	res, err := RunOutput("merry", "faucet", "--to", dummyAddr)
	if err != nil {
		return err
	}
	// todo : maybe it's better to use a regex to parse it. we might configure the node in a different port.
	txid := strings.TrimSpace(strings.TrimPrefix(string(res), "Successfully submitted at http://localhost:5050/tx/"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return WaitMined(ctx, indexer, WaitTx(txid))
}

// RunOutput the command and catch the output
func RunOutput(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

type WaitMinedFunc func(ctx context.Context, indexer btc.IndexerClient) error

func WaitMined(ctx context.Context, indexer btc.IndexerClient, f WaitMinedFunc) error {
	for time.Sleep(time.Second); ; time.Sleep(time.Second) {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout : %w", ctx.Err())
		default:
		}

		if err := f(ctx, indexer); err == nil {
			return nil
		}
	}
}

func WaitTx(txid string) WaitMinedFunc {
	return func(ctx context.Context, indexer btc.IndexerClient) error {
		_, err := indexer.GetTx(ctx, txid)
		return err
	}
}

func WaitFunds(addr btcutil.Address) WaitMinedFunc {
	return func(ctx context.Context, indexer btc.IndexerClient) error {
		utxos, err := indexer.GetUTXOs(ctx, addr)
		if err != nil {
			return err
		}
		if len(utxos) == 0 {
			return fmt.Errorf("utxo not found")
		}
		return nil
	}
}

func WaitBlock(height uint64) WaitMinedFunc {
	return func(ctx context.Context, indexer btc.IndexerClient) error {
		latest, err := indexer.GetTipBlockHeight(ctx)
		if err != nil {
			return err
		}
		if latest < height {
			return errors.New("block not mined")
		}
		return nil
	}
}
