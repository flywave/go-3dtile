package tile3d

import (
	"encoding/binary"
	"math"
)

func calcDataSize(data interface{}) int {
	switch data := data.(type) {
	case bool, int8, uint8, *bool, *int8, *uint8:
		return 1
	case []bool:
		return len(data)
	case []int8:
		return len(data)
	case []uint8:
		return len(data)
	case int16, uint16, *int16, *uint16:
		return 2
	case []int16:
		return 2 * len(data)
	case []uint16:
		return 2 * len(data)
	case int32, uint32, *int32, *uint32:
		return 4
	case []int32:
		return 4 * len(data)
	case []uint32:
		return 4 * len(data)
	case int64, uint64, *int64, *uint64:
		return 8
	case []int64:
		return 8 * len(data)
	case []uint64:
		return 8 * len(data)
	}
	return 0
}

func toByte(bs []byte) byte {
	return bs[0]
}

func writeByte(dst []byte, v byte) {
	dst[0] = v
}

func toByteArray(bs []byte, length int) []byte {
	data := make([]byte, length)
	for i := range bs {
		data[i] = bs[i]
	}
	return data
}

func writeByteArray(dst []byte, bs []byte) {
	for i := range bs {
		dst[i] = bs[i]
	}
}

func toShort(bs []byte, order binary.ByteOrder) int16 {
	return int16(order.Uint16(bs))
}

func writeShort(dst []byte, v int16, order binary.ByteOrder) {
	order.PutUint16(dst, uint16(v))
}

func toUnsignedShort(bs []byte, order binary.ByteOrder) uint16 {
	return order.Uint16(bs)
}

func writeUnsignedShort(dst []byte, v uint16, order binary.ByteOrder) {
	order.PutUint16(dst, v)
}

func toShortArray(bs []byte, length int, order binary.ByteOrder) []int16 {
	data := make([]int16, length)
	for i := 0; i < length; i++ {
		data[i] = int16(order.Uint16(bs[i*2 : (i+1)*2]))
	}
	return data
}

func writeShortArray(dst []byte, v []int16, order binary.ByteOrder) {
	_ = dst[len(v)*2-1]
	for i := 0; i < len(v); i++ {
		writeShort(dst[i*2:(i+1)*2], v[i], order)
	}
}

func toUnsignedShortArray(bs []byte, length int, order binary.ByteOrder) []uint16 {
	data := make([]uint16, length)
	for i := 0; i < length; i++ {
		data[i] = order.Uint16(bs[i*2 : (i+1)*2])
	}
	return data
}

func writeUnsignedShortArray(dst []byte, v []uint16, order binary.ByteOrder) {
	_ = dst[len(v)*2-1]
	for i := 0; i < len(v); i++ {
		writeUnsignedShort(dst[i*2:(i+1)*2], v[i], order)
	}
}

func toInt(bs []byte, order binary.ByteOrder) int32 {
	return int32(order.Uint32(bs))
}

func writeInt(dst []byte, v int32, order binary.ByteOrder) {
	order.PutUint32(dst, uint32(v))
}

func toUnsignedInt(bs []byte, order binary.ByteOrder) uint32 {
	return order.Uint32(bs)
}

func writeUnsignedInt(dst []byte, v uint32, order binary.ByteOrder) {
	order.PutUint32(dst, v)
}

func toIntArray(bs []byte, length int, order binary.ByteOrder) []int32 {
	data := make([]int32, length)
	for i := 0; i < length; i++ {
		data[i] = int32(order.Uint32(bs[i*4 : (i+1)*4]))
	}
	return data
}

func writeIntArray(dst []byte, v []int32, order binary.ByteOrder) {
	_ = dst[len(v)*4-1]
	for i := 0; i < len(v); i++ {
		writeInt(dst[i*4:(i+1)*4+1], v[i], order)
	}
}

func toUnsignedIntArray(bs []byte, length int, order binary.ByteOrder) []uint32 {
	data := make([]uint32, length)
	for i := 0; i < length; i++ {
		data[i] = order.Uint32(bs[i*4 : (i+1)*4+1])
	}
	return data
}

func writeUnsignedIntArray(dst []byte, v []uint32, order binary.ByteOrder) {
	_ = dst[len(v)*4-1]
	for i := 0; i < len(v); i++ {
		writeUnsignedInt(dst[i*4:(i+1)*4], v[i], order)
	}
}

func toLong(bs []byte, order binary.ByteOrder) int64 {
	return int64(order.Uint64(bs))
}

func writeLong(dst []byte, v int64, order binary.ByteOrder) {
	order.PutUint64(dst, uint64(v))
}

func toLongArray(bs []byte, length int, order binary.ByteOrder) []int64 {
	data := make([]int64, length)
	for i := 0; i < length; i++ {
		data[i] = int64(order.Uint64(bs[i*8 : (i+1)*8+1]))
	}
	return data
}

func writeLongArray(dst []byte, v []int64, order binary.ByteOrder) {
	_ = dst[len(v)*8-1]
	for i := 0; i < len(v); i++ {
		writeLong(dst[i*8:(i+1)*8+1], v[i], order)
	}
}

func toFloat(bs []byte) float32 {
	bits := littleEndian.Uint32(bs)
	return math.Float32frombits(bits)
}

func writeFloat(dst []byte, v float32) {
	bs := math.Float32bits(v)
	littleEndian.PutUint32(dst, bs)
}

func toFloatArray(bs []byte, length int) []float32 {
	data := make([]float32, length)
	for i := 0; i < length; i++ {
		data[i] = toFloat(bs[i*4 : (i+1)*4+1])
	}
	return data
}

func writeFloatArray(dst []byte, v []float32) {
	_ = dst[len(v)*4-1]
	for i := 0; i < len(v); i++ {
		writeFloat(dst[i*4:(i+1)*4+1], v[i])
	}
}

func toDouble(bs []byte) float64 {
	bits := binary.LittleEndian.Uint64(bs)
	return math.Float64frombits(bits)
}

func writeDouble(dst []byte, v float64) {
	bs := math.Float64bits(v)
	littleEndian.PutUint64(dst, bs)
}

func toDoubleArray(bs []byte, length int) []float64 {
	data := make([]float64, length)
	for i := 0; i < length; i++ {
		data[i] = toDouble(bs[i*4 : (i+1)*4+1])
	}
	return data
}

func writeDoubleArray(dst []byte, v []float64) {
	_ = dst[len(v)*8-1]
	for i := 0; i < len(v); i++ {
		writeDouble(dst[i*8:(i+1)*8+1], v[i])
	}
}

func toBoolean(bs []byte) bool {
	return bs[0] != 0
}

func writeBoolean(dst []byte, v bool) {
	_ = dst[0]
	if v {
		dst[0] = 1
	} else {
		dst[0] = 0
	}
}

func writeStringFix(dst []byte, v string, len int) {
	_ = dst[len]
	for i := 0; i < len; i++ {
		dst[i] = v[i]
	}
}
