package subtree

import (
	"bytes"
	"encoding/json"
	"io"
)

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
	var subtreeJson SubtreeJson
	err = json.Unmarshal(jsonData, &subtreeJson)
	if err != nil {
		return nil, err
	}

	// Process tile availability
	if subtreeJson.TileAvailability != nil {
		if subtreeJson.TileAvailability.Bitstream != nil {
			bufferView := subtreeJson.BufferViews[*subtreeJson.TileAvailability.Bitstream]
			subtree.TileAvailability = Read(subtree.SubtreeBinary, bufferView.ByteOffset, bufferView.ByteLength)
		} else if subtreeJson.TileAvailability.Constant != nil {
			subtree.TileAvailabilityConstant = *subtreeJson.TileAvailability.Constant
		}
	}

	// Process content availability (using first element as in C# code)
	if len(subtreeJson.ContentAvailability) > 0 {
		contentAvailability := subtreeJson.ContentAvailability[0]
		if contentAvailability.Bitstream != nil {
			bufferView := subtreeJson.BufferViews[*contentAvailability.Bitstream]
			subtree.ContentAvailability = Read(subtree.SubtreeBinary, bufferView.ByteOffset, bufferView.ByteLength)
		} else if contentAvailability.Constant != nil {
			subtree.ContentAvailabilityConstant = *contentAvailability.Constant
		}
	}

	// Process child subtree availability
	if subtreeJson.ChildSubtreeAvailability != nil && subtreeJson.ChildSubtreeAvailability.Bitstream != nil {
		bufferView := subtreeJson.BufferViews[*subtreeJson.ChildSubtreeAvailability.Bitstream]
		subtree.ChildSubtreeAvailability = Read(subtree.SubtreeBinary, bufferView.ByteOffset, bufferView.ByteLength)
	}

	return subtree, nil
}

// ReadSubtree reads a subtree from a byte array
func ReadSubtree(data []byte) (*Subtree, error) {
	reader := bytes.NewReader(data)
	return ReadSubtreeFromReader(reader)
}
