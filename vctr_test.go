package tile3d

import (
	"fmt"
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	path := "./data/parent_batchtable.vctr"
	rd, _ := os.Open(path)
	vt := &Vctr{}
	vt.Read(rd)
	fmt.Println(vt)
}
