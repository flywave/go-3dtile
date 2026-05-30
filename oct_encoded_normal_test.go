package tile3d

import (
	"math"
	"testing"
)

func TestOctEncodedNormalRoundTrip(t *testing.T) {
	normals := [][3]float32{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
		{-1, 0, 0},
		{0, -1, 0},
		{0, 0, -1},
		{0.577, 0.577, 0.577},
		{-0.707, 0.707, 0},
		{0.707, -0.707, 0},
	}
	for _, n := range normals {
		// Normalize first
		mag := float32(math.Sqrt(float64(n[0]*n[0] + n[1]*n[1] + n[2]*n[2])))
		if mag > 0 {
			n[0] /= mag
			n[1] /= mag
			n[2] /= mag
		}
		enc := NewOctEncodedNormal(n)
		dec := enc.Decode()
		dot := n[0]*dec[0] + n[1]*dec[1] + n[2]*dec[2]
		if dot < 0.99 {
			t.Errorf("normal %v: encoded/decoded to %v, dot=%f", n, dec, dot)
		}
	}
}

func TestEncodeXYZEdgeCases(t *testing.T) {
	// Zero vector - should handle gracefully
	val := encodeXYZ(0, 0, 1)
	if val == 0 {
		t.Error("encodeXYZ(0,0,1) should not be 0")
	}
}

func TestClamp(t *testing.T) {
	if clamp(5, 0, 1) != 1 {
		t.Error("clamp(5,0,1) should be 1")
	}
	if clamp(-1, 0, 1) != 0 {
		t.Error("clamp(-1,0,1) should be 0")
	}
	if clamp(0.5, 0, 1) != 0.5 {
		t.Error("clamp(0.5,0,1) should be 0.5")
	}
}

func TestSignNotZero(t *testing.T) {
	if signNotZero(5) != 1 {
		t.Error("signNotZero(5) should be 1")
	}
	if signNotZero(-3) != -1 {
		t.Error("signNotZero(-3) should be -1")
	}
	if signNotZero(0) != 1 {
		t.Error("signNotZero(0) should be 1")
	}
}

func TestNormalizeInPlace(t *testing.T) {
	n := [3]float32{3, 4, 0}
	result := normalizeInPlace(n)
	if result[0] != 0.6 || result[1] != 0.8 || result[2] != 0 {
		t.Errorf("normalizeInPlace(3,4,0) = %v, want [0.6 0.8 0]", result)
	}
}
