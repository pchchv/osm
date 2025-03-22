package osm

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
