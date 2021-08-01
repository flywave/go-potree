package potree

import (
	vec3d "github.com/flywave/go3d/float64/vec3"
)

type Chunk struct {
	min  vec3d.T
	max  vec3d.T
	file string
	id   string
}

type Chunks struct {
	list       []*Chunk
	min        vec3d.T
	max        vec3d.T
	attributes []Attribute
}

type HierarchyChunk struct {
	name  string
	nodes []*Node
}
