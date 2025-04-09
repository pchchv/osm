package annotate

import "github.com/pchchv/osm"

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
