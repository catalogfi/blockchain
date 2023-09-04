package btc_test

import (
	"os"
	"testing"

	"github.com/catalogfi/multichain/btc"
	"github.com/catalogfi/multichain/testutil"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
	// Expect(os.Getenv("BTC_INDEXER_QUICKNODE")).ShouldNot(BeEmpty())
})

// RegtestClient initialises a btc.Client could be used with a local regression testnet.
// This uses some default setting of the client options and assumes all ENVs exist and not null.
func RegtestClient(logger *zap.Logger) btc.Client {
	user := testutil.ParseStringEnv("BTC_USER", "")
	password := testutil.ParseStringEnv("BTC_PASSWORD", "")
	opts := btc.DefaultClientOptions().WithUser(user).WithPassword(password)
	indexClient := btc.NewElectrsIndexerClient(logger, btc.DefaultElectrsIndexerURL)
	client := btc.NewClient(opts, logger, indexClient)
	return client
}
