// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tredd

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

// TreddABI is the input ABI used to generate the binding from.
const TreddABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"tokenType\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"collateral\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"clearRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"cipherRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"revealDeadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"refundDeadline\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"decryptionKey\",\"type\":\"bytes32\"}],\"name\":\"evDecryptionKey\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"claimPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mBuyer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mCipherRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mClearRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mDecryptionKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRefundDeadline\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRevealDeadline\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRevealed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mSeller\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mTokenType\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reclaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"cipherChunk\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"clearHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"cipherProof\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"clearProof\",\"type\":\"tuple[]\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"decryptionKey\",\"type\":\"bytes32\"}],\"name\":\"reveal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TreddFuncSigs maps the 4-byte function signature to its string representation.
var TreddFuncSigs = map[string]string{
	"c7dea2f2": "claimPayment()",
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
	"80e9071b": "reclaim()",
	"ac280f3d": "refund(uint64,bytes,bytes32,(bytes,bool)[],(bytes,bool)[])",
	"701fd0f1": "reveal(bytes32)",
}

// TreddBin is the compiled bytecode used for deploying new contracts.
var TreddBin = "0x60806040523480156200001157600080fd5b5060405162000cf038038062000cf083398101604081905262000034916200013c565b60008054336001600160a01b031991821617909155600180549091166001600160a01b038a161790556002879055855162000077906003906020890190620000a0565b506004949094556005929092556006556007556008555050600a805460ff191690555062000258565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620000e357805160ff191683800117855562000113565b8280016001018555821562000113579182015b8281111562000113578251825591602001919060010190620000f6565b506200012192915062000125565b5090565b5b8082111562000121576000815560010162000126565b600080600080600080600080610100898b03121562000159578384fd5b88516001600160a01b038116811462000170578485fd5b60208a81015160408c0151929a509850906001600160401b038082111562000196578687fd5b818c0191508c601f830112620001aa578687fd5b815181811115620001b757fe5b604051601f8201601f1916810185018381118282101715620001d557fe5b60405281815283820185018f1015620001ec578889fd5b8892505b818310156200020f5783830185015181840186015291840191620001f0565b818311156200022057888583830101525b809a505050505050606089015194506080890151935060a0890151925060c0890151915060e089015190509295985092959890939650565b610a8880620002686000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c8063649bfb36116100975780638bae87ba116100665780638bae87ba1461019c5780639067c7a9146101a4578063ac280f3d146101ac578063c7dea2f2146101bf576100f5565b8063649bfb3614610162578063701fd0f1146101775780637d966e7d1461018c57806380e9071b14610194576100f5565b806321b0ae82116100d357806321b0ae82146101355780632df6a9da1461013d57806354b534361461014557806361a5ab221461015a576100f5565b8063095e4c20146100fa5780630c590dce146101185780631d595ee71461012d575b600080fd5b6101026101c7565b60405161010f91906109c2565b60405180910390f35b6101206101cd565b60405161010f91906109cb565b61010261025b565b610102610261565b610102610267565b61014d61026d565b60405161010f91906109b7565b610102610276565b61016a61027c565b60405161010f91906109a3565b61018a6101853660046107f5565b61028b565b005b61010261030d565b61018a610313565b61016a61034c565b61010261035b565b61018a6101ba366004610825565b610361565b61018a6104bd565b60045481565b6003805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156102535780601f1061022857610100808354040283529160200191610253565b820191906000526020600020905b81548152906001019060200180831161023657829003601f168201915b505050505081565b60065481565b60055481565b60085481565b600a5460ff1681565b60075481565b6000546001600160a01b031681565b6001546001600160a01b031633146102a257600080fd5b60075442106102b057600080fd5b600a5460ff16156102c057600080fd5b6009819055600a805460ff191660011790556040517f34292d279a4eb74e15e8b454d2d45ea717fe4914773188f0540fd6fbe67db081906103029083906109c2565b60405180910390a150565b60025481565b6000546001600160a01b0316331461032a57600080fd5b60075442101561033957600080fd5b600a5460ff161561034957600080fd5b33ff5b6001546001600160a01b031681565b60095481565b6000546001600160a01b0316331461037857600080fd5b600854421061038657600080fd5b600a5460ff1661039557600080fd5b61041082600287876040516020016103ae92919061096d565b60408051601f19818403018152908290526103c891610951565b602060405180830381855afa1580156103e5573d6000803e3d6000fd5b5050506040513d601f19601f82011682018060405250810190610408919061080d565b6006546104e3565b61041957600080fd5b61042681846005546104e3565b61042f57600080fd5b82600286610447878967ffffffffffffffff1661069b565b60405160200161045892919061096d565b60408051601f198184030181529082905261047291610951565b602060405180830381855afa15801561048f573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906104b2919061080d565b141561034957600080fd5b6001546001600160a01b031633146104d457600080fd5b60085442101561034957600080fd5b6040516000908190600160f81b90829060029061050690839089906020016108cc565b60408051601f198184030181529082905261052091610951565b602060405180830381855afa15801561053d573d6000803e3d6000fd5b5050506040513d601f19601f82011682018060405250810190610560919061080d565b905060005b875181101561068e576105766106a3565b88828151811061058257fe5b602002602001015190508060200151156106105780516040516002916105af918791908790602001610919565b60408051601f19818403018152908290526105c991610951565b602060405180830381855afa1580156105e6573d6000803e3d6000fd5b5050506040513d601f19601f82011682018060405250810190610609919061080d565b9250610685565b805160405160029161062891879187916020016108e6565b60408051601f198184030181529082905261064291610951565b602060405180830381855afa15801561065f573d6000803e3d6000fd5b5050506040513d601f19601f82011682018060405250810190610682919061080d565b92505b50600101610565565b5090931495945050505050565b606092915050565b60408051808201909152606081526000602082015290565b600082601f8301126106cb578081fd5b813567ffffffffffffffff808211156106e057fe5b60206106ef81828502016109fe565b838152935080840185820160005b858110156107825781358801604080601f19838d0301121561071e57600080fd5b8051818101818110898211171561073157fe5b8252828701358881111561074457600080fd5b6107528d898387010161078e565b8252509181013591821515831461076857600080fd5b8087019290925250835291830191908301906001016106fd565b50505050505092915050565b600082601f83011261079e578081fd5b813567ffffffffffffffff8111156107b257fe5b6107c5601f8201601f19166020016109fe565b91508082528360208285010111156107dc57600080fd5b8060208401602084013760009082016020015292915050565b600060208284031215610806578081fd5b5035919050565b60006020828403121561081e578081fd5b5051919050565b600080600080600060a0868803121561083c578081fd5b853567ffffffffffffffff8082168214610854578283fd5b90955060208701359080821115610869578283fd5b61087589838a0161078e565b9550604088013594506060880135915080821115610891578283fd5b61089d89838a016106bb565b935060808801359150808211156108b2578283fd5b506108bf888289016106bb565b9150509295509295909350565b6001600160f81b0319929092168252600182015260210190565b600060ff60f81b85168252836001830152825161090a816021850160208701610a22565b91909101602101949350505050565b6001600160f81b031984168152825160009061093c816001850160208801610a22565b60019201918201929092526021019392505050565b60008251610963818460208701610a22565b9190910192915050565b600067ffffffffffffffff60c01b8460c01b1682528251610995816008850160208701610a22565b919091016008019392505050565b6001600160a01b0391909116815260200190565b901515815260200190565b90815260200190565b60006020825282518060208401526109ea816040850160208701610a22565b601f01601f19169190910160400192915050565b60405181810167ffffffffffffffff81118282101715610a1a57fe5b604052919050565b60005b83811015610a3d578181015183820152602001610a25565b83811115610a4c576000848401525b5050505056fea2646970667358221220c57549352b842352c99e4e983c1e76068b1c5fcbc961ba2671f054995713500264736f6c63430007020033"

