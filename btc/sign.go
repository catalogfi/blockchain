package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
)

var (
	// BaseSizeP2PKH is the worst case (largest) serialize size of a transaction input script that redeems a compressed
	// P2PKH output. This assumes we always use a low-s value signature. The signature size are usually 71 (60%) or
	// 72(40%), with a very small chance of 70 or 69 (<1%). It is calculated as :
	// sigLength(1) + sig(72) + pubKeyLength(1) + compressedPubKey(33)
	BaseSizeP2PKH = 1 + 72 + 1 + 33

	BaseSizeP2WPKH = 0

	BaseSizeP2TR = 0

	SegwitSizeP2PKH = 0

	// SegwitSizeP2WPKH is the worst case weight of a witness for spending P2WPKH outputs. It is calculated as :
	// number of items(1) + sigLength(1) + sig(72) + pubKeyLength(1) + compressedPubKey(33)
	SegwitSizeP2WPKH = 1 + 1 + 72 + 1 + 33

	// SegwitSizeP2TR is the worst case weight of a witness for spending P2TR outputs. It is calculated as :
	// number of items(1) + sigLength(1) + sig(64)
	SegwitSizeP2TR = 1 + 1 + 64

	// SegwitSizeP2trDefault is the witness size when the schnorr signature is signed using the default sighash flag.
	SegwitSizeP2trDefault = 1 + 1 + 64
)

type SigOptions func(*sigOptions)

type sigOptions struct {
	compressed        bool
	sigHashes         *txscript.TxSigHashes
	sigHashType       txscript.SigHashType
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
		sigHashType:       txscript.SigHashDefault,
		tapScriptRootHash: nil,
	}
}

func WithCompressed(compressed bool) SigOptions {
	return func(o *sigOptions) {
		o.compressed = compressed
	}
}

func WithSighashes(sighashes *txscript.TxSigHashes) SigOptions {
	return func(o *sigOptions) {
		o.sigHashes = sighashes
	}
}

func WithSighashType(sighashType txscript.SigHashType) SigOptions {
	return func(o *sigOptions) {
		o.sigHashType = sighashType
	}
}

func WithTapScriptRootHash(hash []byte) SigOptions {
	return func(o *sigOptions) {
		o.tapScriptRootHash = hash
	}
}

type Signers interface {
	AddUtxo(signer Signer, utxos ...UTXO)

	Sign(tx *wire.MsgTx) error

	EstimateTxWeight(tx *wire.MsgTx) (int, error)

	EstimateTxVirtualSize(tx *wire.MsgTx) (int, error)
}

type signers struct {
	fetcher *txscript.MultiPrevOutFetcher
	signers map[string]Signer
}

func NewSigners(signer Signer, utxos ...UTXO) (Signers, error) {
	sm := map[string]Signer{}
	fetcher := txscript.NewMultiPrevOutFetcher(nil)

	for _, utxo := range utxos {
		if len(utxo.PkScript) == 0 {
			return nil, fmt.Errorf("utxo %v has no pkscript", utxo)
		}

		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}
		fetcher.AddPrevOut(wire.OutPoint{
			Hash:  *hash,
			Index: utxo.Vout,
		}, wire.NewTxOut(utxo.Amount, utxo.PkScript))
		sm[utxo.String()] = signer
	}

	return &signers{
		fetcher: fetcher,
		signers: sm,
	}, nil
}

func NewSignersByAddrType(addrType waddrmgr.AddressType, key *btcec.PrivateKey, utxos []UTXO, sigOpts ...SigOptions) (Signers, Signer, error) {
	opts := defaultSigOptions()
	opts.Parse(sigOpts...)

	var signer Signer
	switch addrType {
	case waddrmgr.PubKeyHash:
		signer = NewP2pkhSigner(key, sigOpts...)
	case waddrmgr.WitnessPubKey:
		signer = NewP2wpkSigner(key, sigOpts...)
	case waddrmgr.TaprootPubKey:
		signer = NewP2trSigner(key, sigOpts...)
	default:
		return nil, nil, fmt.Errorf("unsupported address type = %v", addrType)
	}
	s, err := NewSigners(signer, utxos...)
	if err != nil {
		return nil, nil, err
	}
	return s, signer, nil
}

func (signers signers) AddUtxo(signer Signer, utxos ...UTXO) {
	for _, utxo := range utxos {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return
		}
		signers.fetcher.AddPrevOut(wire.OutPoint{
			Hash:  *hash,
			Index: utxo.Vout,
		}, wire.NewTxOut(utxo.Amount, utxo.PkScript))
		signers.signers[utxo.String()] = signer
	}
}

func (signers signers) Sign(tx *wire.MsgTx) error {
	// Make sure we have the txOut info for all the inputs
	for _, input := range tx.TxIn {
		outpiont := signers.fetcher.FetchPrevOutput(input.PreviousOutPoint)
		if outpiont == nil {
			return fmt.Errorf("no output found for txid=%v", tx.TxHash())
		}
	}
	sigHashes := txscript.NewTxSigHashes(tx, signers.fetcher)

	for i, input := range tx.TxIn {
		// Make sure we have the txOut info for all the inputs
		outpoint := signers.fetcher.FetchPrevOutput(input.PreviousOutPoint)
		if outpoint == nil {
			return fmt.Errorf("no output found for txid=%v", tx.TxHash())
		}

		// Sign the input
		signer, ok := signers.signers[input.PreviousOutPoint.String()]
		if !ok {
			return fmt.Errorf("Signers.Sign: don't know how to sign %v", input.PreviousOutPoint.String())
		}
		if err := signer.Sign(tx, i, outpoint, sigHashes); err != nil {
			return err
		}
	}
	return nil
}

