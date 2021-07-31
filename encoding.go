package potree

// #include <stdlib.h>
// #include <string.h>
// #include "brotli/decode.h"
// #include "brotli/encode.h"
// #cgo CFLAGS: -I ./  -I ./libs
// #cgo CXXFLAGS: -I ./ -I ./libs
// #cgo linux LDFLAGS: -L ./libs  -Wl,--start-group -lbrotlicommon -lbrotlidec -lm -lbrotlienc -Wl,--end-group
// #cgo darwin LDFLAGS: -L ./libs  -lbrotlicommon -lbrotlidec -lbrotlienc
// #cgo windows LDFLAGS: -L ./libs -lbrotlicommon -lbrotlidec -lbrotlienc  -fPIC
import "C"

type BrotliEnc struct {
}

type BrotliDec struct {
}
