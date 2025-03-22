package osm

// NodeID corresponds the primary key of a node.
// The node id + version uniquely identify a node.
type NodeID int64

// ObjectID is a helper returning the object id for this node id.
func (id NodeID) ObjectID(v int) ObjectID {
	return ObjectID(id.ElementID(v))
}

// FeatureID is a helper returning the feature id for this node id.
func (id NodeID) FeatureID() FeatureID {
	return FeatureID(nodeMask | (id << versionBits))
}

// ElementID is a helper to convert the id to an element id.
func (id NodeID) ElementID(v int) ElementID {
	return id.FeatureID().ElementID(v)
}
