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
)

// AtomicSwapMetaData contains all meta data concerning the AtomicSwap contract.
var AtomicSwapMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initiatedAt\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Initiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"secrectHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"secret\",\"type\":\"bytes\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"atomicSwapOrders\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initiatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isFulfilled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_secretHash\",\"type\":\"bytes32\"}],\"name\":\"initiate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_orderId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_secret\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_orderId\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610f3c380380610f3c83398101604081905261002f91610040565b6001600160a01b0316608052610070565b60006020828403121561005257600080fd5b81516001600160a01b038116811461006957600080fd5b9392505050565b608051610e9c6100a06000396000818161012e015281816102ab015281816106a801526108f00152610e9c6000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80633f7b9c381461005c5780637249fbb6146100ee57806397ffc7ae14610103578063f7ff720714610116578063fc0c546a14610129575b600080fd5b6100af61006a366004610c2c565b6000602081905290815260409020805460018201546002830154600384015460048501546005909501546001600160a01b039485169593909416939192909160ff1686565b604080516001600160a01b0397881681529690951660208701529385019290925260608401526080830152151560a082015260c0015b60405180910390f35b6101016100fc366004610c2c565b610168565b005b610101610111366004610c45565b6102d9565b610101610124366004610c8c565b6106df565b6101507f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020016100e5565b600081815260208190526040902080546001600160a01b03166101d25760405162461bcd60e51b815260206004820152601e60248201527f41746f6d6963537761703a206f72646572206e6f7420696e697461746564000060448201526064015b60405180910390fd5b600581015460ff16156101f75760405162461bcd60e51b81526004016101c990610d08565b438160020154826003015461020c9190610d4b565b106102595760405162461bcd60e51b815260206004820152601d60248201527f41746f6d6963537761703a206f72646572206e6f74206578706972656400000060448201526064016101c9565b60058101805460ff1916600117905560405182907ffe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf090600090a2600181015460048201546102d5916001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000811692911690610921565b5050565b833384846001600160a01b03841661033f5760405162461bcd60e51b8152602060048201526024808201527f41746f6d6963537761703a20696e76616c69642072656465656d6572206164646044820152637265737360e01b60648201526084016101c9565b836001600160a01b0316836001600160a01b0316036103be5760405162461bcd60e51b815260206004820152603560248201527f41746f6d6963537761703a2072656465656d657220616e6420696e69746961746044820152746f722063616e6e6f74206265207468652073616d6560581b60648201526084016101c9565b600082116104255760405162461bcd60e51b815260206004820152602e60248201527f41746f6d6963537761703a206578706972792073686f756c642062652067726560448201526d61746572207468616e207a65726f60901b60648201526084016101c9565b6000811161047f5760405162461bcd60e51b815260206004820152602160248201527f41746f6d6963537761703a20616d6f756e742063616e6e6f74206265207a65726044820152606f60f81b60648201526084016101c9565b6000600286336040516020016104a89291909182526001600160a01b0316602082015260400190565b60408051601f19818403018152908290526104c291610d96565b602060405180830381855afa1580156104df573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906105029190610db2565b60008181526020818152604091829020825160c08101845281546001600160a01b0390811680835260018401549091169382019390935260028201549381019390935260038101546060840152600481015460808401526005015460ff16151560a083015291925090156105b85760405162461bcd60e51b815260206004820152601b60248201527f41746f6d6963537761703a206475706c6963617465206f72646572000000000060448201526064016101c9565b6040805160c0810182526001600160a01b038c811682523360208084019182528385018e81524360608601908152608086018f8152600060a088018181528b825281865290899020885181546001600160a01b0319908116918a16919091178255965160018201805490981698169790971790955591516002860155516003850181905590516004850181905592516005909401805460ff19169415159490941790935584519283528201529091899185917f3dd1f59c2a4b236fc1e76892b9a4b62de617c6a44a56ed208a3ba79c589823ab910160405180910390a360808101516106d2906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169033903090610989565b5050505050505050505050565b600083815260208190526040902080546001600160a01b03166107445760405162461bcd60e51b815260206004820152601e60248201527f41746f6d6963537761703a206f72646572206e6f7420696e697461746564000060448201526064016101c9565b600581015460ff16156107695760405162461bcd60e51b81526004016101c990610d08565b60006002848460405161077d929190610dcb565b602060405180830381855afa15801561079a573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906107bd9190610db2565b600183015460408051602081018490526001600160a01b0390921690820152909150859060029060600160408051601f198184030181529082905261080191610d96565b602060405180830381855afa15801561081e573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906108419190610db2565b1461088e5760405162461bcd60e51b815260206004820152601a60248201527f41746f6d6963537761703a20696e76616c69642073656372657400000000000060448201526064016101c9565b60058201805460ff19166001179055604051819086907f4c9a044220477b4e94dbb0d07ff6ff4ac30d443bef59098c4541b006954778e2906108d39088908890610ddb565b60405180910390a38154600483015461091a916001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000811692911690610921565b5050505050565b6040516001600160a01b03831660248201526044810182905261098490849063a9059cbb60e01b906064015b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091526109c7565b505050565b6040516001600160a01b03808516602483015283166044820152606481018290526109c19085906323b872dd60e01b9060840161094d565b50505050565b6000610a1c826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b0316610a9c9092919063ffffffff16565b9050805160001480610a3d575080806020019051810190610a3d9190610e0a565b6109845760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b60648201526084016101c9565b6060610aab8484600085610ab3565b949350505050565b606082471015610b145760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b60648201526084016101c9565b600080866001600160a01b03168587604051610b309190610d96565b60006040518083038185875af1925050503d8060008114610b6d576040519150601f19603f3d011682016040523d82523d6000602084013e610b72565b606091505b5091509150610b8387838387610b8e565b979650505050505050565b60608315610bfd578251600003610bf6576001600160a01b0385163b610bf65760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016101c9565b5081610aab565b610aab8383815115610c125781518083602001fd5b8060405162461bcd60e51b81526004016101c99190610e33565b600060208284031215610c3e57600080fd5b5035919050565b60008060008060808587031215610c5b57600080fd5b84356001600160a01b0381168114610c7257600080fd5b966020860135965060408601359560600135945092505050565b600080600060408486031215610ca157600080fd5b83359250602084013567ffffffffffffffff80821115610cc057600080fd5b818601915086601f830112610cd457600080fd5b813581811115610ce357600080fd5b876020828501011115610cf557600080fd5b6020830194508093505050509250925092565b60208082526023908201527f41746f6d6963537761703a206f7264657220616c72656164792066756c66696c6040820152621b195960ea1b606082015260800190565b80820180821115610d6c57634e487b7160e01b600052601160045260246000fd5b92915050565b60005b83811015610d8d578181015183820152602001610d75565b50506000910152565b60008251610da8818460208701610d72565b9190910192915050565b600060208284031215610dc457600080fd5b5051919050565b8183823760009101908152919050565b60208152816020820152818360408301376000818301604090810191909152601f909201601f19160101919050565b600060208284031215610e1c57600080fd5b81518015158114610e2c57600080fd5b9392505050565b6020815260008251806020840152610e52816040850160208701610d72565b601f01601f1916919091016040019291505056fea264697066735822122083a334bcafdce1a49fe2cff587175a99496fe0ce5cdcb963bad8cf57e85424bc64736f6c63430008120033",
}

