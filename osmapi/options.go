package osmapi

// FeatureOption used when fetching a feature or a set of different features.
type FeatureOption interface {
	applyFeature([]string) ([]string, error)
}
