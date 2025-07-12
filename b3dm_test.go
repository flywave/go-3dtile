package tile3d

import (
	"os"
	"testing"

	"github.com/flywave/gltf"
)

func TestB3dm(t *testing.T) {
	ph := "/home/hj/0-93.b3dm"
	f, _ := os.Open(ph)
	defer f.Close()
	b3d := NewB3dm()
	b3d.Read(f)
	gltf.SaveBinary(b3d.Model, "data/0-0.glb")
}
