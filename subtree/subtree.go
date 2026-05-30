package subtree

import (
	"fmt"
	"math"
	"strings"
)

// Subtree represents a 3D Tiles subtree
type Subtree struct {
	SubtreeHeader                    *SubtreeHeader
	SubtreeJson                      string
	SubtreeBinary                    []byte
	ChildSubtreeAvailability         []bool
	ChildSubtreeAvailabilityConstant uint8
	TileAvailability                 []bool
	TileAvailabilityConstant         uint8
	ContentAvailability              []bool
	ContentAvailabilityConstant      uint8
}

// NewSubtree creates a new Subtree instance
func NewSubtree() *Subtree {
	return &Subtree{
		SubtreeHeader:                    NewSubtreeHeader(),
		TileAvailabilityConstant:         0,
		ContentAvailabilityConstant:      0,
		ChildSubtreeAvailabilityConstant: 0,
	}
}

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

// Encode3D encodes 3D coordinates into Morton order
func (mo MortonOrder) Encode3D(x, y, z uint) uint {
	// Interleave bits of x, y, and z
	result := uint(0)
	for i := 0; i < 21; i++ { // 64 bits / 3 ≈ 21 iterations
		if (x & (1 << i)) != 0 {
			result |= 1 << (3 * i)
		}
		if (y & (1 << i)) != 0 {
			result |= 1 << (3*i + 1)
		}
		if (z & (1 << i)) != 0 {
			result |= 1 << (3*i + 2)
		}
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

// Decode3D decodes a Morton index into 3D coordinates
func (mo MortonOrder) Decode3D(mortonIndex uint) (uint, uint, uint) {
	x := uint(0)
	y := uint(0)
	z := uint(0)

	for i := 0; i < 21; i++ { // 64 bits / 3 ≈ 21 iterations
		x |= ((mortonIndex >> (3 * i)) & 1) << i
		y |= ((mortonIndex >> (3*i + 1)) & 1) << i
		z |= ((mortonIndex >> (3*i + 2)) & 1) << i
	}

	return x, y, z
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

// NewAvailabilityLevel creates a new AvailabilityLevel for the given subdivision scheme
func NewAvailabilityLevel(level int, scheme ImplicitSubdivisionScheme) *AvailabilityLevel {
	dim := 1 << uint(level) // 2^level tiles per dimension: sqrt(4^level) for quadtree, cbrt(8^level) for octree
	width := dim
	height := dim

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
	for i := range s {
		s[i] = '0'
	}
	mo := MortonOrder{}

	for x := 0; x < al.width; x++ {
		for y := 0; y < al.height; y++ {
			index := mo.Encode2D(uint(x), uint(y))
			if index < uint(len(s)) {
				if al.BitArray2D.Get(x, y) {
					s[index] = '1'
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

// GetLevel 获取指定级别的 AvailabilityLevel
func (al AvailabilityLevels) GetLevel(level int) *AvailabilityLevel {
	for _, l := range al {
		if l.Level == level {
			return l
		}
	}
	return nil
}

// MaxLevel 获取最大级别数
func (al AvailabilityLevels) MaxLevel() int {
	max := -1
	for _, l := range al {
		if l.Level > max {
			max = l.Level
		}
	}
	return max
}

// GetTileAvailabilityLevels calculates tile availability levels based on content availability levels
func GetTileAvailabilityLevels(contentAvailabilityLevels AvailabilityLevels) AvailabilityLevels {
	tileAvailabilityLevels := make(AvailabilityLevels, 0)

	// 获取最大级别数
	maxLevelNumber := contentAvailabilityLevels.MaxLevel()

	// 为每个级别创建新的 AvailabilityLevel
	for i := 0; i <= maxLevelNumber; i++ {
		tileAvailabilityLevels = append(tileAvailabilityLevels, NewAvailabilityLevel(i, Quadtree))
	}

	// 特殊情况：如果最大级别为0且内容可用性在(0,0)位置为true
	if maxLevelNumber == 0 {
		contentLevel0 := contentAvailabilityLevels.GetLevel(0)
		tileLevel0 := tileAvailabilityLevels.GetLevel(0)
		if contentLevel0 != nil && contentLevel0.BitArray2D != nil &&
			tileLevel0 != nil && tileLevel0.BitArray2D != nil &&
			contentLevel0.BitArray2D.Get(0, 0) {
			tileLevel0.BitArray2D.Set(0, 0, true)
		}
	}

	// 从最高级别向下处理
	for l := maxLevelNumber; l > 0; l-- {
		currentContentLevel := contentAvailabilityLevels.GetLevel(l)
		currentTileLevel := tileAvailabilityLevels.GetLevel(l)
		parentTileLevel := tileAvailabilityLevels.GetLevel(l - 1)

		// 确保所有级别都存在且 BitArray2D 不为 nil
		if currentContentLevel == nil || currentContentLevel.BitArray2D == nil ||
			currentTileLevel == nil || currentTileLevel.BitArray2D == nil ||
			parentTileLevel == nil || parentTileLevel.BitArray2D == nil {
			continue
		}

		w := currentTileLevel.BitArray2D.Width()
		h := currentTileLevel.BitArray2D.Height()

		// 遍历当前级别的所有坐标
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				// 如果内容可用或瓦片已可用
				if currentContentLevel.BitArray2D.Get(x, y) || currentTileLevel.BitArray2D.Get(x, y) {
					// 设置当前瓦片为可用
					currentTileLevel.BitArray2D.Set(x, y, true)

					// 计算父级坐标
					parentX := x >> 1
					parentY := y >> 1

					// 设置父级瓦片为可用
					parentTileLevel.BitArray2D.Set(parentX, parentY, true)
				}
			}
		}
	}

	return tileAvailabilityLevels
}
