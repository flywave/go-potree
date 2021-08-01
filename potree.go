package potree

const (
	MaxPointsPerChunk = 10000
	ENCODING_BROTLI   = "BROTLI"
	HierarchyStepSize = 4
	MetadataName      = "metadata.json"
	HierarchyName     = "hierarchy.bin"
	OctreeName        = "octree.bin"
)

type PotreeBin struct {
	Path     string
	Root     *Node
	Nodes    map[string]*Node
	Metadata Metadata
}

func (b *PotreeBin) ReadMetadata() error {
	return nil
}

func (b *PotreeBin) WriteMetadata() error {
	return nil
}

func (b *PotreeBin) checkHierarchy() error {
	return nil
}

func (b *PotreeBin) readHierarchy() error {
	return nil
}

func (b *PotreeBin) openOctree() error {
	return nil
}

func (b *PotreeBin) closeOctree() error {
	return nil
}

func (b *PotreeBin) readOctreeNode(node *Node) error {
	return nil
}
