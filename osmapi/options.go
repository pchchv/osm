package osmapi

import (
	"strings"
	"time"
)

// FeatureOption used when fetching a feature or a set of different features.
type FeatureOption interface {
	applyFeature([]string) ([]string, error)
}

type at struct {
	t time.Time
}

func (o *at) applyFeature(p []string) ([]string, error) {
	return append(p, "at="+o.t.UTC().Format("2006-01-02T15:04:05Z")), nil
}

// At adds an `at=2006-01-02T15:04:05Z` parameter to the request.
// The osm.fyi supports requesting features and maps as they were at the given time.
func At(t time.Time) FeatureOption {
	return &at{t}
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
