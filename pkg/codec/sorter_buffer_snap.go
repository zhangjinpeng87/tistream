package codec

import (
	"encoding/binary"
	"io"

	"google.golang.org/protobuf/proto"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
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

type SorterBufferSnapDecoder struct {
	// The tenant id.
	TenantID uint64

	// range of the buffer.
	Range *pb.Task_Range
}

func NewSorterBufferSnapDecoder(tenantID uint64, range_ *pb.Task_Range) *SorterBufferSnapDecoder {
	return &SorterBufferSnapDecoder{
		TenantID: tenantID,
		Range:    range_,
	}
}

func (d *SorterBufferSnapDecoder) Decode(r io.Reader) ([]*pb.EventRow, error) {
	// Read and check the header.
	var magicNumber uint32
	if err := binary.Read(r, binary.LittleEndian, &magicNumber); err != nil {
		return nil, err
	}
	if magicNumber != SorterSnapMagicNumber {
		return nil, utils.ErrInvalidSorterBufferSnapFile
	}

	// Read and check the version.
	var fileVersion uint32
	if err := binary.Read(r, binary.LittleEndian, &fileVersion); err != nil {
		return nil, err
	}
	if fileVersion != SorterSnapFileVersion {
		return nil, utils.ErrInvalidSorterBufferSnapFileChecksum
	}

	// Read and check the tenant id.
	var tenantID uint64
	if err := binary.Read(r, binary.LittleEndian, &tenantID); err != nil {
		return nil, err
	}
	if tenantID != d.TenantID {
		return nil, utils.ErrUnmatchedTenantID
	}

	// Read the range.
	var startLen, endLen uint32
	if err := binary.Read(r, binary.LittleEndian, &startLen); err != nil {
		return nil, err
	}
	var start []byte
	if startLen > 0 {
		start := make([]byte, startLen)
		if _, err := io.ReadFull(r, start); err != nil {
			return nil, err
		}
	}
	if err := binary.Read(r, binary.LittleEndian, &endLen); err != nil {
		return nil, err
	}
	var end []byte
	if endLen > 0 {
		end := make([]byte, endLen)
		if _, err := io.ReadFull(r, end); err != nil {
			return nil, err
		}
	}
	if string(start) != string(d.Range.Start) || string(end) != string(d.Range.End) {
		return nil, utils.ErrUnmatchedRange
	}

	// Read all events.
	var eventCount uint32
	if err := binary.Read(r, binary.LittleEndian, &eventCount); err != nil {
		return nil, err
	}

	events := make([]*pb.EventRow, 0, eventCount)
	for i := uint32(0); i < eventCount; i++ {
		var eventLen uint32
		if err := binary.Read(r, binary.LittleEndian, &eventLen); err != nil {
			return nil, err
		}
		eventBytes := make([]byte, eventLen)
		if _, err := io.ReadFull(r, eventBytes); err != nil {
			return nil, err
		}
		event := &pb.EventRow{}
		if err := proto.Unmarshal(eventBytes, event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
