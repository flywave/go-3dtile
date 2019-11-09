package tile3d

import (
	"encoding/binary"
	"io"

	"github.com/qmuntal/gltf"
)

const (
	I3DM_MAGIC         = "i3dm"
	I3DM_GLTF_URI      = 0
	I3DM_GLTF_EMBEDDED = 1
)

const (
	I3DM_PROP_POSITION                = "POSITION"
	I3DM_PROP_POSITION_QUANTIZED      = "POSITION_QUANTIZED"
	I3DM_PROP_QUANTIZED_VOLUME_OFFSET = "QUANTIZED_VOLUME_OFFSET"
	I3DM_PROP_QUANTIZED_VOLUME_SCALE  = "QUANTIZED_VOLUME_SCALE"
	I3DM_PROP_NORMAL_UP               = "NORMAL_UP"
	I3DM_PROP_NORMAL_RIGHT            = "NORMAL_RIGHT"
	I3DM_PROP_NORMAL_UP_OCT32P        = "NORMAL_UP_OCT32P"
	I3DM_PROP_NORMAL_RIGHT_OCT32P     = "NORMAL_RIGHT_OCT32P"
	I3DM_PROP_SCALE                   = "SCALE"
	I3DM_PROP_SCALE_NON_UNIFORM       = "SCALE_NON_UNIFORM"
	I3DM_PROP_BATCH_ID                = "BATCH_ID"
	I3DM_PROP_INSTANCES_LENGTH        = "INSTANCES_LENGTH"
	I3DM_PROP_RTC_CENTER              = "RTC_CENTER"
	I3DM_PROP_EAST_NORTH_UP           = "EAST_NORTH_UP"
)

type I3dmHeader struct {
	Header
	Magic                        [4]byte
	Version                      uint32
	ByteLength                   uint32
	FeatureTableJSONByteLength   uint32
	FeatureTableBinaryByteLength uint32
	BatchTableJSONByteLength     uint32
	BatchTableBinaryByteLength   uint32
	GltfFormat                   uint32
}

func (h *I3dmHeader) CalcSize() int64 {
	return 32
}

func (h *I3dmHeader) GetByteLength() uint32 {
	return h.ByteLength
}

func (h *I3dmHeader) GetFeatureTableJSONByteLength() uint32 {
	return h.FeatureTableJSONByteLength
}

func (h *I3dmHeader) GetFeatureTableBinaryByteLength() uint32 {
	return h.FeatureTableBinaryByteLength
}

func (h *I3dmHeader) GetBatchTableJSONByteLength() uint32 {
	return h.BatchTableJSONByteLength
}

func (h *I3dmHeader) GetBatchTableBinaryByteLength() uint32 {
	return h.BatchTableBinaryByteLength
}

func (h *I3dmHeader) SetFeatureTableJSONByteLength(n uint32) {
	h.FeatureTableJSONByteLength = n
}

func (h *I3dmHeader) SetFeatureTableBinaryByteLength(n uint32) {
	h.FeatureTableBinaryByteLength = n
}

func (h *I3dmHeader) SetBatchTableJSONByteLength(n uint32) {
	h.BatchTableJSONByteLength = n
}

func (h *I3dmHeader) SetBatchTableBinaryByteLength(n uint32) {
	h.BatchTableBinaryByteLength = n
}

type I3dmFeatureTableView struct {
	Position              [][3]float64
	PositionQuantized     [][3]uint16
	NormalRight           [][3]float32
	NormalUp              [][3]float32
	NormalRightOCT16P     [][2]uint16
	NormalUpOCT16P        [][2]uint16
	Scale                 []float32
	ScaleNONUniform       [][3]float32
	BatchId               interface{}
	InstanceLength        int
	RtcCenter             []float32
	QuantizedVolumeOffset []float32
	QuantizedVolumeScale  []float32
	EastNorthUp           *bool
}

func I3dmFeatureTableConvert(header map[string]interface{}, buff []byte) map[string]interface{} {
	return nil
}

type I3dm struct {
	TileModel
	Header       I3dmHeader
	FeatureTable FeatureTable
	BatchTable   BatchTable
	GltfUri      string
	Model        *gltf.Document
}

