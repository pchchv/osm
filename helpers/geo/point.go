package geo

var _ Pointer = Point{}

// Point is a Lon/Lat 2d point.
type Point [2]float64

// Point returns itself so it implements the Pointer interface.
func (p Point) Point() Point {
	return p
}

// Lon returns the horizontal, longitude coordinate of the point.
func (p Point) Lon() float64 {
	return p[0]
}

// Equal checks if the point represents the same point or vector.
func (p Point) Equal(point Point) bool {
	return p[0] == point[0] && p[1] == point[1]
}
