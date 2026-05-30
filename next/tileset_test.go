package next

import (
	"encoding/json"
	"testing"
)

func TestDefaultAsset(t *testing.T) {
	a := DefaultAsset()
	if a.Version != "1.1" {
		t.Errorf("Version = %q, want 1.1", a.Version)
	}
}

func TestDefaultTile(t *testing.T) {
	ti := DefaultTile()
	if ti.GeometricError != 0.0 {
		t.Errorf("GeometricError = %v, want 0.0", ti.GeometricError)
	}
}

func TestDefaultStyle(t *testing.T) {
	s := DefaultStyle()
	if s.Show.BooleanExpression == nil {
		t.Error("Show.BooleanExpression should not be nil")
	}
	if s.Color.ColorExpression != "color('#FFFFFF')" {
		t.Errorf("Color.ColorExpression = %q, want color('#FFFFFF')", s.Color.ColorExpression)
	}
}

func TestBoundingVolume(t *testing.T) {
	box := [12]float64{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1}
	bv := NewBoxBoundingVolume(box)
	if bv.Box == nil || *bv.Box != box {
		t.Error("NewBoxBoundingVolume failed")
	}

	region := [6]float64{-1, -1, -1, 1, 1, 1}
	bv2 := NewRegionBoundingVolume(region)
	if bv2.Region == nil || *bv2.Region != region {
		t.Error("NewRegionBoundingVolume failed")
	}

	sphere := [4]float64{0, 0, 0, 1}
	bv3 := NewSphereBoundingVolume(sphere)
	if bv3.Sphere == nil || *bv3.Sphere != sphere {
		t.Error("NewSphereBoundingVolume failed")
	}
}

func TestIsIdentityMatrix(t *testing.T) {
	if !IsIdentityMatrix([16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}) {
		t.Error("identity matrix should be detected")
	}
	if IsIdentityMatrix([16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 1}) {
		t.Error("non-identity matrix should not be detected")
	}
}

func TestTilesetJSONRoundTrip(t *testing.T) {
	ts := Tileset{
		Asset:          DefaultAsset(),
		GeometricError: 500.0,
		Root: Tile{
			BoundingVolume: NewBoxBoundingVolume([12]float64{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1}),
			GeometricError: 100.0,
			Refine:         refPtr(RefineAdd),
			Content: &Content{
				URI: "content.glb",
			},
		},
	}

	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Tileset
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.Asset.Version != "1.1" {
		t.Errorf("Version = %q, want 1.1", decoded.Asset.Version)
	}
	if decoded.GeometricError != 500.0 {
		t.Errorf("GeometricError = %v, want 500.0", decoded.GeometricError)
	}
	if decoded.Root.Content == nil || decoded.Root.Content.URI != "content.glb" {
		t.Errorf("Root.Content.URI = %v, want content.glb", decoded.Root.Content)
	}
}

func TestImplicitTiling(t *testing.T) {
	it := ImplicitTiling{
		SubdivisionScheme: SubdivisionSchemeQuadtree,
		SubtreeLevels:     5,
		AvailableLevels:   4,
		Subtrees: Subtrees{
			URI: "{level}/{x}/{y}.subtree",
		},
	}

	data, err := json.Marshal(it)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded ImplicitTiling
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.SubdivisionScheme != SubdivisionSchemeQuadtree {
		t.Errorf("SubdivisionScheme = %q, want QUADTREE", decoded.SubdivisionScheme)
	}
	if decoded.Subtrees.URI != "{level}/{x}/{y}.subtree" {
		t.Errorf("Subtrees.URI = %q, want {level}/{x}/{y}.subtree", decoded.Subtrees.URI)
	}
}

func TestMetadataEntity(t *testing.T) {
	me := MetadataEntity{
		Class: "testClass",
		Properties: map[string]json.RawMessage{
			"height": json.RawMessage(`100.5`),
			"name":   json.RawMessage(`"building"`),
		},
	}

	data, err := json.Marshal(me)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded MetadataEntity
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.Class != "testClass" {
		t.Errorf("Class = %q, want testClass", decoded.Class)
	}
}

func TestProperties(t *testing.T) {
	ts := Tileset{
		Asset:          DefaultAsset(),
		GeometricError: 500.0,
		Properties: map[string]Property{
			"height": {
				Minimum: float64Ptr(0.0),
				Maximum: float64Ptr(100.0),
			},
		},
		Root: DefaultTile(),
	}

	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Tileset
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.Properties == nil {
		t.Fatal("Properties should not be nil")
	}
	prop, ok := decoded.Properties["height"]
	if !ok {
		t.Fatal("height property missing")
	}
	if prop.Minimum == nil || *prop.Minimum != 0.0 {
		t.Errorf("height.Minimum = %v, want 0.0", prop.Minimum)
	}
	if prop.Maximum == nil || *prop.Maximum != 100.0 {
		t.Errorf("height.Maximum = %v, want 100.0", prop.Maximum)
	}
}

func TestGroupMetadata(t *testing.T) {
	gm := GroupMetadata{
		Class: "groupClass",
	}

	data, err := json.Marshal(gm)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded GroupMetadata
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.Class != "groupClass" {
		t.Errorf("Class = %q, want groupClass", decoded.Class)
	}
}

func TestStatistics(t *testing.T) {
	s := Statistics{
		Classes: map[string]StatisticsClass{
			"building": {
				Count: uint32Ptr(100),
				Properties: map[string]StatisticsProperty{
					"height": {
						Minimum: float64Ptr(0.0),
						Maximum: float64Ptr(50.0),
						Mean:    float64Ptr(25.0),
					},
				},
			},
		},
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Statistics
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	buildingClass, ok := decoded.Classes["building"]
	if !ok {
		t.Fatal("building class missing")
	}
	if buildingClass.Count == nil || *buildingClass.Count != 100 {
		t.Errorf("Count = %v, want 100", buildingClass.Count)
	}
	if heightProp, ok := buildingClass.Properties["height"]; !ok {
		t.Fatal("height property missing")
	} else if heightProp.Minimum == nil || *heightProp.Minimum != 0.0 {
		t.Errorf("height.minimum = %v, want 0.0", heightProp.Minimum)
	}
}

func TestTilesetWithMultipleContents(t *testing.T) {
	ts := Tileset{
		Asset:          DefaultAsset(),
		GeometricError: 500.0,
		Root: Tile{
			BoundingVolume: NewBoxBoundingVolume([12]float64{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1}),
			GeometricError: 0.0,
			Contents: []Content{
				{URI: "building.glb"},
				{URI: "trees.glb"},
			},
		},
	}

	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Tileset
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if len(decoded.Root.Contents) != 2 {
		t.Errorf("len(Contents) = %d, want 2", len(decoded.Root.Contents))
	}
}

func TestTilesetExtensions(t *testing.T) {
	ts := Tileset{
		Asset:              DefaultAsset(),
		GeometricError:     500.0,
		Root:               DefaultTile(),
		ExtensionsUsed:     []string{"3DTILES_multiple_contents"},
		ExtensionsRequired: []string{"3DTILES_multiple_contents"},
	}

	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Tileset
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if len(decoded.ExtensionsUsed) != 1 || decoded.ExtensionsUsed[0] != "3DTILES_multiple_contents" {
		t.Errorf("ExtensionsUsed = %v, want [3DTILES_multiple_contents]", decoded.ExtensionsUsed)
	}
}

func refPtr(r Refine) *Refine {
	return &r
}

func float64Ptr(f float64) *float64 {
	return &f
}

func uint32Ptr(u uint32) *uint32 {
	return &u
}
