package btc_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/localnet"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("BatchWallet:RBF", Ordered, func() {

	chainParams := &chaincfg.RegressionNetParams
	logger, err := zap.NewDevelopment()
	Expect(err).To(BeNil())

	indexer := btc.NewElectrsIndexerClient(logger, os.Getenv("BTC_REGNET_INDEXER"), time.Millisecond*500)

	privateKey, err := btcec.NewPrivateKey()
	Expect(err).To(BeNil())

	requiredFeeRate := int64(10)

	mockFeeEstimator := NewMockFeeEstimator(int(requiredFeeRate))
	cache := NewTestCache(btc.RBF)
	wallet, err := btc.NewBatcherWallet(privateKey, indexer, mockFeeEstimator, chainParams, cache, logger, btc.WithPTI(5*time.Second), btc.WithStrategy(btc.RBF))
	Expect(err).To(BeNil())

	faucet, err := btc.NewSimpleWallet(privateKey, chainParams, indexer, mockFeeEstimator, btc.HighFee)
	Expect(err).To(BeNil())

	defaultAmount := int64(100000)

	p2wshSigCheckScript, p2wshSigCheckScriptAddr, err := sigCheckScript(*chainParams, privateKey)
	Expect(err).To(BeNil())

	p2wshAdditionScript, p2wshScriptAddr, err := additionScript(*chainParams)
	Expect(err).To(BeNil())

	p2trAdditionScript, p2trScriptAddr, cb, err := additionTapscript(*chainParams)
	Expect(err).To(BeNil())

	checkSigScript, checkSigScriptAddr, checkSigScriptCb, err := sigCheckTapScript(*chainParams, schnorr.SerializePubKey(privateKey.PubKey()))
	Expect(err).To(BeNil())

	p2wshSigCheckScript2, p2wshSigCheckScriptAddr2, err := sigCheckScript(*chainParams, privateKey)
	Expect(err).To(BeNil())

	randAddr, err := randomP2wpkhAddress(*chainParams)
	Expect(err).To(BeNil())

	var sacp []byte

	pk1, err := btcec.NewPrivateKey()
	Expect(err).To(BeNil())
	address1, err := btc.PublicKeyAddress(chainParams, waddrmgr.WitnessPubKey, pk1.PubKey())
	Expect(err).To(BeNil())

	pk2, err := btcec.NewPrivateKey()
	Expect(err).To(BeNil())
	address2, err := btc.PublicKeyAddress(chainParams, waddrmgr.WitnessPubKey, pk2.PubKey())
	Expect(err).To(BeNil())

	BeforeAll(func() {
		_, err := localnet.FundBitcoin(wallet.Address().EncodeAddress(), indexer)
		Expect(err).To(BeNil())

		_, err = localnet.FundBitcoin(faucet.Address().EncodeAddress(), indexer)
		Expect(err).To(BeNil())

		sacp, err = generateSACP(faucet, *chainParams, privateKey, randAddr, 10000, 100)
		Expect(err).To(BeNil())

		faucetTx, err := faucet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: defaultAmount,
				To:     p2wshSigCheckScriptAddr,
			},
			{
				Amount: defaultAmount,
				To:     p2wshScriptAddr,
			},
			{
				Amount: defaultAmount,
				To:     p2trScriptAddr,
			},
			{
				Amount: defaultAmount,
				To:     checkSigScriptAddr,
			},
			{
				Amount: defaultAmount,
				To:     p2wshSigCheckScriptAddr2,
			},
		}, nil, nil)
		Expect(err).To(BeNil())
		fmt.Println("funded scripts", "txid :", faucetTx)

		err = localnet.MineBitcoinBlocks(1, indexer)
		Expect(err).To(BeNil())

		err = wallet.Start(context.Background())
		Expect(err).To(BeNil())

	})

	AfterAll(func() {
		err := wallet.Stop()
		Expect(err).To(BeNil())
	})

	It("should be able to send funds", func() {
		req := []btc.SendRequest{
			{
				Amount: defaultAmount,
				To:     wallet.Address(),
			},
		}

		id, err := wallet.Send(context.Background(), req, nil, nil)
		Expect(err).To(BeNil())

		var tx btc.Transaction
		var ok bool

		for {
			fmt.Println("waiting for tx", id)
			tx, ok, err = wallet.Status(context.Background(), id)
			Expect(err).To(BeNil())
			if ok {
				Expect(tx).ShouldNot(BeNil())
				break
			}
			time.Sleep(5 * time.Second)
		}

		// to address
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))

		time.Sleep(10 * time.Second)
	})

	It("should be able to update fee with RBF", func() {
		mockFeeEstimator.UpdateFee(int(requiredFeeRate) + 10)
		time.Sleep(10 * time.Second)
		lb, err := cache.ReadLatestBatch(context.Background())
		Expect(err).To(BeNil())

		feeRate := (lb.Tx.Fee * blockchain.WitnessScaleFactor) / int64(lb.Tx.Weight)
		Expect(feeRate).Should(BeNumerically(">=", int(requiredFeeRate)+10))

		requiredFeeRate += 10
	})

	It("should be able to spend multiple scripts and send to multiple parties", func() {
		defaultAmount := int64(defaultAmount)

		By("Send funds to Bob and Dave by spending the scripts")
		id, err := wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: defaultAmount,
				To:     address1,
			},
			{
				Amount: defaultAmount,
				To:     address1,
			},
			{
				Amount: defaultAmount,
				To:     address1,
			},
			{
				Amount: defaultAmount,
				To:     address2,
			},
		}, []btc.SpendRequest{
			{
				Witness: [][]byte{
					{0x1},
					{0x1},
					p2wshAdditionScript,
				},
				Script:        p2wshAdditionScript,
				ScriptAddress: p2wshScriptAddr,
				HashType:      txscript.SigHashAll,
			},
			{
				Witness: [][]byte{
					btc.AddSignatureSegwitOp,
					p2wshSigCheckScript,
				},
				Script:        p2wshSigCheckScript,
				ScriptAddress: p2wshSigCheckScriptAddr,
				HashType:      txscript.SigHashAll,
			},
			{
				Witness: [][]byte{
					{0x1},
					{0x1},
					p2trAdditionScript,
					cb,
				},
				Leaf:          txscript.NewTapLeaf(0xc0, p2trAdditionScript),
				ScriptAddress: p2trScriptAddr,
			},
			{
				Witness: [][]byte{
					btc.AddSignatureSchnorrOp,
					checkSigScript,
					checkSigScriptCb,
				},
				Leaf:          txscript.NewTapLeaf(0xc0, checkSigScript),
				ScriptAddress: checkSigScriptAddr,
				HashType:      txscript.SigHashAll,
			},
		}, nil)
		Expect(err).To(BeNil())
		Expect(id).ShouldNot(BeEmpty())

		var tx btc.Transaction
		var ok bool
		for {
			fmt.Println("waiting for tx", id)
			tx, ok, err = wallet.Status(context.Background(), id)
			Expect(err).To(BeNil())
			if ok {
				Expect(tx).ShouldNot(BeNil())
				break
			}
			time.Sleep(5 * time.Second)
		}

		By("The tx should have 5 outputs")
		Expect(tx).ShouldNot(BeNil())
		Expect(tx.VOUTs).Should(HaveLen(5))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(address1.EncodeAddress()))
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(address1.EncodeAddress()))
		Expect(tx.VOUTs[2].ScriptPubKeyAddress).Should(Equal(address1.EncodeAddress()))
		Expect(tx.VOUTs[3].ScriptPubKeyAddress).Should(Equal(address2.EncodeAddress()))
		Expect(tx.VOUTs[4].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))

		lb, err := cache.ReadLatestBatch(context.Background())
		Expect(err).To(BeNil())
		feeRate := (lb.Tx.Fee * blockchain.WitnessScaleFactor) / int64(lb.Tx.Weight)
		Expect(feeRate).Should(BeNumerically(">=", requiredFeeRate+10))
	})

	It("should be able to update fee with RBF", func() {
		mockFeeEstimator.UpdateFee(int(requiredFeeRate) + 10)
		time.Sleep(10 * time.Second)
		lb, err := cache.ReadLatestBatch(context.Background())
		Expect(err).To(BeNil())

		feeRate := (lb.Tx.Fee * blockchain.WitnessScaleFactor) / int64(lb.Tx.Weight)
		Expect(feeRate).Should(BeNumerically(">=", int(requiredFeeRate)+10))

		requiredFeeRate += 10
	})

	It("should do nothing if fee decreases", func() {
		mockFeeEstimator.UpdateFee(int(requiredFeeRate) - 10)
		time.Sleep(10 * time.Second)
		lb, err := cache.ReadLatestBatch(context.Background())
		Expect(err).To(BeNil())

		feeRate := (lb.Tx.Fee * blockchain.WitnessScaleFactor) / int64(lb.Tx.Weight)
		Expect(feeRate).Should(BeNumerically(">=", int(requiredFeeRate)))
	})

	It("should be able to mix SACPs with spend requests", func() {
		id, err := wallet.Send(context.Background(), nil, []btc.SpendRequest{
			{
				Witness: [][]byte{
					btc.AddSignatureSegwitOp,
					p2wshSigCheckScript2,
				},
				Script:        p2wshSigCheckScript2,
				ScriptAddress: p2wshSigCheckScriptAddr2,
				HashType:      txscript.SigHashAll,
			},
		}, [][]byte{sacp})
		Expect(err).To(BeNil())
		Expect(id).ShouldNot(BeEmpty())

		var tx btc.Transaction
		var ok bool

		for {
			fmt.Println("waiting for tx", id)
			tx, ok, err = wallet.Status(context.Background(), id)
			Expect(err).To(BeNil())
			if ok {
				Expect(tx).ShouldNot(BeNil())
				break
			}
			time.Sleep(5 * time.Second)
		}

		// Three outputs, one for the script, one for the SACP recipient
		Expect(tx.VOUTs).Should(HaveLen(6))

		// SACP recipient is the wallet itself
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(randAddr.EncodeAddress()))
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(address1.EncodeAddress()))
		Expect(tx.VOUTs[2].ScriptPubKeyAddress).Should(Equal(address1.EncodeAddress()))
		Expect(tx.VOUTs[3].ScriptPubKeyAddress).Should(Equal(address1.EncodeAddress()))
		Expect(tx.VOUTs[4].ScriptPubKeyAddress).Should(Equal(address2.EncodeAddress()))
		Expect(tx.VOUTs[5].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))

		lb, err := cache.ReadLatestBatch(context.Background())
		Expect(err).To(BeNil())
		feeRate := (lb.Tx.Fee * blockchain.WitnessScaleFactor) / int64(lb.Tx.Weight)
		Expect(feeRate).Should(BeNumerically(">=", requiredFeeRate+10))
	})
})
