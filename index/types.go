package index

import "tae"

type PrimaryKeyResolver interface {
	ProbeSingleKey(key tae.KeyType) tae.ProbeResult
	HasDuplication(keys []tae.KeyType) bool
}
