package osm

import "encoding/json"

// Tag is a key-value item attached to
// osm nodes, ways and relations.
type Tag struct {
	Key   string `xml:"k,attr"`
	Value string `xml:"v,attr"`
}

// Tags is a collection of Tag objects
// with some helper functions.
type Tags []Tag

// Map returns the tags as a key/value map.
func (ts Tags) Map() map[string]string {
	result := make(map[string]string, len(ts))
	for _, t := range ts {
		result[t.Key] = t.Value
	}

	return result
}

// Find returns the value for the key.
// Returns empty string if not found.
func (ts Tags) Find(k string) string {
	for _, t := range ts {
		if t.Key == k {
			return t.Value
		}
	}

	return ""
}

// FindTag returns the Tag for the given key.
// Can be used to determine if a key exists,
// even with an empty value.
// Returns nil if not found.
func (ts Tags) FindTag(k string) *Tag {
	for _, t := range ts {
		if t.Key == k {
			return &t
		}
	}

	return nil
}

// HasTag returns true if a tag exists for the given key.
func (ts Tags) HasTag(k string) bool {
	for _, t := range ts {
		if t.Key == k {
			return true
		}
	}

	return false
}

// MarshalJSON marshals tags as a key/value object,
// as defined by the overpass osmjson.
func (ts Tags) MarshalJSON() ([]byte, error) {
	return marshalJSON(ts.Map())
}

// UnmarshalJSON unmarshals tags from a key/value object,
// as defined by the overpass osmjson.
func (ts *Tags) UnmarshalJSON(data []byte) error {
	o := make(map[string]string)
	if err := json.Unmarshal(data, &o); err != nil {
		return err
	}

	tags := make(Tags, 0, len(o))
	for k, v := range o {
		tags = append(tags, Tag{Key: k, Value: v})
	}

	*ts = tags
	return nil
}
