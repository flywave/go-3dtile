package tile3d

import (
	"testing"
)

func TestComponentTypeSize(t *testing.T) {
	tests := []struct {
		tp      string
		expected int
	}{
		{"BYTE", 1},
		{"UNSIGNED_BYTE", 1},
		{"SHORT", 2},
		{"UNSIGNED_SHORT", 2},
		{"INT", 4},
		{"UNSIGNED_INT", 4},
		{"FLOAT", 4},
		{"DOUBLE", 8},
		{"UNKNOWN", 0},
	}
	for _, tt := range tests {
		got := ComponentTypeSize(tt.tp)
		if got != tt.expected {
			t.Errorf("ComponentTypeSize(%q) = %d, want %d", tt.tp, got, tt.expected)
		}
	}
}

func TestContainerTypeSize(t *testing.T) {
	tests := []struct {
		tp      string
		expected int
	}{
		{"SCALAR", 1},
		{"VEC2", 2},
		{"VEC3", 3},
		{"VEC4", 4},
		{"UNKNOWN", 0},
	}
	for _, tt := range tests {
		got := ContainerTypeSize(tt.tp)
		if got != tt.expected {
			t.Errorf("ContainerTypeSize(%q) = %d, want %d", tt.tp, got, tt.expected)
		}
	}
}

func TestBinaryBodyReferenceFromMap(t *testing.T) {
	r := &BinaryBodyReference{}
	m := map[string]interface{}{
		"byteOffset":    float64(42),
		"componentType": "FLOAT",
		"type":          "VEC3",
	}
	r.FromMap(m)
	if r.ByteOffset != 42 {
		t.Errorf("ByteOffset = %d, want 42", r.ByteOffset)
	}
	if r.ComponentType != "FLOAT" {
		t.Errorf("ComponentType = %s, want FLOAT", r.ComponentType)
	}
	if r.ContainerType != "VEC3" {
		t.Errorf("ContainerType = %s, want VEC3", r.ContainerType)
	}
}

func TestBinaryBodyReferenceGetMap(t *testing.T) {
	r := &BinaryBodyReference{
		ByteOffset:    42,
		ComponentType: "FLOAT",
		ContainerType: "VEC3",
	}
	m := r.GetMap()
	if m["byteOffset"] != float64(42) {
		t.Errorf("byteOffset = %v, want 42", m["byteOffset"])
	}
}

func TestBinaryBodyReferenceRoundTrip(t *testing.T) {
	r := &BinaryBodyReference{ByteOffset: 99, ComponentType: "UNSIGNED_INT", ContainerType: "SCALAR"}
	m := r.GetMap()
	r2 := &BinaryBodyReference{}
	r2.FromMap(m)
	if r.ByteOffset != r2.ByteOffset || r.ComponentType != r2.ComponentType || r.ContainerType != r2.ContainerType {
		t.Errorf("round trip failed: %+v -> %+v", r, r2)
	}
}

func TestMaxBatchId(t *testing.T) {
	if maxBatchId([]int64{1, 5, 3, 9, 2}) != 9 {
		t.Error("maxBatchId failed")
	}
	if maxBatchId([]int64{0xFF, 0x100}) != 0x100 {
		t.Error("maxBatchId failed for boundary")
	}
}

func TestAddReference(t *testing.T) {
	header := make(map[string]interface{})
	err := addReference(&header, "BATCH_ID", 0, "UNSIGNED_BYTE", "SCALAR", true)
	if err != nil {
		t.Fatal("addReference failed:", err)
	}
	if header["BATCH_ID"] == nil {
		t.Error("BATCH_ID should exist in header")
	}

	header2 := make(map[string]interface{})
	err = addReference(&header2, "BATCH_ID", 0, "UNSIGNED_SHORT", "SCALAR", true)
	if err != nil {
		t.Fatal("addReference failed:", err)
	}

	header3 := make(map[string]interface{})
	err = addReference(&header3, "BATCH_ID", 0, "FLOAT", "SCALAR", true)
	if err == nil {
		t.Error("addReference should fail with FLOAT component type for BATCH_ID")
	}
}
