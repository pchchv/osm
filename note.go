package osm

// NoteID is the unique identifier for an osm note.
type NoteID int64

// ObjectID is a helper returning the object id for this note id.
func (id NoteID) ObjectID() ObjectID {
	return ObjectID(noteMask | (id << versionBits))
}
