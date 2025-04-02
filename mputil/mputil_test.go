package mputil

import (
	"testing"

	"github.com/pchchv/geo"
)

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
