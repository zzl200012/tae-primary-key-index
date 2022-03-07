package pkidx

import "tae/common"

type PrimaryKeyIndex interface {
	ProbeSingleKey(key common.Key) common.ProbeResult
	HasDuplication(batch common.QueryBatch) bool
}

type Mode uint8

const (
	Invalid Mode = iota
	InMemory
	OnDisk
)
