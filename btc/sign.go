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
	ignoredIndexes    map[int]struct{}
}

func defaultSigOptions() *sigOptions {
	return &sigOptions{
		compressed:        true,
		sighashType:       txscript.SigHashAll,
		tapScriptRootHash: nil,
		ignoredIndexes:    map[int]struct{}{},
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

func WithIgnoredIndex(indexes ...int) SigOptions {
	return func(o *sigOptions) {
		for _, index := range indexes {
			o.ignoredIndexes[index] = struct{}{}
		}
	}
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

// SignInput signs the input of the `tx` at index `index` with the private key.
func SignInput(addrType waddrmgr.AddressType, tx *wire.MsgTx, index int, key *btcec.PrivateKey, fetcher *txscript.MultiPrevOutFetcher, sigHashes *txscript.TxSigHashes, sigOpts ...SigOptions) error {
	// Parse the options
	opts := defaultSigOptions()
	for _, sigOpt := range sigOpts {
		sigOpt(opts)
	}
	if _, ok := opts.ignoredIndexes[index]; ok {
		return nil
	}

	outpoint := fetcher.FetchPrevOutput(tx.TxIn[index].PreviousOutPoint)
	switch addrType {
	case waddrmgr.PubKeyHash:
		sigScript, err := txscript.SignatureScript(tx, index, outpoint.PkScript, opts.sighashType, key, opts.compressed)
		if err != nil {
			return err
		}
		tx.TxIn[index].SignatureScript = sigScript
	case waddrmgr.WitnessPubKey:
		sig, err := txscript.RawTxInWitnessSignature(tx, sigHashes, index, outpoint.Value, outpoint.PkScript, opts.sighashType, key)
		if err != nil {
			return err
		}
		tx.TxIn[index].Witness = wire.TxWitness{sig, key.PubKey().SerializeCompressed()}
	case waddrmgr.TaprootPubKey:
		sighashType := opts.sighashType
		if sighashType == txscript.SigHashAll {
			sighashType = txscript.SigHashDefault
		}
		sig, err := txscript.RawTxInTaprootSignature(tx, sigHashes, index, outpoint.Value, outpoint.PkScript, opts.tapScriptRootHash, sighashType, key)
		if err != nil {
			return err
		}
		tx.TxIn[index].Witness = wire.TxWitness{sig}
	default:
		return fmt.Errorf("unknown address type: %v", addrType)
	}

	return nil
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

// SignTx signs the entire transaction `tx` with the private key.
func SignTx(addrType waddrmgr.AddressType, tx *wire.MsgTx, key *btcec.PrivateKey, utxos UTXOs, sigOpts ...SigOptions) error {
	opts := defaultSigOptions()
	for _, sigOpt := range sigOpts {
		sigOpt(opts)
	}

	// Calculate the pkScript basing on the addr type
	externalKey := key.PubKey()
	if addrType == waddrmgr.TaprootPubKey {
		externalKey = txscript.ComputeTaprootOutputKey(key.PubKey(), opts.tapScriptRootHash)
	}
	pkScript, err := PkScript(addrType, externalKey)
	if err != nil {
		return err
	}

	// Calculate the sighashes
	fetcher, err := NewFetcher(pkScript, utxos...)
	if err != nil {
		return err
	}
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)

	// Sign each utxo
	for i := range tx.TxIn {
		if err := SignInput(addrType, tx, i, key, fetcher, sigHashes, sigOpts...); err != nil {
			return err
		}
	}
	return nil
}
