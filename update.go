package osm

import "time"

// Update is a change to children of a way or relation.
// The child type, id, ref and/or role are the same as the child at the given index.
// Lon/Lat are only updated for nodes.
type Update struct {
	Index       int         `xml:"index,attr" json:"index"`
	Version     int         `xml:"version,attr" json:"version"`
	Timestamp   time.Time   `xml:"timestamp,attr" json:"timestamp"` // committed at time if time > CommitInfoStart or the element timestamp if before that date
	ChangesetID ChangesetID `xml:"changeset,attr,omitempty" json:"changeset,omitempty"`
	Lat         float64     `xml:"lat,attr,omitempty" json:"lat,omitempty"`
	Lon         float64     `xml:"lon,attr,omitempty" json:"lon,omitempty"`
	Reverse     bool        `xml:"reverse,attr,omitempty" json:"reverse,omitempty"`
}

// Updates are collections of updates.
type Updates []Update

// UpTo returns the subset of updates taking place upto and on the given time.
func (us Updates) UpTo(t time.Time) (result Updates) {
	for _, u := range us {
		if u.Timestamp.After(t) {
			continue
		}

		result = append(result, u)
	}

	return
}
