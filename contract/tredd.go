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

// TreddProofStep is an auto generated low-level Go binding around an user-defined struct.
type TreddProofStep struct {
	H    []byte
	Left bool
}

// TreddMetaData contains all meta data concerning the Tredd contract.
var TreddMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenType\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collateral\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"clearRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"cipherRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"revealDeadline\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"refundDeadline\",\"type\":\"uint64\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"decryptionKey\",\"type\":\"bytes32\"}],\"name\":\"evDecryptionKey\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"steps\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"leaf\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"want\",\"type\":\"bytes32\"}],\"name\":\"checkProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"steps\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"prefix\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"chunk\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"want\",\"type\":\"bytes32\"}],\"name\":\"checkProofWithPrefixedChunk\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"steps\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"prefix\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"want\",\"type\":\"bytes32\"}],\"name\":\"checkProofWithPrefixedHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"chunk\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"}],\"name\":\"decrypt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mBuyer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mCipherRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mClearRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mDecryptionKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRefundDeadline\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRevealDeadline\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRevealed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mSeller\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mTokenType\",\"outputs\":[{\"internalType\":\"contractERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paid\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"cipherChunk\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"clearHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"cipherProof\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"clearProof\",\"type\":\"tuple[]\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"decryptionKey\",\"type\":\"bytes32\"}],\"name\":\"reveal\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040526040516122b83803806122b883398181016040528101906100259190610278565b335f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508760015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508660025f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550856003819055508460048190555083600581905550826006819055508160075f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555080600760086101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055505f60095f6101000a81548160ff0219169083151502179055505050505050505050610329565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6101a48261017b565b9050919050565b6101b48161019a565b81146101be575f5ffd5b50565b5f815190506101cf816101ab565b92915050565b5f819050919050565b6101e7816101d5565b81146101f1575f5ffd5b50565b5f81519050610202816101de565b92915050565b5f819050919050565b61021a81610208565b8114610224575f5ffd5b50565b5f8151905061023581610211565b92915050565b5f67ffffffffffffffff82169050919050565b6102578161023b565b8114610261575f5ffd5b50565b5f815190506102728161024e565b92915050565b5f5f5f5f5f5f5f5f610100898b03121561029557610294610177565b5b5f6102a28b828c016101c1565b98505060206102b38b828c016101c1565b97505060406102c48b828c016101f4565b96505060606102d58b828c016101f4565b95505060806102e68b828c01610227565b94505060a06102f78b828c01610227565b93505060c06103088b828c01610264565b92505060e06103198b828c01610264565b9150509295985092959890939650565b611f82806103365f395ff3fe608060405260043610610122575f3560e01c8063649bfb361161009f578063a159896811610063578063a1598968146103b9578063ac280f3d146103f5578063c7dea2f21461041d578063ea8a1af014610433578063fc6210c51461044957610129565b8063649bfb36146102f5578063701fd0f11461031f5780637d966e7d1461033b5780638bae87ba146103655780639067c7a91461038f57610129565b8063295b4e17116100e6578063295b4e17146102115780632df6a9da1461023b57806333bbe2a71461026557806354b53436146102a157806361a5ab22146102cb57610129565b8063095e4c201461012d5780630c590dce146101575780631235ffeb146101815780631d595ee7146101bd57806321b0ae82146101e757610129565b3661012957005b5f5ffd5b348015610138575f5ffd5b50610141610485565b60405161014e9190611286565b60405180910390f35b348015610162575f5ffd5b5061016b61048b565b6040516101789190611319565b60405180910390f35b34801561018c575f5ffd5b506101a760048036038101906101a2919061163a565b6104b0565b6040516101b491906116d1565b60405180910390f35b3480156101c8575f5ffd5b506101d16106a3565b6040516101de91906116f9565b60405180910390f35b3480156101f2575f5ffd5b506101fb6106a9565b60405161020891906116f9565b60405180910390f35b34801561021c575f5ffd5b506102256106af565b6040516102329190611286565b60405180910390f35b348015610246575f5ffd5b5061024f610764565b60405161025c9190611734565b60405180910390f35b348015610270575f5ffd5b5061028b60048036038101906102869190611777565b61077e565b60405161029891906116d1565b60405180910390f35b3480156102ac575f5ffd5b506102b56107b5565b6040516102c291906116d1565b60405180910390f35b3480156102d6575f5ffd5b506102df6107c7565b6040516102ec9190611734565b60405180910390f35b348015610300575f5ffd5b506103096107e0565b6040516103169190611833565b60405180910390f35b6103396004803603810190610334919061184c565b610804565b005b348015610346575f5ffd5b5061034f6109e3565b60405161035c9190611286565b60405180910390f35b348015610370575f5ffd5b506103796109e9565b6040516103869190611833565b60405180910390f35b34801561039a575f5ffd5b506103a3610a0e565b6040516103b091906116f9565b60405180910390f35b3480156103c4575f5ffd5b506103df60048036038101906103da9190611877565b610a14565b6040516103ec9190611931565b60405180910390f35b348015610400575f5ffd5b5061041b60048036038101906104169190611951565b610c33565b005b348015610428575f5ffd5b50610431610e56565b005b34801561043e575f5ffd5b50610447610fd4565b005b348015610454575f5ffd5b5061046f600480360381019061046a9190611a1c565b6111f7565b60405161047c91906116d1565b60405180910390f35b60045481565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f5f5f90505f7f010000000000000000000000000000000000000000000000000000000000000090505f600283876040516020016104ef929190611b21565b60405160208183030381529060405260405161050b9190611b48565b602060405180830381855afa158015610526573d5f5f3e3d5ffd5b5050506040513d601f19601f820116820180604052508101906105499190611b72565b90505f5f90505b87518163ffffffff161015610693575f888263ffffffff168151811061057957610578611b9d565b5b6020026020010151905080602001511561060857600284825f0151856040516020016105a793929190611bea565b6040516020818303038152906040526040516105c39190611b48565b602060405180830381855afa1580156105de573d5f5f3e3d5ffd5b5050506040513d601f19601f820116820180604052508101906106019190611b72565b925061067f565b60028484835f015160405160200161062293929190611c22565b60405160208183030381529060405260405161063e9190611b48565b602060405180830381855afa158015610659573d5f5f3e3d5ffd5b5050506040513d601f19601f8201168201806040525081019061067c9190611b72565b92505b50808061068b90611c96565b915050610550565b5084811493505050509392505050565b60065481565b60055481565b5f6106b861122e565b156106c557479050610761565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b815260040161071f9190611833565b602060405180830381865afa15801561073a573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061075e9190611ceb565b90505b90565b600760089054906101000a900467ffffffffffffffff1681565b5f6107ab858585604051602001610796929190611d4a565b604051602081830303815290604052846104b0565b9050949350505050565b60095f9054906101000a900460ff1681565b60075f9054906101000a900467ffffffffffffffff1681565b5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461085c575f5ffd5b60075f9054906101000a900467ffffffffffffffff1667ffffffffffffffff164210610886575f5ffd5b60095f9054906101000a900460ff161561089e575f5ffd5b6108a661122e565b156108be576004543410156108b9575f5ffd5b610988565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16306004546040518463ffffffff1660e01b815260040161093f93929190611d71565b6020604051808303815f875af115801561095b573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061097f9190611dba565b610987575f5ffd5b5b80600881905550600160095f6101000a81548160ff0219169083151502179055507f34292d279a4eb74e15e8b454d2d45ea717fe4914773188f0540fd6fbe67db081816040516109d891906116f9565b60405180910390a150565b60035481565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60085481565b60605f835167ffffffffffffffff811115610a3257610a31611357565b5b6040519080825280601f01601f191660200182016040528015610a645781602001600182028036833780820191505090505b5090505f5f90505b8451602082610a7b9190611de5565b67ffffffffffffffff161015610c28575f602082610a999190611de5565b90505f60026008548785604051602001610ab593929190611e21565b604051602081830303815290604052604051610ad19190611b48565b602060405180830381855afa158015610aec573d5f5f3e3d5ffd5b5050506040513d601f19601f82011682018060405250810190610b0f9190611b72565b90505f5f90505b60208163ffffffff16108015610b47575087518163ffffffff1684610b3b9190611e5d565b67ffffffffffffffff16105b15610c1257818163ffffffff1660208110610b6557610b64611b9d565b5b1a60f81b888263ffffffff1685610b7c9190611e5d565b67ffffffffffffffff1681518110610b9757610b96611b9d565b5b602001015160f81c60f81b18858263ffffffff1685610bb69190611e5d565b67ffffffffffffffff1681518110610bd157610bd0611b9d565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a9053508080610c0a90611c96565b915050610b16565b5050508080610c2090611e98565b915050610a6c565b508091505092915050565b5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610c8a575f5ffd5b600760089054906101000a900467ffffffffffffffff1667ffffffffffffffff164210610cb5575f5ffd5b60095f9054906101000a900460ff16610ccc575f5ffd5b610cda82868660065461077e565b610ce2575f5ffd5b610cf08186856005546111f7565b610cf8575f5ffd5b5f610d038587610a14565b905083600282604051610d169190611b48565b602060405180830381855afa158015610d31573d5f5f3e3d5ffd5b5050506040513d601f19601f82011682018060405250810190610d549190611b72565b03610d5d575f5ffd5b610d6561122e565b610e3d5760025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600454600354610dd79190611ec7565b6040518363ffffffff1660e01b8152600401610df4929190611efa565b6020604051808303815f875af1158015610e10573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610e349190611dba565b610e3c575f5ffd5b5b3373ffffffffffffffffffffffffffffffffffffffff16ff5b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610eae575f5ffd5b600760089054906101000a900467ffffffffffffffff1667ffffffffffffffff16421015610eda575f5ffd5b610ee261122e565b610fbb5760025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600454600354610f559190611ec7565b6040518363ffffffff1660e01b8152600401610f72929190611efa565b6020604051808303815f875af1158015610f8e573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610fb29190611dba565b610fba575f5ffd5b5b3373ffffffffffffffffffffffffffffffffffffffff16ff5b5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461102b575f5ffd5b60075f9054906101000a900467ffffffffffffffff1667ffffffffffffffff16421015611056575f5ffd5b60095f9054906101000a900460ff161561106e575f5ffd5b61107661122e565b6111de575f60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016110d59190611833565b602060405180830381865afa1580156110f0573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906111149190611ceb565b90505f8111156111dc5760025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16836040518363ffffffff1660e01b815260040161119a929190611efa565b6020604051808303815f875af11580156111b6573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906111da9190611dba565b505b505b3373ffffffffffffffffffffffffffffffffffffffff16ff5b5f61122485858560405160200161120f929190611f21565b604051602081830303815290604052846104b0565b9050949350505050565b5f5f60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614905090565b5f819050919050565b6112808161126e565b82525050565b5f6020820190506112995f830184611277565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f819050919050565b5f6112e16112dc6112d78461129f565b6112be565b61129f565b9050919050565b5f6112f2826112c7565b9050919050565b5f611303826112e8565b9050919050565b611313816112f9565b82525050565b5f60208201905061132c5f83018461130a565b92915050565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b61138d82611347565b810181811067ffffffffffffffff821117156113ac576113ab611357565b5b80604052505050565b5f6113be611332565b90506113ca8282611384565b919050565b5f67ffffffffffffffff8211156113e9576113e8611357565b5b602082029050602081019050919050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f67ffffffffffffffff82111561142457611423611357565b5b61142d82611347565b9050602081019050919050565b828183375f83830152505050565b5f61145a6114558461140a565b6113b5565b90508281526020810184848401111561147657611475611406565b5b61148184828561143a565b509392505050565b5f82601f83011261149d5761149c611343565b5b81356114ad848260208601611448565b91505092915050565b5f8115159050919050565b6114ca816114b6565b81146114d4575f5ffd5b50565b5f813590506114e5816114c1565b92915050565b5f60408284031215611500576114ff6113fe565b5b61150a60406113b5565b90505f82013567ffffffffffffffff81111561152957611528611402565b5b61153584828501611489565b5f830152506020611548848285016114d7565b60208301525092915050565b5f611566611561846113cf565b6113b5565b90508083825260208201905060208402830185811115611589576115886113fa565b5b835b818110156115d057803567ffffffffffffffff8111156115ae576115ad611343565b5b8086016115bb89826114eb565b8552602085019450505060208101905061158b565b5050509392505050565b5f82601f8301126115ee576115ed611343565b5b81356115fe848260208601611554565b91505092915050565b5f819050919050565b61161981611607565b8114611623575f5ffd5b50565b5f8135905061163481611610565b92915050565b5f5f5f606084860312156116515761165061133b565b5b5f84013567ffffffffffffffff81111561166e5761166d61133f565b5b61167a868287016115da565b935050602084013567ffffffffffffffff81111561169b5761169a61133f565b5b6116a786828701611489565b92505060406116b886828701611626565b9150509250925092565b6116cb816114b6565b82525050565b5f6020820190506116e45f8301846116c2565b92915050565b6116f381611607565b82525050565b5f60208201905061170c5f8301846116ea565b92915050565b5f67ffffffffffffffff82169050919050565b61172e81611712565b82525050565b5f6020820190506117475f830184611725565b92915050565b61175681611712565b8114611760575f5ffd5b50565b5f813590506117718161174d565b92915050565b5f5f5f5f6080858703121561178f5761178e61133b565b5b5f85013567ffffffffffffffff8111156117ac576117ab61133f565b5b6117b8878288016115da565b94505060206117c987828801611763565b935050604085013567ffffffffffffffff8111156117ea576117e961133f565b5b6117f687828801611489565b925050606061180787828801611626565b91505092959194509250565b5f61181d8261129f565b9050919050565b61182d81611813565b82525050565b5f6020820190506118465f830184611824565b92915050565b5f602082840312156118615761186061133b565b5b5f61186e84828501611626565b91505092915050565b5f5f6040838503121561188d5761188c61133b565b5b5f83013567ffffffffffffffff8111156118aa576118a961133f565b5b6118b685828601611489565b92505060206118c785828601611763565b9150509250929050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f611903826118d1565b61190d81856118db565b935061191d8185602086016118eb565b61192681611347565b840191505092915050565b5f6020820190508181035f83015261194981846118f9565b905092915050565b5f5f5f5f5f60a0868803121561196a5761196961133b565b5b5f61197788828901611763565b955050602086013567ffffffffffffffff8111156119985761199761133f565b5b6119a488828901611489565b94505060406119b588828901611626565b935050606086013567ffffffffffffffff8111156119d6576119d561133f565b5b6119e2888289016115da565b925050608086013567ffffffffffffffff811115611a0357611a0261133f565b5b611a0f888289016115da565b9150509295509295909350565b5f5f5f5f60808587031215611a3457611a3361133b565b5b5f85013567ffffffffffffffff811115611a5157611a5061133f565b5b611a5d878288016115da565b9450506020611a6e87828801611763565b9350506040611a7f87828801611626565b9250506060611a9087828801611626565b91505092959194509250565b5f7fff0000000000000000000000000000000000000000000000000000000000000082169050919050565b5f819050919050565b611ae1611adc82611a9c565b611ac7565b82525050565b5f81905092915050565b5f611afb826118d1565b611b058185611ae7565b9350611b158185602086016118eb565b80840191505092915050565b5f611b2c8285611ad0565b600182019150611b3c8284611af1565b91508190509392505050565b5f611b538284611af1565b915081905092915050565b5f81519050611b6c81611610565b92915050565b5f60208284031215611b8757611b8661133b565b5b5f611b9484828501611b5e565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f819050919050565b611be4611bdf82611607565b611bca565b82525050565b5f611bf58286611ad0565b600182019150611c058285611af1565b9150611c118284611bd3565b602082019150819050949350505050565b5f611c2d8286611ad0565b600182019150611c3d8285611bd3565b602082019150611c4d8284611af1565b9150819050949350505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f63ffffffff82169050919050565b5f611ca082611c87565b915063ffffffff8203611cb657611cb5611c5a565b5b600182019050919050565b611cca8161126e565b8114611cd4575f5ffd5b50565b5f81519050611ce581611cc1565b92915050565b5f60208284031215611d0057611cff61133b565b5b5f611d0d84828501611cd7565b91505092915050565b5f8160c01b9050919050565b5f611d2c82611d16565b9050919050565b611d44611d3f82611712565b611d22565b82525050565b5f611d558285611d33565b600882019150611d658284611af1565b91508190509392505050565b5f606082019050611d845f830186611824565b611d916020830185611824565b611d9e6040830184611277565b949350505050565b5f81519050611db4816114c1565b92915050565b5f60208284031215611dcf57611dce61133b565b5b5f611ddc84828501611da6565b91505092915050565b5f611def82611712565b9150611dfa83611712565b9250828202611e0881611712565b9150808214611e1a57611e19611c5a565b5b5092915050565b5f611e2c8286611bd3565b602082019150611e3c8285611d33565b600882019150611e4c8284611d33565b600882019150819050949350505050565b5f611e6782611712565b9150611e7283611712565b9250828201905067ffffffffffffffff811115611e9257611e91611c5a565b5b92915050565b5f611ea282611712565b915067ffffffffffffffff8203611ebc57611ebb611c5a565b5b600182019050919050565b5f611ed18261126e565b9150611edc8361126e565b9250828201905080821115611ef457611ef3611c5a565b5b92915050565b5f604082019050611f0d5f830185611824565b611f1a6020830184611277565b9392505050565b5f611f2c8285611d33565b600882019150611f3c8284611bd3565b602082019150819050939250505056fea26469706673582212202e5eeff9919cdc3bbae7c3941d47408cba4ee157bce6aa3a94113c53c6373e3764736f6c63430008210033",
}

