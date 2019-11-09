package tile3d

import (
	"encoding/binary"
	"io"
)

const (
	PNTS_MAGIC = "pnts"
)

const (
	PNTS_PROP_BATCH_LENGTH            = "BATCH_LENGTH"
	PNTS_PROP_RTC_CENTER              = "RTC_CENTER"
	PNTS_PROP_POSITION                = "POSITION"
	PNTS_PROP_POSITION_QUANTIZED      = "POSITION_QUANTIZED"
	PNTS_PROP_QUANTIZED_VOLUME_OFFSET = "QUANTIZED_VOLUME_OFFSET"
	PNTS_PROP_QUANTIZED_VOLUME_SCALE  = "QUANTIZED_VOLUME_SCALE"
	PNTS_PROP_NORMAL                  = "NORMAL"
	PNTS_PROP_NORMAL_OCT32P           = "NORMAL_OCT32P"
	PNTS_PROP_POSITION_LENGTH         = "POSITION_LENGTH"
	PNTS_PROP_CONSTANT_RGBA           = "CONSTANT_RGBA"
	PNTS_PROP_RGB                     = "RGB"
	PNTS_PROP_RGBA                    = "RGBA"
	PNTS_PROP_RGB565                  = "RGB565"
	PNTS_PROP_BATCH_ID                = "BATCH_ID"
)

const (
	PNTS_COLOR_FORMAT_RGB    = 0
	PNTS_COLOR_FORMAT_RGBA   = 1
	PNTS_COLOR_FORMAT_RGB565 = 2
	PNTS_COLOR_FORMAT_NONE   = 3
)

type PntsHeader struct {
	Header
	Magic                        [4]byte
	Version                      uint32
	ByteLength                   uint32
	FeatureTableJSONByteLength   uint32
	FeatureTableBinaryByteLength uint32
	BatchTableJSONByteLength     uint32
	BatchTableBinaryByteLength   uint32
}

func (h *PntsHeader) CalcSize() int64 {
	return 28
}

func (h *PntsHeader) GetByteLength() uint32 {
	return h.ByteLength
}

func (h *PntsHeader) GetFeatureTableJSONByteLength() uint32 {
	return h.FeatureTableJSONByteLength
}

func (h *PntsHeader) GetFeatureTableBinaryByteLength() uint32 {
	return h.FeatureTableBinaryByteLength
}

func (h *PntsHeader) GetBatchTableJSONByteLength() uint32 {
	return h.BatchTableJSONByteLength
}

func (h *PntsHeader) GetBatchTableBinaryByteLength() uint32 {
	return h.BatchTableBinaryByteLength
}

func (h *PntsHeader) SetFeatureTableJSONByteLength(n uint32) {
	h.FeatureTableJSONByteLength = n
}

func (h *PntsHeader) SetFeatureTableBinaryByteLength(n uint32) {
	h.FeatureTableBinaryByteLength = n
}

func (h *PntsHeader) SetBatchTableJSONByteLength(n uint32) {
	h.BatchTableJSONByteLength = n
}

func (h *PntsHeader) SetBatchTableBinaryByteLength(n uint32) {
	h.BatchTableBinaryByteLength = n
}

type PntsFeatureTableView struct {
	Position              [][3]float64
	PositionQuantized     [][3]uint16
	RGBA                  [][4]uint8
	RGB                   [][3]uint8
	RGB565                []uint16
	Normal                [][3]float32
	NormalOCT16P          [][2]uint8
	BatchId               interface{}
	PointsLength          uint32
	RtcCenter             []float64
	QuantizedVolumeOffset []float64
	QuantizedVolumeScale  []float64
	ConstantRGBA          []uint8
	BatchLength           *uint32
}

func PntsFeatureTableDecode(header map[string]interface{}, buff []byte) map[string]interface{} {
	return nil
}

func PntsFeatureTableEncode(header map[string]interface{}, data map[string]interface{}) []byte {
	return nil
}

type PointCloud struct {
	TileModel
	Header       PntsHeader
	FeatureTable FeatureTable
	BatchTable   BatchTable
}

