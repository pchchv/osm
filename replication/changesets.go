package replication

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
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


func (ds *Datasource) baseChangesetURL(cn ChangesetSeqNum) string {
	n := cn.Uint64()
	return fmt.Sprintf("%s/replication/%s/%03d/%03d/%03d", ds.baseURL(), cn.Dir(), n/1000000, (n%1000000)/1000, n%1000)
}

func (ds *Datasource) fetchChangesetState(ctx context.Context, n ChangesetSeqNum) (*State, error) {
	var url string
	if n.Uint64() != 0 {
		url = ds.baseChangesetURL(n) + ".state.txt"
	} else {
		url = fmt.Sprintf("%s/replication/%s/state.yaml", ds.baseURL(), n.Dir())
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := ds.client().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, &UnexpectedStatusCodeError{
			Code: resp.StatusCode,
			URL:  url,
		}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s, err := decodeChangesetState(data)
	if err != nil {
		return nil, err
	}

	// starting at 2008004 the changeset sequence number in the state file is one less than the name of the file
	// this is a consistent mistake
	// the correctly paired state and data files have the same name
	// the number in the state file is the one that is off
	if n == 0 {
		s.SeqNum++
	} else {
		s.SeqNum = uint64(n)
	}

	return s, nil
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
