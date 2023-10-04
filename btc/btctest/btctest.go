package btctest

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/testutil"
	"go.uber.org/zap"
)

// RegtestClient initialises a btc.Client could be used with a local regression testnet.
// This uses some default setting of the client options and assumes all ENVs exist and not null.
func RegtestClient() (btc.Client, error) {
	user := testutil.ParseStringEnv("BTC_USER", "")
	password := testutil.ParseStringEnv("BTC_PASSWORD", "")
	config := &rpcclient.ConnConfig{
		Params:       chaincfg.RegressionNetParams.Name,
		Host:         "0.0.0.0:18443",
		User:         user,
		Pass:         password,
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	return btc.NewClient(config)
}

func RegtestIndexer() btc.IndexerClient {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	url := testutil.ParseStringEnv("BTC_INDEXER_ELECTRS_REGNET", "")
	return btc.NewElectrsIndexerClient(logger, url, btc.DefaultRetryInterval)
}
