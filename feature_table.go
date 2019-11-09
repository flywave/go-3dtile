package tile3d

import (
	"encoding/json"
	"io"
)

type FeatureTableConvert func(header map[string]interface{}, buff []byte) map[string]interface{}

type FeatureTable struct {
	Header map[string]interface{}
	Data   map[string]interface{}
}

func (t *FeatureTable) readJSONHeader(data io.ReadSeeker) error {
	dec := json.NewDecoder(data)
	if err := dec.Decode(&t.Header); err != nil {
		return err
	}
	t.Header = transformBinaryBodyReference(t.Header)
	return nil
}

func (t *FeatureTable) writeJSONHeader(wr io.Writer) error {
	enc := json.NewEncoder(wr)
	if err := enc.Encode(t.Header); err != nil {
		return err
	}
	return nil
}

func (h *FeatureTable) calcJSONSize() int64 {
	w := newSizeWriter()
	if err := h.writeJSONHeader(w.writer); err != nil {
		return 0
	}
	return int64(w.Size)
}

func (h *FeatureTable) convertFeatureTableData(f FeatureTableConvert, buff []byte) {
	h.Data = f(h.Header, buff)
}

func (h *FeatureTable) CalcSize() int64 {
	return 0
}

func (h *FeatureTable) readData(reader io.ReadSeeker) error {
	return nil
}

func (h *FeatureTable) writeData(wr io.Writer) error {
	return nil
}

func (h *FeatureTable) Read(reader io.ReadSeeker, header Header) error {
	return nil
}

func (h *FeatureTable) Write(writer io.Writer, header Header) error {
	return nil
}
