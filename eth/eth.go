package eth

import (
	"fmt"
	"math/big"

	"github.com/catalogfi/multichain/eth/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Network int

func (net Network) ChainID() *big.Int {
	return big.NewInt(int64(net))
}

const (
	NetworkMainnet (Network) = 1
	NetworkGoerli  (Network) = 5
	NetworkSepolia (Network) = 11155111
)

func SignHash(data []byte) common.Hash {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256Hash([]byte(msg))
}

func InstantWalletInitData(owner, underwriter common.Address, timelock *big.Int) ([]byte, error) {
	contractABI, err := bindings.InstantWalletMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return contractABI.Pack("__InstantWallet_init", owner, underwriter, timelock)
}

func InstantWalletInitCode(factory, implementation common.Address, initData []byte) ([]byte, error) {
	contractABI, err := bindings.InstantWalletFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	code, err := contractABI.Pack("createInstantWallet", implementation, initData)
	if err != nil {
		return nil, err
	}
	return append(factory.Bytes(), code...), nil
}

func ERC1967ProxyAddress(caller, implementation common.Address, initData []byte) (common.Address, error) {
	// salt
	initDataHash := crypto.Keccak256Hash(initData)
	bytes32Ty, err := abi.NewType("bytes32", "bytes32", nil)
	if err != nil {
		return common.Address{}, err
	}
	args := abi.Arguments{
		{
			Type: bytes32Ty,
		},
	}
	packed, err := args.Pack(initDataHash)
	if err != nil {
		return common.Address{}, err
	}
	salt := crypto.Keccak256Hash(packed)

	// Contract init code hash
	contractABI, err := bindings.ERC1967ProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, err
	}
	binBytes, err := hexutil.Decode(bindings.ERC1967ProxyMetaData.Bin)
	if err != nil {
		return common.Address{}, err
	}
	packed, err = contractABI.Pack("", implementation, initData)
	if err != nil {
		return common.Address{}, err
	}
	creationCode := append(binBytes, packed...)
	initCodeHash := crypto.Keccak256Hash(creationCode)

	// Use create2 to calculate the contract address
	return crypto.CreateAddress2(caller, salt, initCodeHash.Bytes()), nil
}
