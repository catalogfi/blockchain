// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// InstantWalletFactoryMetaData contains all meta data concerning the InstantWalletFactory contract.
var InstantWalletFactoryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_wallet\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_initDataID\",\"type\":\"bytes32\"}],\"name\":\"InstantWalletCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_implementationAddr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_initData\",\"type\":\"bytes\"}],\"name\":\"createInstantWallet\",\"outputs\":[{\"internalType\":\"contractInstantWallet\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_implementationAddr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_initData\",\"type\":\"bytes\"}],\"name\":\"instantWalletAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610b21806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80633ea6ac2f1461003b57806377cb29891461006a575b600080fd5b61004e61004936600461025f565b61007d565b6040516001600160a01b03909116815260200160405180910390f35b61004e61007836600461025f565b610189565b6000808280519060200120905060ff60f81b30826040516020016100a391815260200190565b60405160208183030381529060405280519060200120604051806020016100c99061023c565b601f1982820381018352601f9091011660408190526100ee9089908990602001610353565b60408051601f198184030181529082905261010c9291602001610395565b6040516020818303038152906040528051906020012060405160200161016994939291906001600160f81b031994909416845260609290921b6bffffffffffffffffffffffff191660018401526015830152603582015260550190565b60408051601f198184030181529190528051602090910120949350505050565b600080828051906020012090506000816040516020016101ab91815260200190565b6040516020818303038152906040528051906020012085856040516101cf9061023c565b6101da929190610353565b8190604051809103906000f59050801580156101fa573d6000803e3d6000fd5b50905081816001600160a01b03167fa2acd6c8ec804a7e18acf9f20e2c5ce22c139f0dc5724abeca3ac823c3829dff60405160405180910390a3949350505050565b610727806103c583390190565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561027257600080fd5b82356001600160a01b038116811461028957600080fd5b9150602083013567ffffffffffffffff808211156102a657600080fd5b818501915085601f8301126102ba57600080fd5b8135818111156102cc576102cc610249565b604051601f8201601f19908116603f011681019083821181831017156102f4576102f4610249565b8160405282815288602084870101111561030d57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60005b8381101561034a578181015183820152602001610332565b50506000910152565b60018060a01b0383168152604060208201526000825180604084015261038081606085016020870161032f565b601f01601f1916919091016060019392505050565b600083516103a781846020880161032f565b8351908301906103bb81836020880161032f565b0194935050505056fe608060405260405161072738038061072783398101604081905261002291610319565b61002e82826000610035565b5050610436565b61003e8361006b565b60008251118061004b5750805b156100665761006483836100ab60201b6100291760201c565b505b505050565b610074816100d7565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606100d08383604051806060016040528060278152602001610700602791396101a9565b9392505050565b6100ea8161022260201b6100551760201c565b6101515760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b60648201526084015b60405180910390fd5b806101887f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc60001b61023160201b6100641760201c565b80546001600160a01b0319166001600160a01b039290921691909117905550565b6060600080856001600160a01b0316856040516101c691906103e7565b600060405180830381855af49150503d8060008114610201576040519150601f19603f3d011682016040523d82523d6000602084013e610206565b606091505b50909250905061021886838387610234565b9695505050505050565b6001600160a01b03163b151590565b90565b606083156102a357825160000361029c576001600160a01b0385163b61029c5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610148565b50816102ad565b6102ad83836102b5565b949350505050565b8151156102c55781518083602001fd5b8060405162461bcd60e51b81526004016101489190610403565b634e487b7160e01b600052604160045260246000fd5b60005b838110156103105781810151838201526020016102f8565b50506000910152565b6000806040838503121561032c57600080fd5b82516001600160a01b038116811461034357600080fd5b60208401519092506001600160401b038082111561036057600080fd5b818501915085601f83011261037457600080fd5b815181811115610386576103866102df565b604051601f8201601f19908116603f011681019083821181831017156103ae576103ae6102df565b816040528281528860208487010111156103c757600080fd5b6103d88360208301602088016102f5565b80955050505050509250929050565b600082516103f98184602087016102f5565b9190910192915050565b60208152600082518060208401526104228160408501602087016102f5565b601f01601f19169190910160400192915050565b6102bb806104456000396000f3fe60806040523661001357610011610017565b005b6100115b610027610022610067565b61009f565b565b606061004e838360405180606001604052806027815260200161025f602791396100c3565b9392505050565b6001600160a01b03163b151590565b90565b600061009a7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b905090565b3660008037600080366000845af43d6000803e8080156100be573d6000f35b3d6000fd5b6060600080856001600160a01b0316856040516100e0919061020f565b600060405180830381855af49150503d806000811461011b576040519150601f19603f3d011682016040523d82523d6000602084013e610120565b606091505b50915091506101318683838761013b565b9695505050505050565b606083156101af5782516000036101a8576001600160a01b0385163b6101a85760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064015b60405180910390fd5b50816101b9565b6101b983836101c1565b949350505050565b8151156101d15781518083602001fd5b8060405162461bcd60e51b815260040161019f919061022b565b60005b838110156102065781810151838201526020016101ee565b50506000910152565b600082516102218184602087016101eb565b9190910192915050565b602081526000825180602084015261024a8160408501602087016101eb565b601f01601f1916919091016040019291505056fe416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220b15214f7e07895c9f4804cedbe1f39c3b7bf93f6f92ea36aac6e2248c6d2403064736f6c63430008120033416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220ef773297df7c4a479ef47113a82e4af2ea273e6510d6ec3fefc3376878b46e3f64736f6c63430008120033",
}

// InstantWalletFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use InstantWalletFactoryMetaData.ABI instead.
var InstantWalletFactoryABI = InstantWalletFactoryMetaData.ABI

// InstantWalletFactoryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use InstantWalletFactoryMetaData.Bin instead.
var InstantWalletFactoryBin = InstantWalletFactoryMetaData.Bin

// DeployInstantWalletFactory deploys a new Ethereum contract, binding an instance of InstantWalletFactory to it.
func DeployInstantWalletFactory(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *InstantWalletFactory, error) {
	parsed, err := InstantWalletFactoryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(InstantWalletFactoryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &InstantWalletFactory{InstantWalletFactoryCaller: InstantWalletFactoryCaller{contract: contract}, InstantWalletFactoryTransactor: InstantWalletFactoryTransactor{contract: contract}, InstantWalletFactoryFilterer: InstantWalletFactoryFilterer{contract: contract}}, nil
}

// InstantWalletFactory is an auto generated Go binding around an Ethereum contract.
type InstantWalletFactory struct {
	InstantWalletFactoryCaller     // Read-only binding to the contract
	InstantWalletFactoryTransactor // Write-only binding to the contract
	InstantWalletFactoryFilterer   // Log filterer for contract events
}

// InstantWalletFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type InstantWalletFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InstantWalletFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type InstantWalletFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InstantWalletFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type InstantWalletFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InstantWalletFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type InstantWalletFactorySession struct {
	Contract     *InstantWalletFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// InstantWalletFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type InstantWalletFactoryCallerSession struct {
	Contract *InstantWalletFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// InstantWalletFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type InstantWalletFactoryTransactorSession struct {
	Contract     *InstantWalletFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// InstantWalletFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type InstantWalletFactoryRaw struct {
	Contract *InstantWalletFactory // Generic contract binding to access the raw methods on
}

// InstantWalletFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type InstantWalletFactoryCallerRaw struct {
	Contract *InstantWalletFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// InstantWalletFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type InstantWalletFactoryTransactorRaw struct {
	Contract *InstantWalletFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInstantWalletFactory creates a new instance of InstantWalletFactory, bound to a specific deployed contract.
func NewInstantWalletFactory(address common.Address, backend bind.ContractBackend) (*InstantWalletFactory, error) {
	contract, err := bindInstantWalletFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &InstantWalletFactory{InstantWalletFactoryCaller: InstantWalletFactoryCaller{contract: contract}, InstantWalletFactoryTransactor: InstantWalletFactoryTransactor{contract: contract}, InstantWalletFactoryFilterer: InstantWalletFactoryFilterer{contract: contract}}, nil
}

// NewInstantWalletFactoryCaller creates a new read-only instance of InstantWalletFactory, bound to a specific deployed contract.
func NewInstantWalletFactoryCaller(address common.Address, caller bind.ContractCaller) (*InstantWalletFactoryCaller, error) {
	contract, err := bindInstantWalletFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InstantWalletFactoryCaller{contract: contract}, nil
}

