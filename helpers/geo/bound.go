package geo

// Bound represents a closed box or rectangle.
type Bound struct {
	Min Point
	Max Point
}

// Equal returns if two bounds are equal.
func (b Bound) Equal(c Bound) bool {
	return b.Min == c.Min && b.Max == c.Max
}

// GeoJSONType returns the GeoJSON type for the object.
func (b Bound) GeoJSONType() string {
	return "Polygon"
}

// Dimensions returns 2 because a Bound is a 2d object.
func (b Bound) Dimensions() int {
	return 2
}

// Bound returns the the same bound.
func (b Bound) Bound() Bound {
	return b
}
