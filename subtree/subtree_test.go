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

func TestLevelOffsetGetLevelOffset(t *testing.T) {
	lo := LevelOffset{}

	// Test Quadtree
	if offset := lo.GetLevelOffset(0, Quadtree); offset != 0 {
		t.Errorf("Expected offset 0 for level 0 quadtree, got %d", offset)
	}

	if offset := lo.GetLevelOffset(1, Quadtree); offset != 1 {
		t.Errorf("Expected offset 1 for level 1 quadtree, got %d", offset)
	}

	if offset := lo.GetLevelOffset(2, Quadtree); offset != 5 {
		t.Errorf("Expected offset 5 for level 2 quadtree, got %d", offset)
	}

	// Test Octree
	if offset := lo.GetLevelOffset(0, Octree); offset != 0 {
		t.Errorf("Expected offset 0 for level 0 octree, got %d", offset)
	}

	if offset := lo.GetLevelOffset(1, Octree); offset != 1 {
		t.Errorf("Expected offset 1 for level 1 octree, got %d", offset)
	}

	if offset := lo.GetLevelOffset(2, Octree); offset != 9 {
		t.Errorf("Expected offset 9 for level 2 octree, got %d", offset)
	}
}

func TestBitArray2D(t *testing.T) {
	// Test NewBitArray2D with valid dimensions
	ba, err := NewBitArray2D(2, 2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if ba == nil {
		t.Fatal("Expected non-nil BitArray2D")
	}

	if ba.width != 2 || ba.height != 2 {
		t.Errorf("Expected 2x2 dimensions, got %dx%d", ba.width, ba.height)
	}

	// Test NewBitArray2D with invalid dimensions
	_, err = NewBitArray2D(0, 2)
	if err == nil {
		t.Error("Expected error for width <= 0")
	}

	_, err = NewBitArray2D(2, 0)
	if err == nil {
		t.Error("Expected error for height <= 0")
	}

	// Test Set and Get
	result := ba.Set(0, 0, true)
	if !result {
		t.Error("Expected Set to return true")
	}

	if !ba.Get(0, 0) {
		t.Error("Expected (0,0) to be true")
	}

	result = ba.Set(1, 1, true)
	if !result {
		t.Error("Expected Set to return true")
	}

	if !ba.Get(1, 1) {
		t.Error("Expected (1,1) to be true")
	}

	// Test default value
	if ba.Get(0, 1) {
		t.Error("Expected (0,1) to be false by default")
	}

	// Test out of bounds
	if ba.Get(-1, 0) {
		t.Error("Expected out of bounds access to return false")
	}

	if ba.Get(2, 0) {
		t.Error("Expected out of bounds access to return false")
	}

	// Test Set out of bounds
	if ba.Set(-1, 0, true) {
		t.Error("Expected out of bounds Set to return false")
	}

	if ba.Set(2, 0, true) {
		t.Error("Expected out of bounds Set to return false")
	}
}

func TestBitArray2DWidthHeight(t *testing.T) {
	ba, _ := NewBitArray2D(3, 4)

	if width := ba.Width(); width != 3 {
		t.Errorf("Expected width 3, got %d", width)
	}

	if height := ba.Height(); height != 4 {
		t.Errorf("Expected height 4, got %d", height)
	}
}

func TestBitArray2DIsAvailable(t *testing.T) {
	ba, _ := NewBitArray2D(2, 2)

	// Test with all false values
	if ba.IsAvailable() {
		t.Error("Expected IsAvailable to return false for all false values")
	}

	// Test with at least one true value
	ba.Set(0, 0, true)
	if !ba.IsAvailable() {
		t.Error("Expected IsAvailable to return true when at least one value is true")
	}
}

func TestBitArray2DCount(t *testing.T) {
	ba, _ := NewBitArray2D(2, 2)

	// Test count of false values (all default to false)
	if count := ba.Count(false); count != 4 {
		t.Errorf("Expected count of false values to be 4, got %d", count)
	}

	// Test count of true values (all default to false)
	if count := ba.Count(true); count != 0 {
		t.Errorf("Expected count of true values to be 0, got %d", count)
	}

	// Set some values to true
	ba.Set(0, 0, true)
	ba.Set(1, 1, true)

	// Test count of true values
	if count := ba.Count(true); count != 2 {
		t.Errorf("Expected count of true values to be 2, got %d", count)
	}

	// Test count of false values
	if count := ba.Count(false); count != 2 {
		t.Errorf("Expected count of false values to be 2, got %d", count)
	}
}

func TestBitArray2DGetAvailableFiles(t *testing.T) {
	ba, _ := NewBitArray2D(2, 2)

	// Set some values to true
	ba.Set(0, 0, true)
	ba.Set(1, 1, true)

	// Test GetAvailableFiles
	files := ba.GetAvailableFiles(0, 0, 0)
	if len(files) != 2 {
		t.Errorf("Expected 2 available files, got %d", len(files))
	}

	// Check that we have the right number of files
	// For a 2x2 grid with 2 true values, we should have 2 files
	// The exact file names depend on the level calculation
	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d: %v", len(files), files)
	}
}

