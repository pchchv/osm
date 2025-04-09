package core

import (
	"context"
	"errors"

	"github.com/pchchv/osm"
)

var ErrNotFound             = errors.New("not found")

// TestDS implements a datasource for testing.
type TestDS struct {
	data map[osm.FeatureID]ChildList
}

// Get returns the history in ChildList form.
func (tds *TestDS) Get(ctx context.Context, id osm.FeatureID) (ChildList, error) {
	if tds.data == nil {
		return nil, ErrNotFound
	}

	v := tds.data[id]
	if v == nil {
		return nil, ErrNotFound
	}

	return v, nil
}

// MustGet is used by tests only to simplify some code.
func (tds *TestDS) MustGet(id osm.FeatureID) ChildList {
	v, err := tds.Get(context.TODO(), id)
	if err != nil {
		panic(err)
	}

	return v
}
