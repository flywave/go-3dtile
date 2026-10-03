package tile3d

import (
	"bytes"
	"encoding/binary"
	"math"
	"unsafe"
)

var (
	littleEndian = binary.LittleEndian
)

func getUnsignedShortBatchIDs(header map[string]interface{}, buff []byte, propName string, length int) []uint16 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case BinaryBodyReference:
		offset := oref.ByteOffset
		buf := bytes.NewBuffer(buff[offset:])
		switch oref.ComponentType {
		case "UNSIGNED_SHORT":
			ret := make([]uint16, length)
			err := binary.Read(buf, littleEndian, ret)
			if err != nil {
				return nil
			}
			return ret
		}
	}
	return nil
}

func getBatchLength(header map[string]interface{}, buff []byte, length int) interface{} {
	objValue := header["BATCH_ID"]
	switch oref := objValue.(type) {
	case BinaryBodyReference:
		offset := oref.ByteOffset
		buf := bytes.NewBuffer(buff[offset:])
		switch oref.ComponentType {
		case "UNSIGNED_BYTE":
			ret := make([]uint8, length)
			err := binary.Read(buf, littleEndian, ret)
			if err != nil {
				return nil
			}
			return ret
		case "UNSIGNED_SHORT":
			ret := make([]uint16, length)
			err := binary.Read(buf, littleEndian, ret)
			if err != nil {
				return nil
			}
			return ret
		case "UNSIGNED_INT":
			ret := make([]uint32, length)
			err := binary.Read(buf, littleEndian, ret)
			if err != nil {
				return nil
			}
			return ret
		}
	}
	return nil
}

func getUnsignedByteArrayFeatureValue(header map[string]interface{}, buff []byte, propName string, length int) []byte {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case BinaryBodyReference:
		offset := int(oref.ByteOffset)
		return buff[offset : offset+length]
	case []byte:
		return oref
	}
	return nil
}

func getShortArrayFeatureValue(header map[string]interface{}, buff []byte, propName string, length int) []int16 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case BinaryBodyReference:
		buf := bytes.NewBuffer(buff[oref.ByteOffset:])
		ret := make([]int16, length)
		err := binary.Read(buf, littleEndian, ret)
		if err != nil {
			return nil
		}
		return ret
	case []float64:
		ret := make([]int16, length)
		for i := 0; i < length; i++ {
			ret[i] = int16(oref[i])
		}
		return ret
	}
	return nil
}

func getUnsignedShortArrayFeatureValue(header map[string]interface{}, buff []byte, propName string, length int) []uint16 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case BinaryBodyReference:
		buf := bytes.NewBuffer(buff[oref.ByteOffset:])
		ret := make([]uint16, length)
		err := binary.Read(buf, littleEndian, ret)
		if err != nil {
			return nil
		}
		return ret
	case []float64:
		ret := make([]uint16, length)
		for i := 0; i < length; i++ {
			ret[i] = uint16(oref[i])
		}
		return ret
	}
	return nil
}

func getUnsignedIntArrayFeatureValue(header map[string]interface{}, buff []byte, propName string, length int) []uint32 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case BinaryBodyReference:
		buf := bytes.NewBuffer(buff[oref.ByteOffset:])
		ret := make([]uint32, length)
		err := binary.Read(buf, littleEndian, ret)
		if err != nil {
			return nil
		}
		return ret
	case []float64:
		ret := make([]uint32, length)
		for i := 0; i < length; i++ {
			ret[i] = uint32(oref[i])
		}
		return ret
	case []uint32:
		return oref
	}
	return nil
}

func getFloatVec3ArrayFeatureValue(header map[string]interface{}, buff []byte, propName string, length int) [][3]float32 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case BinaryBodyReference:
		buf := bytes.NewBuffer(buff[oref.ByteOffset:])
		ret := make([][3]float32, length)
		err := binary.Read(buf, littleEndian, ret)
		if err != nil {
			return nil
		}
		return ret
	case []float64:
		ret := make([][3]float32, length)
		for i := 0; i < length; i++ {
			ret[i][0] = float32(oref[i*3])
			ret[i][1] = float32(oref[i*3+1])
			ret[i][2] = float32(oref[i*3+2])
		}
		return ret
	}
	return nil
}

func getUnsignedShortVec3ArrayFeatureValue(header map[string]interface{}, buff []byte, propName string, length int) [][3]uint16 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case BinaryBodyReference:
		buf := bytes.NewBuffer(buff[oref.ByteOffset:])
		ret := make([][3]uint16, length)
		err := binary.Read(buf, littleEndian, ret)
		if err != nil {
			return nil
		}
		return ret
	case []float64:
		ret := make([][3]uint16, length)
		for i := 0; i < length; i++ {
			ret[i][0] = uint16(oref[i*3])
			ret[i][1] = uint16(oref[i*3+1])
			ret[i][2] = uint16(oref[i*3+2])
		}
		return ret
	}
	return nil
}

