package index

import "tae/common"

type PrimaryKeyIndex interface {
	ProbeSingleKey(key common.Key) common.ProbeResult
	DeduplicateBatch(batch common.Batch) bool
}
