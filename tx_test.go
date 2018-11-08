package tedd

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
	"math"
	"os"
	"testing"
	"time"

	"github.com/bobg/merkle"
	"github.com/chain/txvm/crypto/ed25519"
	"github.com/chain/txvm/errors"
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

	var anchor [32]byte
	copy(anchor[:], vm.Log[len(vm.Log)-5][2].(txvm.Bytes))
	outputID := []byte(vm.Log[len(vm.Log)-2][2].(txvm.Bytes))

	r := &Redeem{
		RefundDeadline: refundDeadline,
		Buyer:          buyer,
		Seller:         seller,
		Amount:         20,
		AssetID:        assetID,
		Anchor:         anchor,
		CipherRoot:     cipherRoot,
		ClearRoot:      clearRoot,
		Key:            key,
	}

	claimPaymentProg, err := ClaimPayment(r)
	if err != nil {
		t.Fatal(err)
	}

	vm, err = txvm.Validate(claimPaymentProg, 3, math.MaxInt64)
	if err != nil {
		t.Fatal(err)
	}
	if got := []byte(vm.Log[0][2].(txvm.Bytes)); !bytes.Equal(got, outputID) {
		t.Errorf("on input, got outputID %x, want %x", got, outputID)
	}

	var (
		clearTree, cipherTree *merkle.Tree
		refhash, refchunk     []byte
	)

	f, err := os.Open("testdata/commonsense.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	hasher := sha256.New()

	for index := uint64(0); ; index++ {
		var chunk [chunkSize + binary.MaxVarintLen64]byte
		m := binary.PutUvarint(chunk[:], index)
		n, err := io.ReadFull(f, chunk[m:m+chunkSize])
		if err == io.EOF {
			// "The error is EOF only if no bytes were read."
			break
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			t.Fatal(err)
		}

		var h [32 + binary.MaxVarintLen64]byte
		binary.PutUvarint(h[:], index)
		merkle.LeafHash(hasher, h[:m], chunk[:m+n])

		if index == 0 {
			refhash = make([]byte, 32)
			copy(refhash[:], h[m:m+32])
			clearTree = merkle.NewProofTree(sha256.New(), h[:m+32])
		}
		clearTree.Add(h[:m+32])
		crypt(key, chunk[m:m+n], index)
		if index == 0 {
			refchunk = make([]byte, n)
			copy(refchunk, chunk[m:m+n])
			cipherTree = merkle.NewProofTree(sha256.New(), chunk[:m+n])
		}
		cipherTree.Add(chunk[:m+n])
	}
	clearProof := clearTree.Proof()
	cipherProof := cipherTree.Proof()

	// With the right key and clear hash and cipher chunk,
	// it should not be possible to get a refund.
	claimRefundProg, err := ClaimRefund(r, 0, refchunk, refhash, cipherProof, clearProof)
	if err != nil {
		t.Fatal(err)
	}

	vm, err = txvm.Validate(claimRefundProg, 3, math.MaxInt64)
	if errors.Root(err) != txvm.ErrVerifyFail {
		t.Errorf("got error %v, want %s", err, txvm.ErrVerifyFail)
	}

	// With the wrong key, on the other hand...
	r.Key[0] ^= 1

	claimRefundProg, err = ClaimRefund(r, 0, refchunk, refhash, cipherProof, clearProof)
	if err != nil {
		t.Fatal(err)
	}
	vm, err = txvm.Validate(claimRefundProg, 3, math.MaxInt64)
	if err != nil {
		t.Error(err)
	}
}
