package potree

import (
	"math"
	"reflect"
	"unsafe"
)

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
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Size        int         `json:"size"`
	NumElements int         `json:"numElements"`
	ElementSize int         `json:"elementSize"`
	Type        string      `json:"type"`
	Min         []float64   `json:"min,omitempty"`
	Max         []float64   `json:"max,omitempty"`
	Buffer      []byte      `json:"-"`
	Data        interface{} `json:"-"`
}

func NewAttribute(name string, size, numElements, elementSize int, type_ AttributeType) *Attribute {
	attr := &Attribute{Name: name, Size: size, NumElements: numElements, ElementSize: elementSize, Type: AttributeTypeName[type_]}
	return attr
}

func (a *Attribute) GetType() AttributeType {
	return TypenameToType(a.Type)
}

func unsafeCopy(data unsafe.Pointer, dst []byte) {
	var bufSlice []byte
	bufHeader := (*reflect.SliceHeader)((unsafe.Pointer(&bufSlice)))
	bufHeader.Cap = int(len(dst))
	bufHeader.Len = int(len(dst))
	bufHeader.Data = uintptr(data)

	copy(dst, bufSlice)
}

func unsafeCopyDst(src []byte, dst unsafe.Pointer) {
	var bufSlice []byte
	bufHeader := (*reflect.SliceHeader)((unsafe.Pointer(&bufSlice)))
	bufHeader.Cap = int(len(src))
	bufHeader.Len = int(len(src))
	bufHeader.Data = uintptr(dst)

	copy(bufSlice, src)
}

