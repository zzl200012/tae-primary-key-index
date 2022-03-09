package index

import (
	"github.com/FastFilter/xorfilter"
	"tae"
)

type PrimaryKeyFilter struct {
	inner *xorfilter.BinaryFuse8
}

func (filter *PrimaryKeyFilter) Contains(key tae.KeyType) bool {
	return filter.inner.Contains(tae.Hash(key))
}


