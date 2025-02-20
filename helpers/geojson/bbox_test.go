package geojson

import (
	"testing"

	"github.com/pchchv/osm/helpers/geo"
)

func TestBBoxValid(t *testing.T) {
	cases := []struct {
		name   string
		bbox   BBox
		result bool
	}{
		{
			name:   "true for 4 length array",
			bbox:   []float64{1, 2, 3, 4},
			result: true,
		},
		{
			name:   "true for 3d box",
			bbox:   []float64{1, 2, 3, 4, 5, 6},
			result: true,
		},
		{
			name:   "false for nil box",
			bbox:   nil,
			result: false,
		},
		{
			name:   "false for short array",
			bbox:   []float64{1, 2, 3},
			result: false,
		},
		{
			name:   "false for incorrect length array",
			bbox:   []float64{1, 2, 3, 4, 5},
			result: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if v := tc.bbox.Valid(); v != tc.result {
				t.Errorf("incorrect result: %v != %v", v, tc.result)
			}
		})
	}
}

func TestBBoxBound(t *testing.T) {
	cases := []struct {
		name   string
		bbox   BBox
		result geo.Bound
	}{
		{
			name:   "empty for invalid bbox",
			bbox:   []float64{1, 2, 3},
			result: geo.Bound{},
		},
		{
			name:   "correct order for 2d box",
			bbox:   []float64{1, 2, 3, 4},
			result: geo.Bound{Min: geo.Point{1, 2}, Max: geo.Point{3, 4}},
		},
		{
			name:   "correct order for 3d box",
			bbox:   []float64{1, 2, 3, 4, 5, 6},
			result: geo.Bound{Min: geo.Point{1, 2}, Max: geo.Point{4, 5}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if v := tc.bbox.Bound(); !v.Equal(tc.result) {
				t.Errorf("incorrect result: %v != %v", v, tc.result)
			}
		})
	}
}
