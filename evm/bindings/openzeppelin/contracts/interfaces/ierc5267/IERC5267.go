// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ierc5267

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

// IERC5267MetaData contains all meta data concerning the IERC5267 contract.
var IERC5267MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IERC5267ABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC5267MetaData.ABI instead.
var IERC5267ABI = IERC5267MetaData.ABI

// IERC5267 is an auto generated Go binding around an Ethereum contract.
type IERC5267 struct {
	IERC5267Caller     // Read-only binding to the contract
	IERC5267Transactor // Write-only binding to the contract
	IERC5267Filterer   // Log filterer for contract events
}

// IERC5267Caller is an auto generated read-only Go binding around an Ethereum contract.
type IERC5267Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC5267Transactor is an auto generated write-only Go binding around an Ethereum contract.
type IERC5267Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC5267Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IERC5267Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC5267Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IERC5267Session struct {
	Contract     *IERC5267         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC5267CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IERC5267CallerSession struct {
	Contract *IERC5267Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// IERC5267TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IERC5267TransactorSession struct {
	Contract     *IERC5267Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IERC5267Raw is an auto generated low-level Go binding around an Ethereum contract.
type IERC5267Raw struct {
	Contract *IERC5267 // Generic contract binding to access the raw methods on
}

// IERC5267CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IERC5267CallerRaw struct {
	Contract *IERC5267Caller // Generic read-only contract binding to access the raw methods on
}

// IERC5267TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IERC5267TransactorRaw struct {
	Contract *IERC5267Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC5267 creates a new instance of IERC5267, bound to a specific deployed contract.
func NewIERC5267(address common.Address, backend bind.ContractBackend) (*IERC5267, error) {
	contract, err := bindIERC5267(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC5267{IERC5267Caller: IERC5267Caller{contract: contract}, IERC5267Transactor: IERC5267Transactor{contract: contract}, IERC5267Filterer: IERC5267Filterer{contract: contract}}, nil
}

// NewIERC5267Caller creates a new read-only instance of IERC5267, bound to a specific deployed contract.
func NewIERC5267Caller(address common.Address, caller bind.ContractCaller) (*IERC5267Caller, error) {
	contract, err := bindIERC5267(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC5267Caller{contract: contract}, nil
}

// NewIERC5267Transactor creates a new write-only instance of IERC5267, bound to a specific deployed contract.
func NewIERC5267Transactor(address common.Address, transactor bind.ContractTransactor) (*IERC5267Transactor, error) {
	contract, err := bindIERC5267(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC5267Transactor{contract: contract}, nil
}

// NewIERC5267Filterer creates a new log filterer instance of IERC5267, bound to a specific deployed contract.
func NewIERC5267Filterer(address common.Address, filterer bind.ContractFilterer) (*IERC5267Filterer, error) {
	contract, err := bindIERC5267(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC5267Filterer{contract: contract}, nil
}

// bindIERC5267 binds a generic wrapper to an already deployed contract.
func bindIERC5267(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC5267MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC5267 *IERC5267Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC5267.Contract.IERC5267Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC5267 *IERC5267Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC5267.Contract.IERC5267Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC5267 *IERC5267Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC5267.Contract.IERC5267Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC5267 *IERC5267CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC5267.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC5267 *IERC5267TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC5267.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC5267 *IERC5267TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC5267.Contract.contract.Transact(opts, method, params...)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_IERC5267 *IERC5267Caller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _IERC5267.contract.Call(opts, &out, "eip712Domain")

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
func (_IERC5267 *IERC5267Session) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _IERC5267.Contract.Eip712Domain(&_IERC5267.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_IERC5267 *IERC5267CallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _IERC5267.Contract.Eip712Domain(&_IERC5267.CallOpts)
}

// IERC5267EIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the IERC5267 contract.
type IERC5267EIP712DomainChangedIterator struct {
	Event *IERC5267EIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *IERC5267EIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC5267EIP712DomainChanged)
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
		it.Event = new(IERC5267EIP712DomainChanged)
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
func (it *IERC5267EIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC5267EIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC5267EIP712DomainChanged represents a EIP712DomainChanged event raised by the IERC5267 contract.
type IERC5267EIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_IERC5267 *IERC5267Filterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*IERC5267EIP712DomainChangedIterator, error) {

	logs, sub, err := _IERC5267.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &IERC5267EIP712DomainChangedIterator{contract: _IERC5267.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_IERC5267 *IERC5267Filterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *IERC5267EIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _IERC5267.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC5267EIP712DomainChanged)
				if err := _IERC5267.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_IERC5267 *IERC5267Filterer) ParseEIP712DomainChanged(log types.Log) (*IERC5267EIP712DomainChanged, error) {
	event := new(IERC5267EIP712DomainChanged)
	if err := _IERC5267.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
