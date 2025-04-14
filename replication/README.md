# osm/replication [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/osm)](https://pkg.go.dev/github.com/pchchv/osm/replication)

Package `replication` handles fetching the Minute, Hour, Day and Changeset replication and the associated state value from [Planet OSM](http://planet.osm.org).

I. e., to fetch the current Minute replication state:

```go
num, fullState, err := replication.CurrentMinuteState(ctx)
```

This is the data in [http://planet.osm.org/replication/minute/state.txt](http://planet.osm.org/replication/minute/state.txt) updated every minute.

Once recognized the change number, fetch the change using:

```go
change, err := replication.Minute(ctx, num)
```

## Finding sequences numbers by timestamp

It's also possible to find the sequence number by timestamp.
These calls make multiple requests for state to find the one matching the given timestamp.

```go
MinuteStateAt(ctx context.Context, timestamp time.Time) (MinuteSeqNum, *State, error)
HourStateAt(ctx context.Context, timestamp time.Time) (HourSeqNum, *State, error)
DayStateAt(ctx context.Context, timestamp time.Time) (DaySeqNum, *State, error)
ChangesetStateAt(ctx context.Context, timestamp time.Time) (ChangesetSeqNum, *State, error)
```
