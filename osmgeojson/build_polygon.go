package osmgeojson

import "github.com/pchchv/geo"

func polygonContains(outer geo.Ring, r geo.Ring) bool {
	for _, p := range r {
		var inside bool
		x, y := p[0], p[1]
		i, j := 0, len(outer)-1
		for i < len(outer) {
			xi, yi := outer[i][0], outer[i][1]
			xj, yj := outer[j][0], outer[j][1]
			if ((yi > y) != (yj > y)) &&
				(x < (xj-xi)*(y-yi)/(yj-yi)+xi) {
				inside = !inside
			}

			j = i
			i++
		}

		if inside {
			return true
		}
	}

	return false
}

func addToMultiPolygon(mp geo.MultiPolygon, ring geo.Ring, includeInvalidPolygons bool) geo.MultiPolygon {
	for i := range mp {
		if polygonContains(mp[i][0], ring) {
			mp[i] = append(mp[i], ring)
			return mp
		}
	}

	if !includeInvalidPolygons {
		// inner without its outer
		return mp
	}

	if len(mp) > 0 {
		// if the outer ring of the first polygon is not closed,
		// it is not known whether this inner must be part of it,
		// but it is assumed that it is
		fr := mp[0][0]
		if len(fr) != 0 && fr[0] != fr[len(fr)-1] {
			mp[0] = append(mp[0], ring)
			return mp
		}

		// trying to find an existing "without outer" polygon to add this to
		for i := range mp {
			if len(mp[i][0]) == 0 {
				mp[i] = append(mp[i], ring)
				return mp
			}
		}
	}

	// no polygons with empty outer, so create one
	// create another polygon with empty outer
	return append(mp, geo.Polygon{nil, ring})
}
