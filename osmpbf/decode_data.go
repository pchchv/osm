package osmpbf

import (
	"github.com/pchchv/osm"
	"github.com/pchchv/osm/osmpbf/internal/osmpbf"
	"github.com/pchchv/pbr"
)

// dataDecoder is a decoder for Blob with OSMData (PrimitiveBlock).
type dataDecoder struct {
	scanner        *Scanner
	data           []byte
	q              []osm.Object
	primitiveBlock *osmpbf.PrimitiveBlock // cache objects to save allocations
	keys, vals     *pbr.Iterator
	nodes          *pbr.Iterator // ways
	wlats          *pbr.Iterator
	wlons          *pbr.Iterator
	roles          *pbr.Iterator // relations
	memids         *pbr.Iterator
	types          *pbr.Iterator
	ids            *pbr.Iterator // dense nodes
	versions       *pbr.Iterator
	timestamps     *pbr.Iterator
	changesets     *pbr.Iterator
	uids           *pbr.Iterator
	usids          *pbr.Iterator
	visibles       *pbr.Iterator
	lats           *pbr.Iterator
	lons           *pbr.Iterator
	keyvals        *pbr.Iterator
}
