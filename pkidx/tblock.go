package pkidx

import (
	art "github.com/plar/go-adaptive-radix-tree"
	"sync"
	"tae/common"
)

type TBlockPart struct {
	sync.RWMutex
	ZoneMap BlockZoneMap
	ART         art.Tree
	Initialized bool
}

func NewTBlockPart() *TBlockPart {
	return &TBlockPart{
		ZoneMap: NewBlockZoneMap(),
		ART:     art.New(),
	}
}

func (p *TBlockPart) HasDuplication(batch common.QueryBatch) bool {
	p.RLock()
	defer p.RUnlock()
	if !p.Initialized {
		return false
	}
	var left []common.Key
	for _, key := range batch {
		if common.CompareInterface(key, p.ZoneMap.Max) || common.CompareInterface(p.ZoneMap.Min, key) {
			continue
		}
		left = append(left, key)
	}
	if len(left) == 0 {
		return false
	}
	for _, key := range left {
		if _, found := p.ART.Search(art.Key(key)); found {
			return true
		}
	}
	return false
}

func (p *TBlockPart) ProbeSingleKey(key common.Key) common.ProbeResult {
	p.RLock()
	defer p.RUnlock()
	if !p.Initialized {
		return common.ProbeResult{}
	}
	if common.CompareInterface(key, p.ZoneMap.Max) || common.CompareInterface(p.ZoneMap.Min, key) {
		return common.ProbeResult{}
	}
	if res, found := p.ART.Search(art.Key(key)); found {
		return common.ProbeResult{
			Positive: true,
			Answer:   res.(common.Location),
		}
	}
	return common.ProbeResult{}
}

func (p *TBlockPart) Update(batch common.InsertionBatch) {
	// Require: no duplication within the batch
	if len(batch) == 0 {
		return
	}
	p.Lock()
	defer p.Unlock()
	if !p.Initialized {
		pair := batch[0]
		batch = batch[1:]
		p.ZoneMap.Max = pair.Key
		p.ZoneMap.Min = pair.Key
		p.ART.Insert(art.Key(pair.Key), pair.Location)
		p.Initialized = true
	}
	for _, pair := range batch {
		key := pair.Key
		location := pair.Location
		if common.CompareInterface(key, p.ZoneMap.Max) {
			p.ZoneMap.Max = key
		}
		if common.CompareInterface(p.ZoneMap.Min, key) {
			p.ZoneMap.Min = key
		}
		p.ART.Insert(art.Key(key), location)
	}
}
