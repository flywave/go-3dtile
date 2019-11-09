package tile3d

import "io"

type TileModel interface {
	GetHeader() Header
	GetFeatureTable() *FeatureTable
	GetBatchTable() *BatchTable
	CalcSize() int64
	Read(reader io.ReadSeeker) error
	Write(writer io.Writer) error
}