func getUnsignedIntegerScalarFeatureValue(header map[string]interface{}, buff []byte, propName string) uint32 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case float64:
		return uint32(oref)
	}
	return 0
}

func getIntegerScalarFeatureValue(header map[string]interface{}, buff []byte, propName string) int32 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case float64:
		return int32(oref)
	}
	return 0
}

func getLongScalarFeatureValue(header map[string]interface{}, buff []byte, propName string) int64 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case float64:
		return int64(oref)
	}
	return 0
}

func getFloatScalarFeatureValue(header map[string]interface{}, buff []byte, propName string) float32 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case float64:
		return float32(oref)
	}
	return 0
}

func getDoubleScalarFeatureValue(header map[string]interface{}, buff []byte, propName string) float64 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case float64:
		return oref
	}
	return 0
}

func getUnsignedByteVec3FeatureValue(header map[string]interface{}, buff []byte, propName string) [3]uint8 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case []float64:
		ret := [3]uint8{}
		for i := 0; i < 3; i++ {
			ret[i] = uint8(oref[i])
		}
		return ret
	}
	return [3]uint8{0, 0, 0}
}

func getUnsignedByteVec4FeatureValue(header map[string]interface{}, buff []byte, propName string) [4]uint8 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case []float64:
		ret := [4]uint8{}
		for i := 0; i < 4; i++ {
			ret[i] = uint8(oref[i])
		}
		return ret
	}
	return [4]uint8{0, 0, 0, 0}
}

func getFloatVec3FeatureValue(header map[string]interface{}, buff []byte, propName string) [3]float32 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case []float64:
		ret := [3]float32{}
		for i := 0; i < 3; i++ {
			ret[i] = float32(oref[i])
		}
		return ret
	case []interface{}:
		ret := [3]float32{}
		for i := 0; i < 3; i++ {
			ret[i] = float32(oref[i].(float64))
		}
		return ret
	}
	return [3]float32{0, 0, 0}
}

func getFloat64Vec3FeatureValue(header map[string]interface{}, buff []byte, propName string) [3]float64 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case []float64:
		ret := [3]float64{}
		for i := 0; i < 3; i++ {
			ret[i] = oref[i]
		}
		return ret
	case []interface{}:
		ret := [3]float64{}
		for i := 0; i < 3; i++ {
			ret[i] = oref[i].(float64)
		}
		return ret
	}
	return [3]float64{0, 0, 0}
}

func getFloatVec4FeatureValue(header map[string]interface{}, buff []byte, propName string) [4]float32 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case []float64:
		ret := [4]float32{}
		for i := 0; i < 4; i++ {
			ret[i] = float32(oref[i])
		}
		return ret
	}
	return [4]float32{0, 0, 0, 0}
}

func getFloatArrayFeatureValue(header map[string]interface{}, buff []byte, propName string, length int) []float32 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case BinaryBodyReference:
		buf := bytes.NewBuffer(buff[oref.ByteOffset:])
		ret := make([]float32, length)
		err := binary.Read(buf, littleEndian, ret)
		if err != nil {
			return nil
		}
		return ret
	case []float64:
		ret := make([]float32, length)
		for i := 0; i < length; i++ {
			ret[i] = float32(oref[i])
		}
		return ret
	case []interface{}:
		ret := make([]float32, length)
		for i := 0; i < length; i++ {
			ret[i] = float32(oref[i].(float64))
		}
		return ret
	case float64:
		ret := make([]float32, 1)
		ret[0] = float32(oref)
		return ret
	}
	return nil
}

func getBinaryBodyReference(header map[string]interface{}, propName string) *BinaryBodyReference {
	objValue := header[propName]
	if objValue == nil {
		return nil
	}
	switch t := objValue.(type) {
	case map[string]interface{}:
		bt := new(BinaryBodyReference)
		offset := t[REF_PROP_BYTE_OFFSET]
		if offset == nil {
			return nil
		}
		bt.ByteOffset = offset.(uint32)
		oct := t[REF_PROP_COMPONENT_TYPE]
		if oct != nil {
			bt.ComponentType = oct.(string)
		}
		oct = t[REF_PROP_TYPE]
		if oct != nil {
			bt.ContainerType = oct.(string)
		}
		return bt
	}
	return nil
}

func getBatchTableValue(header map[string]interface{}, buff []byte, propName string, batchLength int) interface{} {
	ref := getBinaryBodyReference(header, propName)
	return getBatchTableValuesFromRef(ref, buff, propName, batchLength)
}

