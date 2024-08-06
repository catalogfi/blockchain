package btc

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/catalogfi/blockchain"
)

type BTCAsset struct {
	address btcutil.Address
	chain   blockchain.Chain
}

func NewBTCAsset(address btcutil.Address, chain blockchain.Chain) BTCAsset {
	return BTCAsset{address: address, chain: chain}
}

func (a BTCAsset) String() string {
	return a.address.EncodeAddress()
}

func (a BTCAsset) Chain() blockchain.Chain {
	return a.chain
}

type HTLCClient interface {
	HTLCEvents(ctx context.Context, asset blockchain.Asset, fromBlock, toBlock uint64) ([]HTLCEvent, error)
}

type htlcClient struct {
	indexer IndexerClient
}

func NewHTLCClient(indexer IndexerClient) HTLCClient {
	return &htlcClient{indexer: indexer}
}

func (client *htlcClient) HTLCEvents(ctx context.Context, asset blockchain.Asset, fromBlock, toBlock uint64) ([]HTLCEvent, error) {
	if toBlock == 0 {
		toBlock = ^uint64(0)
	}

	if fromBlock > toBlock {
		return nil, fmt.Errorf("fromBlock %d is greater than toBlock %d", fromBlock, toBlock)
	}

	assestAddr, err := parseAddress(asset.String())
	if err != nil {
		return nil, err
	}

	txs, err := client.indexer.GetAddressTxs(ctx, assestAddr, "")
	if err != nil {
		return nil, err
	}

	addressStr := assestAddr.EncodeAddress()
	events := make([]HTLCEvent, 0)

	for _, tx := range txs {
		var blockHeight uint64

		if tx.Status.Confirmed {
			blockHeight = *tx.Status.BlockHeight
			if blockHeight < fromBlock || blockHeight > toBlock {
				continue
			}
		}

		for _, VOUT := range tx.VOUTs {
			if VOUT.ScriptPubKeyAddress == addressStr {
				events = append(events, HTLCInitiated{
					ID:                    addressStr,
					InitiateTxBlockNumber: blockHeight,
					InitiateTxHash:        tx.TxID,
					Asset:                 asset,
					Amount:                uint64(VOUT.Value),
				})
			}
		}

		for _, VIN := range tx.VINs {
			if VIN.Prevout.ScriptPubKeyAddress != addressStr {
				continue
			}

			witness := *VIN.Witness
			if len(witness) < 3 {
				continue
			}

			scriptHex := witness[len(witness)-2]
			script, err := hex.DecodeString(scriptHex)
			if err != nil {
				return nil, err
			}

			switch assestAddr.(type) {
			case *btcutil.AddressWitnessScriptHash:
				handleWitnessScriptHashEvents(&events, asset, assestAddr, blockHeight, tx.TxID, witness, script)
			case *btcutil.AddressTaproot:
				handleTaprootEvents(&events, asset, assestAddr, blockHeight, tx.TxID, witness, script)
			}
		}
	}

	return events, nil
}

func handleWitnessScriptHashEvents(events *[]HTLCEvent, asset blockchain.Asset, assestAddr btcutil.Address, blockHeight uint64, txID string, witness []string, script []byte) {
	addressStr := assestAddr.EncodeAddress()
	branch := witness[len(witness)-2]
	if IsHtlc(script) {
		if branch == "01" && len(witness) == 5 {
			*events = append(*events, HTLCRedeemed{
				ID:                  addressStr,
				RedeemTxBlockNumber: blockHeight,
				RedeemTxHash:        txID,
				Asset:               asset,
				Secret:              []byte(witness[2]),
				RedeemerPubkey:      witness[1],
			})
		} else if branch == "" && len(witness) == 4 {
			*events = append(*events, HTLCRefunded{
				ID:                  addressStr,
				RefundTxBlockNumber: blockHeight,
				RefundTxHash:        txID,
				Asset:               asset,
				RefunderPubkey:      witness[1],
			})
		}
	}
}

