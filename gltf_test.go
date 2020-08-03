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

func TestReadB3dm(t *testing.T) {
	f, _ := os.Open("./data/0.b3dm")
	cp := &B3dm{}
	cp.Read(f)

	f1, _ := os.Open("./data/0-0.cmpt")
	cp1 := &Cmpt{}
	cp1.Read(f1)
}
