// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gardenfeeaccountfactory

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

// FeeAccountHTLC is an auto generated low-level Go binding around an user-defined struct.
type FeeAccountHTLC struct {
	SecretHash    [32]byte
	Expiry        *big.Int
	SendAmount    *big.Int
	RecieveAmount *big.Int
}

// GardenFEEAccountFactoryMetaData contains all meta data concerning the GardenFEEAccountFactory contract.
var GardenFEEAccountFactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"token_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeManager_\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"feeAccountName_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"feeAccountVersion_\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"channel\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"name\":\"Claimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"channel\",\"type\":\"address\"}],\"name\":\"Closed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"funder\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"channel\",\"type\":\"address\"}],\"name\":\"Created\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"channels\",\"outputs\":[{\"internalType\":\"contractFeeAccount\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"name\":\"claimed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"closed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"create\",\"outputs\":[{\"internalType\":\"contractFeeAccount\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sendAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"recieveAmount\",\"type\":\"uint256\"}],\"internalType\":\"structFeeAccount.HTLC[]\",\"name\":\"htlcs\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"secrets\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"funderSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"recipientSig\",\"type\":\"bytes\"}],\"name\":\"createAndClaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"funderSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"recipientSig\",\"type\":\"bytes\"}],\"name\":\"createAndClose\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeAccountName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeAccountVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"feeManagerCreate\",\"outputs\":[{\"internalType\":\"contractFeeAccount\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"settle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"template\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// GardenFEEAccountFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use GardenFEEAccountFactoryMetaData.ABI instead.
var GardenFEEAccountFactoryABI = GardenFEEAccountFactoryMetaData.ABI

// GardenFEEAccountFactory is an auto generated Go binding around an Ethereum contract.
type GardenFEEAccountFactory struct {
	GardenFEEAccountFactoryCaller     // Read-only binding to the contract
	GardenFEEAccountFactoryTransactor // Write-only binding to the contract
	GardenFEEAccountFactoryFilterer   // Log filterer for contract events
}

// GardenFEEAccountFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type GardenFEEAccountFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenFEEAccountFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GardenFEEAccountFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenFEEAccountFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GardenFEEAccountFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenFEEAccountFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GardenFEEAccountFactorySession struct {
	Contract     *GardenFEEAccountFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// GardenFEEAccountFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GardenFEEAccountFactoryCallerSession struct {
	Contract *GardenFEEAccountFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// GardenFEEAccountFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GardenFEEAccountFactoryTransactorSession struct {
	Contract     *GardenFEEAccountFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// GardenFEEAccountFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type GardenFEEAccountFactoryRaw struct {
	Contract *GardenFEEAccountFactory // Generic contract binding to access the raw methods on
}

// GardenFEEAccountFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GardenFEEAccountFactoryCallerRaw struct {
	Contract *GardenFEEAccountFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// GardenFEEAccountFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GardenFEEAccountFactoryTransactorRaw struct {
	Contract *GardenFEEAccountFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGardenFEEAccountFactory creates a new instance of GardenFEEAccountFactory, bound to a specific deployed contract.
func NewGardenFEEAccountFactory(address common.Address, backend bind.ContractBackend) (*GardenFEEAccountFactory, error) {
	contract, err := bindGardenFEEAccountFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GardenFEEAccountFactory{GardenFEEAccountFactoryCaller: GardenFEEAccountFactoryCaller{contract: contract}, GardenFEEAccountFactoryTransactor: GardenFEEAccountFactoryTransactor{contract: contract}, GardenFEEAccountFactoryFilterer: GardenFEEAccountFactoryFilterer{contract: contract}}, nil
}

// NewGardenFEEAccountFactoryCaller creates a new read-only instance of GardenFEEAccountFactory, bound to a specific deployed contract.
func NewGardenFEEAccountFactoryCaller(address common.Address, caller bind.ContractCaller) (*GardenFEEAccountFactoryCaller, error) {
	contract, err := bindGardenFEEAccountFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GardenFEEAccountFactoryCaller{contract: contract}, nil
}

// NewGardenFEEAccountFactoryTransactor creates a new write-only instance of GardenFEEAccountFactory, bound to a specific deployed contract.
func NewGardenFEEAccountFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*GardenFEEAccountFactoryTransactor, error) {
	contract, err := bindGardenFEEAccountFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GardenFEEAccountFactoryTransactor{contract: contract}, nil
}

