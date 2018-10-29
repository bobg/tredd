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

func Decrypt(w io.Writer, clearHashes, cipherChunks ChunkStream, key [32]byte) error {
	hasher := sha256.New()

	for index := uint64(0); ; index++ {
		if !clearHashes.Next() {
			break
		}
		wantClearHash, err := clearHashes.Chunk()
		if err != nil {
			return errors.Wrapf(err, "getting clear hash %d", index)
		}

		if !cipherChunks.Next() {
			return errors.Wrapf(errMissingChunk, "getting cipher chunk %d", index)
		}
		chunk, err := cipherChunks.Chunk()
		if err != nil {
			return errors.Wrapf(err, "getting cipher chunk %d", index)
		}
		crypt(key, chunk, index)

		var gotClearHash [32]byte
		merkle.LeafHash(hasher, gotClearHash[:0], chunk)
		if !bytes.Equal(gotClearHash[:], wantClearHash) {
			return BadClearHashError{Index: index}
		}

		var indexBuf [binary.MaxVarintLen64]byte
		offset := binary.PutUvarint(indexBuf[:], index)
		_, err = w.Write(chunk[offset:])
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
