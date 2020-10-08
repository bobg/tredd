package tredd

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/bobg/merkle"
	"github.com/chain/txvm/errors"
)

// Decrypt decrypts the chunks in cipherChunks by xoring with hashes derived from key.
// It writes the concatenated cleartext chunks to w.
// Along the way, it compares each cleartext chunk's hash to the corresponding value in clearHashes.
// If it finds a mismatch, it returns a BadClearHashError.
func Decrypt(w io.Writer, clearHashes, cipherChunks ChunkStore, key [32]byte) error {
	var (
		hasher          = sha256.New()
		chunkWithPrefix [ChunkSize + binary.MaxVarintLen64]byte
		gotClearHash    [32]byte
	)

	nhashes, err := clearHashes.Len()
	if err != nil {
		return errors.Wrap(err, "counting clear hashes")
	}
	for index := int64(0); index < nhashes; index++ {
		wantClearHash, err := clearHashes.Get(index)
		if err != nil {
			return errors.Wrapf(err, "getting clear hash %d", index)
		}

		chunk, err := cipherChunks.Get(index)
		if err != nil {
			return errors.Wrapf(err, "getting cipher chunk %d", index)
		}
		Crypt(key, chunk, index)

		m := binary.PutUvarint(chunkWithPrefix[:], uint64(index))
		copy(chunkWithPrefix[m:], chunk)

		merkle.LeafHash(hasher, gotClearHash[:0], chunkWithPrefix[m:m+len(chunk)])
		if !bytes.Equal(gotClearHash[:], wantClearHash) {
			return BadClearHashError{Index: index}
		}

		_, err = w.Write(chunk)
		if err != nil {
			return errors.Wrapf(err, "writing clear chunk %d", index)
		}
	}

	return nil
}

// BadClearHashError gives the index of a cleartext chunk whose hash doesn't have the expected value.
type BadClearHashError struct {
	// Index is the index of the chunk and of the hash within their respective ChunkStores.
	Index int64
}

func (e BadClearHashError) Error() string {
	return fmt.Sprintf("chunk %d clear hash mismatch", e.Index)
}