// NewInstantWalletFactoryTransactor creates a new write-only instance of InstantWalletFactory, bound to a specific deployed contract.
func NewInstantWalletFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*InstantWalletFactoryTransactor, error) {
	contract, err := bindInstantWalletFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InstantWalletFactoryTransactor{contract: contract}, nil
}

// NewInstantWalletFactoryFilterer creates a new log filterer instance of InstantWalletFactory, bound to a specific deployed contract.
func NewInstantWalletFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*InstantWalletFactoryFilterer, error) {
	contract, err := bindInstantWalletFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InstantWalletFactoryFilterer{contract: contract}, nil
}

// bindInstantWalletFactory binds a generic wrapper to an already deployed contract.
func bindInstantWalletFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := InstantWalletFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InstantWalletFactory *InstantWalletFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _InstantWalletFactory.Contract.InstantWalletFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InstantWalletFactory *InstantWalletFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InstantWalletFactory.Contract.InstantWalletFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InstantWalletFactory *InstantWalletFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InstantWalletFactory.Contract.InstantWalletFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InstantWalletFactory *InstantWalletFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _InstantWalletFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InstantWalletFactory *InstantWalletFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InstantWalletFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InstantWalletFactory *InstantWalletFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InstantWalletFactory.Contract.contract.Transact(opts, method, params...)
}

