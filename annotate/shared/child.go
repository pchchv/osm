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
