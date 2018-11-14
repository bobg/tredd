package main

import (
	"os"
)

type fileChunkStore struct {
	filename  string
	chunksize int64
}

func (s *fileChunkStore) Add(bits []byte) error {
	f, err := os.OpenFile(s.filename, os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(bits)
	return err
}

func (s *fileChunkStore) Get(index uint64) ([]byte, error) {
	f, err := os.Open(s.filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = f.Seek(int64(index)*s.chunksize, os.SEEK_SET) // xxx range check
	if err != nil {
		return nil, err
	}

	result := make([]byte, s.chunksize)
	n, err := f.Read(result) // xxx use ReadFull, allow partial chunks only at eof
	return result[:n], err
}

func (s *fileChunkStore) Len() (int, error) {
	info, err := os.Stat(s.filename)
	if err != nil {
		return 0, err
	}
	size := info.Size()
	n, r := int(size/s.chunksize), size%s.chunksize // xxx range check
	if r > 0 {
		n++
	}
	return n, nil
}
