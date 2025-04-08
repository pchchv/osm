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

// NewChildFirstOrdering creates a new ordering object.
// It is used to provided a child before parent ordering for relations.
// This order must be used when inserting+annotating relations into the datastore.
func NewChildFirstOrdering(ctx context.Context, ids []osm.RelationID, ds RelationHistoryDatasourcer) *ChildFirstOrdering {
	ctx, done := context.WithCancel(ctx)
	o := &ChildFirstOrdering{
		ctx:     ctx,
		done:    done,
		ds:      ds,
		visited: make(map[osm.RelationID]struct{}, len(ids)),
		out:     make(chan osm.RelationID),
	}

	o.wg.Add(1)
	go func() {
		defer o.wg.Done()
		defer close(o.out)

		path := make([]osm.RelationID, 0, 100)
		for i, id := range ids {
			err := o.walk(id, path)
			if err != nil {
				o.err = err
				return
			}

			o.CompletedIndex = i
		}
	}()

	return o
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

func (o *ChildFirstOrdering) walk(id osm.RelationID, path []osm.RelationID) error {
	if _, ok := o.visited[id]; ok {
		return nil
	}

	relations, err := o.ds.RelationHistory(o.ctx, id)
	if o.ds.NotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	for _, r := range relations {
		for _, m := range r.Members {
			if m.Type != osm.TypeRelation {
				continue
			}

			mid := osm.RelationID(m.Ref)
			for _, pid := range path {
				if pid == mid {
					// circular relations are allowed
					// (see https://github.com/openstreetmap/openstreetmap-website/issues/1465#issuecomment-282323187)
					// since this relation is already being worked out higher up the stack, it's okay to just come back here
					return nil
				}
			}

			if err := o.walk(mid, append(path, mid)); err != nil {
				return err
			}
		}
	}

	if o.ctx.Err() != nil {
		return o.ctx.Err()
	}

	o.visited[id] = struct{}{}
	select {
	case o.out <- id:
	case <-o.ctx.Done():
		return o.ctx.Err()
	}

	return nil
}
