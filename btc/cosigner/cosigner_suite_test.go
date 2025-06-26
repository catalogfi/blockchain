package cosigner_test

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/catalogfi/blockchain/btc/cosigner"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	// Envs
	debug string

	// Vars
	network     *chaincfg.Params
	logger      *zap.Logger
	indexer     btc.IndexerClient
	cosignerPub *btcec.PublicKey
	btcClient   btc.Client
	client      *cosigner.Client
	addrTypes   []waddrmgr.AddressType
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

	btcClient = btc.NewClient(network, btctest.DefaultRegtestHost, btctest.RegressionRpcUsername, btctest.RegressionRpcPassword)
	addrTypes = []waddrmgr.AddressType{
		waddrmgr.PubKeyHash,
		waddrmgr.WitnessPubKey,
		waddrmgr.TaprootPubKey,
	}
	cosignerPubStr, err := hex.DecodeString("0321f053cec7917da6213b489994212ed0637dad9df995c107a763fb8e4ee081f1")
	Expect(err).Should(BeNil())
	cosignerPub, err = btcec.ParsePubKey(cosignerPubStr)
	Expect(err).Should(BeNil())
	client = cosigner.NewClient("http://127.0.0.1:8080")
})
