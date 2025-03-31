package osmtest

import "github.com/pchchv/osm"

// Scanner implements the osm.Scanner interface with just a list of objects.
type Scanner struct {
	// ScanError can be used to trigger an error.
	// If non-nil, Next() will return false and Err() will
	// return this error.
	ScanError error
	offset    int
	objects   osm.Objects
}

// NewScanner creates a new test scanner useful for test stubbing.
func NewScanner(objects osm.Objects) *Scanner {
	return &Scanner{
		offset:  -1,
		objects: objects,
	}
}
