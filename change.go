package osm

import "encoding/xml"

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

// AppendCreate appends the object to the Create OSM object.
func (c *Change) AppendCreate(o Object) {
	if c.Create == nil {
		c.Create = &OSM{}
	}

	c.Create.Append(o)
}

// AppendModify appends the object to the Modify OSM object.
func (c *Change) AppendModify(o Object) {
	if c.Modify == nil {
		c.Modify = &OSM{}
	}

	c.Modify.Append(o)
}

// AppendDelete appends the object to the Delete OSM object.
func (c *Change) AppendDelete(o Object) {
	if c.Delete == nil {
		c.Delete = &OSM{}
	}

	c.Delete.Append(o)
}

// MarshalXML implements the xml.Marshaller method to allow for the
// correct wrapper/start element case and attr data.
func (c Change) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	start.Name.Local = "osmChange"
	start.Attr = []xml.Attr{}
	if c.Version != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "version"}, Value: c.Version})
	}

	if c.Generator != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "generator"}, Value: c.Generator})
	}

	if c.Copyright != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "copyright"}, Value: c.Copyright})
	}

	if c.Attribution != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "attribution"}, Value: c.Attribution})
	}

	if c.License != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "license"}, Value: c.License})
	}

	if err = e.EncodeToken(start); err != nil {
		return
	}

	if err = marshalInnerChange(e, "create", c.Create); err != nil {
		return
	}

	if err = marshalInnerChange(e, "modify", c.Modify); err != nil {
		return
	}

	if err = marshalInnerChange(e, "delete", c.Delete); err != nil {
		return
	}

	return e.EncodeToken(start.End())
}

func marshalInnerChange(e *xml.Encoder, name string, o *OSM) (err error) {
	if o == nil {
		return nil
	}

	t := xml.StartElement{Name: xml.Name{Local: name}}
	if err = e.EncodeToken(t); err != nil {
		return
	}

	if err = o.marshalInnerXML(e); err != nil {
		return
	}

	return e.EncodeToken(t.End())
}
