package tredd

import (
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
	Get(int64) ([]byte, error)

	// Len tells the number of chunks in the store.
	Len() (int64, error)
}

var errMissingChunk = errors.New("missing chunk")

var cryptArgTypes = abi.Arguments{
	{Type: mustABIType("byte32")},
	{Type: mustABIType("uint64")},
	{Type: mustABIType("uint64")},
}

func Crypt(key [32]byte, chunk []byte, index int64) error {
	var (
		hasher = sha256.New()
		subkey [32]byte
	)

	for i := 0; 32*i < len(chunk); i++ {
		// compute subchunk key
		hasher.Reset()

		inp, err := cryptArgTypes.Pack(key, uint64(index), uint64(i))
		if err != nil {
			return err
		}

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

func mustABIType(name string) abi.Type {
	typ, err := abi.NewType(name, "", nil)
	if err != nil {
		panic(err)
	}
	return typ
}
