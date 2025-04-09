package annotate

import (
	"time"

	"github.com/pchchv/osm"
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
