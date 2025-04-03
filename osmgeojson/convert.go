package osmgeojson

import (
	"fmt"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/geojson"
	"github.com/pchchv/osm"
	"github.com/pchchv/osm/mputil"
)

type relationSummary struct {
	ID   osm.RelationID    `json:"id"`
	Role string            `json:"role"`
	Tags map[string]string `json:"tags"`
}

type context struct {
	noID                   bool
	noMeta                 bool
	noRelationMembership   bool
	includeInvalidPolygons bool
	osm                    *osm.OSM
	wayMap                 map[osm.WayID]*osm.Way
	skippable              map[osm.WayID]struct{}
	wayMember              map[osm.NodeID]struct{}
	nodeMap                map[osm.NodeID]*osm.Node
	relationMember         map[osm.FeatureID][]*relationSummary
}

func (ctx *context) addMetaProperties(props geojson.Properties, e osm.Element) {
	if !ctx.noRelationMembership {
		relations := ctx.relationMember[e.FeatureID()]
		if len(relations) != 0 {
			props["relations"] = relations
		} else {
			props["relations"] = []*relationSummary{}
		}
	}

	if ctx.noMeta {
		return
	}

	meta := make(map[string]interface{}, 5)
	switch e := e.(type) {
	case *osm.Node:
		if !e.Timestamp.IsZero() {
			meta["timestamp"] = e.Timestamp
		}

		if e.Version != 0 {
			meta["version"] = e.Version
		}

		if e.ChangesetID != 0 {
			meta["changeset"] = e.ChangesetID
		}

		if e.User != "" {
			meta["user"] = e.User
		}

		if e.UserID != 0 {
			meta["uid"] = e.UserID
		}
	case *osm.Way:
		if !e.Timestamp.IsZero() {
			meta["timestamp"] = e.Timestamp
		}

		if e.Version != 0 {
			meta["version"] = e.Version
		}

		if e.ChangesetID != 0 {
			meta["changeset"] = e.ChangesetID
		}

		if e.User != "" {
			meta["user"] = e.User
		}

		if e.UserID != 0 {
			meta["uid"] = e.UserID
		}
	case *osm.Relation:
		if !e.Timestamp.IsZero() {
			meta["timestamp"] = e.Timestamp
		}

		if e.Version != 0 {
			meta["version"] = e.Version
		}

		if e.ChangesetID != 0 {
			meta["changeset"] = e.ChangesetID
		}

		if e.User != "" {
			meta["user"] = e.User
		}

		if e.UserID != 0 {
			meta["uid"] = e.UserID
		}
	default:
		panic("unsupported type")
	}

	props["meta"] = meta
}

// getNode finds a node in the set.
// This allows to lazily create a
// node map only if nodes+ways are not augmented
// (ie. include the lat/lon on them).
func (ctx *context) getNode(id osm.NodeID) *osm.Node {
	if ctx.nodeMap == nil {
		ctx.nodeMap = make(map[osm.NodeID]*osm.Node, len(ctx.osm.Nodes))
		for _, n := range ctx.osm.Nodes {
			ctx.nodeMap[n.ID] = n
		}
	}

	return ctx.nodeMap[id]
}

func (ctx *context) nodeToFeature(n *osm.Node) *geojson.Feature {
	if n.Lon == 0 && n.Lat == 0 && n.Version == 0 {
		return nil
	}

	f := geojson.NewFeature(geo.Point{n.Lon, n.Lat})
	if !ctx.noID {
		f.ID = fmt.Sprintf("node/%d", n.ID)
	}

	f.Properties["id"] = int(n.ID)
	f.Properties["type"] = "node"
	f.Properties["tags"] = n.Tags.Map()
	ctx.addMetaProperties(f.Properties, n)

	return f
}

func (ctx *context) wayToLineString(w *osm.Way) (geo.LineString, bool) {
	var tainted bool
	ls := make(geo.LineString, 0, len(w.Nodes))
	for _, wn := range w.Nodes {
		if wn.Lon != 0 || wn.Lat != 0 {
			ls = append(ls, geo.Point{wn.Lon, wn.Lat})
		} else if n := ctx.getNode(wn.ID); n != nil {
			ls = append(ls, geo.Point{n.Lon, n.Lat})
		} else {
			tainted = true
		}
	}

	return ls, tainted
}

func (ctx *context) buildRouteLineString(relation *osm.Relation) *geojson.Feature {
	var tainted bool
	lines := make([]mputil.Segment, 0, 10)
	for _, m := range relation.Members {
		if m.Type != osm.TypeWay {
			continue
		}

		way := ctx.wayMap[osm.WayID(m.Ref)]
		if way == nil {
			tainted = true
			continue
		}

		if !hasInterestingTags(way.Tags, nil) {
			ctx.skippable[way.ID] = struct{}{}
		}

		ls, t := ctx.wayToLineString(way)
		if t {
			tainted = true
		}

		if len(ls) == 0 {
			continue
		}

		lines = append(lines, mputil.Segment{
			Orientation: m.Orientation,
			Line:        ls,
		})
	}

	if len(lines) == 0 {
		return nil
	}

	var geometry geo.Geometry
	lineSections := mputil.Join(lines)
	if len(lineSections) == 1 {
		geometry = lineSections[0].LineString()
	} else {
		mls := make(geo.MultiLineString, 0, len(lines))
		for _, ls := range lineSections {
			mls = append(mls, ls.LineString())
		}
		geometry = mls
	}

	f := geojson.NewFeature(geometry)
	if !ctx.noID {
		f.ID = fmt.Sprintf("relation/%d", relation.ID)
	}

	f.Properties["id"] = int(relation.ID)
	f.Properties["type"] = "relation"
	if tainted {
		f.Properties["tainted"] = true
	}

	f.Properties["tags"] = relation.Tags.Map()
	ctx.addMetaProperties(f.Properties, relation)

	return f
}

func toRing(ls geo.LineString) geo.Ring {
	if len(ls) < 2 {
		return geo.Ring(ls)
	}

	// duplicate last point
	if ls[0] != ls[len(ls)-1] {
		return geo.Ring(append(ls, ls[0]))
	}

	return geo.Ring(ls)
}

func hasInterestingTags(tags osm.Tags, ignore map[string]string) bool {
	if len(tags) == 0 {
		return false
	}

	for _, tag := range tags {
		k, v := tag.Key, tag.Value
		if !osm.UninterestingTags[k] &&
			(ignore == nil || !(ignore[k] == "true" || ignore[k] == v)) {
			return true
		}
	}

	return false
}
