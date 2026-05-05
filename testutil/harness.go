package testutil

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/bobg/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/ethereum/go-ethereum/node"

	"github.com/bobg/tredd/contract"
)

// Use these constants for reproducibility
// (rather than generating random new keys each time).
const (
	secp256k1JSON     = `{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256}`
	buyerKeyJSON      = `{"X":106137327885008459029433685034979965204777290812390373077765777766929045630616,"Y":64565985154334530541640099111240376671268158415813158379126686844588611988459,"D":90769587056954039490056047683741742231702779454899233049594067387646290264706}`
	sellerKeyJSON     = `{"X":17584145466380143975510816014412290760093596753774943791675900103048620655792,"Y":40623230501215950686519909283241339443538574477742632018201392342778703371797,"D":55123640651322237403179227776230301641416033286033349057271134608478213089253}`
	decryptionKeyHex  = "6dcf7dc83d36b7e36fe66c4bd25f4ac9bec1e4bc231e423030b9ad21024ed7ff"
	udhrClearRootHex  = "bfb11979628ce73032a8ba5f02b27309c14d337b464be8298fdb3413fb52c8ce"
	udhrCipherRootHex = "66a699dd914183add184eab58435b2a2919018e4b749b32e3efe726f65c5887e"
)

var (
	DecryptionKey [32]byte
	ClearRoot     [32]byte
	CipherRoot    [32]byte
)

const (
	RevealDeadlineSecs = 600
	RefundDeadlineSecs = 1200
	StartingBalance    = 1000000000
)

func init() {
	_, err := hex.Decode(DecryptionKey[:], []byte(decryptionKeyHex))
	if err != nil {
		panic(err)
	}

	_, err = hex.Decode(ClearRoot[:], []byte(udhrClearRootHex))
	if err != nil {
		panic(err)
	}

	_, err = hex.Decode(CipherRoot[:], []byte(udhrCipherRootHex))
	if err != nil {
		panic(err)
	}
}

type Harness struct {
	Buyer, Seller                  *bind.TransactOpts
	Sim                            *simulated.Backend
	Client                         *testClient
	RevealDeadline, RefundDeadline time.Time
	ContractAddr                   common.Address // only set after Harness.Deploy is called
	Contract                       *bind.BoundContract
	BuyerBalance, SellerBalance    uint64 // caller updates these then calls CheckBalances
}

// testClient wraps a simulated backend and its client together, exposing
// both the ContractBackend/DeployBackend methods (from the client) and
// the Commit method (from the backend) that is needed by waitMined.
type testClient struct {
	*simulated.Backend
	simulated.Client
}

func NewHarness() (*Harness, error) {
	var curve secp256k1.BitCurve
	err := json.Unmarshal([]byte(secp256k1JSON), &curve)
	if err != nil {
		return nil, err
	}

	var buyerKey, sellerKey ecdsa.PrivateKey

	err = json.Unmarshal([]byte(buyerKeyJSON), &buyerKey)
	if err != nil {
		return nil, err
	}
	buyerKey.Curve = &curve
	buyer := bind.NewKeyedTransactor(&buyerKey, big.NewInt(1337))
	buyer.GasPrice = big.NewInt(1)

	err = json.Unmarshal([]byte(sellerKeyJSON), &sellerKey)
	if err != nil {
		return nil, err
	}
	sellerKey.Curve = &curve
	seller := bind.NewKeyedTransactor(&sellerKey, big.NewInt(1337))
	seller.GasPrice = big.NewInt(1)

	alloc := core.GenesisAlloc{
		buyer.From:  core.GenesisAccount{Balance: big.NewInt(StartingBalance)},
		seller.From: core.GenesisAccount{Balance: big.NewInt(StartingBalance)},
	}

	// Use simulated.NewBackend directly so we can set the genesis base fee to
	// zero.  This lets us use legacy transactions with GasPrice=1 wei, which
	// keeps the per-gas cost at 1 wei and preserves the balance arithmetic in
	// tests.  WithMinerMinTip(1) sets the minimum inclusion tip to 1 wei.
	sim := simulated.NewBackend(
		types.GenesisAlloc(alloc),
		simulated.WithBlockGasLimit(30_000_000),
		simulated.WithMinerMinTip(big.NewInt(1)),
		func(_ *node.Config, ethConf *ethconfig.Config) {
			ethConf.Genesis.BaseFee = big.NewInt(0)
		},
	)

	now := time.Now()

	return &Harness{
		Buyer:          buyer,
		Seller:         seller,
		Sim:            sim,
		Client:         &testClient{Backend: sim, Client: sim.Client()},
		RevealDeadline: now.Add(RevealDeadlineSecs * time.Second),
		RefundDeadline: now.Add(RefundDeadlineSecs * time.Second),
		BuyerBalance:   StartingBalance,
		SellerBalance:  StartingBalance,
	}, nil
}

var (
	big2 = big.NewInt(2)
	big3 = big.NewInt(3)
)

var treddABI = contract.NewTredd()

func (h *Harness) Deploy(ctx context.Context) error {
	constructorInput := treddABI.PackConstructor(h.Seller.From, common.Address{}, big3, big2, ClearRoot, CipherRoot, uint64(h.RevealDeadline.Unix()), uint64(h.RefundDeadline.Unix()))
	addr, deployTx, err := bind.DeployContract(h.Buyer, common.FromHex(contract.TreddMetaData.Bin), h.Client, constructorInput)
	if err != nil {
		return errors.Wrap(err, "deploying tredd contract")
	}
	h.Sim.Commit()

	// Wait for deployment to be mined.
	_, err = bind.WaitMined(ctx, h.Client, deployTx.Hash())
	if err != nil {
		return errors.Wrap(err, "waiting for contract deployment")
	}

	// Transfer the buyer payment to the contract (ETH path: send ETH to contract).
	txOpts := *h.Buyer
	txOpts.Value = big3
	instance := treddABI.Instance(h.Client, addr)
	transferTx, err := instance.Transfer(&txOpts)
	if err != nil {
		return errors.Wrap(err, "transferring buyer payment to contract")
	}
	h.Sim.Commit()

	_, err = bind.WaitMined(ctx, h.Client, transferTx.Hash())
	if err != nil {
		return errors.Wrap(err, "waiting for payment transfer")
	}

	h.ContractAddr = addr
	h.Contract = instance
	return nil
}

func (h *Harness) Balances(ctx context.Context) (buyer, seller *big.Int, err error) {
	buyer, err = h.Client.BalanceAt(ctx, h.Buyer.From, nil)
	if err != nil {
		return nil, nil, err
	}
	seller, err = h.Client.BalanceAt(ctx, h.Seller.From, nil)
	return buyer, seller, err
}

func (h *Harness) CheckBalances(ctx context.Context) error {
	gotBuyer, gotSeller, err := h.Balances(ctx)
	if err != nil {
		return err
	}
	wantBuyer := big.NewInt(int64(h.BuyerBalance))
	if gotBuyer.Cmp(wantBuyer) != 0 {
		return fmt.Errorf("got buyer balance %s, want %s", gotBuyer, wantBuyer)
	}
	wantSeller := big.NewInt(int64(h.SellerBalance))
	if gotSeller.Cmp(wantSeller) != 0 {
		return fmt.Errorf("got seller balance %s, want %s", gotSeller, wantSeller)
	}
	return nil
}
