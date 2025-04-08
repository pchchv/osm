package core

// childLoc references a location of a child in the parents + children.
type childLoc struct {
	Parent int
	Index  int
}

type childLocs []childLoc
