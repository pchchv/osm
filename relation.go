package osm

// RelationID is the primary key of a relation.
// A relation is uniquely identifiable by the id + version.
type RelationID int64

// FeatureID is a helper returning the feature id for this relation id.
func (id RelationID) FeatureID() FeatureID {
	return FeatureID(relationMask | id<<versionBits)
}
