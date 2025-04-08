package annotate

import (
	"time"

	"github.com/pchchv/geo"
	"github.com/pchchv/osm"
	"github.com/pchchv/osm/mputil"
)

// orientation annotates the orientation of multipolygon relation members.
// This makes it possible to reconstruct relations with partial data in the right direction.
// Return value indicates if the result is 'tainted', e.g. not all way members were present.
func orientation(members osm.Members, ways map[osm.WayID]*osm.Way, at time.Time) bool {
	outer, inner, tainted := mputil.Group(members, ways, at)
	outers, inners := mputil.Join(outer), mputil.Join(inner)
	for _, outer := range outers {
		annotateOrientation(members, outer, geo.CCW)
	}

	for _, inner := range inners {
		annotateOrientation(members, inner, geo.CW)
	}

	return tainted
}

func annotateOrientation(members osm.Members, ms mputil.MultiSegment, o geo.Orientation) {
	factor := geo.Orientation(1)
	if ms.Orientation() != o {
		factor = -1
	}

	for _, segment := range ms {
		if segment.Reversed {
			members[segment.Index].Orientation = -1 * factor * o
		} else {
			members[segment.Index].Orientation = factor * o
		}
	}
}
