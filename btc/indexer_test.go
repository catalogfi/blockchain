package btc_test

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/localnet"
	"github.com/fatih/color"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Indexer client", func() {
	Context("electrs HTTP retry", func() {
		It("retries with backoff until the indexer responds successfully", func() {
			var attempts atomic.Int64
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				n := attempts.Add(1)
				if r.URL.Path != "/blocks/tip/height" {
					http.NotFound(w, r)
					return
				}
				if n < 3 {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				_, _ = w.Write([]byte("42"))
			}))
			defer srv.Close()

			c := btc.NewElectrsIndexerClient(zap.NewNop(), srv.URL, 5*time.Millisecond)
			tip, err := c.GetTipBlockHeight(context.Background())
			Expect(err).To(BeNil())
			Expect(tip).To(Equal(uint64(42)))
			Expect(attempts.Load()).To(BeNumerically(">=", 3))
		})

		It("stops waiting and returns when the context is cancelled during backoff", func() {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer srv.Close()

			c := btc.NewElectrsIndexerClient(zap.NewNop(), srv.URL, time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()
			_, err := c.GetTipBlockHeight(ctx)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(context.DeadlineExceeded.Error()))
		})
	})

	Context("GetTx deadline guard", func() {
		It("applies DefaultAPITimeout when context has no deadline", func() {
			// Server returns a valid tx JSON immediately, so the test completes fast
			// even though GetTx internally adds a 300s timeout.
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"txid":"deadbeef","version":2,"locktime":0,"vin":[],"vout":[],"size":0,"weight":0,"fee":0,"status":{"confirmed":false}}`)
			}))
			defer srv.Close()

			c := btc.NewElectrsIndexerClient(zap.NewNop(), srv.URL, 5*time.Millisecond)

			// Pass context.Background() directly — no deadline.
			// GetTx should add its own DefaultAPITimeout and succeed.
			tx, err := c.GetTx(context.Background(), "deadbeef")
			Expect(err).To(BeNil())
			Expect(tx.TxID).To(Equal("deadbeef"))
		})

		It("does not double-wrap when context already carries a deadline", func() {
			// Server returns a valid tx JSON on the first attempt.
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"txid":"abc123","version":2,"locktime":0,"vin":[],"vout":[],"size":0,"weight":0,"fee":0,"status":{"confirmed":false}}`)
			}))
			defer srv.Close()

			c := btc.NewElectrsIndexerClient(zap.NewNop(), srv.URL, 5*time.Millisecond)

			// Caller already sets a tight deadline — GetTx should skip wrapping.
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			tx, err := c.GetTx(ctx, "abc123")
			Expect(err).To(BeNil())
			Expect(tx.TxID).To(Equal("abc123"))
		})
	})

	Context("When using electrs API", func() {
		It("should be able to fetch the utxos of an address ", func() {
			By("GetUTXOs()")
			key, err := btcec.NewPrivateKey()
			Expect(err).To(BeNil())
			addr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(key.PubKey().SerializeCompressed()), network)
			Expect(err).To(BeNil())
			txid, err := localnet.FundBTC(addr.EncodeAddress())
			Expect(err).To(BeNil())
			time.Sleep(5 * time.Second)
			utxos, err := indexer.GetUTXOs(context.Background(), addr)
			Expect(err).To(BeNil())
			Expect(len(utxos)).Should(BeNumerically(">=", 1))
			exist := false
			for _, utxo := range utxos {
				if utxo.TxID == txid.String() {
					exist = true
					break
				}
			}
			Expect(exist).Should(BeTrue())

			By("GetUTXOsForAmount()")
			utxos, utxosValue, err := indexer.GetUTXOsForAmount(context.Background(), addr, 1e8)
			Expect(err).To(BeNil())
			Expect(utxosValue).Should(BeNumerically(">=", 1e8))
			Expect(len(utxos)).Should(BeNumerically(">=", 1))
			exist = false
			for _, utxo := range utxos {
				if utxo.TxID == txid.String() {
					exist = true
					break
				}
			}
			Expect(exist).Should(BeTrue())

			By("GetTipBlockHeight()")
			tip, err := indexer.GetTipBlockHeight(context.Background())
			Expect(err).To(BeNil())
			Expect(tip).Should(BeNumerically(">=", 100))

			By("GetTx()")
			tx, err := indexer.GetTx(context.Background(), txid.String())
			Expect(err).To(BeNil())
			Expect(tx.TxID).Should(Equal(txid.String()))
			Expect(tx.Status.Confirmed).Should(BeTrue())

			By("GetTxHex()")
			txHex, err := indexer.GetTxHex(context.Background(), txid.String())
			Expect(err).To(BeNil())
			Expect(txHex).ShouldNot(BeEmpty())
			txBytes, err := hex.DecodeString(txHex)
			Expect(err).To(BeNil())
			Expect(txBytes).ShouldNot(BeEmpty())
			btcTx, err := btcutil.NewTxFromBytes(txBytes)
			Expect(err).To(BeNil())
			Expect(btcTx.MsgTx().TxHash().String()).Should(Equal(txid.String()))

			By("SubmitTx()")
			amount, feeRate := int64(1e6), 10
			recipients := []btc.Recipient{
				{
					To:     addr.String(),
					Amount: amount,
				},
			}
			rawTx, err := btc.BuildTransaction(network, feeRate, btc.NewRawInputs(), utxos, btc.P2pkhUpdater, recipients, addr)
			Expect(err).To(BeNil())
			for i := range rawTx.TxIn {
				pkScript, err := txscript.PayToAddrScript(addr)
				Expect(err).To(BeNil())
				sigScript, err := txscript.SignatureScript(rawTx, i, pkScript, txscript.SigHashAll, key, true)
				Expect(err).To(BeNil())
				rawTx.TxIn[i].SignatureScript = sigScript
			}
			Expect(client.SubmitTx(context.Background(), rawTx)).Should(Succeed())

			By("FeeEstimate()")
			By("    --local env should not have enough data for the estimate")
			_, err = indexer.FeeEstimate(context.Background())
			Expect(err.Error()).Should(ContainSubstring("not enough data"))

			By("GetAddressTxs()")
			txs, err := indexer.GetAddressTxs(context.Background(), addr, "")
			Expect(err).To(BeNil())
			has := false
			for _, tx := range txs {
				if tx.TxID == txid.String() {
					has = true
				}
			}
			Expect(has).Should(BeTrue())
		})
	})

	Context("errors", func() {
		It("should return specific errors", func() {
			By("New address")
			privKey, err := btcec.NewPrivateKey()
			Expect(err).To(BeNil())
			pubKey := privKey.PubKey()
			pkAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), network)
			Expect(err).To(BeNil())

			By("funding the addresses")
			_, err = localnet.FundBTC(pkAddr.EncodeAddress())
			Expect(err).To(BeNil())
			time.Sleep(5 * time.Second)

			By("Construct a new tx")
			utxos, err := indexer.GetUTXOs(context.Background(), pkAddr)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e6), 10
			recipients := []btc.Recipient{
				{
					To:     pkAddr.EncodeAddress(),
					Amount: amount,
				},
			}
			transaction, err := btc.BuildTransaction(network, feeRate, btc.NewRawInputs(), utxos, btc.P2pkhUpdater, recipients, pkAddr)
			Expect(err).To(BeNil())
			for i := range transaction.TxIn {
				pkScript, err := txscript.PayToAddrScript(pkAddr)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(transaction, i, pkScript, txscript.SigHashAll, privKey, true)
				Expect(err).To(BeNil())
				transaction.TxIn[i].SignatureScript = sigScript
			}

			By("Submit the transaction")
			Expect(indexer.SubmitTx(context.Background(), transaction)).Should(Succeed())
			By(fmt.Sprintf("Funding tx hash = %v", color.YellowString(transaction.TxHash().String())))
			time.Sleep(time.Second)

			By("Expect a `ErrAlreadyInChain` error if the tx is already in a block")
			Expect(localnet.MineBTCBlock()).Should(Succeed())
			time.Sleep(1 * time.Second)
			err = indexer.SubmitTx(context.Background(), transaction)
			Expect(errors.Is(err, btc.ErrAlreadyInChain)).Should(BeTrue())

			By("Try construct a new transaction spending the same input")
			recipients1 := []btc.Recipient{
				{
					To:     pkAddr.EncodeAddress(),
					Amount: 2 * amount,
				},
			}
			transaction1, err := btc.BuildTransaction(network, feeRate, btc.NewRawInputs(), utxos, btc.P2pkhUpdater, recipients1, pkAddr)
			Expect(err).To(BeNil())
			for i := range transaction1.TxIn {
				pkScript, err := txscript.PayToAddrScript(pkAddr)
				Expect(err).To(BeNil())

				sigScript, err := txscript.SignatureScript(transaction1, i, pkScript, txscript.SigHashAll, privKey, true)
				Expect(err).To(BeNil())
				transaction1.TxIn[i].SignatureScript = sigScript
			}
			By("Expect a `ErrAlreadyInChain` error if the tx is already in a block")
			err = indexer.SubmitTx(context.Background(), transaction1)
			Expect(errors.Is(err, btc.ErrTxInputsMissingOrSpent)).Should(BeTrue())
		})
	})
})
