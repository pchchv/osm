package annotate

import (
	"fmt"

	"github.com/pchchv/osm"
)

// NoHistoryError is returned if there
// is no entry in the history map for a specific child.
type NoHistoryError struct {
	ID osm.FeatureID
}

// Error returns a pretty string of the error.
func (e *NoHistoryError) Error() string {
	return fmt.Sprintf("element history not found for %v", e.ID)
}
