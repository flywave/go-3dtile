package tile3d

import "math"

func clamp(val float32, minVal float32, maxVal float32) float32 {
	if val < minVal {
		return minVal
	}

	if val > maxVal {
		return maxVal
	}

	return val
}

func clampUint8(val float32) uint16 {
	return roundUint16(0.5 + float64(clamp(val, -1, 1)*0.5+0.5)*255.0)
}

func roundUint16(val float64) uint16 {
	return uint16(val)
}

func signNotZero(val float32) float32 {
	if val < 0.0 {
		return -1.0
	}
	return 1.0
}

type OctEncodedNormal struct {
	value uint16
}

func NewOctEncodedNormal(v [3]float32) *OctEncodedNormal {
	return &OctEncodedNormal{value: encodeXYZ(v[0], v[1], v[2])}
}

func (n *OctEncodedNormal) Decode() [3]float32 {
	return decodeValue(n.value)
}

func encodeXYZ(nx float32, ny float32, nz float32) uint16 {
	denom := float32(math.Abs(float64(nx)) + math.Abs(float64(ny)) + math.Abs(float64(nz)))
	rx := nx / denom
	ry := ny / denom
	if nz < 0 {
		x := rx
		y := ry
		rx = (1 - float32(math.Abs(float64(y)))) * signNotZero(x)
		ry = (1 - float32(math.Abs(float64(x)))) * signNotZero(y)
	}
	return clampUint8(rx) | clampUint8(ry)<<8
}

func decodeValue(val uint16) [3]float32 {
	ex := float32(val & 0xff)
	ey := float32(val >> 8)
	ex = ex/255.0*2.0 - 1.0
	ey = ey/255.0*2.0 - 1.0
	ez := 1 - float32(math.Abs(float64(ex))+math.Abs(float64(ey)))
	var n [3]float32

	n[0] = ex
	n[1] = ey
	n[2] = ez

	if n[2] < 0 {
		x := n[0]
		y := n[1]

		n[0] = (1 - float32(math.Abs(float64(y)))) * signNotZero(x)
		n[1] = (1 - float32(math.Abs(float64(x)))) * signNotZero(y)
	}

	n = normalizeInPlace(n)

	return n
}

const smallMetricDistance = 1.0e-6

func inverseMetricDistance(a float64) *float64 {
	if math.Abs(a) <= smallMetricDistance {
		return nil
	}
	dist := 1.0 / a
	return &dist
}

func normalizeInPlace(n [3]float32) [3]float32 {
	magnitude := math.Sqrt(float64(n[0]*n[0] + n[1]*n[1] + n[2]*n[2]))
	a := inverseMetricDistance(magnitude)
	if a == nil {
		return n
	}
	n[0] *= float32(*a)
	n[1] *= float32(*a)
	n[2] *= float32(*a)
	return n
}