// TreddABI is the input ABI used to generate the binding from.
// Deprecated: Use TreddMetaData.ABI instead.
var TreddABI = TreddMetaData.ABI

// TreddBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TreddMetaData.Bin instead.
var TreddBin = TreddMetaData.Bin

// DeployTredd deploys a new Ethereum contract, binding an instance of Tredd to it.
func DeployTredd(auth *bind.TransactOpts, backend bind.ContractBackend, seller common.Address, tokenType common.Address, amount *big.Int, collateral *big.Int, clearRoot [32]byte, cipherRoot [32]byte, revealDeadline uint64, refundDeadline uint64) (common.Address, *types.Transaction, *Tredd, error) {
	parsed, err := TreddMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TreddBin), backend, seller, tokenType, amount, collateral, clearRoot, cipherRoot, revealDeadline, refundDeadline)
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
	parsed, err := TreddMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Tredd *TreddRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
func (_Tredd *TreddCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "checkProof", steps, leaf, want)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

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
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "checkProofWithPrefixedChunk", steps, prefix, chunk, want)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

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
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "checkProofWithPrefixedHash", steps, prefix, hash, want)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

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
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "decrypt", chunk, index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

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
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "mAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

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
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "mBuyer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

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
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "mCipherRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

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
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "mClearRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

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
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "mCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

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
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "mDecryptionKey")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

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
// Solidity: function mRefundDeadline() view returns(uint64)
func (_Tredd *TreddCaller) MRefundDeadline(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "mRefundDeadline")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MRefundDeadline is a free data retrieval call binding the contract method 0x2df6a9da.
//
// Solidity: function mRefundDeadline() view returns(uint64)
func (_Tredd *TreddSession) MRefundDeadline() (uint64, error) {
	return _Tredd.Contract.MRefundDeadline(&_Tredd.CallOpts)
}

