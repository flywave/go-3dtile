package tile3d

import (
	"io"

	"github.com/qmuntal/gltf"
)

func openGltf(path string) *gltf.Document {
	doc, ok := gltf.Open(path)
	if ok == nil {
		return doc
	}
	return nil
}

func loadGltfFromByte(reader io.Reader) (*gltf.Document, error) {
	dec := gltf.NewDecoder(reader)
	doc := new(gltf.Document)
	if err := dec.Decode(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func saveGltf(path string, doc *gltf.Document) error {
	return gltf.Save(doc, path)
}

func writeGltfBinary(writer io.Writer, doc *gltf.Document) error {
	enc := gltf.NewEncoder(writer)
	enc.AsBinary = true
	if err := enc.Encode(doc); err != nil {
		return err
	}
	return nil
}

type calcSizeWriter struct {
	writer io.Writer
	Size   int
	Data   []byte
}

func newSizeWriter() calcSizeWriter {
	return calcSizeWriter{Size: int(0)}
}

func (w *calcSizeWriter) Write(p []byte) (n int, err error) {
	si := len(p)
	w.Data = append(w.Data, p...)
	w.Size += int(si)
	return si, nil
}

func calcGltfSize(doc *gltf.Document, paddingUnit int) int64 {
	w := newSizeWriter()
	enc := gltf.NewEncoder(w.writer)
	enc.AsBinary = true
	if err := enc.Encode(doc); err != nil {
		return 0
	}
	return int64(calcPadding(w.Size, paddingUnit))
}

func getGltfBinary(doc *gltf.Document, paddingUnit int) ([]byte, error) {
	w := newSizeWriter()
	enc := gltf.NewEncoder(w.writer)
	enc.AsBinary = true
	if err := enc.Encode(doc); err != nil {
		return nil, err
	}
	padding := calcPadding(w.Size, paddingUnit)
	if padding == 0 {
		return w.Data, nil
	}
	pad := make([]byte, padding)
	for i := range pad {
		pad[i] = 0x20
	}
	w.Data = append(w.Data, pad...)
	return w.Data, nil
}
