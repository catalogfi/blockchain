// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gardenfeeaccount

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// GardenFEEAccountHTLC is an auto generated low-level Go binding around an user-defined struct.
type GardenFEEAccountHTLC struct {
	SecretHash    [32]byte
	TimeLock      *big.Int
	SendAmount    *big.Int
	RecieveAmount *big.Int
}

// GardenFEEAccountMetaData contains all meta data concerning the GardenFEEAccount contract.
var GardenFEEAccountMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"token_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"funder_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"feeAccountName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"feeAccountVersion\",\"type\":\"string\"}],\"name\":\"__GardenFEEAccount_init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"amount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce_\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timeLock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sendAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recieveAmount\",\"type\":\"uint256\"}],\"internalType\":\"structGardenFEEAccount.HTLC[]\",\"name\":\"htlcs\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"secrets\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"funderSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"recipientSig\",\"type\":\"bytes\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce_\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timeLock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sendAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recieveAmount\",\"type\":\"uint256\"}],\"internalType\":\"structGardenFEEAccount.HTLC[]\",\"name\":\"htlcs\",\"type\":\"tuple[]\"}],\"name\":\"claimHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"funderSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"recipientSig\",\"type\":\"bytes\"}],\"name\":\"close\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"expiration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"contractIGardenFEEAccountFactory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"funder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recipient\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"secretsProvided\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"settle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// GardenFEEAccountABI is the input ABI used to generate the binding from.
// Deprecated: Use GardenFEEAccountMetaData.ABI instead.
var GardenFEEAccountABI = GardenFEEAccountMetaData.ABI

// GardenFEEAccount is an auto generated Go binding around an Ethereum contract.
type GardenFEEAccount struct {
	GardenFEEAccountCaller     // Read-only binding to the contract
	GardenFEEAccountTransactor // Write-only binding to the contract
	GardenFEEAccountFilterer   // Log filterer for contract events
}

// GardenFEEAccountCaller is an auto generated read-only Go binding around an Ethereum contract.
type GardenFEEAccountCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenFEEAccountTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GardenFEEAccountTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenFEEAccountFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GardenFEEAccountFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenFEEAccountSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GardenFEEAccountSession struct {
	Contract     *GardenFEEAccount // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GardenFEEAccountCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GardenFEEAccountCallerSession struct {
	Contract *GardenFEEAccountCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// GardenFEEAccountTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GardenFEEAccountTransactorSession struct {
	Contract     *GardenFEEAccountTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// GardenFEEAccountRaw is an auto generated low-level Go binding around an Ethereum contract.
type GardenFEEAccountRaw struct {
	Contract *GardenFEEAccount // Generic contract binding to access the raw methods on
}

// GardenFEEAccountCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GardenFEEAccountCallerRaw struct {
	Contract *GardenFEEAccountCaller // Generic read-only contract binding to access the raw methods on
}

// GardenFEEAccountTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GardenFEEAccountTransactorRaw struct {
	Contract *GardenFEEAccountTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGardenFEEAccount creates a new instance of GardenFEEAccount, bound to a specific deployed contract.
