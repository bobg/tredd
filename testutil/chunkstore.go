package testutil

import "fmt"

type ChunkStore struct {
	chunks [][]byte
}

func (s *ChunkStore) Add(chunk []byte) error {
	dup := make([]byte, len(chunk))
	copy(dup, chunk)
	s.chunks = append(s.chunks, dup)
	return nil
}

func (s *ChunkStore) Get(index uint64) ([]byte, error) {
	if index >= uint64(len(s.chunks)) {
		return nil, fmt.Errorf("index %d >= len %d", index, len(s.chunks))
	}
	return s.chunks[index], nil
}

func (s *ChunkStore) Len() (uint64, error) {
	return uint64(len(s.chunks)), nil
}
