// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// TreddProofStep is an auto generated low-level Go binding around an user-defined struct.
type TreddProofStep struct {
	H    []byte
	Left bool
}

// TreddMetaData contains all meta data concerning the Tredd contract.
var TreddMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenType\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collateral\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"clearRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"cipherRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"revealDeadline\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"refundDeadline\",\"type\":\"uint64\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"decryptionKey\",\"type\":\"bytes32\"}],\"name\":\"evDecryptionKey\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"steps\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"leaf\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"want\",\"type\":\"bytes32\"}],\"name\":\"checkProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"steps\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"prefix\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"chunk\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"want\",\"type\":\"bytes32\"}],\"name\":\"checkProofWithPrefixedChunk\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"steps\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"prefix\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"want\",\"type\":\"bytes32\"}],\"name\":\"checkProofWithPrefixedHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"chunk\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"}],\"name\":\"decrypt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mBuyer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mCipherRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mClearRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mDecryptionKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRefundDeadline\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRevealDeadline\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mRevealed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mSeller\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mTokenType\",\"outputs\":[{\"internalType\":\"contractERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paid\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"cipherChunk\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"clearHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"cipherProof\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"h\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"left\",\"type\":\"bool\"}],\"internalType\":\"structTredd.ProofStep[]\",\"name\":\"clearProof\",\"type\":\"tuple[]\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"decryptionKey\",\"type\":\"bytes32\"}],\"name\":\"reveal\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	ID:  "Tredd",
	Bin: "0x60806040526040516122b83803806122b883398181016040528101906100259190610278565b335f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508760015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508660025f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550856003819055508460048190555083600581905550826006819055508160075f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555080600760086101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055505f60095f6101000a81548160ff0219169083151502179055505050505050505050610329565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6101a48261017b565b9050919050565b6101b48161019a565b81146101be575f5ffd5b50565b5f815190506101cf816101ab565b92915050565b5f819050919050565b6101e7816101d5565b81146101f1575f5ffd5b50565b5f81519050610202816101de565b92915050565b5f819050919050565b61021a81610208565b8114610224575f5ffd5b50565b5f8151905061023581610211565b92915050565b5f67ffffffffffffffff82169050919050565b6102578161023b565b8114610261575f5ffd5b50565b5f815190506102728161024e565b92915050565b5f5f5f5f5f5f5f5f610100898b03121561029557610294610177565b5b5f6102a28b828c016101c1565b98505060206102b38b828c016101c1565b97505060406102c48b828c016101f4565b96505060606102d58b828c016101f4565b95505060806102e68b828c01610227565b94505060a06102f78b828c01610227565b93505060c06103088b828c01610264565b92505060e06103198b828c01610264565b9150509295985092959890939650565b611f82806103365f395ff3fe608060405260043610610122575f3560e01c8063649bfb361161009f578063a159896811610063578063a1598968146103b9578063ac280f3d146103f5578063c7dea2f21461041d578063ea8a1af014610433578063fc6210c51461044957610129565b8063649bfb36146102f5578063701fd0f11461031f5780637d966e7d1461033b5780638bae87ba146103655780639067c7a91461038f57610129565b8063295b4e17116100e6578063295b4e17146102115780632df6a9da1461023b57806333bbe2a71461026557806354b53436146102a157806361a5ab22146102cb57610129565b8063095e4c201461012d5780630c590dce146101575780631235ffeb146101815780631d595ee7146101bd57806321b0ae82146101e757610129565b3661012957005b5f5ffd5b348015610138575f5ffd5b50610141610485565b60405161014e9190611286565b60405180910390f35b348015610162575f5ffd5b5061016b61048b565b6040516101789190611319565b60405180910390f35b34801561018c575f5ffd5b506101a760048036038101906101a2919061163a565b6104b0565b6040516101b491906116d1565b60405180910390f35b3480156101c8575f5ffd5b506101d16106a3565b6040516101de91906116f9565b60405180910390f35b3480156101f2575f5ffd5b506101fb6106a9565b60405161020891906116f9565b60405180910390f35b34801561021c575f5ffd5b506102256106af565b6040516102329190611286565b60405180910390f35b348015610246575f5ffd5b5061024f610764565b60405161025c9190611734565b60405180910390f35b348015610270575f5ffd5b5061028b60048036038101906102869190611777565b61077e565b60405161029891906116d1565b60405180910390f35b3480156102ac575f5ffd5b506102b56107b5565b6040516102c291906116d1565b60405180910390f35b3480156102d6575f5ffd5b506102df6107c7565b6040516102ec9190611734565b60405180910390f35b348015610300575f5ffd5b506103096107e0565b6040516103169190611833565b60405180910390f35b6103396004803603810190610334919061184c565b610804565b005b348015610346575f5ffd5b5061034f6109e3565b60405161035c9190611286565b60405180910390f35b348015610370575f5ffd5b506103796109e9565b6040516103869190611833565b60405180910390f35b34801561039a575f5ffd5b506103a3610a0e565b6040516103b091906116f9565b60405180910390f35b3480156103c4575f5ffd5b506103df60048036038101906103da9190611877565b610a14565b6040516103ec9190611931565b60405180910390f35b348015610400575f5ffd5b5061041b60048036038101906104169190611951565b610c33565b005b348015610428575f5ffd5b50610431610e56565b005b34801561043e575f5ffd5b50610447610fd4565b005b348015610454575f5ffd5b5061046f600480360381019061046a9190611a1c565b6111f7565b60405161047c91906116d1565b60405180910390f35b60045481565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f5f5f90505f7f010000000000000000000000000000000000000000000000000000000000000090505f600283876040516020016104ef929190611b21565b60405160208183030381529060405260405161050b9190611b48565b602060405180830381855afa158015610526573d5f5f3e3d5ffd5b5050506040513d601f19601f820116820180604052508101906105499190611b72565b90505f5f90505b87518163ffffffff161015610693575f888263ffffffff168151811061057957610578611b9d565b5b6020026020010151905080602001511561060857600284825f0151856040516020016105a793929190611bea565b6040516020818303038152906040526040516105c39190611b48565b602060405180830381855afa1580156105de573d5f5f3e3d5ffd5b5050506040513d601f19601f820116820180604052508101906106019190611b72565b925061067f565b60028484835f015160405160200161062293929190611c22565b60405160208183030381529060405260405161063e9190611b48565b602060405180830381855afa158015610659573d5f5f3e3d5ffd5b5050506040513d601f19601f8201168201806040525081019061067c9190611b72565b92505b50808061068b90611c96565b915050610550565b5084811493505050509392505050565b60065481565b60055481565b5f6106b861122e565b156106c557479050610761565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b815260040161071f9190611833565b602060405180830381865afa15801561073a573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061075e9190611ceb565b90505b90565b600760089054906101000a900467ffffffffffffffff1681565b5f6107ab858585604051602001610796929190611d4a565b604051602081830303815290604052846104b0565b9050949350505050565b60095f9054906101000a900460ff1681565b60075f9054906101000a900467ffffffffffffffff1681565b5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461085c575f5ffd5b60075f9054906101000a900467ffffffffffffffff1667ffffffffffffffff164210610886575f5ffd5b60095f9054906101000a900460ff161561089e575f5ffd5b6108a661122e565b156108be576004543410156108b9575f5ffd5b610988565b60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16306004546040518463ffffffff1660e01b815260040161093f93929190611d71565b6020604051808303815f875af115801561095b573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061097f9190611dba565b610987575f5ffd5b5b80600881905550600160095f6101000a81548160ff0219169083151502179055507f34292d279a4eb74e15e8b454d2d45ea717fe4914773188f0540fd6fbe67db081816040516109d891906116f9565b60405180910390a150565b60035481565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60085481565b60605f835167ffffffffffffffff811115610a3257610a31611357565b5b6040519080825280601f01601f191660200182016040528015610a645781602001600182028036833780820191505090505b5090505f5f90505b8451602082610a7b9190611de5565b67ffffffffffffffff161015610c28575f602082610a999190611de5565b90505f60026008548785604051602001610ab593929190611e21565b604051602081830303815290604052604051610ad19190611b48565b602060405180830381855afa158015610aec573d5f5f3e3d5ffd5b5050506040513d601f19601f82011682018060405250810190610b0f9190611b72565b90505f5f90505b60208163ffffffff16108015610b47575087518163ffffffff1684610b3b9190611e5d565b67ffffffffffffffff16105b15610c1257818163ffffffff1660208110610b6557610b64611b9d565b5b1a60f81b888263ffffffff1685610b7c9190611e5d565b67ffffffffffffffff1681518110610b9757610b96611b9d565b5b602001015160f81c60f81b18858263ffffffff1685610bb69190611e5d565b67ffffffffffffffff1681518110610bd157610bd0611b9d565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a9053508080610c0a90611c96565b915050610b16565b5050508080610c2090611e98565b915050610a6c565b508091505092915050565b5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610c8a575f5ffd5b600760089054906101000a900467ffffffffffffffff1667ffffffffffffffff164210610cb5575f5ffd5b60095f9054906101000a900460ff16610ccc575f5ffd5b610cda82868660065461077e565b610ce2575f5ffd5b610cf08186856005546111f7565b610cf8575f5ffd5b5f610d038587610a14565b905083600282604051610d169190611b48565b602060405180830381855afa158015610d31573d5f5f3e3d5ffd5b5050506040513d601f19601f82011682018060405250810190610d549190611b72565b03610d5d575f5ffd5b610d6561122e565b610e3d5760025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600454600354610dd79190611ec7565b6040518363ffffffff1660e01b8152600401610df4929190611efa565b6020604051808303815f875af1158015610e10573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610e349190611dba565b610e3c575f5ffd5b5b3373ffffffffffffffffffffffffffffffffffffffff16ff5b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610eae575f5ffd5b600760089054906101000a900467ffffffffffffffff1667ffffffffffffffff16421015610eda575f5ffd5b610ee261122e565b610fbb5760025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600454600354610f559190611ec7565b6040518363ffffffff1660e01b8152600401610f72929190611efa565b6020604051808303815f875af1158015610f8e573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610fb29190611dba565b610fba575f5ffd5b5b3373ffffffffffffffffffffffffffffffffffffffff16ff5b5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461102b575f5ffd5b60075f9054906101000a900467ffffffffffffffff1667ffffffffffffffff16421015611056575f5ffd5b60095f9054906101000a900460ff161561106e575f5ffd5b61107661122e565b6111de575f60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016110d59190611833565b602060405180830381865afa1580156110f0573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906111149190611ceb565b90505f8111156111dc5760025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16836040518363ffffffff1660e01b815260040161119a929190611efa565b6020604051808303815f875af11580156111b6573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906111da9190611dba565b505b505b3373ffffffffffffffffffffffffffffffffffffffff16ff5b5f61122485858560405160200161120f929190611f21565b604051602081830303815290604052846104b0565b9050949350505050565b5f5f60025f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614905090565b5f819050919050565b6112808161126e565b82525050565b5f6020820190506112995f830184611277565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f819050919050565b5f6112e16112dc6112d78461129f565b6112be565b61129f565b9050919050565b5f6112f2826112c7565b9050919050565b5f611303826112e8565b9050919050565b611313816112f9565b82525050565b5f60208201905061132c5f83018461130a565b92915050565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b61138d82611347565b810181811067ffffffffffffffff821117156113ac576113ab611357565b5b80604052505050565b5f6113be611332565b90506113ca8282611384565b919050565b5f67ffffffffffffffff8211156113e9576113e8611357565b5b602082029050602081019050919050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f67ffffffffffffffff82111561142457611423611357565b5b61142d82611347565b9050602081019050919050565b828183375f83830152505050565b5f61145a6114558461140a565b6113b5565b90508281526020810184848401111561147657611475611406565b5b61148184828561143a565b509392505050565b5f82601f83011261149d5761149c611343565b5b81356114ad848260208601611448565b91505092915050565b5f8115159050919050565b6114ca816114b6565b81146114d4575f5ffd5b50565b5f813590506114e5816114c1565b92915050565b5f60408284031215611500576114ff6113fe565b5b61150a60406113b5565b90505f82013567ffffffffffffffff81111561152957611528611402565b5b61153584828501611489565b5f830152506020611548848285016114d7565b60208301525092915050565b5f611566611561846113cf565b6113b5565b90508083825260208201905060208402830185811115611589576115886113fa565b5b835b818110156115d057803567ffffffffffffffff8111156115ae576115ad611343565b5b8086016115bb89826114eb565b8552602085019450505060208101905061158b565b5050509392505050565b5f82601f8301126115ee576115ed611343565b5b81356115fe848260208601611554565b91505092915050565b5f819050919050565b61161981611607565b8114611623575f5ffd5b50565b5f8135905061163481611610565b92915050565b5f5f5f606084860312156116515761165061133b565b5b5f84013567ffffffffffffffff81111561166e5761166d61133f565b5b61167a868287016115da565b935050602084013567ffffffffffffffff81111561169b5761169a61133f565b5b6116a786828701611489565b92505060406116b886828701611626565b9150509250925092565b6116cb816114b6565b82525050565b5f6020820190506116e45f8301846116c2565b92915050565b6116f381611607565b82525050565b5f60208201905061170c5f8301846116ea565b92915050565b5f67ffffffffffffffff82169050919050565b61172e81611712565b82525050565b5f6020820190506117475f830184611725565b92915050565b61175681611712565b8114611760575f5ffd5b50565b5f813590506117718161174d565b92915050565b5f5f5f5f6080858703121561178f5761178e61133b565b5b5f85013567ffffffffffffffff8111156117ac576117ab61133f565b5b6117b8878288016115da565b94505060206117c987828801611763565b935050604085013567ffffffffffffffff8111156117ea576117e961133f565b5b6117f687828801611489565b925050606061180787828801611626565b91505092959194509250565b5f61181d8261129f565b9050919050565b61182d81611813565b82525050565b5f6020820190506118465f830184611824565b92915050565b5f602082840312156118615761186061133b565b5b5f61186e84828501611626565b91505092915050565b5f5f6040838503121561188d5761188c61133b565b5b5f83013567ffffffffffffffff8111156118aa576118a961133f565b5b6118b685828601611489565b92505060206118c785828601611763565b9150509250929050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f611903826118d1565b61190d81856118db565b935061191d8185602086016118eb565b61192681611347565b840191505092915050565b5f6020820190508181035f83015261194981846118f9565b905092915050565b5f5f5f5f5f60a0868803121561196a5761196961133b565b5b5f61197788828901611763565b955050602086013567ffffffffffffffff8111156119985761199761133f565b5b6119a488828901611489565b94505060406119b588828901611626565b935050606086013567ffffffffffffffff8111156119d6576119d561133f565b5b6119e2888289016115da565b925050608086013567ffffffffffffffff811115611a0357611a0261133f565b5b611a0f888289016115da565b9150509295509295909350565b5f5f5f5f60808587031215611a3457611a3361133b565b5b5f85013567ffffffffffffffff811115611a5157611a5061133f565b5b611a5d878288016115da565b9450506020611a6e87828801611763565b9350506040611a7f87828801611626565b9250506060611a9087828801611626565b91505092959194509250565b5f7fff0000000000000000000000000000000000000000000000000000000000000082169050919050565b5f819050919050565b611ae1611adc82611a9c565b611ac7565b82525050565b5f81905092915050565b5f611afb826118d1565b611b058185611ae7565b9350611b158185602086016118eb565b80840191505092915050565b5f611b2c8285611ad0565b600182019150611b3c8284611af1565b91508190509392505050565b5f611b538284611af1565b915081905092915050565b5f81519050611b6c81611610565b92915050565b5f60208284031215611b8757611b8661133b565b5b5f611b9484828501611b5e565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f819050919050565b611be4611bdf82611607565b611bca565b82525050565b5f611bf58286611ad0565b600182019150611c058285611af1565b9150611c118284611bd3565b602082019150819050949350505050565b5f611c2d8286611ad0565b600182019150611c3d8285611bd3565b602082019150611c4d8284611af1565b9150819050949350505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f63ffffffff82169050919050565b5f611ca082611c87565b915063ffffffff8203611cb657611cb5611c5a565b5b600182019050919050565b611cca8161126e565b8114611cd4575f5ffd5b50565b5f81519050611ce581611cc1565b92915050565b5f60208284031215611d0057611cff61133b565b5b5f611d0d84828501611cd7565b91505092915050565b5f8160c01b9050919050565b5f611d2c82611d16565b9050919050565b611d44611d3f82611712565b611d22565b82525050565b5f611d558285611d33565b600882019150611d658284611af1565b91508190509392505050565b5f606082019050611d845f830186611824565b611d916020830185611824565b611d9e6040830184611277565b949350505050565b5f81519050611db4816114c1565b92915050565b5f60208284031215611dcf57611dce61133b565b5b5f611ddc84828501611da6565b91505092915050565b5f611def82611712565b9150611dfa83611712565b9250828202611e0881611712565b9150808214611e1a57611e19611c5a565b5b5092915050565b5f611e2c8286611bd3565b602082019150611e3c8285611d33565b600882019150611e4c8284611d33565b600882019150819050949350505050565b5f611e6782611712565b9150611e7283611712565b9250828201905067ffffffffffffffff811115611e9257611e91611c5a565b5b92915050565b5f611ea282611712565b915067ffffffffffffffff8203611ebc57611ebb611c5a565b5b600182019050919050565b5f611ed18261126e565b9150611edc8361126e565b9250828201905080821115611ef457611ef3611c5a565b5b92915050565b5f604082019050611f0d5f830185611824565b611f1a6020830184611277565b9392505050565b5f611f2c8285611d33565b600882019150611f3c8284611bd3565b602082019150819050939250505056fea26469706673582212202e5eeff9919cdc3bbae7c3941d47408cba4ee157bce6aa3a94113c53c6373e3764736f6c63430008210033",
}

