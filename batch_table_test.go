package tile3d

import (
	"bytes"
	"strings"
	"testing"
)

func TestBatchTableReadJSON(t *testing.T) {
	bt := &BatchTable{}
	jsonStr := `{"names": {"byteOffset": 0, "componentType": "UNSIGNED_BYTE", "type": "SCALAR"}}`
	reader := strings.NewReader(jsonStr)
	err := bt.readJSONHeader(reader)
	if err != nil {
		t.Fatal("readJSONHeader failed:", err)
	}
	if bt.Header["names"] == nil {
		t.Error("names should not be nil")
	}
	ref, ok := bt.Header["names"].(BinaryBodyReference)
	if !ok {
		t.Fatal("names should be BinaryBodyReference")
	}
	if ref.ComponentType != "UNSIGNED_BYTE" {
		t.Errorf("ComponentType = %s, want UNSIGNED_BYTE", ref.ComponentType)
	}
}

func TestBatchTableWriteJSON(t *testing.T) {
	bt := &BatchTable{Header: map[string]interface{}{"test": "value"}}
	var buf bytes.Buffer
	err := bt.writeJSONHeader(&buf)
	if err != nil {
		t.Fatal("writeJSONHeader failed:", err)
	}
	if buf.Len() == 0 {
		t.Error("writeJSONHeader produced empty output")
	}
}

func TestBatchTableCalcSize(t *testing.T) {
	bt := &BatchTable{Header: map[string]interface{}{"test": float64(1)}}
	h := &B3dmHeader{}
	size := bt.CalcSize(h)
	if size <= 0 {
		t.Errorf("CalcSize() = %d, want > 0", size)
	}
}

func TestBatchTableGetProperty(t *testing.T) {
	bt := &BatchTable{Data: map[string]interface{}{
		"name": []uint8{1, 2, 3},
		"id":   []uint32{100, 200},
	}}
	if bt.GetProperty("name", 0) == nil {
		t.Error("GetProperty('name', 0) should not be nil")
	}
	if bt.GetProperty("nonexistent", 0) != nil {
		t.Error("GetProperty('nonexistent', 0) should be nil")
	}
}

func TestTransformBinaryBodyReference(t *testing.T) {
	m := map[string]interface{}{
		"prop": map[string]interface{}{
			"byteOffset":    float64(16),
			"componentType": "FLOAT",
			"type":          "VEC3",
		},
		"scalar": float64(42),
	}
	result := transformBinaryBodyReference(m)
	ref, ok := result["prop"].(BinaryBodyReference)
	if !ok {
		t.Fatal("prop should be BinaryBodyReference")
	}
	if ref.ByteOffset != 16 {
		t.Errorf("ByteOffset = %d, want 16", ref.ByteOffset)
	}
	if result["scalar"] != float64(42) {
		t.Error("scalar should be preserved")
	}
}
