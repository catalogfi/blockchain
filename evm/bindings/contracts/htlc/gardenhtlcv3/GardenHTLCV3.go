// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gardenhtlcv3

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

// GardenHTLCV3MetaData contains all meta data concerning the GardenHTLCV3 contract.
var GardenHTLCV3MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"HTLC__DuplicateOrder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__HTLCAlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__IncorrectSecret\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__InvalidInitiatorSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__InvalidRedeemerSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__OrderFulfilled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__OrderNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__OrderNotInitiated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__SameFunderAndRedeemer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__SameInitiatorAndRedeemer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__ZeroAddressInitiator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__ZeroAddressRedeemer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__ZeroAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HTLC__ZeroTimelock\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidShortString\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Initiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"destinationData\",\"type\":\"bytes\"}],\"name\":\"InitiatedWithDestinationData\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"secret\",\"type\":\"bytes\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"initialise\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timelock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"destinationData\",\"type\":\"bytes\"}],\"name\":\"initiate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timelock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"initiate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timelock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"initiateOnBehalf\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timelock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"destinationData\",\"type\":\"bytes\"}],\"name\":\"initiateOnBehalf\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timelock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"initiateWithSignature\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"instantRefund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"}],\"name\":\"instantRefundDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isInitialized\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"orders\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"initiatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timelock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilledAt\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"secret\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// GardenHTLCV3ABI is the input ABI used to generate the binding from.
// Deprecated: Use GardenHTLCV3MetaData.ABI instead.
var GardenHTLCV3ABI = GardenHTLCV3MetaData.ABI

// GardenHTLCV3 is an auto generated Go binding around an Ethereum contract.
type GardenHTLCV3 struct {
	GardenHTLCV3Caller     // Read-only binding to the contract
	GardenHTLCV3Transactor // Write-only binding to the contract
	GardenHTLCV3Filterer   // Log filterer for contract events
}