func (m *PointCloud) SetFeatureTable(view PntsFeatureTableView) {
	m.FeatureTable.Header[PNTS_PROP_POSITION] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_VEC3}
	m.FeatureTable.Data[PNTS_PROP_POSITION] = view.Position

	m.FeatureTable.Header[PNTS_PROP_POSITION_QUANTIZED] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_VEC3}
	m.FeatureTable.Data[PNTS_PROP_POSITION_QUANTIZED] = view.PositionQuantized

	if view.RGBA != nil {
		m.FeatureTable.Header[PNTS_PROP_RGBA] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_BYTE, ContainerType: CONTAINER_TYPE_VEC4}
		m.FeatureTable.Data[PNTS_PROP_RGBA] = view.RGBA
	}

	if view.RGB != nil {
		m.FeatureTable.Header[PNTS_PROP_RGB] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_BYTE, ContainerType: CONTAINER_TYPE_VEC3}
		m.FeatureTable.Data[PNTS_PROP_RGB] = view.RGB
	}

	if view.RGB565 != nil {
		m.FeatureTable.Header[PNTS_PROP_RGB565] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[PNTS_PROP_RGB565] = view.RGB565
	}

	if view.Normal != nil {
		m.FeatureTable.Header[PNTS_PROP_NORMAL] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_VEC3}
		m.FeatureTable.Data[PNTS_PROP_NORMAL] = view.Normal
	}

	if view.NormalOCT16P != nil {
		m.FeatureTable.Header[PNTS_PROP_NORMAL_OCT32P] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_BYTE, ContainerType: CONTAINER_TYPE_VEC2}
		m.FeatureTable.Data[PNTS_PROP_NORMAL_OCT32P] = view.NormalOCT16P
	}

	if view.BatchId != nil {
		switch t := view.BatchId.(type) {
		case []uint8:
			m.FeatureTable.Header[PNTS_PROP_BATCH_ID] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_BYTE, ContainerType: CONTAINER_TYPE_SCALAR}
			m.FeatureTable.Header[PNTS_PROP_BATCH_ID] = t
		case []uint16:
			m.FeatureTable.Header[PNTS_PROP_BATCH_ID] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
			m.FeatureTable.Header[PNTS_PROP_BATCH_ID] = t
		case []uint32:
			m.FeatureTable.Header[PNTS_PROP_BATCH_ID] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_INT, ContainerType: CONTAINER_TYPE_SCALAR}
			m.FeatureTable.Header[PNTS_PROP_BATCH_ID] = t
		case []int64:
			max := maxBatchId(t)
			if max > 0xFFFF {
				m.FeatureTable.Header[PNTS_PROP_BATCH_ID] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_INT, ContainerType: CONTAINER_TYPE_SCALAR}
				out := make([]uint32, len(t))
				for i := range out {
					out[i] = uint32(t[i])
				}
				m.FeatureTable.Data[PNTS_PROP_BATCH_ID] = out
			} else if max > 0xFF {
				m.FeatureTable.Header[PNTS_PROP_BATCH_ID] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
				out := make([]uint16, len(t))
				for i := range out {
					out[i] = uint16(t[i])
				}
				m.FeatureTable.Data[PNTS_PROP_BATCH_ID] = out
			} else {
				m.FeatureTable.Header[PNTS_PROP_BATCH_ID] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_BYTE, ContainerType: CONTAINER_TYPE_SCALAR}
				out := make([]uint8, len(t))
				for i := range out {
					out[i] = uint8(t[i])
				}
				m.FeatureTable.Data[PNTS_PROP_BATCH_ID] = out
			}
		}
	}

	m.FeatureTable.Header[PNTS_PROP_POSITION_LENGTH] = view.PointsLength
	if view.RtcCenter != nil && len(view.RtcCenter) == 3 {
		m.FeatureTable.Header[PNTS_PROP_RTC_CENTER] = view.RtcCenter
	}
	if view.QuantizedVolumeOffset != nil && len(view.QuantizedVolumeOffset) == 3 {
		m.FeatureTable.Header[PNTS_PROP_QUANTIZED_VOLUME_OFFSET] = view.QuantizedVolumeOffset
	}
	if view.QuantizedVolumeScale != nil && len(view.QuantizedVolumeScale) == 3 {
		m.FeatureTable.Header[PNTS_PROP_QUANTIZED_VOLUME_SCALE] = view.QuantizedVolumeScale
	}
	if view.ConstantRGBA != nil && len(view.ConstantRGBA) == 4 {
		m.FeatureTable.Header[PNTS_PROP_CONSTANT_RGBA] = view.ConstantRGBA
	}
	if view.BatchLength != nil {
		m.FeatureTable.Header[PNTS_PROP_BATCH_LENGTH] = *view.BatchLength
	}
}

