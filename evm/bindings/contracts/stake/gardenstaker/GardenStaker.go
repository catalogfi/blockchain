// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gardenstaker

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

// GardenStakerMetaData contains all meta data concerning the GardenStaker contract.
var GardenStakerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"seed\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"delegateStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fillerStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fillerCooldown\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deregisteredAt\",\"type\":\"uint256\"}],\"name\":\"FillerDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"FillerFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"}],\"name\":\"FillerRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"}],\"name\":\"FillerRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newLockBlocks\",\"type\":\"uint256\"}],\"name\":\"StakeExtended\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"}],\"name\":\"StakeRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newLockBlocks\",\"type\":\"uint256\"}],\"name\":\"StakeRenewed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"}],\"name\":\"Voted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldFiller\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newFiller\",\"type\":\"address\"}],\"name\":\"VotesChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DELEGATE_STAKE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FILLER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FILLER_COOL_DOWN\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FILLER_STAKE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SEED\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"newFiller\",\"type\":\"address\"}],\"name\":\"changeVote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"delegateNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deregister\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"newLockBlocks\",\"type\":\"uint256\"}],\"name\":\"extend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"}],\"name\":\"getFiller\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"feeInBips\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deregisteredAt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"delegateStakeIDs\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"}],\"name\":\"getVotes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"filler_\",\"type\":\"address\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"newLockBlocks\",\"type\":\"uint256\"}],\"name\":\"renew\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"stakes\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"units\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"newFee\",\"type\":\"uint16\"}],\"name\":\"updateFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"units\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockBlocks\",\"type\":\"uint256\"}],\"name\":\"vote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// GardenStakerABI is the input ABI used to generate the binding from.
// Deprecated: Use GardenStakerMetaData.ABI instead.
var GardenStakerABI = GardenStakerMetaData.ABI

// GardenStaker is an auto generated Go binding around an Ethereum contract.
type GardenStaker struct {
	GardenStakerCaller     // Read-only binding to the contract
	GardenStakerTransactor // Write-only binding to the contract
	GardenStakerFilterer   // Log filterer for contract events
}

