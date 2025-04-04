package osmapi

import "strings"

// FeatureOption used when fetching a feature or a set of different features.
type FeatureOption interface {
	applyFeature([]string) ([]string, error)
}

func featureOptions(opts []FeatureOption) (string, error) {
	if len(opts) == 0 {
		return "", nil
	}

	params := make([]string, 0, len(opts))
	for _, o := range opts {
		var err error
		if params, err = o.applyFeature(params); err != nil {
			return "", err
		}
	}

	return strings.Join(params, "&"), nil
}