func (m *PointCloud) GetFeatureTableView() PntsFeatureTableView {
	ret := PntsFeatureTableView{}

	if t := m.FeatureTable.Data[PNTS_PROP_POSITION]; t != nil {
		ret.Position = t.([][3]float64)
	}

	if t := m.FeatureTable.Data[PNTS_PROP_POSITION_QUANTIZED]; t != nil {
		ret.PositionQuantized = t.([][3]uint16)
	}

	if t := m.FeatureTable.Data[PNTS_PROP_RGBA]; t != nil {
		ret.RGBA = t.([][4]uint8)
	}

	if t := m.FeatureTable.Data[PNTS_PROP_RGB]; t != nil {
		ret.RGB = t.([][3]uint8)
	}

	if t := m.FeatureTable.Data[PNTS_PROP_RGB565]; t != nil {
		ret.RGB565 = t.([]uint16)
	}

	if t := m.FeatureTable.Data[PNTS_PROP_NORMAL]; t != nil {
		ret.Normal = t.([][3]float32)
	}

	if t := m.FeatureTable.Data[PNTS_PROP_NORMAL_OCT32P]; t != nil {
		ret.NormalOCT16P = t.([][2]uint8)
	}

	if m.FeatureTable.Data[PNTS_PROP_BATCH_ID] != nil {
		ret.BatchId = m.FeatureTable.Data[PNTS_PROP_BATCH_ID]
	}

	ret.PointsLength = m.FeatureTable.Header[PNTS_PROP_POSITION_LENGTH].(uint32)

	if m.FeatureTable.Header[PNTS_PROP_RTC_CENTER] != nil {
		ret.RtcCenter = m.FeatureTable.Header[PNTS_PROP_RTC_CENTER].([]float64)
	}
	if m.FeatureTable.Header[PNTS_PROP_QUANTIZED_VOLUME_OFFSET] != nil {
		ret.QuantizedVolumeOffset = m.FeatureTable.Header[PNTS_PROP_QUANTIZED_VOLUME_OFFSET].([]float64)
	}
	if m.FeatureTable.Header[PNTS_PROP_QUANTIZED_VOLUME_SCALE] != nil {
		ret.QuantizedVolumeScale = m.FeatureTable.Header[PNTS_PROP_QUANTIZED_VOLUME_SCALE].([]float64)
	}
	if m.FeatureTable.Header[PNTS_PROP_CONSTANT_RGBA] != nil {
		ret.ConstantRGBA = m.FeatureTable.Header[PNTS_PROP_CONSTANT_RGBA].([]uint8)
	}
	if m.FeatureTable.Header[PNTS_PROP_BATCH_LENGTH] != nil {
		d := m.FeatureTable.Header[PNTS_PROP_BATCH_LENGTH].(uint32)
		ret.BatchLength = &d
	}
	return ret
}

func (m *PointCloud) GetHeader() Header {
	return &m.Header
}

func (m *PointCloud) GetFeatureTable() *FeatureTable {
	return &m.FeatureTable
}

func (m *PointCloud) GetBatchTable() *BatchTable {
	return &m.BatchTable
}

func (m *PointCloud) CalcSize() int64 {
	return m.Header.CalcSize() + m.FeatureTable.CalcSize() + m.BatchTable.CalcSize(m.FeatureTable.GetBatchLength())
}

func (m *PointCloud) Read(reader io.ReadSeeker) error {
	err := binary.Read(reader, littleEndian, &m.Header)
	if err != nil {
		return err
	}

	m.FeatureTable.decode = PntsFeatureTableDecode

	if err := m.FeatureTable.Read(reader, m.GetHeader()); err != nil {
		return err
	}

	if err := m.BatchTable.Read(reader, m.GetHeader(), m.FeatureTable.GetBatchLength()); err != nil {
		return err
	}

	return nil
}

func (m *PointCloud) Write(writer io.Writer) error {
	si := m.Header.CalcSize() + m.FeatureTable.CalcSize() + m.BatchTable.CalcSize(m.FeatureTable.GetBatchLength())

	m.Header.ByteLength = uint32(si)

	m.FeatureTable.encode = PntsFeatureTableEncode

	err := binary.Write(writer, littleEndian, m.Header)

	if err != nil {
		return err
	}

	if err := m.FeatureTable.Write(writer, m.GetHeader()); err != nil {
		return err
	}

	if err := m.BatchTable.Write(writer, m.GetHeader()); err != nil {
		return err
	}

	return nil
}
