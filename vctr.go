package tile3d

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	VCTR_MAGIC = "vctr"
)

const (
	VCTR_PROP_POLYGONS_LENGTH         = "POLYGONS_LENGTH"
	VCTR_PROP_POLYLINES_LENGTH        = "POLYLINES_LENGTH"
	VCTR_PROP_POINTS_LENGTH           = "POINTS_LENGTH"
	VCTR_PROP_REGION                  = "REGION"
	VCTR_PROP_RTC_CENTER              = "RTC_CENTER"
	VCTR_PROP_POLYGON_COUNTS          = "POLYGON_COUNTS"
	VCTR_PROP_POLYGON_INDEX_COUNTS    = "POLYGON_INDEX_COUNTS"
	VCTR_PROP_POLYGON_MINIMUM_HEIGHTS = "POLYGON_MINIMUM_HEIGHTS"
	VCTR_PROP_POLYGON_MAXIMUM_HEIGHTS = "POLYGON_MAXIMUM_HEIGHTS"
	VCTR_PROP_POLYGON_BATCH_IDS       = "POLYGON_BATCH_IDS"
	VCTR_PROP_POLYLINE_COUNTS         = "POLYLINE_COUNTS"
	VCTR_PROP_POLYLINE_WIDTHS         = "POLYLINE_WIDTHS"
	VCTR_PROP_POLYLINE_BATCH_IDS      = "POLYLINE_BATCH_IDS"
	VCTR_PROP_POINT_BATCH_IDS         = "POINT_BATCH_IDS"
)

type VctrHeader struct {
	Magic                        [4]byte
	Version                      uint32
	ByteLength                   uint32
	FeatureTableJSONByteLength   uint32
	FeatureTableBinaryByteLength uint32
	BatchTableJSONByteLength     uint32
	BatchTableBinaryByteLength   uint32
	PolygonIndicesByteLength     uint32
	PolygonPositionsByteLength   uint32
	PolylinePositionsByteLength  uint32
	PointPositionsByteLength     uint32
}

func (h *VctrHeader) CalcSize() int64 {
	return 44
}

func (h *VctrHeader) GetByteLength() uint32 {
	return h.ByteLength
}

func (h *VctrHeader) GetFeatureTableJSONByteLength() uint32 {
	return h.FeatureTableJSONByteLength
}

func (h *VctrHeader) GetFeatureTableBinaryByteLength() uint32 {
	return h.FeatureTableBinaryByteLength
}

func (h *VctrHeader) GetBatchTableJSONByteLength() uint32 {
	return h.BatchTableJSONByteLength
}

func (h *VctrHeader) GetBatchTableBinaryByteLength() uint32 {
	return h.BatchTableBinaryByteLength
}

func (h *VctrHeader) GetPolygonIndicesByteLength() uint32 {
	return h.PolygonIndicesByteLength
}

func (h *VctrHeader) GetPolygonPositionsByteLength() uint32 {
	return h.PolygonPositionsByteLength
}

func (h *VctrHeader) GetPolylinePositionsByteLength() uint32 {
	return h.PolylinePositionsByteLength
}

func (h *VctrHeader) GetPointPositionsByteLength() uint32 {
	return h.PointPositionsByteLength
}

func (h *VctrHeader) SetFeatureTableJSONByteLength(n uint32) {
	h.FeatureTableJSONByteLength = n
}

func (h *VctrHeader) SetFeatureTableBinaryByteLength(n uint32) {
	h.FeatureTableBinaryByteLength = n
}

func (h *VctrHeader) SetBatchTableJSONByteLength(n uint32) {
	h.BatchTableJSONByteLength = n
}

func (h *VctrHeader) SetBatchTableBinaryByteLength(n uint32) {
	h.BatchTableBinaryByteLength = n
}

func (h *VctrHeader) SetPolygonIndicesByteLength(n uint32) {
	h.PolygonIndicesByteLength = n
}

func (h *VctrHeader) SetPolygonPositionsByteLength(n uint32) {
	h.PolygonPositionsByteLength = n
}

func (h *VctrHeader) SetPolylinePositionsByteLength(n uint32) {
	h.PolylinePositionsByteLength = n
}

func (h *VctrHeader) SetPointPositionsByteLength(n uint32) {
	h.PointPositionsByteLength = n
}

