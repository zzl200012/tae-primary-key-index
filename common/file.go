package common

import "io"

type IVFile interface {
	io.Reader
	Ref()
	Unref()
	RefCount() int64
	Stat() FileInfo
}

type FileInfo interface {
	Name() string
	Size() int64
	OriginSize() int64
	CompressAlgo() int
	Rows() uint64
}

type baseMemFile struct {
	stat baseFileInfo
}

func NewMemFile(size int64, rows uint64) IVFile {
	return &baseMemFile{
		stat: baseFileInfo{
			size: size,
			rows: rows,
		},
	}
}

func (f *baseMemFile) Ref()                             {}
func (f *baseMemFile) Unref()                           {}
func (f *baseMemFile) RefCount() int64                  { return 0 }
func (f *baseMemFile) Read(p []byte) (n int, err error) { return n, err }
func (f *baseMemFile) Stat() FileInfo                   { return &f.stat }

type baseFileInfo struct {
	size int64
	rows uint64
}

func (i *baseFileInfo) Rows() uint64      { return i.rows }
func (i *baseFileInfo) Name() string      { return "" }
func (i *baseFileInfo) Size() int64       { return i.size }
func (i *baseFileInfo) OriginSize() int64 { return i.size }
func (i *baseFileInfo) CompressAlgo() int { return 0 }