// NewGardenFEEAccountFactoryFilterer creates a new log filterer instance of GardenFEEAccountFactory, bound to a specific deployed contract.
func NewGardenFEEAccountFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*GardenFEEAccountFactoryFilterer, error) {
	contract, err := bindGardenFEEAccountFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GardenFEEAccountFactoryFilterer{contract: contract}, nil
}

// bindGardenFEEAccountFactory binds a generic wrapper to an already deployed contract.
func bindGardenFEEAccountFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GardenFEEAccountFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GardenFEEAccountFactory.Contract.GardenFEEAccountFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.GardenFEEAccountFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.GardenFEEAccountFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GardenFEEAccountFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.contract.Transact(opts, method, params...)
}

// Channels is a free data retrieval call binding the contract method 0x7dce34f7.
//
// Solidity: function channels(address ) view returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCaller) Channels(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _GardenFEEAccountFactory.contract.Call(opts, &out, "channels", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Channels is a free data retrieval call binding the contract method 0x7dce34f7.
//
// Solidity: function channels(address ) view returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) Channels(arg0 common.Address) (common.Address, error) {
	return _GardenFEEAccountFactory.Contract.Channels(&_GardenFEEAccountFactory.CallOpts, arg0)
}

// Channels is a free data retrieval call binding the contract method 0x7dce34f7.
//
// Solidity: function channels(address ) view returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCallerSession) Channels(arg0 common.Address) (common.Address, error) {
	return _GardenFEEAccountFactory.Contract.Channels(&_GardenFEEAccountFactory.CallOpts, arg0)
}

// FeeAccountName is a free data retrieval call binding the contract method 0x7f48196d.
//
// Solidity: function feeAccountName() view returns(string)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCaller) FeeAccountName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GardenFEEAccountFactory.contract.Call(opts, &out, "feeAccountName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// FeeAccountName is a free data retrieval call binding the contract method 0x7f48196d.
//
// Solidity: function feeAccountName() view returns(string)
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) FeeAccountName() (string, error) {
	return _GardenFEEAccountFactory.Contract.FeeAccountName(&_GardenFEEAccountFactory.CallOpts)
}

// FeeAccountName is a free data retrieval call binding the contract method 0x7f48196d.
//
// Solidity: function feeAccountName() view returns(string)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCallerSession) FeeAccountName() (string, error) {
	return _GardenFEEAccountFactory.Contract.FeeAccountName(&_GardenFEEAccountFactory.CallOpts)
}

// FeeAccountVersion is a free data retrieval call binding the contract method 0xb10b1335.
//
// Solidity: function feeAccountVersion() view returns(string)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCaller) FeeAccountVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GardenFEEAccountFactory.contract.Call(opts, &out, "feeAccountVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// FeeAccountVersion is a free data retrieval call binding the contract method 0xb10b1335.
//
// Solidity: function feeAccountVersion() view returns(string)
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) FeeAccountVersion() (string, error) {
	return _GardenFEEAccountFactory.Contract.FeeAccountVersion(&_GardenFEEAccountFactory.CallOpts)
}

// FeeAccountVersion is a free data retrieval call binding the contract method 0xb10b1335.
//
// Solidity: function feeAccountVersion() view returns(string)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCallerSession) FeeAccountVersion() (string, error) {
	return _GardenFEEAccountFactory.Contract.FeeAccountVersion(&_GardenFEEAccountFactory.CallOpts)
}

// FeeManager is a free data retrieval call binding the contract method 0xd0fb0203.
//
// Solidity: function feeManager() view returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCaller) FeeManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GardenFEEAccountFactory.contract.Call(opts, &out, "feeManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeManager is a free data retrieval call binding the contract method 0xd0fb0203.
//
// Solidity: function feeManager() view returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) FeeManager() (common.Address, error) {
	return _GardenFEEAccountFactory.Contract.FeeManager(&_GardenFEEAccountFactory.CallOpts)
}

