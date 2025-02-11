package evm

import (
	"crypto/sha256"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Htlc struct {
	ID         [32]byte
	Contract   common.Address
	Initiator  common.Address
	Redeemer   common.Address
	SecretHash common.Hash
	Secret     []byte
	Amount     *big.Int
	Expiry     *big.Int
	ChainID    *big.Int
}

func NewHtlc(initiator, redeemer, contract common.Address, secretHash common.Hash, amount, expiry, chainID *big.Int) Htlc {
	chainIDPadded := leftPadBytes(chainID.Bytes(), 32)
	secretHashPadded := rightPadBytes(secretHash.Bytes(), 32)
	initiatorPadded := rightPadBytes(common.HexToHash(initiator.String()).Bytes(), 32)
	data := append(chainIDPadded, secretHashPadded...)
	data = append(data, initiatorPadded...)
	id := sha256.Sum256(data)

	return Htlc{
		ID:         id,
		Contract:   contract,
		Initiator:  initiator,
		Redeemer:   redeemer,
		SecretHash: secretHash,
		Amount:     amount,
		Expiry:     expiry,
		ChainID:    chainID,
	}
}

func leftPadBytes(slice []byte, size int) []byte {
	if len(slice) >= size {
		return slice
	}
	padding := make([]byte, size-len(slice))
	return append(padding, slice...)
}

func rightPadBytes(slice []byte, size int) []byte {
	if len(slice) >= size {
		return slice
	}
	padding := make([]byte, size-len(slice))
	return append(slice, padding...)
}
