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

// Changeset is a set of metadata around a set of osm changes.
type Changeset struct {
	XMLName       xmlNameJSONTypeCS    `xml:"changeset" json:"type"`
	ID            ChangesetID          `xml:"id,attr" json:"id"`
	User          string               `xml:"user,attr" json:"user,omitempty"`
	UserID        UserID               `xml:"uid,attr" json:"uid,omitempty"`
	CreatedAt     time.Time            `xml:"created_at,attr" json:"created_at"`
	ClosedAt      time.Time            `xml:"closed_at,attr" json:"closed_at"`
	Open          bool                 `xml:"open,attr" json:"open"`
	ChangesCount  int                  `xml:"num_changes,attr,omitempty" json:"num_changes,omitempty"`
	MinLat        float64              `xml:"min_lat,attr" json:"min_lat,omitempty"`
	MaxLat        float64              `xml:"max_lat,attr" json:"max_lat,omitempty"`
	MinLon        float64              `xml:"min_lon,attr" json:"min_lon,omitempty"`
	MaxLon        float64              `xml:"max_lon,attr" json:"max_lon,omitempty"`
	CommentsCount int                  `xml:"comments_count,attr,omitempty" json:"comments_count,omitempty"`
	Tags          Tags                 `xml:"tag" json:"tags,omitempty"`
	Discussion    *ChangesetDiscussion `xml:"discussion,omitempty" json:"discussion,omitempty"`
	Change        *Change              `xml:"-" json:"change,omitempty"`
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
