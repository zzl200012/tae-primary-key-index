package index

import (
	"github.com/FastFilter/xorfilter"
	"tae"
)

type PrimaryKeyFilter struct {
	inner *xorfilter.BinaryFuse8
}

func BuildPrimaryKeyFilter(keys []tae.KeyType) *PrimaryKeyFilter {
	inputs := make([]uint64, 0)
	for _, key := range keys {
		inputs = append(inputs, tae.Hash(key))
	}
	bf, err := xorfilter.PopulateBinaryFuse8(inputs)
	if err != nil {
		panic(err)
	}
	return &PrimaryKeyFilter{inner: bf}
}

func (filter *PrimaryKeyFilter) Contains(key tae.KeyType) bool {
	return filter.inner.Contains(tae.Hash(key))
}