func handleTaprootEvents(events *[]HTLCEvent, asset blockchain.Asset, assestAddr btcutil.Address, blockHeight uint64, txID string, witness []string, script []byte) {
	addressStr := assestAddr.EncodeAddress()
	if IsRedeemLeaf(script) {
		*events = append(*events, HTLCRedeemed{
			ID:                  addressStr,
			RedeemTxBlockNumber: blockHeight,
			RedeemTxHash:        txID,
			Asset:               asset,
			Secret:              []byte(witness[1]),
			RedeemerPubkey:      string(script),
		})
	} else if IsRefundLeaf(script) || IsMultiSigLeaf(script) {
		*events = append(*events, HTLCRefunded{
			ID:                  addressStr,
			RefundTxBlockNumber: blockHeight,
			RefundTxHash:        txID,
			Asset:               asset,
			RefunderPubkey:      string(script),
		})
	}
}

type HTLCEvent interface {
	OrderID() string
	Equal(e HTLCEvent) bool
	BlockNumber() uint64
	TxHash() string
	BlockchainAsset() blockchain.Asset
}

type HTLCInitiated struct {
	ID                    string
	InitiateTxBlockNumber uint64
	InitiateTxHash        string
	Asset                 blockchain.Asset
	Amount                uint64
}

func (e HTLCInitiated) OrderID() string {
	return e.ID
}

func (e HTLCInitiated) BlockNumber() uint64 {
	return e.InitiateTxBlockNumber
}

func (e HTLCInitiated) TxHash() string {
	return e.InitiateTxHash
}

func (e HTLCInitiated) BlockchainAsset() blockchain.Asset {
	return e.Asset
}

func (e HTLCInitiated) Equal(o HTLCEvent) bool {
	other, ok := o.(HTLCInitiated)
	if !ok {
		return ok
	}

	return e.ID == other.ID &&
		e.Amount == other.Amount &&
		e.InitiateTxHash == other.InitiateTxHash &&
		e.InitiateTxBlockNumber == other.InitiateTxBlockNumber
}

type HTLCRedeemed struct {
	ID                  string
	RedeemTxBlockNumber uint64
	RedeemTxHash        string
	Asset               blockchain.Asset
	Secret              []byte
	RedeemerPubkey      string
}

func (e HTLCRedeemed) OrderID() string {
	return e.ID
}

func (e HTLCRedeemed) BlockNumber() uint64 {
	return e.RedeemTxBlockNumber
}

func (e HTLCRedeemed) TxHash() string {
	return e.RedeemTxHash
}

func (e HTLCRedeemed) BlockchainAsset() blockchain.Asset {
	return e.Asset
}

func (e HTLCRedeemed) Equal(o HTLCEvent) bool {
	other, ok := o.(HTLCRedeemed)
	if !ok {
		return ok
	}

	return e.ID == other.ID &&
		bytes.Equal(e.Secret[:], other.Secret[:]) &&
		e.RedeemTxHash == other.RedeemTxHash &&
		e.RedeemTxBlockNumber == other.RedeemTxBlockNumber
}

type HTLCRefunded struct {
	ID                  string
	RefundTxBlockNumber uint64
	RefundTxHash        string
	Asset               blockchain.Asset
	RefunderPubkey      string
}

func (e HTLCRefunded) OrderID() string {
	return e.ID
}

func (e HTLCRefunded) BlockNumber() uint64 {
	return e.RefundTxBlockNumber
}

func (e HTLCRefunded) TxHash() string {
	return e.RefundTxHash
}

func (e HTLCRefunded) BlockchainAsset() blockchain.Asset {
	return e.Asset
}

func (e HTLCRefunded) Equal(o HTLCEvent) bool {
	other, ok := o.(HTLCRefunded)
	if !ok {
		return ok
	}

	return e.ID == other.ID &&
		e.RefundTxHash == other.RefundTxHash &&
		e.RefundTxBlockNumber == other.RefundTxBlockNumber
}
