package tile3d

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	GEOM_MAGIC = "geom"
)

const (
	GEOM_PROP_BOXES               = "BOXES"
	GEOM_PROP_BOX_BATCH_IDS       = "BOX_BATCH_IDS"
	GEOM_PROP_CYLINDERS           = "CYLINDERS"
	GEOM_PROP_CYLINDER_BATCH_IDS  = "CYLINDER_BATCH_IDS"
	GEOM_PROP_ELLIPSOIDS          = "ELLIPSOIDS"
	GEOM_PROP_ELLIPSOID_BATCH_IDS = "ELLIPSOID_BATCH_IDS"
	GEOM_PROP_SPHERES             = "SPHERES"
	GEOM_PROP_SPHERE_BATCH_IDS    = "SPHERE_BATCH_IDS"
	GEOM_PROP_BOXES_LENGTH        = "BOXES_LENGTH"
	GEOM_PROP_CYLINDERS_LENGTH    = "CYLINDERS_LENGTH"
	GEOM_PROP_ELLIPSOIDS_LENGTH   = "ELLIPSOIDS_LENGTH"
	GEOM_PROP_SPHERES_LENGTH      = "SPHERES_LENGTH"
	GEOM_PROP_RTC_CENTER          = "RTC_CENTER"
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
	return 28
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

type GeomBox [16]float32

type GeomCylinder [16]float32

type GeomEllipsoid [16]float32

type GeomSphere [4]float32

type GeomFeatureTableView struct {
	Boxs             []GeomBox
	BoxBatchId       interface{}
	Cylinders        []GeomCylinder
	CylinderBatchId  interface{}
	Ellipsoids       []GeomEllipsoid
	EllipsoidBatchId interface{}
	Spheres          []GeomSphere
	SphereBatchId    interface{}
	RtcCenter        [3]float64
}

func GeomFeatureTableDecode(header map[string]interface{}, buff []byte) map[string]interface{} {
	ret := make(map[string]interface{})
	boxsLength := getIntegerScalarFeatureValue(header, buff, GEOM_PROP_BOXES_LENGTH)
	cylindersLength := getIntegerScalarFeatureValue(header, buff, GEOM_PROP_CYLINDERS_LENGTH)
	ellipsoidsLength := getIntegerScalarFeatureValue(header, buff, GEOM_PROP_ELLIPSOIDS_LENGTH)
	spheresLength := getIntegerScalarFeatureValue(header, buff, GEOM_PROP_SPHERES_LENGTH)

	ret[GEOM_PROP_RTC_CENTER] = getFloatVec3FeatureValue(header, buff, GEOM_PROP_RTC_CENTER)

	if boxsLength > 0 {
		ret[GEOM_PROP_BOX_BATCH_IDS] = getUnsignedShortBatchIDs(header, buff, GEOM_PROP_BOX_BATCH_IDS, int(boxsLength))
		box := getFloatArrayFeatureValue(header, buff, GEOM_PROP_BOXES, int(boxsLength*6))
		ret[GEOM_PROP_BOXES] = box
	}

	if cylindersLength > 0 {
		ret[GEOM_PROP_CYLINDER_BATCH_IDS] = getUnsignedShortBatchIDs(header, buff, GEOM_PROP_CYLINDER_BATCH_IDS, int(cylindersLength))
		cylinder := getFloatArrayFeatureValue(header, buff, GEOM_PROP_CYLINDERS, int(cylindersLength*6))
		ret[GEOM_PROP_CYLINDERS] = cylinder
	}

	if ellipsoidsLength > 0 {
		ret[GEOM_PROP_ELLIPSOID_BATCH_IDS] = getUnsignedShortBatchIDs(header, buff, GEOM_PROP_ELLIPSOID_BATCH_IDS, int(ellipsoidsLength))
		ellipsoid := getFloatArrayFeatureValue(header, buff, GEOM_PROP_ELLIPSOIDS, int(ellipsoidsLength*6))
		ret[GEOM_PROP_ELLIPSOIDS] = ellipsoid
	}

	if spheresLength > 0 {
		ret[GEOM_PROP_SPHERE_BATCH_IDS] = getUnsignedShortBatchIDs(header, buff, GEOM_PROP_SPHERE_BATCH_IDS, int(spheresLength))
		sphere := getFloatArrayFeatureValue(header, buff, GEOM_PROP_SPHERES, int(spheresLength*4))
		ret[GEOM_PROP_SPHERES] = sphere
	}
	return ret
}

func GeomFeatureTableEncode(header map[string]interface{}, data map[string]interface{}) []byte {
	var out []byte
	buf := bytes.NewBuffer(out)
	offset := 0

	if t := data[GEOM_PROP_BOX_BATCH_IDS]; t != nil {
		dt := t.([]uint16)
		binary.Write(buf, littleEndian, dt)
		header[GEOM_PROP_BOX_BATCH_IDS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 2)
	}

	if t := data[GEOM_PROP_BOXES]; t != nil {
		dt := t.([]float32)
		binary.Write(buf, littleEndian, dt)
		header[GEOM_PROP_BOXES] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 4)
	}

	if t := data[GEOM_PROP_CYLINDER_BATCH_IDS]; t != nil {
		dt := t.([]uint16)
		binary.Write(buf, littleEndian, dt)
		header[GEOM_PROP_CYLINDER_BATCH_IDS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 2)
	}

	if t := data[GEOM_PROP_CYLINDERS]; t != nil {
		dt := t.([]float32)
		binary.Write(buf, littleEndian, dt)
		header[GEOM_PROP_CYLINDERS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 4)
	}

	if t := data[GEOM_PROP_ELLIPSOID_BATCH_IDS]; t != nil {
		dt := t.([]uint16)
		binary.Write(buf, littleEndian, dt)
		header[GEOM_PROP_ELLIPSOID_BATCH_IDS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 2)
	}

	if t := data[GEOM_PROP_ELLIPSOIDS]; t != nil {
		dt := t.([]float32)
		binary.Write(buf, littleEndian, dt)
		header[GEOM_PROP_ELLIPSOIDS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 4)
	}

	if t := data[GEOM_PROP_SPHERE_BATCH_IDS]; t != nil {
		dt := t.([]uint16)
		binary.Write(buf, littleEndian, dt)
		header[GEOM_PROP_SPHERE_BATCH_IDS] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 2)
	}

	if t := data[GEOM_PROP_SPHERES]; t != nil {
		dt := t.([]float32)
		binary.Write(buf, littleEndian, dt)
		header[GEOM_PROP_SPHERES] = BinaryBodyReference{ByteOffset: uint32(offset), ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_SCALAR}
		offset += (len(dt) * 4)
	}

	return buf.Bytes()
}

