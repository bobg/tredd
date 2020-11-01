// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// TreddProofStep is an auto generated low-level Go binding around an user-defined struct.
type TreddProofStep struct {
	H    []byte
	Left bool
}

// ERC20ABI is the input ABI used to generate the binding from.
const ERC20ABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ERC20FuncSigs maps the 4-byte function signature to its string representation.
var ERC20FuncSigs = map[string]string{
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
}

// ERC20 is an auto generated Go binding around an Ethereum contract.
type ERC20 struct {
	ERC20Caller     // Read-only binding to the contract
	ERC20Transactor // Write-only binding to the contract
	ERC20Filterer   // Log filterer for contract events
}

// ERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20Session struct {
	Contract     *ERC20            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20CallerSession struct {
	Contract *ERC20Caller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20TransactorSession struct {
	Contract     *ERC20Transactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20Raw struct {
	Contract *ERC20 // Generic contract binding to access the raw methods on
}

// ERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20CallerRaw struct {
	Contract *ERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20TransactorRaw struct {
	Contract *ERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20 creates a new instance of ERC20, bound to a specific deployed contract.
func NewERC20(address common.Address, backend bind.ContractBackend) (*ERC20, error) {
	contract, err := bindERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20{ERC20Caller: ERC20Caller{contract: contract}, ERC20Transactor: ERC20Transactor{contract: contract}, ERC20Filterer: ERC20Filterer{contract: contract}}, nil
}

// NewERC20Caller creates a new read-only instance of ERC20, bound to a specific deployed contract.
func NewERC20Caller(address common.Address, caller bind.ContractCaller) (*ERC20Caller, error) {
	contract, err := bindERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20Caller{contract: contract}, nil
}

// NewERC20Transactor creates a new write-only instance of ERC20, bound to a specific deployed contract.
func NewERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ERC20Transactor, error) {
	contract, err := bindERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20Transactor{contract: contract}, nil
}

// NewERC20Filterer creates a new log filterer instance of ERC20, bound to a specific deployed contract.
func NewERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ERC20Filterer, error) {
	contract, err := bindERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20Filterer{contract: contract}, nil
}

// bindERC20 binds a generic wrapper to an already deployed contract.
func bindERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20 *ERC20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20.Contract.ERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20 *ERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20.Contract.ERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20 *ERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20.Contract.ERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20 *ERC20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20 *ERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20 *ERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20 *ERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20 *ERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20 *ERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20 *ERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "balanceOf", account)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20 *ERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20 *ERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, account)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20 *ERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20 *ERC20Session) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20 *ERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ERC20 *ERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ERC20 *ERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ERC20 *ERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, sender, recipient, amount)
}

// ERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20 contract.
type ERC20ApprovalIterator struct {
	Event *ERC20Approval // Event containing the contract specifics and raw log

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
func (it *ERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20Approval)
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
		it.Event = new(ERC20Approval)
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
func (it *ERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20Approval represents a Approval event raised by the ERC20 contract.
type ERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20 *ERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20ApprovalIterator{contract: _ERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20 *ERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20Approval)
				if err := _ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20 *ERC20Filterer) ParseApproval(log types.Log) (*ERC20Approval, error) {
	event := new(ERC20Approval)
	if err := _ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20 contract.
type ERC20TransferIterator struct {
	Event *ERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20Transfer)
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
		it.Event = new(ERC20Transfer)
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
func (it *ERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20Transfer represents a Transfer event raised by the ERC20 contract.
type ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20 *ERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20TransferIterator{contract: _ERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20 *ERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20Transfer)
				if err := _ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20 *ERC20Filterer) ParseTransfer(log types.Log) (*ERC20Transfer, error) {
	event := new(ERC20Transfer)
	if err := _ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreddABI is the input ABI used to generate the binding from.
const TreddABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenType\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collateral\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"clearRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"cipherRoot\",\"type\":\"bytes32\"},{\"internalType\":\"int64\",\"name\":\"revealDeadline\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"refundDeadline\",\"type\":\"int64\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"decryptionKey\",\"type\":\"bytes32\"}],\"name\":\"evDecryptionKey\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"steps\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"leaf\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"want\",\"type\":\"bytes32\"}],\"name\":\"checkProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"steps\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"prefix\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"chunk\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"want\",\"type\":\"bytes32\"}],\"name\":\"checkProofWithPrefixedChunk\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"steps\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"prefix\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"want\",\"type\":\"bytes32\"}],\"name\":\"checkProofWithPrefixedHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"chunk\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"}],\"name\":\"decrypt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mBuyer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mCipherRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mClearRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mDecryptionKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRefundDeadline\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"\",\"type\":\"int64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRevealDeadline\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"\",\"type\":\"int64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRevealed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mSeller\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mTokenType\",\"outputs\":[{\"internalType\":\"contractERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paid\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"cipherChunk\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"clearHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"cipherProof\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"clearProof\",\"type\":\"tuple[]\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"decryptionKey\",\"type\":\"bytes32\"}],\"name\":\"reveal\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

// TreddFuncSigs maps the 4-byte function signature to its string representation.
var TreddFuncSigs = map[string]string{
	"ea8a1af0": "cancel()",
	"1235ffeb": "checkProof((bytes,bool)[],bytes,bytes32)",
	"33bbe2a7": "checkProofWithPrefixedChunk((bytes,bool)[],uint64,bytes,bytes32)",
	"fc6210c5": "checkProofWithPrefixedHash((bytes,bool)[],uint64,bytes32,bytes32)",
	"c7dea2f2": "claimPayment()",
	"a1598968": "decrypt(bytes,uint64)",
	"7d966e7d": "mAmount()",
	"649bfb36": "mBuyer()",
	"1d595ee7": "mCipherRoot()",
	"21b0ae82": "mClearRoot()",
	"095e4c20": "mCollateral()",
	"9067c7a9": "mDecryptionKey()",
	"2df6a9da": "mRefundDeadline()",
	"61a5ab22": "mRevealDeadline()",
	"54b53436": "mRevealed()",
	"8bae87ba": "mSeller()",
	"0c590dce": "mTokenType()",
	"295b4e17": "paid()",
	"ac280f3d": "refund(uint64,bytes,bytes32,(bytes,bool)[],(bytes,bool)[])",
	"701fd0f1": "reveal(bytes32)",
}

// TreddBin is the compiled bytecode used for deploying new contracts.
var TreddBin = "0x60806040526040516200148e3803806200148e8339810160408190526200002691620000f9565b600080546001600160a01b03199081163317909155600180546001600160a01b039a8b1690831617905560028054989099169716969096179096556003939093556004919091556005556006556007805491810b6001600160401b039081166801000000000000000002600160401b600160801b031994830b9091166001600160401b031990931692909217929092161790556009805460ff191690556200017c565b80516001600160a01b0381168114620000e157600080fd5b919050565b8051600781900b8114620000e157600080fd5b600080600080600080600080610100898b03121562000116578384fd5b6200012189620000c9565b97506200013160208a01620000c9565b965060408901519550606089015194506080890151935060a089015192506200015d60c08a01620000e6565b91506200016d60e08a01620000e6565b90509295985092959890939650565b611302806200018c6000396000f3fe6080604052600436106101235760003560e01c8063649bfb36116100a0578063a159896811610064578063a1598968146102bd578063ac280f3d146102ea578063c7dea2f21461030a578063ea8a1af01461031f578063fc6210c5146103345761012a565b8063649bfb3614610254578063701fd0f1146102695780637d966e7d1461027e5780638bae87ba146102935780639067c7a9146102a85761012a565b8063295b4e17116100e7578063295b4e17146101d35780632df6a9da146101e857806333bbe2a71461020a57806354b534361461022a57806361a5ab221461023f5761012a565b8063095e4c201461012f5780630c590dce1461015a5780631235ffeb1461017c5780631d595ee7146101a957806321b0ae82146101be5761012a565b3661012a57005b600080fd5b34801561013b57600080fd5b50610144610354565b604051610151919061121e565b60405180910390f35b34801561016657600080fd5b5061016f61035a565b60405161015191906111c2565b34801561018857600080fd5b5061019c610197366004610e1b565b610369565b6040516101519190611213565b3480156101b557600080fd5b5061014461052d565b3480156101ca57600080fd5b50610144610533565b3480156101df57600080fd5b50610144610539565b3480156101f457600080fd5b506101fd6105d5565b604051610151919061125a565b34801561021657600080fd5b5061019c610225366004610ede565b6105e5565b34801561023657600080fd5b5061019c61061c565b34801561024b57600080fd5b506101fd610625565b34801561026057600080fd5b5061016f61062e565b61027c610277366004610f79565b61063d565b005b34801561028a57600080fd5b5061014461077a565b34801561029f57600080fd5b5061016f610780565b3480156102b457600080fd5b5061014461078f565b3480156102c957600080fd5b506102dd6102d8366004610fa9565b610795565b6040516101519190611227565b3480156102f657600080fd5b5061027c610305366004610ff4565b610931565b34801561031657600080fd5b5061027c610ab1565b34801561032b57600080fd5b5061027c610b2d565b34801561034057600080fd5b5061019c61034f366004610e84565b610c8b565b60045481565b6002546001600160a01b031681565b6040516000908190600160f81b90829060029061038c90839089906020016110c5565b60408051601f19818403018152908290526103a691611154565b602060405180830381855afa1580156103c3573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906103e69190610f91565b905060005b87518163ffffffff16101561052057610402610cb4565b888263ffffffff168151811061041457fe5b602002602001015190508060200151156104a25780516040516002916104419187919087906020016110f6565b60408051601f198184030181529082905261045b91611154565b602060405180830381855afa158015610478573d6000803e3d6000fd5b5050506040513d601f19601f8201168201806040525081019061049b9190610f91565b9250610517565b80516040516002916104ba9187918791602001611092565b60408051601f19818403018152908290526104d491611154565b602060405180830381855afa1580156104f1573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906105149190610f91565b92505b506001016103eb565b5090931495945050505050565b60065481565b60055481565b6000610543610ca4565b1561054f5750476105d2565b6002546040516370a0823160e01b81526001600160a01b03909116906370a082319061057f9030906004016111c2565b60206040518083038186803b15801561059757600080fd5b505afa1580156105ab573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105cf9190610f91565b90505b90565b60078054600160401b9004900b81565b60006106138585856040516020016105fe92919061118d565b60405160208183030381529060405284610369565b95945050505050565b60095460ff1681565b60078054900b81565b6000546001600160a01b031681565b6001546001600160a01b0316331461065457600080fd5b60078054810b900b421061066757600080fd5b60095460ff161561067757600080fd5b61067f610ca4565b156106985760045434101561069357600080fd5b61072d565b600254600154600480546040516323b872dd60e01b81526001600160a01b03948516946323b872dd946106d29491169230929091016111d6565b602060405180830381600087803b1580156106ec57600080fd5b505af1158015610700573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107249190610f56565b61072d57600080fd5b60088190556009805460ff191660011790556040517f34292d279a4eb74e15e8b454d2d45ea717fe4914773188f0540fd6fbe67db0819061076f90839061121e565b60405180910390a150565b60035481565b6001546001600160a01b031681565b60085481565b60608083516001600160401b03811180156107af57600080fd5b506040519080825280601f01601f1916602001820160405280156107da576020820181803683370190505b50905060005b8451816020026001600160401b031610156109295760008160200290506000600260085487856040516020016108189392919061112e565b60408051601f198184030181529082905261083291611154565b602060405180830381855afa15801561084f573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906108729190610f91565b905060005b60208163ffffffff1610801561089e575087518163ffffffff1684016001600160401b0316105b1561091e57818163ffffffff16602081106108b557fe5b1a60f81b888263ffffffff1685016001600160401b0316815181106108d657fe5b602001015160f81c60f81b18858263ffffffff1685016001600160401b0316815181106108ff57fe5b60200101906001600160f81b031916908160001a905350600101610877565b5050506001016107e0565b509392505050565b6000546001600160a01b0316331461094857600080fd5b60078054600160401b9004810b900b421061096257600080fd5b60095460ff1661097157600080fd5b61097f8286866006546105e5565b61098857600080fd5b610996818685600554610c8b565b61099f57600080fd5b60606109ab8587610795565b9050836002826040516109be9190611154565b602060405180830381855afa1580156109db573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906109fe9190610f91565b1415610a0957600080fd5b610a11610ca4565b610aae576002546000546004805460035460405163a9059cbb60e01b81526001600160a01b039586169563a9059cbb95610a53959116939290920191016111fa565b602060405180830381600087803b158015610a6d57600080fd5b505af1158015610a81573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610aa59190610f56565b610aae57600080fd5b33ff5b6001546001600160a01b03163314610ac857600080fd5b60078054600160401b9004810b900b421015610ae357600080fd5b610aeb610ca4565b610aae576002546001546004805460035460405163a9059cbb60e01b81526001600160a01b039586169563a9059cbb95610a53959116939290920191016111fa565b6000546001600160a01b03163314610b4457600080fd5b60078054810b900b421015610b5857600080fd5b60095460ff1615610b6857600080fd5b610b70610ca4565b610aae576002546040516370a0823160e01b81526000916001600160a01b0316906370a0823190610ba59030906004016111c2565b60206040518083038186803b158015610bbd57600080fd5b505afa158015610bd1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bf59190610f91565b90508015610c875760025460005460405163a9059cbb60e01b81526001600160a01b039283169263a9059cbb92610c339291169085906004016111fa565b602060405180830381600087803b158015610c4d57600080fd5b505af1158015610c61573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c859190610f56565b505b5033ff5b60006106138585856040516020016105fe929190611170565b6002546001600160a01b03161590565b60408051808201909152606081526000602082015290565b600082601f830112610cdc578081fd5b81356001600160401b0380821115610cf057fe5b6020610cff8182850201611268565b838152935080840185820160005b85811015610d8d5781358801604080601f19838d03011215610d2e57600080fd5b80518181018181108982111715610d4157fe5b82528287013588811115610d5457600080fd5b610d628d8983870101610d99565b8252509181013591610d73836112bb565b808701929092525083529183019190830190600101610d0d565b50505050505092915050565b600082601f830112610da9578081fd5b81356001600160401b03811115610dbc57fe5b610dcf601f8201601f1916602001611268565b9150808252836020828501011115610de657600080fd5b8060208401602084013760009082016020015292915050565b80356001600160401b0381168114610e1657600080fd5b919050565b600080600060608486031215610e2f578283fd5b83356001600160401b0380821115610e45578485fd5b610e5187838801610ccc565b94506020860135915080821115610e66578384fd5b50610e7386828701610d99565b925050604084013590509250925092565b60008060008060808587031215610e99578081fd5b84356001600160401b03811115610eae578182fd5b610eba87828801610ccc565b945050610ec960208601610dff565b93969395505050506040820135916060013590565b60008060008060808587031215610ef3578384fd5b84356001600160401b0380821115610f09578586fd5b610f1588838901610ccc565b9550610f2360208801610dff565b94506040870135915080821115610f38578384fd5b50610f4587828801610d99565b949793965093946060013593505050565b600060208284031215610f67578081fd5b8151610f72816112bb565b9392505050565b600060208284031215610f8a578081fd5b5035919050565b600060208284031215610fa2578081fd5b5051919050565b60008060408385031215610fbb578182fd5b82356001600160401b03811115610fd0578283fd5b610fdc85828601610d99565b925050610feb60208401610dff565b90509250929050565b600080600080600060a0868803121561100b578081fd5b61101486610dff565b945060208601356001600160401b038082111561102f578283fd5b61103b89838a01610d99565b9550604088013594506060880135915080821115611057578283fd5b61106389838a01610ccc565b93506080880135915080821115611078578283fd5b5061108588828901610ccc565b9150509295509295909350565b600060ff60f81b8516825283600183015282516110b681602185016020870161128b565b91909101602101949350505050565b6001600160f81b03198316815281516000906110e881600185016020870161128b565b919091016001019392505050565b6001600160f81b031984168152825160009061111981600185016020880161128b565b60019201918201929092526021019392505050565b9283526001600160c01b031960c092831b81166020850152911b16602882015260300190565b6000825161116681846020870161128b565b9190910192915050565b60c09290921b6001600160c01b0319168252600882015260280190565b60006001600160401b0360c01b8460c01b16825282516111b481600885016020870161128b565b919091016008019392505050565b6001600160a01b0391909116815260200190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6001600160a01b03929092168252602082015260400190565b901515815260200190565b90815260200190565b600060208252825180602084015261124681604085016020870161128b565b601f01601f19169190910160400192915050565b60079190910b815260200190565b6040518181016001600160401b038111828210171561128357fe5b604052919050565b60005b838110156112a657818101518382015260200161128e565b838111156112b5576000848401525b50505050565b80151581146112c957600080fd5b5056fea2646970667358221220c69f5c642f61be181b2adea27d13782cecf19e53e23039f7ed518ba34c7a484f64736f6c63430007040033"

// DeployTredd deploys a new Ethereum contract, binding an instance of Tredd to it.
func DeployTredd(auth *bind.TransactOpts, backend bind.ContractBackend, seller common.Address, tokenType common.Address, amount *big.Int, collateral *big.Int, clearRoot [32]byte, cipherRoot [32]byte, revealDeadline int64, refundDeadline int64) (common.Address, *types.Transaction, *Tredd, error) {
	parsed, err := abi.JSON(strings.NewReader(TreddABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TreddBin), backend, seller, tokenType, amount, collateral, clearRoot, cipherRoot, revealDeadline, refundDeadline)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Tredd{TreddCaller: TreddCaller{contract: contract}, TreddTransactor: TreddTransactor{contract: contract}, TreddFilterer: TreddFilterer{contract: contract}}, nil
}

// Tredd is an auto generated Go binding around an Ethereum contract.
type Tredd struct {
	TreddCaller     // Read-only binding to the contract
	TreddTransactor // Write-only binding to the contract
	TreddFilterer   // Log filterer for contract events
}

// TreddCaller is an auto generated read-only Go binding around an Ethereum contract.
type TreddCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TreddTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TreddTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TreddFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TreddFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TreddSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TreddSession struct {
	Contract     *Tredd            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TreddCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TreddCallerSession struct {
	Contract *TreddCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TreddTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TreddTransactorSession struct {
	Contract     *TreddTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TreddRaw is an auto generated low-level Go binding around an Ethereum contract.
type TreddRaw struct {
	Contract *Tredd // Generic contract binding to access the raw methods on
}

// TreddCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TreddCallerRaw struct {
	Contract *TreddCaller // Generic read-only contract binding to access the raw methods on
}

// TreddTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TreddTransactorRaw struct {
	Contract *TreddTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTredd creates a new instance of Tredd, bound to a specific deployed contract.
func NewTredd(address common.Address, backend bind.ContractBackend) (*Tredd, error) {
	contract, err := bindTredd(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Tredd{TreddCaller: TreddCaller{contract: contract}, TreddTransactor: TreddTransactor{contract: contract}, TreddFilterer: TreddFilterer{contract: contract}}, nil
}

// NewTreddCaller creates a new read-only instance of Tredd, bound to a specific deployed contract.
func NewTreddCaller(address common.Address, caller bind.ContractCaller) (*TreddCaller, error) {
	contract, err := bindTredd(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TreddCaller{contract: contract}, nil
}

// NewTreddTransactor creates a new write-only instance of Tredd, bound to a specific deployed contract.
func NewTreddTransactor(address common.Address, transactor bind.ContractTransactor) (*TreddTransactor, error) {
	contract, err := bindTredd(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TreddTransactor{contract: contract}, nil
}

// NewTreddFilterer creates a new log filterer instance of Tredd, bound to a specific deployed contract.
func NewTreddFilterer(address common.Address, filterer bind.ContractFilterer) (*TreddFilterer, error) {
	contract, err := bindTredd(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TreddFilterer{contract: contract}, nil
}

// bindTredd binds a generic wrapper to an already deployed contract.
func bindTredd(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TreddABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Tredd *TreddRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Tredd.Contract.TreddCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Tredd *TreddRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Tredd.Contract.TreddTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Tredd *TreddRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Tredd.Contract.TreddTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Tredd *TreddCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Tredd.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Tredd *TreddTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Tredd.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Tredd *TreddTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Tredd.Contract.contract.Transact(opts, method, params...)
}

// CheckProof is a free data retrieval call binding the contract method 0x1235ffeb.
//
// Solidity: function checkProof((bytes,bool)[] steps, bytes leaf, bytes32 want) pure returns(bool)
func (_Tredd *TreddCaller) CheckProof(opts *bind.CallOpts, steps []TreddProofStep, leaf []byte, want [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "checkProof", steps, leaf, want)
	return *ret0, err
}

// CheckProof is a free data retrieval call binding the contract method 0x1235ffeb.
//
// Solidity: function checkProof((bytes,bool)[] steps, bytes leaf, bytes32 want) pure returns(bool)
func (_Tredd *TreddSession) CheckProof(steps []TreddProofStep, leaf []byte, want [32]byte) (bool, error) {
	return _Tredd.Contract.CheckProof(&_Tredd.CallOpts, steps, leaf, want)
}

// CheckProof is a free data retrieval call binding the contract method 0x1235ffeb.
//
// Solidity: function checkProof((bytes,bool)[] steps, bytes leaf, bytes32 want) pure returns(bool)
func (_Tredd *TreddCallerSession) CheckProof(steps []TreddProofStep, leaf []byte, want [32]byte) (bool, error) {
	return _Tredd.Contract.CheckProof(&_Tredd.CallOpts, steps, leaf, want)
}

// CheckProofWithPrefixedChunk is a free data retrieval call binding the contract method 0x33bbe2a7.
//
// Solidity: function checkProofWithPrefixedChunk((bytes,bool)[] steps, uint64 prefix, bytes chunk, bytes32 want) pure returns(bool)
func (_Tredd *TreddCaller) CheckProofWithPrefixedChunk(opts *bind.CallOpts, steps []TreddProofStep, prefix uint64, chunk []byte, want [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "checkProofWithPrefixedChunk", steps, prefix, chunk, want)
	return *ret0, err
}

// CheckProofWithPrefixedChunk is a free data retrieval call binding the contract method 0x33bbe2a7.
//
// Solidity: function checkProofWithPrefixedChunk((bytes,bool)[] steps, uint64 prefix, bytes chunk, bytes32 want) pure returns(bool)
func (_Tredd *TreddSession) CheckProofWithPrefixedChunk(steps []TreddProofStep, prefix uint64, chunk []byte, want [32]byte) (bool, error) {
	return _Tredd.Contract.CheckProofWithPrefixedChunk(&_Tredd.CallOpts, steps, prefix, chunk, want)
}

// CheckProofWithPrefixedChunk is a free data retrieval call binding the contract method 0x33bbe2a7.
//
// Solidity: function checkProofWithPrefixedChunk((bytes,bool)[] steps, uint64 prefix, bytes chunk, bytes32 want) pure returns(bool)
func (_Tredd *TreddCallerSession) CheckProofWithPrefixedChunk(steps []TreddProofStep, prefix uint64, chunk []byte, want [32]byte) (bool, error) {
	return _Tredd.Contract.CheckProofWithPrefixedChunk(&_Tredd.CallOpts, steps, prefix, chunk, want)
}

// CheckProofWithPrefixedHash is a free data retrieval call binding the contract method 0xfc6210c5.
//
// Solidity: function checkProofWithPrefixedHash((bytes,bool)[] steps, uint64 prefix, bytes32 hash, bytes32 want) pure returns(bool)
func (_Tredd *TreddCaller) CheckProofWithPrefixedHash(opts *bind.CallOpts, steps []TreddProofStep, prefix uint64, hash [32]byte, want [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "checkProofWithPrefixedHash", steps, prefix, hash, want)
	return *ret0, err
}

// CheckProofWithPrefixedHash is a free data retrieval call binding the contract method 0xfc6210c5.
//
// Solidity: function checkProofWithPrefixedHash((bytes,bool)[] steps, uint64 prefix, bytes32 hash, bytes32 want) pure returns(bool)
func (_Tredd *TreddSession) CheckProofWithPrefixedHash(steps []TreddProofStep, prefix uint64, hash [32]byte, want [32]byte) (bool, error) {
	return _Tredd.Contract.CheckProofWithPrefixedHash(&_Tredd.CallOpts, steps, prefix, hash, want)
}

// CheckProofWithPrefixedHash is a free data retrieval call binding the contract method 0xfc6210c5.
//
// Solidity: function checkProofWithPrefixedHash((bytes,bool)[] steps, uint64 prefix, bytes32 hash, bytes32 want) pure returns(bool)
func (_Tredd *TreddCallerSession) CheckProofWithPrefixedHash(steps []TreddProofStep, prefix uint64, hash [32]byte, want [32]byte) (bool, error) {
	return _Tredd.Contract.CheckProofWithPrefixedHash(&_Tredd.CallOpts, steps, prefix, hash, want)
}

// Decrypt is a free data retrieval call binding the contract method 0xa1598968.
//
// Solidity: function decrypt(bytes chunk, uint64 index) view returns(bytes)
func (_Tredd *TreddCaller) Decrypt(opts *bind.CallOpts, chunk []byte, index uint64) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "decrypt", chunk, index)
	return *ret0, err
}

// Decrypt is a free data retrieval call binding the contract method 0xa1598968.
//
// Solidity: function decrypt(bytes chunk, uint64 index) view returns(bytes)
func (_Tredd *TreddSession) Decrypt(chunk []byte, index uint64) ([]byte, error) {
	return _Tredd.Contract.Decrypt(&_Tredd.CallOpts, chunk, index)
}

// Decrypt is a free data retrieval call binding the contract method 0xa1598968.
//
// Solidity: function decrypt(bytes chunk, uint64 index) view returns(bytes)
func (_Tredd *TreddCallerSession) Decrypt(chunk []byte, index uint64) ([]byte, error) {
	return _Tredd.Contract.Decrypt(&_Tredd.CallOpts, chunk, index)
}

// MAmount is a free data retrieval call binding the contract method 0x7d966e7d.
//
// Solidity: function mAmount() view returns(uint256)
func (_Tredd *TreddCaller) MAmount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mAmount")
	return *ret0, err
}

// MAmount is a free data retrieval call binding the contract method 0x7d966e7d.
//
// Solidity: function mAmount() view returns(uint256)
func (_Tredd *TreddSession) MAmount() (*big.Int, error) {
	return _Tredd.Contract.MAmount(&_Tredd.CallOpts)
}

// MAmount is a free data retrieval call binding the contract method 0x7d966e7d.
//
// Solidity: function mAmount() view returns(uint256)
func (_Tredd *TreddCallerSession) MAmount() (*big.Int, error) {
	return _Tredd.Contract.MAmount(&_Tredd.CallOpts)
}

// MBuyer is a free data retrieval call binding the contract method 0x649bfb36.
//
// Solidity: function mBuyer() view returns(address)
func (_Tredd *TreddCaller) MBuyer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mBuyer")
	return *ret0, err
}

// MBuyer is a free data retrieval call binding the contract method 0x649bfb36.
//
// Solidity: function mBuyer() view returns(address)
func (_Tredd *TreddSession) MBuyer() (common.Address, error) {
	return _Tredd.Contract.MBuyer(&_Tredd.CallOpts)
}

// MBuyer is a free data retrieval call binding the contract method 0x649bfb36.
//
// Solidity: function mBuyer() view returns(address)
func (_Tredd *TreddCallerSession) MBuyer() (common.Address, error) {
	return _Tredd.Contract.MBuyer(&_Tredd.CallOpts)
}

// MCipherRoot is a free data retrieval call binding the contract method 0x1d595ee7.
//
// Solidity: function mCipherRoot() view returns(bytes32)
func (_Tredd *TreddCaller) MCipherRoot(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mCipherRoot")
	return *ret0, err
}

// MCipherRoot is a free data retrieval call binding the contract method 0x1d595ee7.
//
// Solidity: function mCipherRoot() view returns(bytes32)
func (_Tredd *TreddSession) MCipherRoot() ([32]byte, error) {
	return _Tredd.Contract.MCipherRoot(&_Tredd.CallOpts)
}

// MCipherRoot is a free data retrieval call binding the contract method 0x1d595ee7.
//
// Solidity: function mCipherRoot() view returns(bytes32)
func (_Tredd *TreddCallerSession) MCipherRoot() ([32]byte, error) {
	return _Tredd.Contract.MCipherRoot(&_Tredd.CallOpts)
}

// MClearRoot is a free data retrieval call binding the contract method 0x21b0ae82.
//
// Solidity: function mClearRoot() view returns(bytes32)
func (_Tredd *TreddCaller) MClearRoot(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mClearRoot")
	return *ret0, err
}

// MClearRoot is a free data retrieval call binding the contract method 0x21b0ae82.
//
// Solidity: function mClearRoot() view returns(bytes32)
func (_Tredd *TreddSession) MClearRoot() ([32]byte, error) {
	return _Tredd.Contract.MClearRoot(&_Tredd.CallOpts)
}

// MClearRoot is a free data retrieval call binding the contract method 0x21b0ae82.
//
// Solidity: function mClearRoot() view returns(bytes32)
func (_Tredd *TreddCallerSession) MClearRoot() ([32]byte, error) {
	return _Tredd.Contract.MClearRoot(&_Tredd.CallOpts)
}

// MCollateral is a free data retrieval call binding the contract method 0x095e4c20.
//
// Solidity: function mCollateral() view returns(uint256)
func (_Tredd *TreddCaller) MCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mCollateral")
	return *ret0, err
}

// MCollateral is a free data retrieval call binding the contract method 0x095e4c20.
//
// Solidity: function mCollateral() view returns(uint256)
func (_Tredd *TreddSession) MCollateral() (*big.Int, error) {
	return _Tredd.Contract.MCollateral(&_Tredd.CallOpts)
}

// MCollateral is a free data retrieval call binding the contract method 0x095e4c20.
//
// Solidity: function mCollateral() view returns(uint256)
func (_Tredd *TreddCallerSession) MCollateral() (*big.Int, error) {
	return _Tredd.Contract.MCollateral(&_Tredd.CallOpts)
}

// MDecryptionKey is a free data retrieval call binding the contract method 0x9067c7a9.
//
// Solidity: function mDecryptionKey() view returns(bytes32)
func (_Tredd *TreddCaller) MDecryptionKey(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mDecryptionKey")
	return *ret0, err
}

// MDecryptionKey is a free data retrieval call binding the contract method 0x9067c7a9.
//
// Solidity: function mDecryptionKey() view returns(bytes32)
func (_Tredd *TreddSession) MDecryptionKey() ([32]byte, error) {
	return _Tredd.Contract.MDecryptionKey(&_Tredd.CallOpts)
}

// MDecryptionKey is a free data retrieval call binding the contract method 0x9067c7a9.
//
// Solidity: function mDecryptionKey() view returns(bytes32)
func (_Tredd *TreddCallerSession) MDecryptionKey() ([32]byte, error) {
	return _Tredd.Contract.MDecryptionKey(&_Tredd.CallOpts)
}

// MRefundDeadline is a free data retrieval call binding the contract method 0x2df6a9da.
//
// Solidity: function mRefundDeadline() view returns(int64)
func (_Tredd *TreddCaller) MRefundDeadline(opts *bind.CallOpts) (int64, error) {
	var (
		ret0 = new(int64)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mRefundDeadline")
	return *ret0, err
}

// MRefundDeadline is a free data retrieval call binding the contract method 0x2df6a9da.
//
// Solidity: function mRefundDeadline() view returns(int64)
func (_Tredd *TreddSession) MRefundDeadline() (int64, error) {
	return _Tredd.Contract.MRefundDeadline(&_Tredd.CallOpts)
}

// MRefundDeadline is a free data retrieval call binding the contract method 0x2df6a9da.
//
// Solidity: function mRefundDeadline() view returns(int64)
func (_Tredd *TreddCallerSession) MRefundDeadline() (int64, error) {
	return _Tredd.Contract.MRefundDeadline(&_Tredd.CallOpts)
}

// MRevealDeadline is a free data retrieval call binding the contract method 0x61a5ab22.
//
// Solidity: function mRevealDeadline() view returns(int64)
func (_Tredd *TreddCaller) MRevealDeadline(opts *bind.CallOpts) (int64, error) {
	var (
		ret0 = new(int64)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mRevealDeadline")
	return *ret0, err
}

// MRevealDeadline is a free data retrieval call binding the contract method 0x61a5ab22.
//
// Solidity: function mRevealDeadline() view returns(int64)
func (_Tredd *TreddSession) MRevealDeadline() (int64, error) {
	return _Tredd.Contract.MRevealDeadline(&_Tredd.CallOpts)
}

// MRevealDeadline is a free data retrieval call binding the contract method 0x61a5ab22.
//
// Solidity: function mRevealDeadline() view returns(int64)
func (_Tredd *TreddCallerSession) MRevealDeadline() (int64, error) {
	return _Tredd.Contract.MRevealDeadline(&_Tredd.CallOpts)
}

// MRevealed is a free data retrieval call binding the contract method 0x54b53436.
//
// Solidity: function mRevealed() view returns(bool)
func (_Tredd *TreddCaller) MRevealed(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mRevealed")
	return *ret0, err
}

// MRevealed is a free data retrieval call binding the contract method 0x54b53436.
//
// Solidity: function mRevealed() view returns(bool)
func (_Tredd *TreddSession) MRevealed() (bool, error) {
	return _Tredd.Contract.MRevealed(&_Tredd.CallOpts)
}

// MRevealed is a free data retrieval call binding the contract method 0x54b53436.
//
// Solidity: function mRevealed() view returns(bool)
func (_Tredd *TreddCallerSession) MRevealed() (bool, error) {
	return _Tredd.Contract.MRevealed(&_Tredd.CallOpts)
}

// MSeller is a free data retrieval call binding the contract method 0x8bae87ba.
//
// Solidity: function mSeller() view returns(address)
func (_Tredd *TreddCaller) MSeller(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mSeller")
	return *ret0, err
}

// MSeller is a free data retrieval call binding the contract method 0x8bae87ba.
//
// Solidity: function mSeller() view returns(address)
func (_Tredd *TreddSession) MSeller() (common.Address, error) {
	return _Tredd.Contract.MSeller(&_Tredd.CallOpts)
}

// MSeller is a free data retrieval call binding the contract method 0x8bae87ba.
//
// Solidity: function mSeller() view returns(address)
func (_Tredd *TreddCallerSession) MSeller() (common.Address, error) {
	return _Tredd.Contract.MSeller(&_Tredd.CallOpts)
}

// MTokenType is a free data retrieval call binding the contract method 0x0c590dce.
//
// Solidity: function mTokenType() view returns(address)
func (_Tredd *TreddCaller) MTokenType(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mTokenType")
	return *ret0, err
}

// MTokenType is a free data retrieval call binding the contract method 0x0c590dce.
//
// Solidity: function mTokenType() view returns(address)
func (_Tredd *TreddSession) MTokenType() (common.Address, error) {
	return _Tredd.Contract.MTokenType(&_Tredd.CallOpts)
}

// MTokenType is a free data retrieval call binding the contract method 0x0c590dce.
//
// Solidity: function mTokenType() view returns(address)
func (_Tredd *TreddCallerSession) MTokenType() (common.Address, error) {
	return _Tredd.Contract.MTokenType(&_Tredd.CallOpts)
}

// Paid is a free data retrieval call binding the contract method 0x295b4e17.
//
// Solidity: function paid() view returns(uint256)
func (_Tredd *TreddCaller) Paid(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "paid")
	return *ret0, err
}

// Paid is a free data retrieval call binding the contract method 0x295b4e17.
//
// Solidity: function paid() view returns(uint256)
func (_Tredd *TreddSession) Paid() (*big.Int, error) {
	return _Tredd.Contract.Paid(&_Tredd.CallOpts)
}

// Paid is a free data retrieval call binding the contract method 0x295b4e17.
//
// Solidity: function paid() view returns(uint256)
func (_Tredd *TreddCallerSession) Paid() (*big.Int, error) {
	return _Tredd.Contract.Paid(&_Tredd.CallOpts)
}

// Cancel is a paid mutator transaction binding the contract method 0xea8a1af0.
//
// Solidity: function cancel() returns()
func (_Tredd *TreddTransactor) Cancel(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Tredd.contract.Transact(opts, "cancel")
}

// Cancel is a paid mutator transaction binding the contract method 0xea8a1af0.
//
// Solidity: function cancel() returns()
func (_Tredd *TreddSession) Cancel() (*types.Transaction, error) {
	return _Tredd.Contract.Cancel(&_Tredd.TransactOpts)
}

// Cancel is a paid mutator transaction binding the contract method 0xea8a1af0.
//
// Solidity: function cancel() returns()
func (_Tredd *TreddTransactorSession) Cancel() (*types.Transaction, error) {
	return _Tredd.Contract.Cancel(&_Tredd.TransactOpts)
}

// ClaimPayment is a paid mutator transaction binding the contract method 0xc7dea2f2.
//
// Solidity: function claimPayment() returns()
func (_Tredd *TreddTransactor) ClaimPayment(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Tredd.contract.Transact(opts, "claimPayment")
}

// ClaimPayment is a paid mutator transaction binding the contract method 0xc7dea2f2.
//
// Solidity: function claimPayment() returns()
func (_Tredd *TreddSession) ClaimPayment() (*types.Transaction, error) {
	return _Tredd.Contract.ClaimPayment(&_Tredd.TransactOpts)
}

// ClaimPayment is a paid mutator transaction binding the contract method 0xc7dea2f2.
//
// Solidity: function claimPayment() returns()
func (_Tredd *TreddTransactorSession) ClaimPayment() (*types.Transaction, error) {
	return _Tredd.Contract.ClaimPayment(&_Tredd.TransactOpts)
}

// Refund is a paid mutator transaction binding the contract method 0xac280f3d.
//
// Solidity: function refund(uint64 index, bytes cipherChunk, bytes32 clearHash, (bytes,bool)[] cipherProof, (bytes,bool)[] clearProof) returns()
func (_Tredd *TreddTransactor) Refund(opts *bind.TransactOpts, index uint64, cipherChunk []byte, clearHash [32]byte, cipherProof []TreddProofStep, clearProof []TreddProofStep) (*types.Transaction, error) {
	return _Tredd.contract.Transact(opts, "refund", index, cipherChunk, clearHash, cipherProof, clearProof)
}

// Refund is a paid mutator transaction binding the contract method 0xac280f3d.
//
// Solidity: function refund(uint64 index, bytes cipherChunk, bytes32 clearHash, (bytes,bool)[] cipherProof, (bytes,bool)[] clearProof) returns()
func (_Tredd *TreddSession) Refund(index uint64, cipherChunk []byte, clearHash [32]byte, cipherProof []TreddProofStep, clearProof []TreddProofStep) (*types.Transaction, error) {
	return _Tredd.Contract.Refund(&_Tredd.TransactOpts, index, cipherChunk, clearHash, cipherProof, clearProof)
}

// Refund is a paid mutator transaction binding the contract method 0xac280f3d.
//
// Solidity: function refund(uint64 index, bytes cipherChunk, bytes32 clearHash, (bytes,bool)[] cipherProof, (bytes,bool)[] clearProof) returns()
func (_Tredd *TreddTransactorSession) Refund(index uint64, cipherChunk []byte, clearHash [32]byte, cipherProof []TreddProofStep, clearProof []TreddProofStep) (*types.Transaction, error) {
	return _Tredd.Contract.Refund(&_Tredd.TransactOpts, index, cipherChunk, clearHash, cipherProof, clearProof)
}

// Reveal is a paid mutator transaction binding the contract method 0x701fd0f1.
//
// Solidity: function reveal(bytes32 decryptionKey) payable returns()
func (_Tredd *TreddTransactor) Reveal(opts *bind.TransactOpts, decryptionKey [32]byte) (*types.Transaction, error) {
	return _Tredd.contract.Transact(opts, "reveal", decryptionKey)
}

// Reveal is a paid mutator transaction binding the contract method 0x701fd0f1.
//
// Solidity: function reveal(bytes32 decryptionKey) payable returns()
func (_Tredd *TreddSession) Reveal(decryptionKey [32]byte) (*types.Transaction, error) {
	return _Tredd.Contract.Reveal(&_Tredd.TransactOpts, decryptionKey)
}

// Reveal is a paid mutator transaction binding the contract method 0x701fd0f1.
//
// Solidity: function reveal(bytes32 decryptionKey) payable returns()
func (_Tredd *TreddTransactorSession) Reveal(decryptionKey [32]byte) (*types.Transaction, error) {
	return _Tredd.Contract.Reveal(&_Tredd.TransactOpts, decryptionKey)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Tredd *TreddTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Tredd.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Tredd *TreddSession) Receive() (*types.Transaction, error) {
	return _Tredd.Contract.Receive(&_Tredd.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Tredd *TreddTransactorSession) Receive() (*types.Transaction, error) {
	return _Tredd.Contract.Receive(&_Tredd.TransactOpts)
}

// TreddEvDecryptionKeyIterator is returned from FilterEvDecryptionKey and is used to iterate over the raw logs and unpacked data for EvDecryptionKey events raised by the Tredd contract.
type TreddEvDecryptionKeyIterator struct {
	Event *TreddEvDecryptionKey // Event containing the contract specifics and raw log

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
func (it *TreddEvDecryptionKeyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreddEvDecryptionKey)
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
		it.Event = new(TreddEvDecryptionKey)
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
func (it *TreddEvDecryptionKeyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreddEvDecryptionKeyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreddEvDecryptionKey represents a EvDecryptionKey event raised by the Tredd contract.
type TreddEvDecryptionKey struct {
	DecryptionKey [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterEvDecryptionKey is a free log retrieval operation binding the contract event 0x34292d279a4eb74e15e8b454d2d45ea717fe4914773188f0540fd6fbe67db081.
//
// Solidity: event evDecryptionKey(bytes32 decryptionKey)
func (_Tredd *TreddFilterer) FilterEvDecryptionKey(opts *bind.FilterOpts) (*TreddEvDecryptionKeyIterator, error) {

	logs, sub, err := _Tredd.contract.FilterLogs(opts, "evDecryptionKey")
	if err != nil {
		return nil, err
	}
	return &TreddEvDecryptionKeyIterator{contract: _Tredd.contract, event: "evDecryptionKey", logs: logs, sub: sub}, nil
}

// WatchEvDecryptionKey is a free log subscription operation binding the contract event 0x34292d279a4eb74e15e8b454d2d45ea717fe4914773188f0540fd6fbe67db081.
//
// Solidity: event evDecryptionKey(bytes32 decryptionKey)
func (_Tredd *TreddFilterer) WatchEvDecryptionKey(opts *bind.WatchOpts, sink chan<- *TreddEvDecryptionKey) (event.Subscription, error) {

	logs, sub, err := _Tredd.contract.WatchLogs(opts, "evDecryptionKey")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreddEvDecryptionKey)
				if err := _Tredd.contract.UnpackLog(event, "evDecryptionKey", log); err != nil {
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

// ParseEvDecryptionKey is a log parse operation binding the contract event 0x34292d279a4eb74e15e8b454d2d45ea717fe4914773188f0540fd6fbe67db081.
//
// Solidity: event evDecryptionKey(bytes32 decryptionKey)
func (_Tredd *TreddFilterer) ParseEvDecryptionKey(log types.Log) (*TreddEvDecryptionKey, error) {
	event := new(TreddEvDecryptionKey)
	if err := _Tredd.contract.UnpackLog(event, "evDecryptionKey", log); err != nil {
		return nil, err
	}
	return event, nil
}
