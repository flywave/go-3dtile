package tile3d

func calcPadding(offset, paddingUnit uint32) uint32 {
	padding := offset % paddingUnit
	if padding != 0 {
		padding = paddingUnit - padding
	}
	return padding
}

func paddingBytes(bytes []byte, srcLen int, paddingUnit uint32, paddingCode byte) {
	padding := calcPadding(uint32(srcLen), paddingUnit)

	for i := 0; i < int(padding); i++ {
		bytes[(srcLen)+i] = paddingCode
	}
}

func createPaddingBytes(bytes []byte, offset, paddingUnit uint32, paddingCode byte) []byte {
	padding := calcPadding(offset, paddingUnit)
	if padding == 0 {
		return bytes
	}
	for i := 0; i < int(padding); i++ {
		bytes = append(bytes, paddingCode)
	}
	return bytes
}
