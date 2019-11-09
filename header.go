package tile3d

type Header interface {
	CalcSize() int64

	GetByteLength() uint32

	GetFeatureTableJSONByteLength() uint32
	GetFeatureTableBinaryByteLength() uint32

	GetBatchTableJSONByteLength() uint32
	GetBatchTableBinaryByteLength() uint32

	SetFeatureTableJSONByteLength(uint32)
	SetFeatureTableBinaryByteLength(uint32)

	SetBatchTableJSONByteLength(uint32)
	SetBatchTableBinaryByteLength(uint32)
}
