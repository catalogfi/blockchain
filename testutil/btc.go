package testutil

import (
	"os"
	"os/exec"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func NigiriStart() error {
	_, err := RunOutput("nigiri", "start")
	return err
}

func NigiriStop() error {
	_, err := RunOutput("nigiri", "stop")
	return err
}

// NigiriFaucet funds the given address using the `nigiri faucet` command. It will transfer 1 BTC to the target address
// and automatically generate a new block for the tx. Beware the result will be something like `txId: xxxxx` which is
// not exactly the transaction hash. You'll need to trim the prefix when parsing the hash.
func NigiriFaucet(addr string) (*chainhash.Hash, error) {
	res, err := RunOutput("nigiri", "faucet", addr)
	if err != nil {
		return nil, err
	}
	txid := strings.TrimSpace(strings.TrimPrefix(string(res), "txId:"))
	return chainhash.NewHashFromStr(txid)
}

// NigiriNewBlock will mine a new block in the reg testnet. This is usually useful when we need to test something with
// confirmations. It uses the `nigiri faucet` command to generate a new block, it shouldn't the things we want to test
// since it's sending funds a random non-related address.
func NigiriNewBlock() error {
	addr := "mwt4FeMsGv6Ua3WrfuhypPtqDUse9CoJev"
	_, err := NigiriFaucet(addr)
	return err
}

// RunOutput the command and catch the output
func RunOutput(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

// RandomBtcAddressP2PKH returns a random p2pkh address of the given network.
func RandomBtcAddressP2PKH(network *chaincfg.Params) (btcutil.Address, error) {
	key, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, err
	}
	return btcutil.NewAddressPubKeyHash(btcutil.Hash160(key.PubKey().SerializeCompressed()), network)
}
