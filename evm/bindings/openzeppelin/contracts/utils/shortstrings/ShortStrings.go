// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package shortstrings

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

// ShortStringsMetaData contains all meta data concerning the ShortStrings contract.
var ShortStringsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidShortString\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"}]",
}

// ShortStringsABI is the input ABI used to generate the binding from.
// Deprecated: Use ShortStringsMetaData.ABI instead.
var ShortStringsABI = ShortStringsMetaData.ABI

// ShortStrings is an auto generated Go binding around an Ethereum contract.
type ShortStrings struct {
	ShortStringsCaller     // Read-only binding to the contract
	ShortStringsTransactor // Write-only binding to the contract
	ShortStringsFilterer   // Log filterer for contract events
}

// ShortStringsCaller is an auto generated read-only Go binding around an Ethereum contract.
type ShortStringsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ShortStringsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ShortStringsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ShortStringsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ShortStringsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ShortStringsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ShortStringsSession struct {
	Contract     *ShortStrings     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ShortStringsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ShortStringsCallerSession struct {
	Contract *ShortStringsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ShortStringsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ShortStringsTransactorSession struct {
	Contract     *ShortStringsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ShortStringsRaw is an auto generated low-level Go binding around an Ethereum contract.
type ShortStringsRaw struct {
	Contract *ShortStrings // Generic contract binding to access the raw methods on
}

// ShortStringsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ShortStringsCallerRaw struct {
	Contract *ShortStringsCaller // Generic read-only contract binding to access the raw methods on
}

// ShortStringsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ShortStringsTransactorRaw struct {
	Contract *ShortStringsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewShortStrings creates a new instance of ShortStrings, bound to a specific deployed contract.
func NewShortStrings(address common.Address, backend bind.ContractBackend) (*ShortStrings, error) {
	contract, err := bindShortStrings(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ShortStrings{ShortStringsCaller: ShortStringsCaller{contract: contract}, ShortStringsTransactor: ShortStringsTransactor{contract: contract}, ShortStringsFilterer: ShortStringsFilterer{contract: contract}}, nil
}

// NewShortStringsCaller creates a new read-only instance of ShortStrings, bound to a specific deployed contract.
func NewShortStringsCaller(address common.Address, caller bind.ContractCaller) (*ShortStringsCaller, error) {
	contract, err := bindShortStrings(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ShortStringsCaller{contract: contract}, nil
}

// NewShortStringsTransactor creates a new write-only instance of ShortStrings, bound to a specific deployed contract.
func NewShortStringsTransactor(address common.Address, transactor bind.ContractTransactor) (*ShortStringsTransactor, error) {
	contract, err := bindShortStrings(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ShortStringsTransactor{contract: contract}, nil
}

// NewShortStringsFilterer creates a new log filterer instance of ShortStrings, bound to a specific deployed contract.
func NewShortStringsFilterer(address common.Address, filterer bind.ContractFilterer) (*ShortStringsFilterer, error) {
	contract, err := bindShortStrings(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ShortStringsFilterer{contract: contract}, nil
}

// bindShortStrings binds a generic wrapper to an already deployed contract.
func bindShortStrings(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ShortStringsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ShortStrings *ShortStringsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ShortStrings.Contract.ShortStringsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ShortStrings *ShortStringsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ShortStrings.Contract.ShortStringsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ShortStrings *ShortStringsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ShortStrings.Contract.ShortStringsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ShortStrings *ShortStringsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ShortStrings.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ShortStrings *ShortStringsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ShortStrings.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ShortStrings *ShortStringsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ShortStrings.Contract.contract.Transact(opts, method, params...)
}
