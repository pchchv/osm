package osmpbf

import (
	"context"
	"io"

	"github.com/pchchv/osm"
)

// Scanner provides a convenient interface reading a stream of osm data from a file or url.
// Successive calls to the Scan method will step through the data.
//
// Scanning stops unrecoverably at EOF, the first I/O error,
// the first xml error or the context being cancelled.
// When a scan stops,
// the reader may have advanced arbitrarily far past the last token.
//
// Scanner API is based on [bufio.Scanner](https://golang.org/pkg/bufio/#Scanner)
type Scanner struct {
	SkipNodes      bool // Skip element types that are not needed.
	SkipWays       bool // The data is skipped at the encoded protobuf level,
	SkipRelations  bool // but each block still needs to be decompressed.
	started        bool
	closed         bool
	FilterNode     func(*osm.Node) bool     // Filter functions must be fast, they block the decoder, there are `procs` number of concurrent decoders.
	FilterWay      func(*osm.Way) bool      // Elements can be stored if the function returns true, or skipped if false.
	FilterRelation func(*osm.Relation) bool // Memory is reused if Filter returns false.
	ctx            context.Context
	decoder        *decoder
	procs          int
	next           osm.Object
	err            error
}

// New returns a new Scanner to read from r.
// procs indicates amount of paralellism,
// when reading blocks which will off load the
// unzipping/decoding to multiple cpus.
func New(ctx context.Context, r io.Reader, procs int) *Scanner {
	if ctx == nil {
		ctx = context.Background()
	}

	s := &Scanner{
		ctx:   ctx,
		procs: procs,
	}
	s.decoder = newDecoder(ctx, s, r)
	return s
}
