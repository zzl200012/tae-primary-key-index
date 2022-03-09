package layout

type TransientBlock struct {
	host *Segment
}

func (tblk *TransientBlock) GetId() uint32 {
	return 0
}

func (tblk *TransientBlock) GetSegmentId() uint32 {
	return 0
}

func (tblk *TransientBlock) GetBlockOffset() uint32 {
	return 0
}
