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

func TestBatchTableWriteReadRoundTripAllBatches(t *testing.T) {
	const batchCount = 3

	classification := []uint8{3, 7, 9}
	colors := []float32{
		1.5, 2.5, 3.5,
		4.5, 5.5, 6.5,
		7.5, 8.5, 9.5,
	}

	src := &BatchTable{
		Header: map[string]interface{}{
			"CLASSIFICATION": BinaryBodyReference{
				ComponentType: COMPONENT_TYPE_UNSIGNED_BYTE,
				ContainerType: CONTAINER_TYPE_SCALAR,
			},
			"COLOR": BinaryBodyReference{
				ComponentType: COMPONENT_TYPE_FLOAT,
				ContainerType: CONTAINER_TYPE_VEC3,
			},
		},
		Data: map[string]interface{}{
			"CLASSIFICATION": classification,
			"COLOR":          colors,
		},
	}

	h := &PntsHeader{}
	var buf bytes.Buffer
	if err := src.Write(&buf, h); err != nil {
		t.Fatal("Write failed:", err)
	}
	if h.GetBatchTableJSONByteLength() == 0 || h.GetBatchTableBinaryByteLength() == 0 {
		t.Fatalf("Write did not set batch table lengths: json=%d binary=%d",
			h.GetBatchTableJSONByteLength(), h.GetBatchTableBinaryByteLength())
	}
	wantClassification := []uint8{3, 7, 9}
	wantColors := []float32{1.5, 2.5, 3.5, 4.5, 5.5, 6.5, 7.5, 8.5, 9.5}
	for i := range classification {
		if classification[i] != wantClassification[i] {
			t.Fatalf("Write mutated Data[CLASSIFICATION]: got %v, want %v", classification, wantClassification)
		}
	}
	for i := range colors {
		if colors[i] != wantColors[i] {
			t.Fatalf("Write mutated Data[COLOR]: got %v, want %v", colors, wantColors)
		}
	}

	dst := &BatchTable{}
	if err := dst.Read(bytes.NewReader(buf.Bytes()), h, batchCount); err != nil {
		t.Fatal("Read failed:", err)
	}

	got, ok := dst.Data["CLASSIFICATION"].([]uint8)
	if !ok {
		t.Fatalf("CLASSIFICATION = %T (%v), want []uint8 covering all batches",
			dst.Data["CLASSIFICATION"], dst.Data["CLASSIFICATION"])
	}
	if len(got) != batchCount {
		t.Fatalf("CLASSIFICATION len = %d, want %d (one value per batch)", len(got), batchCount)
	}
	for i, want := range wantClassification {
		if got[i] != want {
			t.Errorf("CLASSIFICATION[%d] = %d, want %d", i, got[i], want)
		}
	}

	gotF, ok := dst.Data["COLOR"].([]float32)
	if !ok {
		t.Fatalf("COLOR = %T (%v), want []float32 covering all batches",
			dst.Data["COLOR"], dst.Data["COLOR"])
	}
	if len(gotF) != batchCount*3 {
		t.Fatalf("COLOR len = %d, want %d (3 components × %d batches)", len(gotF), batchCount*3, batchCount)
	}
	for i := 0; i < batchCount; i++ {
		for c := 0; c < 3; c++ {
			if gotF[i*3+c] != wantColors[i*3+c] {
				t.Errorf("COLOR batch %d component %d = %v, want %v", i, c, gotF[i*3+c], wantColors[i*3+c])
			}
		}
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
