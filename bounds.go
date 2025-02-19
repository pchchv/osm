package osm

// Bounds are the bounds of osm data as defined in the xml file.
type Bounds struct {
	MinLat float64 `xml:"minlat,attr"`
	MaxLat float64 `xml:"maxlat,attr"`
	MinLon float64 `xml:"minlon,attr"`
	MaxLon float64 `xml:"maxlon,attr"`
}

// ObjectID returns the bounds type but with 0 id.
// Since id doesn't make sense.
// This is here to implement the Object interface since it technically is an osm object type.
// It also allows bounds to be returned via the osmxml.Scanner.
func (b *Bounds) ObjectID() ObjectID {
	return ObjectID(boundsMask)
}
