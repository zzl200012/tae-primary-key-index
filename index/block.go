package index

import (
	"sync"
	"tae"
	"tae/layout"
)

type BlockPart struct {
	mu sync.RWMutex
	host *layout.Block
	filter *PrimaryKeyFilter
	zoneMap *BlockZoneMap
}

func (part *BlockPart) ProbeSingleKey(key tae.KeyType) tae.ProbeResult {
	if !part.zoneMap.Match(key) {
		return tae.ProbeResult{}
	}
	if !part.filter.Contains(key) {
		return tae.ProbeResult{}
	}
	return part.host.ProbeSingleKey(key)
}

func (part *BlockPart) HasDuplication(keys []tae.KeyType) bool {
	left := part.zoneMap.MatchBatch(keys)
	if len(left) == 0 {
		return false
	}

	for _, key := range left {
		if !part.filter.Contains(key) {
			continue
		}
		if res := part.host.ProbeSingleKey(key); res.Presented {
			return true
		}
	}
	return false
}
