// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mathupgradeable

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

// MathUpgradeableMetaData contains all meta data concerning the MathUpgradeable contract.
var MathUpgradeableMetaData = &bind.MetaData{
	ABI: "[]",
}

// MathUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use MathUpgradeableMetaData.ABI instead.
var MathUpgradeableABI = MathUpgradeableMetaData.ABI

// MathUpgradeable is an auto generated Go binding around an Ethereum contract.
type MathUpgradeable struct {
	MathUpgradeableCaller     // Read-only binding to the contract
	MathUpgradeableTransactor // Write-only binding to the contract
	MathUpgradeableFilterer   // Log filterer for contract events
}

// MathUpgradeableCaller is an auto generated read-only Go binding around an Ethereum contract.
type MathUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MathUpgradeableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MathUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MathUpgradeableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MathUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MathUpgradeableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MathUpgradeableSession struct {
	Contract     *MathUpgradeable  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MathUpgradeableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MathUpgradeableCallerSession struct {
	Contract *MathUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// MathUpgradeableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MathUpgradeableTransactorSession struct {
	Contract     *MathUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// MathUpgradeableRaw is an auto generated low-level Go binding around an Ethereum contract.
type MathUpgradeableRaw struct {
	Contract *MathUpgradeable // Generic contract binding to access the raw methods on
}

// MathUpgradeableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MathUpgradeableCallerRaw struct {
	Contract *MathUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// MathUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MathUpgradeableTransactorRaw struct {
	Contract *MathUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMathUpgradeable creates a new instance of MathUpgradeable, bound to a specific deployed contract.
func NewMathUpgradeable(address common.Address, backend bind.ContractBackend) (*MathUpgradeable, error) {
	contract, err := bindMathUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MathUpgradeable{MathUpgradeableCaller: MathUpgradeableCaller{contract: contract}, MathUpgradeableTransactor: MathUpgradeableTransactor{contract: contract}, MathUpgradeableFilterer: MathUpgradeableFilterer{contract: contract}}, nil
}

// NewMathUpgradeableCaller creates a new read-only instance of MathUpgradeable, bound to a specific deployed contract.
func NewMathUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*MathUpgradeableCaller, error) {
	contract, err := bindMathUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MathUpgradeableCaller{contract: contract}, nil
}

// NewMathUpgradeableTransactor creates a new write-only instance of MathUpgradeable, bound to a specific deployed contract.
func NewMathUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*MathUpgradeableTransactor, error) {
	contract, err := bindMathUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MathUpgradeableTransactor{contract: contract}, nil
}

// NewMathUpgradeableFilterer creates a new log filterer instance of MathUpgradeable, bound to a specific deployed contract.
func NewMathUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*MathUpgradeableFilterer, error) {
	contract, err := bindMathUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MathUpgradeableFilterer{contract: contract}, nil
}

// bindMathUpgradeable binds a generic wrapper to an already deployed contract.
func bindMathUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MathUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MathUpgradeable *MathUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MathUpgradeable.Contract.MathUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MathUpgradeable *MathUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MathUpgradeable.Contract.MathUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MathUpgradeable *MathUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MathUpgradeable.Contract.MathUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MathUpgradeable *MathUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MathUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MathUpgradeable *MathUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MathUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MathUpgradeable *MathUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MathUpgradeable.Contract.contract.Transact(opts, method, params...)
}
