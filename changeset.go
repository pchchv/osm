package osm

import (
	"encoding/xml"
	"time"
)

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

// ChangesetDiscussion is a conversation about a changeset.
type ChangesetDiscussion struct {
	Comments []*ChangesetComment `xml:"comment" json:"comments"`
}

// MarshalXML implements the xml.Marshaller method to exclude this
// whole element if the comments are empty.
func (csd ChangesetDiscussion) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	if len(csd.Comments) == 0 {
		return nil
	}

	if err = e.EncodeToken(start); err != nil {
		return
	}

	t := xml.StartElement{Name: xml.Name{Local: "comment"}}
	if err = e.EncodeElement(csd.Comments, t); err != nil {
		return
	}

	return e.EncodeToken(start.End())
}
