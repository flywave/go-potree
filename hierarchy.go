package potree

type Hierarchy struct {
	StepSize       int64  `json:"stepSize"`
	FirstChunkSize int64  `json:"firstChunkSize"`
	Buffer         []byte `json:"-"`
	Depth          *int   `json:"depth,omitempty"`
}
