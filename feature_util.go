package tile3d

import (
	"encoding/binary"
)

var (
	littleEndian = binary.LittleEndian
)

func getBatchId(header map[string]interface{}, buff []byte, length int) interface{} {
	objValue := header["BATCH_ID"]
	switch oref := objValue.(type) {
	case BinaryBodyReference:
		offset := oref.ByteOffset
		switch oref.ComponentType {
		case "UNSIGNED_BYTE":
			ret := make([]uint8, length)
			for i := 0; i < length; i++ {
				ret[i] = uint8(buff[offset+i])
			}
			return ret
		case "UNSIGNED_SHORT":
			ret := make([]uint16, length)
			for i := 0; i < length; i++ {
				ret[i] = littleEndian.Uint16(buff[offset+i*2 : offset+(i+1)*2])
			}
			return ret
		case "UNSIGNED_INT":
			ret := make([]uint32, length)
			for i := 0; i < length; i++ {
				ret[i] = littleEndian.Uint32(buff[offset+i*4 : offset+(i+1)*4])
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
		offset := oref.ByteOffset
		ret := make([]byte, length)
		for i := 0; i < length; i++ {
			ret[i] = buff[offset+i]
		}
		return ret
	case []byte:
		return oref
	}
	return nil
}

func getShortArrayFeatureValue(header map[string]interface{}, buff []byte, propName string, length int) []int16 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case BinaryBodyReference:
		offset := oref.ByteOffset
		ret := make([]int16, length)
		for i := 0; i < length; i++ {
			ret[i] = toShort(buff[offset+i*2:offset+i*2+2], littleEndian)
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
		offset := oref.ByteOffset
		ret := make([]uint16, length)
		for i := 0; i < length; i++ {
			ret[i] = toUnsignedShort(buff[offset+i*2:offset+i*2+2], littleEndian)
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

func getFloatVec3ArrayFeatureValue(header map[string]interface{}, buff []byte, propName string, length int) [][3]float32 {
	objValue := header[propName]
	switch oref := objValue.(type) {
	case BinaryBodyReference:
		offset := oref.ByteOffset
		ret := make([][3]float32, length)
		for i := 0; i < length; i++ {
			//ret[i][0] = toFloat(buff[offset+i*12 : offset+i*12+5])
			//ret[i][1] = toFloat(buff[offset+i*12+4 : offset+i*12+10])
			//ret[i][2] = toFloat(buff[offset+i*12+10 : offset+i*12+15])
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
	}
	return [3]float32{0, 0, 0}
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
	return nil
}

func getBinaryBodyReference(header map[string]interface{}, buff []byte, propName string) *BinaryBodyReference {
	return nil
}

func getBatchTableValue(header map[string]interface{}, buff []byte, propName string, batchId int) interface{} {
	return nil
}

func getBatchTableValuesFromRef(reference *BinaryBodyReference, buff []byte, propName string, batchLength int) []interface{} {
	return nil
}

func toReference(val interface{}) *BinaryBodyReference {
	return nil
}

func getBatchTableValues(header map[string]interface{}, buff []byte, propName string, batchLength int) []interface{} {
	return nil
}

func getBatchTableBinaryByte(reference *BinaryBodyReference, data []interface{}) []byte {
	return nil
}
