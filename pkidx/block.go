package pkidx

import (
	bf "github.com/FastFilter/xorfilter"
	"tae/common"
)

type BlockPart struct {
	ZoneMap BlockZoneMap
	Filter *bf.BinaryFuse8
	Mode Mode
}

func NewBlockPart(data []common.Key) *BlockPart {

}
