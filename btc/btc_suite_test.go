package btc_test

import (
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/catalogfi/multichain/btc"
	"github.com/catalogfi/multichain/testutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

func TestBtc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Btc Suite")
}

var _ = BeforeSuite(func() {
	// Check the ENVS are set for the tests.
	By("These are the requirements for all tests in this suite. ")
	By("You may want to disable some assertion when forcing running a specific test")
	Expect(os.Getenv("PRIV_KEY_1")).ShouldNot(BeEmpty())
	Expect(os.Getenv("PRIV_KEY_2")).ShouldNot(BeEmpty())
	Expect(os.Getenv("BTC_RPC")).ShouldNot(BeEmpty())
	Expect(os.Getenv("BTC_USER")).ShouldNot(BeEmpty())
	Expect(os.Getenv("BTC_PASSWORD")).ShouldNot(BeEmpty())
	Expect(os.Getenv("BTC_INDEXER_ELECTRS_REGNET")).ShouldNot(BeEmpty())
})

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
