package osm

var (
	conditionAll       conditionType = "all"
	conditionBlacklist conditionType = "blacklist"
	conditionWhitelist conditionType = "whitelist"
)

type conditionType string
