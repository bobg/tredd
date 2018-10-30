package tedd

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestServeGetDecrypt(t *testing.T) {
	const (
		keyHex            = "17f9d2125c385c2b7626034a506e524b971d9487daeb688538101c4d7d6d1f2a"
		clearRootHex      = "d78b83cba3f32b8eb56831c834f6819d383c97637c2ef95cfc71339a2da2d94f"
		wantCipherRootHex = "684d6d5652e44d45452d3c56ae5d229f701c67205a03d5c61de5a2a2134e5a0e"
	)

	var key [32]byte
	_, err := hex.Decode(key[:], []byte(keyHex))
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

func (t *testChunkStore) Len() int {
	return len(t.chunks)
}
