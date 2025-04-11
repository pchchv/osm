package replication

import "time"

// State returns information about the current replication state.
type State struct {
	SeqNum        uint64    `json:"seq_num"`
	Timestamp     time.Time `json:"timestamp"`
	TxnMax        int       `json:"txn_max,omitempty"`
	TxnMaxQueried int       `json:"txn_max_queries,omitempty"`
}

// MinuteSeqNum indicates the sequence of the minutely diff replication found here:
// http://planet.osm.org/replication/minute
type MinuteSeqNum uint64
