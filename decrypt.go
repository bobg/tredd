package tedd

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/bobg/merkle"
	"github.com/chain/txvm/errors"
)

func Decrypt(w io.Writer, clearHashes, cipherChunks ChunkStore, key [32]byte) error {
	var (
		hasher          = sha256.New()
		chunkWithPrefix [chunkSize + binary.MaxVarintLen64]byte
		gotClearHash    [32]byte
	)

	for index := uint64(0); index < uint64(clearHashes.Len()); index++ {
		wantClearHash, err := clearHashes.Get(index)
		if err != nil {
			return errors.Wrapf(err, "getting clear hash %d", index)
		}

		chunk, err := cipherChunks.Get(index)
		if err != nil {
			return errors.Wrapf(err, "getting cipher chunk %d", index)
		}
		crypt(key, chunk, index)

		m := binary.PutUvarint(chunkWithPrefix[:], index)
		copy(chunkWithPrefix[m:], chunk)

		merkle.LeafHash(hasher, gotClearHash[:0], chunkWithPrefix[:m+len(chunk)])
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

type BadClearHashError struct {
	Index uint64
}

func (e BadClearHashError) Error() string {
	return fmt.Sprintf("chunk %d clear hash mismatch", e.Index)
}
