package index

import (
	"tae"
)

type SegmentPart struct {

}

func (part *SegmentPart) ProbeSingleKey(key tae.KeyType) tae.ProbeResult {
	panic("implement me")
}

func (part *SegmentPart) HasDuplication(keys []tae.KeyType) bool {
	panic("implement me")
}

