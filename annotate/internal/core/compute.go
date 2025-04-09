package core

import (
	"time"

	"github.com/pchchv/osm"
	"github.com/pchchv/osm/annotate/shared"
)

// Options allow for passing som parameters to the matching process.
type Options struct {
	Threshold             time.Duration
	IgnoreInconsistency   bool
	IgnoreMissingChildren bool
	ChildFilter           func(osm.FeatureID) bool
}

// childLoc references a location of a child in the parents + children.
type childLoc struct {
	Parent int
	Index  int
}

type childLocs []childLoc

func (locs childLocs) GroupByParent() (result []childLocs) {
	for len(locs) > 0 {
		var end int
		p := locs[0].Parent
		for end < len(locs) && locs[end].Parent == p {
			end++
		}

		result = append(result, locs[:end])
		locs = locs[end:]
	}

	return result
}

// mapChildLocs builds a cache of a where a child is in a set of parents.
func mapChildLocs(parents []Parent, filter func(osm.FeatureID) bool) map[osm.FeatureID]childLocs {
	result := make(map[osm.FeatureID]childLocs)
	for i, p := range parents {
		refs, annotated := p.Refs()
		for j, fid := range refs {
			if annotated[j] && filter != nil && !filter(fid) {
				continue
			}

			if result[fid] == nil {
				result[fid] = make([]childLoc, 0, len(parents))
			}

			result[fid] = append(result[fid], childLoc{Parent: i, Index: j})
		}
	}

	return result
}

func nextVersionIndex(current *shared.Child, child ChildList, nextParent Parent, opts *Options) int {
	if nextParent == nil {
		// no next parent version,
		// so is needed to include all future versions of this child
		return child[len(child)-1].VersionIndex + 1
	}

	next := child.FindVisible(
		nextParent.ChangesetID(),
		timeThresholdParent(nextParent, 0),
		opts.Threshold,
	)

	if next != nil {
		// if the child was updated enough before the next parent include it in the minor versions
		if timeThreshold(next, 0).Before(timeThresholdParent(nextParent, -opts.Threshold)) {
			return next.VersionIndex + 1
		}

		return next.VersionIndex
	}

	// child is one of:
	// - not in the next parent version,
	// - next parent is deleted,
	// - data inconsistency and not visible for the next parent
	// so is needed to know what was the last available before the next parent

	// this timestamp helps when creating updates, and need to make sure this it is:
	// - 1 threshold before the next parent
	// - not before the current child timestamp
	ts := timeThresholdParent(nextParent, -opts.Threshold)
	if current != nil && !ts.After(timeThreshold(current, 0)) { // before and equal still matches
		// visible in current but child and next parent are within the same threshold, no updates
		// i.e. next version is same as current version
		return 0 // no updates
	}

	// current child and next parent are far apart
	next = child.VersionBefore(ts)
	if next == nil {
		// missing at current and next parent
		return 0 // no updates
	}

	// visible or not, it needs to be included
	// nonvisible versions of this child will be filtered below
	return next.VersionIndex + 1
}
