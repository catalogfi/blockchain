// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package delegatemanager

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

// DelegateManagerMetaData contains all meta data concerning the DelegateManager contract.
var DelegateManagerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newLockBlocks\",\"type\":\"uint256\"}],\"name\":\"StakeExtended\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"}],\"name\":\"StakeRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newLockBlocks\",\"type\":\"uint256\"}],\"name\":\"StakeRenewed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"}],\"name\":\"Voted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldFiller\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newFiller\",\"type\":\"address\"}],\"name\":\"VotesChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DELEGATE_STAKE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FILLER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FILLER_COOL_DOWN\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FILLER_STAKE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SEED\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"newFiller\",\"type\":\"address\"}],\"name\":\"changeVote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"delegateNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"newLockBlocks\",\"type\":\"uint256\"}],\"name\":\"extend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"}],\"name\":\"getFiller\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"feeInBips\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deregisteredAt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"delegateStakeIDs\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"}],\"name\":\"getVotes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"newLockBlocks\",\"type\":\"uint256\"}],\"name\":\"renew\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"stakes\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"units\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"units\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockBlocks\",\"type\":\"uint256\"}],\"name\":\"vote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"stakeID\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DelegateManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use DelegateManagerMetaData.ABI instead.
var DelegateManagerABI = DelegateManagerMetaData.ABI

// DelegateManager is an auto generated Go binding around an Ethereum contract.
type DelegateManager struct {
	DelegateManagerCaller     // Read-only binding to the contract
	DelegateManagerTransactor // Write-only binding to the contract
	DelegateManagerFilterer   // Log filterer for contract events
}

// DelegateManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegateManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegateManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegateManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegateManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegateManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegateManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegateManagerSession struct {
	Contract     *DelegateManager  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DelegateManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegateManagerCallerSession struct {
	Contract *DelegateManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// DelegateManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegateManagerTransactorSession struct {
	Contract     *DelegateManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// DelegateManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegateManagerRaw struct {
	Contract *DelegateManager // Generic contract binding to access the raw methods on
}

// DelegateManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegateManagerCallerRaw struct {
	Contract *DelegateManagerCaller // Generic read-only contract binding to access the raw methods on
}

// DelegateManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegateManagerTransactorRaw struct {
	Contract *DelegateManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegateManager creates a new instance of DelegateManager, bound to a specific deployed contract.
func NewDelegateManager(address common.Address, backend bind.ContractBackend) (*DelegateManager, error) {
	contract, err := bindDelegateManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DelegateManager{DelegateManagerCaller: DelegateManagerCaller{contract: contract}, DelegateManagerTransactor: DelegateManagerTransactor{contract: contract}, DelegateManagerFilterer: DelegateManagerFilterer{contract: contract}}, nil
}

// NewDelegateManagerCaller creates a new read-only instance of DelegateManager, bound to a specific deployed contract.
func NewDelegateManagerCaller(address common.Address, caller bind.ContractCaller) (*DelegateManagerCaller, error) {
	contract, err := bindDelegateManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegateManagerCaller{contract: contract}, nil
}

// NewDelegateManagerTransactor creates a new write-only instance of DelegateManager, bound to a specific deployed contract.
func NewDelegateManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegateManagerTransactor, error) {
	contract, err := bindDelegateManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegateManagerTransactor{contract: contract}, nil
}

// NewDelegateManagerFilterer creates a new log filterer instance of DelegateManager, bound to a specific deployed contract.
func NewDelegateManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegateManagerFilterer, error) {
	contract, err := bindDelegateManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegateManagerFilterer{contract: contract}, nil
}

// bindDelegateManager binds a generic wrapper to an already deployed contract.
func bindDelegateManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DelegateManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelegateManager *DelegateManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DelegateManager.Contract.DelegateManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelegateManager *DelegateManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegateManager.Contract.DelegateManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelegateManager *DelegateManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelegateManager.Contract.DelegateManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelegateManager *DelegateManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DelegateManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelegateManager *DelegateManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegateManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelegateManager *DelegateManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelegateManager.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DelegateManager *DelegateManagerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DelegateManager.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DelegateManager *DelegateManagerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _DelegateManager.Contract.DEFAULTADMINROLE(&_DelegateManager.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DelegateManager *DelegateManagerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _DelegateManager.Contract.DEFAULTADMINROLE(&_DelegateManager.CallOpts)
}