// FeeManager is a free data retrieval call binding the contract method 0xd0fb0203.
//
// Solidity: function feeManager() view returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCallerSession) FeeManager() (common.Address, error) {
	return _GardenFEEAccountFactory.Contract.FeeManager(&_GardenFEEAccountFactory.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GardenFEEAccountFactory.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _GardenFEEAccountFactory.Contract.Nonces(&_GardenFEEAccountFactory.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _GardenFEEAccountFactory.Contract.Nonces(&_GardenFEEAccountFactory.CallOpts, arg0)
}

// Template is a free data retrieval call binding the contract method 0x6f2ddd93.
//
// Solidity: function template() view returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCaller) Template(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GardenFEEAccountFactory.contract.Call(opts, &out, "template")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Template is a free data retrieval call binding the contract method 0x6f2ddd93.
//
// Solidity: function template() view returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) Template() (common.Address, error) {
	return _GardenFEEAccountFactory.Contract.Template(&_GardenFEEAccountFactory.CallOpts)
}

// Template is a free data retrieval call binding the contract method 0x6f2ddd93.
//
// Solidity: function template() view returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCallerSession) Template() (common.Address, error) {
	return _GardenFEEAccountFactory.Contract.Template(&_GardenFEEAccountFactory.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GardenFEEAccountFactory.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) Token() (common.Address, error) {
	return _GardenFEEAccountFactory.Contract.Token(&_GardenFEEAccountFactory.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryCallerSession) Token() (common.Address, error) {
	return _GardenFEEAccountFactory.Contract.Token(&_GardenFEEAccountFactory.CallOpts)
}

// Claimed is a paid mutator transaction binding the contract method 0xcfd94956.
//
// Solidity: function claimed(address recipient, uint256 amount, uint256 nonce, uint256 expiration) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactor) Claimed(opts *bind.TransactOpts, recipient common.Address, amount *big.Int, nonce *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.contract.Transact(opts, "claimed", recipient, amount, nonce, expiration)
}

// Claimed is a paid mutator transaction binding the contract method 0xcfd94956.
//
// Solidity: function claimed(address recipient, uint256 amount, uint256 nonce, uint256 expiration) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) Claimed(recipient common.Address, amount *big.Int, nonce *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.Claimed(&_GardenFEEAccountFactory.TransactOpts, recipient, amount, nonce, expiration)
}

// Claimed is a paid mutator transaction binding the contract method 0xcfd94956.
//
// Solidity: function claimed(address recipient, uint256 amount, uint256 nonce, uint256 expiration) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactorSession) Claimed(recipient common.Address, amount *big.Int, nonce *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.Claimed(&_GardenFEEAccountFactory.TransactOpts, recipient, amount, nonce, expiration)
}

// Closed is a paid mutator transaction binding the contract method 0x52942373.
//
// Solidity: function closed(address recipient) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactor) Closed(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.contract.Transact(opts, "closed", recipient)
}

// Closed is a paid mutator transaction binding the contract method 0x52942373.
//
// Solidity: function closed(address recipient) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) Closed(recipient common.Address) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.Closed(&_GardenFEEAccountFactory.TransactOpts, recipient)
}

// Closed is a paid mutator transaction binding the contract method 0x52942373.
//
// Solidity: function closed(address recipient) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactorSession) Closed(recipient common.Address) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.Closed(&_GardenFEEAccountFactory.TransactOpts, recipient)
}

// Create is a paid mutator transaction binding the contract method 0xefc81a8c.
//
// Solidity: function create() returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactor) Create(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.contract.Transact(opts, "create")
}

// Create is a paid mutator transaction binding the contract method 0xefc81a8c.
//
// Solidity: function create() returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) Create() (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.Create(&_GardenFEEAccountFactory.TransactOpts)
}

// Create is a paid mutator transaction binding the contract method 0xefc81a8c.
//
// Solidity: function create() returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactorSession) Create() (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.Create(&_GardenFEEAccountFactory.TransactOpts)
}

// CreateAndClaim is a paid mutator transaction binding the contract method 0x3589aefd.
//
// Solidity: function createAndClaim(uint256 amount, uint256 nonce, (bytes32,uint256,uint256,uint256)[] htlcs, bytes[] secrets, bytes funderSig, bytes recipientSig) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactor) CreateAndClaim(opts *bind.TransactOpts, amount *big.Int, nonce *big.Int, htlcs []FeeAccountHTLC, secrets [][]byte, funderSig []byte, recipientSig []byte) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.contract.Transact(opts, "createAndClaim", amount, nonce, htlcs, secrets, funderSig, recipientSig)
}

