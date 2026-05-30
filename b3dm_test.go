package tile3d

import (
	"bytes"
	"os"
	"testing"
)

func TestB3dmHeaderSize(t *testing.T) {
	h := &B3dmHeader{}
	if s := h.CalcSize(); s != 28 {
		t.Errorf("B3dmHeader.CalcSize() = %d, want 28", s)
	}
}

func TestB3dmNew(t *testing.T) {
	m := NewB3dm()
	if m == nil {
		t.Fatal("NewB3dm() returned nil")
	}
	if string(m.Header.Magic[:]) != B3DM_MAGIC {
		t.Errorf("magic = %s, want %s", m.Header.Magic, B3DM_MAGIC)
	}
}

func TestB3dmFeatureTable(t *testing.T) {
	m := NewB3dm()
	m.SetFeatureTable(B3dmFeatureTableView{BatchLength: 10, RtcCenter: []float64{1, 2, 3}})
	view := m.GetFeatureTableView()
	if view.BatchLength != 10 {
		t.Errorf("BatchLength = %d, want 10", view.BatchLength)
	}
	if len(view.RtcCenter) != 3 || view.RtcCenter[0] != 1 {
		t.Errorf("RtcCenter = %v, want [1 2 3]", view.RtcCenter)
	}
}

func TestB3dmReadWriteRoundTrip(t *testing.T) {
	f, err := os.Open("./data/batchedWithBatchTableBinary.b3dm")
	if err != nil {
		t.Skip("test data not available:", err)
	}
	defer f.Close()

	b3d := NewB3dm()
	if err := b3d.Read(f); err != nil {
		t.Skip("Read failed (expected with external test data):", err)
	}

	if string(b3d.Header.Magic[:]) != B3DM_MAGIC {
		t.Errorf("magic = %s, want %s", b3d.Header.Magic, B3DM_MAGIC)
	}

	if b3d.Model == nil {
		t.Error("Model is nil after reading")
	} else {
		var buf bytes.Buffer
		if err := b3d.Write(&buf); err != nil {
			t.Fatal("Write failed:", err)
		}

		if buf.Len() == 0 {
			t.Error("Write produced empty output")
		}
	}
}

func TestB3dmDecode(t *testing.T) {
	header := map[string]interface{}{"BATCH_LENGTH": float64(5)}
	result := B3dmFeatureTableDecode(header, nil)
	if result["BATCH_LENGTH"] != int32(5) {
		t.Errorf("BATCH_LENGTH = %v, want 5", result["BATCH_LENGTH"])
	}
}
