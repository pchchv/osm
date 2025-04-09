package annotate

import "github.com/pchchv/osm/annotate/internal/core"

// Option is a parameter that can be used for annotating.
type Option func(*core.Options) error
