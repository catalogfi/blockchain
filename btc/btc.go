package btc

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"golang.org/x/crypto/pbkdf2"
)

const (
	// DefaultBTCVersion is the version field set in the bitcoin transaction. Since some of the script ops are only
	// supported in version 2. All functions in this package to build a tx will use this version number.
	DefaultBTCVersion = 2

	// DustAmount is the minimum amount we'll accept for an UTXO.
	DustAmount = 546

	// SigHashSingleAnyoneCanPay is an alias for the `txscript.SigHashSingle | txscript.SigHashAnyOneCanPay`
	SigHashSingleAnyoneCanPay = txscript.SigHashSingle | txscript.SigHashAnyOneCanPay
)

type UTXOs []UTXO

type UTXO struct {
	TxID   string `json:"txid"`
	Vout   uint32 `json:"vout"`
	Amount int64  `json:"value"`
}

type Recipient struct {
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

func GenerateSystemPrivKey(mnemonic string, userPubkeyHex []byte) (*btcec.PrivateKey, error) {
	seed := pbkdf2.Key([]byte(mnemonic), append([]byte("mnemonic"), userPubkeyHex...), 2048, 64, sha512.New)
	var masterKey *hdkeychain.ExtendedKey
	var err error
	for {
		// Network parameter doesn't affect key generation, so we use mainnet params here.
		masterKey, err = hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
		if err != nil {
			// Very small chance to happen
			if err == hdkeychain.ErrUnusableSeed {
				continue
			}
			return nil, err
		}
		break
	}

	return masterKey.ECPrivKey()
}

// PickUTXOs will find enough utxos which are enough to cover the target amount and fee.
func PickUTXOs(utxos []UTXO, amount, fee int64) ([]UTXO, error) {
	var inputs UTXOs
	inputAmount := int64(0)
	for _, utxo := range utxos {
		inputAmount += utxo.Amount
		inputs = append(inputs, utxo)
		if inputAmount >= amount+fee {
			break
		}
	}
	if inputAmount < amount+fee {
		return nil, fmt.Errorf("not enough funds")
	}
	return inputs, nil
}

// BuildTx creates an unsigned txs with the given inputs and outputs. It will add an output to the `fromAddress` if
// there's any change left. The returned boolean will indicate if there's a change utxo.
func BuildTx(network *chaincfg.Params, inputs UTXOs, recipients []Recipient, fee int64, fromAddress btcutil.Address) (*wire.MsgTx, bool, error) {
	tx := wire.NewMsgTx(DefaultBTCVersion)
	totalIn, totalOut := int64(0), int64(0)

	// Inputs
	for _, utxo := range inputs {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, false, err
		}
		txIn := wire.NewTxIn(wire.NewOutPoint(hash, utxo.Vout), nil, nil)
		tx.AddTxIn(txIn)
		totalIn += utxo.Amount
	}

	// Outputs
	for _, recipient := range recipients {
		toAddress, err := btcutil.DecodeAddress(recipient.To, network)
		if err != nil {
			return nil, false, err
		}
		toScript, err := txscript.PayToAddrScript(toAddress)
		if err != nil {
			return nil, false, err
		}
		tx.AddTxOut(wire.NewTxOut(recipient.Amount, toScript))
		totalOut += recipient.Amount
	}

	// If there's any leftover, we send it back to fromAddress.
	changeTx := false
	if totalIn-totalOut-fee > DustAmount {
		fromScript, err := txscript.PayToAddrScript(fromAddress)
		if err != nil {
			return nil, false, err
		}
		changeTx = true
		tx.AddTxOut(wire.NewTxOut(totalIn-totalOut-fee, fromScript))
	}

	return tx, changeTx, nil
}

// P2wshAddress returns the P2WSH address of the give script
func P2wshAddress(script []byte, network *chaincfg.Params) (string, error) {
	scriptHash := sha256.Sum256(script)
	instantWalletAddr, err := btcutil.NewAddressWitnessScriptHash(scriptHash[:], network)
	if err != nil {
		return "", err
	}
	return instantWalletAddr.EncodeAddress(), nil
}
