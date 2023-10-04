package btc

import (
	"crypto/sha512"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
	"golang.org/x/crypto/pbkdf2"
)

const (
	// DefaultBTCVersion is the Bitcoin transaction version used by this package. Some script operations are only
	// supported in version 2.
	DefaultBTCVersion = 2

	// DustAmount is the minimum transaction amount accepted by Bitcoin miners.
	DustAmount = 546

	// SigHashSingleAnyoneCanPay is an alias for the signature hash types: `txscript.SigHashSingle |
	// txscript.SigHashAnyOneCanPay`.
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

// BuildTransaction is helper function for building a bitcoin transaction. It uses the `feeEstimator` to get the current
// fee rates. `inputs` will be a list of utxos that required to be included in the transaction. `extraBaseSize` and
// `extraSegwitSize` will be the extra base/segwit size of the `inputs`. `utxos` is a list of transaction will be picked
// to cover the output amount and fees. If there's any change, it will be sent back to the `changeAddr`.
func BuildTransaction(feeRate int, network *chaincfg.Params, inputs, utxos []UTXO,
	recipients []Recipient, extraBaseSize, extraSegwitSize int, changeAddr btcutil.Address) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(DefaultBTCVersion)
	totalIn, totalOut := int64(0), int64(0)

	// Adding required inputs
	for _, utxo := range inputs {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}
		txIn := wire.NewTxIn(wire.NewOutPoint(hash, utxo.Vout), nil, nil)
		tx.AddTxIn(txIn)
		totalIn += utxo.Amount
	}
	for _, recipient := range recipients {
		toAddress, err := btcutil.DecodeAddress(recipient.To, network)
		if err != nil {
			return nil, err
		}
		toScript, err := txscript.PayToAddrScript(toAddress)
		if err != nil {
			return nil, err
		}
		tx.AddTxOut(wire.NewTxOut(recipient.Amount, toScript))
		totalOut += recipient.Amount
	}

	// Keep adding utxos until we have enough funds to cover the output amount
	for _, utxo := range utxos {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}
		txIn := wire.NewTxIn(wire.NewOutPoint(hash, utxo.Vout), nil, nil)
		tx.AddTxIn(txIn)
		totalIn += utxo.Amount
		extraBaseSize += txsizes.RedeemP2PKHSigScriptSize

		// Check the fees when we have enough inputs to cover the outputs
		if totalIn > totalOut {
			vs := EstimateVirtualSize(tx, extraBaseSize, extraSegwitSize)
			fees := int64(vs * feeRate)

			// If the amount is enough to cover the outputs and fees
			if totalIn > totalOut+fees {
				// Add a change utxo to the output if the change amount is greater than the dust
				if totalIn-totalOut-fees > DustAmount {
					changeScript, err := txscript.PayToAddrScript(changeAddr)
					if err != nil {
						return nil, err
					}
					tx.AddTxOut(wire.NewTxOut(0, changeScript))

					// Estimate the fees again as we add a new output
					vs := EstimateVirtualSize(tx, extraBaseSize, extraSegwitSize)
					fees := int64(vs * feeRate)

					// Adjust the change utxo amount if it's still enough, delete it otherwise
					if totalIn-totalOut-fees > DustAmount {
						tx.TxOut[len(tx.TxOut)-1].Value = totalIn - totalOut - fees
					} else {
						tx.TxOut = tx.TxOut[:len(tx.TxOut)-1]
					}
				}
				return tx, nil
			}
		}
	}
	return nil, fmt.Errorf("funds not enough")
}
