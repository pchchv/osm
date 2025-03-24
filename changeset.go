package osm

import "time"

// ChangesetID is the primary key for an osm changeset.
type ChangesetID int64

// ObjectID is a helper returning the object id for this changeset id.
func (id ChangesetID) ObjectID() ObjectID {
	return ObjectID(changesetMask | (id << versionBits))
}

// ChangesetComment is a specific comment in a changeset discussion.
type ChangesetComment struct {
	User      string    `xml:"user,attr" json:"user"`
	UserID    UserID    `xml:"uid,attr" json:"uid"`
	Timestamp time.Time `xml:"date,attr" json:"date"`
	Text      string    `xml:"text" json:"text"`
}
