package btc_test

import (
	"math/rand"
	"reflect"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/multichain/btc"
	"github.com/catalogfi/multichain/testutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bitcoin client options ", func() {
	Context("when using the default option", func() {
		It("should have all values set to default", func() {
			opts := btc.DefaultClientOptions()

			Expect(reflect.DeepEqual(opts.Net, &chaincfg.RegressionNetParams)).Should(BeTrue())
			Expect(opts.Timeout).Should(Equal(btc.DefaultClientTimeout))
			Expect(opts.TimeoutRetry).Should(Equal(btc.DefaultClientTimeoutRetry))
			Expect(opts.Host).Should(Equal(btc.DefaultClientHost))
			Expect(opts.User).Should(Equal(btc.DefaultClientUser))
			Expect(opts.Password).Should(Equal(btc.DefaultClientPassword))
		})
	})

	Context("when customizing the options", func() {
		It("should properly set the option", func() {
			opts := btc.DefaultClientOptions()

			By("Set the network")
			opts = opts.WithNet(&chaincfg.TestNet3Params)
			Expect(reflect.DeepEqual(opts.Net, &chaincfg.RegressionNetParams)).Should(BeFalse())
			Expect(reflect.DeepEqual(opts.Net, &chaincfg.TestNet3Params)).Should(BeTrue())

			By("Set the Timeout")
			rTimeout := time.Duration(rand.Int()) * time.Second
			opts = opts.WithTimeout(rTimeout)
			Expect(opts.Timeout).Should(Equal(rTimeout))

			By("Set the TimeoutRetry")
			rTimeoutRetry := time.Duration(rand.Int()) * time.Second
			opts = opts.WithTimeoutRetry(rTimeoutRetry)
			Expect(opts.TimeoutRetry).Should(Equal(rTimeoutRetry))

			By("Set the Host")
			rHost := testutil.RandomString()
			opts = opts.WithHost(rHost)
			Expect(opts.Host).Should(Equal(rHost))

			By("Set the User")
			rUser := testutil.RandomString()
			opts = opts.WithUser(rUser)
			Expect(opts.User).Should(Equal(rUser))

			By("Set the Password")
			rPassword := testutil.RandomString()
			opts = opts.WithPassword(rPassword)
			Expect(opts.Password).Should(Equal(rPassword))
		})
	})
})
