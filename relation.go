package osm

import "github.com/pchchv/geo"

// RelationID is the primary key of a relation.
// Relation is uniquely identifiable by the id + version.
type RelationID int64

// FeatureID is a helper returning the feature id for this relation id.
func (id RelationID) FeatureID() FeatureID {
	return FeatureID(relationMask | id<<versionBits)
}

// ObjectID is a helper returning the object id for this relation id.
func (id RelationID) ObjectID(v int) ObjectID {
	return ObjectID(id.ElementID(v))
}

// ElementID is a helper to convert the id to an element id.
func (id RelationID) ElementID(v int) ElementID {
	return id.FeatureID().ElementID(v)
}

// Member is a member of a relation.
type Member struct {
	Type        Type        `xml:"type,attr" json:"type"`
	Ref         int64       `xml:"ref,attr" json:"ref"`
	Role        string      `xml:"role,attr" json:"role"`
	Version     int         `xml:"version,attr,omitempty" json:"version,omitempty"`
	ChangesetID ChangesetID `xml:"changeset,attr,omitempty" json:"changeset,omitempty"`
	// Node location if Type == Node
	// Closest vertex to centroid if Type == Way
	// Empty/invalid if Type == Relation
	Lat float64 `xml:"lat,attr,omitempty" json:"lat,omitempty"`
	Lon float64 `xml:"lon,attr,omitempty" json:"lon,omitempty"`
	// Orientation is the direction of the way around a ring of a multipolygon.
	// Only valid for multipolygon or boundary relations.
	Orientation geo.Orientation `xml:"orientation,attr,omitempty" json:"orientation,omitempty"`
	// Nodes are sometimes included in members of type way to include the lat/lon
	// path of the way. Overpass returns xml like this.
	Nodes WayNodes `xml:"nd" json:"nodes,omitempty"`
}

// FeatureID returns the feature id of the member.
func (m Member) FeatureID() FeatureID {
	switch m.Type {
	case TypeNode:
		return NodeID(m.Ref).FeatureID()
	case TypeWay:
		return WayID(m.Ref).FeatureID()
	case TypeRelation:
		return RelationID(m.Ref).FeatureID()
	default:
		panic("unknown type")
	}
}

// ElementID returns the element id of the member.
func (m Member) ElementID() ElementID {
	return m.FeatureID().ElementID(m.Version)
}

// Point returns the geo.Point location for the member.
// Will be (0, 0) if the relation is not annotated.
// For way members this location is annotated as the "surface point".
func (m Member) Point() geo.Point {
	return geo.Point{m.Lon, m.Lat}
}
