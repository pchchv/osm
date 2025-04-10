package annotate

import (
	"context"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/planar"
	"github.com/pchchv/osm"
	"github.com/pchchv/osm/annotate/internal/core"
	"github.com/pchchv/osm/annotate/shared"
)

type relationDatasource struct {
	osm.HistoryDatasourcer
}

func (rds *relationDatasource) Get(ctx context.Context, id osm.FeatureID) (core.ChildList, error) {
	switch id.Type() {
	case osm.TypeNode:
		nodes, err := rds.NodeHistory(ctx, id.NodeID())
		if err != nil {
			return nil, err
		}

		return nodesToChildList(nodes), nil
	case osm.TypeWay:
		ways, err := rds.WayHistory(ctx, id.WayID())
		if err != nil {
			return nil, err
		}

		return waysToChildList(ways), nil
	case osm.TypeRelation:
		relations, err := rds.RelationHistory(ctx, id.RelationID())
		if err != nil {
			return nil, err
		}

		return relationsToChildList(relations), nil
	}

	return nil, &UnsupportedMemberTypeError{
		MemberType: id.Type(),
	}
}

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

func newWayDatasourcer(ds NodeHistoryDatasourcer) core.Datasourcer {
	if d, ok := ds.(NodeHistoryAsChildrenDatasourcer); ok {
		return &wayChildDatasource{d}
	}

	return &wayDatasource{ds}
}

func (wds *wayDatasource) Get(ctx context.Context, id osm.FeatureID) (core.ChildList, error) {
	if id.Type() != osm.TypeNode {
		panic("only node types supported")
	}

	nodes, err := wds.NodeHistory(ctx, id.NodeID())
	if err != nil {
		return nil, err
	}

	return nodesToChildList(nodes), nil
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

func relationsToChildList(relations osm.Relations) core.ChildList {
	if len(relations) == 0 {
		return nil
	}

	list := make(core.ChildList, len(relations))
	relations.SortByIDVersion()
	for i, r := range relations {
		c := shared.FromRelation(r)
		c.VersionIndex = i
		list[i] = c
	}

	return list
}

func nodesToChildList(nodes osm.Nodes) core.ChildList {
	if len(nodes) == 0 {
		return nil
	}

	list := make(core.ChildList, len(nodes))
	nodes.SortByIDVersion()
	for i, n := range nodes {
		c := shared.FromNode(n)
		c.VersionIndex = i
		list[i] = c
	}

	return list
}

func waysToChildList(ways osm.Ways) core.ChildList {
	if len(ways) == 0 {
		return nil
	}

	list := make(core.ChildList, len(ways))
	ways.SortByIDVersion()
	for i, w := range ways {
		c := shared.FromWay(w)
		c.VersionIndex = i

		if i != 0 {
			c.ReverseOfPrevious = IsReverse(w, ways[i-1])
		}

		list[i] = c
	}

	return list
}
