package tredd

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"testing"

	"github.com/bobg/merkle"
)

func TestTxVMMerkleCheck(t *testing.T) {
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

	hasher := sha256.New()
	for _, refchunk := range chunks {
		tree := merkle.NewProofTree(hasher, refchunk)
		for _, chunk := range chunks {
			tree.Add(chunk)
		}
		root := tree.Root()
		proof := tree.Proof()

		testMerkleCheck(t, proof, root, refchunk)
	}
}

func testMerkleCheck(t *testing.T, proof merkle.Proof, wantRoot, refchunk []byte) {
	// TODO: test solidity validation of the proof
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
