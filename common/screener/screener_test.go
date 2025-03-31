package screener_test

import (
	"context"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/common/screener"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Screener Tests", func() {
	Context("when testing a single address for screening", func() {
		It("should return whether the address is blacklisted", func(ctx context.Context) {
			By("Test a sanctioned address")
			chain, err := blockchain.ParseChainName(blockchain.Bitcoin)
			Expect(err).Should(BeNil())
			addr := "bc1qng0keqn7cq6p8qdt4rjnzdxrygnzq7nd0pju8q"
			addrs1 := map[string]blockchain.Chain{
				addr: chain,
			}

			scn := screener.New(url)
			res, err := scn.IsBlacklisted(ctx, addrs1)
			Expect(err).Should(BeNil())
			Expect(len(res)).Should(Equal(1))
			Expect(res[addr]).Should(BeTrue())
		})
	})

	Context("when testing multiple addresses for screening", func() {
		It("should return whether the addresses are blacklisted", func(ctx context.Context) {
			By("Test a sanctioned address")
			chain, err := blockchain.ParseChainName(blockchain.Bitcoin)
			Expect(err).Should(BeNil())
			addrs1 := map[string]blockchain.Chain{
				"bc1qhww67feqfdf6xasjat88x5stqa6vzx0c6fgtnj": chain,
				"bc1qng0keqn7cq6p8qdt4rjnzdxrygnzq7nd0pju8q": chain,
			}

			scn := screener.New(url)
			res, err := scn.IsBlacklisted(ctx, addrs1)
			Expect(err).Should(BeNil())
			Expect(len(res)).Should(Equal(1))
			for addr := range addrs1 {
				Expect(res[addr]).Should(BeTrue())
			}
		})
	})

	Context("when sending an invalid address to the API", func() {
		It("should return false", func(ctx context.Context) {
			By("Test a invalid address")
			chain, err := blockchain.ParseChainName(blockchain.Bitcoin)
			Expect(err).Should(BeNil())
			addr := "asd"
			addrs1 := map[string]blockchain.Chain{
				addr: chain,
			}

			scn := screener.New(url)
			_, err = scn.IsBlacklisted(ctx, addrs1)
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("when sending an empty list of addresses to the API", func() {
		It("should return false", func(ctx context.Context) {
			scn := screener.New(url)
			res, err := scn.IsBlacklisted(ctx, map[string]blockchain.Chain{})
			Expect(err).ShouldNot(BeNil())
			Expect(len(res)).Should(Equal(0))
		})
	})
})
