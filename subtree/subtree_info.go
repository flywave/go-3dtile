package subtree

import (
	"fmt"
)

// SubtreeInfo represents the information extracted from a subtree file
type SubtreeInfo struct {
	HeaderMagic              string
	HeaderVersion            uint32
	TileAvailability         string
	ContentAvailability      string
	ChildSubtreeAvailability string
	SubdivisionScheme        ImplicitSubdivisionScheme
}

// GetSubtreeInfo extracts information from a subtree byte array
func GetSubtreeInfo(subtreeBytes []byte, scheme ImplicitSubdivisionScheme) (*SubtreeInfo, error) {
	// Read the subtree
	subtree, err := ReadSubtree(subtreeBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read subtree: %v", err)
	}

	info := &SubtreeInfo{
		HeaderMagic:       subtree.SubtreeHeader.GetMagic(),
		HeaderVersion:     subtree.SubtreeHeader.GetVersion(),
		SubdivisionScheme: scheme,
	}

	// Extract tile availability
	if subtree.TileAvailability != nil {
		info.TileAvailability = AsString(subtree.TileAvailability)
	} else if subtree.TileAvailabilityConstant != 0 {
		// Handle constant tile availability
		info.TileAvailability = fmt.Sprintf("%d", subtree.TileAvailabilityConstant)
	}

	// Extract content availability
	if subtree.ContentAvailability != nil {
		info.ContentAvailability = AsString(subtree.ContentAvailability)
	} else if subtree.ContentAvailabilityConstant != 0 {
		// Handle constant content availability
		info.ContentAvailability = fmt.Sprintf("%d", subtree.ContentAvailabilityConstant)
	}

	// Extract child subtree availability
	if subtree.ChildSubtreeAvailability != nil {
		info.ChildSubtreeAvailability = AsString(subtree.ChildSubtreeAvailability)
	} else if subtree.ChildSubtreeAvailabilityConstant != 0 {
		// Handle constant child subtree availability
		info.ChildSubtreeAvailability = fmt.Sprintf("%d", subtree.ChildSubtreeAvailabilityConstant)
	}

	return info, nil
}

// PrintSubtreeInfo prints detailed information about a subtree
func PrintSubtreeInfo(info *SubtreeInfo) {
	fmt.Println("Subtree info")
	fmt.Println("Action: Info")
	fmt.Printf("Subdivision scheme: %v\n", info.SubdivisionScheme)
	fmt.Printf("Header magic: %s\n", info.HeaderMagic)
	fmt.Printf("Header version: %d\n", info.HeaderVersion)

	fmt.Println()
	fmt.Println("1] Tile availability: ")
	if info.TileAvailability != "" {
		fmt.Printf("TileAvailability: %s\n", info.TileAvailability)
		PrintAvailability(info.TileAvailability, info.SubdivisionScheme)
	}

	fmt.Println("2] Content availability: ")
	if info.ContentAvailability != "" {
		fmt.Printf("ContentAvailability: %s\n", info.ContentAvailability)
		PrintAvailability(info.ContentAvailability, info.SubdivisionScheme)
	}

	fmt.Println()
	fmt.Println("3] Child subtree availability: ")
	if info.ChildSubtreeAvailability != "" {
		fmt.Printf("Availability: %s\n", info.ChildSubtreeAvailability)
		// Note: In the Go implementation, we would need to implement additional functions
		// to match the full functionality of the C# version
	}
}

// PrintAvailability prints detailed availability information
func PrintAvailability(availability string, scheme ImplicitSubdivisionScheme) {
	l := GetNumberOfLevels(availability, scheme)
	fmt.Printf("Number of levels: %d\n", l)

	total := 0
	for i := 0; i < l; i++ {
		ba := Availability{}.GetLevel(availability, i, scheme)
		if ba != nil {
			levelAvailable := ba.Count(true)
			levelTotal := ba.Width() * ba.Height()
			total += levelTotal
			fmt.Printf("Level: %d, available %d/%d\n", i, levelAvailable, levelTotal)
		}
	}
	fmt.Printf("Total: %d\n", total)

	if scheme == Quadtree {
		maxLevel := l
		if l > 4 {
			maxLevel = 4
			fmt.Printf("Printing level [0..%d] of %d...\n", maxLevel-1, l)
		} else {
			fmt.Printf("Printing level [0..%d]...\n", maxLevel-1)
		}
		fmt.Println("")
		for j := 0; j < maxLevel; j++ {
			lo := LevelOffset{}
			offset := lo.GetLevelOffset(j, scheme)
			offset1 := lo.GetLevelOffset(j+1, scheme)
			if offset < len(availability) {
				end := offset1
				if end > len(availability) {
					end = len(availability)
				}
				levelAvailability := availability[offset:end]
				bac := BitArray2DCreator{}
				availabilityArray := bac.GetBitArray2D(levelAvailability)
				if availabilityArray != nil {
					PrintBitArray2D(availabilityArray)
				}
			}
		}
	}
}

// PrintBitArray2D prints a 2D bit array
func PrintBitArray2D(bitArray2D *BitArray2D) {
	if bitArray2D == nil {
		return
	}

	for y := bitArray2D.Height() - 1; y >= 0; y-- {
		for x := 0; x < bitArray2D.Width(); x++ {
			if bitArray2D.Get(x, y) {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
			if x == bitArray2D.Width()-1 {
				fmt.Print(";")
			} else {
				fmt.Print("-")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// GetNumberOfLevels calculates the number of levels in an availability string
func GetNumberOfLevels(availability string, scheme ImplicitSubdivisionScheme) int {
	level := 0
	length := len(availability)
	cont := true

	lo := LevelOffset{}
	for cont {
		offset := lo.GetLevelOffset(level, scheme)
		offsetNext := lo.GetLevelOffset(level+1, scheme)

		if offset < length && offsetNext > length {
			cont = false
		} else {
			level++
		}
	}

	return level
}
