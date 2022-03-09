package index

import (
	"sync"
	"tae"
	"tae/layout"
)

type TransientBlockPart struct {
	mu sync.RWMutex
	host *layout.TransientBlock
	zoneMap TransientBlockZoneMap
	art TransientBlockART
}

func (part *TransientBlockPart) ProbeSingleKey(key tae.KeyType) tae.ProbeResult {
	if !part.zoneMap.Match(key) {
		return tae.ProbeResult{}
	}
	return part.art.ProbeSingleKey(key)
}

func (part *TransientBlockPart) HasDuplication(keys []tae.KeyType) bool {
	left := part.zoneMap.MatchBatch(keys)
	if len(left) == 0 {
		return false
	}
	return part.art.CheckKeysExist(left)
}

func (part *TransientBlockPart) Update(key tae.KeyType, offset uint32) {
	part.zoneMap.Update(key)

	part.art.Insert(key, offset)
}

func (part *TransientBlockPart) UpdateBatch(keys []tae.KeyType, offsets []uint32) {
	part.zoneMap.UpdateBatch(keys)

	part.art.BatchInsert(keys, offsets)
}

//func (part *TransientBlockPart) Freeze() *PrimaryKeyFilter {
//
//}