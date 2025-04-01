package osmpbf

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/pchchv/osm"
)

const (
	// Originally downloaded from https://download.geofabrik.de/north-america/us/delaware.html
	Delaware = "../testdata/delaware-latest.osm.pbf"
	// Originally downloaded from http://download.geofabrik.de/europe/great-britain/england/greater-london.html
	London    = "../testdata/greater-london-140324.osm.pbf"
	LondonURL = "https://github.com/pchchv/osm/raw/refs/heads/main/testdata/greater-london-140324.osm.pbf"
	// Created based on the above file, by running `osmium add-locations-to-ways`.
	LondonLocations    = "../testdata/greater-london-140324-low.osm.pbf"
	LondonLocationsURL = "https://github.com/pchchv/osm/raw/refs/heads/main/testdata/greater-london-140324-low.osm.pbf"
)

var (
	erc  uint64 = 12833
	ewc  uint64 = 459055
	encl uint64 = 244523
	enc  uint64 = 2729006
	ew          = stripCoordinates(ewl)
	en          = &osm.Node{
		ID:  18088578,
		Lat: 51.5442632,
		Lon: -0.2010027,
		Tags: osm.Tags([]osm.Tag{
			{Key: "alt_name", Value: "The King's Head"},
			{Key: "amenity", Value: "pub"},
			{Key: "created_by", Value: "JOSM"},
			{Key: "name", Value: "The Luminaire"},
			{Key: "note", Value: "Live music venue too"},
		}),
		Version:     2,
		Timestamp:   parseTime("2009-05-20T10:28:54Z"),
		ChangesetID: 1260468,
		UserID:      508,
		User:        "Welshie",
		Visible:     true,
	}
	ewl = &osm.Way{
		ID: 4257116,
		Nodes: osm.WayNodes{
			{ID: 21544864, Lat: 51.5230531, Lon: -0.1408525},
			{ID: 333731851, Lat: 51.5224309, Lon: -0.1402297},
			{ID: 333731852, Lat: 51.5224107, Lon: -0.1401878},
			{ID: 333731850, Lat: 51.522422, Lon: -0.1401375},
			{ID: 333731855, Lat: 51.522792, Lon: -0.1392477},
			{ID: 333731858, Lat: 51.5228209, Lon: -0.1392124},
			{ID: 333731854, Lat: 51.5228579, Lon: -0.1392339},
			{ID: 108047, Lat: 51.5234407, Lon: -0.1398771},
			{ID: 769984352, Lat: 51.5232469, Lon: -0.1403648},
			{ID: 21544864, Lat: 51.5230531, Lon: -0.1408525},
		},
		Tags: osm.Tags([]osm.Tag{
			{Key: "area", Value: "yes"},
			{Key: "highway", Value: "pedestrian"},
			{Key: "name", Value: "Fitzroy Square"},
		}),
		Version:     7,
		Timestamp:   parseTime("2013-08-07T12:08:39Z"),
		ChangesetID: 17253164,
		UserID:      1016290,
		User:        "Amaroussi",
		Visible:     true,
	}
	er = &osm.Relation{
		ID: 7677,
		Members: osm.Members{
			{Ref: 4875932, Type: osm.TypeWay, Role: "outer"},
			{Ref: 4894305, Type: osm.TypeWay, Role: "inner"},
		},
		Tags: osm.Tags([]osm.Tag{
			{Key: "created_by", Value: "Potlatch 0.9c"},
			{Key: "type", Value: "multipolygon"},
		}),
		Version:     4,
		Timestamp:   parseTime("2008-07-19T15:04:03Z"),
		ChangesetID: 540201,
		UserID:      3876,
		User:        "Edgemaster",
		Visible:     true,
	}
)

func TestDecode(t *testing.T) {
	ft := &OSMFileTest{
		T:            t,
		FileName:     London,
		FileURL:      LondonURL,
		ExpNode:      en,
		ExpWay:       ew,
		ExpRel:       er,
		ExpNodeCount: enc,
		ExpWayCount:  ewc,
		ExpRelCount:  erc,
		IDsExpOrder:  IDsExpectedOrder,
	}
	ft.testDecode()
}

func TestDecodeLocations(t *testing.T) {
	ft := &OSMFileTest{
		T:            t,
		FileName:     LondonLocations,
		FileURL:      LondonLocationsURL,
		ExpNode:      en,
		ExpWay:       ewl,
		ExpRel:       er,
		ExpNodeCount: encl,
		ExpWayCount:  ewc,
		ExpRelCount:  erc,
		IDsExpOrder:  IDsExpectedOrderNoNodes,
	}
	ft.testDecode()
}

func TestDecode_Close(t *testing.T) {
	f, err := os.Open(Delaware)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// should close at start
	f.Seek(0, 0)
	d := newDecoder(context.Background(), &Scanner{}, f)
	d.Start(5)

	if err = d.Close(); err != nil {
		t.Errorf("close error: %e", err)
	}

	// should close after partial read
	f.Seek(0, 0)
	d = newDecoder(context.Background(), &Scanner{}, f)
	d.Start(5)
	d.Next()
	d.Next()
	if err = d.Close(); err != nil {
		t.Errorf("close error: %e", err)
	}

	// should close after full read
	f.Seek(0, 0)
	d = newDecoder(context.Background(), &Scanner{}, f)
	d.Start(5)
	var elements int
	for {
		if _, err := d.Next(); err == io.EOF {
			break
		} else if err != nil {
			t.Errorf("next error: %e", err)
		}

		elements++
	}

	if elements < 2 {
		t.Errorf("did not read enough elements: %v", elements)
	}

	// should close at end of read
	if err = d.Close(); err != nil {
		t.Errorf("close error: %e", err)
	}
}

func parseTime(s string) time.Time {
	if t, err := time.Parse(time.RFC3339, s); err != nil {
		panic(err)
	} else {
		return t
	}
}

func stripCoordinates(w *osm.Way) *osm.Way {
	if w == nil {
		return nil
	}

	ws := new(osm.Way)
	*ws = *w
	ws.Nodes = make(osm.WayNodes, len(w.Nodes))
	for i, n := range w.Nodes {
		n.Lat, n.Lon = 0, 0
		ws.Nodes[i] = n
	}

	return ws
}
