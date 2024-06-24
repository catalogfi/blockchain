// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eip712upgradeable

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

// EIP712UpgradeableMetaData contains all meta data concerning the EIP712Upgradeable contract.
var EIP712UpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// EIP712UpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use EIP712UpgradeableMetaData.ABI instead.
var EIP712UpgradeableABI = EIP712UpgradeableMetaData.ABI

// EIP712Upgradeable is an auto generated Go binding around an Ethereum contract.
type EIP712Upgradeable struct {
	EIP712UpgradeableCaller     // Read-only binding to the contract
	EIP712UpgradeableTransactor // Write-only binding to the contract
	EIP712UpgradeableFilterer   // Log filterer for contract events
}

// EIP712UpgradeableCaller is an auto generated read-only Go binding around an Ethereum contract.
type EIP712UpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP712UpgradeableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EIP712UpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP712UpgradeableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EIP712UpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP712UpgradeableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EIP712UpgradeableSession struct {
	Contract     *EIP712Upgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// EIP712UpgradeableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EIP712UpgradeableCallerSession struct {
	Contract *EIP712UpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// EIP712UpgradeableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EIP712UpgradeableTransactorSession struct {
	Contract     *EIP712UpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// EIP712UpgradeableRaw is an auto generated low-level Go binding around an Ethereum contract.
type EIP712UpgradeableRaw struct {
	Contract *EIP712Upgradeable // Generic contract binding to access the raw methods on
}

// EIP712UpgradeableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EIP712UpgradeableCallerRaw struct {
	Contract *EIP712UpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// EIP712UpgradeableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EIP712UpgradeableTransactorRaw struct {
	Contract *EIP712UpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEIP712Upgradeable creates a new instance of EIP712Upgradeable, bound to a specific deployed contract.
func NewEIP712Upgradeable(address common.Address, backend bind.ContractBackend) (*EIP712Upgradeable, error) {
	contract, err := bindEIP712Upgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EIP712Upgradeable{EIP712UpgradeableCaller: EIP712UpgradeableCaller{contract: contract}, EIP712UpgradeableTransactor: EIP712UpgradeableTransactor{contract: contract}, EIP712UpgradeableFilterer: EIP712UpgradeableFilterer{contract: contract}}, nil
}

// NewEIP712UpgradeableCaller creates a new read-only instance of EIP712Upgradeable, bound to a specific deployed contract.
func NewEIP712UpgradeableCaller(address common.Address, caller bind.ContractCaller) (*EIP712UpgradeableCaller, error) {
	contract, err := bindEIP712Upgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EIP712UpgradeableCaller{contract: contract}, nil
}

// NewEIP712UpgradeableTransactor creates a new write-only instance of EIP712Upgradeable, bound to a specific deployed contract.
func NewEIP712UpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*EIP712UpgradeableTransactor, error) {
	contract, err := bindEIP712Upgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EIP712UpgradeableTransactor{contract: contract}, nil
}

// NewEIP712UpgradeableFilterer creates a new log filterer instance of EIP712Upgradeable, bound to a specific deployed contract.
func NewEIP712UpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*EIP712UpgradeableFilterer, error) {
	contract, err := bindEIP712Upgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EIP712UpgradeableFilterer{contract: contract}, nil
}

// bindEIP712Upgradeable binds a generic wrapper to an already deployed contract.
func bindEIP712Upgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EIP712UpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EIP712Upgradeable *EIP712UpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EIP712Upgradeable.Contract.EIP712UpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EIP712Upgradeable *EIP712UpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP712Upgradeable.Contract.EIP712UpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EIP712Upgradeable *EIP712UpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EIP712Upgradeable.Contract.EIP712UpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EIP712Upgradeable *EIP712UpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EIP712Upgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EIP712Upgradeable *EIP712UpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP712Upgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EIP712Upgradeable *EIP712UpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EIP712Upgradeable.Contract.contract.Transact(opts, method, params...)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_EIP712Upgradeable *EIP712UpgradeableCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _EIP712Upgradeable.contract.Call(opts, &out, "eip712Domain")

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
func (_EIP712Upgradeable *EIP712UpgradeableSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _EIP712Upgradeable.Contract.Eip712Domain(&_EIP712Upgradeable.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_EIP712Upgradeable *EIP712UpgradeableCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _EIP712Upgradeable.Contract.Eip712Domain(&_EIP712Upgradeable.CallOpts)
}

// EIP712UpgradeableEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the EIP712Upgradeable contract.
type EIP712UpgradeableEIP712DomainChangedIterator struct {
	Event *EIP712UpgradeableEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *EIP712UpgradeableEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EIP712UpgradeableEIP712DomainChanged)
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
		it.Event = new(EIP712UpgradeableEIP712DomainChanged)
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
func (it *EIP712UpgradeableEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EIP712UpgradeableEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EIP712UpgradeableEIP712DomainChanged represents a EIP712DomainChanged event raised by the EIP712Upgradeable contract.
type EIP712UpgradeableEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_EIP712Upgradeable *EIP712UpgradeableFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*EIP712UpgradeableEIP712DomainChangedIterator, error) {

	logs, sub, err := _EIP712Upgradeable.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &EIP712UpgradeableEIP712DomainChangedIterator{contract: _EIP712Upgradeable.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_EIP712Upgradeable *EIP712UpgradeableFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *EIP712UpgradeableEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _EIP712Upgradeable.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EIP712UpgradeableEIP712DomainChanged)
				if err := _EIP712Upgradeable.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_EIP712Upgradeable *EIP712UpgradeableFilterer) ParseEIP712DomainChanged(log types.Log) (*EIP712UpgradeableEIP712DomainChanged, error) {
	event := new(EIP712UpgradeableEIP712DomainChanged)
	if err := _EIP712Upgradeable.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EIP712UpgradeableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the EIP712Upgradeable contract.
type EIP712UpgradeableInitializedIterator struct {
	Event *EIP712UpgradeableInitialized // Event containing the contract specifics and raw log

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
func (it *EIP712UpgradeableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EIP712UpgradeableInitialized)
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
		it.Event = new(EIP712UpgradeableInitialized)
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
func (it *EIP712UpgradeableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EIP712UpgradeableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EIP712UpgradeableInitialized represents a Initialized event raised by the EIP712Upgradeable contract.
type EIP712UpgradeableInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_EIP712Upgradeable *EIP712UpgradeableFilterer) FilterInitialized(opts *bind.FilterOpts) (*EIP712UpgradeableInitializedIterator, error) {

	logs, sub, err := _EIP712Upgradeable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &EIP712UpgradeableInitializedIterator{contract: _EIP712Upgradeable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_EIP712Upgradeable *EIP712UpgradeableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *EIP712UpgradeableInitialized) (event.Subscription, error) {

	logs, sub, err := _EIP712Upgradeable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EIP712UpgradeableInitialized)
				if err := _EIP712Upgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_EIP712Upgradeable *EIP712UpgradeableFilterer) ParseInitialized(log types.Log) (*EIP712UpgradeableInitialized, error) {
	event := new(EIP712UpgradeableInitialized)
	if err := _EIP712Upgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
