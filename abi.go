package tredd

import "encoding/binary"

func Prefix(n uint64, chunk []byte) []byte {
	result := make([]byte, len(chunk)+8)
	binary.BigEndian.PutUint64(result, n)
	copy(result[8:], chunk)
	return result
}

func SubchunkKeyParams(key [32]byte, index, n uint64) []byte {
	var result [48]byte
	copy(result[:], key[:])
	binary.BigEndian.PutUint64(result[32:], index)
	binary.BigEndian.PutUint64(result[40:], n)
	return result[:]
}
