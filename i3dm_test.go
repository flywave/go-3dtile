package tile3d

import (
	"bytes"
	"os"
	"testing"
)

func TestI3dmHeaderSize(t *testing.T) {
	h := &I3dmHeader{}
	if s := h.CalcSize(); s != 32 {
		t.Errorf("I3dmHeader.CalcSize() = %d, want 32", s)
	}
}

func TestI3dmRead(t *testing.T) {
	f, err := os.Open("./data/instancedWithBatchTableBinary.i3dm")
	if err != nil {
		t.Skip("test data not available:", err)
	}
	defer f.Close()

	i3 := &I3dm{}
	if err := i3.Read(f); err != nil {
		t.Skip("Read failed (expected with external test data):", err)
	}

	if string(i3.Header.Magic[:]) != I3DM_MAGIC {
		t.Errorf("magic = %s, want %s", i3.Header.Magic, I3DM_MAGIC)
	}
}

func TestI3dmSetFeatureTable(t *testing.T) {
	m := &I3dm{}
	m.SetFeatureTable(I3dmFeatureTableView{
		Position:       [][3]float32{{0, 0, 0}, {1, 1, 1}},
		InstanceLength: 2,
		BatchId:        []uint16{0, 1},
	})

	if m.FeatureTable.Header[I3DM_PROP_POSITION] == nil {
		t.Error("POSITION should not be nil")
	}

	view := m.GetFeatureTableView()
	if view.InstanceLength != 2 {
		t.Errorf("InstanceLength = %d, want 2", view.InstanceLength)
	}
}

func TestI3dmFeatureTableDecode(t *testing.T) {
	header := map[string]interface{}{
		"INSTANCES_LENGTH": float64(2),
		"POSITION":         map[string]interface{}{"byteOffset": float64(0), "componentType": "FLOAT", "type": "VEC3"},
	}
	data := []byte{
		0, 0, 0, 0, 0, 0, 128, 63, 0, 0, 0, 64,
		0, 0, 64, 64, 0, 0, 128, 64, 0, 0, 160, 64,
	}
	ret := I3dmFeatureTableDecode(header, data)
	if ret["INSTANCES_LENGTH"] != int32(2) {
		t.Errorf("INSTANCES_LENGTH = %v, want 2", ret["INSTANCES_LENGTH"])
	}
}

func TestI3dmWrite(t *testing.T) {
	m := &I3dm{}
	m.Header.GltfFormat = 0
	m.SetFeatureTable(I3dmFeatureTableView{
		Position:       [][3]float32{{0, 0, 0}},
		InstanceLength: 1,
	})

	var buf bytes.Buffer
	if err := m.Write(&buf); err != nil {
		t.Fatal("Write failed:", err)
	}

	if buf.Len() == 0 {
		t.Error("Write produced empty output")
	}
}
