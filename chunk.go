package tedd

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/chain/txvm/errors"
	"github.com/chain/txvm/protocol/txvm"
)

const chunkSize = 8192

type ChunkStore interface {
	Add([]byte) error
	Get(uint64) ([]byte, error)
	Len() int
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
		end := pos + 32
		if end > len(chunk) {
			end = len(chunk)
		}

		for j := offset; pos+j < end; j++ {
			chunk[pos+j] ^= subkey[j]
		}
		offset = 0
	}
}
