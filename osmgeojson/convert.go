package osmgeojson

import (
	"fmt"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/geojson"
	"github.com/pchchv/osm"
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
