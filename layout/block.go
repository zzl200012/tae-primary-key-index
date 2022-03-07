package layout

import "tae/common"

type baseBlock struct {
	Id uint64
	Data []common.Key
}

type TransientBlock struct {
	inner baseBlock
}

func NewTransientBlock(id uint64) *TransientBlock {
	return nil
}

type SortedBlock struct {
	inner baseBlock
}

type MergeSortedBlock struct {
	inner baseBlock
}


