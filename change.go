package osm

// Change is the structure of a changeset to be uploaded or downloaded from the osm api server.
type Change struct {
	Version     string `xml:"version,attr,omitempty" json:"version,omitempty"`
	Generator   string `xml:"generator,attr,omitempty" json:"generator,omitempty"`
	Copyright   string `xml:"copyright,attr,omitempty" json:"copyright,omitempty"` // to indicate the origin of the data
	Attribution string `xml:"attribution,attr,omitempty" json:"attribution,omitempty"`
	License     string `xml:"license,attr,omitempty" json:"license,omitempty"`
	Create      *OSM   `xml:"create" json:"create,omitempty"`
	Modify      *OSM   `xml:"modify" json:"modify,omitempty"`
	Delete      *OSM   `xml:"delete" json:"delete,omitempty"`
}
