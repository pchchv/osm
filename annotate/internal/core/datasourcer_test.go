package core

import "github.com/pchchv/osm"

// TestDS implements a datasource for testing.
type TestDS struct {
	data map[osm.FeatureID]ChildList
}
