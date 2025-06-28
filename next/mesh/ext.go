package mesh

import "encoding/json"

// ExtStructuralMetadata represents the EXT_structural_metadata glTF Mesh Primitive extension
type ExtStructuralMetadata struct {
	PropertyTextures   []uint32                   `json:"propertyTextures,omitempty"`
	PropertyAttributes []uint32                   `json:"propertyAttributes,omitempty"`
	Extensions         map[string]json.RawMessage `json:"extensions,omitempty"`
	Extras             json.RawMessage            `json:"extras,omitempty"`
}
