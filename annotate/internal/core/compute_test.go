package core

import (
	"reflect"
	"testing"

	"github.com/pchchv/osm"
)

func compareParents(t *testing.T, parents []Parent, expected []ChildList) {
	t.Helper()
	for i, p := range parents {
		parent := p.(*testParent)
		if expected[i] == nil {
			if parent.children != nil {
				t.Errorf("expected no children for %d", i)
				t.Logf("got: %+v", parent.children)
			}

			continue
		}

		if expected[i] != nil && parent.children == nil {
			t.Errorf("got no children for %d", i)
			t.Logf("expected: %+v", expected[i])
			continue
		}

		if parent.children[0] != expected[i][0] {
			t.Errorf("incorrect at parent %d", i)
			t.Logf("%+v", parent.children)
			t.Logf("%+v", expected[i])
		}
	}
}

func compareUpdates(t *testing.T, updates, expected []osm.Updates) {
	t.Helper()
	if !reflect.DeepEqual(updates, expected) {
		t.Errorf("updates not equal")
		if len(updates) != len(expected) {
			// length should be the length of the parents
			t.Fatalf("length of updates mismatch, %d != %d", len(updates), len(expected))
		}

		for i := range updates {
			if !reflect.DeepEqual(updates[i], expected[i]) {
				t.Errorf("index %d not equal", i)
				for j := range updates[i] {
					if !reflect.DeepEqual(updates[i][j], expected[i][j]) {
						t.Errorf("sub-index %d not equal", j)
						t.Logf("%+v", updates[i][j])
						t.Logf("%+v", expected[i][j])
					}
				}
			}
		}
	}
}
