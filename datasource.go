package osm

import (
	"context"
	"errors"
)

var errNotFound = errors.New("osm: feature not found")

// HistoryDatasource wraps maps to
// implement the HistoryDataSource interface.
type HistoryDatasource struct {
	Nodes     map[NodeID]Nodes
	Ways      map[WayID]Ways
	Relations map[RelationID]Relations
}

// NodeHistory returns the history for the given id from the map.
func (ds *HistoryDatasource) NodeHistory(ctx context.Context, id NodeID) (Nodes, error) {
	if ds.Nodes == nil {
		return nil, errNotFound
	}

	v := ds.Nodes[id]
	if v == nil {
		return nil, errNotFound
	}

	return v, nil
}

// WayHistory returns the history for the given id from the map.
func (ds *HistoryDatasource) WayHistory(ctx context.Context, id WayID) (Ways, error) {
	if ds.Ways == nil {
		return nil, errNotFound
	}

	v := ds.Ways[id]
	if v == nil {
		return nil, errNotFound
	}

	return v, nil
}

// RelationHistory returns the history for the given id from the map.
func (ds *HistoryDatasource) RelationHistory(ctx context.Context, id RelationID) (Relations, error) {
	if ds.Relations == nil {
		return nil, errNotFound
	}

	v := ds.Relations[id]
	if v == nil {
		return nil, errNotFound
	}

	return v, nil
}
