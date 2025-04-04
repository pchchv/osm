package osmapi

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pchchv/osm"
)

// Node returns the latest version of the node from the osm rest api.
func (ds *Datasource) Node(ctx context.Context, id osm.NodeID, opts ...FeatureOption) (*osm.Node, error) {
	params, err := featureOptions(opts)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/node/%d?%s", ds.baseURL(), id, params)
	o := &osm.OSM{}
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	if l := len(o.Nodes); l != 1 {
		return nil, fmt.Errorf("wrong number of nodes, expected 1, got %v", l)
	}

	return o.Nodes[0], nil
}

// Nodes returns the latest version of the nodes from the osm rest api.
// Returns 404 if any node is missing.
func (ds *Datasource) Nodes(ctx context.Context, ids []osm.NodeID, opts ...FeatureOption) (osm.Nodes, error) {
	params, err := featureOptions(opts)
	if err != nil {
		return nil, err
	}

	data := make([]byte, 0, 11*len(ids))
	for i, id := range ids {
		if i != 0 {
			data = append(data, byte(','))
		}
		data = strconv.AppendInt(data, int64(id), 10)
	}

	url := ds.baseURL() + "/nodes?nodes=" + string(data)
	if len(params) > 0 {
		url += "&" + params
	}

	o := &osm.OSM{}
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	return o.Nodes, nil
}