func (a *Attribute) unpack(isBrotliEncoded bool) {
	tp := TypenameToType(a.Type)
	elsize := len(a.Buffer) / a.ElementSize
	switch tp {
	case ATTR_INT8:
		rawdata := make([]int8, elsize)
		unsafeCopyDst(a.Buffer, unsafe.Pointer(&rawdata[0]))
	case ATTR_INT16:
		rawdata := make([]int16, elsize)
		unsafeCopyDst(a.Buffer, unsafe.Pointer(&rawdata[0]))
	case ATTR_INT32:
		rawdata := make([]int32, elsize)
		unsafeCopyDst(a.Buffer, unsafe.Pointer(&rawdata[0]))
	case ATTR_INT64:
		rawdata := make([]int64, elsize)
		unsafeCopyDst(a.Buffer, unsafe.Pointer(&rawdata[0]))
	case ATTR_UINT8:
		rawdata := make([]uint8, elsize)
		unsafeCopyDst(a.Buffer, unsafe.Pointer(&rawdata[0]))
	case ATTR_UINT16:
		rawdata := make([]uint16, elsize)
		unsafeCopyDst(a.Buffer, unsafe.Pointer(&rawdata[0]))
	case ATTR_UINT32:
		rawdata := make([]uint32, elsize)
		unsafeCopyDst(a.Buffer, unsafe.Pointer(&rawdata[0]))
	case ATTR_UINT64:
		rawdata := make([]uint64, elsize)
		unsafeCopyDst(a.Buffer, unsafe.Pointer(&rawdata[0]))
	case ATTR_FLOAT:
		rawdata := make([]float32, elsize)
		unsafeCopyDst(a.Buffer, unsafe.Pointer(&rawdata[0]))
	case ATTR_DOUBLE:
		rawdata := make([]float64, elsize)
		unsafeCopyDst(a.Buffer, unsafe.Pointer(&rawdata[0]))
	case ATTR_UNDEFINED:
	}

	if a.Name == "position_morton" && isBrotliEncoded {
		raw := a.Data.([]uint64)

		si := len(raw) / a.Size

		var pos []uint32
		posHeader := (*reflect.SliceHeader)((unsafe.Pointer(&pos)))
		posHeader.Cap = int(len(raw) * 2)
		posHeader.Len = int(len(raw) * 2)
		posHeader.Data = uintptr(unsafe.Pointer(&raw[0]))

		newPos := make([]int32, si*3)

		for i := 0; i < si; i++ {
			mc_0, mc_1, mc_2, mc_3 := uint32(pos[i*4]), uint32(pos[i*4+1]), uint32(pos[i*4+2]), uint32(pos[i*4+3])

			X := dealign24b((mc_3&0x00FFFFFF)>>0) | (dealign24b(((mc_3>>24)|(mc_2<<8))>>0) << 8)
			Y := dealign24b((mc_3&0x00FFFFFF)>>1) | (dealign24b(((mc_3>>24)|(mc_2<<8))>>1) << 8)
			Z := dealign24b((mc_3&0x00FFFFFF)>>2) | (dealign24b(((mc_3>>24)|(mc_2<<8))>>2) << 8)

			if mc_1 != 0 || mc_2 != 0 {
				X = X | (dealign24b((mc_1&0x00FFFFFF)>>0) << 16) | (dealign24b(((mc_1>>24)|(mc_0<<8))>>0) << 24)
				Y = Y | (dealign24b((mc_1&0x00FFFFFF)>>1) << 16) | (dealign24b(((mc_1>>24)|(mc_0<<8))>>1) << 24)
				Z = Z | (dealign24b((mc_1&0x00FFFFFF)>>2) << 16) | (dealign24b(((mc_1>>24)|(mc_0<<8))>>2) << 24)
			}

			newPos[i*3] = int32(X)
			newPos[i*3+1] = int32(Y)
			newPos[i*3+2] = int32(Z)
		}

		a.Name = "position"
		a.Type = "int32"
		a.NumElements = 3
		a.ElementSize = 4
		a.Size = 12
		a.Data = newPos
	}

	if a.Name == "rgb_morton" && isBrotliEncoded {
		raw := a.Data.([]uint64)

		si := len(raw) / a.Size

		var pos []uint32
		posHeader := (*reflect.SliceHeader)((unsafe.Pointer(&pos)))
		posHeader.Cap = int(len(raw) * 2)
		posHeader.Len = int(len(raw) * 2)
		posHeader.Data = uintptr(unsafe.Pointer(&raw[0]))

		newColor := make([]uint16, si*3)

		for i := 0; i < si; i++ {
			mc_0, mc_1 := uint32(pos[i*2]), uint32(pos[i*2+1])

			r := dealign24b((mc_1&0x00FFFFFF)>>0) | (dealign24b(((mc_1>>24)|(mc_0<<8))>>0) << 8)
			g := dealign24b((mc_1&0x00FFFFFF)>>1) | (dealign24b(((mc_1>>24)|(mc_0<<8))>>1) << 8)
			b := dealign24b((mc_1&0x00FFFFFF)>>2) | (dealign24b(((mc_1>>24)|(mc_0<<8))>>2) << 8)

			newColor[i*3] = uint16(r)
			newColor[i*3+1] = uint16(g)
			newColor[i*3+2] = uint16(b)
		}

		a.Name = "rgb"
		a.Type = "uint16"
		a.NumElements = 3
		a.ElementSize = 2
		a.Size = 6
		a.Data = newColor
	}
}

