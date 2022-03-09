package tae

import "bytes"

type ProbeResult struct {
	Presented bool
	SegmentId uint32
	BlockOffset uint32
	RowOffset uint32
}

type KeyType []byte

func Hash(key KeyType) uint64 {
	return 0
}

func Compare(a, b KeyType) int {
	return bytes.Compare(a, b)
}


