package osm

var (
	polyConditions     []polyCondition
	conditionAll       conditionType = "all"
	conditionBlacklist conditionType = "blacklist"
	conditionWhitelist conditionType = "whitelist"
)

type conditionType string

type polyCondition struct {
	Key       string        `json:"key"`
	Condition conditionType `json:"polygon"`
	Values    []string      `json:"values"`
}
