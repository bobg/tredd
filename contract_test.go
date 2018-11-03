package tedd

import (
	"crypto/sha256"
	"io"
	"math"
	"os"
	"testing"

	"github.com/bobg/merkle"
	"github.com/chain/txvm/protocol/txvm/op"
	"github.com/chain/txvm/protocol/txvm/txvmutil"

	"i10r.io/protocol/txvm"
)

func TestMerkleCheck(t *testing.T) {
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
