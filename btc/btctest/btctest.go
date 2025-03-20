package btctest

import (
	"context"
	crand "crypto/rand"
	"crypto/sha256"
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
	DefaultRegtestHost    = "http://0.0.0.0:18443"
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

// NewBtcAddrWithFunds generates a new bitcoin private key and funds it using the `merry faucet` command. The `indexer`
// parameter is optional and can be used to wait for the tx to be mined.
func NewBtcAddrWithFunds(network *chaincfg.Params, addrType waddrmgr.AddressType, indexer btc.IndexerClient) (*btcec.PrivateKey, btcutil.Address, error) {
	key, addr, err := NewBtcKey(network, addrType)
	if err != nil {
		return nil, nil, err
	}
	_, err = Faucet(addr.EncodeAddress())
	if err != nil {
		return nil, nil, err
	}
	if indexer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = WaitMined(ctx, indexer, WaitFunds(addr))
		return key, addr, err
	}
	return key, addr, nil
}

// NewWallet returns a new wallet. If `funds` is true, the wallet will be funded with the `merry faucet` command.
func NewWallet(network *chaincfg.Params, addrType waddrmgr.AddressType, indexer btc.IndexerClient, client btc.Client, feeEstimator btc.FeeEstimator, waitMined bool) (btc.Wallet, error) {
	waitMinedIndexer := indexer
	if !waitMined {
		waitMinedIndexer = nil
	}
	key, _, err := NewBtcAddrWithFunds(network, addrType, waitMinedIndexer)
	if err != nil {
		return nil, err
	}
	return btc.NewWallet(network, addrType, key, indexer, client, feeEstimator)
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

func RandomAmount(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomFeeRate will generates a random fee rate between 1 sat/vb to 1000 sat/vb.
func RandomFeeRate() btc.SatoshiPerKb {
	return btc.SatoshiPerKb(1000) + btc.SatoshiPerKb(rand.Int63n(99*1000+1))
}

// NewHtlc creates a new HTLC with random secret.
func NewHtlc(initiatorPubKey, redeemerPubKey *btcec.PublicKey, timelock, amount int64) (*btc.HTLC, error) {
	initiatorPubBytes := schnorr.SerializePubKey(initiatorPubKey)
	redeemerPubBytes := schnorr.SerializePubKey(redeemerPubKey)
	secret, secretHash := RandomSecret()
	htlc, err := btc.NewHTLC(initiatorPubBytes, redeemerPubBytes, secretHash[:], timelock, amount)
	if err != nil {
		return nil, err
	}
	htlc.SetSecret(secret)
	return htlc, err
}

// PrepareActions generates `n` number of htlcs and make them ready for the given action type. `wal1` and `wal2` will
// be the initiator's wallet and redeemer's wallet.
func PrepareActions(ctx context.Context, n int, wal1, wal2 btc.Wallet, indexer btc.IndexerClient, actionType btc.HtlcActionType) ([]btc.HtlcAction, error) {
	actions := make([]btc.HtlcAction, n)
	initiatorActions := make([]btc.HtlcAction, 0, n)
	timelock, amount := int64(6), RandomAmount(1e5, 1e7)
	miningBlocks := false
	for i := 0; i < n; i++ {
		htlc, err := NewHtlc(wal1.PublicKey(), wal2.PublicKey(), timelock, amount)
		if err != nil {
			return nil, err
		}

		switch actionType {
		case btc.HtlcActionInitiate:
			actions[i] = btc.HtlcAction{
				Htlc:       htlc,
				ActionType: btc.HtlcActionInitiate,
			}
		case btc.HtlcActionRedeem, btc.HtlcActionRefund:
			initiatorActions = append(initiatorActions, btc.HtlcAction{
				Htlc:       htlc,
				ActionType: btc.HtlcActionInitiate,
			})

			actions[i] = btc.HtlcAction{
				Htlc:       htlc,
				ActionType: actionType,
			}
			if actionType == btc.HtlcActionRefund {
				miningBlocks = true
			}
		case btc.HtlcActionInstantRefund:
			_, instantRefundTx, err := wal1.Initiate(ctx, htlc)
			if err != nil {
				return nil, err
			}
			actions[i] = btc.HtlcAction{
				Htlc:            htlc,
				ActionType:      btc.HtlcActionInstantRefund,
				InstantRefundTx: instantRefundTx,
			}
		}
	}

	if len(initiatorActions) != 0 {
		_, err := wal1.Execute(ctx, initiatorActions, "")
		if err != nil {
			return nil, err
		}

		// Mining certain blocks to make the
		if miningBlocks {
			for i := 0; i < int(timelock)-1; i++ {
				if err := NewBlock(); err != nil {
					return nil, err
				}
			}
		}
	}

	if actionType != btc.HtlcActionInitiate {
		// Mine a new block and wait for it to be confirmed
		if err := NewBlockWaitMined(indexer); err != nil {
			return nil, err
		}
	}

	return actions, nil
}
