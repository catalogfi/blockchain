package localnet

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
	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/btc"
	"github.com/catalogfi/blockchain/evm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
)

func EVMClient() evm.Client {
	client, err := evm.NewClient(evm.Config{RPC: map[string]string{
		string(blockchain.EthereumLocalnet): "http://localhost:8545",
		string(blockchain.ArbitrumLocalnet): "http://localhost:8546",
	}})
	if err != nil {
		panic(fmt.Errorf("failed to create localnet client: %v", err))
	}
	return client
}

func EVMWallet() evm.Wallet {
	key, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		panic(fmt.Errorf("failed to parse localnet privatekey: %v", err))
	}
	return evm.NewWallet(EVMClient(), key)
}

func ETH() blockchain.EVMAsset {
	return blockchain.NewETH(blockchain.NewEvmChain(blockchain.EthereumLocalnet), common.HexToAddress(""))
}

func WBTC() blockchain.EVMAsset {
	return blockchain.NewERC20(blockchain.NewEvmChain(blockchain.EthereumLocalnet), common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"), common.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"))
}

func ArbitrumETH() blockchain.EVMAsset {
	return blockchain.NewETH(blockchain.NewEvmChain(blockchain.ArbitrumLocalnet), common.HexToAddress(""))
}

func ArbitrumWBTC() blockchain.EVMAsset {
	return blockchain.NewERC20(blockchain.NewEvmChain(blockchain.ArbitrumLocalnet), common.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"), common.HexToAddress("0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"))
}

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
	res, err := RunOutput("merry", "faucet", "--to", addr)
	if err != nil {
		return nil, err
	}
	txid := strings.TrimSpace(strings.TrimPrefix(string(res), "Successfully submitted at http://localhost:5000/tx/"))
	color.Green(fmt.Sprintf("Funding address1 %v , txid = %v", addr, txid))
	return chainhash.NewHashFromStr(txid)
}

// MerryNewBlock will mine a new block in the reg testnet. This is usually useful when we need to test something with
// confirmations. It uses the `merry faucet` command to generate a new block, the receiver address is a dummy address
// which shouldn't affect our testing
func MerryNewBlock() error {
	addr := "mwt4FeMsGv6Ua3WrfuhypPtqDUse9CoJev"
	_, err := RunOutput("merry", "faucet", "--to", addr)
	return err
}

// RunOutput the command and catch the output
func RunOutput(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Output()
}
