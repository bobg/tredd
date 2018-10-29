package tedd

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"

	"github.com/bobg/merkle"
	"github.com/chain/txvm/errors"
)

var (
	errBadClearRoot = errors.New("bad clear merkle root hash")
	errBadPrefix    = errors.New("bad chunk prefix")
	errPartial      = errors.New("partial non-final chunk")
)

func Get(r io.Reader, clearRoot [32]byte, clearHashes, cipherChunks ChunkStore) ([]byte, error) {
	var (
		wasPartial bool // only the final chunk may have a partial length.
		clearMT    = merkle.NewHTree(sha256.New())
		cipherMT   = merkle.NewTree(sha256.New())
	)

	for index := uint64(0); ; index++ {
		var clearHash1, clearHash2 [32]byte

		_, err := io.ReadFull(r, clearHash1[:])
		if err == io.EOF {
			// "The error is EOF only if no bytes were read."
			break
		}
		if err != nil { // including io.ErrUnexpectedEOF
			return nil, errors.Wrapf(err, "reading clear hash %d", index)
		}

		copy(clearHash2[:], clearHash1[:])
		err = clearHashes.Store(index, clearHash1[:])
		if err != nil {
			return nil, errors.Wrapf(err, "storing clear hash %d", index)
		}
		clearMT.Add(clearHash2[:])

		var cipherChunk [chunkSize]byte
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

		gotIndex, offset := binary.Uvarint(cipherChunk[:n])
		if offset <= 0 || gotIndex != index {
			return nil, errors.Wrapf(errBadPrefix, "parsing cipher chunk %d", index)
		}

		err = cipherChunks.Store(index, cipherChunk[:n])
		if err != nil {
			return nil, errors.Wrapf(err, "storing cipher chunk %d", index)
		}
		cipherMT.Add(cipherChunk[:n])
	}

	gotClearRoot := clearMT.Root()
	if !bytes.Equal(gotClearRoot, clearRoot[:]) {
		return nil, errBadClearRoot
	}

	return cipherMT.Root(), nil
}
