package tile3d

import (
	"fmt"
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	path := "./data/vctr/0.vctr"
	rd, _ := os.Open(path)
	vt := &Vctr{}
	vt.Read(rd)
	ft := vt.GetFeatureTable()
	fmt.Println(ft)
	bt := vt.GetBatchTable()
	fmt.Println(bt)

	path2 := "./data/vctr/tile.vctr"
	rd2, _ := os.Open(path2)
	vt2 := &Vctr{}
	vt2.Read(rd2)
}
