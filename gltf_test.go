package tile3d

import (
	"os"
	"testing"
)

func TestLoadGlb(t *testing.T) {
	g := openGltf("./data/tree.glb")

	if g == nil {
		t.Error("error")
	}
}

func TestReadCmpt(t *testing.T) {
	f, _ := os.Open("./data/0-0.cmpt")
	cp := &Cmpt{}
	cp.Read(f)
}
