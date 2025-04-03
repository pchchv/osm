package osmgeojson

import (
	"testing"

	"github.com/pchchv/geo"
	"github.com/pchchv/osm"
)

func TestBuildRouteLineString(t *testing.T) {
	ctx := &context{
		osm:       &osm.OSM{},
		skippable: map[osm.WayID]struct{}{},
		wayMap: map[osm.WayID]*osm.Way{
			2: {
				ID: 2,
				Nodes: osm.WayNodes{
					{ID: 1, Lat: 1, Lon: 2},
					{ID: 2},
					{ID: 3, Lat: 3, Lon: 4},
				},
			},
		},
	}
	relation := &osm.Relation{
		ID: 1,
		Members: osm.Members{
			{Type: osm.TypeNode, Ref: 1},
			{Type: osm.TypeWay, Ref: 2},
			{Type: osm.TypeWay, Ref: 3},
		},
	}
	feature := ctx.buildRouteLineString(relation)
	if !geo.Equal(feature.Geometry, geo.LineString{{2, 1}, {4, 3}}) {
		t.Errorf("incorrect geometry: %v", feature.Geometry)
	}

	relation = &osm.Relation{
		ID: 1,
		Members: osm.Members{
			{Type: osm.TypeWay, Ref: 20},
			{Type: osm.TypeWay, Ref: 30},
		},
	}
	feature = ctx.buildRouteLineString(relation)
	if feature != nil {
		t.Errorf("should not return feature if no ways present: %v", feature)
	}
}
