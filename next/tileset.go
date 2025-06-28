package next

import (
	"encoding/json"

	"github.com/flywave/go-3dtile/next/gltf"
)

// Asset represents metadata about the entire tileset
type Asset struct {
	Version        string                     `json:"version"`
	TilesetVersion *string                    `json:"tilesetVersion,omitempty"`
	Extensions     map[string]json.RawMessage `json:"extensions,omitempty"`
	Extra          json.RawMessage            `json:"extra,omitempty"`
}

// DefaultAsset returns a new Asset with default values
func DefaultAsset() Asset {
	return Asset{
		Version: "1.1",
	}
}

// MetadataEntity represents metadata that conforms to a class
type MetadataEntity struct {
	Class      string                     `json:"class"`
	Properties map[string]json.RawMessage `json:"properties,omitempty"`
	Extensions map[string]json.RawMessage `json:"extensions,omitempty"`
	Extra      json.RawMessage            `json:"extra,omitempty"`
}

// BoundingVolume represents a volume that encloses a tile or content
type BoundingVolume struct {
	Box        *[12]float64               `json:"box,omitempty"`
	Region     *[6]float64                `json:"region,omitempty"`
	Sphere     *[4]float64                `json:"sphere,omitempty"`
	Extensions map[string]json.RawMessage `json:"extensions,omitempty"`
	Extra      json.RawMessage            `json:"extra,omitempty"`
}

// NewBoxBoundingVolume creates a new box bounding volume
func NewBoxBoundingVolume(box [12]float64) BoundingVolume {
	return BoundingVolume{
		Box: &box,
	}
}

// NewRegionBoundingVolume creates a new region bounding volume
func NewRegionBoundingVolume(region [6]float64) BoundingVolume {
	return BoundingVolume{
		Region: &region,
	}
}

// NewSphereBoundingVolume creates a new sphere bounding volume
func NewSphereBoundingVolume(sphere [4]float64) BoundingVolume {
	return BoundingVolume{
		Sphere: &sphere,
	}
}

// Subtrees describes the location of subtree files
type Subtrees struct {
	URI        string                     `json:"uri"`
	Extra      json.RawMessage            `json:"extra,omitempty"`
	Extensions map[string]json.RawMessage `json:"extensions,omitempty"`
}

// SubdivisionScheme defines the subdivision scheme
type SubdivisionScheme string

const (
	SubdivisionSchemeQuadtree SubdivisionScheme = "QUADTREE"
	SubdivisionSchemeOctree   SubdivisionScheme = "OCTREE"
)

// ImplicitTiling describes implicit subdivision of a tile
type ImplicitTiling struct {
	SubdivisionScheme SubdivisionScheme          `json:"subdivisionScheme"`
	SubtreeLevels     uint32                     `json:"subtreeLevels"`
	AvailableLevels   uint32                     `json:"availableLevels"`
	Subtrees          Subtrees                   `json:"subtrees"`
	Extra             json.RawMessage            `json:"extra,omitempty"`
	Extensions        map[string]json.RawMessage `json:"extensions,omitempty"`
}

// Refine defines the refinement type
type Refine string

const (
	RefineAdd     Refine = "ADD"
	RefineReplace Refine = "REPLACE"
)

// Content represents tile content metadata
type Content struct {
	BoundingVolume *BoundingVolume            `json:"boundingVolume,omitempty"`
	URI            string                     `json:"uri"`
	Metadata       *MetadataEntity            `json:"metadata,omitempty"`
	Group          *uint32                    `json:"group,omitempty"`
	Extensions     map[string]json.RawMessage `json:"extensions,omitempty"`
	Extra          json.RawMessage            `json:"extra,omitempty"`
}

// Tile represents a tile in the tileset
type Tile struct {
	BoundingVolume      BoundingVolume             `json:"boundingVolume"`
	ViewerRequestVolume *BoundingVolume            `json:"viewerRequestVolume,omitempty"`
	GeometricError      float64                    `json:"geometricError"`
	Refine              *Refine                    `json:"refine,omitempty"`
	Transform           [16]float64                `json:"transform,omitempty"`
	Content             *Content                   `json:"content,omitempty"`
	Contents            []Content                  `json:"contents,omitempty"`
	Metadata            *MetadataEntity            `json:"metadata,omitempty"`
	ImplicitTiling      *ImplicitTiling            `json:"implicitTiling,omitempty"`
	Children            []Tile                     `json:"children,omitempty"`
	Extensions          map[string]json.RawMessage `json:"extensions,omitempty"`
	Extra               json.RawMessage            `json:"extra,omitempty"`
}

// DefaultTile returns a Tile with default values
func DefaultTile() Tile {
	return Tile{
		BoundingVolume: BoundingVolume{},
		GeometricError: 0.0,
		Transform:      [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
	}
}

// GroupMetadata represents metadata about a group
type GroupMetadata struct {
	Class      string                     `json:"class"`
	Properties map[string]json.RawMessage `json:"properties,omitempty"`
	Extensions map[string]json.RawMessage `json:"extensions,omitempty"`
	Extra      json.RawMessage            `json:"extra,omitempty"`
}

// Tileset represents the root tileset object
type Tileset struct {
	Asset              Asset                      `json:"asset"`
	Properties         map[string]json.RawMessage `json:"properties,omitempty"`
	Schema             *gltf.Schema               `json:"schema,omitempty"`
	SchemaURI          *string                    `json:"schemaUri,omitempty"`
	Statistics         *Statistics                `json:"statistics,omitempty"`
	Groups             []GroupMetadata            `json:"groups,omitempty"`
	Metadata           *MetadataEntity            `json:"metadata,omitempty"`
	GeometricError     float64                    `json:"geometricError"`
	Root               Tile                       `json:"root"`
	ExtensionsUsed     []string                   `json:"extensionsUsed,omitempty"`
	ExtensionsRequired []string                   `json:"extensionsRequired,omitempty"`
	Extensions         map[string]json.RawMessage `json:"extensions,omitempty"`
	Extra              json.RawMessage            `json:"extra,omitempty"`
}

// Statistics contains statistics about entities
type Statistics struct {
	Classes    map[string]StatisticsClass `json:"classes,omitempty"`
	Extensions map[string]json.RawMessage `json:"extensions,omitempty"`
	Extra      json.RawMessage            `json:"extra,omitempty"`
}

// StatisticsClass contains statistics about a class
type StatisticsClass struct {
	Count      *uint32                       `json:"count,omitempty"`
	Properties map[string]StatisticsProperty `json:"properties,omitempty"`
	Extensions map[string]json.RawMessage    `json:"extensions,omitempty"`
	Extra      json.RawMessage               `json:"extra,omitempty"`
}

// StatisticsProperty contains statistics about a property
type StatisticsProperty struct {
	Min               *float64                   `json:"min,omitempty"`
	Max               *float64                   `json:"max,omitempty"`
	Mean              *float64                   `json:"mean,omitempty"`
	Median            *float64                   `json:"median,omitempty"`
	StandardDeviation *float64                   `json:"standardDeviation,omitempty"`
	Variance          *float64                   `json:"variance,omitempty"`
	Sum               *float64                   `json:"sum,omitempty"`
	Occurrences       json.RawMessage            `json:"occurrences,omitempty"`
	Extensions        map[string]json.RawMessage `json:"extensions,omitempty"`
	Extra             json.RawMessage            `json:"extra,omitempty"`
}

// IsIdentityMatrix checks if a 4x4 matrix is the identity matrix
func IsIdentityMatrix(m [16]float64) bool {
	return m == [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}
