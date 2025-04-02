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