// Tredd is an auto generated Go binding around an Ethereum contract.
type Tredd struct {
	abi abi.ABI
}

// NewTredd creates a new instance of Tredd.
func NewTredd() *Tredd {
	parsed, err := TreddMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Tredd{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Tredd) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address seller, address tokenType, uint256 amount, uint256 collateral, bytes32 clearRoot, bytes32 cipherRoot, uint64 revealDeadline, uint64 refundDeadline) payable returns()
func (tredd *Tredd) PackConstructor(seller common.Address, tokenType common.Address, amount *big.Int, collateral *big.Int, clearRoot [32]byte, cipherRoot [32]byte, revealDeadline uint64, refundDeadline uint64) []byte {
	enc, err := tredd.abi.Pack("", seller, tokenType, amount, collateral, clearRoot, cipherRoot, revealDeadline, refundDeadline)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCancel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xea8a1af0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancel() returns()
func (tredd *Tredd) PackCancel() []byte {
	enc, err := tredd.abi.Pack("cancel")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCancel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xea8a1af0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cancel() returns()
func (tredd *Tredd) TryPackCancel() ([]byte, error) {
	return tredd.abi.Pack("cancel")
}

// PackCheckProof is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1235ffeb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkProof((bytes,bool)[] steps, bytes leaf, bytes32 want) pure returns(bool)
func (tredd *Tredd) PackCheckProof(steps []TreddProofStep, leaf []byte, want [32]byte) []byte {
	enc, err := tredd.abi.Pack("checkProof", steps, leaf, want)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckProof is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1235ffeb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkProof((bytes,bool)[] steps, bytes leaf, bytes32 want) pure returns(bool)
func (tredd *Tredd) TryPackCheckProof(steps []TreddProofStep, leaf []byte, want [32]byte) ([]byte, error) {
	return tredd.abi.Pack("checkProof", steps, leaf, want)
}

// UnpackCheckProof is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1235ffeb.
//
// Solidity: function checkProof((bytes,bool)[] steps, bytes leaf, bytes32 want) pure returns(bool)
func (tredd *Tredd) UnpackCheckProof(data []byte) (bool, error) {
	out, err := tredd.abi.Unpack("checkProof", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackCheckProofWithPrefixedChunk is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x33bbe2a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkProofWithPrefixedChunk((bytes,bool)[] steps, uint64 prefix, bytes chunk, bytes32 want) pure returns(bool)
func (tredd *Tredd) PackCheckProofWithPrefixedChunk(steps []TreddProofStep, prefix uint64, chunk []byte, want [32]byte) []byte {
	enc, err := tredd.abi.Pack("checkProofWithPrefixedChunk", steps, prefix, chunk, want)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckProofWithPrefixedChunk is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x33bbe2a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkProofWithPrefixedChunk((bytes,bool)[] steps, uint64 prefix, bytes chunk, bytes32 want) pure returns(bool)
func (tredd *Tredd) TryPackCheckProofWithPrefixedChunk(steps []TreddProofStep, prefix uint64, chunk []byte, want [32]byte) ([]byte, error) {
	return tredd.abi.Pack("checkProofWithPrefixedChunk", steps, prefix, chunk, want)
}

// UnpackCheckProofWithPrefixedChunk is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x33bbe2a7.
//
// Solidity: function checkProofWithPrefixedChunk((bytes,bool)[] steps, uint64 prefix, bytes chunk, bytes32 want) pure returns(bool)
func (tredd *Tredd) UnpackCheckProofWithPrefixedChunk(data []byte) (bool, error) {
	out, err := tredd.abi.Unpack("checkProofWithPrefixedChunk", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackCheckProofWithPrefixedHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfc6210c5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkProofWithPrefixedHash((bytes,bool)[] steps, uint64 prefix, bytes32 hash, bytes32 want) pure returns(bool)
func (tredd *Tredd) PackCheckProofWithPrefixedHash(steps []TreddProofStep, prefix uint64, hash [32]byte, want [32]byte) []byte {
	enc, err := tredd.abi.Pack("checkProofWithPrefixedHash", steps, prefix, hash, want)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckProofWithPrefixedHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfc6210c5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkProofWithPrefixedHash((bytes,bool)[] steps, uint64 prefix, bytes32 hash, bytes32 want) pure returns(bool)
func (tredd *Tredd) TryPackCheckProofWithPrefixedHash(steps []TreddProofStep, prefix uint64, hash [32]byte, want [32]byte) ([]byte, error) {
	return tredd.abi.Pack("checkProofWithPrefixedHash", steps, prefix, hash, want)
}

// UnpackCheckProofWithPrefixedHash is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfc6210c5.
//
// Solidity: function checkProofWithPrefixedHash((bytes,bool)[] steps, uint64 prefix, bytes32 hash, bytes32 want) pure returns(bool)
func (tredd *Tredd) UnpackCheckProofWithPrefixedHash(data []byte) (bool, error) {
	out, err := tredd.abi.Unpack("checkProofWithPrefixedHash", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackClaimPayment is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7dea2f2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function claimPayment() returns()
func (tredd *Tredd) PackClaimPayment() []byte {
	enc, err := tredd.abi.Pack("claimPayment")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClaimPayment is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7dea2f2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function claimPayment() returns()
func (tredd *Tredd) TryPackClaimPayment() ([]byte, error) {
	return tredd.abi.Pack("claimPayment")
}

// PackDecrypt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1598968.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decrypt(bytes chunk, uint64 index) view returns(bytes)
func (tredd *Tredd) PackDecrypt(chunk []byte, index uint64) []byte {
	enc, err := tredd.abi.Pack("decrypt", chunk, index)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDecrypt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1598968.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function decrypt(bytes chunk, uint64 index) view returns(bytes)
func (tredd *Tredd) TryPackDecrypt(chunk []byte, index uint64) ([]byte, error) {
	return tredd.abi.Pack("decrypt", chunk, index)
}

// UnpackDecrypt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa1598968.
//
// Solidity: function decrypt(bytes chunk, uint64 index) view returns(bytes)
func (tredd *Tredd) UnpackDecrypt(data []byte) ([]byte, error) {
	out, err := tredd.abi.Unpack("decrypt", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackMAmount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d966e7d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mAmount() view returns(uint256)
func (tredd *Tredd) PackMAmount() []byte {
	enc, err := tredd.abi.Pack("mAmount")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAmount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d966e7d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mAmount() view returns(uint256)
func (tredd *Tredd) TryPackMAmount() ([]byte, error) {
	return tredd.abi.Pack("mAmount")
}

// UnpackMAmount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7d966e7d.
//
// Solidity: function mAmount() view returns(uint256)
func (tredd *Tredd) UnpackMAmount(data []byte) (*big.Int, error) {
	out, err := tredd.abi.Unpack("mAmount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMBuyer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x649bfb36.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mBuyer() view returns(address)
func (tredd *Tredd) PackMBuyer() []byte {
	enc, err := tredd.abi.Pack("mBuyer")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMBuyer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x649bfb36.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mBuyer() view returns(address)
func (tredd *Tredd) TryPackMBuyer() ([]byte, error) {
	return tredd.abi.Pack("mBuyer")
}

// UnpackMBuyer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x649bfb36.
//
// Solidity: function mBuyer() view returns(address)
func (tredd *Tredd) UnpackMBuyer(data []byte) (common.Address, error) {
	out, err := tredd.abi.Unpack("mBuyer", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackMCipherRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1d595ee7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mCipherRoot() view returns(bytes32)
func (tredd *Tredd) PackMCipherRoot() []byte {
	enc, err := tredd.abi.Pack("mCipherRoot")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMCipherRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1d595ee7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mCipherRoot() view returns(bytes32)
func (tredd *Tredd) TryPackMCipherRoot() ([]byte, error) {
	return tredd.abi.Pack("mCipherRoot")
}

// UnpackMCipherRoot is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1d595ee7.
//
// Solidity: function mCipherRoot() view returns(bytes32)
func (tredd *Tredd) UnpackMCipherRoot(data []byte) ([32]byte, error) {
	out, err := tredd.abi.Unpack("mCipherRoot", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackMClearRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21b0ae82.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mClearRoot() view returns(bytes32)
func (tredd *Tredd) PackMClearRoot() []byte {
	enc, err := tredd.abi.Pack("mClearRoot")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMClearRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21b0ae82.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mClearRoot() view returns(bytes32)
func (tredd *Tredd) TryPackMClearRoot() ([]byte, error) {
	return tredd.abi.Pack("mClearRoot")
}

// UnpackMClearRoot is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x21b0ae82.
//
// Solidity: function mClearRoot() view returns(bytes32)
func (tredd *Tredd) UnpackMClearRoot(data []byte) ([32]byte, error) {
	out, err := tredd.abi.Unpack("mClearRoot", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackMCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095e4c20.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mCollateral() view returns(uint256)
func (tredd *Tredd) PackMCollateral() []byte {
	enc, err := tredd.abi.Pack("mCollateral")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095e4c20.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mCollateral() view returns(uint256)
func (tredd *Tredd) TryPackMCollateral() ([]byte, error) {
	return tredd.abi.Pack("mCollateral")
}

// UnpackMCollateral is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095e4c20.
//
// Solidity: function mCollateral() view returns(uint256)
func (tredd *Tredd) UnpackMCollateral(data []byte) (*big.Int, error) {
	out, err := tredd.abi.Unpack("mCollateral", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMDecryptionKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9067c7a9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mDecryptionKey() view returns(bytes32)
func (tredd *Tredd) PackMDecryptionKey() []byte {
	enc, err := tredd.abi.Pack("mDecryptionKey")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMDecryptionKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9067c7a9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mDecryptionKey() view returns(bytes32)
func (tredd *Tredd) TryPackMDecryptionKey() ([]byte, error) {
	return tredd.abi.Pack("mDecryptionKey")
}

// UnpackMDecryptionKey is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9067c7a9.
//
// Solidity: function mDecryptionKey() view returns(bytes32)
func (tredd *Tredd) UnpackMDecryptionKey(data []byte) ([32]byte, error) {
	out, err := tredd.abi.Unpack("mDecryptionKey", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackMRefundDeadline is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2df6a9da.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mRefundDeadline() view returns(uint64)
func (tredd *Tredd) PackMRefundDeadline() []byte {
	enc, err := tredd.abi.Pack("mRefundDeadline")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMRefundDeadline is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2df6a9da.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mRefundDeadline() view returns(uint64)
func (tredd *Tredd) TryPackMRefundDeadline() ([]byte, error) {
	return tredd.abi.Pack("mRefundDeadline")
}

// UnpackMRefundDeadline is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2df6a9da.
//
// Solidity: function mRefundDeadline() view returns(uint64)
func (tredd *Tredd) UnpackMRefundDeadline(data []byte) (uint64, error) {
	out, err := tredd.abi.Unpack("mRefundDeadline", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackMRevealDeadline is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x61a5ab22.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mRevealDeadline() view returns(uint64)
func (tredd *Tredd) PackMRevealDeadline() []byte {
	enc, err := tredd.abi.Pack("mRevealDeadline")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMRevealDeadline is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x61a5ab22.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mRevealDeadline() view returns(uint64)
func (tredd *Tredd) TryPackMRevealDeadline() ([]byte, error) {
	return tredd.abi.Pack("mRevealDeadline")
}

// UnpackMRevealDeadline is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x61a5ab22.
//
// Solidity: function mRevealDeadline() view returns(uint64)
func (tredd *Tredd) UnpackMRevealDeadline(data []byte) (uint64, error) {
	out, err := tredd.abi.Unpack("mRevealDeadline", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackMRevealed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x54b53436.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mRevealed() view returns(bool)
func (tredd *Tredd) PackMRevealed() []byte {
	enc, err := tredd.abi.Pack("mRevealed")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMRevealed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x54b53436.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mRevealed() view returns(bool)
func (tredd *Tredd) TryPackMRevealed() ([]byte, error) {
	return tredd.abi.Pack("mRevealed")
}

// UnpackMRevealed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x54b53436.
//
// Solidity: function mRevealed() view returns(bool)
func (tredd *Tredd) UnpackMRevealed(data []byte) (bool, error) {
	out, err := tredd.abi.Unpack("mRevealed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackMSeller is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8bae87ba.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mSeller() view returns(address)
func (tredd *Tredd) PackMSeller() []byte {
	enc, err := tredd.abi.Pack("mSeller")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMSeller is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8bae87ba.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mSeller() view returns(address)
func (tredd *Tredd) TryPackMSeller() ([]byte, error) {
	return tredd.abi.Pack("mSeller")
}

// UnpackMSeller is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8bae87ba.
//
// Solidity: function mSeller() view returns(address)
func (tredd *Tredd) UnpackMSeller(data []byte) (common.Address, error) {
	out, err := tredd.abi.Unpack("mSeller", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackMTokenType is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0c590dce.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mTokenType() view returns(address)
func (tredd *Tredd) PackMTokenType() []byte {
	enc, err := tredd.abi.Pack("mTokenType")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMTokenType is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0c590dce.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mTokenType() view returns(address)
func (tredd *Tredd) TryPackMTokenType() ([]byte, error) {
	return tredd.abi.Pack("mTokenType")
}

// UnpackMTokenType is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0c590dce.
//
// Solidity: function mTokenType() view returns(address)
func (tredd *Tredd) UnpackMTokenType(data []byte) (common.Address, error) {
	out, err := tredd.abi.Unpack("mTokenType", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPaid is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x295b4e17.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function paid() view returns(uint256)
func (tredd *Tredd) PackPaid() []byte {
	enc, err := tredd.abi.Pack("paid")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPaid is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x295b4e17.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function paid() view returns(uint256)
func (tredd *Tredd) TryPackPaid() ([]byte, error) {
	return tredd.abi.Pack("paid")
}

// UnpackPaid is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x295b4e17.
//
// Solidity: function paid() view returns(uint256)
func (tredd *Tredd) UnpackPaid(data []byte) (*big.Int, error) {
	out, err := tredd.abi.Unpack("paid", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRefund is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xac280f3d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function refund(uint64 index, bytes cipherChunk, bytes32 clearHash, (bytes,bool)[] cipherProof, (bytes,bool)[] clearProof) returns()
func (tredd *Tredd) PackRefund(index uint64, cipherChunk []byte, clearHash [32]byte, cipherProof []TreddProofStep, clearProof []TreddProofStep) []byte {
	enc, err := tredd.abi.Pack("refund", index, cipherChunk, clearHash, cipherProof, clearProof)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRefund is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xac280f3d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function refund(uint64 index, bytes cipherChunk, bytes32 clearHash, (bytes,bool)[] cipherProof, (bytes,bool)[] clearProof) returns()
func (tredd *Tredd) TryPackRefund(index uint64, cipherChunk []byte, clearHash [32]byte, cipherProof []TreddProofStep, clearProof []TreddProofStep) ([]byte, error) {
	return tredd.abi.Pack("refund", index, cipherChunk, clearHash, cipherProof, clearProof)
}

// PackReveal is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x701fd0f1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function reveal(bytes32 decryptionKey) payable returns()
func (tredd *Tredd) PackReveal(decryptionKey [32]byte) []byte {
	enc, err := tredd.abi.Pack("reveal", decryptionKey)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackReveal is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x701fd0f1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function reveal(bytes32 decryptionKey) payable returns()
func (tredd *Tredd) TryPackReveal(decryptionKey [32]byte) ([]byte, error) {
	return tredd.abi.Pack("reveal", decryptionKey)
}

// TreddEvDecryptionKey represents a evDecryptionKey event raised by the Tredd contract.
type TreddEvDecryptionKey struct {
	DecryptionKey [32]byte
	Raw           *types.Log // Blockchain specific contextual infos
}

const TreddEvDecryptionKeyEventName = "evDecryptionKey"

// ContractEventName returns the user-defined event name.
func (TreddEvDecryptionKey) ContractEventName() string {
	return TreddEvDecryptionKeyEventName
}

// UnpackEvDecryptionKeyEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event evDecryptionKey(bytes32 decryptionKey)
func (tredd *Tredd) UnpackEvDecryptionKeyEvent(log *types.Log) (*TreddEvDecryptionKey, error) {
	event := "evDecryptionKey"
	if len(log.Topics) == 0 || log.Topics[0] != tredd.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(TreddEvDecryptionKey)
	if len(log.Data) > 0 {
		if err := tredd.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range tredd.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}