type Geom struct {
	Header       GeomHeader
	FeatureTable FeatureTable
	BatchTable   BatchTable
}

func (m *Geom) SetFeatureTable(view GeomFeatureTableView) {
	if m.FeatureTable.Header == nil {
		m.FeatureTable.Header = make(map[string]interface{})
	}
	if m.FeatureTable.Data == nil {
		m.FeatureTable.Data = make(map[string]interface{})
	}

	m.FeatureTable.Header[GEOM_PROP_BOXES_LENGTH] = int32(len(view.Boxs))
	m.FeatureTable.Header[GEOM_PROP_CYLINDERS_LENGTH] = int32(len(view.Cylinders))
	m.FeatureTable.Header[GEOM_PROP_ELLIPSOIDS_LENGTH] = int32(len(view.Ellipsoids))
	m.FeatureTable.Header[GEOM_PROP_SPHERES_LENGTH] = int32(len(view.Spheres))

	m.FeatureTable.Header[GEOM_PROP_RTC_CENTER] = view.RtcCenter

	if view.BoxBatchId != nil {
		m.FeatureTable.Header[GEOM_PROP_BOX_BATCH_IDS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[GEOM_PROP_BOX_BATCH_IDS] = view.BoxBatchId
	}

	if view.Boxs != nil {
		m.FeatureTable.Header[GEOM_PROP_BOXES] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_SCALAR}
		raw := make([]float32, len(view.Boxs)*6)
		for i := 0; i < len(view.Boxs); i++ {
			copy(raw[i*6:i*6+6], view.Boxs[i][:])
		}
		m.FeatureTable.Data[GEOM_PROP_BOXES] = raw
	}

	if view.CylinderBatchId != nil {
		m.FeatureTable.Header[GEOM_PROP_CYLINDER_BATCH_IDS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[GEOM_PROP_CYLINDER_BATCH_IDS] = view.CylinderBatchId
	}

	if view.Cylinders != nil {
		m.FeatureTable.Header[GEOM_PROP_CYLINDERS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_SCALAR}
		raw := make([]float32, len(view.Cylinders)*6)
		for i := 0; i < len(view.Cylinders); i++ {
			copy(raw[i*6:i*6+6], view.Cylinders[i][:])
		}
		m.FeatureTable.Data[GEOM_PROP_CYLINDERS] = raw
	}

	if view.EllipsoidBatchId != nil {
		m.FeatureTable.Header[GEOM_PROP_ELLIPSOID_BATCH_IDS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[GEOM_PROP_ELLIPSOID_BATCH_IDS] = view.EllipsoidBatchId
	}

	if view.Ellipsoids != nil {
		m.FeatureTable.Header[GEOM_PROP_ELLIPSOIDS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_SCALAR}
		raw := make([]float32, len(view.Ellipsoids)*6)
		for i := 0; i < len(view.Ellipsoids); i++ {
			copy(raw[i*6:i*6+6], view.Ellipsoids[i][:])
		}
		m.FeatureTable.Data[GEOM_PROP_ELLIPSOIDS] = raw
	}

	if view.SphereBatchId != nil {
		m.FeatureTable.Header[GEOM_PROP_SPHERE_BATCH_IDS] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_UNSIGNED_SHORT, ContainerType: CONTAINER_TYPE_SCALAR}
		m.FeatureTable.Data[GEOM_PROP_SPHERE_BATCH_IDS] = view.SphereBatchId
	}

	if view.Spheres != nil {
		m.FeatureTable.Header[GEOM_PROP_SPHERES] = BinaryBodyReference{ComponentType: COMPONENT_TYPE_FLOAT, ContainerType: CONTAINER_TYPE_SCALAR}
		raw := make([]float32, len(view.Spheres)*6)
		for i := 0; i < len(view.Spheres); i++ {
			copy(raw[i*4:i*4+4], view.Spheres[i][:])
		}
		m.FeatureTable.Data[GEOM_PROP_SPHERES] = raw
	}
}

func (m *Geom) GetFeatureTableView() GeomFeatureTableView {
	ret := GeomFeatureTableView{}

	if t := m.FeatureTable.Data[GEOM_PROP_RTC_CENTER]; t != nil {
		copy(ret.RtcCenter[:], t.([]float64))
	}

	if t := m.FeatureTable.Data[GEOM_PROP_BOXES_LENGTH]; t != nil {
		ret.Boxs = make([]GeomBox, t.(uint32))
	}

	if t := m.FeatureTable.Data[GEOM_PROP_BOXES]; t != nil {
		boxarr := t.([]float32)
		for i := 0; i < len(ret.Boxs); i++ {
			copy(ret.Boxs[i][:], boxarr[i*6:i*6+6])
		}
	}

	if t := m.FeatureTable.Data[GEOM_PROP_BOX_BATCH_IDS]; t != nil {
		ret.BoxBatchId = t.([]uint16)
	}

	if t := m.FeatureTable.Data[GEOM_PROP_CYLINDERS_LENGTH]; t != nil {
		ret.Cylinders = make([]GeomCylinder, t.(uint32))
	}

	if t := m.FeatureTable.Data[GEOM_PROP_CYLINDERS]; t != nil {
		cylinderarr := t.([]float32)
		for i := 0; i < len(ret.Cylinders); i++ {
			copy(ret.Cylinders[i][:], cylinderarr[i*6:i*6+6])
		}
	}

	if t := m.FeatureTable.Data[GEOM_PROP_CYLINDER_BATCH_IDS]; t != nil {
		ret.CylinderBatchId = t.([]uint16)
	}

	if t := m.FeatureTable.Data[GEOM_PROP_ELLIPSOIDS_LENGTH]; t != nil {
		ret.Ellipsoids = make([]GeomEllipsoid, t.(uint32))
	}

	if t := m.FeatureTable.Data[GEOM_PROP_ELLIPSOIDS]; t != nil {
		ellipsoidarr := t.([]float32)
		for i := 0; i < len(ret.Ellipsoids); i++ {
			copy(ret.Ellipsoids[i][:], ellipsoidarr[i*6:i*6+6])
		}
	}

	if t := m.FeatureTable.Data[GEOM_PROP_ELLIPSOID_BATCH_IDS]; t != nil {
		ret.EllipsoidBatchId = t.([]uint16)
	}

	if t := m.FeatureTable.Data[GEOM_PROP_SPHERES_LENGTH]; t != nil {
		ret.Spheres = make([]GeomSphere, t.(uint32))
	}

	if t := m.FeatureTable.Data[GEOM_PROP_SPHERES]; t != nil {
		spherearr := t.([]float32)
		for i := 0; i < len(ret.Spheres); i++ {
			copy(ret.Spheres[i][:], spherearr[i*4:i*4+4])
		}
	}

	if t := m.FeatureTable.Data[GEOM_PROP_SPHERE_BATCH_IDS]; t != nil {
		ret.SphereBatchId = t.([]uint16)
	}

	return ret
}

func (m *Geom) GetHeader() Header {
	return &m.Header
}

func (m *Geom) CalcSize() int64 {
	return m.Header.CalcSize() + m.FeatureTable.CalcSize(m.GetHeader()) + m.BatchTable.CalcSize(m.GetHeader())
}

func (m *Geom) GetFeatureTable() *FeatureTable {
	return &m.FeatureTable
}

func (m *Geom) GetBatchTable() *BatchTable {
	return &m.BatchTable
}

func (m *Geom) Read(reader io.ReadSeeker) error {
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

func (m *Geom) Write(writer io.Writer) error {
	m.FeatureTable.encode = GeomFeatureTableEncode
	_ = GeomFeatureTableEncode(m.FeatureTable.Header, m.FeatureTable.Data)
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

	return nil
}
