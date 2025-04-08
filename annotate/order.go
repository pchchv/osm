package annotate

import (
	"context"
	"sync"

	"github.com/pchchv/osm"
)

var _ RelationHistoryDatasourcer = &osm.HistoryDatasource{}

// RelationHistoryDatasourcer is a stricter interface for cases where only relationship history is needed.
type RelationHistoryDatasourcer interface {
	RelationHistory(context.Context, osm.RelationID) (osm.Relations, error)
	NotFound(error) bool
}

// ChildFirstOrdering allows to process a set of relations in a dept first order.
// Since relations can refer to other relations,
// it must be ensured that children are added before parents.
type ChildFirstOrdering struct {
	// CompletedIndex is the number of relation ids in the provided array that have been finished.
	// This can be used as a good restart position.
	CompletedIndex int
	ctx            context.Context
	done           context.CancelFunc
	ds             RelationHistoryDatasourcer
	visited        map[osm.RelationID]struct{}
	out            chan osm.RelationID
	wg             sync.WaitGroup
	id             osm.RelationID
	err            error
}
