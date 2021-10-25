package potree

import (
	"errors"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
)

const (
	MaxPointsPerChunk = 10000
	MetadataName      = "metadata.json"
	HierarchyName     = "hierarchy.bin"
	OctreeName        = "octree.bin"
	HierarchyStepSize = 4
	BytesPerNode      = 22
)

type chunknodelist []*Node

func (s chunknodelist) Len() int { return len(s) }

func (s chunknodelist) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s chunknodelist) Less(i, j int) bool {
	if len(s[i].Name) != len(s[j].Name) {
		return len(s[i].Name) < len(s[j].Name)
	} else {
		return s[i].Name < s[j].Name
	}
}

type hierarchyChunk struct {
	name  string
	nodes chunknodelist
}

func (c *hierarchyChunk) chunkSize() int64 {
	return int64(len(c.nodes) * BytesPerNode)
}

func (c *hierarchyChunk) sortNodes() {
	sort.Sort(c.nodes)
}

type PotreeArchive struct {
	path         string
	root         *Node
	nodeMaps     map[string]*Node
	metadata     *Metadata
	octree       *os.File
	octreeOffset int64
}

func NewArchive(path string) *PotreeArchive {
	return &PotreeArchive{path: path}
}

func (b *PotreeArchive) Load() error {
	err := b.readMetadata()
	if err != nil {
		return err
	}
	err = b.readHierarchy()
	if err != nil {
		return err
	}
	for _, n := range b.nodeMaps {
		err = b.unpackNode(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *PotreeArchive) Save() error {
	b.octreeOffset = 0
	err := b.writeOctree(b.root)
	if err != nil {
		return err
	}
	err = b.writeHierarchy()
	if err != nil {
		return err
	}
	_, err = b.writeMetadata()
	if err != nil {
		return err
	}
	return nil
}

func (b *PotreeArchive) getMetadataPath() string {
	return path.Join(b.path, MetadataName)
}

func (b *PotreeArchive) getHierarchyPath() string {
	return path.Join(b.path, HierarchyName)
}

func (b *PotreeArchive) getOctreePath() string {
	return path.Join(b.path, OctreeName)
}

func (b *PotreeArchive) readMetadata() error {
	p := b.getMetadataPath()
	if !FileExists(p) {
		return errors.New("metadata.json not found!")
	}
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	return b.metadata.readMetadata(f)
}

func (b *PotreeArchive) writeMetadata() (int, error) {
	p := b.getMetadataPath()
	f, err := os.Open(p)
	if err != nil {
		return -1, err
	}
	return b.metadata.writeMetadata(f)
}

func (b *PotreeArchive) checkHierarchy() error {
	if b.metadata == nil {
		err := b.readMetadata()
		if err != nil {
			return err
		}
	}
	return nil
}

func nodeExists(nodes []*Node, node *Node) bool {
	for _, e := range nodes {
		if e != nil && e.Name == node.Name {
			return true
		}
	}
	return false
}

func (b *PotreeArchive) parseNode(n *Node, buffer io.ReadSeeker, chunkSize int64, nodes []*Node) {
	nodes[0] = n

	for i, current := range nodes {
		if current == nil {
			current = &Node{}
			current.Name = n.Name
		}
		temp := &node{}
		temp.readNode(io.LimitReader(buffer, chunkSize))

		current.NumPoints = temp.NumPoints

		if current.Type == NT_PROXY {
			current.ByteOffset = temp.ByteOffset
			current.ByteSize = temp.ByteSize
			b.nodeMaps[current.Name] = current
			nodes[i] = current
		} else {
			current.ByteOffset = temp.ByteOffset
			current.ByteSize = temp.ByteSize
		}

		current.Type = temp.Type

		if current.Type == NT_PROXY {
			if current.Name != n.Name {
				buffer.Seek(current.ByteOffset, io.SeekStart)
				_child_num_nodes := int(current.ByteSize / BytesPerNode)
				if _child_num_nodes == 0 {
					continue
				}
				_nodes := make([]*Node, _child_num_nodes)
				b.parseNode(current, buffer, current.ByteSize, _nodes)
			}
			continue
		}

		for child_index := 0; child_index < 8; child_index++ {
			child_exists := ((1 << uint32(child_index)) & uint32(current.ChildMask)) != 0
			if !child_exists {
				continue
			}

			child_name := current.Name + strconv.Itoa(child_index)
			child := &Node{}
			child.Name = child_name
			child.Parent = n

			if !nodeExists(nodes, child) {
				n.Childs[child_index] = child
			}
			b.nodeMaps[child.Name] = child
		}
	}
}

func (b *PotreeArchive) gatherChunk(start *Node, levels int) hierarchyChunk {
	startLevel := len(start.Name) - 1

	chunk := hierarchyChunk{}
	chunk.name = start.Name

	stack := NewStack()
	stack.Push(start)

	for !stack.Empty() {
		node := stack.Top()
		stack.Pop()

		chunk.nodes = append(chunk.nodes, node)

		childLevel := len(start.Name)
		if childLevel <= startLevel+levels {
			for _, child := range node.Childs {
				if child == nil {
					continue
				}

				stack.Push(child)
			}
		}
	}

	return chunk
}

func (b *PotreeArchive) createHierarchyChunks(hierarchyStepSize int) []hierarchyChunk {
	hierarchyChunks := []hierarchyChunk{}
	stack := NewStack()
	stack.Push(b.root)

	for !stack.Empty() {
		chunkRoot := stack.Top()
		stack.Pop()

		chunk := b.gatherChunk(chunkRoot, hierarchyStepSize)

		for _, node := range chunk.nodes {
			isProxy := node.Level() == chunkRoot.Level()+hierarchyStepSize

			if isProxy {
				stack.Push(node)
			}

		}

		hierarchyChunks = append(hierarchyChunks, chunk)
	}

	return hierarchyChunks
}

func (b *PotreeArchive) writeHierarchyChunk(c *hierarchyChunk, offset int64, buffer io.Writer, hierarchyStepSize int, chunks []hierarchyChunk, chunkByteOffsets []int64, chunkPointers map[string]int) error {
	chunkLevel := len(c.name) - 1
	for _, n := range c.nodes {
		isProxy := n.Level() == chunkLevel+hierarchyStepSize

		n.ChildMask = ChildMaskOf(n)

		if isProxy {
			targetChunkIndex := chunkPointers[n.Name]
			targetChunk := chunks[targetChunkIndex]

			proxy := *n
			proxy.Type = NT_PROXY
			proxy.ByteOffset = chunkByteOffsets[targetChunkIndex]
			proxy.ByteSize = targetChunk.chunkSize()

			err := proxy.writeNode(buffer)

			if err != nil {
				return err
			}
		} else {

			err := n.writeNode(buffer)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *PotreeArchive) writeNode(n *Node, buffer io.Writer) error {
	if n == nil {
		return errors.New("node is nil")
	}

	err := n.writeNode(buffer)
	if err != nil {
		return err
	}
	return nil
}

func (b *PotreeArchive) writeHierarchy() error {
	hierarchy_file := b.getHierarchyPath()
	f, err := os.Open(hierarchy_file)
	if err != nil {
		return err
	}
	defer f.Close()

	chunks := b.createHierarchyChunks(HierarchyStepSize)

	chunkPointers := make(map[string]int)
	chunkByteOffsets := make([]int64, len(chunks))
	hierarchyBufferSize := int64(0)

	for i, chunk := range chunks {
		chunkPointers[chunk.name] = i

		chunk.sortNodes()

		if i >= 1 {
			chunkByteOffsets[i] = chunkByteOffsets[i-1] + chunks[i-1].chunkSize()
		}

		hierarchyBufferSize += chunk.chunkSize()
	}

	offset := int64(0)
	for _, c := range chunks {
		c.sortNodes()
		si := c.chunkSize()

		err := b.writeHierarchyChunk(&c, offset, f, HierarchyStepSize, chunks, chunkByteOffsets, chunkPointers)
		if err != nil {
			return err
		}
		offset += si
	}

	hierarchy := &Hierarchy{}
	hierarchy.StepSize = HierarchyStepSize
	hierarchy.FirstChunkSize = int64(len(chunks[0].nodes) * BytesPerNode)

	b.metadata.Hierarchy = hierarchy
	return nil
}

func (b *PotreeArchive) readHierarchy() error {
	if b.metadata == nil {
		err := b.readMetadata()
		if err != nil {
			return err
		}
	}
	if b.metadata.Hierarchy != nil {
		first_chunk_size := b.metadata.Hierarchy.FirstChunkSize
		hierarchy_file := b.getHierarchyPath()
		if !FileExists(hierarchy_file) {
			return errors.New("hierarchy.bin not found!")
		}
		f, err := os.Open(hierarchy_file)
		if err != nil {
			return err
		}
		defer f.Close()
		b.root = &Node{}
		b.root.Type = NT_PROXY

		numNodes := int(first_chunk_size / BytesPerNode)

		_nodes := make([]*Node, numNodes)
		b.nodeMaps = make(map[string]*Node)
		f.Seek(0, io.SeekStart)
		b.parseNode(b.root, f, first_chunk_size, _nodes)
	}
	return nil
}

func (b *PotreeArchive) openOctree() error {
	p := b.getOctreePath()
	if !FileExists(p) {
		return errors.New("octree.bin not found!")
	}
	var err error
	b.octree, err = os.Open(p)
	if err != nil {
		return err
	}
	return nil
}

func (b *PotreeArchive) closeOctree() error {
	if b.octree != nil {
		return b.octree.Close()
	}
	return nil
}

func (b *PotreeArchive) writeOctree(node *Node) error {
	err := b.writeOctreeNode(node)
	if err != nil {
		return err
	}
	for _, c := range node.Childs {
		if c != nil {
			err := b.writeOctree(c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *PotreeArchive) writeOctreeNode(node *Node) error {
	if b.octree == nil {
		err := b.openOctree()
		if err != nil {
			return err
		}
	}
	if node.Buffer == nil {
		if b.metadata.IsBrotliEncoded() {
			node.Buffer = node.compress(b.metadata.Attrs[node.AttrStart:node.AttrEnd])
		} else {
			node.Buffer = node.compact(b.metadata.Attrs[node.AttrStart:node.AttrEnd], false)
		}
	}
	node.ByteOffset = b.octreeOffset
	node.ByteSize = int64(len(node.Buffer))

	b.octreeOffset += node.ByteSize
	return node.write(b.octreeOffset, node.Buffer, b.octree)
}

func (b *PotreeArchive) readOctreeNode(node *Node) error {
	if b.metadata == nil {
		err := b.readMetadata()
		if err != nil {
			return err
		}
	}
	if b.octree == nil {
		err := b.openOctree()
		if err != nil {
			return err
		}
	}
	if node != nil && node.ByteSize > 0 {
		data, err := node.read(b.octree)
		if err != nil {
			return err
		}
		node.Buffer = data
	}
	return nil
}

func (b *PotreeArchive) unpackNode(node *Node) error {
	if node.Buffer == nil {
		err := b.readOctreeNode(node)
		if err != nil {
			return err
		}
	}
	if b.metadata.IsBrotliEncoded() {
		node.uncompress(node.Buffer, b.metadata.Attrs)
	} else {
		node.uncompact(node.Buffer, b.metadata.Attrs, false)
	}
	return nil
}
