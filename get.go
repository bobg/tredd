package tredd

import (
	"bytes"
	"crypto/sha256"
	"io"

	"github.com/bobg/merkle"
	"github.com/pkg/errors"
)

var (
	errBadClearRoot = errors.New("bad clear merkle root hash")
	errBadPrefix    = errors.New("bad chunk prefix")
	errPartial      = errors.New("partial non-final chunk")
)

// Get parses a stream of interleaved <clearhash><cipherchunk> pairs,
// placing them in their respective ChunkStores.
// Along the way it compares the clear hashes' Merkle root hash to the expected value in clearRoot.
// If it finds a mismatch it returns errBadClearRoot.
// If there is no error, the Merkle root hash of the cipher chunks is returned.
// Both Merkle root hashes are computed from values prepended with the chunk index number.
func Get(r io.Reader, clearRoot [32]byte, clearHashes, cipherChunks ChunkStore) ([]byte, error) {
	var (
		wasPartial bool // only the final chunk may have a partial length.
		clearMT    = merkle.NewTree(sha256.New())
		cipherMT   = merkle.NewTree(sha256.New())
	)

	for index := uint64(0); ; index++ {
		var clearHash [32]byte
		_, err := io.ReadFull(r, clearHash[:])
		if err == io.EOF {
			// "The error is EOF only if no bytes were read."
			break
		}
		if err != nil { // including io.ErrUnexpectedEOF
			return nil, errors.Wrapf(err, "reading clear hash %d", index)
		}

		err = clearHashes.Add(clearHash[:])
		if err != nil {
			return nil, errors.Wrapf(err, "storing clear hash %d", index)
		}

		prefixedClearHash := PrefixHash(index, clearHash)

		clearMT.Add(prefixedClearHash)

		var cipherChunk [ChunkSize]byte
		n, err := io.ReadFull(r, cipherChunk[:])
		if err == io.EOF {
			// "The error is EOF only if no bytes were read."
			return nil, errors.Wrapf(errMissingChunk, "reading chunk %d", index)
		}
		if err == io.ErrUnexpectedEOF {
			if wasPartial {
				return nil, errPartial
			}
			wasPartial = true
		} else if err != nil {
			return nil, errors.Wrapf(err, "reading cipher chunk %d", index)
		}

		err = cipherChunks.Add(cipherChunk[:n])
		if err != nil {
			return nil, errors.Wrapf(err, "storing cipher chunk %d", index)
		}

		prefixedCipherChunk := PrefixChunk(index, cipherChunk[:n])

		cipherMT.Add(prefixedCipherChunk)
	}

	gotClearRoot := clearMT.Root()
	if !bytes.Equal(gotClearRoot, clearRoot[:]) {
		return nil, errors.Wrapf(errBadClearRoot, "got %x, want %x", gotClearRoot, clearRoot[:])
	}

	return cipherMT.Root(), nil
}
