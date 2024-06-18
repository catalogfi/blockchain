// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stringsupgradeable

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

// StringsUpgradeableMetaData contains all meta data concerning the StringsUpgradeable contract.
var StringsUpgradeableMetaData = &bind.MetaData{
	ABI: "[]",
}

// StringsUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use StringsUpgradeableMetaData.ABI instead.
var StringsUpgradeableABI = StringsUpgradeableMetaData.ABI

// StringsUpgradeable is an auto generated Go binding around an Ethereum contract.
type StringsUpgradeable struct {
	StringsUpgradeableCaller     // Read-only binding to the contract
	StringsUpgradeableTransactor // Write-only binding to the contract
	StringsUpgradeableFilterer   // Log filterer for contract events
}

// StringsUpgradeableCaller is an auto generated read-only Go binding around an Ethereum contract.
type StringsUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StringsUpgradeableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StringsUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StringsUpgradeableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StringsUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StringsUpgradeableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StringsUpgradeableSession struct {
	Contract     *StringsUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// StringsUpgradeableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StringsUpgradeableCallerSession struct {
	Contract *StringsUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// StringsUpgradeableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StringsUpgradeableTransactorSession struct {
	Contract     *StringsUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// StringsUpgradeableRaw is an auto generated low-level Go binding around an Ethereum contract.
type StringsUpgradeableRaw struct {
	Contract *StringsUpgradeable // Generic contract binding to access the raw methods on
}

// StringsUpgradeableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StringsUpgradeableCallerRaw struct {
	Contract *StringsUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// StringsUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StringsUpgradeableTransactorRaw struct {
	Contract *StringsUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStringsUpgradeable creates a new instance of StringsUpgradeable, bound to a specific deployed contract.
func NewStringsUpgradeable(address common.Address, backend bind.ContractBackend) (*StringsUpgradeable, error) {
	contract, err := bindStringsUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StringsUpgradeable{StringsUpgradeableCaller: StringsUpgradeableCaller{contract: contract}, StringsUpgradeableTransactor: StringsUpgradeableTransactor{contract: contract}, StringsUpgradeableFilterer: StringsUpgradeableFilterer{contract: contract}}, nil
}

// NewStringsUpgradeableCaller creates a new read-only instance of StringsUpgradeable, bound to a specific deployed contract.
func NewStringsUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*StringsUpgradeableCaller, error) {
	contract, err := bindStringsUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StringsUpgradeableCaller{contract: contract}, nil
}

// NewStringsUpgradeableTransactor creates a new write-only instance of StringsUpgradeable, bound to a specific deployed contract.
func NewStringsUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*StringsUpgradeableTransactor, error) {
	contract, err := bindStringsUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StringsUpgradeableTransactor{contract: contract}, nil
}

// NewStringsUpgradeableFilterer creates a new log filterer instance of StringsUpgradeable, bound to a specific deployed contract.
func NewStringsUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*StringsUpgradeableFilterer, error) {
	contract, err := bindStringsUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StringsUpgradeableFilterer{contract: contract}, nil
}

// bindStringsUpgradeable binds a generic wrapper to an already deployed contract.
func bindStringsUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StringsUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StringsUpgradeable *StringsUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StringsUpgradeable.Contract.StringsUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StringsUpgradeable *StringsUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StringsUpgradeable.Contract.StringsUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StringsUpgradeable *StringsUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StringsUpgradeable.Contract.StringsUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StringsUpgradeable *StringsUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StringsUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StringsUpgradeable *StringsUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StringsUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StringsUpgradeable *StringsUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StringsUpgradeable.Contract.contract.Transact(opts, method, params...)
}
