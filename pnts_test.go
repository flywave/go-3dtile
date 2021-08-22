package tile3d

import (
	"os"
	"testing"
)

func TestPnts(t *testing.T) {
	p := &Pnts{}
	f, _ := os.Open("./data/7.pnts")
	p.Read(f)
}
