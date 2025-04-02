package mputil

import (
	"reflect"
	"testing"
	"time"

	"github.com/pchchv/geo"
	"github.com/pchchv/osm"
)

func TestMultiSegment_Ring_noAnnotation(t *testing.T) {
	cases := []struct {
		name        string
		orientation geo.Orientation
		input       MultiSegment
		output      geo.Ring
	}{
		{
			name:        "ring is direction requested",
			orientation: geo.CW,
			input: MultiSegment{
				{
					Line: geo.LineString{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
				},
			},
			output: geo.Ring{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
		},
		{
			name:        "ring opposite direction of requested",
			orientation: geo.CW,
			input: MultiSegment{
				{
					Line: geo.LineString{{0, 0}, {1, 0}, {1, 1}, {0, 0}},
				},
			},
			output: geo.Ring{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
		},
		{
			name:        "multi segments in direction of requested",
			orientation: geo.CW,
			input: MultiSegment{
				{
					Line: geo.LineString{{0, 0}, {1, 1}},
				},
				{
					Line: geo.LineString{{1, 0}, {0, 0}},
				},
			},
			output: geo.Ring{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
		},
		{
			name:        "multi segments in opposite direction of requested",
			orientation: geo.CW,
			input: MultiSegment{
				{
					Line: geo.LineString{{0, 0}, {1, 0}},
				},
				{
					Line: geo.LineString{{1, 1}, {0, 0}},
				},
			},
			output: geo.Ring{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testRing(t, tc.input, tc.output, tc.orientation)
		})
	}
}

func TestMultiSegment_Ring_annotation(t *testing.T) {
	cases := []struct {
		name        string
		orientation geo.Orientation
		input       MultiSegment
		output      geo.Ring
	}{
		{
			name:        "ring is direction requested",
			orientation: geo.CW,
			input: MultiSegment{
				{
					Orientation: geo.CW,
					Line:        geo.LineString{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
				},
			},
			output: geo.Ring{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
		},
		{
			name:        "ring opposite direction of requested",
			orientation: geo.CW,
			input: MultiSegment{
				{
					Orientation: geo.CCW,
					Line:        geo.LineString{{0, 0}, {1, 0}, {1, 1}, {0, 0}},
				},
			},
			output: geo.Ring{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
		},
		{
			name:        "multi segments in direction of requested",
			orientation: geo.CW,
			input: MultiSegment{
				{
					Orientation: geo.CW,
					Line:        geo.LineString{{0, 0}, {1, 1}},
				},
				{
					Orientation: geo.CW,
					Line:        geo.LineString{{1, 0}, {0, 0}},
				},
			},
			output: geo.Ring{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
		},
		{
			name:        "multi segments in opposite direction of requested",
			orientation: geo.CW,
			input: MultiSegment{
				{
					Orientation: geo.CCW,
					Line:        geo.LineString{{0, 0}, {1, 0}},
				},
				{
					Orientation: geo.CCW,
					Line:        geo.LineString{{1, 1}, {0, 0}},
				},
			},
			output: geo.Ring{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
		},
		{
			name:        "reversed to correct",
			orientation: geo.CW,
			input: MultiSegment{
				{
					Orientation: geo.CW,
					Line:        geo.LineString{{0, 0}, {1, 1}},
				},
				{
					Orientation: geo.CCW,
					Reversed:    true,
					Line:        geo.LineString{{1, 0}, {0, 0}},
				},
			},
			output: geo.Ring{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
		},
		{
			name:        "reversed to wrong direction",
			orientation: geo.CW,
			input: MultiSegment{
				{
					Orientation: geo.CCW,
					Line:        geo.LineString{{0, 0}, {1, 0}},
				},
				{
					Orientation: geo.CW,
					Reversed:    true,
					Line:        geo.LineString{{1, 1}, {0, 0}},
				},
			},
			output: geo.Ring{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testRing(t, tc.input, tc.output, tc.orientation)
		})
	}
}

func TestMultiSegment_Orientation(t *testing.T) {
	ms := MultiSegment{
		{
			Line: geo.LineString{{0, 0}, {1, 0}},
		},
		{
			Line: geo.LineString{{1, 1}, {0, 1}},
		},
	}

	if o := ms.Orientation(); o != geo.CCW {
		t.Errorf("incorrect orientation: %v != %v", o, geo.CCW)
	}
}

func TestMultiSegment_LineString(t *testing.T) {
	ms := MultiSegment{
		{
			Line: geo.LineString{{1, 1}, {2, 2}},
		},
		{
			Line: geo.LineString{{3, 3}, {4, 4}},
		},
	}

	ls := ms.LineString()
	expected := geo.LineString{{1, 1}, {2, 2}, {3, 3}, {4, 4}}
	if !ls.Equal(expected) {
		t.Errorf("incorrect line string: %v", ls)
	}
}

func TestGroup(t *testing.T) {
	members := osm.Members{
		{Type: osm.TypeNode, Ref: 1},
		{Type: osm.TypeWay, Ref: 1, Role: "outer", Orientation: geo.CW},
		{Type: osm.TypeWay, Ref: 2, Role: "inner", Orientation: geo.CCW},
		{Type: osm.TypeWay, Ref: 3, Role: "inner", Orientation: geo.CCW},
		{Type: osm.TypeRelation, Ref: 3},
	}
	ways := map[osm.WayID]*osm.Way{
		1: {ID: 1, Nodes: osm.WayNodes{
			{Lat: 1.0, Lon: 2.0},
			{Lat: 2.0, Lon: 3.0},
		}},
		2: {ID: 1, Nodes: osm.WayNodes{
			{Lat: 3.0, Lon: 4.0},
			{Lat: 4.0, Lon: 5.0},
		}},
	}
	outer, inner, tainted := Group(members, ways, time.Time{})
	if !tainted {
		t.Errorf("should be tainted")
	}

	// outer
	expected := []Segment{
		{
			Index: 1, Orientation: geo.CW, Reversed: true,
			Line: geo.LineString{{3, 2}, {2, 1}},
		},
	}
	if !reflect.DeepEqual(outer, expected) {
		t.Errorf("incorrect outer: %+v", inner)
	}

	// inner
	expected = []Segment{
		{
			Index: 2, Orientation: geo.CCW, Reversed: true,
			Line: geo.LineString{{5, 4}, {4, 3}},
		},
	}
	if !reflect.DeepEqual(inner, expected) {
		t.Errorf("incorrect inner: %+v", inner)
	}
}

func TestGroup_zeroLengthWays(t *testing.T) {
	// should not panic
	Group(
		osm.Members{
			{Type: osm.TypeWay, Ref: 1, Role: "outer", Orientation: geo.CW},
			{Type: osm.TypeWay, Ref: 1, Role: "inner", Orientation: geo.CCW},
		},
		map[osm.WayID]*osm.Way{
			1: {ID: 1},
		},
		time.Time{},
	)
}

func testRing(t testing.TB, input MultiSegment, expected geo.Ring, orient geo.Orientation) {
	t.Helper()
	ring := input.Ring(orient)
	if o := ring.Orientation(); o != orient {
		t.Errorf("different orientation: %v != %v", o, orient)
	}

	if !ring.Equal(expected) {
		t.Errorf("wrong ring")
		t.Logf("%v", ring)
		t.Logf("%v", expected)
	}

	// with reverse orientation
	orient *= -1
	expected.Reverse()
	ring = input.Ring(orient)
	if o := ring.Orientation(); o != orient {
		t.Errorf("reversed, different orientation: %v != %v", o, orient)
	}

	if !ring.Equal(expected) {
		t.Errorf("reversed, wrong ring")
		t.Logf("%v", ring)
		t.Logf("%v", expected)
	}
}
