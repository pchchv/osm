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

// **NOTE**, it is assumed that keys and vals have the
// same length and that the index is within the range of stringTable.
func scanTags(stringTable []string, keys, vals *pbr.Iterator) (osm.Tags, error) {
	var index int
	tags := make(osm.Tags, keys.Count(pbr.WireTypeVarint))
	for keys.HasNext() {
		k, err := keys.Uint32()
		if err != nil {
			return nil, err
		}

		v, err := vals.Uint32()
		if err != nil {
			return nil, err
		}

		tags[index] = osm.Tag{
			Key:   stringTable[k],
			Value: stringTable[v],
		}
		index++
	}

	return tags, nil
}

// Makes relation members from a stringtable and three parallel arrays of IDs.
func extractMembers(st []string, roles *pbr.Iterator, memids *pbr.Iterator, types *pbr.Iterator) (osm.Members, error) {
	var index, memID int64
	members := make(osm.Members, types.Count(pbr.WireTypeVarint))
	for roles.HasNext() {
		r, err := roles.Int32()
		if err != nil {
			return nil, err
		}

		members[index].Role = st[r]
		m, err := memids.Sint64()
		if err != nil {
			return nil, err
		}

		memID += m
		members[index].Ref = memID
		t, err := types.Int32()
		if err != nil {
			return nil, err
		}

		switch osmpbf.Relation_MemberType(t) {
		case osmpbf.Relation_NODE:
			members[index].Type = osm.TypeNode
		case osmpbf.Relation_WAY:
			members[index].Type = osm.TypeWay
		case osmpbf.Relation_RELATION:
			members[index].Type = osm.TypeRelation
		}

		index++
	}

	return members, nil
}
