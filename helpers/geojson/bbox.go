package geojson

import "github.com/pchchv/osm/helpers/geo"

// BBox is for the geojson bbox attribute which is an
// array with all axes of the most southwesterly point
// followed by all axes of the more northeasterly point.
type BBox []float64

// NewBBox creates a bbox from a a bound.
func NewBBox(b geo.Bound) BBox {
	return []float64{
		b.Min[0], b.Min[1],
		b.Max[0], b.Max[1],
	}
}
