package layout

import (
	"sync"
	"tae"
	"tae/index"
)

type Segment struct {
	mu sync.RWMutex
	blocks []BlockHandle
	typ SegmentType
	primaryKeyPart *index.SegmentPart
}

func (seg *Segment) ProbeSingleKey(key tae.KeyType) tae.ProbeResult {
	seg.mu.RLock()
	defer seg.mu.RUnlock()
	if seg.typ == Appendable {
		iter := seg.NewBlockIterator()
		defer iter.Close()
		for iter.Valid() {
			blk := iter.Curr()
			if res := blk.ProbeSingleKey(key); res.Presented {
				return res
			}
			iter.Next()
		}
		return tae.ProbeResult{}
	} else if seg.typ == NonAppendable {
		return seg.primaryKeyPart.ProbeSingleKey(key)
	} else {
		panic("")
	}
}

func (seg *Segment) HasDuplication(keys []tae.KeyType) bool {
	seg.mu.RLock()
	defer seg.mu.RUnlock()
	if seg.typ == Appendable {
		iter := seg.NewBlockIterator()
		defer iter.Close()
		for iter.Valid() {
			blk := iter.Curr()
			if found := blk.HasDuplication(keys); found {
				return true
			}
			iter.Next()
		}
		return false
	} else if seg.typ == NonAppendable {
		return seg.primaryKeyPart.HasDuplication(keys)
	} else {
		panic("")
	}
}

func (seg *Segment) NewBlockIterator() *BlockIterator {
	seg.mu.RLock()
	defer seg.mu.RUnlock()
	return NewBlockIterator(seg.blocks)
}

type BlockIterator struct {
	blks []BlockHandle
	pos int
}

func NewBlockIterator(blks []BlockHandle) *BlockIterator {
	return &BlockIterator{blks: blks, pos: 0}
}

func (it *BlockIterator) Next() {
	it.pos++
}

func (it *BlockIterator) Curr() BlockHandle {
	return it.blks[it.pos]
}

func (it *BlockIterator) Valid() bool {
	return it.pos < len(it.blks)
}

func (it *BlockIterator) Close() {

}