func getBatchTableValuesFromRef(ref *BinaryBodyReference, buff []byte, propName string, batchLength int) interface{} {
	if ref != nil {
		offset := int(ref.ByteOffset)
		containerSize := ContainerTypeSize(ref.ContainerType)
		switch ref.ComponentType {
		case COMPONENT_TYPE_BYTE:
			if containerSize == 1 {
				return buff[offset+batchLength]
			}
			return buff[offset+batchLength*containerSize : offset+(batchLength+1)*containerSize]
		case COMPONENT_TYPE_UNSIGNED_BYTE:
			if containerSize == 1 {
				return uint8(buff[offset+batchLength])
			}
			out := make([]uint8, containerSize)
			for i := 0; i < containerSize; i++ {
				out[i] = uint8(buff[offset+batchLength*containerSize+i])
			}
			return out
		case COMPONENT_TYPE_SHORT:
			if containerSize == 1 {
				return int16(littleEndian.Uint16(buff[offset+batchLength : offset+batchLength+2]))
			}
			out := make([]int16, containerSize)
			for i := 0; i < containerSize; i++ {
				out[i] = int16(littleEndian.Uint16(buff[offset+batchLength*containerSize+i : offset+batchLength*containerSize+i+2]))
			}
			return out
		case COMPONENT_TYPE_UNSIGNED_SHORT:
			if containerSize == 1 {
				return littleEndian.Uint16(buff[offset+batchLength : offset+batchLength+2])
			}
			out := make([]uint16, containerSize)
			for i := 0; i < containerSize; i++ {
				out[i] = littleEndian.Uint16(buff[offset+batchLength*containerSize+i : offset+batchLength*containerSize+i+2])
			}
			return out
		case COMPONENT_TYPE_INT:
			if containerSize == 1 {
				return int32(littleEndian.Uint32(buff[offset+batchLength : offset+batchLength+4]))
			}
			out := make([]int32, containerSize)
			for i := 0; i < containerSize; i++ {
				out[i] = int32(littleEndian.Uint32(buff[offset+batchLength*containerSize+i : offset+batchLength*containerSize+i+4]))
			}
			return out
		case COMPONENT_TYPE_UNSIGNED_INT:
			if containerSize == 1 {
				return littleEndian.Uint32(buff[offset+batchLength : offset+batchLength+4])
			}
			out := make([]uint32, containerSize)
			for i := 0; i < containerSize; i++ {
				out[i] = littleEndian.Uint32(buff[offset+batchLength*containerSize+i : offset+batchLength*containerSize+i+4])
			}
			return out
		case COMPONENT_TYPE_FLOAT:
			if containerSize == 1 {
				i := littleEndian.Uint32(buff[offset+batchLength : offset+batchLength+4])
				return math.Float32frombits(i)
			}
			out := make([]float32, containerSize)
			for i := 0; i < containerSize; i++ {
				inte := littleEndian.Uint32(buff[offset+batchLength*containerSize+i : offset+batchLength*containerSize+i+4])
				out[i] = math.Float32frombits(inte)
			}
			return out
		case COMPONENT_TYPE_DOUBLE:
			if containerSize == 1 {
				i := littleEndian.Uint64(buff[offset+batchLength : offset+batchLength+8])
				return math.Float64frombits(i)
			}
			out := make([]float64, containerSize)
			for i := 0; i < containerSize; i++ {
				inte := littleEndian.Uint64(buff[offset+batchLength*containerSize+i : offset+batchLength*containerSize+i+8])
				out[i] = math.Float64frombits(inte)
			}
			return out
		}
	}
	return nil
}

