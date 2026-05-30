package tile3d

import (
	"bytes"
	"testing"
)

func TestGeomHeaderSize(t *testing.T) {
	h := &GeomHeader{}
	if s := h.CalcSize(); s != 28 {
		t.Errorf("GeomHeader.CalcSize() = %d, want 28", s)
	}
}

func TestGeomSetFeatureTable(t *testing.T) {
	m := &Geom{}
	m.SetFeatureTable(GeomFeatureTableView{
		Boxs:       []GeomBox{{1, 2, 3, 4, 5, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		Spheres:    []GeomSphere{{0, 0, 0, 1}},
		BoxBatchId: []uint16{0},
	})

	if m.FeatureTable.Header[GEOM_PROP_BOXES] == nil {
		t.Error("BOXES should not be nil")
	}
}

func TestGeomFeatureTableDecode(t *testing.T) {
	header := map[string]interface{}{
		"BOXES_LENGTH":      float64(0),
		"CYLINDERS_LENGTH":  float64(0),
		"ELLIPSOIDS_LENGTH": float64(0),
		"SPHERES_LENGTH":    float64(0),
	}
	ret := GeomFeatureTableDecode(header, nil)
	if ret == nil {
		t.Error("GeomFeatureTableDecode returned nil")
	}
}

func TestGeomFeatureTableEncode(t *testing.T) {
	header := make(map[string]interface{})
	data := map[string]interface{}{
		GEOM_PROP_BOX_BATCH_IDS: []uint16{0, 1, 2},
		GEOM_PROP_BOXES:         []float32{0, 0, 0, 1, 1, 1},
	}

	result := GeomFeatureTableEncode(header, data)
	if len(result) == 0 {
		t.Error("GeomFeatureTableEncode returned empty")
	}
}

func TestGeomWriteRead(t *testing.T) {
	m := &Geom{}
	m.SetFeatureTable(GeomFeatureTableView{
		Boxs:       []GeomBox{{1, 2, 3, 4, 5, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		BoxBatchId: []uint16{0},
	})

	var buf bytes.Buffer
	if err := m.Write(&buf); err != nil {
		t.Fatal("Write failed:", err)
	}
	if buf.Len() == 0 {
		t.Error("Write produced empty output")
	}
}
