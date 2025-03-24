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
