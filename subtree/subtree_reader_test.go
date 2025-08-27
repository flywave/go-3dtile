package subtree

import (
	"testing"
)

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
