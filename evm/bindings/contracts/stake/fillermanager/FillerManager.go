// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package fillermanager

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

// FillerManagerMetaData contains all meta data concerning the FillerManager contract.
var FillerManagerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deregisteredAt\",\"type\":\"uint256\"}],\"name\":\"FillerDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"FillerFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"}],\"name\":\"FillerRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"}],\"name\":\"FillerRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DELEGATE_STAKE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FILLER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FILLER_COOL_DOWN\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FILLER_STAKE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SEED\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"delegateNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deregister\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"}],\"name\":\"getFiller\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"feeInBips\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deregisteredAt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"delegateStakeIDs\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"filler_\",\"type\":\"address\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"stakes\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"units\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"newFee\",\"type\":\"uint16\"}],\"name\":\"updateFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// FillerManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use FillerManagerMetaData.ABI instead.
var FillerManagerABI = FillerManagerMetaData.ABI

// FillerManager is an auto generated Go binding around an Ethereum contract.
type FillerManager struct {
	FillerManagerCaller     // Read-only binding to the contract
	FillerManagerTransactor // Write-only binding to the contract
	FillerManagerFilterer   // Log filterer for contract events
}

// FillerManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type FillerManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FillerManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FillerManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FillerManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FillerManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FillerManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FillerManagerSession struct {
	Contract     *FillerManager    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FillerManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FillerManagerCallerSession struct {
	Contract *FillerManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// FillerManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FillerManagerTransactorSession struct {
	Contract     *FillerManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// FillerManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type FillerManagerRaw struct {
	Contract *FillerManager // Generic contract binding to access the raw methods on
}

// FillerManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FillerManagerCallerRaw struct {
	Contract *FillerManagerCaller // Generic read-only contract binding to access the raw methods on
}

// FillerManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FillerManagerTransactorRaw struct {
	Contract *FillerManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFillerManager creates a new instance of FillerManager, bound to a specific deployed contract.
func NewFillerManager(address common.Address, backend bind.ContractBackend) (*FillerManager, error) {
	contract, err := bindFillerManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FillerManager{FillerManagerCaller: FillerManagerCaller{contract: contract}, FillerManagerTransactor: FillerManagerTransactor{contract: contract}, FillerManagerFilterer: FillerManagerFilterer{contract: contract}}, nil
}

// NewFillerManagerCaller creates a new read-only instance of FillerManager, bound to a specific deployed contract.
func NewFillerManagerCaller(address common.Address, caller bind.ContractCaller) (*FillerManagerCaller, error) {
	contract, err := bindFillerManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FillerManagerCaller{contract: contract}, nil
}

// NewFillerManagerTransactor creates a new write-only instance of FillerManager, bound to a specific deployed contract.
func NewFillerManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*FillerManagerTransactor, error) {
	contract, err := bindFillerManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FillerManagerTransactor{contract: contract}, nil
}

// NewFillerManagerFilterer creates a new log filterer instance of FillerManager, bound to a specific deployed contract.
func NewFillerManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*FillerManagerFilterer, error) {
	contract, err := bindFillerManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FillerManagerFilterer{contract: contract}, nil
}

// bindFillerManager binds a generic wrapper to an already deployed contract.
func bindFillerManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FillerManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FillerManager *FillerManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FillerManager.Contract.FillerManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FillerManager *FillerManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FillerManager.Contract.FillerManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FillerManager *FillerManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FillerManager.Contract.FillerManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FillerManager *FillerManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FillerManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FillerManager *FillerManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FillerManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FillerManager *FillerManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FillerManager.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_FillerManager *FillerManagerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FillerManager.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_FillerManager *FillerManagerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _FillerManager.Contract.DEFAULTADMINROLE(&_FillerManager.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_FillerManager *FillerManagerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _FillerManager.Contract.DEFAULTADMINROLE(&_FillerManager.CallOpts)
}

