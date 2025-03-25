package osm

// Tag is a key-value item attached to
// osm nodes, ways and relations.
type Tag struct {
	Key   string `xml:"k,attr"`
	Value string `xml:"v,attr"`
}

// Tags is a collection of Tag objects
// with some helper functions.
type Tags []Tag

