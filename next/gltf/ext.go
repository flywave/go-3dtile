package gltf

import "encoding/json"

// ExtStructuralMetadata represents the root EXT_structural_metadata object
type ExtStructuralMetadata struct {
	Schema             *Schema             `json:"schema,omitempty"`
	SchemaURI          *string             `json:"schemaUri,omitempty"`
	PropertyTables     []PropertyTable     `json:"propertyTables,omitempty"`
	PropertyTextures   []PropertyTexture   `json:"propertyTextures,omitempty"`
	PropertyAttributes []PropertyAttribute `json:"propertyAttributes,omitempty"`
}

// Schema defines classes and enums
type Schema struct {
	ID          string                     `json:"id"`
	Name        *string                    `json:"name,omitempty"`
	Description *string                    `json:"description,omitempty"`
	Version     *string                    `json:"version,omitempty"`
	Classes     map[string]Class           `json:"classes,omitempty"`
	Enums       map[string]Enum            `json:"enums,omitempty"`
	Extensions  map[string]json.RawMessage `json:"extensions,omitempty"`
	Extras      json.RawMessage            `json:"extras,omitempty"`
}

// Class defines a class with properties
type Class struct {
	Name        *string                    `json:"name,omitempty"`
	Description *string                    `json:"description,omitempty"`
	Properties  map[string]ClassProperty   `json:"properties,omitempty"`
	Extensions  map[string]json.RawMessage `json:"extensions,omitempty"`
	Extras      json.RawMessage            `json:"extras,omitempty"`
}

// ClassPropertyType defines the element type
type ClassPropertyType string

const (
	ClassPropertyTypeScalar  ClassPropertyType = "SCALAR"
	ClassPropertyTypeVec2    ClassPropertyType = "VEC2"
	ClassPropertyTypeVec3    ClassPropertyType = "VEC3"
	ClassPropertyTypeVec4    ClassPropertyType = "VEC4"
	ClassPropertyTypeMat2    ClassPropertyType = "MAT2"
	ClassPropertyTypeMat3    ClassPropertyType = "MAT3"
	ClassPropertyTypeMat4    ClassPropertyType = "MAT4"
	ClassPropertyTypeString  ClassPropertyType = "STRING"
	ClassPropertyTypeBoolean ClassPropertyType = "BOOLEAN"
	ClassPropertyTypeEnum    ClassPropertyType = "ENUM"
)

// ClassPropertyComponentType defines the component type
type ClassPropertyComponentType string

const (
	ClassPropertyComponentTypeInt8    ClassPropertyComponentType = "INT8"
	ClassPropertyComponentTypeUint8   ClassPropertyComponentType = "UINT8"
	ClassPropertyComponentTypeInt16   ClassPropertyComponentType = "INT16"
	ClassPropertyComponentTypeUint16  ClassPropertyComponentType = "UINT16"
	ClassPropertyComponentTypeInt32   ClassPropertyComponentType = "INT32"
	ClassPropertyComponentTypeUint32  ClassPropertyComponentType = "UINT32"
	ClassPropertyComponentTypeInt64   ClassPropertyComponentType = "INT64"
	ClassPropertyComponentTypeUint64  ClassPropertyComponentType = "UINT64"
	ClassPropertyComponentTypeFloat32 ClassPropertyComponentType = "FLOAT32"
	ClassPropertyComponentTypeFloat64 ClassPropertyComponentType = "FLOAT64"
)

// ClassProperty defines a property within a class
type ClassProperty struct {
	Name          *string                     `json:"name,omitempty"`
	Description   *string                     `json:"description,omitempty"`
	Type          ClassPropertyType           `json:"type"`
	ComponentType *ClassPropertyComponentType `json:"componentType,omitempty"`
	EnumType      *string                     `json:"enumType,omitempty"`
	Array         bool                        `json:"array,omitempty"`
	Count         *uint32                     `json:"count,omitempty"`
	Normalized    bool                        `json:"normalized,omitempty"`
	Offset        json.RawMessage             `json:"offset,omitempty"`
	Scale         json.RawMessage             `json:"scale,omitempty"`
	Max           json.RawMessage             `json:"max,omitempty"`
	Min           json.RawMessage             `json:"min,omitempty"`
	Required      bool                        `json:"required,omitempty"`
	NoData        json.RawMessage             `json:"noData,omitempty"`
	Default       json.RawMessage             `json:"default,omitempty"`
	Semantic      *string                     `json:"semantic,omitempty"`
	Extensions    map[string]json.RawMessage  `json:"extensions,omitempty"`
	Extras        json.RawMessage             `json:"extras,omitempty"`
}

