package potree

import "io"

type NodeType uint8

const (
	NT_NORMAL NodeType = 0
	NT_LEAF   NodeType = 1
	NT_PROXY  NodeType = 2
)

type node struct {
	Type       NodeType
	ChildMask  uint8
	NumPoints  uint32
	ByteOffset int64
	ByteSize   int64
}

func (n *node) CalcSize() int64 {
	return 0
}

func (n *node) Read(reader io.ReadSeeker) error {
	return nil
}

func (n *node) Write(writer io.Writer) error {
	return nil
}

type Node struct {
	node
	Name   string
	Parent *Node
	Childs [8]*Node
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
