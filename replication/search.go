package replication

import "context"

type stater struct {
	Min     uint64
	Current func(context.Context) (*State, error)
	State   func(context.Context, uint64) (*State, error)
}
