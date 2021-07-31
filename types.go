package potree

type AABB struct {
	Min [3]float64 `json:"min,omitempty"`
	Max [3]float64 `json:"max,omitempty"`
}
