package subtree

import (
	"encoding/json"
)

// HandleBitArray processes a bit array and returns the bytes, true bit count, and buffer view
func HandleBitArray(bitArray []bool) ([]byte, int, *BufferView) {
	trueBits := Count(bitArray, true)
	bits := ToByteArray(bitArray)
	bytes := AddBinaryPadding(bits, 0)

	bufferView := &BufferView{
		Buffer:     0,
		ByteLength: len(bits),
		ByteOffset: 0,
	}

	return bytes, trueBits, bufferView
}

// ToSubtreeBinary converts a subtree to its binary representation
func ToSubtreeBinary(subtree *Subtree) ([]byte, *SubtreeJson) {
	substreamBinary := []byte{}
	subtreeJson := &SubtreeJson{}
	bufferViews := []BufferView{}

	// Process tile availability
	if subtree.TileAvailability != nil {
		resultTileAvailability, trueBits, bufferView := HandleBitArray(subtree.TileAvailability)
		bufferViews = append(bufferViews, *bufferView)
		substreamBinary = append(substreamBinary, resultTileAvailability...)

		subtreeJson.TileAvailability = &TileAvailability{
			Bitstream:      intPtr(0),
			AvailableCount: intPtr(trueBits),
		}
	} else {
		subtreeJson.TileAvailability = &TileAvailability{
			Constant: intPtr(subtree.TileAvailabilityConstant),
		}
	}

	// Process content availability
	if subtree.ContentAvailability != nil {
		resultContentAvailability, trueBits, bufferView := HandleBitArray(subtree.ContentAvailability)
		bufferView.ByteOffset = len(substreamBinary)

		subtreeJson.ContentAvailability = []ContentAvailability{
			{
				Bitstream:      intPtr(len(bufferViews)),
				AvailableCount: trueBits,
			},
		}

		bufferViews = append(bufferViews, *bufferView)
		substreamBinary = append(substreamBinary, resultContentAvailability...)
	} else {
		subtreeJson.ContentAvailability = []ContentAvailability{
			{
				Constant: intPtr(subtree.ContentAvailabilityConstant),
			},
		}
	}

	// Process child subtree availability
	if subtree.ChildSubtreeAvailability != nil {
		resultSubstreamAvailability, trueBits, bufferView := HandleBitArray(subtree.ChildSubtreeAvailability)
		bufferView.ByteOffset = len(substreamBinary)

		subtreeJson.ChildSubtreeAvailability = &ChildSubtreeAvailability{
			Bitstream:      intPtr(len(bufferViews)),
			AvailableCount: trueBits,
		}

		bufferViews = append(bufferViews, *bufferView)
		substreamBinary = append(substreamBinary, resultSubstreamAvailability...)
	} else {
		subtreeJson.ChildSubtreeAvailability = &ChildSubtreeAvailability{
			Constant: intPtr(0),
		}
	}

	// Set buffers and buffer views
	subtreeJson.Buffers = []Buffer{
		{
			ByteLength: len(substreamBinary),
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
func intPtr(i int) *int {
	return &i
}
