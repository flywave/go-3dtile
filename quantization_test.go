package tile3d

import (
	"testing"

	"github.com/flywave/go3d/vec3"
)

func TestQuantizeUnQuantize(t *testing.T) {
	origin := float32(0)
	scale := float32(100)
	q := Quantize(50.0, origin, scale, rangeScale16)
	if q != 32767 {
		t.Errorf("Quantize(50, 0, 100) = %d, want 32767", q)
	}
	uq := UnQuantize(q, origin, scale, rangeScale16)
	if uq < 49.9 || uq > 50.1 {
		t.Errorf("UnQuantize(%d) = %f, want ~50", q, uq)
	}
}

func TestQuantizeEdgeCases(t *testing.T) {
	origin := float32(10)
	scale := float32(20)

	// Below origin
	q := Quantize(5.0, origin, scale, rangeScale16)
	if q != 0 {
		t.Errorf("Quantize(5, 10, 20) = %d, want 0", q)
	}

	// Above range
	q = Quantize(100.0, origin, scale, rangeScale16)
	if q != rangeScale16 {
		t.Errorf("Quantize(100, 10, 20) = %d, want %d", q, rangeScale16)
	}
}

func TestIsQuantizable(t *testing.T) {
	if !IsQuantizable(50, 0, 100, rangeScale16) {
		t.Error("IsQuantizable(50,0,100) should be true")
	}
}

func TestQParams3d(t *testing.T) {
	p := &QParams3d{}
	box := vec3.Box{Min: [3]float32{0, 0, 0}, Max: [3]float32{100, 200, 300}}
	p.SetFromRange(&box, rangeScale16)

	pos := [3]float64{50, 100, 150}
	qpos := QuantizePoint3d(pos, p)
	decoded := UnQuantizePoint3d(qpos, p)

	for i := 0; i < 3; i++ {
		if decoded[i] < pos[i]-1 || decoded[i] > pos[i]+1 {
			t.Errorf("component %d: original=%f, decoded=%f", i, pos[i], decoded[i])
		}
	}
}

func TestIsQuantized(t *testing.T) {
	if !IsQuantized(uint16(100)) {
		t.Error("IsQuantized(100) should be true")
	}
}

func TestRangeDiagonal(t *testing.T) {
	p := &QParams3d{
		Origin: [3]float32{0, 0, 0},
		Scale:  [3]float32{100, 200, 0},
	}
	diag := p.rangeDiagonal()
	if diag[0] <= 0 || diag[1] <= 0 {
		t.Error("rangeDiagonal components should be > 0 for non-zero scale")
	}
	if diag[2] != 0 {
		t.Error("rangeDiagonal component for zero scale should be 0")
	}
}

func TestComputeScale(t *testing.T) {
	if computeScale(0, 0xffff) != 1 {
		t.Error("computeScale(0) should return 1")
	}
	if computeScale(100, 0xffff) != 100 {
		t.Errorf("computeScale(100) = %f, want 100", computeScale(100, 0xffff))
	}
}
