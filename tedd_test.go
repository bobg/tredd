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
		clearRootHex      = "fd9c327e21d6e216690fe9f1b6463db8c8619afed7ae2cebfb55300abf026110"
		wantCipherRootHex = "8857673cd5291bb005dc266705e0e35467d933e6fbe768be8f3b1438efb43380"
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
