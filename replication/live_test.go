package replication

import (
	"context"
	"os"
	"testing"
)

func TestCurrentState(t *testing.T) {
	liveOnly(t)
	ctx := context.Background()
	if _, _, err := CurrentMinuteState(ctx); err != nil {
		t.Fatalf("request error: %e", err)
	}

	if _, _, err := CurrentHourState(ctx); err != nil {
		t.Fatalf("request error: %e", err)
	}

	if _, _, err := CurrentDayState(ctx); err != nil {
		t.Fatalf("request error: %e", err)
	}
}

func TestDownloadChanges(t *testing.T) {
	liveOnly(t)
	ctx := context.Background()
	if _, err := Minute(ctx, 10); err != nil {
		t.Fatalf("request error: %e", err)
	}

	if _, err := Hour(ctx, 10); err != nil {
		t.Fatalf("request error: %e", err)
	}

	if _, err := Day(ctx, 1); err != nil {
		t.Fatalf("request error: %e", err)
	}
}

func TestCurrentChangesetState(t *testing.T) {
	liveOnly(t)
	ctx := context.Background()
	if _, _, err := CurrentChangesetState(ctx); err != nil {
		t.Fatalf("request error: %e", err)
	}
}

func TestChangesets(t *testing.T) {
	liveOnly(t)
	ctx := context.Background()
	sets, err := Changesets(ctx, 100)
	if err != nil {
		t.Fatalf("request error: %e", err)
	}

	if l := len(sets); l != 12 {
		t.Errorf("incorrect number of changesets: %v", l)
	}
}

func TestChangesetState(t *testing.T) {
	liveOnly(t)
	ctx := context.Background()
	state, err := ChangesetState(ctx, 5001990)
	if err != nil {
		t.Fatalf("request error: %e", err)
	}

	if state.SeqNum != 5001990 {
		t.Errorf("incorrect state: %+v", state)
	}

	// current state
	n, state, err := CurrentChangesetState(ctx)
	if err != nil {
		t.Fatalf("request error: %e", err)
	}

	changes, err := Changesets(ctx, n)
	if err != nil {
		t.Fatalf("request error: %e", err)
	}

	for _, c := range changes {
		if c.CreatedAt.After(state.Timestamp) {
			t.Errorf("data is after the state file?")
		}
	}
}

func liveOnly(t testing.TB) {
	if os.Getenv("LIVE_TEST") != "true" {
		t.Skipf("skipping live test, set LIVE_TEST=true to enable")
	}
}
