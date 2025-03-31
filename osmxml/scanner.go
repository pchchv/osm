package osmxml

import (
	"context"
	"encoding/xml"
	"io"

	"github.com/pchchv/osm"
)

// Scanner provides a convenient interface for reading a stream of osm data from a file or url.
// Successive calls to the Scan method will step through the data.
//
// Scanning is irrevocably stopped at EOF, first I/O error, first xml error or context cancel.
// When a scan stops, the reader may have advanced arbitrarily far past the last token.
//
// The Scanner API is based on [bufio.Scanner](https://golang.org/pkg/bufio/#Scanner)
type Scanner struct {
	ctx     context.Context
	done    context.CancelFunc
	closed  bool
	decoder *xml.Decoder
	next    osm.Object
	error   error
}

// New returns a new Scanner to read from r.
func New(ctx context.Context, r io.Reader) *Scanner {
	if ctx == nil {
		ctx = context.Background()
	}

	s := &Scanner{
		decoder: xml.NewDecoder(r),
	}
	s.ctx, s.done = context.WithCancel(ctx)
	return s
}
