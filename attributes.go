package potree

type AttributeType uint32

const (
	ATTR_INT8      AttributeType = 0
	ATTR_INT16     AttributeType = 1
	ATTR_INT32     AttributeType = 2
	ATTR_INT64     AttributeType = 3
	ATTR_UINT8     AttributeType = 10
	ATTR_UINT16    AttributeType = 11
	ATTR_UINT32    AttributeType = 12
	ATTR_UINT64    AttributeType = 13
	ATTR_FLOAT     AttributeType = 20
	ATTR_DOUBLE    AttributeType = 21
	ATTR_UNDEFINED AttributeType = 123456
)

var (
	AttributeTypeSize = map[AttributeType]int{
		ATTR_UNDEFINED: 0,
		ATTR_UINT8:     1,
		ATTR_UINT16:    2,
		ATTR_UINT32:    4,
		ATTR_UINT64:    8,
		ATTR_INT8:      1,
		ATTR_INT16:     2,
		ATTR_INT32:     4,
		ATTR_INT64:     8,
		ATTR_FLOAT:     4,
		ATTR_DOUBLE:    8,
	}

	AttributeTypeName = map[AttributeType]string{
		ATTR_INT8:      "int8",
		ATTR_INT16:     "int16",
		ATTR_INT32:     "int32",
		ATTR_INT64:     "int64",
		ATTR_UINT8:     "uint8",
		ATTR_UINT16:    "uint16",
		ATTR_UINT32:    "uint32",
		ATTR_UINT64:    "uint64",
		ATTR_FLOAT:     "float",
		ATTR_DOUBLE:    "double",
		ATTR_UNDEFINED: "undefined",
	}
)

func TypenameToType(name string) AttributeType {
	if name == "int8" {
		return ATTR_INT8
	} else if name == "int16" {
		return ATTR_INT16
	} else if name == "int32" {
		return ATTR_INT32
	} else if name == "int64" {
		return ATTR_INT64
	} else if name == "uint8" {
		return ATTR_UINT8
	} else if name == "uint16" {
		return ATTR_UINT16
	} else if name == "uint32" {
		return ATTR_UINT32
	} else if name == "uint64" {
		return ATTR_UINT64
	} else if name == "float" {
		return ATTR_FLOAT
	} else if name == "double" {
		return ATTR_DOUBLE
	} else if name == "undefined" {
		return ATTR_UNDEFINED
	}
	return ATTR_UNDEFINED
}

type Attribute struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Size        int    `json:"size"`
	NumElements int    `json:"numElements"`
	ElementSize int    `json:"elementSize"`
	Type        string `json:"type"`
}

func NewAttribute(name string, size, numElements, elementSize int, type_ AttributeType) *Attribute {
	attr := &Attribute{Name: name, Size: size, NumElements: numElements, ElementSize: elementSize, Type: AttributeTypeName[type_]}
	return attr
}

func (a *Attribute) GetType() AttributeType {
	return TypenameToType(a.Type)
}
