// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package signedmathupgradeable

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

// SignedMathUpgradeableMetaData contains all meta data concerning the SignedMathUpgradeable contract.
var SignedMathUpgradeableMetaData = &bind.MetaData{
	ABI: "[]",
}

// SignedMathUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use SignedMathUpgradeableMetaData.ABI instead.
var SignedMathUpgradeableABI = SignedMathUpgradeableMetaData.ABI

// SignedMathUpgradeable is an auto generated Go binding around an Ethereum contract.
type SignedMathUpgradeable struct {
	SignedMathUpgradeableCaller     // Read-only binding to the contract
	SignedMathUpgradeableTransactor // Write-only binding to the contract
	SignedMathUpgradeableFilterer   // Log filterer for contract events
}

// SignedMathUpgradeableCaller is an auto generated read-only Go binding around an Ethereum contract.
type SignedMathUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignedMathUpgradeableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SignedMathUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignedMathUpgradeableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SignedMathUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignedMathUpgradeableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SignedMathUpgradeableSession struct {
	Contract     *SignedMathUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// SignedMathUpgradeableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SignedMathUpgradeableCallerSession struct {
	Contract *SignedMathUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// SignedMathUpgradeableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SignedMathUpgradeableTransactorSession struct {
	Contract     *SignedMathUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// SignedMathUpgradeableRaw is an auto generated low-level Go binding around an Ethereum contract.
type SignedMathUpgradeableRaw struct {
	Contract *SignedMathUpgradeable // Generic contract binding to access the raw methods on
}

// SignedMathUpgradeableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SignedMathUpgradeableCallerRaw struct {
	Contract *SignedMathUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// SignedMathUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SignedMathUpgradeableTransactorRaw struct {
	Contract *SignedMathUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSignedMathUpgradeable creates a new instance of SignedMathUpgradeable, bound to a specific deployed contract.
func NewSignedMathUpgradeable(address common.Address, backend bind.ContractBackend) (*SignedMathUpgradeable, error) {
	contract, err := bindSignedMathUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SignedMathUpgradeable{SignedMathUpgradeableCaller: SignedMathUpgradeableCaller{contract: contract}, SignedMathUpgradeableTransactor: SignedMathUpgradeableTransactor{contract: contract}, SignedMathUpgradeableFilterer: SignedMathUpgradeableFilterer{contract: contract}}, nil
}

// NewSignedMathUpgradeableCaller creates a new read-only instance of SignedMathUpgradeable, bound to a specific deployed contract.
func NewSignedMathUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*SignedMathUpgradeableCaller, error) {
	contract, err := bindSignedMathUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SignedMathUpgradeableCaller{contract: contract}, nil
}

// NewSignedMathUpgradeableTransactor creates a new write-only instance of SignedMathUpgradeable, bound to a specific deployed contract.
func NewSignedMathUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*SignedMathUpgradeableTransactor, error) {
	contract, err := bindSignedMathUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SignedMathUpgradeableTransactor{contract: contract}, nil
}

// NewSignedMathUpgradeableFilterer creates a new log filterer instance of SignedMathUpgradeable, bound to a specific deployed contract.
func NewSignedMathUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*SignedMathUpgradeableFilterer, error) {
	contract, err := bindSignedMathUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SignedMathUpgradeableFilterer{contract: contract}, nil
}

// bindSignedMathUpgradeable binds a generic wrapper to an already deployed contract.
func bindSignedMathUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SignedMathUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SignedMathUpgradeable *SignedMathUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SignedMathUpgradeable.Contract.SignedMathUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SignedMathUpgradeable *SignedMathUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignedMathUpgradeable.Contract.SignedMathUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SignedMathUpgradeable *SignedMathUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignedMathUpgradeable.Contract.SignedMathUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SignedMathUpgradeable *SignedMathUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SignedMathUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SignedMathUpgradeable *SignedMathUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignedMathUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SignedMathUpgradeable *SignedMathUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignedMathUpgradeable.Contract.contract.Transact(opts, method, params...)
}
