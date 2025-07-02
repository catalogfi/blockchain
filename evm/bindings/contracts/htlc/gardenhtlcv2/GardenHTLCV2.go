// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gardenhtlcv2

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

// GardenHTLCV2MetaData contains all meta data concerning the GardenHTLCV2 contract.
var GardenHTLCV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token_\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"HTLC__DuplicateOrder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__IncorrectSecret\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__InvalidInitiatorSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__InvalidRedeemerSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__OrderFulfilled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__OrderNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__OrderNotInitiated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__SameFunderAndRedeemer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__SameInitiatorAndRedeemer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__ZeroAddressInitiator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__ZeroAddressRedeemer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__ZeroAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__ZeroTimelock\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidShortString\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"}],\"name\":\"Initiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"secret\",\"type\":\"bytes\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timelock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"initiate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timelock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"initiateOnBehalf\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timelock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"initiateWithSignature\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"instantRefund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"}],\"name\":\"instantRefundDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"orders\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"initiatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timelock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilledAt\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"secret\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// GardenHTLCV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use GardenHTLCV2MetaData.ABI instead.
var GardenHTLCV2ABI = GardenHTLCV2MetaData.ABI

// GardenHTLCV2 is an auto generated Go binding around an Ethereum contract.
type GardenHTLCV2 struct {
	GardenHTLCV2Caller     // Read-only binding to the contract
	GardenHTLCV2Transactor // Write-only binding to the contract
	GardenHTLCV2Filterer   // Log filterer for contract events
}

// GardenHTLCV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type GardenHTLCV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenHTLCV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type GardenHTLCV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenHTLCV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GardenHTLCV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenHTLCV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GardenHTLCV2Session struct {
	Contract     *GardenHTLCV2     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GardenHTLCV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GardenHTLCV2CallerSession struct {
	Contract *GardenHTLCV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// GardenHTLCV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GardenHTLCV2TransactorSession struct {
	Contract     *GardenHTLCV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// GardenHTLCV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type GardenHTLCV2Raw struct {
	Contract *GardenHTLCV2 // Generic contract binding to access the raw methods on
}

// GardenHTLCV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GardenHTLCV2CallerRaw struct {
	Contract *GardenHTLCV2Caller // Generic read-only contract binding to access the raw methods on
}

// GardenHTLCV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GardenHTLCV2TransactorRaw struct {
	Contract *GardenHTLCV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewGardenHTLCV2 creates a new instance of GardenHTLCV2, bound to a specific deployed contract.
func NewGardenHTLCV2(address common.Address, backend bind.ContractBackend) (*GardenHTLCV2, error) {
	contract, err := bindGardenHTLCV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV2{GardenHTLCV2Caller: GardenHTLCV2Caller{contract: contract}, GardenHTLCV2Transactor: GardenHTLCV2Transactor{contract: contract}, GardenHTLCV2Filterer: GardenHTLCV2Filterer{contract: contract}}, nil
}

// NewGardenHTLCV2Caller creates a new read-only instance of GardenHTLCV2, bound to a specific deployed contract.
func NewGardenHTLCV2Caller(address common.Address, caller bind.ContractCaller) (*GardenHTLCV2Caller, error) {
	contract, err := bindGardenHTLCV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV2Caller{contract: contract}, nil
}

// NewGardenHTLCV2Transactor creates a new write-only instance of GardenHTLCV2, bound to a specific deployed contract.
func NewGardenHTLCV2Transactor(address common.Address, transactor bind.ContractTransactor) (*GardenHTLCV2Transactor, error) {
	contract, err := bindGardenHTLCV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV2Transactor{contract: contract}, nil
}

// NewGardenHTLCV2Filterer creates a new log filterer instance of GardenHTLCV2, bound to a specific deployed contract.
func NewGardenHTLCV2Filterer(address common.Address, filterer bind.ContractFilterer) (*GardenHTLCV2Filterer, error) {
	contract, err := bindGardenHTLCV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV2Filterer{contract: contract}, nil
}

// bindGardenHTLCV2 binds a generic wrapper to an already deployed contract.
func bindGardenHTLCV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GardenHTLCV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GardenHTLCV2 *GardenHTLCV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GardenHTLCV2.Contract.GardenHTLCV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GardenHTLCV2 *GardenHTLCV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.GardenHTLCV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GardenHTLCV2 *GardenHTLCV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.GardenHTLCV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GardenHTLCV2 *GardenHTLCV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GardenHTLCV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GardenHTLCV2 *GardenHTLCV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GardenHTLCV2 *GardenHTLCV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.contract.Transact(opts, method, params...)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_GardenHTLCV2 *GardenHTLCV2Caller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _GardenHTLCV2.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_GardenHTLCV2 *GardenHTLCV2Session) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _GardenHTLCV2.Contract.Eip712Domain(&_GardenHTLCV2.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_GardenHTLCV2 *GardenHTLCV2CallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _GardenHTLCV2.Contract.Eip712Domain(&_GardenHTLCV2.CallOpts)
}

