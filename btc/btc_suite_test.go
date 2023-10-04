package btc_test

import (
	"os"
	"testing"

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
	Expect(os.Getenv("BTC_USER")).ShouldNot(BeEmpty())
	Expect(os.Getenv("BTC_PASSWORD")).ShouldNot(BeEmpty())
	Expect(os.Getenv("BTC_INDEXER_ELECTRS_REGNET")).ShouldNot(BeEmpty())
})
