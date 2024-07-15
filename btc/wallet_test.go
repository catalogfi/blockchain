package btc_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/localnet"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("Wallets", Ordered, func() {

	chainParams := chaincfg.RegressionNetParams
	logger, err := zap.NewDevelopment()
	Expect(err).To(BeNil())

	indexer := btc.NewElectrsIndexerClient(logger, os.Getenv("BTC_REGNET_INDEXER"), time.Millisecond*500)

	privateKey, err := btcec.NewPrivateKey()
	Expect(err).To(BeNil())

	fixedFeeEstimator := btc.NewFixFeeEstimator(10)
	feeLevel := btc.HighFee
	wallet, err := btc.NewSimpleWallet(privateKey, &chainParams, indexer, fixedFeeEstimator, feeLevel)
	Expect(err).To(BeNil())

	BeforeAll(func() {
		_, err := localnet.FundBitcoin(wallet.Address().EncodeAddress(), indexer)
		Expect(err).To(BeNil())
	})

	It("should be able to send funds", func() {
		req := []btc.SendRequest{
			{
				Amount: 100000,
				To:     wallet.Address(),
			},
		}

		txid, err := wallet.Send(context.Background(), req, nil, nil)
		Expect(err).To(BeNil())

		tx, ok, err := wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())
		Expect(ok).Should(BeTrue())
		Expect(tx).ShouldNot(BeNil())
		// to address
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
		// change address
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
	})

	It("should be able to send funds to multiple addresses", func() {
		pk, err := btcec.NewPrivateKey()
		Expect(err).To(BeNil())
		bobWallet, err := btc.NewSimpleWallet(pk, &chainParams, indexer, fixedFeeEstimator, feeLevel)
		Expect(err).To(BeNil())
		bobAddr := bobWallet.Address()

		req := []btc.SendRequest{
			{
				Amount: 100000,
				To:     wallet.Address(),
			},
			{
				Amount: 100000,
				To:     bobAddr,
			},
		}

		txid, err := wallet.Send(context.Background(), req, nil, nil)
		Expect(err).To(BeNil())

		tx, _, err := wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())
		Expect(tx).ShouldNot(BeNil())

		// first vout address
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
		// second vout address
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(bobAddr.EncodeAddress()))
		// change address
		Expect(tx.VOUTs[2].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))

	})

	It("should not be able to send funds if the amount + fee is greater than the balance", func() {

		balance, err := getBalance(indexer, wallet.Address())
		Expect(err).To(BeNil())

		Expect(balance).Should(BeNumerically(">", 0))
		req := []btc.SendRequest{
			{
				Amount: 100000000000,
				To:     wallet.Address(),
			},
		}
		_, err = wallet.Send(context.Background(), req, nil, nil)

		Expect(err).ShouldNot(BeNil())
		req = []btc.SendRequest{
			{
				Amount: 100000000 - 1000,
				To:     wallet.Address(),
			},
		}
		_, err = wallet.Send(context.Background(), req, nil, nil)
		Expect(err).ShouldNot(BeNil())
	})

	It("should be able to send funds without change if change is less than the dust limit", func() {
		By("Create new Bob Wallet")
		pk, err := btcec.NewPrivateKey()
		Expect(err).To(BeNil())
		bobWallet, err := btc.NewSimpleWallet(pk, &chainParams, indexer, fixedFeeEstimator, feeLevel)

		By("Send funds to Bob Wallet")
		bobBalance := int64(100000)
		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: bobBalance,
				To:     bobWallet.Address(),
			},
		}, nil, nil)
		Expect(err).To(BeNil())

		By("Build the send request")
		req := []btc.SendRequest{
			{
				// 1090 is the fee set by the p2wpkh wallet for regtest for single input and single output
				Amount: bobBalance - 1090,
				// self transfer
				To: bobWallet.Address(),
			},
		}
		txid, err := bobWallet.Send(context.Background(), req, nil, nil)
		Expect(err).To(BeNil())

		tx, _, err := wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())
		By("Tx should not contain a change output")
		Expect(tx).ShouldNot(BeNil())
		// there should be only one output
		Expect(tx.VOUTs).Should(HaveLen(1))

		bobBalance, err = getBalance(indexer, bobWallet.Address())
		Expect(err).To(BeNil())

		By("Send funds in such a way that change is less than the dust limit")
		req = []btc.SendRequest{
			{
				// 1090 is the fee set by the p2wpkh wallet for regtest for single input and single output
				// we subtract 400 from the amount to make the change less than the dust limit
				// such that the change is not included in the transaction but used for fee
				Amount: bobBalance - 1090 - 400,
				To:     wallet.Address(),
			},
		}

		txid, err = bobWallet.Send(context.Background(), req, nil, nil)
		Expect(err).To(BeNil())

		By("Tx should contain only one output without change")
		tx, _, err = wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())
		Expect(tx).ShouldNot(BeNil())
		Expect(tx.VOUTs).Should(HaveLen(1))

		By("Excess amount should be used for fee")
		// Fee paid should be 1090 + 400 = 1490
		Expect(tx.Fee).Should(Equal(int64(1490)))
		Expect(tx.VOUTs[0].Value).Should(Equal(int(bobBalance - (1090 + 400))))
	})

	It("should be able to spend funds from p2wpkh script", func() {
		address := wallet.Address()
		pkScript, err := txscript.PayToAddrScript(address)
		Expect(err).To(BeNil())

		txId, err := wallet.Send(context.Background(), nil, []btc.SpendRequest{
			{
				Witness: [][]byte{
					btc.AddSignatureSegwitOp,
					btc.AddPubkeyCompressedOp,
				},
				Script:        pkScript,
				ScriptAddress: address,
				HashType:      txscript.SigHashAll,
			},
		}, nil)
		Expect(err).To(BeNil())

		tx, _, err := wallet.Status(context.Background(), txId)
		Expect(err).To(BeNil())
		Expect(tx).ShouldNot(BeNil())

		Expect(tx.VOUTs).Should(HaveLen(1))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(address.EncodeAddress()))

	})

	It("should be able to spend funds from a simple p2wsh script", func() {

		script, scriptAddr, err := additionScript(chainParams)
		Expect(err).To(BeNil())

		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: 100000,
				To:     scriptAddr,
			},
		}, nil, nil)
		Expect(err).To(BeNil())

		// spend the script
		txId, err := wallet.Send(context.Background(), nil,
			[]btc.SpendRequest{
				{
					Witness: [][]byte{
						{0x1},
						{0x1},
						script,
					},
					Script:        script,
					ScriptAddress: scriptAddr,
					HashType:      txscript.SigHashAll,
				},
			}, nil)
		Expect(err).To(BeNil())
		Expect(txId).ShouldNot(BeEmpty())

	})

	It("should be able to spend funds from a simple p2tr script", func() {
		script, scriptAddr, cb, err := additionTapscript(chainParams)
		Expect(err).To(BeNil())

		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: 100000,
				To:     scriptAddr,
			},
		}, nil, nil)
		Expect(err).To(BeNil())

		// spend the script
		txId, err := wallet.Send(context.Background(), nil, []btc.SpendRequest{
			{
				Witness: [][]byte{
					{0x1},
					{0x1},
					script,
					cb,
				},
				Script:        script,
				Leaf:          txscript.NewTapLeaf(0xc0, script),
				ScriptAddress: scriptAddr,
			},
		}, nil,
		)
		Expect(err).To(BeNil())
		Expect(txId).ShouldNot(BeEmpty())
	})

	It("should be able to spend funds from signature-check script p2tr", func() {
		script, scriptAddr, cb, err := sigCheckTapScript(chainParams, schnorr.SerializePubKey(privateKey.PubKey()))
		Expect(err).To(BeNil())

		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: 100000,
				To:     scriptAddr,
			},
		}, nil, nil)
		Expect(err).To(BeNil())

		txId, err := wallet.Send(context.Background(), nil, []btc.SpendRequest{
			{
				// we don't pass the script here as it is not needed for taproot
				// instead we pass the leaf
				Witness: [][]byte{
					btc.AddSignatureSchnorrOp,
					script,
					cb,
				},
				Leaf:          txscript.NewTapLeaf(0xc0, script),
				ScriptAddress: scriptAddr,
			},
		}, nil)
		Expect(err).To(BeNil())
		Expect(txId).ShouldNot(BeEmpty())

	})

	It("should be able to spend funds from signature-check script p2wsh", func() {
		script, scriptAddr, err := sigCheckScript(chainParams, privateKey)
		Expect(err).To(BeNil())
		scriptBalance := int64(100000)
		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: scriptBalance,
				To:     scriptAddr,
			},
		}, nil, nil)
		Expect(err).To(BeNil())

		// spend the script
		txId, err := wallet.Send(context.Background(), nil, []btc.SpendRequest{
			{
				Witness: [][]byte{
					btc.AddSignatureSegwitOp,
					script,
				},
				Script:        script,
				ScriptAddress: scriptAddr,
				HashType:      txscript.SigHashAll,
			},
		}, nil)
		Expect(err).To(BeNil())
		Expect(txId).ShouldNot(BeEmpty())

		scriptBalance, err = getBalance(indexer, scriptAddr)
		Expect(err).To(BeNil())
		Expect(scriptBalance).Should(Equal(int64(0)))

	})

	It("should be able to spend funds from a two p2wsh scripts", func() {
		p2wshScript1, p2wshAddr1, err := additionScript(chainParams)
		Expect(err).To(BeNil())

		p2wshScript2, p2wshAddr2, err := additionScript(chainParams)
		Expect(err).To(BeNil())

		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: 100000,
				To:     p2wshAddr1,
			},
			{
				Amount: 100000,
				To:     p2wshAddr2,
			},
		}, nil, nil)

		txId, err := wallet.Send(context.Background(), nil, []btc.SpendRequest{
			{
				Witness: [][]byte{
					{0x1},
					{0x1},
					p2wshScript1,
				},
				Script:        p2wshScript1,
				ScriptAddress: p2wshAddr1,
				HashType:      txscript.SigHashAll,
			},
			{
				Witness: [][]byte{
					{0x1},
					{0x1},
					p2wshScript2,
				},
				Script:        p2wshScript2,
				ScriptAddress: p2wshAddr2,
				HashType:      txscript.SigHashAll,
			},
		}, nil)

		Expect(err).To(BeNil())
		Expect(txId).ShouldNot(BeEmpty())
	})

	It("should be able to spend funds from different (p2wsh and p2tr) scripts", func() {
		p2wshAdditionScript, p2wshAddr, err := additionScript(chainParams)
		Expect(err).To(BeNil())

		p2trAdditionScript, p2trAddr, cb, err := additionTapscript(chainParams)
		Expect(err).To(BeNil())

		p2trSigCheckScript, p2trSigCheckScriptAddr, sigCheckCb, err := sigCheckTapScript(chainParams, schnorr.SerializePubKey(privateKey.PubKey()))
		Expect(err).To(BeNil())

		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: 100000,
				To:     p2wshAddr,
			},
			{
				Amount: 100000,
				To:     p2trAddr,
			},
			{
				Amount: 100000,
				To:     p2trSigCheckScriptAddr,
			},
		}, nil, nil)

		By("Spend p2wsh and p2tr scripts")
		txId, err := wallet.Send(context.Background(), nil, []btc.SpendRequest{
			{
				Witness: [][]byte{
					{0x1},
					{0x1},
					p2wshAdditionScript,
				},
				Script:        p2wshAdditionScript,
				ScriptAddress: p2wshAddr,
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
				ScriptAddress: p2trAddr,
			},
			{
				Witness: [][]byte{
					btc.AddSignatureSchnorrOp,
					p2trSigCheckScript,
					sigCheckCb,
				},
				Leaf:          txscript.NewTapLeaf(0xc0, p2trSigCheckScript),
				ScriptAddress: p2trSigCheckScriptAddr,
				HashType:      txscript.SigHashAll,
			},
		}, nil)

		Expect(err).To(BeNil())
		Expect(txId).ShouldNot(BeEmpty())

		By("Balances of both scripts should be zero")
		balance, err := getBalance(indexer, p2wshAddr)
		Expect(err).To(BeNil())
		Expect(balance).Should(Equal(int64(0)))

		balance, err = getBalance(indexer, p2trAddr)
		Expect(err).To(BeNil())
		Expect(balance).Should(Equal(int64(0)))

		By("Validate the tx")
		tx, _, err := wallet.Status(context.Background(), txId)
		Expect(err).To(BeNil())
		Expect(tx).ShouldNot(BeNil())
		Expect(tx.VOUTs).Should(HaveLen(1))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
	})

	It("should not be able to spend if the script has no balance", func() {
		additionScript, additionAddr, err := additionScript(chainParams)
		Expect(err).To(BeNil())

		_, err = wallet.Send(context.Background(), nil, []btc.SpendRequest{
			{
				Witness: [][]byte{
					{0x1},
					{0x1},
					additionScript,
				},
				Script:        additionScript,
				ScriptAddress: additionAddr,
				HashType:      txscript.SigHashAll,
			},
		}, nil)
		Expect(err).ShouldNot(BeNil())
	})

	It("should not be able to spend with invalid Inputs", func() {
		_, err := wallet.Send(context.Background(), nil, []btc.SpendRequest{
			{
				Witness: [][]byte{
					{0x1},
					{0x1},
					{0x1},
				},
				Script:        []byte{0x1},
				ScriptAddress: wallet.Address(),
				HashType:      txscript.SigHashAll,
			},
		}, nil)
		Expect(err).ShouldNot(BeNil())

		By("Lets give proper witness but bad script")
		_, err = wallet.Send(context.Background(), nil, []btc.SpendRequest{
			{
				Witness: [][]byte{
					btc.AddSignatureSegwitOp,
					btc.AddPubkeyCompressedOp,
				},
				Script:        []byte{0x1},
				ScriptAddress: wallet.Address(),
				HashType:      txscript.SigHashAll,
			},
		}, nil)
		Expect(err).ShouldNot(BeNil())
	})

	It("should be able to fetch the status and tx details", func() {
		req := []btc.SendRequest{
			{
				Amount: 100000,
				To:     wallet.Address(),
			},
		}

		txid, err := wallet.Send(context.Background(), req, nil, nil)
		Expect(err).To(BeNil())

		tx, ok, err := wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())
		Expect(ok).Should(BeTrue())
		Expect(tx).ShouldNot(BeNil())

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		tx, ok, err = wallet.Status(ctx, "ffffff"+txid[6:])
		Expect(err).To(BeNil())
		Expect(ok).Should(BeFalse())
		Expect(tx.TxID).Should(BeEmpty())
	})

	It("should be able to spend and send at the same time", func() {
		amount := int64(100000)
		// p2wpkh script
		script, err := txscript.PayToAddrScript(wallet.Address())
		Expect(err).To(BeNil())

		By("Let's create a random wallet")
		pk, err := btcec.NewPrivateKey()
		Expect(err).To(BeNil())
		randomWallet, err := btc.NewSimpleWallet(pk, &chainParams, indexer, fixedFeeEstimator, feeLevel)
		Expect(err).To(BeNil())

		By("Send funds to the random address by spending the script")
		txId, err := wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: amount,
				To:     randomWallet.Address(),
			},
		}, []btc.SpendRequest{
			{
				Witness: [][]byte{
					btc.AddSignatureSegwitOp,
					btc.AddPubkeyCompressedOp,
				},
				Script:        script,
				ScriptAddress: wallet.Address(),
				HashType:      txscript.SigHashAll,
			},
		}, nil)
		Expect(err).To(BeNil())
		Expect(txId).ShouldNot(BeEmpty())

		By("Check the balance of the random wallet")
		balance, err := getBalance(indexer, randomWallet.Address())
		Expect(err).To(BeNil())
		Expect(balance).Should(Equal(int64(amount)))

		tx, _, err := wallet.Status(context.Background(), txId)
		Expect(err).To(BeNil())
		Expect(tx).ShouldNot(BeNil())
		// one for the random wallet and one for the change
		Expect(tx.VOUTs).Should(HaveLen(2))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(randomWallet.Address().EncodeAddress()))
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
	})

	It("should be able to spend multiple scripts and send to multiple parties", func() {
		amount := int64(100000)

		p2wshSigCheckScript, p2wshSigCheckScriptAddr, err := sigCheckScript(chainParams, privateKey)

		p2wshAdditionScript, p2wshScriptAddr, err := additionScript(chainParams)
		Expect(err).To(BeNil())

		p2trAdditionScript, p2trScriptAddr, cb, err := additionTapscript(chainParams)
		Expect(err).To(BeNil())

		checkSigScript, checkSigScriptAddr, checkSigScriptCb, err := sigCheckTapScript(chainParams, schnorr.SerializePubKey(privateKey.PubKey()))
		Expect(err).To(BeNil())

		By("Fund the scripts")
		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: amount,
				To:     p2wshScriptAddr,
			},
			{
				Amount: amount,
				To:     p2wshSigCheckScriptAddr,
			},
			{
				Amount: amount,
				To:     p2trScriptAddr,
			},
			{
				Amount: amount,
				To:     checkSigScriptAddr,
			},
		}, nil, nil)
		Expect(err).To(BeNil())

		By("Let's create Bob and Dave wallets")
		pk, err := btcec.NewPrivateKey()
		Expect(err).To(BeNil())
		bobWallet, err := btc.NewSimpleWallet(pk, &chainParams, indexer, fixedFeeEstimator, feeLevel)
		Expect(err).To(BeNil())
		pk, err = btcec.NewPrivateKey()
		Expect(err).To(BeNil())
		daveWallet, err := btc.NewSimpleWallet(pk, &chainParams, indexer, fixedFeeEstimator, feeLevel)
		Expect(err).To(BeNil())

		By("Send funds to Bob and Dave by spending the scripts")
		txId, err := wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: amount,
				To:     bobWallet.Address(),
			},
			{
				Amount: amount,
				To:     daveWallet.Address(),
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
		Expect(txId).ShouldNot(BeEmpty())

		By("The tx should have 3 outputs")
		tx, _, err := wallet.Status(context.Background(), txId)
		Expect(err).To(BeNil())
		Expect(tx).ShouldNot(BeNil())
		Expect(tx.VOUTs).Should(HaveLen(3))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(bobWallet.Address().EncodeAddress()))
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(daveWallet.Address().EncodeAddress()))
		Expect(tx.VOUTs[2].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))

		//Validate whether dave and bob received the right amount
		Expect(tx.VOUTs[0].Value).Should(Equal(int(amount)))
		Expect(tx.VOUTs[1].Value).Should(Equal(int(amount)))

	})

	It("should be able to send SACPs", func() {
		sacp, err := generateSACP(wallet, chainParams, privateKey, wallet.Address(), 10000, 1000)
		Expect(err).To(BeNil())

		txid, err := wallet.Send(context.Background(), nil, nil, [][]byte{sacp})
		Expect(err).To(BeNil())
		Expect(txid).ShouldNot(BeEmpty())

		tx, _, err := wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())
		Expect(tx).ShouldNot(BeNil())

		Expect(tx.VOUTs).Should(HaveLen(2))
		// Actual recipient in generateSACP is the wallet itself
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
		// Change address from wallet.Send
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
	})

	It("should be able to send multiple SACPs", func() {
		err = localnet.MineBitcoinBlocks(1, indexer)
		Expect(err).To(BeNil())

		bobPk, err := btcec.NewPrivateKey()
		Expect(err).To(BeNil())
		bobWallet, err := btc.NewSimpleWallet(bobPk, &chainParams, indexer, fixedFeeEstimator, feeLevel)
		Expect(err).To(BeNil())

		By("Send funds to Bob")
		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: 10000000,
				To:     bobWallet.Address(),
			},
		}, nil, nil)

		sacp1, err := generateSACP(bobWallet, chainParams, bobPk, wallet.Address(), 5000, 0)
		Expect(err).To(BeNil())

		sacp2, err := generateSACP(wallet, chainParams, privateKey, bobWallet.Address(), 10000, 1000)
		Expect(err).To(BeNil())

		txid, err := wallet.Send(context.Background(), nil, nil, [][]byte{sacp1, sacp2})
		Expect(err).To(BeNil())
		Expect(txid).ShouldNot(BeEmpty())

		tx, _, err := wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())

		// One is bob's SACP recipient and one is from Alice (wallet) and other is change
		Expect(tx.VOUTs).Should(HaveLen(3))

		// Bob's SACP recipient
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
		// Alice's SACP recipient
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(bobWallet.Address().EncodeAddress()))
		// Change address
		Expect(tx.VOUTs[2].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
	})

	It("should be able to mix SACPs with send requests", func() {

		bobPk, err := btcec.NewPrivateKey()
		Expect(err).To(BeNil())
		bobWallet, err := btc.NewSimpleWallet(bobPk, &chainParams, indexer, fixedFeeEstimator, feeLevel)
		Expect(err).To(BeNil())

		By("Send funds to Bob")
		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: 10000000,
				To:     bobWallet.Address(),
			},
		}, nil, nil)

		sacp1, err := generateSACP(bobWallet, chainParams, bobPk, wallet.Address(), 10000, 100)
		Expect(err).To(BeNil())

		txid, err := wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: 10000000,
				To:     bobWallet.Address(),
			},
		}, nil, [][]byte{sacp1})
		Expect(err).To(BeNil())
		Expect(txid).ShouldNot(BeEmpty())

		tx, _, err := wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())

		// One is bob's SACP recipient and one is from Alice (wallet) to bob and other is change
		Expect(tx.VOUTs).Should(HaveLen(3))

		// Bob's SACP recipient
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
		// Alice's SACP recipient
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(bobWallet.Address().EncodeAddress()))
		// Change address
		Expect(tx.VOUTs[2].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
	})

	It("should be able to mix SACPs with spend requests", func() {

		script, scriptAddr, err := sigCheckScript(chainParams, privateKey)
		Expect(err).To(BeNil())

		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: 100000,
				To:     scriptAddr,
			},
		}, nil, nil)
		Expect(err).To(BeNil())

		randAddr, err := randomP2wpkhAddress(chainParams)
		Expect(err).To(BeNil())

		sacp, err := generateSACP(wallet, chainParams, privateKey, randAddr, 10000, 100)
		Expect(err).To(BeNil())

		txid, err := wallet.Send(context.Background(), nil, []btc.SpendRequest{
			{
				Witness: [][]byte{
					btc.AddSignatureSegwitOp,
					script,
				},
				Script:        script,
				ScriptAddress: scriptAddr,
				HashType:      txscript.SigHashAll,
			},
		}, [][]byte{sacp})
		Expect(err).To(BeNil())
		Expect(txid).ShouldNot(BeEmpty())

		tx, _, err := wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())
		Expect(tx).ShouldNot(BeNil())

		// Three outputs, one for the script, one for the SACP recipient
		Expect(tx.VOUTs).Should(HaveLen(2))

		// SACP recipient is the wallet itself
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(randAddr.EncodeAddress()))
		// Script recipient is the wallet itself
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
	})

	It("should be able to mix SACPs with spend requests and send requests", func() {

		script, scriptAddr, err := sigCheckScript(chainParams, privateKey)
		Expect(err).To(BeNil())

		sigCheckP2trScript, sigCheckP2trScriptAddr, cb, err := sigCheckTapScript(chainParams, schnorr.SerializePubKey(privateKey.PubKey()))
		Expect(err).To(BeNil())
		sendAmount := int64(100000)
		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: sendAmount,
				To:     scriptAddr,
			},
			{
				Amount: sendAmount,
				To:     sigCheckP2trScriptAddr,
			},
		}, nil, nil)
		Expect(err).To(BeNil())

		randAddr, err := randomP2wpkhAddress(chainParams)
		Expect(err).To(BeNil())

		sacp1, err := generateSACP(wallet, chainParams, privateKey, randAddr, sendAmount, 0)
		Expect(err).To(BeNil())

		bobAddr, err := randomP2wpkhAddress(chainParams)
		Expect(err).To(BeNil())

		sacp2, err := generateSACP(wallet, chainParams, privateKey, bobAddr, sendAmount, 1000)
		Expect(err).To(BeNil())
		txid, err := wallet.Send(context.Background(), []btc.SendRequest{
			{
				To:     bobAddr,
				Amount: sendAmount,
			},
		}, []btc.SpendRequest{
			{
				Witness: [][]byte{
					btc.AddSignatureSegwitOp,
					script,
				},
				Script:        script,
				ScriptAddress: scriptAddr,
				HashType:      txscript.SigHashAll,
			},
			{
				Witness: [][]byte{
					btc.AddSignatureSchnorrOp,
					sigCheckP2trScript,
					cb,
				},
				Leaf:          txscript.NewTapLeaf(0xc0, sigCheckP2trScript),
				ScriptAddress: sigCheckP2trScriptAddr,
				HashType:      txscript.SigHashAll,
			},
		},
			[][]byte{sacp1, sacp2},
		)
		Expect(err).To(BeNil())
		Expect(txid).ShouldNot(BeEmpty())

		tx, _, err := wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())
		Expect(tx).ShouldNot(BeNil())

		// Three outputs, one for the script, one for the SACP recipient
		Expect(tx.VOUTs).Should(HaveLen(4))

		// SACP recipient is the random address
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(randAddr.EncodeAddress()))
		Expect(tx.VOUTs[0].Value).Should(Equal(int(sendAmount)))
		// SACP2 recipient is the bob address
		Expect(tx.VOUTs[1].ScriptPubKeyAddress).Should(Equal(bobAddr.EncodeAddress()))
		Expect(tx.VOUTs[1].Value).Should(Equal(int(sendAmount - 1000))) // 1000 is the fee for SACP2
		// Script recipient is the wallet itself
		Expect(tx.VOUTs[2].ScriptPubKeyAddress).Should(Equal(bobAddr.EncodeAddress()))
		Expect(tx.VOUTs[2].Value).Should(Equal(int(sendAmount)))
		// Change address
		Expect(tx.VOUTs[3].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
	})

	It("should be able to adjust fees based on SACP's fee", func() {
		// Inside the generateSACP function, we are setting the fee to 1000
		sacp, err := generateSACP(wallet, chainParams, privateKey, wallet.Address(), 1000, 1)
		Expect(err).To(BeNil())

		txid, err := wallet.Send(context.Background(), nil, nil, [][]byte{sacp})
		Expect(err).To(BeNil())
		Expect(txid).ShouldNot(BeEmpty())

		tx, _, err := wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())
		Expect(tx).ShouldNot(BeNil())

		txHex, err := indexer.GetTxHex(context.Background(), txid)
		Expect(err).To(BeNil())

		txBytes, err := hex.DecodeString(txHex)
		Expect(err).To(BeNil())

		transaction, err := btcutil.NewTxFromBytes(txBytes)
		Expect(err).To(BeNil())

		feePaid, err := btc.EstimateSegwitFee(transaction.MsgTx(), fixedFeeEstimator, feeLevel)
		Expect(err).To(BeNil())

		Expect(tx.Fee).Should(BeEquivalentTo(feePaid))
	})

	It("should be able to generate an SACP", func() {
		script, scriptAddr, err := sigCheckScript(chainParams, privateKey)
		Expect(err).To(BeNil())

		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: 100000,
				To:     scriptAddr,
			},
		}, nil, nil)
		Expect(err).To(BeNil())

		txBytes, err := wallet.GenerateSACP(context.Background(), btc.SpendRequest{
			Witness: [][]byte{
				btc.AddSignatureSegwitOp,
				script,
			},
			Script:        script,
			ScriptAddress: scriptAddr,
			HashType:      btc.SigHashSingleAnyoneCanPay,
		}, nil)
		Expect(err).To(BeNil())
		Expect(txBytes).ShouldNot(BeEmpty())

		By("Generated SACP should be submittable")

		txid, err := wallet.Send(context.Background(), nil, nil, [][]byte{txBytes})
		Expect(err).To(BeNil())
		Expect(txid).ShouldNot(BeEmpty())

		tx, _, err := wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())
		Expect(tx).ShouldNot(BeNil())

		Expect(tx.VOUTs).Should(HaveLen(1))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
		Expect(tx.VOUTs[0].Value).Should(Equal(int(100000 - tx.Fee)))

	})

	It("should be able to generate p2tr SACP", func() {
		script, scriptAddr, cb, err := sigCheckTapScript(chainParams, schnorr.SerializePubKey(privateKey.PubKey()))
		Expect(err).To(BeNil())

		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: 100000,
				To:     scriptAddr,
			},
		}, nil, nil)
		Expect(err).To(BeNil())

		txBytes, err := wallet.GenerateSACP(context.Background(), btc.SpendRequest{
			Witness: [][]byte{
				btc.AddSignatureSchnorrOp,
				script,
				cb,
			},
			Leaf:          txscript.NewTapLeaf(0xc0, script),
			ScriptAddress: scriptAddr,
		}, nil)
		Expect(err).To(BeNil())
		Expect(txBytes).ShouldNot(BeEmpty())

		txid, err := wallet.Send(context.Background(), nil, nil, [][]byte{txBytes})
		Expect(err).To(BeNil())
		Expect(txid).ShouldNot(BeEmpty())

		tx, _, err := wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())
		Expect(tx).ShouldNot(BeNil())

		Expect(tx.VOUTs).Should(HaveLen(1))
		Expect(tx.VOUTs[0].ScriptPubKeyAddress).Should(Equal(wallet.Address().EncodeAddress()))
		Expect(tx.VOUTs[0].Value).Should(Equal(int(100000 - tx.Fee)))
	})

	It("should be able to generate SACP signature", func() {

		script, scriptAddr, cb, err := sigCheckTapScript(chainParams, schnorr.SerializePubKey(privateKey.PubKey()))
		Expect(err).To(BeNil())

		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: 100000,
				To:     scriptAddr,
			},
		}, nil, nil)
		Expect(err).To(BeNil())

		txBytes, err := wallet.GenerateSACP(context.Background(), btc.SpendRequest{
			Witness: [][]byte{
				btc.AddSignatureSchnorrOp,
				script,
				cb,
			},
			Leaf:          txscript.NewTapLeaf(0xc0, script),
			ScriptAddress: scriptAddr,
			HashType:      btc.SigHashSingleAnyoneCanPay,
		}, nil)
		Expect(err).To(BeNil())
		Expect(txBytes).ShouldNot(BeEmpty())

		btcTx, err := btcutil.NewTxFromBytes(txBytes)
		Expect(err).Should(BeNil())

		tx := btcTx.MsgTx()

		witness, err := wallet.SignSACPTx(tx, 0, 100000, txscript.NewTapLeaf(0xc0, script), scriptAddr, [][]byte{
			btc.AddSignatureSchnorrOp,
			script,
			cb,
		})
		Expect(err).To(BeNil())
		Expect(bytes.Equal(tx.TxIn[0].Witness[0], witness[0])).Should(BeTrue())

	})

})

