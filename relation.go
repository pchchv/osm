package osm

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
