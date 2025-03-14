// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package orderbook

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

// OrderbookMetaData contains all meta data concerning the Orderbook contract.
var OrderbookMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"createOrder\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"fillOrder\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Created\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Filled\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
	Bin: "0x6080604052348015600f57600080fd5b506101aa8061001f6000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063bc2e8e5f1461003b578063d15b385814610050575b600080fd5b61004e6100493660046100d1565b610063565b005b61004e61005e3660046100d1565b6100a0565b7f3b4a808d54c495ced21b3f493818c8c5a90f41e600c8b12208571402dcd23cec8282604051610094929190610145565b60405180910390a15050565b7f4bf70fd1dfd385f5ae9a107d7fff4995db1197c3bcb51f665f8a667df32aabb28282604051610094929190610145565b600080602083850312156100e457600080fd5b823567ffffffffffffffff8111156100fb57600080fd5b8301601f8101851361010c57600080fd5b803567ffffffffffffffff81111561012357600080fd5b85602082840101111561013557600080fd5b6020919091019590945092505050565b60208152816020820152818360408301376000818301604090810191909152601f909201601f1916010191905056fea26469706673582212204e5f1c8d9607250ca385d7a6481a800c74d709c723241fb068779fdb9ad85b9964736f6c634300081a0033",
}

// OrderbookABI is the input ABI used to generate the binding from.
// Deprecated: Use OrderbookMetaData.ABI instead.
var OrderbookABI = OrderbookMetaData.ABI

// OrderbookBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OrderbookMetaData.Bin instead.
var OrderbookBin = OrderbookMetaData.Bin

// DeployOrderbook deploys a new Ethereum contract, binding an instance of Orderbook to it.
func DeployOrderbook(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Orderbook, error) {
	parsed, err := OrderbookMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OrderbookBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Orderbook{OrderbookCaller: OrderbookCaller{contract: contract}, OrderbookTransactor: OrderbookTransactor{contract: contract}, OrderbookFilterer: OrderbookFilterer{contract: contract}}, nil
}

// Orderbook is an auto generated Go binding around an Ethereum contract.
type Orderbook struct {
	OrderbookCaller     // Read-only binding to the contract
	OrderbookTransactor // Write-only binding to the contract
	OrderbookFilterer   // Log filterer for contract events
}

