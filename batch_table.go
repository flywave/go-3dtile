package tile3d

import (
	"bytes"
	"encoding/json"
	"io"
)

type BatchTable struct {
	Header map[string]interface{}
	Data   map[string]interface{}
}

func transformBinaryBodyReference(m map[string]interface{}) map[string]interface{} {
	ref := make(map[string]interface{})
	for k, v := range m {
		switch tp := v.(type) {
		case map[string]interface{}:
			r := new(BinaryBodyReference)
			r.FromMap(tp)
			ref[k] = *r
		default:
			ref[k] = v
		}
	}
	return ref
}

func (t *BatchTable) readJSONHeader(data io.Reader) error {
	t.Header = make(map[string]interface{})
	dec := json.NewDecoder(data)
	if err := dec.Decode(&t.Header); err != nil {
		return err
	}
	t.Header = transformBinaryBodyReference(t.Header)
	return nil
}

func (t *BatchTable) writeJSONHeader(wr io.Writer) error {
	enc := json.NewEncoder(wr)
	if err := enc.Encode(t.Header); err != nil {
		return err
	}
	return nil
}

func (h *BatchTable) calcJSONSize() int64 {
	w := newSizeWriter()
	if err := h.writeJSONHeader(w.writer); err != nil {
		return 0
	}
	return int64(w.Size)
}

func (h *BatchTable) CalcSize(batchLength int) int64 {
	outJSONHeader := make(map[string]interface{})
	helper := BinaryBodySizeHelper{Header: &outJSONHeader}
	for k, v := range h.Header {
		switch t := v.(type) {
		case BinaryBodyReference:
			helper.addProperty(k, h.Data[k], t.ComponentType, t.ContainerType, false)
		default:
			outJSONHeader[k] = v
		}
	}
	helper.finished()
	return int64(helper.calcHeaderSize(0) + helper.Size)
}

func (h *BatchTable) GetProperty(property string, batchId int) interface{} {
	ret := h.Data[property]
	if ret == nil {
		return nil
	}
	switch t := h.Data[property].(type) {
	case byte:
		return t
	case []byte:
		return t
	case int8:
		return t
	case []int8:
		return t
	case int16:
		return t
	case []int16:
		return t
	case uint16:
		return t
	case []uint16:
		return t
	case int32:
		return t
	case []int32:
		return t
	case uint32:
		return t
	case []uint32:
		return t
	case float32:
		return t
	case []float32:
		return t
	case float64:
		return t
	case []float64:
		return t
	}
	return nil
}

func (h *BatchTable) Read(reader io.ReadSeeker, header Header, batchLength int) error {
	if header.GetBatchTableJSONByteLength() <= 0 {
		return nil
	}
	jsonb := make([]byte, header.GetBatchTableJSONByteLength())

	if _, err := reader.Read(jsonb); err != nil {
		return err
	}
	jsonr := bytes.NewReader(jsonb)
	if err := h.readJSONHeader(jsonr); err != nil {
		return err
	}

	batchdata := make([]byte, header.GetBatchTableBinaryByteLength())
	if _, err := reader.Read(batchdata); err != nil {
		return err
	}
	h.Data = make(map[string]interface{})
	for k, v := range h.Header {
		switch t := v.(type) {
		case BinaryBodyReference:
			h.Data[k] = getBatchTableValuesFromRef(&t, batchdata, k, batchLength)
		case []interface{}:
			h.Data[k] = t
		default:
			continue
		}
	}

	return nil
}

func (h *BatchTable) Write(writer io.Writer, header Header) error {
	outJSONHeader := make(map[string]interface{})
	var outBinaryBytes [][]byte
	var JSONLenght int
	offset := 0
	for k, v := range h.Header {
		switch t := v.(type) {
		case BinaryBodyReference:
			t.ByteOffset = offset
			outJSONHeader[k] = t.GetMap()
			bts := getBatchTableBinaryByte(&t, h.Data[k])
			offset += len(bts)
			outBinaryBytes = append(outBinaryBytes, bts)
		default:
			outJSONHeader[k] = v
		}
	}

	if bts, err := json.Marshal(outJSONHeader); err != nil {
		return err
	} else {
		n := len(bts)
		bts := createPaddingBytes(bts, n, 8, 0x20)
		if _, err := writer.Write(bts); err != nil {
			return err
		}
		JSONLenght = len(bts)
	}
	BinaryLenght := 0
	for i := range outBinaryBytes {
		BinaryLenght += len(outBinaryBytes[i])
		if _, err := writer.Write(outBinaryBytes[i]); err != nil {
			return err
		}
	}

	if header != nil {
		header.SetBatchTableJSONByteLength(uint32(JSONLenght))
		header.SetBatchTableBinaryByteLength(uint32(BinaryLenght))
	}
	return nil
}
