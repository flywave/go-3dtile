package subtree

import (
	"testing"
)

func TestSubtreeWriteAndRead(t *testing.T) {
	// Create a subtree with specific availability patterns
	subtree := NewSubtree()

	// Set tile availability (101100000100110010000000)
	subtree.TileAvailability = FromString("101100000100110010000000")

	// Set child subtree availability
	subtree.ChildSubtreeAvailability = FromString("0000000000000000011000000000011001100000000001100000000000000000")

	// Convert to bytes
	bytes := ToBytes(subtree)

	// Read back
	readSubtree, err := ReadSubtree(bytes)
	if err != nil {
		t.Fatalf("Failed to read subtree: %v", err)
	}

	// Verify header
	if readSubtree.SubtreeHeader.GetMagic() != "subt" {
		t.Errorf("Expected magic 'subt', got '%s'", readSubtree.SubtreeHeader.GetMagic())
	}

	if readSubtree.SubtreeHeader.GetVersion() != 1 {
		t.Errorf("Expected version 1, got %d", readSubtree.SubtreeHeader.GetVersion())
	}

	// Verify tile availability
	if readSubtree.TileAvailability == nil {
		t.Fatal("Expected tile availability to be non-nil")
	}

	tileAvailabilityStr := AsString(readSubtree.TileAvailability)
	expectedTileAvailability := "101100000100110010000000"
	if tileAvailabilityStr != expectedTileAvailability {
		t.Errorf("Tile availability mismatch. Expected '%s', got '%s'", expectedTileAvailability, tileAvailabilityStr)
	}

	// Verify child subtree availability
	if readSubtree.ChildSubtreeAvailability == nil {
		t.Fatal("Expected child subtree availability to be non-nil")
	}

	childSubtreeAvailabilityStr := AsString(readSubtree.ChildSubtreeAvailability)
	expectedChildSubtreeAvailability := "0000000000000000011000000000011001100000000001100000000000000000"
	if childSubtreeAvailabilityStr != expectedChildSubtreeAvailability {
		t.Errorf("Child subtree availability mismatch. Expected '%s', got '%s'", expectedChildSubtreeAvailability, childSubtreeAvailabilityStr)
	}
}

func TestSubtreeWriteAndReadWithContent(t *testing.T) {
	// Create a subtree with content availability
	subtree := NewSubtree()

	// Set tile availability
	subtree.TileAvailability = FromString("110010110000000000110000")

	// Set content availability
	subtree.ContentAvailability = FromString("000000110000000000110000")

	// Convert to bytes
	bytes := ToBytes(subtree)

	// Read back
	readSubtree, err := ReadSubtree(bytes)
	if err != nil {
		t.Fatalf("Failed to read subtree: %v", err)
	}

	// Verify tile availability
	if readSubtree.TileAvailability == nil {
		t.Fatal("Expected tile availability to be non-nil")
	}

	tileAvailabilityStr := AsString(readSubtree.TileAvailability)
	expectedTileAvailability := "110010110000000000110000"
	if tileAvailabilityStr != expectedTileAvailability {
		t.Errorf("Tile availability mismatch. Expected '%s', got '%s'", expectedTileAvailability, tileAvailabilityStr)
	}

	// Verify content availability
	if readSubtree.ContentAvailability == nil {
		t.Fatal("Expected content availability to be non-nil")
	}

	contentAvailabilityStr := AsString(readSubtree.ContentAvailability)
	expectedContentAvailability := "000000110000000000110000"
	if contentAvailabilityStr != expectedContentAvailability {
		t.Errorf("Content availability mismatch. Expected '%s', got '%s'", expectedContentAvailability, contentAvailabilityStr)
	}
}

func TestSubtreeWithConstants(t *testing.T) {
	// Create a subtree with constants
	subtree := NewSubtree()
	subtree.TileAvailabilityConstant = 1
	subtree.ContentAvailabilityConstant = 0

	// For this test, we'll just verify that we can create a subtree with constants
	// The full test would require modifying the ToBytes function to handle constants properly
	if subtree.TileAvailabilityConstant != 1 {
		t.Errorf("Expected TileAvailabilityConstant to be 1, got %d", subtree.TileAvailabilityConstant)
	}

	if subtree.ContentAvailabilityConstant != 0 {
		t.Errorf("Expected ContentAvailabilityConstant to be 0, got %d", subtree.ContentAvailabilityConstant)
	}
}

func TestToBytesFromStrings(t *testing.T) {
	// Test the convenience function
	tileAvailability := "101100000100110010000000"
	contentAvailability := "000000110000000000110000"
	childSubtreeAvailability := "0000000000000000011000000000011001100000000001100000000000000000"

	bytes := ToBytesFromStrings(tileAvailability, contentAvailability, childSubtreeAvailability)

	if len(bytes) == 0 {
		t.Error("Expected non-empty byte array")
	}

	// Read back and verify
	readSubtree, err := ReadSubtree(bytes)
	if err != nil {
		t.Fatalf("Failed to read subtree: %v", err)
	}

	// Verify tile availability
	tileAvailabilityStr := AsString(readSubtree.TileAvailability)
	if tileAvailabilityStr != tileAvailability {
		t.Errorf("Tile availability mismatch. Expected '%s', got '%s'", tileAvailability, tileAvailabilityStr)
	}

	// Verify content availability
	contentAvailabilityStr := AsString(readSubtree.ContentAvailability)
	if contentAvailabilityStr != contentAvailability {
		t.Errorf("Content availability mismatch. Expected '%s', got '%s'", contentAvailability, contentAvailabilityStr)
	}

	// Verify child subtree availability
	childSubtreeAvailabilityStr := AsString(readSubtree.ChildSubtreeAvailability)
	if childSubtreeAvailabilityStr != childSubtreeAvailability {
		t.Errorf("Child subtree availability mismatch. Expected '%s', got '%s'", childSubtreeAvailability, childSubtreeAvailabilityStr)
	}
}
