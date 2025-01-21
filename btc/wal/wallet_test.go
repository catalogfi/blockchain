package wallet_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/btc/btctest"
	wallet "github.com/catalogfi/blockchain/btc/wal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wallet", func() {
	Context("Operation on an HTLC", func() {
		It("should be able to initiate and redeem an HTLC", func() {
			By("Init keys and wallet")
			addrType := waddrmgr.WitnessPubKey
			key1, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).Should(BeNil())
			key2, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).Should(BeNil())
			wal1, err := wallet.NewWallet(network, addrType, key1, indexer, feeEstimator)
			Expect(err).Should(BeNil())
			wal2, err := wallet.NewWallet(network, addrType, key2, indexer, feeEstimator)
			Expect(err).Should(BeNil())

			By("Initiate an HTLC")
			secret := btctest.RandomSecret()
			secretHash := sha256.Sum256(secret)
			timelock := int64(144)
			key1PubBytes := schnorr.SerializePubKey(key1.PubKey())
			key2PubBytes := schnorr.SerializePubKey(key2.PubKey())
			htlc, err := btc.NewHTLC(key1PubBytes, key2PubBytes, secretHash[:], timelock, 1e7)
			Expect(err).Should(BeNil())
			_, err = wal1.Initiate(context.Background(), htlc)
			Expect(err).Should(BeNil())
			Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

			By("Redeem an HTLC")
			txid, err := wal2.Redeem(context.Background(), htlc, secret)
			Expect(err).Should(BeNil())
			log.Print("txid = ", txid)
		})

		It("should be able to refund an HTLC after it expires", func() {
			By("Init keys and wallet")
			addrType := waddrmgr.WitnessPubKey
			key1, _, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).Should(BeNil())
			key2, _, err := btctest.NewBtcKey(network, addrType)
			Expect(err).Should(BeNil())
			wal1, err := wallet.NewWallet(network, addrType, key1, indexer, feeEstimator)
			Expect(err).Should(BeNil())

			By("Initiate an HTLC")
			secret := btctest.RandomSecret()
			secretHash := sha256.Sum256(secret)
			timelock := int64(6)
			key1PubBytes := schnorr.SerializePubKey(key1.PubKey())
			key2PubBytes := schnorr.SerializePubKey(key2.PubKey())
			htlc, err := btc.NewHTLC(key1PubBytes, key2PubBytes, secretHash[:], timelock, 1e7)
			Expect(err).Should(BeNil())
			_, err = wal1.Initiate(context.Background(), htlc)
			Expect(err).Should(BeNil())

			By("Mine expiry no of blocks")
			for i := 0; i < int(timelock)-1; i++ {
				err = btctest.NewBlock()
				Expect(err).Should(BeNil())
			}
			Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

			By("Refund an HTLC")
			txid, err := wal1.Refund(context.Background(), htlc)
			Expect(err).Should(BeNil())
			log.Print("txid = ", txid)
		})

		FIt("should be able to instant refund an HTLC", func() {
			By("Init keys and wallet")
			addrType := waddrmgr.WitnessPubKey
			key1, addr1, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).Should(BeNil())
			key2, addr2, err := btctest.NewBtcAddrWithFunds(network, addrType, indexer)
			Expect(err).Should(BeNil())

			wal1, err := wallet.NewWallet(network, addrType, key1, indexer, feeEstimator)
			Expect(err).Should(BeNil())
			wal2, err := wallet.NewWallet(network, addrType, key2, indexer, feeEstimator)
			Expect(err).Should(BeNil())

			By("Initiate an HTLC")
			secret := btctest.RandomSecret()
			secretHash := sha256.Sum256(secret)
			timelock := int64(6)
			key1PubBytes := schnorr.SerializePubKey(key1.PubKey())
			key2PubBytes := schnorr.SerializePubKey(key2.PubKey())
			htlc, err := btc.NewHTLC(key1PubBytes, key2PubBytes, secretHash[:], timelock, 1e7)
			Expect(err).Should(BeNil())
			_, err = wal1.Initiate(context.Background(), htlc)
			Expect(err).Should(BeNil())
			Expect(btctest.NewBlockWaitMined(indexer)).Should(Succeed())

			By("Use the instant refund leaf to refund the HTLC")
			addr, err := htlc.Address(network)
			Expect(err).Should(BeNil())
			htlcUtxos, err := indexer.GetUTXOs(context.Background(), addr)
			Expect(err).Should(BeNil())
			utxos1, err := indexer.GetUTXOs(context.Background(), wal1.Address())
			Expect(err).Should(BeNil())

			// Construct the transaction
			sizer := btc.NewSizeEstimator(utxos1, btc.BaseSizeP2WPKH, btc.SegwitSizeP2WPKH)
			sizer.AddUtxos(htlcUtxos, btc.BaseSizeP2WPKH, btc.SegwitSizeP2WPKH)
			recipients := []btc.Recipient{
				{
					To:     addr1.EncodeAddress(),
					Amount: 1e7,
				},
			}
			tx1, err := btc.BuildTransaction(network, 100, htlcUtxos, utxos1, sizer, recipients, wal1.Address())
			Expect(err).Should(BeNil())

			// Sign the tx1
			pkScript, err := txscript.PayToAddrScript(addr1)
			Expect(err).Should(BeNil())

			fetcher, err := btc.InitFetcher(utxos1, pkScript)
			Expect(err).Should(BeNil())
			hash, err := chainhash.NewHashFromStr(htlcUtxos[0].TxID)
			Expect(err).Should(BeNil())

			script, err := htlc.P2trScript()
			Expect(err).Should(BeNil())

			fetcher.AddPrevOut(wire.OutPoint{
				Hash:  *hash,
				Index: htlcUtxos[0].Vout,
			}, wire.NewTxOut(htlcUtxos[0].Amount, script))

			sigHashes := txscript.NewTxSigHashes(tx1, fetcher)
			leaf, ctrBlk := htlc.InstantRefundLeaf()
			ctrBlkBytes, err := ctrBlk.ToBytes()
			if err != nil {
				Expect(err).Should(BeNil())
			}
			for i := range tx1.TxIn {
				if i == 1 {
					output := fetcher.FetchPrevOutput(tx1.TxIn[i].PreviousOutPoint)
					sig, err := txscript.RawTxInWitnessSignature(tx1, sigHashes, i, output.Value, output.PkScript, txscript.SigHashAll, key1)
					Expect(err).Should(BeNil())
					tx1.TxIn[i].Witness = wire.TxWitness{sig, key1.PubKey().SerializeCompressed()}
				} else {

					for i, utxo := range tx1.TxIn {
						out := fetcher.FetchPrevOutput(utxo.PreviousOutPoint)
						sig1, err := txscript.RawTxInTapscriptSignature(tx1, sigHashes, i, out.Value, out.PkScript, leaf, btc.SigHashSingleAnyoneCanPay, key1)
						if err != nil {
							Expect(err).Should(BeNil())
						}
						sig2, err := txscript.RawTxInTapscriptSignature(tx1, sigHashes, i, out.Value, out.PkScript, leaf, txscript.SigHashAll, key2)
						if err != nil {
							Expect(err).Should(BeNil())
						}
						tx1.TxIn[i].Witness = append(tx1.TxIn[i].Witness, sig2, sig1, leaf.Script, ctrBlkBytes)
					}
				}
			}

			// Submit tx1
			log.Print(tx1.TxHash().String())
			Expect(indexer.SubmitTx(context.Background(), tx1)).Should(Succeed())
			time.Sleep(10 * time.Second)

			// Construct a replacement transaction
			utxos2, err := indexer.GetUTXOs(context.Background(), addr2)
			sizer2 := btc.NewSizeEstimator(utxos2, btc.BaseSizeP2WPKH, btc.SegwitSizeP2WPKH)
			sizer2.AddUtxos(htlcUtxos, btc.BaseSizeP2WPKH, btc.SegwitSizeP2WPKH)
			tx2, err := btc.BuildTransaction(network, 200, htlcUtxos, utxos2, sizer2, recipients, wal2.Address())
			Expect(err).Should(BeNil())

			pkScript2, err := txscript.PayToAddrScript(addr2)
			Expect(err).Should(BeNil())
			fetcher2, err := btc.InitFetcher(utxos2, pkScript2)
			Expect(err).Should(BeNil())
			fetcher2.AddPrevOut(wire.OutPoint{
				Hash:  *hash,
				Index: htlcUtxos[0].Vout,
			}, wire.NewTxOut(htlcUtxos[0].Amount, script))
			sigHashes = txscript.NewTxSigHashes(tx2, fetcher2)
			for i, utxo := range tx2.TxIn {
				if i == 1 {
					output := fetcher2.FetchPrevOutput(tx2.TxIn[i].PreviousOutPoint)
					sig, err := txscript.RawTxInWitnessSignature(tx2, sigHashes, i, output.Value, output.PkScript, txscript.SigHashAll, key2)
					Expect(err).Should(BeNil())
					tx2.TxIn[i].Witness = wire.TxWitness{sig, key2.PubKey().SerializeCompressed()}
				} else {
					out := fetcher2.FetchPrevOutput(utxo.PreviousOutPoint)
					sig1, err := txscript.RawTxInTapscriptSignature(tx2, sigHashes, i, out.Value, out.PkScript, leaf, btc.SigHashSingleAnyoneCanPay, key1)
					if err != nil {
						Expect(err).Should(BeNil())
					}
					sig2, err := txscript.RawTxInTapscriptSignature(tx2, sigHashes, i, out.Value, out.PkScript, leaf, txscript.SigHashAll, key2)
					if err != nil {
						Expect(err).Should(BeNil())
					}
					tx2.TxIn[i].Witness = append(tx2.TxIn[i].Witness, sig2, sig1, leaf.Script, ctrBlkBytes)

					// tx2.TxIn[i].Witness = tx1.TxIn[i].Witness
					log.Printf("have %v witness", len(tx2.TxIn[i].Witness))
				}
			}
			raw, err := btc.TxRawBytes(tx2)
			Expect(err).Should(BeNil())

			log.Print(hex.EncodeToString(raw))

			log.Print(tx2.TxHash().String())
			Expect(indexer.SubmitTx(context.Background(), tx2)).Should(Succeed())
		})
	})
})
