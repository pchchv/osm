package osm

import (
	"time"

	"github.com/pchchv/geo"
)

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

// Members represents an ordered list of relation members.
type Members []Member

// FeatureIDs returns the a list of feature ids for the members.
func (ms Members) FeatureIDs() FeatureIDs {
	ids := make(FeatureIDs, len(ms), len(ms)+1)
	for i, m := range ms {
		ids[i] = m.FeatureID()
	}

	return ids
}

// ElementIDs returns the a list of element ids for the members.
func (ms Members) ElementIDs() ElementIDs {
	ids := make(ElementIDs, len(ms), len(ms)+1)
	for i, m := range ms {
		ids[i] = m.ElementID()
	}

	return ids
}

// MarshalJSON allows the members to be marshalled
// as defined by the overpass osmjson.
// This function is a wrapper to marshal null as [].
func (ms Members) MarshalJSON() ([]byte, error) {
	if len(ms) == 0 {
		return []byte(`[]`), nil
	}

	return marshalJSON([]Member(ms))
}

// Relation is an collection of nodes,
// ways and other relations with some defining attributes.
type Relation struct {
	XMLName     xmlNameJSONTypeRel `xml:"relation" json:"type"`
	ID          RelationID         `xml:"id,attr" json:"id"`
	User        string             `xml:"user,attr" json:"user,omitempty"`
	UserID      UserID             `xml:"uid,attr" json:"uid,omitempty"`
	Visible     bool               `xml:"visible,attr" json:"visible"`
	Version     int                `xml:"version,attr" json:"version,omitempty"`
	ChangesetID ChangesetID        `xml:"changeset,attr" json:"changeset,omitempty"`
	Timestamp   time.Time          `xml:"timestamp,attr" json:"timestamp,omitempty"`
	Tags        Tags               `xml:"tag" json:"tags,omitempty"`
	Members     Members            `xml:"member" json:"members"`
	// Committed, is the estimated time this object was committed
	// and made visible in the central OSM database.
	Committed *time.Time `xml:"committed,attr,omitempty" json:"committed,omitempty"`
	// Updates are changes to the members of this relation independent
	// of an update to the relation itself. The OSM api allows a child
	// to be updated without any changes to the parent.
	Updates Updates `xml:"update,omitempty" json:"updates,omitempty"`
	// Bounds are included by overpass, and maybe others
	Bounds *Bounds `xml:"bounds,omitempty" json:"bounds,omitempty"`
}

// ObjectID returns the object id of the relation.
func (r *Relation) ObjectID() ObjectID {
	return r.ID.ObjectID(r.Version)
}

// FeatureID returns the feature id of the relation.
func (r *Relation) FeatureID() FeatureID {
	return r.ID.FeatureID()
}

// ElementID returns the element id of the relation.
func (r *Relation) ElementID() ElementID {
	return r.ID.ElementID(r.Version)
}

// ApplyUpdatesUpTo applies the updates to this object upto and including the given time.
func (r *Relation) ApplyUpdatesUpTo(t time.Time) error {
	var notApplied []Update
	for _, u := range r.Updates {
		if u.Timestamp.After(t) {
			notApplied = append(notApplied, u)
			continue
		}

		if err := r.applyUpdate(u); err != nil {
			return err
		}
	}

	r.Updates = notApplied
	return nil
}

// applyUpdate modifies the current relation and dictated by the given update.
// Will return UpdateIndexOutOfRangeError if the update index is too large.
func (r *Relation) applyUpdate(u Update) error {
	if u.Index >= len(r.Members) {
		return &UpdateIndexOutOfRangeError{Index: u.Index}
	}

	r.Members[u.Index].Version = u.Version
	r.Members[u.Index].ChangesetID = u.ChangesetID
	r.Members[u.Index].Lat = u.Lat
	r.Members[u.Index].Lon = u.Lon
	if u.Reverse {
		r.Members[u.Index].Orientation *= -1
	}

	return nil
}

// TagMap returns the element tags as a key/value map.
func (r *Relation) TagMap() map[string]string {
	return r.Tags.Map()
}

// CommittedAt returns the best estimate on when this
// element became was written/committed into the database.
func (r *Relation) CommittedAt() time.Time {
	if r.Committed != nil {
		return *r.Committed
	}

	return r.Timestamp
}

// Relations is a list of relations with some helper functions attached.
type Relations []*Relation

type relationsSort Relations

func (rs relationsSort) Len() int {
	return len(rs)
}

func (rs relationsSort) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs relationsSort) Less(i, j int) bool {
	if rs[i].ID == rs[j].ID {
		return rs[i].Version < rs[j].Version
	}

	return rs[i].ID < rs[j].ID
}
