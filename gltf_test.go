package tile3d

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/flywave/go-mst"
)

func TestLoadGlb(t *testing.T) {
	// f, _ := os.Open("./data/100.glb")
	// bt, _ := ioutil.ReadAll(f)
	// fmt.Println(len(bt))
	g := openGltf("./data/building.glb")
	if g == nil {
		t.Error("error")
	}
}

func TestReadB3dm(t *testing.T) {
	f1, _ := os.Open("./data/0.cmpt")
	// bt := make([]byte, 69)
	// f1.Read(bt)
	// fmt.Println(string(bt))
	cp1 := &Cmpt{}
	cp1.Read(f1)
	b3d := cp1.Tiles[0].(*B3dm)
	bts, _ := mst.GetGltfBinary(b3d.Model, 8)
	ioutil.WriteFile("./data/0.glb", bts, 0755)
}

func TestReadi3dm(t *testing.T) {
	f, _ := os.Open("/home/hj/workspace/GISCore/build/public/Upload/Tilesets/yanshi/11-12/14/13533/6366/0-34.b3dm")
	// bt := make([]byte, 267)
	// f.Read(bt)
	// fmt.Println(string(bt))
	b3d := &B3dm{}
	b3d.Read(f)
	bts, _ := mst.GetGltfBinary(b3d.Model, 8)
	ioutil.WriteFile("./data/0.glb", bts, 0755)

}
