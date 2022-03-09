package index

import (
	"sync"
	"tae"
)

type TransientBlockZoneMap struct {
	mu sync.RWMutex
	min tae.KeyType
	max tae.KeyType
	initialized bool
}

func (zm *TransientBlockZoneMap) Match(key tae.KeyType) bool {
	zm.mu.RLock()
	defer zm.mu.RUnlock()
	if !zm.initialized {
		return false
	}
	if tae.Compare(key, zm.min) >= 0 && tae.Compare(key, zm.max) <= 0 {
		return true
	}
	return false
}

func (zm *TransientBlockZoneMap) MatchBatch(keys []tae.KeyType) []tae.KeyType {
	zm.mu.RLock()
	defer zm.mu.RUnlock()
	if !zm.initialized {
		return []tae.KeyType{}
	}
	left := make([]tae.KeyType, 0)
	for _, key := range keys {
		if tae.Compare(key, zm.min) >= 0 && tae.Compare(key, zm.max) <= 0 {
			left = append(left, key)
		}
	}
	return left
}

func (zm *TransientBlockZoneMap) Update(key tae.KeyType) {
	zm.mu.Lock()
	defer zm.mu.Unlock()
	if !zm.initialized {
		zm.max = key
		zm.min = key
		zm.initialized = true
		return
	}
	if tae.Compare(key, zm.max) > 0 {
		zm.max = key
	}
	if tae.Compare(key, zm.min) < 0 {
		zm.min = key
	}
}

func (zm *TransientBlockZoneMap) UpdateBatch(keys []tae.KeyType) {
	zm.mu.Lock()
	defer zm.mu.Unlock()
	min, max := keys[0], keys[0]
	for _, key := range keys {
		if tae.Compare(key, max) > 0 {
			max = key
		}
		if tae.Compare(key, min) < 0 {
			min = key
		}
	}
	if !zm.initialized {
		zm.max = max
		zm.min = min
		zm.initialized = true
		return
	}
	if tae.Compare(max, zm.max) > 0 {
		zm.max = max
	}
	if tae.Compare(min, zm.min) < 0 {
		zm.min = min
	}
}

func (zm *TransientBlockZoneMap) SetMax(key tae.KeyType) {
	zm.mu.Lock()
	defer zm.mu.Unlock()
	if !zm.initialized || tae.Compare(zm.max, key) > 0 {
		return
	}
	zm.max = key
}

func (zm *TransientBlockZoneMap) SetMin(key tae.KeyType) {
	zm.mu.Lock()
	defer zm.mu.Unlock()
	if !zm.initialized || tae.Compare(key, zm.min) > 0 {
		return
	}
	zm.min = key
}

func (zm *TransientBlockZoneMap) GetMax() tae.KeyType {
	zm.mu.RLock()
	defer zm.mu.RUnlock()
	return zm.max
}

func (zm *TransientBlockZoneMap) GetMin() tae.KeyType {
	zm.mu.RLock()
	defer zm.mu.RUnlock()
	return zm.min
}

type BlockZoneMap struct {
	min tae.KeyType
	max tae.KeyType
}

func (zm *BlockZoneMap) Match(key tae.KeyType) bool {
	if tae.Compare(key, zm.min) >= 0 && tae.Compare(key, zm.max) <= 0 {
		return true
	}
	return false
}

func (zm *BlockZoneMap) MatchBatch(keys []tae.KeyType) []tae.KeyType {
	left := make([]tae.KeyType, 0)
	for _, key := range keys {
		if zm.Match(key) {
			left = append(left, key)
		}
	}
	return left
}

type SegmentZoneMap struct {
	min tae.KeyType
	max tae.KeyType
	subMaps []Range
}

func (zm *SegmentZoneMap) Match(key tae.KeyType) (bool, uint32) {
	if tae.Compare(key, zm.max) > 0 || tae.Compare(key, zm.min) < 0 {
		return false, 0
	}
	start, end := 0, len(zm.subMaps) - 1
	for start <= end {
		middle := start + (end - start) / 2
		ans := zm.subMaps[middle].match(key)
		if ans == 0 {
			return true, zm.subMaps[middle].rowOffset
		} else if ans > 0 {
			start = middle + 1
		} else {
			end = middle - 1
		}
	}
	return false, 0
}

type Range struct {
	rowOffset uint32
	min tae.KeyType
	max tae.KeyType
}

func (r *Range) match(key tae.KeyType) int {
	if tae.Compare(key, r.max) > 0 {
		return 1
	}
	if tae.Compare(key, r.min) < 0 {
		return -1
	}
	return 0
}
