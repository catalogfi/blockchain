package btc

import (
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/catalogfi/blockchain"
)

type HTLCWallet interface {
	Initiate(chain blockchain.UtxoChain, initiator, redeemer *btcec.PublicKey, secretHash [32]byte, expiry uint64, amount *big.Int) (string, error)
	Redeem(chain blockchain.UtxoChain, initiator, redeemer *btcec.PublicKey, expiry uint64, amount *big.Int, secret []byte) (string, error)
	Refund(chain blockchain.UtxoChain, initiator, redeemer *btcec.PublicKey, secretHash [32]byte, expiry uint64, amount *big.Int, signature []byte) (string, error)
}

type htlcWalletBatcher struct {
	batcher Batcher
}

func (w *htlcWalletBatcher) Initiate(chain blockchain.UtxoChain, initiator, redeemer *btcec.PublicKey, secretHash [32]byte, expiry uint64, amount *big.Int) (string, error) {
	// build script
	return "", nil
}

type Batcher interface {
	AddOutput(r Recipient) (string, error)
	AddInput()
}
