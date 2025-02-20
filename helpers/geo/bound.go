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
