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

	o := &osm.OSM{}
	url := fmt.Sprintf("%s/node/%d?%s", ds.baseURL(), id, params)
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

// NodeRelations returns all relations a node is part of.
// There is no error if the element does not exist.
func (ds *Datasource) NodeRelations(ctx context.Context, id osm.NodeID, opts ...FeatureOption) (osm.Relations, error) {
	params, err := featureOptions(opts)
	if err != nil {
		return nil, err
	}

	o := &osm.OSM{}
	url := fmt.Sprintf("%s/node/%d/relations?%s", ds.baseURL(), id, params)
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	return o.Relations, nil
}

// NodeWays returns all ways a node is part of.
// There is no error if the element does not exist.
func (ds *Datasource) NodeWays(ctx context.Context, id osm.NodeID, opts ...FeatureOption) (osm.Ways, error) {
	params, err := featureOptions(opts)
	if err != nil {
		return nil, err
	}

	o := &osm.OSM{}
	url := fmt.Sprintf("%s/node/%d/ways?%s", ds.baseURL(), id, params)
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	return o.Ways, nil
}

// NodeVersion returns the specific version of the node from the osm rest api.
func (ds *Datasource) NodeVersion(ctx context.Context, id osm.NodeID, v int) (*osm.Node, error) {
	o := &osm.OSM{}
	url := fmt.Sprintf("%s/node/%d/%d", ds.baseURL(), id, v)
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	if l := len(o.Nodes); l != 1 {
		return nil, fmt.Errorf("wrong number of nodes, expected 1, got %v", l)
	}

	return o.Nodes[0], nil
}

// NodeHistory returns all the versions of the node from the osm rest api.
func (ds *Datasource) NodeHistory(ctx context.Context, id osm.NodeID) (osm.Nodes, error) {
	o := &osm.OSM{}
	url := fmt.Sprintf("%s/node/%d/history", ds.baseURL(), id)
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	return o.Nodes, nil
}

// Node returns the latest version of the node from the osm rest api.
// Delegates to the DefaultDatasource and uses its http.Client to make the request.
func Node(ctx context.Context, id osm.NodeID, opts ...FeatureOption) (*osm.Node, error) {
	return DefaultDatasource.Node(ctx, id, opts...)
}

// Nodes returns the latest version of the nodes from the osm rest api.
// Delegates to the DefaultDatasource and uses its http.Client to make the request.
func Nodes(ctx context.Context, ids []osm.NodeID, opts ...FeatureOption) (osm.Nodes, error) {
	return DefaultDatasource.Nodes(ctx, ids, opts...)
}

// NodeRelations returns all relations a node is part of.
// There is no error if the element does not exist.
// Delegates to the DefaultDatasource and uses its http.Client to make the request.
func NodeRelations(ctx context.Context, id osm.NodeID, opts ...FeatureOption) (osm.Relations, error) {
	return DefaultDatasource.NodeRelations(ctx, id, opts...)
}