// MRefundDeadline is a free data retrieval call binding the contract method 0x2df6a9da.
//
// Solidity: function mRefundDeadline() view returns(uint64)
func (_Tredd *TreddCallerSession) MRefundDeadline() (uint64, error) {
	return _Tredd.Contract.MRefundDeadline(&_Tredd.CallOpts)
}

// MRevealDeadline is a free data retrieval call binding the contract method 0x61a5ab22.
//
// Solidity: function mRevealDeadline() view returns(uint64)
func (_Tredd *TreddCaller) MRevealDeadline(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "mRevealDeadline")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MRevealDeadline is a free data retrieval call binding the contract method 0x61a5ab22.
//
// Solidity: function mRevealDeadline() view returns(uint64)
func (_Tredd *TreddSession) MRevealDeadline() (uint64, error) {
	return _Tredd.Contract.MRevealDeadline(&_Tredd.CallOpts)
}

// MRevealDeadline is a free data retrieval call binding the contract method 0x61a5ab22.
//
// Solidity: function mRevealDeadline() view returns(uint64)
func (_Tredd *TreddCallerSession) MRevealDeadline() (uint64, error) {
	return _Tredd.Contract.MRevealDeadline(&_Tredd.CallOpts)
}

// MRevealed is a free data retrieval call binding the contract method 0x54b53436.
//
// Solidity: function mRevealed() view returns(bool)
func (_Tredd *TreddCaller) MRevealed(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "mRevealed")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

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
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "mSeller")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

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
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "mTokenType")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

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
	var out []interface{}
	err := _Tredd.contract.Call(opts, &out, "paid")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

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
	event.Raw = log
	return event, nil
}