type VctrFeatureTableView struct {
	Region               *[6]float32
	RtcCenter            [3]float32
	PointsLength         uint32
	PointBatchId         interface{}
	PolygonsLength       uint32
	PolylinesLength      uint32
	PolylineBatchId      interface{}
	PolylineWidths       []uint16
	PolylineCounts       []uint32
	PolygonBatchId       interface{}
	PolygonMinimumHeight []float32
	PolygonMaximumHeight []float32
	PolygonIndexCounts   []uint32
	PolygonCounts        []uint32
}

func VctrFeatureTableDecode(header map[string]interface{}, buff []byte) map[string]interface{} {
	ret := make(map[string]interface{})
	polygonsLength := getIntegerScalarFeatureValue(header, buff, VCTR_PROP_POLYGONS_LENGTH)
	polylinesLength := getIntegerScalarFeatureValue(header, buff, VCTR_PROP_POLYLINES_LENGTH)
	pointsLength := getIntegerScalarFeatureValue(header, buff, VCTR_PROP_POINTS_LENGTH)

	ret[VCTR_PROP_REGION] = getFloatArrayFeatureValue(header, buff, VCTR_PROP_REGION, 6)
	ret[VCTR_PROP_RTC_CENTER] = getFloatVec3FeatureValue(header, buff, VCTR_PROP_RTC_CENTER)

	if pointsLength > 0 {
		ret[VCTR_PROP_POINTS_LENGTH] = pointsLength
		ret[VCTR_PROP_POINT_BATCH_IDS] = getUnsignedShortBatchIDs(header, buff, VCTR_PROP_POINT_BATCH_IDS, int(pointsLength))
	}

	if polylinesLength > 0 {
		ret[VCTR_PROP_POLYLINES_LENGTH] = polylinesLength
		ret[VCTR_PROP_POLYLINE_BATCH_IDS] = getUnsignedShortBatchIDs(header, buff, VCTR_PROP_POLYLINE_BATCH_IDS, int(polylinesLength))
		ret[VCTR_PROP_POLYLINE_COUNTS] = getUnsignedIntArrayFeatureValue(header, buff, VCTR_PROP_POLYLINE_COUNTS, 1)
		ret[VCTR_PROP_POLYLINE_WIDTHS] = getUnsignedShortArrayFeatureValue(header, buff, VCTR_PROP_POLYLINE_WIDTHS, int(polygonsLength))
	}

	if polygonsLength > 0 {
		ret[VCTR_PROP_POLYGONS_LENGTH] = polygonsLength
		ret[VCTR_PROP_POLYGON_BATCH_IDS] = getUnsignedShortBatchIDs(header, buff, VCTR_PROP_POLYGON_BATCH_IDS, int(polygonsLength))
		ret[VCTR_PROP_POLYGON_COUNTS] = getUnsignedIntArrayFeatureValue(header, buff, VCTR_PROP_POLYGON_COUNTS, int(polygonsLength))
		ret[VCTR_PROP_POLYGON_INDEX_COUNTS] = getUnsignedIntArrayFeatureValue(header, buff, VCTR_PROP_POLYGON_INDEX_COUNTS, int(polygonsLength))
		ret[VCTR_PROP_POLYGON_MAXIMUM_HEIGHTS] = getFloatArrayFeatureValue(header, buff, VCTR_PROP_POLYGON_MAXIMUM_HEIGHTS, int(polygonsLength))
		ret[VCTR_PROP_POLYGON_MINIMUM_HEIGHTS] = getFloatArrayFeatureValue(header, buff, VCTR_PROP_POLYGON_MINIMUM_HEIGHTS, int(polygonsLength))
	}

	return ret
}