// DELEGATESTAKE is a free data retrieval call binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() view returns(uint256)
func (_DelegateManager *DelegateManagerCaller) DELEGATESTAKE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DelegateManager.contract.Call(opts, &out, "DELEGATE_STAKE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DELEGATESTAKE is a free data retrieval call binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() view returns(uint256)
func (_DelegateManager *DelegateManagerSession) DELEGATESTAKE() (*big.Int, error) {
	return _DelegateManager.Contract.DELEGATESTAKE(&_DelegateManager.CallOpts)
}

// DELEGATESTAKE is a free data retrieval call binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() view returns(uint256)
func (_DelegateManager *DelegateManagerCallerSession) DELEGATESTAKE() (*big.Int, error) {
	return _DelegateManager.Contract.DELEGATESTAKE(&_DelegateManager.CallOpts)
}

// FILLER is a free data retrieval call binding the contract method 0x64645bae.
//
// Solidity: function FILLER() view returns(bytes32)
func (_DelegateManager *DelegateManagerCaller) FILLER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DelegateManager.contract.Call(opts, &out, "FILLER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FILLER is a free data retrieval call binding the contract method 0x64645bae.
//
// Solidity: function FILLER() view returns(bytes32)
func (_DelegateManager *DelegateManagerSession) FILLER() ([32]byte, error) {
	return _DelegateManager.Contract.FILLER(&_DelegateManager.CallOpts)
}

// FILLER is a free data retrieval call binding the contract method 0x64645bae.
//
// Solidity: function FILLER() view returns(bytes32)
func (_DelegateManager *DelegateManagerCallerSession) FILLER() ([32]byte, error) {
	return _DelegateManager.Contract.FILLER(&_DelegateManager.CallOpts)
}

// FILLERCOOLDOWN is a free data retrieval call binding the contract method 0x015a44e8.
//
// Solidity: function FILLER_COOL_DOWN() view returns(uint256)
func (_DelegateManager *DelegateManagerCaller) FILLERCOOLDOWN(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DelegateManager.contract.Call(opts, &out, "FILLER_COOL_DOWN")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FILLERCOOLDOWN is a free data retrieval call binding the contract method 0x015a44e8.
//
// Solidity: function FILLER_COOL_DOWN() view returns(uint256)
func (_DelegateManager *DelegateManagerSession) FILLERCOOLDOWN() (*big.Int, error) {
	return _DelegateManager.Contract.FILLERCOOLDOWN(&_DelegateManager.CallOpts)
}

// FILLERCOOLDOWN is a free data retrieval call binding the contract method 0x015a44e8.
//
// Solidity: function FILLER_COOL_DOWN() view returns(uint256)
func (_DelegateManager *DelegateManagerCallerSession) FILLERCOOLDOWN() (*big.Int, error) {
	return _DelegateManager.Contract.FILLERCOOLDOWN(&_DelegateManager.CallOpts)
}

// FILLERSTAKE is a free data retrieval call binding the contract method 0xecfe923e.
//
// Solidity: function FILLER_STAKE() view returns(uint256)
func (_DelegateManager *DelegateManagerCaller) FILLERSTAKE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DelegateManager.contract.Call(opts, &out, "FILLER_STAKE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FILLERSTAKE is a free data retrieval call binding the contract method 0xecfe923e.
//
// Solidity: function FILLER_STAKE() view returns(uint256)
func (_DelegateManager *DelegateManagerSession) FILLERSTAKE() (*big.Int, error) {
	return _DelegateManager.Contract.FILLERSTAKE(&_DelegateManager.CallOpts)
}

// FILLERSTAKE is a free data retrieval call binding the contract method 0xecfe923e.
//
// Solidity: function FILLER_STAKE() view returns(uint256)
func (_DelegateManager *DelegateManagerCallerSession) FILLERSTAKE() (*big.Int, error) {
	return _DelegateManager.Contract.FILLERSTAKE(&_DelegateManager.CallOpts)
}

// SEED is a free data retrieval call binding the contract method 0x0edc4737.
//
// Solidity: function SEED() view returns(address)
func (_DelegateManager *DelegateManagerCaller) SEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DelegateManager.contract.Call(opts, &out, "SEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SEED is a free data retrieval call binding the contract method 0x0edc4737.
//
// Solidity: function SEED() view returns(address)
func (_DelegateManager *DelegateManagerSession) SEED() (common.Address, error) {
	return _DelegateManager.Contract.SEED(&_DelegateManager.CallOpts)
}

// SEED is a free data retrieval call binding the contract method 0x0edc4737.
//
// Solidity: function SEED() view returns(address)
func (_DelegateManager *DelegateManagerCallerSession) SEED() (common.Address, error) {
	return _DelegateManager.Contract.SEED(&_DelegateManager.CallOpts)
}

// DelegateNonce is a free data retrieval call binding the contract method 0x92e1f102.
//
// Solidity: function delegateNonce(address ) view returns(uint256)
func (_DelegateManager *DelegateManagerCaller) DelegateNonce(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _DelegateManager.contract.Call(opts, &out, "delegateNonce", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DelegateNonce is a free data retrieval call binding the contract method 0x92e1f102.
//
// Solidity: function delegateNonce(address ) view returns(uint256)
func (_DelegateManager *DelegateManagerSession) DelegateNonce(arg0 common.Address) (*big.Int, error) {
	return _DelegateManager.Contract.DelegateNonce(&_DelegateManager.CallOpts, arg0)
}

// DelegateNonce is a free data retrieval call binding the contract method 0x92e1f102.
//
// Solidity: function delegateNonce(address ) view returns(uint256)
func (_DelegateManager *DelegateManagerCallerSession) DelegateNonce(arg0 common.Address) (*big.Int, error) {
	return _DelegateManager.Contract.DelegateNonce(&_DelegateManager.CallOpts, arg0)
}

// GetFiller is a free data retrieval call binding the contract method 0x0b3ff53f.
//
// Solidity: function getFiller(address filler) view returns(uint16 feeInBips, uint256 stake, uint256 deregisteredAt, bytes32[] delegateStakeIDs)
func (_DelegateManager *DelegateManagerCaller) GetFiller(opts *bind.CallOpts, filler common.Address) (struct {
	FeeInBips        uint16
	Stake            *big.Int
	DeregisteredAt   *big.Int
	DelegateStakeIDs [][32]byte
}, error) {
	var out []interface{}
	err := _DelegateManager.contract.Call(opts, &out, "getFiller", filler)

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
func (_DelegateManager *DelegateManagerSession) GetFiller(filler common.Address) (struct {
	FeeInBips        uint16
	Stake            *big.Int
	DeregisteredAt   *big.Int
	DelegateStakeIDs [][32]byte
}, error) {
	return _DelegateManager.Contract.GetFiller(&_DelegateManager.CallOpts, filler)
}

// GetFiller is a free data retrieval call binding the contract method 0x0b3ff53f.
//
// Solidity: function getFiller(address filler) view returns(uint16 feeInBips, uint256 stake, uint256 deregisteredAt, bytes32[] delegateStakeIDs)
func (_DelegateManager *DelegateManagerCallerSession) GetFiller(filler common.Address) (struct {
	FeeInBips        uint16
	Stake            *big.Int
	DeregisteredAt   *big.Int
	DelegateStakeIDs [][32]byte
}, error) {
	return _DelegateManager.Contract.GetFiller(&_DelegateManager.CallOpts, filler)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DelegateManager *DelegateManagerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _DelegateManager.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DelegateManager *DelegateManagerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _DelegateManager.Contract.GetRoleAdmin(&_DelegateManager.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DelegateManager *DelegateManagerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _DelegateManager.Contract.GetRoleAdmin(&_DelegateManager.CallOpts, role)
}

// GetVotes is a free data retrieval call binding the contract method 0x9ab24eb0.
//
// Solidity: function getVotes(address filler) view returns(uint256 voteCount)
func (_DelegateManager *DelegateManagerCaller) GetVotes(opts *bind.CallOpts, filler common.Address) (*big.Int, error) {
	var out []interface{}
	err := _DelegateManager.contract.Call(opts, &out, "getVotes", filler)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotes is a free data retrieval call binding the contract method 0x9ab24eb0.
//
// Solidity: function getVotes(address filler) view returns(uint256 voteCount)
func (_DelegateManager *DelegateManagerSession) GetVotes(filler common.Address) (*big.Int, error) {
	return _DelegateManager.Contract.GetVotes(&_DelegateManager.CallOpts, filler)
}

// GetVotes is a free data retrieval call binding the contract method 0x9ab24eb0.
//
// Solidity: function getVotes(address filler) view returns(uint256 voteCount)
func (_DelegateManager *DelegateManagerCallerSession) GetVotes(filler common.Address) (*big.Int, error) {
	return _DelegateManager.Contract.GetVotes(&_DelegateManager.CallOpts, filler)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DelegateManager *DelegateManagerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _DelegateManager.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DelegateManager *DelegateManagerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _DelegateManager.Contract.HasRole(&_DelegateManager.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DelegateManager *DelegateManagerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _DelegateManager.Contract.HasRole(&_DelegateManager.CallOpts, role, account)
}

// Stakes is a free data retrieval call binding the contract method 0x8fee6407.
//
// Solidity: function stakes(bytes32 ) view returns(address owner, uint256 stake, uint256 units, uint256 votes, address filler, uint256 expiry)
func (_DelegateManager *DelegateManagerCaller) Stakes(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Owner  common.Address
	Stake  *big.Int
	Units  *big.Int
	Votes  *big.Int
	Filler common.Address
	Expiry *big.Int
}, error) {
	var out []interface{}
	err := _DelegateManager.contract.Call(opts, &out, "stakes", arg0)

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
func (_DelegateManager *DelegateManagerSession) Stakes(arg0 [32]byte) (struct {
	Owner  common.Address
	Stake  *big.Int
	Units  *big.Int
	Votes  *big.Int
	Filler common.Address
	Expiry *big.Int
}, error) {
	return _DelegateManager.Contract.Stakes(&_DelegateManager.CallOpts, arg0)
}

// Stakes is a free data retrieval call binding the contract method 0x8fee6407.
//
// Solidity: function stakes(bytes32 ) view returns(address owner, uint256 stake, uint256 units, uint256 votes, address filler, uint256 expiry)
func (_DelegateManager *DelegateManagerCallerSession) Stakes(arg0 [32]byte) (struct {
	Owner  common.Address
	Stake  *big.Int
	Units  *big.Int
	Votes  *big.Int
	Filler common.Address
	Expiry *big.Int
}, error) {
	return _DelegateManager.Contract.Stakes(&_DelegateManager.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DelegateManager *DelegateManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _DelegateManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DelegateManager *DelegateManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DelegateManager.Contract.SupportsInterface(&_DelegateManager.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DelegateManager *DelegateManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DelegateManager.Contract.SupportsInterface(&_DelegateManager.CallOpts, interfaceId)
}

// ChangeVote is a paid mutator transaction binding the contract method 0xb9b0b22b.
//
// Solidity: function changeVote(bytes32 stakeID, address newFiller) returns()
func (_DelegateManager *DelegateManagerTransactor) ChangeVote(opts *bind.TransactOpts, stakeID [32]byte, newFiller common.Address) (*types.Transaction, error) {
	return _DelegateManager.contract.Transact(opts, "changeVote", stakeID, newFiller)
}

// ChangeVote is a paid mutator transaction binding the contract method 0xb9b0b22b.
//
// Solidity: function changeVote(bytes32 stakeID, address newFiller) returns()
func (_DelegateManager *DelegateManagerSession) ChangeVote(stakeID [32]byte, newFiller common.Address) (*types.Transaction, error) {
	return _DelegateManager.Contract.ChangeVote(&_DelegateManager.TransactOpts, stakeID, newFiller)
}

// ChangeVote is a paid mutator transaction binding the contract method 0xb9b0b22b.
//
// Solidity: function changeVote(bytes32 stakeID, address newFiller) returns()
func (_DelegateManager *DelegateManagerTransactorSession) ChangeVote(stakeID [32]byte, newFiller common.Address) (*types.Transaction, error) {
	return _DelegateManager.Contract.ChangeVote(&_DelegateManager.TransactOpts, stakeID, newFiller)
}

// Extend is a paid mutator transaction binding the contract method 0x6df3e83f.
//
// Solidity: function extend(bytes32 stakeID, uint256 newLockBlocks) returns()
func (_DelegateManager *DelegateManagerTransactor) Extend(opts *bind.TransactOpts, stakeID [32]byte, newLockBlocks *big.Int) (*types.Transaction, error) {
	return _DelegateManager.contract.Transact(opts, "extend", stakeID, newLockBlocks)
}

// Extend is a paid mutator transaction binding the contract method 0x6df3e83f.
//
// Solidity: function extend(bytes32 stakeID, uint256 newLockBlocks) returns()
func (_DelegateManager *DelegateManagerSession) Extend(stakeID [32]byte, newLockBlocks *big.Int) (*types.Transaction, error) {
	return _DelegateManager.Contract.Extend(&_DelegateManager.TransactOpts, stakeID, newLockBlocks)
}

// Extend is a paid mutator transaction binding the contract method 0x6df3e83f.
//
// Solidity: function extend(bytes32 stakeID, uint256 newLockBlocks) returns()
func (_DelegateManager *DelegateManagerTransactorSession) Extend(stakeID [32]byte, newLockBlocks *big.Int) (*types.Transaction, error) {
	return _DelegateManager.Contract.Extend(&_DelegateManager.TransactOpts, stakeID, newLockBlocks)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DelegateManager *DelegateManagerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegateManager.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DelegateManager *DelegateManagerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegateManager.Contract.GrantRole(&_DelegateManager.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DelegateManager *DelegateManagerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegateManager.Contract.GrantRole(&_DelegateManager.TransactOpts, role, account)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 stakeID) returns()
func (_DelegateManager *DelegateManagerTransactor) Refund(opts *bind.TransactOpts, stakeID [32]byte) (*types.Transaction, error) {
	return _DelegateManager.contract.Transact(opts, "refund", stakeID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 stakeID) returns()
func (_DelegateManager *DelegateManagerSession) Refund(stakeID [32]byte) (*types.Transaction, error) {
	return _DelegateManager.Contract.Refund(&_DelegateManager.TransactOpts, stakeID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 stakeID) returns()
func (_DelegateManager *DelegateManagerTransactorSession) Refund(stakeID [32]byte) (*types.Transaction, error) {
	return _DelegateManager.Contract.Refund(&_DelegateManager.TransactOpts, stakeID)
}

// Renew is a paid mutator transaction binding the contract method 0xf544d82f.
//
// Solidity: function renew(bytes32 stakeID, uint256 newLockBlocks) returns()
func (_DelegateManager *DelegateManagerTransactor) Renew(opts *bind.TransactOpts, stakeID [32]byte, newLockBlocks *big.Int) (*types.Transaction, error) {
	return _DelegateManager.contract.Transact(opts, "renew", stakeID, newLockBlocks)
}

// Renew is a paid mutator transaction binding the contract method 0xf544d82f.
//
// Solidity: function renew(bytes32 stakeID, uint256 newLockBlocks) returns()
func (_DelegateManager *DelegateManagerSession) Renew(stakeID [32]byte, newLockBlocks *big.Int) (*types.Transaction, error) {
	return _DelegateManager.Contract.Renew(&_DelegateManager.TransactOpts, stakeID, newLockBlocks)
}

// Renew is a paid mutator transaction binding the contract method 0xf544d82f.
//
// Solidity: function renew(bytes32 stakeID, uint256 newLockBlocks) returns()
func (_DelegateManager *DelegateManagerTransactorSession) Renew(stakeID [32]byte, newLockBlocks *big.Int) (*types.Transaction, error) {
	return _DelegateManager.Contract.Renew(&_DelegateManager.TransactOpts, stakeID, newLockBlocks)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_DelegateManager *DelegateManagerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegateManager.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_DelegateManager *DelegateManagerSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegateManager.Contract.RenounceRole(&_DelegateManager.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_DelegateManager *DelegateManagerTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegateManager.Contract.RenounceRole(&_DelegateManager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DelegateManager *DelegateManagerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegateManager.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DelegateManager *DelegateManagerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegateManager.Contract.RevokeRole(&_DelegateManager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DelegateManager *DelegateManagerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegateManager.Contract.RevokeRole(&_DelegateManager.TransactOpts, role, account)
}

// Vote is a paid mutator transaction binding the contract method 0x2a4a1b73.
//
// Solidity: function vote(address filler, uint256 units, uint256 lockBlocks) returns(bytes32 stakeID)
func (_DelegateManager *DelegateManagerTransactor) Vote(opts *bind.TransactOpts, filler common.Address, units *big.Int, lockBlocks *big.Int) (*types.Transaction, error) {
	return _DelegateManager.contract.Transact(opts, "vote", filler, units, lockBlocks)
}

// Vote is a paid mutator transaction binding the contract method 0x2a4a1b73.
//
// Solidity: function vote(address filler, uint256 units, uint256 lockBlocks) returns(bytes32 stakeID)
func (_DelegateManager *DelegateManagerSession) Vote(filler common.Address, units *big.Int, lockBlocks *big.Int) (*types.Transaction, error) {
	return _DelegateManager.Contract.Vote(&_DelegateManager.TransactOpts, filler, units, lockBlocks)
}

// Vote is a paid mutator transaction binding the contract method 0x2a4a1b73.
//
// Solidity: function vote(address filler, uint256 units, uint256 lockBlocks) returns(bytes32 stakeID)
func (_DelegateManager *DelegateManagerTransactorSession) Vote(filler common.Address, units *big.Int, lockBlocks *big.Int) (*types.Transaction, error) {
	return _DelegateManager.Contract.Vote(&_DelegateManager.TransactOpts, filler, units, lockBlocks)
}

// DelegateManagerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the DelegateManager contract.
type DelegateManagerRoleAdminChangedIterator struct {
	Event *DelegateManagerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *DelegateManagerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateManagerRoleAdminChanged)
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
		it.Event = new(DelegateManagerRoleAdminChanged)
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
func (it *DelegateManagerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateManagerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateManagerRoleAdminChanged represents a RoleAdminChanged event raised by the DelegateManager contract.
type DelegateManagerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_DelegateManager *DelegateManagerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*DelegateManagerRoleAdminChangedIterator, error) {

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

	logs, sub, err := _DelegateManager.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &DelegateManagerRoleAdminChangedIterator{contract: _DelegateManager.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_DelegateManager *DelegateManagerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *DelegateManagerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _DelegateManager.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateManagerRoleAdminChanged)
				if err := _DelegateManager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_DelegateManager *DelegateManagerFilterer) ParseRoleAdminChanged(log types.Log) (*DelegateManagerRoleAdminChanged, error) {
	event := new(DelegateManagerRoleAdminChanged)
	if err := _DelegateManager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegateManagerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the DelegateManager contract.
type DelegateManagerRoleGrantedIterator struct {
	Event *DelegateManagerRoleGranted // Event containing the contract specifics and raw log

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
func (it *DelegateManagerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateManagerRoleGranted)
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
		it.Event = new(DelegateManagerRoleGranted)
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
func (it *DelegateManagerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateManagerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateManagerRoleGranted represents a RoleGranted event raised by the DelegateManager contract.
type DelegateManagerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_DelegateManager *DelegateManagerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*DelegateManagerRoleGrantedIterator, error) {

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

	logs, sub, err := _DelegateManager.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &DelegateManagerRoleGrantedIterator{contract: _DelegateManager.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_DelegateManager *DelegateManagerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *DelegateManagerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _DelegateManager.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateManagerRoleGranted)
				if err := _DelegateManager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_DelegateManager *DelegateManagerFilterer) ParseRoleGranted(log types.Log) (*DelegateManagerRoleGranted, error) {
	event := new(DelegateManagerRoleGranted)
	if err := _DelegateManager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegateManagerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the DelegateManager contract.
type DelegateManagerRoleRevokedIterator struct {
	Event *DelegateManagerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *DelegateManagerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateManagerRoleRevoked)
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
		it.Event = new(DelegateManagerRoleRevoked)
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
func (it *DelegateManagerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateManagerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateManagerRoleRevoked represents a RoleRevoked event raised by the DelegateManager contract.
type DelegateManagerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_DelegateManager *DelegateManagerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*DelegateManagerRoleRevokedIterator, error) {

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

	logs, sub, err := _DelegateManager.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &DelegateManagerRoleRevokedIterator{contract: _DelegateManager.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_DelegateManager *DelegateManagerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *DelegateManagerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _DelegateManager.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateManagerRoleRevoked)
				if err := _DelegateManager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_DelegateManager *DelegateManagerFilterer) ParseRoleRevoked(log types.Log) (*DelegateManagerRoleRevoked, error) {
	event := new(DelegateManagerRoleRevoked)
	if err := _DelegateManager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegateManagerStakeExtendedIterator is returned from FilterStakeExtended and is used to iterate over the raw logs and unpacked data for StakeExtended events raised by the DelegateManager contract.
type DelegateManagerStakeExtendedIterator struct {
	Event *DelegateManagerStakeExtended // Event containing the contract specifics and raw log

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
func (it *DelegateManagerStakeExtendedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateManagerStakeExtended)
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
		it.Event = new(DelegateManagerStakeExtended)
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
func (it *DelegateManagerStakeExtendedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateManagerStakeExtendedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateManagerStakeExtended represents a StakeExtended event raised by the DelegateManager contract.
type DelegateManagerStakeExtended struct {
	StakeID       [32]byte
	NewLockBlocks *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterStakeExtended is a free log retrieval operation binding the contract event 0x521ebd8078bb6d9f423366c9ae8a0d4dcec7887f1db2d3c0cd0923fe00ed2ed7.
//
// Solidity: event StakeExtended(bytes32 indexed stakeID, uint256 newLockBlocks)
func (_DelegateManager *DelegateManagerFilterer) FilterStakeExtended(opts *bind.FilterOpts, stakeID [][32]byte) (*DelegateManagerStakeExtendedIterator, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}

	logs, sub, err := _DelegateManager.contract.FilterLogs(opts, "StakeExtended", stakeIDRule)
	if err != nil {
		return nil, err
	}
	return &DelegateManagerStakeExtendedIterator{contract: _DelegateManager.contract, event: "StakeExtended", logs: logs, sub: sub}, nil
}

// WatchStakeExtended is a free log subscription operation binding the contract event 0x521ebd8078bb6d9f423366c9ae8a0d4dcec7887f1db2d3c0cd0923fe00ed2ed7.
//
// Solidity: event StakeExtended(bytes32 indexed stakeID, uint256 newLockBlocks)
func (_DelegateManager *DelegateManagerFilterer) WatchStakeExtended(opts *bind.WatchOpts, sink chan<- *DelegateManagerStakeExtended, stakeID [][32]byte) (event.Subscription, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}

	logs, sub, err := _DelegateManager.contract.WatchLogs(opts, "StakeExtended", stakeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateManagerStakeExtended)
				if err := _DelegateManager.contract.UnpackLog(event, "StakeExtended", log); err != nil {
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
func (_DelegateManager *DelegateManagerFilterer) ParseStakeExtended(log types.Log) (*DelegateManagerStakeExtended, error) {
	event := new(DelegateManagerStakeExtended)
	if err := _DelegateManager.contract.UnpackLog(event, "StakeExtended", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegateManagerStakeRefundedIterator is returned from FilterStakeRefunded and is used to iterate over the raw logs and unpacked data for StakeRefunded events raised by the DelegateManager contract.
type DelegateManagerStakeRefundedIterator struct {
	Event *DelegateManagerStakeRefunded // Event containing the contract specifics and raw log

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
func (it *DelegateManagerStakeRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateManagerStakeRefunded)
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
		it.Event = new(DelegateManagerStakeRefunded)
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
func (it *DelegateManagerStakeRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateManagerStakeRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateManagerStakeRefunded represents a StakeRefunded event raised by the DelegateManager contract.
type DelegateManagerStakeRefunded struct {
	StakeID [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterStakeRefunded is a free log retrieval operation binding the contract event 0xba33cc6a4502f7f80bfe643ffa925ba7ab5e0af61c061a9aceabef0349c9dce4.
//
// Solidity: event StakeRefunded(bytes32 indexed stakeID)
func (_DelegateManager *DelegateManagerFilterer) FilterStakeRefunded(opts *bind.FilterOpts, stakeID [][32]byte) (*DelegateManagerStakeRefundedIterator, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}

	logs, sub, err := _DelegateManager.contract.FilterLogs(opts, "StakeRefunded", stakeIDRule)
	if err != nil {
		return nil, err
	}
	return &DelegateManagerStakeRefundedIterator{contract: _DelegateManager.contract, event: "StakeRefunded", logs: logs, sub: sub}, nil
}

// WatchStakeRefunded is a free log subscription operation binding the contract event 0xba33cc6a4502f7f80bfe643ffa925ba7ab5e0af61c061a9aceabef0349c9dce4.
//
// Solidity: event StakeRefunded(bytes32 indexed stakeID)
func (_DelegateManager *DelegateManagerFilterer) WatchStakeRefunded(opts *bind.WatchOpts, sink chan<- *DelegateManagerStakeRefunded, stakeID [][32]byte) (event.Subscription, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}

	logs, sub, err := _DelegateManager.contract.WatchLogs(opts, "StakeRefunded", stakeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateManagerStakeRefunded)
				if err := _DelegateManager.contract.UnpackLog(event, "StakeRefunded", log); err != nil {
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
func (_DelegateManager *DelegateManagerFilterer) ParseStakeRefunded(log types.Log) (*DelegateManagerStakeRefunded, error) {
	event := new(DelegateManagerStakeRefunded)
	if err := _DelegateManager.contract.UnpackLog(event, "StakeRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegateManagerStakeRenewedIterator is returned from FilterStakeRenewed and is used to iterate over the raw logs and unpacked data for StakeRenewed events raised by the DelegateManager contract.
type DelegateManagerStakeRenewedIterator struct {
	Event *DelegateManagerStakeRenewed // Event containing the contract specifics and raw log

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
func (it *DelegateManagerStakeRenewedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateManagerStakeRenewed)
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
		it.Event = new(DelegateManagerStakeRenewed)
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
func (it *DelegateManagerStakeRenewedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateManagerStakeRenewedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateManagerStakeRenewed represents a StakeRenewed event raised by the DelegateManager contract.
type DelegateManagerStakeRenewed struct {
	StakeID       [32]byte
	NewLockBlocks *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterStakeRenewed is a free log retrieval operation binding the contract event 0xe2c4eee10fe6a93f81bf01601e62cba732126941bccca75f3460e6f2f845513d.
//
// Solidity: event StakeRenewed(bytes32 indexed stakeID, uint256 newLockBlocks)
func (_DelegateManager *DelegateManagerFilterer) FilterStakeRenewed(opts *bind.FilterOpts, stakeID [][32]byte) (*DelegateManagerStakeRenewedIterator, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}

	logs, sub, err := _DelegateManager.contract.FilterLogs(opts, "StakeRenewed", stakeIDRule)
	if err != nil {
		return nil, err
	}
	return &DelegateManagerStakeRenewedIterator{contract: _DelegateManager.contract, event: "StakeRenewed", logs: logs, sub: sub}, nil
}

// WatchStakeRenewed is a free log subscription operation binding the contract event 0xe2c4eee10fe6a93f81bf01601e62cba732126941bccca75f3460e6f2f845513d.
//
// Solidity: event StakeRenewed(bytes32 indexed stakeID, uint256 newLockBlocks)
func (_DelegateManager *DelegateManagerFilterer) WatchStakeRenewed(opts *bind.WatchOpts, sink chan<- *DelegateManagerStakeRenewed, stakeID [][32]byte) (event.Subscription, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}

	logs, sub, err := _DelegateManager.contract.WatchLogs(opts, "StakeRenewed", stakeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateManagerStakeRenewed)
				if err := _DelegateManager.contract.UnpackLog(event, "StakeRenewed", log); err != nil {
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
func (_DelegateManager *DelegateManagerFilterer) ParseStakeRenewed(log types.Log) (*DelegateManagerStakeRenewed, error) {
	event := new(DelegateManagerStakeRenewed)
	if err := _DelegateManager.contract.UnpackLog(event, "StakeRenewed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegateManagerStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the DelegateManager contract.
type DelegateManagerStakedIterator struct {
	Event *DelegateManagerStaked // Event containing the contract specifics and raw log

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
func (it *DelegateManagerStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateManagerStaked)
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
		it.Event = new(DelegateManagerStaked)
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
func (it *DelegateManagerStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateManagerStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateManagerStaked represents a Staked event raised by the DelegateManager contract.
type DelegateManagerStaked struct {
	StakeID [32]byte
	Owner   common.Address
	Stake   *big.Int
	Expiry  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0xa1fdccfe567643a44425efdd141171e8d992854a81e5c819c1432b0de47c9a11.
//
// Solidity: event Staked(bytes32 indexed stakeID, address indexed owner, uint256 stake, uint256 expiry)
func (_DelegateManager *DelegateManagerFilterer) FilterStaked(opts *bind.FilterOpts, stakeID [][32]byte, owner []common.Address) (*DelegateManagerStakedIterator, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _DelegateManager.contract.FilterLogs(opts, "Staked", stakeIDRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &DelegateManagerStakedIterator{contract: _DelegateManager.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0xa1fdccfe567643a44425efdd141171e8d992854a81e5c819c1432b0de47c9a11.
//
// Solidity: event Staked(bytes32 indexed stakeID, address indexed owner, uint256 stake, uint256 expiry)
func (_DelegateManager *DelegateManagerFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *DelegateManagerStaked, stakeID [][32]byte, owner []common.Address) (event.Subscription, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _DelegateManager.contract.WatchLogs(opts, "Staked", stakeIDRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateManagerStaked)
				if err := _DelegateManager.contract.UnpackLog(event, "Staked", log); err != nil {
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
func (_DelegateManager *DelegateManagerFilterer) ParseStaked(log types.Log) (*DelegateManagerStaked, error) {
	event := new(DelegateManagerStaked)
	if err := _DelegateManager.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegateManagerVotedIterator is returned from FilterVoted and is used to iterate over the raw logs and unpacked data for Voted events raised by the DelegateManager contract.
type DelegateManagerVotedIterator struct {
	Event *DelegateManagerVoted // Event containing the contract specifics and raw log

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
func (it *DelegateManagerVotedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateManagerVoted)
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
		it.Event = new(DelegateManagerVoted)
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
func (it *DelegateManagerVotedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateManagerVotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateManagerVoted represents a Voted event raised by the DelegateManager contract.
type DelegateManagerVoted struct {
	StakeID [32]byte
	Filler  common.Address
	Votes   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterVoted is a free log retrieval operation binding the contract event 0xe4abc5380fa6939d1dc23b5e90b3a8a0e328f0f1a82a5f42bfb795bf9c717505.
//
// Solidity: event Voted(bytes32 indexed stakeID, address indexed filler, uint256 votes)
func (_DelegateManager *DelegateManagerFilterer) FilterVoted(opts *bind.FilterOpts, stakeID [][32]byte, filler []common.Address) (*DelegateManagerVotedIterator, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}
	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _DelegateManager.contract.FilterLogs(opts, "Voted", stakeIDRule, fillerRule)
	if err != nil {
		return nil, err
	}
	return &DelegateManagerVotedIterator{contract: _DelegateManager.contract, event: "Voted", logs: logs, sub: sub}, nil
}

// WatchVoted is a free log subscription operation binding the contract event 0xe4abc5380fa6939d1dc23b5e90b3a8a0e328f0f1a82a5f42bfb795bf9c717505.
//
// Solidity: event Voted(bytes32 indexed stakeID, address indexed filler, uint256 votes)
func (_DelegateManager *DelegateManagerFilterer) WatchVoted(opts *bind.WatchOpts, sink chan<- *DelegateManagerVoted, stakeID [][32]byte, filler []common.Address) (event.Subscription, error) {

	var stakeIDRule []interface{}
	for _, stakeIDItem := range stakeID {
		stakeIDRule = append(stakeIDRule, stakeIDItem)
	}
	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _DelegateManager.contract.WatchLogs(opts, "Voted", stakeIDRule, fillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateManagerVoted)
				if err := _DelegateManager.contract.UnpackLog(event, "Voted", log); err != nil {
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
func (_DelegateManager *DelegateManagerFilterer) ParseVoted(log types.Log) (*DelegateManagerVoted, error) {
	event := new(DelegateManagerVoted)
	if err := _DelegateManager.contract.UnpackLog(event, "Voted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegateManagerVotesChangedIterator is returned from FilterVotesChanged and is used to iterate over the raw logs and unpacked data for VotesChanged events raised by the DelegateManager contract.
type DelegateManagerVotesChangedIterator struct {
	Event *DelegateManagerVotesChanged // Event containing the contract specifics and raw log

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
func (it *DelegateManagerVotesChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateManagerVotesChanged)
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
		it.Event = new(DelegateManagerVotesChanged)
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
func (it *DelegateManagerVotesChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateManagerVotesChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateManagerVotesChanged represents a VotesChanged event raised by the DelegateManager contract.
type DelegateManagerVotesChanged struct {
	StakeID   [32]byte
	OldFiller common.Address
	NewFiller common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterVotesChanged is a free log retrieval operation binding the contract event 0xf3a02e43bc9e793a4fe98273c318b73e54fbbd0c599adb2b620ade22b2235f6e.
//
// Solidity: event VotesChanged(bytes32 indexed stakeID, address indexed oldFiller, address indexed newFiller)
func (_DelegateManager *DelegateManagerFilterer) FilterVotesChanged(opts *bind.FilterOpts, stakeID [][32]byte, oldFiller []common.Address, newFiller []common.Address) (*DelegateManagerVotesChangedIterator, error) {

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

	logs, sub, err := _DelegateManager.contract.FilterLogs(opts, "VotesChanged", stakeIDRule, oldFillerRule, newFillerRule)
	if err != nil {
		return nil, err
	}
	return &DelegateManagerVotesChangedIterator{contract: _DelegateManager.contract, event: "VotesChanged", logs: logs, sub: sub}, nil
}

// WatchVotesChanged is a free log subscription operation binding the contract event 0xf3a02e43bc9e793a4fe98273c318b73e54fbbd0c599adb2b620ade22b2235f6e.
//
// Solidity: event VotesChanged(bytes32 indexed stakeID, address indexed oldFiller, address indexed newFiller)
func (_DelegateManager *DelegateManagerFilterer) WatchVotesChanged(opts *bind.WatchOpts, sink chan<- *DelegateManagerVotesChanged, stakeID [][32]byte, oldFiller []common.Address, newFiller []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _DelegateManager.contract.WatchLogs(opts, "VotesChanged", stakeIDRule, oldFillerRule, newFillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateManagerVotesChanged)
				if err := _DelegateManager.contract.UnpackLog(event, "VotesChanged", log); err != nil {
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
func (_DelegateManager *DelegateManagerFilterer) ParseVotesChanged(log types.Log) (*DelegateManagerVotesChanged, error) {
	event := new(DelegateManagerVotesChanged)
	if err := _DelegateManager.contract.UnpackLog(event, "VotesChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
