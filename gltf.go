package tile3d

import (
	"bytes"
	"io"

	"github.com/qmuntal/gltf"
)

func openGltf(path string) *gltf.Document {
	doc, err := gltf.Open(path)
	if err != nil {
		return nil
	}
	return doc
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
	Size   uint32
}

func newSizeWriter() calcSizeWriter {
	wt := bytes.NewBuffer([]byte{})
	return calcSizeWriter{Size: 0, writer: wt}
}

func (w *calcSizeWriter) Write(p []byte) (n int, err error) {
	si := len(p)
	w.writer.Write(p)
	w.Size += uint32(si)
	return si, nil
}

func (w *calcSizeWriter) Bytes() []byte {
	return w.writer.(*bytes.Buffer).Bytes()
}

func (w *calcSizeWriter) GetSize() int {
	return len(w.Bytes())
}

func calcGltfSize(doc *gltf.Document, paddingUnit uint32) int64 {
	w := newSizeWriter()
	enc := gltf.NewEncoder(w.writer)
	enc.AsBinary = true
	if err := enc.Encode(doc); err != nil {
		return 0
	}
	return int64(calcPadding(w.Size, paddingUnit))
}

func getGltfBinary(doc *gltf.Document, paddingUnit uint32) ([]byte, error) {
	w := newSizeWriter()
	enc := gltf.NewEncoder(w.writer)
	enc.AsBinary = true
	if err := enc.Encode(doc); err != nil {
		return nil, err
	}
	padding := calcPadding(w.Size, paddingUnit)
	if padding == 0 {
		return w.Bytes(), nil
	}
	pad := make([]byte, padding)
	for i := range pad {
		pad[i] = 0x20
	}
	w.Write(pad)
	return w.Bytes(), nil
}