// GardenStakerCaller is an auto generated read-only Go binding around an Ethereum contract.
type GardenStakerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenStakerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GardenStakerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenStakerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GardenStakerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenStakerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GardenStakerSession struct {
	Contract     *GardenStaker     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GardenStakerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GardenStakerCallerSession struct {
	Contract *GardenStakerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// GardenStakerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GardenStakerTransactorSession struct {
	Contract     *GardenStakerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// GardenStakerRaw is an auto generated low-level Go binding around an Ethereum contract.
type GardenStakerRaw struct {
	Contract *GardenStaker // Generic contract binding to access the raw methods on
}

// GardenStakerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GardenStakerCallerRaw struct {
	Contract *GardenStakerCaller // Generic read-only contract binding to access the raw methods on
}

// GardenStakerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GardenStakerTransactorRaw struct {
	Contract *GardenStakerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGardenStaker creates a new instance of GardenStaker, bound to a specific deployed contract.
func NewGardenStaker(address common.Address, backend bind.ContractBackend) (*GardenStaker, error) {
	contract, err := bindGardenStaker(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GardenStaker{GardenStakerCaller: GardenStakerCaller{contract: contract}, GardenStakerTransactor: GardenStakerTransactor{contract: contract}, GardenStakerFilterer: GardenStakerFilterer{contract: contract}}, nil
}

// NewGardenStakerCaller creates a new read-only instance of GardenStaker, bound to a specific deployed contract.
func NewGardenStakerCaller(address common.Address, caller bind.ContractCaller) (*GardenStakerCaller, error) {
	contract, err := bindGardenStaker(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GardenStakerCaller{contract: contract}, nil
}

// NewGardenStakerTransactor creates a new write-only instance of GardenStaker, bound to a specific deployed contract.
func NewGardenStakerTransactor(address common.Address, transactor bind.ContractTransactor) (*GardenStakerTransactor, error) {
	contract, err := bindGardenStaker(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GardenStakerTransactor{contract: contract}, nil
}

// NewGardenStakerFilterer creates a new log filterer instance of GardenStaker, bound to a specific deployed contract.
func NewGardenStakerFilterer(address common.Address, filterer bind.ContractFilterer) (*GardenStakerFilterer, error) {
	contract, err := bindGardenStaker(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GardenStakerFilterer{contract: contract}, nil
}

// bindGardenStaker binds a generic wrapper to an already deployed contract.
func bindGardenStaker(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GardenStakerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GardenStaker *GardenStakerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GardenStaker.Contract.GardenStakerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GardenStaker *GardenStakerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenStaker.Contract.GardenStakerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GardenStaker *GardenStakerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GardenStaker.Contract.GardenStakerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GardenStaker *GardenStakerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GardenStaker.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GardenStaker *GardenStakerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenStaker.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GardenStaker *GardenStakerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GardenStaker.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GardenStaker *GardenStakerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GardenStaker.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GardenStaker *GardenStakerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _GardenStaker.Contract.DEFAULTADMINROLE(&_GardenStaker.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GardenStaker *GardenStakerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _GardenStaker.Contract.DEFAULTADMINROLE(&_GardenStaker.CallOpts)
}

// DELEGATESTAKE is a free data retrieval call binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() view returns(uint256)
func (_GardenStaker *GardenStakerCaller) DELEGATESTAKE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GardenStaker.contract.Call(opts, &out, "DELEGATE_STAKE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DELEGATESTAKE is a free data retrieval call binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() view returns(uint256)
func (_GardenStaker *GardenStakerSession) DELEGATESTAKE() (*big.Int, error) {
	return _GardenStaker.Contract.DELEGATESTAKE(&_GardenStaker.CallOpts)
}

// DELEGATESTAKE is a free data retrieval call binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() view returns(uint256)
func (_GardenStaker *GardenStakerCallerSession) DELEGATESTAKE() (*big.Int, error) {
	return _GardenStaker.Contract.DELEGATESTAKE(&_GardenStaker.CallOpts)
}

// FILLER is a free data retrieval call binding the contract method 0x64645bae.
//
// Solidity: function FILLER() view returns(bytes32)
func (_GardenStaker *GardenStakerCaller) FILLER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GardenStaker.contract.Call(opts, &out, "FILLER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FILLER is a free data retrieval call binding the contract method 0x64645bae.
//
// Solidity: function FILLER() view returns(bytes32)
func (_GardenStaker *GardenStakerSession) FILLER() ([32]byte, error) {
	return _GardenStaker.Contract.FILLER(&_GardenStaker.CallOpts)
}

// FILLER is a free data retrieval call binding the contract method 0x64645bae.
//
// Solidity: function FILLER() view returns(bytes32)
func (_GardenStaker *GardenStakerCallerSession) FILLER() ([32]byte, error) {
	return _GardenStaker.Contract.FILLER(&_GardenStaker.CallOpts)
}

// FILLERCOOLDOWN is a free data retrieval call binding the contract method 0x015a44e8.
//
// Solidity: function FILLER_COOL_DOWN() view returns(uint256)
func (_GardenStaker *GardenStakerCaller) FILLERCOOLDOWN(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GardenStaker.contract.Call(opts, &out, "FILLER_COOL_DOWN")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FILLERCOOLDOWN is a free data retrieval call binding the contract method 0x015a44e8.
//
// Solidity: function FILLER_COOL_DOWN() view returns(uint256)
func (_GardenStaker *GardenStakerSession) FILLERCOOLDOWN() (*big.Int, error) {
	return _GardenStaker.Contract.FILLERCOOLDOWN(&_GardenStaker.CallOpts)
}

// FILLERCOOLDOWN is a free data retrieval call binding the contract method 0x015a44e8.
//
// Solidity: function FILLER_COOL_DOWN() view returns(uint256)
func (_GardenStaker *GardenStakerCallerSession) FILLERCOOLDOWN() (*big.Int, error) {
	return _GardenStaker.Contract.FILLERCOOLDOWN(&_GardenStaker.CallOpts)
}

// FILLERSTAKE is a free data retrieval call binding the contract method 0xecfe923e.
//
// Solidity: function FILLER_STAKE() view returns(uint256)
func (_GardenStaker *GardenStakerCaller) FILLERSTAKE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GardenStaker.contract.Call(opts, &out, "FILLER_STAKE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FILLERSTAKE is a free data retrieval call binding the contract method 0xecfe923e.
//
// Solidity: function FILLER_STAKE() view returns(uint256)
func (_GardenStaker *GardenStakerSession) FILLERSTAKE() (*big.Int, error) {
	return _GardenStaker.Contract.FILLERSTAKE(&_GardenStaker.CallOpts)
}

// FILLERSTAKE is a free data retrieval call binding the contract method 0xecfe923e.
//
// Solidity: function FILLER_STAKE() view returns(uint256)
func (_GardenStaker *GardenStakerCallerSession) FILLERSTAKE() (*big.Int, error) {
	return _GardenStaker.Contract.FILLERSTAKE(&_GardenStaker.CallOpts)
}

// SEED is a free data retrieval call binding the contract method 0x0edc4737.
//
// Solidity: function SEED() view returns(address)
func (_GardenStaker *GardenStakerCaller) SEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GardenStaker.contract.Call(opts, &out, "SEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SEED is a free data retrieval call binding the contract method 0x0edc4737.
//
// Solidity: function SEED() view returns(address)
func (_GardenStaker *GardenStakerSession) SEED() (common.Address, error) {
	return _GardenStaker.Contract.SEED(&_GardenStaker.CallOpts)
}

// SEED is a free data retrieval call binding the contract method 0x0edc4737.
//
// Solidity: function SEED() view returns(address)
func (_GardenStaker *GardenStakerCallerSession) SEED() (common.Address, error) {
	return _GardenStaker.Contract.SEED(&_GardenStaker.CallOpts)
}

// DelegateNonce is a free data retrieval call binding the contract method 0x92e1f102.
//
// Solidity: function delegateNonce(address ) view returns(uint256)
func (_GardenStaker *GardenStakerCaller) DelegateNonce(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GardenStaker.contract.Call(opts, &out, "delegateNonce", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DelegateNonce is a free data retrieval call binding the contract method 0x92e1f102.
//
// Solidity: function delegateNonce(address ) view returns(uint256)
func (_GardenStaker *GardenStakerSession) DelegateNonce(arg0 common.Address) (*big.Int, error) {
	return _GardenStaker.Contract.DelegateNonce(&_GardenStaker.CallOpts, arg0)
}

// DelegateNonce is a free data retrieval call binding the contract method 0x92e1f102.
//
// Solidity: function delegateNonce(address ) view returns(uint256)
func (_GardenStaker *GardenStakerCallerSession) DelegateNonce(arg0 common.Address) (*big.Int, error) {
	return _GardenStaker.Contract.DelegateNonce(&_GardenStaker.CallOpts, arg0)
}

// GetFiller is a free data retrieval call binding the contract method 0x0b3ff53f.
//
// Solidity: function getFiller(address filler) view returns(uint16 feeInBips, uint256 stake, uint256 deregisteredAt, bytes32[] delegateStakeIDs)
func (_GardenStaker *GardenStakerCaller) GetFiller(opts *bind.CallOpts, filler common.Address) (struct {
	FeeInBips        uint16
	Stake            *big.Int
	DeregisteredAt   *big.Int
	DelegateStakeIDs [][32]byte
}, error) {
	var out []interface{}
	err := _GardenStaker.contract.Call(opts, &out, "getFiller", filler)

	outstruct := new(struct {
		FeeInBips        uint16
		Stake            *big.Int
		DeregisteredAt   *big.Int
		DelegateStakeIDs [][32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.FeeInBips = *abi.ConvertType(out[0], new(uint16)).(*uint16)
	outstruct.Stake = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.DeregisteredAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.DelegateStakeIDs = *abi.ConvertType(out[3], new([][32]byte)).(*[][32]byte)

	return *outstruct, err

}

// GetFiller is a free data retrieval call binding the contract method 0x0b3ff53f.
//
// Solidity: function getFiller(address filler) view returns(uint16 feeInBips, uint256 stake, uint256 deregisteredAt, bytes32[] delegateStakeIDs)
func (_GardenStaker *GardenStakerSession) GetFiller(filler common.Address) (struct {
	FeeInBips        uint16
	Stake            *big.Int
	DeregisteredAt   *big.Int
	DelegateStakeIDs [][32]byte
}, error) {
	return _GardenStaker.Contract.GetFiller(&_GardenStaker.CallOpts, filler)
}

// GetFiller is a free data retrieval call binding the contract method 0x0b3ff53f.
//
// Solidity: function getFiller(address filler) view returns(uint16 feeInBips, uint256 stake, uint256 deregisteredAt, bytes32[] delegateStakeIDs)
func (_GardenStaker *GardenStakerCallerSession) GetFiller(filler common.Address) (struct {
	FeeInBips        uint16
	Stake            *big.Int
	DeregisteredAt   *big.Int
	DelegateStakeIDs [][32]byte
}, error) {
	return _GardenStaker.Contract.GetFiller(&_GardenStaker.CallOpts, filler)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GardenStaker *GardenStakerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _GardenStaker.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GardenStaker *GardenStakerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _GardenStaker.Contract.GetRoleAdmin(&_GardenStaker.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GardenStaker *GardenStakerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _GardenStaker.Contract.GetRoleAdmin(&_GardenStaker.CallOpts, role)
}

// GetVotes is a free data retrieval call binding the contract method 0x9ab24eb0.
//
// Solidity: function getVotes(address filler) view returns(uint256 voteCount)
func (_GardenStaker *GardenStakerCaller) GetVotes(opts *bind.CallOpts, filler common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GardenStaker.contract.Call(opts, &out, "getVotes", filler)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotes is a free data retrieval call binding the contract method 0x9ab24eb0.
//
// Solidity: function getVotes(address filler) view returns(uint256 voteCount)
func (_GardenStaker *GardenStakerSession) GetVotes(filler common.Address) (*big.Int, error) {
	return _GardenStaker.Contract.GetVotes(&_GardenStaker.CallOpts, filler)
}

// GetVotes is a free data retrieval call binding the contract method 0x9ab24eb0.
//
// Solidity: function getVotes(address filler) view returns(uint256 voteCount)
func (_GardenStaker *GardenStakerCallerSession) GetVotes(filler common.Address) (*big.Int, error) {
	return _GardenStaker.Contract.GetVotes(&_GardenStaker.CallOpts, filler)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GardenStaker *GardenStakerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _GardenStaker.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GardenStaker *GardenStakerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _GardenStaker.Contract.HasRole(&_GardenStaker.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GardenStaker *GardenStakerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _GardenStaker.Contract.HasRole(&_GardenStaker.CallOpts, role, account)
}

// Stakes is a free data retrieval call binding the contract method 0x8fee6407.
//
// Solidity: function stakes(bytes32 ) view returns(address owner, uint256 stake, uint256 units, uint256 votes, address filler, uint256 expiry)
func (_GardenStaker *GardenStakerCaller) Stakes(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Owner  common.Address
	Stake  *big.Int
	Units  *big.Int
	Votes  *big.Int
	Filler common.Address
	Expiry *big.Int
}, error) {
	var out []interface{}
	err := _GardenStaker.contract.Call(opts, &out, "stakes", arg0)

	outstruct := new(struct {
		Owner  common.Address
		Stake  *big.Int
		Units  *big.Int
		Votes  *big.Int
		Filler common.Address
		Expiry *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Owner = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Stake = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Units = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Votes = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Filler = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Expiry = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Stakes is a free data retrieval call binding the contract method 0x8fee6407.
//
// Solidity: function stakes(bytes32 ) view returns(address owner, uint256 stake, uint256 units, uint256 votes, address filler, uint256 expiry)
func (_GardenStaker *GardenStakerSession) Stakes(arg0 [32]byte) (struct {
	Owner  common.Address
	Stake  *big.Int
	Units  *big.Int
	Votes  *big.Int
	Filler common.Address
	Expiry *big.Int
}, error) {
	return _GardenStaker.Contract.Stakes(&_GardenStaker.CallOpts, arg0)
}

// Stakes is a free data retrieval call binding the contract method 0x8fee6407.
//
// Solidity: function stakes(bytes32 ) view returns(address owner, uint256 stake, uint256 units, uint256 votes, address filler, uint256 expiry)
func (_GardenStaker *GardenStakerCallerSession) Stakes(arg0 [32]byte) (struct {
	Owner  common.Address
	Stake  *big.Int
	Units  *big.Int
	Votes  *big.Int
	Filler common.Address
	Expiry *big.Int
}, error) {
	return _GardenStaker.Contract.Stakes(&_GardenStaker.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GardenStaker *GardenStakerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _GardenStaker.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GardenStaker *GardenStakerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GardenStaker.Contract.SupportsInterface(&_GardenStaker.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GardenStaker *GardenStakerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GardenStaker.Contract.SupportsInterface(&_GardenStaker.CallOpts, interfaceId)
}

// ChangeVote is a paid mutator transaction binding the contract method 0xb9b0b22b.
//
// Solidity: function changeVote(bytes32 stakeID, address newFiller) returns()
func (_GardenStaker *GardenStakerTransactor) ChangeVote(opts *bind.TransactOpts, stakeID [32]byte, newFiller common.Address) (*types.Transaction, error) {
	return _GardenStaker.contract.Transact(opts, "changeVote", stakeID, newFiller)
}

// ChangeVote is a paid mutator transaction binding the contract method 0xb9b0b22b.
//
// Solidity: function changeVote(bytes32 stakeID, address newFiller) returns()
func (_GardenStaker *GardenStakerSession) ChangeVote(stakeID [32]byte, newFiller common.Address) (*types.Transaction, error) {
	return _GardenStaker.Contract.ChangeVote(&_GardenStaker.TransactOpts, stakeID, newFiller)
}

// ChangeVote is a paid mutator transaction binding the contract method 0xb9b0b22b.
//
// Solidity: function changeVote(bytes32 stakeID, address newFiller) returns()
func (_GardenStaker *GardenStakerTransactorSession) ChangeVote(stakeID [32]byte, newFiller common.Address) (*types.Transaction, error) {
	return _GardenStaker.Contract.ChangeVote(&_GardenStaker.TransactOpts, stakeID, newFiller)
}

// Deregister is a paid mutator transaction binding the contract method 0xaff5edb1.
//
// Solidity: function deregister() returns()
func (_GardenStaker *GardenStakerTransactor) Deregister(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenStaker.contract.Transact(opts, "deregister")
}

// Deregister is a paid mutator transaction binding the contract method 0xaff5edb1.
//
// Solidity: function deregister() returns()
func (_GardenStaker *GardenStakerSession) Deregister() (*types.Transaction, error) {
	return _GardenStaker.Contract.Deregister(&_GardenStaker.TransactOpts)
}

// Deregister is a paid mutator transaction binding the contract method 0xaff5edb1.
//
// Solidity: function deregister() returns()
func (_GardenStaker *GardenStakerTransactorSession) Deregister() (*types.Transaction, error) {
	return _GardenStaker.Contract.Deregister(&_GardenStaker.TransactOpts)
}

// Extend is a paid mutator transaction binding the contract method 0x6df3e83f.
//
// Solidity: function extend(bytes32 stakeID, uint256 newLockBlocks) returns()
func (_GardenStaker *GardenStakerTransactor) Extend(opts *bind.TransactOpts, stakeID [32]byte, newLockBlocks *big.Int) (*types.Transaction, error) {
	return _GardenStaker.contract.Transact(opts, "extend", stakeID, newLockBlocks)
}

// Extend is a paid mutator transaction binding the contract method 0x6df3e83f.
//
// Solidity: function extend(bytes32 stakeID, uint256 newLockBlocks) returns()
func (_GardenStaker *GardenStakerSession) Extend(stakeID [32]byte, newLockBlocks *big.Int) (*types.Transaction, error) {
	return _GardenStaker.Contract.Extend(&_GardenStaker.TransactOpts, stakeID, newLockBlocks)
}

// Extend is a paid mutator transaction binding the contract method 0x6df3e83f.
//
// Solidity: function extend(bytes32 stakeID, uint256 newLockBlocks) returns()
func (_GardenStaker *GardenStakerTransactorSession) Extend(stakeID [32]byte, newLockBlocks *big.Int) (*types.Transaction, error) {
	return _GardenStaker.Contract.Extend(&_GardenStaker.TransactOpts, stakeID, newLockBlocks)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GardenStaker *GardenStakerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GardenStaker.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GardenStaker *GardenStakerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GardenStaker.Contract.GrantRole(&_GardenStaker.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GardenStaker *GardenStakerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GardenStaker.Contract.GrantRole(&_GardenStaker.TransactOpts, role, account)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 stakeID) returns()
func (_GardenStaker *GardenStakerTransactor) Refund(opts *bind.TransactOpts, stakeID [32]byte) (*types.Transaction, error) {
	return _GardenStaker.contract.Transact(opts, "refund", stakeID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 stakeID) returns()
func (_GardenStaker *GardenStakerSession) Refund(stakeID [32]byte) (*types.Transaction, error) {
	return _GardenStaker.Contract.Refund(&_GardenStaker.TransactOpts, stakeID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 stakeID) returns()
func (_GardenStaker *GardenStakerTransactorSession) Refund(stakeID [32]byte) (*types.Transaction, error) {
	return _GardenStaker.Contract.Refund(&_GardenStaker.TransactOpts, stakeID)
}

// Refund0 is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address filler_) returns()
func (_GardenStaker *GardenStakerTransactor) Refund0(opts *bind.TransactOpts, filler_ common.Address) (*types.Transaction, error) {
	return _GardenStaker.contract.Transact(opts, "refund0", filler_)
}

// Refund0 is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address filler_) returns()
func (_GardenStaker *GardenStakerSession) Refund0(filler_ common.Address) (*types.Transaction, error) {
	return _GardenStaker.Contract.Refund0(&_GardenStaker.TransactOpts, filler_)
}

// Refund0 is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address filler_) returns()
func (_GardenStaker *GardenStakerTransactorSession) Refund0(filler_ common.Address) (*types.Transaction, error) {
	return _GardenStaker.Contract.Refund0(&_GardenStaker.TransactOpts, filler_)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_GardenStaker *GardenStakerTransactor) Register(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenStaker.contract.Transact(opts, "register")
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_GardenStaker *GardenStakerSession) Register() (*types.Transaction, error) {
	return _GardenStaker.Contract.Register(&_GardenStaker.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_GardenStaker *GardenStakerTransactorSession) Register() (*types.Transaction, error) {
	return _GardenStaker.Contract.Register(&_GardenStaker.TransactOpts)
}

// Renew is a paid mutator transaction binding the contract method 0xf544d82f.
//
// Solidity: function renew(bytes32 stakeID, uint256 newLockBlocks) returns()
func (_GardenStaker *GardenStakerTransactor) Renew(opts *bind.TransactOpts, stakeID [32]byte, newLockBlocks *big.Int) (*types.Transaction, error) {
	return _GardenStaker.contract.Transact(opts, "renew", stakeID, newLockBlocks)
}

// Renew is a paid mutator transaction binding the contract method 0xf544d82f.
//
// Solidity: function renew(bytes32 stakeID, uint256 newLockBlocks) returns()
func (_GardenStaker *GardenStakerSession) Renew(stakeID [32]byte, newLockBlocks *big.Int) (*types.Transaction, error) {
	return _GardenStaker.Contract.Renew(&_GardenStaker.TransactOpts, stakeID, newLockBlocks)
}

// Renew is a paid mutator transaction binding the contract method 0xf544d82f.
//
// Solidity: function renew(bytes32 stakeID, uint256 newLockBlocks) returns()
func (_GardenStaker *GardenStakerTransactorSession) Renew(stakeID [32]byte, newLockBlocks *big.Int) (*types.Transaction, error) {
	return _GardenStaker.Contract.Renew(&_GardenStaker.TransactOpts, stakeID, newLockBlocks)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_GardenStaker *GardenStakerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GardenStaker.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_GardenStaker *GardenStakerSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GardenStaker.Contract.RenounceRole(&_GardenStaker.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_GardenStaker *GardenStakerTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GardenStaker.Contract.RenounceRole(&_GardenStaker.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GardenStaker *GardenStakerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GardenStaker.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GardenStaker *GardenStakerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GardenStaker.Contract.RevokeRole(&_GardenStaker.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GardenStaker *GardenStakerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GardenStaker.Contract.RevokeRole(&_GardenStaker.TransactOpts, role, account)
}

// UpdateFee is a paid mutator transaction binding the contract method 0x2c6cda93.
//
// Solidity: function updateFee(uint16 newFee) returns()
func (_GardenStaker *GardenStakerTransactor) UpdateFee(opts *bind.TransactOpts, newFee uint16) (*types.Transaction, error) {
	return _GardenStaker.contract.Transact(opts, "updateFee", newFee)
}

// UpdateFee is a paid mutator transaction binding the contract method 0x2c6cda93.
//
// Solidity: function updateFee(uint16 newFee) returns()
func (_GardenStaker *GardenStakerSession) UpdateFee(newFee uint16) (*types.Transaction, error) {
	return _GardenStaker.Contract.UpdateFee(&_GardenStaker.TransactOpts, newFee)
}

// UpdateFee is a paid mutator transaction binding the contract method 0x2c6cda93.
//
// Solidity: function updateFee(uint16 newFee) returns()
func (_GardenStaker *GardenStakerTransactorSession) UpdateFee(newFee uint16) (*types.Transaction, error) {
	return _GardenStaker.Contract.UpdateFee(&_GardenStaker.TransactOpts, newFee)
}

// Vote is a paid mutator transaction binding the contract method 0x2a4a1b73.
//
// Solidity: function vote(address filler, uint256 units, uint256 lockBlocks) returns(bytes32 stakeID)
func (_GardenStaker *GardenStakerTransactor) Vote(opts *bind.TransactOpts, filler common.Address, units *big.Int, lockBlocks *big.Int) (*types.Transaction, error) {
	return _GardenStaker.contract.Transact(opts, "vote", filler, units, lockBlocks)
}

// Vote is a paid mutator transaction binding the contract method 0x2a4a1b73.
//
// Solidity: function vote(address filler, uint256 units, uint256 lockBlocks) returns(bytes32 stakeID)
func (_GardenStaker *GardenStakerSession) Vote(filler common.Address, units *big.Int, lockBlocks *big.Int) (*types.Transaction, error) {
	return _GardenStaker.Contract.Vote(&_GardenStaker.TransactOpts, filler, units, lockBlocks)
}

// Vote is a paid mutator transaction binding the contract method 0x2a4a1b73.
//
// Solidity: function vote(address filler, uint256 units, uint256 lockBlocks) returns(bytes32 stakeID)
func (_GardenStaker *GardenStakerTransactorSession) Vote(filler common.Address, units *big.Int, lockBlocks *big.Int) (*types.Transaction, error) {
	return _GardenStaker.Contract.Vote(&_GardenStaker.TransactOpts, filler, units, lockBlocks)
}

// GardenStakerFillerDeregisteredIterator is returned from FilterFillerDeregistered and is used to iterate over the raw logs and unpacked data for FillerDeregistered events raised by the GardenStaker contract.
type GardenStakerFillerDeregisteredIterator struct {
	Event *GardenStakerFillerDeregistered // Event containing the contract specifics and raw log

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
func (it *GardenStakerFillerDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenStakerFillerDeregistered)
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
		it.Event = new(GardenStakerFillerDeregistered)
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
func (it *GardenStakerFillerDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenStakerFillerDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenStakerFillerDeregistered represents a FillerDeregistered event raised by the GardenStaker contract.
type GardenStakerFillerDeregistered struct {
	Filler         common.Address
	DeregisteredAt *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterFillerDeregistered is a free log retrieval operation binding the contract event 0x83f9b3ce2be58c675065c1675b20665690bc38864c4a0335ee681a57f9f88227.
//
// Solidity: event FillerDeregistered(address indexed filler, uint256 deregisteredAt)
func (_GardenStaker *GardenStakerFilterer) FilterFillerDeregistered(opts *bind.FilterOpts, filler []common.Address) (*GardenStakerFillerDeregisteredIterator, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _GardenStaker.contract.FilterLogs(opts, "FillerDeregistered", fillerRule)
	if err != nil {
		return nil, err
	}
	return &GardenStakerFillerDeregisteredIterator{contract: _GardenStaker.contract, event: "FillerDeregistered", logs: logs, sub: sub}, nil
}

// WatchFillerDeregistered is a free log subscription operation binding the contract event 0x83f9b3ce2be58c675065c1675b20665690bc38864c4a0335ee681a57f9f88227.
//
// Solidity: event FillerDeregistered(address indexed filler, uint256 deregisteredAt)
func (_GardenStaker *GardenStakerFilterer) WatchFillerDeregistered(opts *bind.WatchOpts, sink chan<- *GardenStakerFillerDeregistered, filler []common.Address) (event.Subscription, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _GardenStaker.contract.WatchLogs(opts, "FillerDeregistered", fillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenStakerFillerDeregistered)
				if err := _GardenStaker.contract.UnpackLog(event, "FillerDeregistered", log); err != nil {
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

// ParseFillerDeregistered is a log parse operation binding the contract event 0x83f9b3ce2be58c675065c1675b20665690bc38864c4a0335ee681a57f9f88227.
//
// Solidity: event FillerDeregistered(address indexed filler, uint256 deregisteredAt)
func (_GardenStaker *GardenStakerFilterer) ParseFillerDeregistered(log types.Log) (*GardenStakerFillerDeregistered, error) {
	event := new(GardenStakerFillerDeregistered)
	if err := _GardenStaker.contract.UnpackLog(event, "FillerDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenStakerFillerFeeUpdatedIterator is returned from FilterFillerFeeUpdated and is used to iterate over the raw logs and unpacked data for FillerFeeUpdated events raised by the GardenStaker contract.
type GardenStakerFillerFeeUpdatedIterator struct {
	Event *GardenStakerFillerFeeUpdated // Event containing the contract specifics and raw log

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
func (it *GardenStakerFillerFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenStakerFillerFeeUpdated)
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
		it.Event = new(GardenStakerFillerFeeUpdated)
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
func (it *GardenStakerFillerFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenStakerFillerFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenStakerFillerFeeUpdated represents a FillerFeeUpdated event raised by the GardenStaker contract.
type GardenStakerFillerFeeUpdated struct {
	Filler common.Address
	Fee    *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFillerFeeUpdated is a free log retrieval operation binding the contract event 0x7318fde364eaeddb75aa88119c39d65d30c300d1e4279711a132b11922b9aa64.
//
// Solidity: event FillerFeeUpdated(address indexed filler, uint256 fee)
func (_GardenStaker *GardenStakerFilterer) FilterFillerFeeUpdated(opts *bind.FilterOpts, filler []common.Address) (*GardenStakerFillerFeeUpdatedIterator, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _GardenStaker.contract.FilterLogs(opts, "FillerFeeUpdated", fillerRule)
	if err != nil {
		return nil, err
	}
	return &GardenStakerFillerFeeUpdatedIterator{contract: _GardenStaker.contract, event: "FillerFeeUpdated", logs: logs, sub: sub}, nil
}

// WatchFillerFeeUpdated is a free log subscription operation binding the contract event 0x7318fde364eaeddb75aa88119c39d65d30c300d1e4279711a132b11922b9aa64.
//
// Solidity: event FillerFeeUpdated(address indexed filler, uint256 fee)
func (_GardenStaker *GardenStakerFilterer) WatchFillerFeeUpdated(opts *bind.WatchOpts, sink chan<- *GardenStakerFillerFeeUpdated, filler []common.Address) (event.Subscription, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _GardenStaker.contract.WatchLogs(opts, "FillerFeeUpdated", fillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenStakerFillerFeeUpdated)
				if err := _GardenStaker.contract.UnpackLog(event, "FillerFeeUpdated", log); err != nil {
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

// ParseFillerFeeUpdated is a log parse operation binding the contract event 0x7318fde364eaeddb75aa88119c39d65d30c300d1e4279711a132b11922b9aa64.
//
// Solidity: event FillerFeeUpdated(address indexed filler, uint256 fee)
func (_GardenStaker *GardenStakerFilterer) ParseFillerFeeUpdated(log types.Log) (*GardenStakerFillerFeeUpdated, error) {
	event := new(GardenStakerFillerFeeUpdated)
	if err := _GardenStaker.contract.UnpackLog(event, "FillerFeeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenStakerFillerRefundedIterator is returned from FilterFillerRefunded and is used to iterate over the raw logs and unpacked data for FillerRefunded events raised by the GardenStaker contract.
type GardenStakerFillerRefundedIterator struct {
	Event *GardenStakerFillerRefunded // Event containing the contract specifics and raw log

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
func (it *GardenStakerFillerRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenStakerFillerRefunded)
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
		it.Event = new(GardenStakerFillerRefunded)
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
func (it *GardenStakerFillerRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenStakerFillerRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenStakerFillerRefunded represents a FillerRefunded event raised by the GardenStaker contract.
type GardenStakerFillerRefunded struct {
	Filler common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFillerRefunded is a free log retrieval operation binding the contract event 0xbe2dbe5153e94feab433d0b81928c0c44ffadf824ea8b4692f7523e0f4c16cd6.
//
// Solidity: event FillerRefunded(address indexed filler)
func (_GardenStaker *GardenStakerFilterer) FilterFillerRefunded(opts *bind.FilterOpts, filler []common.Address) (*GardenStakerFillerRefundedIterator, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _GardenStaker.contract.FilterLogs(opts, "FillerRefunded", fillerRule)
	if err != nil {
		return nil, err
	}
	return &GardenStakerFillerRefundedIterator{contract: _GardenStaker.contract, event: "FillerRefunded", logs: logs, sub: sub}, nil
}

// WatchFillerRefunded is a free log subscription operation binding the contract event 0xbe2dbe5153e94feab433d0b81928c0c44ffadf824ea8b4692f7523e0f4c16cd6.
//
// Solidity: event FillerRefunded(address indexed filler)
func (_GardenStaker *GardenStakerFilterer) WatchFillerRefunded(opts *bind.WatchOpts, sink chan<- *GardenStakerFillerRefunded, filler []common.Address) (event.Subscription, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _GardenStaker.contract.WatchLogs(opts, "FillerRefunded", fillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenStakerFillerRefunded)
				if err := _GardenStaker.contract.UnpackLog(event, "FillerRefunded", log); err != nil {
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

// ParseFillerRefunded is a log parse operation binding the contract event 0xbe2dbe5153e94feab433d0b81928c0c44ffadf824ea8b4692f7523e0f4c16cd6.
//
// Solidity: event FillerRefunded(address indexed filler)
func (_GardenStaker *GardenStakerFilterer) ParseFillerRefunded(log types.Log) (*GardenStakerFillerRefunded, error) {
	event := new(GardenStakerFillerRefunded)
	if err := _GardenStaker.contract.UnpackLog(event, "FillerRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenStakerFillerRegisteredIterator is returned from FilterFillerRegistered and is used to iterate over the raw logs and unpacked data for FillerRegistered events raised by the GardenStaker contract.
type GardenStakerFillerRegisteredIterator struct {
	Event *GardenStakerFillerRegistered // Event containing the contract specifics and raw log

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
func (it *GardenStakerFillerRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenStakerFillerRegistered)
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
		it.Event = new(GardenStakerFillerRegistered)
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
func (it *GardenStakerFillerRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenStakerFillerRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenStakerFillerRegistered represents a FillerRegistered event raised by the GardenStaker contract.
type GardenStakerFillerRegistered struct {
	Filler common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFillerRegistered is a free log retrieval operation binding the contract event 0x3d5e4b9cd5a0b3c03402df7b2eb1b594bed42e75163fc8d3f08ff1d1751bfdfd.
//
// Solidity: event FillerRegistered(address indexed filler)
func (_GardenStaker *GardenStakerFilterer) FilterFillerRegistered(opts *bind.FilterOpts, filler []common.Address) (*GardenStakerFillerRegisteredIterator, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _GardenStaker.contract.FilterLogs(opts, "FillerRegistered", fillerRule)
	if err != nil {
		return nil, err
	}
	return &GardenStakerFillerRegisteredIterator{contract: _GardenStaker.contract, event: "FillerRegistered", logs: logs, sub: sub}, nil
}

// WatchFillerRegistered is a free log subscription operation binding the contract event 0x3d5e4b9cd5a0b3c03402df7b2eb1b594bed42e75163fc8d3f08ff1d1751bfdfd.
//
// Solidity: event FillerRegistered(address indexed filler)
func (_GardenStaker *GardenStakerFilterer) WatchFillerRegistered(opts *bind.WatchOpts, sink chan<- *GardenStakerFillerRegistered, filler []common.Address) (event.Subscription, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _GardenStaker.contract.WatchLogs(opts, "FillerRegistered", fillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenStakerFillerRegistered)
				if err := _GardenStaker.contract.UnpackLog(event, "FillerRegistered", log); err != nil {
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

// ParseFillerRegistered is a log parse operation binding the contract event 0x3d5e4b9cd5a0b3c03402df7b2eb1b594bed42e75163fc8d3f08ff1d1751bfdfd.
//
// Solidity: event FillerRegistered(address indexed filler)
func (_GardenStaker *GardenStakerFilterer) ParseFillerRegistered(log types.Log) (*GardenStakerFillerRegistered, error) {
	event := new(GardenStakerFillerRegistered)
	if err := _GardenStaker.contract.UnpackLog(event, "FillerRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenStakerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the GardenStaker contract.
type GardenStakerRoleAdminChangedIterator struct {
	Event *GardenStakerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *GardenStakerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenStakerRoleAdminChanged)
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
		it.Event = new(GardenStakerRoleAdminChanged)
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
func (it *GardenStakerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenStakerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenStakerRoleAdminChanged represents a RoleAdminChanged event raised by the GardenStaker contract.
type GardenStakerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GardenStaker *GardenStakerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*GardenStakerRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _GardenStaker.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &GardenStakerRoleAdminChangedIterator{contract: _GardenStaker.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GardenStaker *GardenStakerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *GardenStakerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _GardenStaker.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenStakerRoleAdminChanged)
				if err := _GardenStaker.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GardenStaker *GardenStakerFilterer) ParseRoleAdminChanged(log types.Log) (*GardenStakerRoleAdminChanged, error) {
	event := new(GardenStakerRoleAdminChanged)
	if err := _GardenStaker.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenStakerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the GardenStaker contract.
type GardenStakerRoleGrantedIterator struct {
	Event *GardenStakerRoleGranted // Event containing the contract specifics and raw log

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
func (it *GardenStakerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenStakerRoleGranted)
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
		it.Event = new(GardenStakerRoleGranted)
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
func (it *GardenStakerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenStakerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenStakerRoleGranted represents a RoleGranted event raised by the GardenStaker contract.
type GardenStakerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GardenStaker *GardenStakerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*GardenStakerRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GardenStaker.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &GardenStakerRoleGrantedIterator{contract: _GardenStaker.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GardenStaker *GardenStakerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *GardenStakerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GardenStaker.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenStakerRoleGranted)
				if err := _GardenStaker.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GardenStaker *GardenStakerFilterer) ParseRoleGranted(log types.Log) (*GardenStakerRoleGranted, error) {
	event := new(GardenStakerRoleGranted)
	if err := _GardenStaker.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenStakerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the GardenStaker contract.
type GardenStakerRoleRevokedIterator struct {
	Event *GardenStakerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *GardenStakerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenStakerRoleRevoked)
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
		it.Event = new(GardenStakerRoleRevoked)
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
func (it *GardenStakerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenStakerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenStakerRoleRevoked represents a RoleRevoked event raised by the GardenStaker contract.
type GardenStakerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GardenStaker *GardenStakerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*GardenStakerRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GardenStaker.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &GardenStakerRoleRevokedIterator{contract: _GardenStaker.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GardenStaker *GardenStakerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *GardenStakerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GardenStaker.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenStakerRoleRevoked)
				if err := _GardenStaker.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GardenStaker *GardenStakerFilterer) ParseRoleRevoked(log types.Log) (*GardenStakerRoleRevoked, error) {
	event := new(GardenStakerRoleRevoked)
	if err := _GardenStaker.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenStakerStakeExtendedIterator is returned from FilterStakeExtended and is used to iterate over the raw logs and unpacked data for StakeExtended events raised by the GardenStaker contract.
type GardenStakerStakeExtendedIterator struct {
	Event *GardenStakerStakeExtended // Event containing the contract specifics and raw log

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
func (it *GardenStakerStakeExtendedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenStakerStakeExtended)
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
		it.Event = new(GardenStakerStakeExtended)
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
func (it *GardenStakerStakeExtendedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenStakerStakeExtendedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenStakerStakeExtended represents a StakeExtended event raised by the GardenStaker contract.
type GardenStakerStakeExtended struct {
	StakeID       [32]byte
	NewLockBlocks *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterStakeExtended is a free log retrieval operation binding the contract event 0x521ebd8078bb6d9f423366c9ae8a0d4dcec7887f1db2d3c0cd0923fe00ed2ed7.
//
// Solidity: event StakeExtended(bytes32 indexed stakeID, uint256 newLockBlocks)
func (_GardenStaker *GardenStakerFilterer) FilterStakeExtended(opts *bind.FilterOpts, stakeID [][32]byte) (*GardenStakerStakeExtendedIterator, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}

	logs, sub, err := _GardenStaker.contract.FilterLogs(opts, "StakeExtended", stakeIDRule)
	if err != nil {
		return nil, err
	}
	return &GardenStakerStakeExtendedIterator{contract: _GardenStaker.contract, event: "StakeExtended", logs: logs, sub: sub}, nil
}

// WatchStakeExtended is a free log subscription operation binding the contract event 0x521ebd8078bb6d9f423366c9ae8a0d4dcec7887f1db2d3c0cd0923fe00ed2ed7.
//
// Solidity: event StakeExtended(bytes32 indexed stakeID, uint256 newLockBlocks)
func (_GardenStaker *GardenStakerFilterer) WatchStakeExtended(opts *bind.WatchOpts, sink chan<- *GardenStakerStakeExtended, stakeID [][32]byte) (event.Subscription, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}

	logs, sub, err := _GardenStaker.contract.WatchLogs(opts, "StakeExtended", stakeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenStakerStakeExtended)
				if err := _GardenStaker.contract.UnpackLog(event, "StakeExtended", log); err != nil {
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

// ParseStakeExtended is a log parse operation binding the contract event 0x521ebd8078bb6d9f423366c9ae8a0d4dcec7887f1db2d3c0cd0923fe00ed2ed7.
//
// Solidity: event StakeExtended(bytes32 indexed stakeID, uint256 newLockBlocks)
func (_GardenStaker *GardenStakerFilterer) ParseStakeExtended(log types.Log) (*GardenStakerStakeExtended, error) {
	event := new(GardenStakerStakeExtended)
	if err := _GardenStaker.contract.UnpackLog(event, "StakeExtended", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenStakerStakeRefundedIterator is returned from FilterStakeRefunded and is used to iterate over the raw logs and unpacked data for StakeRefunded events raised by the GardenStaker contract.
type GardenStakerStakeRefundedIterator struct {
	Event *GardenStakerStakeRefunded // Event containing the contract specifics and raw log

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
func (it *GardenStakerStakeRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenStakerStakeRefunded)
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
		it.Event = new(GardenStakerStakeRefunded)
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
func (it *GardenStakerStakeRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenStakerStakeRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenStakerStakeRefunded represents a StakeRefunded event raised by the GardenStaker contract.
type GardenStakerStakeRefunded struct {
	StakeID [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterStakeRefunded is a free log retrieval operation binding the contract event 0xba33cc6a4502f7f80bfe643ffa925ba7ab5e0af61c061a9aceabef0349c9dce4.
//
// Solidity: event StakeRefunded(bytes32 indexed stakeID)
func (_GardenStaker *GardenStakerFilterer) FilterStakeRefunded(opts *bind.FilterOpts, stakeID [][32]byte) (*GardenStakerStakeRefundedIterator, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}

	logs, sub, err := _GardenStaker.contract.FilterLogs(opts, "StakeRefunded", stakeIDRule)
	if err != nil {
		return nil, err
	}
	return &GardenStakerStakeRefundedIterator{contract: _GardenStaker.contract, event: "StakeRefunded", logs: logs, sub: sub}, nil
}

// WatchStakeRefunded is a free log subscription operation binding the contract event 0xba33cc6a4502f7f80bfe643ffa925ba7ab5e0af61c061a9aceabef0349c9dce4.
//
// Solidity: event StakeRefunded(bytes32 indexed stakeID)
func (_GardenStaker *GardenStakerFilterer) WatchStakeRefunded(opts *bind.WatchOpts, sink chan<- *GardenStakerStakeRefunded, stakeID [][32]byte) (event.Subscription, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}

	logs, sub, err := _GardenStaker.contract.WatchLogs(opts, "StakeRefunded", stakeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenStakerStakeRefunded)
				if err := _GardenStaker.contract.UnpackLog(event, "StakeRefunded", log); err != nil {
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

// ParseStakeRefunded is a log parse operation binding the contract event 0xba33cc6a4502f7f80bfe643ffa925ba7ab5e0af61c061a9aceabef0349c9dce4.
//
// Solidity: event StakeRefunded(bytes32 indexed stakeID)
func (_GardenStaker *GardenStakerFilterer) ParseStakeRefunded(log types.Log) (*GardenStakerStakeRefunded, error) {
	event := new(GardenStakerStakeRefunded)
	if err := _GardenStaker.contract.UnpackLog(event, "StakeRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenStakerStakeRenewedIterator is returned from FilterStakeRenewed and is used to iterate over the raw logs and unpacked data for StakeRenewed events raised by the GardenStaker contract.
type GardenStakerStakeRenewedIterator struct {
	Event *GardenStakerStakeRenewed // Event containing the contract specifics and raw log

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
func (it *GardenStakerStakeRenewedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenStakerStakeRenewed)
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
		it.Event = new(GardenStakerStakeRenewed)
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
func (it *GardenStakerStakeRenewedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenStakerStakeRenewedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenStakerStakeRenewed represents a StakeRenewed event raised by the GardenStaker contract.
type GardenStakerStakeRenewed struct {
	StakeID       [32]byte
	NewLockBlocks *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterStakeRenewed is a free log retrieval operation binding the contract event 0xe2c4eee10fe6a93f81bf01601e62cba732126941bccca75f3460e6f2f845513d.
//
// Solidity: event StakeRenewed(bytes32 indexed stakeID, uint256 newLockBlocks)
func (_GardenStaker *GardenStakerFilterer) FilterStakeRenewed(opts *bind.FilterOpts, stakeID [][32]byte) (*GardenStakerStakeRenewedIterator, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}

	logs, sub, err := _GardenStaker.contract.FilterLogs(opts, "StakeRenewed", stakeIDRule)
	if err != nil {
		return nil, err
	}
	return &GardenStakerStakeRenewedIterator{contract: _GardenStaker.contract, event: "StakeRenewed", logs: logs, sub: sub}, nil
}

// WatchStakeRenewed is a free log subscription operation binding the contract event 0xe2c4eee10fe6a93f81bf01601e62cba732126941bccca75f3460e6f2f845513d.
//
// Solidity: event StakeRenewed(bytes32 indexed stakeID, uint256 newLockBlocks)
func (_GardenStaker *GardenStakerFilterer) WatchStakeRenewed(opts *bind.WatchOpts, sink chan<- *GardenStakerStakeRenewed, stakeID [][32]byte) (event.Subscription, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}

	logs, sub, err := _GardenStaker.contract.WatchLogs(opts, "StakeRenewed", stakeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenStakerStakeRenewed)
				if err := _GardenStaker.contract.UnpackLog(event, "StakeRenewed", log); err != nil {
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

// ParseStakeRenewed is a log parse operation binding the contract event 0xe2c4eee10fe6a93f81bf01601e62cba732126941bccca75f3460e6f2f845513d.
//
// Solidity: event StakeRenewed(bytes32 indexed stakeID, uint256 newLockBlocks)
func (_GardenStaker *GardenStakerFilterer) ParseStakeRenewed(log types.Log) (*GardenStakerStakeRenewed, error) {
	event := new(GardenStakerStakeRenewed)
	if err := _GardenStaker.contract.UnpackLog(event, "StakeRenewed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenStakerStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the GardenStaker contract.
type GardenStakerStakedIterator struct {
	Event *GardenStakerStaked // Event containing the contract specifics and raw log

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
func (it *GardenStakerStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenStakerStaked)
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
		it.Event = new(GardenStakerStaked)
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
func (it *GardenStakerStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenStakerStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenStakerStaked represents a Staked event raised by the GardenStaker contract.
type GardenStakerStaked struct {
	StakeID [32]byte
	Owner   common.Address
	Stake   *big.Int
	Expiry  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0xa1fdccfe567643a44425efdd141171e8d992854a81e5c819c1432b0de47c9a11.
//
// Solidity: event Staked(bytes32 indexed stakeID, address indexed owner, uint256 stake, uint256 expiry)
func (_GardenStaker *GardenStakerFilterer) FilterStaked(opts *bind.FilterOpts, stakeID [][32]byte, owner []common.Address) (*GardenStakerStakedIterator, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _GardenStaker.contract.FilterLogs(opts, "Staked", stakeIDRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &GardenStakerStakedIterator{contract: _GardenStaker.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0xa1fdccfe567643a44425efdd141171e8d992854a81e5c819c1432b0de47c9a11.
//
// Solidity: event Staked(bytes32 indexed stakeID, address indexed owner, uint256 stake, uint256 expiry)
func (_GardenStaker *GardenStakerFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *GardenStakerStaked, stakeID [][32]byte, owner []common.Address) (event.Subscription, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _GardenStaker.contract.WatchLogs(opts, "Staked", stakeIDRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenStakerStaked)
				if err := _GardenStaker.contract.UnpackLog(event, "Staked", log); err != nil {
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

// ParseStaked is a log parse operation binding the contract event 0xa1fdccfe567643a44425efdd141171e8d992854a81e5c819c1432b0de47c9a11.
//
// Solidity: event Staked(bytes32 indexed stakeID, address indexed owner, uint256 stake, uint256 expiry)
func (_GardenStaker *GardenStakerFilterer) ParseStaked(log types.Log) (*GardenStakerStaked, error) {
	event := new(GardenStakerStaked)
	if err := _GardenStaker.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenStakerVotedIterator is returned from FilterVoted and is used to iterate over the raw logs and unpacked data for Voted events raised by the GardenStaker contract.
type GardenStakerVotedIterator struct {
	Event *GardenStakerVoted // Event containing the contract specifics and raw log

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
func (it *GardenStakerVotedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenStakerVoted)
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
		it.Event = new(GardenStakerVoted)
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
func (it *GardenStakerVotedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenStakerVotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenStakerVoted represents a Voted event raised by the GardenStaker contract.
type GardenStakerVoted struct {
	StakeID [32]byte
	Filler  common.Address
	Votes   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterVoted is a free log retrieval operation binding the contract event 0xe4abc5380fa6939d1dc23b5e90b3a8a0e328f0f1a82a5f42bfb795bf9c717505.
//
// Solidity: event Voted(bytes32 indexed stakeID, address indexed filler, uint256 votes)
func (_GardenStaker *GardenStakerFilterer) FilterVoted(opts *bind.FilterOpts, stakeID [][32]byte, filler []common.Address) (*GardenStakerVotedIterator, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}
	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _GardenStaker.contract.FilterLogs(opts, "Voted", stakeIDRule, fillerRule)
	if err != nil {
		return nil, err
	}
	return &GardenStakerVotedIterator{contract: _GardenStaker.contract, event: "Voted", logs: logs, sub: sub}, nil
}

// WatchVoted is a free log subscription operation binding the contract event 0xe4abc5380fa6939d1dc23b5e90b3a8a0e328f0f1a82a5f42bfb795bf9c717505.
//
// Solidity: event Voted(bytes32 indexed stakeID, address indexed filler, uint256 votes)
func (_GardenStaker *GardenStakerFilterer) WatchVoted(opts *bind.WatchOpts, sink chan<- *GardenStakerVoted, stakeID [][32]byte, filler []common.Address) (event.Subscription, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}
	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _GardenStaker.contract.WatchLogs(opts, "Voted", stakeIDRule, fillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenStakerVoted)
				if err := _GardenStaker.contract.UnpackLog(event, "Voted", log); err != nil {
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

// ParseVoted is a log parse operation binding the contract event 0xe4abc5380fa6939d1dc23b5e90b3a8a0e328f0f1a82a5f42bfb795bf9c717505.
//
// Solidity: event Voted(bytes32 indexed stakeID, address indexed filler, uint256 votes)
func (_GardenStaker *GardenStakerFilterer) ParseVoted(log types.Log) (*GardenStakerVoted, error) {
	event := new(GardenStakerVoted)
	if err := _GardenStaker.contract.UnpackLog(event, "Voted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenStakerVotesChangedIterator is returned from FilterVotesChanged and is used to iterate over the raw logs and unpacked data for VotesChanged events raised by the GardenStaker contract.
type GardenStakerVotesChangedIterator struct {
	Event *GardenStakerVotesChanged // Event containing the contract specifics and raw log

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
func (it *GardenStakerVotesChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenStakerVotesChanged)
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
		it.Event = new(GardenStakerVotesChanged)
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
func (it *GardenStakerVotesChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenStakerVotesChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenStakerVotesChanged represents a VotesChanged event raised by the GardenStaker contract.
type GardenStakerVotesChanged struct {
	StakeID   [32]byte
	OldFiller common.Address
	NewFiller common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterVotesChanged is a free log retrieval operation binding the contract event 0xf3a02e43bc9e793a4fe98273c318b73e54fbbd0c599adb2b620ade22b2235f6e.
//
// Solidity: event VotesChanged(bytes32 indexed stakeID, address indexed oldFiller, address indexed newFiller)
func (_GardenStaker *GardenStakerFilterer) FilterVotesChanged(opts *bind.FilterOpts, stakeID [][32]byte, oldFiller []common.Address, newFiller []common.Address) (*GardenStakerVotesChangedIterator, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}
	var oldFillerRule []interface{}
	for _, oldFillerItem := range oldFiller {
		oldFillerRule = append(oldFillerRule, oldFillerItem)
	}
	var newFillerRule []interface{}
	for _, newFillerItem := range newFiller {
		newFillerRule = append(newFillerRule, newFillerItem)
	}

	logs, sub, err := _GardenStaker.contract.FilterLogs(opts, "VotesChanged", stakeIDRule, oldFillerRule, newFillerRule)
	if err != nil {
		return nil, err
	}
	return &GardenStakerVotesChangedIterator{contract: _GardenStaker.contract, event: "VotesChanged", logs: logs, sub: sub}, nil
}

// WatchVotesChanged is a free log subscription operation binding the contract event 0xf3a02e43bc9e793a4fe98273c318b73e54fbbd0c599adb2b620ade22b2235f6e.
//
// Solidity: event VotesChanged(bytes32 indexed stakeID, address indexed oldFiller, address indexed newFiller)
func (_GardenStaker *GardenStakerFilterer) WatchVotesChanged(opts *bind.WatchOpts, sink chan<- *GardenStakerVotesChanged, stakeID [][32]byte, oldFiller []common.Address, newFiller []common.Address) (event.Subscription, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}
	var oldFillerRule []interface{}
	for _, oldFillerItem := range oldFiller {
		oldFillerRule = append(oldFillerRule, oldFillerItem)
	}
	var newFillerRule []interface{}
	for _, newFillerItem := range newFiller {
		newFillerRule = append(newFillerRule, newFillerItem)
	}

	logs, sub, err := _GardenStaker.contract.WatchLogs(opts, "VotesChanged", stakeIDRule, oldFillerRule, newFillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenStakerVotesChanged)
				if err := _GardenStaker.contract.UnpackLog(event, "VotesChanged", log); err != nil {
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

// ParseVotesChanged is a log parse operation binding the contract event 0xf3a02e43bc9e793a4fe98273c318b73e54fbbd0c599adb2b620ade22b2235f6e.
//
// Solidity: event VotesChanged(bytes32 indexed stakeID, address indexed oldFiller, address indexed newFiller)
func (_GardenStaker *GardenStakerFilterer) ParseVotesChanged(log types.Log) (*GardenStakerVotesChanged, error) {
	event := new(GardenStakerVotesChanged)
	if err := _GardenStaker.contract.UnpackLog(event, "VotesChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