// getBatchTableAllValuesFromRef 解码二进制批表属性的全部 batch 值。
// 与 getBatchTableValuesFromRef（第四个参数是单个批索引、只取该索引一个 batch 的值）不同，
// batchCount 是批数量：把 [offset, offset+batchCount*containerSize*compSize) 整段解码为
// 组件类型的完整切片；批 i 的分量位于下标 [i*containerSize, (i+1)*containerSize)（批主序平铺），
// 如 FLOAT VEC3 返回长度 3*batchCount 的 []float32。越界或未知类型返回 nil（不 panic）。
func getBatchTableAllValuesFromRef(ref *BinaryBodyReference, buff []byte, batchCount int) interface{} {
	if ref == nil || batchCount <= 0 {
		return nil
	}
	compSize := ComponentTypeSize(ref.ComponentType)
	if compSize <= 0 {
		return nil
	}
	containerSize := ContainerTypeSize(ref.ContainerType)
	if containerSize <= 0 {
		containerSize = 1 // 批表 type 缺省视为 SCALAR
	}
	offset := int(ref.ByteOffset)
	total := batchCount * containerSize
	if offset < 0 || offset+total*compSize > len(buff) {
		return nil
	}
	switch ref.ComponentType {
	case COMPONENT_TYPE_BYTE, COMPONENT_TYPE_UNSIGNED_BYTE:
		out := make([]byte, total)
		copy(out, buff[offset:offset+total])
		return out
	case COMPONENT_TYPE_SHORT:
		out := make([]int16, total)
		for i := range out {
			out[i] = int16(littleEndian.Uint16(buff[offset+2*i : offset+2*i+2]))
		}
		return out
	case COMPONENT_TYPE_UNSIGNED_SHORT:
		out := make([]uint16, total)
		for i := range out {
			out[i] = littleEndian.Uint16(buff[offset+2*i : offset+2*i+2])
		}
		return out
	case COMPONENT_TYPE_INT:
		out := make([]int32, total)
		for i := range out {
			out[i] = int32(littleEndian.Uint32(buff[offset+4*i : offset+4*i+4]))
		}
		return out
	case COMPONENT_TYPE_UNSIGNED_INT:
		out := make([]uint32, total)
		for i := range out {
			out[i] = littleEndian.Uint32(buff[offset+4*i : offset+4*i+4])
		}
		return out
	case COMPONENT_TYPE_FLOAT:
		out := make([]float32, total)
		for i := range out {
			out[i] = math.Float32frombits(littleEndian.Uint32(buff[offset+4*i : offset+4*i+4]))
		}
		return out
	case COMPONENT_TYPE_DOUBLE:
		out := make([]float64, total)
		for i := range out {
			out[i] = math.Float64frombits(littleEndian.Uint64(buff[offset+8*i : offset+8*i+8]))
		}
		return out
	}
	return nil
}

func getBatchTableBinaryByte(ref *BinaryBodyReference, data interface{}) []byte {
	if ref != nil {
		switch ref.ComponentType {
		case COMPONENT_TYPE_BYTE:
			switch d := data.(type) {
			case byte:
			case int8:
				return []byte{byte(d)}
			case []byte:
				return d
			case []int8:
				return *(*[]byte)(unsafe.Pointer(&d[0]))
			}
		case COMPONENT_TYPE_UNSIGNED_BYTE:
			switch d := data.(type) {
			case uint8:
				return []byte{byte(d)}
			case []uint8:
				return d
			}
		case COMPONENT_TYPE_SHORT:
			switch d := data.(type) {
			case int16:
				ret := make([]byte, 2)
				littleEndian.PutUint16(ret, uint16(d))
				return ret
			case []int16:
				ret := make([]byte, 2*len(d))
				binary.Write(bytes.NewBuffer(ret[:0]), littleEndian, d)
				return ret
			}
		case COMPONENT_TYPE_UNSIGNED_SHORT:
			switch d := data.(type) {
			case uint16:
				ret := make([]byte, 2)
				littleEndian.PutUint16(ret, uint16(d))
				return ret
			case []uint16:
				ret := make([]byte, 2*len(d))
				binary.Write(bytes.NewBuffer(ret[:0]), littleEndian, d)
				return ret
			}
		case COMPONENT_TYPE_INT:
			switch d := data.(type) {
			case int32:
				ret := make([]byte, 4)
				littleEndian.PutUint32(ret, uint32(d))
				return ret
			case []int32:
				ret := make([]byte, 4*len(d))
				binary.Write(bytes.NewBuffer(ret[:0]), littleEndian, d)
				return ret
			}
		case COMPONENT_TYPE_UNSIGNED_INT:
			switch d := data.(type) {
			case uint32:
				ret := make([]byte, 4)
				littleEndian.PutUint32(ret, d)
				return ret
			case []uint32:
				ret := make([]byte, 4*len(d))
				binary.Write(bytes.NewBuffer(ret[:0]), littleEndian, d)
				return ret
			}
		case COMPONENT_TYPE_FLOAT:
			switch d := data.(type) {
			case float32:
				ret := make([]byte, 4)
				littleEndian.PutUint32(ret, math.Float32bits(d))
				return ret
			case []float32:
				ret := make([]byte, 4*len(d))
				binary.Write(bytes.NewBuffer(ret[:0]), littleEndian, d)
				return ret
			}
		case COMPONENT_TYPE_DOUBLE:
			switch d := data.(type) {
			case float64:
				ret := make([]byte, 8)
				littleEndian.PutUint64(ret, math.Float64bits(d))
				return ret
			case []float64:
				ret := make([]byte, 8*len(d))
				binary.Write(bytes.NewBuffer(ret[:0]), littleEndian, d)
				return ret
			}
		}
	}
	return nil
}
