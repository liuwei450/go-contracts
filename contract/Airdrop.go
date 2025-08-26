// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// AirdropMetaData contains all meta data concerning the Airdrop contract.
var AirdropMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AirdropBNB\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AirdropERC20\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"airdropBNB\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"airdropERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gov\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newGov\",\"type\":\"address\"}],\"name\":\"setGov\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// AirdropABI is the input ABI used to generate the binding from.
// Deprecated: Use AirdropMetaData.ABI instead.
var AirdropABI = AirdropMetaData.ABI

// Airdrop is an auto generated Go binding around an Ethereum contract.
type Airdrop struct {
	AirdropCaller     // Read-only binding to the contract
	AirdropTransactor // Write-only binding to the contract
	AirdropFilterer   // Log filterer for contract events
}

// AirdropCaller is an auto generated read-only Go binding around an Ethereum contract.
type AirdropCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AirdropTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AirdropTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AirdropFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AirdropFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AirdropSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AirdropSession struct {
	Contract     *Airdrop          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AirdropCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AirdropCallerSession struct {
	Contract *AirdropCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// AirdropTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AirdropTransactorSession struct {
	Contract     *AirdropTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// AirdropRaw is an auto generated low-level Go binding around an Ethereum contract.
type AirdropRaw struct {
	Contract *Airdrop // Generic contract binding to access the raw methods on
}

// AirdropCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AirdropCallerRaw struct {
	Contract *AirdropCaller // Generic read-only contract binding to access the raw methods on
}

// AirdropTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AirdropTransactorRaw struct {
	Contract *AirdropTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAirdrop creates a new instance of Airdrop, bound to a specific deployed contract.
func NewAirdrop(address common.Address, backend bind.ContractBackend) (*Airdrop, error) {
	contract, err := bindAirdrop(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Airdrop{AirdropCaller: AirdropCaller{contract: contract}, AirdropTransactor: AirdropTransactor{contract: contract}, AirdropFilterer: AirdropFilterer{contract: contract}}, nil
}

// NewAirdropCaller creates a new read-only instance of Airdrop, bound to a specific deployed contract.
func NewAirdropCaller(address common.Address, caller bind.ContractCaller) (*AirdropCaller, error) {
	contract, err := bindAirdrop(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AirdropCaller{contract: contract}, nil
}

// NewAirdropTransactor creates a new write-only instance of Airdrop, bound to a specific deployed contract.
func NewAirdropTransactor(address common.Address, transactor bind.ContractTransactor) (*AirdropTransactor, error) {
	contract, err := bindAirdrop(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AirdropTransactor{contract: contract}, nil
}

// NewAirdropFilterer creates a new log filterer instance of Airdrop, bound to a specific deployed contract.
func NewAirdropFilterer(address common.Address, filterer bind.ContractFilterer) (*AirdropFilterer, error) {
	contract, err := bindAirdrop(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AirdropFilterer{contract: contract}, nil
}

// bindAirdrop binds a generic wrapper to an already deployed contract.
func bindAirdrop(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AirdropMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Airdrop *AirdropRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Airdrop.Contract.AirdropCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Airdrop *AirdropRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Airdrop.Contract.AirdropTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Airdrop *AirdropRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Airdrop.Contract.AirdropTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Airdrop *AirdropCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Airdrop.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Airdrop *AirdropTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Airdrop.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Airdrop *AirdropTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Airdrop.Contract.contract.Transact(opts, method, params...)
}

// Gov is a free data retrieval call binding the contract method 0x12d43a51.
//
// Solidity: function gov() view returns(address)
func (_Airdrop *AirdropCaller) Gov(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Airdrop.contract.Call(opts, &out, "gov")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Gov is a free data retrieval call binding the contract method 0x12d43a51.
//
// Solidity: function gov() view returns(address)
func (_Airdrop *AirdropSession) Gov() (common.Address, error) {
	return _Airdrop.Contract.Gov(&_Airdrop.CallOpts)
}

// Gov is a free data retrieval call binding the contract method 0x12d43a51.
//
// Solidity: function gov() view returns(address)
func (_Airdrop *AirdropCallerSession) Gov() (common.Address, error) {
	return _Airdrop.Contract.Gov(&_Airdrop.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Airdrop *AirdropCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Airdrop.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Airdrop *AirdropSession) Token() (common.Address, error) {
	return _Airdrop.Contract.Token(&_Airdrop.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Airdrop *AirdropCallerSession) Token() (common.Address, error) {
	return _Airdrop.Contract.Token(&_Airdrop.CallOpts)
}

// AirdropBNB is a paid mutator transaction binding the contract method 0x677e88c3.
//
// Solidity: function airdropBNB(address[] recipients, uint256[] amounts) payable returns()
func (_Airdrop *AirdropTransactor) AirdropBNB(opts *bind.TransactOpts, recipients []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Airdrop.contract.Transact(opts, "airdropBNB", recipients, amounts)
}

// AirdropBNB is a paid mutator transaction binding the contract method 0x677e88c3.
//
// Solidity: function airdropBNB(address[] recipients, uint256[] amounts) payable returns()
func (_Airdrop *AirdropSession) AirdropBNB(recipients []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Airdrop.Contract.AirdropBNB(&_Airdrop.TransactOpts, recipients, amounts)
}

// AirdropBNB is a paid mutator transaction binding the contract method 0x677e88c3.
//
// Solidity: function airdropBNB(address[] recipients, uint256[] amounts) payable returns()
func (_Airdrop *AirdropTransactorSession) AirdropBNB(recipients []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Airdrop.Contract.AirdropBNB(&_Airdrop.TransactOpts, recipients, amounts)
}

// AirdropERC20 is a paid mutator transaction binding the contract method 0xa3b86230.
//
// Solidity: function airdropERC20(address[] recipients, uint256[] amounts) returns()
func (_Airdrop *AirdropTransactor) AirdropERC20(opts *bind.TransactOpts, recipients []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Airdrop.contract.Transact(opts, "airdropERC20", recipients, amounts)
}

// AirdropERC20 is a paid mutator transaction binding the contract method 0xa3b86230.
//
// Solidity: function airdropERC20(address[] recipients, uint256[] amounts) returns()
func (_Airdrop *AirdropSession) AirdropERC20(recipients []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Airdrop.Contract.AirdropERC20(&_Airdrop.TransactOpts, recipients, amounts)
}

// AirdropERC20 is a paid mutator transaction binding the contract method 0xa3b86230.
//
// Solidity: function airdropERC20(address[] recipients, uint256[] amounts) returns()
func (_Airdrop *AirdropTransactorSession) AirdropERC20(recipients []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Airdrop.Contract.AirdropERC20(&_Airdrop.TransactOpts, recipients, amounts)
}

// SetGov is a paid mutator transaction binding the contract method 0xcfad57a2.
//
// Solidity: function setGov(address _newGov) returns()
func (_Airdrop *AirdropTransactor) SetGov(opts *bind.TransactOpts, _newGov common.Address) (*types.Transaction, error) {
	return _Airdrop.contract.Transact(opts, "setGov", _newGov)
}

// SetGov is a paid mutator transaction binding the contract method 0xcfad57a2.
//
// Solidity: function setGov(address _newGov) returns()
func (_Airdrop *AirdropSession) SetGov(_newGov common.Address) (*types.Transaction, error) {
	return _Airdrop.Contract.SetGov(&_Airdrop.TransactOpts, _newGov)
}

// SetGov is a paid mutator transaction binding the contract method 0xcfad57a2.
//
// Solidity: function setGov(address _newGov) returns()
func (_Airdrop *AirdropTransactorSession) SetGov(_newGov common.Address) (*types.Transaction, error) {
	return _Airdrop.Contract.SetGov(&_Airdrop.TransactOpts, _newGov)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Airdrop *AirdropTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Airdrop.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Airdrop *AirdropSession) Receive() (*types.Transaction, error) {
	return _Airdrop.Contract.Receive(&_Airdrop.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Airdrop *AirdropTransactorSession) Receive() (*types.Transaction, error) {
	return _Airdrop.Contract.Receive(&_Airdrop.TransactOpts)
}

// AirdropAirdropBNBIterator is returned from FilterAirdropBNB and is used to iterate over the raw logs and unpacked data for AirdropBNB events raised by the Airdrop contract.
type AirdropAirdropBNBIterator struct {
	Event *AirdropAirdropBNB // Event containing the contract specifics and raw log

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
func (it *AirdropAirdropBNBIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AirdropAirdropBNB)
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
		it.Event = new(AirdropAirdropBNB)
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
func (it *AirdropAirdropBNBIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AirdropAirdropBNBIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AirdropAirdropBNB represents a AirdropBNB event raised by the Airdrop contract.
type AirdropAirdropBNB struct {
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAirdropBNB is a free log retrieval operation binding the contract event 0xad69676d22bff286c2c75292a594cdb796eb0183314fb9892ab1c9182823a348.
//
// Solidity: event AirdropBNB(address indexed recipient, uint256 amount)
func (_Airdrop *AirdropFilterer) FilterAirdropBNB(opts *bind.FilterOpts, recipient []common.Address) (*AirdropAirdropBNBIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Airdrop.contract.FilterLogs(opts, "AirdropBNB", recipientRule)
	if err != nil {
		return nil, err
	}
	return &AirdropAirdropBNBIterator{contract: _Airdrop.contract, event: "AirdropBNB", logs: logs, sub: sub}, nil
}

// WatchAirdropBNB is a free log subscription operation binding the contract event 0xad69676d22bff286c2c75292a594cdb796eb0183314fb9892ab1c9182823a348.
//
// Solidity: event AirdropBNB(address indexed recipient, uint256 amount)
func (_Airdrop *AirdropFilterer) WatchAirdropBNB(opts *bind.WatchOpts, sink chan<- *AirdropAirdropBNB, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Airdrop.contract.WatchLogs(opts, "AirdropBNB", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AirdropAirdropBNB)
				if err := _Airdrop.contract.UnpackLog(event, "AirdropBNB", log); err != nil {
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

// ParseAirdropBNB is a log parse operation binding the contract event 0xad69676d22bff286c2c75292a594cdb796eb0183314fb9892ab1c9182823a348.
//
// Solidity: event AirdropBNB(address indexed recipient, uint256 amount)
func (_Airdrop *AirdropFilterer) ParseAirdropBNB(log types.Log) (*AirdropAirdropBNB, error) {
	event := new(AirdropAirdropBNB)
	if err := _Airdrop.contract.UnpackLog(event, "AirdropBNB", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AirdropAirdropERC20Iterator is returned from FilterAirdropERC20 and is used to iterate over the raw logs and unpacked data for AirdropERC20 events raised by the Airdrop contract.
type AirdropAirdropERC20Iterator struct {
	Event *AirdropAirdropERC20 // Event containing the contract specifics and raw log

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
func (it *AirdropAirdropERC20Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AirdropAirdropERC20)
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
		it.Event = new(AirdropAirdropERC20)
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
func (it *AirdropAirdropERC20Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AirdropAirdropERC20Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AirdropAirdropERC20 represents a AirdropERC20 event raised by the Airdrop contract.
type AirdropAirdropERC20 struct {
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAirdropERC20 is a free log retrieval operation binding the contract event 0x9df1485392d28cf778b2fa04dcf0a936cb7703bea0ed28dbe1d0b17f728fa5de.
//
// Solidity: event AirdropERC20(address indexed recipient, uint256 amount)
func (_Airdrop *AirdropFilterer) FilterAirdropERC20(opts *bind.FilterOpts, recipient []common.Address) (*AirdropAirdropERC20Iterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Airdrop.contract.FilterLogs(opts, "AirdropERC20", recipientRule)
	if err != nil {
		return nil, err
	}
	return &AirdropAirdropERC20Iterator{contract: _Airdrop.contract, event: "AirdropERC20", logs: logs, sub: sub}, nil
}

// WatchAirdropERC20 is a free log subscription operation binding the contract event 0x9df1485392d28cf778b2fa04dcf0a936cb7703bea0ed28dbe1d0b17f728fa5de.
//
// Solidity: event AirdropERC20(address indexed recipient, uint256 amount)
func (_Airdrop *AirdropFilterer) WatchAirdropERC20(opts *bind.WatchOpts, sink chan<- *AirdropAirdropERC20, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Airdrop.contract.WatchLogs(opts, "AirdropERC20", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AirdropAirdropERC20)
				if err := _Airdrop.contract.UnpackLog(event, "AirdropERC20", log); err != nil {
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

// ParseAirdropERC20 is a log parse operation binding the contract event 0x9df1485392d28cf778b2fa04dcf0a936cb7703bea0ed28dbe1d0b17f728fa5de.
//
// Solidity: event AirdropERC20(address indexed recipient, uint256 amount)
func (_Airdrop *AirdropFilterer) ParseAirdropERC20(log types.Log) (*AirdropAirdropERC20, error) {
	event := new(AirdropAirdropERC20)
	if err := _Airdrop.contract.UnpackLog(event, "AirdropERC20", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
