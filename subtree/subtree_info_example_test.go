package subtree_test

import (
	"fmt"
	"log"

	"github.com/flywave/go-3dtile/subtree"
)

// Example of using the subtree info functionality
func ExampleGetSubtreeInfo() {
	// Create a sample subtree
	tileAvailability := "101100000100110010000000"
	contentAvailability := "000000110000000000110000"

	subtreeBytes := subtree.ToBytesFromStrings(tileAvailability, contentAvailability, "")

	// Get subtree info
	info, err := subtree.GetSubtreeInfo(subtreeBytes, subtree.Quadtree)
	if err != nil {
		log.Fatalf("Failed to get subtree info: %v", err)
	}

	// Print basic info
	fmt.Printf("Header magic: %s\n", info.HeaderMagic)
	fmt.Printf("Header version: %d\n", info.HeaderVersion)
	fmt.Printf("Tile availability: %s\n", info.TileAvailability)
	fmt.Printf("Content availability: %s\n", info.ContentAvailability)

	// Output:
	// Header magic: subt
	// Header version: 1
	// Tile availability: 101100000100110010000000
	// Content availability: 000000110000000000110000
}
