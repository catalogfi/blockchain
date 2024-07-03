// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ecdsaupgradeable

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

// ECDSAUpgradeableMetaData contains all meta data concerning the ECDSAUpgradeable contract.
var ECDSAUpgradeableMetaData = &bind.MetaData{
	ABI: "[]",
}

// ECDSAUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use ECDSAUpgradeableMetaData.ABI instead.
var ECDSAUpgradeableABI = ECDSAUpgradeableMetaData.ABI

// ECDSAUpgradeable is an auto generated Go binding around an Ethereum contract.
type ECDSAUpgradeable struct {
	ECDSAUpgradeableCaller     // Read-only binding to the contract
	ECDSAUpgradeableTransactor // Write-only binding to the contract
	ECDSAUpgradeableFilterer   // Log filterer for contract events
}

// ECDSAUpgradeableCaller is an auto generated read-only Go binding around an Ethereum contract.
type ECDSAUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECDSAUpgradeableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ECDSAUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECDSAUpgradeableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ECDSAUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECDSAUpgradeableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ECDSAUpgradeableSession struct {
	Contract     *ECDSAUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ECDSAUpgradeableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ECDSAUpgradeableCallerSession struct {
	Contract *ECDSAUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ECDSAUpgradeableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ECDSAUpgradeableTransactorSession struct {
	Contract     *ECDSAUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ECDSAUpgradeableRaw is an auto generated low-level Go binding around an Ethereum contract.
type ECDSAUpgradeableRaw struct {
	Contract *ECDSAUpgradeable // Generic contract binding to access the raw methods on
}

// ECDSAUpgradeableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ECDSAUpgradeableCallerRaw struct {
	Contract *ECDSAUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// ECDSAUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ECDSAUpgradeableTransactorRaw struct {
	Contract *ECDSAUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewECDSAUpgradeable creates a new instance of ECDSAUpgradeable, bound to a specific deployed contract.
func NewECDSAUpgradeable(address common.Address, backend bind.ContractBackend) (*ECDSAUpgradeable, error) {
	contract, err := bindECDSAUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ECDSAUpgradeable{ECDSAUpgradeableCaller: ECDSAUpgradeableCaller{contract: contract}, ECDSAUpgradeableTransactor: ECDSAUpgradeableTransactor{contract: contract}, ECDSAUpgradeableFilterer: ECDSAUpgradeableFilterer{contract: contract}}, nil
}

// NewECDSAUpgradeableCaller creates a new read-only instance of ECDSAUpgradeable, bound to a specific deployed contract.
func NewECDSAUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*ECDSAUpgradeableCaller, error) {
	contract, err := bindECDSAUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ECDSAUpgradeableCaller{contract: contract}, nil
}

// NewECDSAUpgradeableTransactor creates a new write-only instance of ECDSAUpgradeable, bound to a specific deployed contract.
func NewECDSAUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*ECDSAUpgradeableTransactor, error) {
	contract, err := bindECDSAUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ECDSAUpgradeableTransactor{contract: contract}, nil
}

// NewECDSAUpgradeableFilterer creates a new log filterer instance of ECDSAUpgradeable, bound to a specific deployed contract.
func NewECDSAUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*ECDSAUpgradeableFilterer, error) {
	contract, err := bindECDSAUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ECDSAUpgradeableFilterer{contract: contract}, nil
}

// bindECDSAUpgradeable binds a generic wrapper to an already deployed contract.
func bindECDSAUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ECDSAUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ECDSAUpgradeable *ECDSAUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ECDSAUpgradeable.Contract.ECDSAUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ECDSAUpgradeable *ECDSAUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ECDSAUpgradeable.Contract.ECDSAUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ECDSAUpgradeable *ECDSAUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ECDSAUpgradeable.Contract.ECDSAUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ECDSAUpgradeable *ECDSAUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ECDSAUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ECDSAUpgradeable *ECDSAUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ECDSAUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ECDSAUpgradeable *ECDSAUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ECDSAUpgradeable.Contract.contract.Transact(opts, method, params...)
}
