package btcswap

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/blockchain/btc"
)

type Swap struct {
	Network    *chaincfg.Params
	Amount     int64
	Secret     []byte
	SecretHash []byte
	WaitBlock  int64
	Address    btcutil.Address
	Initiator  btcutil.Address
	Redeemer   btcutil.Address
	Script     []byte
}

func NewSwap(network *chaincfg.Params, initiatorAddr, redeemer btcutil.Address, amount int64, secretHash []byte, waitBlock int64) (Swap, error) {
	htlc, err := btc.HtlcScript(initiatorAddr.ScriptAddress(), redeemer.ScriptAddress(), secretHash, waitBlock)
	if err != nil {
		return Swap{}, err
	}
	addr, err := btc.P2wshAddress(htlc, network)
	if err != nil {
		return Swap{}, err
	}

	return Swap{
		Network:    network,
		Amount:     amount,
		SecretHash: secretHash,
		WaitBlock:  waitBlock,
		Address:    addr,
		Initiator:  initiatorAddr,
		Redeemer:   redeemer,
		Script:     htlc,
	}, nil
}

func (swap *Swap) IsInitiator(address string) bool {
	return swap.Initiator.EncodeAddress() == address
}

func (swap *Swap) IsRedeemer(address string) bool {
	return swap.Redeemer.EncodeAddress() == address
}

// Initiated returns if the swap has been initiated. It will also return an uint64 which is the block height of the last
// confirmed initiated tx. The swap doesn't have an idea about block confirmations. It will let the caller decide if the
// swap initiation has reached enough confirmation.
func (swap *Swap) Initiated(ctx context.Context, client btc.IndexerClient) (bool, uint64, error) {
	// Fetch all utxos
	utxos, err := client.GetUTXOs(ctx, swap.Address)
	if err != nil {
		return false, 0, fmt.Errorf("failed to get UTXOs: %w", err)
	}

	// Too many utxos will cause us losing many fees
	if len(utxos) > 3 {
		return false, 0, fmt.Errorf("too many partial fill")
	}

	// Check we have enough confirmed utxos (total Amount >= required Amount)
	total, blockHeight := int64(0), uint64(0)
	txs := make([]string, 0, len(utxos))
	for _, utxo := range utxos {
		if utxo.Status != nil && utxo.Status.Confirmed {
			total += utxo.Amount
			txs = append(txs, utxo.TxID)
			if *utxo.Status.BlockHeight > blockHeight {
				blockHeight = *utxo.Status.BlockHeight
			}
		}
	}

	return total >= swap.Amount, blockHeight, nil
}

func (swap *Swap) Initiators(ctx context.Context, client btc.IndexerClient) ([]string, error) {
	// Fetch all utxos
	utxos, err := client.GetUTXOs(ctx, swap.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get UTXOs: %w", err)
	}

	// Check we have enough confirmed utxos (total Amount >= required Amount)
	total, confirmedBlock := int64(0), uint64(0)
	txhashes := make([]string, 0, len(utxos))
	for _, utxo := range utxos {
		if utxo.Status != nil && utxo.Status.Confirmed {
			if *utxo.Status.BlockHeight > confirmedBlock {
				confirmedBlock = *utxo.Status.BlockHeight
			}
			txhashes = append(txhashes, utxo.TxID)
			total += utxo.Amount
		}
	}

	if total >= swap.Amount {
		txSendersMap := map[string]bool{}
		for _, hash := range txhashes {
			rawTx, err := client.GetTx(ctx, hash)
			if err != nil {
				return nil, err
			}
			for _, vin := range rawTx.VINs {
				txSendersMap[vin.Prevout.ScriptPubKeyAddress] = true
			}
		}

		// Convert it to a slice
		txSenders := make([]string, 0, len(txSendersMap))
		for sender := range txSendersMap {
			txSenders = append(txSenders, sender)
		}
		return txSenders, nil
	}
	return nil, nil
}

func (swap *Swap) Redeemed(ctx context.Context, client btc.IndexerClient) (bool, []byte, error) {
	if len(swap.Secret) != 0 {
		return true, swap.Secret, nil
	}

	txs, err := client.GetAddressTxs(ctx, swap.Address, "")
	if err != nil {
		return false, nil, err
	}
	for _, tx := range txs {
		for _, vin := range tx.VINs {
			if vin.Prevout.ScriptPubKeyAddress == swap.Address.EncodeAddress() {
				if len(*vin.Witness) == 5 {
					// witness format
					// [
					//   0 : sig,
					//   1 : spender's public key,
					//   2 : secret,
					//   3 : []byte{0x1},
					//   4 : Script
					// ]
					secretString := (*vin.Witness)[2]
					secretBytes := make([]byte, hex.DecodedLen(len(secretString)))
					_, err := hex.Decode(secretBytes, []byte(secretString))
					if err != nil {
						return false, nil, err
					}

					swap.Secret = secretBytes
					return true, swap.Secret, nil
				}
			}
		}
	}
	return false, nil, nil
}

// Expired returns true if there's any utxo of the swap address has more confirmations than the timelock and is ready
// for refunding.
func (swap *Swap) Expired(ctx context.Context, client btc.IndexerClient) (bool, error) {
	// Fetch all utxos
	utxos, err := client.GetUTXOs(ctx, swap.Address)
	if err != nil {
		return false, fmt.Errorf("failed to get UTXOs: %w", err)
	}

	// Get the number of block height
	current, err := client.GetTipBlockHeight(ctx)
	if err != nil {
		return false, err
	}

	for _, utxo := range utxos {
		if utxo.Status != nil && utxo.Status.Confirmed {
			if current-*utxo.Status.BlockHeight >= uint64(swap.WaitBlock) {
				return true, nil
			}
		}
	}
	return false, nil
}
