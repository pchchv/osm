package replication

import (
	"context"
	"time"
)

type stater struct {
	Min     uint64
	Current func(context.Context) (*State, error)
	State   func(context.Context, uint64) (*State, error)
}

func findBound(ctx context.Context, s *stater, upper *State, timestamp time.Time) (*State, *State, error) {
	var lowerID uint64 = 1
	var lower *State
	var err error
	// is needed to find the lower bound
	for lower == nil {
		lower, err = s.State(ctx, lowerID)
		if err != nil && !NotFound(err) {
			return nil, nil, err
		}

		if lower != nil && lower.Timestamp.After(timestamp) {
			if lower.SeqNum+1 >= upper.SeqNum {
				return lower, upper, nil // edge case if there are only two sequence numbers
			}

			// in our search for lower we found a new upper bound
			upper = lower
			lower = nil
			lowerID = 1
		}

		if lower != nil {
			break
		}

		// no lower yet, so try a higher id (binary search wise)
		newID := (lowerID + upper.SeqNum) / 2
		if newID <= lowerID {
			// nothing suitable found, so upper is probably the best that can be done
			return upper, upper, nil
		}

		lowerID = newID
	}

	return lower, upper, nil
}

func findInRange(ctx context.Context, s *stater, lower, upper *State, timestamp time.Time) (*State, error) {
	// binary search through the range to find the sequence number
	for lower.SeqNum+1 < upper.SeqNum {
		// could do better here
		splitID := (lower.SeqNum + upper.SeqNum) / 2
		split, err := s.State(ctx, splitID)
		if err != nil && !NotFound(err) {
			return nil, err
		}

		if split == nil {
			// file missing, search the next towards lower
			sID := splitID - 1
			for split == nil && lower.SeqNum < splitID {
				split, err = s.State(ctx, sID)
				if err != nil && !NotFound(err) {
					return nil, err
				}

				sID--
			}
		}

		if split == nil {
			// still missing?
			// search the next towards upper
			sID := splitID + 1
			for split == nil && splitID < upper.SeqNum {
				split, err = s.State(ctx, sID)
				if err != nil && !NotFound(err) {
					return nil, err
				}

				sID++
			}
		}

		if split == nil {
			// still nothing
			return lower, nil
		}

		// set the new boundary
		if timestamp.After(split.Timestamp) {
			lower = split
		} else {
			upper = split
		}
	}

	// timestamp is now between lower and upper, return the upper
	return upper, nil
}
