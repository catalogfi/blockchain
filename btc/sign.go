package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
)

type SigOptions func(*sigOptions)

type sigOptions struct {
	compressed        bool
	sighashType       txscript.SigHashType
	tapScriptRootHash []byte
}

func (so *sigOptions) Parse(sigOpts ...SigOptions) {
	for _, sigOpt := range sigOpts {
		sigOpt(so)
	}
}

func defaultSigOptions() *sigOptions {
	return &sigOptions{
		compressed:        true,
		sighashType:       txscript.SigHashDefault,
		tapScriptRootHash: nil,
	}
}

func WithCompressed(compressed bool) SigOptions {
	return func(o *sigOptions) {
		o.compressed = compressed
	}
}

func WithSighashType(sighashType txscript.SigHashType) SigOptions {
	return func(o *sigOptions) {
		o.sighashType = sighashType
	}
}

func WithTapScriptRootHash(hash []byte) SigOptions {
	return func(o *sigOptions) {
		o.tapScriptRootHash = hash
	}
}

// SignFunc defines the function to sign a particular utxo in the tx. It takes a few commonly needed params for signing
// and adds the signature/witness to the tx.
type SignFunc func(tx *wire.MsgTx, index int) error

// SignFuncPubKeyHash returns the SignFunc for a p2pkh utxo.
func SignFuncPubKeyHash(key *btcec.PrivateKey, fetcher txscript.PrevOutputFetcher, sigOpts ...SigOptions) SignFunc {
	opts := defaultSigOptions()
	opts.Parse(sigOpts...)

	return func(tx *wire.MsgTx, index int) error {
		outpoint := fetcher.FetchPrevOutput(tx.TxIn[index].PreviousOutPoint)
		if outpoint == nil {
			return fmt.Errorf("no output found for txid: %s", tx.TxHash().String())
		}
		if opts.sighashType == txscript.SigHashDefault {
			opts.sighashType = txscript.SigHashAll
		}
		sigScript, err := txscript.SignatureScript(tx, index, outpoint.PkScript, opts.sighashType, key, opts.compressed)
		if err != nil {
			return err
		}

		tx.TxIn[index].SignatureScript = sigScript
		return nil
	}
}

// SignFuncWitnessPubKey returns the SignFunc for a p2wpkh utxo.
func SignFuncWitnessPubKey(key *btcec.PrivateKey, fetcher txscript.PrevOutputFetcher, sigHashes *txscript.TxSigHashes, sigOpts ...SigOptions) SignFunc {
	opts := defaultSigOptions()
	opts.Parse(sigOpts...)

	return func(tx *wire.MsgTx, index int) error {
		outpoint := fetcher.FetchPrevOutput(tx.TxIn[index].PreviousOutPoint)
		if outpoint == nil {
			return fmt.Errorf("no output found for txid: %s", tx.TxHash().String())
		}
		if opts.sighashType == txscript.SigHashDefault {
			opts.sighashType = txscript.SigHashAll
		}
		sig, err := txscript.RawTxInWitnessSignature(tx, sigHashes, index, outpoint.Value, outpoint.PkScript, opts.sighashType, key)
		if err != nil {
			return err
		}
		tx.TxIn[index].Witness = wire.TxWitness{sig, key.PubKey().SerializeCompressed()}
		return nil
	}
}

// SignFuncTaprootPublicKey returns the SignFunc for a p2tr utxo.
func SignFuncTaprootPublicKey(key *btcec.PrivateKey, fetcher txscript.PrevOutputFetcher, sigHashes *txscript.TxSigHashes, sigOpts ...SigOptions) SignFunc {
	opts := defaultSigOptions()
	opts.Parse(sigOpts...)

	return func(tx *wire.MsgTx, index int) error {
		outpoint := fetcher.FetchPrevOutput(tx.TxIn[index].PreviousOutPoint)
		if outpoint == nil {
			return fmt.Errorf("no output found for txid: %s", tx.TxHash().String())
		}
		sig, err := txscript.RawTxInTaprootSignature(tx, sigHashes, index, outpoint.Value, outpoint.PkScript, opts.tapScriptRootHash, opts.sighashType, key)
		if err != nil {
			return err
		}
		tx.TxIn[index].Witness = wire.TxWitness{sig}
		return nil
	}
}

