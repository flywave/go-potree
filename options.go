package potree

const (
	ENCODING_BROTLI       = "BROTLI"
	ENCODING_DEFAULT      = "DEFAULT"
	ENCODING_UNCOMPRESSED = "UNCOMPRESSED"
)

type Options struct {
	Encoding string
	Outdir   string
	Name     string
}
