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
	g := openGltf("./data/0-103.glb")
	if g == nil {
		t.Error("error")
	}
}

func TestReadB3dm(t *testing.T) {
	f1, _ := os.Open("./data/house/0-1.i3dm")
	// bt := make([]byte, 69)
	// f1.Read(bt)
	// fmt.Println(string(bt))
	cp1 := &I3dm{}
	cp1.Read(f1)
	bts, _ := mst.GetGltfBinary(cp1.Model, 8)
	ioutil.WriteFile("./data/house/0-1.glb", bts, 0755)
}

func TestReadi3dm(t *testing.T) {
	f, _ := os.Open("data/1-0.b3dm")
	// bt := make([]byte, 267)
	// f.Read(bt)
	// fmt.Println(string(bt))
	b3d := &B3dm{}
	b3d.Read(f)
	bts, _ := mst.GetGltfBinary(b3d.Model, 8)
	ioutil.WriteFile("./data/1-0.glb", bts, 0755)
}
