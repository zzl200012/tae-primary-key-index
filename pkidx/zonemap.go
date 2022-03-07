package pkidx

import "tae/common"

type SegmentZoneMap struct {
	// TODO
}

type BlockZoneMap struct {
	Min interface{}
	Max interface{}
}

func NewBlockZoneMap() BlockZoneMap {
	return BlockZoneMap{}
}

func (bzm *BlockZoneMap) SetMax(val interface{}) {
	if !common.CompareInterface(val, bzm.Max) {
		panic("logic error")
	}
	bzm.Max = val
}

func (bzm *BlockZoneMap) SetMin(val interface{}) {
	if common.CompareInterface(val, bzm.Min) {
		panic("logic error")
	}
	bzm.Min = val
}
