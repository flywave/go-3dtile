package mesh

import "encoding/json"

// ExtMeshFeatures represents the EXT_mesh_features glTF Mesh Primitive extension
type ExtMeshFeatures struct {
	FeatureIDs []FeatureID                `json:"featureIds"`
	Extensions map[string]json.RawMessage `json:"extensions,omitempty"`
	Extras     json.RawMessage            `json:"extras,omitempty"`
}

// FeatureID represents a feature ID set
type FeatureID struct {
	FeatureCount  uint32                     `json:"featureCount"`
	NullFeatureID *uint32                    `json:"nullFeatureId,omitempty"`
	Label         *string                    `json:"label,omitempty"`
	Attribute     *uint32                    `json:"attribute,omitempty"`
	Texture       *FeatureIDTexture          `json:"texture,omitempty"`
	PropertyTable *uint32                    `json:"propertyTable,omitempty"`
	Extensions    map[string]json.RawMessage `json:"extensions,omitempty"`
	Extras        json.RawMessage            `json:"extras,omitempty"`
}

// FeatureIDTexture represents a texture containing feature IDs
type FeatureIDTexture struct {
	Channels   []uint32                   `json:"channels,omitempty"`
	Index      uint32                     `json:"index"`
	TexCoord   uint32                     `json:"texCoord,omitempty"`
	Extensions map[string]json.RawMessage `json:"extensions,omitempty"`
	Extras     json.RawMessage            `json:"extras,omitempty"`
}

// DefaultChannels returns the default channels value
func DefaultChannels() []uint32 {
	return []uint32{0}
}

// IsDefaultChannels checks if channels has the default value
func IsDefaultChannels(channels []uint32) bool {
	return len(channels) == 1 && channels[0] == 0
}

// IsDefault checks if a value equals its type's default value
func IsDefault[T comparable](value T) bool {
	var zero T
	return value == zero
}
