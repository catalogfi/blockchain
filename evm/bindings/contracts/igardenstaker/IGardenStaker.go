// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package igardenstaker

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

// IGardenStakerMetaData contains all meta data concerning the IGardenStaker contract.
var IGardenStakerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"DELEGATE_STAKE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SEED\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"newFiller\",\"type\":\"address\"}],\"name\":\"changeVote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"units\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockBlocks\",\"type\":\"uint256\"}],\"name\":\"vote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IGardenStakerABI is the input ABI used to generate the binding from.
// Deprecated: Use IGardenStakerMetaData.ABI instead.
var IGardenStakerABI = IGardenStakerMetaData.ABI

// IGardenStaker is an auto generated Go binding around an Ethereum contract.
type IGardenStaker struct {
	IGardenStakerCaller     // Read-only binding to the contract
	IGardenStakerTransactor // Write-only binding to the contract
	IGardenStakerFilterer   // Log filterer for contract events
}

// IGardenStakerCaller is an auto generated read-only Go binding around an Ethereum contract.
type IGardenStakerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IGardenStakerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IGardenStakerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IGardenStakerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IGardenStakerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IGardenStakerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IGardenStakerSession struct {
	Contract     *IGardenStaker    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IGardenStakerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IGardenStakerCallerSession struct {
	Contract *IGardenStakerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// IGardenStakerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IGardenStakerTransactorSession struct {
	Contract     *IGardenStakerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// IGardenStakerRaw is an auto generated low-level Go binding around an Ethereum contract.
type IGardenStakerRaw struct {
	Contract *IGardenStaker // Generic contract binding to access the raw methods on
}

// IGardenStakerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IGardenStakerCallerRaw struct {
	Contract *IGardenStakerCaller // Generic read-only contract binding to access the raw methods on
}

// IGardenStakerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IGardenStakerTransactorRaw struct {
	Contract *IGardenStakerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIGardenStaker creates a new instance of IGardenStaker, bound to a specific deployed contract.
func NewIGardenStaker(address common.Address, backend bind.ContractBackend) (*IGardenStaker, error) {
	contract, err := bindIGardenStaker(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IGardenStaker{IGardenStakerCaller: IGardenStakerCaller{contract: contract}, IGardenStakerTransactor: IGardenStakerTransactor{contract: contract}, IGardenStakerFilterer: IGardenStakerFilterer{contract: contract}}, nil
}

// NewIGardenStakerCaller creates a new read-only instance of IGardenStaker, bound to a specific deployed contract.
func NewIGardenStakerCaller(address common.Address, caller bind.ContractCaller) (*IGardenStakerCaller, error) {
	contract, err := bindIGardenStaker(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IGardenStakerCaller{contract: contract}, nil
}

// NewIGardenStakerTransactor creates a new write-only instance of IGardenStaker, bound to a specific deployed contract.
func NewIGardenStakerTransactor(address common.Address, transactor bind.ContractTransactor) (*IGardenStakerTransactor, error) {
	contract, err := bindIGardenStaker(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IGardenStakerTransactor{contract: contract}, nil
}

// NewIGardenStakerFilterer creates a new log filterer instance of IGardenStaker, bound to a specific deployed contract.
func NewIGardenStakerFilterer(address common.Address, filterer bind.ContractFilterer) (*IGardenStakerFilterer, error) {
	contract, err := bindIGardenStaker(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IGardenStakerFilterer{contract: contract}, nil
}

// bindIGardenStaker binds a generic wrapper to an already deployed contract.
func bindIGardenStaker(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IGardenStakerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IGardenStaker *IGardenStakerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IGardenStaker.Contract.IGardenStakerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IGardenStaker *IGardenStakerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IGardenStaker.Contract.IGardenStakerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IGardenStaker *IGardenStakerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IGardenStaker.Contract.IGardenStakerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IGardenStaker *IGardenStakerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IGardenStaker.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IGardenStaker *IGardenStakerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IGardenStaker.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IGardenStaker *IGardenStakerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IGardenStaker.Contract.contract.Transact(opts, method, params...)
}

// DELEGATESTAKE is a paid mutator transaction binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() returns(uint256)
func (_IGardenStaker *IGardenStakerTransactor) DELEGATESTAKE(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IGardenStaker.contract.Transact(opts, "DELEGATE_STAKE")
}

// DELEGATESTAKE is a paid mutator transaction binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() returns(uint256)
func (_IGardenStaker *IGardenStakerSession) DELEGATESTAKE() (*types.Transaction, error) {
	return _IGardenStaker.Contract.DELEGATESTAKE(&_IGardenStaker.TransactOpts)
}

// DELEGATESTAKE is a paid mutator transaction binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() returns(uint256)
func (_IGardenStaker *IGardenStakerTransactorSession) DELEGATESTAKE() (*types.Transaction, error) {
	return _IGardenStaker.Contract.DELEGATESTAKE(&_IGardenStaker.TransactOpts)
}

// SEED is a paid mutator transaction binding the contract method 0x0edc4737.
//
// Solidity: function SEED() returns(address)
func (_IGardenStaker *IGardenStakerTransactor) SEED(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IGardenStaker.contract.Transact(opts, "SEED")
}

// SEED is a paid mutator transaction binding the contract method 0x0edc4737.
//
// Solidity: function SEED() returns(address)
func (_IGardenStaker *IGardenStakerSession) SEED() (*types.Transaction, error) {
	return _IGardenStaker.Contract.SEED(&_IGardenStaker.TransactOpts)
}

// SEED is a paid mutator transaction binding the contract method 0x0edc4737.
//
// Solidity: function SEED() returns(address)
func (_IGardenStaker *IGardenStakerTransactorSession) SEED() (*types.Transaction, error) {
	return _IGardenStaker.Contract.SEED(&_IGardenStaker.TransactOpts)
}

// ChangeVote is a paid mutator transaction binding the contract method 0xb9b0b22b.
//
// Solidity: function changeVote(bytes32 stakeID, address newFiller) returns()
func (_IGardenStaker *IGardenStakerTransactor) ChangeVote(opts *bind.TransactOpts, stakeID [32]byte, newFiller common.Address) (*types.Transaction, error) {
	return _IGardenStaker.contract.Transact(opts, "changeVote", stakeID, newFiller)
}

// ChangeVote is a paid mutator transaction binding the contract method 0xb9b0b22b.
//
// Solidity: function changeVote(bytes32 stakeID, address newFiller) returns()
func (_IGardenStaker *IGardenStakerSession) ChangeVote(stakeID [32]byte, newFiller common.Address) (*types.Transaction, error) {
	return _IGardenStaker.Contract.ChangeVote(&_IGardenStaker.TransactOpts, stakeID, newFiller)
}

// ChangeVote is a paid mutator transaction binding the contract method 0xb9b0b22b.
//
// Solidity: function changeVote(bytes32 stakeID, address newFiller) returns()
func (_IGardenStaker *IGardenStakerTransactorSession) ChangeVote(stakeID [32]byte, newFiller common.Address) (*types.Transaction, error) {
	return _IGardenStaker.Contract.ChangeVote(&_IGardenStaker.TransactOpts, stakeID, newFiller)
}

// Vote is a paid mutator transaction binding the contract method 0x2a4a1b73.
//
// Solidity: function vote(address filler, uint256 units, uint256 lockBlocks) returns(bytes32)
func (_IGardenStaker *IGardenStakerTransactor) Vote(opts *bind.TransactOpts, filler common.Address, units *big.Int, lockBlocks *big.Int) (*types.Transaction, error) {
	return _IGardenStaker.contract.Transact(opts, "vote", filler, units, lockBlocks)
}

// Vote is a paid mutator transaction binding the contract method 0x2a4a1b73.
//
// Solidity: function vote(address filler, uint256 units, uint256 lockBlocks) returns(bytes32)
func (_IGardenStaker *IGardenStakerSession) Vote(filler common.Address, units *big.Int, lockBlocks *big.Int) (*types.Transaction, error) {
	return _IGardenStaker.Contract.Vote(&_IGardenStaker.TransactOpts, filler, units, lockBlocks)
}

// Vote is a paid mutator transaction binding the contract method 0x2a4a1b73.
//
// Solidity: function vote(address filler, uint256 units, uint256 lockBlocks) returns(bytes32)
func (_IGardenStaker *IGardenStakerTransactorSession) Vote(filler common.Address, units *big.Int, lockBlocks *big.Int) (*types.Transaction, error) {
	return _IGardenStaker.Contract.Vote(&_IGardenStaker.TransactOpts, filler, units, lockBlocks)
}
