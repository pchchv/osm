package replication

import "fmt"

// ChangesetSeqNum indicates the sequence of the changeset replication found here:
// https://planet.osm.org/replication/changesets/
type ChangesetSeqNum uint64

// Dir returns the directory of this data on planet osm.
func (n ChangesetSeqNum) Dir() string {
	return "changesets"
}

// String returns 'changeset/%d'.
func (n ChangesetSeqNum) String() string {
	return fmt.Sprintf("changeset/%d", n)
}

// Uint64 returns the seq num as a uint64 type.
func (n ChangesetSeqNum) Uint64() uint64 {
	return uint64(n)
}