//------------------------------ Helper functions for tests --------------------------------

func randomP2wpkhAddress(chainParams chaincfg.Params) (btcutil.Address, error) {
	pk, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, err
	}
	return btc.PublicKeyAddress(&chainParams, waddrmgr.WitnessPubKey, pk.PubKey())
}

// generateSACP generates an SACP transaction with single input and output.
// It uses wallet's address as recipient and a script as the sender.
func generateSACP(wallet btc.Wallet, chainParams chaincfg.Params, privKey *secp256k1.PrivateKey, recipient btcutil.Address, amount int64, fee int) ([]byte, error) {
	script, scriptAddr, err := sigCheckScript(chainParams, privKey)
	if err != nil {
		return nil, err
	}

	// fund the script
	txid, err := wallet.Send(context.Background(), []btc.SendRequest{
		{
			Amount: amount,
			To:     scriptAddr,
		},
	}, nil, nil)
	if err != nil {
		return nil, err
	}
	tx := wire.NewMsgTx(btc.DefaultTxVersion)

	// add input
	txHash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		return nil, err
	}
	txIn := wire.NewTxIn(wire.NewOutPoint(txHash, 0), nil, nil)
	tx.AddTxIn(txIn)

	// add output
	scriptPubKey, err := txscript.PayToAddrScript(recipient)
	if err != nil {
		return nil, err
	}

	txOut := wire.NewTxOut((amount - int64(fee)), scriptPubKey)
	tx.AddTxOut(txOut)

	// sign the script
	fetcher := txscript.NewCannedPrevOutputFetcher(script, amount)
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)
	sig, err := txscript.RawTxInWitnessSignature(tx, sigHashes, 0, amount, script, btc.SigHashSingleAnyoneCanPay, privKey)
	if err != nil {
		return nil, err
	}

	tx.TxIn[0].Witness = [][]byte{
		sig,
		script,
	}

	// serialize the tx
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getBalance(indexer btc.IndexerClient, address btcutil.Address) (int64, error) {
	utxos, err := indexer.GetUTXOs(context.Background(), address)
	if err != nil {
		return 0, err
	}
	balance := int64(0)
	for _, utxo := range utxos {
		balance += utxo.Amount
	}
	return balance, nil
}

