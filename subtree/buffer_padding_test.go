package subtree

import (
	"testing"
)

func TestAddPadding(t *testing.T) {
	// Test with empty byte array
	empty := []byte{}
	padded := AddPadding(empty, 0)
	if len(padded) != 0 {
		t.Errorf("Expected 0 bytes, got %d", len(padded))
	}
	
	// Test with 5 bytes (should add 3 padding bytes to make it 8)
	data := []byte{1, 2, 3, 4, 5}
	padded = AddPadding(data, 0)
	if len(padded) != 8 {
		t.Errorf("Expected 8 bytes, got %d", len(padded))
	}
	
	// Check that first 5 bytes are unchanged
	for i := 0; i < 5; i++ {
		if padded[i] != data[i] {
			t.Errorf("Data byte %d mismatch: expected %d, got %d", i, data[i], padded[i])
		}
	}
	
	// Check that last 3 bytes are spaces
	for i := 5; i < 8; i++ {
		if padded[i] != ' ' {
			t.Errorf("Padding byte %d should be space, got %d", i, padded[i])
		}
	}
}

func TestAddBinaryPadding(t *testing.T) {
	// Test with 5 bytes (should add 3 padding bytes to make it 8)
	data := []byte{1, 2, 3, 4, 5}
	padded := AddBinaryPadding(data, 0)
	if len(padded) != 8 {
		t.Errorf("Expected 8 bytes, got %d", len(padded))
	}
	
	// Check that first 5 bytes are unchanged
	for i := 0; i < 5; i++ {
		if padded[i] != data[i] {
			t.Errorf("Data byte %d mismatch: expected %d, got %d", i, data[i], padded[i])
		}
	}
	
	// Check that last 3 bytes are zeros
	for i := 5; i < 8; i++ {
		if padded[i] != 0 {
			t.Errorf("Padding byte %d should be zero, got %d", i, padded[i])
		}
	}
}

func TestGetByteArray(t *testing.T) {
	// Test with length 0
	arr := GetByteArray(0)
	if len(arr) != 0 {
		t.Errorf("Expected 0 bytes, got %d", len(arr))
	}
	
	// Test with length 5
	arr = GetByteArray(5)
	if len(arr) != 5 {
		t.Errorf("Expected 5 bytes, got %d", len(arr))
	}
	
	// Check that all bytes are zeros
	for i, b := range arr {
		if b != 0 {
			t.Errorf("Byte %d should be zero, got %d", i, b)
		}
	}
}