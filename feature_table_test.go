package tile3d

import (
	"bytes"
	"testing"
)

func TestFeatureTableGetBatchLength(t *testing.T) {
	ft := &FeatureTable{Data: map[string]interface{}{"BATCH_LENGTH": 42}}
	if bl := ft.GetBatchLength(); bl != 42 {
		t.Errorf("GetBatchLength() = %d, want 42", bl)
	}

	ft2 := &FeatureTable{Data: map[string]interface{}{"BATCH_LENGTH": float64(42)}}
	if bl := ft2.GetBatchLength(); bl != 42 {
		t.Errorf("GetBatchLength() = %d, want 42", bl)
	}

	ft3 := &FeatureTable{Data: map[string]interface{}{}}
	if bl := ft3.GetBatchLength(); bl != 0 {
		t.Errorf("GetBatchLength() = %d, want 0", bl)
	}
}

func TestFeatureTableCalcSize(t *testing.T) {
	ft := &FeatureTable{Header: map[string]interface{}{"test": float64(1)}}
	h := &B3dmHeader{}

	size := ft.CalcSize(h)
	if size <= 0 {
		t.Errorf("CalcSize() = %d, want > 0", size)
	}
}

func TestFeatureTableWriteJSON(t *testing.T) {
	ft := &FeatureTable{Header: map[string]interface{}{"BATCH_LENGTH": float64(5)}}
	h := &B3dmHeader{}

	var buf bytes.Buffer
	err := ft.Write(&buf, h)
	if err != nil {
		t.Fatal("Write failed:", err)
	}
	if buf.Len() == 0 {
		t.Error("Write produced empty output")
	}
}

func TestFeatureTableReadJSONHeader(t *testing.T) {
	jsonStr := `{"BATCH_LENGTH": 10, "RTC_CENTER": [1.0, 2.0, 3.0]}`
	ft := &FeatureTable{}
	buf := bytes.NewReader([]byte(jsonStr))
	err := ft.readJSONHeader(buf, len(jsonStr))
	if err != nil {
		t.Fatal("readJSONHeader failed:", err)
	}
	if ft.Header["BATCH_LENGTH"] == nil {
		t.Error("BATCH_LENGTH should not be nil")
	}
}
