package subtree

import (
	"testing"
)

func TestMortonOrder_Encode3D(t *testing.T) {
	mo := MortonOrder{}

	// Test case 1: (0, 0, 0) should encode to 0
	result := mo.Encode3D(0, 0, 0)
	expected := uint(0)
	if result != expected {
		t.Errorf("Encode3D(0, 0, 0) = %d; expected %d", result, expected)
	}

	// Test case 2: (1, 0, 0) should encode to 1
	result = mo.Encode3D(1, 0, 0)
	expected = uint(1)
	if result != expected {
		t.Errorf("Encode3D(1, 0, 0) = %d; expected %d", result, expected)
	}

	// Test case 3: (0, 1, 0) should encode to 2
	result = mo.Encode3D(0, 1, 0)
	expected = uint(2)
	if result != expected {
		t.Errorf("Encode3D(0, 1, 0) = %d; expected %d", result, expected)
	}

	// Test case 4: (0, 0, 1) should encode to 4
	result = mo.Encode3D(0, 0, 1)
	expected = uint(4)
	if result != expected {
		t.Errorf("Encode3D(0, 0, 1) = %d; expected %d", result, expected)
	}

	// Test case 5: (1, 1, 1) should encode to 7
	result = mo.Encode3D(1, 1, 1)
	expected = uint(7)
	if result != expected {
		t.Errorf("Encode3D(1, 1, 1) = %d; expected %d", result, expected)
	}
}

func TestMortonOrder_Decode3D(t *testing.T) {
	mo := MortonOrder{}

	// Test case 1: 0 should decode to (0, 0, 0)
	x, y, z := mo.Decode3D(0)
	if x != 0 || y != 0 || z != 0 {
		t.Errorf("Decode3D(0) = (%d, %d, %d); expected (0, 0, 0)", x, y, z)
	}

	// Test case 2: 1 should decode to (1, 0, 0)
	x, y, z = mo.Decode3D(1)
	if x != 1 || y != 0 || z != 0 {
		t.Errorf("Decode3D(1) = (%d, %d, %d); expected (1, 0, 0)", x, y, z)
	}

	// Test case 3: 2 should decode to (0, 1, 0)
	x, y, z = mo.Decode3D(2)
	if x != 0 || y != 1 || z != 0 {
		t.Errorf("Decode3D(2) = (%d, %d, %d); expected (0, 1, 0)", x, y, z)
	}

	// Test case 4: 4 should decode to (0, 0, 1)
	x, y, z = mo.Decode3D(4)
	if x != 0 || y != 0 || z != 1 {
		t.Errorf("Decode3D(4) = (%d, %d, %d); expected (0, 0, 1)", x, y, z)
	}

	// Test case 5: 7 should decode to (1, 1, 1)
	x, y, z = mo.Decode3D(7)
	if x != 1 || y != 1 || z != 1 {
		t.Errorf("Decode3D(7) = (%d, %d, %d); expected (1, 1, 1)", x, y, z)
	}
}

func TestMortonOrder_EncodeDecode3D(t *testing.T) {
	mo := MortonOrder{}

	// Test that encoding and then decoding returns the original values
	testCases := []struct {
		x, y, z uint
	}{
		{0, 0, 0},
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
		{1, 1, 1},
		{2, 3, 4},
		{7, 5, 3},
	}

	for _, tc := range testCases {
		encoded := mo.Encode3D(tc.x, tc.y, tc.z)
		x, y, z := mo.Decode3D(encoded)
		if x != tc.x || y != tc.y || z != tc.z {
			t.Errorf("Encode3D(%d, %d, %d) = %d; Decode3D(%d) = (%d, %d, %d); expected (%d, %d, %d)",
				tc.x, tc.y, tc.z, encoded, encoded, x, y, z, tc.x, tc.y, tc.z)
		}
	}
}