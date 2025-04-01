package osmpbf

import (
	"io"
	"math"
	"net/http"
	"os"
	"testing"

	"github.com/pchchv/osm"
)

const coordinatesPrecision = 1e7

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
