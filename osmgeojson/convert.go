package osmgeojson

import "github.com/pchchv/osm"

type relationSummary struct {
	ID   osm.RelationID    `json:"id"`
	Role string            `json:"role"`
	Tags map[string]string `json:"tags"`
}