// DeployTredd deploys a new Ethereum contract, binding an instance of Tredd to it.
func DeployTredd(auth *bind.TransactOpts, backend bind.ContractBackend, seller common.Address, amount *big.Int, tokenType []byte, collateral *big.Int, clearRoot [32]byte, cipherRoot [32]byte, revealDeadline *big.Int, refundDeadline *big.Int) (common.Address, *types.Transaction, *Tredd, error) {
	parsed, err := abi.JSON(strings.NewReader(TreddABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TreddBin), backend, seller, amount, tokenType, collateral, clearRoot, cipherRoot, revealDeadline, refundDeadline)
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
// Solidity: function mRefundDeadline() view returns(uint256)
func (_Tredd *TreddCaller) MRefundDeadline(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mRefundDeadline")
	return *ret0, err
}

// MRefundDeadline is a free data retrieval call binding the contract method 0x2df6a9da.
//
// Solidity: function mRefundDeadline() view returns(uint256)
func (_Tredd *TreddSession) MRefundDeadline() (*big.Int, error) {
	return _Tredd.Contract.MRefundDeadline(&_Tredd.CallOpts)
}

// MRefundDeadline is a free data retrieval call binding the contract method 0x2df6a9da.
//
// Solidity: function mRefundDeadline() view returns(uint256)
func (_Tredd *TreddCallerSession) MRefundDeadline() (*big.Int, error) {
	return _Tredd.Contract.MRefundDeadline(&_Tredd.CallOpts)
}

// MRevealDeadline is a free data retrieval call binding the contract method 0x61a5ab22.
//
// Solidity: function mRevealDeadline() view returns(uint256)
func (_Tredd *TreddCaller) MRevealDeadline(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mRevealDeadline")
	return *ret0, err
}

// MRevealDeadline is a free data retrieval call binding the contract method 0x61a5ab22.
//
// Solidity: function mRevealDeadline() view returns(uint256)
func (_Tredd *TreddSession) MRevealDeadline() (*big.Int, error) {
	return _Tredd.Contract.MRevealDeadline(&_Tredd.CallOpts)
}

// MRevealDeadline is a free data retrieval call binding the contract method 0x61a5ab22.
//
// Solidity: function mRevealDeadline() view returns(uint256)
func (_Tredd *TreddCallerSession) MRevealDeadline() (*big.Int, error) {
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
// Solidity: function mTokenType() view returns(bytes)
func (_Tredd *TreddCaller) MTokenType(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _Tredd.contract.Call(opts, out, "mTokenType")
	return *ret0, err
}

// MTokenType is a free data retrieval call binding the contract method 0x0c590dce.
//
// Solidity: function mTokenType() view returns(bytes)
func (_Tredd *TreddSession) MTokenType() ([]byte, error) {
	return _Tredd.Contract.MTokenType(&_Tredd.CallOpts)
}

// MTokenType is a free data retrieval call binding the contract method 0x0c590dce.
//
// Solidity: function mTokenType() view returns(bytes)
func (_Tredd *TreddCallerSession) MTokenType() ([]byte, error) {
	return _Tredd.Contract.MTokenType(&_Tredd.CallOpts)
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

// Reclaim is a paid mutator transaction binding the contract method 0x80e9071b.
//
// Solidity: function reclaim() returns()
func (_Tredd *TreddTransactor) Reclaim(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Tredd.contract.Transact(opts, "reclaim")
}

// Reclaim is a paid mutator transaction binding the contract method 0x80e9071b.
//
// Solidity: function reclaim() returns()
func (_Tredd *TreddSession) Reclaim() (*types.Transaction, error) {
	return _Tredd.Contract.Reclaim(&_Tredd.TransactOpts)
}

// Reclaim is a paid mutator transaction binding the contract method 0x80e9071b.
//
// Solidity: function reclaim() returns()
func (_Tredd *TreddTransactorSession) Reclaim() (*types.Transaction, error) {
	return _Tredd.Contract.Reclaim(&_Tredd.TransactOpts)
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
// Solidity: function reveal(bytes32 decryptionKey) returns()
func (_Tredd *TreddTransactor) Reveal(opts *bind.TransactOpts, decryptionKey [32]byte) (*types.Transaction, error) {
	return _Tredd.contract.Transact(opts, "reveal", decryptionKey)
}

// Reveal is a paid mutator transaction binding the contract method 0x701fd0f1.
//
// Solidity: function reveal(bytes32 decryptionKey) returns()
func (_Tredd *TreddSession) Reveal(decryptionKey [32]byte) (*types.Transaction, error) {
	return _Tredd.Contract.Reveal(&_Tredd.TransactOpts, decryptionKey)
}

// Reveal is a paid mutator transaction binding the contract method 0x701fd0f1.
//
// Solidity: function reveal(bytes32 decryptionKey) returns()
func (_Tredd *TreddTransactorSession) Reveal(decryptionKey [32]byte) (*types.Transaction, error) {
	return _Tredd.Contract.Reveal(&_Tredd.TransactOpts, decryptionKey)
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
