package osmpbf

import (
	"time"

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

func (dec *dataDecoder) extractDenseNodes() error {
	var uid, usid int32
	var id, lat, lon, timestamp, changeset int64
	st := dec.primitiveBlock.GetStringtable().GetS()
	granularity := int64(dec.primitiveBlock.GetGranularity())
	dateGranularity := int64(dec.primitiveBlock.GetDateGranularity())
	latOffset := dec.primitiveBlock.GetLatOffset()
	lonOffset := dec.primitiveBlock.GetLonOffset()
	// NOTE: do not try pre-allocating an array of nodes because
	// saving just one will stop the GC from cleaning up the
	// whole pre-allocated array
	n := &osm.Node{Visible: true}
	for dec.ids.HasNext() {
		// ID
		v1, err := dec.ids.Sint64()
		if err != nil {
			return err
		}

		id += v1
		n.ID = osm.NodeID(id)
		// version
		if dec.versions != nil {
			v2, err := dec.versions.Int32()
			if err != nil {
				return err
			}
			n.Version = int(v2)
		}

		// timestamp
		if dec.timestamps != nil {
			v3, err := dec.timestamps.Sint64()
			if err != nil {
				return err
			}

			timestamp += v3
			millisec := time.Duration(timestamp*dateGranularity) * time.Millisecond
			n.Timestamp = time.Unix(0, millisec.Nanoseconds()).UTC()
		}

		// changeset
		if dec.changesets != nil {
			v4, err := dec.changesets.Sint64()
			if err != nil {
				return err
			}

			changeset += v4
			n.ChangesetID = osm.ChangesetID(changeset)
		}

		// uid
		if dec.uids != nil {
			v5, err := dec.uids.Sint32()
			if err != nil {
				return err
			}

			uid += v5
			n.UserID = osm.UserID(uid)
		}

		// usid
		if dec.usids != nil {
			v6, err := dec.usids.Sint32()
			if err != nil {
				return err
			}

			usid += v6
			n.User = st[usid]
		}

		// visible
		if dec.visibles != nil {
			v7, err := dec.visibles.Bool()
			if err != nil {
				return err
			}

			n.Visible = v7
		}

		// lat
		v8, err := dec.lats.Sint64()
		if err != nil {
			return err
		}

		lat += v8
		n.Lat = 1e-9 * float64(latOffset+(granularity*lat))

		// lon
		v9, err := dec.lons.Sint64()
		if err != nil {
			return err
		}

		lon += v9
		n.Lon = 1e-9 * float64(lonOffset+(granularity*lon))

		// tags, could be missing if all nodes are tagless
		if dec.keyvals != nil {
			var count int
			for i := dec.keyvals.Index; i < len(dec.keyvals.Data); i++ {
				b := dec.keyvals.Data[i]
				if b == 0 {
					break
				}

				if b < 128 {
					count++
				}
			}

			if cap(n.Tags) < count/2 {
				n.Tags = make(osm.Tags, 0, count/2)
			}

			for {
				if k, err := dec.keyvals.Int32(); err != nil {
					return err
				} else if k == 0 {
					break
				} else {
					v, err := dec.keyvals.Int32()
					if err != nil {
						return err
					}

					n.Tags = append(n.Tags, osm.Tag{Key: st[k], Value: st[v]})
				}
			}
		}

		if dec.scanner.FilterNode == nil || dec.scanner.FilterNode(n) {
			dec.q = append(dec.q, n)
			n = &osm.Node{Visible: true}
		} else {
			// skip unwanted nodes
			*n = osm.Node{Visible: true, Tags: n.Tags[:0]}
		}
	}

	return nil
}

func (dec *dataDecoder) scanWays(data []byte, way *osm.Way) (*osm.Way, error) {
	st := dec.primitiveBlock.GetStringtable().GetS()
	granularity := int64(dec.primitiveBlock.GetGranularity())
	dateGranularity := int64(dec.primitiveBlock.GetDateGranularity())
	latOffset := dec.primitiveBlock.GetLatOffset()
	lonOffset := dec.primitiveBlock.GetLonOffset()
	msg := pbr.New(data)
	if way == nil {
		way = &osm.Way{Visible: true}
	}

	var foundKeys, foundVals bool
	for msg.Next() {
		var i64 int64
		var err error
		switch msg.FieldNumber() {
		case 1:
			i64, err = msg.Int64()
			way.ID = osm.WayID(i64)
		case 2:
			dec.keys, err = msg.Iterator(dec.keys)
			foundKeys = true
		case 3:
			dec.vals, err = msg.Iterator(dec.vals)
			foundVals = true
		case 4: // info
			d, err := msg.MessageData()
			if err != nil {
				return nil, err
			}

			info := pbr.New(d)
			for info.Next() {
				switch info.FieldNumber() {
				case 1:
					v, err := info.Int32()
					if err != nil {
						return nil, err
					}
					way.Version = int(v)
				case 2:
					v, err := info.Int64()
					if err != nil {
						return nil, err
					}
					millisec := time.Duration(v*dateGranularity) * time.Millisecond
					way.Timestamp = time.Unix(0, millisec.Nanoseconds()).UTC()
				case 3:
					v, err := info.Int64()
					if err != nil {
						return nil, err
					}
					way.ChangesetID = osm.ChangesetID(v)
				case 4:
					v, err := info.Int32()
					if err != nil {
						return nil, err
					}
					way.UserID = osm.UserID(v)
				case 5:
					v, err := info.Uint32()
					if err != nil {
						return nil, err
					}
					way.User = st[v]
				case 6:
					v, err := info.Bool()
					if err != nil {
						return nil, err
					}
					way.Visible = v
				default:
					info.Skip()
				}
			}

			if info.Error() != nil {
				return nil, info.Error()
			}
		case 8: // refs or nodes
			dec.nodes, err = msg.Iterator(dec.nodes)
			if err != nil {
				return nil, err
			}

			var prev, index int64
			if len(way.Nodes) == 0 {
				way.Nodes = make(osm.WayNodes, dec.nodes.Count(pbr.WireTypeVarint))
			}

			for dec.nodes.HasNext() {
				v, err := dec.nodes.Sint64()
				if err != nil {
					return nil, err
				}

				prev = v + prev // delta encoding
				way.Nodes[index].ID = osm.NodeID(prev)
				index++
			}
		case 9: // lat
			dec.wlats, err = msg.Iterator(dec.wlats)
			if err != nil {
				return nil, err
			}

			var prev, index int64
			if len(way.Nodes) == 0 {
				way.Nodes = make(osm.WayNodes, dec.wlats.Count(pbr.WireTypeVarint))
			}

			for dec.wlats.HasNext() {
				v, err := dec.wlats.Sint64()
				if err != nil {
					return nil, err
				}

				prev = v + prev // delta encoding
				way.Nodes[index].Lat = 1e-9 * float64(latOffset+(granularity*prev))
				index++
			}
		case 10: // lon
			dec.wlons, err = msg.Iterator(dec.wlons)
			if err != nil {
				return nil, err
			}

			var prev, index int64
			if len(way.Nodes) == 0 {
				way.Nodes = make(osm.WayNodes, dec.wlons.Count(pbr.WireTypeVarint))
			}

			for dec.wlons.HasNext() {
				v, err := dec.wlons.Sint64()
				if err != nil {
					return nil, err
				}

				prev = v + prev // delta encoding
				way.Nodes[index].Lon = 1e-9 * float64(lonOffset+(granularity*prev))
				index++
			}
		default:
			msg.Skip()
		}

		if err != nil {
			return nil, err
		}
	}

	if msg.Error() != nil {
		return nil, msg.Error()
	}

	if foundKeys && foundVals {
		var err error
		way.Tags, err = scanTags(st, dec.keys, dec.vals)
		if err != nil {
			return nil, err
		}
	}

	return way, nil
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