func (a *Attribute) pack(isBrotliEncoded bool) {
	tp := TypenameToType(a.Type)
	if a.Data == nil {
		return
	}
	numPoints := 0
	switch tp {
	case ATTR_INT8:
		rawdata, ok := a.Data.([]int8)
		if !ok {
			return
		}
		numPoints := len(rawdata)
		a.Buffer = make([]byte, a.Size*numPoints)
		unsafeCopy(unsafe.Pointer(&rawdata[0]), a.Buffer)
	case ATTR_INT16:
		rawdata, ok := a.Data.([]int16)
		if !ok {
			return
		}
		numPoints := len(rawdata)
		a.Buffer = make([]byte, a.Size*numPoints)
		unsafeCopy(unsafe.Pointer(&rawdata[0]), a.Buffer)
	case ATTR_INT32:
		rawdata, ok := a.Data.([]int32)
		if !ok {
			return
		}
		numPoints := len(rawdata)
		a.Buffer = make([]byte, a.Size*numPoints)
		unsafeCopy(unsafe.Pointer(&rawdata[0]), a.Buffer)
	case ATTR_INT64:
		rawdata, ok := a.Data.([]int64)
		if !ok {
			return
		}
		numPoints := len(rawdata)
		a.Buffer = make([]byte, a.Size*numPoints)
		unsafeCopy(unsafe.Pointer(&rawdata[0]), a.Buffer)
	case ATTR_UINT8:
		rawdata, ok := a.Data.([]uint8)
		if !ok {
			return
		}
		numPoints := len(rawdata)
		a.Buffer = make([]byte, a.Size*numPoints)
		unsafeCopy(unsafe.Pointer(&rawdata[0]), a.Buffer)
	case ATTR_UINT16:
		rawdata, ok := a.Data.([]uint16)
		if !ok {
			return
		}
		numPoints := len(rawdata)
		a.Buffer = make([]byte, a.Size*numPoints)
		unsafeCopy(unsafe.Pointer(&rawdata[0]), a.Buffer)
	case ATTR_UINT32:
		rawdata, ok := a.Data.([]uint32)
		if !ok {
			return
		}
		numPoints := len(rawdata)
		a.Buffer = make([]byte, a.Size*numPoints)
		unsafeCopy(unsafe.Pointer(&rawdata[0]), a.Buffer)
	case ATTR_UINT64:
		rawdata, ok := a.Data.([]uint64)
		if !ok {
			return
		}
		numPoints := len(rawdata)
		a.Buffer = make([]byte, a.Size*numPoints)
		unsafeCopy(unsafe.Pointer(&rawdata[0]), a.Buffer)
	case ATTR_FLOAT:
		rawdata, ok := a.Data.([]float32)
		if !ok {
			return
		}
		numPoints := len(rawdata)
		a.Buffer = make([]byte, a.Size*numPoints)
		unsafeCopy(unsafe.Pointer(&rawdata[0]), a.Buffer)
	case ATTR_DOUBLE:
		rawdata, ok := a.Data.([]float64)
		if !ok {
			return
		}
		numPoints := len(rawdata)
		a.Buffer = make([]byte, a.Size*numPoints)
		unsafeCopy(unsafe.Pointer(&rawdata[0]), a.Buffer)
	case ATTR_UNDEFINED:
	}

	if a.Name == "position" && isBrotliEncoded {
		rawdata, _ := a.Data.([]int32)

		mcs := []MortonCode{}

		type P struct {
			x, y, z int32
		}
		ps := []P{}
		min := P{}
		min.x = math.MaxInt32
		min.y = math.MaxInt32
		min.z = math.MaxInt32

		for i := 0; i < numPoints; i++ {

			var p P
			p.x = rawdata[i*3]
			p.y = rawdata[i*3+1]
			p.z = rawdata[i*3+2]

			min.x = MinInt(min.x, p.x)
			min.y = MinInt(min.y, p.y)
			min.z = MinInt(min.z, p.z)

			ps = append(ps, p)
		}

		i := uint64(0)
		for _, p := range ps {

			mx := uint32(p.x - min.x)
			my := uint32(p.y - min.y)
			mz := uint32(p.z - min.z)

			mx_l := uint32(mx & 0x0000ffff)
			my_l := uint32(my & 0x0000ffff)
			mz_l := uint32(mz & 0x0000ffff)

			mx_h := mx >> 16
			my_h := my >> 16
			mz_h := mz >> 16

			mc_l := MortonEncodeMagicBits(mx_l, my_l, mz_l)
			mc_h := MortonEncodeMagicBits(mx_h, my_h, mz_h)

			var mc MortonCode
			mc.lower = mc_l
			mc.upper = mc_h
			mc.whatever = MortonEncodeMagicBits(mx, my, mz)
			mc.index = i

			mcs = append(mcs, mc)

			i++
		}

		pos := make([]uint64, numPoints*2)

		for i := 0; i < numPoints; i++ {
			mc := mcs[i]

			pos = append(pos, mc.upper)
			pos = append(pos, mc.lower)
		}
		a.Name = "position_morton"
		a.Type = "uint64"
		a.NumElements = 2
		a.ElementSize = 8
		a.Size = 16
		a.Data = pos
	}

	if a.Name == "rgb" && isBrotliEncoded {
		rawdata, _ := a.Data.([]uint16)

		color := make([]uint64, numPoints)

		for i := 0; i < numPoints; i++ {
			r, g, b := rawdata[i*3], rawdata[i*3+1], rawdata[i*3+2]
			mc := MortonEncodeMagicBits(uint32(r), uint32(g), uint32(b))
			color = append(color, mc)
		}

		a.Name = "rgb_morton"
		a.Type = "uint64"
		a.NumElements = 1
		a.ElementSize = 8
		a.Size = 8
		a.Data = color
	}
}

