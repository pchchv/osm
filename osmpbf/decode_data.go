package osmpbf

import (
	"errors"
	"time"

	"github.com/pchchv/osm"
	"github.com/pchchv/osm/osmpbf/internal/osmpbf"
	"github.com/pchchv/pbr"
	"google.golang.org/protobuf/proto"
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

func (dec *dataDecoder) Decode(blob *osmpbf.Blob) ([]osm.Object, error) {
	var err error
	dec.q = make([]osm.Object, 0, 8000) // typical PrimitiveBlock contains 8k OSM entities
	if dec.data, err = getData(blob, dec.data); err != nil {
		return nil, err
	}

	if err = dec.scanPrimitiveBlock(dec.data); err != nil {
		return nil, err
	}

	return dec.q, nil
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

func (dec *dataDecoder) scanRelations(data []byte, relation *osm.Relation) (*osm.Relation, error) {
	st := dec.primitiveBlock.GetStringtable().GetS()
	dateGranularity := int64(dec.primitiveBlock.GetDateGranularity())
	msg := pbr.New(data)
	if relation == nil {
		relation = &osm.Relation{Visible: true}
	}

	var foundKeys, foundVals, foundRoles, foundMemids, foundTypes bool
	for msg.Next() {
		var i64 int64
		var err error
		switch msg.FieldNumber() {
		case 1:
			i64, err = msg.Int64()
			relation.ID = osm.RelationID(i64)
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
					relation.Version = int(v)
				case 2:
					v, err := info.Int64()
					if err != nil {
						return nil, err
					}
					millisec := time.Duration(v*dateGranularity) * time.Millisecond
					relation.Timestamp = time.Unix(0, millisec.Nanoseconds()).UTC()
				case 3:
					v, err := info.Int64()
					if err != nil {
						return nil, err
					}
					relation.ChangesetID = osm.ChangesetID(v)
				case 4:
					v, err := info.Int32()
					if err != nil {
						return nil, err
					}
					relation.UserID = osm.UserID(v)
				case 5:
					v, err := info.Uint32()
					if err != nil {
						return nil, err
					}
					relation.User = st[v]
				case 6:
					v, err := info.Bool()
					if err != nil {
						return nil, err
					}
					relation.Visible = v
				default:
					info.Skip()
				}
			}

			if info.Error() != nil {
				return nil, info.Error()
			}
		case 8: // refs or nodes
			dec.roles, err = msg.Iterator(dec.roles)
			foundRoles = true
		case 9:
			dec.memids, err = msg.Iterator(dec.memids)
			foundMemids = true
		case 10:
			dec.types, err = msg.Iterator(dec.types)
			foundTypes = true
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

	var err error
	// possible for relation to not have tags
	if foundKeys && foundVals {
		relation.Tags, err = scanTags(st, dec.keys, dec.vals)
		if err != nil {
			return nil, err
		}
	}

	// possible for relation to not have any members
	if foundRoles && foundMemids && foundTypes {
		relation.Members, err = extractMembers(st, dec.roles, dec.memids, dec.types)
		if err != nil {
			return nil, err
		}
	}

	return relation, nil
}

func (dec *dataDecoder) scanDenseNodes(data []byte) (err error) {
	var foundIds, foundInfo, foundLats, foundLons, foundKeyVals bool
	msg := pbr.New(data)
	for msg.Next() {
		switch msg.FieldNumber() {
		case 1: // ids
			dec.ids, err = msg.Iterator(dec.ids)
			foundIds = true
		case 5: // dense info
			d, err := msg.MessageData()
			if err != nil {
				return err
			}

			// verify that all fields are “found” because the
			// object from the previous block is reused and
			// cannot simply be checked for nil
			var foundVersions, foundTimestamps, foundChangesets, foundUids, foundUsids, foundVisibles bool
			info := pbr.New(d)
			for info.Next() {
				switch info.FieldNumber() {
				case 1: // version
					dec.versions, err = info.Iterator(dec.versions)
					foundVersions = true
				case 2: // timestamp
					dec.timestamps, err = info.Iterator(dec.timestamps)
					foundTimestamps = true
				case 3: // changeset
					dec.changesets, err = info.Iterator(dec.changesets)
					foundChangesets = true
				case 4: // uid
					dec.uids, err = info.Iterator(dec.uids)
					foundUids = true
				case 5: // user_sid
					dec.usids, err = info.Iterator(dec.usids)
					foundUsids = true
				case 6: // visible, optional, default true
					dec.visibles, err = info.Iterator(dec.visibles)
					foundVisibles = true
				default:
					info.Skip()
				}

				if err != nil {
					return err
				}
			}

			if info.Error() != nil {
				return info.Error()
			}

			if !foundVersions {
				dec.versions = nil
			}

			if !foundTimestamps {
				dec.timestamps = nil
			}

			if !foundChangesets {
				dec.changesets = nil
			}

			if !foundUids {
				dec.uids = nil
			}

			if !foundUsids {
				dec.usids = nil
			}

			// visibles are optional, default is true
			if !foundVisibles {
				dec.visibles = nil
			}

			foundInfo = true
		case 8: // lat
			dec.lats, err = msg.Iterator(dec.lats)
			foundLats = true
		case 9: // lon
			dec.lons, err = msg.Iterator(dec.lons)
			foundLons = true
		case 10: // keys_vals
			dec.keyvals, err = msg.Iterator(dec.keyvals)
			foundKeyVals = true
		default:
			msg.Skip()
		}

		if err != nil {
			return err
		}
	}

	if msg.Error() != nil {
		return msg.Error()
	}

	if !foundIds {
		return errors.New("osmpbf: dense node did not contain ids")
	}

	if !foundLats {
		return errors.New("osmpbf: dense node did not contain latitudes")
	}

	if !foundLons {
		return errors.New("osmpbf: dense node did not contain longitudes")
	}

	// keyvals could be empty if all nodes are tagless
	if !foundKeyVals {
		dec.keyvals = nil
	}

	if !foundInfo {
		dec.versions = nil
		dec.timestamps = nil
		dec.changesets = nil
		dec.uids = nil
		dec.usids = nil
		dec.visibles = nil
	}

	return dec.extractDenseNodes()
}

func (dec *dataDecoder) scanPrimitiveGroup(data []byte) error {
	msg := pbr.New(data)
	way := &osm.Way{Visible: true}
	relation := &osm.Relation{Visible: true}
	for msg.Next() {
		fn := msg.FieldNumber()
		if fn == 1 {
			panic("nodes are not supported, currently untested")
		}

		if fn == 2 && !dec.scanner.SkipNodes {
			data, err := msg.MessageData()
			if err != nil {
				return err
			}

			if err = dec.scanDenseNodes(data); err != nil {
				return err
			}

			continue
		}

		if fn == 3 && !dec.scanner.SkipWays {
			data, err := msg.MessageData()
			if err != nil {
				return err
			}

			way, err = dec.scanWays(data, way)
			if err != nil {
				return err
			}

			if dec.scanner.FilterWay == nil || dec.scanner.FilterWay(way) {
				dec.q = append(dec.q, way)
				way = &osm.Way{Visible: true}
			} else {
				tags := way.Tags
				nodes := way.Nodes
				*way = osm.Way{Visible: true, Nodes: nodes[:0], Tags: tags[:0]}
			}

			continue
		}

		if fn == 4 && !dec.scanner.SkipRelations {
			data, err := msg.MessageData()
			if err != nil {
				return err
			}

			relation, err = dec.scanRelations(data, relation)
			if err != nil {
				return err
			}

			if dec.scanner.FilterRelation == nil || dec.scanner.FilterRelation(relation) {
				dec.q = append(dec.q, relation)
				relation = &osm.Relation{Visible: true}
			} else {
				tags := relation.Tags
				members := relation.Members
				*relation = osm.Relation{Visible: true, Members: members[:0], Tags: tags[:0]}
			}

			continue
		}

		msg.Skip()
	}

	return msg.Error()
}

func (dec *dataDecoder) scanPrimitiveBlock(data []byte) error {
	msg := pbr.New(data)
	if dec.primitiveBlock == nil {
		dec.primitiveBlock = &osmpbf.PrimitiveBlock{
			Stringtable: &osmpbf.StringTable{},
		}
	} else {
		dec.primitiveBlock.Stringtable.S = dec.primitiveBlock.Stringtable.S[:0]
		dec.primitiveBlock.Primitivegroup = dec.primitiveBlock.Primitivegroup[:0]
		dec.primitiveBlock.Granularity = nil
		dec.primitiveBlock.LatOffset = nil
		dec.primitiveBlock.LonOffset = nil
		dec.primitiveBlock.DateGranularity = nil
	}

	for msg.Next() {
		switch msg.FieldNumber() {
		case 1:
			d, err := msg.MessageData()
			if err != nil {
				return err
			}

			if err = proto.Unmarshal(d, dec.primitiveBlock.Stringtable); err != nil {
				return err
			}
		case 17:
			v, err := msg.Int32()
			dec.primitiveBlock.Granularity = &v
			if err != nil {
				return err
			}
		case 18:
			v, err := msg.Int32()
			dec.primitiveBlock.DateGranularity = &v
			if err != nil {
				return err
			}
		case 19:
			v, err := msg.Int64()
			dec.primitiveBlock.LatOffset = &v
			if err != nil {
				return err
			}
		case 20:
			v, err := msg.Int64()
			dec.primitiveBlock.LonOffset = &v
			if err != nil {
				return err
			}
		default:
			msg.Skip()
		}
	}

	if msg.Error() != nil {
		return msg.Error()
	}

	// is needed the offsets and granularities for the group decoding
	msg.Reset(nil)
	for msg.Next() {
		switch msg.FieldNumber() {
		case 2:
			d, err := msg.MessageData()
			if err != nil {
				return err
			}
			err = dec.scanPrimitiveGroup(d)
			if err != nil {
				return err
			}
		default:
			msg.Skip()
		}
	}

	return msg.Error()
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

// extractMembers makes relation members from stringtable and three parallel arrays of IDs.
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
