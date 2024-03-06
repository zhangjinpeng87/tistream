package codec

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

const (
	// The magic number of the committed range index file.
	CommittedRangeIndexMagicNumber = uint32(0x30541990)
	CommittedRangeIndexVersion     = uint32(1)
)

type CommittedRangeIndexEncoder struct {
	// The tenant id.
	TenantID uint64

	// The range of the file.
	Range *pb.Task_Range

	// All sub folders
	TsArray []uint64
}

func NewCommittedRangeIndexEncoder(tenantID uint64, range_ *pb.Task_Range, tsArray []uint64) *CommittedRangeIndexEncoder {
	return &CommittedRangeIndexEncoder{
		TenantID: tenantID,
		Range:    range_,
		TsArray:  tsArray,
	}
}

func (c *CommittedRangeIndexEncoder) Encode(w io.Writer) error {
	// Write header.
	if err := binary.Write(w, binary.LittleEndian, CommittedRangeIndexMagicNumber); err != nil {
		return err
	}

	// Write version.
	if err := binary.Write(w, binary.LittleEndian, CommittedRangeIndexVersion); err != nil {
		return err
	}

	// Write tenant id.
	if err := binary.Write(w, binary.LittleEndian, c.TenantID); err != nil {
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

	// Write the tsArray
	if err := binary.Write(w, binary.LittleEndian, uint32(len(c.TsArray))); err != nil {
		return err
	}
	for _, ts := range c.TsArray {
		if err := binary.Write(w, binary.LittleEndian, ts); err != nil {
			return err
		}
	}

	return nil
}

type CommittedRangeIndexDecoder struct {
	// The tenant id.
	TenantID uint64

	// The range of the file.
	Range *pb.Task_Range
}

func NewCommittedRangeIndexDecoder(tenantID uint64, range_ *pb.Task_Range) *CommittedRangeIndexDecoder {
	return &CommittedRangeIndexDecoder{
		TenantID: tenantID,
		Range:    range_,
	}
}

func (c *CommittedRangeIndexDecoder) Decode(r io.Reader) ([]uint64, error) {
	// Read the magic number.
	var magicNumber uint32
	if err := binary.Read(r, binary.LittleEndian, &magicNumber); err != nil {
		return nil, err
	}
	if magicNumber != CommittedRangeIndexMagicNumber {
		return nil, utils.ErrInvalidCommittedRangeIndexFile
	}

	// Read the version.
	var version uint32
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return nil, err
	}
	// Todo: decode body according to the version.

	// Read tenant id.
	if err := binary.Read(r, binary.LittleEndian, &c.TenantID); err != nil {
		return nil, err
	}

	// Read the range.
	var startLen uint32
	if err := binary.Read(r, binary.LittleEndian, &startLen); err != nil {
		return nil, err
	}
	var start []byte
	if startLen > 0 {
		start = make([]byte, startLen)
		if _, err := io.ReadFull(r, start); err != nil {
			return nil, err
		}
	}
	var endLen uint32
	if err := binary.Read(r, binary.LittleEndian, &endLen); err != nil {
		return nil, err
	}
	var end []byte
	if endLen > 0 {
		end = make([]byte, endLen)
		if _, err := io.ReadFull(r, end); err != nil {
			return nil, err
		}
	}
	if !bytes.Equal(start, c.Range.Start) || !bytes.Equal(end, c.Range.End) {
		return nil, utils.ErrInvalidCommittedRangeIndexFile
	}

	// Read the tsArray
	var tsArrayLen uint32
	if err := binary.Read(r, binary.LittleEndian, &tsArrayLen); err != nil {
		return nil, err
	}
	tsArray := make([]uint64, tsArrayLen)
	for i := uint32(0); i < tsArrayLen; i++ {
		if err := binary.Read(r, binary.LittleEndian, &tsArray[i]); err != nil {
			return nil, err
		}
	}

	return tsArray, nil
}
