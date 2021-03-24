package tile3d

import (
	"encoding/json"
	"errors"
)

const (
	COMPONENT_TYPE_BYTE           = "BYTE"
	COMPONENT_TYPE_UNSIGNED_BYTE  = "UNSIGNED_BYTE"
	COMPONENT_TYPE_SHORT          = "SHORT"
	COMPONENT_TYPE_UNSIGNED_SHORT = "UNSIGNED_SHORT"
	COMPONENT_TYPE_INT            = "INT"
	COMPONENT_TYPE_UNSIGNED_INT   = "UNSIGNED_INT"
	COMPONENT_TYPE_FLOAT          = "FLOAT"
	COMPONENT_TYPE_DOUBLE         = "DOUBLE"
)

const (
	CONTAINER_TYPE_SCALAR = "SCALAR"
	CONTAINER_TYPE_VEC2   = "VEC2"
	CONTAINER_TYPE_VEC3   = "VEC3"
	CONTAINER_TYPE_VEC4   = "VEC4"
)

const (
	REF_PROP_BYTE_OFFSET    = "byteOffset"
	REF_PROP_COMPONENT_TYPE = "componentType"
	REF_PROP_TYPE           = "type"
)

func ComponentTypeSize(tp string) int {
	switch tp {
	case COMPONENT_TYPE_BYTE:
		return 1
	case COMPONENT_TYPE_UNSIGNED_BYTE:
		return 1
	case COMPONENT_TYPE_SHORT:
		return 2
	case COMPONENT_TYPE_UNSIGNED_SHORT:
		return 2
	case COMPONENT_TYPE_INT:
		return 4
	case COMPONENT_TYPE_UNSIGNED_INT:
		return 4
	case COMPONENT_TYPE_FLOAT:
		return 4
	case COMPONENT_TYPE_DOUBLE:
		return 8
	default:
		return 0
	}
}

func ContainerTypeSize(tp string) int {
	switch tp {
	case CONTAINER_TYPE_SCALAR:
		return 1
	case CONTAINER_TYPE_VEC2:
		return 2
	case CONTAINER_TYPE_VEC3:
		return 3
	case CONTAINER_TYPE_VEC4:
		return 4
	default:
		return 0
	}
}

type BinaryBodyReference struct {
	ByteOffset    uint32 `json:"byteOffset"`
	ComponentType string `json:"componentType,omitempty"`
	ContainerType string `json:"type,omitempty"`
}

func (r *BinaryBodyReference) GetMap() map[string]interface{} {
	ret := make(map[string]interface{})
	ret[REF_PROP_BYTE_OFFSET] = r.ByteOffset
	if len(r.ComponentType) > 0 {
		ret[REF_PROP_COMPONENT_TYPE] = r.ComponentType
	}
	if len(r.ContainerType) > 0 {
		ret[REF_PROP_TYPE] = r.ContainerType
	}
	return ret
}

func (r *BinaryBodyReference) FromMap(d map[string]interface{}) {
	if d[REF_PROP_BYTE_OFFSET] != nil {
		r.ByteOffset = uint32(d[REF_PROP_BYTE_OFFSET].(float64))
	}
	if d[REF_PROP_COMPONENT_TYPE] != nil {
		r.ComponentType = d[REF_PROP_COMPONENT_TYPE].(string)
	}
	if d[REF_PROP_TYPE] != nil {
		r.ContainerType = d[REF_PROP_TYPE].(string)
	}
}

func createReference(offset uint32, componentType *string, containerType *string) map[string]interface{} {
	reference := make(map[string]interface{})
	reference[REF_PROP_BYTE_OFFSET] = offset
	if componentType != nil {
		reference[REF_PROP_COMPONENT_TYPE] = *componentType
	}
	if containerType != nil {
		reference[REF_PROP_TYPE] = *containerType
	}
	return reference
}

func addReference(jsonHeader *map[string]interface{}, property string, offset uint32, componentType string, containerType string, feature bool) error {
	var reference map[string]interface{}
	if feature {
		if "BATCH_ID" == property {
			if containerType != CONTAINER_TYPE_SCALAR {
				return errors.New("Invalid container type for BATCH_ID: " + containerType + ".")
			}
			if componentType != COMPONENT_TYPE_UNSIGNED_BYTE && componentType != COMPONENT_TYPE_UNSIGNED_SHORT && componentType != COMPONENT_TYPE_UNSIGNED_INT {
				return errors.New("Invalid component type for BATCH_ID: " + componentType + ".")
			}
			if componentType != COMPONENT_TYPE_UNSIGNED_SHORT {
				reference = createReference(offset, &componentType, nil)
			} else {
				reference = createReference(offset, nil, nil)
			}
		} else {
			reference = createReference(offset, nil, nil)
		}
	} else {
		reference = createReference(offset, &componentType, &containerType)
	}
	(*jsonHeader)[property] = reference
	return nil
}

type BinaryBodySizeHelper struct {
	Header *map[string]interface{}
	Size   uint32
}

func (b *BinaryBodySizeHelper) addProperty(property string, data interface{}, componentType string, containerType string, feature bool) {
	width := ComponentTypeSize(componentType)
	b.Size += calcPadding(b.Size, uint32(width))
	addReference(b.Header, property, b.Size, componentType, containerType, feature)
	unit := ContainerTypeSize(containerType)
	b.Size += uint32(width) * uint32(unit) * uint32(len(data.([]interface{})))
}

func (b *BinaryBodySizeHelper) finished() {
	b.Size += calcPadding(b.Size, 8)
}

func maxBatchId(array []int64) int64 {
	maxElement := array[0]
	for _, element := range array {
		if maxElement > element {
			maxElement = element
		}
	}
	return maxElement
}

func (b *BinaryBodySizeHelper) addBatchId(array []int64) {
	max := maxBatchId(array)
	if max > 0xFFFF {
		b.addProperty("BATCH_ID", array, COMPONENT_TYPE_UNSIGNED_INT, CONTAINER_TYPE_SCALAR, true)
	} else if max > 0xFF {
		b.addProperty("BATCH_ID", array, COMPONENT_TYPE_UNSIGNED_SHORT, CONTAINER_TYPE_SCALAR, true)
	} else {
		b.addProperty("BATCH_ID", array, COMPONENT_TYPE_UNSIGNED_BYTE, CONTAINER_TYPE_SCALAR, true)
	}
}

func (b *BinaryBodySizeHelper) calcHeaderSize(offset uint32) uint32 {
	var err error
	var bs []byte
	if bs, err = json.Marshal(b.Header); err != nil {
		return 0
	}
	length := uint32(len(bs))
	return length + calcPadding(offset+length, 8)
}