// DELEGATESTAKE is a free data retrieval call binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() view returns(uint256)
func (_FillerManager *FillerManagerCaller) DELEGATESTAKE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FillerManager.contract.Call(opts, &out, "DELEGATE_STAKE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DELEGATESTAKE is a free data retrieval call binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() view returns(uint256)
func (_FillerManager *FillerManagerSession) DELEGATESTAKE() (*big.Int, error) {
	return _FillerManager.Contract.DELEGATESTAKE(&_FillerManager.CallOpts)
}

// DELEGATESTAKE is a free data retrieval call binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() view returns(uint256)
func (_FillerManager *FillerManagerCallerSession) DELEGATESTAKE() (*big.Int, error) {
	return _FillerManager.Contract.DELEGATESTAKE(&_FillerManager.CallOpts)
}

// FILLER is a free data retrieval call binding the contract method 0x64645bae.
//
// Solidity: function FILLER() view returns(bytes32)
func (_FillerManager *FillerManagerCaller) FILLER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FillerManager.contract.Call(opts, &out, "FILLER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FILLER is a free data retrieval call binding the contract method 0x64645bae.
//
// Solidity: function FILLER() view returns(bytes32)
func (_FillerManager *FillerManagerSession) FILLER() ([32]byte, error) {
	return _FillerManager.Contract.FILLER(&_FillerManager.CallOpts)
}

// FILLER is a free data retrieval call binding the contract method 0x64645bae.
//
// Solidity: function FILLER() view returns(bytes32)
func (_FillerManager *FillerManagerCallerSession) FILLER() ([32]byte, error) {
	return _FillerManager.Contract.FILLER(&_FillerManager.CallOpts)
}

// FILLERCOOLDOWN is a free data retrieval call binding the contract method 0x015a44e8.
//
// Solidity: function FILLER_COOL_DOWN() view returns(uint256)
func (_FillerManager *FillerManagerCaller) FILLERCOOLDOWN(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FillerManager.contract.Call(opts, &out, "FILLER_COOL_DOWN")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FILLERCOOLDOWN is a free data retrieval call binding the contract method 0x015a44e8.
//
// Solidity: function FILLER_COOL_DOWN() view returns(uint256)
func (_FillerManager *FillerManagerSession) FILLERCOOLDOWN() (*big.Int, error) {
	return _FillerManager.Contract.FILLERCOOLDOWN(&_FillerManager.CallOpts)
}

// FILLERCOOLDOWN is a free data retrieval call binding the contract method 0x015a44e8.
//
// Solidity: function FILLER_COOL_DOWN() view returns(uint256)
func (_FillerManager *FillerManagerCallerSession) FILLERCOOLDOWN() (*big.Int, error) {
	return _FillerManager.Contract.FILLERCOOLDOWN(&_FillerManager.CallOpts)
}

// FILLERSTAKE is a free data retrieval call binding the contract method 0xecfe923e.
//
// Solidity: function FILLER_STAKE() view returns(uint256)
func (_FillerManager *FillerManagerCaller) FILLERSTAKE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FillerManager.contract.Call(opts, &out, "FILLER_STAKE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FILLERSTAKE is a free data retrieval call binding the contract method 0xecfe923e.
//
// Solidity: function FILLER_STAKE() view returns(uint256)
func (_FillerManager *FillerManagerSession) FILLERSTAKE() (*big.Int, error) {
	return _FillerManager.Contract.FILLERSTAKE(&_FillerManager.CallOpts)
}

// FILLERSTAKE is a free data retrieval call binding the contract method 0xecfe923e.
//
// Solidity: function FILLER_STAKE() view returns(uint256)
func (_FillerManager *FillerManagerCallerSession) FILLERSTAKE() (*big.Int, error) {
	return _FillerManager.Contract.FILLERSTAKE(&_FillerManager.CallOpts)
}

// SEED is a free data retrieval call binding the contract method 0x0edc4737.
//
// Solidity: function SEED() view returns(address)
func (_FillerManager *FillerManagerCaller) SEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FillerManager.contract.Call(opts, &out, "SEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SEED is a free data retrieval call binding the contract method 0x0edc4737.
//
// Solidity: function SEED() view returns(address)
func (_FillerManager *FillerManagerSession) SEED() (common.Address, error) {
	return _FillerManager.Contract.SEED(&_FillerManager.CallOpts)
}

// SEED is a free data retrieval call binding the contract method 0x0edc4737.
//
// Solidity: function SEED() view returns(address)
func (_FillerManager *FillerManagerCallerSession) SEED() (common.Address, error) {
	return _FillerManager.Contract.SEED(&_FillerManager.CallOpts)
}

// DelegateNonce is a free data retrieval call binding the contract method 0x92e1f102.
//
// Solidity: function delegateNonce(address ) view returns(uint256)
func (_FillerManager *FillerManagerCaller) DelegateNonce(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FillerManager.contract.Call(opts, &out, "delegateNonce", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DelegateNonce is a free data retrieval call binding the contract method 0x92e1f102.
//
// Solidity: function delegateNonce(address ) view returns(uint256)
func (_FillerManager *FillerManagerSession) DelegateNonce(arg0 common.Address) (*big.Int, error) {
	return _FillerManager.Contract.DelegateNonce(&_FillerManager.CallOpts, arg0)
}

// DelegateNonce is a free data retrieval call binding the contract method 0x92e1f102.
//
// Solidity: function delegateNonce(address ) view returns(uint256)
func (_FillerManager *FillerManagerCallerSession) DelegateNonce(arg0 common.Address) (*big.Int, error) {
	return _FillerManager.Contract.DelegateNonce(&_FillerManager.CallOpts, arg0)
}

// GetFiller is a free data retrieval call binding the contract method 0x0b3ff53f.
//
// Solidity: function getFiller(address filler) view returns(uint16 feeInBips, uint256 stake, uint256 deregisteredAt, bytes32[] delegateStakeIDs)
func (_FillerManager *FillerManagerCaller) GetFiller(opts *bind.CallOpts, filler common.Address) (struct {
	FeeInBips        uint16
	Stake            *big.Int
	DeregisteredAt   *big.Int
	DelegateStakeIDs [][32]byte
}, error) {
	var out []interface{}
	err := _FillerManager.contract.Call(opts, &out, "getFiller", filler)

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
func (_FillerManager *FillerManagerSession) GetFiller(filler common.Address) (struct {
	FeeInBips        uint16
	Stake            *big.Int
	DeregisteredAt   *big.Int
	DelegateStakeIDs [][32]byte
}, error) {
	return _FillerManager.Contract.GetFiller(&_FillerManager.CallOpts, filler)
}

// GetFiller is a free data retrieval call binding the contract method 0x0b3ff53f.
//
// Solidity: function getFiller(address filler) view returns(uint16 feeInBips, uint256 stake, uint256 deregisteredAt, bytes32[] delegateStakeIDs)
func (_FillerManager *FillerManagerCallerSession) GetFiller(filler common.Address) (struct {
	FeeInBips        uint16
	Stake            *big.Int
	DeregisteredAt   *big.Int
	DelegateStakeIDs [][32]byte
}, error) {
	return _FillerManager.Contract.GetFiller(&_FillerManager.CallOpts, filler)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_FillerManager *FillerManagerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _FillerManager.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_FillerManager *FillerManagerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _FillerManager.Contract.GetRoleAdmin(&_FillerManager.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_FillerManager *FillerManagerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _FillerManager.Contract.GetRoleAdmin(&_FillerManager.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_FillerManager *FillerManagerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _FillerManager.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_FillerManager *FillerManagerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _FillerManager.Contract.HasRole(&_FillerManager.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_FillerManager *FillerManagerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _FillerManager.Contract.HasRole(&_FillerManager.CallOpts, role, account)
}

// Stakes is a free data retrieval call binding the contract method 0x8fee6407.
//
// Solidity: function stakes(bytes32 ) view returns(address owner, uint256 stake, uint256 units, uint256 votes, address filler, uint256 expiry)
func (_FillerManager *FillerManagerCaller) Stakes(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Owner  common.Address
	Stake  *big.Int
	Units  *big.Int
	Votes  *big.Int
	Filler common.Address
	Expiry *big.Int
}, error) {
	var out []interface{}
	err := _FillerManager.contract.Call(opts, &out, "stakes", arg0)

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
func (_FillerManager *FillerManagerSession) Stakes(arg0 [32]byte) (struct {
	Owner  common.Address
	Stake  *big.Int
	Units  *big.Int
	Votes  *big.Int
	Filler common.Address
	Expiry *big.Int
}, error) {
	return _FillerManager.Contract.Stakes(&_FillerManager.CallOpts, arg0)
}

// Stakes is a free data retrieval call binding the contract method 0x8fee6407.
//
// Solidity: function stakes(bytes32 ) view returns(address owner, uint256 stake, uint256 units, uint256 votes, address filler, uint256 expiry)
func (_FillerManager *FillerManagerCallerSession) Stakes(arg0 [32]byte) (struct {
	Owner  common.Address
	Stake  *big.Int
	Units  *big.Int
	Votes  *big.Int
	Filler common.Address
	Expiry *big.Int
}, error) {
	return _FillerManager.Contract.Stakes(&_FillerManager.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_FillerManager *FillerManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _FillerManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_FillerManager *FillerManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _FillerManager.Contract.SupportsInterface(&_FillerManager.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_FillerManager *FillerManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _FillerManager.Contract.SupportsInterface(&_FillerManager.CallOpts, interfaceId)
}

// Deregister is a paid mutator transaction binding the contract method 0xaff5edb1.
//
// Solidity: function deregister() returns()
func (_FillerManager *FillerManagerTransactor) Deregister(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FillerManager.contract.Transact(opts, "deregister")
}

// Deregister is a paid mutator transaction binding the contract method 0xaff5edb1.
//
// Solidity: function deregister() returns()
func (_FillerManager *FillerManagerSession) Deregister() (*types.Transaction, error) {
	return _FillerManager.Contract.Deregister(&_FillerManager.TransactOpts)
}

// Deregister is a paid mutator transaction binding the contract method 0xaff5edb1.
//
// Solidity: function deregister() returns()
func (_FillerManager *FillerManagerTransactorSession) Deregister() (*types.Transaction, error) {
	return _FillerManager.Contract.Deregister(&_FillerManager.TransactOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_FillerManager *FillerManagerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FillerManager.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_FillerManager *FillerManagerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FillerManager.Contract.GrantRole(&_FillerManager.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_FillerManager *FillerManagerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FillerManager.Contract.GrantRole(&_FillerManager.TransactOpts, role, account)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address filler_) returns()
func (_FillerManager *FillerManagerTransactor) Refund(opts *bind.TransactOpts, filler_ common.Address) (*types.Transaction, error) {
	return _FillerManager.contract.Transact(opts, "refund", filler_)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address filler_) returns()
func (_FillerManager *FillerManagerSession) Refund(filler_ common.Address) (*types.Transaction, error) {
	return _FillerManager.Contract.Refund(&_FillerManager.TransactOpts, filler_)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address filler_) returns()
func (_FillerManager *FillerManagerTransactorSession) Refund(filler_ common.Address) (*types.Transaction, error) {
	return _FillerManager.Contract.Refund(&_FillerManager.TransactOpts, filler_)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_FillerManager *FillerManagerTransactor) Register(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FillerManager.contract.Transact(opts, "register")
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_FillerManager *FillerManagerSession) Register() (*types.Transaction, error) {
	return _FillerManager.Contract.Register(&_FillerManager.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_FillerManager *FillerManagerTransactorSession) Register() (*types.Transaction, error) {
	return _FillerManager.Contract.Register(&_FillerManager.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_FillerManager *FillerManagerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FillerManager.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_FillerManager *FillerManagerSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FillerManager.Contract.RenounceRole(&_FillerManager.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_FillerManager *FillerManagerTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FillerManager.Contract.RenounceRole(&_FillerManager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_FillerManager *FillerManagerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FillerManager.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_FillerManager *FillerManagerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FillerManager.Contract.RevokeRole(&_FillerManager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_FillerManager *FillerManagerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FillerManager.Contract.RevokeRole(&_FillerManager.TransactOpts, role, account)
}

// UpdateFee is a paid mutator transaction binding the contract method 0x2c6cda93.
//
// Solidity: function updateFee(uint16 newFee) returns()
func (_FillerManager *FillerManagerTransactor) UpdateFee(opts *bind.TransactOpts, newFee uint16) (*types.Transaction, error) {
	return _FillerManager.contract.Transact(opts, "updateFee", newFee)
}

// UpdateFee is a paid mutator transaction binding the contract method 0x2c6cda93.
//
// Solidity: function updateFee(uint16 newFee) returns()
func (_FillerManager *FillerManagerSession) UpdateFee(newFee uint16) (*types.Transaction, error) {
	return _FillerManager.Contract.UpdateFee(&_FillerManager.TransactOpts, newFee)
}

// UpdateFee is a paid mutator transaction binding the contract method 0x2c6cda93.
//
// Solidity: function updateFee(uint16 newFee) returns()
func (_FillerManager *FillerManagerTransactorSession) UpdateFee(newFee uint16) (*types.Transaction, error) {
	return _FillerManager.Contract.UpdateFee(&_FillerManager.TransactOpts, newFee)
}

// FillerManagerFillerDeregisteredIterator is returned from FilterFillerDeregistered and is used to iterate over the raw logs and unpacked data for FillerDeregistered events raised by the FillerManager contract.
type FillerManagerFillerDeregisteredIterator struct {
	Event *FillerManagerFillerDeregistered // Event containing the contract specifics and raw log

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
func (it *FillerManagerFillerDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FillerManagerFillerDeregistered)
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
		it.Event = new(FillerManagerFillerDeregistered)
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
func (it *FillerManagerFillerDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FillerManagerFillerDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FillerManagerFillerDeregistered represents a FillerDeregistered event raised by the FillerManager contract.
type FillerManagerFillerDeregistered struct {
	Filler         common.Address
	DeregisteredAt *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterFillerDeregistered is a free log retrieval operation binding the contract event 0x83f9b3ce2be58c675065c1675b20665690bc38864c4a0335ee681a57f9f88227.
//
// Solidity: event FillerDeregistered(address indexed filler, uint256 deregisteredAt)
func (_FillerManager *FillerManagerFilterer) FilterFillerDeregistered(opts *bind.FilterOpts, filler []common.Address) (*FillerManagerFillerDeregisteredIterator, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _FillerManager.contract.FilterLogs(opts, "FillerDeregistered", fillerRule)
	if err != nil {
		return nil, err
	}
	return &FillerManagerFillerDeregisteredIterator{contract: _FillerManager.contract, event: "FillerDeregistered", logs: logs, sub: sub}, nil
}

// WatchFillerDeregistered is a free log subscription operation binding the contract event 0x83f9b3ce2be58c675065c1675b20665690bc38864c4a0335ee681a57f9f88227.
//
// Solidity: event FillerDeregistered(address indexed filler, uint256 deregisteredAt)
func (_FillerManager *FillerManagerFilterer) WatchFillerDeregistered(opts *bind.WatchOpts, sink chan<- *FillerManagerFillerDeregistered, filler []common.Address) (event.Subscription, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _FillerManager.contract.WatchLogs(opts, "FillerDeregistered", fillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FillerManagerFillerDeregistered)
				if err := _FillerManager.contract.UnpackLog(event, "FillerDeregistered", log); err != nil {
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
func (_FillerManager *FillerManagerFilterer) ParseFillerDeregistered(log types.Log) (*FillerManagerFillerDeregistered, error) {
	event := new(FillerManagerFillerDeregistered)
	if err := _FillerManager.contract.UnpackLog(event, "FillerDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FillerManagerFillerFeeUpdatedIterator is returned from FilterFillerFeeUpdated and is used to iterate over the raw logs and unpacked data for FillerFeeUpdated events raised by the FillerManager contract.
type FillerManagerFillerFeeUpdatedIterator struct {
	Event *FillerManagerFillerFeeUpdated // Event containing the contract specifics and raw log

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
func (it *FillerManagerFillerFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FillerManagerFillerFeeUpdated)
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
		it.Event = new(FillerManagerFillerFeeUpdated)
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
func (it *FillerManagerFillerFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FillerManagerFillerFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FillerManagerFillerFeeUpdated represents a FillerFeeUpdated event raised by the FillerManager contract.
type FillerManagerFillerFeeUpdated struct {
	Filler common.Address
	Fee    *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFillerFeeUpdated is a free log retrieval operation binding the contract event 0x7318fde364eaeddb75aa88119c39d65d30c300d1e4279711a132b11922b9aa64.
//
// Solidity: event FillerFeeUpdated(address indexed filler, uint256 fee)
func (_FillerManager *FillerManagerFilterer) FilterFillerFeeUpdated(opts *bind.FilterOpts, filler []common.Address) (*FillerManagerFillerFeeUpdatedIterator, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _FillerManager.contract.FilterLogs(opts, "FillerFeeUpdated", fillerRule)
	if err != nil {
		return nil, err
	}
	return &FillerManagerFillerFeeUpdatedIterator{contract: _FillerManager.contract, event: "FillerFeeUpdated", logs: logs, sub: sub}, nil
}

// WatchFillerFeeUpdated is a free log subscription operation binding the contract event 0x7318fde364eaeddb75aa88119c39d65d30c300d1e4279711a132b11922b9aa64.
//
// Solidity: event FillerFeeUpdated(address indexed filler, uint256 fee)
func (_FillerManager *FillerManagerFilterer) WatchFillerFeeUpdated(opts *bind.WatchOpts, sink chan<- *FillerManagerFillerFeeUpdated, filler []common.Address) (event.Subscription, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _FillerManager.contract.WatchLogs(opts, "FillerFeeUpdated", fillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FillerManagerFillerFeeUpdated)
				if err := _FillerManager.contract.UnpackLog(event, "FillerFeeUpdated", log); err != nil {
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
func (_FillerManager *FillerManagerFilterer) ParseFillerFeeUpdated(log types.Log) (*FillerManagerFillerFeeUpdated, error) {
	event := new(FillerManagerFillerFeeUpdated)
	if err := _FillerManager.contract.UnpackLog(event, "FillerFeeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FillerManagerFillerRefundedIterator is returned from FilterFillerRefunded and is used to iterate over the raw logs and unpacked data for FillerRefunded events raised by the FillerManager contract.
type FillerManagerFillerRefundedIterator struct {
	Event *FillerManagerFillerRefunded // Event containing the contract specifics and raw log

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
func (it *FillerManagerFillerRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FillerManagerFillerRefunded)
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
		it.Event = new(FillerManagerFillerRefunded)
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
func (it *FillerManagerFillerRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FillerManagerFillerRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FillerManagerFillerRefunded represents a FillerRefunded event raised by the FillerManager contract.
type FillerManagerFillerRefunded struct {
	Filler common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFillerRefunded is a free log retrieval operation binding the contract event 0xbe2dbe5153e94feab433d0b81928c0c44ffadf824ea8b4692f7523e0f4c16cd6.
//
// Solidity: event FillerRefunded(address indexed filler)
func (_FillerManager *FillerManagerFilterer) FilterFillerRefunded(opts *bind.FilterOpts, filler []common.Address) (*FillerManagerFillerRefundedIterator, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _FillerManager.contract.FilterLogs(opts, "FillerRefunded", fillerRule)
	if err != nil {
		return nil, err
	}
	return &FillerManagerFillerRefundedIterator{contract: _FillerManager.contract, event: "FillerRefunded", logs: logs, sub: sub}, nil
}

// WatchFillerRefunded is a free log subscription operation binding the contract event 0xbe2dbe5153e94feab433d0b81928c0c44ffadf824ea8b4692f7523e0f4c16cd6.
//
// Solidity: event FillerRefunded(address indexed filler)
func (_FillerManager *FillerManagerFilterer) WatchFillerRefunded(opts *bind.WatchOpts, sink chan<- *FillerManagerFillerRefunded, filler []common.Address) (event.Subscription, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _FillerManager.contract.WatchLogs(opts, "FillerRefunded", fillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FillerManagerFillerRefunded)
				if err := _FillerManager.contract.UnpackLog(event, "FillerRefunded", log); err != nil {
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
func (_FillerManager *FillerManagerFilterer) ParseFillerRefunded(log types.Log) (*FillerManagerFillerRefunded, error) {
	event := new(FillerManagerFillerRefunded)
	if err := _FillerManager.contract.UnpackLog(event, "FillerRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FillerManagerFillerRegisteredIterator is returned from FilterFillerRegistered and is used to iterate over the raw logs and unpacked data for FillerRegistered events raised by the FillerManager contract.
type FillerManagerFillerRegisteredIterator struct {
	Event *FillerManagerFillerRegistered // Event containing the contract specifics and raw log

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
func (it *FillerManagerFillerRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FillerManagerFillerRegistered)
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
		it.Event = new(FillerManagerFillerRegistered)
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
func (it *FillerManagerFillerRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FillerManagerFillerRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FillerManagerFillerRegistered represents a FillerRegistered event raised by the FillerManager contract.
type FillerManagerFillerRegistered struct {
	Filler common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFillerRegistered is a free log retrieval operation binding the contract event 0x3d5e4b9cd5a0b3c03402df7b2eb1b594bed42e75163fc8d3f08ff1d1751bfdfd.
//
// Solidity: event FillerRegistered(address indexed filler)
func (_FillerManager *FillerManagerFilterer) FilterFillerRegistered(opts *bind.FilterOpts, filler []common.Address) (*FillerManagerFillerRegisteredIterator, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _FillerManager.contract.FilterLogs(opts, "FillerRegistered", fillerRule)
	if err != nil {
		return nil, err
	}
	return &FillerManagerFillerRegisteredIterator{contract: _FillerManager.contract, event: "FillerRegistered", logs: logs, sub: sub}, nil
}

// WatchFillerRegistered is a free log subscription operation binding the contract event 0x3d5e4b9cd5a0b3c03402df7b2eb1b594bed42e75163fc8d3f08ff1d1751bfdfd.
//
// Solidity: event FillerRegistered(address indexed filler)
func (_FillerManager *FillerManagerFilterer) WatchFillerRegistered(opts *bind.WatchOpts, sink chan<- *FillerManagerFillerRegistered, filler []common.Address) (event.Subscription, error) {

	var fillerRule []interface{}
	for _, fillerItem := range filler {
		fillerRule = append(fillerRule, fillerItem)
	}

	logs, sub, err := _FillerManager.contract.WatchLogs(opts, "FillerRegistered", fillerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FillerManagerFillerRegistered)
				if err := _FillerManager.contract.UnpackLog(event, "FillerRegistered", log); err != nil {
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
func (_FillerManager *FillerManagerFilterer) ParseFillerRegistered(log types.Log) (*FillerManagerFillerRegistered, error) {
	event := new(FillerManagerFillerRegistered)
	if err := _FillerManager.contract.UnpackLog(event, "FillerRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FillerManagerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the FillerManager contract.
type FillerManagerRoleAdminChangedIterator struct {
	Event *FillerManagerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *FillerManagerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FillerManagerRoleAdminChanged)
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
		it.Event = new(FillerManagerRoleAdminChanged)
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
func (it *FillerManagerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FillerManagerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FillerManagerRoleAdminChanged represents a RoleAdminChanged event raised by the FillerManager contract.
type FillerManagerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_FillerManager *FillerManagerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*FillerManagerRoleAdminChangedIterator, error) {

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

	logs, sub, err := _FillerManager.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &FillerManagerRoleAdminChangedIterator{contract: _FillerManager.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_FillerManager *FillerManagerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *FillerManagerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _FillerManager.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FillerManagerRoleAdminChanged)
				if err := _FillerManager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_FillerManager *FillerManagerFilterer) ParseRoleAdminChanged(log types.Log) (*FillerManagerRoleAdminChanged, error) {
	event := new(FillerManagerRoleAdminChanged)
	if err := _FillerManager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FillerManagerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the FillerManager contract.
type FillerManagerRoleGrantedIterator struct {
	Event *FillerManagerRoleGranted // Event containing the contract specifics and raw log

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
func (it *FillerManagerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FillerManagerRoleGranted)
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
		it.Event = new(FillerManagerRoleGranted)
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
func (it *FillerManagerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FillerManagerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FillerManagerRoleGranted represents a RoleGranted event raised by the FillerManager contract.
type FillerManagerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_FillerManager *FillerManagerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*FillerManagerRoleGrantedIterator, error) {

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

	logs, sub, err := _FillerManager.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &FillerManagerRoleGrantedIterator{contract: _FillerManager.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_FillerManager *FillerManagerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *FillerManagerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _FillerManager.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FillerManagerRoleGranted)
				if err := _FillerManager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_FillerManager *FillerManagerFilterer) ParseRoleGranted(log types.Log) (*FillerManagerRoleGranted, error) {
	event := new(FillerManagerRoleGranted)
	if err := _FillerManager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FillerManagerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the FillerManager contract.
type FillerManagerRoleRevokedIterator struct {
	Event *FillerManagerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *FillerManagerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FillerManagerRoleRevoked)
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
		it.Event = new(FillerManagerRoleRevoked)
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
func (it *FillerManagerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FillerManagerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FillerManagerRoleRevoked represents a RoleRevoked event raised by the FillerManager contract.
type FillerManagerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_FillerManager *FillerManagerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*FillerManagerRoleRevokedIterator, error) {

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

	logs, sub, err := _FillerManager.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &FillerManagerRoleRevokedIterator{contract: _FillerManager.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_FillerManager *FillerManagerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *FillerManagerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _FillerManager.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FillerManagerRoleRevoked)
				if err := _FillerManager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_FillerManager *FillerManagerFilterer) ParseRoleRevoked(log types.Log) (*FillerManagerRoleRevoked, error) {
	event := new(FillerManagerRoleRevoked)
	if err := _FillerManager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