func VctrFeatureTableEncode(header map[string]interface{}, data map[string]interface{}) []byte {
	var out []byte
	buf := bytes.NewBuffer(out)
	offset := 0

	if t := data[VCTR_PROP_POINT_BATCH_IDS]; t != nil {
		dt := t.([]uint16)
		binary.Write(buf, littleEndian, dt)
		header[VCTR_PROP_POINT_BATCH_IDS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 2)
	}

	if t := data[VCTR_PROP_POLYLINE_BATCH_IDS]; t != nil {
		dt := t.([]uint16)
		binary.Write(buf, littleEndian, dt)
		header[VCTR_PROP_POLYLINE_BATCH_IDS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 2)
	}

	if t := data[VCTR_PROP_POLYLINE_COUNTS]; t != nil {
		dt := t.([]uint32)
		binary.Write(buf, littleEndian, dt)
		header[VCTR_PROP_POLYLINE_COUNTS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_UNSIGNED_INT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 4)
	}

	if t := data[VCTR_PROP_POLYLINE_WIDTHS]; t != nil {
		dt := t.([]uint16)
		binary.Write(buf, littleEndian, dt)
		header[VCTR_PROP_POLYLINE_WIDTHS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 2)
	}

	if t := data[VCTR_PROP_POLYGON_BATCH_IDS]; t != nil {
		dt := t.([]uint16)
		binary.Write(buf, littleEndian, dt)
		header[VCTR_PROP_POLYGON_BATCH_IDS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 2)
	}

	if t := data[VCTR_PROP_POLYGON_COUNTS]; t != nil {
		dt := t.([]uint32)
		binary.Write(buf, littleEndian, dt)
		header[VCTR_PROP_POLYGON_COUNTS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_UNSIGNED_INT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 4)
	}

	if t := data[VCTR_PROP_POLYGON_INDEX_COUNTS]; t != nil {
		dt := t.([]uint32)
		binary.Write(buf, littleEndian, dt)
		header[VCTR_PROP_POLYGON_INDEX_COUNTS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_UNSIGNED_INT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 4)
	}

	if t := data[VCTR_PROP_POLYGON_MAXIMUM_HEIGHTS]; t != nil {
		dt := t.([]float32)
		binary.Write(buf, littleEndian, dt)
		header[VCTR_PROP_POLYGON_MAXIMUM_HEIGHTS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 4)
	}

	if t := data[VCTR_PROP_POLYGON_MINIMUM_HEIGHTS]; t != nil {
		dt := t.([]float32)
		binary.Write(buf, littleEndian, dt)
		header[VCTR_PROP_POLYGON_MINIMUM_HEIGHTS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 4)
	}

	return buf.Bytes()
}

type VctrIndices [][3]uint32

func (m VctrIndices) CalcSize(header Header) int64 {
	return int64(len(m) * 3 * 4)
}

func (m VctrIndices) Read(reader io.ReadSeeker, header Header) error {
	ch := header.(*VctrHeader)
	m = make([][3]uint32, int(ch.GetPolygonIndicesByteLength()/3/4))
	err := binary.Read(reader, littleEndian, m)
	if err != nil {
		return err
	}
	return nil
}

func (m VctrIndices) Write(writer io.Writer, header Header) error {
	err := binary.Write(writer, littleEndian, m)

	if err != nil {
		return err
	}
	if header != nil {
		ch := header.(*VctrHeader)
		ch.SetPolygonIndicesByteLength(uint32(m.CalcSize(header)))
	}
	return nil
}

type VctrPolygons struct {
	p [][2]int
}

func (m *VctrPolygons) Add(pt [2]int) {
	m.p = append(m.p, pt)
}

func (m *VctrPolygons) encode() (us, vs []uint16) {
	us, vs = encodePolygonPoints(m.p)
	return
}

func (m *VctrPolygons) decode(us, vs []uint16) {
	m.p = decodePolygonPoints(us, vs)
}

func (m *VctrPolygons) CalcSize(header Header) int64 {
	return int64(len(m.p) * 2 * 2)
}

func (m *VctrPolygons) Read(reader io.ReadSeeker, header Header) error {
	ch := header.(*VctrHeader)
	us := make([]uint16, int(ch.GetPolygonPositionsByteLength()/2/4))
	vs := make([]uint16, int(ch.GetPolygonPositionsByteLength()/2/4))

	err := binary.Read(reader, littleEndian, us)
	if err != nil {
		return err
	}
	err = binary.Read(reader, littleEndian, vs)
	if err != nil {
		return err
	}
	m.decode(us, vs)
	return nil
}