func sigCheckScript(params chaincfg.Params, privkey *secp256k1.PrivateKey) ([]byte, *btcutil.AddressWitnessScriptHash, error) {
	builder := txscript.NewScriptBuilder()
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, nil, err
	}
	builder.AddData(randomBytes)
	builder.AddOp(txscript.OP_DROP)
	builder.AddData(privkey.PubKey().SerializeCompressed())
	builder.AddOp(txscript.OP_CHECKSIG)

	script, err := builder.Script()
	if err != nil {
		return nil, nil, err
	}

	hash := sha256.Sum256(script)
	scriptAdd, err := btcutil.NewAddressWitnessScriptHash(hash[:], &params)
	if err != nil {
		return nil, nil, err
	}
	return script, scriptAdd, nil
}

func sigCheckTapScript(params chaincfg.Params, pubkey []byte) ([]byte, *btcutil.AddressTaproot, []byte, error) {

	builder := txscript.NewScriptBuilder()
	builder.AddData(pubkey)
	builder.AddOp(txscript.OP_CHECKSIG)

	script, err := builder.Script()
	if err != nil {
		return nil, nil, nil, err
	}

	tapLeaf := txscript.NewBaseTapLeaf(script)
	tapScriptTree := txscript.AssembleTaprootScriptTree(tapLeaf)
	internalKey, err := btc.GardenNUMS()
	if err != nil {
		return nil, nil, nil, err
	}

	ctrlBlock := tapScriptTree.LeafMerkleProofs[0].ToControlBlock(
		internalKey,
	)
	tapScriptRootHash := tapScriptTree.RootNode.TapHash()
	outputKey := txscript.ComputeTaprootOutputKey(
		internalKey, tapScriptRootHash[:],
	)

	addr, err := btcutil.NewAddressTaproot(outputKey.X().Bytes(), &params)

	cbBytes, err := ctrlBlock.ToBytes()
	if err != nil {
		return nil, nil, nil, err
	}

	return script, addr, cbBytes, nil
}