func TestBitArray2DGetLevel(t *testing.T) {
	ba, _ := NewBitArray2D(4, 4)

	// For a 4x4 grid (16 elements), level should be 2 (since 4^2 = 16)
	// But our calculation is log4(16) = log2(16) / 2 = 4 / 2 = 2
	level := ba.getLevel(16)
	if level != 2 {
		t.Errorf("Expected level 2 for 16 elements, got %d", level)
	}

	// For a 2x2 grid (4 elements), level should be 1
	level = ba.getLevel(4)
	if level != 1 {
		t.Errorf("Expected level 1 for 4 elements, got %d", level)
	}
}

func TestBitArray2DCreatorGetBitArray2D(t *testing.T) {
	bac := BitArray2DCreator{}

	// Test nil case
	if ba := bac.GetBitArray2D(""); ba != nil {
		t.Error("Expected nil for empty string")
	}

	// Test non-empty case
	data := "1010"
	if ba := bac.GetBitArray2D(data); ba == nil {
		t.Error("Expected non-nil BitArray2D for non-empty string")
	} else {
		if !ba.Get(0, 0) {
			t.Error("Expected first bit to be true")
		}

		if ba.Get(1, 0) {
			t.Error("Expected second bit to be false")
		}

		if !ba.Get(2, 0) {
			t.Error("Expected third bit to be true")
		}

		if ba.Get(3, 0) {
			t.Error("Expected fourth bit to be false")
		}
	}
}

func TestAvailabilityGetLevelAvailability(t *testing.T) {
	a := Availability{}
	availability := "101010101010" // Sample availability string

	// Test with Quadtree scheme
	lo := LevelOffset{}
	offset0 := lo.GetLevelOffset(0, Quadtree) // 0
	offset1 := lo.GetLevelOffset(1, Quadtree) // 1
	offset2 := lo.GetLevelOffset(2, Quadtree) // 5

	// Level 0: from offset 0 to 1
	expectedLevel0 := availability[offset0:offset1]
	if result := a.GetLevelAvailability(availability, 0, Quadtree); result != expectedLevel0 {
		t.Errorf("Expected level 0 availability %s, got %s", expectedLevel0, result)
	}

	// Level 1: from offset 1 to 5
	expectedLevel1 := availability[offset1:offset2]
	if result := a.GetLevelAvailability(availability, 1, Quadtree); result != expectedLevel1 {
		t.Errorf("Expected level 1 availability %s, got %s", expectedLevel1, result)
	}
}

func TestAvailabilityGetLevel(t *testing.T) {
	a := Availability{}
	availability := "101010101010" // Sample availability string

	// Test with Quadtree scheme
	if ba := a.GetLevel(availability, 0, Quadtree); ba == nil {
		t.Error("Expected non-nil BitArray2D for level 0")
	}

	// Test with Octree scheme
	if ba := a.GetLevel(availability, 0, Octree); ba == nil {
		t.Error("Expected non-nil BitArray2D for level 0 with octree scheme")
	}
}

