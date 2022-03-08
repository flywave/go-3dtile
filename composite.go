package tile3d

import (
	"encoding/binary"
	"io"
)

const (
	CMPT_MAGIC = "cmpt"
)

func NewCmpt() *Cmpt {
	m := &Cmpt{}
	mg := []byte(CMPT_MAGIC)
	m.Header.Magic[0] = mg[0]
	m.Header.Magic[1] = mg[1]
	m.Header.Magic[2] = mg[2]
	m.Header.Magic[3] = mg[3]
	return m
}

type CmptHeader struct {
	Magic       [4]byte
	Version     uint32
	ByteLength  uint32
	TilesLength uint32
}

func (h *CmptHeader) CalcSize() int64 {
	return 16
}

func (h *CmptHeader) GetByteLength() uint32 {
	return h.ByteLength
}

func (h *CmptHeader) GetTilesLength() uint32 {
	return h.TilesLength
}

func (h *CmptHeader) SetTilesLength(n uint32) {
	h.TilesLength = n
}

func (*CmptHeader) GetFeatureTableJSONByteLength() uint32   { return 0 }
func (*CmptHeader) GetFeatureTableBinaryByteLength() uint32 { return 0 }

func (*CmptHeader) GetBatchTableJSONByteLength() uint32   { return 0 }
func (*CmptHeader) GetBatchTableBinaryByteLength() uint32 { return 0 }

func (*CmptHeader) SetFeatureTableJSONByteLength(uint32)   {}
func (*CmptHeader) SetFeatureTableBinaryByteLength(uint32) {}

func (*CmptHeader) SetBatchTableJSONByteLength(uint32)   {}
func (*CmptHeader) SetBatchTableBinaryByteLength(uint32) {}

type Cmpt struct {
	Header CmptHeader
	Tiles  []TileModel
}

func (m *Cmpt) GetHeader() Header {
	return &m.Header
}

func (m *Cmpt) CalcSize() int64 {
	return m.Header.CalcSize()
}

func (m *Cmpt) Read(reader io.ReadSeeker) error {
	err := binary.Read(reader, littleEndian, &m.Header)
	if err != nil {
		return err
	}

	for i := 0; i < int(m.Header.TilesLength); i++ {
		magic := make([]byte, 4)
		_, err := reader.Read(magic)
		if err != nil {
			return err
		}
		reader.Seek(-4, io.SeekCurrent)
		switch string(magic) {
		case B3DM_MAGIC:
			b3 := new(B3dm)
			err := b3.Read(reader)
			if err != nil {
				return err
			}
			m.Tiles = append(m.Tiles, b3)
		case I3DM_MAGIC:
			i3 := new(I3dm)
			err := i3.Read(reader)
			if err != nil {
				return err
			}
			m.Tiles = append(m.Tiles, i3)
		case PNTS_MAGIC:
			pn := new(Pnts)
			err := pn.Read(reader)
			if err != nil {
				return err
			}
			m.Tiles = append(m.Tiles, pn)
		}
	}

	return nil
}

func (m *Cmpt) Write(writer io.Writer) error {
	m.Header.TilesLength = uint32(len(m.Tiles))
	for i := range m.Tiles {
		m.Header.ByteLength += uint32(m.Tiles[i].CalcSize())
	}

	err := binary.Write(writer, littleEndian, m.Header)

	if err != nil {
		return err
	}

	for i := range m.Tiles {
		err := m.Tiles[i].Write(writer)
		if err != nil {
			return err
		}
	}

	return nil
}
