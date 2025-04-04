package osmapi

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pchchv/osm"
)

// Relation returns the latest version of the relation from the osm rest api.
func (ds *Datasource) Relation(ctx context.Context, id osm.RelationID, opts ...FeatureOption) (*osm.Relation, error) {
	params, err := featureOptions(opts)
	if err != nil {
		return nil, err
	}

	o := &osm.OSM{}
	url := fmt.Sprintf("%s/relation/%d?%s", ds.baseURL(), id, params)
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	if l := len(o.Relations); l != 1 {
		return nil, fmt.Errorf("wrong number of relations, expected 1, got %v", l)
	}

	return o.Relations[0], nil
}

// Relations returns the latest version of the relations from the osm rest api.
// Returns 404 if any node is missing.
func (ds *Datasource) Relations(ctx context.Context, ids []osm.RelationID, opts ...FeatureOption) (osm.Relations, error) {
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

	url := ds.baseURL() + "/relations?relations=" + string(data)
	if len(params) > 0 {
		url += "&" + params
	}

	o := &osm.OSM{}
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	return o.Relations, nil
}

// RelationRelations returns all relations a relation is part of.
// There is no error if the element does not exist.
func (ds *Datasource) RelationRelations(ctx context.Context, id osm.RelationID, opts ...FeatureOption) (osm.Relations, error) {
	params, err := featureOptions(opts)
	if err != nil {
		return nil, err
	}

	o := &osm.OSM{}
	url := fmt.Sprintf("%s/relation/%d/relations?%s", ds.baseURL(), id, params)
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	return o.Relations, nil
}

// RelationFull returns the relation and its nodes for the latest version the relation.
func (ds *Datasource) RelationFull(ctx context.Context, id osm.RelationID, opts ...FeatureOption) (*osm.OSM, error) {
	params, err := featureOptions(opts)
	if err != nil {
		return nil, err
	}

	o := &osm.OSM{}
	url := fmt.Sprintf("%s/relation/%d/full?%s", ds.baseURL(), id, params)
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	return o, nil
}

// RelationVersion returns the specific version of the relation from the osm rest api.
func (ds *Datasource) RelationVersion(ctx context.Context, id osm.RelationID, v int) (*osm.Relation, error) {
	o := &osm.OSM{}
	url := fmt.Sprintf("%s/relation/%d/%d", ds.baseURL(), id, v)
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	if l := len(o.Relations); l != 1 {
		return nil, fmt.Errorf("wrong number of relations, expected 1, got %v", l)
	}

	return o.Relations[0], nil
}

// RelationHistory returns all the versions of the relation from the osm rest api.
func (ds *Datasource) RelationHistory(ctx context.Context, id osm.RelationID) (osm.Relations, error) {
	o := &osm.OSM{}
	url := fmt.Sprintf("%s/relation/%d/history", ds.baseURL(), id)
	if err := ds.getFromAPI(ctx, url, &o); err != nil {
		return nil, err
	}

	return o.Relations, nil
}

// Relation returns the latest version of the relation from the osm rest api.
// Delegates to the DefaultDatasource and uses its http.Client to make the request.
func Relation(ctx context.Context, id osm.RelationID, opts ...FeatureOption) (*osm.Relation, error) {
	return DefaultDatasource.Relation(ctx, id, opts...)
}

// Relations returns the latest version of the relations from the osm rest api.
// Delegates to the DefaultDatasource and uses its http.Client to make the request.
func Relations(ctx context.Context, ids []osm.RelationID, opts ...FeatureOption) (osm.Relations, error) {
	return DefaultDatasource.Relations(ctx, ids, opts...)
}

// RelationRelations returns all relations a relation is part of.
// There is no error if the element does not exist.
// Delegates to the DefaultDatasource and uses its http.Client to make the request.
func RelationRelations(ctx context.Context, id osm.RelationID, opts ...FeatureOption) (osm.Relations, error) {
	return DefaultDatasource.RelationRelations(ctx, id, opts...)
}

// RelationFull returns the relation and its nodes for the latest version the relation.
// Delegates to the DefaultDatasource and uses its http.Client to make the request.
func RelationFull(ctx context.Context, id osm.RelationID, opts ...FeatureOption) (*osm.OSM, error) {
	return DefaultDatasource.RelationFull(ctx, id, opts...)
}

// RelationVersion returns the specific version of the relation from the osm rest api.
// Delegates to the DefaultDatasource and uses its http.Client to make the request.
func RelationVersion(ctx context.Context, id osm.RelationID, v int) (*osm.Relation, error) {
	return DefaultDatasource.RelationVersion(ctx, id, v)
}

// RelationHistory returns all the versions of the relation from the osm rest api.
// Delegates to the DefaultDatasource and uses its http.Client to make the request.
func RelationHistory(ctx context.Context, id osm.RelationID) (osm.Relations, error) {
	return DefaultDatasource.RelationHistory(ctx, id)
}
