package osm

import (
	"encoding/xml"
	"time"
)

const dateLayout = "2006-01-02 15:04:05 MST"

var (
	NoteOpen           NoteStatus        = "open"
	NoteClosed         NoteStatus        = "closed"
	NoteCommentOpened  NoteCommentAction = "opened"
	NoteCommentClosed  NoteCommentAction = "closed"
	NoteCommentComment NoteCommentAction = "commented"
)

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

// MarshalJSON returns null if the date is empty.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte(`null`), nil
	}

	return marshalJSON(d.Time)
}

// MarshalXML is meant to encode the time.Time into the
// osm note date formation of '2006-01-02 15:04:05 MST'.
func (d Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(d.Format(dateLayout), start)
}

// UnmarshalXML is meant to decode the osm note date formation of
// '2006-01-02 15:04:05 MST' into a time.Time object.
func (d *Date) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) (err error) {
	var s string
	if err = dec.DecodeElement(&s, &start); err != nil {
		return
	}

	d.Time, err = time.Parse(dateLayout, s)
	return err
}

// NoteCommentAction are actions that a note comment took.
type NoteCommentAction string

// NoteStatus is the status of the note.
type NoteStatus string

// NoteComment is a comment on a note.
type NoteComment struct {
	XMLName xml.Name          `xml:"comment" json:"-"`
	Date    Date              `xml:"date" json:"date"`
	UserID  UserID            `xml:"uid" json:"uid,omitempty"`
	User    string            `xml:"user" json:"user,omitempty"`
	UserURL string            `xml:"user_url" json:"user_url,omitempty"`
	Action  NoteCommentAction `xml:"action" json:"action"`
	Text    string            `xml:"text" json:"text"`
	HTML    string            `xml:"html" json:"html"`
}

// Note is information for other mappers dropped at a map location.
type Note struct {
	XMLName     xmlNameJSONTypeNote `xml:"note" json:"type"`
	ID          NoteID              `xml:"id" json:"id"`
	Lat         float64             `xml:"lat,attr" json:"lat"`
	Lon         float64             `xml:"lon,attr" json:"lon"`
	URL         string              `xml:"url" json:"url,omitempty"`
	CommentURL  string              `xml:"comment_url" json:"comment_url,omitempty"`
	CloseURL    string              `xml:"close_url" json:"close_url,omitempty"`
	ReopenURL   string              `xml:"reopen_url" json:"reopen_url,omitempty"`
	DateCreated Date                `xml:"date_created" json:"date_created"`
	DateClosed  Date                `xml:"date_closed" json:"date_closed,omitempty"`
	Status      NoteStatus          `xml:"status" json:"status,omitempty"`
	Comments    []*NoteComment      `xml:"comments>comment" json:"comments"`
}

// ObjectID returns the object id of the note.
func (n *Note) ObjectID() ObjectID {
	return n.ID.ObjectID()
}

// Notes is a collection of notes with some helpers attached.
type Notes []*Note
