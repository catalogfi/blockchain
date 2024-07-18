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

var _ = Describe("BatchWallet:CPFP", Ordered, func() {

	chainParams := &chaincfg.RegressionNetParams
	logger, err := zap.NewDevelopment()
	Expect(err).To(BeNil())

	indexer := btc.NewElectrsIndexerClient(logger, os.Getenv("BTC_REGNET_INDEXER"), time.Millisecond*500)

	privateKey, err := btcec.NewPrivateKey()
	Expect(err).To(BeNil())

	mockFeeEstimator := NewMockFeeEstimator(10)
	cache := NewTestCache()
	wallet, err := btc.NewBatcherWallet(privateKey, indexer, mockFeeEstimator, chainParams, cache, logger, btc.WithPTI(5*time.Second), btc.WithStrategy(btc.CPFP))
	Expect(err).To(BeNil())

	BeforeAll(func() {
		_, err := localnet.FundBitcoin(wallet.Address().EncodeAddress(), indexer)
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
				Amount: 100000,
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
		// change address
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))

		time.Sleep(10 * time.Second)
	})

	It("should be able to update fee with CPFP", func() {
		neeFeeRate := 20
		mockFeeEstimator.UpdateFee(neeFeeRate)

		time.Sleep(10 * time.Second)

		pendingBatches, err := cache.ReadPendingBatches(context.Background(), btc.CPFP)
		Expect(err).To(BeNil())

		for _, batch := range pendingBatches {
			feeRate := (batch.Tx.Fee * blockchain.WitnessScaleFactor) / int64(batch.Tx.Weight)
			Expect(feeRate).Should(BeNumerically(">=", neeFeeRate))
		}
	})

	It("should be able to send funds to multiple addresses", func() {
		pk1, err := btcec.NewPrivateKey()
		Expect(err).To(BeNil())
		address1, err := btc.PublicKeyAddress(chainParams, waddrmgr.WitnessPubKey, pk1.PubKey())
		Expect(err).To(BeNil())

		pk2, err := btcec.NewPrivateKey()
		Expect(err).To(BeNil())
		address2, err := btc.PublicKeyAddress(chainParams, waddrmgr.WitnessPubKey, pk2.PubKey())
		Expect(err).To(BeNil())

		req := []btc.SendRequest{
			{
				Amount: 100000,
				To:     address1,
			},
			{
				Amount: 100000,
				To:     address2,
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

		// first vout address
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(address1.EncodeAddress()))
		// second vout address
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(address2.EncodeAddress()))
		// change address
		Expect(tx.VOUTs[2].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
	})

	It("should be able to spend multiple scripts and send to multiple parties", func() {
		amount := int64(100000)

		p2wshSigCheckScript, p2wshSigCheckScriptAddr, err := sigCheckScript(*chainParams, privateKey)
		Expect(err).To(BeNil())

		p2wshAdditionScript, p2wshScriptAddr, err := additionScript(*chainParams)
		Expect(err).To(BeNil())

		p2trAdditionScript, p2trScriptAddr, cb, err := additionTapscript(*chainParams)
		Expect(err).To(BeNil())

		checkSigScript, checkSigScriptAddr, checkSigScriptCb, err := sigCheckTapScript(*chainParams, schnorr.SerializePubKey(privateKey.PubKey()))
		Expect(err).To(BeNil())

		err = fundScripts([]string{
			p2wshScriptAddr.EncodeAddress(),
			p2wshSigCheckScriptAddr.EncodeAddress(),
			p2trScriptAddr.EncodeAddress(),
			checkSigScriptAddr.EncodeAddress(),
		})
		Expect(err).To(BeNil())

		By("Let's create recipients")
		pk1, err := btcec.NewPrivateKey()
		Expect(err).To(BeNil())
		address1, err := btc.PublicKeyAddress(chainParams, waddrmgr.WitnessPubKey, pk1.PubKey())
		Expect(err).To(BeNil())

		pk2, err := btcec.NewPrivateKey()
		Expect(err).To(BeNil())
		address2, err := btc.PublicKeyAddress(chainParams, waddrmgr.WitnessPubKey, pk2.PubKey())
		Expect(err).To(BeNil())

		By("Send funds to Bob and Dave by spending the scripts")
		id, err := wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: amount,
				To:     address1,
			},
			{
				Amount: amount,
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

		for {
			fmt.Println("waiting for tx", id)
			tx, ok, err := wallet.Status(context.Background(), id)
			Expect(err).To(BeNil())
			if ok {
				Expect(tx).ShouldNot(BeNil())
				break
			}
			time.Sleep(5 * time.Second)
		}

		By("The tx should have 3 outputs")
		tx, _, err := wallet.Status(context.Background(), id)
		Expect(err).To(BeNil())
		Expect(tx).ShouldNot(BeNil())
		Expect(tx.VOUTs).Should(HaveLen(3))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(address1.EncodeAddress()))
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(address2.EncodeAddress()))
		Expect(tx.VOUTs[2].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))

		//Validate whether dave and bob received the right amount
		Expect(tx.VOUTs[0].Value).Should(Equal(int(amount)))
		Expect(tx.VOUTs[1].Value).Should(Equal(int(amount)))

	})
})

func fundScripts(addresses []string) error {
	for _, address := range addresses {
		_, err := localnet.FundBitcoin(address, indexer)
		if err != nil {
			return err
		}
	}
	return nil
}
