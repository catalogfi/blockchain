package btc

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
)

const (
	// DefaultTxVersion is the Bitcoin transaction version used by this package. Some script operations are only
	// supported in version 2.
	DefaultTxVersion = 2

	// DustAmount is the minimum transaction amount accepted by Bitcoin miners.
	DustAmount = 546

	// MinRelayFeeRate is the minimum fee rate (in sat/vb) a transaction must meet in order to be broadcast by the node.
	// This is default value used by Bitcore Core. Different node may have different setting for this.
	MinRelayFeeRate = 1

	// MaxRelayFeeRate is not something in the bitcoin protocol, but more of a defensive check to make sure we're not
	// using an unreasonable value. (i.e. third-party api) If the fee rate we choose is greater than this value, this
	// usually means something is wrong.
	MaxRelayFeeRate = 1000

	// SigHashSingleAnyoneCanPay is an alias for the signature hash types: `txscript.SigHashSingle |
	// txscript.SigHashAnyOneCanPay`.
	SigHashSingleAnyoneCanPay = txscript.SigHashSingle | txscript.SigHashAnyOneCanPay
)

// UTXOs is a list of UTXOs.
type UTXOs []UTXO

// UTXO is an unspent transaction output.
type UTXO struct {
	TxID   string  `json:"txid"`
	Vout   uint32  `json:"vout"`
	Amount int64   `json:"value"`
	Status *Status `json:"status"`
}

// String returns a string that identifies this UTXO, it is in the format of "<txid>:<vout>".
func (utxo UTXO) String() string {
	return fmt.Sprintf("%v:%v", utxo.TxID, utxo.Vout)
}

// Recipient is a recipient of a transaction. It contains the address and the amount to be sent.
type Recipient struct {
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

// SingleRecipient is a helper function to initiate a recipient list with only one element.
func SingleRecipient(to string, amount int64) []Recipient {
	return []Recipient{
		{
			To:     to,
			Amount: amount,
		},
	}
}

// PublicKeyAddress generates a Bitcoin address from a given public key, depending on the specified address type.
// If an unsupported address type is provided, it returns an error.
func PublicKeyAddress(network *chaincfg.Params, addrType waddrmgr.AddressType, pub *btcec.PublicKey) (btcutil.Address, error) {
	switch addrType {
	case waddrmgr.RawPubKey:
		return btcutil.NewAddressPubKey(pub.SerializeCompressed(), network)
	case waddrmgr.PubKeyHash:
		return btcutil.NewAddressPubKeyHash(btcutil.Hash160(pub.SerializeCompressed()), network)
	case waddrmgr.WitnessPubKey:
		return btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pub.SerializeCompressed()), network)
	case waddrmgr.TaprootPubKey:
		// We assume the key is already tweaked. You'll need to tweak the key first if it's the internal key.
		return btcutil.NewAddressTaproot(schnorr.SerializePubKey(pub), network)
	default:
		return nil, fmt.Errorf("unsupported address type = %v", addrType)
	}
}

