package subtree

import (
	"testing"
)

func TestHandleBitArray(t *testing.T) {
	// Test with simple bit array
	bits := []bool{true, false, true, false}

	bytes, trueBits, bufferView := HandleBitArray(bits)

	// Check that we got results
	if bytes == nil {
		t.Error("Expected bytes to be non-nil")
	}

	if bufferView == nil {
		t.Error("Expected bufferView to be non-nil")
		return
	}

	// Check true bit count
	if trueBits != 2 {
		t.Errorf("Expected 2 true bits, got %d", trueBits)
	}

	// Check buffer view properties
	if bufferView.Buffer != 0 {
		t.Errorf("Expected buffer index 0, got %d", bufferView.Buffer)
	}

	if bufferView.ByteLength <= 0 {
		t.Error("Expected positive byte length")
	}

	if bufferView.ByteOffset != 0 {
		t.Errorf("Expected offset 0, got %d", bufferView.ByteOffset)
	}
}
