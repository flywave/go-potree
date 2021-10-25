package potree

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"
)

type NodeType uint8

const (
	NT_NORMAL NodeType = 0
	NT_LEAF   NodeType = 1
	NT_PROXY  NodeType = 2
)

var (
	POTREE_BYTEORDER = binary.LittleEndian
)

type node struct {
	Type       NodeType
	ChildMask  uint8
	NumPoints  uint32
	ByteOffset int64
	ByteSize   int64
	AttrStart  int
	AttrEnd    int
}

func (n *node) readNode(reader io.Reader) error {
	return binary.Read(reader, POTREE_BYTEORDER, n)
}

func (n *node) writeNode(writer io.Writer) error {
	return binary.Write(writer, POTREE_BYTEORDER, *n)
}

func (n *node) size() int64 {
	return n.ByteSize
}

func (n *node) read(reader io.ReadSeeker) ([]byte, error) {
	_, err := reader.Seek(n.ByteOffset, io.SeekStart)
	if err != nil {
		return nil, err
	}
	ret := make([]byte, n.ByteSize)
	_, err = reader.Read(ret)
	return ret, err
}

func (n *node) write(off int64, buf []byte, writer io.Writer) error {
	n.ByteOffset = off
	si, err := writer.Write(buf)
	if err != nil {
		return err
	}
	n.ByteSize = int64(si)
	return nil
}

func (n *node) compact(attributes []Attribute, isBrotliEncoded bool) []byte {
	buf := &bytes.Buffer{}
	for i := range attributes {
		attributes[i].pack(isBrotliEncoded)
		buf.Write(attributes[i].Buffer)
	}
	return buf.Bytes()
}

func (n *node) uncompact(data []byte, attributes []Attribute, isBrotliEncoded bool) {
	offset := 0
	attrOff := n.AttrStart
	for {
		size := int(n.NumPoints) * attributes[attrOff].Size
		attributes[attrOff].Buffer = data[offset : offset+size]
		attributes[attrOff].unpack(isBrotliEncoded)
		offset += size
		attrOff++
		if offset == len(data) {
			n.AttrEnd = attrOff
			return
		}
	}
}

func (n *node) compress(attributes []Attribute) []byte {
	uncomress := n.compact(attributes, true)
	ctx := &Brotli{}
	ret := ctx.Encode(nil, uncomress)
	return ret
}

func (n *node) uncompress(data []byte, attributes []Attribute) {
	ctx := &Brotli{}
	uncomress := ctx.Decode(nil, data)
	n.uncompact(uncomress, attributes, true)
}

type Node struct {
	node
	Box      AABB
	Name     string
	Parent   *Node
	Childs   [8]*Node
	Buffer   []byte
	genProxy bool
}

func (n *Node) Level() int {
	return len(n.Name) - 1
}

func (n *Node) IsLeaf() bool {
	return n.Type == NT_LEAF
}

func (n *Node) Traverse(callback func(*Node) bool) bool {
	if !callback(n) {
		return false
	}

	for _, child := range n.Childs {
		if child != nil {
			if !child.Traverse(callback) {
				return false
			}
		}
	}

	return true
}

func ChildMaskOf(node *Node) uint8 {
	mask := uint8(0)

	for i := 0; i < 8; i++ {
		child := node.Childs[i]

		if child != nil {
			mask = mask | (1 << i)
		}
	}

	return mask
}

func CaclNodeID(level, gridSize int, x, y, z int64) string {
	id := "r"

	currentGridSize := gridSize
	lx := x
	ly := y
	lz := z

	for i := 0; i < level; i++ {
		index := 0

		if lx >= int64(currentGridSize/2) {
			index = index + 0b100
			lx = lx - int64(currentGridSize/2)
		}

		if ly >= int64(currentGridSize/2) {
			index = index + 0b010
			ly = ly - int64(currentGridSize/2)
		}

		if lz >= int64(currentGridSize/2) {
			index = index + 0b001
			lz = lz - int64(currentGridSize/2)
		}

		id = id + strconv.Itoa(index)
		currentGridSize = currentGridSize / 2
	}

	return id
}