func (m *I3dm) SetFeatureTable(view I3dmFeatureTableView) {
	m.FeatureTable.Header[I3DM_PROP_POSITION] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_VEC3}
	m.FeatureTable.Data[I3DM_PROP_POSITION] = view.Position

	m.FeatureTable.Header[I3DM_PROP_POSITION_QUANTIZED] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_VEC3}
	m.FeatureTable.Data[I3DM_PROP_POSITION_QUANTIZED] = view.PositionQuantized

	if view.NormalUp != nil {
		m.FeatureTable.Header[I3DM_PROP_NORMAL_UP] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_VEC3}
		m.FeatureTable.Data[I3DM_PROP_NORMAL_UP] = view.NormalUp
	}

	if view.NormalRight != nil {
		m.FeatureTable.Header[I3DM_PROP_NORMAL_RIGHT] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_VEC3}
		m.FeatureTable.Data[I3DM_PROP_NORMAL_UP] = view.NormalRight
	}

	if view.NormalUpOCT16P != nil {
		m.FeatureTable.Header[I3DM_PROP_NORMAL_UP_OCT32P] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_VEC3}
		m.FeatureTable.Data[I3DM_PROP_NORMAL_UP_OCT32P] = view.NormalUpOCT16P
	}

	if view.NormalRightOCT16P != nil {
		m.FeatureTable.Header[I3DM_PROP_NORMAL_RIGHT_OCT32P] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_VEC3}
		m.FeatureTable.Data[I3DM_PROP_NORMAL_RIGHT_OCT32P] = view.NormalRightOCT16P
	}

	if view.Scale != nil {
		m.FeatureTable.Header[I3DM_PROP_SCALE] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[I3DM_PROP_SCALE] = view.Scale
	}

	if view.ScaleNONUniform != nil {
		m.FeatureTable.Header[I3DM_PROP_SCALE_NON_UNIFORM] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_VEC3}
		m.FeatureTable.Data[I3DM_PROP_SCALE_NON_UNIFORM] = view.ScaleNONUniform
	}

	if view.BatchId != nil {
		switch t := view.BatchId.(type) {
		case []uint8:
			m.FeatureTable.Header[I3DM_PROP_BATCH_ID] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_BYTE, ContainerType: CONTAINER_TYPE_SCALAR}
			m.FeatureTable.Data[I3DM_PROP_BATCH_ID] = t
		case []uint16:
			m.FeatureTable.Header[I3DM_PROP_BATCH_ID] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
			m.FeatureTable.Data[I3DM_PROP_BATCH_ID] = t
		case []uint32:
			m.FeatureTable.Header[I3DM_PROP_BATCH_ID] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_INT, ContainerType: CONTAINER_TYPE_SCALAR}
			m.FeatureTable.Data[I3DM_PROP_BATCH_ID] = t
		case []int64:
			max := maxBatchId(t)
			if max > 0xFFFF {
				m.FeatureTable.Header[I3DM_PROP_BATCH_ID] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_INT, ContainerType: CONTAINER_TYPE_SCALAR}
				out := make([]uint32, len(t))
				for i := range out {
					out[i] = uint32(t[i])
				}
				m.FeatureTable.Data[I3DM_PROP_BATCH_ID] = out
			} else if max > 0xFF {
				m.FeatureTable.Header[I3DM_PROP_BATCH_ID] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
				out := make([]uint16, len(t))
				for i := range out {
					out[i] = uint16(t[i])
				}
				m.FeatureTable.Data[I3DM_PROP_BATCH_ID] = out
			} else {
				m.FeatureTable.Header[I3DM_PROP_BATCH_ID] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_BYTE, ContainerType: CONTAINER_TYPE_SCALAR}
				out := make([]uint8, len(t))
				for i := range out {
					out[i] = uint8(t[i])
				}
				m.FeatureTable.Data[I3DM_PROP_BATCH_ID] = out
			}
		}
	}

	m.FeatureTable.Header[I3DM_PROP_INSTANCES_LENGTH] = view.InstanceLength
	if view.RtcCenter != nil && len(view.RtcCenter) == 3 {
		m.FeatureTable.Header[I3DM_PROP_RTC_CENTER] = view.RtcCenter
	}
	if view.QuantizedVolumeOffset != nil && len(view.QuantizedVolumeOffset) == 3 {
		m.FeatureTable.Header[I3DM_PROP_QUANTIZED_VOLUME_OFFSET] = view.QuantizedVolumeOffset
	}
	if view.QuantizedVolumeScale != nil && len(view.QuantizedVolumeScale) == 3 {
		m.FeatureTable.Header[I3DM_PROP_QUANTIZED_VOLUME_SCALE] = view.QuantizedVolumeScale
	}
	if view.EastNorthUp != nil {
		m.FeatureTable.Header[I3DM_PROP_EAST_NORTH_UP] = *view.EastNorthUp
	}
}

