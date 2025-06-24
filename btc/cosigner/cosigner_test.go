package cosigner_test

import (
	"context"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/catalogfi/blockchain/btc/cosigner"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cosigner", func() {
	Context("Spending from the cosigner address", func() {
		It("should able to spend funds in the cosigner address", func(ctx context.Context) {
			By("Initialize keys and the cosigner address")
			key1, err := btcec.NewPrivateKey()
			Expect(err).To(BeNil())
			key2, err := btcec.NewPrivateKey()
			Expect(err).To(BeNil())
			script, err := cosigner.Script(key1.PubKey().SerializeCompressed(), key2.PubKey().SerializeCompressed(), 144*180)
			Expect(err).To(BeNil())
			addr, err := btc.P2wshAddress(script, network)
			Expect(err).To(BeNil())

			By("Fund the cosigner address")
			fundTx, err := btctest.Faucet(addr.EncodeAddress())
			Expect(err).To(BeNil())
			Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(fundTx.String()))).Should(Succeed())

			By("Construct a tx to spend the funds")
			utxos, err := indexer.GetUTXOs(ctx, addr)
			Expect(err).To(BeNil())
			sizer := btc.NewSizeEstimator(cosigner.BaseSizeSpend, cosigner.SegwitSizeSpend, utxos...)
			feeMode := btc.MinFeeRateMode(1e3, sizer)
			tx1, err := btc.BuildTx(network, feeMode, utxos, nil, nil, addr)
			Expect(err).To(BeNil())

			By("Sign and submit the tx")
			fetcher, err := btc.NewFetcher(script, utxos...)
			Expect(err).To(BeNil())
			for i := range tx1.TxIn {
				sig1, err := cosigner.Sign(script, tx1, i, fetcher, key1)
				Expect(err).To(BeNil())
				sig2, err := cosigner.Sign(script, tx1, i, fetcher, key2)
				Expect(err).To(BeNil())
				tx1.TxIn[i].Witness = cosigner.Witness(script, sig1, sig2, false)
			}
			Expect(indexer.SubmitTx(ctx, tx1)).Should(Succeed())
		})

		It("should able to refund funds after certain timestamp", func(ctx context.Context) {
			By("Initialize keys and the cosigner address")
			waitTime := int64(6)
			key1, err := btcec.NewPrivateKey()
			Expect(err).To(BeNil())
			key2, err := btcec.NewPrivateKey()
			Expect(err).To(BeNil())
			script, err := cosigner.Script(key1.PubKey().SerializeCompressed(), key2.PubKey().SerializeCompressed(), waitTime)
			Expect(err).To(BeNil())
			addr, err := btc.P2wshAddress(script, network)
			Expect(err).To(BeNil())

			By("Fund the cosigner address")
			fundTx, err := btctest.Faucet(addr.EncodeAddress())
			Expect(err).To(BeNil())
			Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(fundTx.String()))).Should(Succeed())

			By("Construct a tx to refund the funds")
			utxos, err := indexer.GetUTXOs(ctx, addr)
			Expect(err).To(BeNil())
			sizer := btc.NewSizeEstimator(cosigner.BaseSizeRefund, cosigner.SegwitSizeRefund(waitTime), utxos...)
			feeMode := btc.MinFeeRateMode(1e3, sizer)
			tx1, err := btc.BuildTx(network, feeMode, utxos, nil, nil, addr)
			Expect(err).To(BeNil())
			Expect(btctest.NewBlockWaitMined(6, indexer)).Should(Succeed())

			By("Sign and submit the tx")
			fetcher, err := btc.NewFetcher(script, utxos...)
			Expect(err).To(BeNil())
			for i := range tx1.TxIn {
				tx1.TxIn[i].Sequence = uint32(waitTime)
			}
			for i := range tx1.TxIn {
				sig2, err := cosigner.Sign(script, tx1, i, fetcher, key2)
				Expect(err).To(BeNil())
				tx1.TxIn[i].Witness = cosigner.Witness(script, nil, sig2, true)
			}
			Expect(indexer.SubmitTx(ctx, tx1)).Should(Succeed())
		})
	})
})
