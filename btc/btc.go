package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
)

const (
	// DefaultTxVersion is the Bitcoin transaction version used by this package. Some script operations are only
	// supported in version 2.
	DefaultTxVersion = 2

	// DustAmount is the minimum transaction amount accepted by Bitcoin miners.
	DustAmount = 546

	// SigHashSingleAnyoneCanPay is an alias for the signature hash types: `txscript.SigHashSingle |
	// txscript.SigHashAnyOneCanPay`.
	SigHashSingleAnyoneCanPay = txscript.SigHashSingle | txscript.SigHashAnyOneCanPay
)

// SizeUpdater returns the base and segwit size of signing a particular type of utxo. This is used by the tx building
// function to estimate the fee.
type SizeUpdater func() (int, int)

var (
	P2pkhUpdater = func() (int, int) {
		return txsizes.RedeemP2PKHSigScriptSize, 0
	}

	P2wpkhUpdater = func() (int, int) {
		return 0, txsizes.RedeemP2WPKHInputWitnessWeight
	}

	HtlcUpdater = func(secretSize int) func() (int, int) {
		return func() (int, int) {
			if secretSize == 0 {
				return 0, RedeemHtlcRefundSigScriptSize
			}
			return 0, RedeemHtlcRedeemSigScriptSize(secretSize)
		}
	}

	MultisigUpdater = func() (int, int) {
		return 0, RedeemMultisigSigScriptSize
	}
)

type UTXOs []UTXO

type UTXO struct {
	TxID   string  `json:"txid"`
	Vout   uint32  `json:"vout"`
	Amount int64   `json:"value"`
	Status *Status `json:"status"`
}

type Recipient struct {
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

type RawInputs struct {
	VIN        []UTXO
	BaseSize   int
	SegwitSize int
}

func NewRawInputs() RawInputs {
	return RawInputs{
		VIN:        nil,
		BaseSize:   0,
		SegwitSize: 0,
	}
}

// BuildTransaction is a helper function for building a bitcoin transaction. It uses the given `feeRate` to calculate
// fees. `inputs` will be a list of utxos that required to be included in the transaction, it comes with the base and
// segwit size of the signature for fee-estimation purpose. `utxos` is a list of transaction will be picked
// to cover the output amount and fees. We assume the utxos all comes from a single address. The `sizeUpdater` function
// returns the base and segwit size of each utxo from the `utxos`. If there's any change, it will be sent back to the
// `changeAddr`.
func BuildTransaction(network *chaincfg.Params, feeRate int, inputs RawInputs, utxos []UTXO, sizeUpdater SizeUpdater, recipients []Recipient, changeAddr btcutil.Address) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(DefaultTxVersion)
	totalIn, totalOut := int64(0), int64(0)
	base, segwit := inputs.BaseSize, inputs.SegwitSize

	// Adding required inputs and output
	for _, utxo := range inputs.VIN {
		// // Skip the utxo if the amount is not large enough.
		// if utxo.Amount <= int64(minUtxoValue) {
		// 	continue
		// }
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

		vs := EstimateVirtualSize(tx, base, segwit)
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
					vs := EstimateVirtualSize(tx, base, segwit)
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

	// Check if the existing inputs are enough and we might not need to add extra utxos
	enough, err := valueCheck()
	if err != nil {
		return nil, err
	}
	if enough {
		return tx, nil
	}

	// Calculate the minimum utxo value we want to add to the tx.
	// This is to prevent adding a dust utxo and cause the tx use more fees.
	minUtxoValue := 0
	if sizeUpdater != nil {
		minBase, minSegwit := sizeUpdater()
		minVS := minBase + (minSegwit+3)/blockchain.WitnessScaleFactor
		minUtxoValue = minVS * feeRate
	}

	// Keep adding utxos until we have enough funds to cover the output amount
	for _, utxo := range utxos {
		// Skip dust utxo
		if utxo.Amount < int64(minUtxoValue) {
			continue
		}

		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}
		tx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(hash, utxo.Vout), nil, nil))
		totalIn += utxo.Amount
		if sizeUpdater != nil {
			additionalBaseSize, additionalSegwitSize := sizeUpdater()
			base += additionalBaseSize
			segwit += additionalSegwitSize
		}

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
func BuildRbfTransaction(network *chaincfg.Params, feeRate int, inputs RawInputs, utxos []UTXO, sizeUpdater SizeUpdater, recipients []Recipient, changeAddr btcutil.Address) (*wire.MsgTx, error) {
	tx, err := BuildTransaction(network, feeRate, inputs, utxos, sizeUpdater, recipients, changeAddr)
	if err != nil {
		return nil, err
	}
	for i := range tx.TxIn {
		tx.TxIn[i].Sequence = mempool.MaxRBFSequence
	}
	return tx, nil
}

// SignP2pkhTx is a helper function to sign inputs from a p2pkh address. It uses `txscript.SigHashAll` and compressed
// public key as default.
func SignP2pkhTx(network *chaincfg.Params, key *btcec.PrivateKey, tx *wire.MsgTx) error {
	addr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(key.PubKey().SerializeCompressed()), network)
	if err != nil {
		return err
	}
	for i := range tx.TxIn {
		pkScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return err
		}

		sigScript, err := txscript.SignatureScript(tx, i, pkScript, txscript.SigHashAll, key, true)
		if err != nil {
			return err
		}
		tx.TxIn[i].SignatureScript = sigScript
	}
	return nil
}

// PublicKeyAddress is a helper function which calculates the address of the given public key. This function only
// supports address types are derived from a public key.
func PublicKeyAddress(network *chaincfg.Params, addrType waddrmgr.AddressType, pub *btcec.PublicKey) (btcutil.Address, error) {
	switch addrType {
	case waddrmgr.RawPubKey:
		return btcutil.NewAddressPubKey(pub.SerializeCompressed(), network)
	case waddrmgr.PubKeyHash:
		return btcutil.NewAddressPubKeyHash(btcutil.Hash160(pub.SerializeCompressed()), network)
	case waddrmgr.WitnessPubKey:
		return btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pub.SerializeCompressed()), network)
	case waddrmgr.TaprootPubKey:
		tapKey := txscript.ComputeTaprootKeyNoScript(pub)
		return btcutil.NewAddressTaproot(schnorr.SerializePubKey(tapKey), network)
	default:
		return nil, fmt.Errorf("unsupported address type")
	}
}
