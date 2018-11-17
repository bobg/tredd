package tredd

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"testing"

	"github.com/bobg/merkle"
	"github.com/chain/txvm/errors"
	"github.com/chain/txvm/protocol/txvm"
	"github.com/chain/txvm/protocol/txvm/asm"
	"github.com/chain/txvm/protocol/txvm/op"
	"github.com/chain/txvm/protocol/txvm/txvmutil"
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

		prog := testMerkleCheckProg(proof, root, refchunk)

		_, err := txvm.Validate(prog, 3, math.MaxInt64)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testMerkleCheckProg(proof merkle.Proof, wantRoot, refchunk []byte) []byte {
	b := new(txvmutil.Builder)
	b.PushdataBytes(wantRoot)
	b.Tuple(func(b *txvmutil.TupleBuilder) {
		for i := len(proof) - 1; i >= 0; i-- {
			b.PushdataBytes(proof[i].H)
			var isLeft int64
			if proof[i].Left {
				isLeft = 1
			}
			b.PushdataInt64(isLeft)
		}
	})
	b.PushdataBytes(refchunk)
	b.PushdataBytes(merkleCheckProg).Op(op.Exec)
	b.PushdataBytes([]byte{}).PushdataInt64(0).Op(op.Nonce)
	b.Op(op.Finalize)

	return b.Build()
}

func TestTxVMDecrypt(t *testing.T) {
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

	src := fmt.Sprintf("x'%x' 0 x'%x'\n%s", key[:], cipher[:], decryptSrc)
	prog, err := asm.Assemble(src)
	if err != nil {
		t.Fatal(err)
	}
	vm, err := txvm.Validate(prog, 3, math.MaxInt64)
	if errors.Root(err) != txvm.ErrResidue {
		t.Fatalf("expected ErrResidue, got %v", err)
	}
	tuple := vm.StackItem(vm.StackLen() - 1).(txvm.Tuple)
	if typecode := string(tuple[0].(txvm.Bytes)); typecode != "S" {
		t.Fatalf("top of VM stack is item with type code %s, want S (for string)", typecode)
	}
	b := tuple[1].(txvm.Bytes)
	if !bytes.Equal(b, clear[:]) {
		t.Errorf("got %x, want %x", []byte(b), clear[:])
	}
}
