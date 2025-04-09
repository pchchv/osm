package annotate

import (
	"time"

	"github.com/pchchv/osm"
	"github.com/pchchv/osm/annotate/shared"
)

// parentWay wraps a osm.Way into the
// core.Parent interface so that updates can be computed.
type parentWay struct {
	Way *osm.Way
}

func (w *parentWay) Version() int {
	return w.Way.Version
}

func (w *parentWay) ID() osm.FeatureID {
	return w.Way.FeatureID()
}

func (w *parentWay) ChangesetID() osm.ChangesetID {
	return w.Way.ChangesetID
}

func (w *parentWay) Timestamp() time.Time {
	return w.Way.Timestamp
}

func (w *parentWay) Committed() time.Time {
	if w.Way.Committed == nil {
		return time.Time{}
	}

	return *w.Way.Committed
}

func (w *parentWay) Visible() bool {
	return w.Way.Visible
}

func (w *parentWay) SetChild(idx int, child *shared.Child) {
	if child == nil {
		return
	}

	w.Way.Nodes[idx].Version = child.Version
	w.Way.Nodes[idx].ChangesetID = child.ChangesetID
	w.Way.Nodes[idx].Lat = child.Lat
	w.Way.Nodes[idx].Lon = child.Lon
}

func (w *parentWay) Refs() (osm.FeatureIDs, []bool) {
	ids := make(osm.FeatureIDs, len(w.Way.Nodes))
	annotated := make([]bool, len(w.Way.Nodes))
	for i := range w.Way.Nodes {
		ids[i] = w.Way.Nodes[i].FeatureID()
		annotated[i] = w.Way.Nodes[i].Version != 0
	}

	return ids, annotated
}
