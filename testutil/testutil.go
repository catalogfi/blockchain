package testutil

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

// RandomSecret creates a random secret with size [0,32)
func RandomSecret() []byte {
	length := rand.Intn(32)
	data := make([]byte, length)

	_, err := crand.Read(data)
	if err != nil {
		panic(err)
	}
	return data
}

func RandomString() string {
	id := make([]byte, 32)

	if _, err := io.ReadFull(crand.Reader, id); err != nil {
		panic(err) // This shouldn't happen
	}
	return hex.EncodeToString(id)
}

// RandomHash creates a random hash of 32 bytes
func RandomHash() [32]byte {
	secret := RandomSecret()
	return sha256.Sum256(secret)
}

// RandomHashString returns the hex encoded string of a random hash
func RandomHashString() string {
	rHash := RandomHash()
	hash, err := chainhash.NewHash(rHash[:])
	if err != nil {
		panic(err)
	}
	return hash.String()
}

func ParseKeys(name string, params *chaincfg.Params) (*btcec.PrivateKey, *btcec.PublicKey, btcutil.Address, error) {
	keyStr := os.Getenv(name)
	privKeyBytes, err := hex.DecodeString(keyStr)
	if err != nil {
		return nil, nil, nil, err
	}
	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes)
	addr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), params)
	if err != nil {
		return nil, nil, nil, err
	}
	return privKey, privKey.PubKey(), addr, nil
}

func ParseStringEnv(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		if defaultValue == "" {
			panic(fmt.Errorf("%v is not set", name))
		} else {
			return defaultValue
		}
	}
	return value
}