// SignMap is a few types can be accepted by the Sign function
type SignMap interface {
	SignFunc | map[int]SignFunc | map[string]SignFunc
}

// SignTx will sign the entire `tx` basing on the sign method defined by the `signMap`.
func SignTx[K SignMap](tx *wire.MsgTx, signMap K) error {
	signMapAsAny := any(signMap)
	switch s := signMapAsAny.(type) {

	// All inputs will be signed using the same SignFunc
	case SignFunc:
		for i := range tx.TxIn {
			if err := s(tx, i); err != nil {
				return err
			}
		}
	// The SignFunc is mapped by the utxo's index
	case map[int]SignFunc:
		for i := range tx.TxIn {
			sf := s[i]
			if err := sf(tx, i); err != nil {
				return err
			}
		}
	// The SignFunc is mapped by the utxo's string
	case map[string]SignFunc:
		for i := range tx.TxIn {
			sf := s[tx.TxIn[i].PreviousOutPoint.String()]
			if err := sf(tx, i); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown sign map type: %T", s)
	}
	return nil
}

// QuickSign can be used to sign some simple tx with known script type.
func QuickSign(addrType waddrmgr.AddressType, tx *wire.MsgTx, key *btcec.PrivateKey, utxos []UTXO, sigOpts ...SigOptions) error {
	opts := defaultSigOptions()
	opts.Parse(sigOpts...)

	// Calculate the pkScript basing on the addr type
	externalKey := key.PubKey()
	if addrType == waddrmgr.TaprootPubKey {
		externalKey = txscript.ComputeTaprootOutputKey(key.PubKey(), opts.tapScriptRootHash)
	}
	pkScript, err := PkScript(addrType, externalKey)
	if err != nil {
		return err
	}

	// Initiate a new fetcher
	fetcher, err := NewFetcher(pkScript, utxos...)
	if err != nil {
		return err
	}

	var sf SignFunc
	switch addrType {
	case waddrmgr.PubKeyHash:
		sf = SignFuncPubKeyHash(key, fetcher, sigOpts...)
	case waddrmgr.WitnessPubKey:
		sigHashes := txscript.NewTxSigHashes(tx, fetcher)
		sf = SignFuncWitnessPubKey(key, fetcher, sigHashes, sigOpts...)
	case waddrmgr.TaprootPubKey:
		sigHashes := txscript.NewTxSigHashes(tx, fetcher)
		sf = SignFuncTaprootPublicKey(key, fetcher, sigHashes, sigOpts...)
	default:
		return fmt.Errorf("unsupported address type: %v", addrType)
	}

	return SignTx(tx, sf)
}

// PayToPubKeyHashScript creates a new script to pay a transaction
// output to a 20-byte pubkey hash. It is expected that the input is a valid
// hash.
func PayToPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
		AddData(pubKeyHash).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).
		Script()
}

// PayToWitnessPubKeyHashScript creates a new script to pay to a version 0
// pubkey hash witness program. The passed hash is expected to be valid.
func PayToWitnessPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
}

// PayToScriptHashScript creates a new script to pay a transaction output to a
// script hash. It is expected that the input is a valid hash.
func PayToScriptHashScript(scriptHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_HASH160).AddData(scriptHash).
		AddOp(txscript.OP_EQUAL).Script()
}

// PayToWitnessScriptHashScript creates a new script to pay to a version 0
// script hash witness program. The passed hash is expected to be valid.
func PayToWitnessScriptHashScript(scriptHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(scriptHash).Script()
}

// PayToWitnessTaprootScript creates a new script to pay to a version 1
// (taproot) witness program. The passed hash is expected to be valid.
func PayToWitnessTaprootScript(rawKey []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_1).AddData(rawKey).Script()
}

// PkScript returns the pkScript basing on the addr type.
func PkScript(addrType waddrmgr.AddressType, key *btcec.PublicKey) ([]byte, error) {
	switch addrType {
	case waddrmgr.PubKeyHash:
		return PayToPubKeyHashScript(btcutil.Hash160(key.SerializeCompressed()))
	case waddrmgr.WitnessPubKey:
		return PayToWitnessPubKeyHashScript(btcutil.Hash160(key.SerializeCompressed()))
	case waddrmgr.TaprootPubKey:
		return PayToWitnessTaprootScript(schnorr.SerializePubKey(key))
	default:
		return nil, fmt.Errorf("unknown address type: %v", addrType)
	}
}
