package osmpbf

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
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

var parseCapabilities = map[string]bool{
	"OsmSchema-V0.6":        true,
	"DenseNodes":            true,
	"HistoricalInformation": true,
}

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

func (dec *decoder) readFileBlock(sizeBuf, headerBuf, blobBuf []byte) (*osmpbf.BlobHeader, *osmpbf.Blob, error) {
	blobHeaderSize, err := dec.readBlobHeaderSize(sizeBuf)
	if err != nil {
		return nil, nil, err
	}

	headerBuf = headerBuf[:blobHeaderSize]
	blobHeader, err := dec.readBlobHeader(headerBuf)
	if err != nil {
		return nil, nil, err
	}

	blobBuf = blobBuf[:blobHeader.GetDatasize()]
	blob, err := dec.readBlob(blobBuf)
	if err != nil {
		return nil, nil, err
	}

	dec.bytesRead += 4 + int64(blobHeaderSize) + int64(blobHeader.GetDatasize())
	return blobHeader, blob, nil
}

func getData(blob *osmpbf.Blob, data []byte) ([]byte, error) {
	switch {
	case blob.RawSize != nil:
		return blob.GetRaw(), nil
	case blob.Data != nil:
		r, err := zlibReader(blob.GetZlibData())
		if err != nil {
			return nil, err
		}

		// using the bytes.Buffer allows for the preallocation of the necessary space.
		l := blob.GetRawSize() + bytes.MinRead
		if cap(data) < int(l) {
			data = make([]byte, 0, l+l/10)
		} else {
			data = data[:0]
		}

		buf := bytes.NewBuffer(data)
		if _, err = buf.ReadFrom(r); err != nil {
			return nil, err
		}

		if buf.Len() != int(blob.GetRawSize()) {
			return nil, fmt.Errorf("raw blob data size %d but expected %d", buf.Len(), blob.GetRawSize())
		}

		return buf.Bytes(), nil
	default:
		return nil, errors.New("unknown blob data")
	}
}

func decodeOSMHeader(blob *osmpbf.Blob) (*Header, error) {
	data, err := getData(blob, nil)
	if err != nil {
		return nil, err
	}

	headerBlock := &osmpbf.HeaderBlock{}
	if err := proto.Unmarshal(data, headerBlock); err != nil {
		return nil, err
	}

	// capability check
	requiredFeatures := headerBlock.GetRequiredFeatures()
	for _, feature := range requiredFeatures {
		if !parseCapabilities[feature] {
			return nil, fmt.Errorf("parser does not have %s capability", feature)
		}
	}

	// read the header
	header := &Header{
		RequiredFeatures:   headerBlock.GetRequiredFeatures(),
		OptionalFeatures:   headerBlock.GetOptionalFeatures(),
		WritingProgram:     headerBlock.GetWritingprogram(),
		Source:             headerBlock.GetSource(),
		ReplicationBaseURL: headerBlock.GetOsmosisReplicationBaseUrl(),
		ReplicationSeqNum:  uint64(headerBlock.GetOsmosisReplicationSequenceNumber()),
	}

	// convert timestamp epoch seconds to golang time structure if it exists
	if headerBlock.OsmosisReplicationTimestamp != nil {
		header.ReplicationTimestamp = time.Unix(*headerBlock.OsmosisReplicationTimestamp, 0).UTC()
	}

	// read bounding box if it exists
	if headerBlock.Bbox != nil {
		// units are always in nanodegree and do not obey granularity rules
		// see osmformat.proto
		header.Bounds = &osm.Bounds{
			MinLon: 1e-9 * float64(*headerBlock.Bbox.Left),
			MaxLon: 1e-9 * float64(*headerBlock.Bbox.Right),
			MinLat: 1e-9 * float64(*headerBlock.Bbox.Bottom),
			MaxLat: 1e-9 * float64(*headerBlock.Bbox.Top),
		}
	}

	return header, nil
}