// OrderbookCaller is an auto generated read-only Go binding around an Ethereum contract.
type OrderbookCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrderbookTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OrderbookTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrderbookFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OrderbookFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrderbookSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OrderbookSession struct {
	Contract     *Orderbook        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OrderbookCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OrderbookCallerSession struct {
	Contract *OrderbookCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// OrderbookTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OrderbookTransactorSession struct {
	Contract     *OrderbookTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// OrderbookRaw is an auto generated low-level Go binding around an Ethereum contract.
type OrderbookRaw struct {
	Contract *Orderbook // Generic contract binding to access the raw methods on
}

// OrderbookCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OrderbookCallerRaw struct {
	Contract *OrderbookCaller // Generic read-only contract binding to access the raw methods on
}

// OrderbookTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OrderbookTransactorRaw struct {
	Contract *OrderbookTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOrderbook creates a new instance of Orderbook, bound to a specific deployed contract.
func NewOrderbook(address common.Address, backend bind.ContractBackend) (*Orderbook, error) {
	contract, err := bindOrderbook(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Orderbook{OrderbookCaller: OrderbookCaller{contract: contract}, OrderbookTransactor: OrderbookTransactor{contract: contract}, OrderbookFilterer: OrderbookFilterer{contract: contract}}, nil
}

// NewOrderbookCaller creates a new read-only instance of Orderbook, bound to a specific deployed contract.
func NewOrderbookCaller(address common.Address, caller bind.ContractCaller) (*OrderbookCaller, error) {
	contract, err := bindOrderbook(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OrderbookCaller{contract: contract}, nil
}

// NewOrderbookTransactor creates a new write-only instance of Orderbook, bound to a specific deployed contract.
func NewOrderbookTransactor(address common.Address, transactor bind.ContractTransactor) (*OrderbookTransactor, error) {
	contract, err := bindOrderbook(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OrderbookTransactor{contract: contract}, nil
}

// NewOrderbookFilterer creates a new log filterer instance of Orderbook, bound to a specific deployed contract.
func NewOrderbookFilterer(address common.Address, filterer bind.ContractFilterer) (*OrderbookFilterer, error) {
	contract, err := bindOrderbook(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OrderbookFilterer{contract: contract}, nil
}

// bindOrderbook binds a generic wrapper to an already deployed contract.
func bindOrderbook(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OrderbookMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Orderbook *OrderbookRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Orderbook.Contract.OrderbookCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Orderbook *OrderbookRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Orderbook.Contract.OrderbookTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Orderbook *OrderbookRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Orderbook.Contract.OrderbookTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Orderbook *OrderbookCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Orderbook.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Orderbook *OrderbookTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Orderbook.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Orderbook *OrderbookTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Orderbook.Contract.contract.Transact(opts, method, params...)
}

// CreateOrder is a paid mutator transaction binding the contract method 0xbc2e8e5f.
//
// Solidity: function createOrder(bytes data) returns()
func (_Orderbook *OrderbookTransactor) CreateOrder(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _Orderbook.contract.Transact(opts, "createOrder", data)
}

// CreateOrder is a paid mutator transaction binding the contract method 0xbc2e8e5f.
//
// Solidity: function createOrder(bytes data) returns()
func (_Orderbook *OrderbookSession) CreateOrder(data []byte) (*types.Transaction, error) {
	return _Orderbook.Contract.CreateOrder(&_Orderbook.TransactOpts, data)
}

// CreateOrder is a paid mutator transaction binding the contract method 0xbc2e8e5f.
//
// Solidity: function createOrder(bytes data) returns()
func (_Orderbook *OrderbookTransactorSession) CreateOrder(data []byte) (*types.Transaction, error) {
	return _Orderbook.Contract.CreateOrder(&_Orderbook.TransactOpts, data)
}

// FillOrder is a paid mutator transaction binding the contract method 0xd15b3858.
//
// Solidity: function fillOrder(bytes data) returns()
func (_Orderbook *OrderbookTransactor) FillOrder(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _Orderbook.contract.Transact(opts, "fillOrder", data)
}

// FillOrder is a paid mutator transaction binding the contract method 0xd15b3858.
//
// Solidity: function fillOrder(bytes data) returns()
func (_Orderbook *OrderbookSession) FillOrder(data []byte) (*types.Transaction, error) {
	return _Orderbook.Contract.FillOrder(&_Orderbook.TransactOpts, data)
}

// FillOrder is a paid mutator transaction binding the contract method 0xd15b3858.
//
// Solidity: function fillOrder(bytes data) returns()
func (_Orderbook *OrderbookTransactorSession) FillOrder(data []byte) (*types.Transaction, error) {
	return _Orderbook.Contract.FillOrder(&_Orderbook.TransactOpts, data)
}

// OrderbookCreatedIterator is returned from FilterCreated and is used to iterate over the raw logs and unpacked data for Created events raised by the Orderbook contract.
type OrderbookCreatedIterator struct {
	Event *OrderbookCreated // Event containing the contract specifics and raw log

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
func (it *OrderbookCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OrderbookCreated)
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
		it.Event = new(OrderbookCreated)
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
func (it *OrderbookCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OrderbookCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OrderbookCreated represents a Created event raised by the Orderbook contract.
type OrderbookCreated struct {
	Data []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterCreated is a free log retrieval operation binding the contract event 0x3b4a808d54c495ced21b3f493818c8c5a90f41e600c8b12208571402dcd23cec.
//
// Solidity: event Created(bytes data)
func (_Orderbook *OrderbookFilterer) FilterCreated(opts *bind.FilterOpts) (*OrderbookCreatedIterator, error) {

	logs, sub, err := _Orderbook.contract.FilterLogs(opts, "Created")
	if err != nil {
		return nil, err
	}
	return &OrderbookCreatedIterator{contract: _Orderbook.contract, event: "Created", logs: logs, sub: sub}, nil
}

// WatchCreated is a free log subscription operation binding the contract event 0x3b4a808d54c495ced21b3f493818c8c5a90f41e600c8b12208571402dcd23cec.
//
// Solidity: event Created(bytes data)
func (_Orderbook *OrderbookFilterer) WatchCreated(opts *bind.WatchOpts, sink chan<- *OrderbookCreated) (event.Subscription, error) {

	logs, sub, err := _Orderbook.contract.WatchLogs(opts, "Created")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OrderbookCreated)
				if err := _Orderbook.contract.UnpackLog(event, "Created", log); err != nil {
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

// ParseCreated is a log parse operation binding the contract event 0x3b4a808d54c495ced21b3f493818c8c5a90f41e600c8b12208571402dcd23cec.
//
// Solidity: event Created(bytes data)
func (_Orderbook *OrderbookFilterer) ParseCreated(log types.Log) (*OrderbookCreated, error) {
	event := new(OrderbookCreated)
	if err := _Orderbook.contract.UnpackLog(event, "Created", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OrderbookFilledIterator is returned from FilterFilled and is used to iterate over the raw logs and unpacked data for Filled events raised by the Orderbook contract.
type OrderbookFilledIterator struct {
	Event *OrderbookFilled // Event containing the contract specifics and raw log

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
func (it *OrderbookFilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OrderbookFilled)
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
		it.Event = new(OrderbookFilled)
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
func (it *OrderbookFilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OrderbookFilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OrderbookFilled represents a Filled event raised by the Orderbook contract.
type OrderbookFilled struct {
	Data []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterFilled is a free log retrieval operation binding the contract event 0x4bf70fd1dfd385f5ae9a107d7fff4995db1197c3bcb51f665f8a667df32aabb2.
//
// Solidity: event Filled(bytes data)
func (_Orderbook *OrderbookFilterer) FilterFilled(opts *bind.FilterOpts) (*OrderbookFilledIterator, error) {

	logs, sub, err := _Orderbook.contract.FilterLogs(opts, "Filled")
	if err != nil {
		return nil, err
	}
	return &OrderbookFilledIterator{contract: _Orderbook.contract, event: "Filled", logs: logs, sub: sub}, nil
}

// WatchFilled is a free log subscription operation binding the contract event 0x4bf70fd1dfd385f5ae9a107d7fff4995db1197c3bcb51f665f8a667df32aabb2.
//
// Solidity: event Filled(bytes data)
func (_Orderbook *OrderbookFilterer) WatchFilled(opts *bind.WatchOpts, sink chan<- *OrderbookFilled) (event.Subscription, error) {

	logs, sub, err := _Orderbook.contract.WatchLogs(opts, "Filled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OrderbookFilled)
				if err := _Orderbook.contract.UnpackLog(event, "Filled", log); err != nil {
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

// ParseFilled is a log parse operation binding the contract event 0x4bf70fd1dfd385f5ae9a107d7fff4995db1197c3bcb51f665f8a667df32aabb2.
//
// Solidity: event Filled(bytes data)
func (_Orderbook *OrderbookFilterer) ParseFilled(log types.Log) (*OrderbookFilled, error) {
	event := new(OrderbookFilled)
	if err := _Orderbook.contract.UnpackLog(event, "Filled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
