package geo

// Geometry represents the shared attributes of a geometry.
type Geometry interface {
	GeoJSONType() string
	Dimensions() int // i.e., 0d, 1d, 2d
	Bound() Bound
	private() // requiring because sub package type switch over all possible types
}
