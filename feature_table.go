package tile3d

import (
	"bytes"
	"encoding/json"
	"io"
)

type featureTableDecode func(header map[string]interface{}, buff []byte) map[string]interface{}
type featureTableEncode func(header map[string]interface{}, data map[string]interface{}) []byte

type FeatureTable struct {
	Header map[string]interface{}
	Data   map[string]interface{}

	decode featureTableDecode
	encode featureTableEncode
}

func (t *FeatureTable) readJSONHeader(data io.ReadSeeker, jsonLength int) error {
	jdata := make([]byte, jsonLength)
	_, err := data.Read(jdata)
	dec := json.NewDecoder(bytes.NewBuffer(jdata))
	if err != nil {
		return nil
	}
	t.Header = make(map[string]interface{})
	if err := dec.Decode(&t.Header); err != nil {
		return err
	}
	t.Header = transformBinaryBodyReference(t.Header)
	return nil
}

func (t *FeatureTable) writeJSONHeader(wr io.Writer) (int, error) {
	var jdata []byte
	buf := bytes.NewBuffer(jdata)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(t.Header); err != nil {
		return 0, err
	}
	jdata = buf.Bytes()
	n, err := wr.Write(jdata)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (h *FeatureTable) calcJSONSize() int64 {
	w := newSizeWriter()
	if _, err := h.writeJSONHeader(w.writer); err != nil {
		return 0
	}
	return int64(w.GetSize())
}

func (h *FeatureTable) CalcSize(header Header) int64 {
	w := newSizeWriter()
	if err := h.Write(w.writer, header); err != nil {
		return 0
	}
	return int64(w.GetSize())
}

func (h *FeatureTable) GetBatchLength() int {
	if h.Header["BATCH_LENGTH"] != nil {
		switch d := h.Header["BATCH_LENGTH"].(type) {
		case int:
			return d
		case float64:
			return int(d)
		}
	}
	return 0
}

func (h *FeatureTable) readData(reader io.ReadSeeker, buffLength int) error {
	bdata := make([]byte, buffLength)
	_, err := reader.Read(bdata)
	if err != nil {
		return err
	}
	h.Data = h.decode(h.Header, bdata)
	return nil
}

func (h *FeatureTable) writeData(wr io.Writer) (int, error) {
	buff := h.encode(h.Header, h.Data)
	if buff != nil {
		n, err := wr.Write(buff)
		if err != nil {
			return 0, err
		}
		return n, nil
	}
	return 0, nil
}

func (h *FeatureTable) Read(reader io.ReadSeeker, header Header) error {
	err := h.readJSONHeader(reader, int(header.GetFeatureTableJSONByteLength()))
	if err != nil {
		return err
	}
	err = h.readData(reader, int(header.GetFeatureTableBinaryByteLength()))
	if err != nil {
		return err
	}
	return nil
}

func (h *FeatureTable) Write(writer io.Writer, header Header) error {
	JSONLenght, err := h.writeJSONHeader(writer)
	if err != nil {
		return err
	}
	BinaryLenght, err := h.writeData(writer)
	if err != nil {
		return err
	}
	if header != nil {
		header.SetFeatureTableJSONByteLength(uint32(JSONLenght))
		header.SetFeatureTableBinaryByteLength(uint32(BinaryLenght))
	}
	return nil
}
