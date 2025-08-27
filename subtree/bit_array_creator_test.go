package subtree

import (
	"testing"
)

func TestFromString(t *testing.T) {
	// Test with empty string
	empty := FromString("")
	if len(empty) != 0 {
		t.Errorf("Expected empty array, got length %d", len(empty))
	}

	// Test with "1010"
	result := FromString("1010")
	expected := []bool{true, false, true, false}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Index %d: expected %v, got %v", i, v, result[i])
		}
	}

	// Test with "10110000" as in C# test
	ba := FromString("10110000")
	if len(ba) != 8 {
		t.Errorf("Expected length 8, got %d", len(ba))
	}

	if !ba[0] {
		t.Error("Expected bit 0 to be true")
	}

	if ba[1] {
		t.Error("Expected bit 1 to be false")
	}
}

func TestToByteArray(t *testing.T) {
	// Test with empty array
	empty := ToByteArray([]bool{})
	if len(empty) != 0 {
		t.Errorf("Expected empty array, got length %d", len(empty))
	}

	// Test with simple case
	bits := []bool{true, false, true, false}
	result := ToByteArray(bits)
	if len(result) != 1 {
		t.Errorf("Expected 1 byte, got %d", len(result))
	}

	// First byte should have bits 0 and 2 set (1 + 4 = 5)
	if result[0] != 5 {
		t.Errorf("Expected 5, got %d", result[0])
	}
}

func TestAsString(t *testing.T) {
	// Test with empty array
	empty := AsString([]bool{})
	if empty != "" {
		t.Errorf("Expected empty string, got '%s'", empty)
	}

	// Test with mixed values
	bits := []bool{true, false, true, false}
	result := AsString(bits)
	expected := "1010"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestCount(t *testing.T) {
	// Test with empty array
	empty := Count([]bool{}, true)
	if empty != 0 {
		t.Errorf("Expected 0, got %d", empty)
	}

	// Test counting true values
	bits := []bool{true, false, true, false, true}
	count := Count(bits, true)
	if count != 3 {
		t.Errorf("Expected 3 true values, got %d", count)
	}

	// Test counting false values
	count = Count(bits, false)
	if count != 2 {
		t.Errorf("Expected 2 false values, got %d", count)
	}
}