func (signers signers) EstimateTxWeight(tx *wire.MsgTx) (int, error) {
	totalBase, totalSegwit := tx.SerializeSizeStripped(), 0
	legacy := 0
	for _, input := range tx.TxIn {
		key := input.PreviousOutPoint.String()
		base, segwit := signers.SigSize(key)
		if base == 0 && segwit == 0 {
			return 0, fmt.Errorf("unknown utxo = %v", key)
		}
		totalBase += base
		totalSegwit += segwit
		if base != 0 {
			legacy++
		}
	}

	// Additional 2 weight units for segwit marker + flag if tx has any witness input
	if totalSegwit > 0 {
		totalSegwit += 2
	}

	// When including both legacy and segwit inputs
	if totalSegwit > 0 && legacy > 0 {
		totalSegwit += legacy
	}

	return totalBase*4 + totalSegwit, nil
}

func (signers signers) EstimateTxVirtualSize(tx *wire.MsgTx) (int, error) {
	weight, err := signers.EstimateTxWeight(tx)
	if err != nil {
		return 0, err
	}
	return (weight + 3) / blockchain.WitnessScaleFactor, nil
}

func (signers signers) SigSize(id string) (int, int) {
	signer, ok := signers.signers[id]
	if !ok {
		return 0, 0
	}
	return signer.SigSize()
}

// Signer specify a particular way to spend the script.
type Signer interface {
	Sign(tx *wire.MsgTx, index int, outpoint *wire.TxOut, sigHashes *txscript.TxSigHashes) error

	SigSize() (int, int)
}

type P2pkhSigner struct {
	opts *sigOptions
	key  *btcec.PrivateKey
}

func NewP2pkhSigner(key *btcec.PrivateKey, sigOpts ...SigOptions) Signer {
	opts := defaultSigOptions()
	opts.Parse(sigOpts...)

	return &P2pkhSigner{
		opts: opts,
		key:  key,
	}
}

func (signer *P2pkhSigner) Sign(tx *wire.MsgTx, index int, outpoint *wire.TxOut, sigHashes *txscript.TxSigHashes) error {
	if signer.opts.sigHashType == txscript.SigHashDefault {
		signer.opts.sigHashType = txscript.SigHashAll
	}
	sigScript, err := txscript.SignatureScript(tx, index, outpoint.PkScript, signer.opts.sigHashType, signer.key, signer.opts.compressed)
	if err != nil {
		return err
	}
	tx.TxIn[index].SignatureScript = sigScript

	return nil
}

func (signer *P2pkhSigner) SigSize() (int, int) {
	return BaseSizeP2PKH, SegwitSizeP2PKH
}

type P2wpkSigner struct {
	opts *sigOptions
	key  *btcec.PrivateKey
}

func NewP2wpkSigner(key *btcec.PrivateKey, sigOpts ...SigOptions) Signer {
	opts := defaultSigOptions()
	opts.Parse(sigOpts...)

	return &P2wpkSigner{
		opts: opts,
		key:  key,
	}
}

func (signer *P2wpkSigner) Sign(tx *wire.MsgTx, index int, outpoint *wire.TxOut, sigHashes *txscript.TxSigHashes) error {
	if signer.opts.sigHashType == txscript.SigHashDefault {
		signer.opts.sigHashType = txscript.SigHashAll
	}

	sig, err := txscript.RawTxInWitnessSignature(tx, sigHashes, index, outpoint.Value, outpoint.PkScript, signer.opts.sigHashType, signer.key)
	if err != nil {
		return err
	}
	tx.TxIn[index].Witness = wire.TxWitness{sig, signer.key.PubKey().SerializeCompressed()}
	return nil
}

func (signer *P2wpkSigner) SigSize() (int, int) {
	return BaseSizeP2WPKH, SegwitSizeP2WPKH
}

type P2trSigner struct {
	opts *sigOptions
	key  *btcec.PrivateKey
}

func NewP2trSigner(key *btcec.PrivateKey, sigOpts ...SigOptions) Signer {
	opts := defaultSigOptions()
	opts.Parse(sigOpts...)

	return &P2trSigner{
		opts: opts,
		key:  key,
	}
}

func (signer *P2trSigner) Sign(tx *wire.MsgTx, index int, outpoint *wire.TxOut, sigHashes *txscript.TxSigHashes) error {
	sig, err := txscript.RawTxInTaprootSignature(tx, sigHashes, index, outpoint.Value, outpoint.PkScript, signer.opts.tapScriptRootHash, signer.opts.sigHashType, signer.key)
	if err != nil {
		return err
	}
	tx.TxIn[index].Witness = wire.TxWitness{sig}
	return nil
}

func (signer *P2trSigner) SigSize() (int, int) {
	return BaseSizeP2TR, SegwitSizeP2TR
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
