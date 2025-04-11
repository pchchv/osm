package replication

import (
	"fmt"
	"time"
)

var (
	_        = SeqNum(MinuteSeqNum(0)).private // for linters
	_ SeqNum = MinuteSeqNum(0)
	_ SeqNum = HourSeqNum(0)
	_ SeqNum = DaySeqNum(0)
)

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

// String returns 'minute/%d'.
func (n MinuteSeqNum) String() string {
	return fmt.Sprintf("minute/%d", n)
}

// Dir returns the directory of this data on planet osm.
func (n MinuteSeqNum) Dir() string {
	return "minute"
}

// Uint64 returns the seq num as a uint64 type.
func (n MinuteSeqNum) Uint64() uint64 {
	return uint64(n)
}

func (n MinuteSeqNum) private() {}

// HourSeqNum indicates the sequence of the hourly diff replication found here:
// http://planet.osm.org/replication/hour
type HourSeqNum uint64

// String returns 'hour/%d'.
func (n HourSeqNum) String() string {
	return fmt.Sprintf("hour/%d", n)
}

// Dir returns the directory of this data on planet osm.
func (n HourSeqNum) Dir() string {
	return "hour"
}

// Uint64 returns the seq num as a uint64 type.
func (n HourSeqNum) Uint64() uint64 {
	return uint64(n)
}

func (n HourSeqNum) private() {}

// DaySeqNum indicates the sequence of the daily diff replication found here:
// http://planet.osm.org/replication/day
type DaySeqNum uint64

// String returns 'day/%d'.
func (n DaySeqNum) String() string {
	return fmt.Sprintf("day/%d", n)
}

// Dir returns the directory of this data on planet osm.
func (n DaySeqNum) Dir() string {
	return "day"
}

// Uint64 returns the seq num as a uint64 type.
func (n DaySeqNum) Uint64() uint64 {
	return uint64(n)
}

func (n DaySeqNum) private() {}

// SeqNum is an interface type that includes MinuteSeqNum,
// HourSeqNum and DaySeqNum.
// This is an experiment to implement a sum type,
// a type that can be one of several things only.
type SeqNum interface {
	fmt.Stringer
	Dir() string
	Uint64() uint64
	private()
}
