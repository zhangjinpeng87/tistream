package codec

import (
	"encoding/binary"
	"io"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
	"google.golang.org/protobuf/proto"
)

const (
	// The magic number of the committed data file.
	CommittedDataFileMagicNumber = uint32(0x30541989)
	CommittedDataFileVersion     = uint32(1)
)

type CommittedDataEncoder struct {
	// The tenant id.
	tenantID uint64

	// The range of the file.
	Range *pb.Task_Range

	lowWatermark  uint64
	highWatermark uint64
}

func NewCommittedDataEncoder(tenantID uint64, range_ *pb.Task_Range, lowWatermark, highWatermark uint64) *CommittedDataEncoder {
	return &CommittedDataEncoder{
		tenantID:      tenantID,
		Range:         range_,
		lowWatermark:  lowWatermark,
		highWatermark: highWatermark,
	}
}

func (c *CommittedDataEncoder) Encode(w io.Writer, events []*pb.EventRow) error {
	// Write header.
	if err := binary.Write(w, binary.LittleEndian, CommittedDataFileMagicNumber); err != nil {
		return err
	}

	// Write version.
	if err := binary.Write(w, binary.LittleEndian, CommittedDataFileVersion); err != nil {
		return err
	}

	// Write tenant id.
	if err := binary.Write(w, binary.LittleEndian, c.tenantID); err != nil {
		return err
	}

	// Write the range.
	if err := binary.Write(w, binary.LittleEndian, uint32(len(c.Range.Start))); err != nil {
		return err
	}
	if _, err := w.Write(c.Range.Start); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint32(len(c.Range.End))); err != nil {
		return err
	}
	if _, err := w.Write(c.Range.End); err != nil {
		return err
	}

	// Write watermarks.
	if err := binary.Write(w, binary.LittleEndian, c.lowWatermark); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, c.highWatermark); err != nil {
		return err
	}

	// Write all events.
	if err := binary.Write(w, binary.LittleEndian, uint32(len(events))); err != nil {
		return err
	}

	for _, event := range events {
		b, err := proto.Marshal(event)
		if err != nil {
			return err
		}
		l := uint32(len(b))
		if l > 0 {
			if err := binary.Write(w, binary.LittleEndian, l); err != nil {
				return err
			}
			if _, err := w.Write(b); err != nil {
				return err
			}
		}
	}

	// Todo: append the checksum.

	return nil
}

type CommittedDataDecoder struct {
	// The tenant id.
	TenantID uint64

	// range of the file.
	Range *pb.Task_Range
}

func NewCommittedDataDecoder(tenantID uint64, range_ *pb.Task_Range) *CommittedDataDecoder {
	return &CommittedDataDecoder{
		TenantID: tenantID,
		Range:    range_,
	}
}

func (c *CommittedDataDecoder) DecodeFrom(r io.Reader) ([]*pb.EventRow, error) {
	// Read the header.
	var magicNumber uint32
	if err := binary.Read(r, binary.LittleEndian, &magicNumber); err != nil {
		return nil, err
	}
	if magicNumber != CommittedDataFileMagicNumber {
		return nil, utils.ErrInvalidCommittedDataFile
	}

	// Read the version.
	var fileVersion uint32
	if err := binary.Read(r, binary.LittleEndian, &fileVersion); err != nil {
		return nil, err
	}

	// Todo : decode the file according to the version.
	switch fileVersion {
	case 1:
		return c.decodeBodyV1(r)
	default:
		return nil, utils.ErrInvalidCommittedDataFileVersion
	}
}

func (c *CommittedDataDecoder) decodeBodyV1(r io.Reader) ([]*pb.EventRow, error) {
	// Read the tenant id.
	var tenantID uint64
	if err := binary.Read(r, binary.LittleEndian, &tenantID); err != nil {
		return nil, err
	}
	if tenantID != c.TenantID {
		return nil, utils.ErrUnmatchedTenantID
	}

	// Read the range.
	var startLen uint32
	if err := binary.Read(r, binary.LittleEndian, &startLen); err != nil {
		return nil, err
	}
	start := make([]byte, startLen)
	if startLen > 0 {
		if _, err := io.ReadFull(r, start); err != nil {
			return nil, err
		}
	}
	var endLen uint32
	if err := binary.Read(r, binary.LittleEndian, &endLen); err != nil {
		return nil, err
	}
	end := make([]byte, endLen)
	if endLen > 0 {
		if _, err := io.ReadFull(r, end); err != nil {
			return nil, err
		}
	}
	if string(start) != string(c.Range.Start) || string(end) != string(c.Range.End) {
		return nil, utils.ErrUnmatchedRange
	}

	// Read the watermarks.
	var lowWatermark uint64
	if err := binary.Read(r, binary.LittleEndian, &lowWatermark); err != nil {
		return nil, err
	}
	var highWatermark uint64
	if err := binary.Read(r, binary.LittleEndian, &highWatermark); err != nil {
		return nil, err
	}

	// Read all events.
	var eventCount uint32
	if err := binary.Read(r, binary.LittleEndian, &eventCount); err != nil {
		return nil, err
	}

	events := make([]*pb.EventRow, 0, eventCount)
	for i := uint32(0); i < eventCount; i++ {
		var event pb.EventRow
		var eventSize uint32
		if err := binary.Read(r, binary.LittleEndian, &eventSize); err != nil {
			return nil, err
		}
		if eventSize > 0 {
			eventBytes := make([]byte, eventSize)
			if _, err := io.ReadFull(r, eventBytes); err != nil {
				return nil, err
			}
			if err := proto.Unmarshal(eventBytes, &event); err != nil {
				return nil, err
			}
			events = append(events, &event)
		}
	}

	return events, nil
}
