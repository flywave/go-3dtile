package subtree

import (
	"fmt"
	"math"
	"strings"
)

// ImplicitSubdivisionScheme represents the subdivision scheme
type ImplicitSubdivisionScheme int

const (
	// Quadtree subdivision scheme
	Quadtree ImplicitSubdivisionScheme = iota
	// Octree subdivision scheme
	Octree
)

// LevelOffset provides level offset calculations
type LevelOffset struct{}

// GetLevelOffset calculates the level offset based on the subdivision scheme
func (lo LevelOffset) GetLevelOffset(level int, scheme ImplicitSubdivisionScheme) int {
	switch scheme {
	case Quadtree:
		// For quadtree: sum of 4^i for i from 0 to level-1
		offset := 0
		for i := 0; i < level; i++ {
			power := 1
			for j := 0; j < i; j++ {
				power *= 4
			}
			offset += power
		}
		return offset
	case Octree:
		// For octree: sum of 8^i for i from 0 to level-1
		offset := 0
		for i := 0; i < level; i++ {
			power := 1
			for j := 0; j < i; j++ {
				power *= 8
			}
			offset += power
		}
		return offset
	default:
		return 0
	}
}

// MortonOrder provides Morton order encoding functionality
type MortonOrder struct{}

// Encode2D encodes 2D coordinates into Morton order
func (mo MortonOrder) Encode2D(x, y uint) uint {
	// Interleave bits of x and y
	result := uint(0)
	for i := 0; i < 16; i++ {
		result |= ((x & (1 << i)) << i) | ((y & (1 << i)) << (i + 1))
	}
	return result
}

// Decode2D decodes a Morton index into 2D coordinates
func (mo MortonOrder) Decode2D(mortonIndex uint) (uint, uint) {
	x := uint(0)
	y := uint(0)

	for i := 0; i < 16; i++ {
		x |= ((mortonIndex >> (2 * i)) & 1) << i
		y |= ((mortonIndex >> (2*i + 1)) & 1) << i
	}

	return x, y
}

// BitArray2D represents a 2D bit array
type BitArray2D struct {
	data   []bool
	width  int
	height int
}

// NewBitArray2D creates a new BitArray2D with the specified dimensions
func NewBitArray2D(width, height int) (*BitArray2D, error) {
	if width <= 0 {
		return nil, fmt.Errorf("width must be greater than 0, got %d", width)
	}
	if height <= 0 {
		return nil, fmt.Errorf("height must be greater than 0, got %d", height)
	}

	return &BitArray2D{
		data:   make([]bool, width*height),
		width:  width,
		height: height,
	}, nil
}

// Get gets the bit value at the specified coordinates
func (b *BitArray2D) Get(x, y int) bool {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return false
	}
	index := y*b.width + x
	return b.data[index]
}

// Set sets the bit value at the specified coordinates and returns the set value
func (b *BitArray2D) Set(x, y int, value bool) bool {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return false
	}
	index := y*b.width + x
	b.data[index] = value
	return value
}

// Width returns the width of the BitArray2D
func (b *BitArray2D) Width() int {
	return b.width
}

// Height returns the height of the BitArray2D
func (b *BitArray2D) Height() int {
	return b.height
}

// IsAvailable checks if any bits are set to true
func (b *BitArray2D) IsAvailable() bool {
	for _, bit := range b.data {
		if bit {
			return true
		}
	}
	return false
}

// Count counts the number of bits set to the specified value
func (b *BitArray2D) Count(value bool) int {
	count := 0
	for _, bit := range b.data {
		if bit == value {
			count++
		}
	}
	return count
}

// GetAvailableFiles gets the available files based on the bit array
func (b *BitArray2D) GetAvailableFiles(rootZ, rootX, rootY int) []string {
	level := b.getLevel(b.height * b.width)
	newLevel := rootZ + level
	baseX := rootX * 2 * level
	baseY := rootY * 2 * level

	var availableFiles []string

	for x := 0; x < b.width; x++ {
		for y := 0; y < b.height; y++ {
			if b.Get(x, y) {
				file := fmt.Sprintf("%d.%d.%d", newLevel, baseX+x, baseY+y)
				availableFiles = append(availableFiles, file)
			}
		}
	}

	return availableFiles
}

// getLevel calculates the level based on the size
func (b *BitArray2D) getLevel(size int) int {
	if size <= 0 {
		return 0
	}
	// Calculate level as log4(size) = log2(size) / 2
	return int(math.Log2(float64(size)) / 2)
}

// BitArray2DCreator creates BitArray2D instances
type BitArray2DCreator struct{}

// GetBitArray2D creates a BitArray2D from a string representation
func (bac BitArray2DCreator) GetBitArray2D(data string) *BitArray2D {
	if data == "" {
		return nil
	}

	// Create a BitArray2D from string data
	width := len(data)
	height := 1
	ba, _ := NewBitArray2D(width, height)

	for i, char := range data {
		ba.Set(i, 0, char == '1')
	}

	return ba
}

// Availability provides availability-related functionality
type Availability struct{}

// GetLevelAvailability gets the availability for a specific level
func (a Availability) GetLevelAvailability(availability string, level int, scheme ImplicitSubdivisionScheme) string {
	lo := LevelOffset{}
	offset := lo.GetLevelOffset(level, scheme)
	offset1 := lo.GetLevelOffset(level+1, scheme)

	if offset >= len(availability) {
		return ""
	}

	end := offset1
	if end > len(availability) {
		end = len(availability)
	}

	levelAvailability := availability[offset:end]
	return levelAvailability
}

// GetLevel gets the BitArray2D for a specific level
func (a Availability) GetLevel(availability string, level int, scheme ImplicitSubdivisionScheme) *BitArray2D {
	levelAvailability := a.GetLevelAvailability(availability, level, scheme)
	bac := BitArray2DCreator{}
	return bac.GetBitArray2D(levelAvailability)
}

// AvailabilityLevel represents an availability level
type AvailabilityLevel struct {
	Level      int
	BitArray2D *BitArray2D
	width      int
	height     int
}

// NewAvailabilityLevel creates a new AvailabilityLevel
func NewAvailabilityLevel(level int) *AvailabilityLevel {
	width := int(math.Sqrt(math.Pow(4, float64(level))))
	height := int(math.Sqrt(math.Pow(4, float64(level))))

	ba, _ := NewBitArray2D(width, height)

	return &AvailabilityLevel{
		Level:      level,
		BitArray2D: ba,
		width:      width,
		height:     height,
	}
}

// ToMortonIndex converts the availability level to Morton index representation
func (al *AvailabilityLevel) ToMortonIndex() string {
	if al.BitArray2D == nil {
		return ""
	}

	s := make([]byte, al.width*al.height)
	mo := MortonOrder{}

	for x := 0; x < al.width; x++ {
		for y := 0; y < al.height; y++ {
			index := mo.Encode2D(uint(x), uint(y))
			if index < uint(len(s)) {
				if al.BitArray2D.Get(x, y) {
					s[index] = '1'
				} else {
					s[index] = '0'
				}
			}
		}
	}

	return string(s)
}

// AvailabilityLevels represents a list of availability levels
type AvailabilityLevels []*AvailabilityLevel

// ToMortonIndex converts all availability levels to Morton index representation
func (al AvailabilityLevels) ToMortonIndex() string {
	var result strings.Builder
	for _, level := range al {
		result.WriteString(level.ToMortonIndex())
	}
	return result.String()
}
