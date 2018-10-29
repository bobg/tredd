package tedd

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/chain/txvm/errors"
	"github.com/chain/txvm/protocol/txvm"
)

const chunkSize = 8192

type ChunkStore interface {
	Store(uint64, []byte) error
}

type ChunkStream interface {
	Next() bool
	Chunk() ([]byte, error)
}

var errMissingChunk = errors.New("missing chunk")

func crypt(key [32]byte, chunk []byte, index uint64) {
	var indexBuf [binary.MaxVarintLen64]byte
	offset := binary.PutUvarint(indexBuf[:], index)

	var (
		hasher = sha256.New()
		subkey [32]byte
	)

	for i := 0; 32*i < len(chunk); i++ {
		// compute subchunk key
		hasher.Reset()
		hasher.Write(key[:])
		hasher.Write(txvm.Encode(txvm.Int(index)))
		hasher.Write(txvm.Encode(txvm.Int(i)))
		hasher.Sum(subkey[:0])

		pos := 32 * i
		for j := offset; j < 32; j++ {
			chunk[pos+j] ^= subkey[j]
		}
		offset = 0
	}
}
