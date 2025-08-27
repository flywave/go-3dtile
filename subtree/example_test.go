package subtree_test

import (
	"fmt"
	"log"

	"github.com/flywave/go-3dtile/subtree"
)

// Example of creating and reading a subtree
func ExampleSubtree() {
	// Create a subtree from string representations
	tileAvailability := "101100000100110010000000"
	contentAvailability := "000000110000000000110000"
	childSubtreeAvailability := "0000000000000000011000000000011001100000000001100000000000000000"

	// Convert to bytes
	subtreeBytes := subtree.ToBytesFromStrings(tileAvailability, contentAvailability, childSubtreeAvailability)

	fmt.Printf("Generated subtree: %d bytes\n", len(subtreeBytes))

	// Read the subtree back
	readSubtree, err := subtree.ReadSubtree(subtreeBytes)
	if err != nil {
		log.Fatalf("Error reading subtree: %v", err)
	}

	// Print information about the subtree
	fmt.Printf("Subtree header magic: %s\n", readSubtree.SubtreeHeader.GetMagic())
	fmt.Printf("Subtree header version: %d\n", readSubtree.SubtreeHeader.GetVersion())
	fmt.Printf("Tile availability: %s\n", subtree.AsString(readSubtree.TileAvailability))
	fmt.Printf("Content availability: %s\n", subtree.AsString(readSubtree.ContentAvailability))
	fmt.Printf("Child subtree availability: %s\n", subtree.AsString(readSubtree.ChildSubtreeAvailability))

	// Output:
	// Generated subtree: 400 bytes
	// Subtree header magic: subt
	// Subtree header version: 1
	// Tile availability: 101100000100110010000000
	// Content availability: 000000110000000000110000
	// Child subtree availability: 0000000000000000011000000000011001100000000001100000000000000000
}
