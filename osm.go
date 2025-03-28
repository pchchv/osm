package osm

import "fmt"

const (
	// should be returned if the osm data is actual
	// osm data to give some information about the source and license
	Copyright   = "OpenStreetMap and contributors"
	Attribution = "http://www.openstreetmap.org/copyright"
	License     = "http://opendatacommons.org/licenses/odbl/1-0/"
)

// OSM represents the core osm data designed
// to parse [OSM XML](http://wiki.openstreetmap.org/wiki/OSM_XML).
type OSM struct {
	// JSON APIs can return version as a string or number,
	// converted to string for consistency.
	Version   string `xml:"version,attr,omitempty"`
	Generator string `xml:"generator,attr,omitempty"`
	// These three attributes are returned by the osm api.
	// The Copyright, Attribution and License constants contain
	// suggested values that match those returned by the official api.
	Copyright   string    `xml:"copyright,attr,omitempty"`
	Attribution string    `xml:"attribution,attr,omitempty"`
	License     string    `xml:"license,attr,omitempty"`
	Bounds      *Bounds   `xml:"bounds,omitempty"`
	Nodes       Nodes     `xml:"node"`
	Ways        Ways      `xml:"way"`
	Relations   Relations `xml:"relation"`
	// Changesets typically not be included with actual data,
	// but all this stuff is technically all under the osm xml
	Changesets Changesets `xml:"changeset"`
	Notes      Notes      `xml:"note"`
	Users      Users      `xml:"user"`
}

// FeatureIDs returns the slice of feature ids for
// all the nodes, ways and relations.
func (o *OSM) FeatureIDs() FeatureIDs {
	if o == nil {
		return nil
	}

	result := make(FeatureIDs, 0, len(o.Nodes)+len(o.Ways)+len(o.Relations))
	for _, e := range o.Nodes {
		result = append(result, e.FeatureID())
	}

	for _, e := range o.Ways {
		result = append(result, e.FeatureID())
	}

	for _, e := range o.Relations {
		result = append(result, e.FeatureID())
	}

	return result
}

// ElementIDs returns the slice of element ids for
// all the nodes, ways and relations.
func (o *OSM) ElementIDs() ElementIDs {
	if o == nil {
		return nil
	}

	result := make(ElementIDs, 0, len(o.Nodes)+len(o.Ways)+len(o.Relations))
	for _, e := range o.Nodes {
		result = append(result, e.ElementID())
	}

	for _, e := range o.Ways {
		result = append(result, e.ElementID())
	}

	for _, e := range o.Relations {
		result = append(result, e.ElementID())
	}

	return result
}

type typeS struct {
	Type string `json:"type"`
}

func findType(index int, data []byte) (string, error) {
	ts := typeS{}
	if err := unmarshalJSON(data, &ts); err != nil {
		// should not happened due to previous decoding succeeded
		return "", err
	}

	if ts.Type == "" {
		return "", fmt.Errorf("could not find type in element index %d", index)
	}

	return ts.Type, nil
}
