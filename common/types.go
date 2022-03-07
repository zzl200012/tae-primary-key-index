package common

// Key stands for a primary key
type Key []byte

// Batch stands for a batch of primary key during an insertion
type Batch []Key

// Location stands for the location of a row
type Location struct {
	SegmentId uint64
	BlockOffset uint64
	RowOffset uint64
}

// ProbeResult stands for the result of primary key probing
type ProbeResult struct {
	Positive bool
	Answer Location
}
