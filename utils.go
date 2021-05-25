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

func zigZag(value uint32) uint16 {
	return uint16((value << 1) ^ (value>>15)&0xFFFF)
}

func zigZagDecode(value uint16) uint32 {
	return uint32(value>>1) ^ uint32(-(value & 1))
}

func encodePolygonPoints(points [][2]uint32) (us, vs []uint16) {
	us = make([]uint16, len(points))
	vs = make([]uint16, len(points))

	lastU := uint32(0)
	lastV := uint32(0)

	for i := 0; i < len(points); i++ {
		u := points[i][0]
		v := points[i][1]

		us[i] = zigZag(u - lastU)
		vs[i] = zigZag(v - lastV)

		lastU = u
		lastV = v
	}
	return
}

func encodePoints(points [][3]uint32) (us, vs, hs []uint16) {
	us = make([]uint16, len(points))
	vs = make([]uint16, len(points))
	hs = make([]uint16, len(points))

	lastU := uint32(0)
	lastV := uint32(0)
	lastH := uint32(0)

	for i := 0; i < len(points); i++ {
		u := points[i][0]
		v := points[i][1]
		h := points[i][2]

		us[i] = zigZag(u - lastU)
		vs[i] = zigZag(v - lastV)
		hs[i] = zigZag(h - lastH)

		lastU = u
		lastV = v
		lastH = h
	}
	return
}

func decodePolygonPoints(us, vs []uint16, pos [][2]uint32) {
	u := uint32(0)
	v := uint32(0)
	pos = make([][2]uint32, len(us))

	for i := 0; i < len(us); i++ {
		u += zigZagDecode(us[i])
		v += zigZagDecode(vs[i])

		pos[i][0] = u
		pos[i][1] = v
	}
}

func decodePoints(us, vs, hs []uint16, pos [][3]uint32) {
	u := uint32(0)
	v := uint32(0)
	height := uint32(0)

	pos = make([][3]uint32, len(us))

	for i := 0; i < len(us); i++ {
		u += zigZagDecode(us[i])
		v += zigZagDecode(vs[i])
		height += zigZagDecode(hs[i])

		pos[i][0] = u
		pos[i][1] = v
		pos[i][2] = height
	}
	return
}
