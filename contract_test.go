package tredd

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/bobg/merkle"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

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

	buyer, seller, client, err := testutil.Harness()
	if err != nil {
		t.Fatal(err)
	}

	_, tx, con, err := DeployTredd(buyer, client, seller.From, common.Address{}, big.NewInt(1), big.NewInt(1), zeroes, zeroes, time.Now().Add(time.Hour).Unix(), time.Now().Add(2*time.Hour).Unix())
	if err != nil {
		t.Fatal(err)
	}

	client.Commit()

	ctx := context.Background()

	_, err = client.TransactionReceipt(ctx, tx.Hash())
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

		ok, err := con.CheckProof(callopts, toTreddProof(proof), refchunk, wantRootBuf)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Error("proof validation failed")
		}

		refchunk[0] ^= 1
		ok, err = con.CheckProof(callopts, toTreddProof(proof), refchunk, wantRootBuf)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Error("proof validation succeeded unexpectedly")
		}
	}
}

func TestDecrypt(t *testing.T) {
	var key [32]byte
	_, err := hex.Decode(key[:], []byte(testKeyHex))
	if err != nil {
		t.Fatal(err)
	}

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

	Crypt(key, cipher[:], 0)
	if bytes.Equal(cipher[:], clear[:]) {
		t.Fatal("encrypting did nothing?!")
	}

	// TODO: test solidity decryption of `cipher` produces `clear`
}
