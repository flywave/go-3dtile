package subtree

import (
	"bytes"
	"testing"
)

func TestSubtreeHeaderCreation(t *testing.T) {
	header := NewSubtreeHeader()

	if header.GetMagic() != SUBTREE_MAGIC {
		t.Errorf("Expected magic %s, got %s", SUBTREE_MAGIC, header.GetMagic())
	}

	if header.GetVersion() != 1 {
		t.Errorf("Expected version 1, got %d", header.GetVersion())
	}

	if header.GetJsonByteLength() != 0 {
		t.Errorf("Expected JSON byte length 0, got %d", header.GetJsonByteLength())
	}

	if header.GetBinaryByteLength() != 0 {
		t.Errorf("Expected binary byte length 0, got %d", header.GetBinaryByteLength())
	}
}

func TestSubtreeHeaderSetters(t *testing.T) {
	header := NewSubtreeHeader()

	header.SetJsonByteLength(100)
	header.SetBinaryByteLength(200)

	if header.GetJsonByteLength() != 100 {
		t.Errorf("Expected JSON byte length 100, got %d", header.GetJsonByteLength())
	}

	if header.GetBinaryByteLength() != 200 {
		t.Errorf("Expected binary byte length 200, got %d", header.GetBinaryByteLength())
	}
}

func TestSubtreeHeaderAsBinary(t *testing.T) {
	header := NewSubtreeHeader()
	header.SetJsonByteLength(100)
	header.SetBinaryByteLength(200)

	binaryData := header.AsBinary()

	if len(binaryData) != int(header.CalcSize()) {
		t.Errorf("Expected binary data length %d, got %d", header.CalcSize(), len(binaryData))
	}
}

func TestSubtreeHeaderReadAndWrite(t *testing.T) {
	// Create original header
	originalHeader := NewSubtreeHeader()
	originalHeader.SetJsonByteLength(1000)
	originalHeader.SetBinaryByteLength(2000)

	// Write to buffer
	buf := new(bytes.Buffer)
	err := originalHeader.Write(buf)
	if err != nil {
		t.Fatalf("Failed to write header: %v", err)
	}

	// Read from buffer
	newHeader := &SubtreeHeader{}
	err = newHeader.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read header: %v", err)
	}

	// Verify values
	if newHeader.GetMagic() != originalHeader.GetMagic() {
		t.Errorf("Magic mismatch: expected %s, got %s", originalHeader.GetMagic(), newHeader.GetMagic())
	}

	if newHeader.GetVersion() != originalHeader.GetVersion() {
		t.Errorf("Version mismatch: expected %d, got %d", originalHeader.GetVersion(), newHeader.GetVersion())
	}

	if newHeader.GetJsonByteLength() != originalHeader.GetJsonByteLength() {
		t.Errorf("JSON byte length mismatch: expected %d, got %d", originalHeader.GetJsonByteLength(), newHeader.GetJsonByteLength())
	}

	if newHeader.GetBinaryByteLength() != originalHeader.GetBinaryByteLength() {
		t.Errorf("Binary byte length mismatch: expected %d, got %d", originalHeader.GetBinaryByteLength(), newHeader.GetBinaryByteLength())
	}
}

func TestNewSubtreeHeaderFromReader(t *testing.T) {
	// Create original header
	originalHeader := NewSubtreeHeader()
	originalHeader.SetJsonByteLength(500)
	originalHeader.SetBinaryByteLength(1500)

	// Write to buffer
	buf := new(bytes.Buffer)
	err := originalHeader.Write(buf)
	if err != nil {
		t.Fatalf("Failed to write header: %v", err)
	}

	// Create header from reader
	newHeader, err := NewSubtreeHeaderFromReader(buf)
	if err != nil {
		t.Fatalf("Failed to create header from reader: %v", err)
	}

	// Verify values
	if newHeader.GetMagic() != originalHeader.GetMagic() {
		t.Errorf("Magic mismatch: expected %s, got %s", originalHeader.GetMagic(), newHeader.GetMagic())
	}

	if newHeader.GetVersion() != originalHeader.GetVersion() {
		t.Errorf("Version mismatch: expected %d, got %d", originalHeader.GetVersion(), newHeader.GetVersion())
	}

	if newHeader.GetJsonByteLength() != originalHeader.GetJsonByteLength() {
		t.Errorf("JSON byte length mismatch: expected %d, got %d", originalHeader.GetJsonByteLength(), newHeader.GetJsonByteLength())
	}

	if newHeader.GetBinaryByteLength() != originalHeader.GetBinaryByteLength() {
		t.Errorf("Binary byte length mismatch: expected %d, got %d", originalHeader.GetBinaryByteLength(), newHeader.GetBinaryByteLength())
	}
}
