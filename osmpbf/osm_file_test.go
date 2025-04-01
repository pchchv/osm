package osmpbf

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/pchchv/osm"
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
