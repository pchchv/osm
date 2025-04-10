package annotate

import (
	"context"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/planar"
	"github.com/pchchv/osm"
	"github.com/pchchv/osm/annotate/internal/core"
)

type wayChildDatasource struct {
	NodeHistoryAsChildrenDatasourcer
}

func (wds *wayChildDatasource) Get(ctx context.Context, id osm.FeatureID) (core.ChildList, error) {
	if id.Type() != osm.TypeNode {
		panic("only node types supported")
	}

	return wds.NodeHistoryAsChildren(ctx, id.NodeID())
}

type wayDatasource struct {
	NodeHistoryDatasourcer
}

// IsReverse checks if the update to this way was “reversal”.
// This is very tricky to answer in the general case,
// but easier for a minor update to a relation.
// Since the relation has not been updated,
// assume things are still connected and may just check the endpoints.
func IsReverse(w1, w2 *osm.Way) bool {
	if len(w1.Nodes) < 2 || len(w2.Nodes) < 2 {
		return false
	}

	// check if either is a ring
	if w1.Nodes[0].ID == w1.Nodes[len(w1.Nodes)-1].ID || w2.Nodes[0].ID == w2.Nodes[len(w2.Nodes)-1].ID {
		r1 := geo.Ring(w1.LineString())
		r2 := geo.Ring(w2.LineString())
		return planar.Area(r1)*planar.Area(r2) < 0
	}

	// not a ring so see if endpoint were flipped
	return w1.Nodes[0].ID == w2.Nodes[len(w2.Nodes)-1].ID && w2.Nodes[0].ID == w1.Nodes[len(w1.Nodes)-1].ID
}
