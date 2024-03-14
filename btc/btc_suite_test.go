package btc_test

import (
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var (
	// Envs
	btcUsername string
	btcPassword string
	indexerHost string

	// Vars
	network *chaincfg.Params
	logger  *zap.Logger
	indexer btc.IndexerClient
	client  btc.Client
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
	indexerHost, ok = os.LookupEnv("BTC_REGNET_INDEXER")
	Expect(ok).Should(BeTrue())

	By("Initialise some variables used across tests")
	var err error
	network = &chaincfg.RegressionNetParams
	logger, err = zap.NewDevelopment()
	Expect(err).Should(BeNil())
	indexer = btc.NewElectrsIndexerClient(logger, indexerHost, btc.DefaultRetryInterval)
	config := &rpcclient.ConnConfig{
		Params:       chaincfg.RegressionNetParams.Name,
		Host:         btctest.DefaultRegtestHost,
		User:         btcUsername,
		Pass:         btcPassword,
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	client, err = btc.NewClient(config)
	Expect(err).Should(BeNil())
})