func NewGardenFEEAccount(address common.Address, backend bind.ContractBackend) (*GardenFEEAccount, error) {
	contract, err := bindGardenFEEAccount(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GardenFEEAccount{GardenFEEAccountCaller: GardenFEEAccountCaller{contract: contract}, GardenFEEAccountTransactor: GardenFEEAccountTransactor{contract: contract}, GardenFEEAccountFilterer: GardenFEEAccountFilterer{contract: contract}}, nil
}

// NewGardenFEEAccountCaller creates a new read-only instance of GardenFEEAccount, bound to a specific deployed contract.
func NewGardenFEEAccountCaller(address common.Address, caller bind.ContractCaller) (*GardenFEEAccountCaller, error) {
	contract, err := bindGardenFEEAccount(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GardenFEEAccountCaller{contract: contract}, nil
}

// NewGardenFEEAccountTransactor creates a new write-only instance of GardenFEEAccount, bound to a specific deployed contract.
func NewGardenFEEAccountTransactor(address common.Address, transactor bind.ContractTransactor) (*GardenFEEAccountTransactor, error) {
	contract, err := bindGardenFEEAccount(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GardenFEEAccountTransactor{contract: contract}, nil
}

// NewGardenFEEAccountFilterer creates a new log filterer instance of GardenFEEAccount, bound to a specific deployed contract.
func NewGardenFEEAccountFilterer(address common.Address, filterer bind.ContractFilterer) (*GardenFEEAccountFilterer, error) {
	contract, err := bindGardenFEEAccount(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GardenFEEAccountFilterer{contract: contract}, nil
}

// bindGardenFEEAccount binds a generic wrapper to an already deployed contract.
func bindGardenFEEAccount(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GardenFEEAccountMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GardenFEEAccount *GardenFEEAccountRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GardenFEEAccount.Contract.GardenFEEAccountCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GardenFEEAccount *GardenFEEAccountRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenFEEAccount.Contract.GardenFEEAccountTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GardenFEEAccount *GardenFEEAccountRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GardenFEEAccount.Contract.GardenFEEAccountTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GardenFEEAccount *GardenFEEAccountCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GardenFEEAccount.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GardenFEEAccount *GardenFEEAccountTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenFEEAccount.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GardenFEEAccount *GardenFEEAccountTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GardenFEEAccount.Contract.contract.Transact(opts, method, params...)
}

// Amount is a free data retrieval call binding the contract method 0xaa8c217c.
//
// Solidity: function amount() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountCaller) Amount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GardenFEEAccount.contract.Call(opts, &out, "amount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Amount is a free data retrieval call binding the contract method 0xaa8c217c.
//
// Solidity: function amount() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountSession) Amount() (*big.Int, error) {
	return _GardenFEEAccount.Contract.Amount(&_GardenFEEAccount.CallOpts)
}

// Amount is a free data retrieval call binding the contract method 0xaa8c217c.
//
// Solidity: function amount() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountCallerSession) Amount() (*big.Int, error) {
	return _GardenFEEAccount.Contract.Amount(&_GardenFEEAccount.CallOpts)
}

// ClaimHash is a free data retrieval call binding the contract method 0x9044ef10.
//
// Solidity: function claimHash(uint256 amount_, uint256 nonce_, (bytes32,uint256,uint256,uint256)[] htlcs) view returns(bytes32)
func (_GardenFEEAccount *GardenFEEAccountCaller) ClaimHash(opts *bind.CallOpts, amount_ *big.Int, nonce_ *big.Int, htlcs []GardenFEEAccountHTLC) ([32]byte, error) {
	var out []interface{}
	err := _GardenFEEAccount.contract.Call(opts, &out, "claimHash", amount_, nonce_, htlcs)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ClaimHash is a free data retrieval call binding the contract method 0x9044ef10.
//
// Solidity: function claimHash(uint256 amount_, uint256 nonce_, (bytes32,uint256,uint256,uint256)[] htlcs) view returns(bytes32)
func (_GardenFEEAccount *GardenFEEAccountSession) ClaimHash(amount_ *big.Int, nonce_ *big.Int, htlcs []GardenFEEAccountHTLC) ([32]byte, error) {
	return _GardenFEEAccount.Contract.ClaimHash(&_GardenFEEAccount.CallOpts, amount_, nonce_, htlcs)
}

// ClaimHash is a free data retrieval call binding the contract method 0x9044ef10.
//
// Solidity: function claimHash(uint256 amount_, uint256 nonce_, (bytes32,uint256,uint256,uint256)[] htlcs) view returns(bytes32)
func (_GardenFEEAccount *GardenFEEAccountCallerSession) ClaimHash(amount_ *big.Int, nonce_ *big.Int, htlcs []GardenFEEAccountHTLC) ([32]byte, error) {
	return _GardenFEEAccount.Contract.ClaimHash(&_GardenFEEAccount.CallOpts, amount_, nonce_, htlcs)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_GardenFEEAccount *GardenFEEAccountCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _GardenFEEAccount.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_GardenFEEAccount *GardenFEEAccountSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _GardenFEEAccount.Contract.Eip712Domain(&_GardenFEEAccount.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_GardenFEEAccount *GardenFEEAccountCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _GardenFEEAccount.Contract.Eip712Domain(&_GardenFEEAccount.CallOpts)
}

// Expiration is a free data retrieval call binding the contract method 0x4665096d.
//
// Solidity: function expiration() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountCaller) Expiration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GardenFEEAccount.contract.Call(opts, &out, "expiration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Expiration is a free data retrieval call binding the contract method 0x4665096d.
//
// Solidity: function expiration() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountSession) Expiration() (*big.Int, error) {
	return _GardenFEEAccount.Contract.Expiration(&_GardenFEEAccount.CallOpts)
}

// Expiration is a free data retrieval call binding the contract method 0x4665096d.
//
// Solidity: function expiration() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountCallerSession) Expiration() (*big.Int, error) {
	return _GardenFEEAccount.Contract.Expiration(&_GardenFEEAccount.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_GardenFEEAccount *GardenFEEAccountCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GardenFEEAccount.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_GardenFEEAccount *GardenFEEAccountSession) Factory() (common.Address, error) {
	return _GardenFEEAccount.Contract.Factory(&_GardenFEEAccount.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_GardenFEEAccount *GardenFEEAccountCallerSession) Factory() (common.Address, error) {
	return _GardenFEEAccount.Contract.Factory(&_GardenFEEAccount.CallOpts)
}

// Funder is a free data retrieval call binding the contract method 0x041ae880.
//
// Solidity: function funder() view returns(address)
func (_GardenFEEAccount *GardenFEEAccountCaller) Funder(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GardenFEEAccount.contract.Call(opts, &out, "funder")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Funder is a free data retrieval call binding the contract method 0x041ae880.
//
// Solidity: function funder() view returns(address)
func (_GardenFEEAccount *GardenFEEAccountSession) Funder() (common.Address, error) {
	return _GardenFEEAccount.Contract.Funder(&_GardenFEEAccount.CallOpts)
}

// Funder is a free data retrieval call binding the contract method 0x041ae880.
//
// Solidity: function funder() view returns(address)
func (_GardenFEEAccount *GardenFEEAccountCallerSession) Funder() (common.Address, error) {
	return _GardenFEEAccount.Contract.Funder(&_GardenFEEAccount.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountCaller) Nonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GardenFEEAccount.contract.Call(opts, &out, "nonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountSession) Nonce() (*big.Int, error) {
	return _GardenFEEAccount.Contract.Nonce(&_GardenFEEAccount.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountCallerSession) Nonce() (*big.Int, error) {
	return _GardenFEEAccount.Contract.Nonce(&_GardenFEEAccount.CallOpts)
}

// Recipient is a free data retrieval call binding the contract method 0x66d003ac.
//
// Solidity: function recipient() view returns(address)
func (_GardenFEEAccount *GardenFEEAccountCaller) Recipient(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GardenFEEAccount.contract.Call(opts, &out, "recipient")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Recipient is a free data retrieval call binding the contract method 0x66d003ac.
//
// Solidity: function recipient() view returns(address)
func (_GardenFEEAccount *GardenFEEAccountSession) Recipient() (common.Address, error) {
	return _GardenFEEAccount.Contract.Recipient(&_GardenFEEAccount.CallOpts)
}

// Recipient is a free data retrieval call binding the contract method 0x66d003ac.
//
// Solidity: function recipient() view returns(address)
func (_GardenFEEAccount *GardenFEEAccountCallerSession) Recipient() (common.Address, error) {
	return _GardenFEEAccount.Contract.Recipient(&_GardenFEEAccount.CallOpts)
}

// SecretsProvided is a free data retrieval call binding the contract method 0x103074f9.
//
// Solidity: function secretsProvided() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountCaller) SecretsProvided(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GardenFEEAccount.contract.Call(opts, &out, "secretsProvided")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SecretsProvided is a free data retrieval call binding the contract method 0x103074f9.
//
// Solidity: function secretsProvided() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountSession) SecretsProvided() (*big.Int, error) {
	return _GardenFEEAccount.Contract.SecretsProvided(&_GardenFEEAccount.CallOpts)
}

// SecretsProvided is a free data retrieval call binding the contract method 0x103074f9.
//
// Solidity: function secretsProvided() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountCallerSession) SecretsProvided() (*big.Int, error) {
	return _GardenFEEAccount.Contract.SecretsProvided(&_GardenFEEAccount.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenFEEAccount *GardenFEEAccountCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GardenFEEAccount.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenFEEAccount *GardenFEEAccountSession) Token() (common.Address, error) {
	return _GardenFEEAccount.Contract.Token(&_GardenFEEAccount.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenFEEAccount *GardenFEEAccountCallerSession) Token() (common.Address, error) {
	return _GardenFEEAccount.Contract.Token(&_GardenFEEAccount.CallOpts)
}

// TotalAmount is a free data retrieval call binding the contract method 0x1a39d8ef.
//
// Solidity: function totalAmount() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountCaller) TotalAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GardenFEEAccount.contract.Call(opts, &out, "totalAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalAmount is a free data retrieval call binding the contract method 0x1a39d8ef.
//
// Solidity: function totalAmount() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountSession) TotalAmount() (*big.Int, error) {
	return _GardenFEEAccount.Contract.TotalAmount(&_GardenFEEAccount.CallOpts)
}

// TotalAmount is a free data retrieval call binding the contract method 0x1a39d8ef.
//
// Solidity: function totalAmount() view returns(uint256)
func (_GardenFEEAccount *GardenFEEAccountCallerSession) TotalAmount() (*big.Int, error) {
	return _GardenFEEAccount.Contract.TotalAmount(&_GardenFEEAccount.CallOpts)
}

// GardenFEEAccountInit is a paid mutator transaction binding the contract method 0xfa809e7a.
//
// Solidity: function __GardenFEEAccount_init(address token_, address funder_, address recipient_, string feeAccountName, string feeAccountVersion) returns()
func (_GardenFEEAccount *GardenFEEAccountTransactor) GardenFEEAccountInit(opts *bind.TransactOpts, token_ common.Address, funder_ common.Address, recipient_ common.Address, feeAccountName string, feeAccountVersion string) (*types.Transaction, error) {
	return _GardenFEEAccount.contract.Transact(opts, "__GardenFEEAccount_init", token_, funder_, recipient_, feeAccountName, feeAccountVersion)
}

// GardenFEEAccountInit is a paid mutator transaction binding the contract method 0xfa809e7a.
//
// Solidity: function __GardenFEEAccount_init(address token_, address funder_, address recipient_, string feeAccountName, string feeAccountVersion) returns()
func (_GardenFEEAccount *GardenFEEAccountSession) GardenFEEAccountInit(token_ common.Address, funder_ common.Address, recipient_ common.Address, feeAccountName string, feeAccountVersion string) (*types.Transaction, error) {
	return _GardenFEEAccount.Contract.GardenFEEAccountInit(&_GardenFEEAccount.TransactOpts, token_, funder_, recipient_, feeAccountName, feeAccountVersion)
}

// GardenFEEAccountInit is a paid mutator transaction binding the contract method 0xfa809e7a.
//
// Solidity: function __GardenFEEAccount_init(address token_, address funder_, address recipient_, string feeAccountName, string feeAccountVersion) returns()
func (_GardenFEEAccount *GardenFEEAccountTransactorSession) GardenFEEAccountInit(token_ common.Address, funder_ common.Address, recipient_ common.Address, feeAccountName string, feeAccountVersion string) (*types.Transaction, error) {
	return _GardenFEEAccount.Contract.GardenFEEAccountInit(&_GardenFEEAccount.TransactOpts, token_, funder_, recipient_, feeAccountName, feeAccountVersion)
}

// Claim is a paid mutator transaction binding the contract method 0x74cb7eaf.
//
// Solidity: function claim(uint256 amount_, uint256 nonce_, (bytes32,uint256,uint256,uint256)[] htlcs, bytes[] secrets, bytes funderSig, bytes recipientSig) returns()
func (_GardenFEEAccount *GardenFEEAccountTransactor) Claim(opts *bind.TransactOpts, amount_ *big.Int, nonce_ *big.Int, htlcs []GardenFEEAccountHTLC, secrets [][]byte, funderSig []byte, recipientSig []byte) (*types.Transaction, error) {
	return _GardenFEEAccount.contract.Transact(opts, "claim", amount_, nonce_, htlcs, secrets, funderSig, recipientSig)
}

// Claim is a paid mutator transaction binding the contract method 0x74cb7eaf.
//
// Solidity: function claim(uint256 amount_, uint256 nonce_, (bytes32,uint256,uint256,uint256)[] htlcs, bytes[] secrets, bytes funderSig, bytes recipientSig) returns()
func (_GardenFEEAccount *GardenFEEAccountSession) Claim(amount_ *big.Int, nonce_ *big.Int, htlcs []GardenFEEAccountHTLC, secrets [][]byte, funderSig []byte, recipientSig []byte) (*types.Transaction, error) {
	return _GardenFEEAccount.Contract.Claim(&_GardenFEEAccount.TransactOpts, amount_, nonce_, htlcs, secrets, funderSig, recipientSig)
}

// Claim is a paid mutator transaction binding the contract method 0x74cb7eaf.
//
// Solidity: function claim(uint256 amount_, uint256 nonce_, (bytes32,uint256,uint256,uint256)[] htlcs, bytes[] secrets, bytes funderSig, bytes recipientSig) returns()
func (_GardenFEEAccount *GardenFEEAccountTransactorSession) Claim(amount_ *big.Int, nonce_ *big.Int, htlcs []GardenFEEAccountHTLC, secrets [][]byte, funderSig []byte, recipientSig []byte) (*types.Transaction, error) {
	return _GardenFEEAccount.Contract.Claim(&_GardenFEEAccount.TransactOpts, amount_, nonce_, htlcs, secrets, funderSig, recipientSig)
}

// Close is a paid mutator transaction binding the contract method 0x44275289.
//
// Solidity: function close(uint256 amount_, bytes funderSig, bytes recipientSig) returns()
func (_GardenFEEAccount *GardenFEEAccountTransactor) Close(opts *bind.TransactOpts, amount_ *big.Int, funderSig []byte, recipientSig []byte) (*types.Transaction, error) {
	return _GardenFEEAccount.contract.Transact(opts, "close", amount_, funderSig, recipientSig)
}

// Close is a paid mutator transaction binding the contract method 0x44275289.
//
// Solidity: function close(uint256 amount_, bytes funderSig, bytes recipientSig) returns()
func (_GardenFEEAccount *GardenFEEAccountSession) Close(amount_ *big.Int, funderSig []byte, recipientSig []byte) (*types.Transaction, error) {
	return _GardenFEEAccount.Contract.Close(&_GardenFEEAccount.TransactOpts, amount_, funderSig, recipientSig)
}

// Close is a paid mutator transaction binding the contract method 0x44275289.
//
// Solidity: function close(uint256 amount_, bytes funderSig, bytes recipientSig) returns()
func (_GardenFEEAccount *GardenFEEAccountTransactorSession) Close(amount_ *big.Int, funderSig []byte, recipientSig []byte) (*types.Transaction, error) {
	return _GardenFEEAccount.Contract.Close(&_GardenFEEAccount.TransactOpts, amount_, funderSig, recipientSig)
}

// Settle is a paid mutator transaction binding the contract method 0x11da60b4.
//
// Solidity: function settle() returns()
func (_GardenFEEAccount *GardenFEEAccountTransactor) Settle(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenFEEAccount.contract.Transact(opts, "settle")
}

// Settle is a paid mutator transaction binding the contract method 0x11da60b4.
//
// Solidity: function settle() returns()
func (_GardenFEEAccount *GardenFEEAccountSession) Settle() (*types.Transaction, error) {
	return _GardenFEEAccount.Contract.Settle(&_GardenFEEAccount.TransactOpts)
}

// Settle is a paid mutator transaction binding the contract method 0x11da60b4.
//
// Solidity: function settle() returns()
func (_GardenFEEAccount *GardenFEEAccountTransactorSession) Settle() (*types.Transaction, error) {
	return _GardenFEEAccount.Contract.Settle(&_GardenFEEAccount.TransactOpts)
}

// GardenFEEAccountEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the GardenFEEAccount contract.
type GardenFEEAccountEIP712DomainChangedIterator struct {
	Event *GardenFEEAccountEIP712DomainChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GardenFEEAccountEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenFEEAccountEIP712DomainChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GardenFEEAccountEIP712DomainChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GardenFEEAccountEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenFEEAccountEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenFEEAccountEIP712DomainChanged represents a EIP712DomainChanged event raised by the GardenFEEAccount contract.
type GardenFEEAccountEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_GardenFEEAccount *GardenFEEAccountFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*GardenFEEAccountEIP712DomainChangedIterator, error) {

	logs, sub, err := _GardenFEEAccount.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &GardenFEEAccountEIP712DomainChangedIterator{contract: _GardenFEEAccount.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_GardenFEEAccount *GardenFEEAccountFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *GardenFEEAccountEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _GardenFEEAccount.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenFEEAccountEIP712DomainChanged)
				if err := _GardenFEEAccount.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_GardenFEEAccount *GardenFEEAccountFilterer) ParseEIP712DomainChanged(log types.Log) (*GardenFEEAccountEIP712DomainChanged, error) {
	event := new(GardenFEEAccountEIP712DomainChanged)
	if err := _GardenFEEAccount.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenFEEAccountInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the GardenFEEAccount contract.
type GardenFEEAccountInitializedIterator struct {
	Event *GardenFEEAccountInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GardenFEEAccountInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenFEEAccountInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GardenFEEAccountInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GardenFEEAccountInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenFEEAccountInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenFEEAccountInitialized represents a Initialized event raised by the GardenFEEAccount contract.
type GardenFEEAccountInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GardenFEEAccount *GardenFEEAccountFilterer) FilterInitialized(opts *bind.FilterOpts) (*GardenFEEAccountInitializedIterator, error) {

	logs, sub, err := _GardenFEEAccount.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GardenFEEAccountInitializedIterator{contract: _GardenFEEAccount.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GardenFEEAccount *GardenFEEAccountFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GardenFEEAccountInitialized) (event.Subscription, error) {

	logs, sub, err := _GardenFEEAccount.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenFEEAccountInitialized)
				if err := _GardenFEEAccount.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GardenFEEAccount *GardenFEEAccountFilterer) ParseInitialized(log types.Log) (*GardenFEEAccountInitialized, error) {
	event := new(GardenFEEAccountInitialized)
	if err := _GardenFEEAccount.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
