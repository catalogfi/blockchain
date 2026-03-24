package btc_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/blockchain/btc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

// delayProxy is a reverse proxy that adds configurable delays to matching requests.
type delayProxy struct {
	target *url.URL

	mu       sync.RWMutex
	rules    []delayRule
	reqCount atomic.Int64
}

type delayRule struct {
	// match returns true if this rule should apply to the request.
	match func(r *http.Request, body []byte) bool
	delay time.Duration
}

func newDelayProxy(target string) *delayProxy {
	u, err := url.Parse(target)
	Expect(err).To(BeNil())
	return &delayProxy{target: u}
}

func (p *delayProxy) addRule(match func(r *http.Request, body []byte) bool, delay time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rules = append(p.rules, delayRule{match: match, delay: delay})
}

func (p *delayProxy) clearRules() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rules = nil
}

func (p *delayProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.reqCount.Add(1)

	// Read the body so we can inspect it for JSON-RPC method matching,
	// then replay it to the upstream.
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close()

	// Check delay rules
	p.mu.RLock()
	for _, rule := range p.rules {
		if rule.match(r, bodyBytes) {
			p.mu.RUnlock()
			time.Sleep(rule.delay)
			p.mu.RLock()
			break
		}
	}
	p.mu.RUnlock()

	// Build upstream request
	upstreamURL := *p.target
	upstreamURL.Path = r.URL.Path
	upstreamURL.RawQuery = r.URL.RawQuery

	bodyReader := bytes.NewReader(bodyBytes)
	upReq, err := http.NewRequestWithContext(r.Context(), r.Method, upstreamURL.String(), bodyReader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	upReq.Header = r.Header.Clone()
	upReq.ContentLength = int64(len(bodyBytes))

	resp, err := http.DefaultClient.Do(upReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, vs := range resp.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// mineAndWait mines n blocks and waits for the indexer to index them.
func mineAndWait(proxyIndexer btc.IndexerClient, n int) {
	heightBefore, err := proxyIndexer.GetTipBlockHeight(context.Background())
	Expect(err).To(BeNil())

	out, cmdErr := exec.Command("merry", "rpc", "--generate", fmt.Sprintf("%d", n)).CombinedOutput()
	Expect(cmdErr).To(BeNil(), "merry rpc --generate failed: %s", string(out))

	Eventually(func() bool {
		h, err := proxyIndexer.GetTipBlockHeight(context.Background())
		if err != nil {
			return false
		}
		return h >= heightBefore+uint64(n)
	}, 30*time.Second, 1*time.Second).Should(BeTrue(), "indexer should catch up to new block height")
}

var _ = Describe("BatchWallet:RBF:DelayProxy", Ordered, func() {
	chainParams := &chaincfg.RegressionNetParams
	logger, err := zap.NewDevelopment()
	Expect(err).To(BeNil())

	privateKey, err := btcec.NewPrivateKey()
	Expect(err).To(BeNil())

	mockFeeEstimator := NewMockFeeEstimator(1)
	dbPath := "./test_rbf_proxy"

	// Real backends
	const realIndexerURL = "http://0.0.0.0:30000"
	const realNodeURL = "http://0.0.0.0:18443"
	const nodeUser = "admin1"
	const nodePass = "123"

	// Proxies
	indexerProxy := newDelayProxy(realIndexerURL)
	nodeProxy := newDelayProxy(realNodeURL)

	var indexerServer *httptest.Server
	var nodeServer *httptest.Server
	var wallet btc.BatcherWallet
	var cache btc.Cache
	var proxyIndexer btc.IndexerClient

	BeforeAll(func() {
		// Start proxy servers
		indexerServer = httptest.NewServer(indexerProxy)
		nodeServer = httptest.NewServer(nodeProxy)

		// Create indexer client pointing at the proxy
		proxyIndexer = btc.NewElectrsIndexerClient(logger, indexerServer.URL, 2*time.Second)

		// Create RPC client pointing at the node proxy
		nodeProxyURL := nodeServer.URL
		bitcoinRPC := btc.NewBitcoinClient(nodeUser, nodePass, nodeProxyURL)

		db, err := leveldb.OpenFile(dbPath, nil)
		Expect(err).To(BeNil())

		cache = btc.NewBatcherCache(db, "proxy_test", btc.RBF)
		wallet, err = btc.NewBatcherWallet(
			privateKey,
			proxyIndexer,
			mockFeeEstimator,
			chainParams,
			cache,
			logger,
			&bitcoinRPC,
			btc.WithPTI(5*time.Second),
			btc.WithStrategy(btc.RBF),
		)
		Expect(err).To(BeNil())

		// Fund the wallet via merry faucet directly
		out, cmdErr := exec.Command("merry", "faucet", "--to", wallet.Address().EncodeAddress()).CombinedOutput()
		Expect(cmdErr).To(BeNil(), "merry faucet failed: %s", string(out))
		fmt.Printf("[DelayProxy] Funded wallet: %s\n", strings.TrimSpace(string(out)))

		// Mine a block so the funding tx confirms
		mineAndWait(proxyIndexer, 1)

		// Wait for the indexer to have confirmed UTXOs for the wallet.
		Eventually(func() bool {
			utxos, err := proxyIndexer.GetUTXOs(context.Background(), wallet.Address())
			if err != nil {
				return false
			}
			return len(utxos) > 0
		}, 60*time.Second, 2*time.Second).Should(BeTrue(), "wallet should have UTXOs after funding+mining")

		// Start the wallet once — it stays running across all tests
		err = wallet.Start(context.Background())
		Expect(err).To(BeNil())
	})

	AfterAll(func() {
		if wallet != nil {
			wallet.Stop()
		}
		indexerServer.Close()
		nodeServer.Close()
		os.RemoveAll(dbPath)
	})

	// AfterEach: mine blocks to confirm the current batch so the batcher
	// doesn't get stuck on updateRBF/getmempooldescendants for the next test.
	AfterEach(func() {
		indexerProxy.clearRules()
		nodeProxy.clearRules()
		mineAndWait(proxyIndexer, 1)
		// Give the batcher a couple of ticks to recognize the confirmation
		time.Sleep(10 * time.Second)
	})

	Context("when GetTx is delayed beyond a single retry interval", func() {
		It("should still persist the batch after retrying", func() {
			// Delay GetTx responses by 5s — longer than the 2s retry interval,
			// but within the 300s DefaultAPITimeout. The per-request context
			// timeout (retryInterval=2s) should cause it to retry, and eventually
			// the request will succeed when we drop the delay.
			getTxCallCount := atomic.Int64{}
			indexerProxy.addRule(func(r *http.Request, body []byte) bool {
				// Match GET /tx/<txid> (not POST /tx which is SubmitTx)
				if r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/tx/") {
					count := getTxCallCount.Add(1)
					// Delay the first 3 calls, then let subsequent ones pass fast.
					// This simulates the indexer being slow to index initially.
					if count <= 3 {
						return true
					}
				}
				return false
			}, 5*time.Second)

			// Send a transaction
			req := []btc.SendRequest{
				{
					Amount: 10000,
					To:     wallet.Address(),
				},
			}

			id, err := wallet.Send(context.Background(), req, nil, nil)
			Expect(err).To(BeNil())

			// Wait for the batch to be processed
			var tx btc.Transaction
			var ok bool
			Eventually(func() bool {
				tx, ok, err = wallet.Status(context.Background(), id)
				Expect(err).To(BeNil())
				return ok
			}, 120*time.Second, 3*time.Second).Should(BeTrue())
			Expect(tx.TxID).ShouldNot(BeEmpty())

			// Verify the batch was persisted in cache
			batch, err := cache.ReadLatestBatch(context.Background())
			Expect(err).To(BeNil())
			Expect(batch.Tx.TxID).ShouldNot(BeEmpty())

			fmt.Printf("[DelayProxy] Test1: Batch persisted with txid=%s after %d GetTx calls\n", batch.Tx.TxID, getTxCallCount.Load())
		})
	})

	Context("when GetTx is completely unreachable but node has the tx", func() {
		It("should find the tx via node mempool and keep retrying indexer", func() {
			// Block ALL GetTx responses with a huge delay (simulating indexer down)
			// but let SubmitTx pass through. The node will have it, so the retry
			// loop should detect it via GetMempoolEntry and keep trying.
			indexerProxy.addRule(func(r *http.Request, body []byte) bool {
				return r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/tx/")
			}, 300*time.Second) // effectively blocks the call — will hit the per-request 2s timeout

			req := []btc.SendRequest{
				{
					Amount: 10000,
					To:     wallet.Address(),
				},
			}

			id, err := wallet.Send(context.Background(), req, nil, nil)
			Expect(err).To(BeNil())

			// Give it enough time for the first indexer attempt to fail and hit the node check
			time.Sleep(15 * time.Second)

			// Now remove the delay — simulating indexer recovery
			indexerProxy.clearRules()

			// The retry loop should eventually succeed
			var tx btc.Transaction
			var ok bool
			Eventually(func() bool {
				tx, ok, err = wallet.Status(context.Background(), id)
				Expect(err).To(BeNil())
				return ok
			}, 120*time.Second, 3*time.Second).Should(BeTrue())
			Expect(tx.TxID).ShouldNot(BeEmpty())

			batch, err := cache.ReadLatestBatch(context.Background())
			Expect(err).To(BeNil())
			Expect(batch.Tx.TxID).ShouldNot(BeEmpty())

			fmt.Printf("[DelayProxy] Test2: Batch persisted after indexer recovery, txid=%s\n", batch.Tx.TxID)
		})
	})

	Context("when both indexer and node are slow", func() {
		It("should re-submit and eventually persist the batch", func() {
			// Delay both GetTx and getmempoolentry initially
			callCount := atomic.Int64{}

			indexerProxy.addRule(func(r *http.Request, body []byte) bool {
				if r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/tx/") {
					count := callCount.Add(1)
					return count <= 2 // delay first 2 GetTx calls
				}
				return false
			}, 5*time.Second)

			nodeProxy.addRule(func(r *http.Request, body []byte) bool {
				return strings.Contains(string(body), "getmempoolentry")
			}, 4*time.Second)

			req := []btc.SendRequest{
				{
					Amount: 10000,
					To:     wallet.Address(),
				},
			}

			id, err := wallet.Send(context.Background(), req, nil, nil)
			Expect(err).To(BeNil())

			// Wait and then drop delays
			time.Sleep(10 * time.Second)
			indexerProxy.clearRules()
			nodeProxy.clearRules()

			var tx btc.Transaction
			var ok bool
			Eventually(func() bool {
				tx, ok, err = wallet.Status(context.Background(), id)
				Expect(err).To(BeNil())
				return ok
			}, 120*time.Second, 3*time.Second).Should(BeTrue())
			Expect(tx.TxID).ShouldNot(BeEmpty())

			batch, err := cache.ReadLatestBatch(context.Background())
			Expect(err).To(BeNil())
			Expect(batch.Tx.TxID).ShouldNot(BeEmpty())

			fmt.Printf("[DelayProxy] Test3: Batch persisted despite dual delays, txid=%s\n", batch.Tx.TxID)
		})
	})
})
