package potree

import (
	"bytes"
	"encoding/json"
	"io"
)

const (
	POTREE_VERSION = "2.0"
)

type Metadata struct {
	Version         string      `json:"version"`
	Points          *int64      `json:"points,omitempty"`
	PointsProcessed *int64      `json:"pointsProcessed,omitempty"`
	NodesProcessed  *int64      `json:"nodesProcessed,omitempty"`
	DurationMS      *float64    `json:"durationMS,omitempty"`
	BoundingBox     AABB        `json:"boundingBox"`
	Attrs           []Attribute `json:"attributes"`
	BytesPerPoint   int         `json:"bytesPerPoint"`
	Scale           [3]float64  `json:"scale"`
	Spacing         *float64    `json:"spacing,omitempty"`
	Offset          *[3]float64 `json:"offset,omitempty"`
	Hierarchy       *Hierarchy  `json:"hierarchy,omitempty"`
	Encoding        *string     `json:"encoding,omitempty"`
	Projection      *string     `json:"projection,omitempty"`
}

func NewMetadata(attributes []Attribute) *Metadata {
	ret := &Metadata{Version: POTREE_VERSION}
	ret.Attrs = attributes

	for _, attribute := range attributes {
		ret.BytesPerPoint += attribute.Size
	}
	return ret
}

func (l *Metadata) readMetadata(data io.ReadSeeker, jsonLength int) error {
	jdata := make([]byte, jsonLength)
	_, err := data.Read(jdata)
	dec := json.NewDecoder(bytes.NewBuffer(jdata))
	if err != nil {
		return nil
	}
	if err := dec.Decode(l); err != nil {
		return err
	}
	return nil
}

func (l *Metadata) writeMetadata(wr io.Writer) (int, error) {
	var jdata []byte
	buf := bytes.NewBuffer(jdata)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(l); err != nil {
		return 0, err
	}
	jdata = buf.Bytes()
	n, err := wr.Write(jdata)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (l *Metadata) Add(attribute *Attribute) {
	l.Attrs = append(l.Attrs, *attribute)
	l.BytesPerPoint += attribute.Size
}

func (l *Metadata) GetOffset(name string) int {
	offset := 0

	for _, attribute := range l.Attrs {
		if attribute.Name == name {
			return offset
		}

		offset += attribute.Size
	}

	return -1
}

func (l *Metadata) Get(name string) *Attribute {
	for _, attribute := range l.Attrs {
		if attribute.Name == name {
			return &attribute
		}
	}
	return nil
}

func (l *Metadata) IsBrotliEncoded() bool {
	if l.Encoding != nil {
		return *l.Encoding == "BROTLI"
	}
	return false
}
