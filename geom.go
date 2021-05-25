package tile3d

import (
	"io"
)

const (
	GEOM_MAGIC = "geom"
)

type GeomHeader struct {
	Magic                        [4]byte
	Version                      uint32
	ByteLength                   uint32
	FeatureTableJSONByteLength   uint32
	FeatureTableBinaryByteLength uint32
	BatchTableJSONByteLength     uint32
	BatchTableBinaryByteLength   uint32
}

func (h *GeomHeader) CalcSize() int64 {
	return 16
}

func (h *GeomHeader) GetByteLength() uint32 {
	return h.ByteLength
}

func (h *GeomHeader) GetFeatureTableJSONByteLength() uint32 {
	return h.FeatureTableJSONByteLength
}

func (h *GeomHeader) GetFeatureTableBinaryByteLength() uint32 {
	return h.FeatureTableBinaryByteLength
}

func (h *GeomHeader) GetBatchTableJSONByteLength() uint32 {
	return h.BatchTableJSONByteLength
}

func (h *GeomHeader) GetBatchTableBinaryByteLength() uint32 {
	return h.BatchTableBinaryByteLength
}

func (h *GeomHeader) SetFeatureTableJSONByteLength(n uint32) {
	h.FeatureTableJSONByteLength = n
}

func (h *GeomHeader) SetFeatureTableBinaryByteLength(n uint32) {
	h.FeatureTableBinaryByteLength = n
}

func (h *GeomHeader) SetBatchTableJSONByteLength(n uint32) {
	h.BatchTableJSONByteLength = n
}

func (h *GeomHeader) SetBatchTableBinaryByteLength(n uint32) {
	h.BatchTableBinaryByteLength = n
}

type GeomBoxs struct {
}

type GeomCylinders struct {
}

type GeomEllipsoids struct {
}

type GeomSpheres struct {
}

type Geom struct {
	Header       GeomHeader
	FeatureTable FeatureTable
	BatchTable   BatchTable
}

func (m *Geom) GetHeader() Header {
	return &m.Header
}

func (m *Geom) CalcSize() int64 {
	return m.Header.CalcSize()
}

func (m *Geom) GetFeatureTable() *FeatureTable {
	return &m.FeatureTable
}

func (m *Geom) GetBatchTable() *BatchTable {
	return &m.BatchTable
}

func (m *Geom) Read(reader io.ReadSeeker) error {
	return nil
}

func (m *Geom) Write(writer io.Writer) error {
	return nil
}
