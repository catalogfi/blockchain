package btc_test

import (
	"context"
	"errors"
	"log"
	"sort"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	"github.com/fatih/color"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Bitcoin", func() {
	Context("Build a transaction", func() {
		It("should be able to build a transaction", func(ctx context.Context) {
			By("Initialize keys and addresses")
			addrType := waddrmgr.PubKeyHash
			privKey1, p2pkhAddr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).To(BeNil())
			_, p2pkhAddr2, err := btctest.NewBtcKey(network, addrType)
			Expect(err).To(BeNil())

			By("Construct a transaction which sends money from p2pkhAddr1 to p2pkhAddr2")
			utxos, err := indexer.GetUTXOs(ctx, p2pkhAddr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e7), 4
			recipients := btc.SingleRecipient(p2pkhAddr2.EncodeAddress(), amount)
			sizer := btc.NewSizeEstimator(utxos, btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH)
			transaction, err := btc.BuildTransaction(network, feeRate, nil, utxos, sizer, recipients, p2pkhAddr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.SignTx(addrType, transaction, privKey1, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
			By(color.GreenString("tx hash = %v", transaction.TxHash().String()))
		})

		It("should not use a utxo if its uneconomical", func(ctx context.Context) {
			By("Initialize keys and addresses")
			addrType := waddrmgr.TaprootPubKey
			privKey1, p2trAddr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).To(BeNil())
			_, p2trAddr2, err := btctest.NewBtcKey(network, addrType)
			Expect(err).To(BeNil())

			By("Create a utxo which is less than the dust amount")
			utxos, err := indexer.GetUTXOs(ctx, p2trAddr1)
			Expect(err).To(BeNil())
			recipients := btc.SingleRecipient(p2trAddr1.EncodeAddress(), btc.DustAmount)
			sizer := btc.NewSizeEstimator(utxos, btc.BaseSizeP2TR, btc.SegwitSizeP2TR)
			transaction, err := btc.BuildTransaction(network, 3, nil, utxos, sizer, recipients, p2trAddr1)
			Expect(err).To(BeNil())
			Expect(btc.SignTx(addrType, transaction, privKey1, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
			Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

			By("Construct a transaction try using the dust amount")
			utxos1, err := indexer.GetUTXOs(ctx, p2trAddr1)
			Expect(err).To(BeNil())
			sort.Slice(utxos1, func(i, j int) bool {
				return utxos1[i].Amount < utxos1[j].Amount
			})
			Expect(utxos1[0].TxID).Should(Equal(transaction.TxHash().String()))
			amount, feeRate := int64(1e7), 3
			recipients1 := btc.SingleRecipient(p2trAddr2.EncodeAddress(), amount)
			sizer1 := btc.NewSizeEstimator(utxos1, btc.BaseSizeP2TR, btc.SegwitSizeP2TR)
			transaction1, err := btc.BuildTransaction(network, feeRate, nil, utxos1, sizer1, recipients1, p2trAddr1)
			Expect(err).To(BeNil())

			By("The new tx should not use the uneconomical utxo")
			for _, input := range transaction1.TxIn {
				if input.PreviousOutPoint.String() == utxos1[0].String() {
					Fail("The new tx should not use the utxo which is less than the dust amount")
				}
			}
		})
	})

	Context("RBF", func() {
		It("should be able to build and replace an RBF transaction", func(ctx context.Context) {
			By("Initialize keys and addresses")
			privKey1, p2pkhAddr1, err := btctest.NewBtcAddrWithFunds(network, waddrmgr.PubKeyHash, indexer)
			Expect(err).To(BeNil())
			privKey2, p2pkhAddr2, err := btctest.NewBtcKey(network, waddrmgr.TaprootPubKey)
			Expect(err).To(BeNil())

			By("Construct a RBF tx which sends money from p2pkhAddr1 to p2pkhAddr2")
			utxos, err := indexer.GetUTXOs(ctx, p2pkhAddr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e7), 10
			recipients := btc.SingleRecipient(p2pkhAddr2.EncodeAddress(), amount)
			sizer := btc.NewSizeEstimator(utxos, btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH)
			transaction, err := btc.BuildTransaction(network, feeRate, nil, utxos, sizer, recipients, p2pkhAddr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.SignTx(waddrmgr.PubKeyHash, transaction, privKey1, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
			By(color.GreenString("RBF tx hash = %v", transaction.TxHash().String()))

			By("Build a descent tx with higher fee")
			time.Sleep(5 * time.Second)
			utxos1, err := indexer.GetUTXOs(ctx, p2pkhAddr2)
			Expect(err).To(BeNil())
			sizer.AddUtxos(utxos1, btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH)
			transaction1, err := btc.BuildTransaction(network, 14, utxos1, nil, sizer, nil, p2pkhAddr2)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.SignTx(waddrmgr.TaprootPubKey, transaction1, privKey2, utxos1)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction1)).Should(Succeed())
			By(color.GreenString("RBF tx hash = %v", transaction.TxHash().String()))
			time.Sleep(5 * time.Second)

			By("Build a replacement tx with higher fee")
			feeRate += 12
			replaceTx, err := btc.BuildTransaction(network, feeRate, nil, utxos, sizer, recipients, p2pkhAddr1)
			Expect(err).To(BeNil())
			replaceTx.TxOut[len(replaceTx.TxOut)-1].Value = replaceTx.TxOut[len(replaceTx.TxOut)-1].Value - 235

			By("Sign and submit the replacement tx")
			Expect(btc.SignTx(waddrmgr.PubKeyHash, replaceTx, privKey1, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, replaceTx)).Should(Succeed())
			By(color.GreenString("Replaced RBF tx hash = %v", replaceTx.TxHash().String()))
		})

		It("should be able to build and replace an RBF transaction", func(ctx context.Context) {
			By("Initialize keys and addresses")
			addrType := waddrmgr.TaprootPubKey
			privKey1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, nil)
			Expect(err).To(BeNil())
			_, err = btctest.Faucet(addr1.EncodeAddress())
			Expect(err).To(BeNil())
			_, addr2, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).To(BeNil())

			By("Construct a RBF tx which sends money from addr1 to addr2")
			utxos, err := indexer.GetUTXOs(ctx, addr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e7), 10
			recipients := btc.SingleRecipient(addr2.EncodeAddress(), amount)
			sizer := btc.NewSizeEstimator(utxos, btc.BaseSizeP2TR, btc.SegwitSizeP2TR)
			transaction, err := btc.BuildTransaction(network, feeRate, nil, utxos, sizer, recipients, addr1)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			transaction.TxOut[len(transaction.TxOut)-1].Value = transaction.TxOut[len(transaction.TxOut)-1].Value - 22
			Expect(btc.SignTx(addrType, transaction, privKey1, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
			By(color.GreenString("RBF tx hash = %v", transaction.TxHash().String()))

			By("Add one utxo")
			queryTx, err := indexer.GetTx(ctx, transaction.TxHash().String())
			Expect(err).To(BeNil())
			vsize := (queryTx.Weight + 3) / 4
			prevFeeRate := int(queryTx.Fee) * 1000 / vsize
			log.Printf("min fee rate = %v  fee = %v ", prevFeeRate, queryTx.Fee)
			transaction1, err := btc.BuildRbfTransaction(network, prevFeeRate, queryTx.Fee, utxos, nil, sizer, recipients, addr1)
			Expect(err).To(BeNil())

			By("The fee we choose should be optimal, which means it should be rejected if it's using 1 less sat")
			transaction1.TxOut[len(transaction1.TxOut)-1].Value = transaction1.TxOut[len(transaction1.TxOut)-1].Value + 1
			Expect(btc.SignTx(addrType, transaction1, privKey1, utxos)).Should(Succeed())
			err = indexer.SubmitTx(ctx, transaction1)
			log.Printf("1 err = %v", err)
			Expect(err).ShouldNot(BeNil())

			By("Pump fee by 1 and it should be accepted")
			transaction1.TxOut[len(transaction1.TxOut)-1].Value = transaction1.TxOut[len(transaction1.TxOut)-1].Value - 1
			Expect(btc.SignTx(addrType, transaction1, privKey1, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction1)).Should(Succeed())
			log.Printf("tx2 = %v", transaction1.TxHash().String())
		})

		It("should get an error when trying to replace a mined tx", func(ctx context.Context) {
			By("Initialize keys and addresses")
			addrType := waddrmgr.PubKeyHash
			privKey, pkAddr, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).To(BeNil())
			_, toAddr, err := btctest.NewBtcKey(network, addrType)
			Expect(err).To(BeNil())

			By("Construct a RBF tx which sends money from pkAddr to toAddr")
			utxos, err := indexer.GetUTXOs(ctx, pkAddr)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e7), 5
			recipients := btc.SingleRecipient(toAddr.EncodeAddress(), amount)
			sizer := btc.NewSizeEstimator(utxos, btc.BaseSizeP2PKH, btc.SegwitSizeP2PKH)
			transaction, err := btc.BuildTransaction(network, feeRate, nil, utxos, sizer, recipients, pkAddr)
			Expect(err).To(BeNil())

			By("Sign and submit the fund tx")
			Expect(btc.SignTx(addrType, transaction, privKey, utxos)).Should(Succeed())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
			By(color.GreenString("RBF tx hash = %v", transaction.TxHash().String()))

			By("Mine a new block")
			Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

			By("Build a replacement tx with higher fee")
			queryTx, err := indexer.GetTx(ctx, transaction.TxHash().String())
			Expect(err).To(BeNil())
			vsize := (queryTx.Weight + 3) / 4
			prevFeeRate := int(queryTx.Fee) * 1000 / vsize
			replaceTx, err := btc.BuildRbfTransaction(network, prevFeeRate, queryTx.Fee, nil, utxos, sizer, recipients, pkAddr)
			Expect(err).To(BeNil())

			By("Sign and attempt to submit the replacement tx")
			Expect(btc.SignTx(addrType, replaceTx, privKey, utxos)).Should(Succeed())
			err = indexer.SubmitTx(ctx, replaceTx)
			Expect(errors.Is(err, btc.ErrTxInputsMissingOrSpent)).Should(BeTrue())
		})
	})

	Context("Single Anyone Can Pay flag", func() {
		It("should allow us to reuse the signature", func(ctx context.Context) {
			By("Initialize keys and addresses")
			addrType := waddrmgr.WitnessPubKey
			pk1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).To(BeNil())
			txid, err := btctest.Faucet(addr1.EncodeAddress())
			Expect(err).To(BeNil())
			Expect(btctest.WaitMined(ctx, indexer, btctest.WaitTx(txid.String()))).Should(Succeed())
			_, addr2, err := btctest.NewBtcKey(network, addrType)
			Expect(err).To(BeNil())

			By("Construct a transaction which sends money from addr1 to addr2")
			utxos, err := indexer.GetUTXOs(ctx, addr1)
			Expect(err).To(BeNil())
			amount, feeRate := int64(1e8-500), 2
			recipients := btc.SingleRecipient(addr2.EncodeAddress(), amount)
			sizer := btc.NewSizeEstimator(utxos, btc.BaseSizeP2WPKH, btc.SegwitSizeP2WPKH)
			transaction, err := btc.BuildTransaction(network, feeRate, nil, utxos, sizer, recipients, addr1)
			Expect(err).To(BeNil())

			By("Sign the transaction")
			pkScript, err := txscript.PayToAddrScript(addr1)
			Expect(err).Should(BeNil())
			fetcher, err := btc.InitFetcher(utxos, pkScript)
			sighashes := txscript.NewTxSigHashes(transaction, fetcher)
			err = btc.SignUtxos(waddrmgr.WitnessPubKey, transaction, 0, pk1, fetcher, sighashes, btc.WithSighashType(btc.SigHashSingleAnyoneCanPay))
			Expect(err).Should(BeNil())

			By("Add extra input and output to the transaction")
			extraUtxo := utxos[0]
			if utxos[0].String() == transaction.TxIn[0].PreviousOutPoint.String() {
				extraUtxo = utxos[1]
			}
			hash, err := chainhash.NewHashFromStr(extraUtxo.TxID)
			Expect(err).Should(BeNil())
			txIn := wire.NewTxIn(wire.NewOutPoint(hash, extraUtxo.Vout), nil, nil)
			transaction.TxIn = append([]*wire.TxIn{txIn}, transaction.TxIn...)
			toScript, err := txscript.PayToAddrScript(addr2)
			Expect(err).Should(BeNil())
			transaction.TxOut = append([]*wire.TxOut{wire.NewTxOut(extraUtxo.Amount, toScript)}, transaction.TxOut...)

			By("Sign and submit the transaction")
			sighashes = txscript.NewTxSigHashes(transaction, fetcher)
			err = btc.SignUtxos(waddrmgr.WitnessPubKey, transaction, 0, pk1, fetcher, sighashes)
			Expect(err).Should(BeNil())
			Expect(indexer.SubmitTx(ctx, transaction)).Should(Succeed())
		})
	})
})
