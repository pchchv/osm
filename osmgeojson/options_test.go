package osmgeojson

import (
	"encoding/xml"
	"testing"

	"github.com/pchchv/geo/geojson"
	"github.com/pchchv/osm"
)

func convertXML(t *testing.T, data string, opts ...Option) *geojson.FeatureCollection {
	o := &osm.OSM{}
	err := xml.Unmarshal([]byte(data), &o)
	if err != nil {
		t.Fatalf("failed to unmarshal xml: %v", err)
	}

	fc, err := Convert(o, opts...)
	if err != nil {
		t.Fatalf("failed to convert: %v", err)
	}

	return fc
}
