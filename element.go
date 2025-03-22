package osm

import "fmt"

// ElementID is a unique key for an osm element.
// It contains the type, id and version information.
type ElementID int64

// Type returns the Type for the element.
func (id ElementID) Type() Type {
	switch id & typeMask {
	case nodeMask:
		return TypeNode
	case wayMask:
		return TypeWay
	case relationMask:
		return TypeRelation
	default:
		panic("unknown type")
	}
}

// Version returns the version of the element.
func (id ElementID) Version() int {
	return int(id & (versionMask))
}

// Ref returns the ID reference for the element.
// Not unique without the type.
func (id ElementID) Ref() int64 {
	return int64((id & refMask) >> versionBits)
}

// String returns "type/ref:version" for the element.
func (id ElementID) String() (v string) {
	if id.Version() == 0 {
		v = "-"
	} else {
		v = fmt.Sprintf("%d", id.Version())
	}

	return fmt.Sprintf("%s/%d:", id.Type(), id.Ref()) + v
}

// ObjectID is a helper to convert the id to an object id.
func (id ElementID) ObjectID() ObjectID {
	return ObjectID(id)
}

// FeatureID returns the feature id for the element id. i.e removing the version.
func (id ElementID) FeatureID() FeatureID {
	return FeatureID(id & featureMask)
}

// NodeID returns the id of this feature as a node id.
// The function will panic if this element is not of TypeNode.
func (id ElementID) NodeID() NodeID {
	if id&nodeMask != nodeMask {
		panic(fmt.Sprintf("not a node: %v", id))
	}

	return NodeID(id.Ref())
}

// WayID returns the id of this feature as a way id.
// The function will panic if this element is not of TypeWay.
func (id ElementID) WayID() WayID {
	if id&wayMask != wayMask {
		panic(fmt.Sprintf("not a way: %v", id))
	}

	return WayID(id.Ref())
}

// RelationID returns the id of this feature as a relation id.
// The function will panic if this element is not of TypeRelation.
func (id ElementID) RelationID() RelationID {
	if int64(id)&relationMask != relationMask {
		panic(fmt.Sprintf("not a relation: %v", id))
	}

	return RelationID(id.Ref())
}
