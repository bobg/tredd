package tredd

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io"
	"math/big"
	"os"
	"testing"

	"github.com/bobg/merkle/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"

	"github.com/bobg/tredd/contract"
	"github.com/bobg/tredd/testutil"
)

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

	for i, refchunk := range chunks {
		var (
			chunkTree = merkle.NewProofTree(sha256.New(), Prefix(uint64(i), refchunk))
			refhash   = sha256.Sum256(refchunk)
			hashTree  = merkle.NewProofTree(sha256.New(), Prefix(uint64(i), refhash[:]))
		)
		for j, chunk := range chunks {
			chunkTree.Add(Prefix(uint64(j), chunk))
			hash := sha256.Sum256(chunk)
			hashTree.Add(Prefix(uint64(j), hash[:]))
		}

		var chunkRoot [32]byte
		copy(chunkRoot[:], chunkTree.Root())
		chunkProof := chunkTree.Proof()

		var hashRoot [32]byte
		copy(hashRoot[:], hashTree.Root())
		hashProof := hashTree.Proof()

		callopts := new(bind.CallOpts)

		ok, err := bind.Call(harness.Contract, callopts, treddABI.PackCheckProofWithPrefixedChunk(contract.Proof(chunkProof), uint64(i), refchunk, chunkRoot), treddABI.UnpackCheckProofWithPrefixedChunk)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Error("chunkTree proof validation failed")
		}

		ok, err = bind.Call(harness.Contract, callopts, treddABI.PackCheckProofWithPrefixedHash(contract.Proof(hashProof), uint64(i), refhash, hashRoot), treddABI.UnpackCheckProofWithPrefixedHash)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Error("hashTree proof validation failed")
		}

		refchunk[0] ^= 1
		ok, err = bind.Call(harness.Contract, callopts, treddABI.PackCheckProofWithPrefixedChunk(contract.Proof(chunkProof), uint64(i), refchunk, chunkRoot), treddABI.UnpackCheckProofWithPrefixedChunk)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Error("chunkTree proof validation succeeded unexpectedly")
		}

		refhash[0] ^= 1
		ok, err = bind.Call(harness.Contract, callopts, treddABI.PackCheckProofWithPrefixedHash(contract.Proof(hashProof), uint64(i), refhash, hashRoot), treddABI.UnpackCheckProofWithPrefixedHash)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Error("hashTree proof validation succeeded unexpectedly")
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
	_, err = bind.Transact(harness.Contract, &txOpts, treddABI.PackReveal(testutil.DecryptionKey))
	if err != nil {
		t.Fatal(err)
	}
	harness.Sim.Commit()

	callopts := new(bind.CallOpts)

	got, err := bind.Call(harness.Contract, callopts, treddABI.PackDecrypt(cipher[:], 0), treddABI.UnpackDecrypt)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, clear[:]) {
		t.Error("mismatch")
	}
}
