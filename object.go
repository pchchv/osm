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

// ObjectIDs is a slice of ObjectIDs with some helpers on top.
type ObjectIDs []ObjectID

// Object represents a Node, Way, Relation, Changeset, Note or User only.
type Object interface {
	ObjectID() ObjectID
	private() // to ensure that **ID types do not implement this interface
}

// Objects is a set of objects with some helpers
type Objects []Object

// ObjectIDs returns a slice of the object ids of the osm objects.
func (os Objects) ObjectIDs() ObjectIDs {
	if len(os) == 0 {
		return nil
	}

	ids := make(ObjectIDs, 0, len(os))
	for _, o := range os {
		ids = append(ids, o.ObjectID())
	}

	return ids
}

// Scanner reads osm data from planet dump files.
// It is based on the bufio.Scanner, common usage.
// Scanners are not safe for parallel use.
// One should feed the objects into their
// own channel and have workers read from that.
//
//	s := scanner.New(r)
//	defer s.Close()
//
//	for s.Next() {
//		o := s.Object()
//		// do something
//	}
//
//	if s.Err() != nil {
//		// scanner did not complete fully
//	}
type Scanner interface {
	Scan() bool
	Object() Object
	Err() error
	Close() error
}
