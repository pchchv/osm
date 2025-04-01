package osmpbf

import (
	"time"

	"github.com/pchchv/osm"
)

// Header contains the contents of the header in the pbf file.
type Header struct {
	Bounds               *osm.Bounds
	RequiredFeatures     []string
	OptionalFeatures     []string
	WritingProgram       string
	Source               string
	ReplicationTimestamp time.Time
	ReplicationSeqNum    uint64
	ReplicationBaseURL   string
}

// oPair is a group sent over the channel from the decoder goroutines.
// It will contain the list of all objects.
type oPair struct {
	Offset  int64
	Objects []osm.Object
	Err     error
}