// InstantRefundDigest is a free data retrieval call binding the contract method 0x4882a380.
//
// Solidity: function instantRefundDigest(bytes32 orderID) view returns(bytes32)
func (_GardenHTLCV2 *GardenHTLCV2Caller) InstantRefundDigest(opts *bind.CallOpts, orderID [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _GardenHTLCV2.contract.Call(opts, &out, "instantRefundDigest", orderID)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// InstantRefundDigest is a free data retrieval call binding the contract method 0x4882a380.
//
// Solidity: function instantRefundDigest(bytes32 orderID) view returns(bytes32)
func (_GardenHTLCV2 *GardenHTLCV2Session) InstantRefundDigest(orderID [32]byte) ([32]byte, error) {
	return _GardenHTLCV2.Contract.InstantRefundDigest(&_GardenHTLCV2.CallOpts, orderID)
}

// InstantRefundDigest is a free data retrieval call binding the contract method 0x4882a380.
//
// Solidity: function instantRefundDigest(bytes32 orderID) view returns(bytes32)
func (_GardenHTLCV2 *GardenHTLCV2CallerSession) InstantRefundDigest(orderID [32]byte) ([32]byte, error) {
	return _GardenHTLCV2.Contract.InstantRefundDigest(&_GardenHTLCV2.CallOpts, orderID)
}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders(bytes32 ) view returns(address initiator, address redeemer, uint256 initiatedAt, uint256 timelock, uint256 amount, uint256 fulfilledAt)
func (_GardenHTLCV2 *GardenHTLCV2Caller) Orders(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Initiator   common.Address
	Redeemer    common.Address
	InitiatedAt *big.Int
	Timelock    *big.Int
	Amount      *big.Int
	FulfilledAt *big.Int
}, error) {
	var out []interface{}
	err := _GardenHTLCV2.contract.Call(opts, &out, "orders", arg0)

	outstruct := new(struct {
		Initiator   common.Address
		Redeemer    common.Address
		InitiatedAt *big.Int
		Timelock    *big.Int
		Amount      *big.Int
		FulfilledAt *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Initiator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Redeemer = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.InitiatedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Timelock = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Amount = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.FulfilledAt = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders(bytes32 ) view returns(address initiator, address redeemer, uint256 initiatedAt, uint256 timelock, uint256 amount, uint256 fulfilledAt)
func (_GardenHTLCV2 *GardenHTLCV2Session) Orders(arg0 [32]byte) (struct {
	Initiator   common.Address
	Redeemer    common.Address
	InitiatedAt *big.Int
	Timelock    *big.Int
	Amount      *big.Int
	FulfilledAt *big.Int
}, error) {
	return _GardenHTLCV2.Contract.Orders(&_GardenHTLCV2.CallOpts, arg0)
}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders(bytes32 ) view returns(address initiator, address redeemer, uint256 initiatedAt, uint256 timelock, uint256 amount, uint256 fulfilledAt)
func (_GardenHTLCV2 *GardenHTLCV2CallerSession) Orders(arg0 [32]byte) (struct {
	Initiator   common.Address
	Redeemer    common.Address
	InitiatedAt *big.Int
	Timelock    *big.Int
	Amount      *big.Int
	FulfilledAt *big.Int
}, error) {
	return _GardenHTLCV2.Contract.Orders(&_GardenHTLCV2.CallOpts, arg0)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenHTLCV2 *GardenHTLCV2Caller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GardenHTLCV2.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenHTLCV2 *GardenHTLCV2Session) Token() (common.Address, error) {
	return _GardenHTLCV2.Contract.Token(&_GardenHTLCV2.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenHTLCV2 *GardenHTLCV2CallerSession) Token() (common.Address, error) {
	return _GardenHTLCV2.Contract.Token(&_GardenHTLCV2.CallOpts)
}

// Initiate is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLCV2 *GardenHTLCV2Transactor) Initiate(opts *bind.TransactOpts, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV2.contract.Transact(opts, "initiate", redeemer, timelock, amount, secretHash)
}

// Initiate is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLCV2 *GardenHTLCV2Session) Initiate(redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.Initiate(&_GardenHTLCV2.TransactOpts, redeemer, timelock, amount, secretHash)
}

// Initiate is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLCV2 *GardenHTLCV2TransactorSession) Initiate(redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.Initiate(&_GardenHTLCV2.TransactOpts, redeemer, timelock, amount, secretHash)
}

// InitiateOnBehalf is a paid mutator transaction binding the contract method 0x13d4a787.
//
// Solidity: function initiateOnBehalf(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLCV2 *GardenHTLCV2Transactor) InitiateOnBehalf(opts *bind.TransactOpts, initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV2.contract.Transact(opts, "initiateOnBehalf", initiator, redeemer, timelock, amount, secretHash)
}

// InitiateOnBehalf is a paid mutator transaction binding the contract method 0x13d4a787.
//
// Solidity: function initiateOnBehalf(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLCV2 *GardenHTLCV2Session) InitiateOnBehalf(initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.InitiateOnBehalf(&_GardenHTLCV2.TransactOpts, initiator, redeemer, timelock, amount, secretHash)
}

