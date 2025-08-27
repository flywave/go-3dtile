package subtree

import (
	"testing"
)

func TestNewSubtree(t *testing.T) {
	subtree := NewSubtree()

	if subtree == nil {
		t.Error("NewSubtree returned nil")
		return
	}

	if subtree.SubtreeHeader == nil {
		t.Error("SubtreeHeader should not be nil")
	}

	if subtree.TileAvailabilityConstant != 0 {
		t.Errorf("Expected TileAvailabilityConstant to be 0, got %d", subtree.TileAvailabilityConstant)
	}

	if subtree.ContentAvailabilityConstant != 0 {
		t.Errorf("Expected ContentAvailabilityConstant to be 0, got %d", subtree.ContentAvailabilityConstant)
	}
}
