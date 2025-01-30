package btctest

import (
	"context"
	crand "crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
)

const (
	DefaultRegtestHost    = "0.0.0.0:18443"
	DefaultRegtestIndexer = "http://localhost:30000"
)

type WaitMinedFunc func(ctx context.Context, indexer btc.IndexerClient) error

func WaitMined(ctx context.Context, indexer btc.IndexerClient, f WaitMinedFunc) error {
	for time.Sleep(time.Second); ; time.Sleep(time.Second) {
		select {
		case <-ctx.Done():
			return fmt.Errorf("tx not found: %w", ctx.Err())
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

// NewBtcAddrWithFunds generates a new bitcoin private key and funds it using the `merry faucet` command. The `indexer`
// parameter is optional and can be used to wait for the tx to be mined.
func NewBtcAddrWithFunds(network *chaincfg.Params, addrType waddrmgr.AddressType, indexer btc.IndexerClient) (*btcec.PrivateKey, btcutil.Address, error) {
	key, addr, err := NewBtcKey(network, addrType)
	if err != nil {
		return nil, nil, err
	}
	txid, err := Faucet(addr.EncodeAddress())
	if err != nil {
		return nil, nil, err
	}
	if indexer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = WaitMined(ctx, indexer, WaitTx(txid.String()))
		return key, addr, err
	}
	return key, addr, nil
}

// RandomSecret creates a random secret with size [1,32)
func RandomSecret() ([]byte, [32]byte) {
	length := rand.Intn(31) + 1
	data := make([]byte, length)

	_, err := crand.Read(data)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(data)
	return data, hash
}

// NewHtlc creates a new HTLC with random secret.
func NewHtlc(initiatorPubKey, redeemerPubKey *btcec.PublicKey, timelock, amount int64) (*btc.HTLC, []byte, error) {
	initiatorPubBytes := schnorr.SerializePubKey(initiatorPubKey)
	redeemerPubBytes := schnorr.SerializePubKey(redeemerPubKey)
	secret, secretHash := RandomSecret()
	htlc, err := btc.NewHTLC(initiatorPubBytes, redeemerPubBytes, secretHash[:], timelock, amount)
	return htlc, secret, err
}
