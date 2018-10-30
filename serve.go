package tedd

import (
	"crypto/sha256"
	"encoding/binary"
	"io"

	"github.com/bobg/merkle"
	"github.com/chain/txvm/errors"
)

// r is the unchunked cleartext of the data being served,
// w is where to serve it to
// TODO: cleartext chunks and their hashes should be precomputed
func Serve(w io.Writer, r io.Reader, key [32]byte) ([]byte, error) {
	var (
		cipherMT = merkle.NewTree(sha256.New())
		hasher   = sha256.New()
	)

	for index := uint64(0); ; index++ {
		var chunk [chunkSize + binary.MaxVarintLen64]byte

		n1 := binary.PutUvarint(chunk[:], index)
		n2, err := io.ReadFull(r, chunk[n1:n1+chunkSize])
		if err == io.EOF {
			// "The error is EOF only if no bytes were read."
			break
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			return nil, errors.Wrapf(err, "reading clear chunk %d", index)
		}
		n := n1 + n2

		// The clearHash is prefixed with the index of the chunk.
		// If a buyer has to prove that a cipherchunk is wrong,
		// they'll have to show:
		//   - the cipherchunk has index prefix M
		//   - it is one of the cipherchunks sent by the seller
		//     (via merkle proof)
		//   - decrypting it with key K gives a clear chunk with hash H'
		//   - H' does not match H
		//   - H, with prefix M, is in the clear merkle tree
		var clearHash [32 + binary.MaxVarintLen64]byte
		n3 := binary.PutUvarint(clearHash[:], index)
		merkle.LeafHash(hasher, clearHash[:n3], chunk[:n])

		_, err = w.Write(clearHash[:n3+32])
		if err != nil {
			return nil, errors.Wrapf(err, "writing clear hash %d", index)
		}

		crypt(key, chunk[n1:n], index) // n.b. overwrites the contents of chunk
		_, err = w.Write(chunk[:n])
		if err != nil {
			return nil, errors.Wrapf(err, "writing cipher chunk %d", index)
		}
		cipherMT.Add(chunk[:n])
	}

	return cipherMT.Root(), nil
}
