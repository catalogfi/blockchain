package testutil

import (
	"os"
	"os/exec"
	"strings"
)

// NigiriFaucet funds the given address using the `nigiri faucet` command. It will transfer 1 BTC to the target address
// and automatically generate a new block for the tx. Beware the result will be something like `txId: xxxxx` which is
// not exactly the transaction hash. You'll need to trim the prefix when parsing the hash.
func NigiriFaucet(addr string) (string, error) {
	res, err := RunOutput("nigiri", "faucet", addr)
	if err != nil {
		return "", err
	}
	txid := strings.TrimSpace(strings.TrimPrefix(string(res), "txId:"))
	return txid, nil
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