var (
	POSITION                   = Attribute{Name: "position", Type: "int32", NumElements: 3, ElementSize: 4, Size: 12}
	POSITION_MORTON            = Attribute{Name: "position_morton", Type: "uint64", NumElements: 2, ElementSize: 8, Size: 16}
	COLOR                      = Attribute{Name: "rgb", Type: "uint16", NumElements: 3, ElementSize: 2, Size: 6}
	COLOR_MORTON               = Attribute{Name: "rgb_morton", Type: "uint64", NumElements: 1, ElementSize: 8, Size: 8}
	INTENSITY                  = Attribute{Name: "intensity", Type: "uint16", NumElements: 1, ElementSize: 2, Size: 2}
	CLASSIFICATION             = Attribute{Name: "classification", Type: "uint8", NumElements: 1, ElementSize: 1, Size: 1}
	RETURN_NUMBER              = Attribute{Name: "return number", Type: "uint8", NumElements: 1, ElementSize: 1, Size: 1}
	NUMBER_OF_RETURNS          = Attribute{Name: "number of returns", Type: "uint8", NumElements: 1, ElementSize: 1, Size: 1}
	POINT_SOURCE_ID            = Attribute{Name: "point source id", Type: "uint16", NumElements: 1, ElementSize: 2, Size: 2}
	GPS_TIME                   = Attribute{Name: "gps-time", Type: "double", NumElements: 1, ElementSize: 8, Size: 8}
	NORMAL                     = Attribute{Name: "normal", Type: "float", NumElements: 3, ElementSize: 4, Size: 12}
	NORMAL_X                   = Attribute{Name: "NormalX", Type: "double", NumElements: 1, ElementSize: 8, Size: 8}
	NORMAL_Y                   = Attribute{Name: "NormalY", Type: "double", NumElements: 1, ElementSize: 8, Size: 8}
	NORMAL_Z                   = Attribute{Name: "NormalZ", Type: "double", NumElements: 1, ElementSize: 8, Size: 8}
	SCAN_ANGLE                 = Attribute{Name: "scan angle", Type: "int16", NumElements: 1, ElementSize: 2, Size: 2}
	SCAN_ANGLE_RANK            = Attribute{Name: "scan angle rank", Type: "uint8", NumElements: 1, ElementSize: 1, Size: 1}
	USER_DATA                  = Attribute{Name: "user data", Type: "uint8", NumElements: 1, ElementSize: 1, Size: 1}
	CLASSIFICATION_FLAGS       = Attribute{Name: "classification flags", Type: "uint8", NumElements: 1, ElementSize: 1, Size: 1}
	POSITION_PROJECTED_PROFILE = Attribute{Name: "position_projected_profile", Type: "int32", NumElements: 2, ElementSize: 4, Size: 8}
)
