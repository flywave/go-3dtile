package next

import (
	"encoding/json"
	"testing"
)

func TestSubtreeJSON(t *testing.T) {
	bv := BufferView{
		Buffer:     0,
		ByteOffset: 0,
		ByteLength: 16,
	}

	s := Subtree{
		Buffers: []Buffer{
			{ByteLength: 16},
		},
		BufferViews: []BufferView{bv},
		TileAvailability: &Availability{
			Bitstream:      uint32Ptr(0),
			AvailableCount: uint32Ptr(5),
		},
		ContentAvailability: []Availability{
			{
				Bitstream:      uint32Ptr(0),
				AvailableCount: uint32Ptr(3),
			},
		},
		ChildSubtreeAvailability: &Availability{
			Bitstream:      uint32Ptr(0),
			AvailableCount: uint32Ptr(2),
		},
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Subtree
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.TileAvailability == nil {
		t.Fatal("TileAvailability should not be nil")
	}
	if decoded.TileAvailability.Bitstream == nil || *decoded.TileAvailability.Bitstream != 0 {
		t.Errorf("TileAvailability.Bitstream = %v, want 0", decoded.TileAvailability.Bitstream)
	}
	if len(decoded.ContentAvailability) != 1 {
		t.Errorf("len(ContentAvailability) = %d, want 1", len(decoded.ContentAvailability))
	}
}

func TestSubtreeWithConstantAvailability(t *testing.T) {
	s := Subtree{
		Buffers: []Buffer{},
		TileAvailability: &Availability{
			Constant: uint8Ptr(1),
		},
		ContentAvailability: []Availability{
			{Constant: uint8Ptr(0)},
		},
		ChildSubtreeAvailability: &Availability{
			Constant: uint8Ptr(1),
		},
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Subtree
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.TileAvailability.Constant == nil || *decoded.TileAvailability.Constant != 1 {
		t.Errorf("TileAvailability.Constant = %v, want 1", decoded.TileAvailability.Constant)
	}
}

func TestSubtreeWithMultipleContent(t *testing.T) {
	s := Subtree{
		Buffers: []Buffer{{ByteLength: 32}},
		BufferViews: []BufferView{
			{Buffer: 0, ByteOffset: 0, ByteLength: 8},
			{Buffer: 0, ByteOffset: 8, ByteLength: 8},
		},
		TileAvailability: &Availability{
			Bitstream:      uint32Ptr(0),
			AvailableCount: uint32Ptr(10),
		},
		ContentAvailability: []Availability{
			{Bitstream: uint32Ptr(1), AvailableCount: uint32Ptr(5)},
			{Constant: uint8Ptr(1)},
		},
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Subtree
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if len(decoded.ContentAvailability) != 2 {
		t.Errorf("len(ContentAvailability) = %d, want 2", len(decoded.ContentAvailability))
	}
}

func TestBuffer(t *testing.T) {
	uri := "external.bin"
	b := Buffer{
		URI:        &uri,
		ByteLength: 1024,
		Name:       strPtr("test"),
	}

	data, err := json.Marshal(b)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Buffer
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.URI == nil || *decoded.URI != "external.bin" {
		t.Errorf("URI = %v, want external.bin", decoded.URI)
	}
	if decoded.ByteLength != 1024 {
		t.Errorf("ByteLength = %d, want 1024", decoded.ByteLength)
	}
}

func TestBufferView(t *testing.T) {
	bv := BufferView{
		Buffer:     1,
		ByteOffset: 16,
		ByteLength: 32,
	}

	data, err := json.Marshal(bv)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded BufferView
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.Buffer != 1 {
		t.Errorf("Buffer = %d, want 1", decoded.Buffer)
	}
	if decoded.ByteOffset != 16 {
		t.Errorf("ByteOffset = %d, want 16", decoded.ByteOffset)
	}
	if decoded.ByteLength != 32 {
		t.Errorf("ByteLength = %d, want 32", decoded.ByteLength)
	}
}

func TestSubtreeMetadata(t *testing.T) {
	s := Subtree{
		Buffers: []Buffer{},
		SubtreeMetadata: &MetadataEntity{
			Class: "metadataClass",
			Properties: map[string]json.RawMessage{
				"version": json.RawMessage(`"1.0"`),
			},
		},
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Subtree
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.SubtreeMetadata == nil {
		t.Fatal("SubtreeMetadata should not be nil")
	}
	if decoded.SubtreeMetadata.Class != "metadataClass" {
		t.Errorf("SubtreeMetadata.Class = %q, want metadataClass", decoded.SubtreeMetadata.Class)
	}
}

func TestSubtreeExtensions(t *testing.T) {
	s := Subtree{
		Buffers: []Buffer{},
		Extensions: map[string]json.RawMessage{
			"3DTILES_extension": json.RawMessage(`{"data": true}`),
		},
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}

	var decoded Subtree
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	if decoded.Extensions == nil {
		t.Fatal("Extensions should not be nil")
	}
}

func uint8Ptr(u uint8) *uint8 {
	return &u
}

func strPtr(s string) *string {
	return &s
}
