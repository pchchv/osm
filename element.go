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
