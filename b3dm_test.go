package tile3d

import (
	"os"
	"testing"
)

func TestB3dm(t *testing.T) {
	ph := "/home/hj/workspace/cesium/Specs/Data/Cesium3DTiles/Hierarchy/BatchTableHierarchyBinary/tile.b3dm"
	f, _ := os.Open(ph)
	defer f.Close()
	b3d := NewB3dm()
	b3d.Read(f)
}
