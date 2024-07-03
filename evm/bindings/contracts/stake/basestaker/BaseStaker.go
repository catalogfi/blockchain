// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package basestaker

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

// BaseStakerMetaData contains all meta data concerning the BaseStaker contract.
var BaseStakerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DELEGATE_STAKE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FILLER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FILLER_COOL_DOWN\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FILLER_STAKE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SEED\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"delegateNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"}],\"name\":\"getFiller\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"feeInBips\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deregisteredAt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"delegateStakeIDs\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"stakes\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"units\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"filler\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// BaseStakerABI is the input ABI used to generate the binding from.
// Deprecated: Use BaseStakerMetaData.ABI instead.
var BaseStakerABI = BaseStakerMetaData.ABI

// BaseStaker is an auto generated Go binding around an Ethereum contract.
type BaseStaker struct {
	BaseStakerCaller     // Read-only binding to the contract
	BaseStakerTransactor // Write-only binding to the contract
	BaseStakerFilterer   // Log filterer for contract events
}

// BaseStakerCaller is an auto generated read-only Go binding around an Ethereum contract.
type BaseStakerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseStakerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BaseStakerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseStakerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BaseStakerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseStakerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BaseStakerSession struct {
	Contract     *BaseStaker       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BaseStakerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BaseStakerCallerSession struct {
	Contract *BaseStakerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// BaseStakerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BaseStakerTransactorSession struct {
	Contract     *BaseStakerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// BaseStakerRaw is an auto generated low-level Go binding around an Ethereum contract.
type BaseStakerRaw struct {
	Contract *BaseStaker // Generic contract binding to access the raw methods on
}

// BaseStakerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BaseStakerCallerRaw struct {
	Contract *BaseStakerCaller // Generic read-only contract binding to access the raw methods on
}

// BaseStakerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BaseStakerTransactorRaw struct {
	Contract *BaseStakerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBaseStaker creates a new instance of BaseStaker, bound to a specific deployed contract.
func NewBaseStaker(address common.Address, backend bind.ContractBackend) (*BaseStaker, error) {
	contract, err := bindBaseStaker(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BaseStaker{BaseStakerCaller: BaseStakerCaller{contract: contract}, BaseStakerTransactor: BaseStakerTransactor{contract: contract}, BaseStakerFilterer: BaseStakerFilterer{contract: contract}}, nil
}

// NewBaseStakerCaller creates a new read-only instance of BaseStaker, bound to a specific deployed contract.
func NewBaseStakerCaller(address common.Address, caller bind.ContractCaller) (*BaseStakerCaller, error) {
	contract, err := bindBaseStaker(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BaseStakerCaller{contract: contract}, nil
}

// NewBaseStakerTransactor creates a new write-only instance of BaseStaker, bound to a specific deployed contract.
func NewBaseStakerTransactor(address common.Address, transactor bind.ContractTransactor) (*BaseStakerTransactor, error) {
	contract, err := bindBaseStaker(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BaseStakerTransactor{contract: contract}, nil
}

// NewBaseStakerFilterer creates a new log filterer instance of BaseStaker, bound to a specific deployed contract.
func NewBaseStakerFilterer(address common.Address, filterer bind.ContractFilterer) (*BaseStakerFilterer, error) {
	contract, err := bindBaseStaker(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BaseStakerFilterer{contract: contract}, nil
}

// bindBaseStaker binds a generic wrapper to an already deployed contract.
func bindBaseStaker(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BaseStakerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BaseStaker *BaseStakerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BaseStaker.Contract.BaseStakerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BaseStaker *BaseStakerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BaseStaker.Contract.BaseStakerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BaseStaker *BaseStakerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BaseStaker.Contract.BaseStakerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BaseStaker *BaseStakerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BaseStaker.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BaseStaker *BaseStakerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BaseStaker.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BaseStaker *BaseStakerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BaseStaker.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_BaseStaker *BaseStakerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BaseStaker.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_BaseStaker *BaseStakerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _BaseStaker.Contract.DEFAULTADMINROLE(&_BaseStaker.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_BaseStaker *BaseStakerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _BaseStaker.Contract.DEFAULTADMINROLE(&_BaseStaker.CallOpts)
}

// DELEGATESTAKE is a free data retrieval call binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() view returns(uint256)
func (_BaseStaker *BaseStakerCaller) DELEGATESTAKE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BaseStaker.contract.Call(opts, &out, "DELEGATE_STAKE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DELEGATESTAKE is a free data retrieval call binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() view returns(uint256)
func (_BaseStaker *BaseStakerSession) DELEGATESTAKE() (*big.Int, error) {
	return _BaseStaker.Contract.DELEGATESTAKE(&_BaseStaker.CallOpts)
}

// DELEGATESTAKE is a free data retrieval call binding the contract method 0x08155f3e.
//
// Solidity: function DELEGATE_STAKE() view returns(uint256)
func (_BaseStaker *BaseStakerCallerSession) DELEGATESTAKE() (*big.Int, error) {
	return _BaseStaker.Contract.DELEGATESTAKE(&_BaseStaker.CallOpts)
}

// FILLER is a free data retrieval call binding the contract method 0x64645bae.
//
// Solidity: function FILLER() view returns(bytes32)
func (_BaseStaker *BaseStakerCaller) FILLER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BaseStaker.contract.Call(opts, &out, "FILLER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FILLER is a free data retrieval call binding the contract method 0x64645bae.
//
// Solidity: function FILLER() view returns(bytes32)
func (_BaseStaker *BaseStakerSession) FILLER() ([32]byte, error) {
	return _BaseStaker.Contract.FILLER(&_BaseStaker.CallOpts)
}

// FILLER is a free data retrieval call binding the contract method 0x64645bae.
//
// Solidity: function FILLER() view returns(bytes32)
func (_BaseStaker *BaseStakerCallerSession) FILLER() ([32]byte, error) {
	return _BaseStaker.Contract.FILLER(&_BaseStaker.CallOpts)
}

// FILLERCOOLDOWN is a free data retrieval call binding the contract method 0x015a44e8.
//
// Solidity: function FILLER_COOL_DOWN() view returns(uint256)
func (_BaseStaker *BaseStakerCaller) FILLERCOOLDOWN(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BaseStaker.contract.Call(opts, &out, "FILLER_COOL_DOWN")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FILLERCOOLDOWN is a free data retrieval call binding the contract method 0x015a44e8.
//
// Solidity: function FILLER_COOL_DOWN() view returns(uint256)
func (_BaseStaker *BaseStakerSession) FILLERCOOLDOWN() (*big.Int, error) {
	return _BaseStaker.Contract.FILLERCOOLDOWN(&_BaseStaker.CallOpts)
}

// FILLERCOOLDOWN is a free data retrieval call binding the contract method 0x015a44e8.
//
// Solidity: function FILLER_COOL_DOWN() view returns(uint256)
func (_BaseStaker *BaseStakerCallerSession) FILLERCOOLDOWN() (*big.Int, error) {
	return _BaseStaker.Contract.FILLERCOOLDOWN(&_BaseStaker.CallOpts)
}

// FILLERSTAKE is a free data retrieval call binding the contract method 0xecfe923e.
//
// Solidity: function FILLER_STAKE() view returns(uint256)
func (_BaseStaker *BaseStakerCaller) FILLERSTAKE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BaseStaker.contract.Call(opts, &out, "FILLER_STAKE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FILLERSTAKE is a free data retrieval call binding the contract method 0xecfe923e.
//
// Solidity: function FILLER_STAKE() view returns(uint256)
func (_BaseStaker *BaseStakerSession) FILLERSTAKE() (*big.Int, error) {
	return _BaseStaker.Contract.FILLERSTAKE(&_BaseStaker.CallOpts)
}

// FILLERSTAKE is a free data retrieval call binding the contract method 0xecfe923e.
//
// Solidity: function FILLER_STAKE() view returns(uint256)
func (_BaseStaker *BaseStakerCallerSession) FILLERSTAKE() (*big.Int, error) {
	return _BaseStaker.Contract.FILLERSTAKE(&_BaseStaker.CallOpts)
}

// SEED is a free data retrieval call binding the contract method 0x0edc4737.
//
// Solidity: function SEED() view returns(address)
func (_BaseStaker *BaseStakerCaller) SEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BaseStaker.contract.Call(opts, &out, "SEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SEED is a free data retrieval call binding the contract method 0x0edc4737.
//
// Solidity: function SEED() view returns(address)
func (_BaseStaker *BaseStakerSession) SEED() (common.Address, error) {
	return _BaseStaker.Contract.SEED(&_BaseStaker.CallOpts)
}

// SEED is a free data retrieval call binding the contract method 0x0edc4737.
//
// Solidity: function SEED() view returns(address)
func (_BaseStaker *BaseStakerCallerSession) SEED() (common.Address, error) {
	return _BaseStaker.Contract.SEED(&_BaseStaker.CallOpts)
}

// DelegateNonce is a free data retrieval call binding the contract method 0x92e1f102.
//
// Solidity: function delegateNonce(address ) view returns(uint256)
func (_BaseStaker *BaseStakerCaller) DelegateNonce(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BaseStaker.contract.Call(opts, &out, "delegateNonce", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DelegateNonce is a free data retrieval call binding the contract method 0x92e1f102.
//
// Solidity: function delegateNonce(address ) view returns(uint256)
func (_BaseStaker *BaseStakerSession) DelegateNonce(arg0 common.Address) (*big.Int, error) {
	return _BaseStaker.Contract.DelegateNonce(&_BaseStaker.CallOpts, arg0)
}

// DelegateNonce is a free data retrieval call binding the contract method 0x92e1f102.
//
// Solidity: function delegateNonce(address ) view returns(uint256)
func (_BaseStaker *BaseStakerCallerSession) DelegateNonce(arg0 common.Address) (*big.Int, error) {
	return _BaseStaker.Contract.DelegateNonce(&_BaseStaker.CallOpts, arg0)
}

// GetFiller is a free data retrieval call binding the contract method 0x0b3ff53f.
//
// Solidity: function getFiller(address filler) view returns(uint16 feeInBips, uint256 stake, uint256 deregisteredAt, bytes32[] delegateStakeIDs)
func (_BaseStaker *BaseStakerCaller) GetFiller(opts *bind.CallOpts, filler common.Address) (struct {
	FeeInBips        uint16
	Stake            *big.Int
	DeregisteredAt   *big.Int
	DelegateStakeIDs [][32]byte
}, error) {
	var out []interface{}
	err := _BaseStaker.contract.Call(opts, &out, "getFiller", filler)

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
func (_BaseStaker *BaseStakerSession) GetFiller(filler common.Address) (struct {
	FeeInBips        uint16
	Stake            *big.Int
	DeregisteredAt   *big.Int
	DelegateStakeIDs [][32]byte
}, error) {
	return _BaseStaker.Contract.GetFiller(&_BaseStaker.CallOpts, filler)
}

// GetFiller is a free data retrieval call binding the contract method 0x0b3ff53f.
//
// Solidity: function getFiller(address filler) view returns(uint16 feeInBips, uint256 stake, uint256 deregisteredAt, bytes32[] delegateStakeIDs)
func (_BaseStaker *BaseStakerCallerSession) GetFiller(filler common.Address) (struct {
	FeeInBips        uint16
	Stake            *big.Int
	DeregisteredAt   *big.Int
	DelegateStakeIDs [][32]byte
}, error) {
	return _BaseStaker.Contract.GetFiller(&_BaseStaker.CallOpts, filler)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_BaseStaker *BaseStakerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _BaseStaker.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_BaseStaker *BaseStakerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _BaseStaker.Contract.GetRoleAdmin(&_BaseStaker.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_BaseStaker *BaseStakerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _BaseStaker.Contract.GetRoleAdmin(&_BaseStaker.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BaseStaker *BaseStakerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _BaseStaker.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BaseStaker *BaseStakerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _BaseStaker.Contract.HasRole(&_BaseStaker.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BaseStaker *BaseStakerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _BaseStaker.Contract.HasRole(&_BaseStaker.CallOpts, role, account)
}

// Stakes is a free data retrieval call binding the contract method 0x8fee6407.
//
// Solidity: function stakes(bytes32 ) view returns(address owner, uint256 stake, uint256 units, uint256 votes, address filler, uint256 expiry)
func (_BaseStaker *BaseStakerCaller) Stakes(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Owner  common.Address
	Stake  *big.Int
	Units  *big.Int
	Votes  *big.Int
	Filler common.Address
	Expiry *big.Int
}, error) {
	var out []interface{}
	err := _BaseStaker.contract.Call(opts, &out, "stakes", arg0)

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
func (_BaseStaker *BaseStakerSession) Stakes(arg0 [32]byte) (struct {
	Owner  common.Address
	Stake  *big.Int
	Units  *big.Int
	Votes  *big.Int
	Filler common.Address
	Expiry *big.Int
}, error) {
	return _BaseStaker.Contract.Stakes(&_BaseStaker.CallOpts, arg0)
}

// Stakes is a free data retrieval call binding the contract method 0x8fee6407.
//
// Solidity: function stakes(bytes32 ) view returns(address owner, uint256 stake, uint256 units, uint256 votes, address filler, uint256 expiry)
func (_BaseStaker *BaseStakerCallerSession) Stakes(arg0 [32]byte) (struct {
	Owner  common.Address
	Stake  *big.Int
	Units  *big.Int
	Votes  *big.Int
	Filler common.Address
	Expiry *big.Int
}, error) {
	return _BaseStaker.Contract.Stakes(&_BaseStaker.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BaseStaker *BaseStakerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BaseStaker.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BaseStaker *BaseStakerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BaseStaker.Contract.SupportsInterface(&_BaseStaker.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BaseStaker *BaseStakerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BaseStaker.Contract.SupportsInterface(&_BaseStaker.CallOpts, interfaceId)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_BaseStaker *BaseStakerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BaseStaker.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_BaseStaker *BaseStakerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BaseStaker.Contract.GrantRole(&_BaseStaker.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_BaseStaker *BaseStakerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BaseStaker.Contract.GrantRole(&_BaseStaker.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_BaseStaker *BaseStakerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BaseStaker.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_BaseStaker *BaseStakerSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BaseStaker.Contract.RenounceRole(&_BaseStaker.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_BaseStaker *BaseStakerTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BaseStaker.Contract.RenounceRole(&_BaseStaker.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_BaseStaker *BaseStakerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BaseStaker.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_BaseStaker *BaseStakerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BaseStaker.Contract.RevokeRole(&_BaseStaker.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_BaseStaker *BaseStakerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BaseStaker.Contract.RevokeRole(&_BaseStaker.TransactOpts, role, account)
}

// BaseStakerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the BaseStaker contract.
type BaseStakerRoleAdminChangedIterator struct {
	Event *BaseStakerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *BaseStakerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseStakerRoleAdminChanged)
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
		it.Event = new(BaseStakerRoleAdminChanged)
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
func (it *BaseStakerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseStakerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseStakerRoleAdminChanged represents a RoleAdminChanged event raised by the BaseStaker contract.
type BaseStakerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_BaseStaker *BaseStakerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*BaseStakerRoleAdminChangedIterator, error) {

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

	logs, sub, err := _BaseStaker.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &BaseStakerRoleAdminChangedIterator{contract: _BaseStaker.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_BaseStaker *BaseStakerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *BaseStakerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _BaseStaker.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseStakerRoleAdminChanged)
				if err := _BaseStaker.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_BaseStaker *BaseStakerFilterer) ParseRoleAdminChanged(log types.Log) (*BaseStakerRoleAdminChanged, error) {
	event := new(BaseStakerRoleAdminChanged)
	if err := _BaseStaker.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BaseStakerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the BaseStaker contract.
type BaseStakerRoleGrantedIterator struct {
	Event *BaseStakerRoleGranted // Event containing the contract specifics and raw log

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
func (it *BaseStakerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseStakerRoleGranted)
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
		it.Event = new(BaseStakerRoleGranted)
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
func (it *BaseStakerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseStakerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseStakerRoleGranted represents a RoleGranted event raised by the BaseStaker contract.
type BaseStakerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_BaseStaker *BaseStakerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BaseStakerRoleGrantedIterator, error) {

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

	logs, sub, err := _BaseStaker.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BaseStakerRoleGrantedIterator{contract: _BaseStaker.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_BaseStaker *BaseStakerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *BaseStakerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _BaseStaker.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseStakerRoleGranted)
				if err := _BaseStaker.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_BaseStaker *BaseStakerFilterer) ParseRoleGranted(log types.Log) (*BaseStakerRoleGranted, error) {
	event := new(BaseStakerRoleGranted)
	if err := _BaseStaker.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BaseStakerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the BaseStaker contract.
type BaseStakerRoleRevokedIterator struct {
	Event *BaseStakerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *BaseStakerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseStakerRoleRevoked)
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
		it.Event = new(BaseStakerRoleRevoked)
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
func (it *BaseStakerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseStakerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseStakerRoleRevoked represents a RoleRevoked event raised by the BaseStaker contract.
type BaseStakerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_BaseStaker *BaseStakerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BaseStakerRoleRevokedIterator, error) {

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

	logs, sub, err := _BaseStaker.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BaseStakerRoleRevokedIterator{contract: _BaseStaker.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_BaseStaker *BaseStakerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *BaseStakerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _BaseStaker.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseStakerRoleRevoked)
				if err := _BaseStaker.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_BaseStaker *BaseStakerFilterer) ParseRoleRevoked(log types.Log) (*BaseStakerRoleRevoked, error) {
	event := new(BaseStakerRoleRevoked)
	if err := _BaseStaker.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
