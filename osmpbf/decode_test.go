package osmpbf

import (
	"time"

	"github.com/pchchv/osm"
)

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
