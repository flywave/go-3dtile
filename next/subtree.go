package next

import (
	"encoding/json"

	ext_gltf "github.com/flywave/gltf/ext/3dtile/gltf"
)

// Subtree describes the availability of tiles and content in a subtree
type Subtree struct {
	Buffers                  []Buffer                   `json:"buffers,omitempty"`
	BufferViews              []BufferView               `json:"bufferViews,omitempty"`
	PropertyTables           []ext_gltf.PropertyTable   `json:"propertyTables,omitempty"`
	TileAvailability         Availability               `json:"tileAvailability"`
	ContentAvailability      []Availability             `json:"contentAvailability,omitempty"`
	ChildSubtreeAvailability Availability               `json:"childSubtreeAvailability"`
	TileMetadata             *uint32                    `json:"tileMetadata,omitempty"`
	ContentMetadata          []uint32                   `json:"contentMetadata,omitempty"`
	SubtreeMetadata          *MetadataEntity            `json:"subtreeMetadata,omitempty"`
	Extensions               map[string]json.RawMessage `json:"extensions,omitempty"`
	Extra                    json.RawMessage            `json:"extra,omitempty"`
}

// Availability describes the availability of a set of elements
type Availability struct {
	Bitstream      *uint32                    `json:"bitstream,omitempty"`
	AvailableCount *uint32                    `json:"availableCount,omitempty"`
	Constant       *uint8                     `json:"constant,omitempty"`
	Extensions     map[string]json.RawMessage `json:"extensions,omitempty"`
	Extra          json.RawMessage            `json:"extra,omitempty"`
}

// Buffer represents a binary blob
type Buffer struct {
	URI        *string                    `json:"uri,omitempty"`
	ByteLength uint32                     `json:"byteLength"`
	Name       *string                    `json:"name,omitempty"`
	Extensions map[string]json.RawMessage `json:"extensions,omitempty"`
	Extra      json.RawMessage            `json:"extra,omitempty"`
}

// BufferView represents a contiguous subset of a buffer
type BufferView struct {
	Buffer     uint32                     `json:"buffer"`
	ByteOffset uint32                     `json:"byteOffset"`
	ByteLength uint32                     `json:"byteLength"`
	Name       *string                    `json:"name,omitempty"`
	Extensions map[string]json.RawMessage `json:"extensions,omitempty"`
	Extra      json.RawMessage            `json:"extra,omitempty"`
}
