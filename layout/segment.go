package layout

import "sync"

type Segment struct {
	mu sync.RWMutex
	blocks []*Block
}

