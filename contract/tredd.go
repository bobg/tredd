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
const TreddABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenType\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collateral\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"clearRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"cipherRoot\",\"type\":\"bytes32\"},{\"internalType\":\"int64\",\"name\":\"revealDeadline\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"refundDeadline\",\"type\":\"int64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"decryptionKey\",\"type\":\"bytes32\"}],\"name\":\"evDecryptionKey\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"steps\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32\",\"name\":\"leaf\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"want\",\"type\":\"bytes32\"}],\"name\":\"checkProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"chunk\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"}],\"name\":\"decrypt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mBuyer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mCipherRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mClearRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mDecryptionKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRefundDeadline\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"\",\"type\":\"int64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRevealDeadline\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"\",\"type\":\"int64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRevealed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mSeller\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mTokenType\",\"outputs\":[{\"internalType\":\"contractERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paid\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"cipherChunk\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"clearHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"cipherProof\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"clearProof\",\"type\":\"tuple[]\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"decryptionKey\",\"type\":\"bytes32\"}],\"name\":\"reveal\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

// TreddFuncSigs maps the 4-byte function signature to its string representation.
var TreddFuncSigs = map[string]string{
	"ea8a1af0": "cancel()",
	"8343f445": "checkProof((bytes,bool)[],bytes32,bytes32)",
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
var TreddBin = "0x60806040523480156200001157600080fd5b50604051620012ec380380620012ec833981016040819052620000349162000107565b600080546001600160a01b03199081163317909155600180546001600160a01b039a8b1690831617905560028054989099169716969096179096556003939093556004919091556005556006556007805491810b6001600160401b039081166801000000000000000002600160401b600160801b031994830b9091166001600160401b031990931692909217929092161790556009805460ff191690556200018a565b80516001600160a01b0381168114620000ef57600080fd5b919050565b8051600781900b8114620000ef57600080fd5b600080600080600080600080610100898b03121562000124578384fd5b6200012f89620000d7565b97506200013f60208a01620000d7565b965060408901519550606089015194506080890151935060a089015192506200016b60c08a01620000f4565b91506200017b60e08a01620000f4565b90509295985092959890939650565b611152806200019a6000396000f3fe60806040526004361061010d5760003560e01c8063701fd0f1116100955780639067c7a9116100645780639067c7a914610272578063a159896814610287578063ac280f3d146102b4578063c7dea2f2146102d4578063ea8a1af0146102e957610114565b8063701fd0f1146102135780637d966e7d146102285780638343f4451461023d5780638bae87ba1461025d57610114565b8063295b4e17116100dc578063295b4e17146101905780632df6a9da146101a557806354b53436146101c757806361a5ab22146101e9578063649bfb36146101fe57610114565b8063095e4c20146101195780630c590dce146101445780631d595ee71461016657806321b0ae821461017b57610114565b3661011457005b600080fd5b34801561012557600080fd5b5061012e6102fe565b60405161013b919061106d565b60405180910390f35b34801561015057600080fd5b50610159610304565b60405161013b9190611011565b34801561017257600080fd5b5061012e610313565b34801561018757600080fd5b5061012e610319565b34801561019c57600080fd5b5061012e61031f565b3480156101b157600080fd5b506101ba6103bb565b60405161013b91906110a9565b3480156101d357600080fd5b506101dc6103cb565b60405161013b9190611062565b3480156101f557600080fd5b506101ba6103d4565b34801561020a57600080fd5b506101596103dd565b610226610221366004610e03565b6103ec565b005b34801561023457600080fd5b5061012e610529565b34801561024957600080fd5b506101dc610258366004610d94565b61052f565b34801561026957600080fd5b5061015961067a565b34801561027e57600080fd5b5061012e610689565b34801561029357600080fd5b506102a76102a2366004610e33565b61068f565b60405161013b9190611076565b3480156102c057600080fd5b506102266102cf366004610e7f565b610830565b3480156102e057600080fd5b50610226610a40565b3480156102f557600080fd5b50610226610abc565b60045481565b6002546001600160a01b031681565b60065481565b60055481565b6000610329610c1a565b156103355750476103b8565b6002546040516370a0823160e01b81526001600160a01b03909116906370a0823190610365903090600401611011565b60206040518083038186803b15801561037d57600080fd5b505afa158015610391573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103b59190610e1b565b90505b90565b60078054600160401b9004900b81565b60095460ff1681565b60078054900b81565b6000546001600160a01b031681565b6001546001600160a01b0316331461040357600080fd5b60078054810b900b421061041657600080fd5b60095460ff161561042657600080fd5b61042e610c1a565b156104475760045434101561044257600080fd5b6104dc565b600254600154600480546040516323b872dd60e01b81526001600160a01b03948516946323b872dd94610481949116923092909101611025565b602060405180830381600087803b15801561049b57600080fd5b505af11580156104af573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104d39190610de0565b6104dc57600080fd5b60088190556009805460ff191660011790556040517f34292d279a4eb74e15e8b454d2d45ea717fe4914773188f0540fd6fbe67db0819061051e90839061106d565b60405180910390a150565b60035481565b6000600160f81b83825b86518163ffffffff16101561066e57610550610c2a565b878263ffffffff168151811061056257fe5b602002602001015190508060200151156105f057805160405160029161058f918791908790602001610f51565b60408051601f19818403018152908290526105a991610ff5565b602060405180830381855afa1580156105c6573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906105e99190610e1b565b9250610665565b80516040516002916106089187918791602001610f1e565b60408051601f198184030181529082905261062291610ff5565b602060405180830381855afa15801561063f573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906106629190610e1b565b92505b50600101610539565b50909214949350505050565b6001546001600160a01b031681565b60085481565b606080835167ffffffffffffffff811180156106aa57600080fd5b506040519080825280601f01601f1916602001820160405280156106d5576020820181803683370190505b50905060005b84518160200267ffffffffffffffff16101561082857600081602002905060006002600854878560405160200161071493929190610fcf565b60408051601f198184030181529082905261072e91610ff5565b602060405180830381855afa15801561074b573d6000803e3d6000fd5b5050506040513d601f19601f8201168201806040525081019061076e9190610e1b565b905060005b60208163ffffffff1610801561079b575087518163ffffffff16840167ffffffffffffffff16105b1561081d57818163ffffffff16602081106107b257fe5b1a60f81b888263ffffffff16850167ffffffffffffffff16815181106107d457fe5b602001015160f81c60f81b18858263ffffffff16850167ffffffffffffffff16815181106107fe57fe5b60200101906001600160f81b031916908160001a905350600101610773565b5050506001016106db565b509392505050565b6000546001600160a01b0316331461084757600080fd5b60078054600160401b9004810b900b421061086157600080fd5b60095460ff1661087057600080fd5b600080600282888860405160200161088a93929190610f89565b60408051601f19818403018152908290526108a491610ff5565b602060405180830381855afa1580156108c1573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906108e49190610e1b565b90506108f3848260065461052f565b6108fc57600080fd5b610909838660055461052f565b61091257600080fd5b84600283896109218a8c61068f565b60405160200161093393929190610f89565b60408051601f198184030181529082905261094d91610ff5565b602060405180830381855afa15801561096a573d6000803e3d6000fd5b5050506040513d601f19601f8201168201806040525081019061098d9190610e1b565b141561099857600080fd5b6109a0610c1a565b610a3d576002546000546004805460035460405163a9059cbb60e01b81526001600160a01b039586169563a9059cbb956109e295911693929092019101611049565b602060405180830381600087803b1580156109fc57600080fd5b505af1158015610a10573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a349190610de0565b610a3d57600080fd5b33ff5b6001546001600160a01b03163314610a5757600080fd5b60078054600160401b9004810b900b421015610a7257600080fd5b610a7a610c1a565b610a3d576002546001546004805460035460405163a9059cbb60e01b81526001600160a01b039586169563a9059cbb956109e295911693929092019101611049565b6000546001600160a01b03163314610ad357600080fd5b60078054810b900b421015610ae757600080fd5b60095460ff1615610af757600080fd5b610aff610c1a565b610a3d576002546040516370a0823160e01b81526000916001600160a01b0316906370a0823190610b34903090600401611011565b60206040518083038186803b158015610b4c57600080fd5b505afa158015610b60573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b849190610e1b565b90508015610c165760025460005460405163a9059cbb60e01b81526001600160a01b039283169263a9059cbb92610bc2929116908590600401611049565b602060405180830381600087803b158015610bdc57600080fd5b505af1158015610bf0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c149190610de0565b505b5033ff5b6002546001600160a01b03161590565b60408051808201909152606081526000602082015290565b600082601f830112610c52578081fd5b813567ffffffffffffffff80821115610c6757fe5b6020610c7681828502016110b7565b838152935080840185820160005b85811015610d045781358801604080601f19838d03011215610ca557600080fd5b80518181018181108982111715610cb857fe5b82528287013588811115610ccb57600080fd5b610cd98d8983870101610d10565b8252509181013591610cea8361110b565b808701929092525083529183019190830190600101610c84565b50505050505092915050565b600082601f830112610d20578081fd5b813567ffffffffffffffff811115610d3457fe5b610d47601f8201601f19166020016110b7565b9150808252836020828501011115610d5e57600080fd5b8060208401602084013760009082016020015292915050565b803567ffffffffffffffff81168114610d8f57600080fd5b919050565b600080600060608486031215610da8578283fd5b833567ffffffffffffffff811115610dbe578384fd5b610dca86828701610c42565b9660208601359650604090950135949350505050565b600060208284031215610df1578081fd5b8151610dfc8161110b565b9392505050565b600060208284031215610e14578081fd5b5035919050565b600060208284031215610e2c578081fd5b5051919050565b60008060408385031215610e45578182fd5b823567ffffffffffffffff811115610e5b578283fd5b610e6785828601610d10565b925050610e7660208401610d77565b90509250929050565b600080600080600060a08688031215610e96578081fd5b610e9f86610d77565b9450602086013567ffffffffffffffff80821115610ebb578283fd5b610ec789838a01610d10565b9550604088013594506060880135915080821115610ee3578283fd5b610eef89838a01610c42565b93506080880135915080821115610f04578283fd5b50610f1188828901610c42565b9150509295509295909350565b600060ff60f81b851682528360018301528251610f428160218501602087016110db565b91909101602101949350505050565b6001600160f81b0319841681528251600090610f748160018501602088016110db565b60019201918201929092526021019392505050565b6001600160f81b03198416815260c083901b6001600160c01b03191660018201528151600090610fc08160098501602087016110db565b91909101600901949350505050565b9283526001600160c01b031960c092831b81166020850152911b16602882015260300190565b600082516110078184602087016110db565b9190910192915050565b6001600160a01b0391909116815260200190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6001600160a01b03929092168252602082015260400190565b901515815260200190565b90815260200190565b60006020825282518060208401526110958160408501602087016110db565b601f01601f19169190910160400192915050565b60079190910b815260200190565b60405181810167ffffffffffffffff811182821017156110d357fe5b604052919050565b60005b838110156110f65781810151838201526020016110de565b83811115611105576000848401525b50505050565b801515811461111957600080fd5b5056fea264697066735822122024ce373d2e65c57fc3bc39b0613ff31dd2284f8cf0eddb6d81c2618ddd0330c164736f6c63430007030033"

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

// CheckProof is a free data retrieval call binding the contract method 0x8343f445.
//
// Solidity: function checkProof((bytes,bool)[] steps, bytes32 leaf, bytes32 want) pure returns(bool)
func (_Tredd *TreddCaller) CheckProof(opts *bind.CallOpts, steps []TreddProofStep, leaf [32]byte, want [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "checkProof", steps, leaf, want)
	return *ret0, err
}

// CheckProof is a free data retrieval call binding the contract method 0x8343f445.
//
// Solidity: function checkProof((bytes,bool)[] steps, bytes32 leaf, bytes32 want) pure returns(bool)
func (_Tredd *TreddSession) CheckProof(steps []TreddProofStep, leaf [32]byte, want [32]byte) (bool, error) {
	return _Tredd.Contract.CheckProof(&_Tredd.CallOpts, steps, leaf, want)
}

// CheckProof is a free data retrieval call binding the contract method 0x8343f445.
//
// Solidity: function checkProof((bytes,bool)[] steps, bytes32 leaf, bytes32 want) pure returns(bool)
func (_Tredd *TreddCallerSession) CheckProof(steps []TreddProofStep, leaf [32]byte, want [32]byte) (bool, error) {
	return _Tredd.Contract.CheckProof(&_Tredd.CallOpts, steps, leaf, want)
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
