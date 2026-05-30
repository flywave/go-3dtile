package tile3d

import (
	"os"
	"testing"
)

func TestVctrHeaderSize(t *testing.T) {
	h := &VctrHeader{}
	if s := h.CalcSize(); s != 44 {
		t.Errorf("VctrHeader.CalcSize() = %d, want 44", s)
	}
}

func TestVctrRead(t *testing.T) {
	paths := []string{
		"./data/tile.vctr",
		"./data/ll.vctr",
		"./data/with_batchtable.vctr",
		"./data/polygon.vctr",
		"./data/polygon_lr.vctr",
		"./data/polygon_parent.vctr",
		"./data/polygon_children.vctr",
		"./data/parent_batchtable.vctr",
	}
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			t.Skip("test data not available:", err)
		}
		vt := &Vctr{}
		if err := vt.Read(f); err != nil {
			t.Errorf("Read(%s) failed: %v", path, err)
		}
		f.Close()
		if string(vt.Header.Magic[:]) != VCTR_MAGIC {
			t.Errorf("%s: magic = %s, want %s", path, vt.Header.Magic, VCTR_MAGIC)
		}
	}
}

func TestVctrSetFeatureTable(t *testing.T) {
	m := &Vctr{}
	region := [6]float32{1, 2, 3, 4, 5, 6}
	m.SetFeatureTable(VctrFeatureTableView{
		PolygonsLength:  5,
		PolylinesLength: 3,
		PointsLength:    2,
		Region:          &region,
		PolygonBatchId:  []uint16{0, 1, 2, 3, 4},
	})

	if m.FeatureTable.Header[VCTR_PROP_POLYGONS_LENGTH] != uint32(5) {
		t.Errorf("POLYGONS_LENGTH = %v, want 5", m.FeatureTable.Header[VCTR_PROP_POLYGONS_LENGTH])
	}
}

func TestVctrFeatureTableDecode(t *testing.T) {
	header := map[string]interface{}{
		"POLYGONS_LENGTH":  float64(0),
		"POLYLINES_LENGTH": float64(0),
		"POINTS_LENGTH":    float64(0),
	}
	ret := VctrFeatureTableDecode(header, nil)
	if ret == nil {
		t.Error("VctrFeatureTableDecode returned nil")
	}
}
