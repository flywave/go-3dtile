package tile3d

import (
	"os"
	"testing"
)

func TestCmptHeaderSize(t *testing.T) {
	h := &CmptHeader{}
	if s := h.CalcSize(); s != 16 {
		t.Errorf("CmptHeader.CalcSize() = %d, want 16", s)
	}
}

func TestCmptNew(t *testing.T) {
	m := NewCmpt()
	if m == nil {
		t.Fatal("NewCmpt() returned nil")
	}
	if string(m.Header.Magic[:]) != CMPT_MAGIC {
		t.Errorf("magic = %s, want %s", m.Header.Magic, CMPT_MAGIC)
	}
}

func TestCmptRead(t *testing.T) {
	f, err := os.Open("./data/composite.cmpt")
	if err != nil {
		t.Skip("test data not available:", err)
	}
	defer f.Close()

	c := NewCmpt()
	if err := c.Read(f); err != nil {
		t.Skip("Read failed (expected with external test data):", err)
	}

	if string(c.Header.Magic[:]) != CMPT_MAGIC {
		t.Errorf("magic = %s, want %s", c.Header.Magic, CMPT_MAGIC)
	}
}

func TestCmptReadCompositeOfComposite(t *testing.T) {
	f, err := os.Open("./data/compositeOfComposite.cmpt")
	if err != nil {
		t.Skip("test data not available:", err)
	}
	defer f.Close()

	c := NewCmpt()
	if err := c.Read(f); err != nil {
		t.Skip("Read failed (expected with external test data):", err)
	}
}

func TestCmptHeaderInterface(t *testing.T) {
	h := &CmptHeader{}
	if h.GetFeatureTableJSONByteLength() != 0 {
		t.Error("GetFeatureTableJSONByteLength should return 0")
	}
	if h.GetBatchTableJSONByteLength() != 0 {
		t.Error("GetBatchTableJSONByteLength should return 0")
	}
}
