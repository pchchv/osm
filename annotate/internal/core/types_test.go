package core

import (
	"time"

	"github.com/pchchv/osm"
	"github.com/pchchv/osm/annotate/shared"
)

var _ Parent = &testParent{}

type testParent struct {
	changesetID osm.ChangesetID
	version     int
	visible     bool
	timestamp   time.Time
	committed   time.Time
	refs        osm.FeatureIDs
	children    ChildList
}

func (t testParent) ID() osm.FeatureID {
	return osm.FeatureID(0) // this is only used for logging.
}

func (t testParent) ChangesetID() osm.ChangesetID {
	return t.changesetID
}

func (t testParent) Version() int {
	return t.version
}

func (t testParent) Visible() bool {
	return t.visible
}

func (t testParent) Timestamp() time.Time {
	return t.timestamp
}

func (t testParent) Committed() time.Time {
	return t.committed
}

func (t testParent) Refs() (osm.FeatureIDs, []bool) {
	annotated := make([]bool, len(t.refs))
	for i := range annotated {
		annotated[i] = true
	}
	return t.refs, annotated
}

func (t *testParent) SetChild(idx int, c *shared.Child) {
	if idx >= len(t.children) {
		nc := make(ChildList, idx+1)
		copy(nc, t.children)
		t.children = nc
	}
	t.children[idx] = c
}
