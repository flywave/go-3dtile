package tile3d

import (
	"fmt"
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	path := "./data/vctr/1/1/0.vctr"
	rd, _ := os.Open(path)
	vt := &Vctr{}
	vt.Read(rd)
	ft := vt.GetFeatureTable()
	fmt.Println(ft)
	bt := vt.GetBatchTable()
	fmt.Println(bt)
}
