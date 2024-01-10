package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

type AddressType string

var (
	AddressP2PK    AddressType = "P2PK"
	AddressP2PKH   AddressType = "P2PKH"
	AddressP2SH    AddressType = "P2SH"
	AddressP2WPKH  AddressType = "P2WPKH"
	AddressSegwit  AddressType = "Bech32"
	AddressTaproot AddressType = "P2TR"
)

func PublicKeyAddress(pub btcec.PublicKey, network *chaincfg.Params, addressType AddressType) (btcutil.Address, error) {
	switch addressType {
	case AddressP2PK:
		return btcutil.NewAddressPubKey(pub.SerializeCompressed(), network)
	case AddressP2PKH:
		return btcutil.NewAddressPubKeyHash(btcutil.Hash160(pub.SerializeCompressed()), network)
	case AddressP2WPKH:
		return btcutil.NewAddressWitnessScriptHash(btcutil.Hash160(pub.SerializeCompressed()), network)
	case AddressTaproot:
		// todo: add taproot address
		panic("todo")
	default:
		return nil, fmt.Errorf("unsupported address type")
	}
}
