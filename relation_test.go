package osm

import (
	"reflect"
	"testing"
)

func TestMember_ids(t *testing.T) {
	cases := []struct {
		name string
		m    Member
		fid  FeatureID
		eid  ElementID
	}{
		{
			name: "node",
			m:    Member{Type: TypeNode, Ref: 12, Version: 2},
			fid:  NodeID(12).FeatureID(),
			eid:  NodeID(12).ElementID(2),
		},
		{
			name: "way",
			m:    Member{Type: TypeWay, Ref: 12, Version: 2},
			fid:  WayID(12).FeatureID(),
			eid:  WayID(12).ElementID(2),
		},
		{
			name: "relation",
			m:    Member{Type: TypeRelation, Ref: 12, Version: 2},
			fid:  RelationID(12).FeatureID(),
			eid:  RelationID(12).ElementID(2),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if id := tc.m.FeatureID(); id != tc.fid {
				t.Errorf("incorrect feature id: %v", id)
			}

			if id := tc.m.ElementID(); id != tc.eid {
				t.Errorf("incorrect element id: %v", id)
			}
		})
	}
}
func TestMembers_ids(t *testing.T) {
	ms := Members{
		{Type: TypeNode, Ref: 1, Version: 3},
		{Type: TypeWay, Ref: 2, Version: 4},
		{Type: TypeRelation, Ref: 3, Version: 5},
	}

	eids := ElementIDs{
		NodeID(1).ElementID(3),
		WayID(2).ElementID(4),
		RelationID(3).ElementID(5),
	}
	if ids := ms.ElementIDs(); !reflect.DeepEqual(ids, eids) {
		t.Errorf("incorrect element ids: %v", ids)
	}

	fids := FeatureIDs{
		NodeID(1).FeatureID(),
		WayID(2).FeatureID(),
		RelationID(3).FeatureID(),
	}
	if ids := ms.FeatureIDs(); !reflect.DeepEqual(ids, fids) {
		t.Errorf("incorrect feature ids: %v", ids)
	}
}
