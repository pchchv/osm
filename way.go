package osm

import "github.com/pchchv/geo"

// WayID is the primary key of a way.
// Way is uniquely identifiable by the id + version.
type WayID int64

// ObjectID is a helper returning the object id for this way id.
func (id WayID) ObjectID(v int) ObjectID {
	return ObjectID(id.ElementID(v))
}

// FeatureID is a helper returning the feature id for this way id.
func (id WayID) FeatureID() FeatureID {
	return FeatureID(wayMask | (id << versionBits))
}

// ElementID is a helper to convert the id to an element id.
func (id WayID) ElementID(v int) ElementID {
	return id.FeatureID().ElementID(v)
}

// WayNode is a short node used as part of ways and relations in the osm xml.
type WayNode struct {
	ID          NodeID      `xml:"ref,attr,omitempty"`
	Version     int         `xml:"version,attr,omitempty"`
	ChangesetID ChangesetID `xml:"changeset,attr,omitempty"`
	Lat         float64     `xml:"lat,attr,omitempty"`
	Lon         float64     `xml:"lon,attr,omitempty"`
}

// Point returns the geo.Point location for the way node.
// Will be (0, 0) if the way is not annotated.
func (wn WayNode) Point() geo.Point {
	return geo.Point{wn.Lon, wn.Lat}
}

// FeatureID returns the feature id of the way node.
func (wn WayNode) FeatureID() FeatureID {
	return wn.ID.FeatureID()
}

// ElementID returns the element id of the way node.
func (wn WayNode) ElementID() ElementID {
	return wn.ID.ElementID(wn.Version)
}

// WayNodes represents a collection of way nodes.
type WayNodes []WayNode
