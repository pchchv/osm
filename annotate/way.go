package annotate

import "github.com/pchchv/osm"

// parentWay wraps a osm.Way into the
// core.Parent interface so that updates can be computed.
type parentWay struct {
	Way *osm.Way
}
