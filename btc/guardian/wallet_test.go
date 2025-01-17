package guardian_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/guardian"
	"github.com/catalogfi/blockchain/localnet"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

type testWallets struct {
	guardian    *guardian.Wallet
	simple      btc.Wallet
	indexer     btc.IndexerClient
	chainParams *chaincfg.Params
	feeRate     int
}

func setupTest(t *testing.T) (*testWallets, context.Context) {
	// Start mining blocks in background
	go func() {
		for {
			localnet.MineBTCBlock()
			time.Sleep(10 * time.Second)
		}
	}()

	ctx := context.Background()
	privKey, pubKey := randomPrivKey()

	// Setup dependencies
	guardianClient := guardian.NewGuardianClient(ctx, "http://0.0.0.0:11818", pubKey)
	levelDB, err := leveldb.OpenFile(t.TempDir(), nil)
	require.NoError(t, err)

	feeRate := 1

	chainParams := chaincfg.RegressionNetParams
	feeEstimator := btc.NewFixFeeEstimator(feeRate)
	cache := guardian.NewCache(levelDB)
	indexer := indexerClient()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	// Create wallets
	guardianWallet, err := guardian.NewWallet(guardianClient, privKey, cache, indexer, &chainParams, feeEstimator, logger)
	require.NoError(t, err)

	simpleWallet, err := btc.NewSimpleWallet(privKey, &chainParams, indexer, feeEstimator, btc.MediumFee)
	require.NoError(t, err)

	// Fund the guardian wallet
	addr, err := guardianWallet.Address()
	require.NoError(t, err)
	_, err = localnet.FundBitcoin(addr.String(), indexer)
	require.NoError(t, err)

	return &testWallets{
		guardian:    guardianWallet,
		simple:      simpleWallet,
		indexer:     indexer,
		chainParams: &chainParams,
		feeRate:     feeRate,
	}, ctx
}

type AddressesToCheckFundsFor struct {
	address btcutil.Address
	amount  int64
}

func TestGuardianWallet(t *testing.T) {

	setup, ctx := setupTest(t)
	wallet := setup.guardian
	simple := setup.simple
	amount := int64(100000)

	t.Run("should send funds to a random addresses", func(t *testing.T) {

		randomP2PKHAddr1 := randomP2PKHAddr()

		tx, err := wallet.Send(ctx, []btc.SendRequest{
			btc.NewSendRequest("1", amount, randomP2PKHAddr1),
		})
		require.NoError(t, err)
		require.NotEmpty(t, tx)

		// let's make sure the tx fee is 1 v/b
		transaction, err := setup.indexer.GetTx(ctx, tx.String())
		require.NoError(t, err)
		require.Equal(t, int(transaction.Fee)/(transaction.Weight/4), setup.feeRate)

		// let's make sure tx outputs are correct
		require.Equal(t, len(transaction.VOUTs), 2)
		require.Equal(t, transaction.VOUTs[0].Value, int(amount))
		require.Equal(t, transaction.VOUTs[0].ScriptPubKeyAddress, randomP2PKHAddr1.String())

		addr, err := wallet.Address()
		require.NoError(t, err)
		require.Equal(t, addr.String(), transaction.VOUTs[1].ScriptPubKeyAddress)

		randomP2PKHAddr2 := randomP2PKHAddr()

		tx, err = wallet.Send(ctx, []btc.SendRequest{
			btc.NewSendRequest("2", int64(amount), randomP2PKHAddr2),
		})
		require.NoError(t, err)
		require.NotEmpty(t, tx)
	})

	t.Run("should be able to merge merge txs", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			tx, err := wallet.Send(ctx, []btc.SendRequest{
				btc.NewSendRequest(fmt.Sprintf("3_%d", i), amount, simple.Address()),
			})
			require.NoError(t, err)
			require.NotEmpty(t, tx)
			time.Sleep(5 * time.Second)

			randomP2PKHAddr2 := randomP2PKHAddr()
			fee := int64(1000)
			txString, err := simple.Send(ctx, []btc.SendRequest{
				// this amount is after minus fee
				btc.NewSendRequest(fmt.Sprintf("1_%d", i), amount-fee, randomP2PKHAddr2),
			}, nil, nil)
			require.NoError(t, err)
			require.NotEmpty(t, txString)

			txHex, err := setup.indexer.GetTxHex(ctx, txString)
			require.NoError(t, err)

			txx, err := setup.indexer.GetTx(ctx, txString)
			require.NoError(t, err)
			require.Equal(t, txx.Fee, fee)

			newReqs := []btc.SendRequest{}

			for _, out := range txx.VOUTs {
				addr, err := btcutil.DecodeAddress(out.ScriptPubKeyAddress, setup.chainParams)
				require.NoError(t, err)
				newReqs = append(newReqs, btc.NewSendRequestWithInvalidateID(fmt.Sprintf("4_%d", i), int64(out.Value), addr, fmt.Sprintf("3_%d", i), txHex))
			}

			tx, err = wallet.Send(ctx, newReqs)
			if err != nil && err.Error() == "no remaining requests to process" {
				fmt.Println("no remaining requests to process")
				continue
			}
			require.NoError(t, err)
			require.NotEmpty(t, tx)

		}
	})

	t.Run("randomly generate addresses and send funds to them", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			addr := randomP2PKHAddr()
			tx, err := wallet.Send(ctx, []btc.SendRequest{
				btc.NewSendRequest(fmt.Sprintf("5_%d", i), amount, addr),
			})
			require.NoError(t, err)
			require.NotEmpty(t, tx)
		}
	})

}

func randomPrivKey() (*btcec.PrivateKey, *btcec.PublicKey) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		panic(err)
	}
	return privKey, privKey.PubKey()
}

func randomP2PKHAddr() btcutil.Address {
	_, pubKey := randomPrivKey()
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.RegressionNetParams)
	if err != nil {
		panic(err)
	}
	return addr
}

func indexerClient() btc.IndexerClient {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return btc.NewElectrsIndexerClient(logger, "http://localhost:30000", 2*time.Second)
}
