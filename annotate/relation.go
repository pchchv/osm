package annotate

import (
	"context"

	"github.com/pchchv/osm"
	"github.com/pchchv/osm/annotate/shared"
)

// HistoryAsChildrenDatasourcer is an advanced data source that
// returns the needed elements as children directly.
type HistoryAsChildrenDatasourcer interface {
	osm.HistoryDatasourcer
	NodeHistoryAsChildren(context.Context, osm.NodeID) ([]*shared.Child, error)
	WayHistoryAsChildren(context.Context, osm.WayID) ([]*shared.Child, error)
	RelationHistoryAsChildren(context.Context, osm.RelationID) ([]*shared.Child, error)
}

// parentRelation wraps a osm.Relation into the
// core.Parent interface so that updates can be computed.
type parentRelation struct {
	Relation *osm.Relation
	ways     map[osm.WayID]*osm.Way
}

func (r *parentRelation) Version() int {
	return r.Relation.Version
}

func (r *parentRelation) ID() osm.FeatureID {
	return r.Relation.FeatureID()
}

func (r *parentRelation) ChangesetID() osm.ChangesetID {
	return r.Relation.ChangesetID
}
