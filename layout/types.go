package layout

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

type Segment interface {
	Type() SegmentType
}

type Block interface {
	Type() BlockType
}




