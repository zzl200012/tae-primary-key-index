package index

import (
	art "github.com/plar/go-adaptive-radix-tree"
	"sync"
	"tae"
	"tae/layout"
)

type TransientBlockART struct {
	mu sync.RWMutex
	inner art.Tree
	host *layout.TransientBlock
}

func (tree *TransientBlockART) Insert(key tae.KeyType, row uint32) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.inner.Insert(art.Key(key), row)
}

func (tree *TransientBlockART) BatchInsert(keys []tae.KeyType, rows []uint32) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for i := 0; i < len(keys); i++ {
		tree.inner.Insert(art.Key(keys[i]), rows[i])
	}
}

func (tree *TransientBlockART) CheckKeysExist(keys []tae.KeyType) bool {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	for _, key := range keys {
		if _, found := tree.inner.Search(art.Key(key)); found {
			return true
		}
	}
	return false
}

func (tree *TransientBlockART) ProbeSingleKey(key tae.KeyType) tae.ProbeResult {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	rowOffset, found := tree.inner.Search(art.Key(key))
	if !found {
		return tae.ProbeResult{}
	}
	return tae.ProbeResult{
		Presented:   true,
		SegmentId:   tree.host.GetSegmentId(),
		BlockOffset: tree.host.GetBlockOffset(),
		RowOffset:   rowOffset.(uint32),
	}
}


