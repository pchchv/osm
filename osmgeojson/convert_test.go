package osmgeojson

import (
	"encoding/json"
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/geojson"
	"github.com/pchchv/osm"
)

type rawFC struct {
	Type     string        `json:"type"`
	Features []interface{} `json:"features"`
}

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

func jsonLoop(t *testing.T, fc *geojson.FeatureCollection) *rawFC {
	data, err := json.Marshal(fc)
	if err != nil {
		t.Fatalf("unable to marshal fc: %e", err)
	}

	result := &rawFC{}
	if err = json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unable to unmarshal data: %e", err)
	}

	return result
}

func jsonMarshalIndent(t *testing.T, raw interface{}) string {
	data, err := json.MarshalIndent(raw, "", " ")
	if err != nil {
		t.Fatalf("unable to marshal json: %e", err)
	}

	return string(data)
}

func testConvert(t *testing.T, rawXML string, expected *geojson.FeatureCollection, opts ...Option) {
	t.Helper()
	o := &osm.OSM{}
	if err := xml.Unmarshal([]byte(rawXML), &o); err != nil {
		t.Fatalf("failed to unmarshal xml: %e", err)
	}

	// clean up expected a bit
	for _, f := range expected.Features {
		if f.Properties["tags"] == nil {
			f.Properties["tags"] = map[string]string{}
		}

		if f.Properties["meta"] == nil {
			f.Properties["meta"] = map[string]interface{}{}
		}

		if f.Properties["relations"] == nil {
			f.Properties["relations"] = []*relationSummary{}
		} else {
			for _, rs := range f.Properties["relations"].([]*relationSummary) {
				if rs.Tags == nil {
					rs.Tags = map[string]string{}
				}
			}
		}
	}

	fc, err := Convert(o, opts...)
	if err != nil {
		t.Fatalf("convert error: %e", err)
	}

	raw := jsonLoop(t, fc)
	expectedRaw := jsonLoop(t, expected)
	if !reflect.DeepEqual(raw, expectedRaw) {
		if len(raw.Features) != len(expectedRaw.Features) {
			t.Logf("%v", jsonMarshalIndent(t, raw))
			t.Logf("%v", jsonMarshalIndent(t, expectedRaw))
			t.Errorf("not equal")
		} else {
			for i := range expectedRaw.Features {
				if !reflect.DeepEqual(raw.Features[i], expectedRaw.Features[i]) {
					t.Logf("%v", jsonMarshalIndent(t, raw.Features[i]))
					t.Logf("%v", jsonMarshalIndent(t, expectedRaw.Features[i]))
					t.Errorf("Feature %d not equal", i)
				}
			}
		}
	}
}
