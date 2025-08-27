# Subtree Package

This package provides functionality for working with 3D Tiles Subtree files, which are used in implicit tiling schemes.

## Features

- Read and write 3D Tiles Subtree binary format
- Convert between string representations and binary data
- Handle tile availability, content availability, and child subtree availability
- Support for both quadtree and octree subdivision schemes
- Morton order encoding and decoding
- Bit array operations
- Subtree information extraction and printing

## Installation

```bash
go get github.com/flywave/go-3dtile/subtree
```

## Usage

### Creating a Subtree from String Representations

```go
package main

import (
    "fmt"
    "github.com/flywave/go-3dtile/subtree"
)

func main() {
    // Create a subtree from string representations
    tileAvailability := "101100000100110010000000"
    contentAvailability := "000000110000000000110000"
    childSubtreeAvailability := "0000000000000000011000000000011001100000000001100000000000000000"
    
    // Convert to bytes
    subtreeBytes := subtree.ToBytesFromStrings(tileAvailability, contentAvailability, childSubtreeAvailability)
    
    fmt.Printf("Generated subtree: %d bytes\n", len(subtreeBytes))
}
```

### Reading a Subtree

```go
// Read the subtree back
readSubtree, err := subtree.ReadSubtree(subtreeBytes)
if err != nil {
    log.Fatalf("Error reading subtree: %v", err)
}

// Access subtree properties
fmt.Printf("Subtree header magic: %s\n", readSubtree.SubtreeHeader.GetMagic())
fmt.Printf("Tile availability: %s\n", subtree.AsString(readSubtree.TileAvailability))
```

### Getting Subtree Information

```go
// Get detailed information about a subtree
info, err := subtree.GetSubtreeInfo(subtreeBytes, subtree.Quadtree)
if err != nil {
    log.Fatalf("Failed to get subtree info: %v", err)
}

// Print the information
subtree.PrintSubtreeInfo(info)
```

### Working with Bit Arrays

```go
// Convert string to bit array
bits := subtree.FromString("10110000")

// Convert bit array to string
bitString := subtree.AsString(bits)

// Convert bit array to byte array
bytes := subtree.ToByteArray(bits)

// Count bits with specific value
count := subtree.Count(bits, true) // Count true bits
```

### Morton Order Operations

```go
// Encode 2D coordinates to Morton order
mortonIndex := subtree.MortonOrder{}.Encode2D(5, 1)

// Decode Morton index to 2D coordinates
x, y := subtree.MortonOrder{}.Decode2D(mortonIndex)
```

## API Reference

### Subtree Structure

The main [Subtree](file:///Users/xuning/Work/go-3dtile/subtree/subtree_model.go#L3-L16) structure contains:

- `SubtreeHeader`: Header information
- `TileAvailability`: Tile availability bits
- `ContentAvailability`: Content availability bits
- `ChildSubtreeAvailability`: Child subtree availability bits
- `TileAvailabilityConstant`: Constant for tile availability
- `ContentAvailabilityConstant`: Constant for content availability

### Main Functions

- `ToBytesFromStrings()`: Create subtree bytes from string representations
- `ToBytes()`: Convert subtree structure to bytes
- `ReadSubtree()`: Read subtree from bytes
- `GetSubtreeInfo()`: Extract information from subtree bytes
- `PrintSubtreeInfo()`: Print detailed subtree information
- `FromString()`: Convert string to bit array
- `AsString()`: Convert bit array to string
- `ToByteArray()`: Convert bit array to byte array
- `Read()`: Read bit array from byte array with offset and length

## Testing

To run tests:

```bash
cd subtree
go test -v
```

## License

MIT