package tredd

import (
	"crypto/sha256"
	"io"

	"github.com/bobg/merkle"
	"github.com/pkg/errors"
)

// Serve produces a stream of interleaved <clearhash><cipherchunk> pairs from the content in r.
// It writes the stream to w, encrypting the chunks by xoring with hashes derived from key.
// The return value is the Merkle root hash of the cipher chunks, each prepended with its chunk index.
// TODO: Cleartext chunks and their hashes can be precomputed and supplied as ChunkStores.
func Serve(w io.Writer, r io.Reader, key [32]byte) ([]byte, error) {
	var (
		cipherMT = merkle.NewTree(sha256.New())
		hasher   = sha256.New()
	)

	for index := int64(0); ; index++ {
		var chunk [ChunkSize]byte
		n, err := io.ReadFull(r, chunk[:])
		if err == io.EOF {
			break
		} else if err != nil && err != io.ErrUnexpectedEOF {
			return nil, errors.Wrapf(err, "reading clear chunk %d", index)
		}

		var clearChunkHash [32]byte
		merkle.LeafHash(hasher, clearChunkHash[:0], chunk[:n])
		_, err = w.Write(clearChunkHash[:])
		if err != nil {
			return nil, errors.Wrapf(err, "writing clear hash %d", index)
		}

		Crypt(key, chunk[:n], index)
		_, err = w.Write(chunk[:n])
		if err != nil {
			return nil, errors.Wrapf(err, "writing cipher chunk %d", index)
		}

		prefixedCipherChunk, err := PrefixChunk(uint64(index), chunk[:n])
		if err != nil {
			return nil, errors.Wrapf(err, "packing prefixed cipher chunk %d", index)
		}
		cipherMT.Add(prefixedCipherChunk)
	}

	return cipherMT.Root(), nil
}
