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

type reader interface {
	io.Reader
	io.ByteReader
}

func Get(r reader, clearRoot [32]byte, clearHashes, cipherChunks ChunkStore) ([]byte, error) {
	var (
		wasPartial bool // only the final chunk may have a partial length.
		clearMT    = merkle.NewTree(sha256.New())
		cipherMT   = merkle.NewTree(sha256.New())
	)

	for index := uint64(0); ; index++ {
		clearHashPrefix, err := binary.ReadUvarint(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrapf(err, "reading clear hash prefix %d", index)
		}
		if clearHashPrefix != index {
			return nil, errors.Wrapf(errBadPrefix, "reading clear hash prefix %d", index)
		}

		var (
			clearHash           [32]byte
			clearHashWithPrefix [32 + binary.MaxVarintLen64]byte
		)

		_, err = io.ReadFull(r, clearHash[:])
		if err != nil { // including io.ErrUnexpectedEOF
			return nil, errors.Wrapf(err, "reading clear hash %d", index)
		}

		n := binary.PutUvarint(clearHashWithPrefix[:], index)
		copy(clearHashWithPrefix[n:], clearHash[:])

		err = clearHashes.Add(clearHash[:])
		if err != nil {
			return nil, errors.Wrapf(err, "storing clear hash %d", index)
		}
		clearMT.Add(clearHashWithPrefix[:n+32])

		var cipherChunk [chunkSize + binary.MaxVarintLen64]byte
		gotIndex, err := binary.ReadUvarint(r)
		if err != nil {
			return nil, errors.Wrapf(err, "reading cipher chunk prefix %d", index)
		}
		if gotIndex != index {
			return nil, errors.Wrapf(errBadPrefix, "reading cipher chunk %d", index)
		}
		n1 := binary.PutUvarint(cipherChunk[:], index)
		n2, err := io.ReadFull(r, cipherChunk[n1:n1+chunkSize])
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
		n = n1 + n2

		err = cipherChunks.Add(cipherChunk[:n])
		if err != nil {
			return nil, errors.Wrapf(err, "storing cipher chunk %d", index)
		}
		cipherMT.Add(cipherChunk[:n])
	}

	gotClearRoot := clearMT.Root()
	if !bytes.Equal(gotClearRoot, clearRoot[:]) {
		return nil, errors.Wrapf(errBadClearRoot, "got %x, want %x", gotClearRoot, clearRoot[:])
	}

	return cipherMT.Root(), nil
}
