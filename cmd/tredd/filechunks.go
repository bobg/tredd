package main

import (
	"io"
	"os"
)

type fileChunkStore struct {
	filename  string
	chunksize int64
	size      int64 // size of file in bytes
}

func newFileChunkStore(filename string, chunksize int64) (*fileChunkStore, error) {
	result := &fileChunkStore{
		filename:  filename,
		chunksize: chunksize,
	}
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// ok
	} else if err != nil {
		return nil, err
	} else {
		result.size = info.Size()
	}
	return result, nil
}

func (s *fileChunkStore) Add(bits []byte) error {
	f, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(bits)
	if err != nil {
		return err
	}
	s.size += int64(len(bits))
	return nil
}

func (s *fileChunkStore) Get(index int64) ([]byte, error) {
	f, err := os.Open(s.filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = f.Seek(index*s.chunksize, os.SEEK_SET) // TODO: range check
	if err != nil {
		return nil, err
	}

	result := make([]byte, s.chunksize)

	n, err := io.ReadFull(f, result)
	if err == io.ErrUnexpectedEOF {
		// Partial chunk allowed only at EOF.
		if index*s.chunksize+int64(n) == s.size {
			return result[:n], nil
		}
	}
	return result[:n], err
}

func (s *fileChunkStore) Len() (int64, error) {
	n, r := s.size/s.chunksize, s.size%s.chunksize
	if r > 0 {
		n++
	}
	return n, nil
}
