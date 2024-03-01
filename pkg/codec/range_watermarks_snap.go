package codec

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
	"google.golang.org/protobuf/proto"
)

// The magic number of the range watermarks snap file.
const (
	RangeWatermarksSnapMagicNumber = uint32(0x81081793)
	RangeWatermarksSnapFileVersion = uint32(1)
)

// RangeWatermarksSnapEncoder is used to encode the range watermarks snap file.
type RangeWatermarksSnapEncoder struct {
	// The tenant id.
	TenantID uint64

	// range of the buffer.
	Range *pb.Task_Range
}

// NewRangeWatermarksSnapEncoder creates a new RangeWatermarksSnapEncoder.
func NewRangeWatermarksSnapEncoder(tenantID uint64, range_ *pb.Task_Range) *RangeWatermarksSnapEncoder {
	return &RangeWatermarksSnapEncoder{
		TenantID: tenantID,
		Range:    range_,
	}
}

// Encode encodes the range watermarks snap file.
func (r *RangeWatermarksSnapEncoder) Encode(w io.Writer, watermarks []*pb.EventWatermark) error {
	// Write header
	if err := binary.Write(w, binary.LittleEndian, RangeWatermarksSnapMagicNumber); err != nil {
		return err
	}

	// Write version.
	if err := binary.Write(w, binary.LittleEndian, RangeWatermarksSnapFileVersion); err != nil {
		return err
	}

	// Write tenant id.
	if err := binary.Write(w, binary.LittleEndian, r.TenantID); err != nil {
		return err
	}

	// Write the range.
	if err := binary.Write(w, binary.LittleEndian, uint32(len(r.Range.Start))); err != nil {
		return err
	}
	if _, err := w.Write(r.Range.Start); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint32(len(r.Range.End))); err != nil {
		return err
	}
	if _, err := w.Write(r.Range.End); err != nil {
		return err
	}

	// Write watermarks.
	if err := binary.Write(w, binary.LittleEndian, uint32(len(watermarks))); err != nil {
		return err
	}
	for _, watermark := range watermarks {
		data, err := proto.Marshal(watermark)
		if err != nil {
			return err
		}
		if err := binary.Write(w, binary.LittleEndian, uint32(len(data))); err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
	}

	return nil
}

// RangeWatermarksSnapDecoder is used to decode the range watermarks snap file.
type RangeWatermarksSnapDecoder struct {
	// The tenant id.
	TenantID uint64

	// range of the buffer.
	Range *pb.Task_Range
}

// NewRangeWatermarksSnapDecoder creates a new RangeWatermarksSnapDecoder.
func NewRangeWatermarksSnapDecoder(tenantID uint64, range_ *pb.Task_Range) *RangeWatermarksSnapDecoder {
	return &RangeWatermarksSnapDecoder{
		TenantID: tenantID,
		Range:    range_,
	}
}

// Decode decodes the range watermarks snap file.
func (r *RangeWatermarksSnapDecoder) Decode(reader io.Reader) ([]*pb.EventWatermark, error) {
	// Read the header.
	var magicNumber uint32
	if err := binary.Read(reader, binary.LittleEndian, &magicNumber); err != nil {
		return nil, err
	}
	if magicNumber != RangeWatermarksSnapMagicNumber {
		return nil, utils.ErrInvalidRangeWatermarksSnapFile
	}

	// Read the version.
	var fileVersion uint32
	if err := binary.Read(reader, binary.LittleEndian, &fileVersion); err != nil {
		return nil, err
	}
	if fileVersion != RangeWatermarksSnapFileVersion {
		return nil, utils.ErrInvalidRangeWatermarksSnapFileVersion
	}

	// Read the tenant id.
	var tenantID uint64
	if err := binary.Read(reader, binary.LittleEndian, &tenantID); err != nil {
		return nil, err
	}
	if tenantID != r.TenantID {
		return nil, utils.ErrInvalidRangeWatermarksSnapTenantID
	}

	// Read the range.
	var startLen uint32
	if err := binary.Read(reader, binary.LittleEndian, &startLen); err != nil {
		return nil, err
	}
	var start []byte
	if startLen > 0 {
		start := make([]byte, startLen)
		if _, err := io.ReadFull(reader, start); err != nil {
			return nil, err
		}
	}
	var endLen uint32
	if err := binary.Read(reader, binary.LittleEndian, &endLen); err != nil {
		return nil, err
	}
	var end []byte
	if endLen > 0 {
		end := make([]byte, endLen)
		if _, err := io.ReadFull(reader, end); err != nil {
			return nil, err
		}
	}
	if bytes.Compare(r.Range.Start, start) != 0 || bytes.Compare(r.Range.End, end) != 0 {
		return nil, utils.ErrInvalidRangeWatermarksSnapRange
	}

	// Read watermarks.
	var watermarksCnt uint32
	if err := binary.Read(reader, binary.LittleEndian, &watermarksCnt); err != nil {
		return nil, err
	}
	watermarks := make([]*pb.EventWatermark, 0, watermarksCnt)
	for i := uint32(0); i < watermarksCnt; i++ {
		var dataLen uint32
		if err := binary.Read(reader, binary.LittleEndian, &dataLen); err != nil {
			return nil, err
		}
		data := make([]byte, dataLen)
		if _, err := io.ReadFull(reader, data); err != nil {
			return nil, err
		}
		watermark := &pb.EventWatermark{}
		if err := proto.Unmarshal(data, watermark); err != nil {
			return nil, err
		}
		watermarks = append(watermarks, watermark)
	}

	return watermarks, nil
}
