package tedd

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"math"
	"os"
	"testing"
	"time"

	"github.com/chain/txvm/crypto/ed25519"
	"github.com/chain/txvm/protocol/bc"
	"github.com/chain/txvm/protocol/txvm"
)

var nexthashval = make([]byte, 32)

func nexthash() []byte {
	buf := sha256.Sum256(nexthashval)
	nexthashval = buf[:]
	return nexthashval
}

func TestTx(t *testing.T) {
	const (
		clearRootHex  = "d78b83cba3f32b8eb56831c834f6819d383c97637c2ef95cfc71339a2da2d94f"
		cipherRootHex = "684d6d5652e44d45452d3c56ae5d229f701c67205a03d5c61de5a2a2134e5a0e"

		buyerPrvHex  = "1a66ef435b1bd836a3ef4cf4fc8ef9e08c83a01fbcfc72c165054f4e9edd56abe3a41fdffd70ee3fdf5f8561d497b5c3735802aad4e782e29f3ed3162c325b5a"
		sellerPrvHex = "528963ef0aeb416f29206807e1bdb11e94fbbfb67cd9b119495b422cfb173c2b02710442b1eb0206c7228ce6a6ceb72a93d0d3bb6d89de5fc0d7e5bf869c437e"
	)

	ctx := context.Background()

	var key [32]byte
	_, err := hex.Decode(key[:], []byte(testKeyHex))
	if err != nil {
		t.Fatal(err)
	}

	var clearRoot [32]byte
	_, err = hex.Decode(clearRoot[:], []byte(clearRootHex))
	if err != nil {
		t.Fatal(err)
	}

	var cipherRoot [32]byte
	_, err = hex.Decode(cipherRoot[:], []byte(cipherRootHex))
	if err != nil {
		t.Fatal(err)
	}

	prvBytes, err := hex.DecodeString(buyerPrvHex)
	if err != nil {
		t.Fatal(err)
	}
	buyerPrv := ed25519.PrivateKey(prvBytes)
	buyer := buyerPrv.Public().(ed25519.PublicKey)

	prvBytes, err = hex.DecodeString(sellerPrvHex)
	if err != nil {
		t.Fatal(err)
	}
	sellerPrv := ed25519.PrivateKey(prvBytes)
	seller := sellerPrv.Public().(ed25519.PublicKey)

	assetID := bc.HashFromBytes(nexthash())

	reserver := &testReserver{
		utxos: []UTXO{
			&testUTXO{
				amount:  7,
				assetID: assetID,
				anchor:  nexthash(),
			},
			&testUTXO{
				amount:  5,
				assetID: assetID,
				anchor:  nexthash(),
			},
		},
	}

	revealDeadline := time.Unix(233400000, 0)
	refundDeadline := revealDeadline.Add(time.Hour)

	signer := func(msg []byte) ([]byte, error) {
		return ed25519.Sign(buyerPrv, msg), nil
	}

	partial, err := ProposePayment(ctx, buyer, 10, assetID, clearRoot, cipherRoot, revealDeadline, refundDeadline, reserver, signer)
	if err != nil {
		t.Fatal(err)
	}

	reserver = &testReserver{
		utxos: []UTXO{
			&testUTXO{
				amount:  9,
				assetID: assetID,
				anchor:  nexthash(),
			},
			&testUTXO{
				amount:  5,
				assetID: assetID,
				anchor:  nexthash(),
			},
		},
	}

	signer = func(msg []byte) ([]byte, error) {
		return ed25519.Sign(sellerPrv, msg), nil
	}

	complete, err := RevealKey(ctx, partial, seller, key, 10, assetID, reserver, signer, clearRoot, cipherRoot, revealDeadline, refundDeadline)
	if err != nil {
		t.Fatal(err)
	}

	vm, err := txvm.Validate(complete, 3, math.MaxInt64, txvm.Trace(os.Stdout))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("complete tx: %x\n", complete)
	t.Logf("runlimit consumed: %d\n", math.MaxInt64-vm.Runlimit())
}
