package wallet_test

import (
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/catalogfi/blockchain/localnet"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	// Envs
	btcUsername string
	btcPassword string
	debug       string

	// Vars
	network      *chaincfg.Params
	logger       *zap.Logger
	indexer      btc.IndexerClient
	client       btc.Client
	feeEstimator btc.FeeEstimator
)

func TestBtc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Btc Suite")
}

var _ = BeforeSuite(func() {
	By("Check if required ENVs are set.")
	By("You may want to disable some assertion when forcing running a specific test.")

	var ok bool
	btcUsername, ok = os.LookupEnv("BTC_REGNET_USERNAME")
	Expect(ok).Should(BeTrue())
	btcPassword, ok = os.LookupEnv("BTC_REGNET_PASSWORD")
	Expect(ok).Should(BeTrue())
	debug = os.Getenv("debug")

	By("Initialise some variables used across tests")
	var err error
	network = &chaincfg.RegressionNetParams
	loggerConfig := zap.NewDevelopmentConfig()
	if debug == "" {
		loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	logger, err = loggerConfig.Build()
	Expect(err).Should(BeNil())
	indexer = btc.NewElectrsIndexerClient(logger, btctest.DefaultRegtestIndexer, btc.DefaultRetryInterval)
	config := &rpcclient.ConnConfig{
		Params:       chaincfg.RegressionNetParams.Name,
		Host:         localnet.DefaultRegtestHost,
		User:         btcUsername,
		Pass:         btcPassword,
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	client, err = btc.NewClient(config)
	Expect(err).Should(BeNil())
	feeEstimator = btc.NewFixFeeEstimator(10)
})
