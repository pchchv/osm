package osmxml

import (
	"context"
	"encoding/xml"
	"io"
	"strings"

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

// Scan advances the Scanner to the next element,
// which will then be available through the Object method.
// It returns false when the scan stops, either by reaching the end of the input,
// an io error, an xml error or the context being cancelled.
// After Scan returns false,
// the Error method will return any error that occurred during scanning,
// except if it was io.EOF, Error will return nil.
func (s *Scanner) Scan() bool {
	if s.error != nil {
		return false
	}

Loop:
	for {
		if s.ctx.Err() != nil {
			return false
		}

		t, err := s.decoder.Token()
		if err != nil {
			s.error = err
			return false
		}

		se, ok := t.(xml.StartElement)
		if !ok {
			continue
		}

		s.next = nil
		switch strings.ToLower(se.Name.Local) {
		case "bounds":
			bounds := &osm.Bounds{}
			err = s.decoder.DecodeElement(&bounds, &se)
			s.next = bounds
		case "node":
			node := &osm.Node{}
			err = s.decoder.DecodeElement(&node, &se)
			s.next = node
		case "way":
			way := &osm.Way{}
			err = s.decoder.DecodeElement(&way, &se)
			s.next = way
		case "relation":
			relation := &osm.Relation{}
			err = s.decoder.DecodeElement(&relation, &se)
			s.next = relation
		case "changeset":
			cs := &osm.Changeset{}
			err = s.decoder.DecodeElement(&cs, &se)
			s.next = cs
		case "note":
			n := &osm.Note{}
			err = s.decoder.DecodeElement(&n, &se)
			s.next = n
		case "user":
			u := &osm.User{}
			err = s.decoder.DecodeElement(&u, &se)
			s.next = u
		default:
			continue Loop
		}

		if err != nil {
			s.error = err
			return false
		}

		return true
	}
}
