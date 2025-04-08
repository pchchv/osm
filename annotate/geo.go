package annotate

import (
	"math"
	"time"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/geometries"
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

func wayCentroid(w *osm.Way) geo.Point {
	var dist float64
	point := geo.Point{}
	seg := [2]geo.Point{}
	for i := 0; i < len(w.Nodes)-1; i++ {
		seg[0] = w.Nodes[i].Point()
		seg[1] = w.Nodes[i+1].Point()
		d := geometries.Distance(seg[0], seg[1])
		point[0] += (seg[0][0] + seg[1][0]) / 2.0 * d
		point[1] += (seg[0][1] + seg[1][1]) / 2.0 * d
		dist += d
	}

	point[0] /= dist
	point[1] /= dist

	return point
}

// wayPointOnSurface finds closest node to centroid.
func wayPointOnSurface(w *osm.Way) geo.Point {
	var index int
	centroid := wayCentroid(w)
	min := math.MaxFloat64
	for i, n := range w.Nodes {
		if d := geometries.Distance(centroid, n.Point()); d < min {
			index, min = i, d
		}
	}

	return w.Nodes[index].Point()
}
