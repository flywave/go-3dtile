package tile3d

import "testing"

func TestLoadGlb(t *testing.T) {
	g := openGltf("./data/tree.glb")

	if g == nil {
		t.Error("error")
	}
}
