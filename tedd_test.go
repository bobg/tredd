package tedd

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const testKeyHex = "17f9d2125c385c2b7626034a506e524b971d9487daeb688538101c4d7d6d1f2a"

func TestServeGetDecrypt(t *testing.T) {
	const (
		clearRootHex      = "689b09a91f8a3a52fa83f076084878688242222b997a25c62e2ef03d58d50bfc"
		wantCipherRootHex = "684d6d5652e44d45452d3c56ae5d229f701c67205a03d5c61de5a2a2134e5a0e"
	)

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

	f, err := os.Open("testdata/commonsense.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	text, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	served := new(bytes.Buffer)
	cipherRoot, err := Serve(served, bytes.NewReader(text), key)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(cipherRoot) != wantCipherRootHex {
		t.Errorf("got cipher root %x, want %s", cipherRoot, wantCipherRootHex)
	}

	var (
		clearHashes  = new(testChunkStore)
		cipherChunks = new(testChunkStore)
	)

	cipherRoot, err = Get(served, clearRoot, clearHashes, cipherChunks)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(cipherRoot) != wantCipherRootHex {
		t.Errorf("got cipher root %x, want %s", cipherRoot, wantCipherRootHex)
	}

	decrypted := new(bytes.Buffer)
	err = Decrypt(decrypted, clearHashes, cipherChunks, key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decrypted.Bytes(), text) {
		t.Error("text mismatch")
	}
}

type testChunkStore struct {
	chunks [][]byte
}

func (t *testChunkStore) Add(chunk []byte) error {
	dup := make([]byte, len(chunk))
	copy(dup, chunk)
	t.chunks = append(t.chunks, dup)
	return nil
}

func (t *testChunkStore) Get(index uint64) ([]byte, error) {
	if index >= uint64(len(t.chunks)) {
		return nil, fmt.Errorf("index %d >= len %d", index, len(t.chunks))
	}
	return t.chunks[index], nil
}

func (t *testChunkStore) Len() (int64, error) {
	return int64(len(t.chunks)), nil
}

func BenchmarkCrypt(b *testing.B) {
	var key [32]byte
	_, err := hex.Decode(key[:], []byte(testKeyHex))
	if err != nil {
		b.Fatal(err)
	}

	const chunkHex = "507265616d626c650a0a57686572656173207265636f676e6974696f6e206f662074686520696e686572656e74206469676e69747920616e64206f662074686520657175616c20616e640a696e616c69656e61626c6520726967687473206f6620616c6c206d656d62657273206f66207468652068756d616e2066616d696c79206973207468650a666f756e646174696f6e206f662066726565646f6d2c206a75737469636520616e6420706561636520696e2074686520776f726c642c0a0a576865726561732064697372656761726420616e6420636f6e74656d707420666f722068756d616e20726967687473206861766520726573756c74656420696e"
	chunk, err := hex.DecodeString(chunkHex)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Crypt(key, chunk, 0)
	}
}
