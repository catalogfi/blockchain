package evm_test

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm"
	"github.com/catalogfi/blockchain/evm/evmtest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Htlc", func() {

	chains := []blockchain.EvmChain{
		blockchain.NewEvmChain(blockchain.EthereumLocalnet),
		blockchain.NewEvmChain(blockchain.ArbitrumLocalnet),
	}

	for _, chain := range chains {
		chain := chain
		Context(fmt.Sprintf("%v", chain), func() {
			rpc := evmtest.MerryRpc(chain)
			client, err := ethclient.Dial(rpc)
			Expect(err).Should(BeNil())
			keys := evmtest.MerryKeys(2)
			htlcAddr := evmtest.MerryHtlcAddress(chain)
			walletOpts := evm.NewOptions(chain, []common.Address{htlcAddr}, 5*time.Second)
			wallet1, err := evm.NewWallet(walletOpts, keys[0], client)
			Expect(err).Should(BeNil())
			wallet2, err := evm.NewWallet(walletOpts, keys[1], client)
			Expect(err).Should(BeNil())

			It("should be able to initiate and redeem an htlc", func(ctx context.Context) {
				By("Init")
				secret, secretHash := evmtest.RandomSecret()
				htlc := evm.NewHtlc(wallet1.Address(), wallet2.Address(), htlcAddr, secretHash, big.NewInt(1e8), big.NewInt(100), chain.ChainID())
				tx, err := wallet1.Initiate(ctx, htlc)
				Expect(err).Should(BeNil())

				By("Wait for tx to be mined")
				_, err = bind.WaitMined(ctx, client, tx)
				Expect(err).Should(BeNil())

				By("Redeem")
				time.Sleep(time.Second)
				_, err = wallet2.Redeem(ctx, htlc, secret)
				Expect(err).Should(BeNil())
			})

			It("should be able to initiate and refund an htlc", func(ctx context.Context) {
				By("Init")
				_, secretHash := evmtest.RandomSecret()
				htlc := evm.NewHtlc(wallet1.Address(), wallet2.Address(), htlcAddr, secretHash, big.NewInt(1e8), big.NewInt(1), chain.ChainID())
				tx, err := wallet1.Initiate(ctx, htlc)
				Expect(err).Should(BeNil())

				By("Wait for a new block")
				receipt, err := bind.WaitMined(ctx, client, tx)
				Expect(err).Should(BeNil())
				for {
					header, err := client.BlockByNumber(ctx, nil)
					Expect(err).Should(BeNil())
					if header.NumberU64() > receipt.BlockNumber.Uint64() {
						break
					}
					time.Sleep(3 * time.Second)
				}

				By("Refund")
				_, err = wallet1.Refund(ctx, htlc)
				Expect(err).Should(BeNil())
			})

			It("should be able to do an instant refund", func(ctx context.Context) {
				By("Init")
				_, secretHash := evmtest.RandomSecret()
				htlc := evm.NewHtlc(wallet1.Address(), wallet2.Address(), htlcAddr, secretHash, big.NewInt(1e8), big.NewInt(100), chain.ChainID())
				tx, err := wallet1.Initiate(ctx, htlc)
				Expect(err).Should(BeNil())

				By("Wait for tx to be mined")
				_, err = bind.WaitMined(ctx, client, tx)
				Expect(err).Should(BeNil())

				By("Instant refund")
				_, err = wallet2.InstantRefund(ctx, htlc, nil)
				Expect(err).Should(BeNil())
			})
		})
	}
})
