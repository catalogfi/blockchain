package btc_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/localnet"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var (
	// Envs
	btcUsername string
	btcPassword string
	indexerHost string
	mode        MODE

	// Vars
	network *chaincfg.Params
	logger  *zap.Logger
	indexer btc.IndexerClient
	client  btc.Client

	modeFlag = flag.String("mode", string(SIMPLE), "Mode to run the tests: simple, batcher_rbf, batcher_cpfp")
)

type MODE string

const (
	// MODES
	SIMPLE       MODE = "simple"
	BATCHER_RBF  MODE = "batcher_rbf"
	BATCHER_CPFP MODE = "batcher_cpfp"
)

func parseMode(mode string) (MODE, error) {
	switch mode {
	case string(SIMPLE):
		return SIMPLE, nil
	case string(BATCHER_RBF):
		return BATCHER_RBF, nil
	case string(BATCHER_CPFP):
		return BATCHER_CPFP, nil
	default:
		return "", fmt.Errorf("unknown mode %s", mode)
	}
}

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

	By("Select the mode to run the tests")

	var err error
	mode, err = parseMode(*modeFlag)
	Expect(err).Should(BeNil())

	By("Initialise some variables used across tests")
	network = &chaincfg.RegressionNetParams
	logger, err = zap.NewDevelopment()
	Expect(err).Should(BeNil())
	indexer = btc.NewElectrsIndexerClient(logger, indexerHost, btc.DefaultRetryInterval)
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
})
