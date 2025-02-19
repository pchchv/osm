package osm

const (
	// Constants for the different object types.
	TypeNode      Type = "node"
	TypeWay       Type = "way"
	TypeRelation  Type = "relation"
	TypeChangeset Type = "changeset"
	TypeNote      Type = "note"
	TypeUser      Type = "user"
	TypeBounds    Type = "bounds"

	versionBits   = 16
	versionMask   = 0x000000000000FFFF
	refMask       = 0x00FFFFFFFFFF0000
	featureMask   = 0x7FFFFFFFFFFF0000
	typeMask      = 0x7F00000000000000
	boundsMask    = 0x0800000000000000
	nodeMask      = 0x1000000000000000
	wayMask       = 0x2000000000000000
	relationMask  = 0x3000000000000000
	changesetMask = 0x4000000000000000
	noteMask      = 0x5000000000000000
	userMask      = 0x6000000000000000
)

// Type is the type of different osm objects,
// ie. node, way, relation, changeset, note, user.
type Type string