func TestAvailabilityLevel(t *testing.T) {
	// Test NewAvailabilityLevel
	al := NewAvailabilityLevel(1)
	if al == nil {
		t.Fatal("Expected non-nil AvailabilityLevel")
	}

	if al.Level != 1 {
		t.Errorf("Expected level 1, got %d", al.Level)
	}

	// Test width and height calculation
	// For level 1: sqrt(4^1) = sqrt(4) = 2
	if al.width != 2 || al.height != 2 {
		t.Errorf("Expected 2x2 dimensions for level 1, got %dx%d", al.width, al.height)
	}

	// Test ToMortonIndex with empty BitArray2D
	result := al.ToMortonIndex()
	if len(result) != 4 {
		t.Errorf("Expected result length 4, got %d: %s", len(result), result)
	}
}

func TestAvailabilityLevelsToMortonIndex(t *testing.T) {
	// Create test levels
	level0 := NewAvailabilityLevel(0)
	level1 := NewAvailabilityLevel(1)

	// Set some values
	level0.BitArray2D.Set(0, 0, true)
	level1.BitArray2D.Set(0, 0, true)
	level1.BitArray2D.Set(1, 1, true)

	levels := AvailabilityLevels{level0, level1}
	result := levels.ToMortonIndex()

	// Level 0: 1 bit -> "1"
	// Level 1: 4 bits -> "1001" (based on Morton order)
	// Expected: "1" + "1001" = "11001"
	if len(result) != 5 {
		t.Errorf("Expected result length 5, got %d: %s", len(result), result)
	}
}

func TestMortonOrderEncode2D(t *testing.T) {
	// Test basic cases
	tests := []struct {
		x, y     uint
		expected uint
	}{
		{0, 0, 0},  // binary: 000
		{1, 0, 1},  // binary: 001
		{0, 1, 2},  // binary: 010
		{1, 1, 3},  // binary: 011
		{2, 0, 4},  // binary: 100
		{3, 0, 5},  // binary: 101
		{2, 1, 6},  // binary: 110
		{3, 1, 7},  // binary: 111
		{7, 5, 55}, // binary: 110111
	}

	for _, test := range tests {
		result := MortonOrder{}.Encode2D(test.x, test.y)
		if result != test.expected {
			t.Errorf("Encode2D(%d, %d) = %d, expected %d", test.x, test.y, result, test.expected)
		}
	}
}

func TestMortonOrderRoundTrip(t *testing.T) {
	// Test round trip conversion
	// Sample from: https://github.com/CesiumGS/3d-tiles/blob/draft-1.1/specification/ImplicitTiling/AVAILABILITY.adoc#implicittiling-availability-indexing
	mortonIndex := uint(0b010011)
	mo := MortonOrder{}
	// Decode the morton index
	x, y := mo.Decode2D(mortonIndex)

	// Check if we get the expected values
	if x != 5 || y != 1 {
		t.Errorf("Decode2D(0b010011) = (%d, %d), expected (5, 1)", x, y)
	}

	// Re-encode and check if we get the same value
	reencoded := mo.Encode2D(x, y)
	if reencoded != mortonIndex {
		t.Errorf("Encode2D(%d, %d) = %d, expected %d", x, y, reencoded, mortonIndex)
	}
}

func TestMortonOrderDecode2D(t *testing.T) {
	// Test decoding cases
	tests := []struct {
		mortonIndex uint
		expectedX   uint
		expectedY   uint
	}{
		{0, 0, 0},  // binary: 000
		{1, 1, 0},  // binary: 001
		{2, 0, 1},  // binary: 010
		{3, 1, 1},  // binary: 011
		{4, 2, 0},  // binary: 100
		{5, 3, 0},  // binary: 101
		{6, 2, 1},  // binary: 110
		{7, 3, 1},  // binary: 111
		{55, 7, 5}, // binary: 110111
	}

	mo := MortonOrder{}
	for _, test := range tests {
		x, y := mo.Decode2D(test.mortonIndex)
		if x != test.expectedX || y != test.expectedY {
			t.Errorf("Decode2D(%d) = (%d, %d), expected (%d, %d)", test.mortonIndex, x, y, test.expectedX, test.expectedY)
		}
	}
}
