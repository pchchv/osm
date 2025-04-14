package replication

import "time"

var baseTime = time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)

// buildState is helper to build "valid" state files.
func buildState(n int) *State {
	d := time.Duration(n)
	return &State{SeqNum: uint64(n), Timestamp: baseTime.Add(d * time.Hour)}
}
