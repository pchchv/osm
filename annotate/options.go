package annotate

import (
	"time"

	"github.com/pchchv/osm"
	"github.com/pchchv/osm/annotate/internal/core"
)

// Option is a parameter that can be used for annotating.
type Option func(*core.Options) error

// ChildFilter allows for only a subset of children to be annotated on the parent.
// This can greatly improve update speed by only worrying about the children updated in the same batch.
// All unannotated children will be annotated regardless of the results of the filter function.
func ChildFilter(filter func(osm.FeatureID) bool) Option {
	return func(o *core.Options) error {
		o.ChildFilter = filter
		return nil
	}
}

// Threshold is used if the "committed at" time is unknown and deals with the flexibility of commit orders,
// e.g. nodes in the same commit as the way can have a timestamp after the way.
// Threshold defines the time range to "forward group" these changes.
// Default 30 minutes.
func Threshold(t time.Duration) Option {
	return func(o *core.Options) error {
		o.Threshold = t
		return nil
	}
}