// InitiateOnBehalf is a paid mutator transaction binding the contract method 0x13d4a787.
//
// Solidity: function initiateOnBehalf(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLCV2 *GardenHTLCV2TransactorSession) InitiateOnBehalf(initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.InitiateOnBehalf(&_GardenHTLCV2.TransactOpts, initiator, redeemer, timelock, amount, secretHash)
}

// InitiateWithSignature is a paid mutator transaction binding the contract method 0xd4705e9e.
//
// Solidity: function initiateWithSignature(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash, bytes signature) returns()
func (_GardenHTLCV2 *GardenHTLCV2Transactor) InitiateWithSignature(opts *bind.TransactOpts, initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLCV2.contract.Transact(opts, "initiateWithSignature", initiator, redeemer, timelock, amount, secretHash, signature)
}

// InitiateWithSignature is a paid mutator transaction binding the contract method 0xd4705e9e.
//
// Solidity: function initiateWithSignature(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash, bytes signature) returns()
func (_GardenHTLCV2 *GardenHTLCV2Session) InitiateWithSignature(initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.InitiateWithSignature(&_GardenHTLCV2.TransactOpts, initiator, redeemer, timelock, amount, secretHash, signature)
}

// InitiateWithSignature is a paid mutator transaction binding the contract method 0xd4705e9e.
//
// Solidity: function initiateWithSignature(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash, bytes signature) returns()
func (_GardenHTLCV2 *GardenHTLCV2TransactorSession) InitiateWithSignature(initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.InitiateWithSignature(&_GardenHTLCV2.TransactOpts, initiator, redeemer, timelock, amount, secretHash, signature)
}

