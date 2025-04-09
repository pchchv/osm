package annotate

import (
	"context"
	"time"

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

func (r *parentRelation) Timestamp() time.Time {
	return r.Relation.Timestamp
}

func (r *parentRelation) Committed() time.Time {
	if r.Relation.Committed == nil {
		return time.Time{}
	}

	return *r.Relation.Committed
}

func (r *parentRelation) Visible() bool {
	return r.Relation.Visible
}

func (r *parentRelation) SetChild(idx int, child *shared.Child) {
	if r.Relation.Polygon() && r.ways == nil {
		r.ways = make(map[osm.WayID]*osm.Way, len(r.Relation.Members))
	}

	if child == nil {
		return
	}

	r.Relation.Members[idx].Version = child.Version
	r.Relation.Members[idx].ChangesetID = child.ChangesetID
	r.Relation.Members[idx].Lat = child.Lat
	r.Relation.Members[idx].Lon = child.Lon
	if r.ways != nil && child.Way != nil {
		r.ways[child.Way.ID] = child.Way
	}
}

func (r *parentRelation) Refs() (osm.FeatureIDs, []bool) {
	ids := make(osm.FeatureIDs, len(r.Relation.Members))
	annotated := make([]bool, len(r.Relation.Members))
	for i := range r.Relation.Members {
		ids[i] = r.Relation.Members[i].FeatureID()
		annotated[i] = r.Relation.Members[i].Version != 0
	}

	return ids, annotated
}
