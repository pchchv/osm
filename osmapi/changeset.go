package osmapi

import (
	"context"
	"fmt"

	"github.com/pchchv/osm"
)

// Changeset returns a given changeset from the osm rest api.
func (ds *Datasource) Changeset(ctx context.Context, id osm.ChangesetID) (*osm.Changeset, error) {
	url := fmt.Sprintf("%s/changeset/%d", ds.baseURL(), id)
	return ds.getChangeset(ctx, url)
}

func (ds *Datasource) getChangeset(ctx context.Context, url string) (*osm.Changeset, error) {
	css := &osm.OSM{}
	if err := ds.getFromAPI(ctx, url, &css); err != nil {
		return nil, err
	}

	if l := len(css.Changesets); l != 1 {
		return nil, fmt.Errorf("wrong number of changesets, expected 1, got %v", l)
	}

	return css.Changesets[0], nil
}
