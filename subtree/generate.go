package subtree

import (
	"math"
	"strings"
)

// Tile represents a 3D tile with coordinates and availability
type Tile struct {
	Z         int
	X         int
	Y         int
	Available bool
}

// HasChild checks if the given tile is a child of this tile
func (t Tile) HasChild(child Tile) bool {
	// Child must be at a higher level (greater Z)
	if child.Z <= t.Z {
		return false
	}

	// Calculate the level difference
	deltaZ := child.Z - t.Z

	// Calculate the expected parent coordinates of the child
	expectedParentX := child.X >> deltaZ
	expectedParentY := child.Y >> deltaZ

	// Check if the expected parent coordinates match this tile's coordinates
	return expectedParentX == t.X && expectedParentY == t.Y
}

// GetMortonIndices calculates Morton indices for tile and content availability
func GetMortonIndices(tiles []Tile) (tileAvailability string, contentAvailability string) {
	contentAvailabilityLevels := make(AvailabilityLevels, 0)

	// Find maximum Z level
	maxZ := -1
	for _, tile := range tiles {
		if tile.Z > maxZ {
			maxZ = tile.Z
		}
	}

	// Process each level
	for z := 0; z <= maxZ; z++ {
		// Filter tiles for current level that are available
		levelTiles := make([]Tile, 0)
		for _, tile := range tiles {
			if tile.Z == z && tile.Available {
				levelTiles = append(levelTiles, tile)
			}
		}

		// Create availability level
		availabilityLevel := NewAvailabilityLevel(z, Quadtree)

		// Set bits for available tiles
		for _, levelTile := range levelTiles {
			if availabilityLevel.BitArray2D != nil {
				availabilityLevel.BitArray2D.Set(levelTile.X, levelTile.Y, true)
			}
		}

		contentAvailabilityLevels = append(contentAvailabilityLevels, availabilityLevel)
	}

	// Calculate tile availability levels
	tileAvailabilityLevels := GetTileAvailabilityLevels(contentAvailabilityLevels)

	// Convert to Morton index strings
	tileAvailability = tileAvailabilityLevels.ToMortonIndex()
	contentAvailability = contentAvailabilityLevels.ToMortonIndex()

	return tileAvailability, contentAvailability
}

// GetMortonIndexAsBytes calculates Morton indices and returns them as byte arrays
func GetMortonIndexAsBytes(tiles []Tile) (tileAvailability []byte, contentAvailability []byte) {
	// Get Morton indices as strings
	tileAvailabilityStr, contentAvailabilityStr := GetMortonIndices(tiles)

	// Convert strings to byte arrays
	tileAvailability = ToByteArray(FromString(tileAvailabilityStr))
	contentAvailability = ToByteArray(FromString(contentAvailabilityStr))

	return tileAvailability, contentAvailability
}

// GenerateSubtreefile generates a single subtree file from a list of tiles
func GenerateSubtreefile(tiles []Tile) []byte {
	tileAvailability, contentAvailability := GetMortonIndices(tiles)
	subtreeBytes := ToBytesFromStrings(tileAvailability, contentAvailability, "")
	return subtreeBytes
}

// GenerateSubtreefiles generates multiple subtree files from a list of tiles
func GenerateSubtreefiles(tiles []Tile) map[Tile][]byte {
	subtreeFiles := make(map[Tile][]byte)

	// Find maximum level
	maxLevel := -1
	for _, tile := range tiles {
		if tile.Z > maxLevel {
			maxLevel = tile.Z
		}
	}

	// Generate child subtree files at halfway the levels
	// This formula could be adjusted for specific cases
	subtreeLevel := int(math.Ceil(float64(maxLevel+1) / 2))

	if subtreeLevel == 1 {
		subtreeRoot := GenerateSubtreefile(tiles)
		subtreeFiles[Tile{Z: 0, X: 0, Y: 0}] = subtreeRoot
		return subtreeFiles
	}

	tileAvailability, contentAvailability := GetMortonIndices(tiles)

	// Get child subtree availability
	availability := Availability{}
	childSubtreeAvailability := availability.GetLevelAvailability(tileAvailability, subtreeLevel, Quadtree)

	// Calculate offset
	lo := LevelOffset{}
	offset := lo.GetLevelOffset(subtreeLevel, Quadtree)

	// Extract tile and content availability up to offset
	var tileAvailabilityRoot, contentAvailabilityRoot string
	if offset < len(tileAvailability) {
		tileAvailabilityRoot = tileAvailability[:offset]
	} else {
		tileAvailabilityRoot = tileAvailability
	}

	if offset < len(contentAvailability) {
		contentAvailabilityRoot = contentAvailability[:offset]
	} else {
		contentAvailabilityRoot = contentAvailability
	}

	availabilityLength := len(tileAvailabilityRoot)

	// Write the root subtree file
	subtreeRootBytes := ToBytesFromStrings(tileAvailabilityRoot, contentAvailabilityRoot, childSubtreeAvailability)
	subtreeFiles[Tile{Z: 0, X: 0, Y: 0}] = subtreeRootBytes

	// Now create the subtree files
	bac := BitArray2DCreator{}
	ba := bac.GetBitArray2D(childSubtreeAvailability)
	if ba != nil {
		for x := 0; x < ba.Width(); x++ {
			for y := 0; y < ba.Height(); y++ {
				if ba.Get(x, y) {
					t := Tile{Z: subtreeLevel, X: x, Y: y}
					subtreeTiles := GetSubtreeTiles(tiles, t)
					tileAvailabilitySubtree, contentAvailabilitySubtree := GetMortonIndices(subtreeTiles)
					subtreeBytes := ToBytesFromStrings(
						Fill(tileAvailabilitySubtree, availabilityLength),
						Fill(contentAvailabilitySubtree, availabilityLength),
						"")
					subtreeFiles[t] = subtreeBytes
				}
			}
		}
	}

	return subtreeFiles
}

// Fill fills an availability string with zeros to reach the target length
func Fill(availability string, targetLength int) string {
	l := len(availability)
	if l >= targetLength {
		return availability
	}

	// Pad with zeros
	padding := targetLength - l
	zeros := strings.Repeat("0", padding)
	return availability + zeros
}

// GetSubtreeTiles gets the tiles for a specific subtree
func GetSubtreeTiles(tiles []Tile, tile Tile) []Tile {
	res := make([]Tile, 0)

	// Add root tile
	rootTile := Tile{Z: 0, X: 0, Y: 0}

	// Find the subtree tile in the list
	var subtreeTile *Tile
	for i := range tiles {
		if tiles[i].Z == tile.Z && tiles[i].X == tile.X && tiles[i].Y == tile.Y {
			subtreeTile = &tiles[i]
			break
		}
	}

	if subtreeTile != nil {
		rootTile.Available = subtreeTile.Available
	} else {
		rootTile.Available = false
	}

	res = append(res, rootTile)

	// Add children
	for _, child := range tiles {
		if tile.HasChild(child) {
			rel := GetRelativeTile(tile, child)
			rel.Available = child.Available
			res = append(res, rel)
		}
	}

	return res
}

// GetRelativeTile gets the relative position of a tile compared to another tile
func GetRelativeTile(from, to Tile) Tile {
	deltaZ := to.Z - from.Z

	// Calculate the base coordinates for the from tile at the to tile's level
	baseX := from.X
	baseY := from.Y

	for i := 0; i < deltaZ; i++ {
		baseX *= 2
		baseY *= 2
	}

	return Tile{Z: deltaZ, X: to.X - baseX, Y: to.Y - baseY}
}
