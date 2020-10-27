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

// Receive reads the output of Serve and alternately call hashFn and chunkFn on each hash/cipherchunk pair.
func Receive(r io.Reader, hashFn func([32]byte, uint64) error, chunkFn func([]byte, uint64) error) error {
	var wasPartial bool

	for i := uint64(0); ; i++ {
		var h [32]byte
		_, err := io.ReadFull(r, h[:])
		if err == io.EOF {
			return nil
		}
		if err != nil { // including io.ErrUnexpectedEOF
			return errors.Wrapf(err, "reading hash %d", i)
		}
		err = hashFn(h, i)
		if err != nil {
			return errors.Wrapf(err, "processing hash %d", i)
		}

		var chunk [ChunkSize]byte
		n, err := io.ReadFull(r, chunk[:])
		if err == io.EOF {
			return errors.Wrapf(errMissingChunk, "reading chunk %d", i)
		}
		if err == io.ErrUnexpectedEOF {
			if wasPartial {
				return errors.Wrapf(errPartial, "reading chunk %d", i)
			}
			wasPartial = true
		} else if err != nil {
			return errors.Wrapf(err, "reading chunk %d", i)
		}
		err = chunkFn(chunk[:n], i)
		if err != nil {
			return errors.Wrapf(err, "processing chunk %d", i)
		}
	}
}

// Get parses a stream of interleaved <clearhash><cipherchunk> pairs,
// placing them in their respective ChunkStores.
// It uses Receive to handle the output of Serve.
// Along the way it compares the clear hashes' Merkle root hash to the expected value in clearRoot.
// If it finds a mismatch it returns errBadClearRoot.
// If there is no error, the Merkle root hash of the cipher chunks is returned.
// Both Merkle root hashes are computed from values prepended with the chunk index number.
func Get(r io.Reader, clearRoot [32]byte, clearHashes, cipherChunks ChunkStore) ([]byte, error) {
	var (
		clearMT  = merkle.NewHTree(sha256.New())
		cipherMT = merkle.NewTree(sha256.New())
	)

	err := Receive(
		r,
		func(clearHash [32]byte, index uint64) error {
			err := clearHashes.Add(clearHash[:])
			if err != nil {
				return errors.Wrap(err, "adding hash to ChunkStore")
			}
			clearMT.Add(clearHash[:])
			return nil
		},
		func(cipherChunk []byte, index uint64) error {
			err := cipherChunks.Add(cipherChunk)
			if err != nil {
				return errors.Wrap(err, "adding chunk to ChunkStore")
			}
			prefixedCipherChunk := PrefixChunk(index, cipherChunk)
			cipherMT.Add(prefixedCipherChunk)
			return nil
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "receiving data")
	}

	if gotClearRoot := clearMT.Root(); !bytes.Equal(gotClearRoot, clearRoot[:]) {
		return nil, errors.Wrapf(errBadClearRoot, "got %x, want %x", gotClearRoot, clearRoot[:])
	}

	return cipherMT.Root(), err
}

// Serve produces a stream of interleaved <clearhash><cipherchunk> pairs from the content in r.
// It writes the stream to w, encrypting the chunks by xoring with hashes derived from key.
// The return value is the Merkle root hash of the cipher chunks, each prepended with its chunk index.
// TODO: Cleartext chunks and their hashes can be precomputed and supplied as ChunkStores.
func Serve(w io.Writer, r io.Reader, key [32]byte) ([]byte, error) {
	var (
		cipherMT = merkle.NewTree(sha256.New())
		hasher   = sha256.New()
	)

	for index := uint64(0); ; index++ {
		var chunk [ChunkSize]byte
		n, err := io.ReadFull(r, chunk[:])
		if err == io.EOF {
			break
		} else if err != nil && err != io.ErrUnexpectedEOF {
			return nil, errors.Wrapf(err, "reading clear chunk %d", index)
		}

		prefixedClearChunk := PrefixChunk(index, chunk[:n])
		var clearChunkHash [32]byte
		merkle.LeafHash(hasher, clearChunkHash[:0], prefixedClearChunk)
		_, err = w.Write(clearChunkHash[:])
		if err != nil {
			return nil, errors.Wrapf(err, "writing clear hash %d", index)
		}

		Crypt(key, chunk[:n], index)
		_, err = w.Write(chunk[:n])
		if err != nil {
			return nil, errors.Wrapf(err, "writing cipher chunk %d", index)
		}

		prefixedCipherChunk := PrefixChunk(index, chunk[:n])
		cipherMT.Add(prefixedCipherChunk)
	}

	return cipherMT.Root(), nil
}
