package core

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
