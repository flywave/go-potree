package potree

import vec3d "github.com/flywave/go3d/float64/vec3"

type Metadata struct {
	BoundingBox AABB
	Scale       vec3d.T
	Offset      vec3d.T
}
