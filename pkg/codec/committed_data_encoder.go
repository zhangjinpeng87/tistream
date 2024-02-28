package codec

import (
	"encoding/binary"
	"io"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
	"google.golang.org/protobuf/proto"
)

const (
	// The magic number of the committed data file.
	magicNumber = uint32(0x30541989)
	fileVersion = uint32(1)
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
	if err := binary.Write(w, binary.LittleEndian, magicNumber); err != nil {
		return err
	}

	// Write version.
	if err := binary.Write(w, binary.LittleEndian, fileVersion); err != nil {
		return err
	}

	// Write tenant id.
	if err := binary.Write(w, binary.LittleEndian, c.tenantID); err != nil {
		return err
	}

	// Write watermarks.
	if err := binary.Write(w, binary.LittleEndian, c.lowWatermark); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, c.highWatermark); err != nil {
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
