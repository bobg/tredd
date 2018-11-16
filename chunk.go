package tedd

import (
	"crypto/sha256"

	"github.com/chain/txvm/errors"
	"github.com/chain/txvm/protocol/txvm"
)

// ChunkSize is the size of a chunk of Tedd data.
const ChunkSize = 8192

// ChunkStore stores and retrieves data in chunks.
// The chunk size need not be ChunkSize.
type ChunkStore interface {
	// Add adds a chunk to the end of the ChunkStore.
	Add([]byte) error

	// Get gets the chunk with the given index (0-based).
	Get(uint64) ([]byte, error)

	// Len tells the number of chunks in the store.
	Len() (int64, error)
}

var errMissingChunk = errors.New("missing chunk")

func Crypt(key [32]byte, chunk []byte, index uint64) {
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

		for j := 0; pos+j < end; j++ {
			chunk[pos+j] ^= subkey[j]
		}
	}
}
