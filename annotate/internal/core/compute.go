package core

import "github.com/pchchv/osm"

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
