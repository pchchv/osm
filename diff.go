package osm

// The different types of diff actions.
const (
	ActionCreate ActionType = "create"
	ActionModify ActionType = "modify"
	ActionDelete ActionType = "delete"
)

// ActionType is a strong type for the different diff actions.
type ActionType string

// Action is an explicit create,
// modify or delete action with old and new data if applicable.
// Different properties of this struct will be populated depending on the action.
//
//	Create: da.OSM will contain the new element
//	Modify: da.Old and da.New will contain the old and new elements.
//	Delete: da.Old and da.New will contain the old and new elements.
type Action struct {
	Type ActionType `xml:"type,attr"`
	*OSM `xml:",omitempty"`
	Old  *OSM `xml:"old,omitempty"`
	New  *OSM `xml:"new,omitempty"`
}