func (m *VctrPolygons) Write(writer io.Writer, header Header) error {
	us, vs := m.encode()
	err := binary.Write(writer, littleEndian, us)

	if err != nil {
		return err
	}
	err = binary.Write(writer, littleEndian, vs)

	if err != nil {
		return err
	}
	if header != nil {
		ch := header.(*VctrHeader)
		ch.SetPolygonPositionsByteLength(uint32(m.CalcSize(header)))
	}
	return nil
}

type VctrPolylines struct {
	p [][3]int
}

func (m *VctrPolylines) Add(pt [3]int) {
	m.p = append(m.p, pt)
}

func (m *VctrPolylines) encode() (us, vs, hs []uint16) {
	us, vs, hs = encodePoints(m.p)
	return
}

func (m *VctrPolylines) decode(us, vs, hs []uint16) {
	m.p = decodePoints(us, vs, hs)
}

func (m *VctrPolylines) CalcSize(header Header) int64 {
	return int64(len(m.p) * 3 * 2)
}

func (m *VctrPolylines) Read(reader io.ReadSeeker, header Header) error {
	ch := header.(*VctrHeader)
	us := make([]uint16, int(ch.GetPolylinePositionsByteLength()/2/4))
	vs := make([]uint16, int(ch.GetPolylinePositionsByteLength()/2/4))
	hs := make([]uint16, int(ch.GetPolylinePositionsByteLength()/2/4))

	err := binary.Read(reader, littleEndian, us)
	if err != nil {
		return err
	}
	err = binary.Read(reader, littleEndian, vs)
	if err != nil {
		return err
	}
	err = binary.Read(reader, littleEndian, hs)
	if err != nil {
		return err
	}
	m.decode(us, vs, hs)
	return nil
}

func (m *VctrPolylines) Write(writer io.Writer, header Header) error {
	us, vs, hs := m.encode()
	err := binary.Write(writer, littleEndian, us)

	if err != nil {
		return err
	}
	err = binary.Write(writer, littleEndian, vs)

	if err != nil {
		return err
	}
	err = binary.Write(writer, littleEndian, hs)

	if err != nil {
		return err
	}
	if header != nil {
		ch := header.(*VctrHeader)
		ch.SetPolylinePositionsByteLength(uint32(m.CalcSize(header)))
	}
	return nil
}

type VctrPoints struct {
	p [][3]int
}

func (m *VctrPoints) Add(pt [3]int) {
	m.p = append(m.p, pt)
}

func (m *VctrPoints) encode() (us, vs, hs []uint16) {
	us, vs, hs = encodePoints(m.p)
	return
}

func (m *VctrPoints) decode(us, vs, hs []uint16) {
	m.p = decodePoints(us, vs, hs)
}

func (m *VctrPoints) CalcSize(header Header) int64 {
	return int64(len(m.p) * 3 * 2)
}

func (m *VctrPoints) Read(reader io.ReadSeeker, header Header) error {
	ch := header.(*VctrHeader)
	us := make([]uint16, int(ch.PointPositionsByteLength/2/3))
	vs := make([]uint16, int(ch.PointPositionsByteLength/2/3))
	hs := make([]uint16, int(ch.PointPositionsByteLength/2/3))

	err := binary.Read(reader, littleEndian, us)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = binary.Read(reader, littleEndian, vs)
	if err != nil {
		return err
	}
	err = binary.Read(reader, littleEndian, hs)
	if err != nil {
		return err
	}
	m.decode(us, vs, hs)
	return nil
}

func (m *VctrPoints) Write(writer io.Writer, header Header) error {
	us, vs, hs := m.encode()
	err := binary.Write(writer, littleEndian, us)

	if err != nil {
		return err
	}
	err = binary.Write(writer, littleEndian, vs)

	if err != nil {
		return err
	}
	err = binary.Write(writer, littleEndian, hs)

	if err != nil {
		return err
	}
	if header != nil {
		ch := header.(*VctrHeader)
		ch.PointPositionsByteLength = uint32(m.CalcSize(header))
	}
	return nil
}

type Vctr struct {
	Header       VctrHeader
	FeatureTable FeatureTable
	BatchTable   BatchTable
	Indices      VctrIndices
	Polygons     VctrPolygons
	Polylines    VctrPolylines
	Points       VctrPoints
}

