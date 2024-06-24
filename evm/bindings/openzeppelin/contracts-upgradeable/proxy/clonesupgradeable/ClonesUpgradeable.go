// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package clonesupgradeable

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

// ClonesUpgradeableMetaData contains all meta data concerning the ClonesUpgradeable contract.
var ClonesUpgradeableMetaData = &bind.MetaData{
	ABI: "[]",
}

// ClonesUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use ClonesUpgradeableMetaData.ABI instead.
var ClonesUpgradeableABI = ClonesUpgradeableMetaData.ABI

// ClonesUpgradeable is an auto generated Go binding around an Ethereum contract.
type ClonesUpgradeable struct {
	ClonesUpgradeableCaller     // Read-only binding to the contract
	ClonesUpgradeableTransactor // Write-only binding to the contract
	ClonesUpgradeableFilterer   // Log filterer for contract events
}

// ClonesUpgradeableCaller is an auto generated read-only Go binding around an Ethereum contract.
type ClonesUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClonesUpgradeableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ClonesUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClonesUpgradeableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ClonesUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClonesUpgradeableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ClonesUpgradeableSession struct {
	Contract     *ClonesUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ClonesUpgradeableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ClonesUpgradeableCallerSession struct {
	Contract *ClonesUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ClonesUpgradeableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ClonesUpgradeableTransactorSession struct {
	Contract     *ClonesUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ClonesUpgradeableRaw is an auto generated low-level Go binding around an Ethereum contract.
type ClonesUpgradeableRaw struct {
	Contract *ClonesUpgradeable // Generic contract binding to access the raw methods on
}

// ClonesUpgradeableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ClonesUpgradeableCallerRaw struct {
	Contract *ClonesUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// ClonesUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ClonesUpgradeableTransactorRaw struct {
	Contract *ClonesUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewClonesUpgradeable creates a new instance of ClonesUpgradeable, bound to a specific deployed contract.
func NewClonesUpgradeable(address common.Address, backend bind.ContractBackend) (*ClonesUpgradeable, error) {
	contract, err := bindClonesUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ClonesUpgradeable{ClonesUpgradeableCaller: ClonesUpgradeableCaller{contract: contract}, ClonesUpgradeableTransactor: ClonesUpgradeableTransactor{contract: contract}, ClonesUpgradeableFilterer: ClonesUpgradeableFilterer{contract: contract}}, nil
}

// NewClonesUpgradeableCaller creates a new read-only instance of ClonesUpgradeable, bound to a specific deployed contract.
func NewClonesUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*ClonesUpgradeableCaller, error) {
	contract, err := bindClonesUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ClonesUpgradeableCaller{contract: contract}, nil
}

// NewClonesUpgradeableTransactor creates a new write-only instance of ClonesUpgradeable, bound to a specific deployed contract.
func NewClonesUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*ClonesUpgradeableTransactor, error) {
	contract, err := bindClonesUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ClonesUpgradeableTransactor{contract: contract}, nil
}

// NewClonesUpgradeableFilterer creates a new log filterer instance of ClonesUpgradeable, bound to a specific deployed contract.
func NewClonesUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*ClonesUpgradeableFilterer, error) {
	contract, err := bindClonesUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ClonesUpgradeableFilterer{contract: contract}, nil
}

// bindClonesUpgradeable binds a generic wrapper to an already deployed contract.
func bindClonesUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ClonesUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ClonesUpgradeable *ClonesUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ClonesUpgradeable.Contract.ClonesUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ClonesUpgradeable *ClonesUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClonesUpgradeable.Contract.ClonesUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ClonesUpgradeable *ClonesUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ClonesUpgradeable.Contract.ClonesUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ClonesUpgradeable *ClonesUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ClonesUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ClonesUpgradeable *ClonesUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClonesUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ClonesUpgradeable *ClonesUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ClonesUpgradeable.Contract.contract.Transact(opts, method, params...)
}
