package tredd

import (
	"crypto/sha256"

	"github.com/bobg/merkle/v2"
	"github.com/pkg/errors"
)

// ChunkSize is the size of a chunk of Tredd data.
const ChunkSize = 8192

// ChunkStore stores and retrieves data in chunks.
// The chunk size need not be ChunkSize.
type ChunkStore interface {
	// Add adds a chunk to the end of the ChunkStore.
	Add([]byte) error

	// Get gets the chunk with the given index (0-based).
	Get(uint64) ([]byte, error)

	// Len tells the number of chunks in the store.
	Len() (uint64, error)
}

var errMissingChunk = errors.New("missing chunk")

func Crypt(key [32]byte, chunk []byte, index uint64) error {
	var (
		hasher = sha256.New()
		subkey [32]byte
	)

	for i := 0; 32*i < len(chunk); i++ {
		// compute subchunk key
		hasher.Reset()

		inp := SubchunkKeyParams(key, index, uint64(i))

		hasher.Write(inp)
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
	return nil
}

func PrepareForRefund(index uint64, clearHashes, cipherChunks ChunkStore) (clearHashN [32]byte, cipherChunkN []byte, clearProof, cipherProof merkle.Proof, err error) {
	var n uint64
	n, err = clearHashes.Len()
	if err != nil {
		return
	}

	var clearHashNBytes []byte
	clearHashNBytes, err = clearHashes.Get(index)
	if err != nil {
		return
	}
	prefixedClearHashN := Prefix(index, clearHashNBytes)

	cipherChunkN, err = cipherChunks.Get(index)
	if err != nil {
		return
	}
	prefixedCipherChunkN := Prefix(index, cipherChunkN)

	var (
		clearMT  = merkle.NewProofTree(sha256.New(), prefixedClearHashN)
		cipherMT = merkle.NewProofTree(sha256.New(), prefixedCipherChunkN)
	)

	for i := uint64(0); index < n; index++ {
		var clearHash []byte
		clearHash, err = clearHashes.Get(i)
		if err != nil {
			return
		}
		clearMT.Add(clearHash)

		var cipherChunk []byte
		cipherChunk, err = cipherChunks.Get(i)
		if err != nil {
			return
		}
		prefixedCipherChunk := Prefix(i, cipherChunk)
		cipherMT.Add(prefixedCipherChunk)
	}

	copy(clearHashN[:], clearHashNBytes)
	return clearHashN, cipherChunkN, clearMT.Proof(), cipherMT.Proof(), nil
}
