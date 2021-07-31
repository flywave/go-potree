package potree

type NodeType uint8

const (
	NORMAL NodeType = 0
	LEAF   NodeType = 1
	PROXY  NodeType = 2
)

type Node struct {
	ByteOffset uint64
	ByteSize   uint64
	NumPoints  uint32
	Type       NodeType
	ChildMask  uint8
	Name       string
	Parent     *Node
	Childs     [8]*Node
}

func (n *Node) Level() int {
	return len(n.Name) - 1
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
