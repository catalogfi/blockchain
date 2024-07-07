package btc_test

import (
	"context"
	"crypto/sha256"
	"math/rand"
	"os"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/localnet"
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

		txid, err := wallet.Send(context.Background(), req)
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

		txid, err := wallet.Send(context.Background(), req)
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
		_, err = wallet.Send(context.Background(), req)

		Expect(err).ShouldNot(BeNil())
		req = []btc.SendRequest{
			{
				Amount: 100000000 - 1000,
				To:     wallet.Address(),
			},
		}
		_, err = wallet.Send(context.Background(), req)
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
		})
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
		txid, err := bobWallet.Send(context.Background(), req)
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

		txid, err = bobWallet.Send(context.Background(), req)
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

		txId, err := wallet.Spend(context.Background(), &btc.SpendRequest{
			Inputs: []btc.SpendScripts{
				{
					Witness: [][]byte{
						btc.AddSignatureSegwitOp,
						btc.AddPubkeyCompressedOp,
					},
					Script:        pkScript,
					ScriptAddress: address,
					HashType:      txscript.SigHashAll,
				},
			},
			ToAddress: address,
		})
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
		})
		Expect(err).To(BeNil())

		// spend the script
		txId, err := wallet.Spend(context.Background(), &btc.SpendRequest{
			Inputs: []btc.SpendScripts{
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
			},
			ToAddress: wallet.Address(),
		})
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
		})
		Expect(err).To(BeNil())

		// spend the script
		txId, err := wallet.Spend(context.Background(), &btc.SpendRequest{
			Inputs: []btc.SpendScripts{
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
			},
			ToAddress: wallet.Address(),
		})
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
		})
		Expect(err).To(BeNil())

		txId, err := wallet.Spend(context.Background(), &btc.SpendRequest{
			Inputs: []btc.SpendScripts{
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
			},
			ToAddress: wallet.Address(),
		})
		Expect(err).To(BeNil())
		Expect(txId).ShouldNot(BeEmpty())

	})

	It("should be able to spend funds from signature-check script p2wsh", func() {
		script, scriptAddr, err := sigCheckScript(chainParams)
		Expect(err).To(BeNil())
		scriptBalance := int64(100000)
		_, err = wallet.Send(context.Background(), []btc.SendRequest{
			{
				Amount: scriptBalance,
				To:     scriptAddr,
			},
		})
		Expect(err).To(BeNil())

		// spend the script
		txId, err := wallet.Spend(context.Background(), &btc.SpendRequest{
			Inputs: []btc.SpendScripts{
				{
					Witness: [][]byte{
						btc.AddSignatureSegwitOp,
						btc.AddPubkeyCompressedOp,
						script,
					},
					Script:        script,
					ScriptAddress: scriptAddr,
					HashType:      txscript.SigHashAll,
				},
			},
			ToAddress: wallet.Address(),
		})
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
		})

		txId, err := wallet.Spend(context.Background(), &btc.SpendRequest{
			Inputs: []btc.SpendScripts{
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
			},
			ToAddress: wallet.Address(),
		})

		Expect(err).To(BeNil())
		Expect(txId).ShouldNot(BeEmpty())
	})

	It("should be able to spend funds from different (p2wsh and p2tr) scripts", func() {
		p2wshAdditionScript, p2wshAddr, err := additionScript(chainParams)
		Expect(err).To(BeNil())

		p2trAdditionScript, p2trAddr, cb, err := additionTapscript(chainParams)
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
		})
		By("Spend p2wsh and p2tr scripts")
		txId, err := wallet.Spend(context.Background(), &btc.SpendRequest{
			Inputs: []btc.SpendScripts{
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
					Script:        p2trAdditionScript,
					Leaf:          txscript.NewTapLeaf(0xc0, p2trAdditionScript),
					ScriptAddress: p2trAddr,
				},
			},
			ToAddress: wallet.Address(),
		})

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

		_, err = wallet.Spend(context.Background(), &btc.SpendRequest{
			Inputs: []btc.SpendScripts{
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
			},
			ToAddress: wallet.Address(),
		})
		Expect(err).ShouldNot(BeNil())
	})

	It("should not be able to spend with invalid Inputs", func() {
		_, err := wallet.Spend(context.Background(), &btc.SpendRequest{
			Inputs: []btc.SpendScripts{
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
			},
			ToAddress: wallet.Address(),
		})
		Expect(err).ShouldNot(BeNil())

		By("Lets give proper witness but bad script")
		_, err = wallet.Spend(context.Background(), &btc.SpendRequest{
			Inputs: []btc.SpendScripts{
				{
					Witness: [][]byte{
						btc.AddSignatureSegwitOp,
						btc.AddPubkeyCompressedOp,
					},
					Script:        []byte{0x1},
					ScriptAddress: wallet.Address(),
					HashType:      txscript.SigHashAll,
				},
			},
			ToAddress: wallet.Address(),
		})
		Expect(err).ShouldNot(BeNil())
	})

	It("should be able to fetch the status and tx details", func() {
		req := []btc.SendRequest{
			{
				Amount: 100000,
				To:     wallet.Address(),
			},
		}

		txid, err := wallet.Send(context.Background(), req)
		Expect(err).To(BeNil())

		tx, ok, err := wallet.Status(context.Background(), txid)
		Expect(err).To(BeNil())
		Expect(ok).Should(BeTrue())
		Expect(tx).ShouldNot(BeNil())

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		tx, ok, err = wallet.Status(ctx, "f"+txid[1:])
		Expect(err).To(BeNil())
		Expect(ok).Should(BeFalse())
		Expect(tx.TxID).Should(BeEmpty())
	})

})

//------------------------------ Helper functions for tests --------------------------------

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

func sigCheckScript(params chaincfg.Params) ([]byte, *btcutil.AddressWitnessScriptHash, error) {
	builder := txscript.NewScriptBuilder()
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
