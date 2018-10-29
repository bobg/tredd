package tedd

import (
	"crypto/sha256"
	"io"

	"github.com/bobg/merkle"
	"github.com/chain/txvm/errors"
)

func Serve(w io.Writer, clearHashes, clearChunks ChunkStream, key [32]byte) ([]byte, error) {
	cipherMT := merkle.NewTree(sha256.New())

	for index := uint64(0); ; index++ {
		if !clearHashes.Next() {
			break
		}
		clearHash, err := clearHashes.Chunk()
		if err != nil {
			return nil, errors.Wrapf(err, "getting clear hash %d", index)
		}
		_, err = w.Write(clearHash)
		if err != nil {
			return nil, errors.Wrapf(err, "writing clear hash %d", index)
		}

		if !clearChunks.Next() {
			return nil, errors.Wrapf(errMissingChunk, "getting clear chunk %d", index)
		}
		chunk, err := clearChunks.Chunk()
		if err != nil {
			return nil, errors.Wrapf(err, "getting clear chunk %d", index)
		}
		crypt(key, chunk, index) // n.b. overwrites the contents of chunk
		_, err = w.Write(chunk)
		if err != nil {
			return nil, errors.Wrapf(err, "writing cipher chunk %d", index)
		}
		cipherMT.Add(chunk)
	}

	return cipherMT.Root(), nil
}
