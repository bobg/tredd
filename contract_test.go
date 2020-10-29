package tredd

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io"
	"math/big"
	"os"
	"testing"

	"github.com/bobg/merkle"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/bobg/tredd/contract"
	"github.com/bobg/tredd/testutil"
)

var zeroes [32]byte

func TestSolidityMerkleCheck(t *testing.T) {
	f, err := os.Open("testdata/udhr.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	const chunksize = 256
	var chunks [][]byte
	for {
		var buf [chunksize]byte
		n, err := io.ReadFull(f, buf[:])
		if err == io.EOF {
			// "The error is EOF only if no bytes were read."
			break
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			t.Fatal(err)
		}
		chunks = append(chunks, buf[:n])
	}

	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	err = harness.Deploy(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, refchunk := range chunks {
		tree := merkle.NewProofTree(sha256.New(), refchunk)
		for _, chunk := range chunks {
			tree.Add(chunk)
		}

		var root [32]byte
		copy(root[:], tree.Root())
		proof := tree.Proof()

		callopts := new(bind.CallOpts)

		var leaf [32]byte

		merkle.LeafHash(sha256.New(), leaf[:0], refchunk)
		ok, err := harness.Contract.CheckProof(callopts, contract.Proof(proof), leaf, root)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Error("proof validation failed")
		}

		leaf[0] ^= 1
		ok, err = harness.Contract.CheckProof(callopts, contract.Proof(proof), leaf, root)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Error("proof validation succeeded unexpectedly")
		}
	}
}

func TestDecrypt(t *testing.T) {
	f, err := os.Open("testdata/udhr.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	const chunksize = 256
	var clear, cipher [chunksize]byte
	_, err = io.ReadFull(f, clear[:])
	if err != nil {
		t.Fatal(err)
	}

	copy(cipher[:], clear[:])

	err = Crypt(testutil.DecryptionKey, cipher[:], 0)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(cipher[:], clear[:]) {
		t.Fatal("encrypting did nothing?!")
	}

	err = Crypt(testutil.DecryptionKey, cipher[:], 0)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(cipher[:], clear[:]) {
		t.Fatal("Crypt(Crypt(clear)) != clear ?!")
	}

	err = Crypt(testutil.DecryptionKey, cipher[:], 0)
	if err != nil {
		t.Fatal(err)
	}

	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	err = harness.Deploy(ctx)
	if err != nil {
		t.Fatal(err)
	}

	txOpts := *harness.Seller
	txOpts.Value = big.NewInt(2)
	_, err = harness.Contract.Reveal(&txOpts, testutil.DecryptionKey)
	if err != nil {
		t.Fatal(err)
	}
	harness.Client.Commit()

	callopts := new(bind.CallOpts)

	got, err := harness.Contract.Decrypt(callopts, cipher[:], 0)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, clear[:]) {
		t.Error("mismatch")
	}
}
