package annotate

import (
	"context"

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