// EnumValueType defines the type of enum values
type EnumValueType string

const (
	EnumValueTypeInt8   EnumValueType = "INT8"
	EnumValueTypeUint8  EnumValueType = "UINT8"
	EnumValueTypeInt16  EnumValueType = "INT16"
	EnumValueTypeUint16 EnumValueType = "UINT16"
	EnumValueTypeInt32  EnumValueType = "INT32"
	EnumValueTypeUint32 EnumValueType = "UINT32"
	EnumValueTypeInt64  EnumValueType = "INT64"
	EnumValueTypeUint64 EnumValueType = "UINT64"
)

// Enum defines an enumeration
type Enum struct {
	Name        *string                    `json:"name,omitempty"`
	Description *string                    `json:"description,omitempty"`
	ValueType   EnumValueType              `json:"valueType,omitempty"`
	Values      []EnumValue                `json:"values"`
	Extensions  map[string]json.RawMessage `json:"extensions,omitempty"`
	Extras      json.RawMessage            `json:"extras,omitempty"`
}

// EnumValue defines a value within an enum
type EnumValue struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Value       int32   `json:"value"`
}

// PropertyTable contains property values
type PropertyTable struct {
	Name       *string                          `json:"name,omitempty"`
	Class      string                           `json:"class"`
	Count      uint32                           `json:"count"`
	Properties map[string]PropertyTableProperty `json:"properties,omitempty"`
	Extensions map[string]json.RawMessage       `json:"extensions,omitempty"`
	Extras     json.RawMessage                  `json:"extras,omitempty"`
}

// PropertyTexture contains texture-based property values
type PropertyTexture struct {
	Name       *string                            `json:"name,omitempty"`
	Class      string                             `json:"class"`
	Properties map[string]PropertyTextureProperty `json:"properties,omitempty"`
}

// PropertyTextureProperty defines texture-based property storage
type PropertyTextureProperty struct {
	Channels []uint32        `json:"channels,omitempty"`
	Offset   json.RawMessage `json:"offset,omitempty"`
	Scale    json.RawMessage `json:"scale,omitempty"`
	Max      json.RawMessage `json:"max,omitempty"`
	Min      json.RawMessage `json:"min,omitempty"`
}

// PropertyAttribute contains attribute-based property values
type PropertyAttribute struct {
	Name       *string                              `json:"name,omitempty"`
	Class      string                               `json:"class"`
	Properties map[string]PropertyAttributeProperty `json:"properties,omitempty"`
}

// PropertyAttributeProperty defines attribute-based property storage
type PropertyAttributeProperty struct {
	Attribute string          `json:"attribute"`
	Offset    json.RawMessage `json:"offset,omitempty"`
	Scale     json.RawMessage `json:"scale,omitempty"`
	Max       json.RawMessage `json:"max,omitempty"`
	Min       json.RawMessage `json:"min,omitempty"`
}

// OffsetType defines the type of offset values
type OffsetType string

const (
	OffsetTypeUint8  OffsetType = "UINT8"
	OffsetTypeUint16 OffsetType = "UINT16"
	OffsetTypeUint32 OffsetType = "UINT32"
	OffsetTypeUint64 OffsetType = "UINT64"
)

// PropertyTableProperty defines table-based property storage
type PropertyTableProperty struct {
	Values           uint32                     `json:"values"`
	ArrayOffsets     *uint32                    `json:"arrayOffsets,omitempty"`
	StringOffsets    *uint32                    `json:"stringOffsets,omitempty"`
	ArrayOffsetType  OffsetType                 `json:"arrayOffsetType,omitempty"`
	StringOffsetType OffsetType                 `json:"stringOffsetType,omitempty"`
	Offset           json.RawMessage            `json:"offset,omitempty"`
	Scale            json.RawMessage            `json:"scale,omitempty"`
	Max              json.RawMessage            `json:"max,omitempty"`
	Min              json.RawMessage            `json:"min,omitempty"`
	Extensions       map[string]json.RawMessage `json:"extensions,omitempty"`
	Extras           json.RawMessage            `json:"extras,omitempty"`
}

// DefaultChannels returns the default channels value
func DefaultChannels() []uint32 {
	return []uint32{0}
}

// IsDefaultChannels checks if channels has the default value
func IsDefaultChannels(channels []uint32) bool {
	return len(channels) == 1 && channels[0] == 0
}
