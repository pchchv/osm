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

// MultiSegment is an ordered set of
// segments that form a continuous
// section of a multipolygon.
type MultiSegment []Segment

// First returns the first point in the list of linestrings.
func (ms MultiSegment) First() geo.Point {
	return ms[0].Line[0]
}

// Last returns the last point in the list of linestrings.
func (ms MultiSegment) Last() geo.Point {
	line := ms[len(ms)-1].Line
	return line[len(line)-1]
}

// Ring converts the multisegment to a ring of the given orientation.
// Ring uses the orientation on the members if possible.
func (ms MultiSegment) Ring(o geo.Orientation) geo.Ring {
	var length int
	for _, s := range ms {
		length += len(s.Line)
	}

	var haveOrient, reversed bool
	ring := make(geo.Ring, 0, length)
	for _, s := range ms {
		if s.Orientation != 0 {
			haveOrient = true
			if (s.Orientation == o) == s.Reversed {
				reversed = true
			}
		}

		ring = append(ring, s.Line...)
	}

	if (haveOrient && reversed) || (!haveOrient && ring.Orientation() != o) {
		ring.Reverse()
	}

	return ring
}

// Orientation computes the orientation of a multisegment like if it was ring.
func (ms MultiSegment) Orientation() geo.Orientation {
	var area float64
	prev := ms.First()
	// implicitly move everything to near the origin to help with roundoff
	offset := prev
	for _, segment := range ms {
		for _, point := range segment.Line {
			area += (prev[0]-offset[0])*(point[1]-offset[1]) - (point[0]-offset[0])*(prev[1]-offset[1])
			prev = point
		}
	}

	if area > 0 {
		return geo.CCW
	}

	return geo.CW
}

// LineString converts a multisegment into a geo linestring object.
func (ms MultiSegment) LineString() geo.LineString {
	var length int
	for _, s := range ms {
		length += len(s.Line)
	}

	line := make(geo.LineString, 0, length)
	for _, s := range ms {
		line = append(line, s.Line...)
	}

	return line
}
