package osm

import (
	"fmt"
	"sort"
)

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

// FeatureID is an identifier for a feature in OSM.
// It is meant to represent all the versions of a given element.
type FeatureID int64

// Type returns the Type of the feature, or empty string for invalid type.
func (id FeatureID) Type() Type {
	switch id & typeMask {
	case nodeMask:
		return TypeNode
	case wayMask:
		return TypeWay
	case relationMask:
		return TypeRelation
	}

	return ""
}

// String returns "type/ref" for the feature.
func (id FeatureID) String() string {
	t := Type("unknown")
	switch id & typeMask {
	case nodeMask:
		t = TypeNode
	case wayMask:
		t = TypeWay
	case relationMask:
		t = TypeRelation
	}
	return fmt.Sprintf("%s/%d", t, id.Ref())
}

// Ref return the ID reference for the feature.
// Not unique without the type.
func (id FeatureID) Ref() int64 {
	return int64((id & refMask) >> versionBits)
}

// FeatureIDs is a slice of FeatureIDs with some helpers on top.
type FeatureIDs []FeatureID

// Counts returns the number of each type of feature in the set of ids.
func (ids FeatureIDs) Counts() (nodes, ways, relations int) {
	for _, id := range ids {
		switch id.Type() {
		case TypeNode:
			nodes++
		case TypeWay:
			ways++
		case TypeRelation:
			relations++
		}
	}
	return
}

type featureIDsSort FeatureIDs

// Sort will order the ids by type, node, way, relation, changeset and then id.
func (ids FeatureIDs) Sort() {
	sort.Sort(featureIDsSort(ids))
}

func (ids featureIDsSort) Len() int {
	return len(ids)
}

func (ids featureIDsSort) Swap(i, j int) {
	ids[i], ids[j] = ids[j], ids[i]
}

func (ids featureIDsSort) Less(i, j int) bool {
	return ids[i] < ids[j]
}
