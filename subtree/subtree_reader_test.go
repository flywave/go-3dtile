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
func TestReadSubtree(t *testing.T) {
	// Create a simple subtree with string representations
	tileAvailability := "1100"
	contentAvailability := "1010"

	// Generate subtree bytes
	subtreeBytes := ToBytesFromStrings(tileAvailability, contentAvailability, "")

	// Try to read it back
	subtree, err := ReadSubtree(subtreeBytes)

	// Check that we got a subtree without error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if subtree == nil {
		t.Error("Expected subtree to be non-nil")
		return
	}

	// Check that the header was read
	if subtree.SubtreeHeader == nil {
		t.Error("Expected SubtreeHeader to be non-nil")
	}

	// Check that JSON was read
	if subtree.SubtreeJson == "" {
		t.Error("Expected SubtreeJson to be non-empty")
	}

	// Check that binary data was read
	if len(subtree.SubtreeBinary) == 0 {
		t.Error("Expected SubtreeBinary to be non-empty")
	}
}
