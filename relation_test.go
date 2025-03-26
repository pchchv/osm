package osm

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/pchchv/geo"
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

func TestRelation_ids(t *testing.T) {
	r := Relation{ID: 12, Version: 2}
	if id := r.FeatureID(); id != RelationID(12).FeatureID() {
		t.Errorf("incorrect feature id: %v", id)
	}

	if id := r.ElementID(); id != RelationID(12).ElementID(2) {
		t.Errorf("incorrect element id: %v", id)
	}
}

func TestRelation_MarshalJSON(t *testing.T) {
	r := Relation{
		ID: 123,
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if !bytes.Equal(data, []byte(`{"type":"relation","id":123,"visible":false,"timestamp":"0001-01-01T00:00:00Z","members":[]}`)) {
		t.Errorf("incorrect json: %v", string(data))
	}

	// with members
	r = Relation{
		ID:      123,
		Members: Members{{Type: "node", Ref: 123, Role: "outer", Version: 1}},
	}

	data, err = json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if !bytes.Equal(data, []byte(`{"type":"relation","id":123,"visible":false,"timestamp":"0001-01-01T00:00:00Z","members":[{"type":"node","ref":123,"role":"outer","version":1}]}`)) {
		t.Errorf("incorrect json: %v", string(data))
	}
}

func TestRelation_ApplyUpdatesUpTo(t *testing.T) {
	updates := Updates{
		{Index: 0, Timestamp: time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC), Version: 11},
		{Index: 1, Timestamp: time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC), Version: 12},
		{Index: 2, Timestamp: time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC), Version: 13, Lat: 10, Lon: 20},
	}
	r := Relation{
		ID:      123,
		Members: Members{{Version: 1}, {Version: 2}, {Version: 3}},
	}

	r.Updates = updates
	r.ApplyUpdatesUpTo(time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC))
	if r.Members[0].Version != 1 || r.Members[1].Version != 2 || r.Members[2].Version != 3 {
		t.Errorf("incorrect members, got %v", r.Members)
	}

	r.Updates = updates
	r.ApplyUpdatesUpTo(time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC))
	if r.Members[0].Version != 11 || r.Members[1].Version != 2 || r.Members[2].Version != 13 {
		t.Errorf("incorrect members, got %v", r.Members)
	}

	if r.Members[2].Lat != 10 {
		t.Errorf("did not apply lat data")
	}

	if r.Members[2].Lon != 20 {
		t.Errorf("did not apply lon data")
	}

	if l := len(r.Updates); l != 1 {
		t.Errorf("incorrect number of updates: %v", l)
	}

	if r.Updates[0].Index != 1 {
		t.Errorf("incorrect updates: %v", r.Updates)
	}
}

func TestRelation_ApplyUpdate(t *testing.T) {
	r := Relation{
		ID:      123,
		Members: Members{{Ref: 1, Type: TypeWay, Orientation: geo.CW}},
	}
	err := r.applyUpdate(Update{
		Index:       0,
		Version:     1,
		ChangesetID: 2,
		Reverse:     true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := Member{
		Ref:         1,
		Type:        TypeWay,
		Version:     1,
		ChangesetID: 2,
		Orientation: geo.CCW,
	}

	if !reflect.DeepEqual(r.Members[0], expected) {
		t.Errorf("incorrect update, got %v", r.Members[0])
	}
}

func TestRelation_ApplyUpdate_error(t *testing.T) {
	r := Relation{
		ID:      123,
		Members: Members{{Ref: 1, Type: TypeNode}},
	}
	err := r.applyUpdate(Update{
		Index: 1,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if e, ok := err.(*UpdateIndexOutOfRangeError); !ok {
		t.Errorf("incorrect error, got %v", e)
	}
}
