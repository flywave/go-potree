package potree

type Node struct {
	ByteOffset uint64
	ByteSize   uint64
	NumPoints  uint32
	NodeType   uint8
	ChildMask  uint8
	Name       string
	Childs     map[string]*Node
}