func additionTapscript(params chaincfg.Params) ([]byte, *btcutil.AddressTaproot, []byte, error) {
	script, err := additionScriptBytes()
	if err != nil {
		return nil, nil, nil, err
	}

	tapLeaf := txscript.NewBaseTapLeaf(script)
	tapScriptTree := txscript.AssembleTaprootScriptTree(tapLeaf)
	internalKey, err := btc.GardenNUMS()
	if err != nil {
		return nil, nil, nil, err
	}

	// 0 is the leaf index
	ctrlBlock := tapScriptTree.LeafMerkleProofs[0].ToControlBlock(
		internalKey,
	)
	tapScriptRootHash := tapScriptTree.RootNode.TapHash()
	outputKey := txscript.ComputeTaprootOutputKey(
		internalKey, tapScriptRootHash[:],
	)

	addr, err := btcutil.NewAddressTaproot(outputKey.X().Bytes(), &params)

	cbBytes, err := ctrlBlock.ToBytes()
	if err != nil {
		return nil, nil, nil, err
	}

	return script, addr, cbBytes, nil
}

func additionScript(params chaincfg.Params) ([]byte, *btcutil.AddressWitnessScriptHash, error) {
	script, err := additionScriptBytes()
	if err != nil {
		return nil, nil, err
	}
	hash := sha256.Sum256(script)
	scriptAdd, err := btcutil.NewAddressWitnessScriptHash(hash[:], &params)
	if err != nil {
		return nil, nil, err
	}
	return script, scriptAdd, nil
}

func additionScriptBytes() ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	// add random bytes
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	builder.AddData(randomBytes)
	builder.AddOp(txscript.OP_DROP)
	builder.AddOp(txscript.OP_ADD)
	builder.AddOp(txscript.OP_2)
	builder.AddOp(txscript.OP_EQUAL)

	script, err := builder.Script()
	if err != nil {
		return nil, err
	}
	return script, nil
}
