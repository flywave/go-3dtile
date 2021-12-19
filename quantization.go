package tile3d

import (
	"math"

	"github.com/flywave/go3d/vec3"
)

const rangeScale16 = 0xffff
const rangeScale8 = 0xff

func computeScale(extent float32, rangeScale uint16) float32 {
	if 0.0 == extent {
		return extent
	}
	return float32(rangeScale) / extent
}

func isInRange(qpos uint16, rangeScale uint16) bool {
	return qpos >= 0 && qpos < rangeScale+1
}

func Quantize(pos float32, origin float32, scale float32, rangeScale uint16) uint16 {
	return uint16(math.Floor(math.Max(0.0, math.Min(float64(rangeScale), float64(pos)*float64(scale)))))
}

func IsQuantizable(pos float32, origin float32, scale float32, rangeScale uint16) bool {
	return isInRange(Quantize(pos, origin, scale, rangeScale), rangeScale16)
}

func UnQuantize(qpos uint16, origin float32, scale float32) float64 {
	if 0.0 == scale {
		return float64(origin)
	}
	return float64(origin) + float64(qpos)/float64(scale)
}

func IsQuantized(qpos uint16) bool {
	return isInRange(qpos, rangeScale16) && qpos == uint16(math.Floor(float64(qpos)))
}

type QParams3d struct {
	Origin [3]float32
	Scale  [3]float32
}

func (p *QParams3d) SetFromRange(range_ *vec3.Box, rangeScale uint16) {
	p.Origin[0] = range_.Min[0]
	p.Origin[1] = range_.Min[1]
	p.Origin[2] = range_.Min[2]

	p.Scale[0] = computeScale(range_.Max[0]-range_.Min[0], rangeScale)
	p.Scale[1] = computeScale(range_.Max[1]-range_.Min[1], rangeScale)
	p.Scale[2] = computeScale(range_.Max[2]-range_.Min[2], rangeScale)
}

func (p *QParams3d) rangeDiagonal() [3]float32 {
	var x float32
	var y float32
	var z float32

	if p.Scale[0] == 0 {
		x = 0
	} else {
		x = rangeScale16 / p.Scale[0]
	}

	if p.Scale[1] == 0 {
		y = 0
	} else {
		y = rangeScale16 / p.Scale[1]
	}

	if p.Scale[2] == 0 {
		z = 0
	} else {
		z = rangeScale16 / p.Scale[2]
	}

	return [3]float32{x, y, z}
}

func QuantizePoint3d(pos [3]float32, params *QParams3d) [3]uint16 {
	var out [3]uint16
	out[0] = Quantize(pos[0], params.Origin[0], params.Scale[0], rangeScale16)
	out[1] = Quantize(pos[1], params.Origin[1], params.Scale[1], rangeScale16)
	out[2] = Quantize(pos[2], params.Origin[2], params.Scale[2], rangeScale16)
	return out
}

func UnQuantizePoint3d(qpos [3]uint16, params *QParams3d) [3]float64 {
	var out [3]float64
	out[0] = UnQuantize(qpos[0], params.Origin[0], params.Scale[0])
	out[1] = UnQuantize(qpos[1], params.Origin[1], params.Scale[1])
	out[2] = UnQuantize(qpos[2], params.Origin[2], params.Scale[2])
	return out
}
