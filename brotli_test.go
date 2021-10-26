package potree

import "testing"

func TestReadPoTree(t *testing.T) {
	arch := NewArchive("./tests")
	arch.Load()
}
