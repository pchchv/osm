package osm

import "fmt"

// ObjectID encodes the type and ref of an osm object,
// e.g. nodes, ways, relations, changesets, notes and users.
type ObjectID int64

// Version returns the version of the object.
// Return 0 if the object doesn't have versions like users,
// notes and changesets.
func (id ObjectID) Version() int {
	return int(id & (versionMask))
}

// Type returns the Type of the object.
func (id ObjectID) Type() Type {
	switch id & typeMask {
	case nodeMask:
		return TypeNode
	case wayMask:
		return TypeWay
	case relationMask:
		return TypeRelation
	case changesetMask:
		return TypeChangeset
	case noteMask:
		return TypeNote
	case userMask:
		return TypeUser
	case boundsMask:
		return TypeBounds
	default:
		panic("unknown type")
	}
}

// Ref returns the ID reference for the object.
// Not unique without the type.
func (id ObjectID) Ref() int64 {
	return int64((id & refMask) >> versionBits)
}

// String returns "type/ref:version" for the object.
func (id ObjectID) String() string {
	if id.Version() == 0 {
		return fmt.Sprintf("%s/%d:-", id.Type(), id.Ref())
	}

	return fmt.Sprintf("%s/%d:%d", id.Type(), id.Ref(), id.Version())
}
