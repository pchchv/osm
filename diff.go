package osm

// The different types of diff actions.
const (
	ActionCreate ActionType = "create"
	ActionModify ActionType = "modify"
	ActionDelete ActionType = "delete"
)

// ActionType is a strong type for the different diff actions.
type ActionType string
