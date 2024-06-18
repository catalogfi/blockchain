// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gardenhtlc

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

// GardenHTLCMetaData contains all meta data concerning the GardenHTLC contract.
var GardenHTLCMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token_\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidShortString\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Initiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"secret\",\"type\":\"bytes\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"initiate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"initiateWithSignature\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"instantRefund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"orders\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isFulfilled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"initiatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"secret\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// GardenHTLCABI is the input ABI used to generate the binding from.
// Deprecated: Use GardenHTLCMetaData.ABI instead.
var GardenHTLCABI = GardenHTLCMetaData.ABI

// GardenHTLC is an auto generated Go binding around an Ethereum contract.
type GardenHTLC struct {
	GardenHTLCCaller     // Read-only binding to the contract
	GardenHTLCTransactor // Write-only binding to the contract
	GardenHTLCFilterer   // Log filterer for contract events
}

// GardenHTLCCaller is an auto generated read-only Go binding around an Ethereum contract.
type GardenHTLCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenHTLCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GardenHTLCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenHTLCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GardenHTLCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenHTLCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GardenHTLCSession struct {
	Contract     *GardenHTLC       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GardenHTLCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GardenHTLCCallerSession struct {
	Contract *GardenHTLCCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// GardenHTLCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GardenHTLCTransactorSession struct {
	Contract     *GardenHTLCTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// GardenHTLCRaw is an auto generated low-level Go binding around an Ethereum contract.
type GardenHTLCRaw struct {
	Contract *GardenHTLC // Generic contract binding to access the raw methods on
}

// GardenHTLCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GardenHTLCCallerRaw struct {
	Contract *GardenHTLCCaller // Generic read-only contract binding to access the raw methods on
}

// GardenHTLCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GardenHTLCTransactorRaw struct {
	Contract *GardenHTLCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGardenHTLC creates a new instance of GardenHTLC, bound to a specific deployed contract.
func NewGardenHTLC(address common.Address, backend bind.ContractBackend) (*GardenHTLC, error) {
	contract, err := bindGardenHTLC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GardenHTLC{GardenHTLCCaller: GardenHTLCCaller{contract: contract}, GardenHTLCTransactor: GardenHTLCTransactor{contract: contract}, GardenHTLCFilterer: GardenHTLCFilterer{contract: contract}}, nil
}

// NewGardenHTLCCaller creates a new read-only instance of GardenHTLC, bound to a specific deployed contract.
func NewGardenHTLCCaller(address common.Address, caller bind.ContractCaller) (*GardenHTLCCaller, error) {
	contract, err := bindGardenHTLC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCCaller{contract: contract}, nil
}

// NewGardenHTLCTransactor creates a new write-only instance of GardenHTLC, bound to a specific deployed contract.
func NewGardenHTLCTransactor(address common.Address, transactor bind.ContractTransactor) (*GardenHTLCTransactor, error) {
	contract, err := bindGardenHTLC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCTransactor{contract: contract}, nil
}

// NewGardenHTLCFilterer creates a new log filterer instance of GardenHTLC, bound to a specific deployed contract.
func NewGardenHTLCFilterer(address common.Address, filterer bind.ContractFilterer) (*GardenHTLCFilterer, error) {
	contract, err := bindGardenHTLC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCFilterer{contract: contract}, nil
}

// bindGardenHTLC binds a generic wrapper to an already deployed contract.
func bindGardenHTLC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GardenHTLCMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GardenHTLC *GardenHTLCRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GardenHTLC.Contract.GardenHTLCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GardenHTLC *GardenHTLCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenHTLC.Contract.GardenHTLCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GardenHTLC *GardenHTLCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GardenHTLC.Contract.GardenHTLCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GardenHTLC *GardenHTLCCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GardenHTLC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GardenHTLC *GardenHTLCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenHTLC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GardenHTLC *GardenHTLCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GardenHTLC.Contract.contract.Transact(opts, method, params...)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_GardenHTLC *GardenHTLCCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _GardenHTLC.contract.Call(opts, &out, "eip712Domain")

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
func (_GardenHTLC *GardenHTLCSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _GardenHTLC.Contract.Eip712Domain(&_GardenHTLC.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_GardenHTLC *GardenHTLCCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _GardenHTLC.Contract.Eip712Domain(&_GardenHTLC.CallOpts)
}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders(bytes32 ) view returns(bool isFulfilled, address initiator, address redeemer, uint256 initiatedAt, uint256 expiry, uint256 amount)
func (_GardenHTLC *GardenHTLCCaller) Orders(opts *bind.CallOpts, arg0 [32]byte) (struct {
	IsFulfilled bool
	Initiator   common.Address
	Redeemer    common.Address
	InitiatedAt *big.Int
	Expiry      *big.Int
	Amount      *big.Int
}, error) {
	var out []interface{}
	err := _GardenHTLC.contract.Call(opts, &out, "orders", arg0)

	outstruct := new(struct {
		IsFulfilled bool
		Initiator   common.Address
		Redeemer    common.Address
		InitiatedAt *big.Int
		Expiry      *big.Int
		Amount      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsFulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Initiator = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Redeemer = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.InitiatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Expiry = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Amount = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders(bytes32 ) view returns(bool isFulfilled, address initiator, address redeemer, uint256 initiatedAt, uint256 expiry, uint256 amount)
func (_GardenHTLC *GardenHTLCSession) Orders(arg0 [32]byte) (struct {
	IsFulfilled bool
	Initiator   common.Address
	Redeemer    common.Address
	InitiatedAt *big.Int
	Expiry      *big.Int
	Amount      *big.Int
}, error) {
	return _GardenHTLC.Contract.Orders(&_GardenHTLC.CallOpts, arg0)
}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders(bytes32 ) view returns(bool isFulfilled, address initiator, address redeemer, uint256 initiatedAt, uint256 expiry, uint256 amount)
func (_GardenHTLC *GardenHTLCCallerSession) Orders(arg0 [32]byte) (struct {
	IsFulfilled bool
	Initiator   common.Address
	Redeemer    common.Address
	InitiatedAt *big.Int
	Expiry      *big.Int
	Amount      *big.Int
}, error) {
	return _GardenHTLC.Contract.Orders(&_GardenHTLC.CallOpts, arg0)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenHTLC *GardenHTLCCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GardenHTLC.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenHTLC *GardenHTLCSession) Token() (common.Address, error) {
	return _GardenHTLC.Contract.Token(&_GardenHTLC.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenHTLC *GardenHTLCCallerSession) Token() (common.Address, error) {
	return _GardenHTLC.Contract.Token(&_GardenHTLC.CallOpts)
}

// Initiate is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address redeemer, uint256 expiry, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLC *GardenHTLCTransactor) Initiate(opts *bind.TransactOpts, redeemer common.Address, expiry *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLC.contract.Transact(opts, "initiate", redeemer, expiry, amount, secretHash)
}

// Initiate is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address redeemer, uint256 expiry, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLC *GardenHTLCSession) Initiate(redeemer common.Address, expiry *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLC.Contract.Initiate(&_GardenHTLC.TransactOpts, redeemer, expiry, amount, secretHash)
}

// Initiate is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address redeemer, uint256 expiry, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLC *GardenHTLCTransactorSession) Initiate(redeemer common.Address, expiry *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLC.Contract.Initiate(&_GardenHTLC.TransactOpts, redeemer, expiry, amount, secretHash)
}

// InitiateWithSignature is a paid mutator transaction binding the contract method 0x7929d59d.
//
// Solidity: function initiateWithSignature(address redeemer, uint256 expiry, uint256 amount, bytes32 secretHash, bytes signature) returns()
func (_GardenHTLC *GardenHTLCTransactor) InitiateWithSignature(opts *bind.TransactOpts, redeemer common.Address, expiry *big.Int, amount *big.Int, secretHash [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLC.contract.Transact(opts, "initiateWithSignature", redeemer, expiry, amount, secretHash, signature)
}

// InitiateWithSignature is a paid mutator transaction binding the contract method 0x7929d59d.
//
// Solidity: function initiateWithSignature(address redeemer, uint256 expiry, uint256 amount, bytes32 secretHash, bytes signature) returns()
func (_GardenHTLC *GardenHTLCSession) InitiateWithSignature(redeemer common.Address, expiry *big.Int, amount *big.Int, secretHash [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLC.Contract.InitiateWithSignature(&_GardenHTLC.TransactOpts, redeemer, expiry, amount, secretHash, signature)
}

// InitiateWithSignature is a paid mutator transaction binding the contract method 0x7929d59d.
//
// Solidity: function initiateWithSignature(address redeemer, uint256 expiry, uint256 amount, bytes32 secretHash, bytes signature) returns()
func (_GardenHTLC *GardenHTLCTransactorSession) InitiateWithSignature(redeemer common.Address, expiry *big.Int, amount *big.Int, secretHash [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLC.Contract.InitiateWithSignature(&_GardenHTLC.TransactOpts, redeemer, expiry, amount, secretHash, signature)
}

// InstantRefund is a paid mutator transaction binding the contract method 0xedaf5fac.
//
// Solidity: function instantRefund(bytes32 orderID, bytes signature) returns()
func (_GardenHTLC *GardenHTLCTransactor) InstantRefund(opts *bind.TransactOpts, orderID [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLC.contract.Transact(opts, "instantRefund", orderID, signature)
}

// InstantRefund is a paid mutator transaction binding the contract method 0xedaf5fac.
//
// Solidity: function instantRefund(bytes32 orderID, bytes signature) returns()
func (_GardenHTLC *GardenHTLCSession) InstantRefund(orderID [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLC.Contract.InstantRefund(&_GardenHTLC.TransactOpts, orderID, signature)
}

// InstantRefund is a paid mutator transaction binding the contract method 0xedaf5fac.
//
// Solidity: function instantRefund(bytes32 orderID, bytes signature) returns()
func (_GardenHTLC *GardenHTLCTransactorSession) InstantRefund(orderID [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLC.Contract.InstantRefund(&_GardenHTLC.TransactOpts, orderID, signature)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 orderID, bytes secret) returns()
func (_GardenHTLC *GardenHTLCTransactor) Redeem(opts *bind.TransactOpts, orderID [32]byte, secret []byte) (*types.Transaction, error) {
	return _GardenHTLC.contract.Transact(opts, "redeem", orderID, secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 orderID, bytes secret) returns()
func (_GardenHTLC *GardenHTLCSession) Redeem(orderID [32]byte, secret []byte) (*types.Transaction, error) {
	return _GardenHTLC.Contract.Redeem(&_GardenHTLC.TransactOpts, orderID, secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 orderID, bytes secret) returns()
func (_GardenHTLC *GardenHTLCTransactorSession) Redeem(orderID [32]byte, secret []byte) (*types.Transaction, error) {
	return _GardenHTLC.Contract.Redeem(&_GardenHTLC.TransactOpts, orderID, secret)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 orderID) returns()
func (_GardenHTLC *GardenHTLCTransactor) Refund(opts *bind.TransactOpts, orderID [32]byte) (*types.Transaction, error) {
	return _GardenHTLC.contract.Transact(opts, "refund", orderID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 orderID) returns()
func (_GardenHTLC *GardenHTLCSession) Refund(orderID [32]byte) (*types.Transaction, error) {
	return _GardenHTLC.Contract.Refund(&_GardenHTLC.TransactOpts, orderID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 orderID) returns()
func (_GardenHTLC *GardenHTLCTransactorSession) Refund(orderID [32]byte) (*types.Transaction, error) {
	return _GardenHTLC.Contract.Refund(&_GardenHTLC.TransactOpts, orderID)
}

// GardenHTLCEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the GardenHTLC contract.
type GardenHTLCEIP712DomainChangedIterator struct {
	Event *GardenHTLCEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *GardenHTLCEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenHTLCEIP712DomainChanged)
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
		it.Event = new(GardenHTLCEIP712DomainChanged)
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
func (it *GardenHTLCEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenHTLCEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenHTLCEIP712DomainChanged represents a EIP712DomainChanged event raised by the GardenHTLC contract.
type GardenHTLCEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_GardenHTLC *GardenHTLCFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*GardenHTLCEIP712DomainChangedIterator, error) {

	logs, sub, err := _GardenHTLC.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &GardenHTLCEIP712DomainChangedIterator{contract: _GardenHTLC.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_GardenHTLC *GardenHTLCFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *GardenHTLCEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _GardenHTLC.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenHTLCEIP712DomainChanged)
				if err := _GardenHTLC.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_GardenHTLC *GardenHTLCFilterer) ParseEIP712DomainChanged(log types.Log) (*GardenHTLCEIP712DomainChanged, error) {
	event := new(GardenHTLCEIP712DomainChanged)
	if err := _GardenHTLC.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenHTLCInitiatedIterator is returned from FilterInitiated and is used to iterate over the raw logs and unpacked data for Initiated events raised by the GardenHTLC contract.
type GardenHTLCInitiatedIterator struct {
	Event *GardenHTLCInitiated // Event containing the contract specifics and raw log

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
func (it *GardenHTLCInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenHTLCInitiated)
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
		it.Event = new(GardenHTLCInitiated)
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
func (it *GardenHTLCInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenHTLCInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenHTLCInitiated represents a Initiated event raised by the GardenHTLC contract.
type GardenHTLCInitiated struct {
	OrderID    [32]byte
	SecretHash [32]byte
	Initiator  common.Address
	Redeemer   common.Address
	Expiry     *big.Int
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterInitiated is a free log retrieval operation binding the contract event 0x339a507b68879fc593b831ada44efd433300b5d1121d89f6602d8dbe245ce4c3.
//
// Solidity: event Initiated(bytes32 indexed orderID, bytes32 indexed secretHash, address initiator, address redeemer, uint256 expiry, uint256 amount)
func (_GardenHTLC *GardenHTLCFilterer) FilterInitiated(opts *bind.FilterOpts, orderID [][32]byte, secretHash [][32]byte) (*GardenHTLCInitiatedIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}

	logs, sub, err := _GardenHTLC.contract.FilterLogs(opts, "Initiated", orderIDRule, secretHashRule)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCInitiatedIterator{contract: _GardenHTLC.contract, event: "Initiated", logs: logs, sub: sub}, nil
}

// WatchInitiated is a free log subscription operation binding the contract event 0x339a507b68879fc593b831ada44efd433300b5d1121d89f6602d8dbe245ce4c3.
//
// Solidity: event Initiated(bytes32 indexed orderID, bytes32 indexed secretHash, address initiator, address redeemer, uint256 expiry, uint256 amount)
func (_GardenHTLC *GardenHTLCFilterer) WatchInitiated(opts *bind.WatchOpts, sink chan<- *GardenHTLCInitiated, orderID [][32]byte, secretHash [][32]byte) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}

	logs, sub, err := _GardenHTLC.contract.WatchLogs(opts, "Initiated", orderIDRule, secretHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenHTLCInitiated)
				if err := _GardenHTLC.contract.UnpackLog(event, "Initiated", log); err != nil {
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

// ParseInitiated is a log parse operation binding the contract event 0x339a507b68879fc593b831ada44efd433300b5d1121d89f6602d8dbe245ce4c3.
//
// Solidity: event Initiated(bytes32 indexed orderID, bytes32 indexed secretHash, address initiator, address redeemer, uint256 expiry, uint256 amount)
func (_GardenHTLC *GardenHTLCFilterer) ParseInitiated(log types.Log) (*GardenHTLCInitiated, error) {
	event := new(GardenHTLCInitiated)
	if err := _GardenHTLC.contract.UnpackLog(event, "Initiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenHTLCRedeemedIterator is returned from FilterRedeemed and is used to iterate over the raw logs and unpacked data for Redeemed events raised by the GardenHTLC contract.
type GardenHTLCRedeemedIterator struct {
	Event *GardenHTLCRedeemed // Event containing the contract specifics and raw log

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
func (it *GardenHTLCRedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenHTLCRedeemed)
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
		it.Event = new(GardenHTLCRedeemed)
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
func (it *GardenHTLCRedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenHTLCRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenHTLCRedeemed represents a Redeemed event raised by the GardenHTLC contract.
type GardenHTLCRedeemed struct {
	OrderID    [32]byte
	SecretHash [32]byte
	Secret     []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x4c9a044220477b4e94dbb0d07ff6ff4ac30d443bef59098c4541b006954778e2.
//
// Solidity: event Redeemed(bytes32 indexed orderID, bytes32 indexed secretHash, bytes secret)
func (_GardenHTLC *GardenHTLCFilterer) FilterRedeemed(opts *bind.FilterOpts, orderID [][32]byte, secretHash [][32]byte) (*GardenHTLCRedeemedIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}

	logs, sub, err := _GardenHTLC.contract.FilterLogs(opts, "Redeemed", orderIDRule, secretHashRule)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCRedeemedIterator{contract: _GardenHTLC.contract, event: "Redeemed", logs: logs, sub: sub}, nil
}

// WatchRedeemed is a free log subscription operation binding the contract event 0x4c9a044220477b4e94dbb0d07ff6ff4ac30d443bef59098c4541b006954778e2.
//
// Solidity: event Redeemed(bytes32 indexed orderID, bytes32 indexed secretHash, bytes secret)
func (_GardenHTLC *GardenHTLCFilterer) WatchRedeemed(opts *bind.WatchOpts, sink chan<- *GardenHTLCRedeemed, orderID [][32]byte, secretHash [][32]byte) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}

	logs, sub, err := _GardenHTLC.contract.WatchLogs(opts, "Redeemed", orderIDRule, secretHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenHTLCRedeemed)
				if err := _GardenHTLC.contract.UnpackLog(event, "Redeemed", log); err != nil {
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

// ParseRedeemed is a log parse operation binding the contract event 0x4c9a044220477b4e94dbb0d07ff6ff4ac30d443bef59098c4541b006954778e2.
//
// Solidity: event Redeemed(bytes32 indexed orderID, bytes32 indexed secretHash, bytes secret)
func (_GardenHTLC *GardenHTLCFilterer) ParseRedeemed(log types.Log) (*GardenHTLCRedeemed, error) {
	event := new(GardenHTLCRedeemed)
	if err := _GardenHTLC.contract.UnpackLog(event, "Redeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenHTLCRefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the GardenHTLC contract.
type GardenHTLCRefundedIterator struct {
	Event *GardenHTLCRefunded // Event containing the contract specifics and raw log

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
func (it *GardenHTLCRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenHTLCRefunded)
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
		it.Event = new(GardenHTLCRefunded)
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
func (it *GardenHTLCRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenHTLCRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenHTLCRefunded represents a Refunded event raised by the GardenHTLC contract.
type GardenHTLCRefunded struct {
	OrderID [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed orderID)
func (_GardenHTLC *GardenHTLCFilterer) FilterRefunded(opts *bind.FilterOpts, orderID [][32]byte) (*GardenHTLCRefundedIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}

	logs, sub, err := _GardenHTLC.contract.FilterLogs(opts, "Refunded", orderIDRule)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCRefundedIterator{contract: _GardenHTLC.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed orderID)
func (_GardenHTLC *GardenHTLCFilterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *GardenHTLCRefunded, orderID [][32]byte) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}

	logs, sub, err := _GardenHTLC.contract.WatchLogs(opts, "Refunded", orderIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenHTLCRefunded)
				if err := _GardenHTLC.contract.UnpackLog(event, "Refunded", log); err != nil {
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
func (_GardenHTLC *GardenHTLCFilterer) ParseRefunded(log types.Log) (*GardenHTLCRefunded, error) {
	event := new(GardenHTLCRefunded)
	if err := _GardenHTLC.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
