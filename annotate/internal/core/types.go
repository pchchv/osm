package core

import (
	"time"

	"github.com/pchchv/osm"
	"github.com/pchchv/osm/annotate/shared"
)

// Parent holds children,
// i.e. ways have nodes as children and relations can have nodes,
// ways and relations as children.
type Parent interface {
	ID() osm.FeatureID // used for logging
	ChangesetID() osm.ChangesetID
	Version() int
	Visible() bool
	Timestamp() time.Time
	Committed() time.Time
	// Refs returns normalized information about the children.
	// Currently this is the feature ids and if it is already annotated.
	// Note: we auto-annotate all unannotated children if they would have
	// been filtered out.
	Refs() (osm.FeatureIDs, []bool)
	SetChild(idx int, c *shared.Child)
}

type ChildList []*shared.Child

// VersionBefore finds the last child before a given time.
func (cl ChildList) VersionBefore(end time.Time) (latest *shared.Child) {
	for _, c := range cl {
		if !timeThreshold(c, 0).Before(end) {
			break
		}

		latest = c
	}

	return
}

// FindVisible locates the child visible at the given time.
// If 'at' is at or after osm.CommitInfoStart, the committed time is used to determine visibility.
// If 'at' is before, a range +-eps around the give time.
// The closes visible node in that range is returned, or the previous node if it visible.
// Children after 'at', but within eps, must have the same changeset id as provided (parent).
// If the previous node is not visible or does not exit, nil is returned.
func (cl ChildList) FindVisible(cid osm.ChangesetID, at time.Time, eps time.Duration) *shared.Child {
	var nearest *shared.Child
	var diff time.Duration = -1
	start := at.Add(-eps)
	for _, c := range cl {
		if c.Committed.Before(osm.CommitInfoStart) {
			// more complicated logic for early data
			offset := c.Timestamp.Sub(start)
			visible := c.Visible
			// if this node is after the end then it's over
			if offset > 2*eps {
				break
			}

			// if in front of the start set with the latest node
			if offset < 0 {
				if visible {
					nearest = c
				} else {
					nearest = nil
				}

				continue
			}

			// in the range
			d := absDuration(offset - eps)
			if diff < 0 || (d <= diff) {
				// first within range, set if not visible
				if diff == -1 && !visible && offset == 0 {
					nearest = nil
				}

				// update only the closest ones if they are visible,
				// because it is necessary that the closest ones are visible within the range
				if visible {
					if offset <= eps {
						// if before at, pick it
						nearest = c
					} else if c.ChangesetID == cid {
						// if after at, changeset must be same
						nearest = c
					} else {
						// after at, not same changeset, ignore.
						continue
					}
				}

				diff = d
			}
		} else {
			// simpler logic
			// if committed is on or before 'at' consider that element
			if c.Committed.After(at) {
				break
			}

			if c.Visible {
				nearest = c
			} else {
				nearest = nil
			}
		}
	}

	return nearest
}

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

func timeThresholdParent(p Parent, esp time.Duration) time.Time {
	if p.Committed().Before(osm.CommitInfoStart) {
		return p.Timestamp().Add(esp)
	}

	return p.Committed()
}
