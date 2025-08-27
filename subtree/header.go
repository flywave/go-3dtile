package subtree

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	SUBTREE_MAGIC = "subt"
)

type SubtreeHeader struct {
	Magic            [4]byte
	Version          uint32
	JsonByteLength   uint64
	BinaryByteLength uint64
}

func NewSubtreeHeader() *SubtreeHeader {
	h := &SubtreeHeader{}
	mg := []byte(SUBTREE_MAGIC)
	h.Magic[0] = mg[0]
	h.Magic[1] = mg[1]
	h.Magic[2] = mg[2]
	h.Magic[3] = mg[3]
	h.Version = 1
	return h
}

func NewSubtreeHeaderFromReader(reader io.Reader) (*SubtreeHeader, error) {
	h := &SubtreeHeader{}
	err := binary.Read(reader, binary.LittleEndian, h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (h *SubtreeHeader) CalcSize() int64 {
	return int64(binary.Size(h))
}

func (h *SubtreeHeader) GetMagic() string {
	return string(h.Magic[:])
}

func (h *SubtreeHeader) GetVersion() uint32 {
	return h.Version
}

func (h *SubtreeHeader) GetJsonByteLength() uint64 {
	return h.JsonByteLength
}

func (h *SubtreeHeader) GetBinaryByteLength() uint64 {
	return h.BinaryByteLength
}

func (h *SubtreeHeader) SetJsonByteLength(length uint64) {
	h.JsonByteLength = length
}

func (h *SubtreeHeader) SetBinaryByteLength(length uint64) {
	h.BinaryByteLength = length
}

func (h *SubtreeHeader) AsBinary() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, h)
	return buf.Bytes()
}

func (h *SubtreeHeader) Read(reader io.Reader) error {
	return binary.Read(reader, binary.LittleEndian, h)
}

func (h *SubtreeHeader) Write(writer io.Writer) error {
	return binary.Write(writer, binary.LittleEndian, h)
}
