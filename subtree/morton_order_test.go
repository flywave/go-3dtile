package subtree

import (
	"testing"
)

func TestMortonOrderEncode2D(t *testing.T) {
	// Test basic cases
	tests := []struct {
		x, y     uint
		expected uint
	}{
		{0, 0, 0}, // binary: 000
		{1, 0, 1}, // binary: 001
		{0, 1, 2}, // binary: 010
		{1, 1, 3}, // binary: 011
		{2, 0, 4}, // binary: 100
		{3, 0, 5}, // binary: 101
		{2, 1, 6}, // binary: 110
		{3, 1, 7}, // binary: 111
		{7, 5, 55}, // binary: 110111
	}

	for _, test := range tests {
		result := MortonOrder{}.Encode2D(test.x, test.y)
		if result != test.expected {
			t.Errorf("Encode2D(%d, %d) = %d, expected %d", test.x, test.y, result, test.expected)
		}
	}
}

func TestMortonOrderRoundTrip(t *testing.T) {
	// Test round trip conversion
	// Sample from: https://github.com/CesiumGS/3d-tiles/blob/draft-1.1/specification/ImplicitTiling/AVAILABILITY.adoc#implicittiling-availability-indexing
	mortonIndex := uint(0b010011)
	mo := MortonOrder{}
	// Decode the morton index
	x, y := mo.Decode2D(mortonIndex)
	
	// Check if we get the expected values
	if x != 5 || y != 1 {
		t.Errorf("Decode2D(0b010011) = (%d, %d), expected (5, 1)", x, y)
	}
	
	// Re-encode and check if we get the same value
	reencoded := mo.Encode2D(x, y)
	if reencoded != mortonIndex {
		t.Errorf("Encode2D(%d, %d) = %d, expected %d", x, y, reencoded, mortonIndex)
	}
}

func TestMortonOrderDecode2D(t *testing.T) {
	// Test decoding cases
	tests := []struct {
		mortonIndex uint
		expectedX   uint
		expectedY   uint
	}{
		{0, 0, 0},  // binary: 000
		{1, 1, 0},  // binary: 001
		{2, 0, 1},  // binary: 010
		{3, 1, 1},  // binary: 011
		{4, 2, 0},  // binary: 100
		{5, 3, 0},  // binary: 101
		{6, 2, 1},  // binary: 110
		{7, 3, 1},  // binary: 111
		{55, 7, 5}, // binary: 110111
	}

	mo := MortonOrder{}
	for _, test := range tests {
		x, y := mo.Decode2D(test.mortonIndex)
		if x != test.expectedX || y != test.expectedY {
			t.Errorf("Decode2D(%d) = (%d, %d), expected (%d, %d)", test.mortonIndex, x, y, test.expectedX, test.expectedY)
		}
	}
}