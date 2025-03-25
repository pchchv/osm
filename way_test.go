package osm

import (
	"reflect"
	"testing"

	"github.com/pchchv/geo"
)

func TestWayNode_ids(t *testing.T) {
	wn := WayNode{ID: 12, Version: 2}
	if id := wn.FeatureID(); id != NodeID(12).FeatureID() {
		t.Errorf("incorrect feature id: %v", id)
	}

	if id := wn.ElementID(); id != NodeID(12).ElementID(2) {
		t.Errorf("incorrect element id: %v", id)
	}
}

func TestWayNode_Point(t *testing.T) {
	wn := WayNode{ID: 12, Version: 2, Lon: 1, Lat: 2}
	p := wn.Point()
	if p.Lon() != 1 {
		t.Errorf("incorrect point lon: %v", p)
	}

	if p.Lat() != 2 {
		t.Errorf("incorrect point lat: %v", p)
	}
}

func TestWayNodes_Bounds(t *testing.T) {
	wn := WayNodes{
		{Lat: 1, Lon: 2},
		{Lat: 3, Lon: 4},
		{Lat: 2, Lon: 3},
	}
	b := wn.Bounds()
	if !reflect.DeepEqual(b, &Bounds{1, 3, 2, 4}) {
		t.Errorf("incorrect bounds: %v", b)
	}
}

func TestWayNodes_Bound(t *testing.T) {
	wn := WayNodes{
		{Lat: 1, Lon: 2},
		{Lat: 3, Lon: 4},
		{Lat: 2, Lon: 3},
	}
	b := wn.Bound()
	if !reflect.DeepEqual(b, geo.Bound{Min: geo.Point{2, 1}, Max: geo.Point{4, 3}}) {
		t.Errorf("incorrect bound: %v", b)
	}
}

func TestWayNodes_ids(t *testing.T) {
	wns := WayNodes{
		{ID: 1, Version: 3},
		{ID: 2, Version: 4},
	}

	eids := ElementIDs{NodeID(1).ElementID(3), NodeID(2).ElementID(4)}
	if ids := wns.ElementIDs(); !reflect.DeepEqual(ids, eids) {
		t.Errorf("incorrect element ids: %v", ids)
	}

	fids := FeatureIDs{NodeID(1).FeatureID(), NodeID(2).FeatureID()}
	if ids := wns.FeatureIDs(); !reflect.DeepEqual(ids, fids) {
		t.Errorf("incorrect feature ids: %v", ids)
	}

	nids := []NodeID{NodeID(1), NodeID(2)}
	if ids := wns.NodeIDs(); !reflect.DeepEqual(ids, nids) {
		t.Errorf("incorrect node ids: %v", nids)
	}
}

func TestWayNodes_UnmarshalJSON(t *testing.T) {
	wn := WayNodes{}
	if err := wn.UnmarshalJSON([]byte("[asdf,]")); err == nil {
		t.Errorf("should return error when json is invalid")
	}

	json := []byte(`[1,2,3,4]`)
	if err := wn.UnmarshalJSON(json); err != nil {
		t.Errorf("unmarshal error: %e", err)
	}

	expected := []NodeID{1, 2, 3, 4}
	if ids := wn.NodeIDs(); !reflect.DeepEqual(ids, expected) {
		t.Errorf("incorrect ids: %v", ids)
	}
}
