package shared

import (
	"time"

	"github.com/pchchv/osm"
)

// Child represents a node, way or relation that is a
// dependent for annotating ways or relations.
type Child struct {
	ID                osm.FeatureID
	Version           int
	ChangesetID       osm.ChangesetID
	VersionIndex      int // sorted version index (versions do not have to start with 1 or be sequential)
	Timestamp         time.Time
	Committed         time.Time
	Lon               float64 // for nodes
	Lat               float64
	Way               *osm.Way // for ways
	ReverseOfPrevious bool
	Visible           bool
}

// FromNode converts a node to a child.
func FromNode(n *osm.Node) *Child {
	c := &Child{
		ID:          n.FeatureID(),
		Version:     n.Version,
		ChangesetID: n.ChangesetID,
		Visible:     n.Visible,
		Timestamp:   n.Timestamp,
		Lon:         n.Lon,
		Lat:         n.Lat,
	}

	if n.Committed != nil {
		c.Committed = *n.Committed
	}

	return c
}

// FromWay converts a way to a child.
func FromWay(w *osm.Way) *Child {
	c := &Child{
		ID:          w.FeatureID(),
		Version:     w.Version,
		ChangesetID: w.ChangesetID,
		Visible:     w.Visible,
		Timestamp:   w.Timestamp,
		Way:         w,
	}

	if w.Committed != nil {
		c.Committed = *w.Committed
	}

	return c
}

// FromRelation converts a way to a child.
func FromRelation(r *osm.Relation) *Child {
	c := &Child{
		ID:          r.FeatureID(),
		Version:     r.Version,
		ChangesetID: r.ChangesetID,
		Visible:     r.Visible,
		Timestamp:   r.Timestamp,
	}

	if r.Committed != nil {
		c.Committed = *r.Committed
	}

	return c
}
