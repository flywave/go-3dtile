package tile3d

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

func TestGetBatchTableAllValuesFromRef(t *testing.T) {
	encode := func(v interface{}) []byte {
		var buf bytes.Buffer
		if err := binary.Write(&buf, littleEndian, v); err != nil {
			t.Fatal("encode failed:", err)
		}
		return buf.Bytes()
	}
	scalar := func(ct string) BinaryBodyReference {
		return BinaryBodyReference{ComponentType: ct, ContainerType: CONTAINER_TYPE_SCALAR}
	}

	cases := []struct {
		name       string
		ref        BinaryBodyReference
		buff       []byte
		batchCount int
		want       interface{}
	}{
		{
			name:       "BYTE SCALAR",
			ref:        scalar(COMPONENT_TYPE_BYTE),
			buff:       encode([]byte{0xFE, 0x02}),
			batchCount: 2,
			want:       []byte{0xFE, 0x02},
		},
		{
			name:       "UNSIGNED_BYTE SCALAR",
			ref:        scalar(COMPONENT_TYPE_UNSIGNED_BYTE),
			buff:       encode([]uint8{3, 7, 9}),
			batchCount: 3,
			want:       []uint8{3, 7, 9},
		},
		{
			name:       "SHORT SCALAR",
			ref:        scalar(COMPONENT_TYPE_SHORT),
			buff:       encode([]int16{-2, 300}),
			batchCount: 2,
			want:       []int16{-2, 300},
		},
		{
			name:       "UNSIGNED_SHORT SCALAR",
			ref:        scalar(COMPONENT_TYPE_UNSIGNED_SHORT),
			buff:       encode([]uint16{65535, 300}),
			batchCount: 2,
			want:       []uint16{65535, 300},
		},
		{
			name:       "INT SCALAR",
			ref:        scalar(COMPONENT_TYPE_INT),
			buff:       encode([]int32{-1, 70000}),
			batchCount: 2,
			want:       []int32{-1, 70000},
		},
		{
			name:       "UNSIGNED_INT SCALAR",
			ref:        scalar(COMPONENT_TYPE_UNSIGNED_INT),
			buff:       encode([]uint32{0xFFFFFFFF, 70000}),
			batchCount: 2,
			want:       []uint32{0xFFFFFFFF, 70000},
		},
		{
			name:       "FLOAT SCALAR",
			ref:        scalar(COMPONENT_TYPE_FLOAT),
			buff:       encode([]float32{1.5, -2.25}),
			batchCount: 2,
			want:       []float32{1.5, -2.25},
		},
		{
			name:       "DOUBLE SCALAR",
			ref:        scalar(COMPONENT_TYPE_DOUBLE),
			buff:       encode([]float64{1.5, -2.25}),
			batchCount: 2,
			want:       []float64{1.5, -2.25},
		},
		{
			name: "FLOAT VEC3 batch-major",
			ref: BinaryBodyReference{
				ComponentType: COMPONENT_TYPE_FLOAT,
				ContainerType: CONTAINER_TYPE_VEC3,
			},
			buff:       encode([]float32{1, 2, 3, 4, 5, 6}),
			batchCount: 2,
			want:       []float32{1, 2, 3, 4, 5, 6},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ref := c.ref
			got := getBatchTableAllValuesFromRef(&ref, c.buff, c.batchCount)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got %#v (%T), want %#v", got, got, c.want)
			}
		})
	}
}

func TestGetBatchTableAllValuesFromRefGuards(t *testing.T) {
	ref := &BinaryBodyReference{
		ComponentType: COMPONENT_TYPE_UNSIGNED_BYTE,
		ContainerType: CONTAINER_TYPE_SCALAR,
	}
	if v := getBatchTableAllValuesFromRef(nil, []byte{1, 2}, 2); v != nil {
		t.Errorf("nil ref = %#v, want nil", v)
	}
	if v := getBatchTableAllValuesFromRef(ref, []byte{1, 2}, 0); v != nil {
		t.Errorf("batchCount 0 = %#v, want nil", v)
	}
	if v := getBatchTableAllValuesFromRef(ref, []byte{1, 2}, 3); v != nil {
		t.Errorf("batchCount beyond buffer = %#v, want nil (no panic)", v)
	}
	offRef := &BinaryBodyReference{
		ByteOffset:    3,
		ComponentType: COMPONENT_TYPE_UNSIGNED_BYTE,
		ContainerType: CONTAINER_TYPE_SCALAR,
	}
	if v := getBatchTableAllValuesFromRef(offRef, []byte{1, 2, 3, 4}, 2); v != nil {
		t.Errorf("byteOffset beyond buffer = %#v, want nil (no panic)", v)
	}
	bad := &BinaryBodyReference{ComponentType: "BOGUS", ContainerType: CONTAINER_TYPE_SCALAR}
	if v := getBatchTableAllValuesFromRef(bad, []byte{1, 2}, 2); v != nil {
		t.Errorf("unknown componentType = %#v, want nil", v)
	}
}
