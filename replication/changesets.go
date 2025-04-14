package replication

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/pchchv/osm"
	"github.com/pchchv/osm/osmxml"
)

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


func changesetDecoder(ctx context.Context, r io.Reader) (osm.Changesets, error) {
	gzReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	var changesets []*osm.Changeset
	scanner := osmxml.New(ctx, gzReader)
	for scanner.Scan() {
		o := scanner.Object()
		c, ok := o.(*osm.Changeset)
		if !ok {
			return nil, fmt.Errorf("osm replication: object not a changeset: %[1]T: %[1]v", o)
		}
		changesets = append(changesets, c)
	}

	return changesets, scanner.Err()
}

// example
// ---
// last_run: 2016-07-02 22:46:01.422137422 +00:00  (or Z)
// sequence: 1912325
func decodeChangesetState(data []byte) (*State, error) {
	lines := bytes.Split(data, []byte("\n"))
	parts := bytes.Split(lines[1], []byte(":"))
	timeString := string(bytes.TrimSpace(bytes.Join(parts[1:], []byte(":"))))
	t, err := decodeTime(timeString)
	if err != nil {
		return nil, err
	}

	parts = bytes.Split(lines[2], []byte(":"))
	n, err := strconv.ParseUint(string(bytes.TrimSpace(parts[1])), 10, 64)
	if err != nil {
		return nil, err
	}

	return &State{
		SeqNum:    n,
		Timestamp: t,
	}, nil
}
