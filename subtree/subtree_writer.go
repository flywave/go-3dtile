package subtree

import (
	"encoding/json"

	"github.com/flywave/go-3dtile/next"
)

const boundary = 8

// AddPadding adds padding to a byte array
func AddPadding(bytes []byte, offset int) []byte {
	remainder := (offset + len(bytes)) % boundary
	padding := 0
	if remainder != 0 {
		padding = boundary - remainder
	}

	// Create padding bytes (whitespace)
	paddingBytes := make([]byte, padding)
	for i := range paddingBytes {
		paddingBytes[i] = ' '
	}

	// Concatenate original bytes with padding
	result := make([]byte, len(bytes)+padding)
	copy(result, bytes)
	copy(result[len(bytes):], paddingBytes)

	return result
}

// AddPaddingString adds padding to a string
func AddPaddingString(input string, offset int) string {
	bytes := []byte(input)
	paddedBytes := AddPadding(bytes, offset)
	return string(paddedBytes)
}

// AddBinaryPadding adds binary padding to a byte array
func AddBinaryPadding(bytes []byte, offset int) []byte {
	remainder := (offset + len(bytes)) % boundary
	padding := 0
	if remainder != 0 {
		padding = boundary - remainder
	}

	// Create padding bytes (zeros)
	paddingBytes := make([]byte, padding)

	// Concatenate original bytes with padding
	result := make([]byte, len(bytes)+padding)
	copy(result, bytes)
	copy(result[len(bytes):], paddingBytes)

	return result
}

// GetByteArray creates a byte array of the specified length filled with zeros
func GetByteArray(length int) []byte {
	arr := make([]byte, length)
	// Already filled with zeros by default
	return arr
}

// HandleBitArray processes a bit array and returns the bytes, true bit count, and buffer view
func HandleBitArray(bitArray []bool) ([]byte, int, *next.BufferView) {
	trueBits := Count(bitArray, true)
	bits := ToByteArray(bitArray)
	bytes := AddBinaryPadding(bits, 0)

	bufferView := &next.BufferView{
		Buffer:     0,
		ByteLength: uint32(len(bits)),
		ByteOffset: 0,
	}

	return bytes, trueBits, bufferView
}

// ToSubtreeBinary converts a subtree to its binary representation
func ToSubtreeBinary(subtree *Subtree) ([]byte, *next.Subtree) {
	substreamBinary := []byte{}
	subtreeJson := &next.Subtree{}
	bufferViews := []next.BufferView{}

	// Process tile availability
	if subtree.TileAvailability != nil {
		resultTileAvailability, trueBits, bufferView := HandleBitArray(subtree.TileAvailability)
		bufferViews = append(bufferViews, *bufferView)
		substreamBinary = append(substreamBinary, resultTileAvailability...)

		subtreeJson.TileAvailability = &next.Availability{
			Bitstream:      uint32Ptr(0),
			AvailableCount: uint32Ptr(uint32(trueBits)),
		}
	} else {
		subtreeJson.TileAvailability = &next.Availability{
			Constant: uint8Ptr(subtree.TileAvailabilityConstant),
		}
	}

	// Process content availability
	if subtree.ContentAvailability != nil {
		resultContentAvailability, trueBits, bufferView := HandleBitArray(subtree.ContentAvailability)
		bufferView.ByteOffset = uint32(len(substreamBinary))

		subtreeJson.ContentAvailability = []next.Availability{
			{
				Bitstream:      uint32Ptr(uint32(len(bufferViews))),
				AvailableCount: uint32Ptr(uint32(trueBits)),
			},
		}

		bufferViews = append(bufferViews, *bufferView)
		substreamBinary = append(substreamBinary, resultContentAvailability...)
	} else {
		subtreeJson.ContentAvailability = []next.Availability{
			{
				Constant: uint8Ptr(subtree.ContentAvailabilityConstant),
			},
		}
	}

	// Process child subtree availability
	if subtree.ChildSubtreeAvailability != nil {
		resultSubstreamAvailability, trueBits, bufferView := HandleBitArray(subtree.ChildSubtreeAvailability)
		bufferView.ByteOffset = uint32(len(substreamBinary))

		subtreeJson.ChildSubtreeAvailability = &next.Availability{
			Bitstream:      uint32Ptr(uint32(len(bufferViews))),
			AvailableCount: uint32Ptr(uint32(trueBits)),
		}

		bufferViews = append(bufferViews, *bufferView)
		substreamBinary = append(substreamBinary, resultSubstreamAvailability...)
	} else {
		subtreeJson.ChildSubtreeAvailability = &next.Availability{
			Constant: uint8Ptr(0),
		}
	}

	// Set buffers and buffer views
	subtreeJson.Buffers = []next.Buffer{
		{
			ByteLength: uint32(len(substreamBinary)),
		},
	}

	subtreeJson.BufferViews = bufferViews

	return substreamBinary, subtreeJson
}

// ToBytes converts a subtree to its byte representation
func ToBytes(subtree *Subtree) []byte {
	bin, subtreeJson := ToSubtreeBinary(subtree)

	// Serialize subtreeJson to JSON
	jsonData, err := json.Marshal(subtreeJson)
	if err != nil {
		// Handle error appropriately in real implementation
		return []byte{}
	}

	// Add padding to JSON data
	subtreeJsonPadded := AddPaddingString(string(jsonData), 0)
	subtreeBinaryPadded := AddBinaryPadding(bin, 0)

	// Create the final byte array
	var result []byte

	// Add header
	header := subtree.SubtreeHeader
	header.JsonByteLength = uint64(len(subtreeJsonPadded))
	header.BinaryByteLength = uint64(len(subtreeBinaryPadded))

	// Convert header to bytes
	headerBytes := header.AsBinary()
	result = append(result, headerBytes...)

	// Add JSON data
	result = append(result, []byte(subtreeJsonPadded)...)

	// Add binary data
	result = append(result, subtreeBinaryPadded...)

	return result
}

// ToBytesFromStrings creates a subtree from string representations and converts it to bytes
func ToBytesFromStrings(tileAvailability string, contentAvailability string, subtreeAvailability string) []byte {
	subtree := NewSubtree()
	subtree.TileAvailability = FromString(tileAvailability)

	if contentAvailability != "" {
		subtree.ContentAvailability = FromString(contentAvailability)
	}

	if subtreeAvailability != "" {
		subtree.ChildSubtreeAvailability = FromString(subtreeAvailability)
	}

	return ToBytes(subtree)
}

// Helper function to create pointer to int
func uint32Ptr(i uint32) *uint32 {
	return &i
}

func uint8Ptr(i uint8) *uint8 {
	return &i
}
