package potree

type PotreeBin struct {
	Path     string
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