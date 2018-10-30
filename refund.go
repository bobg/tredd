package tedd

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/bobg/merkle"
	"github.com/chain/txvm/errors"
)

func Refund(index uint64, clearHashes, cipherChunks ChunkStore, key [32]byte) error {
	clearHash, err := clearHashes.Get(index)
	if err != nil {
		return errors.Wrapf(err, "getting clear chunk hash %d", index)
	}

	var clearHashWithPrefix [32 + binary.MaxVarintLen64]byte
	n := binary.PutUvarint(clearHashWithPrefix[:], index)
	copy(clearHashWithPrefix[n:], clearHash)

	cipherChunk, err := cipherChunks.Get(index)
	if err != nil {
		return errors.Wrapf(err, "getting cipher chunk %d", index)
	}

	var (
		clearMT  = merkle.NewProofTree(sha256.New(), clearHashWithPrefix[:n+32])
		cipherMT = merkle.NewProofTree(sha256.New(), cipherChunk)
	)

	for i := uint64(0); i < uint64(clearHashes.Len()); i++ {
		clearHash, err := clearHashes.Get(index)
		if err != nil {
			return errors.Wrapf(err, "getting clear chunk hash %d", index)
		}
		var clearHashWithPrefix [32 + binary.MaxVarintLen64]byte
		n := binary.PutUvarint(clearHashWithPrefix[:], i)
		copy(clearHashWithPrefix[n:], clearHash)
		clearMT.Add(clearHashWithPrefix[:n+32])

		cipherMT.Add(cipherChunk)
	}

	// xxx use index, clearMT.Proof(), cipherMT.Proof(), clearHashWithPrefix, and cipherChunk
	// in an invocation of the contract refund clause
	return nil
}
