package tile3d

func calcPadding(offset, paddingUnit int) int {
	padding := offset % paddingUnit
	if padding != 0 {
		padding = paddingUnit - padding
	}
	return padding
}

func paddingBytes(bytes []byte, srcLen, paddingUnit int, paddingCode byte) {
	padding := calcPadding(srcLen, paddingUnit)

	for i := 0; i < padding; i++ {
		bytes[srcLen+i] = paddingCode
	}
}

func createPaddingBytes(bytes []byte, offset, paddingUnit int, paddingCode byte) []byte {
	padding := calcPadding(offset, paddingUnit)
	if padding == 0 {
		return bytes
	}
	for i := 0; i < padding; i++ {
		bytes = append(bytes, paddingCode)
	}
	return bytes
}
