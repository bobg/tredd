package testutil

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pkg/errors"

	"github.com/bobg/tredd/contract"
)

// Use these constants for reproducibility
// (rather than generating random new keys each time).
const (
	secp256k1JSON     = `{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256}`
	buyerKeyJSON      = `{"X":106137327885008459029433685034979965204777290812390373077765777766929045630616,"Y":64565985154334530541640099111240376671268158415813158379126686844588611988459,"D":90769587056954039490056047683741742231702779454899233049594067387646290264706}`
	sellerKeyJSON     = `{"X":17584145466380143975510816014412290760093596753774943791675900103048620655792,"Y":40623230501215950686519909283241339443538574477742632018201392342778703371797,"D":55123640651322237403179227776230301641416033286033349057271134608478213089253}`
	decryptionKeyHex  = "6dcf7dc83d36b7e36fe66c4bd25f4ac9bec1e4bc231e423030b9ad21024ed7ff"
	udhrClearRootHex  = "3c1ec1141f8c59544fbaa3fe5eaf41323bd9e4ced48c7d2c4bd084f8015a83bb"
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
	Client                         *backends.SimulatedBackend
	RevealDeadline, RefundDeadline time.Time
	ContractAddr                   common.Address // only set after Harness.Deploy is called
	Contract                       *contract.Tredd
	BuyerBalance, SellerBalance    uint64 // caller updates these then calls CheckBalances
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
	buyer := bind.NewKeyedTransactor(&buyerKey)

	err = json.Unmarshal([]byte(sellerKeyJSON), &sellerKey)
	if err != nil {
		return nil, err
	}
	sellerKey.Curve = &curve
	seller := bind.NewKeyedTransactor(&sellerKey)

	alloc := core.GenesisAlloc{
		buyer.From:  core.GenesisAccount{Balance: big.NewInt(StartingBalance)},
		seller.From: core.GenesisAccount{Balance: big.NewInt(StartingBalance)},
	}

	client := backends.NewSimulatedBackend(alloc, 4712388) // This number comes from https://goethereumbook.org/client-simulated/

	now := time.Unix(0, 0)

	return &Harness{
		Buyer:          buyer,
		Seller:         seller,
		Client:         client,
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

func (h *Harness) Deploy(ctx context.Context) error {
	addr, _, con, err := contract.DeployTredd(h.Buyer, h.Client, h.Seller.From, common.Address{}, big3, big2, ClearRoot, CipherRoot, h.RevealDeadline.Unix(), h.RefundDeadline.Unix())
	if err != nil {
		return errors.Wrap(err, "deploying tredd contract")
	}

	txOpts := *h.Buyer
	txOpts.Value = big3
	raw := &contract.TreddRaw{Contract: con}

	_, err = raw.Transfer(&txOpts)
	if err != nil {
		return errors.Wrap(err, "transfering buyer payment to contract")
	}
	h.Client.Commit()

	h.ContractAddr = addr
	h.Contract = con
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
