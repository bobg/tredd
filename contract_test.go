package tredd

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io"
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

	con, err := harness.Contract()
	if err != nil {
		t.Fatal(err)
	}

	hasher := sha256.New()
	for _, refchunk := range chunks {
		tree := merkle.NewProofTree(hasher, refchunk)
		for _, chunk := range chunks {
			tree.Add(chunk)
		}
		root := tree.Root()
		proof := tree.Proof()

		var wantRootBuf [32]byte
		copy(wantRootBuf[:], root)

		callopts := new(bind.CallOpts)

		ok, err := con.CheckProof(callopts, contract.Proof(proof), refchunk, wantRootBuf)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Error("proof validation failed")
		}

		refchunk[0] ^= 1
		ok, err = con.CheckProof(callopts, contract.Proof(proof), refchunk, wantRootBuf)
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

	err = harness.Reveal(ctx)
	if err != nil {
		t.Fatal(err)
	}

	callopts := new(bind.CallOpts)

	con, err := harness.Contract()
	if err != nil {
		t.Fatal(err)
	}

	got, err := con.Decrypt(callopts, cipher[:], 0)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, clear[:]) {
		t.Error("mismatch")
	}
}
