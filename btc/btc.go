package btc

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcjson"
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

	// MinRelayFeeRate is the minimum feerate (in sat/kvb) a transaction must meet in order to be broadcast by the node.
	// This is a default value used by Bitcore Core. Different node may have different setting for this.
	MinRelayFeeRate = SatoshiPerKb(1e3)

	// MaxRelayFeeRate is not something in the bitcoin protocol, but more of a defensive check to make sure we're not
	// using an unreasonable value. (i.e. third-party api) If the fee rate we choose is greater than this value, this
	// usually means something is wrong.
	MaxRelayFeeRate = SatoshiPerKb(500e3)

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

// FeeMode defines the rules on how to calculate the fees when building a transaction. It can be called after a tx is
// build and it will return the minimum fees(in sats) required by this mode.
type FeeMode func(tx *wire.MsgTx) (int64, error)

// MinFeeRateMode is used to build a tx with a minimum feeRate. The actual fee rate will be at lease the given `feeRate`
func MinFeeRateMode(feeRate SatoshiPerKb, sizer *SizeEstimator) FeeMode {
	return func(tx *wire.MsgTx) (int64, error) {
		vs, err := sizer.EstimateTxVirtualSize(tx)
		if err != nil {
			return 0, err
		}
		return int64(vs*feeRate.Int()+999) / 1000, nil
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
// `minFeeRate` is an optional parameter which set the minimum fee rate you want to use, similar to `MinFeeRateMode`.
// Use 0 if you don't want to have a minimum fee rate.
// `prevFeeRate` is the maximum fee rate of all directly conflicting transactions. The new fee rate has to be greater
// than this according to the mempool policy.
// `prevFees`is the sum of fees of paid by the original transactions. This includes teh replacement transaction and all
// its descendants.
func RbfMode(minFeeRate, prevFeeRate SatoshiPerKb, prevFees int64, sizer *SizeEstimator) FeeMode {
	return func(tx *wire.MsgTx) (int64, error) {
		vsize, err := sizer.EstimateTxVirtualSize(tx)
		if err != nil {
			return 0, err
		}
		fees1 := math.Ceil(float64((prevFeeRate.Int()+1)*vsize) / 1000)
		fees2 := math.Ceil(float64(prevFees) + float64(vsize))
		fees3 := math.Ceil(float64(minFeeRate.Int()*vsize) / 1000)
		return int64(math.Max(math.Max(fees1, fees2), fees3)), nil
	}
}

func RbfModeFromPrevTx(client Client, txid string, sizer *SizeEstimator) (FeeMode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	entry, err := client.GetMempoolEntry(ctx, txid)
	if err != nil {
		return nil, err
	}
	prevFees, err := btcutil.NewAmount(entry.Fees.Descendant)
	if err != nil {
		return nil, err
	}
	prevFeeRate := NewSatoshiPerKb(int64(prevFees), int(entry.DescendantSize))
	return RbfMode(0, prevFeeRate, int64(prevFees), sizer), nil
}

// BuildTx is a helper function for building a bitcoin transaction. It uses the given `FeeMode` to calculate
// fees. `inputs` will be a list of utxos that guaranteed to be included in the transaction. `utxos` is a list of
// transaction will be picked to cover the output amount and fees. The `recipients` includes a list of target addresses
// and the associated amounts to be sent.  If there's any change, it will be sent back to the `changeAddr`.
// If `changeAddr` is nil, change will be spent as fees.
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

	// Keep adding utxos to make sure it has enough input to cover the output + fees.
	for i := -1; i < len(utxos); i++ {
		// Add the utxo to the transaction input
		if i >= 0 {
			// Skip utxos that are too small
			// Todo : we might want a better way to do this, since sometimes a utxo with (DustAmount +1) could be
			// uneconomical to spend depending on the fee rate.
			if utxos[i].Amount < DustAmount {
				continue
			}
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

	log.Printf("have %v utxos, total in = %v, total out = %v", len(utxos), totalIn, totalOut)

	return nil, fmt.Errorf("funds not enough")
}

// messageToHex serializes a message to the wire protocol encoding using the
// latest protocol version and returns a hex-encoded string of the result.
// Copied from https://github.com/btcsuite/btcd/rpcserver.go and modified
func messageToHex(msg wire.Message) (string, error) {
	maxProtocolVersion := uint32(70002)
	var buf bytes.Buffer
	if err := msg.BtcEncode(&buf, maxProtocolVersion, wire.WitnessEncoding); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf.Bytes()), nil
}

// createVinList returns a slice of JSON objects for the inputs of the passed
// transaction.
// Copied from https://github.com/btcsuite/btcd/rpcserver.go
func createVinList(mtx *wire.MsgTx) []btcjson.Vin {
	// Coinbase transactions only have a single txin by definition.
	vinList := make([]btcjson.Vin, len(mtx.TxIn))
	if blockchain.IsCoinBaseTx(mtx) {
		txIn := mtx.TxIn[0]
		vinList[0].Coinbase = hex.EncodeToString(txIn.SignatureScript)
		vinList[0].Sequence = txIn.Sequence
		vinList[0].Witness = txIn.Witness.ToHexStrings()
		return vinList
	}

	for i, txIn := range mtx.TxIn {
		// The disassembled string will contain [error] inline
		// if the script doesn't fully parse, so ignore the
		// error here.
		disbuf, _ := txscript.DisasmString(txIn.SignatureScript)

		vinEntry := &vinList[i]
		vinEntry.Txid = txIn.PreviousOutPoint.Hash.String()
		vinEntry.Vout = txIn.PreviousOutPoint.Index
		vinEntry.Sequence = txIn.Sequence
		vinEntry.ScriptSig = &btcjson.ScriptSig{
			Asm: disbuf,
			Hex: hex.EncodeToString(txIn.SignatureScript),
		}

		if mtx.HasWitness() {
			vinEntry.Witness = txIn.Witness.ToHexStrings()
		}
	}

	return vinList
}

// createVoutList returns a slice of JSON objects for the outputs of the passed
// transaction.
// Copied from https://github.com/btcsuite/btcd/rpcserver.go
func createVoutList(mtx *wire.MsgTx, chainParams *chaincfg.Params, filterAddrMap map[string]struct{}) []btcjson.Vout {
	voutList := make([]btcjson.Vout, 0, len(mtx.TxOut))
	for i, v := range mtx.TxOut {
		// The disassembled string will contain [error] inline if the
		// script doesn't fully parse, so ignore the error here.
		disbuf, _ := txscript.DisasmString(v.PkScript)

		// Ignore the error here since an error means the script
		// couldn't parse and there is no additional information about
		// it anyways.
		scriptClass, addrs, reqSigs, _ := txscript.ExtractPkScriptAddrs(
			v.PkScript, chainParams)

		// Encode the addresses while checking if the address passes the
		// filter when needed.
		passesFilter := len(filterAddrMap) == 0
		encodedAddrs := make([]string, len(addrs))
		for j, addr := range addrs {
			encodedAddr := addr.EncodeAddress()
			encodedAddrs[j] = encodedAddr

			// No need to check the map again if the filter already
			// passes.
			if passesFilter {
				continue
			}
			if _, exists := filterAddrMap[encodedAddr]; exists {
				passesFilter = true
			}
		}

		if !passesFilter {
			continue
		}

		var vout btcjson.Vout
		vout.N = uint32(i)
		vout.Value = btcutil.Amount(v.Value).ToBTC()
		vout.ScriptPubKey.Addresses = encodedAddrs
		vout.ScriptPubKey.Asm = disbuf
		vout.ScriptPubKey.Hex = hex.EncodeToString(v.PkScript)
		vout.ScriptPubKey.Type = scriptClass.String()
		vout.ScriptPubKey.ReqSigs = int32(reqSigs)

		// Address is defined when there's a single well-defined
		// receiver address. To spend the output a signature for this,
		// and only this, address is required.
		if len(encodedAddrs) == 1 && reqSigs <= 1 {
			vout.ScriptPubKey.Address = encodedAddrs[0]
		}

		voutList = append(voutList, vout)
	}

	return voutList
}

// CreateTxRawResult converts the passed transaction and associated parameters
// to a raw transaction JSON object.
// Copied from https://github.com/btcsuite/btcd/rpcserver.go and modified
func CreateTxRawResult(chainParams *chaincfg.Params, mtx *wire.MsgTx) (*btcjson.TxRawResult, error) {

	mtxHex, err := messageToHex(mtx)
	if err != nil {
		return nil, err
	}

	txReply := &btcjson.TxRawResult{
		Hex:      mtxHex,
		Txid:     mtx.TxHash().String(),
		Hash:     mtx.WitnessHash().String(),
		Size:     int32(mtx.SerializeSize()),
		Vsize:    int32(mempool.GetTxVirtualSize(btcutil.NewTx(mtx))),
		Weight:   int32(blockchain.GetTransactionWeight(btcutil.NewTx(mtx))),
		Vin:      createVinList(mtx),
		Vout:     createVoutList(mtx, chainParams, nil),
		Version:  uint32(mtx.Version),
		LockTime: mtx.LockTime,
	}

	return txReply, nil
}
