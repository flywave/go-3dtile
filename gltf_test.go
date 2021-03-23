package tile3d

import (
	"os"
	"testing"
)

func TestLoadGlb(t *testing.T) {
	// f, _ := os.Open("./data/100.glb")
	// bt, _ := ioutil.ReadAll(f)
	// fmt.Println(len(bt))
	g := openGltf("./data/100.glb")
	if g == nil {
		t.Error("error")
	}
}

func TestReadB3dm(t *testing.T) {
	f1, _ := os.Open("./data/0-0.cmpt")
	// bt := make([]byte, 69)
	// f1.Read(bt)
	// fmt.Println(string(bt))
	cp1 := &Cmpt{}
	cp1.Read(f1)
}

func TestReadi3dm(t *testing.T) {
	f, _ := os.Open("./data/0-0.i3dm")
	// bt := make([]byte, 267)
	// f.Read(bt)
	// fmt.Println(string(bt))
	cp := &I3dm{}
	cp.Read(f)
}
