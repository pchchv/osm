package core

import (
	"time"

	"github.com/pchchv/osm"
	"github.com/pchchv/osm/annotate/shared"
)

type ChildList []*shared.Child

func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}

	return d
}

func timeThreshold(c *shared.Child, esp time.Duration) time.Time {
	if c.Committed.Before(osm.CommitInfoStart) {
		return c.Timestamp.Add(esp)
	}

	return c.Committed
}