func (m *Vctr) SetFeatureTable(view VctrFeatureTableView) {
	if m.FeatureTable.Header == nil {
		m.FeatureTable.Header = make(map[string]interface{})
	}
	if m.FeatureTable.Data == nil {
		m.FeatureTable.Data = make(map[string]interface{})
	}

	m.FeatureTable.Header[VCTR_PROP_POLYGONS_LENGTH] = view.PolygonsLength
	m.FeatureTable.Header[VCTR_PROP_POLYLINES_LENGTH] = view.PolylinesLength
	m.FeatureTable.Header[VCTR_PROP_POINTS_LENGTH] = view.PointsLength

	m.FeatureTable.Header[VCTR_PROP_RTC_CENTER] = view.RtcCenter[:]

	if view.Region != nil && len(view.Region) == 6 {
		m.FeatureTable.Header[VCTR_PROP_REGION] = view.Region[:]
	}

	if view.PointBatchId != nil {
		m.FeatureTable.Header[VCTR_PROP_POINT_BATCH_IDS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[VCTR_PROP_POINT_BATCH_IDS] = view.PointBatchId
	}

	if view.PolylineBatchId != nil {
		m.FeatureTable.Header[VCTR_PROP_POLYLINE_BATCH_IDS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[VCTR_PROP_POLYLINE_BATCH_IDS] = view.PolylineBatchId
	}

	if view.PolylineCounts != nil {
		m.FeatureTable.Header[VCTR_PROP_POLYLINE_COUNTS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_INT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[VCTR_PROP_POLYLINE_COUNTS] = view.PolylineCounts
	}

	if view.PolylineWidths != nil {
		m.FeatureTable.Header[VCTR_PROP_POLYLINE_WIDTHS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[VCTR_PROP_POLYLINE_WIDTHS] = view.PolylineWidths
	}

	if view.PolygonBatchId != nil {
		m.FeatureTable.Header[VCTR_PROP_POLYGON_BATCH_IDS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[VCTR_PROP_POLYGON_BATCH_IDS] = view.PolygonBatchId
	}

	if view.PolygonCounts != nil {
		m.FeatureTable.Header[VCTR_PROP_POLYGON_COUNTS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_INT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[VCTR_PROP_POLYGON_COUNTS] = view.PolygonCounts
	}

	if view.PolygonIndexCounts != nil {
		m.FeatureTable.Header[VCTR_PROP_POLYGON_INDEX_COUNTS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_INT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[VCTR_PROP_POLYGON_INDEX_COUNTS] = view.PolygonIndexCounts
	}

	if view.PolygonMaximumHeight != nil {
		m.FeatureTable.Header[VCTR_PROP_POLYGON_MAXIMUM_HEIGHTS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[VCTR_PROP_POLYGON_MAXIMUM_HEIGHTS] = view.PolygonMaximumHeight
	}

	if view.PolygonMinimumHeight != nil {
		m.FeatureTable.Header[VCTR_PROP_POLYGON_MINIMUM_HEIGHTS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[VCTR_PROP_POLYGON_MINIMUM_HEIGHTS] = view.PolygonMinimumHeight
	}
}

func (m *Vctr) GetFeatureTableView() VctrFeatureTableView {
	ret := VctrFeatureTableView{}

	if t := m.FeatureTable.Data[VCTR_PROP_POLYGONS_LENGTH]; t != nil {
		ret.PolygonsLength = t.(uint32)
	}

	if t := m.FeatureTable.Data[VCTR_PROP_POLYLINES_LENGTH]; t != nil {
		ret.PolylinesLength = t.(uint32)
	}

	if t := m.FeatureTable.Data[VCTR_PROP_POINTS_LENGTH]; t != nil {
		ret.PointsLength = t.(uint32)
	}

	if t := m.FeatureTable.Data[VCTR_PROP_RTC_CENTER]; t != nil {
		copy(ret.RtcCenter[:], t.([]float32))
	}

	if t := m.FeatureTable.Data[VCTR_PROP_REGION]; t != nil {
		copy(ret.Region[:], t.([]float32))
	}

	if t := m.FeatureTable.Data[VCTR_PROP_POINT_BATCH_IDS]; t != nil {
		ret.PointBatchId = t.([]uint16)
	}

	if t := m.FeatureTable.Data[VCTR_PROP_POLYLINE_BATCH_IDS]; t != nil {
		ret.PolylineBatchId = t.([]uint16)
	}

	if t := m.FeatureTable.Data[VCTR_PROP_POLYLINE_COUNTS]; t != nil {
		ret.PolylineCounts = t.([]uint32)
	}

	if t := m.FeatureTable.Data[VCTR_PROP_POLYLINE_WIDTHS]; t != nil {
		ret.PolylineWidths = t.([]uint16)
	}

	if t := m.FeatureTable.Data[VCTR_PROP_POLYGON_BATCH_IDS]; t != nil {
		ret.PolygonBatchId = t.([]uint16)
	}

	if t := m.FeatureTable.Data[VCTR_PROP_POLYGON_COUNTS]; t != nil {
		ret.PolygonCounts = t.([]uint32)
	}

	if t := m.FeatureTable.Data[VCTR_PROP_POLYGON_INDEX_COUNTS]; t != nil {
		ret.PolygonIndexCounts = t.([]uint32)
	}

	if t := m.FeatureTable.Data[VCTR_PROP_POLYGON_MAXIMUM_HEIGHTS]; t != nil {
		ret.PolygonMaximumHeight = t.([]float32)
	}

	if t := m.FeatureTable.Data[VCTR_PROP_POLYGON_MINIMUM_HEIGHTS]; t != nil {
		ret.PolygonMinimumHeight = t.([]float32)
	}

	return ret
}

func (m *Vctr) GetHeader() Header {
	return &m.Header
}

func (m *Vctr) GetFeatureTable() *FeatureTable {
	return &m.FeatureTable
}

func (m *Vctr) GetBatchTable() *BatchTable {
	return &m.BatchTable
}

func (m *Vctr) GetIndices() VctrIndices {
	return m.Indices
}

func (m *Vctr) GetPolygons() VctrPolygons {
	return m.Polygons
}

func (m *Vctr) GetPolylines() VctrPolylines {
	return m.Polylines
}

func (m *Vctr) GetPoints() VctrPoints {
	return m.Points
}

func (m *Vctr) CalcSize() int64 {
	si := m.Header.CalcSize() + m.FeatureTable.CalcSize(m.GetHeader()) + m.BatchTable.CalcSize(m.GetHeader())

	if m.Indices != nil {
		si += m.Indices.CalcSize(m.GetHeader())
	}

	if m.Polygons.p != nil {
		si += m.Polygons.CalcSize(m.GetHeader())
	}

	if m.Polylines.p != nil {
		si += m.Polylines.CalcSize(m.GetHeader())
	}

	if m.Points.p != nil {
		si += m.Points.CalcSize(m.GetHeader())
	}
	return si
}

func (m *Vctr) Read(reader io.ReadSeeker) error {
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

	if err := m.Indices.Read(reader, m.GetHeader()); err != nil {
		return err
	}

	if err := m.Polygons.Read(reader, m.GetHeader()); err != nil {
		return err
	}

	if err := m.Polylines.Read(reader, m.GetHeader()); err != nil {
		return err
	}

	if err := m.Points.Read(reader, m.GetHeader()); err != nil {
		return err
	}

	return nil
}

func (m *Vctr) Write(writer io.Writer) error {
	m.FeatureTable.encode = VctrFeatureTableEncode
	_ = VctrFeatureTableEncode(m.FeatureTable.Header, m.FeatureTable.Data)
	si := m.Header.CalcSize() + m.FeatureTable.CalcSize(m.GetHeader()) + m.BatchTable.CalcSize(m.GetHeader())

	m.Header.ByteLength = uint32(si)

	err := binary.Write(writer, littleEndian, m.Header)

	if err != nil {
		return err
	}

	if err := m.FeatureTable.Write(writer, nil); err != nil {
		return err
	}

	if err := m.BatchTable.Write(writer, nil); err != nil {
		return err
	}

	if m.Indices != nil {
		if err := m.Indices.Write(writer, nil); err != nil {
			return err
		}
	}

	if m.Polygons.p != nil {
		if err := m.Polygons.Write(writer, nil); err != nil {
			return err
		}
	}

	if m.Polylines.p != nil {
		if err := m.Polylines.Write(writer, nil); err != nil {
			return err
		}
	}

	if m.Points.p != nil {
		if err := m.Points.Write(writer, nil); err != nil {
			return err
		}
	}

	return nil
}
