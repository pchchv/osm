package osmgeojson

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/pchchv/osm"
)

const benchGeoJSON = "../testdata/geojson_benchmark.osm"

func BenchmarkConvert(b *testing.B) {
	o := parseFile(b, benchGeoJSON)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := Convert(o); err != nil {
			b.Fatalf("convert error: %e", err)
		}
	}
}

func BenchmarkConvertAnnotated(b *testing.B) {
	o := parseFile(b, benchGeoJSON)
	annotate(o)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := Convert(o); err != nil {
			b.Fatalf("convert error: %e", err)
		}
	}
}

func BenchmarkConvert_NoID(b *testing.B) {
	o := parseFile(b, benchGeoJSON)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := Convert(o, NoID(true)); err != nil {
			b.Fatalf("convert error: %e", err)
		}
	}
}

func BenchmarkConvert_NoMeta(b *testing.B) {
	o := parseFile(b, benchGeoJSON)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := Convert(o, NoMeta(true)); err != nil {
			b.Fatalf("convert error: %e", err)
		}
	}
}

func BenchmarkConvert_NoRelationMembership(b *testing.B) {
	o := parseFile(b, benchGeoJSON)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := Convert(o, NoRelationMembership(true)); err != nil {
			b.Fatalf("convert error: %e", err)
		}
	}
}

func BenchmarkConvert_NoIDsMetaMembership(b *testing.B) {
	o := parseFile(b, benchGeoJSON)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := Convert(o, NoID(true), NoMeta(true), NoRelationMembership(true)); err != nil {
			b.Fatalf("convert error: %e", err)
		}
	}
}

func parseFile(t testing.TB, filename string) (o *osm.OSM) {
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("could not read file: %e", err)
	}

	if err = xml.Unmarshal(data, &o); err != nil {
		t.Fatalf("failed to unmarshal %s: %e", filename, err)
	}

	return
}

func annotate(o *osm.OSM) {
	nodes := make(map[osm.NodeID]*osm.Node)
	for _, n := range o.Nodes {
		nodes[n.ID] = n
	}

	for _, w := range o.Ways {
		for i, wn := range w.Nodes {
			if n := nodes[wn.ID]; n != nil {
				w.Nodes[i].Lat = n.Lat
				w.Nodes[i].Lon = n.Lon
				w.Nodes[i].Version = n.Version
			}
		}
	}
}
