package osm

import (
	"errors"

	"github.com/pchchv/geo/maptile"
)

// Bounds are the bounds of osm data as defined in the xml file.
type Bounds struct {
	MinLat float64 `xml:"minlat,attr"`
	MaxLat float64 `xml:"maxlat,attr"`
	MinLon float64 `xml:"minlon,attr"`
	MaxLon float64 `xml:"maxlon,attr"`
}

// NewBoundsFromTile creates a bound given an online map tile index.
func NewBoundsFromTile(t maptile.Tile) (*Bounds, error) {
	maxIndex := uint32(1 << t.Z)
	if t.X >= maxIndex {
		return nil, errors.New("osm: x index out of range for this zoom")
	} else if t.Y >= maxIndex {
		return nil, errors.New("osm: y index out of range for this zoom")
	}

	b := t.Bound()
	return &Bounds{
		MinLat: b.Min.Lat(),
		MaxLat: b.Max.Lat(),
		MinLon: b.Min.Lon(),
		MaxLon: b.Max.Lon(),
	}, nil
}

// ObjectID returns the bounds type but with 0 id.
// Since id doesn't make sense.
// This is here to implement the Object interface since it technically is an osm object type.
// It also allows bounds to be returned via the osmxml.Scanner.
func (b *Bounds) ObjectID() ObjectID {
	return ObjectID(boundsMask)
}
