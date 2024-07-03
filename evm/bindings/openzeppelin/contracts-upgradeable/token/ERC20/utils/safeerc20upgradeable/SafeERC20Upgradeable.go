// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package safeerc20upgradeable

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

// SafeERC20UpgradeableMetaData contains all meta data concerning the SafeERC20Upgradeable contract.
var SafeERC20UpgradeableMetaData = &bind.MetaData{
	ABI: "[]",
}

// SafeERC20UpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use SafeERC20UpgradeableMetaData.ABI instead.
var SafeERC20UpgradeableABI = SafeERC20UpgradeableMetaData.ABI

// SafeERC20Upgradeable is an auto generated Go binding around an Ethereum contract.
type SafeERC20Upgradeable struct {
	SafeERC20UpgradeableCaller     // Read-only binding to the contract
	SafeERC20UpgradeableTransactor // Write-only binding to the contract
	SafeERC20UpgradeableFilterer   // Log filterer for contract events
}

// SafeERC20UpgradeableCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeERC20UpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeERC20UpgradeableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeERC20UpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeERC20UpgradeableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeERC20UpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeERC20UpgradeableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeERC20UpgradeableSession struct {
	Contract     *SafeERC20Upgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// SafeERC20UpgradeableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeERC20UpgradeableCallerSession struct {
	Contract *SafeERC20UpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// SafeERC20UpgradeableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeERC20UpgradeableTransactorSession struct {
	Contract     *SafeERC20UpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// SafeERC20UpgradeableRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeERC20UpgradeableRaw struct {
	Contract *SafeERC20Upgradeable // Generic contract binding to access the raw methods on
}

// SafeERC20UpgradeableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeERC20UpgradeableCallerRaw struct {
	Contract *SafeERC20UpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// SafeERC20UpgradeableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeERC20UpgradeableTransactorRaw struct {
	Contract *SafeERC20UpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeERC20Upgradeable creates a new instance of SafeERC20Upgradeable, bound to a specific deployed contract.
func NewSafeERC20Upgradeable(address common.Address, backend bind.ContractBackend) (*SafeERC20Upgradeable, error) {
	contract, err := bindSafeERC20Upgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeERC20Upgradeable{SafeERC20UpgradeableCaller: SafeERC20UpgradeableCaller{contract: contract}, SafeERC20UpgradeableTransactor: SafeERC20UpgradeableTransactor{contract: contract}, SafeERC20UpgradeableFilterer: SafeERC20UpgradeableFilterer{contract: contract}}, nil
}

// NewSafeERC20UpgradeableCaller creates a new read-only instance of SafeERC20Upgradeable, bound to a specific deployed contract.
func NewSafeERC20UpgradeableCaller(address common.Address, caller bind.ContractCaller) (*SafeERC20UpgradeableCaller, error) {
	contract, err := bindSafeERC20Upgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeERC20UpgradeableCaller{contract: contract}, nil
}

// NewSafeERC20UpgradeableTransactor creates a new write-only instance of SafeERC20Upgradeable, bound to a specific deployed contract.
func NewSafeERC20UpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeERC20UpgradeableTransactor, error) {
	contract, err := bindSafeERC20Upgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeERC20UpgradeableTransactor{contract: contract}, nil
}

// NewSafeERC20UpgradeableFilterer creates a new log filterer instance of SafeERC20Upgradeable, bound to a specific deployed contract.
func NewSafeERC20UpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeERC20UpgradeableFilterer, error) {
	contract, err := bindSafeERC20Upgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeERC20UpgradeableFilterer{contract: contract}, nil
}

// bindSafeERC20Upgradeable binds a generic wrapper to an already deployed contract.
func bindSafeERC20Upgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SafeERC20UpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeERC20Upgradeable *SafeERC20UpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeERC20Upgradeable.Contract.SafeERC20UpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeERC20Upgradeable *SafeERC20UpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeERC20Upgradeable.Contract.SafeERC20UpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeERC20Upgradeable *SafeERC20UpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeERC20Upgradeable.Contract.SafeERC20UpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeERC20Upgradeable *SafeERC20UpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeERC20Upgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeERC20Upgradeable *SafeERC20UpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeERC20Upgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeERC20Upgradeable *SafeERC20UpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeERC20Upgradeable.Contract.contract.Transact(opts, method, params...)
}
