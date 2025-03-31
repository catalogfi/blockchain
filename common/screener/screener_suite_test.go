package screener_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestScreener(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Screener Suite")
}

var (
	// Envs
	url string
)

var _ = BeforeSuite(func() {
	By("Check if required ENVs are set.")
	By("You may want to disable some assertion when forcing running a specific test.")
	url = os.Getenv("")
	Expect(url).ShouldNot(BeEmpty())
})
