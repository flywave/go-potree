package potree

import (
	"math"
	"os"

	vec3d "github.com/flywave/go3d/float64/vec3"
)

func splitBy3(a uint32) uint64 {
	x := uint64(a & 0x1fffff)           // we only look at the first 21 bits
	x = (x | x<<32) & 0x1f00000000ffff  // shift left 32 bits, OR with self, and 00011111000000000000000000000000000000001111111111111111
	x = (x | x<<16) & 0x1f0000ff0000ff  // shift left 32 bits, OR with self, and 00011111000000000000000011111111000000000000000011111111
	x = (x | x<<8) & 0x100f00f00f00f00f // shift left 32 bits, OR with self, and 0001000000001111000000001111000000001111000000001111000000000000
	x = (x | x<<4) & 0x10c30c30c30c30c3 // shift left 32 bits, OR with self, and 0001000011000011000011000011000011000011000011000011000100000000
	x = (x | x<<2) & 0x1249249249249249
	return x
}

// see https://www.forceflow.be/2013/10/07/morton-encodingdecoding-through-bit-interleaving-implementations/
func MortonEncodeMagicBits(x, y, z uint32) uint64 {
	answer := uint64(0)
	answer |= splitBy3(x) | splitBy3(y)<<1 | splitBy3(z)<<2
	return answer
}

func dealign24b(mortoncode uint32) uint32 {
	// see https://stackoverflow.com/questions/45694690/how-i-can-remove-all-odds-bits-in-c

	// input alignment of desired bits
	// ..a..b..c..d..e..f..g..h..i..j..k..l..m..n..o..p
	x := mortoncode

	//          ..a..b..c..d..e..f..g..h..i..j..k..l..m..n..o..p                     ..a..b..c..d..e..f..g..h..i..j..k..l..m..n..o..p
	//          ..a.....c.....e.....g.....i.....k.....m.....o...                     .....b.....d.....f.....h.....j.....l.....n.....p
	//          ....a.....c.....e.....g.....i.....k.....m.....o.                     .....b.....d.....f.....h.....j.....l.....n.....p
	x = ((x & 0b001000001000001000001000) >> 2) | ((x & 0b000001000001000001000001) >> 0)
	//          ....ab....cd....ef....gh....ij....kl....mn....op                     ....ab....cd....ef....gh....ij....kl....mn....op
	//          ....ab..........ef..........ij..........mn......                     ..........cd..........gh..........kl..........op
	//          ........ab..........ef..........ij..........mn..                     ..........cd..........gh..........kl..........op
	x = ((x & 0b000011000000000011000000) >> 4) | ((x & 0b000000000011000000000011) >> 0)
	//          ........abcd........efgh........ijkl........mnop                     ........abcd........efgh........ijkl........mnop
	//          ........abcd....................ijkl............                     ....................efgh....................mnop
	//          ................abcd....................ijkl....                     ....................efgh....................mnop
	x = ((x & 0b000000001111000000000000) >> 8) | ((x & 0b000000000000000000001111) >> 0)
	//          ................abcdefgh................ijklmnop                     ................abcdefgh................ijklmnop
	//          ................abcdefgh........................                     ........................................ijklmnop
	//          ................................abcdefgh........                     ........................................ijklmnop
	x = ((x & 0b000000000000000000000000) >> 16) | ((x & 0b000000000000000011111111) >> 0)

	// sucessfully realigned!
	//................................abcdefghijklmnop

	return x
}

type MortonCode struct {
	lower    uint64
	upper    uint64
	whatever uint64
	index    uint64
}

func MinInt(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	}
	return true
}

type ScaleOffset struct {
	scale  [3]float64
	offset [3]float64
}

func ComputeScaleOffset(min, max, targetScale vec3d.T) ScaleOffset {
	offset := min
	scale := targetScale
	size := vec3d.Sub(&max, &min)

	min_scale_x := size[0] / math.Pow(2.0, 30.0)
	min_scale_y := size[1] / math.Pow(2.0, 30.0)
	min_scale_z := size[2] / math.Pow(2.0, 30.0)

	scale[0] = math.Max(scale[0], min_scale_x)
	scale[1] = math.Max(scale[1], min_scale_y)
	scale[2] = math.Max(scale[2], min_scale_z)

	var scaleOffset ScaleOffset
	scaleOffset.scale = scale
	scaleOffset.offset = offset

	return scaleOffset
}
