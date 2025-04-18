package osmpbf

import (
	"context"
	"io"
	"sync/atomic"

	"github.com/pchchv/osm"
)

var _ osm.Scanner = &Scanner{}

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

// Scan advances the Scanner to the next element,
// which will then be available through the Element method.
// It returns false when the scan stops,
// either by reaching the end of the input,
// an io error, an xml error or the context being cancelled.
// After Scan returns false,
// the Err method will return any error that occurred during scanning,
// except that if it was io.EOF, Err will return nil.
func (s *Scanner) Scan() bool {
	if !s.started {
		s.started = true
		s.err = s.decoder.Start(s.procs)
	}

	if s.err != nil || s.closed || s.ctx.Err() != nil {
		return false
	}

	s.next, s.err = s.decoder.Next()
	return s.err == nil
}

// FullyScannedBytes returns the number of bytes that have been read and fully scanned.
// OSM protobuf files contain data blocks with 8000 nodes each.
// The returned value contains the bytes for the blocks that have been fully scanned.
//
// A user can use this number of seek forward in a file and begin reading mid-data.
// Note that while elements are usually sorted by Type, ID, Version in OSM protobuf files,
// versions of given element may span blocks.
func (s *Scanner) FullyScannedBytes() int64 {
	return atomic.LoadInt64(&s.decoder.cOffset)
}

// PreviousFullyScannedBytes returns the previous value of FullyScannedBytes.
// This is interesting because it's not totally clear if a feature spans a block.
// For example, if one quits after finding the first relation,
// upon restarting there is no way of knowing if the first relation is complete, so skip it.
// But if this relation is the first relation in the file we'll skip a full relation.
func (s *Scanner) PreviousFullyScannedBytes() int64 {
	return atomic.LoadInt64(&s.decoder.pOffset)
}

// Object returns the most recent token generated by a call to Scan as a new osm.Object.
// Currently osm.pbf files only contain nodes, ways and relations.
// This method returns an object so match the
// osm.Scanner interface and allows this Scanner to share an
// interface with osmxml.Scanner.
func (s *Scanner) Object() osm.Object {
	return s.next
}

// Header returns the pbf file header with
// interesting information about how it was created.
func (s *Scanner) Header() (*Header, error) {
	if !s.started {
		s.started = true
		// header gets read before Start returns
		s.err = s.decoder.Start(s.procs)
	}

	return s.decoder.header, s.err
}

// Close cleans up all the reading goroutines,
// it does not close the underlying reader.
func (s *Scanner) Close() error {
	s.closed = true
	return s.decoder.Close()
}

// Err returns the first non-EOF error that was encountered by the Scanner.
func (s *Scanner) Err() error {
	if s.err != nil {
		if s.err == io.EOF {
			return nil
		}
		return s.err
	}

	if s.closed {
		return osm.ErrScannerClosed
	}

	return s.ctx.Err()
}
