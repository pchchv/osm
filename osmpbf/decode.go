package osmpbf

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/pchchv/osm"
	"github.com/pchchv/osm/osmpbf/internal/osmpbf"
	"google.golang.org/protobuf/proto"
)

const (
	maxBlobSize       = 32 * 1024 * 1024
	maxBlobHeaderSize = 64 * 1024
)

// Header contains the contents of the header in the pbf file.
type Header struct {
	Bounds               *osm.Bounds
	RequiredFeatures     []string
	OptionalFeatures     []string
	WritingProgram       string
	Source               string
	ReplicationTimestamp time.Time
	ReplicationSeqNum    uint64
	ReplicationBaseURL   string
}

// oPair is a group sent over the channel from the decoder goroutines.
// It will contain the list of all objects.
type oPair struct {
	Offset  int64
	Objects []osm.Object
	Err     error
}

// iPair is a group sent over the channel to the
// decoder goroutines that unzip and decode the
// pbf from the headerblock.
type iPair struct {
	Offset int64
	Blob   *osmpbf.Blob
	Err    error
}

// Decoder reads and decodes OpenStreetMap PBF data from an input stream.
type decoder struct {
	scanner    *Scanner
	header     *Header
	r          io.Reader
	bytesRead  int64
	ctx        context.Context
	cancel     func()
	wg         sync.WaitGroup
	inputs     []chan<- iPair // for data decoders
	outputs    []<-chan oPair
	serializer chan oPair
	pOffset    int64
	cOffset    int64
	cData      oPair
	cIndex     int
}

// newDecoder returns a new decoder that reads from r.
func newDecoder(ctx context.Context, s *Scanner, r io.Reader) *decoder {
	c, cancel := context.WithCancel(ctx)
	return &decoder{
		scanner: s,
		ctx:     c,
		cancel:  cancel,
		r:       r,
	}
}

func (dec *decoder) readBlob(buf []byte) (*osmpbf.Blob, error) {
	if _, err := io.ReadFull(dec.r, buf); err != nil {
		return nil, err
	}

	blob := &osmpbf.Blob{}
	if err := proto.Unmarshal(buf, blob); err != nil {
		return nil, err
	}

	return blob, nil
}

func (dec *decoder) readBlobHeader(buf []byte) (*osmpbf.BlobHeader, error) {
	if _, err := io.ReadFull(dec.r, buf); err != nil {
		return nil, err
	}

	blobHeader := &osmpbf.BlobHeader{}
	if err := proto.Unmarshal(buf, blobHeader); err != nil {
		return nil, err
	}

	if blobHeader.GetDatasize() >= maxBlobSize {
		return nil, errors.New("blob size >= 32Mb")
	}

	return blobHeader, nil
}

func (dec *decoder) readBlobHeaderSize(buf []byte) (uint32, error) {
	if _, err := io.ReadFull(dec.r, buf); err != nil {
		return 0, err
	}

	size := binary.BigEndian.Uint32(buf)
	if size >= maxBlobHeaderSize {
		return 0, errors.New("blobHeader size >= 64Kb")
	}

	return size, nil
}
