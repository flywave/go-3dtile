package tile3d

import (
	"testing"
)

func TestCalcPadding(t *testing.T) {
	tests := []struct {
		offset, unit, expected uint32
	}{
		{0, 8, 0},
		{1, 8, 7},
		{8, 8, 0},
		{9, 8, 7},
		{4, 4, 0},
		{5, 4, 3},
		{0, 1, 0},
	}
	for _, tt := range tests {
		got := calcPadding(tt.offset, tt.unit)
		if got != tt.expected {
			t.Errorf("calcPadding(%d, %d) = %d, want %d", tt.offset, tt.unit, got, tt.expected)
		}
	}
}

func TestCreatePaddingBytes(t *testing.T) {
	r := createPaddingBytes([]byte{1, 2, 3}, 5, 8, 0x20)
	if len(r) != 6 {
		t.Errorf("len = %d, want 6", len(r))
	}
	if r[3] != 0x20 || r[4] != 0x20 || r[5] != 0x20 {
		t.Error("padding bytes not 0x20")
	}
}

func TestZigZag(t *testing.T) {
	tests := []struct {
		in  int
		out uint16
	}{
		{0, 0},
		{-1, 1},
		{1, 2},
		{-2, 3},
		{2, 4},
	}
	for _, tt := range tests {
		encoded := encodeZigZag(tt.in)
		if encoded != tt.out {
			t.Errorf("encodeZigZag(%d) = %d, want %d", tt.in, encoded, tt.out)
		}
		decoded := decodeZigZag(encoded)
		if decoded != tt.in {
			t.Errorf("decodeZigZag(%d) = %d, want %d", encoded, decoded, tt.in)
		}
	}
}

func TestEncodeDecodePolygonPoints(t *testing.T) {
	points := [][2]int{{10, 20}, {30, 40}, {50, 60}}
	us, vs := encodePolygonPoints(points)
	decoded := decodePolygonPoints(us, vs)
	if len(decoded) != len(points) {
		t.Fatalf("len mismatch: %d vs %d", len(decoded), len(points))
	}
	for i := range points {
		if decoded[i] != points[i] {
			t.Errorf("point[%d]: got %v, want %v", i, decoded[i], points[i])
		}
	}
}

func TestEncodeDecodePoints(t *testing.T) {
	points := [][3]int{{10, 20, 30}, {40, 50, 60}, {70, 80, 90}}
	us, vs, hs := encodePoints(points)
	decoded := decodePoints(us, vs, hs)
	if len(decoded) != len(points) {
		t.Fatalf("len mismatch: %d vs %d", len(decoded), len(points))
	}
	for i := range points {
		if decoded[i] != points[i] {
			t.Errorf("point[%d]: got %v, want %v", i, decoded[i], points[i])
		}
	}
}
