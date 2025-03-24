package osm

// ChangesetID is the primary key for an osm changeset.
type ChangesetID int64

// ObjectID is a helper returning the object id for this changeset id.
func (id ChangesetID) ObjectID() ObjectID {
	return ObjectID(changesetMask | (id << versionBits))
}
