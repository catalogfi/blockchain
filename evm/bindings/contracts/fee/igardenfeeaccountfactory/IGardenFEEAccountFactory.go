// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package igardenfeeaccountfactory

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

// IGardenFEEAccountFactoryMetaData contains all meta data concerning the IGardenFEEAccountFactory contract.
var IGardenFEEAccountFactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"name\":\"claimed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"closed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IGardenFEEAccountFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use IGardenFEEAccountFactoryMetaData.ABI instead.
var IGardenFEEAccountFactoryABI = IGardenFEEAccountFactoryMetaData.ABI

// IGardenFEEAccountFactory is an auto generated Go binding around an Ethereum contract.
type IGardenFEEAccountFactory struct {
	IGardenFEEAccountFactoryCaller     // Read-only binding to the contract
	IGardenFEEAccountFactoryTransactor // Write-only binding to the contract
	IGardenFEEAccountFactoryFilterer   // Log filterer for contract events
}

// IGardenFEEAccountFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type IGardenFEEAccountFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IGardenFEEAccountFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IGardenFEEAccountFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IGardenFEEAccountFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IGardenFEEAccountFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IGardenFEEAccountFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IGardenFEEAccountFactorySession struct {
	Contract     *IGardenFEEAccountFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IGardenFEEAccountFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IGardenFEEAccountFactoryCallerSession struct {
	Contract *IGardenFEEAccountFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// IGardenFEEAccountFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IGardenFEEAccountFactoryTransactorSession struct {
	Contract     *IGardenFEEAccountFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// IGardenFEEAccountFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type IGardenFEEAccountFactoryRaw struct {
	Contract *IGardenFEEAccountFactory // Generic contract binding to access the raw methods on
}

// IGardenFEEAccountFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IGardenFEEAccountFactoryCallerRaw struct {
	Contract *IGardenFEEAccountFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// IGardenFEEAccountFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IGardenFEEAccountFactoryTransactorRaw struct {
	Contract *IGardenFEEAccountFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIGardenFEEAccountFactory creates a new instance of IGardenFEEAccountFactory, bound to a specific deployed contract.
func NewIGardenFEEAccountFactory(address common.Address, backend bind.ContractBackend) (*IGardenFEEAccountFactory, error) {
	contract, err := bindIGardenFEEAccountFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IGardenFEEAccountFactory{IGardenFEEAccountFactoryCaller: IGardenFEEAccountFactoryCaller{contract: contract}, IGardenFEEAccountFactoryTransactor: IGardenFEEAccountFactoryTransactor{contract: contract}, IGardenFEEAccountFactoryFilterer: IGardenFEEAccountFactoryFilterer{contract: contract}}, nil
}

// NewIGardenFEEAccountFactoryCaller creates a new read-only instance of IGardenFEEAccountFactory, bound to a specific deployed contract.
func NewIGardenFEEAccountFactoryCaller(address common.Address, caller bind.ContractCaller) (*IGardenFEEAccountFactoryCaller, error) {
	contract, err := bindIGardenFEEAccountFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IGardenFEEAccountFactoryCaller{contract: contract}, nil
}

// NewIGardenFEEAccountFactoryTransactor creates a new write-only instance of IGardenFEEAccountFactory, bound to a specific deployed contract.
func NewIGardenFEEAccountFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*IGardenFEEAccountFactoryTransactor, error) {
	contract, err := bindIGardenFEEAccountFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IGardenFEEAccountFactoryTransactor{contract: contract}, nil
}

// NewIGardenFEEAccountFactoryFilterer creates a new log filterer instance of IGardenFEEAccountFactory, bound to a specific deployed contract.
func NewIGardenFEEAccountFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*IGardenFEEAccountFactoryFilterer, error) {
	contract, err := bindIGardenFEEAccountFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IGardenFEEAccountFactoryFilterer{contract: contract}, nil
}

// bindIGardenFEEAccountFactory binds a generic wrapper to an already deployed contract.
func bindIGardenFEEAccountFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IGardenFEEAccountFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IGardenFEEAccountFactory *IGardenFEEAccountFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IGardenFEEAccountFactory.Contract.IGardenFEEAccountFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IGardenFEEAccountFactory *IGardenFEEAccountFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IGardenFEEAccountFactory.Contract.IGardenFEEAccountFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IGardenFEEAccountFactory *IGardenFEEAccountFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IGardenFEEAccountFactory.Contract.IGardenFEEAccountFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IGardenFEEAccountFactory *IGardenFEEAccountFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IGardenFEEAccountFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IGardenFEEAccountFactory *IGardenFEEAccountFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IGardenFEEAccountFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IGardenFEEAccountFactory *IGardenFEEAccountFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IGardenFEEAccountFactory.Contract.contract.Transact(opts, method, params...)
}

// Claimed is a paid mutator transaction binding the contract method 0xcfd94956.
//
// Solidity: function claimed(address recipient, uint256 amount, uint256 nonce, uint256 expiration) returns()
func (_IGardenFEEAccountFactory *IGardenFEEAccountFactoryTransactor) Claimed(opts *bind.TransactOpts, recipient common.Address, amount *big.Int, nonce *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _IGardenFEEAccountFactory.contract.Transact(opts, "claimed", recipient, amount, nonce, expiration)
}

// Claimed is a paid mutator transaction binding the contract method 0xcfd94956.
//
// Solidity: function claimed(address recipient, uint256 amount, uint256 nonce, uint256 expiration) returns()
func (_IGardenFEEAccountFactory *IGardenFEEAccountFactorySession) Claimed(recipient common.Address, amount *big.Int, nonce *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _IGardenFEEAccountFactory.Contract.Claimed(&_IGardenFEEAccountFactory.TransactOpts, recipient, amount, nonce, expiration)
}

// Claimed is a paid mutator transaction binding the contract method 0xcfd94956.
//
// Solidity: function claimed(address recipient, uint256 amount, uint256 nonce, uint256 expiration) returns()
func (_IGardenFEEAccountFactory *IGardenFEEAccountFactoryTransactorSession) Claimed(recipient common.Address, amount *big.Int, nonce *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _IGardenFEEAccountFactory.Contract.Claimed(&_IGardenFEEAccountFactory.TransactOpts, recipient, amount, nonce, expiration)
}

// Closed is a paid mutator transaction binding the contract method 0x52942373.
//
// Solidity: function closed(address recipient) returns()
func (_IGardenFEEAccountFactory *IGardenFEEAccountFactoryTransactor) Closed(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _IGardenFEEAccountFactory.contract.Transact(opts, "closed", recipient)
}

// Closed is a paid mutator transaction binding the contract method 0x52942373.
//
// Solidity: function closed(address recipient) returns()
func (_IGardenFEEAccountFactory *IGardenFEEAccountFactorySession) Closed(recipient common.Address) (*types.Transaction, error) {
	return _IGardenFEEAccountFactory.Contract.Closed(&_IGardenFEEAccountFactory.TransactOpts, recipient)
}

// Closed is a paid mutator transaction binding the contract method 0x52942373.
//
// Solidity: function closed(address recipient) returns()
func (_IGardenFEEAccountFactory *IGardenFEEAccountFactoryTransactorSession) Closed(recipient common.Address) (*types.Transaction, error) {
	return _IGardenFEEAccountFactory.Contract.Closed(&_IGardenFEEAccountFactory.TransactOpts, recipient)
}
