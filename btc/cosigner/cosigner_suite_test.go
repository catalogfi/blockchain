package cosigner_test

import (
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	// Envs
	debug string

	// Vars
	network   *chaincfg.Params
	logger    *zap.Logger
	indexer   btc.IndexerClient
	client    btc.Client
	addrTypes []waddrmgr.AddressType
)

func TestBtc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Btc Suite")
}

var _ = BeforeSuite(func() {
	By("Check if required ENVs are set.")
	By("You may want to disable some assertion when forcing running a specific test.")
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

	client = btc.NewClient(network, btctest.DefaultRegtestHost, btctest.RegressionRpcUsername, btctest.RegressionRpcPassword)
	addrTypes = []waddrmgr.AddressType{
		waddrmgr.PubKeyHash,
		waddrmgr.WitnessPubKey,
		waddrmgr.TaprootPubKey,
	}
})
