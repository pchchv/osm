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

// Next locates the next relation id that can be used.
// Returns false if the context is closed,
// something went wrong or the full tree has been walked.
func (o *ChildFirstOrdering) Next() bool {
	if o.err != nil || o.ctx.Err() != nil {
		return false
	}

	select {
	case id := <-o.out:
		if id == 0 {
			return false
		}
		o.id = id
		return true
	case <-o.ctx.Done():
		return false
	}
}

// Close terminates the scanning process before all ids have been walked.
func (o *ChildFirstOrdering) Close() {
	o.done()
	o.wg.Wait()
}

// Err returns a non-nil error
// if something went wrong with search,
// like a loop, or a datasource error.
func (o *ChildFirstOrdering) Err() error {
	if o.err != nil {
		return o.err
	}

	return o.ctx.Err()
}

// RelationID returns the id found by the previous scan.
func (o *ChildFirstOrdering) RelationID() osm.RelationID {
	return o.id
}