// AtomicSwapABI is the input ABI used to generate the binding from.
// Deprecated: Use AtomicSwapMetaData.ABI instead.
var AtomicSwapABI = AtomicSwapMetaData.ABI

// AtomicSwapBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AtomicSwapMetaData.Bin instead.
var AtomicSwapBin = AtomicSwapMetaData.Bin

// DeployAtomicSwap deploys a new Ethereum contract, binding an instance of AtomicSwap to it.
func DeployAtomicSwap(auth *bind.TransactOpts, backend bind.ContractBackend, _token common.Address) (common.Address, *types.Transaction, *AtomicSwap, error) {
	parsed, err := AtomicSwapMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AtomicSwapBin), backend, _token)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtomicSwap{AtomicSwapCaller: AtomicSwapCaller{contract: contract}, AtomicSwapTransactor: AtomicSwapTransactor{contract: contract}, AtomicSwapFilterer: AtomicSwapFilterer{contract: contract}}, nil
}

// AtomicSwap is an auto generated Go binding around an Ethereum contract.
type AtomicSwap struct {
	AtomicSwapCaller     // Read-only binding to the contract
	AtomicSwapTransactor // Write-only binding to the contract
	AtomicSwapFilterer   // Log filterer for contract events
}

// AtomicSwapCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtomicSwapCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomicSwapTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtomicSwapFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomicSwapSession struct {
	Contract     *AtomicSwap       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomicSwapCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomicSwapCallerSession struct {
	Contract *AtomicSwapCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// AtomicSwapTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomicSwapTransactorSession struct {
	Contract     *AtomicSwapTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// AtomicSwapRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtomicSwapRaw struct {
	Contract *AtomicSwap // Generic contract binding to access the raw methods on
}

// AtomicSwapCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomicSwapCallerRaw struct {
	Contract *AtomicSwapCaller // Generic read-only contract binding to access the raw methods on
}

// AtomicSwapTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomicSwapTransactorRaw struct {
	Contract *AtomicSwapTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomicSwap creates a new instance of AtomicSwap, bound to a specific deployed contract.
func NewAtomicSwap(address common.Address, backend bind.ContractBackend) (*AtomicSwap, error) {
	contract, err := bindAtomicSwap(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtomicSwap{AtomicSwapCaller: AtomicSwapCaller{contract: contract}, AtomicSwapTransactor: AtomicSwapTransactor{contract: contract}, AtomicSwapFilterer: AtomicSwapFilterer{contract: contract}}, nil
}

// NewAtomicSwapCaller creates a new read-only instance of AtomicSwap, bound to a specific deployed contract.
func NewAtomicSwapCaller(address common.Address, caller bind.ContractCaller) (*AtomicSwapCaller, error) {
	contract, err := bindAtomicSwap(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapCaller{contract: contract}, nil
}

// NewAtomicSwapTransactor creates a new write-only instance of AtomicSwap, bound to a specific deployed contract.
func NewAtomicSwapTransactor(address common.Address, transactor bind.ContractTransactor) (*AtomicSwapTransactor, error) {
	contract, err := bindAtomicSwap(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapTransactor{contract: contract}, nil
}

// NewAtomicSwapFilterer creates a new log filterer instance of AtomicSwap, bound to a specific deployed contract.
func NewAtomicSwapFilterer(address common.Address, filterer bind.ContractFilterer) (*AtomicSwapFilterer, error) {
	contract, err := bindAtomicSwap(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapFilterer{contract: contract}, nil
}

// bindAtomicSwap binds a generic wrapper to an already deployed contract.
func bindAtomicSwap(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomicSwapABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomicSwap *AtomicSwapRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AtomicSwap.Contract.AtomicSwapCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomicSwap *AtomicSwapRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomicSwap.Contract.AtomicSwapTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomicSwap *AtomicSwapRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomicSwap.Contract.AtomicSwapTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomicSwap *AtomicSwapCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AtomicSwap.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomicSwap *AtomicSwapTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomicSwap.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomicSwap *AtomicSwapTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomicSwap.Contract.contract.Transact(opts, method, params...)
}

// AtomicSwapOrders is a free data retrieval call binding the contract method 0x3f7b9c38.
//
// Solidity: function atomicSwapOrders(bytes32 ) view returns(address redeemer, address initiator, uint256 expiry, uint256 initiatedAt, uint256 amount, bool isFulfilled)
func (_AtomicSwap *AtomicSwapCaller) AtomicSwapOrders(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Redeemer    common.Address
	Initiator   common.Address
	Expiry      *big.Int
	InitiatedAt *big.Int
	Amount      *big.Int
	IsFulfilled bool
}, error) {
	var out []interface{}
	err := _AtomicSwap.contract.Call(opts, &out, "atomicSwapOrders", arg0)

	outstruct := new(struct {
		Redeemer    common.Address
		Initiator   common.Address
		Expiry      *big.Int
		InitiatedAt *big.Int
		Amount      *big.Int
		IsFulfilled bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Redeemer = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Initiator = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Expiry = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.InitiatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Amount = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.IsFulfilled = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

// AtomicSwapOrders is a free data retrieval call binding the contract method 0x3f7b9c38.
//
// Solidity: function atomicSwapOrders(bytes32 ) view returns(address redeemer, address initiator, uint256 expiry, uint256 initiatedAt, uint256 amount, bool isFulfilled)
func (_AtomicSwap *AtomicSwapSession) AtomicSwapOrders(arg0 [32]byte) (struct {
	Redeemer    common.Address
	Initiator   common.Address
	Expiry      *big.Int
	InitiatedAt *big.Int
	Amount      *big.Int
	IsFulfilled bool
}, error) {
	return _AtomicSwap.Contract.AtomicSwapOrders(&_AtomicSwap.CallOpts, arg0)
}

// AtomicSwapOrders is a free data retrieval call binding the contract method 0x3f7b9c38.
//
// Solidity: function atomicSwapOrders(bytes32 ) view returns(address redeemer, address initiator, uint256 expiry, uint256 initiatedAt, uint256 amount, bool isFulfilled)
func (_AtomicSwap *AtomicSwapCallerSession) AtomicSwapOrders(arg0 [32]byte) (struct {
	Redeemer    common.Address
	Initiator   common.Address
	Expiry      *big.Int
	InitiatedAt *big.Int
	Amount      *big.Int
	IsFulfilled bool
}, error) {
	return _AtomicSwap.Contract.AtomicSwapOrders(&_AtomicSwap.CallOpts, arg0)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_AtomicSwap *AtomicSwapCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AtomicSwap.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_AtomicSwap *AtomicSwapSession) Token() (common.Address, error) {
	return _AtomicSwap.Contract.Token(&_AtomicSwap.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_AtomicSwap *AtomicSwapCallerSession) Token() (common.Address, error) {
	return _AtomicSwap.Contract.Token(&_AtomicSwap.CallOpts)
}

// Initiate is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address _redeemer, uint256 _expiry, uint256 _amount, bytes32 _secretHash) returns()
func (_AtomicSwap *AtomicSwapTransactor) Initiate(opts *bind.TransactOpts, _redeemer common.Address, _expiry *big.Int, _amount *big.Int, _secretHash [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.contract.Transact(opts, "initiate", _redeemer, _expiry, _amount, _secretHash)
}

// Initiate is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address _redeemer, uint256 _expiry, uint256 _amount, bytes32 _secretHash) returns()
func (_AtomicSwap *AtomicSwapSession) Initiate(_redeemer common.Address, _expiry *big.Int, _amount *big.Int, _secretHash [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Initiate(&_AtomicSwap.TransactOpts, _redeemer, _expiry, _amount, _secretHash)
}

// Initiate is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address _redeemer, uint256 _expiry, uint256 _amount, bytes32 _secretHash) returns()
func (_AtomicSwap *AtomicSwapTransactorSession) Initiate(_redeemer common.Address, _expiry *big.Int, _amount *big.Int, _secretHash [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Initiate(&_AtomicSwap.TransactOpts, _redeemer, _expiry, _amount, _secretHash)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 _orderId, bytes _secret) returns()
func (_AtomicSwap *AtomicSwapTransactor) Redeem(opts *bind.TransactOpts, _orderId [32]byte, _secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.contract.Transact(opts, "redeem", _orderId, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 _orderId, bytes _secret) returns()
func (_AtomicSwap *AtomicSwapSession) Redeem(_orderId [32]byte, _secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Redeem(&_AtomicSwap.TransactOpts, _orderId, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 _orderId, bytes _secret) returns()
func (_AtomicSwap *AtomicSwapTransactorSession) Redeem(_orderId [32]byte, _secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Redeem(&_AtomicSwap.TransactOpts, _orderId, _secret)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _orderId) returns()
func (_AtomicSwap *AtomicSwapTransactor) Refund(opts *bind.TransactOpts, _orderId [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.contract.Transact(opts, "refund", _orderId)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _orderId) returns()
func (_AtomicSwap *AtomicSwapSession) Refund(_orderId [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Refund(&_AtomicSwap.TransactOpts, _orderId)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _orderId) returns()
func (_AtomicSwap *AtomicSwapTransactorSession) Refund(_orderId [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Refund(&_AtomicSwap.TransactOpts, _orderId)
}

// AtomicSwapInitiatedIterator is returned from FilterInitiated and is used to iterate over the raw logs and unpacked data for Initiated events raised by the AtomicSwap contract.
type AtomicSwapInitiatedIterator struct {
	Event *AtomicSwapInitiated // Event containing the contract specifics and raw log

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
func (it *AtomicSwapInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomicSwapInitiated)
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
		it.Event = new(AtomicSwapInitiated)
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
func (it *AtomicSwapInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomicSwapInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomicSwapInitiated represents a Initiated event raised by the AtomicSwap contract.
type AtomicSwapInitiated struct {
	OrderId     [32]byte
	SecretHash  [32]byte
	InitiatedAt *big.Int
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterInitiated is a free log retrieval operation binding the contract event 0x3dd1f59c2a4b236fc1e76892b9a4b62de617c6a44a56ed208a3ba79c589823ab.
//
// Solidity: event Initiated(bytes32 indexed orderId, bytes32 indexed secretHash, uint256 initiatedAt, uint256 amount)
func (_AtomicSwap *AtomicSwapFilterer) FilterInitiated(opts *bind.FilterOpts, orderId [][32]byte, secretHash [][32]byte) (*AtomicSwapInitiatedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}

	logs, sub, err := _AtomicSwap.contract.FilterLogs(opts, "Initiated", orderIdRule, secretHashRule)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapInitiatedIterator{contract: _AtomicSwap.contract, event: "Initiated", logs: logs, sub: sub}, nil
}

// WatchInitiated is a free log subscription operation binding the contract event 0x3dd1f59c2a4b236fc1e76892b9a4b62de617c6a44a56ed208a3ba79c589823ab.
//
// Solidity: event Initiated(bytes32 indexed orderId, bytes32 indexed secretHash, uint256 initiatedAt, uint256 amount)
func (_AtomicSwap *AtomicSwapFilterer) WatchInitiated(opts *bind.WatchOpts, sink chan<- *AtomicSwapInitiated, orderId [][32]byte, secretHash [][32]byte) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}

	logs, sub, err := _AtomicSwap.contract.WatchLogs(opts, "Initiated", orderIdRule, secretHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomicSwapInitiated)
				if err := _AtomicSwap.contract.UnpackLog(event, "Initiated", log); err != nil {
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

// ParseInitiated is a log parse operation binding the contract event 0x3dd1f59c2a4b236fc1e76892b9a4b62de617c6a44a56ed208a3ba79c589823ab.
//
// Solidity: event Initiated(bytes32 indexed orderId, bytes32 indexed secretHash, uint256 initiatedAt, uint256 amount)
func (_AtomicSwap *AtomicSwapFilterer) ParseInitiated(log types.Log) (*AtomicSwapInitiated, error) {
	event := new(AtomicSwapInitiated)
	if err := _AtomicSwap.contract.UnpackLog(event, "Initiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtomicSwapRedeemedIterator is returned from FilterRedeemed and is used to iterate over the raw logs and unpacked data for Redeemed events raised by the AtomicSwap contract.
type AtomicSwapRedeemedIterator struct {
	Event *AtomicSwapRedeemed // Event containing the contract specifics and raw log

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
func (it *AtomicSwapRedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomicSwapRedeemed)
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
		it.Event = new(AtomicSwapRedeemed)
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
func (it *AtomicSwapRedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomicSwapRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomicSwapRedeemed represents a Redeemed event raised by the AtomicSwap contract.
type AtomicSwapRedeemed struct {
	OrderId     [32]byte
	SecrectHash [32]byte
	Secret      []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x4c9a044220477b4e94dbb0d07ff6ff4ac30d443bef59098c4541b006954778e2.
//
// Solidity: event Redeemed(bytes32 indexed orderId, bytes32 indexed secrectHash, bytes secret)
func (_AtomicSwap *AtomicSwapFilterer) FilterRedeemed(opts *bind.FilterOpts, orderId [][32]byte, secrectHash [][32]byte) (*AtomicSwapRedeemedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var secrectHashRule []interface{}
	for _, secrectHashItem := range secrectHash {
		secrectHashRule = append(secrectHashRule, secrectHashItem)
	}

	logs, sub, err := _AtomicSwap.contract.FilterLogs(opts, "Redeemed", orderIdRule, secrectHashRule)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapRedeemedIterator{contract: _AtomicSwap.contract, event: "Redeemed", logs: logs, sub: sub}, nil
}

// WatchRedeemed is a free log subscription operation binding the contract event 0x4c9a044220477b4e94dbb0d07ff6ff4ac30d443bef59098c4541b006954778e2.
//
// Solidity: event Redeemed(bytes32 indexed orderId, bytes32 indexed secrectHash, bytes secret)
func (_AtomicSwap *AtomicSwapFilterer) WatchRedeemed(opts *bind.WatchOpts, sink chan<- *AtomicSwapRedeemed, orderId [][32]byte, secrectHash [][32]byte) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var secrectHashRule []interface{}
	for _, secrectHashItem := range secrectHash {
		secrectHashRule = append(secrectHashRule, secrectHashItem)
	}

	logs, sub, err := _AtomicSwap.contract.WatchLogs(opts, "Redeemed", orderIdRule, secrectHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomicSwapRedeemed)
				if err := _AtomicSwap.contract.UnpackLog(event, "Redeemed", log); err != nil {
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
// Solidity: event Redeemed(bytes32 indexed orderId, bytes32 indexed secrectHash, bytes secret)
func (_AtomicSwap *AtomicSwapFilterer) ParseRedeemed(log types.Log) (*AtomicSwapRedeemed, error) {
	event := new(AtomicSwapRedeemed)
	if err := _AtomicSwap.contract.UnpackLog(event, "Redeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtomicSwapRefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the AtomicSwap contract.
type AtomicSwapRefundedIterator struct {
	Event *AtomicSwapRefunded // Event containing the contract specifics and raw log

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
func (it *AtomicSwapRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomicSwapRefunded)
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
		it.Event = new(AtomicSwapRefunded)
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
func (it *AtomicSwapRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomicSwapRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomicSwapRefunded represents a Refunded event raised by the AtomicSwap contract.
type AtomicSwapRefunded struct {
	OrderId [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed orderId)
func (_AtomicSwap *AtomicSwapFilterer) FilterRefunded(opts *bind.FilterOpts, orderId [][32]byte) (*AtomicSwapRefundedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}

	logs, sub, err := _AtomicSwap.contract.FilterLogs(opts, "Refunded", orderIdRule)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapRefundedIterator{contract: _AtomicSwap.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed orderId)
func (_AtomicSwap *AtomicSwapFilterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *AtomicSwapRefunded, orderId [][32]byte) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}

	logs, sub, err := _AtomicSwap.contract.WatchLogs(opts, "Refunded", orderIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomicSwapRefunded)
				if err := _AtomicSwap.contract.UnpackLog(event, "Refunded", log); err != nil {
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
// Solidity: event Refunded(bytes32 indexed orderId)
func (_AtomicSwap *AtomicSwapFilterer) ParseRefunded(log types.Log) (*AtomicSwapRefunded, error) {
	event := new(AtomicSwapRefunded)
	if err := _AtomicSwap.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
