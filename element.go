package osm

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
