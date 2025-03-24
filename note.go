package osm

import "time"

// NoteID is the unique identifier for an osm note.
type NoteID int64

// ObjectID is a helper returning the object id for this note id.
func (id NoteID) ObjectID() ObjectID {
	return ObjectID(noteMask | (id << versionBits))
}

// Date is an object to decode the date format used in the osm notes xml api.
// Format: '2006-01-02 15:04:05 MST'.
type Date struct {
	time.Time
}
