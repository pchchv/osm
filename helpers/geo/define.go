package geo

// Pointer is something that can be represented by a point.
type Pointer interface {
	Point() Point
}
