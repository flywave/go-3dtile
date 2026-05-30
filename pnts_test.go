package tile3d

import (
	"bytes"
	"os"
	"testing"
)

func TestPntsHeaderSize(t *testing.T) {
	h := &PntsHeader{}
	if s := h.CalcSize(); s != 28 {
		t.Errorf("PntsHeader.CalcSize() = %d, want 28", s)
	}
}

func TestPntsNew(t *testing.T) {
	m := NewPnts()
	if m == nil {
		t.Fatal("NewPnts() returned nil")
	}
	if string(m.Header.Magic[:]) != PNTS_MAGIC {
		t.Errorf("magic = %s, want %s", m.Header.Magic, PNTS_MAGIC)
	}
}

func TestPntsRead(t *testing.T) {
	f, err := os.Open("./data/7.pnts")
	if err != nil {
		t.Skip("test data not available:", err)
	}
	defer f.Close()

	p := NewPnts()
	if err := p.Read(f); err != nil {
		t.Fatal("Read failed:", err)
	}

	if string(p.Header.Magic[:]) != PNTS_MAGIC {
		t.Errorf("magic = %s, want %s", p.Header.Magic, PNTS_MAGIC)
	}
}

func TestPntsSetFeatureTable(t *testing.T) {
	m := NewPnts()
	m.SetFeatureTable(PntsFeatureTableView{
		Position:      [][3]float32{{1, 2, 3}, {4, 5, 6}},
		PointsLength:  2,
		ConstantRGBA:  []uint8{255, 0, 0, 255},
	})

	if m.FeatureTable.Header[PNTS_PROP_POINTS_LENGTH] != uint32(2) {
		t.Errorf("POINTS_LENGTH = %v, want 2", m.FeatureTable.Header[PNTS_PROP_POINTS_LENGTH])
	}

	if m.FeatureTable.Header[PNTS_PROP_CONSTANT_RGBA] == nil {
		t.Error("CONSTANT_RGBA should not be nil")
	}
}

func TestPntsFeatureTableDecode(t *testing.T) {
	header := map[string]interface{}{
		"POINTS_LENGTH": float64(3),
		"POSITION":      map[string]interface{}{"byteOffset": float64(0), "componentType": "FLOAT", "type": "VEC3"},
	}
	data := []byte{
		0, 0, 0, 0, 0, 0, 128, 63, 0, 0, 0, 64,
		0, 0, 64, 64, 0, 0, 128, 64, 0, 0, 160, 64,
		0, 0, 192, 64, 0, 0, 224, 64, 0, 0, 0, 65,
	}
	ret := PntsFeatureTableDecode(header, data)
	if ret["POINTS_LENGTH"] != int32(3) {
		t.Errorf("POINTS_LENGTH = %v, want 3", ret["POINTS_LENGTH"])
	}
}

func TestPntsWrite(t *testing.T) {
	m := NewPnts()
	m.SetFeatureTable(PntsFeatureTableView{
		Position:     [][3]float32{{0, 0, 0}, {1, 1, 1}},
		PointsLength: 2,
	})

	var buf bytes.Buffer
	if err := m.Write(&buf); err != nil {
		t.Fatal("Write failed:", err)
	}

	if buf.Len() == 0 {
		t.Error("Write produced empty output")
	}
}
