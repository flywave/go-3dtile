package tile3d

import (
	"bytes"
	"encoding/json"
	"io"
)

const (
	BATCH_TABLE_HIERARCHY     = "3DTILES_batch_table_hierarchy"
	HIERARCHY_INSTANCE_LENGTH = "instancesLength"

	HIERARCHY_CLASSES           = "classes"
	HIERARCHY_CLASSIDS          = "classIds"
	HIERARCHY_CLASSE_NAME       = "name"
	HIERARCHY_CLASSINDEXES      = "classIndexes"
	HIERARCHY_CLASSES_INSTANCES = "instances"
	HIERARCHY_CLASSES_LENGTH    = "length"

	HIERARCHY_PARENT_COUNTS = "parentCounts"
	HIERARCHY_PARENTIDS     = "parentIds"
	HIERARCHY_PARENTINDEXES = "parentIndexes"
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
	return int64(w.GetSize())
}

func (h *BatchTable) CalcSize(header Header) int64 {
	w := newSizeWriter()
	if err := h.Write(w.writer, header); err != nil {
		return 0
	}
	return int64(w.GetSize())
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
	jsonLen := header.GetBatchTableJSONByteLength()

	if jsonLen == 0 {
		h.Data = make(map[string]interface{})
		h.Header = make(map[string]interface{})
		return nil
	}

	jsonb := make([]byte, jsonLen)
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
	var outBinaryBytes [][]byte
	JSONLenght := 0
	offset := 0
	outJSONHeader := make(map[string]interface{})
	for k, v := range h.Header {
		switch t := v.(type) {
		case BinaryBodyReference:
			t.ByteOffset = uint32(offset)
			outJSONHeader[k] = t.GetMap()
			bts := getBatchTableBinaryByte(&t, h.Data[k])
			offset += len(bts)
			outBinaryBytes = append(outBinaryBytes, bts)
		default:
			outJSONHeader[k] = v
		}
	}
	var BinaryLenght int
	if len(outJSONHeader) > 0 {
		if bts, err := json.Marshal(outJSONHeader); err != nil {
			return err
		} else {
			n := len(bts)
			bts := createPaddingBytes(bts, uint32(n), 8, 0x20)
			if _, err := writer.Write(bts); err != nil {
				return err
			}
			JSONLenght = len(bts)
		}
		for i := range outBinaryBytes {
			BinaryLenght += len(outBinaryBytes[i])
			if _, err := writer.Write(outBinaryBytes[i]); err != nil {
				return err
			}
		}
	}

	if header != nil {
		header.SetBatchTableJSONByteLength(uint32(JSONLenght))
		header.SetBatchTableBinaryByteLength(uint32(BinaryLenght))
	}
	return nil
}
