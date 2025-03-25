package osm

import (
	"testing"

	"github.com/pchchv/geo/maptile"
)

func mustBounds(t *testing.T, x, y uint32, z maptile.Zoom) *Bounds {
	bounds, err := NewBoundsFromTile(maptile.New(x, y, z))
	if err != nil {
		t.Fatalf("invalid bounds: %e", err)
	}

	return bounds
}

func centroid(b *Bounds) *Node {
	return &Node{
		Lon: (b.MinLon + b.MaxLon) / 2,
		Lat: (b.MinLat + b.MaxLat) / 2,
	}
}
