package localnet

import (
	"crypto/ecdsa"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

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
	"go.uber.org/zap"
)

var localnetKeys = []string{
	"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
	"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
	"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
	"7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
	"47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a",
	"8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba",
	"92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e",
	"4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356",
	"dbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97",
	"2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6",
}

var EVMLocalnetConfig = evm.Config{RPC: map[string]string{
	string(blockchain.EthereumLocalnet): "http://localhost:8545",
	string(blockchain.ArbitrumLocalnet): "http://localhost:8546",
}}

func ECDSAKey(i uint) *ecdsa.PrivateKey {
	key, err := crypto.HexToECDSA(localnetKeys[i])
	if err != nil {
		panic(err)
	}
	return key
}

func BTCECKey(i uint) *btcec.PrivateKey {
	keyBytes, err := hex.DecodeString(localnetKeys[i])
	if err != nil {
		panic(err)
	}
	privKey, _ := btcec.PrivKeyFromBytes(keyBytes)
	return privKey
}

func EVMHTLCWallet(i uint) (evm.HTLCWallet, error) {
	client, err := EVMHTLCClient()
	if err != nil {
		return nil, err
	}
	return evm.NewHTLCWallet(client, ECDSAKey(i)), nil
}

func EVMHTLCClient() (evm.HTLCClient, error) {
	return evm.NewHTLCClient(EVMLocalnetConfig)
}

func EVMClient() (evm.Client, error) {
	return evm.NewClient(EVMLocalnetConfig)
}

func EVMWallet(i uint) (evm.Wallet, error) {
	client, err := EVMClient()
	if err != nil {
		return nil, err
	}
	return evm.NewWallet(client, ECDSAKey(i)), nil
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

func BTCIndexer() btc.IndexerClient {
	return btc.NewElectrsIndexerClient(zap.NewNop(), "http://127.0.0.1:30000", 5*time.Second)
}

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

// FundBTC funds the given address using the `merry faucet` command. It will transfer 1 BTC to the target address
// and automatically generate a new block for the tx. It returns the txid of the funding transaction.
func FundBTC(addr string) (*chainhash.Hash, error) {
	res, err := RunOutput("merry", "faucet", "--to", addr)
	if err != nil {
		return nil, err
	}
	txid := strings.TrimSpace(strings.TrimPrefix(string(res), "Successfully submitted at http://localhost:5050/tx/"))
	color.Green(fmt.Sprintf("Funding address1 %v , txid = %v", addr, txid))
	return chainhash.NewHashFromStr(txid)
}

// MineBTCBlock will mine a new block in the reg testnet. This is usually useful when we need to test something with
// confirmations. It uses the `merry faucet` command to generate a new block, the receiver address is a dummy address
// which shouldn't affect our testing
func MineBTCBlock() error {
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
