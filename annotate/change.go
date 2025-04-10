package annotate

import (
	"context"

	"github.com/pchchv/osm"
)

func findPreviousNode(ctx context.Context, n *osm.Node, ds osm.HistoryDatasourcer, ignoreMissing bool) (*osm.Node, error) {
	nodes, err := ds.NodeHistory(ctx, n.ID)
	if err != nil {
		return nil, err
	}

	loc, max := -1, -1
	for i, node := range nodes {
		if v := node.Version; v < n.Version && v > max {
			loc, max = i, v
		}
	}

	if loc == -1 {
		// no version before ours
		if ignoreMissing {
			return nil, nil
		}
		return nil, &NoVisibleChildError{ID: n.FeatureID()}
	}

	return nodes[loc], nil
}

func findPreviousWay(ctx context.Context, w *osm.Way, ds osm.HistoryDatasourcer, ignoreMissing bool) (*osm.Way, error) {
	ways, err := ds.WayHistory(ctx, w.ID)
	if err != nil {
		return nil, err
	}

	loc, max := -1, -1
	for i, way := range ways {
		if v := way.Version; v < w.Version && v > max {
			loc, max = i, v
		}
	}

	if loc == -1 {
		// no version before ours
		if ignoreMissing {
			return nil, nil
		}
		return nil, &NoVisibleChildError{ID: w.FeatureID()}
	}

	return ways[loc], nil
}

func findPreviousRelation(ctx context.Context, r *osm.Relation, ds osm.HistoryDatasourcer, ignoreMissing bool) (*osm.Relation, error) {
	relations, err := ds.RelationHistory(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	loc, max := -1, -1
	for i, relation := range relations {
		if v := relation.Version; v < r.Version && v > max {
			loc, max = i, v
		}
	}

	if loc == -1 {
		// no version before ours
		if ignoreMissing {
			return nil, nil
		}
		return nil, &NoVisibleChildError{ID: r.FeatureID()}
	}

	return relations[loc], nil
}

func osmCount(o *osm.OSM) int {
	if o == nil {
		return 0
	}

	return len(o.Nodes) + len(o.Ways) + len(o.Relations)
}

func checkErr(ds osm.HistoryDatasourcer, ignoreMissing bool, err error, id osm.FeatureID) error {
	if err != nil && ds.NotFound(err) {
		if ignoreMissing {
			return nil
		}

		return &NoVisibleChildError{ID: id}
	}

	return nil
}
