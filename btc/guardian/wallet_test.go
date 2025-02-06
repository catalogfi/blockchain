package guardian_test

import (
	"context"
	"fmt"
	"math/rand"
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

func setupSimpleWallets(t *testing.T, count int) ([]btc.Wallet, context.Context) {
	var wallets []btc.Wallet
	for i := 0; i < count; i++ {
		privKey, _ := randomPrivKey()

		feeRate := 1

		chainParams := chaincfg.RegressionNetParams
		feeEstimator := btc.NewFixFeeEstimator(feeRate)
		indexer := indexerClient()

		simpleWallet, err := btc.NewSimpleWallet(privKey, &chainParams, indexer, feeEstimator, btc.MediumFee)
		require.NoError(t, err)
		wallets = append(wallets, simpleWallet)
	}
	ctx := context.Background()
	return wallets, ctx
}

func setupTest(t *testing.T) (*testWallets, context.Context) {
	// Start mining blocks in background
	go func() {
		for {
			localnet.MineBTCBlock()
			rand.NewSource(time.Now().Unix())
			seconds := time.Duration(rand.Intn(20))
			time.Sleep(seconds * time.Second)
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

	bitcoinRPC := btc.CreateBitcoinRPCClient("admin1", "123", "http://0.0.0.0:18443")
	// Create wallets
	guardianWallet, err := guardian.NewWallet(guardianClient, privKey, cache, indexer, &chainParams, feeEstimator, logger, bitcoinRPC)
	require.NoError(t, err)

	simpleWallet, err := btc.NewSimpleWallet(privKey, &chainParams, indexer, feeEstimator, btc.MediumFee)
	require.NoError(t, err)

	// Fund the guardian wallet
	addr := guardianWallet.Address()
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

func setupSimpleFunded(t *testing.T) (*testWallets, context.Context) {
	// Start mining blocks in background
	go func() {
		for {
			localnet.MineBTCBlock()
			rand.NewSource(time.Now().Unix())
			seconds := time.Duration(rand.Intn(5) + 5)
			time.Sleep(seconds * time.Second)
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

	bitcoinRPC := btc.CreateBitcoinRPCClient("admin1", "123", "http://localhost:18443")

	// Create walletstxfrom
	guardianWallet, err := guardian.NewWallet(guardianClient, privKey, cache, indexer, &chainParams, feeEstimator, logger, bitcoinRPC)
	require.NoError(t, err)

	simpleWallet, err := btc.NewSimpleWallet(privKey, &chainParams, indexer, feeEstimator, btc.MediumFee)
	require.NoError(t, err)

	// Fund the guardian wallet
	addr := simpleWallet.Address()
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

		require.Equal(t, wallet.Address().String(), transaction.VOUTs[1].ScriptPubKeyAddress)

		randomP2PKHAddr2 := randomP2PKHAddr()

		tx, err = wallet.Send(ctx, []btc.SendRequest{
			btc.NewSendRequest("2", int64(amount), randomP2PKHAddr2),
		})
		require.NoError(t, err)
		require.NotEmpty(t, tx)
	})

	t.Run("should be able to merge merge txs", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			inits := rand.Intn(25) + 25
			redeems := (rand.Intn(inits) % 25) + 1
			mergeReqs := rand.Intn(redeems) + 1
			normalReqs := rand.Intn(10)

			simpleWallets, _ := setupSimpleWallets(t, inits)
			mergeTxHexes := make(map[string]string)
			utxos := make(map[string][]btc.Prevout)

			sendReqs := []btc.SendRequest{}
			for j := 0; j < inits; j++ {
				sendReqs = append(sendReqs, btc.NewSendRequest(fmt.Sprintf("3_%d_%d", i, j), amount, simpleWallets[j].Address()))
			}

			tx, err := wallet.Send(ctx, sendReqs)
			require.NoError(t, err)
			require.NotEmpty(t, tx)

			time.Sleep(time.Duration(5) * time.Second)

			for k := 0; k < redeems; k++ {
				randomP2PKHAddr2 := randomP2PKHAddr()
				fee := int64(1000)

				txString, err := simpleWallets[k].Send(ctx, []btc.SendRequest{
					// this amount is after minus fee
					btc.NewSendRequest(fmt.Sprintf("1_%d", k), amount-fee, randomP2PKHAddr2),
				}, nil, nil)

				require.NoError(t, err)
				require.NotEmpty(t, txString)
				txHex, err := setup.indexer.GetTxHex(ctx, txString)
				require.NoError(t, err)

				txx, err := setup.indexer.GetTx(ctx, txString)

				mergeTxHexes[fmt.Sprintf("1_%d", k)] = txHex
				utxos[fmt.Sprintf("1_%d", k)] = txx.VOUTs

				require.NoError(t, err)
				require.Equal(t, txx.Fee, fee)
			}

			mergeRe := []btc.SendRequest{}
			for l := 0; l < mergeReqs; l++ {
				for _, out := range utxos[fmt.Sprintf("1_%d", l)] {
					addr, err := btcutil.DecodeAddress(out.ScriptPubKeyAddress, setup.chainParams)
					require.NoError(t, err)
					mergeRe = append(mergeRe, btc.NewSendRequestWithInvalidateID(fmt.Sprintf("4_%d_%d", i, l), int64(out.Value), addr, fmt.Sprintf("3_%d_%d", i, l), mergeTxHexes[fmt.Sprintf("1_%d", l)]))
				}
			}
			if len(mergeRe) > 0 {
				tx, err := wallet.Send(ctx, mergeRe)
				if err != nil && err.Error() == "no remaining requests to process" {
					fmt.Println("no remaining requests to process")
					continue
				}
				require.NoError(t, err)
				require.NotEmpty(t, tx)

			}

			normalRe := []btc.SendRequest{}
			for p := 0; p < normalReqs; p++ {
				randomAddr := randomP2PKHAddr()
				normalRe = append(normalRe, btc.NewSendRequest(fmt.Sprintf("2_%d_%d", i, p), amount-1000, randomAddr))
			}

			if len(normalRe) > 0 {
				tx, err := wallet.Send(ctx, normalRe)
				fmt.Println("after normal req : ", tx)
				if err != nil && err.Error() == "no remaining requests to process" {
					fmt.Println("no remaining requests to process")
					continue
				}
				if err != nil {
					fmt.Println(err)
				}

				require.NoError(t, err)
				require.NotEmpty(t, tx)
			}

		}
	})

	t.Run("randomly generate addresses and send funds to them", func(t *testing.T) {
		for i := 0; i < 15; i++ {
			addr := randomP2PKHAddr()
			tx, err := wallet.Send(ctx, []btc.SendRequest{
				btc.NewSendRequest(fmt.Sprintf("5_%d", i), amount, addr),
			})
			require.NoError(t, err)
			require.NotEmpty(t, tx)
		}
	})

	t.Run("sending empty request", func(t *testing.T) {
		setup, ctx := setupTest(t)
		gWallet := setup.guardian

		newReqs := []btc.SendRequest{}
		tx, err := gWallet.Send(ctx, newReqs)
		require.Empty(t, tx)
		require.Error(t, fmt.Errorf("no request found"), err)
	})
}

func TestPossibleFailure(t *testing.T) {
	setup, ctx := setupTest(t)
	gWallet := setup.guardian
	simple := setup.simple

	t.Run("duplicate_request", func(t *testing.T) {
		tx, err := gWallet.Send(ctx, []btc.SendRequest{
			btc.NewSendRequest("1_0", 10000, randomP2PKHAddr()),
		})
		fmt.Println("first :", tx.String())
		require.NoError(t, err)

		tx, err = gWallet.Send(ctx, []btc.SendRequest{
			btc.NewSendRequest("1_0", 78942, randomP2PKHAddr()),
		})
		fmt.Println("second :", tx.String())

		require.NoError(t, err)
	})
	t.Run("unconfirmed descendanant", func(t *testing.T) {
		tx, err := gWallet.Send(ctx, []btc.SendRequest{
			btc.NewSendRequest("3_0", 10000, simple.Address()),
		})

		fmt.Println("first :", tx.String())
		require.NoError(t, err)

		randomaddr := randomP2PKHAddr()
		txString, err := simple.Send(ctx, []btc.SendRequest{
			btc.NewSendRequest("1_0", 10000-1000, randomaddr),
		}, nil, nil)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Sending from %s to %s\n", simple.Address().String(), randomaddr.String())
		fmt.Println("after simple to randomaddr", txString)
		// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		// defer cancel()
		txx, err := setup.indexer.GetTx(ctx, txString)
		if err != nil {
			fmt.Println(err)
		}

		txHex, err := setup.indexer.GetTxHex(ctx, txString)
		if err != nil {
			fmt.Println(err)
		}

		mergeRe := []btc.SendRequest{}
		for _, out := range txx.VOUTs {
			addr, err := btcutil.DecodeAddress(out.ScriptPubKeyAddress, setup.chainParams)
			require.NoError(t, err)
			mergeRe = append(mergeRe, btc.NewSendRequestWithInvalidateID("4_0", int64(out.Value), addr, "3_0", txHex))
		}

		tx1, err := gWallet.Send(ctx, mergeRe)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Sending from %s to %s\n", gWallet.Address().String(), randomaddr.String())
		fmt.Println("after merge :", tx1.String())
		require.NoError(t, err)
	})
	t.Run("invalidate_nonexistent_request", func(t *testing.T) {
		tx, err := gWallet.Send(ctx, []btc.SendRequest{
			btc.NewSendRequest("3_0", 10000, simple.Address()),
		})

		fmt.Println("first :", tx.String())
		require.NoError(t, err)

		randomaddr := randomP2PKHAddr()
		txString, err := simple.Send(ctx, []btc.SendRequest{
			btc.NewSendRequest("1_0", 10000-1000, randomaddr),
		}, nil, nil)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Sending from %s to %s\n", simple.Address().String(), randomaddr.String())
		fmt.Println("after simple to randomaddr", txString)
		// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		// defer cancel()
		txx, err := setup.indexer.GetTx(ctx, txString)
		if err != nil {
			fmt.Println(err)
		}

		txHex, err := setup.indexer.GetTxHex(ctx, txString)
		if err != nil {
			fmt.Println(err)
		}

		mergeRe := []btc.SendRequest{}
		for _, out := range txx.VOUTs {
			addr, err := btcutil.DecodeAddress(out.ScriptPubKeyAddress, setup.chainParams)
			require.NoError(t, err)
			mergeRe = append(mergeRe, btc.NewSendRequestWithInvalidateID("4_0", int64(out.Value), addr, "3_1", txHex))
		}

		tx1, err := gWallet.Send(ctx, mergeRe)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Sending from %s to %s\n", gWallet.Address().String(), randomaddr.String())
		fmt.Println("after merge :", tx1.String())
		require.NoError(t, err)
	})
	t.Run("invalidate_the_same_request_twice", func(t *testing.T) {
		tx, err := gWallet.Send(ctx, []btc.SendRequest{
			btc.NewSendRequest("3_0", 10000, simple.Address()),
		})

		fmt.Println("first :", tx.String())
		require.NoError(t, err)

		randomaddr := randomP2PKHAddr()
		txString, err := simple.Send(ctx, []btc.SendRequest{
			btc.NewSendRequest("1_0", 10000-1000, randomaddr),
		}, nil, nil)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Sending from %s to %s\n", simple.Address().String(), randomaddr.String())
		fmt.Println("after simple to randomaddr", txString)
		// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		// defer cancel()
		txx, err := setup.indexer.GetTx(ctx, txString)
		if err != nil {
			fmt.Println(err)
		}

		txHex, err := setup.indexer.GetTxHex(ctx, txString)
		if err != nil {
			fmt.Println(err)
		}

		mergeRe := []btc.SendRequest{}
		for _, out := range txx.VOUTs {
			addr, err := btcutil.DecodeAddress(out.ScriptPubKeyAddress, setup.chainParams)
			require.NoError(t, err)
			mergeRe = append(mergeRe, btc.NewSendRequestWithInvalidateID("4_0", int64(out.Value), addr, "3_0", txHex))
		}
		tx1, err := gWallet.Send(ctx, mergeRe)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("after 1st merge4_0: ", tx1.String())

		mergeRe = []btc.SendRequest{}
		for _, out := range txx.VOUTs {
			addr, err := btcutil.DecodeAddress(out.ScriptPubKeyAddress, setup.chainParams)
			require.NoError(t, err)
			mergeRe = append(mergeRe, btc.NewSendRequestWithInvalidateID("4_1", int64(out.Value), addr, "3_0", txHex))
		}

		tx1, err = gWallet.Send(ctx, mergeRe)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Sending from %s to %s\n", gWallet.Address().String(), randomaddr.String())
		fmt.Println("after merge :", tx1.String())
		require.Error(t, fmt.Errorf("invalid merge request found"), err)
	})
}

func TestUtxoUtilization(t *testing.T) {
	setup,ctx := setupSimpleFunded(t)
	gWallet := setup.guardian
	simple := setup.simple

	tx,err := simple.Send(ctx, []btc.SendRequest{
		btc.NewSendRequest("1_0", 10000, gWallet.Address()),
	}, nil,nil)

	require.NotEmpty(t, tx)
	require.NoError(t, err)

	tx,err = simple.Send(ctx, []btc.SendRequest{
		btc.NewSendRequest("2_0", 20000, gWallet.Address()),
	}, nil,nil)

	require.NotEmpty(t, tx)
	require.NoError(t, err)

	randomAddr := randomP2PKHAddr()
	tx1,err := gWallet.Send(ctx, []btc.SendRequest{
		btc.NewSendRequest("1_0", 5000, randomAddr),
	})
	
	require.NotEmpty(t, tx1)
	require.NoError(t, err)
	
	tx1,err = gWallet.Send(ctx, []btc.SendRequest{
		btc.NewSendRequest("2_0", 15000, randomAddr),
	})
	
	require.NotEmpty(t, tx1)
	require.NoError(t, err)
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
