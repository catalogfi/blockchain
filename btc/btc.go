package btc

import (
	"bytes"
	"fmt"
	"math"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
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

	// MinRelayFeeRate is the minimum feerate (in sat/kvb) a transaction must meet in order to be broadcast by the node.
	// This is a default value used by Bitcore Core. Different node may have different setting for this.
	MinRelayFeeRate = 1 * 1000

	// MaxRelayFeeRate is not something in the bitcoin protocol, but more of a defensive check to make sure we're not
	// using an unreasonable value. (i.e. third-party api) If the fee rate we choose is greater than this value, this
	// usually means something is wrong.
	MaxRelayFeeRate = 500 * 1000

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

// NewUtxo returns a new Utxo object with given parameters. It won't have the status field.
func NewUtxo(txid string, vout uint32, amount int64) UTXO {
	return UTXO{
		TxID:   txid,
		Vout:   vout,
		Amount: amount,
	}
}

// String returns a string that identifies this UTXO, it is in the format of "<txid>:<vout>".
func (utxo UTXO) String() string {
	return fmt.Sprintf("%v:%v", utxo.TxID, utxo.Vout)
}

// ToOutPoint converts the UTXO to the type of `*wire.OutPoint`
func (utxo UTXO) ToOutPoint() (*wire.OutPoint, error) {
	hash, err := chainhash.NewHashFromStr(utxo.TxID)
	if err != nil {
		return nil, err
	}
	return wire.NewOutPoint(hash, utxo.Vout), nil
}

// ToTxIn converts the UTXO to a `*wire.TxIn` so it can be easily added to a new tx. It won't have any witness data.
func (utxo UTXO) ToTxIn() (*wire.TxIn, error) {
	outpoint, err := utxo.ToOutPoint()
	if err != nil {
		return nil, err
	}
	return wire.NewTxIn(outpoint, nil, nil), nil
}

// Recipient is a recipient of a transaction. It contains the address and the amount to be sent.
type Recipient struct {
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

// NewRecipient constructs a new Recipient object with the given address and amount.
func NewRecipient(to string, amount int64) Recipient {
	return Recipient{
		To:     to,
		Amount: amount,
	}
}

// ToTxOut converts the Recipient to a `*wire.TxOut` so it can be easily added to a new tx.
func (recipient Recipient) ToTxOut(network *chaincfg.Params) (*wire.TxOut, error) {
	toAddress, err := btcutil.DecodeAddress(recipient.To, network)
	if err != nil {
		return nil, err
	}
	toScript, err := txscript.PayToAddrScript(toAddress)
	if err != nil {
		return nil, err
	}
	return wire.NewTxOut(recipient.Amount, toScript), nil
}

// PublicKeyAddress generates a Bitcoin address from a given public key, depending on the specified address type.
// If an unsupported address type is provided, it returns an error.
func PublicKeyAddress(network *chaincfg.Params, addrType waddrmgr.AddressType, pub *btcec.PublicKey) (btcutil.Address, error) {
	switch addrType {
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

// TxRawBytes returns the raw bytes of a transaction.
func TxRawBytes(tx *wire.MsgTx) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// FeeMode is a function type that calculates the required fees for a given transaction under different scenario.
type FeeMode func(tx *wire.MsgTx) (int64, error)

// MinFeeRateMode is used to build a tx with a minimum feeRate. The actual fee rate will be greater than or equal
// to the given `feeRate` (sats/kvb).
func MinFeeRateMode(feeRate int, sizer *SizeEstimator) FeeMode {
	return func(tx *wire.MsgTx) (int64, error) {
		vs, err := sizer.EstimateTxVirtualSize(tx)
		if err != nil {
			return 0, err
		}
		return int64(vs*feeRate+999) / 1000, nil
	}
}

// FixedFeesMode is used to build a tx with a fixed amount of fees.
func FixedFeesMode(fees int64) FeeMode {
	return func(tx *wire.MsgTx) (int64, error) {
		return fees, nil
	}
}

// GaslessMode will always return 0 fees. This is usually used to build a gas-less transaction.
func GaslessMode() FeeMode {
	return func(tx *wire.MsgTx) (int64, error) {
		return 0, nil
	}
}

// RbfMode computes the required fees for a transaction to be replace-by-fee (RBF) compliant.
// It considers the minimum fee rate, previous fee rate, and previous fees to calculate the maximum
// fees needed for the transaction. Both `minFeeRate` and `prevFeeRate` will be in sats/kvb.
func RbfMode(minFeeRate, prevFeeRate int, prevFees int64, sizer *SizeEstimator) FeeMode {
	return func(tx *wire.MsgTx) (int64, error) {
		weight, err := sizer.EstimateTxWeight(tx)
		if err != nil {
			return 0, err
		}
		vsize := (weight + 3) / 4
		fees1 := math.Ceil(float64((prevFeeRate+1)*vsize) / 1000)
		fees2 := math.Ceil(float64(prevFees) + float64(weight)/4)
		fees3 := math.Ceil(float64(minFeeRate * vsize / 1000))
		return int64(math.Max(math.Max(fees1, fees2), fees3)), nil
	}
}

// BuildTx is a helper function for building a bitcoin transaction. It uses the given `FeeMode` to calculate
// fees. `inputs` will be a list of utxos that guaranteed to be included in the transaction. `utxos` is a list of
// transaction will be picked to cover the output amount and fees. The `recipients` includes a list of target addresses
// and the associated amounts to be sent.  If there's any change, it will be sent back to the `changeAddr`.
func BuildTx(network *chaincfg.Params, feeReq FeeMode, inputs, utxos []UTXO, recipients []Recipient, changeAddr btcutil.Address) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(DefaultTxVersion)
	totalIn, totalOut := int64(0), int64(0)

	// Adding required inputs and output
	for _, utxo := range inputs {
		if utxo.Amount == 0 {
			return nil, fmt.Errorf("utxo %v amount is not set", utxo.String())
		}
		txIn, err := utxo.ToTxIn()
		if err != nil {
			return nil, err
		}
		tx.AddTxIn(txIn)
		totalIn += utxo.Amount
	}
	for _, recipient := range recipients {
		txOut, err := recipient.ToTxOut(network)
		if err != nil {
			return nil, err
		}
		tx.AddTxOut(txOut)
		totalOut += recipient.Amount
	}

	for i := -1; i < len(utxos); i++ {
		// Add the utxo to the transaction input
		if i >= 0 {
			txin, err := utxos[i].ToTxIn()
			if err != nil {
				return nil, err
			}
			tx.AddTxIn(txin)
			totalIn += utxos[i].Amount
		}

		// Check if the input amount is enough to cover the fees and output
		fees, err := feeReq(tx)
		if err != nil {
			return nil, err
		}
		if totalIn >= totalOut+fees {
			// Add a change utxo to the output if the change amount is greater than the dust
			if totalIn-totalOut-fees > DustAmount {
				if changeAddr != nil {
					changeScript, err := txscript.PayToAddrScript(changeAddr)
					if err != nil {
						return nil, err
					}
					tx.AddTxOut(wire.NewTxOut(0, changeScript)) // adjust the amount later

					// Fees will be changed since we add a new output field to the tx
					fees, err := feeReq(tx)
					if err != nil {
						return nil, err
					}

					// Adjust the change utxo amount if it's still enough, delete it otherwise
					if totalIn-totalOut-fees > DustAmount {
						tx.TxOut[len(tx.TxOut)-1].Value = totalIn - totalOut - fees
					} else {
						tx.TxOut = tx.TxOut[:len(tx.TxOut)-1]
					}
				}
			}
			return tx, nil
		}
	}

	return nil, fmt.Errorf("funds not enough")
}
