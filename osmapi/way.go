package osmapi

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pchchv/osm"
)

// Way returns the latest version of the way from the osm rest api.
func (ds *Datasource) Way(ctx context.Context, id osm.WayID, opts ...FeatureOption) (*osm.Way, error) {
	params, err := featureOptions(opts)
	if err != nil {
		return nil, err
	}

	o := &osm.OSM{}
	url := fmt.Sprintf("%s/way/%d?%s", ds.baseURL(), id, params)
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	if l := len(o.Ways); l != 1 {
		return nil, fmt.Errorf("wrong number of ways, expected 1, got %v", l)
	}

	return o.Ways[0], nil
}

// Ways returns the latest version of the ways from the osm rest api.
// Returns 404 if any way is missing.
func (ds *Datasource) Ways(ctx context.Context, ids []osm.WayID, opts ...FeatureOption) (osm.Ways, error) {
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

	url := ds.baseURL() + "/ways?ways=" + string(data)
	if len(params) > 0 {
		url += "&" + params
	}

	o := &osm.OSM{}
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	return o.Ways, nil
}

// WayRelations returns all relations a way is part of.
// There is no error if the element does not exist.
func (ds *Datasource) WayRelations(ctx context.Context, id osm.WayID, opts ...FeatureOption) (osm.Relations, error) {
	params, err := featureOptions(opts)
	if err != nil {
		return nil, err
	}

	o := &osm.OSM{}
	url := fmt.Sprintf("%s/way/%d/relations?%s", ds.baseURL(), id, params)
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	return o.Relations, nil
}
