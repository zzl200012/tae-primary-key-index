package layout

import "tae/index"

type SegmentType uint8
type BlockType uint8

const (
	InvalidSeg SegmentType = iota
	Appendable
	NonAppendable
)

const (
	InvalidBlk BlockType = iota
	Transient
	Sorted
	MergeSorted
)

type BlockHandle interface {
	index.PrimaryKeyResolver
}
