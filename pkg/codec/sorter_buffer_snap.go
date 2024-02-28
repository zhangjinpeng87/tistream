package codec

import (
	"encoding/binary"
	"io"
	"google.golang.org/protobuf/proto"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

const (
	// The magic number of the sorter buffer snap file.
	SorterSnapMagicNumber = uint32(0x81081792)
	SorterSnapFileVersion = uint32(1)
)

type SorterBufferSnapEncoder struct {
	// The tenant id.
	TenantID uint64

	// range of the buffer.
	Range *pb.Task_Range
}

func NewSorterBufferSnapEncoder(tenantID uint64, range_ *pb.Task_Range) *SorterBufferSnapEncoder {
	return &SorterBufferSnapEncoder{
		TenantID: tenantID,
		Range:    range_,
	}
}

func (s *SorterBufferSnapEncoder) Encode(w io.Writer, events []*pb.EventRow) error {
	// Write header
	if err := binary.Write(w, binary.LittleEndian, SorterSnapMagicNumber); err != nil {
		return err
	}

	// Write version.
	if err := binary.Write(w, binary.LittleEndian, SorterSnapFileVersion); err != nil {
		return err
	}

	// Write tenant id.
	if err := binary.Write(w, binary.LittleEndian, s.TenantID); err != nil {
		return err
	}

	// Write the range.
	if err := binary.Write(w, binary.LittleEndian, uint32(len(s.Range.Start))); err != nil {
		return err
	}
	if _, err := w.Write(s.Range.Start); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint32(len(s.Range.End))); err != nil {
		return err
	}
	if _, err := w.Write(s.Range.End); err != nil {
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

	return nil
}
