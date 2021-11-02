package tile3d

import (
	"encoding/json"
	"errors"
	"io"
)

const (
	TILE_REFINE_ADD     = "ADD"
	TILE_REFINE_REPLACE = "REPLACE"
)

var (
	TileDefaultTransform = [16]float64{1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0}
)

type Asset struct {
	Version        string `json:"version"`
	TilesetVersion string `json:"tilesetVersion,omitempty"`
	GltfUpAxis     string `json:"gltfUpAxis,omitempty"`
}

type Content struct {
	Url            string          `json:"uri"`
	BoundingVolume *BoundingVolume `json:"boundingVolume,omitempty"`
}

type Schema struct {
	Maximum float64 `json:"maximum,omitempty"`
	Minimum float64 `json:"minimum,omitempty"`
}

type BoundingVolume struct {
	Region *[]float64 `json:"region,omitempty"`
	Box    *[]float64 `json:"box,omitempty"`
	Sphere *[]float64 `json:"sphere,omitempty"`
}

func (b *BoundingVolume) SetBox(box []float64) error {
	if len(box) != 12 {
		return errors.New("box must 12 element!")
	}
	if b.Region != nil || b.Sphere != nil {
		b.Region = nil
		b.Sphere = nil
	}
	b.Box = &box
	return nil
}

func (b *BoundingVolume) SetRegion(region []float64) error {
	if len(region) != 6 {
		return errors.New("region must 6 element!")
	}
	if b.Box != nil || b.Sphere != nil {
		b.Box = nil
		b.Sphere = nil
	}
	b.Region = &region
	return nil
}

func (b *BoundingVolume) SetSphere(sphere []float64) error {
	if len(sphere) != 4 {
		return errors.New("sphere must 4 element!")
	}
	if b.Box != nil || b.Region != nil {
		b.Box = nil
		b.Region = nil
	}
	b.Sphere = &sphere
	return nil
}

func (b *BoundingVolume) GetRegion() []float64 {
	return *b.Region
}

func (b *BoundingVolume) GetBox() []float64 {
	return *b.Box
}

func (b *BoundingVolume) GetSphere() []float64 {
	return *b.Sphere
}

func (b *BoundingVolume) GetData() []float64 {
	if b.Region != nil {
		return *b.Region
	}
	if b.Box != nil {
		return *b.Box
	}
	if b.Sphere != nil {
		return *b.Sphere
	}
	return nil
}

type Tile struct {
	Content             *Content        `json:"content,omitempty"`
	BoundingVolume      BoundingVolume  `json:"boundingVolume,omitempty"`
	ViewerRequestVolume *BoundingVolume `json:"viewerRequestVolume,omitempty"`
	GeometricError      float64         `json:"geometricError"`
	Refine              string          `json:"refine"`
	Transform           *[16]float64    `json:"transform,omitempty"`
	Children            []Tile          `json:"children,omitempty"`
}

type Tileset struct {
	Asset              Asset              `json:"asset"`
	GeometricError     float64            `json:"geometricError"`
	Root               Tile               `json:"root"`
	Properties         *map[string]Schema `json:"properties,omitempty"`
	ExtensionsUsed     *[]string          `json:"extensionsUsed,omitempty"`
	ExtensionsRequired *[]string          `json:"extensionsRequired,omitempty"`
}

func (ts *Tileset) ToJson() (string, error) {
	b, e := json.Marshal(ts)
	return string(b), e
}

func TilesetFromJson(data io.Reader) *Tileset {
	var ts *Tileset
	json.NewDecoder(data).Decode(&ts)
	return ts
}
