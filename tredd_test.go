package tredd

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bobg/tredd/testutil"
)

const testKeyHex = "17f9d2125c385c2b7626034a506e524b971d9487daeb688538101c4d7d6d1f2a"

func TestServeGetDecrypt(t *testing.T) {
	const (
		clearRootHex      = "abd68b1ada7fbc926f4a7b5dc28f0187ff2d34a4c73a20632743891d5511a204"
		wantCipherRootHex = "36f2a7918a9f710dbbaed6f53444ed780cf3fc070165734d9a46f13992ed74a1"
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
		clearHashes  = new(testutil.ChunkStore)
		cipherChunks = new(testutil.ChunkStore)
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
