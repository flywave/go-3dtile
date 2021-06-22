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

func encodeZigZag(i int) uint16 {
	return uint16((i >> 15) ^ (i << 1))
}

func decodeZigZag(encoded uint16) int {
	unsignedEncoded := int(encoded)
	return unsignedEncoded>>1 ^ -(unsignedEncoded & 1)
}

func encodePolygonPoints(points [][2]int) (us, vs []uint16) {
	us = make([]uint16, len(points))
	vs = make([]uint16, len(points))

	lastU := int(0)
	lastV := int(0)

	for i := 0; i < len(points); i++ {
		u := points[i][0]
		v := points[i][1]

		us[i] = encodeZigZag(u - lastU)
		vs[i] = encodeZigZag(v - lastV)

		lastU = u
		lastV = v
	}
	return
}

func encodePoints(points [][3]int) (us, vs, hs []uint16) {
	us = make([]uint16, len(points))
	vs = make([]uint16, len(points))
	hs = make([]uint16, len(points))

	lastU := int(0)
	lastV := int(0)
	lastH := int(0)

	for i := 0; i < len(points); i++ {
		u := points[i][0]
		v := points[i][1]
		h := points[i][2]

		us[i] = encodeZigZag(u - lastU)
		vs[i] = encodeZigZag(v - lastV)
		hs[i] = encodeZigZag(h - lastH)

		lastU = u
		lastV = v
		lastH = h
	}
	return
}

func decodePolygonPoints(us, vs []uint16) [][2]int {
	u := int(0)
	v := int(0)
	pos := make([][2]int, len(us))

	for i := 0; i < len(us); i++ {
		u += decodeZigZag(us[i])
		v += decodeZigZag(vs[i])

		pos[i][0] = u
		pos[i][1] = v
	}
	return pos
}

func decodePoints(us, vs, hs []uint16) [][3]int {
	u := int(0)
	v := int(0)
	height := int(0)

	pos := make([][3]int, len(us))

	for i := 0; i < len(us); i++ {
		u += decodeZigZag(us[i])
		v += decodeZigZag(vs[i])
		height += decodeZigZag(hs[i])

		pos[i][0] = u
		pos[i][1] = v
		pos[i][2] = height
	}
	return pos
}
