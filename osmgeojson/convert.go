package osmgeojson

import "github.com/pchchv/osm"

type relationSummary struct {
	ID   osm.RelationID    `json:"id"`
	Role string            `json:"role"`
	Tags map[string]string `json:"tags"`
}

type context struct {
	noID                   bool
	noMeta                 bool
	noRelationMembership   bool
	includeInvalidPolygons bool
	osm                    *osm.OSM
	wayMap                 map[osm.WayID]*osm.Way
	skippable              map[osm.WayID]struct{}
	wayMember              map[osm.NodeID]struct{}
	nodeMap                map[osm.NodeID]*osm.Node
	relationMember         map[osm.FeatureID][]*relationSummary
}
