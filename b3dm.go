package tile3d

import (
	"errors"
	"io"

	"github.com/qmuntal/gltf"
)

const (
	B3DM_MAGIC = "b3dm"
)

const (
	B3DM_PROP_BATCH_LENGTH = "BATCH_LENGTH"
	B3DM_PROP_RTC_CENTER   = "RTC_CENTER"
)

type B3dmHeader struct {
	Header
	Magic                        string
	Version                      uint32
	ByteLength                   uint32
	FeatureTableJSONByteLength   uint32
	FeatureTableBinaryByteLength uint32
	BatchTableJSONByteLength     uint32
	BatchTableBinaryByteLength   uint32
}

func (h *B3dmHeader) CalcSize() int64 {
	return 28
}

func (h *B3dmHeader) Read(r io.Reader) error {
	b := make([]byte, h.CalcSize())
	an, err := io.ReadFull(r, b)

	if err != nil || int64(an) != h.CalcSize() {
		return err
	}
	if int64(an) != h.CalcSize() {
		return errors.New("b3dm header must is 28!")
	}
	offset := 0
	h.Magic = string(b[offset : offset+4])
	if h.Magic != B3DM_MAGIC {
		return errors.New("b3dm magic must is b3dm!")
	}
	offset += 4
	h.Version = toUnsignedInt(b[offset:offset+4], littleEndian)
	offset += 4
	h.ByteLength = toUnsignedInt(b[offset:offset+4], littleEndian)
	offset += 4
	h.FeatureTableJSONByteLength = toUnsignedInt(b[offset:offset+4], littleEndian)
	offset += 4
	h.FeatureTableBinaryByteLength = toUnsignedInt(b[offset:offset+4], littleEndian)
	offset += 4
	h.BatchTableJSONByteLength = toUnsignedInt(b[offset:offset+4], littleEndian)
	offset += 4
	h.BatchTableBinaryByteLength = toUnsignedInt(b[offset:offset+4], littleEndian)
	return nil
}

func (h *B3dmHeader) Write(wr io.Writer) error {
	b := make([]byte, h.CalcSize())
	offset := 0
	writeStringFix(b[offset:offset+4], B3DM_MAGIC, 4)
	offset += 4
	writeUnsignedInt(b[offset:offset+4], h.Version, littleEndian)
	offset += 4
	writeUnsignedInt(b[offset:offset+4], h.ByteLength, littleEndian)
	offset += 4
	writeUnsignedInt(b[offset:offset+4], h.FeatureTableJSONByteLength, littleEndian)
	offset += 4
	writeUnsignedInt(b[offset:offset+4], h.FeatureTableBinaryByteLength, littleEndian)
	offset += 4
	writeUnsignedInt(b[offset:offset+4], h.BatchTableJSONByteLength, littleEndian)
	offset += 4
	writeUnsignedInt(b[offset:], h.BatchTableBinaryByteLength, littleEndian)
	return nil
}

func (h *B3dmHeader) GetByteLength() uint32 {
	return h.ByteLength
}

func (h *B3dmHeader) GetFeatureTableJSONByteLength() uint32 {
	return h.FeatureTableJSONByteLength
}

func (h *B3dmHeader) GetFeatureTableBinaryByteLength() uint32 {
	return h.FeatureTableBinaryByteLength
}

func (h *B3dmHeader) GetBatchTableJSONByteLength() uint32 {
	return h.BatchTableJSONByteLength
}

func (h *B3dmHeader) GetBatchTableBinaryByteLength() uint32 {
	return h.BatchTableBinaryByteLength
}

func (h *B3dmHeader) SetFeatureTableJSONByteLength(n uint32) {
	h.FeatureTableJSONByteLength = n
}

func (h *B3dmHeader) SetFeatureTableBinaryByteLength(n uint32) {
	h.FeatureTableBinaryByteLength = n
}

func (h *B3dmHeader) SetBatchTableJSONByteLength(n uint32) {
	h.BatchTableJSONByteLength = n
}

func (h *B3dmHeader) SetBatchTableBinaryByteLength(n uint32) {
	h.BatchTableBinaryByteLength = n
}

type B3dmFeatureTableView struct {
	BatchLength int
	RtcCenter   []float64
}

func B3dmFeatureTableConvert(header map[string]interface{}, buff []byte) map[string][]interface{} {
	return nil
}

type B3dm struct {
	TileModel
	Header       B3dmHeader
	FeatureTable FeatureTable
	BatchTable   BatchTable
	Model        *gltf.Document
}

func (m *B3dm) SetFeatureTable(view B3dmFeatureTableView) {
	m.FeatureTable.Header[B3DM_PROP_BATCH_LENGTH] = view.BatchLength
	if view.RtcCenter != nil && len(view.RtcCenter) == 3 {
		m.FeatureTable.Header[B3DM_PROP_RTC_CENTER] = view.RtcCenter
	}
}

func (m *B3dm) GetFeatureTableView() B3dmFeatureTableView {
	ret := B3dmFeatureTableView{}
	ret.BatchLength = m.FeatureTable.Header[B3DM_PROP_BATCH_LENGTH].(int)
	if m.FeatureTable.Header[B3DM_PROP_RTC_CENTER] != nil {
		ret.RtcCenter = m.FeatureTable.Header[B3DM_PROP_RTC_CENTER].([]float64)
	}
	return ret
}

func (m *B3dm) GetHeader() *Header {
	return &m.Header.Header
}

func (m *B3dm) GetFeatureTable() *FeatureTable {
	return &m.FeatureTable
}

func (m *B3dm) GetBatchTable() *BatchTable {
	return &m.BatchTable
}

func (m *B3dm) CalcSize() int64 {
	return m.Header.CalcSize() + m.FeatureTable.CalcSize() + m.BatchTable.CalcSize() + calcGltfSize(m.Model, 8)
}

func (m *B3dm) Read(reader io.ReadSeeker) error {
	if err := m.Header.Read(reader); err != nil {
		return err
	}

	if err := m.FeatureTable.Read(reader, *m.GetHeader()); err != nil {
		return err
	}

	//TODO batchLength
	if err := m.BatchTable.Read(reader, *m.GetHeader(), 0); err != nil {
		return err
	}

	var err1 error
	if m.Model, err1 = loadGltfFromByte(reader); err1 != nil {
		return err1
	}

	return nil
}

func (m *B3dm) Write(writer io.Writer) error {
	buf, err := getGltfBinary(m.Model, 8)
	if err != nil {
		return err
	}

	si := m.Header.CalcSize() + m.FeatureTable.CalcSize() + m.BatchTable.CalcSize() + int64(len(buf))

	m.Header.ByteLength = uint32(si)

	if err := m.Header.Write(writer); err != nil {
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
