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

	childCtx, cancel := context.WithTimeout(ctx, 5)
	defer cancel()

	txs, err := client.indexer.GetAddressTxs(childCtx, assestAddr, "")
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
					id:                    addressStr,
					initiateTxBlockNumber: blockHeight,
					initiateTxHash:        tx.TxID,
					asset:                 asset,
					amount:                uint64(VOUT.Value),
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
				id:                  addressStr,
				redeemTxBlockNumber: blockHeight,
				redeemTxHash:        txID,
				asset:               asset,
				secret:              []byte(witness[2]),
				redeemerPubkey:      witness[1],
			})
		} else if branch == "" && len(witness) == 4 {
			*events = append(*events, HTLCRefunded{
				id:                  addressStr,
				refundTxBlockNumber: blockHeight,
				refundTxHash:        txID,
				asset:               asset,
				refunderPubkey:      witness[1],
			})
		}
	}
}

func handleTaprootEvents(events *[]HTLCEvent, asset blockchain.Asset, assestAddr btcutil.Address, blockHeight uint64, txID string, witness []string, script []byte) {
	addressStr := assestAddr.EncodeAddress()

	ok, redeemerPubKey := IsRedeemLeaf(script)
	if ok {
		s, err := hex.DecodeString(witness[1])
		if err != nil {
			return
		}

		*events = append(*events, HTLCRedeemed{
			id:                  addressStr,
			redeemTxBlockNumber: blockHeight,
			redeemTxHash:        txID,
			asset:               asset,
			secret:              s,
			redeemerPubkey:      redeemerPubKey,
		})
		return
	}

	ok, refunderPubKey := IsRefundLeaf(script)
	if ok {
		*events = append(*events, HTLCRefunded{
			id:                  addressStr,
			refundTxBlockNumber: blockHeight,
			refundTxHash:        txID,
			asset:               asset,
			refunderPubkey:      refunderPubKey,
		})
		return
	}

	ok, refunderPubKey = IsMultiSigLeaf(script)
	if ok {
		*events = append(*events, HTLCRefunded{
			id:                  addressStr,
			refundTxBlockNumber: blockHeight,
			refundTxHash:        txID,
			asset:               asset,
			refunderPubkey:      refunderPubKey,
		})
		return
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
	id                    string
	initiateTxBlockNumber uint64
	initiateTxHash        string
	asset                 blockchain.Asset
	amount                uint64
}

func (e HTLCInitiated) OrderID() string {
	return e.id
}

func (e HTLCInitiated) BlockNumber() uint64 {
	return e.initiateTxBlockNumber
}

func (e HTLCInitiated) TxHash() string {
	return e.initiateTxHash
}

func (e HTLCInitiated) BlockchainAsset() blockchain.Asset {
	return e.asset
}

func (e HTLCInitiated) Equal(o HTLCEvent) bool {
	other, ok := o.(HTLCInitiated)
	if !ok {
		return ok
	}

	return e.id == other.id &&
		e.amount == other.amount &&
		e.initiateTxHash == other.initiateTxHash &&
		e.initiateTxBlockNumber == other.initiateTxBlockNumber
}

func (e HTLCInitiated) Amount() uint64 {
	return e.amount
}

type HTLCRedeemed struct {
	id                  string
	redeemTxBlockNumber uint64
	redeemTxHash        string
	asset               blockchain.Asset
	secret              []byte
	redeemerPubkey      string
}

func (e HTLCRedeemed) OrderID() string {
	return e.id
}

func (e HTLCRedeemed) BlockNumber() uint64 {
	return e.redeemTxBlockNumber
}

func (e HTLCRedeemed) TxHash() string {
	return e.redeemTxHash
}

func (e HTLCRedeemed) BlockchainAsset() blockchain.Asset {
	return e.asset
}

func (e HTLCRedeemed) Equal(o HTLCEvent) bool {
	other, ok := o.(HTLCRedeemed)
	if !ok {
		return ok
	}

	return e.id == other.id &&
		bytes.Equal(e.secret[:], other.secret[:]) &&
		e.redeemTxHash == other.redeemTxHash &&
		e.redeemTxBlockNumber == other.redeemTxBlockNumber
}

func (e HTLCRedeemed) Secret() []byte {
	return e.secret
}

func (e HTLCRedeemed) RedeemerPubkey() string {
	return e.redeemerPubkey
}

type HTLCRefunded struct {
	id                  string
	refundTxBlockNumber uint64
	refundTxHash        string
	asset               blockchain.Asset
	refunderPubkey      string
}

func (e HTLCRefunded) OrderID() string {
	return e.id
}

func (e HTLCRefunded) BlockNumber() uint64 {
	return e.refundTxBlockNumber
}

func (e HTLCRefunded) TxHash() string {
	return e.refundTxHash
}

func (e HTLCRefunded) BlockchainAsset() blockchain.Asset {
	return e.asset
}

func (e HTLCRefunded) Equal(o HTLCEvent) bool {
	other, ok := o.(HTLCRefunded)
	if !ok {
		return ok
	}

	return e.id == other.id &&
		e.refundTxHash == other.refundTxHash &&
		e.refundTxBlockNumber == other.refundTxBlockNumber
}

func (e HTLCRefunded) RefunderPubkey() string {
	return e.refunderPubkey
}
