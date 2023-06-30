package btc_test

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/multichain/btc"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("fee estimator", func() {
	Context("estimating fee", func() {
		It("should return the live data from mempool", func() {
			estimator := btc.NewFeeEstimator(&chaincfg.MainNetParams, "https://mempool.space/api/v1/fees/recommended")
			fees, err := estimator.FeeSuggestion()
			Expect(err).Should(BeNil())

			Expect(fees.Minimum).Should(BeNumerically(">", 1))

			Expect(fees.Economy).Should(BeNumerically(">", 1))
			Expect(fees.Economy).Should(BeNumerically(">=", fees.Minimum))

			Expect(fees.Low).Should(BeNumerically(">", 1))
			Expect(fees.Low).Should(BeNumerically(">=", fees.Economy))

			Expect(fees.Medium).Should(BeNumerically(">", 1))
			Expect(fees.Medium).Should(BeNumerically(">=", fees.Low))

			Expect(fees.High).Should(BeNumerically(">", 1))
			Expect(fees.High).Should(BeNumerically(">=", fees.Medium))
		})
	})
})