// CreateAndClaim is a paid mutator transaction binding the contract method 0x3589aefd.
//
// Solidity: function createAndClaim(uint256 amount, uint256 nonce, (bytes32,uint256,uint256,uint256)[] htlcs, bytes[] secrets, bytes funderSig, bytes recipientSig) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) CreateAndClaim(amount *big.Int, nonce *big.Int, htlcs []FeeAccountHTLC, secrets [][]byte, funderSig []byte, recipientSig []byte) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.CreateAndClaim(&_GardenFEEAccountFactory.TransactOpts, amount, nonce, htlcs, secrets, funderSig, recipientSig)
}

// CreateAndClaim is a paid mutator transaction binding the contract method 0x3589aefd.
//
// Solidity: function createAndClaim(uint256 amount, uint256 nonce, (bytes32,uint256,uint256,uint256)[] htlcs, bytes[] secrets, bytes funderSig, bytes recipientSig) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactorSession) CreateAndClaim(amount *big.Int, nonce *big.Int, htlcs []FeeAccountHTLC, secrets [][]byte, funderSig []byte, recipientSig []byte) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.CreateAndClaim(&_GardenFEEAccountFactory.TransactOpts, amount, nonce, htlcs, secrets, funderSig, recipientSig)
}

// CreateAndClose is a paid mutator transaction binding the contract method 0x41a9a014.
//
// Solidity: function createAndClose(uint256 amount, bytes funderSig, bytes recipientSig) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactor) CreateAndClose(opts *bind.TransactOpts, amount *big.Int, funderSig []byte, recipientSig []byte) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.contract.Transact(opts, "createAndClose", amount, funderSig, recipientSig)
}

// CreateAndClose is a paid mutator transaction binding the contract method 0x41a9a014.
//
// Solidity: function createAndClose(uint256 amount, bytes funderSig, bytes recipientSig) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) CreateAndClose(amount *big.Int, funderSig []byte, recipientSig []byte) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.CreateAndClose(&_GardenFEEAccountFactory.TransactOpts, amount, funderSig, recipientSig)
}

// CreateAndClose is a paid mutator transaction binding the contract method 0x41a9a014.
//
// Solidity: function createAndClose(uint256 amount, bytes funderSig, bytes recipientSig) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactorSession) CreateAndClose(amount *big.Int, funderSig []byte, recipientSig []byte) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.CreateAndClose(&_GardenFEEAccountFactory.TransactOpts, amount, funderSig, recipientSig)
}

// FeeManagerCreate is a paid mutator transaction binding the contract method 0xabee6a54.
//
// Solidity: function feeManagerCreate(address recipient) returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactor) FeeManagerCreate(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.contract.Transact(opts, "feeManagerCreate", recipient)
}

// FeeManagerCreate is a paid mutator transaction binding the contract method 0xabee6a54.
//
// Solidity: function feeManagerCreate(address recipient) returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) FeeManagerCreate(recipient common.Address) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.FeeManagerCreate(&_GardenFEEAccountFactory.TransactOpts, recipient)
}

// FeeManagerCreate is a paid mutator transaction binding the contract method 0xabee6a54.
//
// Solidity: function feeManagerCreate(address recipient) returns(address)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactorSession) FeeManagerCreate(recipient common.Address) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.FeeManagerCreate(&_GardenFEEAccountFactory.TransactOpts, recipient)
}

// Settle is a paid mutator transaction binding the contract method 0x6a256b29.
//
// Solidity: function settle(address recipient) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactor) Settle(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.contract.Transact(opts, "settle", recipient)
}

// Settle is a paid mutator transaction binding the contract method 0x6a256b29.
//
// Solidity: function settle(address recipient) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactorySession) Settle(recipient common.Address) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.Settle(&_GardenFEEAccountFactory.TransactOpts, recipient)
}

// Settle is a paid mutator transaction binding the contract method 0x6a256b29.
//
// Solidity: function settle(address recipient) returns()
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryTransactorSession) Settle(recipient common.Address) (*types.Transaction, error) {
	return _GardenFEEAccountFactory.Contract.Settle(&_GardenFEEAccountFactory.TransactOpts, recipient)
}

// GardenFEEAccountFactoryClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the GardenFEEAccountFactory contract.
type GardenFEEAccountFactoryClaimedIterator struct {
	Event *GardenFEEAccountFactoryClaimed // Event containing the contract specifics and raw log

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
func (it *GardenFEEAccountFactoryClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenFEEAccountFactoryClaimed)
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
		it.Event = new(GardenFEEAccountFactoryClaimed)
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
func (it *GardenFEEAccountFactoryClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenFEEAccountFactoryClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenFEEAccountFactoryClaimed represents a Claimed event raised by the GardenFEEAccountFactory contract.
type GardenFEEAccountFactoryClaimed struct {
	Channel    common.Address
	Amount     *big.Int
	Nonce      *big.Int
	Expiration *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0x9cdcf2f7714cca3508c7f0110b04a90a80a3a8dd0e35de99689db74d28c5383e.
//
// Solidity: event Claimed(address indexed channel, uint256 amount, uint256 nonce, uint256 expiration)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryFilterer) FilterClaimed(opts *bind.FilterOpts, channel []common.Address) (*GardenFEEAccountFactoryClaimedIterator, error) {

	var channelRule []interface{}
	for _, channelItem := range channel {
		channelRule = append(channelRule, channelItem)
	}

	logs, sub, err := _GardenFEEAccountFactory.contract.FilterLogs(opts, "Claimed", channelRule)
	if err != nil {
		return nil, err
	}
	return &GardenFEEAccountFactoryClaimedIterator{contract: _GardenFEEAccountFactory.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0x9cdcf2f7714cca3508c7f0110b04a90a80a3a8dd0e35de99689db74d28c5383e.
//
// Solidity: event Claimed(address indexed channel, uint256 amount, uint256 nonce, uint256 expiration)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *GardenFEEAccountFactoryClaimed, channel []common.Address) (event.Subscription, error) {

	var channelRule []interface{}
	for _, channelItem := range channel {
		channelRule = append(channelRule, channelItem)
	}

	logs, sub, err := _GardenFEEAccountFactory.contract.WatchLogs(opts, "Claimed", channelRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenFEEAccountFactoryClaimed)
				if err := _GardenFEEAccountFactory.contract.UnpackLog(event, "Claimed", log); err != nil {
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

// ParseClaimed is a log parse operation binding the contract event 0x9cdcf2f7714cca3508c7f0110b04a90a80a3a8dd0e35de99689db74d28c5383e.
//
// Solidity: event Claimed(address indexed channel, uint256 amount, uint256 nonce, uint256 expiration)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryFilterer) ParseClaimed(log types.Log) (*GardenFEEAccountFactoryClaimed, error) {
	event := new(GardenFEEAccountFactoryClaimed)
	if err := _GardenFEEAccountFactory.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenFEEAccountFactoryClosedIterator is returned from FilterClosed and is used to iterate over the raw logs and unpacked data for Closed events raised by the GardenFEEAccountFactory contract.
type GardenFEEAccountFactoryClosedIterator struct {
	Event *GardenFEEAccountFactoryClosed // Event containing the contract specifics and raw log

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
func (it *GardenFEEAccountFactoryClosedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenFEEAccountFactoryClosed)
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
		it.Event = new(GardenFEEAccountFactoryClosed)
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
func (it *GardenFEEAccountFactoryClosedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenFEEAccountFactoryClosedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenFEEAccountFactoryClosed represents a Closed event raised by the GardenFEEAccountFactory contract.
type GardenFEEAccountFactoryClosed struct {
	Channel common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterClosed is a free log retrieval operation binding the contract event 0x13607bf9d2dd20e1f3a7daf47ab12856f8aad65e6ae7e2c75ace3d0c424a40e8.
//
// Solidity: event Closed(address indexed channel)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryFilterer) FilterClosed(opts *bind.FilterOpts, channel []common.Address) (*GardenFEEAccountFactoryClosedIterator, error) {

	var channelRule []interface{}
	for _, channelItem := range channel {
		channelRule = append(channelRule, channelItem)
	}

	logs, sub, err := _GardenFEEAccountFactory.contract.FilterLogs(opts, "Closed", channelRule)
	if err != nil {
		return nil, err
	}
	return &GardenFEEAccountFactoryClosedIterator{contract: _GardenFEEAccountFactory.contract, event: "Closed", logs: logs, sub: sub}, nil
}

// WatchClosed is a free log subscription operation binding the contract event 0x13607bf9d2dd20e1f3a7daf47ab12856f8aad65e6ae7e2c75ace3d0c424a40e8.
//
// Solidity: event Closed(address indexed channel)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryFilterer) WatchClosed(opts *bind.WatchOpts, sink chan<- *GardenFEEAccountFactoryClosed, channel []common.Address) (event.Subscription, error) {

	var channelRule []interface{}
	for _, channelItem := range channel {
		channelRule = append(channelRule, channelItem)
	}

	logs, sub, err := _GardenFEEAccountFactory.contract.WatchLogs(opts, "Closed", channelRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenFEEAccountFactoryClosed)
				if err := _GardenFEEAccountFactory.contract.UnpackLog(event, "Closed", log); err != nil {
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

// ParseClosed is a log parse operation binding the contract event 0x13607bf9d2dd20e1f3a7daf47ab12856f8aad65e6ae7e2c75ace3d0c424a40e8.
//
// Solidity: event Closed(address indexed channel)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryFilterer) ParseClosed(log types.Log) (*GardenFEEAccountFactoryClosed, error) {
	event := new(GardenFEEAccountFactoryClosed)
	if err := _GardenFEEAccountFactory.contract.UnpackLog(event, "Closed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenFEEAccountFactoryCreatedIterator is returned from FilterCreated and is used to iterate over the raw logs and unpacked data for Created events raised by the GardenFEEAccountFactory contract.
type GardenFEEAccountFactoryCreatedIterator struct {
	Event *GardenFEEAccountFactoryCreated // Event containing the contract specifics and raw log

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
func (it *GardenFEEAccountFactoryCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenFEEAccountFactoryCreated)
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
		it.Event = new(GardenFEEAccountFactoryCreated)
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
func (it *GardenFEEAccountFactoryCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenFEEAccountFactoryCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenFEEAccountFactoryCreated represents a Created event raised by the GardenFEEAccountFactory contract.
type GardenFEEAccountFactoryCreated struct {
	Funder  common.Address
	Channel common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterCreated is a free log retrieval operation binding the contract event 0x587ece4cd19692c5be1a4184503d607d45542d2aca0698c0068f52e09ccb541c.
//
// Solidity: event Created(address indexed funder, address indexed channel)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryFilterer) FilterCreated(opts *bind.FilterOpts, funder []common.Address, channel []common.Address) (*GardenFEEAccountFactoryCreatedIterator, error) {

	var funderRule []interface{}
	for _, funderItem := range funder {
		funderRule = append(funderRule, funderItem)
	}
	var channelRule []interface{}
	for _, channelItem := range channel {
		channelRule = append(channelRule, channelItem)
	}

	logs, sub, err := _GardenFEEAccountFactory.contract.FilterLogs(opts, "Created", funderRule, channelRule)
	if err != nil {
		return nil, err
	}
	return &GardenFEEAccountFactoryCreatedIterator{contract: _GardenFEEAccountFactory.contract, event: "Created", logs: logs, sub: sub}, nil
}

// WatchCreated is a free log subscription operation binding the contract event 0x587ece4cd19692c5be1a4184503d607d45542d2aca0698c0068f52e09ccb541c.
//
// Solidity: event Created(address indexed funder, address indexed channel)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryFilterer) WatchCreated(opts *bind.WatchOpts, sink chan<- *GardenFEEAccountFactoryCreated, funder []common.Address, channel []common.Address) (event.Subscription, error) {

	var funderRule []interface{}
	for _, funderItem := range funder {
		funderRule = append(funderRule, funderItem)
	}
	var channelRule []interface{}
	for _, channelItem := range channel {
		channelRule = append(channelRule, channelItem)
	}

	logs, sub, err := _GardenFEEAccountFactory.contract.WatchLogs(opts, "Created", funderRule, channelRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenFEEAccountFactoryCreated)
				if err := _GardenFEEAccountFactory.contract.UnpackLog(event, "Created", log); err != nil {
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

// ParseCreated is a log parse operation binding the contract event 0x587ece4cd19692c5be1a4184503d607d45542d2aca0698c0068f52e09ccb541c.
//
// Solidity: event Created(address indexed funder, address indexed channel)
func (_GardenFEEAccountFactory *GardenFEEAccountFactoryFilterer) ParseCreated(log types.Log) (*GardenFEEAccountFactoryCreated, error) {
	event := new(GardenFEEAccountFactoryCreated)
	if err := _GardenFEEAccountFactory.contract.UnpackLog(event, "Created", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
