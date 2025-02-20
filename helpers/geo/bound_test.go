package geo

import "testing"

func TestBoundCenter(t *testing.T) {
	bound := Bound{Min: Point{1, 1}, Max: Point{2, 2}}
	if c := bound.Center(); !c.Equal(Point{1.5, 1.5}) {
		t.Errorf("incorrect center: %v", c)
	}
}

func TestBoundContains(t *testing.T) {
	bound := Bound{Min: Point{-2, -1}, Max: Point{2, 1}}
	cases := []struct {
		name   string
		point  Point
		result bool
	}{
		{
			name:   "middle",
			point:  Point{0, 0},
			result: true,
		},
		{
			name:   "left border",
			point:  Point{-1, 0},
			result: true,
		},
		{
			name:   "ne corner",
			point:  Point{2, 1},
			result: true,
		},
		{
			name:   "above",
			point:  Point{0, 3},
			result: false,
		},
		{
			name:   "below",
			point:  Point{0, -3},
			result: false,
		},
		{
			name:   "left",
			point:  Point{-3, 0},
			result: false,
		},
		{
			name:   "right",
			point:  Point{3, 0},
			result: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if v := bound.Contains(tc.point); v != tc.result {
				t.Errorf("incorrect contains: %v != %v", v, tc.result)
			}
		})
	}
}

func TestBoundIntersects(t *testing.T) {
	bound := Bound{Min: Point{0, 2}, Max: Point{1, 3}}
	cases := []struct {
		name   string
		bound  Bound
		result bool
	}{
		{
			name:   "outside, top right",
			bound:  Bound{Min: Point{5, 7}, Max: Point{6, 8}},
			result: false,
		},
		{
			name:   "outside, top left",
			bound:  Bound{Min: Point{-6, 7}, Max: Point{-5, 8}},
			result: false,
		},
		{
			name:   "outside, above",
			bound:  Bound{Min: Point{0, 7}, Max: Point{0.5, 8}},
			result: false,
		},
		{
			name:   "over the middle",
			bound:  Bound{Min: Point{0, 0.5}, Max: Point{1, 4}},
			result: true,
		},
		{
			name:   "over the left",
			bound:  Bound{Min: Point{-1, 2}, Max: Point{1, 4}},
			result: true,
		},
		{
			name:   "completely inside",
			bound:  Bound{Min: Point{0.3, 2.3}, Max: Point{0.6, 2.6}},
			result: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := bound.Intersects(tc.bound)
			if v != tc.result {
				t.Errorf("incorrect result: %v != %v", v, tc.result)
			}
		})
	}

	a := Bound{Min: Point{7, 6}, Max: Point{8, 7}}
	b := Bound{Min: Point{6.1, 6.1}, Max: Point{8.1, 8.1}}
	if !a.Intersects(b) || !b.Intersects(a) {
		t.Errorf("expected to intersect")
	}

	a = Bound{Min: Point{1, 2}, Max: Point{4, 3}}
	b = Bound{Min: Point{2, 1}, Max: Point{3, 4}}
	if !a.Intersects(b) || !b.Intersects(a) {
		t.Errorf("expected to intersect")
	}
}
