package subtree_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/flywave/go-3dtile/subtree"
)

// Example of using the subtree info functionality with generated data
func ExampleGetSubtreeInfo() {
	// Create a sample subtree
	tileAvailability := "101100000100110010000000"
	contentAvailability := "000000110000000000110000"
	childSubtreeAvailability := "0000000000000000011000000000011001100000000001100000000000000000"

	subtreeBytes := subtree.ToBytesFromStrings(tileAvailability, contentAvailability, childSubtreeAvailability)

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
	fmt.Printf("Child subtree availability: %s\n", info.ChildSubtreeAvailability)

	// Output:
	// Header magic: subt
	// Header version: 1
	// Tile availability: 101100000100110010000000
	// Content availability: 000000110000000000110000
	// Child subtree availability: 0000000000000000011000000000011001100000000001100000000000000000
}

// TestRealSubtreeFile tests reading a real subtree file
func TestRealSubtreeFile(t *testing.T) {
	// Check if the file exists
	if _, err := os.Stat("../data/subtree/0.0.0.subtree"); os.IsNotExist(err) {
		fmt.Println("Subtree file does not exist, skipping test")
		return
	}

	// Read the real subtree file
	subtreeBytes, err := os.ReadFile("../data/subtree/0.0.0.subtree")
	if err != nil {
		log.Fatalf("Failed to read subtree file: %v", err)
	}

	fmt.Printf("Read subtree file with %d bytes\n", len(subtreeBytes))

	// Get subtree info
	info, err := subtree.GetSubtreeInfo(subtreeBytes, subtree.Quadtree)
	if err != nil {
		log.Fatalf("Failed to get subtree info: %v", err)
	}

	// Print basic info
	fmt.Printf("Header magic: %s\n", info.HeaderMagic)
	fmt.Printf("Header version: %d\n", info.HeaderVersion)
	fmt.Printf("Tile availability length: %d\n", len(info.TileAvailability))
	fmt.Printf("Content availability length: %d\n", len(info.ContentAvailability))
	fmt.Printf("Child subtree availability length: %d\n", len(info.ChildSubtreeAvailability))

	// Print first 50 characters of each availability string if they exist
	if len(info.TileAvailability) > 0 {
		end := 50
		if len(info.TileAvailability) < 50 {
			end = len(info.TileAvailability)
		}
		fmt.Printf("Tile availability (first %d chars): %s\n", end, info.TileAvailability[:end])
	}

	if len(info.ContentAvailability) > 0 {
		end := 50
		if len(info.ContentAvailability) < 50 {
			end = len(info.ContentAvailability)
		}
		fmt.Printf("Content availability (first %d chars): %s\n", end, info.ContentAvailability[:end])
	}

	if len(info.ChildSubtreeAvailability) > 0 {
		end := 50
		if len(info.ChildSubtreeAvailability) < 50 {
			end = len(info.ChildSubtreeAvailability)
		}
		fmt.Printf("Child subtree availability (first %d chars): %s\n", end, info.ChildSubtreeAvailability[:end])
	} else {
		fmt.Println("Child subtree availability is empty")
	}
}
