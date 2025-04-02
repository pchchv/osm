package mputil

import "github.com/pchchv/geo"

// Segment is a section of a
// multipolygon with some extra information on
// the member it came from.
type Segment struct {
	Index       uint32
	Reversed    bool
	Orientation geo.Orientation
	Line        geo.LineString
}

// First returns the first point in
// the segment linestring.
func (s Segment) First() geo.Point {
	return s.Line[0]
}

// Last returns the last point in
// the segment linestring.
func (s Segment) Last() geo.Point {
	return s.Line[len(s.Line)-1]
}

// Reverse reverses the line string of the segment.
func (s *Segment) Reverse() {
	s.Reversed = !s.Reversed
	s.Line.Reverse()
}
