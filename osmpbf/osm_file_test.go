package osmpbf

import (
	"io"
	"math"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pchchv/osm"
)

const coordinatesPrecision = 1e7

var (
	IDs                     map[string]bool
	IDsExpectedOrder        = append(idsExpectedOrderNodes, IDsExpectedOrderNoNodes...)
	IDsExpectedOrderNoNodes = append(idsExpectedOrderWays, idsExpectedOrderRelations...)
	idsExpectedOrderNodes   = []string{
		"node/44", "node/47", "node/52", "node/58", "node/60",
		"node/79", // because way/79 is already there
		"node/2740703694", "node/2740703695", "node/2740703697",
		"node/2740703699", "node/2740703701",
	}
	idsExpectedOrderWays = []string{
		"way/73", "way/74", "way/75", "way/79", "way/482",
		"way/268745428", "way/268745431", "way/268745434", "way/268745436",
		"way/268745439",
	}
	idsExpectedOrderRelations = []string{
		"relation/69", "relation/94", "relation/152", "relation/245",
		"relation/332", "relation/3593436", "relation/3595575",
		"relation/3595798", "relation/3599126", "relation/3599127",
	}
)

type OSMFileTest struct {
	*testing.T
	FileName                               string
	FileURL                                string
	ExpNode                                *osm.Node
	ExpWay                                 *osm.Way
	ExpRel                                 *osm.Relation
	ExpNodeCount, ExpWayCount, ExpRelCount uint64
	IDsExpOrder                            []string
}

func (ft *OSMFileTest) downloadTestOSMFile() {
	if _, err := os.Stat(ft.FileName); os.IsNotExist(err) {
		out, err := os.Create(ft.FileName)
		if err != nil {
			ft.Fatal(err)
		}
		defer out.Close()

		resp, err := http.Get(ft.FileURL)
		if err != nil {
			ft.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			ft.Fatalf("test status code invalid: %v", resp.StatusCode)
		}

		if _, err := io.Copy(out, resp.Body); err != nil {
			ft.Fatal(err)
		}
	} else if err != nil {
		ft.Fatal(err)
	}
}

func roundCoordinates(w *osm.Way) {
	if w != nil {
		for i := range w.Nodes {
			w.Nodes[i].Lat = math.Round(w.Nodes[i].Lat*coordinatesPrecision) / coordinatesPrecision
			w.Nodes[i].Lon = math.Round(w.Nodes[i].Lon*coordinatesPrecision) / coordinatesPrecision
		}
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
