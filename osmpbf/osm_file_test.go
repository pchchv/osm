package osmpbf

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"testing"

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

func (ft *OSMFileTest) testDecode() {
	ft.downloadTestOSMFile()
	f, err := os.Open(ft.FileName)
	if err != nil {
		ft.Fatal(err)
	}
	defer f.Close()

	d := newDecoder(context.Background(), &Scanner{}, f)
	if err = d.Start(runtime.GOMAXPROCS(-1)); err != nil {
		ft.Fatal(err)
	}

	var n *osm.Node
	var w *osm.Way
	var r *osm.Relation
	var nc, wc, rc uint64
	idsOrder := make([]string, 0, len(IDsExpectedOrder))
	for {
		e, err := d.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			ft.Fatal(err)
		}

		switch v := e.(type) {
		case *osm.Node:
			nc++
			if v.ID == ft.ExpNode.ID {
				n = v
			}
			id := fmt.Sprintf("node/%d", v.ID)
			if _, ok := IDs[id]; ok {
				idsOrder = append(idsOrder, id)
			}
		case *osm.Way:
			wc++
			if v.ID == ft.ExpWay.ID {
				w = v
			}
			id := fmt.Sprintf("way/%d", v.ID)
			if _, ok := IDs[id]; ok {
				idsOrder = append(idsOrder, id)
			}
		case *osm.Relation:
			rc++
			if v.ID == ft.ExpRel.ID {
				r = v
			}
			id := fmt.Sprintf("relation/%d", v.ID)
			if _, ok := IDs[id]; ok {
				idsOrder = append(idsOrder, id)
			}
		}
	}
	d.Close()

	if !reflect.DeepEqual(ft.ExpNode, n) {
		ft.Errorf("\nExpected: %#v\nActual:   %#v", ft.ExpNode, n)
	}

	roundCoordinates(w)
	if !reflect.DeepEqual(ft.ExpWay, w) {
		ft.Errorf("\nExpected: %#v\nActual:   %#v", ft.ExpWay, w)
	}

	if !reflect.DeepEqual(ft.ExpRel, r) {
		ft.Errorf("\nExpected: %#v\nActual:   %#v", ft.ExpRel, r)
	}

	if ft.ExpNodeCount != nc || ft.ExpWayCount != wc || ft.ExpRelCount != rc {
		ft.Errorf("\nExpected %7d nodes, %7d ways, %7d relations\nGot %7d nodes, %7d ways, %7d relations.",
			ft.ExpNodeCount, ft.ExpWayCount, ft.ExpRelCount, nc, wc, rc)
	}

	if !reflect.DeepEqual(ft.IDsExpOrder, idsOrder) {
		ft.Errorf("\nExpected: %v\nGot:      %v", ft.IDsExpOrder, idsOrder)
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

func init() {
	IDs = make(map[string]bool)
	for _, id := range IDsExpectedOrder {
		IDs[id] = false
	}
}
