package core

import (
	"context"
	"fmt"
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

// Datasourcer acts as a datasource,
// allowing fetching of children as needed.
type Datasourcer interface {
	Get(ctx context.Context, id osm.FeatureID) (ChildList, error)
	NotFound(err error) bool
}

// Compute first computes the exact version of the children in each parent,
// then returns a set of updates for each version of the parent.
func Compute(ctx context.Context, parents []Parent, histories Datasourcer, opts *Options) ([]osm.Updates, error) {
	if opts == nil {
		opts = &Options{}
	}

	results := make([]osm.Updates, len(parents))
	for fid, locations := range mapChildLocs(parents, opts.ChildFilter) {
		child, err := histories.Get(ctx, fid)
		if err != nil {
			if !histories.NotFound(err) {
				return nil, err
			}

			if opts.IgnoreMissingChildren {
				continue
			}

			return nil, &NoHistoryError{ChildID: fid}
		}

		for _, locs := range locations.GroupByParent() {
			// figure out the parent and the next parent
			parentIndex := locs[0].Parent
			parent := parents[parentIndex]
			if !parent.Visible() {
				continue
			}

			var nextParent Parent
			if parentIndex < len(parents)-1 {
				nextParent = parents[parentIndex+1]
			}

			// get the current child
			c := child.FindVisible(
				parent.ChangesetID(),
				timeThresholdParent(parent, 0),
				opts.Threshold,
			)
			if c == nil && !opts.IgnoreInconsistency {
				return nil, &NoVisibleChildError{
					ChildID:   fid,
					Timestamp: timeThresholdParent(parent, 0)}
			}

			// straight up set this child on major version
			for _, cl := range locs {
				parent.SetChild(cl.Index, c)
			}

			// nextVersionIndex figures out what version of this child
			// is present in the next parent version
			nextVersion := nextVersionIndex(c, child, nextParent, opts)

			start := 0
			if c != nil {
				start = c.VersionIndex + 1
			} else {
				// current child is not defined, is next child
				next := child.VersionBefore(timeThresholdParent(parent, 0))
				if next == nil {
					start = 0
				} else {
					start = next.VersionIndex + 1
				}
			}

			var updates osm.Updates
			for k := start; k < nextVersion; k++ {
				if child[k].Visible {
					// it's possible for this child to be present at multiple locations in the parent
					for _, cl := range locs {
						u := child[k].Update()
						u.Index = cl.Index
						updates = append(updates, u)
					}
				} else {
					// child has become not-visible between parent version
					// this is a data inconsistency that can happen in old data
					// i.e. pre element versioning
					//
					// see node 321452894, changed 7 times in the same changeset,
					// version 5 was a delete
					// (also node 65172196)
					if !opts.IgnoreInconsistency {
						return nil, fmt.Errorf("%v: %v: child deleted between parent versions", parent.ID(), fid)
					}
				}
			}

			// there is everything needed for this parent version
			results[parentIndex] = append(results[parentIndex], updates...)
		}
	}

	for _, r := range results {
		r.SortByIndex()
	}

	return results, nil
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
