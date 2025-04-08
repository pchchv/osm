package annotate

import (
	"context"

	"github.com/pchchv/osm"
)

var _ RelationHistoryDatasourcer = &osm.HistoryDatasource{}

// RelationHistoryDatasourcer is an more strict interface for when we only need the relation history.
type RelationHistoryDatasourcer interface {
	RelationHistory(context.Context, osm.RelationID) (osm.Relations, error)
	NotFound(error) bool
}
