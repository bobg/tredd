package testutil

import (
	"crypto/ecdsa"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// Use these constants for reproducibility
// (rather than generating random new keys each time).
const (
	secp256k1JSON = `{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256}`
	buyerKeyJSON  = `{"X":106137327885008459029433685034979965204777290812390373077765777766929045630616,"Y":64565985154334530541640099111240376671268158415813158379126686844588611988459,"D":90769587056954039490056047683741742231702779454899233049594067387646290264706}`
	sellerKeyJSON = `{"X":17584145466380143975510816014412290760093596753774943791675900103048620655792,"Y":40623230501215950686519909283241339443538574477742632018201392342778703371797,"D":55123640651322237403179227776230301641416033286033349057271134608478213089253}`
)

func Harness() (buyer, seller *bind.TransactOpts, client *backends.SimulatedBackend, err error) {
	var curve secp256k1.BitCurve
	err = json.Unmarshal([]byte(secp256k1JSON), &curve)
	if err != nil {
		return nil, nil, nil, err
	}

	var buyerKey, sellerKey ecdsa.PrivateKey

	err = json.Unmarshal([]byte(buyerKeyJSON), &buyerKey)
	if err != nil {
		return nil, nil, nil, err
	}
	buyerKey.Curve = &curve
	buyer = bind.NewKeyedTransactor(&buyerKey)

	err = json.Unmarshal([]byte(sellerKeyJSON), &sellerKey)
	if err != nil {
		return nil, nil, nil, err
	}
	sellerKey.Curve = &curve
	seller = bind.NewKeyedTransactor(&sellerKey)

	alloc := core.GenesisAlloc{
		buyer.From:  core.GenesisAccount{Balance: big.NewInt(1000000000)},
		seller.From: core.GenesisAccount{Balance: big.NewInt(1000000000)},
	}
	client = backends.NewSimulatedBackend(alloc, 4712388) // This number comes from https://goethereumbook.org/client-simulated/

	return buyer, seller, client, nil
}
