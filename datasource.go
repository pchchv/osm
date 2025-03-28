package osm

// HistoryDatasource wraps maps to
// implement the HistoryDataSource interface.
type HistoryDatasource struct {
	Nodes     map[NodeID]Nodes
	Ways      map[WayID]Ways
	Relations map[RelationID]Relations
}