func (m *I3dm) GetFeatureTableView() I3dmFeatureTableView {
	ret := I3dmFeatureTableView{}

	if t := m.FeatureTable.Data[I3DM_PROP_POSITION]; t != nil {
		ret.Position = t.([][3]float64)
	}

	if t := m.FeatureTable.Data[I3DM_PROP_POSITION_QUANTIZED]; t != nil {
		ret.PositionQuantized = t.([][3]uint16)
	}

	if t := m.FeatureTable.Data[I3DM_PROP_NORMAL_UP]; t != nil {
		ret.NormalUp = t.([][3]float32)
	}

	if t := m.FeatureTable.Data[I3DM_PROP_NORMAL_RIGHT]; t != nil {
		ret.NormalRight = t.([][3]float32)
	}

	if t := m.FeatureTable.Data[I3DM_PROP_NORMAL_UP_OCT32P]; t != nil {
		ret.NormalUpOCT16P = t.([][2]uint16)
	}

	if t := m.FeatureTable.Data[I3DM_PROP_NORMAL_RIGHT_OCT32P]; t != nil {
		ret.NormalRightOCT16P = t.([][2]uint16)
	}

	if t := m.FeatureTable.Data[I3DM_PROP_SCALE]; t != nil {
		ret.Scale = t.([]float32)
	}

	if t := m.FeatureTable.Data[I3DM_PROP_SCALE_NON_UNIFORM]; t != nil {
		ret.ScaleNONUniform = t.([][3]float32)
	}

	if m.FeatureTable.Data[I3DM_PROP_BATCH_ID] != nil {
		ret.BatchId = m.FeatureTable.Data[I3DM_PROP_BATCH_ID]
	}

	ret.InstanceLength = m.FeatureTable.Header[I3DM_PROP_INSTANCES_LENGTH].(int)

	if m.FeatureTable.Header[I3DM_PROP_RTC_CENTER] != nil {
		ret.RtcCenter = m.FeatureTable.Header[I3DM_PROP_RTC_CENTER].([]float32)
	}

	if m.FeatureTable.Header[I3DM_PROP_QUANTIZED_VOLUME_OFFSET] != nil {
		ret.QuantizedVolumeOffset = m.FeatureTable.Header[I3DM_PROP_QUANTIZED_VOLUME_OFFSET].([]float32)
	}

	if m.FeatureTable.Header[I3DM_PROP_QUANTIZED_VOLUME_SCALE] != nil {
		ret.QuantizedVolumeScale = m.FeatureTable.Header[I3DM_PROP_QUANTIZED_VOLUME_SCALE].([]float32)
	}

	if m.FeatureTable.Header[I3DM_PROP_EAST_NORTH_UP] != nil {
		b := m.FeatureTable.Header[I3DM_PROP_EAST_NORTH_UP].(bool)
		ret.EastNorthUp = &b
	}
	return ret
}

func (m *I3dm) GetHeader() Header {
	return &m.Header
}

func (m *I3dm) GetFeatureTable() *FeatureTable {
	return &m.FeatureTable
}

func (m *I3dm) GetBatchTable() *BatchTable {
	return &m.BatchTable
}

func (m *I3dm) CalcSize() int64 {
	gltfSize := 0
	if m.Header.GltfFormat == 0 {
		gltfSize = len(m.GltfUri)
		gltfSize += calcPadding(gltfSize, 8)
	} else if m.Header.GltfFormat == 1 && m.Model != nil {
		gltfSize = int(calcGltfSize(m.Model, 8))
	} else {
		panic("GltfFormat must 0 or 1")
	}
	return m.Header.CalcSize() + m.FeatureTable.CalcSize() + m.BatchTable.CalcSize() + int64(gltfSize)
}

func (m *I3dm) Read(reader io.ReadSeeker) error {
	err := binary.Read(reader, littleEndian, &m.Header)
	if err != nil {
		return err
	}

	if err := m.FeatureTable.Read(reader, m.GetHeader()); err != nil {
		return err
	}

	//TODO batchLength
	if err := m.BatchTable.Read(reader, m.GetHeader(), 0); err != nil {
		return err
	}

	if m.Header.GltfFormat == 0 {
		var uri []byte
		if _, err := io.ReadAtLeast(reader, uri, 0); err != nil {
			return err
		}
		m.GltfUri = string(uri)
	} else if m.Header.GltfFormat == 1 {
		var err1 error
		if m.Model, err1 = loadGltfFromByte(reader); err1 != nil {
			return err1
		}
	} else {
		panic("GltfFormat must 0 or 1")
	}
	return nil
}

func (m *I3dm) Write(writer io.Writer) error {
	var buf []byte

	if m.Header.GltfFormat == 0 {
		buf = createPaddingBytes([]byte(m.GltfUri), len(m.GltfUri), 8, 0x20)
	} else if m.Header.GltfFormat == 1 {
		var err1 error
		if buf, err1 = getGltfBinary(m.Model, 8); err1 != nil {
			return err1
		}
	}

	si := m.Header.CalcSize() + m.FeatureTable.CalcSize() + m.BatchTable.CalcSize() + int64(len(buf))

	m.Header.ByteLength = uint32(si)

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

	if _, err := writer.Write(buf); err != nil {
		return err
	}

	return nil
}
