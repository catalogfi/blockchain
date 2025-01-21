package btctest

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/fatih/color"
)

const (
	DefaultRegtestHost    = "0.0.0.0:18443"
	DefaultRegtestIndexer = "http://localhost:30000"
)

// NewBtcKey generates a new bitcoin private key.
func NewBtcKey(network *chaincfg.Params, addrType waddrmgr.AddressType) (*btcec.PrivateKey, btcutil.Address, error) {
	key, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	// Tweak the private key first if it's a taproot address
	if addrType == waddrmgr.TaprootPubKey {
		tapPubKey := txscript.ComputeTaprootKeyNoScript(key.PubKey())
		addr, err := btc.PublicKeyAddress(network, addrType, tapPubKey)
		if err != nil {
			return nil, nil, err
		}
		return key, addr, nil
	}
	addr, err := btc.PublicKeyAddress(network, addrType, key.PubKey())
	if err != nil {
		return nil, nil, err
	}
	return key, addr, nil
}

// NewBtcAddrWithFunds generates a new bitcoin private key and funds it using the `merry faucet` command.
func NewBtcAddrWithFunds(network *chaincfg.Params, addrType waddrmgr.AddressType, indexer btc.IndexerClient) (*btcec.PrivateKey, btcutil.Address, error) {
	key, addr, err := NewBtcKey(network, addrType)
	if err != nil {
		return nil, nil, err
	}
	if indexer != nil {
		_, err := FaucetWaitedMined(addr.EncodeAddress(), indexer)
		return key, addr, err
	}
	_, err = Faucet(addr.EncodeAddress())
	return key, addr, err
}

// RandomSecret creates a random secret with size [1,32)
func RandomSecret() []byte {
	length := rand.Intn(31) + 1
	data := make([]byte, length)

	_, err := crand.Read(data)
	if err != nil {
		panic(err)
	}
	return data
}

// Faucet funds the given address using the `merry faucet` command. It will transfer 1 BTC to the target address
// and automatically generate a new block for the tx. It returns the txid of the funding transaction.
func Faucet(addr string) (*chainhash.Hash, error) {
	res, err := RunOutput("merry", "faucet", "--to", addr)
	if err != nil {
		return nil, err
	}
	txid := strings.TrimSpace(strings.TrimPrefix(string(res), "Successfully submitted at http://localhost:5050/tx/"))
	color.Green(fmt.Sprintf("Funding address1 %v , txid = %v", addr, txid))

	return chainhash.NewHashFromStr(txid)
}

// FaucetWaitedMined does the same thing as Faucet, but it will wait until the tx been detected by the indexer.
func FaucetWaitedMined(addr string, indexer btc.IndexerClient) (*chainhash.Hash, error) {
	res, err := RunOutput("merry", "faucet", "--to", addr)
	if err != nil {
		return nil, err
	}
	txid := strings.TrimSpace(strings.TrimPrefix(string(res), "Successfully submitted at http://localhost:5050/tx/"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for time.Sleep(3 * time.Second); ; time.Sleep(time.Second) {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("tx not found: %w", err)
		default:
		}

		_, err := indexer.GetTx(ctx, txid)
		if err != nil {
			continue
		}
		break
	}
	color.Green(fmt.Sprintf("Funding address1 %v , txid = %v", addr, txid))
	return chainhash.NewHashFromStr(txid)
}

// NewBlock will mine a new block in the reg testnet. This is usually useful when we need to test something with
// confirmations. It uses the `merry faucet` command to generate a new block, the receiver address is a dummy address
// which shouldn't affect our testing
func NewBlock() error {
	addr := "mwt4FeMsGv6Ua3WrfuhypPtqDUse9CoJev"
	_, err := RunOutput("merry", "faucet", "--to", addr)
	color.Green("Mined a new block")
	return err
}

// NewBlockWaitMined does the same thing as NewBlock, but it will wait until the tx been detected by the indexer.
func NewBlockWaitMined(indexer btc.IndexerClient) error {
	addr := "mwt4FeMsGv6Ua3WrfuhypPtqDUse9CoJev"
	res, err := RunOutput("merry", "faucet", "--to", addr)
	txid := strings.TrimSpace(strings.TrimPrefix(string(res), "Successfully submitted at http://localhost:5050/tx/"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for time.Sleep(3 * time.Second); ; time.Sleep(time.Second) {
		select {
		case <-ctx.Done():
			return fmt.Errorf("tx not found: %w", err)
		default:
		}

		_, err := indexer.GetTx(ctx, txid)
		if err != nil {
			continue
		}
		break
	}
	color.Green("Mined a new block")
	return nil
}

// RunOutput the command and catch the output
func RunOutput(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Output()
}