// InstantWalletAddress is a free data retrieval call binding the contract method 0x3ea6ac2f.
//
// Solidity: function instantWalletAddress(address _implementationAddr, bytes _initData) view returns(address)
func (_InstantWalletFactory *InstantWalletFactoryCaller) InstantWalletAddress(opts *bind.CallOpts, _implementationAddr common.Address, _initData []byte) (common.Address, error) {
	var out []interface{}
	err := _InstantWalletFactory.contract.Call(opts, &out, "instantWalletAddress", _implementationAddr, _initData)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// InstantWalletAddress is a free data retrieval call binding the contract method 0x3ea6ac2f.
//
// Solidity: function instantWalletAddress(address _implementationAddr, bytes _initData) view returns(address)
func (_InstantWalletFactory *InstantWalletFactorySession) InstantWalletAddress(_implementationAddr common.Address, _initData []byte) (common.Address, error) {
	return _InstantWalletFactory.Contract.InstantWalletAddress(&_InstantWalletFactory.CallOpts, _implementationAddr, _initData)
}

// InstantWalletAddress is a free data retrieval call binding the contract method 0x3ea6ac2f.
//
// Solidity: function instantWalletAddress(address _implementationAddr, bytes _initData) view returns(address)
func (_InstantWalletFactory *InstantWalletFactoryCallerSession) InstantWalletAddress(_implementationAddr common.Address, _initData []byte) (common.Address, error) {
	return _InstantWalletFactory.Contract.InstantWalletAddress(&_InstantWalletFactory.CallOpts, _implementationAddr, _initData)
}

// CreateInstantWallet is a paid mutator transaction binding the contract method 0x77cb2989.
//
// Solidity: function createInstantWallet(address _implementationAddr, bytes _initData) returns(address)
func (_InstantWalletFactory *InstantWalletFactoryTransactor) CreateInstantWallet(opts *bind.TransactOpts, _implementationAddr common.Address, _initData []byte) (*types.Transaction, error) {
	return _InstantWalletFactory.contract.Transact(opts, "createInstantWallet", _implementationAddr, _initData)
}

// CreateInstantWallet is a paid mutator transaction binding the contract method 0x77cb2989.
//
// Solidity: function createInstantWallet(address _implementationAddr, bytes _initData) returns(address)
func (_InstantWalletFactory *InstantWalletFactorySession) CreateInstantWallet(_implementationAddr common.Address, _initData []byte) (*types.Transaction, error) {
	return _InstantWalletFactory.Contract.CreateInstantWallet(&_InstantWalletFactory.TransactOpts, _implementationAddr, _initData)
}

// CreateInstantWallet is a paid mutator transaction binding the contract method 0x77cb2989.
//
// Solidity: function createInstantWallet(address _implementationAddr, bytes _initData) returns(address)
func (_InstantWalletFactory *InstantWalletFactoryTransactorSession) CreateInstantWallet(_implementationAddr common.Address, _initData []byte) (*types.Transaction, error) {
	return _InstantWalletFactory.Contract.CreateInstantWallet(&_InstantWalletFactory.TransactOpts, _implementationAddr, _initData)
}

// InstantWalletFactoryInstantWalletCreatedIterator is returned from FilterInstantWalletCreated and is used to iterate over the raw logs and unpacked data for InstantWalletCreated events raised by the InstantWalletFactory contract.
type InstantWalletFactoryInstantWalletCreatedIterator struct {
	Event *InstantWalletFactoryInstantWalletCreated // Event containing the contract specifics and raw log

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
func (it *InstantWalletFactoryInstantWalletCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InstantWalletFactoryInstantWalletCreated)
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
		it.Event = new(InstantWalletFactoryInstantWalletCreated)
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
func (it *InstantWalletFactoryInstantWalletCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InstantWalletFactoryInstantWalletCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InstantWalletFactoryInstantWalletCreated represents a InstantWalletCreated event raised by the InstantWalletFactory contract.
type InstantWalletFactoryInstantWalletCreated struct {
	Wallet     common.Address
	InitDataID [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterInstantWalletCreated is a free log retrieval operation binding the contract event 0xa2acd6c8ec804a7e18acf9f20e2c5ce22c139f0dc5724abeca3ac823c3829dff.
//
// Solidity: event InstantWalletCreated(address indexed _wallet, bytes32 indexed _initDataID)
func (_InstantWalletFactory *InstantWalletFactoryFilterer) FilterInstantWalletCreated(opts *bind.FilterOpts, _wallet []common.Address, _initDataID [][32]byte) (*InstantWalletFactoryInstantWalletCreatedIterator, error) {

	var _walletRule []interface{}
	for _, _walletItem := range _wallet {
		_walletRule = append(_walletRule, _walletItem)
	}
	var _initDataIDRule []interface{}
	for _, _initDataIDItem := range _initDataID {
		_initDataIDRule = append(_initDataIDRule, _initDataIDItem)
	}

	logs, sub, err := _InstantWalletFactory.contract.FilterLogs(opts, "InstantWalletCreated", _walletRule, _initDataIDRule)
	if err != nil {
		return nil, err
	}
	return &InstantWalletFactoryInstantWalletCreatedIterator{contract: _InstantWalletFactory.contract, event: "InstantWalletCreated", logs: logs, sub: sub}, nil
}

// WatchInstantWalletCreated is a free log subscription operation binding the contract event 0xa2acd6c8ec804a7e18acf9f20e2c5ce22c139f0dc5724abeca3ac823c3829dff.
//
// Solidity: event InstantWalletCreated(address indexed _wallet, bytes32 indexed _initDataID)
func (_InstantWalletFactory *InstantWalletFactoryFilterer) WatchInstantWalletCreated(opts *bind.WatchOpts, sink chan<- *InstantWalletFactoryInstantWalletCreated, _wallet []common.Address, _initDataID [][32]byte) (event.Subscription, error) {

	var _walletRule []interface{}
	for _, _walletItem := range _wallet {
		_walletRule = append(_walletRule, _walletItem)
	}
	var _initDataIDRule []interface{}
	for _, _initDataIDItem := range _initDataID {
		_initDataIDRule = append(_initDataIDRule, _initDataIDItem)
	}

	logs, sub, err := _InstantWalletFactory.contract.WatchLogs(opts, "InstantWalletCreated", _walletRule, _initDataIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InstantWalletFactoryInstantWalletCreated)
				if err := _InstantWalletFactory.contract.UnpackLog(event, "InstantWalletCreated", log); err != nil {
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

// ParseInstantWalletCreated is a log parse operation binding the contract event 0xa2acd6c8ec804a7e18acf9f20e2c5ce22c139f0dc5724abeca3ac823c3829dff.
//
// Solidity: event InstantWalletCreated(address indexed _wallet, bytes32 indexed _initDataID)
func (_InstantWalletFactory *InstantWalletFactoryFilterer) ParseInstantWalletCreated(log types.Log) (*InstantWalletFactoryInstantWalletCreated, error) {
	event := new(InstantWalletFactoryInstantWalletCreated)
	if err := _InstantWalletFactory.contract.UnpackLog(event, "InstantWalletCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
