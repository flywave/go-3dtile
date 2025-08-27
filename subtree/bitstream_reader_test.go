package subtree

import (
	"testing"
)

func TestRead(t *testing.T) {
	// Test with empty byte array
	empty := Read([]byte{}, 0, 0)
	if len(empty) != 0 {
		t.Errorf("Expected empty array, got length %d", len(empty))
	}
	
	// Test with simple case
	// Byte with value 5 has bits 0 and 2 set
	data := []byte{5}
	result := Read(data, 0, 1)
	
	// Should have 8 bits
	if len(result) != 8 {
		t.Errorf("Expected 8 bits, got %d", len(result))
	}
	
	// Check that bits 0 and 2 are true, others are false
	expected := []bool{true, false, true, false, false, false, false, false}
	for i, exp := range expected {
		if result[i] != exp {
			t.Errorf("Bit %d: expected %v, got %v", i, exp, result[i])
		}
	}
	
	// Test with offset
	result = Read(data, 1, 1)
	// Should be empty since offset is beyond data
	if len(result) != 0 {
		t.Errorf("Expected empty array with invalid offset, got length %d", len(result))
	}
}