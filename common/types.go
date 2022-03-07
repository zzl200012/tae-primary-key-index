package common

// Key stands for a primary key
type Key []byte

// QueryBatch stands for a batch used for deduplication
type QueryBatch []Key

// InsertionBatch stands for a batch used for insertion
type InsertionBatch []Pair

// Location stands for the location of a row
type Location struct {
	SegmentId uint64
	BlockOffset uint64
	RowOffset uint64
}

type Pair struct {
	Key Key
	Location Location
}

// ProbeResult stands for the result of primary key probing
type ProbeResult struct {
	Positive bool
	Answer Location
}
