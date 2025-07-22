package cosigner

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/catalogfi/blockchain/btc"
)

var DefaultTimelock = int64(144 * 180)

var (
	BaseSizeSpend  = 0
	BaseSizeRefund = 0

	// SegwitSizeSpend is the estimated size of the witness when spending the cosigner script
	// Stack number + stack size * 5 + 1st stack + 2nd stack(sig) + 3rd stack(sig) + 0x1 + script size
	SegwitSizeSpend = 1 + 5 + 0 + 72 + 72 + 1 + 114

	SegwitSizeRefund = func(timelock int64) int {
		timelockSize := 0
		switch {
		case timelock <= 16:
			timelockSize = 1
		case timelock < 128:
			timelockSize = 2
		case timelock < 32768:
			timelockSize = 3
		default:
			timelockSize = 4
		}
		// Stack number + stack size * 3 + 1st stack(sig) + 2nd stack(empty) + 3rd stack(script)
		return 1 + 3 + 72 + 0 + (111 + timelockSize)
	}
)

func Script(cosignerPub, userPub []byte, timelock int64) ([]byte, error) {
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_IF).
		AddOp(txscript.OP_2).
		AddData(cosignerPub).
		AddData(userPub).
		AddOp(txscript.OP_2).
		AddOp(txscript.OP_CHECKMULTISIG).
		AddOp(txscript.OP_ELSE).
		AddInt64(timelock).
		AddOp(txscript.OP_CHECKSEQUENCEVERIFY).
		AddOp(txscript.OP_DROP).
		AddData(userPub).
		AddOp(txscript.OP_CHECKSIG).
		AddOp(txscript.OP_ENDIF).
		Script()
}

func Sign(script []byte, tx *wire.MsgTx, index int, fetcher txscript.PrevOutputFetcher, key *btcec.PrivateKey) ([]byte, error) {
	outpoint := fetcher.FetchPrevOutput(tx.TxIn[index].PreviousOutPoint)
	return txscript.RawTxInWitnessSignature(tx, txscript.NewTxSigHashes(tx, fetcher), index, outpoint.Value, script, txscript.SigHashAll, key)
}

func Witness(script, cosignerSig, userSig []byte, refund bool) wire.TxWitness {
	var witnessStack wire.TxWitness
	if refund {
		witnessStack = make(wire.TxWitness, 3)
		witnessStack[0] = userSig
		witnessStack[1] = nil
		witnessStack[2] = script
	} else {
		witnessStack = make(wire.TxWitness, 5)
		witnessStack[0] = nil
		witnessStack[1] = cosignerSig
		witnessStack[2] = userSig
		witnessStack[3] = []byte{0x1}
		witnessStack[4] = script
	}

	return witnessStack
}

type SpendSigner struct {
	opts      *btc.SigOptions
	key       *btcec.PrivateKey
	script    []byte
	othersSig []byte
	user      bool
}

func NewSpendSigner(key *btcec.PrivateKey, script, otherSig []byte, user bool, sigOpts ...btc.SigOption) btc.Signer {
	opts := btc.DefaultSigOptions()
	opts.Parse(sigOpts...)

	return &SpendSigner{
		opts:      opts,
		key:       key,
		script:    script,
		othersSig: otherSig,
		user:      user,
	}
}

func (signer *SpendSigner) Sign(tx *wire.MsgTx, index int, outpoint *wire.TxOut, sigHashes *txscript.TxSigHashes) error {
	sig, err := txscript.RawTxInWitnessSignature(tx, sigHashes, index, outpoint.Value, outpoint.PkScript, signer.opts.SigHashType, signer.key)
	if err != nil {
		return err
	}
	coSignerSig, userSig := sig, signer.othersSig
	if signer.user {
		coSignerSig, userSig = signer.othersSig, sig
	}

	witnessStack := make(wire.TxWitness, 5)
	witnessStack[0] = nil
	witnessStack[1] = coSignerSig
	witnessStack[2] = userSig
	witnessStack[3] = []byte{0x1}
	witnessStack[4] = signer.script
	tx.TxIn[index].Witness = witnessStack

	return nil
}

func (signer *SpendSigner) SigSize() (int, int) {
	return BaseSizeSpend, SegwitSizeSpend
}

type RefundSigner struct {
	opts     *btc.SigOptions
	key      *btcec.PrivateKey
	script   []byte
	timelock int64
}

func NewRefundSigner(key *btcec.PrivateKey, script []byte, timelock int64, sigOpts ...btc.SigOption) btc.Signer {
	opts := btc.DefaultSigOptions()
	opts.Parse(sigOpts...)

	return &RefundSigner{
		opts:     opts,
		key:      key,
		script:   script,
		timelock: timelock,
	}
}

func (signer *RefundSigner) Sign(tx *wire.MsgTx, index int, outpoint *wire.TxOut, sigHashes *txscript.TxSigHashes) error {
	sig, err := txscript.RawTxInWitnessSignature(tx, sigHashes, index, outpoint.Value, outpoint.PkScript, signer.opts.SigHashType, signer.key)
	if err != nil {
		return err
	}
	witnessStack := make(wire.TxWitness, 3)
	witnessStack[0] = sig
	witnessStack[1] = nil
	witnessStack[2] = signer.script
	tx.TxIn[index].Witness = witnessStack

	return nil

}

func (signer *RefundSigner) SigSize() (int, int) {
	return BaseSizeRefund, SegwitSizeRefund(signer.timelock)
}