// InstantRefund is a paid mutator transaction binding the contract method 0xedaf5fac.
//
// Solidity: function instantRefund(bytes32 orderID, bytes signature) returns()
func (_GardenHTLCV2 *GardenHTLCV2Transactor) InstantRefund(opts *bind.TransactOpts, orderID [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLCV2.contract.Transact(opts, "instantRefund", orderID, signature)
}

// InstantRefund is a paid mutator transaction binding the contract method 0xedaf5fac.
//
// Solidity: function instantRefund(bytes32 orderID, bytes signature) returns()
func (_GardenHTLCV2 *GardenHTLCV2Session) InstantRefund(orderID [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.InstantRefund(&_GardenHTLCV2.TransactOpts, orderID, signature)
}

// InstantRefund is a paid mutator transaction binding the contract method 0xedaf5fac.
//
// Solidity: function instantRefund(bytes32 orderID, bytes signature) returns()
func (_GardenHTLCV2 *GardenHTLCV2TransactorSession) InstantRefund(orderID [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.InstantRefund(&_GardenHTLCV2.TransactOpts, orderID, signature)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 orderID, bytes secret) returns()
func (_GardenHTLCV2 *GardenHTLCV2Transactor) Redeem(opts *bind.TransactOpts, orderID [32]byte, secret []byte) (*types.Transaction, error) {
	return _GardenHTLCV2.contract.Transact(opts, "redeem", orderID, secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 orderID, bytes secret) returns()
func (_GardenHTLCV2 *GardenHTLCV2Session) Redeem(orderID [32]byte, secret []byte) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.Redeem(&_GardenHTLCV2.TransactOpts, orderID, secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 orderID, bytes secret) returns()
func (_GardenHTLCV2 *GardenHTLCV2TransactorSession) Redeem(orderID [32]byte, secret []byte) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.Redeem(&_GardenHTLCV2.TransactOpts, orderID, secret)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 orderID) returns()
func (_GardenHTLCV2 *GardenHTLCV2Transactor) Refund(opts *bind.TransactOpts, orderID [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV2.contract.Transact(opts, "refund", orderID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 orderID) returns()
func (_GardenHTLCV2 *GardenHTLCV2Session) Refund(orderID [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.Refund(&_GardenHTLCV2.TransactOpts, orderID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 orderID) returns()
func (_GardenHTLCV2 *GardenHTLCV2TransactorSession) Refund(orderID [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV2.Contract.Refund(&_GardenHTLCV2.TransactOpts, orderID)
}

// GardenHTLCV2EIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the GardenHTLCV2 contract.
type GardenHTLCV2EIP712DomainChangedIterator struct {
	Event *GardenHTLCV2EIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *GardenHTLCV2EIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenHTLCV2EIP712DomainChanged)
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
		it.Event = new(GardenHTLCV2EIP712DomainChanged)
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
func (it *GardenHTLCV2EIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenHTLCV2EIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenHTLCV2EIP712DomainChanged represents a EIP712DomainChanged event raised by the GardenHTLCV2 contract.
type GardenHTLCV2EIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_GardenHTLCV2 *GardenHTLCV2Filterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*GardenHTLCV2EIP712DomainChangedIterator, error) {

	logs, sub, err := _GardenHTLCV2.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV2EIP712DomainChangedIterator{contract: _GardenHTLCV2.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_GardenHTLCV2 *GardenHTLCV2Filterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *GardenHTLCV2EIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _GardenHTLCV2.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenHTLCV2EIP712DomainChanged)
				if err := _GardenHTLCV2.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_GardenHTLCV2 *GardenHTLCV2Filterer) ParseEIP712DomainChanged(log types.Log) (*GardenHTLCV2EIP712DomainChanged, error) {
	event := new(GardenHTLCV2EIP712DomainChanged)
	if err := _GardenHTLCV2.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenHTLCV2InitiatedIterator is returned from FilterInitiated and is used to iterate over the raw logs and unpacked data for Initiated events raised by the GardenHTLCV2 contract.
type GardenHTLCV2InitiatedIterator struct {
	Event *GardenHTLCV2Initiated // Event containing the contract specifics and raw log

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
func (it *GardenHTLCV2InitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenHTLCV2Initiated)
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
		it.Event = new(GardenHTLCV2Initiated)
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
func (it *GardenHTLCV2InitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenHTLCV2InitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenHTLCV2Initiated represents a Initiated event raised by the GardenHTLCV2 contract.
type GardenHTLCV2Initiated struct {
	OrderID [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitiated is a free log retrieval operation binding the contract event 0x08007a3b331cd9bd3d1d3667a3724ba04d1b2799b75845215f1944debbdf844f.
//
// Solidity: event Initiated(bytes32 indexed orderID)
func (_GardenHTLCV2 *GardenHTLCV2Filterer) FilterInitiated(opts *bind.FilterOpts, orderID [][32]byte) (*GardenHTLCV2InitiatedIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}

	logs, sub, err := _GardenHTLCV2.contract.FilterLogs(opts, "Initiated", orderIDRule)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV2InitiatedIterator{contract: _GardenHTLCV2.contract, event: "Initiated", logs: logs, sub: sub}, nil
}

// WatchInitiated is a free log subscription operation binding the contract event 0x08007a3b331cd9bd3d1d3667a3724ba04d1b2799b75845215f1944debbdf844f.
//
// Solidity: event Initiated(bytes32 indexed orderID)
func (_GardenHTLCV2 *GardenHTLCV2Filterer) WatchInitiated(opts *bind.WatchOpts, sink chan<- *GardenHTLCV2Initiated, orderID [][32]byte) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}

	logs, sub, err := _GardenHTLCV2.contract.WatchLogs(opts, "Initiated", orderIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenHTLCV2Initiated)
				if err := _GardenHTLCV2.contract.UnpackLog(event, "Initiated", log); err != nil {
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

// ParseInitiated is a log parse operation binding the contract event 0x08007a3b331cd9bd3d1d3667a3724ba04d1b2799b75845215f1944debbdf844f.
//
// Solidity: event Initiated(bytes32 indexed orderID)
func (_GardenHTLCV2 *GardenHTLCV2Filterer) ParseInitiated(log types.Log) (*GardenHTLCV2Initiated, error) {
	event := new(GardenHTLCV2Initiated)
	if err := _GardenHTLCV2.contract.UnpackLog(event, "Initiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenHTLCV2RedeemedIterator is returned from FilterRedeemed and is used to iterate over the raw logs and unpacked data for Redeemed events raised by the GardenHTLCV2 contract.
type GardenHTLCV2RedeemedIterator struct {
	Event *GardenHTLCV2Redeemed // Event containing the contract specifics and raw log

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
func (it *GardenHTLCV2RedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenHTLCV2Redeemed)
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
		it.Event = new(GardenHTLCV2Redeemed)
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
func (it *GardenHTLCV2RedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenHTLCV2RedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenHTLCV2Redeemed represents a Redeemed event raised by the GardenHTLCV2 contract.
type GardenHTLCV2Redeemed struct {
	OrderID [32]byte
	Secret  []byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x866c33f43c7dda3105124ae616b2a42ff25811f48048edbb4ab215c59563b1e6.
//
// Solidity: event Redeemed(bytes32 indexed orderID, bytes secret)
func (_GardenHTLCV2 *GardenHTLCV2Filterer) FilterRedeemed(opts *bind.FilterOpts, orderID [][32]byte) (*GardenHTLCV2RedeemedIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}

	logs, sub, err := _GardenHTLCV2.contract.FilterLogs(opts, "Redeemed", orderIDRule)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV2RedeemedIterator{contract: _GardenHTLCV2.contract, event: "Redeemed", logs: logs, sub: sub}, nil
}

// WatchRedeemed is a free log subscription operation binding the contract event 0x866c33f43c7dda3105124ae616b2a42ff25811f48048edbb4ab215c59563b1e6.
//
// Solidity: event Redeemed(bytes32 indexed orderID, bytes secret)
func (_GardenHTLCV2 *GardenHTLCV2Filterer) WatchRedeemed(opts *bind.WatchOpts, sink chan<- *GardenHTLCV2Redeemed, orderID [][32]byte) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}

	logs, sub, err := _GardenHTLCV2.contract.WatchLogs(opts, "Redeemed", orderIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenHTLCV2Redeemed)
				if err := _GardenHTLCV2.contract.UnpackLog(event, "Redeemed", log); err != nil {
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

// ParseRedeemed is a log parse operation binding the contract event 0x866c33f43c7dda3105124ae616b2a42ff25811f48048edbb4ab215c59563b1e6.
//
// Solidity: event Redeemed(bytes32 indexed orderID, bytes secret)
func (_GardenHTLCV2 *GardenHTLCV2Filterer) ParseRedeemed(log types.Log) (*GardenHTLCV2Redeemed, error) {
	event := new(GardenHTLCV2Redeemed)
	if err := _GardenHTLCV2.contract.UnpackLog(event, "Redeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenHTLCV2RefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the GardenHTLCV2 contract.
type GardenHTLCV2RefundedIterator struct {
	Event *GardenHTLCV2Refunded // Event containing the contract specifics and raw log

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
func (it *GardenHTLCV2RefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenHTLCV2Refunded)
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
		it.Event = new(GardenHTLCV2Refunded)
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
func (it *GardenHTLCV2RefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenHTLCV2RefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenHTLCV2Refunded represents a Refunded event raised by the GardenHTLCV2 contract.
type GardenHTLCV2Refunded struct {
	OrderID [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed orderID)
func (_GardenHTLCV2 *GardenHTLCV2Filterer) FilterRefunded(opts *bind.FilterOpts, orderID [][32]byte) (*GardenHTLCV2RefundedIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}

	logs, sub, err := _GardenHTLCV2.contract.FilterLogs(opts, "Refunded", orderIDRule)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV2RefundedIterator{contract: _GardenHTLCV2.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed orderID)
func (_GardenHTLCV2 *GardenHTLCV2Filterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *GardenHTLCV2Refunded, orderID [][32]byte) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}

	logs, sub, err := _GardenHTLCV2.contract.WatchLogs(opts, "Refunded", orderIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenHTLCV2Refunded)
				if err := _GardenHTLCV2.contract.UnpackLog(event, "Refunded", log); err != nil {
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

// ParseRefunded is a log parse operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed orderID)
func (_GardenHTLCV2 *GardenHTLCV2Filterer) ParseRefunded(log types.Log) (*GardenHTLCV2Refunded, error) {
	event := new(GardenHTLCV2Refunded)
	if err := _GardenHTLCV2.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
