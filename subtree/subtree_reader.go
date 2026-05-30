package subtree

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/flywave/go-3dtile/next"
)

// ReadBitArray reads a bit array from a byte array with the specified offset and length
func ReadBitArray(subtreeBinary []byte, offset int, length int) []bool {
	// Ensure we don't go out of bounds
	if offset < 0 || offset >= len(subtreeBinary) {
		return []bool{}
	}

	end := offset + length
	if end > len(subtreeBinary) {
		end = len(subtreeBinary)
	}

	// Extract the relevant bytes
	slicedBytes := subtreeBinary[offset:end]

	// Convert bytes to bit array
	bitCount := len(slicedBytes) * 8
	bitArray := make([]bool, bitCount)

	for i, b := range slicedBytes {
		for j := 0; j < 8; j++ {
			if (b & (1 << j)) != 0 {
				bitArray[i*8+j] = true
			}
		}
	}

	return bitArray
}

// ReadSubtreeFromReader reads a subtree from an io.Reader
func ReadSubtreeFromReader(reader io.Reader) (*Subtree, error) {
	// Read header
	header, err := NewSubtreeHeaderFromReader(reader)
	if err != nil {
		return nil, err
	}

	// Read JSON data
	jsonData := make([]byte, header.JsonByteLength)
	_, err = io.ReadFull(reader, jsonData)
	if err != nil {
		return nil, err
	}

	// Read binary data
	binaryData := make([]byte, header.BinaryByteLength)
	_, err = io.ReadFull(reader, binaryData)
	if err != nil {
		return nil, err
	}

	// Create subtree
	subtree := &Subtree{
		SubtreeHeader: header,
		SubtreeJson:   string(jsonData),
		SubtreeBinary: binaryData,
	}

	// Parse JSON
	var subtreeJson next.Subtree
	err = json.Unmarshal(jsonData, &subtreeJson)
	if err != nil {
		return nil, err
	}

	// Process tile availability
	if subtreeJson.TileAvailability != nil {
		if subtreeJson.TileAvailability.Bitstream != nil {
			bufferView := subtreeJson.BufferViews[*subtreeJson.TileAvailability.Bitstream]
			subtree.TileAvailability = ReadBitArray(subtree.SubtreeBinary, int(bufferView.ByteOffset), int(bufferView.ByteLength))
		} else if subtreeJson.TileAvailability.Constant != nil {
			subtree.TileAvailabilityConstant = *subtreeJson.TileAvailability.Constant
		}
	}

	// Process content availability (using first element as in C# code)
	if len(subtreeJson.ContentAvailability) > 0 {
		contentAvailability := subtreeJson.ContentAvailability[0]
		if contentAvailability.Bitstream != nil {
			bufferView := subtreeJson.BufferViews[*contentAvailability.Bitstream]
			subtree.ContentAvailability = ReadBitArray(subtree.SubtreeBinary, int(bufferView.ByteOffset), int(bufferView.ByteLength))
		} else if contentAvailability.Constant != nil {
			subtree.ContentAvailabilityConstant = *contentAvailability.Constant
		}
	}

	// Process child subtree availability
	if subtreeJson.ChildSubtreeAvailability != nil {
		if subtreeJson.ChildSubtreeAvailability.Bitstream != nil {
			bufferView := subtreeJson.BufferViews[*subtreeJson.ChildSubtreeAvailability.Bitstream]
			subtree.ChildSubtreeAvailability = ReadBitArray(subtree.SubtreeBinary, int(bufferView.ByteOffset), int(bufferView.ByteLength))
		} else if subtreeJson.ChildSubtreeAvailability.Constant != nil {
			// Handle constant value for child subtree availability
			subtree.ChildSubtreeAvailabilityConstant = *subtreeJson.ChildSubtreeAvailability.Constant
			// For constant 0, we set ChildSubtreeAvailability to an empty slice
			// For constant 1, we would need to know the expected size, but for now we leave it as nil
			if *subtreeJson.ChildSubtreeAvailability.Constant == 0 {
				subtree.ChildSubtreeAvailability = []bool{}
			}
		}
	}

	return subtree, nil
}

// ReadSubtree reads a subtree from a byte array
func ReadSubtree(data []byte) (*Subtree, error) {
	reader := bytes.NewReader(data)
	return ReadSubtreeFromReader(reader)
}