// GardenHTLCV3Caller is an auto generated read-only Go binding around an Ethereum contract.
type GardenHTLCV3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenHTLCV3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type GardenHTLCV3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenHTLCV3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GardenHTLCV3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GardenHTLCV3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GardenHTLCV3Session struct {
	Contract     *GardenHTLCV3     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GardenHTLCV3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GardenHTLCV3CallerSession struct {
	Contract *GardenHTLCV3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// GardenHTLCV3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GardenHTLCV3TransactorSession struct {
	Contract     *GardenHTLCV3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// GardenHTLCV3Raw is an auto generated low-level Go binding around an Ethereum contract.
type GardenHTLCV3Raw struct {
	Contract *GardenHTLCV3 // Generic contract binding to access the raw methods on
}

// GardenHTLCV3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GardenHTLCV3CallerRaw struct {
	Contract *GardenHTLCV3Caller // Generic read-only contract binding to access the raw methods on
}

// GardenHTLCV3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GardenHTLCV3TransactorRaw struct {
	Contract *GardenHTLCV3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewGardenHTLCV3 creates a new instance of GardenHTLCV3, bound to a specific deployed contract.
func NewGardenHTLCV3(address common.Address, backend bind.ContractBackend) (*GardenHTLCV3, error) {
	contract, err := bindGardenHTLCV3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV3{GardenHTLCV3Caller: GardenHTLCV3Caller{contract: contract}, GardenHTLCV3Transactor: GardenHTLCV3Transactor{contract: contract}, GardenHTLCV3Filterer: GardenHTLCV3Filterer{contract: contract}}, nil
}

// NewGardenHTLCV3Caller creates a new read-only instance of GardenHTLCV3, bound to a specific deployed contract.
func NewGardenHTLCV3Caller(address common.Address, caller bind.ContractCaller) (*GardenHTLCV3Caller, error) {
	contract, err := bindGardenHTLCV3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV3Caller{contract: contract}, nil
}

// NewGardenHTLCV3Transactor creates a new write-only instance of GardenHTLCV3, bound to a specific deployed contract.
func NewGardenHTLCV3Transactor(address common.Address, transactor bind.ContractTransactor) (*GardenHTLCV3Transactor, error) {
	contract, err := bindGardenHTLCV3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV3Transactor{contract: contract}, nil
}

// NewGardenHTLCV3Filterer creates a new log filterer instance of GardenHTLCV3, bound to a specific deployed contract.
func NewGardenHTLCV3Filterer(address common.Address, filterer bind.ContractFilterer) (*GardenHTLCV3Filterer, error) {
	contract, err := bindGardenHTLCV3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV3Filterer{contract: contract}, nil
}

// bindGardenHTLCV3 binds a generic wrapper to an already deployed contract.
func bindGardenHTLCV3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GardenHTLCV3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GardenHTLCV3 *GardenHTLCV3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GardenHTLCV3.Contract.GardenHTLCV3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GardenHTLCV3 *GardenHTLCV3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.GardenHTLCV3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GardenHTLCV3 *GardenHTLCV3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.GardenHTLCV3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GardenHTLCV3 *GardenHTLCV3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GardenHTLCV3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GardenHTLCV3 *GardenHTLCV3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GardenHTLCV3 *GardenHTLCV3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.contract.Transact(opts, method, params...)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_GardenHTLCV3 *GardenHTLCV3Caller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _GardenHTLCV3.contract.Call(opts, &out, "eip712Domain")

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
func (_GardenHTLCV3 *GardenHTLCV3Session) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _GardenHTLCV3.Contract.Eip712Domain(&_GardenHTLCV3.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_GardenHTLCV3 *GardenHTLCV3CallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _GardenHTLCV3.Contract.Eip712Domain(&_GardenHTLCV3.CallOpts)
}

// InstantRefundDigest is a free data retrieval call binding the contract method 0x4882a380.
//
// Solidity: function instantRefundDigest(bytes32 orderID) view returns(bytes32)
func (_GardenHTLCV3 *GardenHTLCV3Caller) InstantRefundDigest(opts *bind.CallOpts, orderID [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _GardenHTLCV3.contract.Call(opts, &out, "instantRefundDigest", orderID)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// InstantRefundDigest is a free data retrieval call binding the contract method 0x4882a380.
//
// Solidity: function instantRefundDigest(bytes32 orderID) view returns(bytes32)
func (_GardenHTLCV3 *GardenHTLCV3Session) InstantRefundDigest(orderID [32]byte) ([32]byte, error) {
	return _GardenHTLCV3.Contract.InstantRefundDigest(&_GardenHTLCV3.CallOpts, orderID)
}

// InstantRefundDigest is a free data retrieval call binding the contract method 0x4882a380.
//
// Solidity: function instantRefundDigest(bytes32 orderID) view returns(bytes32)
func (_GardenHTLCV3 *GardenHTLCV3CallerSession) InstantRefundDigest(orderID [32]byte) ([32]byte, error) {
	return _GardenHTLCV3.Contract.InstantRefundDigest(&_GardenHTLCV3.CallOpts, orderID)
}

// IsInitialized is a free data retrieval call binding the contract method 0x392e53cd.
//
// Solidity: function isInitialized() view returns(uint256)
func (_GardenHTLCV3 *GardenHTLCV3Caller) IsInitialized(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GardenHTLCV3.contract.Call(opts, &out, "isInitialized")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IsInitialized is a free data retrieval call binding the contract method 0x392e53cd.
//
// Solidity: function isInitialized() view returns(uint256)
func (_GardenHTLCV3 *GardenHTLCV3Session) IsInitialized() (*big.Int, error) {
	return _GardenHTLCV3.Contract.IsInitialized(&_GardenHTLCV3.CallOpts)
}

// IsInitialized is a free data retrieval call binding the contract method 0x392e53cd.
//
// Solidity: function isInitialized() view returns(uint256)
func (_GardenHTLCV3 *GardenHTLCV3CallerSession) IsInitialized() (*big.Int, error) {
	return _GardenHTLCV3.Contract.IsInitialized(&_GardenHTLCV3.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GardenHTLCV3 *GardenHTLCV3Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GardenHTLCV3.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GardenHTLCV3 *GardenHTLCV3Session) Name() (string, error) {
	return _GardenHTLCV3.Contract.Name(&_GardenHTLCV3.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GardenHTLCV3 *GardenHTLCV3CallerSession) Name() (string, error) {
	return _GardenHTLCV3.Contract.Name(&_GardenHTLCV3.CallOpts)
}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders(bytes32 ) view returns(address initiator, address redeemer, uint256 initiatedAt, uint256 timelock, uint256 amount, uint256 fulfilledAt)
func (_GardenHTLCV3 *GardenHTLCV3Caller) Orders(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Initiator   common.Address
	Redeemer    common.Address
	InitiatedAt *big.Int
	Timelock    *big.Int
	Amount      *big.Int
	FulfilledAt *big.Int
}, error) {
	var out []interface{}
	err := _GardenHTLCV3.contract.Call(opts, &out, "orders", arg0)

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
func (_GardenHTLCV3 *GardenHTLCV3Session) Orders(arg0 [32]byte) (struct {
	Initiator   common.Address
	Redeemer    common.Address
	InitiatedAt *big.Int
	Timelock    *big.Int
	Amount      *big.Int
	FulfilledAt *big.Int
}, error) {
	return _GardenHTLCV3.Contract.Orders(&_GardenHTLCV3.CallOpts, arg0)
}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders(bytes32 ) view returns(address initiator, address redeemer, uint256 initiatedAt, uint256 timelock, uint256 amount, uint256 fulfilledAt)
func (_GardenHTLCV3 *GardenHTLCV3CallerSession) Orders(arg0 [32]byte) (struct {
	Initiator   common.Address
	Redeemer    common.Address
	InitiatedAt *big.Int
	Timelock    *big.Int
	Amount      *big.Int
	FulfilledAt *big.Int
}, error) {
	return _GardenHTLCV3.Contract.Orders(&_GardenHTLCV3.CallOpts, arg0)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenHTLCV3 *GardenHTLCV3Caller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GardenHTLCV3.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenHTLCV3 *GardenHTLCV3Session) Token() (common.Address, error) {
	return _GardenHTLCV3.Contract.Token(&_GardenHTLCV3.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GardenHTLCV3 *GardenHTLCV3CallerSession) Token() (common.Address, error) {
	return _GardenHTLCV3.Contract.Token(&_GardenHTLCV3.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_GardenHTLCV3 *GardenHTLCV3Caller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GardenHTLCV3.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_GardenHTLCV3 *GardenHTLCV3Session) Version() (string, error) {
	return _GardenHTLCV3.Contract.Version(&_GardenHTLCV3.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_GardenHTLCV3 *GardenHTLCV3CallerSession) Version() (string, error) {
	return _GardenHTLCV3.Contract.Version(&_GardenHTLCV3.CallOpts)
}

// Initialise is a paid mutator transaction binding the contract method 0x9d6a890f.
//
// Solidity: function initialise(address _token) returns()
func (_GardenHTLCV3 *GardenHTLCV3Transactor) Initialise(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _GardenHTLCV3.contract.Transact(opts, "initialise", _token)
}

// Initialise is a paid mutator transaction binding the contract method 0x9d6a890f.
//
// Solidity: function initialise(address _token) returns()
func (_GardenHTLCV3 *GardenHTLCV3Session) Initialise(_token common.Address) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.Initialise(&_GardenHTLCV3.TransactOpts, _token)
}

// Initialise is a paid mutator transaction binding the contract method 0x9d6a890f.
//
// Solidity: function initialise(address _token) returns()
func (_GardenHTLCV3 *GardenHTLCV3TransactorSession) Initialise(_token common.Address) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.Initialise(&_GardenHTLCV3.TransactOpts, _token)
}

// Initiate is a paid mutator transaction binding the contract method 0x4ede0ab7.
//
// Solidity: function initiate(address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash, bytes destinationData) returns()
func (_GardenHTLCV3 *GardenHTLCV3Transactor) Initiate(opts *bind.TransactOpts, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte, destinationData []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.contract.Transact(opts, "initiate", redeemer, timelock, amount, secretHash, destinationData)
}

// Initiate is a paid mutator transaction binding the contract method 0x4ede0ab7.
//
// Solidity: function initiate(address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash, bytes destinationData) returns()
func (_GardenHTLCV3 *GardenHTLCV3Session) Initiate(redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte, destinationData []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.Initiate(&_GardenHTLCV3.TransactOpts, redeemer, timelock, amount, secretHash, destinationData)
}

// Initiate is a paid mutator transaction binding the contract method 0x4ede0ab7.
//
// Solidity: function initiate(address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash, bytes destinationData) returns()
func (_GardenHTLCV3 *GardenHTLCV3TransactorSession) Initiate(redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte, destinationData []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.Initiate(&_GardenHTLCV3.TransactOpts, redeemer, timelock, amount, secretHash, destinationData)
}

// Initiate0 is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLCV3 *GardenHTLCV3Transactor) Initiate0(opts *bind.TransactOpts, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV3.contract.Transact(opts, "initiate0", redeemer, timelock, amount, secretHash)
}

// Initiate0 is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLCV3 *GardenHTLCV3Session) Initiate0(redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.Initiate0(&_GardenHTLCV3.TransactOpts, redeemer, timelock, amount, secretHash)
}

// Initiate0 is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLCV3 *GardenHTLCV3TransactorSession) Initiate0(redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.Initiate0(&_GardenHTLCV3.TransactOpts, redeemer, timelock, amount, secretHash)
}

// InitiateOnBehalf is a paid mutator transaction binding the contract method 0x13d4a787.
//
// Solidity: function initiateOnBehalf(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLCV3 *GardenHTLCV3Transactor) InitiateOnBehalf(opts *bind.TransactOpts, initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV3.contract.Transact(opts, "initiateOnBehalf", initiator, redeemer, timelock, amount, secretHash)
}

// InitiateOnBehalf is a paid mutator transaction binding the contract method 0x13d4a787.
//
// Solidity: function initiateOnBehalf(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLCV3 *GardenHTLCV3Session) InitiateOnBehalf(initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.InitiateOnBehalf(&_GardenHTLCV3.TransactOpts, initiator, redeemer, timelock, amount, secretHash)
}

// InitiateOnBehalf is a paid mutator transaction binding the contract method 0x13d4a787.
//
// Solidity: function initiateOnBehalf(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash) returns()
func (_GardenHTLCV3 *GardenHTLCV3TransactorSession) InitiateOnBehalf(initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.InitiateOnBehalf(&_GardenHTLCV3.TransactOpts, initiator, redeemer, timelock, amount, secretHash)
}

// InitiateOnBehalf0 is a paid mutator transaction binding the contract method 0xa66a8641.
//
// Solidity: function initiateOnBehalf(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash, bytes destinationData) returns()
func (_GardenHTLCV3 *GardenHTLCV3Transactor) InitiateOnBehalf0(opts *bind.TransactOpts, initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte, destinationData []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.contract.Transact(opts, "initiateOnBehalf0", initiator, redeemer, timelock, amount, secretHash, destinationData)
}

// InitiateOnBehalf0 is a paid mutator transaction binding the contract method 0xa66a8641.
//
// Solidity: function initiateOnBehalf(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash, bytes destinationData) returns()
func (_GardenHTLCV3 *GardenHTLCV3Session) InitiateOnBehalf0(initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte, destinationData []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.InitiateOnBehalf0(&_GardenHTLCV3.TransactOpts, initiator, redeemer, timelock, amount, secretHash, destinationData)
}

// InitiateOnBehalf0 is a paid mutator transaction binding the contract method 0xa66a8641.
//
// Solidity: function initiateOnBehalf(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash, bytes destinationData) returns()
func (_GardenHTLCV3 *GardenHTLCV3TransactorSession) InitiateOnBehalf0(initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte, destinationData []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.InitiateOnBehalf0(&_GardenHTLCV3.TransactOpts, initiator, redeemer, timelock, amount, secretHash, destinationData)
}

// InitiateWithSignature is a paid mutator transaction binding the contract method 0xd4705e9e.
//
// Solidity: function initiateWithSignature(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash, bytes signature) returns()
func (_GardenHTLCV3 *GardenHTLCV3Transactor) InitiateWithSignature(opts *bind.TransactOpts, initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.contract.Transact(opts, "initiateWithSignature", initiator, redeemer, timelock, amount, secretHash, signature)
}

// InitiateWithSignature is a paid mutator transaction binding the contract method 0xd4705e9e.
//
// Solidity: function initiateWithSignature(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash, bytes signature) returns()
func (_GardenHTLCV3 *GardenHTLCV3Session) InitiateWithSignature(initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.InitiateWithSignature(&_GardenHTLCV3.TransactOpts, initiator, redeemer, timelock, amount, secretHash, signature)
}

// InitiateWithSignature is a paid mutator transaction binding the contract method 0xd4705e9e.
//
// Solidity: function initiateWithSignature(address initiator, address redeemer, uint256 timelock, uint256 amount, bytes32 secretHash, bytes signature) returns()
func (_GardenHTLCV3 *GardenHTLCV3TransactorSession) InitiateWithSignature(initiator common.Address, redeemer common.Address, timelock *big.Int, amount *big.Int, secretHash [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.InitiateWithSignature(&_GardenHTLCV3.TransactOpts, initiator, redeemer, timelock, amount, secretHash, signature)
}

// InstantRefund is a paid mutator transaction binding the contract method 0xedaf5fac.
//
// Solidity: function instantRefund(bytes32 orderID, bytes signature) returns()
func (_GardenHTLCV3 *GardenHTLCV3Transactor) InstantRefund(opts *bind.TransactOpts, orderID [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.contract.Transact(opts, "instantRefund", orderID, signature)
}

// InstantRefund is a paid mutator transaction binding the contract method 0xedaf5fac.
//
// Solidity: function instantRefund(bytes32 orderID, bytes signature) returns()
func (_GardenHTLCV3 *GardenHTLCV3Session) InstantRefund(orderID [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.InstantRefund(&_GardenHTLCV3.TransactOpts, orderID, signature)
}

// InstantRefund is a paid mutator transaction binding the contract method 0xedaf5fac.
//
// Solidity: function instantRefund(bytes32 orderID, bytes signature) returns()
func (_GardenHTLCV3 *GardenHTLCV3TransactorSession) InstantRefund(orderID [32]byte, signature []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.InstantRefund(&_GardenHTLCV3.TransactOpts, orderID, signature)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 orderID, bytes secret) returns()
func (_GardenHTLCV3 *GardenHTLCV3Transactor) Redeem(opts *bind.TransactOpts, orderID [32]byte, secret []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.contract.Transact(opts, "redeem", orderID, secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 orderID, bytes secret) returns()
func (_GardenHTLCV3 *GardenHTLCV3Session) Redeem(orderID [32]byte, secret []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.Redeem(&_GardenHTLCV3.TransactOpts, orderID, secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 orderID, bytes secret) returns()
func (_GardenHTLCV3 *GardenHTLCV3TransactorSession) Redeem(orderID [32]byte, secret []byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.Redeem(&_GardenHTLCV3.TransactOpts, orderID, secret)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 orderID) returns()
func (_GardenHTLCV3 *GardenHTLCV3Transactor) Refund(opts *bind.TransactOpts, orderID [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV3.contract.Transact(opts, "refund", orderID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 orderID) returns()
func (_GardenHTLCV3 *GardenHTLCV3Session) Refund(orderID [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.Refund(&_GardenHTLCV3.TransactOpts, orderID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 orderID) returns()
func (_GardenHTLCV3 *GardenHTLCV3TransactorSession) Refund(orderID [32]byte) (*types.Transaction, error) {
	return _GardenHTLCV3.Contract.Refund(&_GardenHTLCV3.TransactOpts, orderID)
}

// GardenHTLCV3EIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the GardenHTLCV3 contract.
type GardenHTLCV3EIP712DomainChangedIterator struct {
	Event *GardenHTLCV3EIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *GardenHTLCV3EIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenHTLCV3EIP712DomainChanged)
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
		it.Event = new(GardenHTLCV3EIP712DomainChanged)
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
func (it *GardenHTLCV3EIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenHTLCV3EIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenHTLCV3EIP712DomainChanged represents a EIP712DomainChanged event raised by the GardenHTLCV3 contract.
type GardenHTLCV3EIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_GardenHTLCV3 *GardenHTLCV3Filterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*GardenHTLCV3EIP712DomainChangedIterator, error) {

	logs, sub, err := _GardenHTLCV3.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV3EIP712DomainChangedIterator{contract: _GardenHTLCV3.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_GardenHTLCV3 *GardenHTLCV3Filterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *GardenHTLCV3EIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _GardenHTLCV3.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenHTLCV3EIP712DomainChanged)
				if err := _GardenHTLCV3.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_GardenHTLCV3 *GardenHTLCV3Filterer) ParseEIP712DomainChanged(log types.Log) (*GardenHTLCV3EIP712DomainChanged, error) {
	event := new(GardenHTLCV3EIP712DomainChanged)
	if err := _GardenHTLCV3.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenHTLCV3InitiatedIterator is returned from FilterInitiated and is used to iterate over the raw logs and unpacked data for Initiated events raised by the GardenHTLCV3 contract.
type GardenHTLCV3InitiatedIterator struct {
	Event *GardenHTLCV3Initiated // Event containing the contract specifics and raw log

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
func (it *GardenHTLCV3InitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenHTLCV3Initiated)
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
		it.Event = new(GardenHTLCV3Initiated)
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
func (it *GardenHTLCV3InitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenHTLCV3InitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenHTLCV3Initiated represents a Initiated event raised by the GardenHTLCV3 contract.
type GardenHTLCV3Initiated struct {
	OrderID    [32]byte
	SecretHash [32]byte
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterInitiated is a free log retrieval operation binding the contract event 0x01b41cbd4bbcc3c5b968a04d3fbdd8c1648a39ff6d9a3929b4840cea1142bc65.
//
// Solidity: event Initiated(bytes32 indexed orderID, bytes32 indexed secretHash, uint256 indexed amount)
func (_GardenHTLCV3 *GardenHTLCV3Filterer) FilterInitiated(opts *bind.FilterOpts, orderID [][32]byte, secretHash [][32]byte, amount []*big.Int) (*GardenHTLCV3InitiatedIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _GardenHTLCV3.contract.FilterLogs(opts, "Initiated", orderIDRule, secretHashRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV3InitiatedIterator{contract: _GardenHTLCV3.contract, event: "Initiated", logs: logs, sub: sub}, nil
}

// WatchInitiated is a free log subscription operation binding the contract event 0x01b41cbd4bbcc3c5b968a04d3fbdd8c1648a39ff6d9a3929b4840cea1142bc65.
//
// Solidity: event Initiated(bytes32 indexed orderID, bytes32 indexed secretHash, uint256 indexed amount)
func (_GardenHTLCV3 *GardenHTLCV3Filterer) WatchInitiated(opts *bind.WatchOpts, sink chan<- *GardenHTLCV3Initiated, orderID [][32]byte, secretHash [][32]byte, amount []*big.Int) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _GardenHTLCV3.contract.WatchLogs(opts, "Initiated", orderIDRule, secretHashRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenHTLCV3Initiated)
				if err := _GardenHTLCV3.contract.UnpackLog(event, "Initiated", log); err != nil {
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

// ParseInitiated is a log parse operation binding the contract event 0x01b41cbd4bbcc3c5b968a04d3fbdd8c1648a39ff6d9a3929b4840cea1142bc65.
//
// Solidity: event Initiated(bytes32 indexed orderID, bytes32 indexed secretHash, uint256 indexed amount)
func (_GardenHTLCV3 *GardenHTLCV3Filterer) ParseInitiated(log types.Log) (*GardenHTLCV3Initiated, error) {
	event := new(GardenHTLCV3Initiated)
	if err := _GardenHTLCV3.contract.UnpackLog(event, "Initiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenHTLCV3InitiatedWithDestinationDataIterator is returned from FilterInitiatedWithDestinationData and is used to iterate over the raw logs and unpacked data for InitiatedWithDestinationData events raised by the GardenHTLCV3 contract.
type GardenHTLCV3InitiatedWithDestinationDataIterator struct {
	Event *GardenHTLCV3InitiatedWithDestinationData // Event containing the contract specifics and raw log

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
func (it *GardenHTLCV3InitiatedWithDestinationDataIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenHTLCV3InitiatedWithDestinationData)
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
		it.Event = new(GardenHTLCV3InitiatedWithDestinationData)
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
func (it *GardenHTLCV3InitiatedWithDestinationDataIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenHTLCV3InitiatedWithDestinationDataIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenHTLCV3InitiatedWithDestinationData represents a InitiatedWithDestinationData event raised by the GardenHTLCV3 contract.
type GardenHTLCV3InitiatedWithDestinationData struct {
	OrderID         [32]byte
	SecretHash      [32]byte
	Amount          *big.Int
	DestinationData []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterInitiatedWithDestinationData is a free log retrieval operation binding the contract event 0x380328b26f928a9e51ec19f8e2199a2bce34db266d093d1e0104bce139ed7555.
//
// Solidity: event InitiatedWithDestinationData(bytes32 indexed orderID, bytes32 indexed secretHash, uint256 indexed amount, bytes destinationData)
func (_GardenHTLCV3 *GardenHTLCV3Filterer) FilterInitiatedWithDestinationData(opts *bind.FilterOpts, orderID [][32]byte, secretHash [][32]byte, amount []*big.Int) (*GardenHTLCV3InitiatedWithDestinationDataIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _GardenHTLCV3.contract.FilterLogs(opts, "InitiatedWithDestinationData", orderIDRule, secretHashRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV3InitiatedWithDestinationDataIterator{contract: _GardenHTLCV3.contract, event: "InitiatedWithDestinationData", logs: logs, sub: sub}, nil
}

// WatchInitiatedWithDestinationData is a free log subscription operation binding the contract event 0x380328b26f928a9e51ec19f8e2199a2bce34db266d093d1e0104bce139ed7555.
//
// Solidity: event InitiatedWithDestinationData(bytes32 indexed orderID, bytes32 indexed secretHash, uint256 indexed amount, bytes destinationData)
func (_GardenHTLCV3 *GardenHTLCV3Filterer) WatchInitiatedWithDestinationData(opts *bind.WatchOpts, sink chan<- *GardenHTLCV3InitiatedWithDestinationData, orderID [][32]byte, secretHash [][32]byte, amount []*big.Int) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _GardenHTLCV3.contract.WatchLogs(opts, "InitiatedWithDestinationData", orderIDRule, secretHashRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenHTLCV3InitiatedWithDestinationData)
				if err := _GardenHTLCV3.contract.UnpackLog(event, "InitiatedWithDestinationData", log); err != nil {
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

// ParseInitiatedWithDestinationData is a log parse operation binding the contract event 0x380328b26f928a9e51ec19f8e2199a2bce34db266d093d1e0104bce139ed7555.
//
// Solidity: event InitiatedWithDestinationData(bytes32 indexed orderID, bytes32 indexed secretHash, uint256 indexed amount, bytes destinationData)
func (_GardenHTLCV3 *GardenHTLCV3Filterer) ParseInitiatedWithDestinationData(log types.Log) (*GardenHTLCV3InitiatedWithDestinationData, error) {
	event := new(GardenHTLCV3InitiatedWithDestinationData)
	if err := _GardenHTLCV3.contract.UnpackLog(event, "InitiatedWithDestinationData", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenHTLCV3RedeemedIterator is returned from FilterRedeemed and is used to iterate over the raw logs and unpacked data for Redeemed events raised by the GardenHTLCV3 contract.
type GardenHTLCV3RedeemedIterator struct {
	Event *GardenHTLCV3Redeemed // Event containing the contract specifics and raw log

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
func (it *GardenHTLCV3RedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenHTLCV3Redeemed)
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
		it.Event = new(GardenHTLCV3Redeemed)
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
func (it *GardenHTLCV3RedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenHTLCV3RedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenHTLCV3Redeemed represents a Redeemed event raised by the GardenHTLCV3 contract.
type GardenHTLCV3Redeemed struct {
	OrderID    [32]byte
	SecretHash [32]byte
	Secret     []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x4c9a044220477b4e94dbb0d07ff6ff4ac30d443bef59098c4541b006954778e2.
//
// Solidity: event Redeemed(bytes32 indexed orderID, bytes32 indexed secretHash, bytes secret)
func (_GardenHTLCV3 *GardenHTLCV3Filterer) FilterRedeemed(opts *bind.FilterOpts, orderID [][32]byte, secretHash [][32]byte) (*GardenHTLCV3RedeemedIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}

	logs, sub, err := _GardenHTLCV3.contract.FilterLogs(opts, "Redeemed", orderIDRule, secretHashRule)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV3RedeemedIterator{contract: _GardenHTLCV3.contract, event: "Redeemed", logs: logs, sub: sub}, nil
}

// WatchRedeemed is a free log subscription operation binding the contract event 0x4c9a044220477b4e94dbb0d07ff6ff4ac30d443bef59098c4541b006954778e2.
//
// Solidity: event Redeemed(bytes32 indexed orderID, bytes32 indexed secretHash, bytes secret)
func (_GardenHTLCV3 *GardenHTLCV3Filterer) WatchRedeemed(opts *bind.WatchOpts, sink chan<- *GardenHTLCV3Redeemed, orderID [][32]byte, secretHash [][32]byte) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}

	logs, sub, err := _GardenHTLCV3.contract.WatchLogs(opts, "Redeemed", orderIDRule, secretHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenHTLCV3Redeemed)
				if err := _GardenHTLCV3.contract.UnpackLog(event, "Redeemed", log); err != nil {
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
func (_GardenHTLCV3 *GardenHTLCV3Filterer) ParseRedeemed(log types.Log) (*GardenHTLCV3Redeemed, error) {
	event := new(GardenHTLCV3Redeemed)
	if err := _GardenHTLCV3.contract.UnpackLog(event, "Redeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GardenHTLCV3RefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the GardenHTLCV3 contract.
type GardenHTLCV3RefundedIterator struct {
	Event *GardenHTLCV3Refunded // Event containing the contract specifics and raw log

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
func (it *GardenHTLCV3RefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GardenHTLCV3Refunded)
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
		it.Event = new(GardenHTLCV3Refunded)
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
func (it *GardenHTLCV3RefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GardenHTLCV3RefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GardenHTLCV3Refunded represents a Refunded event raised by the GardenHTLCV3 contract.
type GardenHTLCV3Refunded struct {
	OrderID [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed orderID)
func (_GardenHTLCV3 *GardenHTLCV3Filterer) FilterRefunded(opts *bind.FilterOpts, orderID [][32]byte) (*GardenHTLCV3RefundedIterator, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}

	logs, sub, err := _GardenHTLCV3.contract.FilterLogs(opts, "Refunded", orderIDRule)
	if err != nil {
		return nil, err
	}
	return &GardenHTLCV3RefundedIterator{contract: _GardenHTLCV3.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed orderID)
func (_GardenHTLCV3 *GardenHTLCV3Filterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *GardenHTLCV3Refunded, orderID [][32]byte) (event.Subscription, error) {

	var orderIDRule []interface{}
	for _, orderIDItem := range orderID {
		orderIDRule = append(orderIDRule, orderIDItem)
	}

	logs, sub, err := _GardenHTLCV3.contract.WatchLogs(opts, "Refunded", orderIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GardenHTLCV3Refunded)
				if err := _GardenHTLCV3.contract.UnpackLog(event, "Refunded", log); err != nil {
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
func (_GardenHTLCV3 *GardenHTLCV3Filterer) ParseRefunded(log types.Log) (*GardenHTLCV3Refunded, error) {
	event := new(GardenHTLCV3Refunded)
	if err := _GardenHTLCV3.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
