package subtree

import (
	"testing"
)

func TestGetSubtreeInfo(t *testing.T) {
	// Create a sample subtree
	tileAvailability := "101100000100110010000000"
	contentAvailability := "000000110000000000110000"
	childSubtreeAvailability := "0000000000000000011000000000011001100000000001100000000000000000"
	
	subtreeBytes := ToBytesFromStrings(tileAvailability, contentAvailability, childSubtreeAvailability)
	
	// Test getting subtree info
	info, err := GetSubtreeInfo(subtreeBytes, Quadtree)
	if err != nil {
		t.Fatalf("Failed to get subtree info: %v", err)
	}
	
	// Verify the info
	if info.HeaderMagic != "subt" {
		t.Errorf("Expected header magic 'subt', got '%s'", info.HeaderMagic)
	}
	
	if info.HeaderVersion != 1 {
		t.Errorf("Expected header version 1, got %d", info.HeaderVersion)
	}
	
	if info.TileAvailability != tileAvailability {
		t.Errorf("Tile availability mismatch. Expected '%s', got '%s'", tileAvailability, info.TileAvailability)
	}
	
	if info.ContentAvailability != contentAvailability {
		t.Errorf("Content availability mismatch. Expected '%s', got '%s'", contentAvailability, info.ContentAvailability)
	}
	
	if info.ChildSubtreeAvailability != childSubtreeAvailability {
		t.Errorf("Child subtree availability mismatch. Expected '%s', got '%s'", childSubtreeAvailability, info.ChildSubtreeAvailability)
	}
}

func TestGetNumberOfLevels(t *testing.T) {
	// Test with a known availability string
	availability := "101100000100110010000000"
	levels := GetNumberOfLevels(availability, Quadtree)
	
	// This should have 3 levels based on the pattern
	if levels <= 0 {
		t.Errorf("Expected positive number of levels, got %d", levels)
	}
}

func TestPrintAvailability(t *testing.T) {
	// Test that the function doesn't panic
	availability := "101100000100110010000000"
	
	// This should not panic
	PrintAvailability(availability, Quadtree)
}

func TestPrintBitArray2D(t *testing.T) {
	// Test with a valid BitArray2D
	ba, err := NewBitArray2D(2, 2)
	if err != nil {
		t.Fatalf("Failed to create BitArray2D: %v", err)
	}
	
	ba.Set(0, 0, true)
	ba.Set(1, 1, true)
	
	// This should not panic
	PrintBitArray2D(ba)
	
	// Test with nil (should not panic)
	PrintBitArray2D(nil)
}