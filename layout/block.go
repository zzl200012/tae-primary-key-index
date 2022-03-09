package layout

import "tae"

type Block struct {
	host *Segment
}

func (blk *Block) GetId() uint32 {
	return 0
}

func (blk *Block) GetSegmentId() uint32 {
	return 0
}

func (blk *Block) GetBlockOffset() uint32 {
	return 0
}

func (blk *Block) ProbeSingleKey(key tae.KeyType) tae.ProbeResult {
	panic("")
}

func (blk *Block) HasDuplication(keys []tae.KeyType) bool {
	return false
}
