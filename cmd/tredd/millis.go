package main

import "time"

// Duplicated from github.com/chain/txvm/protocol/bc

// Millis converts a time.Time to a number of milliseconds since 1970.
func Millis(t time.Time) uint64 {
	return uint64(t.UnixNano()) / uint64(time.Millisecond)
}

// FromMillis converts a number of milliseconds since 1970 to a time.Time.
func FromMillis(ms uint64) time.Time {
	return time.Unix(0, int64(ms*uint64(time.Millisecond))).UTC()
}