// BuildTransaction is a helper function for building a bitcoin transaction. It uses the given `feeRate` to calculate
// fees. `inputs` will be a list of utxos that required to be included in the transaction. `utxos` is a list of
// transaction will be picked to cover the output amount and fees. The size estimator is to estimate the size of the
// signed tx for fee purpose. The `recipients` includes a list of target addresses and the associated amounts to be
// sent.  If there's any change, it will be sent back to the `changeAddr`.
func BuildTransaction(network *chaincfg.Params, feeRate int, inputs, utxos []UTXO, sizeEstimator *SizeEstimator, recipients []Recipient, changeAddr btcutil.Address) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(DefaultTxVersion)
	totalIn, totalOut := int64(0), int64(0)

	if feeRate < MinRelayFeeRate {
		return nil, fmt.Errorf("fee rate too low, expect %v got %v", MinRelayFeeRate, feeRate)
	}
	if feeRate > MaxRelayFeeRate {
		return nil, fmt.Errorf("fee rate too high, expect %v got %v", MaxRelayFeeRate, feeRate)
	}

	// Adding required inputs and output
	for _, utxo := range inputs {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}
		txIn := wire.NewTxIn(wire.NewOutPoint(hash, utxo.Vout), nil, nil)
		tx.AddTxIn(txIn)
		if utxo.Amount == 0 {
			return nil, fmt.Errorf("utxo amount is not set")
		}
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

	// Function to check if the input amount is greater than or equal to the output amount plus fees
	valueCheck := func() (bool, error) {
		if totalIn <= totalOut {
			return false, nil
		}

		vs, err := sizeEstimator.EstimateTxVirtualSize(tx)
		if err != nil {
			return false, err
		}
		fees := int64(vs * feeRate)

		// If the amount is enough to cover the outputs and fees
		if totalIn > totalOut+fees {
			// Add a change utxo to the output if the change amount is greater than the dust
			if totalIn-totalOut-fees > DustAmount {
				if changeAddr != nil {
					changeScript, err := txscript.PayToAddrScript(changeAddr)
					if err != nil {
						return false, err
					}
					tx.AddTxOut(wire.NewTxOut(0, changeScript)) // adjust the amount later

					// Estimate the fees again as we add a new output
					vs, err := sizeEstimator.EstimateTxVirtualSize(tx)
					if err != nil {
						return false, err
					}
					fees := int64(vs * feeRate)

					// Adjust the change utxo amount if it's still enough, delete it otherwise
					if totalIn-totalOut-fees > DustAmount {
						tx.TxOut[len(tx.TxOut)-1].Value = totalIn - totalOut - fees
					} else {
						tx.TxOut = tx.TxOut[:len(tx.TxOut)-1]
					}
				}
			}

			return true, nil
		}

		return false, nil
	}

	// Check if the existing inputs are enough and we might not need to add any extra utxo
	enough, err := valueCheck()
	if err != nil {
		return nil, err
	}
	if enough {
		return tx, nil
	}

	// Keep adding utxos until we have enough funds to cover the output amount
	for _, utxo := range utxos {

		// Check if it's worth to add the tx by calculating the cost
		// (The utxos need to be confirmed in this case as we don't do package calculation)
		// We use a hard coded value (200) for the size of adding this utxo, this is close to the tx size
		// of spending a P2PKH utxo. If we really want, we can do a more precise estimation, which is
		// base = outpoint(36) + sigscript(1)+ sequence(4) | segwit = marker(1) + flag(1) + witness(x)
		// worstVS := (base + 36 + 1 + 4) + (segwit+2+3)/blockchain.WitnessScaleFactor
		cost := 200 * feeRate
		if int64(cost) > utxo.Amount {
			continue
		}

		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}
		tx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(hash, utxo.Vout), nil, nil))
		totalIn += utxo.Amount

		// Check if we have enough inputs to cover the outputs and fee
		enough, err := valueCheck()
		if err != nil {
			return nil, err
		}
		if enough {
			return tx, nil
		}
	}

	return nil, fmt.Errorf("funds not enough")
}

// BuildRbfTransaction is similar to `BuildTransaction`, the only difference is it updates the sequence of all the tx
// inputs to `mempool.MaxRBFSequence`, so the tx is RBF-compatible.
func BuildRbfTransaction(network *chaincfg.Params, feeRate int, inputs, utxos []UTXO, sizeEstimator *SizeEstimator, recipients []Recipient, changeAddr btcutil.Address) (*wire.MsgTx, error) {
	tx, err := BuildTransaction(network, feeRate, inputs, utxos, sizeEstimator, recipients, changeAddr)
	if err != nil {
		return nil, err
	}
	for i := range tx.TxIn {
		tx.TxIn[i].Sequence = mempool.MaxRBFSequence
	}
	return tx, nil
}

// TxRawBytes returns the raw bytes of a transaction.
func TxRawBytes(tx *wire.MsgTx) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// InitFetcher initializes a txscript.MultiPrevOutFetcher with the given utxos and script.
// The returned fetcher can be used to sign transactions with the given utxos as inputs.
func InitFetcher(utxos []UTXO, script []byte) (*txscript.MultiPrevOutFetcher, error) {
	fetcher := txscript.NewMultiPrevOutFetcher(nil)
	for _, utxo := range utxos {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}
		fetcher.AddPrevOut(wire.OutPoint{
			Hash:  *hash,
			Index: utxo.Vout,
		}, wire.NewTxOut(utxo.Amount, script))
	}
	return fetcher, nil
}
