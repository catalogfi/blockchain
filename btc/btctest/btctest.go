package btctest

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/catalogfi/blockchain/btc"
	"github.com/fatih/color"
)

const (
	DefaultRegtestHost = "0.0.0.0:18443"
)

// NewBtcKey generates a new bitcoin private key.
func NewBtcKey(network *chaincfg.Params, addrType waddrmgr.AddressType) (*btcec.PrivateKey, btcutil.Address, error) {
	key, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	addr, err := btc.PublicKeyAddress(network, addrType, key.PubKey())
	if err != nil {
		return nil, nil, err
	}
	return key, addr, nil
}

// RandomSecret creates a random secret with size [1,32)
func RandomSecret() []byte {
	length := rand.Intn(31) + 1
	data := make([]byte, length)

	_, err := crand.Read(data)
	if err != nil {
		panic(err)
	}
	return data
}

// MerryFaucet funds the given address using the `merry faucet` command. It will transfer 1 BTC to the target address
// and automatically generate a new block for the tx. It returns the txid of the funding transaction.
func MerryFaucet(addr string) (*chainhash.Hash, error) {
	res, err := RunOutput("merry", "faucet", addr)
	if err != nil {
		return nil, err
	}
	txid := strings.TrimSpace(strings.TrimPrefix(string(res), "txId:"))
	color.Green(fmt.Sprintf("Funding address1 %v , txid = %v", addr, txid))
	return chainhash.NewHashFromStr(txid)
}

// MerryNewBlock will mine a new block in the reg testnet. This is usually useful when we need to test something with
// confirmations. It uses the `merry faucet` command to generate a new block, the receiver address is a dummy address
// which shouldn't affect our testing
func MerryNewBlock() error {
	addr := "mwt4FeMsGv6Ua3WrfuhypPtqDUse9CoJev"
	_, err := RunOutput("merry", "faucet", addr)
	return err
}

// RunOutput the command and catch the output
func RunOutput(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Output()
}
